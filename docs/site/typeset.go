package site

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/tooltip"
	h "shadcn.go/libhtml"
)

// The typeset page: app/(app)/(typeset) (/typeset, a designer for the
// typeset.css stylesheet) and its preview iframe, app/(view)/preview/typeset
// (typeset_preview.go). The page is one client component upstream; here the
// server renders it for the URL's params and typeset.js keeps it in sync
// (pickers, shuffle, locks, history, hover previews, the preview iframe).
//
// The data (fonts, fixtures, templates, icons) is extracted from the
// upstream sources by gen/typeset.mjs.

//go:embed data/typeset/typeset.json
var typesetJSON []byte

//go:embed static/typeset.js
var typesetScript string

func init() { runtime.Module("typeset", typesetScript) }

type typesetFont struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type typesetFontDefinition struct {
	Name       string   `json:"name"`
	Family     string   `json:"family"`
	Import     string   `json:"import"`
	Dependency string   `json:"dependency"`
	Subsets    []string `json:"subsets"`
	Weight     []string `json:"weight,omitempty"`
}

type typesetOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Width string `json:"width,omitempty"`
}

var typesetData = func() (d struct {
	Fonts       []typesetFont           `json:"fonts"`
	Definitions []typesetFontDefinition `json:"definitions"`
	Items       []struct {
		Value string `json:"value"`
		Label string `json:"label"`
		HTML  string `json:"html"`
	} `json:"items"`
	ChatQuestion string `json:"chatQuestion"`
	Templates    []struct {
		Value string `json:"value"`
		Title string `json:"title"`
	} `json:"templates"`
	FontClasses string            `json:"fontClasses"`
	Icons       map[string]string `json:"icons"`
}) {
	if err := json.Unmarshal(typesetJSON, &d); err != nil {
		panic(err)
	}
	return d
}()

// The option lists of lib/search-params.ts.
var (
	typesetSizes    = []typesetOption{{"14", "14px", ""}, {"15", "15px", ""}, {"16", "16px", ""}, {"18", "18px", ""}}
	typesetLeadings = []typesetOption{{"1.6", "Tight (1.6)", ""}, {"1.75", "Regular (1.75)", ""}, {"1.9", "Loose (1.9)", ""}}
	typesetFlows    = []typesetOption{{"1em", "Compact (1em)", ""}, {"1.25em", "Regular (1.25em)", ""}, {"2em", "Airy (2em)", ""}}
	typesetMeasures = []typesetOption{{"60", "60ch", "28em"}, {"70", "70ch", "33em"}, {"80", "80ch", "37em"}, {"90", "90ch", "42em"}}
)

// typesetParams are the page's search params (lib/search-params.ts), in
// their declaration order.
type typesetParams struct {
	Body, Heading, Mono, Scale, Measure, Flow, Leading, Item string
}

var typesetParamKeys = []string{"body", "heading", "mono", "scale", "measure", "flow", "leading", "item"}

var typesetDefaults = typesetParams{
	Body: "geist", Heading: "inherit", Mono: "geist-mono", Scale: "15",
	Measure: "80", Flow: "1.25em", Leading: "1.75", Item: "docs",
}

func (p *typesetParams) field(key string) *string {
	switch key {
	case "body":
		return &p.Body
	case "heading":
		return &p.Heading
	case "mono":
		return &p.Mono
	case "scale":
		return &p.Scale
	case "measure":
		return &p.Measure
	case "flow":
		return &p.Flow
	case "leading":
		return &p.Leading
	default:
		return &p.Item
	}
}

// typesetValues are the values each param accepts.
func typesetValues(key string) []string {
	var vs []string
	opts := func(list []typesetOption) {
		for _, o := range list {
			vs = append(vs, o.Value)
		}
	}
	switch key {
	case "heading":
		vs = append(vs, "inherit")
		fallthrough
	case "body", "mono":
		for _, f := range typesetData.Fonts {
			vs = append(vs, f.ID)
		}
	case "scale":
		opts(typesetSizes)
	case "measure":
		opts(typesetMeasures)
	case "flow":
		opts(typesetFlows)
	case "leading":
		opts(typesetLeadings)
	case "item":
		for _, it := range typesetData.Items {
			vs = append(vs, it.Value)
		}
	}
	return vs
}

