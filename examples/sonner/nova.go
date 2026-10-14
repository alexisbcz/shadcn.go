package sonner

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the sonner page (ui/apps/v4/examples/radix/sonner-*.tsx)
// without a new-york-v4 counterpart, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "sonner", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		}), Extra: toaster}
	}
	registry.Register(
		nova("sonner-description", NovaDescription),
		nova("sonner-position", NovaPosition),
	)
}

func NovaDescription() h.Node {
	return h.Group(
		button.Button(button.VariantOutline, h.Class("w-fit"), "Show Toast"),
		h.Template(h.Data("parity-ignore", ""), h.Data("example-sonner", "description")),
	)
}

func NovaPosition() h.Node {
	row := []any{h.Class("flex flex-wrap justify-center gap-2")}
	for _, label := range []string{"Top Left", "Top Center", "Top Right", "Bottom Left", "Bottom Center", "Bottom Right"} {
		row = append(row, button.Button(button.VariantOutline, label))
	}
	return h.Div(append(row, handlers())...)
}
