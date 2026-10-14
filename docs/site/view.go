package site

import (
	"net/http"

	"shadcn.go/docs/site/mdx"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// View serves /view/{style}/{name}, the page block previews load in an
// iframe (app/(view)/view/[style]/[name]/page.tsx): the root layout with
// only the example, inside the view's ComponentPreview wrapper. The legacy
// style new-york-v4 renders the Go example as is; the radix, base and aria
// nova styles render the docs example in radix-nova (the closest Go
// rendering for base and aria).
func (s *Site) View(r *http.Request) (h.Node, error) {
	styleName, name := r.PathValue("style"), r.PathValue("name")
	var preview h.Node
	switch styleName {
	case "new-york-v4":
		if ex, ok := registry.Get(name); ok {
			preview = ex.Render(r)
		}
	case "radix-nova", "base-nova", "aria-nova":
		if ex, ok := mdx.LookupExample(name); ok {
			preview = ex(s.Env(r))
		}
	}
	if preview == nil {
		return nil, errNotFound
	}
	preview = mdx.NextPreview(s.Env(r), name, preview)
	return s.Document(r, Meta{Title: name},
		h.Div(h.Class("bg-background *:data-[slot=card]:has-[[data-slot=chart]]:shadow-none"), preview),
	), nil
}
