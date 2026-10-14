// Package sidebar is the shadcn/ui Sidebar: a composable, themeable and
// customizable sidebar component.
//
//	sidebar.Provider(
//		sidebar.Sidebar(sidebar.CollapsibleIcon,
//			sidebar.Header(...),
//			sidebar.Content(
//				sidebar.Group(
//					sidebar.GroupLabel("Application"),
//					sidebar.GroupContent(
//						sidebar.Menu(
//							sidebar.MenuItem(
//								sidebar.MenuButton(lib.AsChild, sidebar.Tooltip("Home"),
//									h.A(h.Href("/"), lucide.House(), h.Span("Home")),
//								),
//							),
//						),
//					),
//				),
//			),
//			sidebar.Footer(...),
//			sidebar.Rail(),
//		),
//		sidebar.Inset(h.Header(sidebar.Trigger()), ...),
//	)
//
// Provider holds the sidebar's state, like SidebarProvider: expanded or
// collapsed on desktop (DefaultOpen, persisted in the "sidebar_state"
// cookie on every change; read it back with CookieState), and the mobile
// sheet's open state below the md breakpoint (768px), where every Sidebar
// that can collapse is shown in a Sheet instead. Cmd+B / Ctrl+B, Trigger
// and Rail toggle the sidebar.
//
// Behavior (sidebar.js). Scripts can observe and drive a provider:
//
//	wrapper.addEventListener("shadcn:open-change", (e) => e.target === wrapper && e.detail.open) // cancelable
//	import { getSidebar } from "/_shadcn/sidebar.js"
//	const s = getSidebar(wrapperOrAnyElementInside)
//	s.state; s.open; s.setOpen(false); s.openMobile; s.setOpenMobile(true); s.isMobile; s.toggleSidebar()
//
// Cancelling "shadcn:open-change" keeps the state, like a controlled
// SidebarProvider whose parent ignores onOpenChange (the cookie is still
// written, as upstream does).
package sidebar

import (
	_ "embed"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sheet"
	"shadcn.go/components/ui/skeleton"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

//go:embed sidebar.js
var script string

func init() { runtime.Module("sidebar", script) }

// Upstream constants.
const (
	CookieName        = "sidebar_state"
	CookieMaxAge      = 60 * 60 * 24 * 7
	Width             = "16rem"
	WidthMobile       = "18rem"
	WidthIcon         = "3rem"
	KeyboardShortcut  = "b"
	MobileBreakpoint  = 768
	cookieValueOpen   = "true"
	tooltipTemplateAt = "data-sidebar-tooltip"
)

// CookieState reads the sidebar state persisted by the client in the
// "sidebar_state" cookie. ok is false when the cookie is absent. Upstream
// apps pass DefaultOpen(open) (absent meaning collapsed in the shadcn docs:
// `cookies().get("sidebar_state")?.value === "true"`).
func CookieState(r *http.Request) (open, ok bool) {
	c, err := r.Cookie(CookieName)
	if err != nil {
		return false, false
	}
	return c.Value == cookieValueOpen, true
}

// --- context ---------------------------------------------------------------------------------

type context struct {
	open bool
}

var defaultContext = &context{open: true}

type contextual struct{ ctx *context }

func (c *contextual) context() *context {
	if c.ctx == nil {
		return defaultContext
	}
	return c.ctx
}

func (c *contextual) setContext(ctx *context) { c.ctx = ctx }

// spreadProp returns a prop spreading user's attributes (not its children)
// over the element it is applied to, classes merged with cn.
func spreadProp(user *h.Element) h.Prop {
	u := user.Clone().SetChildren()
	return func(e *h.Element) { *e = *h.Spread(e, u, lib.Cn) }
}

// asChild resolves props rendered with lib.AsChild (Radix Slot): props are
// merged into the single child, the child's props winning and classes
// joined.
func asChild(part string, props *h.Element) *h.Element {
	return aschild.Resolve("sidebar "+part, props)
}

// --- Provider --------------------------------------------------------------------------------

// DefaultOpen, passed to Provider, sets whether the sidebar starts
// expanded (default true).
type DefaultOpen bool

// ProviderBuilder is returned by Provider.
type ProviderBuilder struct {
	h.Builder[ProviderBuilder]
	open bool
}

// Provider is SidebarProvider: the wrapper <div> holding the sidebar's
// state for every part inside it. Accepts DefaultOpen.
func Provider(children ...any) *ProviderBuilder {
	b := &ProviderBuilder{open: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes DefaultOpen.
func (b *ProviderBuilder) AcceptChild(c any) bool {
	if v, ok := c.(DefaultOpen); ok {
		b.open = bool(v)
		return true
	}
	return false
}

// DefaultOpen sets whether the sidebar starts expanded.
func (b *ProviderBuilder) DefaultOpen(on bool) *ProviderBuilder { b.open = on; return b }

const wrapperClass = "group/sidebar-wrapper flex min-h-svh w-full has-data-[variant=inset]:bg-sidebar"

// Element resolves the provider to its wrapper <div>.
func (b *ProviderBuilder) Element() *h.Element {
	ctx := &context{open: b.open}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ProviderBuilder:
			return false
		case interface{ setContext(*context) }:
			c.setContext(ctx)
		}
		return true
	})
	user := b.Elem()
	style := "--sidebar-width: " + Width + "; --sidebar-width-icon: " + WidthIcon
	if us, ok := user.GetAttr("style"); ok && us != "" {
		style += "; " + us
	}
	base := h.Div(
		h.Data("slot", "sidebar-wrapper"),
		h.Class(lib.Cn(wrapperClass)),
		h.Data("sidebar-default-open", strconv.FormatBool(b.open)),
	)
	out := h.Spread(base, user, lib.Cn)
	out.Attr("style", style)
	return out
}

