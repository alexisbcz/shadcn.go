// Package selectui is the shadcn/ui Select: a list of options for the user
// to pick from, triggered by a button (Radix Select, @radix-ui/react-select
// 2.2.6). The package is named selectui because select is a Go keyword.
//
//	selectui.Select(
//		selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Select a fruit"))),
//		selectui.Content(
//			selectui.Group(
//				selectui.Label("Fruits"),
//				selectui.Item(selectui.ItemValue("apple"), "Apple"),
//				selectui.Item(selectui.ItemValue("banana"), "Banana"),
//			),
//		),
//	)
//
// Options are children. On Select: DefaultValue("apple"), DefaultOpen,
// Disabled, Required, Name("fruit"), Form("form-id"),
// AutoComplete("..."), Dir("rtl"). On Trigger: SizeSm, Disabled. On Value:
// Placeholder(...). On Content: PositionPopper (default item-aligned) and,
// with the popper position, the positioning options of package popper
// (SideTop, AlignStart, SideOffset(8), ...; align defaults to center like
// shadcn's SelectContent). On Item: ItemValue("x"), Disabled,
// TextValue("typeahead text").
//
// Markup is what Radix renders. Closed content is not in the DOM: Content
// renders an inert <template> that select.js clones into <body> on open
// and removes on close (Radix Select has no exit animation). The hidden
// native <select> Radix renders for forms is rendered after the trigger
// and kept only when the trigger is inside a form (or Form is set).
//
// Behavior (select.js): opening on pointer down / Enter / Space / arrows,
// typeahead on the trigger (changes the value) and in the content (moves
// focus), arrow / Home / End navigation, selection on pointer up, click,
// Enter or Space, Escape and outside pointer down to close, the
// item-aligned positioning (the selected item over the trigger), scroll
// buttons with auto-scroll, focus trap, scroll lock, aria-hidden outside.
// The template dispatches "shadcn:open-change" (cancelable, detail.open)
// and "shadcn:value-change" (detail.value); scripts can drive a select with
// getSelect from select.js.
package selectui

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/popper"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed select.js
var script string

func init() { runtime.Module("select", script) }

// Marker attributes select.js reads and removes.
const (
	markerItem        = "data-radix-select-item"
	markerTextValue   = "data-radix-select-text-value"
	markerIndicator   = "data-radix-select-item-indicator"
	markerScrollUp    = "data-radix-select-scroll-up"
	markerScrollDown  = "data-radix-select-scroll-down"
	markerValue       = "data-radix-select-value"
	markerValueKids   = "data-radix-select-value-children"
	markerBubble      = "data-radix-select-bubble"
	markerPlaceholder = "data-radix-select-placeholder"
)

// Positioning options accepted by Content with PositionPopper (see package
// popper).
type (
	Side             = popper.Side
	Align            = popper.Align
	SideOffset       = popper.SideOffset
	AlignOffset      = popper.AlignOffset
	AvoidCollisions  = popper.AvoidCollisions
	CollisionPadding = popper.CollisionPadding
	Sticky           = popper.Sticky
	HideWhenDetached = popper.HideWhenDetached
)

// Sides and alignments.
const (
	SideTop     = popper.SideTop
	SideRight   = popper.SideRight
	SideBottom  = popper.SideBottom
	SideLeft    = popper.SideLeft
	AlignStart  = popper.AlignStart
	AlignCenter = popper.AlignCenter
	AlignEnd    = popper.AlignEnd
)

// Position selects how the content is positioned.
type Position string

// Positions: item-aligned (the default) places the selected item over the
// trigger; popper places the content below the trigger like a popover.
const (
	PositionItemAligned Position = "item-aligned"
	PositionPopper      Position = "popper"
)

// Size is the trigger's size.
type Size string

// Trigger sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
)

// ItemValue is an item's value (required, not empty).
type ItemValue string

// TextValue is the text an item is found by with typeahead (default: its
// text content).
type TextValue string

// Flag is a boolean option.
type Flag int

