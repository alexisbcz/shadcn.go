// Package navigationmenu is the shadcn/ui NavigationMenu: a collection of
// links for navigating websites, with dropdown panels (Radix
// NavigationMenu, @radix-ui/react-navigation-menu 1.2.14).
//
//	navigationmenu.NavigationMenu(
//		navigationmenu.List(
//			navigationmenu.Item(
//				navigationmenu.Trigger("Getting started"),
//				navigationmenu.Content(...links...),
//			),
//			navigationmenu.Item(
//				navigationmenu.Link(lib.AsChild, h.Class(navigationmenu.TriggerStyle()), h.A(h.Href("/docs"), "Docs")),
//			),
//		),
//	)
//
// Options are children. On NavigationMenu: UseViewport(false) (shadcn's
// viewport prop, default true: the panels render in a shared viewport
// below the list), DefaultValue (the open item's value), DelayDuration and
// SkipDelayDuration (milliseconds, defaults 200 and 300), OrientationVertical
// and Dir("rtl"). On Item: ItemValue (default: a generated id). On Trigger:
// Disabled. On Link: Active.
//
// Markup is what Radix renders. Closed panels are not in the DOM: Content
// renders an inert <template> in its item that navigation-menu.js clones
// into place when the item opens and removes after the exit animation.
// Behavior: hover opens after DelayDuration (immediately while another
// item was open less than SkipDelayDuration ago), leaving closes after
// 150ms, a click toggles, arrow keys / Home / End move focus between
// triggers and links, ArrowDown enters an open panel, Tab moves through it
// and out, Escape and outside pointer down close, panels slide between
// items (data-motion), the Indicator follows the open trigger and the
// viewport takes the size of the open panel. The <nav> dispatches
// "shadcn:value-change" (detail.value).
package navigationmenu

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

//go:embed navigation-menu.js
var script string

func init() { runtime.Module("navigation-menu", script) }

// Marker attributes navigation-menu.js reads and removes.
const (
	markerRoot      = "data-radix-navigation-menu"
	markerItem      = "data-radix-navigation-menu-item"
	markerIndicator = "data-radix-navigation-menu-indicator"
	markerViewport  = "data-radix-navigation-menu-viewport"
)

// Orientation of the menu.
type Orientation string

// Orientations.
const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// UseViewport sets whether panels render in the shared viewport (default
// true).
type UseViewport bool

// DefaultValue is the value of the item open on load.
type DefaultValue string

// DelayDuration is the hover delay in milliseconds before an item opens.
type DelayDuration int

// SkipDelayDuration is how long in milliseconds after closing an item
// another one opens without delay.
type SkipDelayDuration int

// Dir is the reading direction ("ltr" or "rtl").
type Dir string

// ItemValue is an item's value.
type ItemValue string

// Flag is a boolean option.
type Flag int

// Flags: Disabled applies to Trigger, Active to Link.
const (
	Disabled Flag = iota + 1
	Active
)

