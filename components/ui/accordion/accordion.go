// Package accordion is the shadcn/ui Accordion (Radix Accordion): a
// vertically stacked set of interactive headings that each reveal a
// section of content.
//
//	accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("item-1"),
//	    accordion.Item(accordion.Value("item-1"),
//	        accordion.Trigger("Is it accessible?"),
//	        accordion.Content("Yes. It adheres to the WAI-ARIA design pattern."),
//	    ),
//	)
//
// Root children may include a Type (single, the default, or multiple),
// Collapsible (single type: whether the open item can be closed), Value
// (the initially open item; repeat it for several with TypeMultiple),
// Disabled, an Orientation and Dir (the reading direction for horizontal
// arrow keys, "ltr" by default: like Radix, an ancestor's dir attribute does
// not count; Radix does not render it). Item takes its Value and Disabled.
// Classes given to Content style its inner <div>, as upstream.
//
// Behavior (accordion.js): clicking a trigger (or Enter/Space) opens or
// closes its item; arrow keys (per orientation and direction), Home and End
// move focus among the enabled triggers, wrapping. Contents animate with
// --radix-accordion-content-height (Radix Collapsible + Presence). A
// cancelable "shadcn:value-change" event (detail: { value }, a string for
// single, an array for multiple) is dispatched on the root before every
// change.
package accordion

import (
	_ "embed"
	"fmt"
	"io"
	"slices"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed accordion.js
var script string

func init() { runtime.Module("accordion", script) }

// Type selects whether one (single) or several (multiple) items can be open.
type Type string

// Types.
const (
	TypeSingle   Type = "single"
	TypeMultiple Type = "multiple"
)

// Value is an item's value; on the root, an initially open item.
type Value string

// Collapsible, on a single accordion, lets the open item be closed.
type Collapsible bool

// Disabled disables the whole accordion (root) or one item.
type Disabled bool

// Orientation of the accordion (default vertical).
type Orientation string

// Orientations.
const (
	OrientationVertical   Orientation = "vertical"
	OrientationHorizontal Orientation = "horizontal"
)

// Dir is the reading direction used by horizontal keyboard navigation.
type Dir string

// Reading directions.
const (
	DirLTR Dir = "ltr"
	DirRTL Dir = "rtl"
)

type rootContext struct {
	typ         Type
	collapsible bool
	value       []string
	disabled    bool
	orientation Orientation
}

type itemContext struct {
	root      *rootContext
	open      bool
	disabled  bool
	triggerID string
	contentID string
}

// Builder is returned by Accordion.
type Builder struct {
	h.Builder[Builder]
	typ         Type
	collapsible bool
	value       []string
	disabled    bool
	orientation Orientation
	dir         Dir
}

// Accordion returns the root <div>.
func Accordion(children ...any) *Builder {
	b := &Builder{typ: TypeSingle, orientation: OrientationVertical}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Type, Collapsible, Value, Disabled, Orientation and
// Dir children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Type:
		b.typ = v
	case Collapsible:
		b.collapsible = bool(v)
	case Value:
		b.value = append(b.value, string(v))
	case Disabled:
		b.disabled = bool(v)
	case Orientation:
		b.orientation = v
	case Dir:
		b.dir = v
	default:
		return false
	}
	return true
}

// Type sets the accordion's type.
func (b *Builder) Type(t Type) *Builder { b.typ = t; return b }

// Collapsible sets whether a single accordion's open item can be closed.
func (b *Builder) Collapsible(on bool) *Builder { b.collapsible = on; return b }

// Value sets the initially open items.
func (b *Builder) Value(v ...string) *Builder { b.value = v; return b }

// Disabled disables every item.
func (b *Builder) Disabled(on bool) *Builder { b.disabled = on; return b }

// Orientation sets the orientation.
func (b *Builder) Orientation(o Orientation) *Builder { b.orientation = o; return b }

// Dir sets the reading direction.
func (b *Builder) Dir(d Dir) *Builder { b.dir = d; return b }

// Element resolves the accordion to its root <div>.
func (b *Builder) Element() *h.Element {
	typ := b.typ
	if typ != TypeMultiple {
		typ = TypeSingle
	}
	orientation := b.orientation
	if orientation != OrientationHorizontal {
		orientation = OrientationVertical
	}
	value := b.value
	if typ == TypeSingle && len(value) > 1 {
		value = value[:1]
	}
	ctx := &rootContext{typ: typ, collapsible: b.collapsible || typ == TypeMultiple, value: value, disabled: b.disabled, orientation: orientation}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ItemBuilder:
			c.root = ctx
			return false // the item hands the context to its parts
		case *Builder:
			return false
		}
		return true
	})
	base := h.Div(h.Data("slot", "accordion"), h.Data("orientation", string(orientation)))
	// Radix props with no DOM trace, read (and removed) by accordion.js.
	if typ == TypeMultiple {
		base.Data("type", "multiple")
	}
	if ctx.collapsible && typ == TypeSingle {
		base.Data("collapsible", "")
	}
	if b.dir != "" {
		base.Data("dir", string(b.dir))
	}
	return h.Spread(base, b.Elem(), nil)
}

