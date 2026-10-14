// Package resizable ports the upstream resizable-* examples.
package resizable

import (
	"shadcn.go/components/ui/resizable"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "resizable", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("resizable-demo", Demo),
		ex("resizable-demo-with-handle", DemoWithHandle),
		ex("resizable-handle", Handle),
		ex("resizable-vertical", Vertical),
		ex("resizable-collapsible", Collapsible),
		ex("resizable-disabled", DisabledGroup),
	)
}

func label(class, text string) h.Node {
	return h.Div(h.Class(class), h.Span(h.Class("font-semibold"), text))
}

func demo(handle ...any) h.Node {
	return resizable.PanelGroup(resizable.Horizontal, h.Class("max-w-md rounded-lg border md:min-w-[450px]"),
		resizable.Panel(resizable.DefaultSize("50%"), label("flex h-[200px] items-center justify-center p-6", "One")),
		resizable.Handle(handle...),
		resizable.Panel(resizable.DefaultSize("50%"),
			resizable.PanelGroup(resizable.Vertical,
				resizable.Panel(resizable.DefaultSize("25%"), label("flex h-full items-center justify-center p-6", "Two")),
				resizable.Handle(handle...),
				resizable.Panel(resizable.DefaultSize("75%"), label("flex h-full items-center justify-center p-6", "Three")),
			),
		),
	)
}

func Demo() h.Node           { return demo() }
func DemoWithHandle() h.Node { return demo(resizable.WithHandle) }

func Handle() h.Node {
	return resizable.PanelGroup(resizable.Horizontal, h.Class("min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]"),
		resizable.Panel(resizable.DefaultSize("25%"), label("flex h-full items-center justify-center p-6", "Sidebar")),
		resizable.Handle(resizable.WithHandle),
		resizable.Panel(resizable.DefaultSize("75%"), label("flex h-full items-center justify-center p-6", "Content")),
	)
}

func Vertical() h.Node {
	return resizable.PanelGroup(resizable.Vertical, h.Class("min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]"),
		resizable.Panel(resizable.DefaultSize("25%"), label("flex h-full items-center justify-center p-6", "Header")),
		resizable.Handle(),
		resizable.Panel(resizable.DefaultSize("75%"), label("flex h-full items-center justify-center p-6", "Content")),
	)
}

// Collapsible ports the resizable-collapsible fixture
// (parity/fixtures/resizable-collapsible.tsx).
func Collapsible() h.Node {
	return resizable.PanelGroup(resizable.Horizontal, h.Class("min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]"),
		resizable.Panel(resizable.DefaultSize("25%"), resizable.MinSize("15%"), resizable.MaxSize("40%"), resizable.Collapsible(true), resizable.CollapsedSize("5%"),
			label("flex h-full items-center justify-center p-6", "Sidebar")),
		resizable.Handle(resizable.WithHandle),
		resizable.Panel(resizable.DefaultSize("75%"), resizable.MinSize("30%"), label("flex h-full items-center justify-center p-6", "Content")),
	)
}

// DisabledGroup ports the resizable-disabled fixture
// (parity/fixtures/resizable-disabled.tsx).
func DisabledGroup() h.Node {
	return resizable.PanelGroup(resizable.Horizontal, resizable.Disabled(true), h.Class("min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]"),
		resizable.Panel(resizable.DefaultSize("40%"), label("flex h-full items-center justify-center p-6", "One")),
		resizable.Handle(resizable.WithHandle),
		resizable.Panel(resizable.DefaultSize("60%"), label("flex h-full items-center justify-center p-6", "Two")),
	)
}
