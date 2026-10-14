// Package separator is the shadcn/ui Separator (Radix Separator): a
// visual or semantic divider.
//
//	separator.Separator(h.Class("my-4"))
//	separator.Separator(separator.Vertical)
//	separator.Separator(separator.Decorative(false)) // role="separator"
package separator

import (
	"io"
	"shadcn.go/components/style"

	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

// Orientation is the separator's axis.
type Orientation string

// Orientations.
const (
	Horizontal Orientation = "horizontal"
	Vertical   Orientation = "vertical"
)

// Decorative, passed as a child, sets whether the separator is purely
// visual (role="none", the default) or semantic (role="separator").
type Decorative bool

// Builder is returned by Separator.
type Builder struct {
	h.Builder[Builder]
	orientation Orientation
	decorative  bool
}

// Separator returns a horizontal, decorative separator <div>. Children may
// include an Orientation and a Decorative.
func Separator(children ...any) *Builder {
	b := &Builder{orientation: Horizontal, decorative: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Orientation and Decorative children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Orientation:
		b.Orientation(v)
	case Decorative:
		b.decorative = bool(v)
	default:
		return false
	}
	return true
}

// Orientation sets the orientation; invalid values mean horizontal, as in
// Radix.
func (b *Builder) Orientation(o Orientation) *Builder {
	if o != Vertical {
		o = Horizontal
	}
	b.orientation = o
	return b
}

// Decorative sets whether the separator is purely visual.
func (b *Builder) Decorative(on bool) *Builder { b.decorative = on; return b }

const class = "shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px"

var cnClass = lib.Cn(class)

// Element resolves the separator to its <div>.
func (b *Builder) Element() *h.Element { return b.element(cnClass) }

// styledElement resolves the separator in style st
// (registry/bases/radix/ui/separator.tsx).
func (b *Builder) styledElement(st style.Style) *h.Element { return b.element(baseClass) }

var baseClass = lib.Cn("shrink-0 bg-border data-horizontal:h-px data-horizontal:w-full data-vertical:w-px data-vertical:self-stretch")

func (b *Builder) element(class string) *h.Element {
	base := h.Div(h.Data("orientation", string(b.orientation)))
	if b.decorative {
		base.Attr("role", "none")
	} else {
		if b.orientation == Vertical {
			base.Attr("aria-orientation", "vertical")
		}
		base.Attr("role", "separator")
	}
	base.Children(h.Data("slot", "separator"), h.Class(class))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the separator, in the style in scope (style.Scope) if any.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn resolves the element as it renders inside w: in the style in
// scope (style.Scope), if any.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		return b.styledElement(st)
	}
	return b.Element()
}
