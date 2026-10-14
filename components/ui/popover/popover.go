// Package popover is the shadcn/ui Popover: rich content in a portal,
// triggered by a button (Radix Popover, @radix-ui/react-popover 1.1.15).
//
//	popover.Popover(
//		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open popover")),
//		popover.Content(h.Class("w-80"), ...),
//	)
//
// Options are children: popover.Modal(true) and popover.DefaultOpen on
// Popover; positioning (popover.SideTop, popover.AlignStart,
// popover.SideOffset(8), popover.AvoidCollisions(false), ... see package
// popper) on Content, which defaults to align center and a side offset of
// 4 like shadcn's PopoverContent.
//
// Markup is what Radix renders. Closed content is not in the DOM: Content
// renders an inert <template> that the behavior (popover.js) clones into
// <body> when the popover opens and removes after the exit animation.
// Behavior: the trigger toggles; Escape, a pointer down outside and focus
// moving outside close (non-modal); modal popovers lock scroll, hide the
// rest of the page from assistive technology, block outside pointer events
// and trap focus; focus moves into the content when it opens and back to
// the trigger when it closes; the content is positioned against the
// trigger (or Anchor) by the shared popper primitive. The template
// dispatches "shadcn:open-change" (cancelable, detail.open); scripts can
// drive a popover with getPopover from popover.js.
package popover

import (
	_ "embed"
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/dialog/primitive"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/popper"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed popover.js
var script string

func init() { runtime.Module("popover", script) }

// Marker attributes the behavior reads and removes.
const (
	markerAnchor = "data-radix-popover-anchor"
	markerClose  = "data-radix-popover-close"
)

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

var (
	contentClass     = "z-50 w-72 origin-(--radix-popover-content-transform-origin) rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-hidden data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"
	headerClass      = lib.Cn("flex flex-col gap-1 text-sm")
	titleClass       = lib.Cn("font-medium")
	descriptionClass = lib.Cn("text-muted-foreground")
)

// The classes of registry/bases/radix/ui/popover.tsx, used when a style is
// in scope (style.Scope).
var (
	baseFile        = style.NewFile("ui/popover.tsx")
	baseContent     = baseFile.Cn("cn-popover-content z-50 w-72 origin-(--radix-popover-content-transform-origin) outline-hidden", style.Arg)
	baseHeader      = baseFile.Cn("cn-popover-header", style.Arg)
	baseTitle       = baseFile.Cn("cn-popover-title", style.Arg)
	baseDescription = baseFile.Cn("cn-popover-description", style.Arg)
)

// classFor returns the default class of a part: new-york-v4's, or the
// base site resolved in the style in scope.
func classFor(w io.Writer, def string, site *style.Site) string {
	if st, ok := style.FromWriter(w); ok {
		return site.Class(st, "")
	}
	return def
}

// Scope is the state a Popover shares with its parts (Radix's
// PopoverContext).
type Scope struct {
	ContentID   string
	Modal       bool
	DefaultOpen bool
}

// NewScope returns a scope with a fresh content id.
func NewScope() *Scope { return &Scope{ContentID: uid.New("radix")} }

// Scoped is embedded by the parts: Popover hands them its scope.
type Scoped struct{ scope *Scope }

// SetPopoverScope is called by Popover.
func (s *Scoped) SetPopoverScope(sc *Scope) { s.scope = sc }

// Scope returns the popover's scope; a part outside any Popover gets one
// of its own.
func (s *Scoped) Scope() *Scope {
	if s.scope == nil {
		s.scope = NewScope()
	}
	return s.scope
}

type participant interface{ SetPopoverScope(*Scope) }

// Option configures a Popover; pass options as children.
type Option func(*Scope)

// Modal sets whether the popover is modal (default false).
func Modal(on bool) Option { return func(s *Scope) { s.Modal = on } }

// DefaultOpen opens the popover when the page loads.
var DefaultOpen Option = func(s *Scope) { s.DefaultOpen = true }

// --- Popover -------------------------------------------------------------------------------

// RootBuilder is returned by Popover.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts []Option
}

// Popover is the root. It renders its children only, after sharing the
// popover's scope with the parts among them.
func Popover(children ...any) *RootBuilder {
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

// Modal sets whether the popover is modal.
func (b *RootBuilder) Modal(on bool) *RootBuilder { b.opts = append(b.opts, Modal(on)); return b }

// DefaultOpen opens the popover when the page loads.
func (b *RootBuilder) DefaultOpen() *RootBuilder { b.opts = append(b.opts, DefaultOpen); return b }

// Render renders the children.
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
		case participant:
			c.SetPopoverScope(sc)
		}
		return true
	})
	return h.Group(kids...).Render(w)
}

func state(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}

// part holds what the parts share.
type part struct {
	Scoped
	asChild bool
}

