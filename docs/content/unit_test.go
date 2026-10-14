package content

import (
	"strings"
	"testing"

	"shadcn.go/docs/hast"
	h "shadcn.go/libhtml"
)

func TestSlugger(t *testing.T) {
	s := &Slugger{}
	for _, c := range []struct{ in, want string }{
		{"Hello World", "hello-world"},
		{"Hello World", "hello-world-1"},
		{"hello-world", "hello-world-2"},
		{"What's new? (v2)", "whats-new-v2"},
		{"<Button /> & `cn`", "button---cn"}, // expected values from github-slugger 2.0.0
		{"Ünïcödé Ωmega", "ünïcödé-ωmega"},
		{"March 19, 2025 - Deprecate tailwindcss-animate", "march-19-2025---deprecate-tailwindcss-animate"},
	} {
		if got := s.Slug(c.in); got != c.want {
			t.Errorf("Slug(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFrontmatter(t *testing.T) {
	fm, err := ParseFrontmatter("title: \"RTL\"\ndescription: You've got it\ndate: 2026-09-03\nfeatured: true\nlinks:\n  doc: https://a.b/c#d\n  api: https://e.f")
	if err != nil {
		t.Fatal(err)
	}
	if fm.Title != "RTL" || fm.Description != "You've got it" || !fm.Featured || !fm.HasDate || fm.Date.Format("2006-01-02") != "2026-09-03" {
		t.Errorf("unexpected %+v", fm)
	}
	if fm.Links == nil || fm.Links.Doc != "https://a.b/c#d" || fm.Links.API != "https://e.f" {
		t.Errorf("links %+v", fm.Links)
	}
	for _, bad := range []string{"tags: [a, b]", "title: |\n  x", "title: a: b"} {
		if _, err := ParseFrontmatter(bad); err == nil {
			t.Errorf("expected an error for %q", bad)
		}
	}
}

// TestUnsupported checks that input the docs never use is rejected rather
// than rendered differently from upstream.
func TestUnsupported(t *testing.T) {
	for _, src := range []string{
		"Title\n=====\n",                      // setext heading
		"[ref]: https://example.com\n",        // link reference definition
		"- [ ] task\n",                        // task list item
		"```npm\nnpm i x\n```\n",              // remarkNpm
		"```bash tab=\"npm\"\nnpm i x\n```\n", // remarkCodeTab
		"Use `x{:js}` inline.\n",              // rehype-pretty-code inline language
		"<include>file.mdx</include>\n",       // remarkInclude
		"![local](./image.png)\n",             // remarkImage with a file
		"<Foo {...props} />\n",                // JSX spread attribute
		"<Callout>\n\nText\n",                 // unclosed JSX
	} {
		if _, err := ParsePage("x.mdx", src, Options{Highlighter: testHighlighter(t)}); err == nil {
			t.Errorf("expected an error for %q", src)
		}
	}
	// Expressions are evaluated when rendering.
	for _, src := range []string{"{foo()}\n", "{a.b}\n", "<Unknown />\n"} {
		p, err := ParsePage("x.mdx", src, Options{Highlighter: testHighlighter(t)})
		if err != nil {
			t.Errorf("%q: %v", src, err)
			continue
		}
		if _, err := Render(p, nil); err == nil {
			t.Errorf("expected a render error for %q", src)
		}
	}
}

func TestRegistry(t *testing.T) {
	p, err := ParsePage("x.mdx", "---\ntitle: X\ndescription: Y\n---\n\n<Callout title=\"Note\" icon={<InfoIcon />}>\n  Some **text**.\n</Callout>\n\nimport { InfoIcon } from \"lucide-react\"\n", Options{Highlighter: testHighlighter(t)})
	if err != nil {
		t.Fatal(err)
	}
	// Default: placeholder, with the icon rendered by the default registry.
	out, err := Render(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := h.MustRenderString(out)
	if !strings.Contains(s, `<mdx-component data-name="Callout" data-props="[[&#34;title&#34;,&#34;Note&#34;]]"><template data-prop="icon"><svg`) || !strings.Contains(s, "lucide-info") || !strings.Contains(s, "<p>Some <strong>text</strong>.</p>") {
		t.Errorf("unexpected placeholder rendering: %s", s)
	}
	// A registered component replaces the placeholder.
	reg := NewRegistry()
	reg.Register("Callout", func(ctx *Context, props Props, children []h.Node) (h.Node, error) {
		icon, _ := props.Get("icon")
		return h.Div(h.Class("callout"), h.Data("title", props.String("title")), icon, h.Group(toAny(children)...)), nil
	})
	out, err = Render(p, reg)
	if err != nil {
		t.Fatal(err)
	}
	s = h.MustRenderString(out)
	if !strings.HasPrefix(s, `<div class="callout" data-title="Note"><svg`) || strings.Contains(s, "mdx-component") {
		t.Errorf("unexpected rendering: %s", s)
	}
}

func TestPageURL(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"(root)/index.mdx", "/docs"},
		{"(root)/cli.mdx", "/docs/cli"},
		{"components/radix/button.mdx", "/docs/components/radix/button"},
		{"changelog/index.mdx", "/docs/changelog"},
	} {
		p, err := ParsePage(c.in, "---\ntitle: T\ndescription: D\n---\n\nx\n", Options{Highlighter: testHighlighter(t)})
		if err != nil {
			t.Fatal(err)
		}
		if p.URL != c.want {
			t.Errorf("%s: URL %s, want %s", c.in, p.URL, c.want)
		}
	}
}

func TestHastHelpers(t *testing.T) {
	e := hast.El("a", "href", "x")
	e.Children = []hast.Node{&hast.Text{Value: "<&>"}}
	if got := hast.ToHTML(e); got != `<a href="x">&#x3C;&#x26;></a>` {
		t.Errorf("ToHTML = %s", got)
	}
}
