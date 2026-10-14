// Package styled lets a new-york-v4 component also render the markup of a
// style of registry/bases/radix (radix-nova, ...) when it is rendered inside
// style.Scope. The new-york-v4 rendering is unchanged outside a Scope.
//
// A component returns a *Node: its Element (used by Slot / asChild, and the
// rendering outside a Scope) is the new-york-v4 element; Render picks the
// alternative for the style in scope.
package styled

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	h "shadcn.go/libhtml"
)

// Node renders def(props), or alt(st, props) inside a
// style.Scope(st, ...). Both are computed at render time, so props (the
// caller's attributes and children) can still be changed until then.
type Node struct {
	props *h.Element
	def   func(props *h.Element) h.Node
	alt   func(st style.Style, props *h.Element) h.Node
}

// New returns a Node of fixed renderings. alt may return nil to fall back
// to def.
func New(def h.Node, alt func(style.Style) h.Node) *Node {
	return &Node{
		def: func(*h.Element) h.Node { return def },
		alt: func(st style.Style, _ *h.Element) h.Node { return alt(st) },
	}
}

// Lazy returns a Node rendering def(props) or alt(st, props); Props and
// Children change props until the node renders. alt may return nil to
// fall back to def.
func Lazy(props *h.Element, def func(props *h.Element) h.Node, alt func(st style.Style, props *h.Element) h.Node) *Node {
	return &Node{props: props, def: def, alt: alt}
}

// Spread is {...props} on a base element: the new-york-v4 rendering is
// h.Spread(base, props) and, in a style, h.Spread(alt(st), props) (classes
// merged with lib.Cn). alt may return nil to fall back to base.
func Spread(base, props *h.Element, alt func(style.Style) *h.Element) *Node {
	return Lazy(props,
		func(p *h.Element) h.Node { return h.Spread(base, p, lib.Cn) },
		func(st style.Style, p *h.Element) h.Node {
			if a := alt(st); a != nil {
				return h.Spread(a, p, lib.Cn)
			}
			return nil
		})
}

// Props returns the caller's props element (nil for New nodes), to add
// attributes after construction: n.Props().Bool("disabled", true).
func (n *Node) Props() *h.Element { return n.props }

// Children appends children (and attributes) to a Lazy node's props.
func (n *Node) Children(children ...any) *Node {
	if n.props != nil {
		n.props.Children(children...)
	}
	return n
}

func (n *Node) defNode() h.Node { return n.def(n.props) }

func (n *Node) altNode(st style.Style) h.Node {
	if n.alt == nil {
		return nil
	}
	return n.alt(st, n.props)
}

// Element returns the new-york-v4 element.
func (n *Node) Element() *h.Element {
	if e, ok := n.defNode().(h.Elementer); ok {
		return e.Element()
	}
	return nil
}

// Render renders the node in the style in scope, if any.
func (n *Node) Render(w io.Writer) error {
	if st, ok := style.FromWriter(w); ok {
		if a := n.altNode(st); a != nil {
			return a.Render(w)
		}
	}
	return n.defNode().Render(w)
}

// ElementIn returns the element the node renders as inside w
// (h.ElementerIn), for components slotting it as an asChild child.
func (n *Node) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		if a := n.altNode(st); a != nil {
			if e := h.ResolveIn(w, a); e != nil {
				return e
			}
		}
	}
	return h.ResolveIn(w, n.defNode())
}

// Elem exposes the node's children to tree.Walk (compound components
// handing context to their parts): the props of a Lazy node, or the
// new-york-v4 element otherwise.
func (n *Node) Elem() *h.Element {
	if n.props != nil {
		return n.props
	}
	switch d := n.defNode().(type) {
	case *h.Element:
		return d
	case interface{ Elem() *h.Element }:
		return d.Elem()
	case h.Elementer:
		return d.Element()
	}
	return nil
}

// Slot is an element with a data-slot and the class of site, a
// file.Cn("...", style.Arg) site (className is merged by Spread).
func Slot(tag, slot string, site *style.Site, attrs ...any) func(style.Style) *h.Element {
	return func(st style.Style) *h.Element {
		return h.El(tag, append([]any{h.Data("slot", slot), h.Class(site.Class(st, ""))}, attrs...)...)
	}
}
