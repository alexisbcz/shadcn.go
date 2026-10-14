package drawer

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/drawer"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the drawer page (ui/apps/v4/examples/radix/drawer-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "drawer", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("drawer-demo", NovaDemo),
		nova("drawer-scrollable-content", NovaScrollableContent),
		nova("drawer-sides", NovaSides),
	)
}

// NovaDemo is drawer-demo with the docs' bar color.
func NovaDemo() h.Node { return demo(chartBars) }

const novaLorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func novaGoalContent(contentAttrs ...any) *drawer.ContentBuilder {
	ps := []any{h.Class("no-scrollbar overflow-y-auto px-4")}
	for range 10 {
		ps = append(ps, h.P(h.Class("mb-4 leading-normal style-lyra:mb-2 style-lyra:leading-relaxed"), novaLorem))
	}
	return drawer.Content(append(contentAttrs,
		drawer.Header(
			drawer.Title("Move Goal"),
			drawer.Description("Set your daily activity goal."),
		),
		h.Div(ps...),
		drawer.Footer(
			button.Button("Submit"),
			drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
		),
	)...)
}

func NovaScrollableContent() h.Node {
	return drawer.Drawer(drawer.DirectionRight,
		drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Scrollable Content")),
		novaGoalContent(),
	)
}

func NovaSides() h.Node {
	row := []any{h.Class("flex flex-wrap gap-2")}
	for _, side := range []drawer.Direction{drawer.DirectionTop, drawer.DirectionRight, drawer.DirectionBottom, drawer.DirectionLeft} {
		row = append(row, drawer.Drawer(side,
			drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("capitalize"), string(side))),
			novaGoalContent(h.Class("data-[vaul-drawer-direction=bottom]:max-h-[50vh] data-[vaul-drawer-direction=top]:max-h-[50vh]")),
		))
	}
	return h.Div(row...)
}
