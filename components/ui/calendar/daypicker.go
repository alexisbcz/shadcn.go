package calendar

import (
	"slices"
	"strconv"
	"strings"
	"time"

	h "shadcn.go/libhtml"
)

func itoa(i int) string { return strconv.Itoa(i) }

// The server-side DayPicker: renders the first state exactly like
// react-day-picker's first render (calendar.js mirrors it on the client).

type calDay struct {
	date, displayMonth int
	outside            bool
}

type calWeek struct {
	weekNumber int
	days       []calDay
}

type calMonth struct {
	date  int
	weeks []calWeek
}

type dayPicker struct {
	cfg              *config
	lib              *dateLib
	today            int
	navStart, navEnd *int
	firstMonth       int
}

func newDayPicker(cfg *config, l *dateLib) *dayPicker {
	dp := &dayPicker{cfg: cfg, lib: l, today: JDN(time.Now())}
	if cfg.Today != nil {
		dp.today = *cfg.Today
	}
	dp.navMonths()
	dp.firstMonth = dp.initialMonth()
	return dp
}

func (dp *dayPicker) numberOfMonths() int {
	if dp.cfg.NumberOfMonths > 0 {
		return dp.cfg.NumberOfMonths
	}
	return 1
}

func (dp *dayPicker) navMonths() {
	cfg, l := dp.cfg, dp.lib
	hasYearDropdown := cfg.CaptionLayout == string(CaptionDropdown) || cfg.CaptionLayout == string(CaptionDropdownYears)
	if cfg.StartMonth != nil {
		dp.navStart = ptr(l.startOfMonth(*cfg.StartMonth))
	} else if hasYearDropdown {
		dp.navStart = ptr(l.startOfYear(l.addYears(dp.today, -100)))
	}
	if cfg.EndMonth != nil {
		dp.navEnd = ptr(l.endOfMonth(*cfg.EndMonth))
	} else if hasYearDropdown {
		dp.navEnd = ptr(l.endOfYear(dp.today))
	}
}

func (dp *dayPicker) initialMonth() int {
	l, n := dp.lib, dp.numberOfMonths()
	initial := dp.today
	if dp.cfg.DefaultMonth != nil {
		initial = *dp.cfg.DefaultMonth
	}
	if dp.navEnd != nil && l.diffMonths(*dp.navEnd, initial) < n-1 {
		initial = l.addMonths(*dp.navEnd, -1*(n-1))
	}
	if dp.navStart != nil && l.diffMonths(initial, *dp.navStart) < 0 {
		initial = *dp.navStart
	}
	return l.startOfMonth(initial)
}

func (dp *dayPicker) calendar() ([]calMonth, []calDay) {
	cfg, l := dp.cfg, dp.lib
	n := dp.numberOfMonths()
	var displayMonths []int
	for i := 0; i < n; i++ {
		m := l.addMonths(dp.firstMonth, i)
		if dp.navEnd != nil && m > *dp.navEnd {
			break
		}
		displayMonths = append(displayMonths, m)
	}
	first, last := displayMonths[0], displayMonths[len(displayMonths)-1]
	startDate := l.startOfWeek(first)
	endDate := l.endOfWeek(l.endOfMonth(last))
	var dates []int
	for j := startDate; j <= endDate; j++ {
		if cfg.EndMonth != nil && j > l.endOfMonth(*cfg.EndMonth) {
			break
		}
		dates = append(dates, j)
	}
	nMonths := l.diffMonths(last, first) + 1
	if cfg.FixedWeeks && len(dates) < 42*nMonths {
		add := 42*nMonths - len(dates)
		for i := 0; i < add; i++ {
			dates = append(dates, dates[len(dates)-1]+1)
		}
	}
	var months []calMonth
	var days []calDay
	for _, month := range displayMonths {
		firstOfWeek := l.startOfWeek(month)
		lastOfWeek := l.endOfWeek(l.endOfMonth(month))
		var monthDates []int
		for _, d := range dates {
			if d >= firstOfWeek && d <= lastOfWeek {
				monthDates = append(monthDates, d)
			}
		}
		if cfg.FixedWeeks && len(monthDates) < 42 {
			add := 42 - len(monthDates)
			for _, d := range dates {
				if d > lastOfWeek && d <= lastOfWeek+add {
					monthDates = append(monthDates, d)
				}
			}
		}
		var weeks []calWeek
		for _, date := range monthDates {
			wn := l.getWeek(date)
			day := calDay{date: date, displayMonth: month, outside: !l.isSameMonth(date, month)}
			i := slices.IndexFunc(weeks, func(w calWeek) bool { return w.weekNumber == wn })
			if i < 0 {
				weeks = append(weeks, calWeek{weekNumber: wn, days: []calDay{day}})
			} else {
				weeks[i].days = append(weeks[i].days, day)
			}
		}
		months = append(months, calMonth{date: month, weeks: weeks})
		for _, w := range weeks {
			days = append(days, w.days...)
		}
	}
	return months, days
}

