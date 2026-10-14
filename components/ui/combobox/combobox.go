// Package combobox is the shadcn/ui Combobox: an input that filters a list
// of items in a popup (Base UI's Combobox, @base-ui/react 1.6.0, under
// shadcn's combobox.tsx). It selects one item, or several shown as chips.
//
//	combobox.Combobox(
//		combobox.Input(h.Placeholder("Select a framework")),
//		combobox.Content(
//			combobox.Empty("No items found."),
//			combobox.List(
//				combobox.Item(combobox.ItemValue("next.js"), combobox.ItemLabel("Next.js"), "Next.js"),
//				combobox.Item(combobox.ItemValue("remix"), combobox.ItemLabel("Remix"), "Remix"),
//			),
//		),
//	)
//
// Items are rendered on the server (inside the content, which stays in an
// inert template until the popup opens); the behavior (combobox.js)
// filters them by their ItemLabel as the user types, exactly like Base UI
// filters the `items` prop. An item's ItemValue is what the form receives;
// its ItemLabel (default: the value) is what filtering, the input and
// Value show.
//
// Root options (children of Combobox): Multiple, AutoHighlight,
// DefaultValue (an ItemValue; repeat it with Multiple), Name, Disabled.
//
// Scripts observe the combobox through events dispatched on the hidden
// input: "shadcn:value-change" (detail.value: the selected ItemValue, or
// the list of them with Multiple), "shadcn:open-change" (detail.open),
// and "shadcn:combobox-render" on the list each time it is rendered
// (detail.groups: the rendered groups in order).
package combobox