// Flags. Disabled applies to Select, Trigger and Item; Required and
// DefaultOpen to Select.
const (
	Disabled Flag = iota + 1
	Required
	DefaultOpen
)

// Option configures a Select; pass options as children.
type Option func(*Scope)

// DefaultValue is the initially selected value.
func DefaultValue(v string) Option { return func(s *Scope) { s.Value = v; s.HasValue = true } }

// Name is the name of the hidden native select, for forms.
func Name(n string) Option { return func(s *Scope) { s.Name = n } }

// Form associates the select with a form by id.
func Form(id string) Option { return func(s *Scope) { s.Form = id } }

// AutoComplete is the native select's autocomplete hint.
func AutoComplete(v string) Option { return func(s *Scope) { s.AutoComplete = v } }

// Dir is the reading direction ("ltr" or "rtl").
func Dir(d string) Option { return func(s *Scope) { s.Dir = d } }

// PlaceholderOption is returned by Placeholder.
type PlaceholderOption struct{ nodes []any }

// Placeholder is what Value shows while no value is selected.
func Placeholder(nodes ...any) PlaceholderOption { return PlaceholderOption{nodes} }

var (
	triggerClass   = "flex w-fit items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[placeholder]:text-muted-foreground data-[size=default]:h-9 data-[size=sm]:h-8 *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-2 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground"
	contentClass   = "relative z-50 max-h-(--radix-select-content-available-height) min-w-[8rem] origin-(--radix-select-content-transform-origin) overflow-x-hidden overflow-y-auto rounded-md border bg-popover text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"
	contentPopper  = "data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
	viewportPopper = "h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)] scroll-my-1"
	labelClass     = "px-2 py-1.5 text-xs text-muted-foreground"
	itemClass      = "relative flex w-full cursor-default items-center gap-2 rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2"
	separatorClass = "pointer-events-none -mx-1 my-1 h-px bg-border"
	scrollClass    = "flex cursor-default items-center justify-center py-1"
)

