package site

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

// The create designer's preview: /preview/{base}/{name}
// (app/(view)/preview/[base]/[name]/page.tsx), the page its iframe loads.
// The design system comes from the search params; the iframe follows the
// designer's changes through create-preview.js: color, radius, font and
// menu changes apply in place, style and icon library changes reload the
// page with the new params (the Go components render a style's markup on
// the server).

//go:embed static/create-preview.js
var createPreviewScript string

func init() { runtime.Module("create-preview", createPreviewScript) }

// createPreviewItems are the preview items the Go site renders.
var createPreviewItems = map[string]func(designSystem) h.Node{
	"preview-02": preview02,
	"preview":    preview01,
}

// Preview serves /preview/{base}/{name}.
func (s *Site) Preview(r *http.Request) (h.Node, error) {
	base, name := r.PathValue("base"), r.PathValue("name")
	render, ok := createPreviewItems[name]
	if !ok || (base != "base" && base != "radix" && base != "aria") {
		return nil, errNotFound
	}
	ds := designSystemFromQuery(r.URL.Query())
	ds.Base, ds.Item = base, name
	st, ok := style.Parse(ds.Style)
	if !ok {
		st = style.Nova
	}
	return s.Document(r, Meta{Title: name},
		previewHead(ds),
		h.Div(h.Class("relative bg-background"), h.Data("create-preview", previewState(ds)),
			style.Scope(st, withIconLibrary(ds.IconLibrary, render(ds))),
		),
	), nil
}

// previewState is the design system the page was rendered with, for
// create-preview.js.
func previewState(ds designSystem) string {
	b, _ := json.Marshal(map[string]any{
		"base": ds.Base, "item": ds.Item, "style": ds.Style, "baseColor": ds.BaseColor, "theme": ds.Theme,
		"chartColor": ds.ChartColor, "iconLibrary": ds.IconLibrary, "font": ds.Font, "fontHeading": ds.FontHeading,
		"radius": ds.Radius, "menuAccent": ds.MenuAccent, "menuColor": ds.MenuColor, "pointer": ds.Pointer, "rtl": ds.RTL,
	})
	return string(b)
}

// previewHead is what the preview layout and DesignSystemProvider put on the
// document: the font variable classes and --font-sans/--font-heading on
// <html>, the style and base color classes on <body>, the theme variables
// and the hidden scrollbar (PreviewStyle).
func previewHead(ds designSystem) h.Node {
	d := createDataOnce()
	sans, heading := ds.fontFamilies()
	js, _ := json.Marshal(map[string]any{
		"html": d.fontClasses, "body": []string{"style-" + ds.Style, "base-color-" + ds.BaseColor},
		"sans": sans, "heading": heading,
	})
	return h.Group(
		h.Link(h.Rel("stylesheet"), h.Href(staticPrefix+"create-fonts.css")),
		// Generated from the embedded config; no user input reaches it
		// unvalidated (the params are matched against the config).
		h.Script(h.UnsafeJS(`(function(c){var d=document.documentElement;d.classList.add.apply(d.classList,c.html);document.body.classList.add.apply(document.body.classList,c.body);d.style.setProperty("--font-sans",c.sans);d.style.setProperty("--font-heading",c.heading)})(`+string(js)+`)`)),
		h.StyleEl(h.ID("design-system-theme-vars"), h.UnsafeCSS(ds.themeCSS())),
		h.StyleEl(h.UnsafeCSS("html{-ms-overflow-style:none;scrollbar-width:none}html::-webkit-scrollbar{display:none}"+previewCascadeCSS+previewLeadingCSS(ds.Style))),
	)
}

// previewCascadeCSS reproduces what the cascade does in upstream's preview,
// where the cn-* tokens stay in the markup and app/style-registry.css styles
// them: when two tokens of an element set the same property, the rule that
// comes later in the style's stylesheet wins, not the later token as in the
// generated components the Go ones follow. A field label is a Label
// (cn-label: leading-none) with cn-field-label (leading-snug).
const previewCascadeCSS = "label[data-slot=field-label]{line-height:1}"

// previewLeadingCSS is the other cascade effect: a token's leading-* rule
// stays in force under a text-* utility of the element (text-* reads
// --tw-leading), where the generated components lose it to
// tailwind-merge ("cn-card-title text-5xl" keeps leading-snug upstream).
// It gives each slot its token's --tw-leading, in the base layer.
func previewLeadingCSS(styleName string) string {
	st, ok := style.Parse(styleName)
	if !ok {
		return ""
	}
	return previewLeading(st)
}

var (
	previewLeadingCache sync.Map
	cnRuleRe            = regexp.MustCompile(`\.cn-([a-z0-9-]+)\s*\{\s*@apply ([^;]+);`)
	leadingValues       = map[string]string{"none": "1", "tight": "var(--leading-tight)", "snug": "var(--leading-snug)", "normal": "var(--leading-normal)", "relaxed": "var(--leading-relaxed)", "loose": "var(--leading-loose)"}
)

func previewLeading(st style.Style) string {
	if v, ok := previewLeadingCache.Load(st); ok {
		return v.(string)
	}
	var b strings.Builder
	b.WriteString("@layer base{")
	for _, m := range cnRuleRe.FindAllStringSubmatch(st.CSS(), -1) {
		for _, c := range strings.Fields(m[2]) {
			if v, ok := leadingValues[strings.TrimPrefix(c, "leading-")]; ok && strings.HasPrefix(c, "leading-") {
				b.WriteString("[data-slot=" + m[1] + "]{--tw-leading:" + v + "}")
			}
		}
	}
	b.WriteString("}")
	previewLeadingCache.Store(st, b.String())
	return b.String()
}