// Render renders the provider.
func (b *ProviderBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Sidebar ---------------------------------------------------------------------------------

// Side is the edge the sidebar is on.
type Side string

// Sides.
const (
	SideLeft  Side = "left"
	SideRight Side = "right"
)

// Variant is the sidebar's look.
type Variant string

// Variants.
const (
	VariantSidebar  Variant = "sidebar"
	VariantFloating Variant = "floating"
	VariantInset    Variant = "inset"
)

// Collapsible is how the sidebar collapses on desktop.
type Collapsible string

// Collapsible modes.
const (
	CollapsibleOffcanvas Collapsible = "offcanvas"
	CollapsibleIcon      Collapsible = "icon"
	CollapsibleNone      Collapsible = "none"
)

// Builder is returned by Sidebar.
type Builder struct {
	h.Builder[Builder]
	contextual
	side        Side
	variant     Variant
	collapsible Collapsible
}

// Sidebar is the sidebar itself. Accepts a Side (default SideLeft), a
// Variant (default VariantSidebar) and a Collapsible (default
// CollapsibleOffcanvas); other attributes go to the sidebar container,
// like upstream's props.
func Sidebar(children ...any) *Builder {
	b := &Builder{side: SideLeft, variant: VariantSidebar, collapsible: CollapsibleOffcanvas}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Side, Variant and Collapsible.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Side:
		b.side = v
	case Variant:
		b.variant = v
	case Collapsible:
		b.collapsible = v
	default:
		return false
	}
	return true
}

// Side sets the side.
func (b *Builder) Side(s Side) *Builder { b.side = s; return b }

// Variant sets the variant.
func (b *Builder) Variant(v Variant) *Builder { b.variant = v; return b }

// Collapsible sets the collapsible mode.
func (b *Builder) Collapsible(c Collapsible) *Builder { b.collapsible = c; return b }

const (
	noneClass   = "flex h-full w-(--sidebar-width) flex-col bg-sidebar text-sidebar-foreground"
	rootClass   = "group peer hidden text-sidebar-foreground md:block"
	innerClass  = "flex h-full w-full flex-col bg-sidebar group-data-[variant=floating]:rounded-lg group-data-[variant=floating]:border group-data-[variant=floating]:border-sidebar-border group-data-[variant=floating]:shadow-sm"
	mobileClass = "w-(--sidebar-width) bg-sidebar p-0 text-sidebar-foreground [&>button]:hidden"
)