import (
	_ "embed"
	"encoding/json"
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed combobox.js
var script string

func init() { runtime.Module("combobox", script) }

// Root options.
type (
	// MultipleFlag is the type of Multiple.
	MultipleFlag struct{}
	// AutoHighlightFlag is the type of AutoHighlight.
	AutoHighlightFlag struct{}
	// DisabledFlag is the type of Disabled.
	DisabledFlag struct{}
	// DefaultValue selects the item with this ItemValue initially.
	DefaultValue string
	// Name is the form field name of the selected value(s).
	Name string
)

// Multiple lets the combobox select several items (shown as chips).
var Multiple = MultipleFlag{}

// AutoHighlight highlights the first matching item as the user types.
var AutoHighlight = AutoHighlightFlag{}

// Disabled disables the combobox (on Combobox) or one part (on Input,
// Trigger, Clear).
var Disabled = DisabledFlag{}

// Item options.
type (
	// ItemValue is an item's value, submitted with the form and matched by
	// DefaultValue.
	ItemValue string
	// ItemLabel is an item's label: what filtering matches, what the input
	// shows once the item is selected (default: the ItemValue).
	ItemLabel string
)

// Input options.
type (
	// ShowTrigger shows the open button inside the Input (default true).
	ShowTrigger bool
	// ShowClear shows the clear button inside the Input (default false).
	ShowClear bool
)

// Content options.
type (
	// Side is the side of the anchor the content is placed on.
	Side string
	// Align is the content's alignment against the anchor.
	Align string
	// SideOffset is the distance (px) from the anchor (default 6).
	SideOffset float64
	// AlignOffset shifts the content along the anchor (px, default 0).
	AlignOffset float64
)

// Sides and alignments.
const (
	SideTop     Side  = "top"
	SideBottom  Side  = "bottom"
	SideLeft    Side  = "left"
	SideRight   Side  = "right"
	AlignStart  Align = "start"
	AlignCenter Align = "center"
	AlignEnd    Align = "end"
)

// Anchor is the element Content is positioned against, created by
// UseAnchor (upstream useComboboxAnchor). Pass it to Chips (or any part
// that accepts it) and to Content.
type Anchor struct{ name string }

// UseAnchor returns a new Anchor.
func UseAnchor() *Anchor { return &Anchor{name: uid.New("anchor")} }

// ShowRemove shows a Chip's remove button (default true).
type ShowRemove bool

// Placeholder is the text Value shows while nothing is selected.
type Placeholder string

// --- shared state --------------------------------------------------------------------------

type item struct {
	Value string `json:"v"`
	Label string `json:"l"`
}

type state struct {
	id        string
	multiple  bool
	autoHigh  bool
	disabled  bool
	name      string
	hasName   bool
	defaults  []string
	hasDef    bool
	items     []item
	inputIn   bool // an Input renders inside Content
	collected bool
}

func (s *state) selected() []item {
	var out []item
	for _, d := range s.defaults {
		for _, it := range s.items {
			if it.Value == d {
				out = append(out, it)
				break
			}
		}
		if !s.multiple && len(out) > 0 {
			break
		}
	}
	return out
}

// hasSelection reports whether a value is selected (single: an item, even
// one with an empty value; multiple: at least one).
func (s *state) hasSelection() bool { return len(s.selected()) > 0 }

func (s *state) selectedLabel() string {
	if sel := s.selected(); len(sel) > 0 && !s.multiple {
		return sel[0].Label
	}
	return ""
}

type part struct{ st *state }

func (p *part) setState(s *state) { p.st = s }

func (p *part) state() *state {
	if p.st == nil {
		p.st = &state{id: uid.New("base-ui")}
	}
	return p.st
}

type participant interface{ setState(*state) }

func mark(el *h.Element, s *state, name string) {
	el.Attr("data-combobox", s.id)
	el.Attr("data-combobox-part", name)
}

// --- Combobox (root) -----------------------------------------------------------------------

// Builder is returned by Combobox.
type Builder struct {
	children []any
	st       *state
}

// Combobox is the root: it renders no element of its own, only its parts
// followed by the hidden form input(s).
func Combobox(children ...any) *Builder {
	b := &Builder{st: &state{id: uid.New("base-ui")}}
	for _, c := range children {
		switch v := c.(type) {
		case MultipleFlag:
			b.st.multiple = true
		case AutoHighlightFlag:
			b.st.autoHigh = true
		case DisabledFlag:
			b.st.disabled = true
		case DefaultValue:
			b.st.defaults = append(b.st.defaults, string(v))
			b.st.hasDef = true
		case Name:
			b.st.name = string(v)
			b.st.hasName = true
		default:
			b.children = append(b.children, c)
		}
	}
	return b
}

func (b *Builder) collect() {
	s := b.st
	if s.collected {
		return
	}
	s.collected = true
	var visit func(n any) bool
	visit = func(n any) bool {
		switch c := n.(type) {
		case *Builder:
			return false
		case *ItemBuilder:
			s.items = append(s.items, c.item())
		case *ContentBuilder:
			tree.Walk(c.Elem().ChildNodes(), func(m any) bool {
				if _, ok := m.(*InputBuilder); ok {
					s.inputIn = true
				}
				return true
			})
		case *ValueBuilder:
			c.setState(s)
			tree.Walk(c.children, visit)
			return false
		}
		if p, ok := n.(participant); ok {
			p.setState(s)
		}
		return true
	}
	tree.Walk(b.children, visit)
}

const (
	visuallyHidden      = "clip-path: inset(50%); overflow: hidden; white-space: nowrap; border: 0px; padding: 0px; width: 1px; height: 1px; margin: -1px; position: fixed; top: 0px; left: 0px;"
	visuallyHiddenInput = "clip-path: inset(50%); overflow: hidden; white-space: nowrap; border: 0px; padding: 0px; width: 1px; height: 1px; margin: -1px; position: absolute;"
)

// Render renders the parts and the hidden input.
func (b *Builder) Render(w io.Writer) error {
	b.collect()
	s := b.st
	sel := s.selected()
	value := ""
	if !s.multiple && len(sel) > 0 {
		value = sel[0].Value
	}
	hidden := h.Input(h.TabIndex(-1), h.Aria("hidden", "true"), h.Value(value))
	if s.hasName && !s.multiple {
		hidden.Attr("name", s.name).Attr("style", visuallyHiddenInput)
	} else {
		hidden.ID(s.id+"-hidden-input").Attr("style", visuallyHidden)
	}
	if s.disabled {
		hidden.Bool("disabled", true)
	}
	opts := map[string]any{
		"id":            s.id,
		"multiple":      s.multiple,
		"autoHighlight": s.autoHigh,
		"disabled":      s.disabled,
		"items":         s.items,
		"name":          s.name,
		"inputInside":   s.inputIn,
	}
	if s.hasDef {
		opts["defaults"] = s.defaults
	}
	js, _ := json.Marshal(opts)
	hidden.Attr("data-combobox-options", string(js))
	mark(hidden, s, "hidden")
	out := []any{b.children, hidden}
	if s.multiple && s.hasName {
		for _, it := range sel {
			in := h.Input(h.Type("hidden"), h.Name(s.name), h.Value(it.Value))
			if s.disabled {
				in.Bool("disabled", true)
			}
			out = append(out, in)
		}
	}
	return h.Group(out...).Render(w)
}

// --- Input ---------------------------------------------------------------------------------

// InputBuilder is returned by Input.
type InputBuilder struct {
	h.Builder[InputBuilder]
	part
	showTrigger bool
	showClear   bool
	disabled    bool
}

// Input is the text input inside an InputGroup, with the open button
// (ShowTrigger, default true) and the clear button (ShowClear). Its
// attributes (placeholder, aria-invalid, ...) go to the <input>, its class
// to the InputGroup; other children (InputGroup addons) follow the input.
func Input(children ...any) *InputBuilder {
	b := &InputBuilder{showTrigger: true}
	b.Builder = h.NewBuilder(b, "input")
	return b.Children(children...)
}

// AcceptChild consumes ShowTrigger, ShowClear and Disabled.
func (b *InputBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case ShowTrigger:
		b.showTrigger = bool(v)
	case ShowClear:
		b.showClear = bool(v)
	case DisabledFlag:
		b.disabled = true
	default:
		return false
	}
	return true
}

