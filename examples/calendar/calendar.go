// Package calendar ports the upstream calendar-* examples.
package calendar

import (
	_ "embed"
	"shadcn.go/components/lib"
	"time"

	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed empty_class.js
var emptyClassScript string

func init() {
	runtime.Module("example-empty-class", emptyClassScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "calendar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("calendar-demo", Demo),
		ex("calendar-hijri", Hijri),
	)
}

func Demo() h.Node {
	return calendar.Calendar(
		calendar.ModeSingle,
		calendar.Selected(lib.Now()),
		h.Class("rounded-md border shadow-sm"),
		calendar.CaptionDropdown,
	)
}

// Hijri uses the Persian calendar. Upstream wraps it in a div carrying the
// Vazirmatn font's class name, which is empty in the reference build.
func Hijri() h.Node {
	date := time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local)
	return h.Div(h.Data("example-empty-class", ""),
		calendar.Calendar(
			calendar.Persian,
			calendar.ModeSingle,
			calendar.DefaultMonth(date),
			calendar.Selected(date),
			h.Class("rounded-lg border shadow-sm"),
			// The example's own copy of the calendar drops bg-popover.
			calendar.ClassNames{"dropdown": "absolute inset-0 opacity-0 rdp-dropdown"},
		),
	)
}
