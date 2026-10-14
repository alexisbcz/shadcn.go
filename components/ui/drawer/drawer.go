// Package drawer is the shadcn/ui Drawer: a panel sliding in from an edge
// of the screen that can be dragged to dismiss (vaul 1.1.2).
//
//	drawer.Drawer(
//		drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open Drawer")),
//		drawer.Content(
//			drawer.Header(drawer.Title("Move Goal"), drawer.Description("...")),
//			...,
//			drawer.Footer(drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel"))),
//		),
//	)
//
// Options are passed as children of Drawer: DirectionTop (and the other
// directions), drawer.Modal(false), drawer.DefaultOpen, drawer.Dismissible(false),
// drawer.ShouldScaleBackground (with an element carrying
// data-vaul-drawer-wrapper), drawer.SnapPoints(...), drawer.HandleOnly.
//
// vaul wraps the Radix Dialog primitive: the markup is Radix's (package
// dialog/primitive) plus vaul's data-vaul-* attributes. The behavior
// (drawer.js) is the Radix Dialog behavior (focus scope, dismissable layer,
// scroll lock, aria-hidden, presence) with vaul's on top: no auto focus on
// open, drag to dismiss with velocity and close threshold, overlay fade
// while dragging, snap points, background scaling and the slide/fade
// animations of vaul's injected stylesheet (shipped identically, injected
// with the page's CSP nonce). The drawer dispatches "shadcn:open-change"
// (cancelable, detail.open) on its template; scripts can drive it with
// getDrawer from the module (see drawer.js).
package drawer

import (
	_ "embed"
	"encoding/json"
	"io"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/dialog/primitive"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed drawer.js
var script string

func init() { runtime.Module("drawer", script) }

// Direction is the edge the drawer comes from.
type Direction string

// Directions.
const (
	DirectionBottom Direction = "bottom"
	DirectionTop    Direction = "top"
	DirectionLeft   Direction = "left"
	DirectionRight  Direction = "right"
)

// Option configures a Drawer; pass options as children of Drawer.
type Option func(*Config)

// Config is the state a Drawer shares with its parts (vaul's DrawerContext
// plus the Radix dialog scope).
type Config struct {
	Dialog                *primitive.Scope
	Direction             Direction
	Dismissible           bool
	ShouldScaleBackground bool
	// SetBackgroundColorOnScale (default true) paints <body> black while
	// the background is scaled.
	SetBackgroundColorOnScale bool
	SnapPoints                []string // numbers (fractions) or "123px"
	FadeFromIndex             *int
	HandleOnly                bool
	NoBodyStyles              bool
	CloseThreshold            float64
	AutoFocus                 bool
}

func newConfig() *Config {
	return &Config{
		Dialog:                    primitive.NewScope(),
		Direction:                 DirectionBottom,
		Dismissible:               true,
		SetBackgroundColorOnScale: true,
		CloseThreshold:            0.25,
	}
}

// Root options.
var (
	// DefaultOpen opens the drawer when the page loads.
	DefaultOpen Option = func(c *Config) { c.Dialog.DefaultOpen = true }
	// ShouldScaleBackground scales the element marked
	// data-vaul-drawer-wrapper down while the drawer is open.
	ShouldScaleBackground Option = func(c *Config) { c.ShouldScaleBackground = true }
	// HandleOnly restricts dragging to the Handle.
	HandleOnly Option = func(c *Config) { c.HandleOnly = true }
	// NoBodyStyles leaves <body>'s styles alone.
	NoBodyStyles Option = func(c *Config) { c.NoBodyStyles = true }
	// AutoFocus moves focus into the drawer when it opens (vaul's autoFocus).
	AutoFocus Option = func(c *Config) { c.AutoFocus = true }
)

// Modal sets whether the drawer is modal (default true). Like vaul, a
// non-modal drawer has no overlay (so no scroll lock), ignores clicks
// outside and leaves the page clickable; the underlying Radix dialog stays
// modal (focus trap, aria-hidden page), since vaul does not pass modal on.
func Modal(on bool) Option { return func(c *Config) { c.Dialog.Modal = on } }

// Dismissible sets whether the drawer can be closed by dragging, Escape
// or a click outside (default true).
func Dismissible(on bool) Option { return func(c *Config) { c.Dismissible = on } }

// SetBackgroundColorOnScale sets whether <body> turns black behind the
// scaled background (default true).
func SetBackgroundColorOnScale(on bool) Option {
	return func(c *Config) { c.SetBackgroundColorOnScale = on }
}

// CloseThreshold is the fraction of the drawer's size it must be dragged
// to close (default 0.25).
func CloseThreshold(v float64) Option { return func(c *Config) { c.CloseThreshold = v } }

// SnapPoints sets the snap points: fractions of the screen (0.5) or pixel
// strings ("355px"). The drawer opens at the first one.
func SnapPoints(points ...any) Option {
	return func(c *Config) {
		c.SnapPoints = c.SnapPoints[:0]
		for _, p := range points {
			switch v := p.(type) {
			case string:
				c.SnapPoints = append(c.SnapPoints, v)
			case float64:
				c.SnapPoints = append(c.SnapPoints, strconv.FormatFloat(v, 'f', -1, 64))
			case int:
				c.SnapPoints = append(c.SnapPoints, strconv.Itoa(v))
			}
		}
	}
}

// FadeFromIndex sets the snap point index from which the overlay fades
// in (default: the last).
func FadeFromIndex(i int) Option { return func(c *Config) { c.FadeFromIndex = &i } }

// participant is implemented by the parts that read the drawer's config.
type participant interface{ setDrawer(*Config) }

type scoped struct {
	primitive.Scoped
	cfg *Config
}

func (s *scoped) setDrawer(c *Config) { s.cfg = c; s.SetDialogScope(c.Dialog) }

func (s *scoped) config() *Config {
	if s.cfg == nil {
		s.cfg = newConfig()
		s.cfg.Dialog = s.Scope()
	}
	return s.cfg
}

// --- Root --------------------------------------------------------------------------------------

// RootBuilder is returned by Drawer.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts []Option
}

