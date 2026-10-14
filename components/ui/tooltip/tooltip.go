// Package tooltip is the shadcn/ui Tooltip: a popup that displays
// information related to an element when it receives keyboard focus or
// the mouse hovers over it (Radix Tooltip, @radix-ui/react-tooltip 1.2.8).
//
//	tooltip.Provider(                     // optional, see below
//		tooltip.Tooltip(
//			tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Hover")),
//			tooltip.Content(h.P("Add to library")),
//		),
//	)
//
// Provider shares the open delay between tooltips: once one is open,
// moving to another opens it instantly until SkipDelayDuration elapsed
// without any open. Its options: tooltip.DelayDuration(ms) (default 0,
// like shadcn's TooltipProvider), tooltip.SkipDelayDuration(ms) (default
// 300), tooltip.DisableHoverableContent(true). Radix requires a provider;
// tooltips rendered outside one share a page-wide provider with the
// defaults above, like the TooltipProvider the shadcn docs mount in their
// layout. Tooltip accepts DelayDuration, DisableHoverableContent and
// DefaultOpen; Content accepts the positioning options (default side top,
// align center, side offset 0) and AriaLabel.
//
// Markup is what Radix renders: closed content is not in the DOM; Content
// renders an inert <template> cloned into <body> by tooltip.js when the
// tooltip opens.
package tooltip

import (
	_ "embed"
	"fmt"
	"io"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/dialog/primitive"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/popper"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed tooltip.js
var script string

func init() { runtime.Module("tooltip", script) }

const markerTrigger = "data-radix-tooltip-trigger"

// Positioning options accepted by Content (see package popper).
type (
	Side                   = popper.Side
	Align                  = popper.Align
	SideOffset             = popper.SideOffset
	AlignOffset            = popper.AlignOffset
	AvoidCollisions        = popper.AvoidCollisions
	CollisionPadding       = popper.CollisionPadding
	CollisionPaddingSides  = popper.CollisionPaddingSides
	CollisionBoundary      = popper.CollisionBoundary
	Sticky                 = popper.Sticky
	HideWhenDetached       = popper.HideWhenDetached
	UpdatePositionStrategy = popper.UpdatePositionStrategy
	ArrowPadding           = popper.ArrowPadding
)

// Sides, alignments and sticky behaviors.
const (
	SideTop       = popper.SideTop
	SideRight     = popper.SideRight
	SideBottom    = popper.SideBottom
	SideLeft      = popper.SideLeft
	AlignStart    = popper.AlignStart
	AlignCenter   = popper.AlignCenter
	AlignEnd      = popper.AlignEnd
	StickyPartial = popper.StickyPartial
	StickyAlways  = popper.StickyAlways
)

const (
	contentClass = "z-50 w-fit origin-(--radix-tooltip-content-transform-origin) animate-in rounded-md bg-foreground px-3 py-1.5 text-xs text-balance text-background fade-in-0 zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95"
	// TooltipPrimitive.Arrow's className is a literal upstream (no cn).
	arrowClass = "z-50 size-2.5 translate-y-[calc(-50%_-_2px)] rotate-45 rounded-[2px] bg-foreground fill-foreground"
	// VisuallyHidden's style.
	visuallyHiddenStyle = "position: absolute; border: 0px; width: 1px; height: 1px; padding: 0px; margin: -1px; overflow: hidden; clip: rect(0px, 0px, 0px, 0px); white-space: nowrap; word-wrap: normal;"
)

// DelayDuration is the time in milliseconds from the pointer entering the
// trigger until the tooltip opens (Provider or Tooltip).
type DelayDuration int

// SkipDelayDuration is how long in milliseconds a user has to enter
// another trigger without incurring a delay again (Provider).
type SkipDelayDuration int

// DisableHoverableContent closes the tooltip as soon as the pointer
// leaves the trigger, instead of letting it move onto the content
// (Provider or Tooltip).
type DisableHoverableContent bool

// AriaLabel is the accessible label of the content announced by screen
// readers instead of its children (Content).
type AriaLabel string

// Option configures a Tooltip; pass options as children.
type Option func(*Scope)

// DefaultOpen opens the tooltip when the page loads.
var DefaultOpen Option = func(s *Scope) { s.DefaultOpen = true }

// providerConfig is a TooltipProvider's props.
type providerConfig struct {
	id                      string
	delayDuration           int
	skipDelayDuration       int
	disableHoverableContent bool
}

var defaultProvider = &providerConfig{id: "default", delayDuration: 0, skipDelayDuration: 300}

// Scope is the state a Tooltip shares with its parts.
type Scope struct {
	ContentID   string
	DefaultOpen bool
	provider    *providerConfig
	delay       *int
	disableHC   *bool
}

func newScope() *Scope { return &Scope{ContentID: uid.New("radix"), provider: defaultProvider} }

func (s *Scope) delayDuration() int {
	if s.delay != nil {
		return *s.delay
	}
	return s.provider.delayDuration
}

func (s *Scope) disableHoverableContent() bool {
	if s.disableHC != nil {
		return *s.disableHC
	}
	return s.provider.disableHoverableContent
}

// --- Provider ------------------------------------------------------------------------------

// ProviderBuilder is returned by Provider.
type ProviderBuilder struct {
	h.Builder[ProviderBuilder]
	cfg providerConfig
}

// Provider shares delays between the tooltips inside it. It renders its
// children only.
func Provider(children ...any) *ProviderBuilder {
	b := &ProviderBuilder{cfg: providerConfig{delayDuration: 0, skipDelayDuration: 300}}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes DelayDuration, SkipDelayDuration and
// DisableHoverableContent.
func (b *ProviderBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case DelayDuration:
		b.cfg.delayDuration = int(v)
	case SkipDelayDuration:
		b.cfg.skipDelayDuration = int(v)
	case DisableHoverableContent:
		b.cfg.disableHoverableContent = bool(v)
	default:
		return false
	}
	return true
}

// Render renders the children.
func (b *ProviderBuilder) Render(w io.Writer) error {
	cfg := b.cfg
	cfg.id = uid.New("tooltip-provider")
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *ProviderBuilder:
			return false
		case *RootBuilder:
			c.provider = &cfg
		}
		return true
	})
	return h.Group(kids...).Render(w)
}

