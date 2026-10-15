package calendar

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// calendar-rtl (ui/apps/v4/examples/{radix,base}/calendar-rtl.tsx): the
// ui-rtl calendar (logical classes) in Arabic (arSA), today selected. The
// base page renders the same markup, through the nova-- fallback.
func init() {
	registry.Register(registry.Example{Name: "nova--calendar-rtl", Component: "calendar", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, style.ScopeRTL(NovaRtl()))
	})})
}

func NovaRtl() h.Node {
	return calendar.Calendar(calendar.ModeSingle, calendar.Selected(lib.Now()), calendar.CaptionDropdown, calendar.Arabic,
		h.Dir("rtl"), h.Class("rounded-lg border [--cell-size:--spacing(9)]"))
}
