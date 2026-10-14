// Package menu is the Go side of Radix's Menu primitive
// (@radix-ui/react-menu 2.1.16), shared by the shadcn/ui DropdownMenu,
// ContextMenu and Menubar: the markup of the menu content and of every
// part that lives in it (items, checkbox and radio items, labels, groups,
// separators, submenus), styled per component by a Kind.
//
// The behavior lives in components/ui/runtime/js/menu.js; each component
// package adds its root and trigger and registers a module that drives the
// primitive (dropdown-menu.js, context-menu.js, menubar.js).
//
// Markup is what React renders. Closed content is not in the DOM: Content
// renders an inert <template> holding the popper wrapper with the content
// exactly as it is when open, before it is positioned; menu.js clones it
// into <body> when the menu opens and removes it after the exit animation.
// A submenu's content is a nested template, cloned in place (or into
// <body> inside a Portal) when the submenu opens.
//
// Checked state lives in the page: checkbox items toggle and radio groups
// change value when selected, and keep their state across openings,
// unless they are Controlled (then only scripts change it). The behavior
// dispatches, on the live elements:
//
//	"shadcn:select"         item selected (cancelable: keeps the menu open)
//	"shadcn:checked-change" checkbox item (cancelable), detail.checked
//	"shadcn:value-change"   radio group (cancelable), detail.value
package menu

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/popper"
	h "shadcn.go/libhtml"
)

// Marker attributes menu.js reads and removes.
const (
	markerCheckbox   = "data-radix-menu-checkbox"    // checkbox item: state key
	markerRadioGroup = "data-radix-menu-radio-group" // radio group: state key
	markerRadioValue = "data-radix-menu-radio-value" // radio group value / radio item value
	markerControlled = "data-radix-menu-controlled"  // state only changes from scripts
	markerIndicator  = "data-radix-menu-indicator"   // ItemIndicator
	markerTextValue  = "data-radix-menu-text-value"  // item textValue
)

// Kind is what differs between the components built on the primitive:
// data-slot prefix, CSS variable namespace and the shadcn classes.
type Kind struct {
	// Name prefixes data-slot values ("dropdown-menu" gives
	// "dropdown-menu-item") and the content's CSS variables
	// (--radix-dropdown-menu-content-transform-origin, ...).
	Name string
	// Classes, before the user's className.
	Content, SubContent, Item, CheckboxItem, RadioItem, Label, Separator, Shortcut, SubTrigger string
	// SubTriggerIcon is the class of the chevron after a SubTrigger's children.
	SubTriggerIcon string
	// ContentDefaults sets the component's default positioning (applied
	// before the user's options); ContentForced overrides them after
	// (ContextMenu forces side, offset and alignment).
	ContentDefaults func(*popper.Options)
	ContentForced   func(*popper.Options)
	// ContentAttrs are extra attributes on the root and sub content
	// (Menubar's data-radix-menubar-content); SubTriggerAttrs on
	// sub triggers (data-radix-menubar-subtrigger).
	ContentAttrs, SubTriggerAttrs []h.Prop
	// ContentIDs renders id and aria-labelledby on the root content
	// (DropdownMenu, Menubar; ContextMenu has no trigger id).
	ContentIDs bool
	// CheckboxIndicator and RadioIndicator, when set, return the node
	// rendered before a checkbox or radio item's children (a wrapper
	// around an Indicator); SubTriggerChevron the node after a sub
	// trigger's children.
	CheckboxIndicator, RadioIndicator, SubTriggerChevron func() h.Node
	// Styled, when set, returns the kind to use when style st is in scope
	// (style.Scope): the markup and classes of registry/bases/radix.
	Styled func(st style.Style) *Kind
}

// For returns the kind to render with: k, or its styled variant when a
// style is in scope of w (the writer libhtml passes to Render).
func (k *Kind) For(w io.Writer) *Kind {
	if k.Styled == nil {
		return k
	}
	if st, ok := style.FromWriter(w); ok {
		if sk := k.Styled(st); sk != nil {
			return sk
		}
	}
	return k
}

// Cached returns build, memoized per style: use it for Kind.Styled.
func Cached(build func(st style.Style) *Kind) func(st style.Style) *Kind {
	var m sync.Map // style.Style -> *Kind
	return func(st style.Style) *Kind {
		if k, ok := m.Load(st); ok {
			return k.(*Kind)
		}
		k, _ := m.LoadOrStore(st, build(st))
		return k.(*Kind)
	}
}

