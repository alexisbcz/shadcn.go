// Package datepicker ports the upstream date-picker-* examples: a Calendar
// in a Popover.
package datepicker

import (
	_ "embed"
	"fmt"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed date_picker.js
var script string

func init() {
	// Keeps the picked date like the examples' useState: updates the trigger
	// and the calendar the popover mounts next time.
	runtime.Module("example-date-picker", script)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "date-picker", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("date-picker-demo", Demo),
		ex("date-picker-with-range", WithRange),
	)
}

func Demo() h.Node {
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline,
			h.Data("example-date-picker", "single"),
			h.Class(lib.Cn("w-[240px] justify-start text-left font-normal", "text-muted-foreground")),
			lucide.Calendar(),
			h.Span("Pick a date"),
		)),
		popover.Content(h.Class("w-auto p-0"), popover.AlignStart,
			calendar.Calendar(calendar.ModeSingle),
		),
	)
}

// rangeLabel is format(from, "LLL dd, y") + " -" + " " + format(to, ...),
// as separate text nodes like React renders them.
func rangeLabel(from, to time.Time) h.Node {
	f := func(t time.Time) string { return fmt.Sprintf("%s %02d, %d", t.Format("Jan"), t.Day(), t.Year()) }
	return h.Group(f(from), h.Comment(""), " -", h.Comment(""), " ", h.Comment(""), f(to))
}

func WithRange() h.Node {
	from := time.Date(2022, 1, 20, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 20)
	return h.Div(h.Class(lib.Cn("grid gap-2", "")),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline,
				h.ID("date"),
				h.Data("example-date-picker", "range"),
				h.Class(lib.Cn("w-[300px] justify-start text-left font-normal")),
				lucide.Calendar(),
				rangeLabel(from, to),
			)),
			popover.Content(h.Class("w-auto p-0"), popover.AlignStart,
				calendar.Calendar(
					calendar.ModeRange,
					calendar.DefaultMonth(from),
					calendar.SelectedRange(from, to),
					calendar.NumberOfMonths(2),
				),
			),
		),
	)
}