// --- Tooltip -------------------------------------------------------------------------------

// RootBuilder is returned by Tooltip.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts     []Option
	provider *providerConfig
}

// Tooltip is the root. It renders its children only, after sharing the
// tooltip's scope with its parts.
func Tooltip(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Options, DelayDuration and DisableHoverableContent.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Option:
		b.opts = append(b.opts, v)
	case DelayDuration:
		d := int(v)
		b.opts = append(b.opts, func(s *Scope) { s.delay = &d })
	case DisableHoverableContent:
		on := bool(v)
		b.opts = append(b.opts, func(s *Scope) { s.disableHC = &on })
	default:
		return false
	}
	return true
}

// Render renders the children.
func (b *RootBuilder) Render(w io.Writer) error {
	sc := newScope()
	if b.provider != nil {
		sc.provider = b.provider
	}
	for _, o := range b.opts {
		o(sc)
	}
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder, *ProviderBuilder:
			return false
		case interface{ setScope(*Scope) }:
			c.setScope(sc)
		}
		return true
	})
	return h.Group(kids...).Render(w)
}

type scoped struct{ scope *Scope }

func (s *scoped) setScope(sc *Scope) { s.scope = sc }

func (s *scoped) Scope() *Scope {
	if s.scope == nil {
		s.scope = newScope()
	}
	return s.scope
}

// --- Trigger -------------------------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	scoped
	asChild bool
}

// Trigger is the element the tooltip describes. Accepts lib.AsChild.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	if _, ok := c.(lib.AsChildFlag); ok {
		b.asChild = true
		return true
	}
	return false
}