// Element resolves the sidebar.
func (b *Builder) Element() *h.Element {
	user := b.Elem()
	if b.collapsible == CollapsibleNone {
		base := h.Div(h.Data("slot", "sidebar"), h.Class(lib.Cn(noneClass)))
		return h.Spread(base, user, lib.Cn)
	}
	ctx := b.context()
	state := "expanded"
	collapsibleAttr := ""
	if !ctx.open {
		state = "collapsed"
		collapsibleAttr = string(b.collapsible)
	}
	floating := b.variant == VariantFloating || b.variant == VariantInset

	gap := []string{
		"relative w-(--sidebar-width) bg-transparent transition-[width] duration-200 ease-linear",
		"group-data-[collapsible=offcanvas]:w-0",
		"group-data-[side=right]:rotate-180",
	}
	if floating {
		gap = append(gap, "group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4)))]")
	} else {
		gap = append(gap, "group-data-[collapsible=icon]:w-(--sidebar-width-icon)")
	}
	container := []string{"fixed inset-y-0 z-10 hidden h-svh w-(--sidebar-width) transition-[left,right,width] duration-200 ease-linear md:flex"}
	if b.side == SideLeft {
		container = append(container, "left-0 group-data-[collapsible=offcanvas]:left-[calc(var(--sidebar-width)*-1)]")
	} else {
		container = append(container, "right-0 group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)]")
	}
	if floating {
		container = append(container, "p-2 group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4))+2px)]")
	} else {
		container = append(container, "group-data-[collapsible=icon]:w-(--sidebar-width-icon) group-data-[side=left]:border-r group-data-[side=right]:border-l")
	}
	container = append(container, user.ClassName())

	props := user.Clone().SetChildren()
	props.SetClass(lib.Cn(container...))
	cont := h.Spread(h.Div(h.Data("slot", "sidebar-container")), props, nil)
	cont.SetClass(lib.Cn(container...))
	cont.Children(h.Div(
		h.Data("sidebar", "sidebar"),
		h.Data("slot", "sidebar-inner"),
		h.Class(innerClass),
		h.Group(user.ChildNodes()...),
	))

	return h.Div(
		h.Class(rootClass),
		h.Data("state", state),
		h.Data("collapsible", collapsibleAttr),
		h.Data("variant", string(b.variant)),
		h.Data("side", string(b.side)),
		h.Data("slot", "sidebar"),
		h.Div(h.Data("slot", "sidebar-gap"), h.Class(lib.Cn(gap...))),
		cont,
		b.mobileTemplate(),
	)
}

// mobileTemplate is the Sheet the sidebar becomes below the md breakpoint:
// sidebar.js moves the sidebar's children into its placeholder.
func (b *Builder) mobileTemplate() *h.Element {
	side := sheet.SideLeft
	if b.side == SideRight {
		side = sheet.SideRight
	}
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("sidebar-mobile", ""),
		h.Data("sidebar-collapsible", string(b.collapsible)),
		sheet.Sheet(
			sheet.Content(side,
				h.Data("sidebar", "sidebar"),
				h.Data("slot", "sidebar"),
				h.Data("mobile", "true"),
				h.Class(mobileClass),
				h.Style("--sidebar-width: "+WidthMobile),
				sheet.Header(h.Class("sr-only"),
					sheet.Title("Sidebar"),
					sheet.Description("Displays the mobile sidebar."),
				),
				h.Div(h.Class("flex h-full w-full flex-col"), h.Data("sidebar-mobile-body", "")),
			),
		),
	)
}

// Render renders the sidebar.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Trigger, Rail, Inset ----------------------------------------------------------------------

