package label

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the label page (ui/apps/v4/examples/radix/label-*.tsx),
// in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "label", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("label-demo", NovaDemo),
		nova("label-rtl", NovaRtl),
	)
}

func NovaDemo() h.Node {
	return h.Div(h.Class("flex gap-2"),
		checkbox.Checkbox(h.ID("terms")),
		label.Label(h.For("terms"), "Accept terms and conditions"),
	)
}

func NovaRtl() h.Node {
	return style.ScopeRTL(h.Div(h.Class("flex gap-2"), h.Dir("rtl"),
		checkbox.Checkbox(h.ID("terms-rtl"), h.Dir("rtl")),
		label.Label(h.For("terms-rtl"), h.Dir("rtl"), "قبول الشروط والأحكام"),
	))
}
