package tooltip

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/kbd"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the tooltip page (ui/apps/v4/examples/radix/tooltip-*.tsx),
// rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "tooltip", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("tooltip-demo", Demo),
		nova("tooltip-sides", NovaSides),
		nova("tooltip-disabled", NovaDisabled),
		nova("tooltip-keyboard", NovaKeyboard),
		nova("tooltip-rtl", NovaRtl),
	)
}

var novaSides = []tooltip.Side{tooltip.SideLeft, tooltip.SideTop, tooltip.SideBottom, tooltip.SideRight}

func NovaSides() h.Node {
	kids := []any{h.Class("flex flex-wrap gap-2")}
	for _, s := range novaSides {
		kids = append(kids, tooltip.Tooltip(
			tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("w-fit capitalize"), string(s))),
			tooltip.Content(s, h.P("Add to library")),
		))
	}
	return h.Div(kids...)
}

func NovaDisabled() h.Node {
	return tooltip.Tooltip(
		tooltip.Trigger(lib.AsChild, h.Span(h.Class("inline-block w-fit"),
			button.Button(button.VariantOutline, h.Bool("disabled"), "Disabled"))),
		tooltip.Content(h.P("This feature is currently unavailable")),
	)
}

func NovaKeyboard() h.Node {
	return tooltip.Tooltip(
		tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIconSm, lucide.Save())),
		tooltip.Content("Save Changes ", kbd.Kbd("S")),
	)
}

func NovaRtl() h.Node {
	labels := map[tooltip.Side]string{tooltip.SideLeft: "يسار", tooltip.SideTop: "أعلى", tooltip.SideBottom: "أسفل", tooltip.SideRight: "يمين"}
	kids := []any{h.Class("flex flex-wrap gap-2")}
	for _, s := range novaSides {
		kids = append(kids, tooltip.Tooltip(
			tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("w-fit capitalize"), labels[s])),
			tooltip.Content(s, "إضافة إلى المكتبة"),
		))
	}
	return h.Div(kids...)
}
