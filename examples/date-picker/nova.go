package datepicker

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed nova.js
var novaScript string

// Docs previews of the date picker (ui/apps/v4/examples/radix/date-picker-*.tsx
// and data-picker-with-dropdowns.tsx), in radix-nova.
func init() {
	runtime.Module("example-nova-date-picker", novaScript)
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "date-picker", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("date-picker-demo", NovaDemo),
		nova("date-picker-basic", NovaBasic),
		nova("date-picker-dob", NovaDob),
		nova("date-picker-range", NovaRange),
		nova("data-picker-with-dropdowns", NovaWithDropdowns),
		nova("date-picker-time", NovaTime),
		nova("date-picker-input", NovaInput),
		nova("date-picker-natural-language", NovaNaturalLanguage),
		nova("date-picker-rtl", NovaRtl),
	)
}

// dpOptions configures nova.js for one picker.
type dpOptions struct {
	Format       string `json:"format"`
	Empty        string `json:"empty,omitempty"`
	DataEmpty    bool   `json:"dataEmpty,omitempty"`
	Range        bool   `json:"range,omitempty"`
	Close        bool   `json:"close,omitempty"`
	DefaultMonth bool   `json:"defaultMonth,omitempty"`
	Input        bool   `json:"input,omitempty"`
	Echo         bool   `json:"echo,omitempty"`
}

func (o dpOptions) attr() h.Prop {
	b, _ := json.Marshal(o)
	return h.Data("nova-date-picker", string(b))
}

func label(s string) h.Node { return h.Span(h.Data("dp-label", ""), s) }

func NovaDemo() h.Node { return novaDemo("Pick a date") }

func novaDemo(placeholder string, attrs ...any) h.Node {
	opts := dpOptions{Format: "PPP", Empty: placeholder, DataEmpty: true, DefaultMonth: true}
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(append([]any{button.VariantOutline, opts.attr(), h.Data("empty", "true"),
			h.Class("w-[212px] justify-between text-left font-normal data-[empty=true]:text-muted-foreground"),
			label(placeholder), lucide.ChevronDown()}, attrs...)...)),
		popover.Content(append([]any{h.Class("w-auto p-0"), popover.AlignStart, calendar.Calendar(calendar.ModeSingle)}, attrs...)...),
	)
}

func NovaRtl() h.Node {
	return style.ScopeRTL(novaDemo("اختر تاريخًا", h.Dir("rtl")))
}

func NovaBasic() h.Node {
	opts := dpOptions{Format: "PPP", Empty: "Pick a date", DefaultMonth: true}
	return field.Field(h.Class("mx-auto w-44"),
		field.Label(h.For("date-picker-simple"), "Date"),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, opts.attr(), h.ID("date-picker-simple"),
				h.Class("justify-start font-normal"), label("Pick a date"))),
			popover.Content(h.Class("w-auto p-0"), popover.AlignStart, calendar.Calendar(calendar.ModeSingle)),
		),
	)
}

func NovaDob() h.Node {
	opts := dpOptions{Format: "locale", Empty: "Select date", Close: true, DefaultMonth: true}
	return field.Field(h.Class("mx-auto w-44"),
		field.Label(h.For("date"), "Date of birth"),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, opts.attr(), h.ID("date"),
				h.Class("justify-start font-normal"), label("Select date"))),
			popover.Content(h.Class("w-auto overflow-hidden p-0"), popover.AlignStart,
				calendar.Calendar(calendar.ModeSingle, calendar.CaptionDropdown)),
		),
	)
}

func shortDate(t time.Time) string { return fmt.Sprintf("%s %02d, %d", t.Format("Jan"), t.Day(), t.Year()) }

func NovaRange() h.Node {
	from := time.Date(time.Now().Year(), 1, 20, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 20)
	opts := dpOptions{Format: "short", Empty: "Pick a date", Range: true, DefaultMonth: true}
	return field.Field(h.Class("mx-auto w-60"),
		field.Label(h.For("date-picker-range"), "Date Picker Range"),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, opts.attr(), h.ID("date-picker-range"),
				h.Class("justify-start px-2.5 font-normal"),
				lucide.Calendar(), rangeLabel(from, to))),
			popover.Content(h.Class("w-auto p-0"), popover.AlignStart,
				calendar.Calendar(calendar.ModeRange, calendar.DefaultMonth(from), calendar.SelectedRange(from, to), calendar.NumberOfMonths(2))),
		),
	)
}

