// Package calendar is the shadcn/ui Calendar: a date field component that
// lets users enter and edit dates, built on react-day-picker 9.8.1 upstream.
//
//	calendar.Calendar(calendar.ModeSingle, calendar.Selected(time.Now()),
//	    calendar.CaptionDropdown, h.Class("rounded-md border shadow-sm"))
//
// Children may include the DayPicker props as typed values: a Mode, the
// selection (Selected, SelectedDates, SelectedRange), DefaultMonth, Today,
// NumberOfMonths, a CaptionLayout, ShowOutsideDays, FixedWeeks,
// ShowWeekNumber, HideWeekdays, Required, DisableNavigation,
// PagedNavigation, StartMonth, EndMonth, Disabled and Hidden matchers,
// WeekStartsOn, Min, Max, ExcludeDisabled, ButtonVariant, ClassNames and
// Persian (react-day-picker/persian: the Solar Hijri calendar, Persian
// locale and digits, right to left), plus any attribute for the root.
//
// Behavior (calendar.js): react-day-picker's, ported to vanilla JS: month
// navigation (buttons and caption dropdowns), selection in each mode,
// roving focus with the full keyboard map (arrows, Shift+arrows,
// PageUp/PageDown, Shift+PageUp/PageDown, Home, End; mirrored in right to
// left), and the focus target (tabindex) rules. The root dispatches
// "shadcn:select" (detail { selected, date } in ISO dates) and
// "shadcn:month-change" (detail { month }).
//
// Dates are calendar days: a time.Time contributes its year, month and day
// in its own location. Today defaults to the server's local date; the
// browser recomputes it on mount, as React does.
package calendar