func (dp *dayPicker) previousMonth() *int {
	cfg, l := dp.cfg, dp.lib
	if cfg.DisableNavigation {
		return nil
	}
	offset := 1
	if cfg.PagedNavigation {
		offset = dp.numberOfMonths()
	}
	month := l.startOfMonth(dp.firstMonth)
	if dp.navStart != nil && l.diffMonths(month, *dp.navStart) <= 0 {
		return nil
	}
	return ptr(l.addMonths(month, -offset))
}

func (dp *dayPicker) nextMonth() *int {
	cfg, l := dp.cfg, dp.lib
	if cfg.DisableNavigation {
		return nil
	}
	n := dp.numberOfMonths()
	offset := 1
	if cfg.PagedNavigation {
		offset = n
	}
	month := l.startOfMonth(dp.firstMonth)
	if dp.navEnd != nil && l.diffMonths(*dp.navEnd, dp.firstMonth) < n {
		return nil
	}
	return ptr(l.addMonths(month, offset))
}

// --- matchers -----------------------------------------------------------------------

func rangeIncludes(from, to *int, j int, excludeEnds bool) bool {
	e := 0
	if excludeEnds {
		e = 1
	}
	if from != nil && to != nil {
		f, t := *from, *to
		if t-f < 0 {
			f, t = t, f
		}
		return j-f >= e && t-j >= e
	}
	if !excludeEnds && to != nil {
		return *to == j
	}
	if !excludeEnds && from != nil {
		return *from == j
	}
	return false
}

func (dp *dayPicker) matches(j int, matchers []any) bool {
	for _, m := range matchers {
		switch v := m.(type) {
		case bool:
			if v {
				return true
			}
		case int:
			if v == j {
				return true
			}
		case map[string]*int:
			if rangeIncludes(v["from"], v["to"], j, false) {
				return true
			}
		case map[string][]int:
			if slices.Contains(v["dayOfWeek"], dp.lib.weekday(j)) {
				return true
			}
		case map[string]int:
			before, hasBefore := v["before"]
			after, hasAfter := v["after"]
			switch {
			case hasBefore && hasAfter:
				isDayBefore, isDayAfter := before-j > 0, after-j < 0
				if before > after {
					if isDayAfter && isDayBefore {
						return true
					}
				} else if isDayBefore || isDayAfter {
					return true
				}
			case hasAfter:
				if j-after > 0 {
					return true
				}
			case hasBefore:
				if before-j > 0 {
					return true
				}
			}
		}
	}
	return false
}

type modifiers struct {
	disabled, hidden, outside, today, focused bool
	// selected and the range flags are tri-state like in React
	// (undefined, false, true): nil is undefined.
	selected                          *bool
	rangeStart, rangeEnd, rangeMiddle *bool
}

func (dp *dayPicker) modifiersOf(day calDay) modifiers {
	l := dp.lib
	var navStart, navEnd *int
	if dp.navStart != nil {
		navStart = ptr(l.startOfMonth(*dp.navStart))
	}
	if dp.navEnd != nil {
		navEnd = ptr(l.endOfMonth(*dp.navEnd))
	}
	hidden := dp.matches(day.date, dp.cfg.Hidden) ||
		(navStart != nil && day.date < *navStart) ||
		(navEnd != nil && day.date > *navEnd) ||
		(!dp.cfg.ShowOutsideDays && day.outside)
	return modifiers{
		disabled: dp.matches(day.date, dp.cfg.Disabled),
		hidden:   hidden,
		outside:  day.outside,
		today:    day.date == dp.today,
	}
}

