// Package bench compares libhtml with html/template and gomponents on the
// same output. Run: go test -bench . -benchmem
package bench

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// --- data ----------------------------------------------------------------------

type row struct {
	ID     int
	Name   string
	Email  string
	Active bool
}

func rows(n int) []row {
	out := make([]row, n)
	for i := range out {
		out[i] = row{ID: i, Name: fmt.Sprintf("User <%d> & co", i), Email: fmt.Sprintf("user%d@example.com", i), Active: i%3 == 0}
	}
	return out
}

var (
	small = rows(10)
	large = rows(1000)
)

// --- libhtml -------------------------------------------------------------------

func libhtmlPage(title string, rs []row) h.Node {
	return h.Doctype(h.Html(h.Lang("en"),
		h.Head(h.Meta(h.Charset("utf-8")), h.Title(title)),
		h.Body(h.Class("p-4"),
			h.H1(h.Class("text-xl font-bold"), title),
			h.Table(h.Class("w-full"),
				h.Thead(h.Tr(h.Th("ID"), h.Th("Name"), h.Th("Email"))),
				h.Tbody(h.Map(rs, func(r row) h.Node {
					return h.Tr(h.ClassIf(r.Active, "active"),
						h.Td(r.ID),
						h.Td(h.A(h.Href(fmt.Sprintf("/users/%d", r.ID)), r.Name)),
						h.Td(h.A(h.Href("mailto:"+r.Email), r.Email)),
					)
				})),
			),
		),
	))
}

// --- html/template -------------------------------------------------------------

var tmpl = template.Must(template.New("page").Parse(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>{{.Title}}</title></head><body class="p-4"><h1 class="text-xl font-bold">{{.Title}}</h1><table class="w-full"><thead><tr><th>ID</th><th>Name</th><th>Email</th></tr></thead><tbody>{{range .Rows}}<tr{{if .Active}} class="active"{{end}}><td>{{.ID}}</td><td><a href="/users/{{.ID}}">{{.Name}}</a></td><td><a href="mailto:{{.Email}}">{{.Email}}</a></td></tr>{{end}}</tbody></table></body></html>`))

func templatePage(w io.Writer, title string, rs []row) error {
	return tmpl.Execute(w, struct {
		Title string
		Rows  []row
	}{title, rs})
}

// --- gomponents ----------------------------------------------------------------

func gomponentsPage(title string, rs []row) g.Node {
	return gh.Doctype(gh.HTML(gh.Lang("en"),
		gh.Head(gh.Meta(gh.Charset("utf-8")), gh.TitleEl(g.Text(title))),
		gh.Body(gh.Class("p-4"),
			gh.H1(gh.Class("text-xl font-bold"), g.Text(title)),
			gh.Table(gh.Class("w-full"),
				gh.THead(gh.Tr(gh.Th(g.Text("ID")), gh.Th(g.Text("Name")), gh.Th(g.Text("Email")))),
				gh.TBody(g.Map(rs, func(r row) g.Node {
					return gh.Tr(g.If(r.Active, gh.Class("active")),
						gh.Td(g.Textf("%d", r.ID)),
						gh.Td(gh.A(gh.Href(fmt.Sprintf("/users/%d", r.ID)), g.Text(r.Name))),
						gh.Td(gh.A(gh.Href("mailto:"+r.Email), g.Text(r.Email))),
					)
				})),
			),
		),
	))
}

// --- correctness: the three produce equivalent markup -------------------------------

func normalize(s string) string {
	// html/template escapes some characters differently (&#43; vs +, &lt; vs
	// &#60;); compare after unifying the entity forms that differ.
	r := strings.NewReplacer("&#34;", `"`, "&#39;", "'", "&#43;", "+", "&#60;", "&lt;", "&#62;", "&gt;", "&#38;", "&amp;", "&quot;", `"`, "<!doctype html>", "<!DOCTYPE html>")
	return r.Replace(s)
}

func TestOutputsAgree(t *testing.T) {
	var a, b, c bytes.Buffer
	if err := libhtmlPage("Users", small).Render(&a); err != nil {
		t.Fatal(err)
	}
	if err := templatePage(&b, "Users", small); err != nil {
		t.Fatal(err)
	}
	if err := gomponentsPage("Users", small).Render(&c); err != nil {
		t.Fatal(err)
	}
	if normalize(a.String()) != normalize(b.String()) {
		t.Errorf("libhtml vs html/template:\n%s\n%s", a.String(), b.String())
	}
	if normalize(a.String()) != normalize(c.String()) {
		t.Errorf("libhtml vs gomponents:\n%s\n%s", a.String(), c.String())
	}
}

// --- benchmarks ------------------------------------------------------------------

func benchSizes(b *testing.B, fn func(b *testing.B, rs []row)) {
	b.Run("rows=10", func(b *testing.B) { fn(b, small) })
	b.Run("rows=1000", func(b *testing.B) { fn(b, large) })
}

// Each benchmark builds the tree and renders it, as a request handler does.

func BenchmarkPage_libhtml(b *testing.B) {
	benchSizes(b, func(b *testing.B, rs []row) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			if err := libhtmlPage("Users", rs).Render(&buf); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(buf.Len()))
	})
}

func BenchmarkPage_htmlTemplate(b *testing.B) {
	benchSizes(b, func(b *testing.B, rs []row) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			if err := templatePage(&buf, "Users", rs); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(buf.Len()))
	})
}

func BenchmarkPage_gomponents(b *testing.B) {
	benchSizes(b, func(b *testing.B, rs []row) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			if err := gomponentsPage("Users", rs).Render(&buf); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(buf.Len()))
	})
}

// Rendering a prebuilt tree isolates the renderer from tree construction.
func BenchmarkRenderOnly_libhtml(b *testing.B) {
	benchSizes(b, func(b *testing.B, rs []row) {
		page := libhtmlPage("Users", rs)
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			if err := page.Render(&buf); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(buf.Len()))
	})
}

func BenchmarkRenderOnly_gomponents(b *testing.B) {
	benchSizes(b, func(b *testing.B, rs []row) {
		page := gomponentsPage("Users", rs)
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			if err := page.Render(&buf); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(buf.Len()))
	})
}

// --- components ------------------------------------------------------------------

func BenchmarkButton(b *testing.B) {
	b.Run("plain", func(b *testing.B) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			_ = button.Button(button.VariantOutline, button.SizeSm, "Save").Render(&buf)
		}
	})
	b.Run("className-override", func(b *testing.B) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			_ = button.Button(button.VariantOutline, button.SizeIcon, h.Class("rounded-full"), h.Aria("label", "Up")).Render(&buf)
		}
	})
	b.Run("asChild", func(b *testing.B) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			_ = button.Button(lib.AsChild, h.A(h.Href("/login"), "Login")).Render(&buf)
		}
	})
	b.Run("100-buttons", func(b *testing.B) {
		var buf bytes.Buffer
		for b.Loop() {
			buf.Reset()
			_ = h.Div(h.Range(100, func(i int) h.Node {
				return button.Button(button.VariantValues[i%len(button.VariantValues)], "Button ", i)
			})).Render(&buf)
		}
		b.SetBytes(int64(buf.Len()))
	})
}
