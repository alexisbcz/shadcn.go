package site

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"shadcn.go/docs/content"
	"shadcn.go/docs/search"
	"shadcn.go/docs/site/pagetree"
)

// The search API (app/api/search/route.ts): fumadocs' createFromSource,
// run by package search on the structured data of every page. The index is
// built in the background the first time a page renders the command menu
// (about a few seconds), and requests wait for it.

// searchAPIPath is the reference's search endpoint (useDocsSearch's
// "fetch" client calls it with ?query=).
const searchAPIPath = "/api/search"

type searchState struct {
	once   sync.Once
	done   chan struct{}
	engine *search.Engine
	err    error
}

var searchStates sync.Map // *Site -> *searchState

func (s *Site) searchState() *searchState {
	v, _ := searchStates.LoadOrStore(s, &searchState{done: make(chan struct{})})
	return v.(*searchState)
}

// startSearchIndex starts building the index (once).
func (s *Site) startSearchIndex() *searchState {
	st := s.searchState()
	st.once.Do(func() {
		go func() {
			defer close(st.done)
			st.engine, st.err = s.buildSearchIndex()
			if st.err != nil {
				s.logf("docs: search index: %v", st.err)
			}
		}()
	})
	return st
}

func (s *Site) buildSearchIndex() (*search.Engine, error) {
	e, err := search.New()
	if err != nil {
		return nil, err
	}
	var indexes []search.Index
	for _, d := range s.searchPageOrder() {
		b, err := os.ReadFile(filepath.Join(s.contentDir(), filepath.FromSlash(d.Path)))
		if err != nil {
			return nil, err
		}
		root, fm, err := content.RemarkTree(string(b))
		if err != nil {
			s.logf("docs: search index: %s: %v", d.Path, err)
			continue
		}
		sd, err := e.Structure(root)
		if err != nil {
			return nil, err
		}
		title := fm.Title
		if title == "" {
			title = filepath.Base(d.Path[:len(d.Path)-len(filepath.Ext(d.Path))])
		}
		indexes = append(indexes, search.Index{
			ID:             d.URL,
			Title:          title,
			Description:    fm.Description,
			URL:            d.URL,
			Breadcrumbs:    breadcrumbs(s.Tree, d.URL),
			StructuredData: sd,
		})
	}
	if err := e.SetIndex(indexes); err != nil {
		return nil, err
	}
	return e, nil
}

var collectionKeyRe = regexp.MustCompile(`"([^"\\]+\.mdx)":\s*[\w$.]+`)

// searchPageOrder returns the pages in the order of fumadocs' getPages()
// on the reference: the key order of the generated docs collection, which
// is the content glob's traversal order and differs from build to build.
// It is read from the reference build's server chunk holding the
// collection (.next/server/chunks), else from .source/server.ts, else the
// pages are sorted by path. The order only breaks ties between equal
// scores.
func (s *Site) searchPageOrder() []*Doc {
	byPath := map[string]*Doc{}
	var all []*Doc
	for _, d := range s.Docs {
		byPath[d.Path] = d
		all = append(all, d)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Path < all[j].Path })
	var keys []string
	if chunks, _ := filepath.Glob(filepath.Join(s.App, ".next", "server", "chunks", "*.js")); len(chunks) > 0 {
		for _, f := range chunks {
			b, err := os.ReadFile(f)
			if err != nil || !strings.Contains(string(b), `"(root)/index.mdx":`) {
				continue
			}
			for _, m := range collectionKeyRe.FindAllStringSubmatch(string(b), -1) {
				keys = append(keys, m[1])
			}
			if len(keys) > 0 {
				break
			}
		}
	}
	if len(keys) == 0 {
		if b, err := os.ReadFile(filepath.Join(s.App, ".source", "server.ts")); err == nil {
			for _, m := range collectionKeyRe.FindAllStringSubmatch(string(b), -1) {
				keys = append(keys, m[1])
			}
		}
	}
	var out []*Doc
	seen := map[string]bool{}
	for _, k := range keys {
		if d := byPath[k]; d != nil && !seen[k] {
			seen[k] = true
			out = append(out, d)
		}
	}
	for _, d := range all {
		if !seen[d.Path] {
			out = append(out, d)
		}
	}
	return out
}

// breadcrumbs is fumadocs' buildBreadcrumbs: the tree's name and the
// names of the folders (and separators) on the path to the page.
func breadcrumbs(tree *pagetree.Node, url string) []string {
	var run func(nodes []*pagetree.Node) []*pagetree.Node
	run = func(nodes []*pagetree.Node) []*pagetree.Node {
		var sep *pagetree.Node
		for _, n := range nodes {
			if n.Kind == pagetree.KindPage && n.URL == url {
				items := []*pagetree.Node{}
				if sep != nil {
					items = append(items, sep)
				}
				return append(items, n)
			}
			if n.Kind == pagetree.KindSeparator {
				sep = n
				continue
			}
			if n.Kind == pagetree.KindFolder {
				var items []*pagetree.Node
				if n.Index != nil && n.Index.URL == url {
					items = []*pagetree.Node{n.Index}
				} else {
					items = run(n.Children)
				}
				if items != nil {
					items = append([]*pagetree.Node{n}, items...)
					if sep != nil {
						items = append([]*pagetree.Node{sep}, items...)
					}
					return items
				}
			}
		}
		return nil
	}
	path := run(tree.Children)
	if path == nil {
		return nil
	}
	path = path[:len(path)-1]
	out := []string{"Docs"}
	for _, n := range path {
		if n.Name != "" {
			out = append(out, n.Name)
		}
	}
	return out
}

func (s *Site) serveSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("query")
	if q == "" {
		_, _ = w.Write([]byte("[]"))
		return
	}
	st := s.startSearchIndex()
	select {
	case <-st.done:
	case <-r.Context().Done():
		return
	}
	if st.err != nil {
		http.Error(w, st.err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := st.engine.Search(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

// searchJSON is used by tests.
func (s *Site) searchJSON(q string) ([]map[string]any, error) {
	st := s.startSearchIndex()
	<-st.done
	if st.err != nil {
		return nil, st.err
	}
	b, err := st.engine.Search(q)
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	return out, json.Unmarshal(b, &out)
}