var (
	rootClass      = "group/navigation-menu relative flex max-w-max flex-1 items-center justify-center"
	listClass      = "group flex flex-1 list-none items-center justify-center gap-1"
	itemClass      = "relative"
	contentClass   = "top-0 left-0 w-full p-2 pr-2.5 data-[motion=from-end]:slide-in-from-right-52 data-[motion=from-start]:slide-in-from-left-52 data-[motion=to-end]:slide-out-to-right-52 data-[motion=to-start]:slide-out-to-left-52 data-[motion^=from-]:animate-in data-[motion^=from-]:fade-in data-[motion^=to-]:animate-out data-[motion^=to-]:fade-out md:absolute md:w-auto"
	contentClass2  = "group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden group-data-[viewport=false]/navigation-menu:rounded-md group-data-[viewport=false]/navigation-menu:border group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:shadow group-data-[viewport=false]/navigation-menu:duration-200 **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none group-data-[viewport=false]/navigation-menu:data-[state=closed]:animate-out group-data-[viewport=false]/navigation-menu:data-[state=closed]:fade-out-0 group-data-[viewport=false]/navigation-menu:data-[state=closed]:zoom-out-95 group-data-[viewport=false]/navigation-menu:data-[state=open]:animate-in group-data-[viewport=false]/navigation-menu:data-[state=open]:fade-in-0 group-data-[viewport=false]/navigation-menu:data-[state=open]:zoom-in-95"
	viewportOuter  = "absolute top-full left-0 isolate z-50 flex justify-center"
	viewportClass  = "origin-top-center relative mt-1.5 h-[var(--radix-navigation-menu-viewport-height)] w-full overflow-hidden rounded-md border bg-popover text-popover-foreground shadow data-[state=closed]:animate-out data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:zoom-in-90 md:w-[var(--radix-navigation-menu-viewport-width)]"
	linkClass      = "flex flex-col gap-1 rounded-sm p-2 text-sm transition-all outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 data-[active=true]:bg-accent/50 data-[active=true]:text-accent-foreground data-[active=true]:hover:bg-accent data-[active=true]:focus:bg-accent [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground"
	indicatorClass = "top-full z-[1] flex h-1.5 items-end justify-center overflow-hidden data-[state=hidden]:animate-out data-[state=hidden]:fade-out data-[state=visible]:animate-in data-[state=visible]:fade-in"
	chevronClass   = "relative top-[1px] ml-1 size-3 transition duration-300 group-data-[state=open]:rotate-180"
)

var triggerStyle = cva.New(
	"group inline-flex h-9 w-max items-center justify-center rounded-md bg-background px-4 py-2 text-sm font-medium transition-[color,box-shadow] outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent",
	cva.Config{},
)

// The classes of registry/bases/radix/ui/navigation-menu.tsx, used inside a
// style.Scope (see package style).
var (
	file              = style.NewFile("ui/navigation-menu.tsx")
	stTriggerStyle    = file.CVA("cn-navigation-menu-trigger group/navigation-menu-trigger inline-flex h-9 w-max items-center justify-center outline-none disabled:pointer-events-none", nil, nil)
	stRoot            = file.Cn("cn-navigation-menu group/navigation-menu relative flex max-w-max flex-1 items-center justify-center", style.Arg)
	stList            = file.Cn("cn-navigation-menu-list group flex flex-1 list-none items-center justify-center", style.Arg)
	stItem            = file.Cn("cn-navigation-menu-item relative", style.Arg)
	stTrigger         = file.Cn(style.Arg, "group", style.Arg)
	stTriggerIcon     = file.Attr("cn-navigation-menu-trigger-icon")
	stContent         = file.Cn("cn-navigation-menu-content top-0 left-0 w-full group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none md:absolute md:w-auto", style.Arg)
	stViewportWrapper = file.Cn("cn-navigation-menu-viewport-wrapper absolute top-full left-0 isolate z-50 flex justify-center")
	stViewport        = file.Cn("cn-navigation-menu-viewport origin-top-center relative mt-1.5 h-(--radix-navigation-menu-viewport-height) w-full overflow-hidden md:w-(--radix-navigation-menu-viewport-width)", style.Arg)
	stLink            = file.Cn("cn-navigation-menu-link", style.Arg)
	stIndicator       = file.Cn("cn-navigation-menu-indicator top-full z-1 flex h-1.5 items-end justify-center overflow-hidden", style.Arg)
	stIndicatorArrow  = file.Attr("cn-navigation-menu-indicator-arrow relative top-[60%] h-2 w-2 rotate-45")
)

// styled records the enclosing style.Scope ("" for new-york-v4).
type styled struct {
	st  style.Style
	rtl bool // inside style.ScopeRTL: the ui-rtl classes
}

func (s *styled) setStyle(w io.Writer) {
	s.st, _ = style.FromWriter(w)
	s.rtl = s.st != "" && style.RTLFromWriter(w)
}

// pick returns nova (resolved for the style) inside a style.Scope, else
// the new-york-v4 classes.
func (s *styled) pick(newYork string, site *style.Site, values ...string) string {
	if s.st == "" {
		return newYork
	}
	if s.rtl {
		return style.RTLClass(site.Class(s.st, values...))
	}
	return site.Class(s.st, values...)
}