// The classes of registry/bases/radix/ui/select.tsx, used inside a
// style.Scope (see package style).
var (
	file              = style.NewFile("ui/select.tsx")
	stGroup           = file.Cn("cn-select-group", style.Arg)
	stTrigger         = file.Cn("cn-select-trigger flex w-fit items-center justify-between whitespace-nowrap outline-none disabled:cursor-not-allowed disabled:opacity-50 *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	stTriggerIcon     = file.Attr("cn-select-trigger-icon pointer-events-none")
	stContent         = file.Cn("cn-select-content cn-menu-target cn-menu-translucent relative z-50 max-h-(--radix-select-content-available-height) origin-(--radix-select-content-transform-origin) overflow-x-hidden overflow-y-auto data-[align-trigger=true]:animate-none", style.Arg, style.Arg)
	stViewport        = file.Cn("cn-select-viewport data-[position=popper]:h-(--radix-select-trigger-height) data-[position=popper]:w-full data-[position=popper]:min-w-(--radix-select-trigger-width)", style.Arg)
	stLabel           = file.Cn("cn-select-label", style.Arg)
	stItem            = file.Cn("cn-select-item relative flex w-full cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	stItemIndicator   = file.Attr("cn-select-item-indicator")
	stIndicatorIcon   = file.Attr("cn-select-item-indicator-icon pointer-events-none")
	stSeparator       = file.Cn("cn-select-separator pointer-events-none", style.Arg)
	stScrollUpButton  = file.Cn("cn-select-scroll-up-button", style.Arg)
	stScrollDownButon = file.Cn("cn-select-scroll-down-button", style.Arg)
)

// styled is a node whose markup depends on the enclosing style.Scope.
type styled func(st style.Style) *h.Element

func (f styled) Render(w io.Writer) error {
	st, _ := style.FromWriter(w)
	return f(st).Render(w)
}

// visuallyHidden is Radix's VISUALLY_HIDDEN_STYLES as React serializes it.
const visuallyHidden = "position: absolute; border: 0px; width: 1px; height: 1px; padding: 0px; margin: -1px; overflow: hidden; clip: rect(0px, 0px, 0px, 0px); white-space: nowrap; overflow-wrap: normal;"

// Scope is the state a Select shares with its parts (Radix's
// SelectContext).
type Scope struct {
	ContentID    string
	Value        string
	HasValue     bool
	DefaultOpen  bool
	Disabled     bool
	Required     bool
	Name         string
	Form         string
	AutoComplete string
	Dir          string

	items []*ItemBuilder
	value *ValueBuilder
}

// NewScope returns a scope with a fresh content id.
func NewScope() *Scope { return &Scope{ContentID: uid.New("radix")} }

func (s *Scope) dir() string {
	if s.Dir == "" {
		return "ltr"
	}
	return s.Dir
}

func (s *Scope) showPlaceholder() bool { return !s.HasValue || s.Value == "" }

func (s *Scope) selectedItem() *ItemBuilder {
	if s.showPlaceholder() {
		return nil
	}
	for _, it := range s.items {
		if it.value == s.Value {
			return it
		}
	}
	return nil
}

// Scoped is embedded by the parts: Select hands them its scope.
type Scoped struct{ scope *Scope }

// SetSelectScope is called by Select.
func (s *Scoped) SetSelectScope(sc *Scope) { s.scope = sc }

// Scope returns the select's scope; a part outside any Select gets one of
// its own.
func (s *Scoped) Scope() *Scope {
	if s.scope == nil {
		s.scope = NewScope()
	}
	return s.scope
}

type participant interface{ SetSelectScope(*Scope) }

func stateOf(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}

// --- Select --------------------------------------------------------------------------------

// RootBuilder is returned by Select.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts []Option
}

// Select is the root. It renders its children, then the hidden native
// select for forms.
func Select(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Options and Flags.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Option:
		b.opts = append(b.opts, v)
	case Flag:
		b.opts = append(b.opts, func(s *Scope) {
			switch v {
			case Disabled:
				s.Disabled = true
			case Required:
				s.Required = true
			case DefaultOpen:
				s.DefaultOpen = true
			}
		})
	default:
		return false
	}
	return true
}

// DefaultValue sets the initially selected value.
func (b *RootBuilder) DefaultValue(v string) *RootBuilder {
	b.opts = append(b.opts, DefaultValue(v))
	return b
}

// Name sets the native select's name.
func (b *RootBuilder) Name(n string) *RootBuilder { b.opts = append(b.opts, Name(n)); return b }

// Render renders the children and the native select.
func (b *RootBuilder) Render(w io.Writer) error {
	sc := NewScope()
	for _, o := range b.opts {
		o(sc)
	}
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case *ItemBuilder:
			c.SetSelectScope(sc)
			sc.items = append(sc.items, c)
		case *ValueBuilder:
			c.SetSelectScope(sc)
			sc.value = c
		case participant:
			c.SetSelectScope(sc)
		}
		return true
	})
	return h.Group(append(kids, bubbleSelect(sc))...).Render(w)
}

var tagRE = regexp.MustCompile(`<[^>]*>`)

// textContent approximates the DOM textContent of nodes.
func textContent(nodes []any) string {
	var buf bytes.Buffer
	if err := h.Group(nodes...).Render(&buf); err != nil {
		return ""
	}
	return html.UnescapeString(tagRE.ReplaceAllString(buf.String(), ""))
}

// bubbleSelect is Radix's SelectBubbleInput with its native options.
func bubbleSelect(sc *Scope) *h.Element {
	sel := h.Select(
		h.Attr(markerBubble, sc.ContentID),
		h.Aria("hidden", "true"),
		h.If(sc.Required, h.Attr("required", "")),
		h.TabIndex(-1),
		h.If(sc.Name != "", h.Attr("name", sc.Name)),
		h.If(sc.AutoComplete != "", h.Attr("autocomplete", sc.AutoComplete)),
		h.If(sc.Disabled, h.Attr("disabled", "")),
		h.If(sc.Form != "", h.Attr("form", sc.Form)),
		h.Style(visuallyHidden),
	)
	if !sc.HasValue {
		sel.Children(h.Option(h.Value("")))
	}
	for _, it := range sc.items {
		sel.Children(h.Option(h.Value(it.value), h.If(it.disabled, h.Attr("disabled", "")), h.If(sc.HasValue && it.value == sc.Value, h.Attr("selected", "")), textContent(it.Elem().ChildNodes())))
	}
	return sel
}

