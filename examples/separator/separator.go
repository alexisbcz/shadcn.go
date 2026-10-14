// Package separator ports the upstream separator-* examples.
package separator

import (
	"shadcn.go/components/ui/separator"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(
		registry.Example{Name: "separator-demo", Component: "separator", Render: registry.Static(Demo)},
		registry.Example{Name: "separator-semantic", Component: "separator", Render: registry.Static(Semantic)},
	)
}

func Demo() h.Node {
	return h.Div(
		h.Div(h.Class("space-y-1"),
			h.H4(h.Class("text-sm leading-none font-medium"), "Radix Primitives"),
			h.P(h.Class("text-sm text-muted-foreground"), "An open-source UI component library."),
		),
		separator.Separator(h.Class("my-4")),
		h.Div(h.Class("flex h-5 items-center space-x-4 text-sm"),
			h.Div("Blog"),
			separator.Separator(separator.Vertical),
			h.Div("Docs"),
			separator.Separator(separator.Vertical),
			h.Div("Source"),
		),
	)
}

// Semantic ports the separator-semantic fixture
// (parity/fixtures/separator-semantic.tsx): non-decorative separators.
func Semantic() h.Node {
	return h.Div(h.Class("w-64"),
		h.Div(h.Class("text-sm"), "Above"),
		separator.Separator(separator.Decorative(false), h.Class("my-4")),
		h.Div(h.Class("flex h-5 items-center gap-4 text-sm"),
			h.Div("Left"),
			separator.Separator(separator.Decorative(false), separator.Vertical),
			h.Div("Right"),
		),
	)
}