// Element resolves the input group.
func (b *InputBuilder) Element() *h.Element {
	s := b.state()
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	kids := props.ChildNodes()
	props.SetChildren()
	disabled := b.disabled || s.disabled

	inputID := s.id
	if s.inputIn {
		inputID = uid.New("base-ui")
	}
	value := ""
	if !s.inputIn && !s.multiple {
		value = s.selectedLabel()
	}
	in := inputgroup.Input(h.Prop(func(e *h.Element) { e.Bool("disabled", disabled) })).Element()
	in = h.Spread(in, comboboxInputAttrs(inputID, value), nil)
	in = h.Spread(in, props, nil)
	mark(in, s, "input")
	if b.disabled {
		in.Attr("data-combobox-disabled", "")
	}

	var addon []any
	if b.showTrigger {
		addon = append(addon, inputTrigger(s, disabled))
	}
	if b.showClear {
		addon = append(addon, clearButton(s, disabled, b.disabled))
	}
	group := inputgroup.InputGroup(h.Class(lib.Cn("w-auto", className)),
		in,
		inputgroup.Addon(append([]any{inputgroup.AlignInlineEnd}, addon...)...),
		kids,
	)
	return group.Element()
}

// Render renders the input group.
func (b *InputBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

func comboboxInputAttrs(id, value string) *h.Element {
	return h.Input(
		h.Attr("autocomplete", "off"),
		h.Attr("spellcheck", "false"),
		h.Attr("autocorrect", "off"),
		h.Attr("autocapitalize", "none"),
		h.Role("combobox"),
		h.Aria("expanded", "false"),
		h.Aria("haspopup", "listbox"),
		h.Aria("autocomplete", "list"),
		h.ID(id),
		h.Type("text"),
		h.Value(value),
	)
}

const triggerClass = "[&_svg:not([class*='size-'])]:size-4"

func triggerIcon() h.Node {
	return lucide.ChevronDown(h.Data("slot", "combobox-trigger-icon"), h.Class("pointer-events-none size-4 text-muted-foreground"))
}

// inputTrigger is <InputGroupButton size="icon-xs" variant="ghost" asChild
// ...><ComboboxTrigger /></InputGroupButton>.
func inputTrigger(s *state, disabled bool) *h.Element {
	slot := inputgroup.Button(inputgroup.ButtonSizeIconXs, button.VariantGhost,
		h.Class("group-has-data-[slot=combobox-clear]/input-group:hidden data-pressed:bg-transparent")).Element()
	el := h.Button(
		h.Type("button"),
		h.TabIndex(-1),
		h.Aria("expanded", "false"),
		h.Aria("haspopup", "listbox"),
		h.Data("slot", "input-group-button"),
		h.Data("variant", "ghost"),
		h.Data("size", "icon-xs"),
		h.Class(lib.Cn(triggerClass, slot.ClassName())),
		triggerIcon(),
	)
	if !s.multiple && !s.hasSelection() {
		el.Attr("data-placeholder", "")
	}
	if disabled {
		el.Attr("data-disabled", "").Bool("disabled", true)
	}
	mark(el, s, "trigger")
	return el
}

func clearButton(s *state, disabled, local bool) h.Node {
	el := inputgroup.Button(inputgroup.ButtonSizeIconXs, button.VariantGhost,
		h.Data("slot", "combobox-clear"),
		h.TabIndex(-1),
		lucide.X(h.Class("pointer-events-none")),
	).Element()
	if disabled {
		el.Attr("data-disabled", "").Bool("disabled", true)
	}
	if local {
		el.Attr("data-combobox-disabled", "")
	}
	tpl := h.Template(h.Data("parity-ignore", ""), el)
	mark(tpl, s, "clear")
	if !s.hasSelection() {
		return tpl
	}
	live := el.Clone()
	live.Attr("data-visible", "")
	mark(live, s, "clear-live")
	return h.Group(live, tpl)
}

// --- Trigger, Value, Clear -----------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	part
	disabled bool
}