// withChild returns a copy of user with node prepended (or appended) to
// its children; user itself when node is nil.
func withChild(user *h.Element, node func() h.Node, first bool) *h.Element {
	if node == nil {
		return user
	}
	kids := user.ChildNodes()
	n := make([]any, 0, len(kids)+1)
	if first {
		n = append(append(n, node()), kids...)
	} else {
		n = append(append(n, kids...), node())
	}
	return user.Clone().SetChildren(n...)
}

// Scope is the state a component's root shares with the parts of its menu.
type Scope struct {
	Kind        *Kind
	ContentID   string
	TriggerID   string
	Modal       bool
	DefaultOpen bool
	Dir         string // "ltr" (default) or "rtl"
	// Template carries the component's own attributes for its behavior
	// (data-radix-dropdown-menu, ...).
	Template []h.Prop
}

// NewScope returns a scope with fresh ids.
func NewScope(k *Kind) *Scope {
	return &Scope{Kind: k, ContentID: uid.New("radix"), TriggerID: uid.New("radix"), Modal: true}
}

// Direction returns the reading direction, "ltr" when unset.
func (s *Scope) Direction() string {
	if s.Dir == "rtl" {
		return "rtl"
	}
	return "ltr"
}

// Participant is implemented by the parts: Share hands them the scope.
type Participant interface{ SetMenuScope(*Scope) }

// Root is implemented by component roots: Share does not cross them.
type Root interface{ MenuRoot() }

// Share hands sc to every part under children, not crossing nested roots.
func Share(children []any, sc *Scope) {
	tree.Walk(children, func(n any) bool {
		if _, ok := n.(Root); ok {
			return false
		}
		if p, ok := n.(Participant); ok {
			p.SetMenuScope(sc)
		}
		return true
	})
}

// Scoped is embedded by the parts.
type Scoped struct{ scope *Scope }

// SetMenuScope is called by the component root.
func (s *Scoped) SetMenuScope(sc *Scope) { s.scope = sc }

// MenuScope returns the scope; a part outside any root gets a default one.
func (s *Scoped) MenuScope() *Scope {
	if s.scope == nil {
		s.scope = NewScope(&Kind{Name: "menu"})
	}
	return s.scope
}

// --- options -------------------------------------------------------------------------------

// Inset indents an item, label or sub trigger to align with checkbox items.
type Inset bool

// Variant is an item's look.
type Variant string

// Item variants.
const (
	VariantDefault     Variant = "default"
	VariantDestructive Variant = "destructive"
)

// Disabled disables an item: it is skipped by keyboard navigation and
// typeahead and cannot be selected.
type Disabled bool

// Checked is a checkbox item's state (true or false); see Indeterminate.
type Checked bool

// CheckedState values for CheckboxItem.
type indeterminate struct{}

// Indeterminate puts a checkbox item in the indeterminate state
// (aria-checked="mixed").
var Indeterminate = indeterminate{}

// Value is a radio group's selected value or a radio item's value.
type Value string

// TextValue is the text typeahead matches for an item (default: its text
// content).
type TextValue string

// Loop makes keyboard navigation wrap from the last item to the first.
type Loop bool

// controlled marks state that only scripts change.
type controlled struct{}

// Controlled makes a checkbox item's checked state or a radio group's value
// change only through scripts, like a React `checked` / `value` prop without
// an onCheckedChange / onValueChange handler. Selecting the item still
// dispatches the change events.
var Controlled = controlled{}

// ForceMount keeps an ItemIndicator rendered when unchecked.
type ForceMount bool

func cnClass(base string, user *h.Element) string { return lib.Cn(base, user.ClassName()) }

// spread merges the defaults with the user's props; classes are cn(base, user).
func spread(defaults, user *h.Element) *h.Element { return h.Spread(defaults, user, lib.Cn) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	Scoped
	user []any
	loop bool
}

// Content is the menu content, portalled to <body> while open and
// positioned against the trigger. Accepts popper options (side, align,
// offsets, collision handling) and Loop.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	b.user = children
	return b.Children(children...)
}

// AcceptChild consumes popper options, Loop and *Scope.
func (b *ContentBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Loop:
		b.loop = bool(v)
		return true
	case *Scope:
		b.SetMenuScope(v)
		return true
	}
	// Positioning options are read in Element, over the kind's defaults.
	var o popper.Options
	return o.Accept(c)
}

