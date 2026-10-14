package content

import (
	"testing"

	h "shadcn.go/libhtml"
)

// TestRenderReact renders every page with the default registry and
// compares the DOM with React's rendering of the same page (testdata/react.json.gz,
// written by docs/oracle/render.mjs: the compiled MDX with the real
// mdx-components.tsx, icons, cn and React 19, where the app's own
// components are placeholders written like Placeholder).
func TestRenderReact(t *testing.T) {
	pages := parsedPages(t)
	reg := NewRegistry()
	failed, relAdded := 0, 0
	for _, rel := range oraclePages(t) {
		p := pages[rel]
		want, ok := reactData[rel]
		if !ok {
			t.Fatalf("%s: not in testdata/react.json.gz", rel)
		}
		node, err := Render(p, reg)
		if err != nil {
			failed++
			t.Errorf("%s: %v", rel, err)
			continue
		}
		got, err := h.RenderString(node)
		if err != nil {
			failed++
			t.Errorf("%s: %v", rel, err)
			continue
		}
		g, err := parseHTML(got)
		if err != nil {
			t.Fatalf("%s: parsing Go output: %v", rel, err)
		}
		w, err := parseHTML(want)
		if err != nil {
			t.Fatalf("%s: parsing React output: %v", rel, err)
		}
		// React DOM hoists <link rel="preload" as="image"> for every <img>
		// into the document head; it is not part of the content.
		w = dropPreloads(w)
		relAdded += libhtmlRel(g, w)
		if cg, cw := canon(g), canon(w); cg != cw {
			failed++
			if failed <= 6 {
				t.Errorf("%s differs:%s", rel, firstDiff(cg, cw))
			}
		}
	}
	t.Logf("%d pages, %d identical DOM, %d different", len(pages), len(pages)-failed, failed)
	t.Logf("%d links with target=_blank and no rel, where libhtml adds rel=\"noopener noreferrer\" (libhtml security default; see report)", relAdded)
}

// libhtmlRel removes the rel="noopener noreferrer" libhtml adds to links
// with target="_blank" that have no rel upstream, and counts them. It walks
// both trees in parallel and only touches nodes at the same position.
func libhtmlRel(got, want []*dnode) int {
	n := 0
	for i := 0; i < len(got) && i < len(want); i++ {
		g, w := got[i], want[i]
		if g.tag == "a" && w.tag == "a" && g.attrs["target"] == "_blank" && w.attrs["target"] == "_blank" {
			if _, has := w.attrs["rel"]; !has && g.attrs["rel"] == "noopener noreferrer" {
				delete(g.attrs, "rel")
				n++
			}
		}
		n += libhtmlRel(g.children, w.children)
	}
	return n
}

func dropPreloads(ns []*dnode) []*dnode {
	var out []*dnode
	for _, n := range ns {
		if n.tag == "link" && n.attrs["rel"] == "preload" {
			continue
		}
		n.children = dropPreloads(n.children)
		out = append(out, n)
	}
	return out
}

// TestTOCTitlesReact renders every TOC title and compares the DOM with
// React's rendering of the title fumadocs exports (testdata/oracle toc[].title).
func TestTOCTitlesReact(t *testing.T) {
	pages := parsedPages(t)
	total, failed := 0, 0
	for _, rel := range oraclePages(t) {
		o := loadOracle(t, rel)
		p := pages[rel]
		if len(p.TOC) != len(o.TOC) {
			t.Fatalf("%s: TOC length %d, want %d", rel, len(p.TOC), len(o.TOC))
		}
		for i, item := range p.TOC {
			total++
			n, err := RenderTOCTitle(p, item, nil)
			if err != nil {
				t.Fatal(err)
			}
			s, err := h.RenderString(n)
			if err != nil {
				t.Fatal(err)
			}
			g, err := parseHTML(s)
			if err != nil {
				t.Fatal(err)
			}
			w, err := parseHTML(o.TOC[i].Title)
			if err != nil {
				t.Fatal(err)
			}
			if canon(g) != canon(w) {
				failed++
				t.Errorf("%s: TOC %s: %s vs %s", rel, item.URL, canon(g), canon(w))
			}
		}
	}
	t.Logf("%d TOC titles, %d identical, %d different", total, total-failed, failed)
}
