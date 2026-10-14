// Package menubar is the shadcn/ui Menubar: a horizontal bar of menus, as
// in desktop applications (Radix Menubar, @radix-ui/react-menubar 1.1.16,
// on the Menu primitive of package menu).
//
//	menubar.Menubar(
//		menubar.Menu(
//			menubar.Trigger("File"),
//			menubar.Content(
//				menubar.Item("New Tab ", menubar.Shortcut("⌘T")),
//				menubar.Sub(menubar.SubTrigger("Share"), menubar.SubContent(menubar.Item("Email link"))),
//			),
//		),
//	)
//
// Options: Loop(false), Dir("rtl") and DefaultValue on Menubar (the value
// of the menu open on load); Value on Menu (default: a generated id);
// positioning options on Content (default: align start, align offset -4,
// side offset 8); the item options of package menu. Menus are never modal.
//
// Behavior (menubar.js on runtime/js/menu.js): pointer down on a trigger
// opens its menu; moving the pointer over another trigger while a menu is
// open switches to it; Enter/Space toggle and ArrowDown opens from the
// keyboard; ArrowLeft/ArrowRight move between triggers and, in an open
// menu, to the previous/next menu (wrapping unless Loop(false)). The menubar
// dispatches "shadcn:value-change" (cancelable, detail.value: the open
// menu's value, "" when closed).
package menubar

import (
	_ "embed"
	"fmt"
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/menu"
	"shadcn.go/components/ui/popper"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed menubar.js
var script string

func init() { runtime.Module("menubar", script) }

// Marker attributes menubar.js reads and removes from the menubar.
const (
	markerLoop = "data-radix-menubar-loop"
	markerDir  = "data-radix-menubar-dir"
)

var rootClass = lib.Cn("flex h-9 items-center gap-1 rounded-md border bg-background p-1 shadow-xs")
var triggerClass = "flex items-center rounded-sm px-2 py-1 text-sm font-medium outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground"

// Kind holds Menubar's classes and slots.
var Kind = &menu.Kind{
	Name:           "menubar",
	Content:        "z-50 min-w-[12rem] origin-(--radix-menubar-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
	SubContent:     "z-50 min-w-[8rem] origin-(--radix-menubar-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
	Item:           "relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[inset]:pl-8 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!",
	CheckboxItem:   "relative flex cursor-default items-center gap-2 rounded-xs py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	RadioItem:      "relative flex cursor-default items-center gap-2 rounded-xs py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	Label:          "px-2 py-1.5 text-sm font-medium data-[inset]:pl-8",
	Separator:      "-mx-1 my-1 h-px bg-border",
	Shortcut:       "ml-auto text-xs tracking-widest text-muted-foreground",
	SubTrigger:     "flex cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-none select-none focus:bg-accent focus:text-accent-foreground data-[inset]:pl-8 data-[state=open]:bg-accent data-[state=open]:text-accent-foreground",
	SubTriggerIcon: "ml-auto h-4 w-4",
	ContentDefaults: func(o *popper.Options) {
		o.Align = popper.AlignStart
		o.AlignOffset = -4
		o.SideOffset = 8
	},
	ContentAttrs:    []h.Prop{h.Data("radix-menubar-content", "")},
	SubTriggerAttrs: []h.Prop{h.Data("radix-menubar-subtrigger", "")},
	ContentIDs:      true,
}

// Positioning options accepted by Content and SubContent (see package popper).
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

// Part options (see package menu). Value is also a Menu's value.
type (
	Inset          = menu.Inset
	Variant        = menu.Variant
	Disabled       = menu.Disabled
	Checked        = menu.Checked
	Value          = menu.Value
	TextValue      = menu.TextValue
	SubDefaultOpen = menu.SubDefaultOpen
)

// Item variants.
const (
	VariantDefault     = menu.VariantDefault
	VariantDestructive = menu.VariantDestructive
)

var (
	// Indeterminate puts a checkbox item in the indeterminate state.
	Indeterminate = menu.Indeterminate
	// Controlled makes checked state change only through scripts.
	Controlled = menu.Controlled
)

// Loop sets whether keyboard navigation wraps around the triggers and
// menus (default true). On Content or SubContent it is menu.Loop: wrapping
// between items (default false).
type Loop bool

// DefaultValue is the value of the menu open when the page loads.
type DefaultValue string

// Dir sets the reading direction ("ltr" or "rtl").
type Dir string

// --- Menubar -------------------------------------------------------------------------------

// RootBuilder is returned by Menubar.
type RootBuilder struct {
	h.Builder[RootBuilder]
	loop         bool
	dir          string
	defaultValue string
}

// Menubar is the bar holding the menus.
func Menubar(children ...any) *RootBuilder {
	b := &RootBuilder{loop: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Loop, Dir and DefaultValue.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Loop:
		b.loop = bool(v)
	case Dir:
		b.dir = string(v)
	case DefaultValue:
		b.defaultValue = string(v)
	default:
		return false
	}
	return true
}

// Element resolves the menubar.
func (b *RootBuilder) Element() *h.Element {

	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case *MenuBuilder:
			c.bar = b
			return false
		}
		return true
	})
	d := h.Div(
		h.Role("menubar"),
		h.Data("slot", "menubar"),
		h.Class(rootClass),
		h.TabIndex(0),
		h.Data("orientation", "horizontal"),
		h.Style("outline: none;"),
		h.If(!b.loop, h.Attr(markerLoop, "false")),
		h.If(b.dir == "rtl", h.Attr(markerDir, "rtl")),
	)
	return h.Spread(d, b.Elem(), lib.Cn)
}

