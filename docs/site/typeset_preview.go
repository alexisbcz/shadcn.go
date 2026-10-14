package site

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

// The typeset preview: app/(view)/preview/typeset (the iframe of /typeset
// and its "Open in New Tab"), and /typeset.css (app/typeset.css/route.ts,
// the raw stylesheet).

//go:embed static/typeset-preview.js
var typesetPreviewScript string

func init() { runtime.Module("typeset-preview", typesetPreviewScript) }

func (s *Site) typesetCSSPath() string {
	return filepath.Join(s.App, "app", "(app)", "(typeset)", "typeset.css")
}

// TypesetCSS serves /typeset.css, read from disk like upstream's route.
func (s *Site) TypesetCSS(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(s.typesetCSSPath())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(b)
}

// typesetPreviewData is what typeset-preview.js needs: the fonts and the
// measure widths.
var typesetPreviewData = func() string {
	measures := map[string]string{}
	for _, o := range typesetMeasures {
		measures[o.Value] = o.Width
	}
	return jsonString(map[string]any{"fonts": typesetData.Fonts, "measures": measures})
}()

// typesetPreviewVars is preview-provider.tsx's previewValues.
func typesetPreviewVars(p typesetParams) map[string]string {
	vars := map[string]string{
		"--preview-size":    p.Scale + "px",
		"--preview-leading": p.Leading,
		"--preview-flow":    p.Flow,
		"--preview-measure": typesetMeasureWidth(p.Measure),
	}
	body := ""
	if f := findTypesetFont(p.Body); f != nil {
		body = f.Value
	}
	vars["--preview-font"] = body
	if p.Heading == "inherit" {
		vars["--preview-font-heading"] = body
	} else if f := findTypesetFont(p.Heading); f != nil {
		vars["--preview-font-heading"] = f.Value
	}
	if f := findTypesetFont(p.Mono); f != nil {
		vars["--preview-font-mono"] = f.Value
	}
	return vars
}

// TypesetPreview serves /preview/typeset/{name}: the preview layout around
// a fixture.
func (s *Site) TypesetPreview(r *http.Request) (h.Node, error) {
	name := r.PathValue("name")
	var item *struct {
		Value string `json:"value"`
		Label string `json:"label"`
		HTML  string `json:"html"`
	}
	for i := range typesetData.Items {
		if typesetData.Items[i].Value == name {
			item = &typesetData.Items[i]
		}
	}
	if item == nil {
		return nil, errNotFound
	}
	p := parseTypesetParams(r.URL.Query())
	d := typesetDefaults
	font := findTypesetFont(d.Body).Value
	mono := findTypesetFont(d.Mono).Value

	vars, _ := json.Marshal(typesetPreviewVars(p))
	var content any = h.Div(h.Class("typeset w-full"), h.UnsafeHTML(item.HTML))
	if name == "chat" {
		content = h.Div(h.Class("flex w-full flex-col gap-10"),
			h.Div(h.Class("ml-auto w-fit max-w-[65%] rounded-3xl bg-muted px-4 py-2.5"), typesetData.ChatQuestion),
			content,
		)
	}
	css, err := os.ReadFile(s.typesetCSSPath())
	if err != nil {
		return nil, err
	}
	return s.Document(r, Meta{},
		// TypesetPreviewProvider applies the params to <html> before paint;
		// PreviewFontVariables puts the font classes on <html>.
		h.Script(h.UnsafeJS(`(function(v,c){var s=document.documentElement.style;for(var k in v)if(v[k])s.setProperty(k,v[k]);document.documentElement.classList.add.apply(document.documentElement.classList,c.split(" "))})(`+string(vars)+`,`+jsonString(typesetData.FontClasses)+`)`)),
		h.Link(h.Rel("stylesheet"), h.Href(staticPrefix+"typeset-fonts.css")),
		h.StyleEl(h.UnsafeCSS(string(css))),
		h.Div(h.Class(typesetData.FontClasses), h.Data("typeset-preview-root", ""), h.Data("typeset-preview-data", typesetPreviewData),
			h.StyleEl(h.UnsafeCSS(".preview-params .typeset {\n  --typeset-leading: var(--preview-leading, "+d.Leading+");\n  --typeset-flow: var(--preview-flow, "+d.Flow+");\n}")),
			h.Div(h.Class("flex min-h-svh justify-center px-6 pt-8 pb-24 md:pt-24 md:pb-32"),
				h.UnsafeAttr("style", "font-family: var(--preview-font, "+font+"); --font-heading: var(--preview-font-heading, "+font+"); --font-mono: var(--preview-font-mono, "+mono+")"),
				h.Div(h.Class("preview-params w-full"),
					h.UnsafeAttr("style", "font-size: var(--preview-size, "+d.Scale+"px); max-width: var(--preview-measure, "+typesetMeasureWidth(d.Measure)+")"),
					content,
				),
			),
		),
	), nil
}