// parseTypesetParams is nuqs' parseAsStringLiteral(...).withDefault(...):
// invalid values fall back to the default.
func parseTypesetParams(q url.Values) typesetParams {
	p := typesetDefaults
	for _, key := range typesetParamKeys {
		v := q.Get(key)
		for _, ok := range typesetValues(key) {
			if v == ok {
				*p.field(key) = v
			}
		}
	}
	return p
}

// serialize is serializeTypesetSearchParams: base plus the non-default
// params.
func (p typesetParams) serialize(base string) string {
	q := []string{}
	for _, key := range typesetParamKeys {
		if v := *p.field(key); v != *typesetDefaults.field(key) {
			q = append(q, key+"="+url.QueryEscape(v))
		}
	}
	if len(q) == 0 {
		return base
	}
	return base + "?" + strings.Join(q, "&")
}

func findTypesetFont(id string) *typesetFont {
	for i := range typesetData.Fonts {
		if typesetData.Fonts[i].ID == id {
			return &typesetData.Fonts[i]
		}
	}
	return nil
}

func typesetOptionLabel(list []typesetOption, v string) string {
	for _, o := range list {
		if o.Value == v {
			return o.Label
		}
	}
	return ""
}

func typesetMeasureWidth(v string) string {
	for _, o := range typesetMeasures {
		if o.Value == v {
			return o.Width
		}
	}
	return ""
}

// Typeset serves /typeset (typeset/page.tsx).
func (s *Site) Typeset(r *http.Request) (h.Node, error) {
	p := parseTypesetParams(r.URL.Query())
	return s.Document(r, Meta{Title: "Typeset", Description: "Typography for markdown you don't control."},
		s.AppLayout(r,
			h.Link(h.Rel("stylesheet"), h.Href(staticPrefix+"typeset-fonts.css")),
			h.Div(h.Data("typeset", ""), h.Data("typeset-params", typesetParamsJSON(p)), h.Data("typeset-data", typesetClientData),
				h.Class("relative z-10 flex min-h-0 flex-1 flex-col overflow-hidden section-soft [--customizer-width:--spacing(48)] [--gap:--spacing(4)] md:[--gap:--spacing(6)] 2xl:[--customizer-width:--spacing(56)] "+typesetData.FontClasses),
				h.Div(h.Data("slot", "designer"),
					h.Class("flex min-h-0 flex-1 flex-col items-start gap-(--gap) p-(--gap) pt-[calc(var(--gap)*0.25)] md:flex-row-reverse"),
					style.Scope(style.Nova,
						typesetDocsPanel(r, p),
						typesetPreview(p),
						typesetCustomizer(r, p),
					),
				),
			),
		),
	), nil
}

// typesetClientData is what typeset.js needs of the data: the defaults,
// the values each param accepts, the fonts, the option labels and the lock
// icons.
var typesetClientData = func() string {
	values := map[string][]string{}
	for _, k := range typesetParamKeys {
		values[k] = typesetValues(k)
	}
	labels := map[string]map[string]string{}
	for k, list := range map[string][]typesetOption{"scale": typesetSizes, "leading": typesetLeadings, "flow": typesetFlows, "measure": typesetMeasures} {
		labels[k] = map[string]string{}
		for _, o := range list {
			labels[k][o.Value] = o.Label
		}
	}
	defaults := map[string]string{}
	for _, k := range typesetParamKeys {
		defaults[k] = *typesetDefaults.field(k)
	}
	return jsonString(map[string]any{
		"defaults": defaults,
		"values":   values,
		"fonts":    typesetData.Fonts,
		"labels":   labels,
		"icons":    map[string]string{"lock": typesetData.Icons["SquareLock01Icon"], "unlock": typesetData.Icons["SquareUnlock01Icon"]},
	})
}()