// Trigger is a button that opens the popup, followed by a chevron (upstream
// ComboboxTrigger). Render it as another component (a button.Button) with
// lib.AsChild; an Input inside Content then filters the list.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Disabled.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	if _, ok := c.(DisabledFlag); ok {
		b.disabled = true
		return true
	}
	return false
}

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	s := b.state()
	props := b.Elem().Clone()
	kids := props.ChildNodes()
	var render *h.Element
	var rest []any
	for _, k := range kids {
		if e, ok := k.(h.Elementer); ok && render == nil && isRender(k) {
			render = e.Element()
			continue
		}
		rest = append(rest, k)
	}
	props.SetChildren(append(rest, triggerIcon())...)
	className := props.ClassName()
	props.SetClass("")
	var el *h.Element
	if render != nil {
		el = render.Clone()
		el.SetChildren()
		el.SetClass(lib.Cn(triggerClass, lib.Cn(render.ClassName(), className)))
	} else {
		el = h.Button(h.Class(lib.Cn(triggerClass, className)))
	}
	el = h.Spread(el, props, nil)
	el.Attr("data-slot", "combobox-trigger")
	el.Attr("type", "button")
	disabled := b.disabled || s.disabled
	if s.inputIn {
		el.Attr("tabindex", "0").Attr("id", s.id).Attr("role", "combobox").Attr("aria-haspopup", "dialog")
	} else {
		el.Attr("tabindex", "-1").Attr("aria-haspopup", "listbox")
	}
	el.Attr("aria-expanded", "false")
	if !s.multiple && !s.hasSelection() {
		el.Attr("data-placeholder", "")
	}
	if disabled {
		el.Attr("data-disabled", "").Bool("disabled", true)
	}
	mark(el, s, "trigger")
	return el
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// renderFlag wraps the element a Trigger renders as (upstream's render
// prop).
type renderFlag struct{ el any }

func (r renderFlag) Element() *h.Element {
	if e, ok := r.el.(h.Elementer); ok {
		return e.Element()
	}
	return h.Button()
}

func (r renderFlag) Render(w io.Writer) error { return r.Element().Render(w) }

func isRender(n any) bool { _, ok := n.(renderFlag); return ok }

// Render makes a Trigger render as el (upstream's render prop), e.g.
// combobox.Trigger(combobox.Render(button.Button(button.VariantOutline)), combobox.Value()).
func Render(el any) h.Node { return renderFlag{el: el} }

// ValueBuilder is returned by Value.
type ValueBuilder struct {
	part
	placeholder string
	hasPH       bool
	children    []any
}

// Value shows the selected item's label (or Placeholder while nothing is
// selected). With Multiple, given a Chip it renders one chip per selected
// item, followed by its other children (a ChipsInput).
func Value(children ...any) *ValueBuilder {
	b := &ValueBuilder{}
	for _, c := range children {
		if p, ok := c.(Placeholder); ok {
			b.placeholder, b.hasPH = string(p), true
			continue
		}
		b.children = append(b.children, c)
	}
	return b
}

