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

// Node renders def, or alt(st) inside a style.Scope(st, ...).
type Node struct {
	def h.Node
	alt func(style.Style) h.Node

	// Spread nodes: def is h.Spread(base, props), alt h.Spread(altBase(st), props).
	base, props *h.Element
	altBase     func(style.Style) *h.Element
}

// New returns a Node. alt may return nil to fall back to def.
func New(def h.Node, alt func(style.Style) h.Node) *Node {
	return &Node{def: def, alt: alt}
}

// Spread is {...props} on a base element: the new-york-v4 rendering is
// h.Spread(base, props) and, in a style, h.Spread(alt(st), props) (classes
// merged with lib.Cn). alt may return nil to fall back to base.
func Spread(base, props *h.Element, alt func(style.Style) *h.Element) *Node {
	return &Node{base: base, props: props, altBase: alt}
}

// Children appends children (and attributes) to a Spread node's props.
func (n *Node) Children(children ...any) *Node {
	if n.props != nil {
		n.props.Children(children...)
	}
	return n
}

func (n *Node) defNode() h.Node {
	if n.props != nil {
		return h.Spread(n.base, n.props, lib.Cn)
	}
	return n.def
}

func (n *Node) altNode(st style.Style) h.Node {
	if n.props != nil {
		if n.altBase == nil {
			return nil
		}
		if a := n.altBase(st); a != nil {
			return h.Spread(a, n.props, lib.Cn)
		}
		return nil
	}
	if n.alt == nil {
		return nil
	}
	return n.alt(st)
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

// Slot is an element with a data-slot and the class of site, a
// file.Cn("...", style.Arg) site (className is merged by Spread).
func Slot(tag, slot string, site *style.Site, attrs ...any) func(style.Style) *h.Element {
	return func(st style.Style) *h.Element {
		return h.El(tag, append([]any{h.Data("slot", slot), h.Class(site.Class(st, ""))}, attrs...)...)
	}
}

// StyledElement returns the element of style st (scoped.Styled), for
// components slotting this node as an asChild child.
func (n *Node) StyledElement(st style.Style) *h.Element {
	if st != "" {
		if a := n.altNode(st); a != nil {
			switch e := a.(type) {
			case interface{ StyledElement(style.Style) *h.Element }:
				return e.StyledElement(st)
			case h.Elementer:
				return e.Element()
			}
		}
	}
	return n.Element()
}

// Elem exposes the node's children to tree.Walk (compound components
// handing context to their parts): the props of a Spread node, or the
// new-york-v4 element otherwise.
func (n *Node) Elem() *h.Element {
	if n.props != nil {
		return n.props
	}
	switch d := n.def.(type) {
	case *h.Element:
		return d
	case interface{ Elem() *h.Element }:
		return d.Elem()
	case h.Elementer:
		return d.Element()
	}
	return nil
}
