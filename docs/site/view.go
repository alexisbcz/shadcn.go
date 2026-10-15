package site

import (
	"net/http"
	"strings"

	"shadcn.go/components/style"
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
	// An example registered as "view--<name>" overrides <name> here: the
	// site serves some examples differently from a client render (server
	// components resolved at build time).
	if ex, ok := registry.Get("view--" + name); ok && styleName == "new-york-v4" {
		preview = ex.Render(r)
	} else if ok {
		preview = style.Scope(style.Nova, ex.Render(r))
	}
	switch {
	case preview != nil:
	case styleName == "new-york-v4":
		if ex, ok := registry.Get(name); ok {
			preview = ex.Render(r)
		}
	case styleName == "radix-nova" || styleName == "base-nova" || styleName == "aria-nova":
		if ex, ok := mdx.LookupExampleStyle(name, styleName); ok {
			preview = ex(s.Env(r))
			if strings.HasPrefix(styleName, "aria-") {
				// Inside NextPreview, which renders the preview.
				preview = style.WithBase(style.BaseAria, preview)
			}
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