func typesetParamsJSON(p typesetParams) string {
	m := map[string]string{}
	for _, key := range typesetParamKeys {
		m[key] = *p.field(key)
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// --- icons -----------------------------------------------------------------

// typesetHugeicon renders a hugeicon as @hugeicons/react's HugeiconsIcon does
// with strokeWidth={2}.
func typesetHugeicon(inner, class string) h.Node {
	return h.El("svg", h.Attr("xmlns", "http://www.w3.org/2000/svg"), h.Attr("width", "24"), h.Attr("height", "24"),
		h.Attr("viewBox", "0 0 24 24"), h.Attr("fill", "none"), h.Attr("color", "currentColor"),
		h.Attr("stroke-width", "2"), h.Attr("stroke", "currentColor"), h.Class(class), h.UnsafeHTML(inner))
}

func typesetIcon(name, class string) h.Node { return typesetHugeicon(typesetData.Icons[name], class) }

// lineHeightIcon is customizer.tsx's LineHeightIcon.
const lineHeightIcon = `<path d="M4.5 3.5H19.5" stroke="currentColor" stroke-linecap="round" stroke-width="2"></path><path d="M4.5 20.5H19.5" stroke="currentColor" stroke-linecap="round" stroke-width="2"></path><path d="M17 17L14.8905 11.4741C13.9109 8.90801 13.4211 7.625 12.625 7.625C11.8289 7.625 11.3391 8.90801 10.3595 11.4741L8.25 17" stroke="currentColor" stroke-linecap="round" stroke-width="2"></path><path d="M9.5 13H15.75" stroke="currentColor" stroke-linecap="round" stroke-width="2"></path>`

// --- preview (components/preview.tsx, toolbar.tsx) ---------------------------

func typesetPreview(p typesetParams) h.Node {
	var pills []any
	for i, it := range typesetData.Items {
		pills = append(pills, tooltip.Tooltip(
			tooltip.Trigger(lib.AsChild,
				h.Button(h.Type("button"), h.Data("slot", "button"), h.Data("variant", "ghost"), h.Data("size", "sm"),
					h.Data("active", boolString(p.Item == it.Value)), h.Data("typeset-item", it.Value),
					h.Class("inline-flex shrink-0 items-center justify-center gap-1.5 whitespace-nowrap outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-accent dark:hover:bg-accent/50 has-[>svg]:px-2.5 h-7 min-w-7 cursor-pointer rounded-lg px-2 text-xs font-medium text-muted-foreground transition-colors hover:text-foreground data-[active=true]:bg-accent data-[active=true]:text-accent-foreground"),
					twoDigits(i+1),
				),
			),
			tooltip.Content(tooltip.SideTop, tooltip.SideOffset(10), it.Label),
		))
	}
	return h.Div(h.Class("relative isolate flex size-full min-h-0 flex-1 overflow-hidden rounded-2xl bg-background ring-1 ring-foreground/10"),
		h.Iframe(h.Data("typeset-preview", ""), h.Src(p.serialize("/preview/typeset/"+p.Item)), h.Title("typeset preview"), h.Class("size-full")),
		h.Div(h.Class("absolute bottom-3 left-1/2 z-20 flex -translate-x-1/2 items-center justify-center gap-1.5"),
			h.Div(h.Class("dark flex items-center gap-1 rounded-xl bg-card/90 p-1 shadow-xl backdrop-blur-xl"), h.Group(pills...)),
			h.Div(h.Class("dark flex items-center gap-1 rounded-xl bg-card/90 p-1 shadow-xl backdrop-blur-xl"),
				h.A(h.Data("slot", "button"), h.Data("variant", "ghost"), h.Data("size", "sm"), h.Data("typeset-open", ""),
					h.Href(p.serialize("/preview/typeset/"+p.Item)), h.Target("_blank"), h.Rel("noreferrer"),
					h.Class("inline-flex shrink-0 items-center justify-center gap-1.5 whitespace-nowrap outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-accent dark:hover:bg-accent/50 has-[>svg]:px-2.5 h-7 cursor-pointer rounded-lg px-2.5 text-xs font-medium text-muted-foreground transition-colors hover:text-foreground"),
					typesetIcon("LinkSquare02Icon", "size-3.5 md:hidden"),
					h.Span(h.Class("max-md:sr-only"), "Open in New Tab"),
				),
			),
		),
	)
}

func twoDigits(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

// --- customizer (components/customizer.tsx and its pickers) ------------------

const pickerTriggerClass = "relative w-36 shrink-0 touch-manipulation rounded-xl p-3 ring-1 ring-foreground/10 select-none hover:bg-muted focus-visible:ring-foreground/50 focus-visible:outline-none disabled:opacity-50 data-popup-open:bg-muted data-[state=open]:bg-muted md:w-full md:rounded-lg md:px-2.5 md:py-2"

// pickerContent is picker.tsx's PickerContent: a Base UI menu popup, here
// the Radix DropdownMenu content (its CSS variables renamed).
func pickerContent(class string, children ...any) h.Node {
	args := []any{dropdownmenu.SideRight, dropdownmenu.AlignStart, dropdownmenu.SideOffset(20),
		h.Class(lib.Cn("cn-menu-target z-50 no-scrollbar max-h-(--radix-dropdown-menu-content-available-height) w-[calc(var(--radix-dropdown-menu-content-available-width)-(--spacing(6)))] min-w-32 origin-(--radix-dropdown-menu-content-transform-origin) translate-y-2 overflow-x-hidden overflow-y-auto rounded-xl border-0 bg-neutral-950/80 p-1.5 text-neutral-100 shadow-none ring-1 ring-neutral-950/80 backdrop-blur-xl outline-none md:w-52 dark:bg-neutral-800/90 dark:ring-neutral-700/50 data-closed:overflow-hidden", class))}
	return dropdownmenu.Content(append(args, children...)...)
}

const pickerItemClass = "relative flex cursor-default items-center gap-2 rounded-lg px-2 py-1.5 text-sm font-medium outline-hidden select-none **:text-neutral-100 focus:bg-neutral-600 focus:text-neutral-100 focus:**:text-neutral-100 data-inset:pl-8 dark:focus:bg-neutral-700/80 pointer-coarse:gap-3 pointer-coarse:py-2.5 pointer-coarse:pl-3 pointer-coarse:text-base data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

const pickerRadioItemClass = "relative flex cursor-default items-center gap-2 rounded-lg py-1.5 pr-8 pl-2 text-sm font-medium outline-hidden select-none **:text-neutral-100 focus:bg-neutral-600 focus:text-neutral-100 focus:**:text-neutral-100 data-inset:pl-8 dark:focus:bg-neutral-700/80 pointer-coarse:gap-3 pointer-coarse:py-2.5 pointer-coarse:pl-3 pointer-coarse:text-base data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

func pickerRadioItem(value string, label any) h.Node {
	return dropdownmenu.RadioItem(dropdownmenu.Value(value), h.Data("value", value), h.Class(pickerRadioItemClass), label)
}

func pickerLabel(text string) h.Node {
	return dropdownmenu.Label(h.Class("px-2 py-1.5 text-xs font-medium text-neutral-400 data-inset:pl-8"), text)
}

func pickerSeparator() h.Node {
	return dropdownmenu.Separator(h.Class("-mx-1.5 my-1.5 h-px bg-neutral-600 dark:bg-neutral-700"))
}

func pickerShortcut(text string, extra ...any) h.Node {
	return dropdownmenu.Shortcut(h.Class("ml-auto text-xs tracking-widest text-neutral-400! group-focus/dropdown-menu-item:text-neutral-100"), text, h.Group(extra...))
}

// lockButton is lock-button.tsx.
func lockButton(param string) h.Node {
	return h.Button(h.Type("button"), h.TitleAttr("Lock"), h.Aria("label", "Lock"), h.Data("locked", "false"), h.Data("typeset-lock", param),
		h.Class("flex size-4 cursor-pointer items-center justify-center rounded opacity-0 ring-foreground/60 transition-opacity outline-none group-focus-within/picker:opacity-100 group-hover/picker:opacity-100 focus:opacity-100 focus-visible:ring-1 data-[locked=true]:opacity-100 max-md:hidden pointer-coarse:hidden absolute top-1/2 right-8 -translate-y-1/2"),
		typesetIcon("SquareUnlock01Icon", "size-5 text-foreground"),
	)
}

// optionPicker is option-picker.tsx.
func optionPicker(label, param, icon string, options []typesetOption, value, class string) h.Node {
	var items []any
	for _, o := range options {
		items = append(items, pickerRadioItem(o.Value, o.Label))
	}
	return h.Div(h.Class(lib.Cn("group/picker relative", class)),
		dropdownmenu.DropdownMenu(
			dropdownmenu.Trigger(lib.AsChild,
				h.Button(h.Type("button"), h.Class(pickerTriggerClass),
					h.Div(h.Class("flex min-w-0 flex-col justify-start pr-8 text-left"),
						h.Div(h.Class("text-xs text-muted-foreground"), label),
						h.Div(h.Class("line-clamp-1 text-sm font-medium text-foreground"), h.Data("typeset-label", param), typesetOptionLabel(options, value)),
					),
					h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 items-center justify-center text-foreground select-none md:right-2.5"),
						typesetHugeicon(icon, "size-4.5"),
					),
				),
			),
			pickerContent("", h.Data("typeset-picker", param),
				dropdownmenu.RadioGroup(dropdownmenu.Value(value), dropdownmenu.Controlled, h.Data("typeset-group", param), h.Group(items...)),
			),
		),
		lockButton(param),
	)
}

// fontPicker is font-picker.tsx.
func fontPicker(label, param string, p typesetParams) h.Node {
	value := *p.field(param)
	body := findTypesetFont(p.Body)
	if body == nil {
		body = &typesetData.Fonts[0]
	}
	current := body
	if !(param == "heading" && value == "inherit") {
		if f := findTypesetFont(value); f != nil {
			current = f
		}
	}
	var groups []any
	if param == "heading" {
		groups = append(groups,
			dropdownmenu.Group(pickerRadioItem("inherit", h.Span(h.Data("typeset-inherit-label", ""), body.Label))),
			pickerSeparator(),
		)
	}
	for _, g := range []struct{ label, typ string }{{"Sans", "sans"}, {"Serif", "serif"}, {"Mono", "mono"}} {
		items := []any{pickerLabel(g.label)}
		for _, f := range typesetData.Fonts {
			if f.Type == g.typ {
				items = append(items, pickerRadioItem(f.ID, f.Label))
			}
		}
		groups = append(groups, dropdownmenu.Group(items...))
	}
	return h.Div(h.Class("group/picker relative"),
		dropdownmenu.DropdownMenu(
			dropdownmenu.Trigger(lib.AsChild,
				h.Button(h.Type("button"), h.Class(pickerTriggerClass),
					h.Div(h.Class("flex flex-col justify-start text-left"),
						h.Div(h.Class("text-xs text-muted-foreground"), label),
						h.Div(h.Class("line-clamp-1 max-w-[80%] truncate text-sm font-medium text-foreground"), h.Data("typeset-label", param), current.Label),
					),
					h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 items-center justify-center text-base text-foreground select-none md:right-2.5"),
						h.Data("typeset-aa", param), h.UnsafeAttr("style", "font-family: "+current.Value), "Aa"),
				),
			),
			pickerContent("max-h-96", h.Data("typeset-picker", param),
				dropdownmenu.RadioGroup(dropdownmenu.Value(value), dropdownmenu.Controlled, h.Data("typeset-group", param), h.Group(groups...)),
			),
		),
		lockButton(param),
	)
}

// typesetMainMenu is main-menu.tsx.
func typesetMainMenu() h.Node {
	item := func(action string, children ...any) h.Node {
		return dropdownmenu.Item(h.Data("typeset-action", action), h.Class(pickerItemClass), h.Group(children...))
	}
	return dropdownmenu.DropdownMenu(
		dropdownmenu.Trigger(lib.AsChild,
			h.Button(h.Type("button"), h.Class(lib.Cn(pickerTriggerClass, "flex items-center justify-between gap-2 rounded-lg px-1.75 ring-1 ring-foreground/10 focus-visible:ring-1")),
				h.Span(h.Class("font-medium"), "Menu"),
				typesetIcon("Menu09Icon", "size-5"),
			),
		),
		pickerContent("", dropdownmenu.AlignOffset(-8), h.Data("typeset-menu", ""),
			dropdownmenu.Group(
				item("shuffle", "Shuffle ", pickerShortcut("R")),
				item("toggle-theme", "Light/Dark ", pickerShortcut("D")),
			),
			pickerSeparator(),
			dropdownmenu.Group(
				dropdownmenu.Item(dropdownmenu.Disabled(true), h.Data("typeset-action", "undo"), h.Class(pickerItemClass), "Undo ", pickerShortcut("⌘Z", h.Data("typeset-shortcut", "undo"))),
				dropdownmenu.Item(dropdownmenu.Disabled(true), h.Data("typeset-action", "redo"), h.Class(pickerItemClass), "Redo ", pickerShortcut("⇧⌘Z", h.Data("typeset-shortcut", "redo"))),
				pickerSeparator(),
				item("reset", "Reset ", pickerShortcut("⇧R")),
			),
		),
	)
}

const typesetFooterButtonClass = "group/button inline-flex shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 h-8 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 touch-manipulation bg-transparent! px-2! py-0! text-sm! transition-none hover:bg-muted! pointer-coarse:h-10! min-w-0 flex-1 md:flex-none"

func typesetCustomizer(r *http.Request, p typesetParams) h.Node {
	separator := func() h.Node {
		return h.Div(h.Data("slot", "field-separator"), h.Data("content", "false"),
			h.Class("relative -my-2 h-5 text-sm group-data-[variant=outline]/field-group:-mb-2 hidden md:block"),
			h.Div(h.Data("slot", "separator"), h.Data("orientation", "horizontal"), h.Role("separator"), h.Aria("orientation", "horizontal"),
				h.Class("shrink-0 bg-border data-horizontal:h-px data-horizontal:w-full data-vertical:w-px data-vertical:self-stretch absolute inset-0 top-1/2")),
		)
	}
	return h.Div(h.Data("slot", "card"), h.Data("size", "sm"), h.Data("typeset-customizer", ""),
		h.Class("group/card flex flex-col gap-(--card-spacing) overflow-hidden rounded-2xl py-(--card-spacing) text-sm text-card-foreground ring-1 ring-foreground/10 [--card-spacing:--spacing(4)] has-data-[slot=card-footer]:pb-0 has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] data-[size=sm]:has-data-[slot=card-footer]:pb-0 *:[img:first-child]:rounded-t-xl *:[img:last-child]:rounded-b-xl dark isolate z-10 max-h-full min-h-0 w-full self-start bg-card/90 backdrop-blur-xl md:w-(--customizer-width)"),
		h.Div(h.Data("slot", "card-header"),
			h.Class("group/card-header @container/card-header auto-rows-min rounded-t-xl px-(--card-spacing) has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto] [.border-b]:pb-(--card-spacing) hidden items-center justify-between gap-2 border-b md:flex"),
			typesetMainMenu(),
		),
		h.Div(h.Data("slot", "card-content"),
			h.Class("px-(--card-spacing) no-scrollbar min-h-0 flex-1 overflow-x-auto overflow-y-hidden max-md:px-0 md:overflow-y-auto"),
			h.Div(h.Data("slot", "field-group"),
				h.Class("group/field-group @container/field-group flex w-full data-[slot=checkbox-group]:gap-3 *:data-[slot=field-group]:gap-4 flex-row gap-2.5 py-px **:data-[slot=field-separator]:-mx-4 **:data-[slot=field-separator]:w-auto max-md:px-3 md:flex-col md:gap-3.25"),
				optionPicker("Measure", "measure", typesetData.Icons["ArrowHorizontalIcon"], typesetMeasures, p.Measure, "max-[28rem]:hidden"),
				separator(),
				fontPicker("Heading", "heading", p),
				fontPicker("Body", "body", p),
				fontPicker("Mono", "mono", p),
				separator(),
				optionPicker("Size", "scale", typesetData.Icons["TextSmallcapsIcon"], typesetSizes, p.Scale, ""),
				optionPicker("Leading", "leading", lineHeightIcon, typesetLeadings, p.Leading, ""),
				optionPicker("Flow", "flow", typesetData.Icons["ParagraphSpacingIcon"], typesetFlows, p.Flow, ""),
				h.Div(h.Aria("hidden", "true"), h.Class("w-0.5 shrink-0 md:hidden")),
			),
		),
		h.Div(h.Data("slot", "card-footer"),
			h.Class("flex items-center rounded-b-xl border-t bg-muted/50 p-(--card-spacing) min-w-0 flex-row-reverse gap-2 md:flex-col md:**:[button]:w-full"),
			h.Button(h.Type("button"), h.Data("slot", "button"), h.Data("typeset-action", "shuffle"), h.Class(typesetFooterButtonClass),
				h.Span(h.Class("w-full truncate text-center font-medium"), "Shuffle"),
			),
			typesetGetCodeDrawer(r, p),
		),
	)
}
