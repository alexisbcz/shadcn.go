package scrollarea

import (
	"shadcn.go/components/style"
	scrollarea "shadcn.go/components/ui/scroll-area"
	"shadcn.go/components/ui/separator"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/scroll-area-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	registry.Register(registry.Example{Name: "nova--scroll-area-rtl", Component: "scroll-area", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, NovaRtl())
	})})
}

func novaRtl() h.Node {
	return scrollarea.ScrollArea(h.Class("h-72 w-48 rounded-md border"), h.Dir("rtl"),
		h.Div(h.Class("p-4"),
			h.H4(h.Class("mb-4 text-sm leading-none font-medium"), "العلامات"),
			h.Map(tags(50), func(tag string) h.Node {
				return h.Group(h.Div(h.Class("text-sm"), tag), separator.Separator(h.Class("my-2")))
			}),
		),
	)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