// Trigger is the button toggling the sidebar (a ghost icon Button).
func Trigger(children ...any) *h.Element {
	user := h.Div(children...)
	kids := user.ChildNodes()
	if len(kids) == 0 {
		kids = []any{lucide.PanelLeft(), h.Span(h.Class("sr-only"), "Toggle Sidebar")}
	}
	return button.Button(button.VariantGhost, button.SizeIcon,
		h.Data("sidebar", "trigger"),
		h.Data("slot", "sidebar-trigger"),
		spreadProp(withClass(user, lib.Cn("size-7", user.ClassName()))),
		h.Group(kids...),
	).Element()
}

func withClass(el *h.Element, class string) *h.Element {
	out := el.Clone()
	out.SetClass(class)
	return out
}

const railClass = "absolute inset-y-0 z-20 hidden w-4 -translate-x-1/2 transition-all ease-linear group-data-[side=left]:-right-4 group-data-[side=right]:left-0 after:absolute after:inset-y-0 after:left-1/2 after:w-[2px] hover:after:bg-sidebar-border sm:flex " +
	"in-data-[side=left]:cursor-w-resize in-data-[side=right]:cursor-e-resize " +
	"[[data-side=left][data-state=collapsed]_&]:cursor-e-resize [[data-side=right][data-state=collapsed]_&]:cursor-w-resize " +
	"group-data-[collapsible=offcanvas]:translate-x-0 group-data-[collapsible=offcanvas]:after:left-full hover:group-data-[collapsible=offcanvas]:bg-sidebar " +
	"[[data-side=left][data-collapsible=offcanvas]_&]:-right-2 " +
	"[[data-side=right][data-collapsible=offcanvas]_&]:-left-2"

// Rail is the thin hit area along the sidebar's edge toggling it.
func Rail(children ...any) *h.Element {
	base := h.Button(
		h.Data("sidebar", "rail"),
		h.Data("slot", "sidebar-rail"),
		h.Aria("label", "Toggle Sidebar"),
		h.TabIndex(-1),
		h.Attr("title", "Toggle Sidebar"),
		h.Class(lib.Cn(railClass)),
	)
	return h.Spread(base, h.Button(children...), lib.Cn)
}

const insetClass = "relative flex w-full flex-1 flex-col bg-background " +
	"md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2"

// Inset is the <main> next to the sidebar (styled for VariantInset).
func Inset(children ...any) *h.Element {
	return h.Spread(h.Main(h.Data("slot", "sidebar-inset"), h.Class(lib.Cn(insetClass))), h.Main(children...), lib.Cn)
}

// Input is an Input styled for the sidebar.
func Input(children ...any) *h.Element {
	user := h.Input(children...)
	return input.Input(
		h.Data("slot", "sidebar-input"),
		h.Data("sidebar", "input"),
		spreadProp(withClass(user, lib.Cn("h-8 w-full bg-background shadow-none", user.ClassName()))),
	)
}

func simple(el *h.Element, slot, class string, children []any) *h.Element {
	base := el.Clone().Children(h.Data("slot", slot), h.Data("sidebar", slot[len("sidebar-"):]), h.Class(lib.Cn(class)))
	return h.Spread(base, el.Clone().Children(children...), lib.Cn)
}

// Header is the sticky top of the sidebar.
func Header(children ...any) *h.Element {
	return simple(h.Div(), "sidebar-header", "flex flex-col gap-2 p-2", children)
}

// Footer is the sticky bottom of the sidebar.
func Footer(children ...any) *h.Element {
	return simple(h.Div(), "sidebar-footer", "flex flex-col gap-2 p-2", children)
}

// Separator is a Separator styled for the sidebar.
func Separator(children ...any) *separator.Builder {
	user := h.Div(children...)
	return separator.Separator(
		h.Data("slot", "sidebar-separator"),
		h.Data("sidebar", "separator"),
		spreadProp(withClass(user, lib.Cn("mx-2 w-auto bg-sidebar-border", user.ClassName()))),
		h.Group(user.ChildNodes()...),
	)
}

// Content is the scrollable middle of the sidebar.
func Content(children ...any) *h.Element {
	return simple(h.Div(), "sidebar-content", "flex min-h-0 flex-1 flex-col gap-2 overflow-auto group-data-[collapsible=icon]:overflow-hidden", children)
}

