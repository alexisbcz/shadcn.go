// Package collapsible ports the upstream collapsible-* examples and the
// port's collapsible-* fixtures.
package collapsible

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "collapsible", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("collapsible-demo", Demo),
		ex("collapsible-states", States),
		ex("collapsible-nested", Nested),
	)
}

func repo(name string) h.Node {
	return h.Div(h.Class("rounded-md border px-4 py-2 font-mono text-sm"), name)
}

// Demo is collapsible-demo (controlled upstream; the state starts closed).
func Demo() h.Node {
	return collapsible.Collapsible(h.Class("flex w-[350px] flex-col gap-2"),
		h.Div(h.Class("flex items-center justify-between gap-4 px-4"),
			h.H4(h.Class("text-sm font-semibold"), "@peduarte starred 3 repositories"),
			collapsible.Trigger(lib.AsChild,
				button.Button(button.VariantGhost, button.SizeIcon, h.Class("size-8"),
					lucide.ChevronsUpDown(),
					h.Span(h.Class("sr-only"), "Toggle"),
				),
			),
		),
		repo("@radix-ui/primitives"),
		collapsible.Content(h.Class("flex flex-col gap-2"),
			repo("@radix-ui/colors"),
			repo("@stitches/react"),
		),
	)
}

// States is the collapsible-states fixture: defaultOpen with an animated
// content, a disabled collapsible and a plain trigger.
func States() h.Node {
	return h.Div(h.Class("flex w-[350px] flex-col gap-6"),
		collapsible.Collapsible(collapsible.Open(true), h.Class("flex flex-col gap-2"),
			collapsible.Trigger(h.Class("rounded-md border px-4 py-2 text-sm"), "Toggle open"),
			collapsible.Content(h.Class("overflow-hidden data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down"),
				repo("@radix-ui/colors"),
				repo("@stitches/react"),
			),
		),
		collapsible.Collapsible(collapsible.Disabled(true), h.Class("flex flex-col gap-2"),
			collapsible.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Disabled")),
			collapsible.Content(repo("@radix-ui/themes")),
		),
	)
}

// Nested is the collapsible-nested fixture: a collapsible nested in
// another one's content, and a second trigger for the outer one.
func Nested() h.Node {
	box := func(text string) h.Node { return h.Div(h.Class("rounded-md border px-4 py-2 text-sm"), text) }
	trigger := func(text string) h.Node {
		return collapsible.Trigger(h.Class("rounded-md border px-4 py-2 text-sm"), text)
	}
	return h.Div(h.Class("flex w-[350px] flex-col gap-2"),
		collapsible.Collapsible(h.Class("flex flex-col gap-2"),
			h.Div(h.Class("flex gap-2"), trigger("Outer"), trigger("Outer too")),
			collapsible.Content(h.Class("flex flex-col gap-2"),
				box("Outer content"),
				collapsible.Collapsible(collapsible.Open(true), h.Class("flex flex-col gap-2"),
					trigger("Inner"),
					collapsible.Content(box("Inner content")),
				),
			),
		),
	)
}