import (
	_ "embed"
	"encoding/json"
	"io"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed calendar.js
var script string

func init() { runtime.Module("calendar", script) }

// Mode is the selection mode (DayPicker's mode prop). Without a mode the
// calendar is not interactive.
type Mode string

// Modes.
const (
	ModeSingle   Mode = "single"
	ModeMultiple Mode = "multiple"
	ModeRange    Mode = "range"
)

// CaptionLayout selects the month caption: a label (default) or dropdowns.
type CaptionLayout string

// Caption layouts.
const (
	CaptionLabel          CaptionLayout = "label"
	CaptionDropdown       CaptionLayout = "dropdown"
	CaptionDropdownMonths CaptionLayout = "dropdown-months"
	CaptionDropdownYears  CaptionLayout = "dropdown-years"
)

// Numerals is a numbering system for the digits ("latn", "arab",
// "arabext", "deva", "beng").
type Numerals string

// Boolean props.
type (
	ShowOutsideDays   bool
	FixedWeeks        bool
	ShowWeekNumber    bool
	HideWeekdays      bool
	Required          bool
	DisableNavigation bool
	PagedNavigation   bool
	ExcludeDisabled   bool
)

// NumberOfMonths is the number of months displayed.
type NumberOfMonths int

// WeekStartsOn is the first day of the week (0 = Sunday).
type WeekStartsOn int

// Min and Max bound the selection (days in a range, dates in multiple mode).
type (
	Min int
	Max int
)

// ButtonVariant is the variant of the navigation buttons (default ghost).
type ButtonVariant button.Variant

// ClassNames override react-day-picker class names (keys like "day",
// "today", "caption_label"), replacing shadcn's.
type ClassNames map[string]string

type persianFlag struct{}

// Persian renders the Solar Hijri calendar like react-day-picker/persian:
// Persian locale, Eastern Arabic digits and dir="rtl".
var Persian = persianFlag{}

type (
	selectedDate   time.Time
	selectedDates  []time.Time
	defaultMonth   time.Time
	todayDate      time.Time
	startMonth     time.Time
	endMonth       time.Time
	disabledMatch  []Matcher
	hiddenMatch    []Matcher
	numeralsOption Numerals
)

// Selected is the selected day in single mode.
func Selected(t time.Time) any { return selectedDate(t) }

// SelectedDates are the selected days in multiple mode.
func SelectedDates(ts ...time.Time) any { return selectedDates(ts) }

// Range is a range of days; a zero From or To is unset.
type Range struct{ From, To time.Time }

// SelectedRange is the selected range in range mode.
func SelectedRange(from, to time.Time) any { return Range{From: from, To: to} }

// DefaultMonth is the month displayed first.
func DefaultMonth(t time.Time) any { return defaultMonth(t) }

// Today overrides today's date.
func Today(t time.Time) any { return todayDate(t) }

// StartMonth and EndMonth bound navigation.
func StartMonth(t time.Time) any { return startMonth(t) }
func EndMonth(t time.Time) any   { return endMonth(t) }

// Disabled marks the days matching any matcher as disabled.
func Disabled(m ...Matcher) any { return disabledMatch(m) }

// Hidden hides the days matching any matcher.
func Hidden(m ...Matcher) any { return hiddenMatch(m) }

// WithNumerals sets the numbering system.
func WithNumerals(n Numerals) any { return numeralsOption(n) }

// Matcher matches days (react-day-picker's Matcher).
type Matcher interface{ matcher() any }

// Day matches one day.
type Day time.Time

// Before matches the days before a day; After the days after it.
type (
	Before time.Time
	After  time.Time
)

// Interval matches the days strictly between After and Before (or outside
// them when Before is not after After), like { before, after }.
type Interval struct{ Before, After time.Time }

// DayOfWeek matches days of the week (0 = Sunday).
type DayOfWeek []int

// Always matches every day (the boolean matcher true).
type Always struct{}

func (d Day) matcher() any    { return JDN(time.Time(d)) }
func (b Before) matcher() any { return map[string]int{"before": JDN(time.Time(b))} }
func (a After) matcher() any  { return map[string]int{"after": JDN(time.Time(a))} }
func (i Interval) matcher() any {
	return map[string]int{"before": JDN(i.Before), "after": JDN(i.After)}
}
func (d DayOfWeek) matcher() any { return map[string][]int{"dayOfWeek": d} }
func (Always) matcher() any      { return true }
func (r Range) matcher() any     { return rangeJSON(r) }

func rangeJSON(r Range) map[string]*int {
	out := map[string]*int{"from": nil, "to": nil}
	if !r.From.IsZero() {
		j := JDN(r.From)
		out["from"] = &j
	}
	if !r.To.IsZero() {
		j := JDN(r.To)
		out["to"] = &j
	}
	return out
}

// Builder is returned by Calendar.
type Builder struct {
	h.Builder[Builder]
	cfg           config
	buttonVariant button.Variant
	classNames    ClassNames
	persian       bool
	numerals      Numerals
	st            style.Style // the enclosing style.Scope, "" for new-york-v4
}

// The classes of registry/bases/radix/ui/calendar.tsx, used inside a
// style.Scope (see package style). Upstream's transform only resolves the
// tokens of className attributes: those in the classNames object stay as
// they are (cn-calendar-dropdown-root, ...), and so does cn-rtl-flip.
var (
	file        = style.NewFile("ui/calendar.tsx")
	stRoot      = file.Cn("cn-calendar group/calendar bg-background in-data-[slot=card-content]:bg-transparent in-data-[slot=popover-content]:bg-transparent", `rtl:**:[.rdp-button\_next>svg]:rotate-180`, `rtl:**:[.rdp-button\_previous>svg]:rotate-180`, style.Arg)
	stDayButton = file.Cn("cn-calendar-day-button relative isolate z-10 flex aspect-square size-auto w-full min-w-(--cell-size) flex-col gap-1 border-0 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-(--cell-radius) data-[range-end=true]:rounded-r-(--cell-radius) data-[range-end=true]:bg-primary data-[range-end=true]:text-primary-foreground data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-muted data-[range-middle=true]:text-foreground data-[range-start=true]:rounded-(--cell-radius) data-[range-start=true]:rounded-l-(--cell-radius) data-[range-start=true]:bg-primary data-[range-start=true]:text-primary-foreground data-[selected-single=true]:bg-primary data-[selected-single=true]:text-primary-foreground dark:hover:text-foreground [&>span]:text-xs [&>span]:opacity-70", style.Arg, style.Arg)
)

// config is what calendar.js reads (data-calendar).
type config struct {
	Mode              *string           `json:"mode"`
	Required          *bool             `json:"required,omitempty"`
	Selected          any               `json:"selected,omitempty"`
	DefaultMonth      *int              `json:"defaultMonth,omitempty"`
	Today             *int              `json:"today,omitempty"`
	NumberOfMonths    int               `json:"numberOfMonths,omitempty"`
	CaptionLayout     string            `json:"captionLayout,omitempty"`
	ShowOutsideDays   bool              `json:"showOutsideDays"`
	FixedWeeks        bool              `json:"fixedWeeks,omitempty"`
	ShowWeekNumber    bool              `json:"showWeekNumber,omitempty"`
	HideWeekdays      bool              `json:"hideWeekdays,omitempty"`
	DisableNavigation bool              `json:"disableNavigation,omitempty"`
	PagedNavigation   bool              `json:"pagedNavigation,omitempty"`
	ExcludeDisabled   bool              `json:"excludeDisabled,omitempty"`
	StartMonth        *int              `json:"startMonth,omitempty"`
	EndMonth          *int              `json:"endMonth,omitempty"`
	Disabled          []any             `json:"disabled,omitempty"`
	Hidden            []any             `json:"hidden,omitempty"`
	WeekStartsOn      *int              `json:"weekStartsOn,omitempty"`
	Min               int               `json:"min,omitempty"`
	Max               int               `json:"max,omitempty"`
	Calendar          string            `json:"calendar,omitempty"`
	Locale            string            `json:"locale,omitempty"`
	Numerals          string            `json:"numerals,omitempty"`
	Dir               string            `json:"dir,omitempty"`
	ClassNames        map[string]string `json:"classNames"`
	DayButtonClass    string            `json:"dayButtonClass"`
	ChevronClass      string            `json:"chevronClass"`

	// Server side only.
	selectedSingle *int
	selectedMulti  []int
	selectedRange  *Range
}

// Calendar returns the calendar's root <div data-slot="calendar">.
func Calendar(children ...any) *Builder {
	b := &Builder{buttonVariant: button.VariantGhost}
	b.cfg.ShowOutsideDays = true
	b.cfg.CaptionLayout = string(CaptionLabel)
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

func ptr[T any](v T) *T { return &v }

// AcceptChild consumes the typed props.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Mode:
		b.cfg.Mode = ptr(string(v))
	case CaptionLayout:
		b.cfg.CaptionLayout = string(v)
	case selectedDate:
		j := JDN(time.Time(v))
		b.cfg.selectedSingle = &j
	case selectedDates:
		b.cfg.selectedMulti = []int{}
		for _, t := range v {
			b.cfg.selectedMulti = append(b.cfg.selectedMulti, JDN(t))
		}
	case Range:
		b.cfg.selectedRange = &v
	case defaultMonth:
		b.cfg.DefaultMonth = ptr(JDN(time.Time(v)))
	case todayDate:
		b.cfg.Today = ptr(JDN(time.Time(v)))
	case startMonth:
		b.cfg.StartMonth = ptr(JDN(time.Time(v)))
	case endMonth:
		b.cfg.EndMonth = ptr(JDN(time.Time(v)))
	case NumberOfMonths:
		b.cfg.NumberOfMonths = int(v)
	case ShowOutsideDays:
		b.cfg.ShowOutsideDays = bool(v)
	case FixedWeeks:
		b.cfg.FixedWeeks = bool(v)
	case ShowWeekNumber:
		b.cfg.ShowWeekNumber = bool(v)
	case HideWeekdays:
		b.cfg.HideWeekdays = bool(v)
	case Required:
		b.cfg.Required = ptr(bool(v))
	case DisableNavigation:
		b.cfg.DisableNavigation = bool(v)
	case PagedNavigation:
		b.cfg.PagedNavigation = bool(v)
	case ExcludeDisabled:
		b.cfg.ExcludeDisabled = bool(v)
	case disabledMatch:
		for _, m := range v {
			b.cfg.Disabled = append(b.cfg.Disabled, m.matcher())
		}
	case hiddenMatch:
		for _, m := range v {
			b.cfg.Hidden = append(b.cfg.Hidden, m.matcher())
		}
	case WeekStartsOn:
		b.cfg.WeekStartsOn = ptr(int(v))
	case Min:
		b.cfg.Min = int(v)
	case Max:
		b.cfg.Max = int(v)
	case ButtonVariant:
		b.buttonVariant = button.Variant(v)
	case button.Variant:
		b.buttonVariant = v
	case ClassNames:
		b.classNames = v
	case persianFlag:
		b.persian = true
	case numeralsOption:
		b.numerals = Numerals(v)
	default:
		return false
	}
	return true
}