// TriggerStyleFor is TriggerStyle in style st ("" is new-york-v4).
func TriggerStyleFor(st style.Style, className ...string) string {
	if st == "" {
		return TriggerStyle(className...)
	}
	return lib.Cn(stTriggerStyle.Classes(st, cva.Selection{}, className...))
}

// TriggerStyle returns navigationMenuTriggerStyle(): the trigger classes,
// for links styled as triggers.
func TriggerStyle(className ...string) string {
	return lib.Cn(triggerStyle.Classes(cva.Selection{}, className...))
}

// Scope is the state a NavigationMenu shares with its parts.
type Scope struct {
	BaseID      string
	Value       string
	Orientation Orientation
	Dir         string
	Viewport    bool
}

func (s *Scope) dir() string {
	if s.Dir == "" {
		return "ltr"
	}
	return s.Dir
}

func (s *Scope) orientation() string {
	if s.Orientation == "" {
		return string(OrientationHorizontal)
	}
	return string(s.Orientation)
}

func (s *Scope) triggerID(v string) string { return s.BaseID + "-trigger-" + v }
func (s *Scope) contentID(v string) string { return s.BaseID + "-content-" + v }

type scoped struct{ scope *Scope }

func (s *scoped) setScope(sc *Scope) { s.scope = sc }

func (s *scoped) Scope() *Scope {
	if s.scope == nil {
		s.scope = &Scope{BaseID: uid.New("radix"), Viewport: true}
	}
	return s.scope
}

type participant interface{ setScope(*Scope) }

func openState(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}

func resolve(name string, asChild bool, defaults, user *h.Element) *h.Element {
	props := h.Spread(defaults, user, lib.Cn)
	if !asChild {
		return props
	}
	kids := props.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("navigation-menu %s: AsChild needs exactly one element child, got %d children", name, len(kids)))
	}
	out, err := aschild.Slot(props.Clone().SetChildren(), kids[0])
	if err != nil {
		return h.Errored(fmt.Errorf("navigation-menu %s: AsChild: %w", name, err))
	}
	return out
}

// --- NavigationMenu ------------------------------------------------------------------------

// RootBuilder is returned by NavigationMenu.
type RootBuilder struct {
	h.Builder[RootBuilder]
	styled
	scope     Scope
	delay     int
	skipDelay int
}

// NavigationMenu is the root <nav>.
func NavigationMenu(children ...any) *RootBuilder {
	b := &RootBuilder{scope: Scope{Viewport: true}, delay: 200, skipDelay: 300}
	b.Builder = h.NewBuilder(b, "nav")
	return b.Children(children...)
}

// AcceptChild consumes the root options.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case UseViewport:
		b.scope.Viewport = bool(v)
	case DefaultValue:
		b.scope.Value = string(v)
	case DelayDuration:
		b.delay = int(v)
	case SkipDelayDuration:
		b.skipDelay = int(v)
	case Orientation:
		b.scope.Orientation = v
	case Dir:
		b.scope.Dir = string(v)
	default:
		return false
	}
	return true
}

// Element resolves the <nav>.
func (b *RootBuilder) Element() *h.Element {
	sc := b.scope
	sc.BaseID = uid.New("radix")
	user := b.Elem().Clone()
	kids := user.ChildNodes()
	user.SetChildren()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case participant:
			c.setScope(&sc)
		}
		return true
	})
	opts, _ := json.Marshal(map[string]any{"delayDuration": b.delay, "skipDelayDuration": b.skipDelay, "value": sc.Value})
	d := h.Nav(
		h.Attr(markerRoot, string(opts)),
		h.Aria("label", "Main"),
		h.Data("orientation", sc.orientation()),
		h.Dir(sc.dir()),
		h.Data("slot", "navigation-menu"),
		h.Data("viewport", fmt.Sprint(sc.Viewport)),
		h.Class(b.pick(lib.Cn(rootClass), stRoot, "")),
	)
	el := h.Spread(d, user, lib.Cn)
	el.Children(kids...)
	if sc.Viewport {
		vp := Viewport()
		vp.setScope(&sc)
		vp.styled = b.styled
		el.Children(vp)
	}
	return el
}