// Drawer is the root (vaul's Drawer.Root). It renders its children only.
func Drawer(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Options and Directions.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Option:
		b.opts = append(b.opts, v)
	case Direction:
		b.opts = append(b.opts, func(c *Config) { c.Direction = v })
	default:
		return false
	}
	return true
}

// Options adds options.
func (b *RootBuilder) Options(opts ...Option) *RootBuilder {
	b.opts = append(b.opts, opts...)
	return b
}

// Render renders the children after sharing the drawer's state with them.
func (b *RootBuilder) Render(w io.Writer) error {
	cfg := newConfig()
	for _, o := range b.opts {
		o(cfg)
	}
	if cfg.Direction == "" {
		cfg.Direction = DirectionBottom
	}
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder, *primitive.RootBuilder:
			return false
		case participant:
			c.setDrawer(cfg)
		case primitive.Participant:
			c.SetDialogScope(cfg.Dialog)
		}
		return true
	})
	return h.Group(kids...).Render(w)
}

// --- Trigger, Close ----------------------------------------------------------------------------

// Trigger is the button that opens the drawer. Accepts lib.AsChild.
func Trigger(children ...any) *primitive.TriggerBuilder {
	return primitive.Trigger(append([]any{h.Data("slot", "drawer-trigger")}, children...)...)
}

// Close is a button that closes the drawer. Accepts lib.AsChild.
func Close(children ...any) *primitive.CloseBuilder {
	return primitive.Close(append([]any{h.Data("slot", "drawer-close")}, children...)...)
}

// --- Portal ------------------------------------------------------------------------------------

// PortalBuilder is returned by Portal.
type PortalBuilder struct {
	h.Builder[PortalBuilder]
	scoped
}

// Portal mounts its children (Overlay, the content) at the end of <body>
// while the drawer is open. It renders the inert <template> they are
// cloned from.
func Portal(children ...any) *PortalBuilder {
	b := &PortalBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// Element resolves the portal's template.
func (b *PortalBuilder) Element() *h.Element {
	cfg := b.config()
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder, *primitive.RootBuilder:
			return false
		case participant:
			c.setDrawer(cfg)
		case primitive.Participant:
			c.SetDialogScope(cfg.Dialog)
		}
		if p, ok := n.(interface{ SetInPortal() }); ok {
			p.SetInPortal()
		}
		return true
	})
	return template(cfg, kids...)
}

