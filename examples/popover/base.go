package popover

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/popover"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The base and aria popover-rtl previews (ui/apps/v4/examples/{base,aria}/
// popover-rtl.tsx): the physical sides, then a second row with the logical
// ones (inline-start is the right side in RTL, inline-end the left).
func init() {
	rtl := func() h.Node { return style.Scope(style.Nova, BaseRtl()) }
	registry.Register(
		registry.Example{Name: "base--popover-rtl", Component: "popover", Render: registry.Static(rtl)},
		registry.Example{Name: "aria--popover-rtl", Component: "popover", Render: registry.Static(rtl)},
	)
}

// BaseRtl is examples/base/popover-rtl.tsx.
func BaseRtl() h.Node {
	type side struct {
		side  popover.Side
		label string
	}
	row := func(sides []side) h.Node {
		kids := []any{h.Class("flex flex-wrap justify-center gap-2")}
		for _, s := range sides {
			kids = append(kids, popover.Popover(
				popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, s.label)),
				popover.Content(s.side, h.Dir("rtl"),
					popover.Header(
						popover.Title("الأبعاد"),
						popover.Description("تعيين الأبعاد للطبقة."),
					),
				),
			))
		}
		return h.Div(kids...)
	}
	return h.Div(h.Class("grid gap-4"),
		row([]side{{popover.SideLeft, "يسار"}, {popover.SideTop, "أعلى"}, {popover.SideBottom, "أسفل"}, {popover.SideRight, "يمين"}}),
		row([]side{{popover.SideRight, "بداية السطر"}, {popover.SideLeft, "نهاية السطر"}}),
	)
}