// shadcn's class names, merged with react-day-picker's defaults.
func (b *Builder) classNameMap() map[string]string {
	label := b.cfg.CaptionLayout == string(CaptionLabel)
	if b.st != "" {
		return b.styleClassNameMap(label)
	}
	nav := button.Variants(b.buttonVariant, "")
	captionLabel := "text-sm"
	if !label {
		captionLabel = "flex h-8 items-center gap-1 rounded-md pr-1 pl-2 text-sm [&>svg]:size-3.5 [&>svg]:text-muted-foreground"
	}
	day := "[&:first-child[data-selected=true]_button]:rounded-l-md"
	if b.cfg.ShowWeekNumber {
		day = "[&:nth-child(2)[data-selected=true]_button]:rounded-l-md"
	}
	m := map[string]string{
		"root":               lib.Cn("w-fit", "rdp-root"),
		"months":             lib.Cn("relative flex flex-col gap-4 md:flex-row", "rdp-months"),
		"month":              lib.Cn("flex w-full flex-col gap-4", "rdp-month"),
		"nav":                lib.Cn("absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1", "rdp-nav"),
		"button_previous":    lib.Cn(nav, "size-(--cell-size) p-0 select-none aria-disabled:opacity-50", "rdp-button_previous"),
		"button_next":        lib.Cn(nav, "size-(--cell-size) p-0 select-none aria-disabled:opacity-50", "rdp-button_next"),
		"month_caption":      lib.Cn("flex h-(--cell-size) w-full items-center justify-center px-(--cell-size)", "rdp-month_caption"),
		"dropdowns":          lib.Cn("flex h-(--cell-size) w-full items-center justify-center gap-1.5 text-sm font-medium", "rdp-dropdowns"),
		"dropdown_root":      lib.Cn("relative rounded-md border border-input shadow-xs has-focus:border-ring has-focus:ring-[3px] has-focus:ring-ring/50", "rdp-dropdown_root"),
		"dropdown":           lib.Cn("absolute inset-0 bg-popover opacity-0", "rdp-dropdown"),
		"caption_label":      lib.Cn("font-medium select-none", captionLabel, "rdp-caption_label"),
		"month_grid":         lib.Cn("w-full border-collapse", "rdp-month_grid"),
		"weekdays":           lib.Cn("flex", "rdp-weekdays"),
		"weekday":            lib.Cn("flex-1 rounded-md text-[0.8rem] font-normal text-muted-foreground select-none", "rdp-weekday"),
		"week":               lib.Cn("mt-2 flex w-full", "rdp-week"),
		"week_number_header": lib.Cn("w-(--cell-size) select-none", "rdp-week_number_header"),
		"week_number":        lib.Cn("text-[0.8rem] text-muted-foreground select-none", "rdp-week_number"),
		"day":                lib.Cn("group/day relative aspect-square h-full w-full p-0 text-center select-none [&:last-child[data-selected=true]_button]:rounded-r-md", day, "rdp-day"),
		"range_start":        lib.Cn("rounded-l-md bg-accent", "rdp-range_start"),
		"range_middle":       lib.Cn("rounded-none", "rdp-range_middle"),
		"range_end":          lib.Cn("rounded-r-md bg-accent", "rdp-range_end"),
		"today":              lib.Cn("rounded-md bg-accent text-accent-foreground data-[selected=true]:rounded-none", "rdp-today"),
		"outside":            lib.Cn("text-muted-foreground aria-selected:text-muted-foreground", "rdp-outside"),
		"disabled":           lib.Cn("text-muted-foreground opacity-50", "rdp-disabled"),
		"hidden":             lib.Cn("invisible", "rdp-hidden"),
		// react-day-picker defaults shadcn leaves alone.
		"chevron":         "rdp-chevron",
		"day_button":      "rdp-day_button",
		"footer":          "rdp-footer",
		"months_dropdown": "rdp-months_dropdown",
		"years_dropdown":  "rdp-years_dropdown",
		"weeks":           "rdp-weeks",
		"focused":         "rdp-focused",
		"selected":        "rdp-selected",
	}
	for k, v := range b.classNames {
		m[k] = v
	}
	return m
}