// Render renders the value.
func (b *ValueBuilder) Render(w io.Writer) error {
	s := b.state()
	var chip *ChipBuilder
	var rest []any
	for _, c := range b.children {
		if cb, ok := c.(*ChipBuilder); ok && chip == nil {
			chip = cb
			continue
		}
		rest = append(rest, c)
	}
	if chip != nil {
		chip.setState(s)
		tpl := h.Template(h.Data("parity-ignore", ""), chip.Element())
		mark(tpl, s, "chip-template")
		out := []any{tpl}
		for _, it := range s.selected() {
			out = append(out, chip.withLabel(it.Label))
		}
		out = append(out, rest...)
		return h.Group(out...).Render(w)
	}
	text := s.selectedLabel()
	if s.multiple {
		text = ""
		for i, it := range s.selected() {
			if i > 0 {
				text += ", "
			}
			text += it.Label
		}
	}
	if !s.hasSelection() && b.hasPH {
		text = b.placeholder
	}
	start := h.Comment("combobox-value:" + s.id)
	return h.Group(start, text, h.Comment("/combobox-value")).Render(w)
}

// --- Content ---------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	part
	side        Side
	align       Align
	sideOffset  float64
	alignOffset float64
	anchor      *Anchor
}

// Content is the popup holding Empty and List (and an Input when the
// combobox opens from a Trigger), portalled to <body> while open. Children
// may include a Side (default bottom), an Align (default start), a
// SideOffset (default 6), an AlignOffset and an Anchor.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{side: SideBottom, align: AlignStart, sideOffset: 6}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Side, Align, SideOffset, AlignOffset and *Anchor.
func (b *ContentBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Side:
		b.side = v
	case Align:
		b.align = v
	case SideOffset:
		b.sideOffset = float64(v)
	case AlignOffset:
		b.alignOffset = float64(v)
	case *Anchor:
		b.anchor = v
	default:
		return false
	}
	return true
}

const contentClass = "group/combobox-content relative max-h-96 w-(--anchor-width) max-w-(--available-width) min-w-[calc(var(--anchor-width)+--spacing(7))] origin-(--transform-origin) overflow-hidden rounded-md bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 data-[chips=true]:min-w-(--anchor-width) data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 *:data-[slot=input-group]:m-1 *:data-[slot=input-group]:mb-0 *:data-[slot=input-group]:h-8 *:data-[slot=input-group]:border-input/30 *:data-[slot=input-group]:bg-input/30 *:data-[slot=input-group]:shadow-none data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"

// Element resolves the content's template.
func (b *ContentBuilder) Element() *h.Element {
	s := b.state()
	chips := "false"
	if b.anchor != nil {
		chips = "true"
	}
	popup := h.Spread(h.Div(h.Data("slot", "combobox-content"), h.Data("chips", chips), h.Class(lib.Cn(contentClass))), b.Elem(), lib.Cn)
	positioner := h.Div(h.Role("presentation"), h.Class("isolate z-50"), popup)
	opts := map[string]any{
		"side":        string(b.side),
		"align":       string(b.align),
		"sideOffset":  b.sideOffset,
		"alignOffset": b.alignOffset,
	}
	if b.anchor != nil {
		opts["anchor"] = b.anchor.name
	}
	js, _ := json.Marshal(opts)
	tpl := h.Template(h.Data("parity-ignore", ""), h.Attr("data-combobox-content", string(js)), positioner)
	mark(tpl, s, "content")
	return tpl
}

// Render renders the content's template.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- List, Item, Group, Label, Collection, Empty, Separator --------------------------------------

