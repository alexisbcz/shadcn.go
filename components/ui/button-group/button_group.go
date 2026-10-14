// Package buttongroup is the shadcn/ui ButtonGroup: a container that groups
// related buttons together with consistent styling.
//
//	buttongroup.ButtonGroup(
//		button.Button(button.VariantOutline, "Archive"),
//		button.Button(button.VariantOutline, "Report"),
//	)
//	buttongroup.ButtonGroup(buttongroup.Vertical, ...)
package buttongroup

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/separator"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Orientation is the group's axis.
type Orientation string

// Orientations.
const (
	Horizontal Orientation = "horizontal"
	Vertical   Orientation = "vertical"
)

var variants = cva.New(
	"flex w-fit items-stretch has-[>[data-slot=button-group]]:gap-2 [&>*]:focus-visible:relative [&>*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-md [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1",
	cva.Config{
		Variants: cva.Variants{
			"orientation": {
				// Upstream: [&>*:not(:last-child)]:rounded-r-none. A Go dropdown menu
				// or tooltip keeps its content in a <template> after the trigger
				// (React portals it away), so a child followed only by a template
				// is the visual last child.
				"horizontal": "[&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child):not(:has(+template:last-child))]:rounded-r-none",
				"vertical":   "flex-col [&>*:not(:first-child)]:rounded-t-none [&>*:not(:first-child)]:border-t-0 [&>*:not(:last-child):not(:has(+template:last-child))]:rounded-b-none",
			},
		},
		Defaults: cva.Selection{"orientation": "horizontal"},
	},
)

// Variants returns the classes of a group with the given orientation plus
// className, like buttonGroupVariants() upstream. An empty orientation
// selects horizontal.
func Variants(o Orientation, className ...string) string {
	return lib.Cn(variants.Classes(cva.Selection{"orientation": string(o)}, className...))
}

// Builder is returned by ButtonGroup, Text and Separator.
type Builder struct {
	h.Builder[Builder]
	kind        int
	orientation Orientation
	asChild     bool
}

const (
	kindGroup = iota
	kindText
	kindSeparator
)

// ButtonGroup returns the group <div role="group">. Children may include an
// Orientation; like upstream, data-orientation is only rendered when one is
// given.
func ButtonGroup(children ...any) *Builder {
	b := &Builder{kind: kindGroup}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Text is a text segment of the group (upstream ButtonGroupText). Pass
// lib.AsChild to style your own element (e.g. a label) instead of a <div>.
func Text(children ...any) *Builder {
	b := &Builder{kind: kindText}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Separator divides two items of the group (upstream ButtonGroupSeparator);
// it is vertical unless given Horizontal (a separator.Orientation is
// accepted too).
func Separator(children ...any) *Builder {
	b := &Builder{kind: kindSeparator, orientation: Vertical}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Orientation, separator.Orientation and lib.AsChild.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Orientation:
		b.orientation = v
	case separator.Orientation:
		b.orientation = Orientation(v)
	case lib.AsChildFlag:
		if b.kind != kindText {
			return false
		}
		b.asChild = true
	default:
		return false
	}
	return true
}

// Orientation sets the orientation.
func (b *Builder) Orientation(o Orientation) *Builder { b.orientation = o; return b }

// AsChild makes Text render its single element child instead of a <div>.
func (b *Builder) AsChild() *Builder { b.asChild = true; return b }

const (
	textClass      = "flex items-center gap-2 rounded-md border bg-muted px-4 text-sm font-medium shadow-xs [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
	separatorClass = "relative m-0! self-stretch bg-input data-[orientation=vertical]:h-auto"
)

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the component to its element as it renders inside w
// (in the style of an enclosing style.Scope).
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	switch b.kind {
	case kindText:
		el := h.Spread(h.Div(h.Class(textClassFor(st))), b.Elem(), lib.Cn)
		if !b.asChild {
			return el
		}
		kids := el.ChildNodes()
		if len(kids) != 1 {
			return h.Errored(fmt.Errorf("buttongroup: Text AsChild needs exactly one element child, got %d children", len(kids)))
		}
		out, err := h.SlotIn(w, el.Clone().SetChildren(), kids[0], slotClassMerger(kids[0]))
		if err != nil {
			return h.Errored(fmt.Errorf("buttongroup: Text AsChild: %w", err))
		}
		return out
	case kindSeparator:
		return h.ResolveIn(w, b.separator(st))
	}
	base := h.Div(h.Role("group"), h.Data("slot", "button-group"))
	if b.orientation != "" {
		base.Data("orientation", string(b.orientation))
	}
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base.Class(groupClass(st, b.orientation, className))
	return h.Spread(base, props, lib.Cn)
}

func (b *Builder) separator(st style.Style) *separator.Builder {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass(separatorClassFor(st, className))
	return separator.Separator(separator.Orientation(b.orientation), h.Data("slot", "button-group-separator"), spreadProp(props))
}

// slotClassMerger returns how Radix Slot's merged className ends up on the
// rendered child. Slot concatenates its classes with the child element's
// className prop. When the child is a shadcn component (it carries a
// data-slot), that component then runs cn(its base classes, merged), so
// conflicting base classes give way to the Slot's: the Go child is already
// rendered, so its classes are merged with the Slot's last. A plain element
// child keeps the plain concatenation.
func slotClassMerger(child any) func(...string) string {
	var el *h.Element
	switch c := child.(type) {
	case *h.Element:
		el = c
	case h.Elementer:
		el = c.Element()
	}
	if el == nil || !el.Has("data-slot") {
		return nil
	}
	return func(classes ...string) string { return lib.Cn(classes[1], classes[0]) }
}

// spreadProp applies props' attributes, classes and children to the
// element it is passed to, like {...props}.
func spreadProp(props *h.Element) h.Prop {
	return func(e *h.Element) { *e = *h.Spread(e, props, nil) }
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error {
	return b.ElementIn(w).Render(w)
}
