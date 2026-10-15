package tooltip

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The base and aria tooltip-rtl previews (ui/apps/v4/examples/{base,aria}/
// tooltip-rtl.tsx): the physical sides, then a second row with the logical
// ones (inline-start is the right side in RTL, inline-end the left).
func init() {
	rtl := func() h.Node { return style.Scope(style.Nova, BaseRtl()) }
	registry.Register(
		registry.Example{Name: "base--tooltip-rtl", Component: "tooltip", Render: registry.Static(rtl)},
		registry.Example{Name: "aria--tooltip-rtl", Component: "tooltip", Render: registry.Static(rtl)},
	)
}

// BaseRtl is examples/base/tooltip-rtl.tsx.
func BaseRtl() h.Node {
	type side struct {
		side  tooltip.Side
		label string
	}
	row := func(sides []side) h.Node {
		kids := []any{h.Class("flex flex-wrap justify-center gap-2")}
		for _, s := range sides {
			kids = append(kids, tooltip.Tooltip(
				tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, s.label)),
				tooltip.Content(s.side, h.Dir("rtl"), "إضافة إلى المكتبة"),
			))
		}
		return h.Div(kids...)
	}
	return h.Div(h.Class("grid gap-4"),
		row([]side{{tooltip.SideLeft, "يسار"}, {tooltip.SideTop, "أعلى"}, {tooltip.SideBottom, "أسفل"}, {tooltip.SideRight, "يمين"}}),
		row([]side{{tooltip.SideRight, "بداية السطر"}, {tooltip.SideLeft, "نهاية السطر"}}),
	)
}