func (dp *dayPicker) isSelected(j int) bool {
	cfg := dp.cfg
	if cfg.Mode == nil {
		return false
	}
	switch Mode(*cfg.Mode) {
	case ModeSingle:
		return cfg.selectedSingle != nil && *cfg.selectedSingle == j
	case ModeMultiple:
		return slices.Contains(cfg.selectedMulti, j)
	case ModeRange:
		if cfg.selectedRange == nil {
			return false
		}
		r := rangeJSON(*cfg.selectedRange)
		return rangeIncludes(r["from"], r["to"], j, false)
	}
	return false
}

func (dp *dayPicker) focusTarget(days []calDay) *calDay {
	var target *calDay
	priority := -1
	for i, day := range days {
		m := dp.modifiersOf(day)
		if m.disabled || m.hidden || m.outside {
			continue
		}
		if dp.isSelected(day.date) && priority < 1 {
			target, priority = &days[i], 1
		} else if m.today && priority < 0 {
			target, priority = &days[i], 0
		}
	}
	if target == nil {
		for i, day := range days {
			m := dp.modifiersOf(day)
			if !m.disabled && !m.hidden && !m.outside {
				return &days[i]
			}
		}
	}
	return target
}

// --- render -----------------------------------------------------------------------------

const srOnlyStyle = "border: 0px; clip: rect(0px, 0px, 0px, 0px); height: 1px; margin: -1px; overflow: hidden; padding: 0px; position: absolute; width: 1px; white-space: nowrap; overflow-wrap: normal;"

func (dp *dayPicker) render() []any {
	cfg, l := dp.cfg, dp.lib
	cn := cfg.ClassNames
	months, days := dp.calendar()
	target := dp.focusTarget(days)
	interactive := cfg.Mode != nil
	prev, next := dp.previousMonth(), dp.nextMonth()

	navButton := func(class, label, orientation string, month *int) *h.Element {
		btn := h.Button(h.Type("button"), h.Class(class), h.Aria("label", label),
			chevron(cfg, orientation, 24, month == nil))
		if month == nil {
			btn.Attr("tabindex", "-1").Attr("aria-disabled", "true")
		}
		return btn
	}
	nav := h.El("nav", h.Class(cn["nav"]), h.Aria("label", ""),
		navButton(cn["button_previous"], "Go to the Previous Month", "left", prev),
		navButton(cn["button_next"], "Go to the Next Month", "right", next),
	)

	ws := l.startOfWeek(JDN(time.Now()))
	children := []any{nav}
	for _, month := range months {
		children = append(children, dp.renderMonth(month, ws, target, interactive))
	}
	return []any{h.Div(h.Class(cn["months"]), h.Group(children...))}
}

type option struct {
	value    int
	label    string
	disabled bool
}

func (dp *dayPicker) dropdown(kind string, options []option, value int, label string) *h.Element {
	cfg := dp.cfg
	cn := cfg.ClassNames
	kindClass := cn["months_dropdown"]
	if kind == "years" {
		kindClass = cn["years_dropdown"]
	}
	sel := h.El("select", h.Class(cn["dropdown"]+" "+kindClass), h.Aria("label", label))
	if cfg.DisableNavigation {
		sel.Bool("disabled", true)
	}
	selectedLabel := ""
	for _, o := range options {
		opt := h.El("option", h.Attr("value", itoa(o.value)), o.label)
		if o.disabled {
			opt.Bool("disabled", true)
		}
		sel.Children(opt)
		if o.value == value {
			selectedLabel = o.label
		}
	}
	caption := h.Span(h.Class(cn["caption_label"]), h.Aria("hidden", "true"))
	if selectedLabel != "" {
		caption.Children(selectedLabel)
	}
	caption.Children(chevron(cfg, "down", 18, false))
	return h.Span(h.Data("disabled", boolString(cfg.DisableNavigation)), h.Class(cn["dropdown_root"]), sel, caption)
}

