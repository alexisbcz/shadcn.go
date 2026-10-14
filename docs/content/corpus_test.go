package content

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"

	"shadcn.go/docs/hast"
	"shadcn.go/docs/mdx"
	h "shadcn.go/libhtml"
)

// TestCorpusHast runs the whole pipeline (parse, fumadocs' remark and
// rehype plugins, highlighting) on every input of testdata/corpus.json.gz
// (docs/oracle/corpus.mjs: CommonMark spec examples, curated edge cases,
// every block of the docs and seeded mutations of them). The Go pipeline
// may reject an input, but when it accepts one, upstream must accept it
// too and produce the same HTML tree and TOC, byte for byte.
func TestCorpusHast(t *testing.T) {
	f, err := os.Open("testdata/corpus.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var list []struct {
		Kind      string `json:"kind"`
		Src       string `json:"src"`
		Hast      string `json:"hast"`
		HastError string `json:"hastError"`
		TOC       []struct {
			Depth int    `json:"depth"`
			URL   string `json:"url"`
			Hast  string `json:"hast"`
		} `json:"toc"`
	}
	if err := json.NewDecoder(zr).Decode(&list); err != nil {
		t.Fatal(err)
	}
	h := testHighlighter(t)
	type result struct{ accepted, problem string }
	results := make([]result, len(list))
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.GOMAXPROCS(0))
	for i := range list {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			e := list[i]
			root, err := mdx.Parse(e.Src)
			if err != nil {
				return
			}
			tree, toc, err := Transform(root, &Frontmatter{Title: "Corpus"}, Options{Highlighter: h})
			if err != nil {
				return
			}
			results[i].accepted = "yes"
			if e.HastError != "" {
				results[i].problem = "accepted, upstream fails: " + e.HastError
				return
			}
			if got := hast.ToHTML(tree...); got != e.Hast {
				results[i].problem = "hast differs:\n got: " + around(got, e.Hast) + "\nwant: " + around(e.Hast, got)
				return
			}
			if len(toc) != len(e.TOC) {
				results[i].problem = fmt.Sprintf("TOC has %d entries, want %d", len(toc), len(e.TOC))
				return
			}
			for j, it := range toc {
				w := e.TOC[j]
				if it.Depth != w.Depth || it.URL != w.URL || hast.ToHTML(it.Title...) != w.Hast {
					results[i].problem = fmt.Sprintf("TOC entry %d: %d %s %s, want %d %s %s", j, it.Depth, it.URL, hast.ToHTML(it.Title...), w.Depth, w.URL, w.Hast)
					return
				}
			}
		}(i)
	}
	wg.Wait()
	stats := map[string]*[3]int{} // identical, rejected, problems
	failed := 0
	for i, r := range results {
		k := list[i].Kind
		if stats[k] == nil {
			stats[k] = &[3]int{}
		}
		switch {
		case r.accepted == "":
			stats[k][1]++
		case r.problem != "":
			stats[k][2]++
			failed++
			if failed <= 20 || os.Getenv("CORPUS_ALL") != "" {
				t.Errorf("input %d (%s): %s\n%q", i, k, r.problem, short(list[i].Src))
			}
		default:
			stats[k][0]++
		}
	}
	for _, k := range []string{"spec", "curated", "docs", "mutation"} {
		if s := stats[k]; s != nil {
			t.Logf("%-8s identical %5d, rejected %5d, different or accepted where upstream fails %d", k, s[0], s[1], s[2])
		}
	}
}

// TestCorpusRender renders every corpus input the Go pipeline accepts and
// compares the DOM with React's rendering of the compiled MDX with the real
// mdx-components.tsx (testdata/corpus-react.json.gz, written by
// `node docs/oracle/render.mjs --corpus`), like TestRenderReact does for
// the pages. Inputs upstream cannot compile or run (undefined components,
// unresolvable imports) are skipped; a Go render error counts as a
// rejection.
func TestCorpusRender(t *testing.T) {
	var list []struct {
		Kind string `json:"kind"`
		Src  string `json:"src"`
	}
	var react []*string
	for name, v := range map[string]any{"corpus.json.gz": &list, "corpus-react.json.gz": &react} {
		f, err := os.Open("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		zr, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.NewDecoder(zr).Decode(v); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
	if len(react) != len(list) {
		t.Fatalf("corpus-react.json.gz has %d entries, corpus %d: regenerate it", len(react), len(list))
	}
	h := testHighlighter(t)
	reg := NewRegistry()
	type result struct{ state, problem string }
	results := make([]result, len(list))
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.GOMAXPROCS(0))
	for i := range list {
		if react[i] == nil {
			results[i].state = "skipped"
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			p, err := ParsePage("corpus.mdx", "---\ntitle: Corpus\n---\n"+list[i].Src, Options{Highlighter: h})
			if err != nil {
				results[i].state = "rejected"
				return
			}
			node, err := Render(p, reg)
			if err != nil {
				results[i].state = "render error"
				return
			}
			s, err := h2s(node)
			if err != nil {
				results[i].state = "render error"
				return
			}
			g, err := parseHTML(s)
			if err != nil {
				results[i] = result{"different", "parsing Go output: " + err.Error()}
				return
			}
			w, err := parseHTML(*react[i])
			if err != nil {
				results[i] = result{"different", "parsing React output: " + err.Error()}
				return
			}
			w = dropPreloads(w)
			libhtmlRel(g, w)
			if cg, cw := canon(g), canon(w); cg != cw {
				results[i] = result{"different", firstDiff(cg, cw)}
				return
			}
			results[i].state = "identical"
		}(i)
	}
	wg.Wait()
	stats := map[string]map[string]int{}
	failed := 0
	for i, r := range results {
		k := list[i].Kind
		if stats[k] == nil {
			stats[k] = map[string]int{}
		}
		stats[k][r.state]++
		if r.state == "different" {
			failed++
			if failed <= 20 || os.Getenv("CORPUS_ALL") != "" {
				t.Errorf("input %d (%s): %s\n%q", i, k, r.problem, short(list[i].Src))
			}
		}
	}
	for _, k := range []string{"spec", "curated", "docs", "mutation"} {
		if s := stats[k]; s != nil {
			t.Logf("%-8s identical %5d, rejected %5d, render errors %d, skipped (upstream fails) %5d, different %d", k, s["identical"], s["rejected"], s["render error"], s["skipped"], s["different"])
		}
	}
}

func h2s(n h.Node) (string, error) { return h.RenderString(n) }
