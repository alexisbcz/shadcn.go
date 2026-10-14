// Package hovercard is the shadcn/ui HoverCard: for sighted users to
// preview content available behind a link (Radix HoverCard,
// @radix-ui/react-hover-card 1.1.15).
//
//	hovercard.HoverCard(
//		hovercard.Trigger(lib.AsChild, button.Button(button.VariantLink, "@nextjs")),
//		hovercard.Content(h.Class("w-80"), ...),
//	)
//
// HoverCard accepts hovercard.OpenDelay(ms) (default 700),
// hovercard.CloseDelay(ms) (default 300) and hovercard.DefaultOpen;
// Content accepts the positioning options (default side bottom, align
// center, side offset 4).
//
// Markup is what Radix renders: the trigger is an <a> (use lib.AsChild
// for another element); closed content is not in the DOM, Content renders
// an inert <template> that hover-card.js clones into <body> when the card
// opens. The card opens when the pointer rests on the trigger (or it gets
// focus) for OpenDelay, closes CloseDelay after the pointer leaves the
// trigger and the content (or on blur, Escape, a pointer down outside),
// stays open while text in it is being selected, and its focusable
// descendants are taken out of the tab order.
package hovercard

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

//go:embed hover-card.js
var script string

func init() { runtime.Module("hover-card", script) }

const markerTrigger = "data-radix-hover-card-trigger"

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

const contentClass = "z-50 w-64 origin-(--radix-hover-card-content-transform-origin) rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-hidden data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"

// OpenDelay is the time in milliseconds from the pointer entering the
// trigger until the card opens.
type OpenDelay int

// CloseDelay is the time in milliseconds from the pointer leaving the
// trigger or the content until the card closes.
type CloseDelay int

// Option configures a HoverCard; pass options as children.
type Option func(*Scope)

// DefaultOpen opens the card when the page loads.
var DefaultOpen Option = func(s *Scope) { s.DefaultOpen = true }

// Scope is the state a HoverCard shares with its parts.
type Scope struct {
	ContentID   string
	DefaultOpen bool
	OpenDelay   int
	CloseDelay  int
}

func newScope() *Scope { return &Scope{ContentID: uid.New("radix"), OpenDelay: 700, CloseDelay: 300} }

// --- HoverCard -----------------------------------------------------------------------------

// RootBuilder is returned by HoverCard.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts []Option
}

// HoverCard is the root. It renders its children only, after sharing the
// card's scope with its parts.
func HoverCard(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Options, OpenDelay and CloseDelay.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Option:
		b.opts = append(b.opts, v)
	case OpenDelay:
		b.opts = append(b.opts, func(s *Scope) { s.OpenDelay = int(v) })
	case CloseDelay:
		b.opts = append(b.opts, func(s *Scope) { s.CloseDelay = int(v) })
	default:
		return false
	}
	return true
}

// Render renders the children.
func (b *RootBuilder) Render(w io.Writer) error {
	sc := newScope()
	for _, o := range b.opts {
		o(sc)
	}
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
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

// Trigger is the link that opens the card. Accepts lib.AsChild.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "a")
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

// AsChild renders the single child in place of the <a>.
func (b *TriggerBuilder) AsChild() *TriggerBuilder { b.asChild = true; return b }

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the trigger as it renders inside w (an asChild child
// keeps the style in scope).
func (b *TriggerBuilder) ElementIn(w io.Writer) *h.Element {
	sc := b.Scope()
	st := "closed"
	if sc.DefaultOpen {
		st = "open"
	}
	d := h.A(h.Data("state", st), h.Data("slot", "hover-card-trigger"), h.Attr(markerTrigger, sc.ContentID))
	props := h.Spread(d, b.Elem(), lib.Cn)
	if !b.asChild {
		return props
	}
	kids := props.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("hovercard Trigger: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := primitive.SlotIn(w, props.Clone().SetChildren(), kids[0])
	if err != nil {
		return h.Errored(fmt.Errorf("hovercard Trigger: AsChild: %w", err))
	}
	return out
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	scoped
	opts popper.Options
}

// Content is the card, portalled to <body> while open. Accepts the
// positioning options.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{opts: popper.Defaults()}
	b.opts.SideOffset = 4
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes positioning options.
func (b *ContentBuilder) AcceptChild(c any) bool { return b.opts.Accept(c) }

// Options returns the content's positioning options for changes.
func (b *ContentBuilder) Options() *popper.Options { return &b.opts }

// Element resolves the template holding the content as it is when open.
func (b *ContentBuilder) Element() *h.Element { return b.element(contentClass) }

func (b *ContentBuilder) element(class string) *h.Element {
	sc := b.Scope()
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case *ArrowBuilder:
			c.side = b.opts.Side
		}
		return true
	})
	d := h.Div(h.Data("state", "open"), h.Data("slot", "hover-card-content"), h.Class(class))
	el := h.Spread(d, b.Elem(), lib.Cn)
	popper.AppendStyle(el, popper.ContentVars("hover-card"))
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-hover-card", ""),
		h.Data("content-id", sc.ContentID),
		h.Data("open-delay", strconv.Itoa(sc.OpenDelay)),
		h.Data("close-delay", strconv.Itoa(sc.CloseDelay)),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		popper.Wrap(withDir(b.opts, el), el),
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

// The content class of registry/bases/radix/ui/hover-card.tsx, used when a style is
// in scope.
var (
	baseFile    = style.NewFile("ui/hover-card.tsx")
	baseContent = baseFile.Cn("cn-hover-card-content z-50 origin-(--radix-hover-card-content-transform-origin) outline-hidden", style.Arg)
)

// --- Arrow ---------------------------------------------------------------------------------

// ArrowBuilder is returned by Arrow.
type ArrowBuilder struct {
	h.Builder[ArrowBuilder]
	side popper.Side
}

// Arrow is Radix's HoverCardArrow: an arrow pointing at the trigger. Put
// it inside Content.
func Arrow(children ...any) *ArrowBuilder {
	b := &ArrowBuilder{}
	b.Builder = h.NewBuilder(b, "svg")
	return b.Children(children...)
}

// Element resolves the arrow.
func (b *ArrowBuilder) Element() *h.Element { return popper.Arrow(b.side, b.Elem()) }

// Render renders the arrow.
func (b *ArrowBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// withDir is o with the content's dir: Radix's popper wrapper repeats the
// content's dir prop.
func withDir(o popper.Options, el *h.Element) popper.Options {
	if d, ok := el.GetAttr("dir"); ok {
		o.Dir = d
	}
	return o
}