var (
	listBase  = h.Div(h.Data("slot", "combobox-list"), h.Class(lib.Cn("max-h-[min(calc(--spacing(96)---spacing(9)),calc(var(--available-height)---spacing(9)))] scroll-py-1 overflow-y-auto p-1 data-empty:p-0")))
	groupBase = h.Div(h.Data("slot", "combobox-group"))
	labelBase = h.Div(h.Data("slot", "combobox-label"), h.Class(lib.Cn("px-2 py-1.5 text-xs text-muted-foreground pointer-coarse:px-3 pointer-coarse:py-2 pointer-coarse:text-sm")))
	emptyBase = h.Div(h.Role("status"), h.Aria("live", "polite"), h.Aria("atomic", "true"), h.Data("slot", "combobox-empty"), h.Class(lib.Cn("hidden w-full justify-center py-2 text-center text-sm text-muted-foreground group-data-empty/combobox-content:flex")))
	sepBase   = h.Div(h.Data("slot", "combobox-separator"), h.Class(lib.Cn("-mx-1 my-1 h-px bg-border")))
)

// List holds the Items (or Groups of items).
func List(children ...any) *h.Element {
	el := h.Spread(listBase, h.Div(children...), lib.Cn)
	return el.Attr("data-combobox-list", "")
}

// Group groups items under a Label.
func Group(children ...any) *h.Element {
	el := h.Spread(groupBase, h.Div(children...), lib.Cn)
	if c := el.ClassName(); c != "" {
		el.SetClass(lib.Cn(c))
	}
	return el.Attr("data-combobox-group", "")
}

// Label is a Group's heading.
func Label(children ...any) *h.Element {
	return h.Spread(labelBase, h.Div(children...), lib.Cn).Attr("data-combobox-label", "")
}

// Collection renders its children (upstream ComboboxCollection, which maps
// the filtered items of its Group; here the items are given and filtered
// by the behavior).
func Collection(children ...any) h.Node { return h.Group(children...) }

// Empty is shown (its children rendered) while no item matches.
func Empty(children ...any) *h.Element {
	return h.Spread(emptyBase, h.Div(children...), lib.Cn).Attr("data-combobox-empty", "")
}

// Separator separates groups.
func Separator(children ...any) *h.Element {
	return h.Spread(sepBase, h.Div(children...), lib.Cn).Attr("data-combobox-separator", "")
}

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	value    string
	label    string
	hasLabel bool
}

const itemClass = "relative flex w-full cursor-default items-center gap-2 rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none data-highlighted:bg-accent data-highlighted:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

// Item is an option. Children must include its ItemValue and may include
// an ItemLabel.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes ItemValue and ItemLabel.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case ItemValue:
		b.value = string(v)
	case ItemLabel:
		b.label, b.hasLabel = string(v), true
	default:
		return false
	}
	return true
}

func (b *ItemBuilder) item() item {
	l := b.label
	if !b.hasLabel {
		l = b.value
	}
	return item{Value: b.value, Label: l}
}