func (dp *dayPicker) renderMonth(month calMonth, weekStart int, target *calDay, interactive bool) *h.Element {
	cfg, l := dp.cfg, dp.lib
	cn := cfg.ClassNames
	caption := l.caption(month.date)
	var captionEl *h.Element
	if strings.HasPrefix(cfg.CaptionLayout, "dropdown") {
		var monthOptions []option
		sy := l.startOfYear(month.date)
		for i := 0; i < 12; i++ {
			m := l.addMonths(sy, i)
			disabled := (dp.navStart != nil && m < l.startOfMonth(*dp.navStart)) || (dp.navEnd != nil && m > l.startOfMonth(*dp.navEnd))
			monthOptions = append(monthOptions, option{value: l.getMonth(m), label: dp.monthDropdown(m), disabled: disabled})
		}
		var yearOptions []option
		if dp.navStart != nil && dp.navEnd != nil {
			lastYear := l.endOfYear(*dp.navEnd)
			for y := l.startOfYear(*dp.navStart); y < lastYear || l.isSameYear(y, lastYear); y = l.addYears(y, 1) {
				yearOptions = append(yearOptions, option{value: l.getYear(y), label: l.year(y)})
			}
		}
		layout := CaptionLayout(cfg.CaptionLayout)
		var monthsPart, yearsPart *h.Element
		if layout == CaptionDropdown || layout == CaptionDropdownMonths {
			monthsPart = dp.dropdown("months", monthOptions, l.getMonth(month.date), "Choose the Month")
		} else {
			monthsPart = h.Span(dp.monthDropdown(month.date))
		}
		if layout == CaptionDropdown || layout == CaptionDropdownYears {
			yearsPart = dp.dropdown("years", yearOptions, l.getYear(month.date), "Choose the Year")
		} else {
			yearsPart = h.Span(l.year(month.date))
		}
		captionEl = h.Div(h.Class(cn["dropdowns"]), monthsPart, yearsPart,
			h.Span(h.Role("status"), h.Aria("live", "polite"), h.Style(srOnlyStyle), caption))
	} else {
		captionEl = h.Span(h.Class(cn["caption_label"]), h.Role("status"), h.Aria("live", "polite"), caption)
	}

	multi := "false"
	if cfg.Mode != nil && (Mode(*cfg.Mode) == ModeMultiple || Mode(*cfg.Mode) == ModeRange) {
		multi = "true"
	}
	table := h.El("table", h.Role("grid"), h.Aria("multiselectable", multi), h.Class(cn["month_grid"]))
	if caption != "" {
		table.Attr("aria-label", caption)
	}
	if !cfg.HideWeekdays {
		tr := h.El("tr", h.Class(cn["weekdays"]))
		if cfg.ShowWeekNumber {
			tr.Children(h.El("th", h.Aria("label", "Week Number"), h.Class(cn["week_number_header"]), h.Attr("scope", "col")))
		}
		for i := 0; i < 7; i++ {
			w := weekStart + i
			tr.Children(h.El("th", h.Aria("label", l.weekdayLong(w)), h.Class(cn["weekday"]), h.Attr("scope", "col"), l.weekdayShort(w)))
		}
		table.Children(h.El("thead", h.Aria("hidden", "true"), tr))
	}
	tbody := h.El("tbody", h.Class(cn["weeks"]))
	for _, week := range month.weeks {
		tr := h.El("tr", h.Class(cn["week"]))
		if cfg.ShowWeekNumber {
			tr.Children(h.El("td", h.Attr("week", "[object Object]"), h.Aria("label", "Week "+itoa(week.weekNumber)),
				h.Class(cn["week_number"]), h.Attr("scope", "row"), h.Role("rowheader"),
				h.Div(h.Class("flex size-(--cell-size) items-center justify-center text-center"), l.weekNumber(week.weekNumber))))
		}
		for _, day := range week.days {
			tr.Children(dp.renderDay(day, target, interactive))
		}
		tbody.Children(tr)
	}
	table.Children(tbody)
	return h.Div(h.Class(cn["month"]), h.Div(h.Class(cn["month_caption"]), captionEl), table)
}

func dataBool(b *bool) string {
	if b == nil {
		return ""
	}
	return boolString(*b)
}