// Render renders the <nav>.
func (b *RootBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- List, Item ----------------------------------------------------------------------------

// ListBuilder is returned by List.
type ListBuilder struct {
	h.Builder[ListBuilder]
	styled
	scoped
}

// List is the <ul> of items, wrapped in the indicator track.
func List(children ...any) *ListBuilder {
	b := &ListBuilder{}
	b.Builder = h.NewBuilder(b, "ul")
	return b.Children(children...)
}

// Element resolves the track and list.
func (b *ListBuilder) Element() *h.Element {
	sc := b.Scope()
	d := h.Ul(
		h.Data("orientation", sc.orientation()),
		h.Data("slot", "navigation-menu-list"),
		h.Class(b.pick(lib.Cn(listClass), stList, "")),
		h.Dir(sc.dir()),
	)
	return h.Div(h.Style("position: relative;"), h.Spread(d, b.Elem(), lib.Cn))
}

// Render renders the list.
func (b *ListBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	styled
	scoped
	value string
}

// Item is an <li> holding a Trigger and its Content, or a Link.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "li")
	return b.Children(children...)
}

// AcceptChild consumes ItemValue.
func (b *ItemBuilder) AcceptChild(c any) bool {
	if v, ok := c.(ItemValue); ok {
		b.value = string(v)
		return true
	}
	return false
}

// Value sets the item's value.
func (b *ItemBuilder) Value(v string) *ItemBuilder { b.value = v; return b }

// Element resolves the <li>.
func (b *ItemBuilder) Element() *h.Element {
	if b.value == "" {
		b.value = uid.New("radix")
	}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ItemBuilder, *RootBuilder:
			return false
		case *TriggerBuilder:
			c.item = b.value
		case *ContentBuilder:
			c.item = b.value
		}
		return true
	})
	d := h.Li(h.Attr(markerItem, b.value), h.Data("slot", "navigation-menu-item"), h.Class(b.pick(lib.Cn(itemClass), stItem, "")))
	return h.Spread(d, b.Elem(), lib.Cn)
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Trigger -------------------------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	styled
	scoped
	item     string
	disabled bool
}

// Trigger is the button that opens its item's Content.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Disabled.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	if c == Disabled {
		b.disabled = true
		return true
	}
	return false
}

// Disabled disables the trigger.
func (b *TriggerBuilder) Disabled(on bool) *TriggerBuilder { b.disabled = on; return b }

// Element resolves the trigger. The focus proxy and the aria-owns span
// Radix renders after an open trigger are added by navigation-menu.js.
func (b *TriggerBuilder) Element() *h.Element {
	sc := b.Scope()
	open := sc.Value != "" && sc.Value == b.item
	d := h.Button(
		h.ID(sc.triggerID(b.item)),
		h.If(b.disabled, h.Attr("disabled", "")),
		h.If(b.disabled, h.Data("disabled", "")),
		h.Data("state", openState(open)),
		h.Aria("expanded", fmt.Sprint(open)),
		h.Aria("controls", sc.contentID(b.item)),
		h.Data("slot", "navigation-menu-trigger"),
		h.Class(b.pick(lib.Cn(TriggerStyle(), "group"), stTrigger, TriggerStyleFor(b.st), "")),
		h.Data("radix-collection-item", ""),
	)
	el := h.Spread(d, b.Elem(), lib.Cn)
	el.Children(" ", lucide.ChevronDown(h.Class(b.pick(chevronClass, stTriggerIcon)), h.Aria("hidden", "true")))
	return el
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	styled
	scoped
	item string
}

// Content is the panel of an item, shown while the item is open.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the template holding the panel as it renders open.
func (b *ContentBuilder) Element() *h.Element {
	sc := b.Scope()
	d := h.Div(
		h.ID(sc.contentID(b.item)),
		h.Aria("labelledby", sc.triggerID(b.item)),
		h.Data("orientation", sc.orientation()),
		h.If(!sc.Viewport, h.Data("state", "open")),
		h.Data("slot", "navigation-menu-content"),
		h.Class(b.pick(lib.Cn(contentClass, contentClass2), stContent, "")),
		h.Dir(sc.dir()),
	)
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-navigation-menu-content", b.item),
		h.Spread(d, b.Elem(), lib.Cn),
	)
}