// Render renders the template.
func (b *PortalBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

func template(cfg *Config, children ...any) *h.Element {
	sc := cfg.Dialog
	var snap string
	if len(cfg.SnapPoints) > 0 {
		bs, _ := json.Marshal(cfg.SnapPoints)
		snap = string(bs)
	}
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("vaul-root", ""),
		h.Data("content-id", sc.ContentID),
		h.Data("direction", string(cfg.Direction)),
		h.If(!sc.Modal, h.Data("modal", "false")),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		h.If(!cfg.Dismissible, h.Data("dismissible", "false")),
		h.If(cfg.ShouldScaleBackground, h.Data("scale-background", "")),
		h.If(!cfg.SetBackgroundColorOnScale, h.Data("no-background-color", "")),
		h.If(cfg.NoBodyStyles, h.Data("no-body-styles", "")),
		h.If(cfg.HandleOnly, h.Data("handle-only", "")),
		h.If(cfg.AutoFocus, h.Data("auto-focus", "")),
		h.If(cfg.CloseThreshold != 0.25, h.Data("close-threshold", strconv.FormatFloat(cfg.CloseThreshold, 'f', -1, 64))),
		h.If(snap != "", h.Data("snap-points", snap)),
		h.If(cfg.FadeFromIndex != nil, h.Data("fade-from-index", strconv.Itoa(derefInt(cfg.FadeFromIndex)))),
		h.Group(children...),
	)
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// The classes of registry/bases/radix/ui/drawer.tsx, used when a style is
// in scope (style.Scope).
var (
	baseFile        = style.NewFile("ui/drawer.tsx")
	baseOverlay     = baseFile.Cn("cn-drawer-overlay fixed inset-0 z-50 data-open:animate-in data-open:fade-in-0 data-closed:animate-out data-closed:fade-out-0", style.Arg)
	baseContent     = baseFile.Cn("cn-drawer-content group/drawer-content fixed z-50 data-[vaul-drawer-direction=bottom]:inset-x-0 data-[vaul-drawer-direction=bottom]:bottom-0 data-[vaul-drawer-direction=bottom]:mt-24 data-[vaul-drawer-direction=bottom]:max-h-[80vh] data-[vaul-drawer-direction=left]:inset-y-0 data-[vaul-drawer-direction=left]:left-0 data-[vaul-drawer-direction=left]:w-3/4 data-[vaul-drawer-direction=right]:inset-y-0 data-[vaul-drawer-direction=right]:right-0 data-[vaul-drawer-direction=right]:w-3/4 data-[vaul-drawer-direction=top]:inset-x-0 data-[vaul-drawer-direction=top]:top-0 data-[vaul-drawer-direction=top]:mb-24 data-[vaul-drawer-direction=top]:max-h-[80vh] data-[vaul-drawer-direction=left]:sm:max-w-sm data-[vaul-drawer-direction=right]:sm:max-w-sm", style.Arg)
	baseHandle      = baseFile.Attr("cn-drawer-handle mx-auto hidden shrink-0 group-data-[vaul-drawer-direction=bottom]/drawer-content:block")
	baseHeader      = baseFile.Cn("cn-drawer-header flex flex-col group-data-[vaul-drawer-direction=bottom]/drawer-content:text-center group-data-[vaul-drawer-direction=top]/drawer-content:text-center", style.Arg)
	baseFooter      = baseFile.Cn("cn-drawer-footer mt-auto flex flex-col", style.Arg)
	baseTitle       = baseFile.Cn("cn-drawer-title cn-font-heading", style.Arg)
	baseDescription = baseFile.Cn("cn-drawer-description", style.Arg)
)

func styleOf(w io.Writer) style.Style {
	st, _ := style.FromWriter(w)
	return st
}

func classFor(st style.Style, def string, site *style.Site) string {
	if st != "" {
		return site.Class(st, "")
	}
	return def
}

// --- Overlay -----------------------------------------------------------------------------------

var overlayClass = lib.Cn("fixed inset-0 z-50 bg-black/50 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:animate-in data-[state=open]:fade-in-0")

// OverlayBuilder is returned by Overlay.
type OverlayBuilder struct {
	h.Builder[OverlayBuilder]
	scoped
}

// Overlay is the translucent layer behind the drawer. Modal drawers only.
func Overlay(children ...any) *OverlayBuilder {
	b := &OverlayBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the overlay as it is when open.
func (b *OverlayBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *OverlayBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *OverlayBuilder) element(st style.Style) *h.Element {
	cfg := b.config()
	hasSnap := len(cfg.SnapPoints) > 0
	// vaul's Drawer.Overlay, then shadcn's DrawerOverlay props.
	vaul := h.Div(
		h.Data("vaul-overlay", ""),
		h.Data("vaul-snap-points", strconv.FormatBool(hasSnap)),
		h.Data("vaul-snap-points-overlay", strconv.FormatBool(shouldFade(cfg))),
		h.Data("vaul-animate", strconv.FormatBool(!cfg.Dialog.DefaultOpen)),
		h.Data("slot", "drawer-overlay"),
		h.Class(classFor(st, overlayClass, baseOverlay)),
	)
	props := h.Spread(vaul, b.Elem(), lib.Cn)
	return primitive.Overlay(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the overlay.
func (b *OverlayBuilder) Render(w io.Writer) error { return render(b.element(styleOf(w)), w) }

// shouldFade is vaul's shouldFade for the first snap point.
func shouldFade(cfg *Config) bool {
	if len(cfg.SnapPoints) == 0 {
		return true
	}
	fade := len(cfg.SnapPoints) - 1
	if cfg.FadeFromIndex != nil {
		fade = *cfg.FadeFromIndex
	}
	return fade == 0
}

// --- Content -----------------------------------------------------------------------------------

const (
	contentBase   = "group/drawer-content fixed z-50 flex h-auto flex-col bg-background"
	contentTop    = "data-[vaul-drawer-direction=top]:inset-x-0 data-[vaul-drawer-direction=top]:top-0 data-[vaul-drawer-direction=top]:mb-24 data-[vaul-drawer-direction=top]:max-h-[80vh] data-[vaul-drawer-direction=top]:rounded-b-lg data-[vaul-drawer-direction=top]:border-b"
	contentBottom = "data-[vaul-drawer-direction=bottom]:inset-x-0 data-[vaul-drawer-direction=bottom]:bottom-0 data-[vaul-drawer-direction=bottom]:mt-24 data-[vaul-drawer-direction=bottom]:max-h-[80vh] data-[vaul-drawer-direction=bottom]:rounded-t-lg data-[vaul-drawer-direction=bottom]:border-t"
	contentRight  = "data-[vaul-drawer-direction=right]:inset-y-0 data-[vaul-drawer-direction=right]:right-0 data-[vaul-drawer-direction=right]:w-3/4 data-[vaul-drawer-direction=right]:border-l data-[vaul-drawer-direction=right]:sm:max-w-sm"
	contentLeft   = "data-[vaul-drawer-direction=left]:inset-y-0 data-[vaul-drawer-direction=left]:left-0 data-[vaul-drawer-direction=left]:w-3/4 data-[vaul-drawer-direction=left]:border-r data-[vaul-drawer-direction=left]:sm:max-w-sm"
	handleClass   = "mx-auto mt-4 hidden h-2 w-[100px] shrink-0 rounded-full bg-muted group-data-[vaul-drawer-direction=bottom]/drawer-content:block"
)

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	scoped
}

// Content is the drawer panel, portalled with its Overlay. It starts with
// the grab bar shown for bottom drawers.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Classes returns the content's classes plus className, like the cn(...)
// call of DrawerContent.
func Classes(className ...string) string {
	return lib.Cn(append([]string{contentBase, contentTop, contentBottom, contentRight, contentLeft}, className...)...)
}

// Element resolves the portal holding the overlay and the panel.
func (b *ContentBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *ContentBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *ContentBuilder) element(st style.Style) *h.Element {
	cfg := b.config()
	hasSnap := len(cfg.SnapPoints) > 0
	vaul := h.Div(
		h.Data("vaul-drawer-direction", string(cfg.Direction)),
		h.Data("vaul-drawer", ""),
		h.Data("vaul-delayed-snap-points", "false"),
		h.Data("vaul-snap-points", strconv.FormatBool(hasSnap)),
		h.Data("vaul-custom-container", "false"),
		h.Data("vaul-animate", strconv.FormatBool(!cfg.Dialog.DefaultOpen)),
		h.Data("slot", "drawer-content"),
	)
	user := b.Elem()
	props := h.Spread(vaul, user, nil)
	if st != "" {
		props.SetClass(baseContent.Class(st, user.ClassName()))
		props.SetChildren(append([]any{h.Div(h.Class(baseHandle.Class(st)))}, user.ChildNodes()...)...)
	} else {
		props.SetClass(Classes(user.ClassName()))
		props.SetChildren(append([]any{h.Div(h.Class(handleClass))}, user.ChildNodes()...)...)
	}

	ov := Overlay()
	ov.setDrawer(cfg)
	ov.SetInPortal()
	content := primitive.Content(primitive.WithProps(props))
	content.SetDialogScope(cfg.Dialog)
	content.SetInPortal()
	return template(cfg, ov, content)
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// --- Handle --------------------------------------------------------------------------------------

// HandleBuilder is returned by Handle.
type HandleBuilder struct {
	h.Builder[HandleBuilder]
	scoped
}

// Handle is vaul's Drawer.Handle: a grab bar that cycles snap points when
// tapped and drags the drawer (the only drag area with HandleOnly).
func Handle(children ...any) *HandleBuilder {
	b := &HandleBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the handle as it is when open.
func (b *HandleBuilder) Element() *h.Element {
	d := h.Div(h.Data("vaul-drawer-visible", "true"), h.Data("vaul-handle", ""), h.Aria("hidden", "true"))
	user := b.Elem()
	el := h.Spread(d, user.Clone().SetChildren(), nil)
	return el.SetChildren(h.Span(h.Data("vaul-handle-hitarea", ""), h.Aria("hidden", "true"), h.Group(user.ChildNodes()...)))
}

// Render renders the handle.
func (b *HandleBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Header, Footer, Title, Description ----------------------------------------------------------

var (
	headerClass = lib.Cn("flex flex-col gap-0.5 p-4 group-data-[vaul-drawer-direction=bottom]/drawer-content:text-center group-data-[vaul-drawer-direction=top]/drawer-content:text-center md:gap-1.5 md:text-left")
	footerClass = lib.Cn("mt-auto flex flex-col gap-2 p-4")
	titleClass  = lib.Cn("font-semibold text-foreground")
	descClass   = lib.Cn("text-sm text-muted-foreground")
)

// Header lays out the title and description.
func Header(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "drawer-header"), h.Class(classFor(styleOf(w), headerClass, baseHeader))), h.Div(children...), lib.Cn).Render(w)
	})
}

// Footer lays out the actions at the bottom of the drawer.
func Footer(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "drawer-footer"), h.Class(classFor(styleOf(w), footerClass, baseFooter))), h.Div(children...), lib.Cn).Render(w)
	})
}

