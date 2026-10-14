// Package label ports the upstream label-* examples.
package label

import (
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(
		registry.Example{Name: "label-demo", Component: "label", Render: registry.Static(Demo)},
	)
}

func Demo() h.Node {
	return h.Div(
		h.Div(h.Class("flex items-center space-x-2"),
			checkbox.Checkbox(h.ID("terms")),
			label.Label(h.For("terms"), "Accept terms and conditions"),
		),
	)
}