// styleClassNameMap is classNameMap inside a style.Scope.
func (b *Builder) styleClassNameMap(label bool) map[string]string {
	nav := button.StyleVariants(b.st, b.buttonVariant, "")
	captionLabel := "cn-calendar-caption text-sm"
	if !label {
		captionLabel = "cn-calendar-caption-label flex items-center gap-1 rounded-(--cell-radius) text-sm [&>svg]:size-3.5 [&>svg]:text-muted-foreground"
	}
	day := "[&:first-child[data-selected=true]_button]:rounded-l-(--cell-radius)"
	if b.cfg.ShowWeekNumber {
		day = "[&:nth-child(2)[data-selected=true]_button]:rounded-l-(--cell-radius)"
	}
	m := map[string]string{
		"root":               lib.Cn("w-fit", "rdp-root"),
		"months":             lib.Cn("relative flex flex-col gap-4 md:flex-row", "rdp-months"),
		"month":              lib.Cn("flex w-full flex-col gap-4", "rdp-month"),
		"nav":                lib.Cn("absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1", "rdp-nav"),
		"button_previous":    lib.Cn(nav, "size-(--cell-size) p-0 select-none aria-disabled:opacity-50", "rdp-button_previous"),
		"button_next":        lib.Cn(nav, "size-(--cell-size) p-0 select-none aria-disabled:opacity-50", "rdp-button_next"),
		"month_caption":      lib.Cn("flex h-(--cell-size) w-full items-center justify-center px-(--cell-size)", "rdp-month_caption"),
		"dropdowns":          lib.Cn("flex h-(--cell-size) w-full items-center justify-center gap-1.5 text-sm font-medium", "rdp-dropdowns"),
		"dropdown_root":      lib.Cn("cn-calendar-dropdown-root relative rounded-(--cell-radius)", "rdp-dropdown_root"),
		"dropdown":           lib.Cn("absolute inset-0 bg-popover opacity-0", "rdp-dropdown"),
		"caption_label":      lib.Cn("font-medium select-none", captionLabel, "rdp-caption_label"),
		"month_grid":         lib.Cn("w-full border-collapse", "rdp-month_grid"),
		"weekdays":           lib.Cn("flex", "rdp-weekdays"),
		"weekday":            lib.Cn("flex-1 rounded-(--cell-radius) text-[0.8rem] font-normal text-muted-foreground select-none", "rdp-weekday"),
		"week":               lib.Cn("mt-2 flex w-full", "rdp-week"),
		"week_number_header": lib.Cn("w-(--cell-size) select-none", "rdp-week_number_header"),
		"week_number":        lib.Cn("text-[0.8rem] text-muted-foreground select-none", "rdp-week_number"),
		"day":                lib.Cn("group/day relative aspect-square h-full w-full rounded-(--cell-radius) p-0 text-center select-none [&:last-child[data-selected=true]_button]:rounded-r-(--cell-radius)", day, "rdp-day"),
		"range_start":        lib.Cn("relative isolate z-0 rounded-l-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:right-0 after:w-4 after:bg-muted", "rdp-range_start"),
		"range_middle":       lib.Cn("rounded-none", "rdp-range_middle"),
		"range_end":          lib.Cn("relative isolate z-0 rounded-r-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:left-0 after:w-4 after:bg-muted", "rdp-range_end"),
		"today":              lib.Cn("rounded-(--cell-radius) bg-muted text-foreground data-[selected=true]:rounded-none", "rdp-today"),
		"outside":            lib.Cn("text-muted-foreground aria-selected:text-muted-foreground", "rdp-outside"),
		"disabled":           lib.Cn("text-muted-foreground opacity-50", "rdp-disabled"),
		"hidden":             lib.Cn("invisible", "rdp-hidden"),
		"chevron":            "rdp-chevron",
		"day_button":         "rdp-day_button",
		"footer":             "rdp-footer",
		"months_dropdown":    "rdp-months_dropdown",
		"years_dropdown":     "rdp-years_dropdown",
		"weeks":              "rdp-weeks",
		"focused":            "rdp-focused",
		"selected":           "rdp-selected",
	}
	for k, v := range b.classNames {
		m[k] = v
	}
	return m
}

