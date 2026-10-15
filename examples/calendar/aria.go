package calendar

import (
	_ "embed"
	"strconv"
	"strings"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed aria.js
var ariaScript string

// Docs previews of the aria calendar page
// (ui/apps/v4/examples/aria/calendar-*.tsx). The aria-nova Calendar is
// react-aria-components' Calendar/RangeCalendar
// (styles/aria-nova/ui/calendar.tsx), not react-day-picker: a grid of
// single-letter weekdays, Select dropdowns for the caption, and outside
// days shown muted. ariaCalendar renders its markup; aria.js handles
// selection and month navigation.
func init() {
	runtime.Module("example-aria-calendar", ariaScript)
	aria := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "aria--" + name, Component: "calendar", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		aria("calendar-demo", AriaDemo),
		aria("calendar-basic", AriaBasic),
		aria("calendar-range", AriaRange),
		aria("calendar-caption", AriaCaption),
		aria("calendar-presets", AriaPresets),
		aria("calendar-time", AriaTime),
		aria("calendar-booked-dates", AriaBookedDates),
		aria("calendar-custom-days", AriaCustomDays),
		aria("calendar-rtl", AriaRtl),
		aria("calendar-hijri", AriaHijri),
	)
}

func today() time.Time {
	n := lib.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
}

func date(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

func AriaDemo() h.Node {
	t := today()
	return ariaCalendar(ariaOpts{selected: &t, dropdown: true, class: "rounded-lg border"})
}

func AriaBasic() h.Node { return ariaCalendar(ariaOpts{class: "rounded-lg border"}) }

func AriaRange() h.Node {
	from := date(today().Year(), time.January, 12)
	to := from.AddDate(0, 0, 30)
	return ariaCalendar(ariaOpts{rangeStart: &from, rangeEnd: &to, months: 2, class: "rounded-lg border"})
}

func AriaCaption() h.Node {
	return ariaCalendar(ariaOpts{dropdown: true, class: "rounded-lg border"})
}

func AriaPresets() h.Node {
	t := today()
	sel := date(t.Year(), time.February, 12)
	presets := []struct {
		label string
		days  int
	}{{"Today", 0}, {"Tomorrow", 1}, {"In 3 days", 3}, {"In a week", 7}, {"In 2 weeks", 14}}
	buttons := make([]any, len(presets))
	for i, p := range presets {
		buttons[i] = button.Button(button.VariantOutline, button.SizeSm, h.Class("flex-1"), h.Data("aria-calendar-preset", strconv.Itoa(p.days)), p.label)
	}
	return card.Card(h.Class("mx-auto w-fit max-w-[300px]"), card.SizeSm, h.Data("aria-calendar-presets", ""),
		card.Content(ariaCalendar(ariaOpts{selected: &sel, month: date(t.Year(), t.Month(), 1), weeks: 6, class: "p-0 [--cell-size:--spacing(9.5)]"})),
		card.Footer(h.Class("flex flex-wrap gap-2 border-t"), h.Group(buttons...)),
	)
}

func AriaTime() h.Node {
	t := today()
	sel := date(t.Year(), t.Month(), 12)
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
		card.Content(ariaCalendar(ariaOpts{selected: &sel, class: "p-0"})),
		card.Footer(h.Class("border-t bg-card"),
			field.Group(
				timeField("time-from", "Start Time", "10:30:00"),
				timeField("time-to", "End Time", "12:30:00"),
			),
		),
	)
}

func AriaBookedDates() h.Node {
	y := today().Year()
	sel := date(y, time.February, 3)
	var booked []string
	for i := 0; i < 15; i++ {
		booked = append(booked, isoDate(date(y, time.February, 12+i)))
	}
	return flatCard(ariaCalendar(ariaOpts{selected: &sel, unavailable: booked}))
}

func AriaCustomDays() h.Node {
	from := date(today().Year(), time.December, 8)
	to := from.AddDate(0, 0, 10)
	return flatCard(ariaCalendar(ariaOpts{rangeStart: &from, rangeEnd: &to, months: 1, dropdown: true, longMonth: true, prices: true,
		class: "[--cell-size:--spacing(10)] md:[--cell-size:--spacing(12)]"}))
}