func NovaWithDropdowns() h.Node {
	opts := dpOptions{Format: "PPP", Empty: "Pick a date"}
	return field.Field(h.Class("mx-auto w-72"),
		popover.Popover(
			field.Label(h.For("date-picker-with-dropdowns-desktop"), "Date"),
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, opts.attr(), h.ID("date-picker-with-dropdowns-desktop"),
				h.Class("justify-start px-2.5 font-normal"),
				label("Pick a date"), lucide.ChevronDown(h.Class("ml-auto")))),
			popover.Content(h.Class("w-auto p-0"), popover.AlignStart,
				calendar.Calendar(calendar.ModeSingle, calendar.CaptionDropdown),
				h.Div(h.Class("flex gap-2 border-t p-2"),
					button.Button(button.VariantOutline, button.SizeSm, h.Class("w-full"), h.Data("dp-done", ""), "Done")),
			),
		),
	)
}

func NovaTime() h.Node {
	opts := dpOptions{Format: "PPP", Empty: "Select date", Close: true, DefaultMonth: true}
	return field.Group(h.Class("mx-auto max-w-xs flex-row"),
		field.Field(
			field.Label(h.For("date-picker-optional"), "Date"),
			popover.Popover(
				popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, opts.attr(), h.ID("date-picker-optional"),
					h.Class("w-32 justify-between font-normal"), label("Select date"), lucide.ChevronDown())),
				popover.Content(h.Class("w-auto overflow-hidden p-0"), popover.AlignStart,
					calendar.Calendar(calendar.ModeSingle, calendar.CaptionDropdown)),
			),
		),
		field.Field(h.Class("w-32"),
			field.Label(h.For("time-picker-optional"), "Time"),
			input.Input(h.Type("time"), h.ID("time-picker-optional"), h.Attr("step", "1"), h.Value("10:30:00"),
				h.Class("appearance-none bg-background [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none")),
		),
	)
}

func longDate(t time.Time) string { return fmt.Sprintf("%s %02d, %d", t.Format("January"), t.Day(), t.Year()) }

func inputPicker(id, labelText, value, placeholder string, opts dpOptions, content []any, after ...any) h.Node {
	return field.Field(append([]any{opts.attr()}, append([]any{
		field.Label(h.For(id), labelText),
		ig.InputGroup(
			ig.Input(h.ID(id), h.Value(value), h.Placeholder(placeholder)),
			ig.Addon(ig.AlignInlineEnd,
				popover.Popover(
					popover.Trigger(lib.AsChild, ig.Button(h.ID("date-picker"), button.VariantGhost, ig.ButtonSizeIconXs, h.Aria("label", "Select date"),
						lucide.Calendar(), h.Span(h.Class("sr-only"), "Select date"))),
					popover.Content(content...),
				),
			),
		),
	}, after...)...)...)
}

func NovaInput() h.Node {
	date := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	return inputPicker("date-required", "Subscription Date", longDate(date), "June 01, 2025",
		dpOptions{Format: "long", Close: true, Input: true},
		[]any{h.Class("w-auto overflow-hidden p-0"), popover.AlignEnd, popover.AlignOffset(-8), popover.SideOffset(10),
			calendar.Calendar(calendar.ModeSingle, calendar.Selected(date), calendar.DefaultMonth(date))},
		h.Class("mx-auto w-48"))
}

func NovaNaturalLanguage() h.Node {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day()+2, 0, 0, 0, 0, time.Local)
	return inputPicker("date-optional", "Schedule Date", "In 2 days", "Tomorrow or next week",
		dpOptions{Format: "long", Close: true, Input: true, Echo: true, DefaultMonth: true},
		[]any{h.Class("w-auto overflow-hidden p-0"), popover.AlignEnd, popover.SideOffset(8),
			calendar.Calendar(calendar.ModeSingle, calendar.CaptionDropdown, calendar.Selected(date), calendar.DefaultMonth(date))},
		h.Class("mx-auto max-w-xs"),
		h.Div(h.Class("px-1 text-sm text-muted-foreground"),
			"Your post will be published on", " ",
			h.Span(h.Class("font-medium"), h.Data("dp-echo", ""), longDate(date)), "."),
	)
}