// TitleBuilder is returned by Title.
type TitleBuilder struct {
	h.Builder[TitleBuilder]
	scoped
}

// Title is the drawer's accessible name (<h2>).
func Title(children ...any) *TitleBuilder {
	b := &TitleBuilder{}
	b.Builder = h.NewBuilder(b, "h2")
	return b.Children(children...)
}

// Element resolves the title.
func (b *TitleBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *TitleBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *TitleBuilder) element(st style.Style) *h.Element {
	props := h.Spread(h.H2(h.Data("slot", "drawer-title"), h.Class(classFor(st, titleClass, baseTitle))), b.Elem(), lib.Cn)
	return primitive.Title(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the title.
func (b *TitleBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// DescriptionBuilder is returned by Description.
type DescriptionBuilder struct {
	h.Builder[DescriptionBuilder]
	scoped
}

// Description is the drawer's accessible description (<p>).
func Description(children ...any) *DescriptionBuilder {
	b := &DescriptionBuilder{}
	b.Builder = h.NewBuilder(b, "p")
	return b.Children(children...)
}

// Element resolves the description.
func (b *DescriptionBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *DescriptionBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *DescriptionBuilder) element(st style.Style) *h.Element {
	props := h.Spread(h.P(h.Data("slot", "drawer-description"), h.Class(classFor(st, descClass, baseDescription))), b.Elem(), lib.Cn)
	return primitive.Description(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the description.
func (b *DescriptionBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

func render(el *h.Element, w io.Writer) error {
	if el == nil {
		return nil
	}
	return el.Render(w)
}
