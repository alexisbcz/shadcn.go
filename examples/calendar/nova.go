package calendar

import (
	"net/http"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the calendar page
// (ui/apps/v4/examples/radix/calendar-*.tsx), rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "calendar", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("calendar-demo", NovaDemo),
		nova("calendar-basic", NovaBasic),
		nova("calendar-caption", NovaCaption),
		nova("calendar-multiple", NovaMultiple),
		nova("calendar-range", NovaRange),
		nova("calendar-week-numbers", NovaWeekNumbers),
		nova("calendar-booked-dates", NovaBookedDates),
		nova("calendar-custom-days", NovaCustomDays),
		nova("calendar-presets", NovaPresets),
		nova("calendar-time", NovaTime),
		// calendar-hijri renders its own copy of the new-york-v4 calendar.
		registry.Example{Name: "nova--calendar-hijri", Component: "calendar", Render: func(*http.Request) h.Node { return NovaHijri() }},
	)
}

func thisYear(month time.Month, day int) time.Time {
	return time.Date(time.Now().Year(), month, day, 0, 0, 0, 0, time.Local)
}

func flatCard(children ...any) h.Node {
	return card.Card(h.Class("mx-auto w-fit p-0"), card.Content(h.Class("p-0"), h.Group(children...)))
}

func NovaDemo() h.Node {
	return calendar.Calendar(calendar.ModeSingle, calendar.Selected(time.Now()), h.Class("rounded-lg border"), calendar.CaptionDropdown)
}

func NovaBasic() h.Node {
	return calendar.Calendar(calendar.ModeSingle, h.Class("rounded-lg border"))
}

func NovaCaption() h.Node {
	return calendar.Calendar(calendar.ModeSingle, calendar.CaptionDropdown, h.Class("rounded-lg border"))
}

func NovaMultiple() h.Node {
	return flatCard(calendar.Calendar(calendar.ModeMultiple))
}

func NovaRange() h.Node {
	from := thisYear(time.January, 12)
	return flatCard(calendar.Calendar(calendar.ModeRange,
		calendar.DefaultMonth(from),
		calendar.SelectedRange(from, from.AddDate(0, 0, 30)),
		calendar.NumberOfMonths(2),
		calendar.Disabled(calendar.After(time.Now()), calendar.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local))),
	))
}

func NovaWeekNumbers() h.Node {
	d := thisYear(time.February, 3)
	return flatCard(calendar.Calendar(calendar.ModeSingle, calendar.DefaultMonth(d), calendar.Selected(d), calendar.ShowWeekNumber(true)))
}

func NovaBookedDates() h.Node {
	d := thisYear(time.February, 3)
	var booked []calendar.Matcher
	for i := 0; i < 15; i++ {
		booked = append(booked, calendar.Day(thisYear(time.February, 12+i)))
	}
	return flatCard(calendar.Calendar(calendar.ModeSingle, calendar.DefaultMonth(d), calendar.Selected(d), calendar.Disabled(booked...)))
}

func NovaCustomDays() h.Node {
	from := thisYear(time.December, 8)
	return flatCard(calendar.Calendar(calendar.ModeRange,
		calendar.DefaultMonth(from),
		calendar.SelectedRange(from, from.AddDate(0, 0, 10)),
		calendar.NumberOfMonths(1),
		calendar.CaptionDropdown,
		h.Class("[--cell-size:--spacing(10)] md:[--cell-size:--spacing(12)]"),
	))
}

func NovaPresets() h.Node {
	now := time.Now()
	presets := []string{"Today", "Tomorrow", "In 3 days", "In a week", "In 2 weeks"}
	buttons := make([]any, len(presets))
	for i, p := range presets {
		buttons[i] = button.Button(button.VariantOutline, button.SizeSm, h.Class("flex-1"), p)
	}
	return card.Card(h.Class("mx-auto w-fit max-w-[300px]"), card.SizeSm,
		card.Content(
			calendar.Calendar(calendar.ModeSingle,
				calendar.Selected(thisYear(time.February, 12)),
				calendar.DefaultMonth(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)),
				calendar.FixedWeeks(true),
				h.Class("p-0 [--cell-size:--spacing(9.5)]"),
			),
		),
		card.Footer(h.Class("flex flex-wrap gap-2 border-t"), h.Group(buttons...)),
	)
}

func NovaTime() h.Node {
	now := time.Now()
	timeField := func(id, label, value string) h.Node {
		return field.Field(
			field.Label(h.For(id), label),
			inputgroup.InputGroup(
				inputgroup.Input(h.ID(id), h.Type("time"), h.Attr("step", "1"), h.Value(value),
					h.Class("appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none")),
				inputgroup.Addon(lucide.Clock2(h.Class("text-muted-foreground"))),
			),
		)
	}
	return card.Card(card.SizeSm, h.Class("mx-auto w-fit"),
		card.Content(
			calendar.Calendar(calendar.ModeSingle, calendar.Selected(time.Date(now.Year(), now.Month(), 12, 0, 0, 0, 0, time.Local)), h.Class("p-0")),
		),
		card.Footer(h.Class("border-t bg-card"),
			field.Group(
				timeField("time-from", "Start Time", "10:30:00"),
				timeField("time-to", "End Time", "12:30:00"),
			),
		),
	)
}

func NovaHijri() h.Node {
	date := time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local)
	return h.Div(h.Data("example-empty-class", ""),
		calendar.Calendar(
			calendar.Persian,
			calendar.ModeSingle,
			calendar.DefaultMonth(date),
			calendar.Selected(date),
			h.Class("rounded-lg border"),
			calendar.ClassNames{"dropdown": "absolute inset-0 opacity-0 rdp-dropdown"},
		),
	)
}