const dayButtonClass = "flex aspect-square size-auto w-full min-w-(--cell-size) flex-col gap-1 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-md data-[range-end=true]:rounded-r-md data-[range-end=true]:bg-primary data-[range-end=true]:text-primary-foreground data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-accent data-[range-middle=true]:text-accent-foreground data-[range-start=true]:rounded-md data-[range-start=true]:rounded-l-md data-[range-start=true]:bg-primary data-[range-start=true]:text-primary-foreground data-[selected-single=true]:bg-primary data-[selected-single=true]:text-primary-foreground dark:hover:text-accent-foreground [&>span]:text-xs [&>span]:opacity-70"

const rootClass = "group/calendar bg-background p-3 [--cell-size:--spacing(8)] [[data-slot=card-content]_&]:bg-transparent [[data-slot=popover-content]_&]:bg-transparent"

// Element resolves the calendar to its root <div>.
func (b *Builder) Element() *h.Element {
	cfg := b.cfg
	cfg.ClassNames = b.classNameMap()
	cfg.DayButtonClass = button.Variants(button.VariantGhost, button.SizeIcon, lib.Cn(dayButtonClass, "rdp-day", cfg.ClassNames["day_button"]))
	if b.st != "" {
		cfg.DayButtonClass = button.StyleVariants(b.st, button.VariantGhost, button.SizeIcon, stDayButton.Class(b.st, "rdp-day", cfg.ClassNames["day_button"]))
	}
	cfg.ChevronClass = lib.Cn("size-4", cfg.ClassNames["chevron"])
	props := b.Elem()
	dir, hasDir := props.GetAttr("dir")
	if b.persian {
		cfg.Calendar = "persian"
		cfg.Locale = "fa"
		cfg.Numerals = "arabext"
		if !hasDir {
			dir, hasDir = "rtl", true
		}
	}
	if b.numerals != "" {
		cfg.Numerals = string(b.numerals)
	}
	if hasDir {
		cfg.Dir = dir
	}
	switch {
	case cfg.selectedSingle != nil:
		cfg.Selected = *cfg.selectedSingle
	case cfg.selectedMulti != nil:
		cfg.Selected = cfg.selectedMulti
	case cfg.selectedRange != nil:
		cfg.Selected = rangeJSON(*cfg.selectedRange)
	}
	raw, _ := json.Marshal(cfg)

	var weekStartsOn *int
	if cfg.WeekStartsOn != nil {
		weekStartsOn = cfg.WeekStartsOn
	}
	dp := newDayPicker(&cfg, newDateLib(b.persian, Numerals(cfg.Numerals), weekStartsOn))

	base := h.Div(h.Data("slot", "calendar"))
	if hasDir {
		base.Attr("dir", dir)
	}
	if cfg.Mode != nil {
		base.Data("mode", *cfg.Mode)
	}
	if cfg.Required != nil {
		base.Data("required", boolString(*cfg.Required))
	}
	if cfg.NumberOfMonths > 1 {
		base.Data("multiple-months", "true")
	}
	if cfg.ShowWeekNumber {
		base.Data("week-numbers", "true")
	}
	base.Data("calendar", string(raw)) // read and removed by calendar.js
	userClass := lib.Cn(rootClass,
		`rtl:**:[.rdp-button\_next>svg]:rotate-180`,
		`rtl:**:[.rdp-button\_previous>svg]:rotate-180`,
		props.ClassName())
	if b.st != "" {
		userClass = stRoot.Class(b.st, props.ClassName())
	}
	out := h.Spread(base, props.Clone().SetClass("").SetChildren(), nil)
	out.SetClass(lib.Cn(cfg.ClassNames["root"] + " " + userClass))
	out.SetChildren(dp.render()...)
	return out
}

// Render renders the calendar.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn resolves the calendar in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	b.st, _ = style.FromWriter(w)
	defer func() { b.st = "" }()
	return b.Element()
}

func boolString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// chevron is shadcn's Chevron component (lucide icons).
func chevron(cfg *config, orientation string, size int, disabled bool) *h.Element {
	var el *h.Element
	switch orientation {
	case "left":
		el = lucide.ChevronLeft(h.Class(cfg.ChevronClass))
	case "right":
		el = lucide.ChevronRight(h.Class(cfg.ChevronClass))
	default:
		el = lucide.ChevronDown(h.Class(cfg.ChevronClass))
	}
	if size != 24 {
		s := itoa(size)
		el.Attr("width", s).Attr("height", s)
	}
	if disabled {
		el.Bool("disabled", true)
	}
	return el
}