// AriaRtl is examples/aria/calendar-rtl.tsx: the ui-rtl calendar in the
// "ar" locale (weeks from Saturday, Arabic month and weekday names).
func AriaRtl() h.Node {
	t := today()
	return ariaCalendar(ariaOpts{selected: &t, dropdown: true, rtl: true, class: "rounded-lg border [--cell-size:--spacing(9)]"})
}

// AriaHijri is examples/aria/calendar-hijri.tsx: the ui-rtl calendar under
// I18nProvider locale "fa-AF-u-ca-persian", showing Jawza 1404 (May 22 to
// June 21, 2025) with the 22nd (June 12) selected. The Persian calendar is
// rendered for that month only; aria.js leaves it static.
func AriaHijri() h.Node {
	sel := date(2025, time.June, 12)
	return ariaCalendar(ariaOpts{selected: &sel, rtl: true, persian: true, month: date(2025, time.May, 22), class: "rounded-lg border"})
}

// persianDigits writes n with Extended Arabic-Indic digits.
func persianDigits(n int) string {
	var b strings.Builder
	for _, r := range strconv.Itoa(n) {
		b.WriteRune('۰' + (r - '0'))
	}
	return b.String()
}

var arabicMonths = []string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"}

// monthName is the caption's month: Intl's short (or long) name.
func (o ariaOpts) monthName(m time.Time, long bool) string {
	if o.rtl {
		return arabicMonths[m.Month()-1]
	}
	if long {
		return m.Format("January")
	}
	return m.Format("Jan")
}

func (o ariaOpts) heading(m time.Time) string {
	if o.persian {
		return "۱۴۰۴ جوزا"
	}
	return o.monthName(m, true) + " " + m.Format("2006")
}

// rtlClass rewrites the physical classes of the ui-rtl calendar.
func (o ariaOpts) rtlClass(c string) string {
	if !o.rtl {
		return c
	}
	return strings.NewReplacer(
		"rounded-r-(--cell-radius)", "rounded-e-(--cell-radius)",
		"rounded-l-(--cell-radius)", "rounded-s-(--cell-radius)",
		"after:right-0", "after:end-0",
		"after:left-0", "after:start-0",
		"has-data-[icon=inline-end]:pr-2", "has-data-[icon=inline-end]:pe-2",
		"has-data-[icon=inline-start]:pl-2", "has-data-[icon=inline-start]:ps-2",
		"py-2 pr-2 pl-2.5", "py-2 ps-2.5 pe-2",
	).Replace(c)
}

type ariaOpts struct {
	rtl     bool
	persian bool // AriaHijri's month (month is its first day)
	// month is the first visible month (default: the selection's, or
	// today's).
	month                time.Time
	months               int
	selected             *time.Time
	rangeStart, rangeEnd *time.Time
	dropdown             bool
	longMonth            bool // headerFormat {month: "long"}
	weeks                int  // weeksInMonth (0: as many as the month needs)
	unavailable          []string
	prices               bool // calendar-custom-days' renderCell
	class                string
}

func isoDate(t time.Time) string { return t.Format("2006-01-02") }