// Element resolves the template holding the content.
func (b *ContentBuilder) Element() *h.Element {
	sc := b.MenuScope()
	k := sc.Kind
	opts := popper.Defaults()
	if k.ContentDefaults != nil {
		k.ContentDefaults(&opts)
	}
	// User options (re-applied over the component's defaults).
	for _, c := range b.user {
		opts.Accept(c)
	}
	if k.ContentForced != nil {
		k.ContentForced(&opts)
	}
	opts.Dir = sc.Direction()
	user := b.Elem()
	d := h.Div(
		h.Role("menu"),
		h.Aria("orientation", "vertical"),
		h.Data("state", "open"),
		h.Data("radix-menu-content", ""),
		h.Dir(sc.Direction()),
		h.If(k.ContentIDs, h.ID(sc.ContentID)),
		h.If(k.ContentIDs, h.Aria("labelledby", sc.TriggerID)),
		h.Props(k.ContentAttrs...),
		h.Data("slot", k.Name+"-content"),
		h.Class(k.Content),
		h.TabIndex(-1),
		h.Data("orientation", "vertical"),
	)
	el := spread(d, user)
	el.Attr("style", contentStyle(k.Name, user))
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Props(sc.Template...),
		h.Data("content-id", sc.ContentID),
		h.Data("trigger-id", sc.TriggerID),
		h.Data("dir", sc.Direction()),
		h.If(!sc.Modal, h.Data("modal", "false")),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		h.If(b.loop, h.Data("loop", "")),
		popper.Wrap(opts, el),
	)
}

// Render renders the template.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// contentStyle is {outline: none, ...user style, ...re-namespaced vars}.
func contentStyle(name string, user *h.Element) string {
	st := "outline: none;"
	if us, ok := user.GetAttr("style"); ok && strings.TrimSpace(us) != "" {
		us = strings.TrimSpace(us)
		if !strings.HasSuffix(us, ";") {
			us += ";"
		}
		st += " " + us
	}
	return st + " " + popper.ContentVars(name) + ";"
}

// --- Group, Label, Separator, Shortcut ---------------------------------------------------------

// GroupBuilder is returned by Group.
type GroupBuilder struct {
	h.Builder[GroupBuilder]
	Scoped
}

