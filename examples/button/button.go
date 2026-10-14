// Package button ports the upstream button-* examples.
package button

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "button", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("button-demo", Demo),
		ex("button-default", Default),
		ex("button-secondary", Secondary),
		ex("button-destructive", Destructive),
		ex("button-outline", Outline),
		ex("button-ghost", Ghost),
		ex("button-link", Link),
		ex("button-icon", Icon),
		ex("button-with-icon", WithIcon),
		ex("button-size", Size),
		ex("button-rounded", Rounded),
		ex("button-loading", Loading),
		ex("button-as-child", AsChild),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex flex-wrap items-center gap-2 md:flex-row"),
		button.Button(button.VariantOutline, "Button"),
		button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Submit"), lucide.ArrowUp()),
	)
}

func Default() h.Node     { return button.Button("Button") }
func Secondary() h.Node   { return button.Button(button.VariantSecondary, "Secondary") }
func Destructive() h.Node { return button.Button(button.VariantDestructive, "Destructive") }
func Outline() h.Node     { return button.Button(button.VariantOutline, "Outline") }
func Ghost() h.Node       { return button.Button(button.VariantGhost, "Ghost") }
func Link() h.Node        { return button.Button(button.VariantLink, "Link") }

func Icon() h.Node {
	return button.Button(button.VariantOutline, button.SizeIcon, lucide.CircleFadingArrowUp())
}

func WithIcon() h.Node {
	return button.Button(button.VariantOutline, button.SizeSm, tabler.IconGitBranch(), " New Branch")
}

func Size() h.Node {
	return h.Div(h.Class("flex flex-col items-start gap-8 sm:flex-row"),
		h.Div(h.Class("flex items-start gap-2"),
			button.Button(button.SizeSm, button.VariantOutline, "Small"),
			button.Button(button.SizeIconSm, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
		),
		h.Div(h.Class("flex items-start gap-2"),
			button.Button(button.VariantOutline, "Default"),
			button.Button(button.SizeIcon, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
		),
		h.Div(h.Class("flex items-start gap-2"),
			button.Button(button.VariantOutline, button.SizeLg, "Large"),
			button.Button(button.SizeIconLg, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
		),
	)
}

func Rounded() h.Node {
	return h.Div(h.Class("flex flex-col gap-8"),
		button.Button(button.VariantOutline, button.SizeIcon, h.Class("rounded-full"), lucide.ArrowUp()),
	)
}

func Loading() h.Node {
	return button.Button(button.SizeSm, button.VariantOutline, h.Disabled(), spinner.Spinner(), "Submit")
}

// AsChild uses a plain <a> where upstream uses next/link, which renders one.
func AsChild() h.Node {
	return button.Button(lib.AsChild, h.A(h.Href("/login"), "Login"))
}
