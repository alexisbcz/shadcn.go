package calendar

import (
	"time"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the base calendar page that differ from the radix ones
// (ui/apps/v4/examples/base/calendar-*.tsx).
func init() {
	base := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "base--" + name, Component: "calendar", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		base("calendar-range", BaseRange),
		base("calendar-week-numbers", BaseWeekNumbers),
		base("calendar-booked-dates", BaseBookedDates),
	)
}

// BaseRange has no card and no disabled days.
func BaseRange() h.Node {
	from := thisYear(time.January, 12)
	return calendar.Calendar(calendar.ModeRange,
		calendar.DefaultMonth(from),
		calendar.SelectedRange(from, from.AddDate(0, 0, 30)),
		calendar.NumberOfMonths(2),
		h.Class("rounded-lg border"),
	)
}

func BaseWeekNumbers() h.Node {
	d := thisYear(time.January, 12)
	return flatCard(calendar.Calendar(calendar.ModeSingle, calendar.DefaultMonth(d), calendar.Selected(d), calendar.ShowWeekNumber(true)))
}

func BaseBookedDates() h.Node {
	d := thisYear(time.January, 6)
	var booked []calendar.Matcher
	for i := 0; i < 15; i++ {
		booked = append(booked, calendar.Day(thisYear(time.January, 12+i)))
	}
	return flatCard(calendar.Calendar(calendar.ModeSingle, calendar.DefaultMonth(d), calendar.Selected(d), calendar.Disabled(booked...),
		calendar.Modifier("booked", "[&>button]:line-through opacity-100", booked...)))
}
