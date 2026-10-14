// Package collapsible is the shadcn/ui Collapsible (Radix Collapsible): an
// interactive component which expands and collapses a panel.
//
//	collapsible.Collapsible(h.Class("flex flex-col gap-2"),
//	    collapsible.Trigger(lib.AsChild, button.Button(button.VariantGhost, "Toggle")),
//	    collapsible.Content("Hidden until opened"),
//	)
//
// Root children may include Open (the initial state, Radix defaultOpen) and
// Disabled. Trigger accepts lib.AsChild.
//
// Behavior (collapsible.js): clicking the trigger toggles the content. The
// content is mounted only while open (and while its exit animation runs,
// Radix Presence); on every change its size is measured into
// --radix-collapsible-content-height/-width for height animations. A
// cancelable "shadcn:open-change" event (detail: { open }) is dispatched on
// the root before every change.
//
// The content's children are rendered in place when open and always carried
// in a <template> the behavior mounts them from (a fresh copy on every open,
// as React remounts them).
package collapsible

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed collapsible.js
var script string

func init() { runtime.Module("collapsible", script) }

// Open, passed to the root, sets whether the collapsible starts open.
type Open bool

// Disabled, passed to the root, disables the collapsible.
type Disabled bool

type context struct {
	open      bool
	disabled  bool
	contentID string
}

// State returns Radix's data-state value for open.
func State(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}

// Builder is returned by Collapsible.
type Builder struct {
	h.Builder[Builder]
	open     bool
	disabled bool
	asChild  bool
}

// Collapsible returns the root <div>. It accepts lib.AsChild: the root's
// props are then merged into its single child.
func Collapsible(children ...any) *Builder {
	b := &Builder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Open and Disabled children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Open:
		b.open = bool(v)
	case Disabled:
		b.disabled = bool(v)
	case lib.AsChildFlag:
		b.asChild = true
	default:
		return false
	}
	return true
}

// Open sets whether the collapsible starts open.
func (b *Builder) Open(on bool) *Builder { b.open = on; return b }

// Disabled disables the collapsible.
func (b *Builder) Disabled(on bool) *Builder { b.disabled = on; return b }

// Element resolves the collapsible to its root <div>.
func (b *Builder) Element() *h.Element {
	ctx := &context{open: b.open, disabled: b.disabled, contentID: uid.New("radix")}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *TriggerBuilder:
			c.ctx = ctx
		case *ContentBuilder:
			c.ctx = ctx
		case *Builder:
			return false // a nested collapsible provides its own context
		}
		return true
	})
	base := h.Div(h.Data("state", State(b.open)))
	if b.disabled {
		base.Data("disabled", "")
	}
	base.Data("slot", "collapsible")
	if b.asChild {
		// The child is typically a component spreading its props last
		// (e.g. SidebarMenuItem): the root's data-slot wins, as in React.
		out := aschild.Resolve("collapsible", h.Spread(base, b.Elem(), nil))
		out.Data("slot", "collapsible")
		return out
	}
	return h.Spread(base, b.Elem(), nil)
}

// Render renders the collapsible.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	ctx     *context
	asChild bool
}

// Trigger returns the <button> that toggles the collapsible. It accepts
// lib.AsChild.
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
func (b *TriggerBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the trigger as it renders inside w: an AsChild child
// keeps the style of an enclosing style.Scope.
func (b *TriggerBuilder) ElementIn(w io.Writer) *h.Element {
	if b.ctx == nil {
		return h.Errored(fmt.Errorf("collapsible: Trigger must be inside Collapsible"))
	}
	base := h.Button(
		h.Data("slot", "collapsible-trigger"),
		h.Type("button"),
		h.Aria("controls", b.ctx.contentID),
		h.Aria("expanded", fmt.Sprint(b.ctx.open)),
		h.Data("state", State(b.ctx.open)),
	)
	if b.ctx.disabled {
		base.Data("disabled", "")
		base.Disabled(true)
	}
	props := h.Spread(base, b.Elem(), nil)
	if b.asChild {
		return aschild.ResolveIn(w, "collapsible Trigger", props)
	}
	return props
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	ctx *context
}

// Content returns the collapsible panel <div>.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the content.
func (b *ContentBuilder) Element() *h.Element {
	if b.ctx == nil {
		return h.Errored(fmt.Errorf("collapsible: Content must be inside Collapsible"))
	}
	base := h.Div(h.Data("slot", "collapsible-content"))
	return ContentElement(base, b.Elem(), b.ctx.open, b.ctx.disabled, b.ctx.contentID, "collapsible-content-template")
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// MountPreventionStyle is the inline style Radix's CollapsibleContent sets on
// content that is open when it mounts, so it does not animate in.
const MountPreventionStyle = "transition-duration:0s;animation-name:none"

// ContentElement renders a Radix CollapsibleContent the way React's first
// render (plus its mount effect) leaves it: base carries the component's
// own attributes (data-slot, role, ...), props the user's attributes and
// children. Open content has its children and the mount-prevention style;
// closed content is hidden and empty. The children are also carried in a
// <template data-slot=templateSlot> for collapsible-content.js. Accordion
// reuses it.
func ContentElement(base, props *h.Element, open, disabled bool, id, templateSlot string) *h.Element {
	base = base.Clone()
	base.Data("state", State(open))
	if disabled {
		base.Data("disabled", "")
	}
	base.ID(id)
	if !open {
		base.Hidden(true)
	}
	user := props.Clone().SetChildren()
	style, _ := user.GetAttr("style")
	if bs, ok := base.GetAttr("style"); ok && bs != "" {
		style = joinStyle(bs, style)
	}
	if open {
		style = joinStyle(style, MountPreventionStyle)
	}
	user.Attr("style", style)
	out := h.Spread(base, user, nil)
	kids := props.ChildNodes()
	if open {
		out.Children(kids...)
	}
	out.Children(h.Template(append([]any{h.Data("slot", templateSlot)}, kids...)...))
	return out
}

func joinStyle(a, b string) string {
	a = strings.TrimSuffix(strings.TrimSpace(a), ";")
	b = strings.TrimSpace(b)
	switch {
	case a == "":
		return b
	case b == "":
		return a
	}
	return a + ";" + b
}