// Group is a section of the sidebar.
func Group(children ...any) *h.Element {
	return simple(h.Div(), "sidebar-group", "relative flex w-full min-w-0 flex-col p-2", children)
}

// GroupContent holds a group's items.
func GroupContent(children ...any) *h.Element {
	return simple(h.Div(), "sidebar-group-content", "w-full text-sm", children)
}

// Menu is the list of menu items.
func Menu(children ...any) *h.Element {
	return simple(h.Ul(), "sidebar-menu", "flex w-full min-w-0 flex-col gap-1", children)
}

// MenuItem is a menu entry.
func MenuItem(children ...any) *h.Element {
	return simple(h.Li(), "sidebar-menu-item", "group/menu-item relative", children)
}

// MenuBadge is a badge next to a menu button.
func MenuBadge(children ...any) *h.Element {
	return simple(h.Div(), "sidebar-menu-badge",
		"pointer-events-none absolute right-1 flex h-5 min-w-5 items-center justify-center rounded-md px-1 text-xs font-medium text-sidebar-foreground tabular-nums select-none "+
			"peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[active=true]/menu-button:text-sidebar-accent-foreground "+
			"peer-data-[size=sm]/menu-button:top-1 "+
			"peer-data-[size=default]/menu-button:top-1.5 "+
			"peer-data-[size=lg]/menu-button:top-2.5 "+
			"group-data-[collapsible=icon]:hidden", children)
}

// MenuSub is a nested list of sub-items.
func MenuSub(children ...any) *h.Element {
	return simple(h.Ul(), "sidebar-menu-sub",
		"mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-l border-sidebar-border px-2.5 py-0.5 group-data-[collapsible=icon]:hidden", children)
}

// MenuSubItem is a sub-menu entry.
func MenuSubItem(children ...any) *h.Element {
	return simple(h.Li(), "sidebar-menu-sub-item", "group/menu-sub-item relative", children)
}

// --- slottable parts ---------------------------------------------------------------------------

// slottable is a part rendering an element of its own or, with
// lib.AsChild, merging into its single child.
type slottable struct {
	asChild bool
}

func (s *slottable) accept(c any) bool {
	if _, ok := c.(lib.AsChildFlag); ok {
		s.asChild = true
		return true
	}
	return false
}

func (s *slottable) resolve(part string, props *h.Element) *h.Element {
	if s.asChild {
		return asChild(part, props)
	}
	return props
}

// GroupLabelBuilder is returned by GroupLabel.
type GroupLabelBuilder struct {
	h.Builder[GroupLabelBuilder]
	slottable
}