const (
	ariaCalendarClass = "group/calendar w-fit bg-background p-2 [--cell-radius:var(--radius-md)] [--cell-size:--spacing(7)] in-data-[slot=card-content]:bg-transparent in-data-[slot=popover-content]:bg-transparent"
	ariaNavClass      = "group/button inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:hover:bg-muted/50 gap-1.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 size-(--cell-size) p-0 select-none aria-disabled:opacity-50"
	ariaTriggerClass  = "flex w-full items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent py-2 pr-2 pl-2.5 text-sm whitespace-nowrap transition-colors outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 data-placeholder:text-muted-foreground data-[size=default]:h-8 data-[size=sm]:h-7 data-[size=sm]:rounded-[min(var(--radius-md),10px)] *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-1.5 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
	ariaHeadClass     = "rounded-(--cell-radius) text-[0.8rem] font-normal text-muted-foreground select-none"
	ariaCellClass     = "group/day relative mt-2 aspect-square h-full w-full cursor-default rounded-(--cell-radius) p-0 text-center select-none [&:is(:last-child>[data-selected=true])>div]:rounded-r-(--cell-radius) [&:is(:first-child>[data-selected=true])>div]:rounded-l-(--cell-radius)"
	ariaDayClass      = "group/button shrink-0 items-center justify-center rounded-lg border-transparent bg-clip-padding text-sm whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:hover:bg-muted/50 size-8 relative isolate z-10 flex aspect-square h-full w-full min-w-(--cell-size) flex-col gap-1 border-0 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-(--cell-radius) data-[range-end=true]:rounded-r-(--cell-radius) data-[range-end=true]:bg-primary data-[range-end=true]:text-primary-foreground data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-muted data-[range-middle=true]:text-foreground data-[range-start=true]:rounded-(--cell-radius) data-[range-start=true]:rounded-l-(--cell-radius) data-[range-start=true]:bg-primary data-[range-start=true]:text-primary-foreground data-[selected-single=true]:bg-primary data-[selected-single=true]:text-primary-foreground dark:hover:text-foreground [&>span]:text-xs [&>span]:opacity-70"
)

// The cellVariants of styles/aria-nova/ui/calendar.tsx, in cva order.
var ariaCellVariants = map[string]string{
	"today":       "rounded-(--cell-radius) bg-muted text-foreground data-[selected=true]:rounded-none",
	"start":       "relative isolate z-0 rounded-l-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:right-0 after:w-4 after:bg-muted",
	"end":         "relative isolate z-0 rounded-r-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:left-0 after:w-4 after:bg-muted",
	"unavailable": "text-muted-foreground opacity-50 [&>div]:line-through",
	"disabled":    "text-muted-foreground opacity-50",
	"outside":     "text-muted-foreground aria-selected:text-muted-foreground",
}

// ariaCalendar renders the calendar for its first state. The options go to
// aria.js as data attributes so it can re-render after a change.
func ariaCalendar(o ariaOpts) h.Node {
	if o.months == 0 {
		o.months = 1
	}
	if o.month.IsZero() {
		switch {
		case o.selected != nil:
			o.month = *o.selected
		case o.rangeStart != nil:
			o.month = *o.rangeStart
		default:
			o.month = today()
		}
	}
	if !o.persian {
		o.month = date(o.month.Year(), o.month.Month(), 1)
	}
	label := o.heading(o.month)
	if o.months > 1 {
		last := o.month.AddDate(0, o.months-1, 0)
		label = o.month.Format("January") + " to " + last.Format("January 2006")
	}
	root := h.Div(h.Data("slot", "calendar"), h.Class(lib.Cn(ariaCalendarClass, o.class)),
		h.Data("rac", ""), h.Attr("aria-label", label), h.Role("application"),
		h.Data("aria-calendar", ""),
		h.Data("month", isoDate(o.month)), h.Data("months", strconv.Itoa(o.months)),
		h.Data("weeks", strconv.Itoa(o.weeks)),
		h.Data("unavailable", strings.Join(o.unavailable, " ")),
	)
	if o.selected != nil {
		root.Children(h.Data("value", isoDate(*o.selected)))
	}
	if o.rangeStart != nil {
		root.Children(h.Data("range", isoDate(*o.rangeStart)+" "+isoDate(*o.rangeEnd)))
	}
	if o.dropdown {
		root.Children(h.Data("dropdown", ""))
	}
	if o.longMonth {
		root.Children(h.Data("long-month", ""))
	}
	if o.prices {
		root.Children(h.Data("prices", ""))
	}
	switch {
	case o.persian:
		root.Children(h.Data("locale", "fa"))
	case o.rtl:
		root.Children(h.Data("locale", "ar"))
	}
	chevron := "cn-rtl-flip size-4"
	if o.rtl {
		chevron = "size-4 rtl:rotate-180"
	}
	nav := func(slot, label string, icon h.Node) h.Node {
		return h.Button(h.Data("slot", "button"), h.Data("variant", "ghost"), h.Data("size", "default"), h.Class(o.rtlClass(ariaNavClass)),
			h.Data("rac", ""), h.Type("button"), h.Attr("tabindex", "0"), h.Attr("aria-label", label), h.Attr("slot", slot), icon)
	}
	inner := h.Div(h.Class("relative flex flex-col gap-4 md:flex-row"),
		h.El("header").Class("absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1").Children(
			nav("previous", "Previous", lucide.ChevronLeft(h.Class(chevron))),
			nav("next", "Next", lucide.ChevronRight(h.Class(chevron))),
		),
	)
	for i := 0; i < o.months; i++ {
		m := o.month.AddDate(0, i, 0)
		inner.Children(h.Div(h.Class("flex w-full flex-col gap-4"),
			ariaCaption(o, m),
			ariaGrid(o, m),
		))
	}
	return root.Children(
		h.Div(h.Style("border:0;clip:rect(0 0 0 0);clip-path:inset(50%);height:1px;margin:-1px;overflow:hidden;padding:0;position:absolute;width:1px;white-space:nowrap"), h.H2(label)),
		inner,
	)
}