// part holds what the parts share.
type part struct {
	Scoped
	asChild bool
	st      style.Style // the enclosing style.Scope, "" for new-york-v4
	rtl     bool        // inside style.ScopeRTL: the ui-rtl classes
}

func (p *part) setStyle(w io.Writer) {
	p.st, _ = style.FromWriter(w)
	p.rtl = p.st != "" && style.RTLFromWriter(w)
}

// dir maps a style's classes to their ui-rtl variant inside
// style.ScopeRTL.
func (p *part) dir(class string) string {
	if p.rtl {
		return style.RTLClass(class)
	}
	return class
}

func (p *part) accept(c any) bool {
	switch v := c.(type) {
	case lib.AsChildFlag:
		p.asChild = true
	case *Scope:
		p.SetSelectScope(v)
	default:
		return false
	}
	return true
}

func (p *part) resolve(name string, defaults, user *h.Element) *h.Element {
	props := h.Spread(defaults, user, lib.Cn)
	if !p.asChild {
		return props
	}
	kids := props.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("select %s: AsChild needs exactly one element child, got %d children", name, len(kids)))
	}
	out, err := aschild.Slot(props.Clone().SetChildren(), kids[0])
	if err != nil {
		return h.Errored(fmt.Errorf("select %s: AsChild: %w", name, err))
	}
	return out
}

// --- Trigger -------------------------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	part
	size     Size
	disabled bool
}

// Trigger is the button that opens the select. Children may include a
// Size and Disabled.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{size: SizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Size, Disabled, lib.AsChild and *Scope.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Size:
		b.Size(v)
		return true
	case Flag:
		if v == Disabled {
			b.disabled = true
			return true
		}
		return false
	}
	return b.accept(c)
}

// Size sets the trigger's size.
func (b *TriggerBuilder) Size(s Size) *TriggerBuilder {
	if s == "" {
		s = SizeDefault
	}
	b.size = s
	return b
}

// Disabled disables the trigger.
func (b *TriggerBuilder) Disabled(on bool) *TriggerBuilder { b.disabled = on; return b }

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	sc := b.Scope()
	user := b.Elem().Clone()
	disabled := sc.Disabled || b.disabled || user.Has("disabled")
	user.Remove("disabled")
	d := h.Button(
		h.Type("button"),
		h.Role("combobox"),
		h.Aria("controls", sc.ContentID),
		h.Aria("expanded", fmt.Sprint(sc.DefaultOpen)),
		h.If(sc.Required, h.Aria("required", "true")),
		h.Aria("autocomplete", "none"),
		h.Dir(sc.dir()),
		h.Data("state", stateOf(sc.DefaultOpen)),
		h.If(disabled, h.Attr("disabled", "")),
		h.If(disabled, h.Data("disabled", "")),
		h.If(sc.showPlaceholder(), h.Data("placeholder", "")),
		h.Data("slot", "select-trigger"),
		h.Data("size", string(b.size)),
		h.Class(lib.Cn(triggerClass)),
	)
	if b.st != "" {
		d.Attr("class", b.dir(stTrigger.Class(b.st, "")))
		user.Children(lucide.ChevronDown(h.Class(stTriggerIcon.Class(b.st)), h.Aria("hidden", "true")))
	} else {
		user.Children(lucide.ChevronDown(h.Class("size-4 opacity-50"), h.Aria("hidden", "true")))
	}
	return b.resolve("Trigger", d, user)
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Value ---------------------------------------------------------------------------------

