package toggle

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/toggle"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/toggle-*.tsx) whose content
// differs from the new-york-v4 example of the same name, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "toggle", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("toggle-disabled", NovaDisabled),
		nova("toggle-outline", NovaOutline),
		nova("toggle-sizes", NovaSizes),
		nova("toggle-text", NovaText),
		nova("toggle-rtl", NovaRtl),
	)
}

func NovaDisabled() h.Node {
	return h.Div(h.Class("flex flex-wrap items-center gap-2"),
		toggle.Toggle(h.Aria("label", "Toggle disabled"), h.Disabled(), "Disabled"),
		toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle disabled outline"), h.Disabled(), "Disabled"),
	)
}

func NovaOutline() h.Node {
	return h.Div(h.Class("flex flex-wrap items-center gap-2"),
		toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle italic"), lucide.Italic(), "Italic"),
		toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle bold"), lucide.Bold(), "Bold"),
	)
}

func NovaSizes() h.Node {
	return h.Div(h.Class("flex flex-wrap items-center gap-2"),
		toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle small"), toggle.SizeSm, "Small"),
		toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle default"), toggle.SizeDefault, "Default"),
		toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle large"), toggle.SizeLg, "Large"),
	)
}

func NovaText() h.Node {
	return toggle.Toggle(h.Aria("label", "Toggle italic"), lucide.Italic(), "Italic")
}

func NovaRtl() h.Node {
	return toggle.Toggle(h.Aria("label", "Toggle bookmark"), toggle.SizeSm, toggle.VariantOutline, h.Dir("rtl"),
		lucide.Bookmark(h.Class("group-aria-pressed/toggle:fill-foreground")),
		"إشارة مرجعية",
	)
}