func ariaCaption(o ariaOpts, m time.Time) h.Node {
	c := h.Div(h.Class("flex h-(--cell-size) w-full items-center justify-center gap-1 px-(--cell-size)"))
	if !o.dropdown {
		return c.Children(h.H2(h.Attr("aria-hidden", "true"), h.Class("text-sm font-medium select-none"), o.heading(m)))
	}
	month := o.monthName(m, o.longMonth)
	valueClass, indicator := "flex flex-1 text-left data-placeholder:text-muted-foreground", "pointer-events-none absolute right-2 flex size-4 items-center justify-center"
	if o.rtl {
		valueClass, indicator = "flex flex-1 text-start data-placeholder:text-muted-foreground", "pointer-events-none absolute end-2 flex size-4 items-center justify-center"
	}
	sel := func(label, value, kind string) h.Node {
		return h.Div(h.Data("slot", "select"), h.Class("w-fit relative"), h.Data("rac", ""),
			h.Button(h.Data("slot", "select-trigger"), h.Data("size", "default"), h.Class(o.rtlClass(ariaTriggerClass)), h.Data("rac", ""),
				h.Type("button"), h.Attr("tabindex", "0"), h.Attr("aria-label", label), h.Attr("aria-haspopup", "listbox"), h.Attr("aria-expanded", "false"),
				h.Data("aria-calendar-select", kind),
				h.Span(h.Data("slot", "select-value"), h.Class(valueClass), h.Data("rac", ""),
					h.Span(h.Class("flex flex-1 shrink-0 gap-2 whitespace-nowrap"), value),
					h.Span(h.Class(indicator)),
				),
				lucide.ChevronDown(h.Class("pointer-events-none size-4 text-muted-foreground")),
			),
		)
	}
	return c.Children(sel("month", month, "month"), sel("year", m.Format("2006"), "year"))
}

func ariaGrid(o ariaOpts, m time.Time) h.Node {
	head := h.Tr()
	days, weekStart := []string{"S", "M", "T", "W", "T", "F", "S"}, 0
	if o.rtl {
		days, weekStart = []string{"س", "ح", "ن", "ث", "ر", "خ", "ج"}, 6
	}
	if o.persian {
		days = []string{"ش", "ی", "د", "س", "چ", "پ", "ج"}
	}
	for _, d := range days {
		head.Children(h.Th(h.Class(ariaHeadClass), d))
	}
	first := m.AddDate(0, 0, -((int(m.Weekday()) - weekStart + 7) % 7))
	weeks := o.weeks
	visStart := o.month
	visEnd := o.month.AddDate(0, o.months, -1)
	if o.persian {
		visEnd = date(2025, time.June, 21)
	}
	if weeks == 0 {
		end := m.AddDate(0, 1, -1)
		if o.persian {
			end = visEnd
		}
		weeks = (int(end.Sub(first).Hours()/24) + 7) / 7
	}
	unavailable := map[string]bool{}
	for _, d := range o.unavailable {
		unavailable[d] = true
	}
	t := today()
	focus := t
	switch {
	case o.selected != nil:
		focus = *o.selected
	case o.rangeStart != nil:
		focus = *o.rangeStart
	}
	body := h.Tbody(h.Class("react-aria-CalendarGridBody"))
	for w := 0; w < weeks; w++ {
		row := h.Tr()
		for d := 0; d < 7; d++ {
			day := first.AddDate(0, 0, w*7+d)
			row.Children(ariaCell(o, day, m, visStart, visEnd, unavailable[isoDate(day)], t, focus))
		}
		body.Children(row)
	}
	return h.Table(h.Attr("aria-label", o.heading(m)), h.Role("grid"), h.Attr("cellpadding", "0"), h.Class("w-full border-collapse"),
		h.Thead(h.Attr("aria-hidden", "true"), h.Class("react-aria-CalendarGridHeader"), head),
		body,
	)
}

