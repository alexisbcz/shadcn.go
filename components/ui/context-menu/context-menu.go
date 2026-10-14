// Package contextmenu is the shadcn/ui ContextMenu: a menu opened at the
// pointer by a right click (or a long press on touch) in an area (Radix
// ContextMenu, @radix-ui/react-context-menu 2.2.16, on the Menu primitive
// of package menu).
//
//	contextmenu.ContextMenu(
//		contextmenu.Trigger(h.Class("..."), "Right click here"),
//		contextmenu.Content(h.Class("w-52"),
//			contextmenu.Item(contextmenu.Inset(true), "Back", contextmenu.Shortcut("⌘[")),
//			contextmenu.Sub(
//				contextmenu.SubTrigger(contextmenu.Inset(true), "More Tools"),
//				contextmenu.SubContent(contextmenu.Item("Save Page...")),
//			),
//			contextmenu.CheckboxItem(contextmenu.Checked(true), "Show Bookmarks"),
//		),
//	)
//
// Options: Modal(false) and Dir("rtl") on ContextMenu; Disabled on Trigger;
// the item options of package menu. The content always opens on the right
// of the pointer, aligned to its start (Radix forces side, offset and
// alignment); positioning options control collisions only.
//
// Behavior (context-menu.js on runtime/js/menu.js): see package
// dropdownmenu for the menu itself. The content template dispatches
// "shadcn:open-change" (cancelable, detail.open).
package contextmenu

import (
	_ "embed"
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/menu"
	"shadcn.go/components/ui/popper"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed context-menu.js
var script string

func init() { runtime.Module("context-menu", script) }

const markerTrigger = "data-radix-context-menu-trigger"

const subTriggerIcon = "ml-auto"

// Kind holds ContextMenu's classes and slots.
var Kind = &menu.Kind{
	Name:              "context-menu",
	Content:           "z-50 max-h-(--radix-context-menu-content-available-height) min-w-[8rem] origin-(--radix-context-menu-content-transform-origin) overflow-x-hidden overflow-y-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
	SubContent:        "z-50 min-w-[8rem] origin-(--radix-context-menu-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
	Item:              "relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[inset]:pl-8 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!",
	CheckboxItem:      "relative flex cursor-default items-center gap-2 rounded-sm py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	RadioItem:         "relative flex cursor-default items-center gap-2 rounded-sm py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	Label:             "px-2 py-1.5 text-sm font-medium text-foreground data-[inset]:pl-8",
	Separator:         "-mx-1 my-1 h-px bg-border",
	Shortcut:          "ml-auto text-xs tracking-widest text-muted-foreground",
	SubTrigger:        "flex cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[inset]:pl-8 data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground",
	CheckboxIndicator: func() h.Node { return indicator(lucide.Check(h.Class("size-4"))) },
	RadioIndicator:    func() h.Node { return indicator(lucide.Circle(h.Class("size-2 fill-current"))) },
	SubTriggerChevron: func() h.Node { return lucide.ChevronRight(h.Class(subTriggerIcon)) },
	Styled:            styledKind,
	SubTriggerIcon:    subTriggerIcon,
	ContentForced: func(o *popper.Options) {
		o.Side = popper.SideRight
		o.SideOffset = 2
		o.Align = popper.AlignStart
	},
}

// Positioning options accepted by Content and SubContent (see package popper).
type (
	AlignOffset      = popper.AlignOffset
	AvoidCollisions  = popper.AvoidCollisions
	CollisionPadding = popper.CollisionPadding
	Sticky           = popper.Sticky
	HideWhenDetached = popper.HideWhenDetached
)

// Part options (see package menu).
type (
	Inset          = menu.Inset
	Variant        = menu.Variant
	Disabled       = menu.Disabled
	Checked        = menu.Checked
	Value          = menu.Value
	TextValue      = menu.TextValue
	Loop           = menu.Loop
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

// --- root ----------------------------------------------------------------------------------

// Option configures a ContextMenu; pass options as children.
type Option func(*menu.Scope)

// Modal sets whether the menu is modal (default true).
func Modal(on bool) Option { return func(s *menu.Scope) { s.Modal = on } }

// Dir sets the reading direction ("ltr" or "rtl").
func Dir(dir string) Option { return func(s *menu.Scope) { s.Dir = dir } }

// RootBuilder is returned by ContextMenu.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts []Option
}

// MenuRoot marks the root (scopes do not cross it).
func (*RootBuilder) MenuRoot() {}

// ContextMenu is the root. It renders its children only.
func ContextMenu(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Options.
func (b *RootBuilder) AcceptChild(c any) bool {
	if o, ok := c.(Option); ok {
		b.opts = append(b.opts, o)
		return true
	}
	return false
}

// Render renders the children.
func (b *RootBuilder) Render(w io.Writer) error {
	sc := menu.NewScope(Kind.For(w))
	sc.Template = []h.Prop{h.Data("radix-context-menu", "")}
	for _, o := range b.opts {
		o(sc)
	}
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

// Trigger is the area that opens the menu (a <span>). Accepts lib.AsChild
// and Disabled.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "span")
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

// Element resolves the trigger. (Radix's style -webkit-touch-callout: none
// is not rendered: browsers without it drop the declaration, and React
// renders no style attribute then.)
func (b *TriggerBuilder) Element() *h.Element { return b.element("") }

func (b *TriggerBuilder) element(class string) *h.Element {
	sc := b.MenuScope()
	d := h.Span(
		h.If(class != "", h.Class(class)),
		h.Data("state", "closed"),
		h.If(b.disabled, h.Data("disabled", "")),
		h.Data("slot", "context-menu-trigger"),
		h.Attr(markerTrigger, sc.ContentID),
	)
	el := h.Spread(d, b.Elem(), lib.Cn)
	if b.asChild {
		return aschild.Resolve("context-menu Trigger", el)
	}
	return el
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error {
	class := ""
	if st, ok := style.FromWriter(w); ok {
		class = baseTrigger.Class(st, "")
	}
	return b.element(class).Render(w)
}

// --- parts ---------------------------------------------------------------------------------

// Portal mounts a SubContent in <body> (Content is always portalled).
func Portal(children ...any) *menu.PortalBuilder { return menu.Portal(children...) }

// Content is the menu, portalled and positioned at the pointer.
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
	return menu.CheckboxItem(children...)
}

// RadioGroup groups radio items. Accepts Value and Controlled.
func RadioGroup(children ...any) *menu.RadioGroupBuilder { return menu.RadioGroup(children...) }

// RadioItem is an item of a RadioGroup. Accepts Value, Disabled and
// TextValue.
func RadioItem(children ...any) *menu.RadioItemBuilder {
	return menu.RadioItem(children...)
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
	return menu.SubTrigger(children...)
}

// SubContent is a submenu's content, rendered in place unless wrapped in
// Portal.
func SubContent(children ...any) *menu.SubContentBuilder { return menu.SubContent(children...) }