// ValueBuilder is returned by Value.
type ValueBuilder struct {
	h.Builder[ValueBuilder]
	part
	placeholder []any
}

// Value shows the selected item's text, or the Placeholder. Children, if
// any, replace the item's text while a value is selected. Like Radix, the
// class and style props are ignored.
func Value(children ...any) *ValueBuilder {
	b := &ValueBuilder{}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// AcceptChild consumes Placeholder, lib.AsChild and *Scope.
func (b *ValueBuilder) AcceptChild(c any) bool {
	if p, ok := c.(PlaceholderOption); ok {
		b.placeholder = p.nodes
		return true
	}
	return b.accept(c)
}

// Placeholder sets the placeholder.
func (b *ValueBuilder) Placeholder(nodes ...any) *ValueBuilder { b.placeholder = nodes; return b }

// Element resolves the value span.
func (b *ValueBuilder) Element() *h.Element {
	sc := b.Scope()
	user := b.Elem().Clone()
	user.Remove("class")
	user.Remove("style")
	kids := user.ChildNodes()
	hasChildren := len(kids) > 0
	user.SetChildren()
	d := h.Span(
		h.Attr(markerValue, sc.ContentID),
		h.If(hasChildren, h.Attr(markerValueKids, "")),
		h.Data("slot", "select-value"),
		h.Style("pointer-events: none;"),
	)
	el := b.resolve("Value", d, user)
	switch {
	case sc.showPlaceholder():
		el.Children(b.placeholder...)
	case hasChildren:
		el.Children(kids...)
	default:
		if it := sc.selectedItem(); it != nil {
			el.Children(it.Elem().ChildNodes()...)
		}
	}
	return el
}

// Render renders the value span.
func (b *ValueBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	part
	position Position
	opts     popper.Options
}

// Content is the select's listbox, portalled to <body> while open.
// Accepts a Position and, with PositionPopper, positioning options.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{position: PositionItemAligned, opts: popper.Defaults()}
	b.opts.Align = popper.AlignCenter
	b.opts.CollisionPadding = float64(10)
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Position, positioning options and *Scope.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if p, ok := c.(Position); ok {
		b.position = p
		return true
	}
	if b.opts.Accept(c) {
		return true
	}
	return b.accept(c)
}

// Position sets the position mode.
func (b *ContentBuilder) Position(p Position) *ContentBuilder { b.position = p; return b }

// Options returns the popper positioning options for changes.
func (b *ContentBuilder) Options() *popper.Options { return &b.opts }

// Element resolves the template holding the content as it is when open.
func (b *ContentBuilder) Element() *h.Element {
	sc := b.Scope()
	isPopper := b.position == PositionPopper
	user := b.Elem().Clone()
	kids := user.ChildNodes()
	user.SetChildren()
	userStyle, _ := user.GetAttr("style")
	user.Remove("style")

	style := "display: flex; flex-direction: column; outline: none;"
	if isPopper {
		style = "box-sizing: border-box; " + style
	} else {
		style = "box-sizing: border-box; max-height: 100%; " + style
	}
	if userStyle != "" {
		style += " " + userStyle
	}
	d := h.Div(
		h.Role("listbox"),
		h.ID(sc.ContentID),
		h.Data("state", "open"),
		h.Dir(sc.dir()),
		h.Data("slot", "select-content"),
		h.Class(lib.Cn(contentClass, h.If(isPopper, contentPopper))),
		h.TabIndex(-1),
	)
	if b.st != "" {
		d.Attr("data-align-trigger", fmt.Sprint(!isPopper))
		d.Attr("class", b.dir(stContent.Class(b.st, h.IfElse(isPopper, contentPopper, ""), "")))
	}
	content := h.Spread(d, user, lib.Cn)
	if !strings.HasSuffix(strings.TrimSpace(style), ";") {
		style += ";"
	}
	content.Attr("style", style)
	if isPopper {
		popper.AppendStyle(content, popper.ContentVars("select"))
	}
	viewportClass := lib.Cn("p-1", h.If(isPopper, viewportPopper))
	if b.st != "" {
		viewportClass = stViewport.Class(b.st, "")
	}
	viewport := h.Div(
		h.If(b.st != "", h.Data("position", string(b.position))),
		h.Data("radix-select-viewport", ""),
		h.Role("presentation"),
		h.Class(viewportClass),
		h.Style("position: relative; flex: 1 1 0%; overflow: hidden auto;"),
		h.Group(kids...),
	)
	content.Children(scrollUpButton(b.st), viewport, scrollDownButton(b.st))

	var wrapper *h.Element
	if isPopper {
		o := b.opts
		o.Dir = sc.dir()
		wrapper = popper.Wrap(o, content)
	} else {
		wrapper = h.Div(h.Style("display: flex; flex-direction: column; position: fixed;"), content)
	}
	var ph *h.Element
	if sc.value != nil {
		ph = h.Span(h.Attr(markerPlaceholder, ""), h.Group(sc.value.placeholder...))
	}
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-select", ""),
		h.Data("content-id", sc.ContentID),
		h.Data("position", string(b.position)),
		h.Data("dir", sc.dir()),
		h.If(sc.HasValue, h.Data("value", sc.Value)),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		h.If(sc.Disabled, h.Data("disabled", "")),
		wrapper,
		ph,
	)
}