// Render renders the accordion.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	root     *rootContext
	value    string
	disabled bool
}

// Item returns an accordion item <div>. Children may include its Value and
// Disabled.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Value and Disabled children.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		b.value = string(v)
	case Disabled:
		b.disabled = bool(v)
	default:
		return false
	}
	return true
}

// Value sets the item's value.
func (b *ItemBuilder) Value(v string) *ItemBuilder { b.value = v; return b }

// Disabled disables the item.
func (b *ItemBuilder) Disabled(on bool) *ItemBuilder { b.disabled = on; return b }

// Element resolves the item to its <div>.
func (b *ItemBuilder) Element() *h.Element {
	if b.root == nil {
		return h.Errored(fmt.Errorf("accordion: Item must be inside Accordion"))
	}
	ctx := &itemContext{
		root:      b.root,
		open:      b.value != "" && slices.Contains(b.root.value, b.value),
		disabled:  b.root.disabled || b.disabled,
		triggerID: uid.New("radix"),
		contentID: uid.New("radix"),
	}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *TriggerBuilder:
			c.item = ctx
		case *ContentBuilder:
			c.item = ctx
			return false // content may hold another accordion
		case *ItemBuilder, *Builder:
			return false
		}
		return true
	})
	base := h.Div(
		h.Data("state", collapsible.State(ctx.open)),
		h.Data("orientation", string(b.root.orientation)),
	)
	if ctx.disabled {
		base.Data("disabled", "")
	}
	base.Data("slot", "accordion-item")
	base.Data("value", b.value) // read (and removed) by accordion.js
	props := b.Elem()
	return h.Spread(base, props, nil).SetClass(lib.Cn("border-b last:border-b-0", props.ClassName()))
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	item *itemContext
}

// Trigger returns the item's heading (<h3 class="flex">) and its toggle
// <button>, followed by a chevron.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

const triggerClass = "flex flex-1 items-start justify-between gap-4 rounded-md py-4 text-left text-sm font-medium transition-all outline-none hover:underline focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 [&[data-state=open]>svg]:rotate-180"

// Element resolves the trigger to its heading.
func (b *TriggerBuilder) Element() *h.Element {
	it := b.item
	if it == nil {
		return h.Errored(fmt.Errorf("accordion: Trigger must be inside an Item"))
	}
	o := string(it.root.orientation)
	state := collapsible.State(it.open)
	header := h.H3(h.Data("orientation", o), h.Data("state", state))
	if it.disabled {
		header.Data("disabled", "")
	}
	header.Class("flex")

	btn := h.Button(
		h.Type("button"),
		h.Aria("controls", it.contentID),
		h.Aria("expanded", fmt.Sprint(it.open)),
		h.Data("state", state),
	)
	if it.disabled {
		btn.Data("disabled", "")
		btn.Disabled(true)
	}
	if it.open && !it.root.collapsible {
		btn.Aria("disabled", "true")
	}
	btn.Data("orientation", o)
	btn.ID(it.triggerID)
	btn.Data("slot", "accordion-trigger")
	btn.Data("radix-collection-item", "")
	props := b.Elem()
	out := h.Spread(btn, props, nil).SetClass(lib.Cn(triggerClass, props.ClassName()))
	out.Children(lucide.ChevronDown(h.Class("pointer-events-none size-4 shrink-0 translate-y-0.5 text-muted-foreground transition-transform duration-200")))
	return header.Children(out)
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	item *itemContext
}

// Content returns the item's region. Its classes go to the inner <div>.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

const contentClass = "overflow-hidden text-sm data-[state=closed]:animate-accordion-up data-[state=open]:animate-accordion-down"

const contentVars = "--radix-accordion-content-height:var(--radix-collapsible-content-height);--radix-accordion-content-width:var(--radix-collapsible-content-width)"

// Element resolves the content.
func (b *ContentBuilder) Element() *h.Element {
	it := b.item
	if it == nil {
		return h.Errored(fmt.Errorf("accordion: Content must be inside an Item"))
	}
	props := b.Elem()
	inner := h.Div(h.Class(lib.Cn("pt-0 pb-4", props.ClassName())))
	inner.Children(props.ChildNodes()...)
	user := props.Clone().SetChildren().SetClass("")
	user.Children(inner)
	base := h.Div(
		h.Role("region"),
		h.Aria("labelledby", it.triggerID),
		h.Data("orientation", string(it.root.orientation)),
		h.Data("slot", "accordion-content"),
		h.Class(contentClass),
		h.Style(contentVars),
	)
	return collapsible.ContentElement(base, user, it.open, it.disabled, it.contentID, "accordion-content-template")
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