// Group groups items (role="group").
func Group(children ...any) *GroupBuilder {
	b := &GroupBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the group.
func (b *GroupBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	return spread(h.Div(h.Role("group"), h.Data("slot", k.Name+"-group")), b.Elem())
}

// Render renders the group.
func (b *GroupBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// LabelBuilder is returned by Label.
type LabelBuilder struct {
	h.Builder[LabelBuilder]
	Scoped
	inset bool
}

// Label labels a group of items. Accepts Inset.
func Label(children ...any) *LabelBuilder {
	b := &LabelBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Inset.
func (b *LabelBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Inset); ok {
		b.inset = bool(v)
		return true
	}
	return false
}

// Element resolves the label.
func (b *LabelBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	return spread(h.Div(h.Data("slot", k.Name+"-label"), h.If(b.inset, h.Data("inset", "true")), h.Class(k.Label)), b.Elem())
}

// Render renders the label.
func (b *LabelBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// SeparatorBuilder is returned by Separator.
type SeparatorBuilder struct {
	h.Builder[SeparatorBuilder]
	Scoped
}

// Separator separates items.
func Separator(children ...any) *SeparatorBuilder {
	b := &SeparatorBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the separator.
func (b *SeparatorBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	return spread(h.Div(h.Role("separator"), h.Aria("orientation", "horizontal"), h.Data("slot", k.Name+"-separator"), h.Class(k.Separator)), b.Elem())
}

// Render renders the separator.
func (b *SeparatorBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// ShortcutBuilder is returned by Shortcut.
type ShortcutBuilder struct {
	h.Builder[ShortcutBuilder]
	Scoped
}

// Shortcut shows a keyboard shortcut after an item's text.
func Shortcut(children ...any) *ShortcutBuilder {
	b := &ShortcutBuilder{}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// Element resolves the shortcut.
func (b *ShortcutBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	return spread(h.Span(h.Data("slot", k.Name+"-shortcut"), h.Class(k.Shortcut)), b.Elem())
}

// Render renders the shortcut.
func (b *ShortcutBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- items ---------------------------------------------------------------------------------

// itemBase holds what every item shares (MenuItemImpl).
type itemBase struct {
	Scoped
	disabled  bool
	textValue *string
	asChild   bool
}

func (p *itemBase) accept(c any) bool {
	switch v := c.(type) {
	case Disabled:
		p.disabled = bool(v)
	case TextValue:
		s := string(v)
		p.textValue = &s
	case lib.AsChildFlag:
		p.asChild = true
	default:
		return false
	}
	return true
}

// implAttrs are MenuItemImpl's attributes (role is set by the caller).
func (p *itemBase) implAttrs(role string) []any {
	return []any{
		h.Role(role),
		h.If(p.disabled, h.Aria("disabled", "true")),
		h.If(p.disabled, h.Data("disabled", "")),
		h.TabIndex(-1),
		h.Data("orientation", "vertical"),
		h.Data("radix-collection-item", ""),
		h.If(p.textValue != nil, h.Attr(markerTextValue, deref(p.textValue))),
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (p *itemBase) resolve(name string, el *h.Element) *h.Element {
	if !p.asChild {
		return el
	}
	return aschild.Resolve(name, el)
}

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	itemBase
	inset   bool
	variant Variant
}

// Item is a menu item. Accepts Inset, Variant, Disabled, TextValue and
// lib.AsChild (e.g. to render a link).
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{variant: VariantDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the item options.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Inset:
		b.inset = bool(v)
		return true
	case Variant:
		if v == "" {
			v = VariantDefault
		}
		b.variant = v
		return true
	}
	return b.accept(c)
}

// Element resolves the item.
func (b *ItemBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	d := h.Div(b.implAttrs("menuitem")...).Children(
		h.Data("slot", k.Name+"-item"),
		h.If(b.inset, h.Data("inset", "true")),
		h.Data("variant", string(b.variant)),
		h.Class(k.Item),
	)
	return b.resolve(k.Name+" Item", spread(d, b.Elem()))
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// CheckedState is a checkbox item's state.
type CheckedState string

// Checked states.
const (
	StateChecked       CheckedState = "checked"
	StateUnchecked     CheckedState = "unchecked"
	StateIndeterminate CheckedState = "indeterminate"
)

func ariaChecked(s CheckedState) string {
	switch s {
	case StateChecked:
		return "true"
	case StateIndeterminate:
		return "mixed"
	}
	return "false"
}

// CheckboxItemBuilder is returned by CheckboxItem.
type CheckboxItemBuilder struct {
	h.Builder[CheckboxItemBuilder]
	itemBase
	state      CheckedState
	controlled bool
	key        string
}

// CheckboxItem is an item that toggles. Accepts Checked, Indeterminate,
// Controlled, Disabled and TextValue.
func CheckboxItem(children ...any) *CheckboxItemBuilder {
	b := &CheckboxItemBuilder{state: StateUnchecked}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the checkbox item options.
func (b *CheckboxItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Checked:
		b.state = StateUnchecked
		if v {
			b.state = StateChecked
		}
		return true
	case indeterminate:
		b.state = StateIndeterminate
		return true
	case controlled:
		b.controlled = true
		return true
	}
	return b.accept(c)
}

// Element resolves the checkbox item. The indicator is rendered by the
// component (see Indicator).
func (b *CheckboxItemBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	if b.key == "" {
		b.key = uid.New("menu-checkbox")
	}
	user := withChild(b.Elem(), k.CheckboxIndicator, true)
	shareIndicator(user.ChildNodes(), b.state)
	d := h.Div(b.implAttrs("menuitemcheckbox")...).Children(
		h.Aria("checked", ariaChecked(b.state)),
		h.Data("slot", k.Name+"-checkbox-item"),
		h.Class(k.CheckboxItem),
		h.Data("state", string(b.state)),
		h.Attr(markerCheckbox, b.key),
		h.If(b.controlled, h.Attr(markerControlled, "")),
	)
	return b.resolve(k.Name+" CheckboxItem", spread(d, user))
}

// Render renders the checkbox item.
func (b *CheckboxItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// RadioGroupBuilder is returned by RadioGroup.
type RadioGroupBuilder struct {
	h.Builder[RadioGroupBuilder]
	Scoped
	value      string
	controlled bool
}

// RadioGroup groups radio items; Value selects one. Accepts Controlled.
func RadioGroup(children ...any) *RadioGroupBuilder {
	b := &RadioGroupBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Value and Controlled.
func (b *RadioGroupBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		b.value = string(v)
	case controlled:
		b.controlled = true
	default:
		return false
	}
	return true
}

// Element resolves the radio group.
func (b *RadioGroupBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	key := uid.New("menu-radio")
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *RadioGroupBuilder:
			return false
		case *RadioItemBuilder:
			c.group = b
		}
		return true
	})
	d := h.Div(
		h.Role("group"),
		h.Data("slot", k.Name+"-radio-group"),
		h.Attr(markerRadioGroup, key),
		h.Attr(markerRadioValue, b.value),
		h.If(b.controlled, h.Attr(markerControlled, "")),
	)
	return spread(d, b.Elem())
}

// Render renders the radio group.
func (b *RadioGroupBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// RadioItemBuilder is returned by RadioItem.
type RadioItemBuilder struct {
	h.Builder[RadioItemBuilder]
	itemBase
	value string
	group *RadioGroupBuilder
}

// RadioItem is an item of a RadioGroup. Accepts Value (required),
// Disabled and TextValue.
func RadioItem(children ...any) *RadioItemBuilder {
	b := &RadioItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the radio item options.
func (b *RadioItemBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Value); ok {
		b.value = string(v)
		return true
	}
	return b.accept(c)
}

// Element resolves the radio item.
func (b *RadioItemBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	checked := b.group != nil && b.group.value == b.value
	st := StateUnchecked
	if checked {
		st = StateChecked
	}
	user := withChild(b.Elem(), k.RadioIndicator, true)
	shareIndicator(user.ChildNodes(), st)
	d := h.Div(b.implAttrs("menuitemradio")...).Children(
		h.Aria("checked", fmt.Sprint(checked)),
		h.Data("slot", k.Name+"-radio-item"),
		h.Class(k.RadioItem),
		h.Data("state", string(st)),
		h.Attr(markerRadioValue, b.value),
	)
	return b.resolve(k.Name+" RadioItem", spread(d, user))
}

// Render renders the radio item.
func (b *RadioItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// shareIndicator hands the item's state to the indicators among children
// (not crossing nested items).
func shareIndicator(children []any, st CheckedState) {
	tree.Walk(children, func(n any) bool {
		switch c := n.(type) {
		case *CheckboxItemBuilder, *RadioItemBuilder:
			return false
		case *IndicatorBuilder:
			c.state = st
		}
		return true
	})
}

// IndicatorBuilder is returned by Indicator.
type IndicatorBuilder struct {
	h.Builder[IndicatorBuilder]
	state CheckedState
	force bool
}

// Indicator is Radix's ItemIndicator: rendered while its checkbox or radio
// item is checked (or indeterminate). Accepts ForceMount.
func Indicator(children ...any) *IndicatorBuilder {
	b := &IndicatorBuilder{state: StateUnchecked}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// AcceptChild consumes ForceMount.
func (b *IndicatorBuilder) AcceptChild(c any) bool {
	if v, ok := c.(ForceMount); ok {
		b.force = bool(v)
		return true
	}
	return false
}

// Element resolves the indicator. It is always rendered with a marker;
// menu.js removes it while its item is unchecked (unless forced).
func (b *IndicatorBuilder) Element() *h.Element {
	v := "1"
	if b.force {
		v = "force"
	}
	return spread(h.Span(h.Attr(markerIndicator, v), h.Data("state", string(b.state))), b.Elem())
}

// Render renders the indicator.
func (b *IndicatorBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- submenus ------------------------------------------------------------------------------

// SubScope is the state a Sub shares with its trigger and content.
type SubScope struct {
	ContentID   string
	TriggerID   string
	DefaultOpen bool
}

type subParticipant interface{ setSubScope(*SubScope) }

// SubBuilder is returned by Sub.
type SubBuilder struct {
	h.Builder[SubBuilder]
	Scoped
	defaultOpen bool
}

// SubDefaultOpen opens a submenu whenever its menu opens.
type SubDefaultOpen bool

// Sub is a submenu: a SubTrigger and its SubContent. It renders its
// children only.
func Sub(children ...any) *SubBuilder {
	b := &SubBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes SubDefaultOpen.
func (b *SubBuilder) AcceptChild(c any) bool {
	if v, ok := c.(SubDefaultOpen); ok {
		b.defaultOpen = bool(v)
		return true
	}
	return false
}

// Render renders the children.
func (b *SubBuilder) Render(w io.Writer) error {
	ss := &SubScope{ContentID: uid.New("radix"), TriggerID: uid.New("radix"), DefaultOpen: b.defaultOpen}
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *SubBuilder:
			return false
		case subParticipant:
			c.setSubScope(ss)
		}
		return true
	})
	return h.Group(kids...).Render(w)
}

// SubTriggerBuilder is returned by SubTrigger.
type SubTriggerBuilder struct {
	h.Builder[SubTriggerBuilder]
	itemBase
	sub   *SubScope
	inset bool
}

func (b *SubTriggerBuilder) setSubScope(s *SubScope) { b.sub = s }

// SubTrigger opens its submenu. Accepts Inset, Disabled and TextValue.
func SubTrigger(children ...any) *SubTriggerBuilder {
	b := &SubTriggerBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the sub trigger options.
func (b *SubTriggerBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Inset); ok {
		b.inset = bool(v)
		return true
	}
	return b.accept(c)
}

// Element resolves the sub trigger, with the chevron (icon) after its
// children.
func (b *SubTriggerBuilder) Element() *h.Element {
	k := b.MenuScope().Kind
	if b.sub == nil {
		b.sub = &SubScope{ContentID: uid.New("radix"), TriggerID: uid.New("radix")}
	}
	open := b.sub.DefaultOpen
	d := h.Div(b.implAttrs("menuitem")...).Children(
		h.ID(b.sub.TriggerID),
		h.Aria("haspopup", "menu"),
		h.Aria("expanded", fmt.Sprint(open)),
		h.Aria("controls", b.sub.ContentID),
		h.Data("state", state(open)),
		h.Props(k.SubTriggerAttrs...),
		h.Data("slot", k.Name+"-sub-trigger"),
		h.If(b.inset, h.Data("inset", "true")),
		h.Class(k.SubTrigger),
	)
	return b.resolve(k.Name+" SubTrigger", spread(d, withChild(b.Elem(), k.SubTriggerChevron, false)))
}

// Render renders the sub trigger.
func (b *SubTriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

func state(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}

// SubContentBuilder is returned by SubContent.
type SubContentBuilder struct {
	h.Builder[SubContentBuilder]
	Scoped
	sub    *SubScope
	opts   popper.Options
	loop   bool
	portal bool
}

func (b *SubContentBuilder) setSubScope(s *SubScope) { b.sub = s }

// SubContent is a submenu's content, positioned next to its trigger.
// Rendered in place (inside the parent content) unless wrapped in Portal.
// Accepts popper options (except side and align, which Radix sets) and
// Loop.
func SubContent(children ...any) *SubContentBuilder {
	b := &SubContentBuilder{opts: popper.Defaults()}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes popper options and Loop.
func (b *SubContentBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Loop); ok {
		b.loop = bool(v)
		return true
	}
	return b.opts.Accept(c)
}

// Element resolves the template holding the sub content.
func (b *SubContentBuilder) Element() *h.Element {
	sc := b.MenuScope()
	k := sc.Kind
	if b.sub == nil {
		b.sub = &SubScope{ContentID: uid.New("radix"), TriggerID: uid.New("radix")}
	}
	opts := b.opts
	opts.Align = popper.AlignStart
	opts.Side = popper.SideRight
	if sc.Direction() == "rtl" {
		opts.Side = popper.SideLeft
	}
	opts.Dir = sc.Direction()
	user := b.Elem()
	d := h.Div(
		h.Role("menu"),
		h.Aria("orientation", "vertical"),
		h.Data("state", "open"),
		h.Data("radix-menu-content", ""),
		h.Dir(sc.Direction()),
		h.ID(b.sub.ContentID),
		h.Aria("labelledby", b.sub.TriggerID),
		h.Props(k.ContentAttrs...),
		h.Data("slot", k.Name+"-sub-content"),
		h.Class(k.SubContent),
		h.TabIndex(-1),
		h.Data("orientation", "vertical"),
	)
	el := spread(d, user)
	el.Attr("style", contentStyle(k.Name, user))
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-menu-sub", ""),
		h.Data("content-id", b.sub.ContentID),
		h.Data("trigger-id", b.sub.TriggerID),
		h.If(b.portal, h.Data("portal", "")),
		h.If(b.sub.DefaultOpen, h.Data("default-open", "")),
		h.If(b.loop, h.Data("loop", "")),
		popper.Wrap(opts, el),
	)
}

// Render renders the template.
func (b *SubContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// PortalBuilder is returned by Portal.
type PortalBuilder struct {
	h.Builder[PortalBuilder]
}

// Portal mounts the SubContent among its children in <body> instead of in
// place. (Content is always portalled.) It renders its children only.
func Portal(children ...any) *PortalBuilder {
	b := &PortalBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// Render renders the children.
func (b *PortalBuilder) Render(w io.Writer) error {
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		if c, ok := n.(*SubContentBuilder); ok {
			c.portal = true
			return false
		}
		return true
	})
	return h.Group(kids...).Render(w)
}