// Render renders the template.
func (b *ContentBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Group, Label, Separator ---------------------------------------------------------------

// GroupBuilder is returned by Group.
type GroupBuilder struct {
	h.Builder[GroupBuilder]
	part
}

// Group groups items under a Label.
func Group(children ...any) *GroupBuilder {
	b := &GroupBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *GroupBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the group.
func (b *GroupBuilder) Element() *h.Element {
	id := uid.New("radix")
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *GroupBuilder:
			return false
		case *LabelBuilder:
			c.id = id
		}
		return true
	})
	d := h.Div(h.Role("group"), h.Aria("labelledby", id), h.Data("slot", "select-group"))
	if b.st != "" {
		d.Attr("class", stGroup.Class(b.st, ""))
	}
	return b.resolve("Group", d, b.Elem())
}

// Render renders the group.
func (b *GroupBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// LabelBuilder is returned by Label.
type LabelBuilder struct {
	h.Builder[LabelBuilder]
	part
	id string
}

// Label labels its Group.
func Label(children ...any) *LabelBuilder {
	b := &LabelBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *LabelBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the label.
func (b *LabelBuilder) Element() *h.Element {
	if b.id == "" {
		b.id = uid.New("radix")
	}
	d := h.Div(h.ID(b.id), h.Data("slot", "select-label"), h.Class(lib.Cn(labelClass)))
	if b.st != "" {
		d.Attr("class", stLabel.Class(b.st, ""))
	}
	return b.resolve("Label", d, b.Elem())
}

// Render renders the label.
func (b *LabelBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// Separator separates items.
func Separator(children ...any) h.Node {
	return styled(func(st style.Style) *h.Element {
		cls := lib.Cn(separatorClass)
		if st != "" {
			cls = stSeparator.Class(st, "")
		}
		return h.Spread(h.Div(h.Aria("hidden", "true"), h.Data("slot", "select-separator"), h.Class(cls)), h.Div(children...), lib.Cn)
	})
}

// --- Item ----------------------------------------------------------------------------------

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	part
	value     string
	disabled  bool
	textValue *string
}

// Item is an option. Children must include its ItemValue and may include
// Disabled and a TextValue.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes ItemValue, TextValue, Disabled, lib.AsChild and
// *Scope.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case ItemValue:
		b.value = string(v)
		return true
	case TextValue:
		s := string(v)
		b.textValue = &s
		return true
	case Flag:
		if v == Disabled {
			b.disabled = true
			return true
		}
		return false
	}
	return b.accept(c)
}

// Value sets the item's value.
func (b *ItemBuilder) Value(v string) *ItemBuilder { b.value = v; return b }

