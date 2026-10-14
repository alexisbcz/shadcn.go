// Package pagetree builds the docs navigation tree from the content
// directory's meta.json files, like fumadocs' loader (source.pageTree):
// folders with a title and an ordered list of pages, "..." for the remaining
// entries, "!name" exclusions, "[Text](url)" links, "---Text---" separators,
// "(group)" folders that add no URL segment, and index pages that give a
// folder its URL.
//
// It is shared by the site chrome (docs sidebar, pager) and the MDX
// components that need the tree (ComponentsList).
package pagetree

import (
	"encoding/json"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Kind is a node type.
type Kind string

const (
	KindPage      Kind = "page"
	KindFolder    Kind = "folder"
	KindSeparator Kind = "separator"
)

// Node is a page, a folder or a separator.
type Node struct {
	Kind Kind
	// Name is the display name: the page title, the folder's meta title,
	// the link text or the separator text.
	Name string
	// URL is the page URL (pages and links).
	URL string
	// External is set for "[Text](url)" links.
	External bool
	// Path is the content path ("components/radix/button.mdx") for pages,
	// the folder path ("components/radix") for folders.
	Path string
	// Index is the folder's index page, if any.
	Index *Node
	// Children are a folder's entries, in meta.json order.
	Children []*Node
	// Root is set on folders whose meta.json has "root": true.
	Root bool
	// Frontmatter keys useful to the chrome (new badges, base switcher).
	Meta map[string]any
	// Fallback, on the tree root, lists in order the pages no meta.json
	// references (fumadocs' transformerFallback, tree.fallback): they are
	// not in the sidebar but findNeighbour still gives them a pager.
	Fallback []*Node
}

// PageInfo is what the tree needs to know about one page.
type PageInfo struct {
	Path  string // relative to the content root, slashes
	URL   string
	Title string
	Meta  map[string]any
}

type meta struct {
	Title string   `json:"title"`
	Pages []string `json:"pages"`
	Root  bool     `json:"root"`
}

var (
	linkRe = regexp.MustCompile(`^\[(.+)\]\((.+)\)$`)
	sepRe  = regexp.MustCompile(`^---(.*)---$`)
)

// Build builds the tree of the content directory root in fsys from its
// pages.
func Build(fsys fs.FS, root string, pages []PageInfo) *Node {
	byPath := map[string]PageInfo{}
	for _, p := range pages {
		byPath[p.Path] = p
	}
	t := folder(fsys, root, "", byPath)
	inTree := map[string]bool{}
	for _, p := range t.Pages() {
		if p.Path != "" {
			inTree[p.Path] = true
		}
	}
	t.Fallback = fallback(fsys, root, "", byPath, inTree)
	return t
}

// fallback flattens the pages missing from the tree the way fumadocs builds
// and flattens its fallback tree: per folder, the index page, then the
// other pages by path, then the sub-folders by path.
func fallback(fsys fs.FS, root, dir string, pages map[string]PageInfo, inTree map[string]bool) []*Node {
	var index, files, sub []*Node
	var dirs []string
	dirents, _ := fs.ReadDir(fsys, path.Join(root, dir))
	for _, e := range dirents {
		rel := path.Join(dir, e.Name())
		if e.IsDir() {
			dirs = append(dirs, rel)
			continue
		}
		p, ok := pages[rel]
		if !ok || inTree[rel] {
			continue
		}
		pn := &Node{Kind: KindPage, Name: p.Title, URL: p.URL, Path: rel, Meta: p.Meta}
		if e.Name() == "index.mdx" {
			index = append(index, pn)
		} else {
			files = append(files, pn)
		}
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	sort.Strings(dirs)
	for _, d := range dirs {
		sub = append(sub, fallback(fsys, root, d, pages, inTree)...)
	}
	return append(append(index, files...), sub...)
}

func folder(fsys fs.FS, root, dir string, pages map[string]PageInfo) *Node {
	full := path.Join(root, dir)
	n := &Node{Kind: KindFolder, Path: dir, Name: folderName(path.Base(dir))}
	var m meta
	if b, err := fs.ReadFile(fsys, path.Join(full, "meta.json")); err == nil {
		_ = json.Unmarshal(b, &m)
	}
	if m.Title != "" {
		n.Name = m.Title
	}
	n.Root = m.Root

	// Entries of the folder: sub-folders and pages, by file name.
	entries := map[string]*Node{}
	var names []string
	dirents, _ := fs.ReadDir(fsys, full)
	for _, e := range dirents {
		name := e.Name()
		rel := path.Join(dir, name)
		switch {
		case e.IsDir():
			entries[name] = folder(fsys, root, rel, pages)
			names = append(names, name)
		case strings.HasSuffix(name, ".mdx"):
			p, ok := pages[rel]
			if !ok {
				continue
			}
			base := strings.TrimSuffix(name, ".mdx")
			pn := &Node{Kind: KindPage, Name: p.Title, URL: p.URL, Path: rel, Meta: p.Meta}
			if base == "index" {
				n.Index = pn
			}
			entries[base] = pn
			names = append(names, base)
		}
	}
	sort.Strings(names)

	used := map[string]bool{}
	excluded := map[string]bool{}
	for _, item := range m.Pages {
		if strings.HasPrefix(item, "!") {
			excluded[item[1:]] = true
		}
	}
	add := func(name string) {
		if e, ok := entries[name]; ok && !used[name] && !excluded[name] {
			used[name] = true
			n.Children = append(n.Children, e)
		}
	}
	// No "pages" lists everything; an empty list (changelog/meta.json)
	// lists nothing, as in fumadocs.
	if m.Pages == nil {
		m.Pages = []string{"..."}
	}
	for _, item := range m.Pages {
		switch {
		case strings.HasPrefix(item, "!"):
		case item == "...":
			for _, name := range names {
				if name == "index" && n.Index != nil && !containsStr(m.Pages, "index") {
					continue
				}
				if !referenced(m.Pages, name) {
					add(name)
				}
			}
		case strings.HasPrefix(item, "..."):
			if f, ok := entries[item[3:]]; ok && f.Kind == KindFolder {
				used[item[3:]] = true
				n.Children = append(n.Children, f.Children...)
			}
		case linkRe.MatchString(item):
			sm := linkRe.FindStringSubmatch(item)
			n.Children = append(n.Children, &Node{Kind: KindPage, Name: sm[1], URL: sm[2], External: !strings.HasPrefix(sm[2], "/docs")})
		case sepRe.MatchString(item):
			n.Children = append(n.Children, &Node{Kind: KindSeparator, Name: sepRe.FindStringSubmatch(item)[1]})
		default:
			add(item)
		}
	}
	return n
}

func referenced(items []string, name string) bool {
	for _, it := range items {
		if it == name || it == "!"+name || it == "..."+name {
			return true
		}
	}
	return false
}

func containsStr(items []string, s string) bool {
	for _, it := range items {
		if it == s {
			return true
		}
	}
	return false
}

func folderName(base string) string {
	base = strings.Trim(base, "()")
	if base == "" || base == "." {
		return ""
	}
	return strings.ToUpper(base[:1]) + strings.ReplaceAll(base[1:], "-", " ")
}

// Find returns the folder at a content path ("components/radix").
func (n *Node) Find(dir string) *Node {
	if n.Kind == KindFolder && n.Path == dir {
		return n
	}
	for _, c := range n.Children {
		if c.Kind == KindFolder {
			if f := c.Find(dir); f != nil {
				return f
			}
		}
	}
	return nil
}

// Pages returns the tree's pages in order (folder index pages first in
// their folder, "[Text](url)" links included), like fumadocs' flattenTree.
func (n *Node) Pages() []*Node {
	var out []*Node
	var walk func(*Node)
	walk = func(f *Node) {
		if f.Index != nil {
			out = append(out, f.Index)
		}
		for _, c := range f.Children {
			switch c.Kind {
			case KindPage:
				if c != f.Index {
					out = append(out, c)
				}
			case KindFolder:
				walk(c)
			}
		}
	}
	walk(n)
	return out
}

// Neighbours returns the previous and next pages of url (fumadocs'
// findNeighbour), nil when there is none.
func (n *Node) Neighbours(url string) (prev, next *Node) {
	for _, pages := range [][]*Node{n.Pages(), n.Fallback} {
		for i, p := range pages {
			if p.URL == url {
				if i > 0 {
					prev = pages[i-1]
				}
				if i+1 < len(pages) {
					next = pages[i+1]
				}
				return
			}
		}
	}
	return nil, nil
}