func (dp *dayPicker) renderDay(day calDay, target *calDay, interactive bool) *h.Element {
	cfg, l := dp.cfg, dp.lib
	cn := cfg.ClassNames
	m := dp.modifiersOf(day)
	if dp.isSelected(day.date) {
		m.selected = ptr(true)
	}
	if cfg.Mode != nil && Mode(*cfg.Mode) == ModeRange && cfg.selectedRange != nil {
		r := rangeJSON(*cfg.selectedRange)
		from, to := r["from"], r["to"]
		m.rangeStart = ptr(from != nil && to != nil && day.date == *from)
		m.rangeEnd = ptr(from != nil && to != nil && day.date == *to)
		m.rangeMiddle = ptr(rangeIncludes(from, to, day.date, true))
	}
	isTrue := func(b *bool) bool { return b != nil && *b }
	// Modifier classes, in getModifiers' key order.
	classes := []string{cn["day"]}
	add := func(on bool, key string) {
		if on && cn[key] != "" {
			classes = append(classes, cn[key])
		}
	}
	add(m.focused, "focused")
	add(m.disabled, "disabled")
	add(m.hidden, "hidden")
	add(m.outside, "outside")
	add(m.today, "today")
	add(isTrue(m.selected), "selected")
	add(isTrue(m.rangeStart), "range_start")
	add(isTrue(m.rangeEnd), "range_end")
	add(isTrue(m.rangeMiddle), "range_middle")
	for _, mod := range cfg.Modifiers {
		if mod.Class != "" && dp.matches(day.date, mod.Match) {
			classes = append(classes, mod.Class)
		}
	}

	label := l.fullDate(day.date)
	if m.today {
		label = "Today, " + label
	}
	td := h.El("td", h.Class(strings.Join(classes, " ")), h.Role("gridcell"), h.Data("day", l.isoDay(day.date)))
	if isTrue(m.selected) {
		td.Aria("selected", "true").Data("selected", "true")
	}
	if !interactive && !m.hidden {
		td.Aria("label", label)
	}
	if day.outside {
		td.Data("month", l.isoMonth(day.date)).Data("outside", "true")
	}
	if m.disabled {
		td.Data("disabled", "true")
	}
	if m.hidden {
		td.Data("hidden", "true")
	}
	if m.today {
		td.Data("today", "true")
	}
	if m.hidden {
		return td
	}
	if !interactive {
		return td.Children(l.day(day.date))
	}
	// data-selected-single={selected && !range_start && !range_end && !range_middle}
	var selectedSingle *bool
	if m.selected != nil {
		selectedSingle = ptr(*m.selected && !isTrue(m.rangeStart) && !isTrue(m.rangeEnd) && !isTrue(m.rangeMiddle))
	}
	btnLabel := label
	if isTrue(m.selected) {
		btnLabel += ", selected"
	}
	tabindex := "-1"
	if target != nil && target.date == day.date && l.isSameMonth(target.displayMonth, day.displayMonth) {
		tabindex = "0"
	}
	btn := h.Button(
		h.Data("slot", "button"),
		h.Data("variant", "ghost"),
		h.Data("size", "icon"),
		h.Class(cfg.DayButtonClass),
		h.Data("day", l.localeDate(day.date)),
		h.Type("button"),
		h.Attr("tabindex", tabindex),
		h.Aria("label", btnLabel),
		l.day(day.date),
	)
	if !day.outside && len(cfg.DayExtra) == 7 {
		if x := cfg.DayExtra[l.weekday(day.date)]; x != "" {
			btn.Children(h.Span(x))
		}
	}
	for _, kv := range []struct {
		name string
		v    *bool
	}{{"selected-single", selectedSingle}, {"range-start", m.rangeStart}, {"range-end", m.rangeEnd}, {"range-middle", m.rangeMiddle}} {
		if kv.v != nil {
			btn.Data(kv.name, dataBool(kv.v))
		}
	}
	if m.disabled {
		btn.Bool("disabled", true)
	}
	return td.Children(btn)
}

// monthDropdown is formatMonthDropdown: short month names, or long ones
// with MonthDropdownLong.
func (dp *dayPicker) monthDropdown(j int) string {
	if dp.cfg.MonthLong && !dp.lib.fa {
		_, gm, _ := d2g(j)
		return enMonths[gm-1]
	}
	return dp.lib.monthDropdown(j)
}