// Disabled disables the item.
func (b *ItemBuilder) Disabled(on bool) *ItemBuilder { b.disabled = on; return b }

// TextValue sets the typeahead text.
func (b *ItemBuilder) TextValue(v string) *ItemBuilder { b.textValue = &v; return b }

// Element resolves the item as it renders while the select is open.
func (b *ItemBuilder) Element() *h.Element {
	if b.value == "" {
		return h.Errored(fmt.Errorf("select Item: a value that is not an empty string is required"))
	}
	sc := b.Scope()
	selected := sc.HasValue && sc.Value == b.value
	textID := uid.New("radix")
	user := b.Elem().Clone()
	kids := user.ChildNodes()
	user.SetChildren()
	d := h.Div(
		h.Attr(markerItem, b.value),
		h.If(b.textValue != nil, h.Attr(markerTextValue, deref(b.textValue))),
		h.Role("option"),
		h.Aria("labelledby", textID),
		h.Aria("selected", "false"),
		h.Data("state", h.IfElse(selected, "checked", "unchecked")),
		h.If(b.disabled, h.Aria("disabled", "true")),
		h.If(b.disabled, h.Data("disabled", "")),
		h.If(!b.disabled, h.TabIndex(-1)),
		h.Data("slot", "select-item"),
		h.Class(lib.Cn(itemClass)),
		h.Data("radix-collection-item", ""),
	)
	if b.st != "" {
		d.Attr("class", b.dir(stItem.Class(b.st, "")))
	}
	el := b.resolve("Item", d, user)
	if b.st != "" {
		el.Children(
			h.Span(h.Class(b.dir(stItemIndicator.Class(b.st))),
				h.Span(h.Attr(markerIndicator, ""), h.Aria("hidden", "true"), lucide.Check(h.Class(stIndicatorIcon.Class(b.st)))),
			),
			h.Span(h.ID(textID), h.Group(kids...)),
		)
		return el
	}
	el.Children(
		h.Span(h.Data("slot", "select-item-indicator"), h.Class("absolute right-2 flex size-3.5 items-center justify-center"),
			h.Span(h.Attr(markerIndicator, ""), h.Aria("hidden", "true"), lucide.Check(h.Class("size-4"))),
		),
		h.Span(h.ID(textID), h.Group(kids...)),
	)
	return el
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Scroll buttons ------------------------------------------------------------------------

// ScrollUpButton is shown at the top of the content while it can scroll
// up. Content renders it; it is exported for completeness.
func ScrollUpButton(children ...any) *h.Element { return scrollUpButton("", children...) }

// ScrollDownButton is shown at the bottom of the content while it can
// scroll down.
func ScrollDownButton(children ...any) *h.Element { return scrollDownButton("", children...) }

func scrollUpButton(st style.Style, children ...any) *h.Element {
	cls := lib.Cn(scrollClass)
	if st != "" {
		cls = stScrollUpButton.Class(st, "")
	}
	d := h.Div(h.Attr(markerScrollUp, ""), h.Aria("hidden", "true"), h.Data("slot", "select-scroll-up-button"), h.Class(cls), h.Style("flex-shrink: 0;"))
	el := h.Spread(d, h.Div(children...), lib.Cn)
	if len(el.ChildNodes()) == 0 {
		el.Children(lucide.ChevronUp(h.If(st == "", h.Class("size-4"))))
	}
	return el
}

func scrollDownButton(st style.Style, children ...any) *h.Element {
	cls := lib.Cn(scrollClass)
	if st != "" {
		cls = stScrollDownButon.Class(st, "")
	}
	d := h.Div(h.Attr(markerScrollDown, ""), h.Aria("hidden", "true"), h.Data("slot", "select-scroll-down-button"), h.Class(cls), h.Style("flex-shrink: 0;"))
	el := h.Spread(d, h.Div(children...), lib.Cn)
	if len(el.ChildNodes()) == 0 {
		el.Children(lucide.ChevronDown(h.If(st == "", h.Class("size-4"))))
	}
	return el
}
