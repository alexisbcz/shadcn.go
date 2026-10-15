package hovercard

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/hover-card"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The base and aria hover-card-rtl previews (ui/apps/v4/examples/{base,aria}/
// hover-card-rtl.tsx): the physical sides, then a second row with the
// logical ones (inline-start is the right side in RTL, inline-end the left).
func init() {
	rtl := func() h.Node { return style.Scope(style.Nova, BaseRtl()) }
	registry.Register(
		registry.Example{Name: "base--hover-card-rtl", Component: "hover-card", Render: registry.Static(rtl)},
		registry.Example{Name: "aria--hover-card-rtl", Component: "hover-card", Render: registry.Static(rtl)},
	)
}

// BaseRtl is examples/base/hover-card-rtl.tsx.
func BaseRtl() h.Node {
	type side struct {
		side  hovercard.Side
		label string
	}
	row := func(sides []side) h.Node {
		kids := []any{h.Class("flex flex-wrap justify-center gap-2")}
		for _, s := range sides {
			kids = append(kids, hovercard.HoverCard(hovercard.OpenDelay(10), hovercard.CloseDelay(100),
				hovercard.Trigger(lib.AsChild, button.Button(button.VariantOutline, s.label)),
				hovercard.Content(s.side, h.Class("flex w-64 flex-col gap-1"), h.Dir("rtl"),
					h.Div(h.Class("font-semibold"), "سماعات لاسلكية"),
					h.Div(h.Class("text-sm text-muted-foreground"), "٩٩.٩٩ $"),
				),
			))
		}
		return h.Div(kids...)
	}
	return h.Div(h.Class("grid gap-4"),
		row([]side{{hovercard.SideLeft, "يسار"}, {hovercard.SideTop, "أعلى"}, {hovercard.SideBottom, "أسفل"}, {hovercard.SideRight, "يمين"}}),
		row([]side{{hovercard.SideRight, "بداية السطر"}, {hovercard.SideLeft, "نهاية السطر"}}),
	)
}
