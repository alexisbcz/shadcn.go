// Package direction is the shadcn/ui DirectionProvider: it sets the reading
// direction (Radix's DirectionProvider context) of every direction-aware
// component below it, without rendering any element of its own.
//
//	direction.Provider(direction.RTL,
//		tabs.Tabs(...),          // renders dir="rtl", arrow keys flip
//		radiogroup.RadioGroup(...),
//	)
//
// Radix resolves a component's direction as its own dir prop, else the
// nearest DirectionProvider's, else "ltr". The Go components take the dir
// prop as an h.Dir attribute on their root and render the resolved value
// (dir="ltr" by default); Provider passes its direction to every component
// below it that renders a dir attribute and was not given one, so the
// server-rendered markup and the components' behavior match React's. A
// nested Provider sets the direction of its own subtree.
//
// Components that keep the direction elsewhere than a dir attribute on
// their root (a Dir option on a root that renders no element, like
// DropdownMenu, Select or Accordion) receive it by implementing Receiver;
// until they do, they must be given it explicitly (dm.Dir("rtl"), ...).
package direction

import (
	"io"

	"shadcn.go/components/ui/internal/tree"
	h "shadcn.go/libhtml"
)

// Dir is a reading direction.
type Dir string

// Directions.
const (
	LTR Dir = "ltr"
	RTL Dir = "rtl"
)

// Receiver is implemented by components that take the direction from the
// nearest Provider through a method instead of a dir attribute on their
// root. SetDirectionContext is called before rendering with the provider's
// direction; the component uses it only when it was not given its own
// (Radix's useDirection(dir): the prop wins over the context).
type Receiver interface {
	SetDirectionContext(dir string)
}

// Builder is returned by Provider.
type Builder struct {
	children []any
	dir      Dir
}

// Provider returns a DirectionProvider. Children may include a Dir
// (upstream's dir or direction prop); without one the direction is "ltr".
func Provider(children ...any) *Builder {
	b := &Builder{dir: LTR}
	for _, c := range children {
		if d, ok := c.(Dir); ok {
			if d != "" {
				b.dir = d
			}
			continue
		}
		b.children = append(b.children, c)
	}
	return b
}

// Dir returns the provider's direction (upstream useDirection, for the
// components below it).
func (b *Builder) Dir() Dir { return b.dir }

// Render renders the children with the provider's direction applied.
func (b *Builder) Render(w io.Writer) error {
	tree.Walk(b.children, func(n any) bool {
		switch c := n.(type) {
		case *Builder:
			return false // a nested provider applies its own direction
		case Receiver:
			c.SetDirectionContext(string(b.dir))
		case tree.Container:
			el := c.Elem()
			if el == nil {
				return true
			}
			if _, local := el.GetAttr("dir"); local {
				return true // the component's own dir prop wins
			}
			if rendersDir(n) {
				el.Attr("dir", string(b.dir))
			}
		}
		return true
	})
	return h.Group(b.children...).Render(w)
}

// rendersDir reports whether a component's root renders a dir attribute,
// i.e. whether it reads the direction context. Components that cannot be
// resolved out of their parent's context report false.
func rendersDir(n any) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	e, isElementer := n.(h.Elementer)
	if !isElementer {
		return false
	}
	el := e.Element()
	if el == nil {
		return false
	}
	_, ok = el.GetAttr("dir")
	return ok
}