// GroupLabel is a group's title. Accepts lib.AsChild.
func GroupLabel(children ...any) *GroupLabelBuilder {
	b := &GroupLabelBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild.
func (b *GroupLabelBuilder) AcceptChild(c any) bool { return b.accept(c) }

// AsChild merges the label into its single child.
func (b *GroupLabelBuilder) AsChild() *GroupLabelBuilder { b.asChild = true; return b }

// Element resolves the label.
func (b *GroupLabelBuilder) Element() *h.Element {
	base := h.Div(h.Data("slot", "sidebar-group-label"), h.Data("sidebar", "group-label"), h.Class(lib.Cn(
		"flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium text-sidebar-foreground/70 ring-sidebar-ring outline-hidden transition-[margin,opacity] duration-200 ease-linear focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0",
		"group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0",
	)))
	return b.resolve("GroupLabel", h.Spread(base, b.Elem(), lib.Cn))
}

// Render renders the label.
func (b *GroupLabelBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// GroupActionBuilder is returned by GroupAction.
type GroupActionBuilder struct {
	h.Builder[GroupActionBuilder]
	slottable
}

// GroupAction is a button in a group's header. Accepts lib.AsChild.
func GroupAction(children ...any) *GroupActionBuilder {
	b := &GroupActionBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild.
func (b *GroupActionBuilder) AcceptChild(c any) bool { return b.accept(c) }

// AsChild merges the action into its single child.
func (b *GroupActionBuilder) AsChild() *GroupActionBuilder { b.asChild = true; return b }

// Element resolves the action.
func (b *GroupActionBuilder) Element() *h.Element {
	base := h.Button(h.Data("slot", "sidebar-group-action"), h.Data("sidebar", "group-action"), h.Class(lib.Cn(
		"absolute top-3.5 right-3 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0",
		"after:absolute after:-inset-2 md:after:hidden",
		"group-data-[collapsible=icon]:hidden",
	)))
	return b.resolve("GroupAction", h.Spread(base, b.Elem(), lib.Cn))
}

// Render renders the action.
func (b *GroupActionBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- MenuButton --------------------------------------------------------------------------------

// IsActive marks a MenuButton or MenuSubButton as the current page.
type IsActive bool

// ButtonVariant is a MenuButton's look.
type ButtonVariant string

// MenuButton variants.
const (
	ButtonVariantDefault ButtonVariant = "default"
	ButtonVariantOutline ButtonVariant = "outline"
)

// ButtonSize is a MenuButton's size.
type ButtonSize string

// MenuButton sizes.
const (
	ButtonSizeDefault ButtonSize = "default"
	ButtonSizeSm      ButtonSize = "sm"
	ButtonSizeLg      ButtonSize = "lg"
)

var menuButtonVariants = cva.New(
	"peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left text-sm ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default": "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
				"outline": "bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)]",
			},
			"size": {
				"default": "h-8 text-sm",
				"sm":      "h-7 text-xs",
				"lg":      "h-12 text-sm group-data-[collapsible=icon]:p-0!",
			},
		},
		Defaults: cva.Selection{"variant": "default", "size": "default"},
		Order:    []string{"variant", "size"},
	},
)

// MenuButtonVariants returns a menu button's classes plus className, like
// sidebarMenuButtonVariants() upstream. Empty values select the defaults.
func MenuButtonVariants(variant ButtonVariant, size ButtonSize, className ...string) string {
	return lib.Cn(menuButtonVariants.Classes(cva.Selection{"variant": string(variant), "size": string(size)}, className...))
}

// TooltipOption, passed to MenuButton, shows a tooltip on the button's
// right while the sidebar is collapsed (not on mobile). Its children are
// the tooltip content's children and props (see tooltip.Content).
type TooltipOption struct{ children []any }

// Tooltip returns a MenuButton tooltip:
//
//	sidebar.MenuButton(sidebar.Tooltip("Home"), ...)
//	sidebar.MenuButton(sidebar.Tooltip(tooltip.SideBottom, h.Class("..."), "Home"), ...)
func Tooltip(children ...any) TooltipOption { return TooltipOption{children: children} }

// TooltipHidden, passed to Tooltip, overrides when the tooltip is hidden
// (upstream's tooltip={{ hidden: false }}): TooltipHidden(false) shows it
// whatever the sidebar's state.
type TooltipHidden bool

// MenuButtonBuilder is returned by MenuButton.
type MenuButtonBuilder struct {
	h.Builder[MenuButtonBuilder]
	contextual
	slottable
	active  bool
	variant ButtonVariant
	size    ButtonSize
	tooltip *TooltipOption
}

// MenuButton is a menu item's button. Accepts lib.AsChild, IsActive, a
// ButtonVariant, a ButtonSize and a Tooltip.
func MenuButton(children ...any) *MenuButtonBuilder {
	b := &MenuButtonBuilder{variant: ButtonVariantDefault, size: ButtonSizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild, IsActive, ButtonVariant, ButtonSize and
// TooltipOption.
func (b *MenuButtonBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case IsActive:
		b.active = bool(v)
	case ButtonVariant:
		b.variant = v
	case ButtonSize:
		b.size = v
	case TooltipOption:
		b.tooltip = &v
	default:
		return b.accept(c)
	}
	return true
}

// AsChild merges the button into its single child.
func (b *MenuButtonBuilder) AsChild() *MenuButtonBuilder { b.asChild = true; return b }

// IsActive marks the button as the current page.
func (b *MenuButtonBuilder) IsActive(on bool) *MenuButtonBuilder { b.active = on; return b }

func (b *MenuButtonBuilder) button() *h.Element {
	v, s := b.variant, b.size
	if v == "" {
		v = ButtonVariantDefault
	}
	if s == "" {
		s = ButtonSizeDefault
	}
	base := h.Button(
		h.Data("slot", "sidebar-menu-button"),
		h.Data("sidebar", "menu-button"),
		h.Data("size", string(s)),
		h.Data("active", strconv.FormatBool(b.active)),
	)
	props := h.Spread(base, b.Elem(), nil)
	props.SetClass(MenuButtonVariants(v, s, b.Elem().ClassName()))
	return b.resolve("MenuButton", props)
}

// Element resolves the button. With a Tooltip, the tooltip's inert
// template is placed first inside the button (Render places it after).
func (b *MenuButtonBuilder) Element() *h.Element {
	btn := b.button()
	if b.tooltip == nil {
		return btn
	}
	var trig, tpl *h.Element
	b.withTooltip(btn, func(t, c *h.Element) { trig, tpl = t, c })
	out := trig.Clone().SetChildren(append([]any{tpl}, trig.ChildNodes()...)...)
	return out
}

// withTooltip renders the tooltip root sharing one scope between the
// trigger and the content, handing both resolved elements to fn.
func (b *MenuButtonBuilder) withTooltip(btn *h.Element, fn func(trigger, content *h.Element)) {
	ctx := b.context()
	hidden := ctx.open
	fixed := false
	var kids []any
	for _, c := range b.tooltip.children {
		if v, ok := c.(TooltipHidden); ok {
			hidden, fixed = bool(v), true
			continue
		}
		kids = append(kids, c)
	}
	trigger := tooltip.Trigger(lib.AsChild, btn)
	content := tooltip.Content(append([]any{
		tooltip.SideRight,
		tooltip.AlignCenter,
		h.If(hidden, h.Bool("hidden")),
	}, kids...)...)
	root := tooltip.Tooltip(trigger, content)
	// Rendering the root hands the scope to its parts.
	_ = root.Render(io.Discard)
	tpl := content.Element()
	if !fixed {
		tpl.Attr(tooltipTemplateAt, "")
	}
	fn(trigger.Element(), tpl)
}

// Render renders the button, followed by its tooltip when it has one.
func (b *MenuButtonBuilder) Render(w io.Writer) error {
	btn := b.button()
	if b.tooltip == nil {
		return btn.Render(w)
	}
	var out h.Node
	b.withTooltip(btn, func(t, c *h.Element) { out = h.Group(t, c) })
	return out.Render(w)
}

// --- MenuAction --------------------------------------------------------------------------------

// ShowOnHover, passed to MenuAction, shows the action only while its menu
// item is hovered or focused (on desktop).
type ShowOnHover bool

// MenuActionBuilder is returned by MenuAction.
type MenuActionBuilder struct {
	h.Builder[MenuActionBuilder]
	slottable
	showOnHover bool
}

// MenuAction is a button next to a menu button. Accepts lib.AsChild and
// ShowOnHover.
func MenuAction(children ...any) *MenuActionBuilder {
	b := &MenuActionBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and ShowOnHover.
func (b *MenuActionBuilder) AcceptChild(c any) bool {
	if v, ok := c.(ShowOnHover); ok {
		b.showOnHover = bool(v)
		return true
	}
	return b.accept(c)
}

// AsChild merges the action into its single child.
func (b *MenuActionBuilder) AsChild() *MenuActionBuilder { b.asChild = true; return b }

// Element resolves the action.
func (b *MenuActionBuilder) Element() *h.Element {
	classes := []string{
		"absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform peer-hover/menu-button:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0",
		"after:absolute after:-inset-2 md:after:hidden",
		"peer-data-[size=sm]/menu-button:top-1",
		"peer-data-[size=default]/menu-button:top-1.5",
		"peer-data-[size=lg]/menu-button:top-2.5",
		"group-data-[collapsible=icon]:hidden",
	}
	if b.showOnHover {
		classes = append(classes, "group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 peer-data-[active=true]/menu-button:text-sidebar-accent-foreground data-[state=open]:opacity-100 md:opacity-0")
	}
	base := h.Button(h.Data("slot", "sidebar-menu-action"), h.Data("sidebar", "menu-action"), h.Class(lib.Cn(classes...)))
	return b.resolve("MenuAction", h.Spread(base, b.Elem(), lib.Cn))
}

// Render renders the action.
func (b *MenuActionBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- MenuSkeleton ------------------------------------------------------------------------------

// ShowIcon, passed to MenuSkeleton, adds an icon placeholder.
type ShowIcon bool

// MenuSkeleton is a loading placeholder for a menu item, its text a
// random width between 50 and 90% (like upstream). Accepts ShowIcon.
func MenuSkeleton(children ...any) *h.Element {
	showIcon := false
	var rest []any
	for _, c := range children {
		if v, ok := c.(ShowIcon); ok {
			showIcon = bool(v)
			continue
		}
		rest = append(rest, c)
	}
	width := fmt.Sprintf("%d%%", rand.IntN(40)+50)
	base := h.Div(h.Data("slot", "sidebar-menu-skeleton"), h.Data("sidebar", "menu-skeleton"), h.Class(lib.Cn("flex h-8 items-center gap-2 rounded-md px-2")))
	out := h.Spread(base, h.Div(rest...), lib.Cn)
	if showIcon {
		out.Children(skeleton.Skeleton(h.Class("size-4 rounded-md"), h.Data("sidebar", "menu-skeleton-icon")))
	}
	out.Children(skeleton.Skeleton(
		h.Class("h-4 max-w-(--skeleton-width) flex-1"),
		h.Data("sidebar", "menu-skeleton-text"),
		h.Style("--skeleton-width: "+width),
	))
	return out
}

// --- MenuSubButton -----------------------------------------------------------------------------

// SubButtonSize is a MenuSubButton's size.
type SubButtonSize string

// MenuSubButton sizes.
const (
	SubButtonSizeSm SubButtonSize = "sm"
	SubButtonSizeMd SubButtonSize = "md"
)

// MenuSubButtonBuilder is returned by MenuSubButton.
type MenuSubButtonBuilder struct {
	h.Builder[MenuSubButtonBuilder]
	slottable
	size   SubButtonSize
	active bool
}

// MenuSubButton is a sub-menu entry's link (<a>). Accepts lib.AsChild,
// IsActive and a SubButtonSize (default SubButtonSizeMd).
func MenuSubButton(children ...any) *MenuSubButtonBuilder {
	b := &MenuSubButtonBuilder{size: SubButtonSizeMd}
	b.Builder = h.NewBuilder(b, "a")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild, IsActive and SubButtonSize.
func (b *MenuSubButtonBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case IsActive:
		b.active = bool(v)
	case SubButtonSize:
		b.size = v
	default:
		return b.accept(c)
	}
	return true
}

// AsChild merges the button into its single child.
func (b *MenuSubButtonBuilder) AsChild() *MenuSubButtonBuilder { b.asChild = true; return b }

// Element resolves the sub button.
func (b *MenuSubButtonBuilder) Element() *h.Element {
	classes := []string{
		"flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground ring-sidebar-ring outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground",
		"data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground",
	}
	if b.size == SubButtonSizeSm {
		classes = append(classes, "text-xs")
	}
	if b.size == SubButtonSizeMd {
		classes = append(classes, "text-sm")
	}
	classes = append(classes, "group-data-[collapsible=icon]:hidden")
	base := h.A(
		h.Data("slot", "sidebar-menu-sub-button"),
		h.Data("sidebar", "menu-sub-button"),
		h.Data("size", string(b.size)),
		h.Data("active", strconv.FormatBool(b.active)),
		h.Class(lib.Cn(classes...)),
	)
	return b.resolve("MenuSubButton", h.Spread(base, b.Elem(), lib.Cn))
}

// Render renders the sub button.
func (b *MenuSubButtonBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