// AsChild renders the single child in place of the <button>.
func (b *TriggerBuilder) AsChild() *TriggerBuilder { b.asChild = true; return b }

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	sc := b.Scope()
	st := "closed"
	if sc.DefaultOpen {
		st = "instant-open"
	}
	d := h.Button(
		h.If(sc.DefaultOpen, h.Aria("describedby", sc.ContentID)),
		h.Data("state", st),
		h.Data("slot", "tooltip-trigger"),
		h.Attr(markerTrigger, sc.ContentID),
	)
	props := h.Spread(d, b.Elem(), lib.Cn)
	if !b.asChild {
		return props
	}
	kids := props.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("tooltip Trigger: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := primitive.Slot(props.Clone().SetChildren(), kids[0])
	if err != nil {
		return h.Errored(fmt.Errorf("tooltip Trigger: AsChild: %w", err))
	}
	return out
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	scoped
	opts      popper.Options
	ariaLabel *string
}

// Content is the tooltip, portalled to <body> while open, with an arrow
// pointing at the trigger. Accepts the positioning options and AriaLabel.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{opts: popper.Defaults()}
	b.opts.Side = popper.SideTop
	b.opts.SideOffset = 0
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes positioning options and AriaLabel.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if v, ok := c.(AriaLabel); ok {
		s := string(v)
		b.ariaLabel = &s
		return true
	}
	return b.opts.Accept(c)
}

// Options returns the content's positioning options for changes.
func (b *ContentBuilder) Options() *popper.Options { return &b.opts }

// Element resolves the template holding the content as it is when open.
func (b *ContentBuilder) Element() *h.Element { return b.element(contentClass) }

func (b *ContentBuilder) element(class string) *h.Element {
	sc := b.Scope()
	user := b.Elem().Clone()
	label := b.ariaLabel
	if v, ok := user.GetAttr("aria-label"); ok {
		label = &v
		user.Remove("aria-label")
	}
	d := h.Div(
		h.Data("state", "delayed-open"),
		h.Data("slot", "tooltip-content"),
		h.Class(class),
	)
	el := h.Spread(d, user, lib.Cn)
	popper.AppendStyle(el, popper.ContentVars("tooltip"))
	// <TooltipPrimitive.Arrow className=... />, then Radix's visually
	// hidden copy of the children (filled by tooltip.js) or of aria-label.
	el.Children(popper.Arrow(b.opts.Side, h.El("svg", h.Class(arrowClass))))
	vh := h.Span(h.ID(sc.ContentID), h.Role("tooltip"), h.Style(visuallyHiddenStyle))
	if label != nil {
		vh.Children(*label)
	} else {
		vh.Attr("data-radix-tooltip-copy", "")
	}
	el.Children(vh)
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-tooltip", ""),
		h.Data("content-id", sc.ContentID),
		h.Data("provider-id", sc.provider.id),
		h.Data("delay-duration", strconv.Itoa(sc.delayDuration())),
		h.Data("skip-delay-duration", strconv.Itoa(sc.provider.skipDelayDuration)),
		h.If(sc.disableHoverableContent(), h.Data("disable-hoverable-content", "")),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		popper.Wrap(b.opts, el),
	)
}

// Render renders the template, in the style in scope (style.Scope) if any.
func (b *ContentBuilder) Render(w io.Writer) error {
	class := contentClass
	if st, ok := style.FromWriter(w); ok {
		class = baseContent.Class(st, "")
	}
	return b.element(class).Render(w)
}

// The content class of registry/bases/radix/ui/tooltip.tsx, used when a style is
// in scope.
var (
	baseFile    = style.NewFile("ui/tooltip.tsx")
	baseContent = baseFile.Cn("cn-tooltip-content z-50 w-fit max-w-xs origin-(--radix-tooltip-content-transform-origin) bg-foreground text-background", style.Arg)
)
