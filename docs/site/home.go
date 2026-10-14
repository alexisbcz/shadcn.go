package site

import (
	"net/http"

	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// Home renders / (app/(app)/(root)/page.tsx). Placeholder: the page header
// only. TODO: announcement, page actions, the examples/cards section.
func (s *Site) Home(r *http.Request) (h.Node, error) {
	return s.Document(r, Meta{},
		s.AppLayout(r,
			h.Section(h.Class("container-wrapper flex flex-col items-center gap-2 px-6 py-16 text-center"),
				h.H1(h.Class("max-w-2xl text-4xl font-semibold tracking-tight text-balance lg:text-5xl"), "The Foundation for your Design System"),
				h.P(h.Class("max-w-3xl text-base text-balance text-foreground sm:text-lg"), SiteDescription),
				h.Div(h.Class("flex gap-2 pt-2"),
					button.Button(button.Size("sm"), h.A(h.Href("/create"), "Get Started")).AsChild(),
					button.Button(button.Size("sm"), button.Variant("ghost"), h.A(h.Href("/docs/components"), "View Components")).AsChild(),
				),
			),
		),
	), nil
}