func ariaCell(o ariaOpts, day, m, visStart, visEnd time.Time, unavailable bool, today, focus time.Time) h.Node {
	isRange := o.rangeStart != nil
	outside := day.Month() != m.Month() || day.Year() != m.Year()
	if o.persian {
		outside = day.Before(visStart) || day.After(visEnd)
	}
	disabled := day.Before(visStart) || day.After(visEnd)
	isToday := day.Equal(today)
	selected, start, end := false, false, false
	if o.selected != nil && day.Equal(*o.selected) {
		selected = true
	}
	if isRange && !day.Before(*o.rangeStart) && !day.After(*o.rangeEnd) && !unavailable {
		selected = true
		start = day.Equal(*o.rangeStart)
		end = day.Equal(*o.rangeEnd)
	}
	if disabled {
		selected, start, end = false, false, false
	}
	class := ariaCellClass
	add := func(on bool, k string) {
		if on {
			class += " " + ariaCellVariants[k]
		}
	}
	add(isToday, "today")
	add(start, "start")
	add(end, "end")
	add(unavailable, "unavailable")
	add(disabled, "disabled")
	add(outside, "outside")

	cell := h.Div(h.Data("react-aria-pressable", "true"))
	if !disabled {
		tab := "-1"
		if day.Equal(focus) {
			tab = "0"
		}
		cell.Children(h.Attr("tabindex", tab))
	}
	cell.Children(h.Role("button"))
	if disabled || unavailable {
		cell.Children(h.Attr("aria-disabled", "true"))
	}
	cell.Children(h.Attr("aria-label", day.Format("Monday, January 2, 2006")), h.Data("date", isoDate(day)))
	if selected {
		cell.Children(h.Data("selected", "true"))
	}
	if start {
		cell.Children(h.Data("selection-start", "true"))
	}
	if end {
		cell.Children(h.Data("selection-end", "true"))
	}
	if isToday {
		cell.Children(h.Data("today", "true"))
	}
	if disabled {
		cell.Children(h.Data("disabled", "true"), h.Data("outside-visible-range", "true"))
	}
	if outside {
		cell.Children(h.Data("outside-month", "true"))
	}
	if unavailable {
		cell.Children(h.Data("unavailable", "true"))
	}
	cell.Children(h.Class(o.rtlClass(class)), h.Data("rac", ""))

	inner := h.Div(h.Data("selected-single", boolAttr(selected && !isRange)),
		h.Data("range-start", boolAttr(start && isRange)), h.Data("range-end", boolAttr(end && isRange)))
	if !selected || isRange {
		inner.Children(h.Data("range-middle", boolAttr(selected && !start && !end && isRange)))
	}
	label := strconv.Itoa(day.Day())
	if o.persian {
		// Sawr has 31 days; Jawza 1 is May 22.
		n := int(day.Sub(visStart).Hours()/24) + 1
		switch {
		case n < 1:
			n += 31
		case n > 31:
			n -= 31
		}
		label = persianDigits(n)
	}
	inner.Children(h.Class(o.rtlClass(ariaDayClass)), label)
	if o.prices && !outside {
		price := "$100"
		if wd := day.Weekday(); wd == time.Saturday || wd == time.Sunday {
			price = "$120"
		}
		inner.Children(h.Span(price))
	}
	cell.Children(inner)

	td := h.Td(h.Role("gridcell"))
	if disabled || unavailable {
		td.Children(h.Attr("aria-disabled", "true"))
	}
	if selected {
		td.Children(h.Attr("aria-selected", "true"))
	}
	return td.Children(cell)
}

func boolAttr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
