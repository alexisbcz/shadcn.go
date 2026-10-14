package resizable

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/resizable"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/resizable-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "resizable", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("resizable-demo", NovaDemo),
		nova("resizable-vertical", NovaVertical),
		nova("resizable-rtl", NovaRtl),
	)
}

func novaPanel(class, text string) h.Node {
	return h.Div(h.Class(class), h.Span(h.Class("font-semibold"), text))
}

func novaDemo(one, two, three string, attrs ...any) h.Node {
	return resizable.PanelGroup(append(append([]any{resizable.Horizontal, h.Class("max-w-sm rounded-lg border")}, attrs...),
		resizable.Panel(resizable.DefaultSize("50%"), novaPanel("flex h-[200px] items-center justify-center p-6", one)),
		resizable.Handle(resizable.WithHandle),
		resizable.Panel(resizable.DefaultSize("50%"),
			resizable.PanelGroup(append([]any{resizable.Vertical}, append(attrs,
				resizable.Panel(resizable.DefaultSize("25%"), novaPanel("flex h-full items-center justify-center p-6", two)),
				resizable.Handle(resizable.WithHandle),
				resizable.Panel(resizable.DefaultSize("75%"), novaPanel("flex h-full items-center justify-center p-6", three)),
			)...)...),
		),
	)...)
}

func NovaDemo() h.Node { return novaDemo("One", "Two", "Three") }

func novaRtl() h.Node { return novaDemo("واحد", "اثنان", "ثلاثة", h.Dir("rtl")) }

func NovaVertical() h.Node {
	return resizable.PanelGroup(resizable.Vertical, h.Class("min-h-[200px] max-w-sm rounded-lg border"),
		resizable.Panel(resizable.DefaultSize("25%"), novaPanel("flex h-full items-center justify-center p-6", "Header")),
		resizable.Handle(),
		resizable.Panel(resizable.DefaultSize("75%"), novaPanel("flex h-full items-center justify-center p-6", "Content")),
	)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