// Element resolves the item (inside the content's template).
func (b *ItemBuilder) Element() *h.Element {
	base := h.Div(h.Role("option"), h.Data("slot", "combobox-item"), h.Class(lib.Cn(itemClass)))
	el := h.Spread(base, b.Elem(), lib.Cn)
	indicator := h.Span(h.Data("slot", "combobox-item-indicator"), h.Aria("hidden", "true"),
		h.Class("pointer-events-none absolute right-2 flex size-4 items-center justify-center"),
		h.Attr("data-combobox-indicator", ""),
		lucide.Check(h.Class("pointer-events-none size-4 pointer-coarse:size-5")))
	el.SetChildren(append(append([]any{}, el.ChildNodes()...), indicator)...)
	it := b.item()
	js, _ := json.Marshal(it)
	el.Attr("data-combobox-item", string(js))
	return el
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Chips -------------------------------------------------------------------------------------

// ChipsBuilder is returned by Chips.
type ChipsBuilder struct {
	h.Builder[ChipsBuilder]
	part
	anchor *Anchor
}

const chipsClass = "flex min-h-9 flex-wrap items-center gap-1.5 rounded-md border border-input bg-transparent bg-clip-padding px-2.5 py-1.5 text-sm shadow-xs transition-[color,box-shadow] focus-within:border-ring focus-within:ring-[3px] focus-within:ring-ring/50 has-aria-invalid:border-destructive has-aria-invalid:ring-[3px] has-aria-invalid:ring-destructive/20 has-data-[slot=combobox-chip]:px-1.5 dark:bg-input/30 dark:has-aria-invalid:border-destructive/50 dark:has-aria-invalid:ring-destructive/40"

// Chips is the bordered box holding the chips and the ChipsInput (put a
// Value with a Chip and a ChipsInput in it). Children may include the
// *Anchor Content is positioned against.
func Chips(children ...any) *ChipsBuilder {
	b := &ChipsBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes *Anchor.
func (b *ChipsBuilder) AcceptChild(c any) bool {
	if a, ok := c.(*Anchor); ok {
		b.anchor = a
		return true
	}
	return false
}

// Element resolves the chips box.
func (b *ChipsBuilder) Element() *h.Element {
	s := b.state()
	base := h.Div(h.Data("slot", "combobox-chips"), h.Class(lib.Cn(chipsClass)))
	if s.hasSelection() {
		base.Role("toolbar")
	}
	el := h.Spread(base, b.Elem(), lib.Cn)
	if b.anchor != nil {
		el.Attr("data-combobox-anchor", b.anchor.name)
	}
	mark(el, s, "chips")
	return el
}

// Render renders the chips box.
func (b *ChipsBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// ChipBuilder is returned by Chip.
type ChipBuilder struct {
	h.Builder[ChipBuilder]
	part
	showRemove bool
}

const chipClass = "flex h-[calc(--spacing(5.5))] w-fit items-center justify-center gap-1 rounded-sm bg-muted px-1.5 text-xs font-medium whitespace-nowrap text-foreground has-disabled:pointer-events-none has-disabled:cursor-not-allowed has-disabled:opacity-50 has-data-[slot=combobox-chip-remove]:pr-0"

// Chip is the chip of one selected item; give it to Value, which renders
// one per selected item with the item's label. Children may include
// ShowRemove(false).
func Chip(children ...any) *ChipBuilder {
	b := &ChipBuilder{showRemove: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes ShowRemove.
func (b *ChipBuilder) AcceptChild(c any) bool {
	if v, ok := c.(ShowRemove); ok {
		b.showRemove = bool(v)
		return true
	}
	return false
}

func (b *ChipBuilder) withLabel(label string) *h.Element {
	el := b.Element()
	kids := append([]any{label}, el.ChildNodes()...)
	el.SetChildren(kids...)
	return el
}

// Element resolves the chip (without its label).
func (b *ChipBuilder) Element() *h.Element {
	s := b.state()
	base := h.Div(h.TabIndex(-1), h.Data("slot", "combobox-chip"), h.Class(lib.Cn(chipClass)))
	if s.disabled {
		base.Aria("disabled", "true")
	}
	el := h.Spread(base, b.Elem(), lib.Cn)
	if b.showRemove {
		rm := button.Button(button.VariantGhost, button.SizeIconXs,
			h.Class("-ml-1 opacity-50 hover:opacity-100"),
			h.Data("slot", "combobox-chip-remove"),
			h.Type("button"), h.TabIndex(-1), h.Aria("disabled", "false"),
			lucide.X(h.Class("pointer-events-none")),
		).Element()
		if s.disabled {
			rm.Attr("aria-disabled", "true").Attr("data-disabled", "")
		}
		el.SetChildren(append(append([]any{}, el.ChildNodes()...), rm)...)
	}
	return el
}

// Render renders a chip.
func (b *ChipBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// ChipsInputBuilder is returned by ChipsInput.
type ChipsInputBuilder struct {
	h.Builder[ChipsInputBuilder]
	part
}

// ChipsInput is the text input inside Chips.
func ChipsInput(children ...any) *ChipsInputBuilder {
	b := &ChipsInputBuilder{}
	b.Builder = h.NewBuilder(b, "input")
	return b.Children(children...)
}

// Element resolves the input.
func (b *ChipsInputBuilder) Element() *h.Element {
	s := b.state()
	base := comboboxInputAttrs(s.id, "")
	base.Data("slot", "combobox-chip-input").Class(lib.Cn("min-w-16 flex-1 outline-none"))
	el := h.Spread(base, b.Elem(), lib.Cn)
	if s.disabled {
		el.Bool("disabled", true).Attr("data-disabled", "")
	}
	mark(el, s, "input")
	return el
}

// Render renders the input.
func (b *ChipsInputBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