// Render renders the template.
func (b *ContentBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Link ----------------------------------------------------------------------------------

// LinkBuilder is returned by Link.
type LinkBuilder struct {
	h.Builder[LinkBuilder]
	styled
	active  bool
	asChild bool
}

// Link is a navigation link. Accepts Active and lib.AsChild.
func Link(children ...any) *LinkBuilder {
	b := &LinkBuilder{}
	b.Builder = h.NewBuilder(b, "a")
	return b.Children(children...)
}

// AcceptChild consumes Active and lib.AsChild.
func (b *LinkBuilder) AcceptChild(c any) bool {
	switch c.(type) {
	case lib.AsChildFlag:
		b.asChild = true
		return true
	case Flag:
		if c == Active {
			b.active = true
			return true
		}
	}
	return false
}

// Active marks the link as the current page.
func (b *LinkBuilder) Active(on bool) *LinkBuilder { b.active = on; return b }

// Element resolves the link.
func (b *LinkBuilder) Element() *h.Element {
	d := h.A(
		h.If(b.active, h.Data("active", "")),
		h.If(b.active, h.Aria("current", "page")),
		h.Data("slot", "navigation-menu-link"),
		h.Class(b.pick(lib.Cn(linkClass), stLink, "")),
		h.Data("radix-collection-item", ""),
	)
	return resolve("Link", b.asChild, d, b.Elem())
}

// Render renders the link.
func (b *LinkBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Indicator -----------------------------------------------------------------------------

// IndicatorBuilder is returned by Indicator.
type IndicatorBuilder struct {
	h.Builder[IndicatorBuilder]
	styled
	scoped
}

// Indicator is the arrow under the open trigger. Place it in the List.
func Indicator(children ...any) *IndicatorBuilder {
	b := &IndicatorBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the template holding the indicator.
func (b *IndicatorBuilder) Element() *h.Element {
	sc := b.Scope()
	d := h.Div(
		h.Aria("hidden", "true"),
		h.Data("state", "visible"),
		h.Data("orientation", sc.orientation()),
		h.Data("slot", "navigation-menu-indicator"),
		h.Class(b.pick(lib.Cn(indicatorClass), stIndicator, "")),
	)
	el := h.Spread(d, b.Elem(), lib.Cn)
	el.Children(h.Div(h.Class(b.pick("relative top-[60%] h-2 w-2 rotate-45 rounded-tl-sm bg-border shadow-md", stIndicatorArrow))))
	return h.Template(h.Data("parity-ignore", ""), h.Attr(markerIndicator, ""), el)
}

// Render renders the template.
func (b *IndicatorBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }

// --- Viewport ------------------------------------------------------------------------------

// ViewportBuilder is returned by Viewport.
type ViewportBuilder struct {
	h.Builder[ViewportBuilder]
	styled
	scoped
}

// Viewport is where panels render when the menu uses a viewport.
// NavigationMenu renders it unless UseViewport(false).
func Viewport(children ...any) *ViewportBuilder {
	b := &ViewportBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the viewport's positioning wrapper, with the viewport
// (present while an item is open) in a template.
func (b *ViewportBuilder) Element() *h.Element {
	sc := b.Scope()
	d := h.Div(
		h.Data("state", "open"),
		h.Data("orientation", sc.orientation()),
		h.Data("slot", "navigation-menu-viewport"),
		h.Class(b.pick(lib.Cn(viewportClass), stViewport, "")),
	)
	return h.Div(h.Class(b.pick(lib.Cn(viewportOuter), stViewportWrapper)),
		h.Template(h.Data("parity-ignore", ""), h.Attr(markerViewport, ""), h.Spread(d, b.Elem(), lib.Cn)),
	)
}

// Render renders the viewport wrapper.
func (b *ViewportBuilder) Render(w io.Writer) error { b.setStyle(w); return b.Element().Render(w) }