// Render renders the menubar.
func (b *RootBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Menu ----------------------------------------------------------------------------------

// MenuBuilder is returned by Menu.
type MenuBuilder struct {
	h.Builder[MenuBuilder]
	bar   *RootBuilder
	value string
}

// MenuRoot marks the root of a menu's scope.
func (*MenuBuilder) MenuRoot() {}

// Menu is one menu of the bar: a Trigger and its Content. Accepts Value.
// It renders its children only.
func Menu(children ...any) *MenuBuilder {
	b := &MenuBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Value.
func (b *MenuBuilder) AcceptChild(c any) bool {
	if v, ok := c.(menu.Value); ok {
		b.value = string(v)
		return true
	}
	return false
}

// Render renders the children.
func (b *MenuBuilder) Render(w io.Writer) error {
	sc := menu.NewScope(Kind)
	sc.Modal = false
	value := b.value
	if value == "" {
		value = uid.New("menubar-menu")
	}
	if b.bar != nil {
		sc.Dir = b.bar.dir
		sc.DefaultOpen = b.bar.defaultValue != "" && b.bar.defaultValue == value
	}
	sc.Template = []h.Prop{h.Data("radix-menubar-menu", ""), h.Data("value", value)}
	kids := b.Elem().ChildNodes()
	menu.Share(kids, sc)
	return h.Group(kids...).Render(w)
}

// --- Trigger -------------------------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	menu.Scoped
	asChild  bool
	disabled bool
}

// Trigger is the button that opens its menu. Accepts lib.AsChild and
// Disabled.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and Disabled.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case lib.AsChildFlag:
		b.asChild = true
	case menu.Disabled:
		b.disabled = bool(v)
	default:
		return false
	}
	return true
}

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	sc := b.MenuScope()
	open := sc.DefaultOpen
	d := h.Button(
		h.Type("button"),
		h.Role("menuitem"),
		h.ID(sc.TriggerID),
		h.Aria("haspopup", "menu"),
		h.Aria("expanded", fmt.Sprint(open)),
		h.If(open, h.Aria("controls", sc.ContentID)),
		h.Data("state", map[bool]string{true: "open", false: "closed"}[open]),
		h.If(b.disabled, h.Data("disabled", "")),
		h.If(b.disabled, h.Bool("disabled")),
		h.TabIndex(-1), // defaultValue does not set the roving tab stop
		h.Data("orientation", "horizontal"),
		h.Data("radix-collection-item", ""),
		h.Data("slot", "menubar-trigger"),
		h.Class(triggerClass),
	)
	el := h.Spread(d, b.Elem(), lib.Cn)
	if b.asChild {
		return aschild.Resolve("menubar Trigger", el)
	}
	return el
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- parts ---------------------------------------------------------------------------------

// Portal mounts a SubContent in <body> (Content is always portalled).
func Portal(children ...any) *menu.PortalBuilder { return menu.Portal(children...) }

// Content is a menu, portalled and positioned below its trigger.
func Content(children ...any) *menu.ContentBuilder { return menu.Content(children...) }

// Group groups items.
func Group(children ...any) *menu.GroupBuilder { return menu.Group(children...) }

// Label labels a group of items. Accepts Inset.
func Label(children ...any) *menu.LabelBuilder { return menu.Label(children...) }

// Item is a menu item. Accepts Inset, VariantDestructive, Disabled,
// TextValue and lib.AsChild.
func Item(children ...any) *menu.ItemBuilder { return menu.Item(children...) }

// CheckboxItem is an item with a check mark that toggles. Accepts
// Checked, Indeterminate, Controlled, Disabled and TextValue.
func CheckboxItem(children ...any) *menu.CheckboxItemBuilder {
	return menu.CheckboxItem(append([]any{indicator(lucide.Check(h.Class("size-4")))}, children...)...)
}

// RadioGroup groups radio items. Accepts Value and Controlled.
func RadioGroup(children ...any) *menu.RadioGroupBuilder { return menu.RadioGroup(children...) }

// RadioItem is an item of a RadioGroup. Accepts Value, Disabled and
// TextValue.
func RadioItem(children ...any) *menu.RadioItemBuilder {
	return menu.RadioItem(append([]any{indicator(lucide.Circle(h.Class("size-2 fill-current")))}, children...)...)
}

func indicator(icon h.Node) h.Node {
	return h.Span(h.Class("pointer-events-none absolute left-2 flex size-3.5 items-center justify-center"), menu.Indicator(icon))
}

// Separator separates items.
func Separator(children ...any) *menu.SeparatorBuilder { return menu.Separator(children...) }

// Shortcut shows a keyboard shortcut in an item.
func Shortcut(children ...any) *menu.ShortcutBuilder { return menu.Shortcut(children...) }

// Sub is a submenu. Accepts SubDefaultOpen.
func Sub(children ...any) *menu.SubBuilder { return menu.Sub(children...) }

// SubTrigger opens its submenu. Accepts Inset and Disabled.
func SubTrigger(children ...any) *menu.SubTriggerBuilder {
	return menu.SubTrigger(append(children, lucide.ChevronRight(h.Class(Kind.SubTriggerIcon)))...)
}

// SubContent is a submenu's content, rendered in place unless wrapped in
// Portal.
func SubContent(children ...any) *menu.SubContentBuilder { return menu.SubContent(children...) }
