package resizable

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/resizable"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The base and aria resizable-handle previews
// (ui/apps/v4/examples/{base,aria}/resizable-handle.tsx): a narrower group
// than the radix one (max-w-sm, no minimum width).
func init() {
	handle := func() h.Node { return style.Scope(style.Nova, BaseHandle()) }
	registry.Register(
		registry.Example{Name: "base--resizable-handle", Component: "resizable", Render: registry.Static(handle)},
		registry.Example{Name: "aria--resizable-handle", Component: "resizable", Render: registry.Static(handle)},
	)
}

// BaseHandle is examples/base/resizable-handle.tsx.
func BaseHandle() h.Node {
	return resizable.PanelGroup(resizable.Horizontal, h.Class("min-h-[200px] max-w-sm rounded-lg border"),
		resizable.Panel(resizable.DefaultSize("25%"), label("flex h-full items-center justify-center p-6", "Sidebar")),
		resizable.Handle(resizable.WithHandle),
		resizable.Panel(resizable.DefaultSize("75%"), label("flex h-full items-center justify-center p-6", "Content")),
	)
}