func (p *part) accept(c any) bool {
	switch v := c.(type) {
	case lib.AsChildFlag:
		p.asChild = true
	case *Scope:
		p.SetPopoverScope(v)
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
		return h.Errored(fmt.Errorf("popover %s: AsChild needs exactly one element child, got %d children", name, len(kids)))
	}
	out, err := primitive.Slot(props.Clone().SetChildren(), kids[0])
	if err != nil {
		return h.Errored(fmt.Errorf("popover %s: AsChild: %w", name, err))
	}
	return out
}

// --- Trigger -------------------------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	part
}

// Trigger is the button that toggles the popover. Accepts lib.AsChild.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *TriggerBuilder) AcceptChild(c any) bool { return b.accept(c) }

// AsChild renders the single child in place of the <button>.
func (b *TriggerBuilder) AsChild() *TriggerBuilder { b.asChild = true; return b }

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	sc := b.Scope()
	d := h.Button(
		h.Type("button"),
		h.Aria("haspopup", "dialog"),
		h.Aria("expanded", fmt.Sprint(sc.DefaultOpen)),
		h.Aria("controls", sc.ContentID),
		h.Data("state", state(sc.DefaultOpen)),
		h.Data("slot", "popover-trigger"),
	)
	return b.resolve("Trigger", d, b.Elem())
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Anchor --------------------------------------------------------------------------------

// AnchorBuilder is returned by Anchor.
type AnchorBuilder struct {
	h.Builder[AnchorBuilder]
	part
}

// Anchor is an element to position the content against instead of the
// trigger. Accepts lib.AsChild.
func Anchor(children ...any) *AnchorBuilder {
	b := &AnchorBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *AnchorBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the anchor.
func (b *AnchorBuilder) Element() *h.Element {
	d := h.Div(h.Data("slot", "popover-anchor"), h.Attr(markerAnchor, b.Scope().ContentID))
	return b.resolve("Anchor", d, b.Elem())
}

// Render renders the anchor.
func (b *AnchorBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	part
	opts popper.Options
}

// Content is the popover's content, portalled to <body> while open and
// positioned against the trigger. Accepts the positioning options.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{opts: popper.Defaults()}
	b.opts.Align = popper.AlignCenter
	b.opts.SideOffset = 4
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes positioning options, lib.AsChild and *Scope.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if b.opts.Accept(c) {
		return true
	}
	return b.accept(c)
}

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
	d := h.Div(
		h.Data("state", "open"),
		h.Role("dialog"),
		h.ID(sc.ContentID),
		h.Data("slot", "popover-content"),
		h.Class(class),
		h.TabIndex(-1),
	)
	el := b.resolve("Content", d, b.Elem())
	popper.AppendStyle(el, popper.ContentVars("popover"))
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-popover", ""),
		h.Data("content-id", sc.ContentID),
		h.If(sc.Modal, h.Data("modal", "")),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		popper.Wrap(b.opts, el),
	)
}

// Render renders the template.
func (b *ContentBuilder) Render(w io.Writer) error {
	return b.element(classFor(w, contentClass, baseContent)).Render(w)
}

// --- Close ---------------------------------------------------------------------------------

// CloseBuilder is returned by Close.
type CloseBuilder struct {
	h.Builder[CloseBuilder]
	part
}

// Close is Radix's PopoverClose: a button that closes the popover.
// Accepts lib.AsChild.
func Close(children ...any) *CloseBuilder {
	b := &CloseBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *CloseBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the button.
func (b *CloseBuilder) Element() *h.Element {
	d := h.Button(h.Type("button"), h.Attr(markerClose, b.Scope().ContentID))
	return b.resolve("Close", d, b.Elem())
}

// Render renders the button.
func (b *CloseBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Arrow ---------------------------------------------------------------------------------

// ArrowBuilder is returned by Arrow.
type ArrowBuilder struct {
	h.Builder[ArrowBuilder]
	side popper.Side
}

// Arrow is Radix's PopoverArrow: an arrow pointing at the anchor. Put it
// inside Content.
func Arrow(children ...any) *ArrowBuilder {
	b := &ArrowBuilder{}
	b.Builder = h.NewBuilder(b, "svg")
	return b.Children(children...)
}

// Element resolves the arrow.
func (b *ArrowBuilder) Element() *h.Element { return popper.Arrow(b.side, b.Elem()) }

// Render renders the arrow.
func (b *ArrowBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Header, Title, Description ------------------------------------------------------------------

// Header lays out a title and a description.
func Header(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "popover-header"), h.Class(classFor(w, headerClass, baseHeader))), h.Div(children...), lib.Cn).Render(w)
	})
}

// Title is the popover's title (a <div>, like upstream).
func Title(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "popover-title"), h.Class(classFor(w, titleClass, baseTitle))), h.Div(children...), lib.Cn).Render(w)
	})
}

// Description is the popover's description.
func Description(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.P(h.Data("slot", "popover-description"), h.Class(classFor(w, descriptionClass, baseDescription))), h.P(children...), lib.Cn).Render(w)
	})
}
