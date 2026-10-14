// Package scoped lets a new-york-v4 component render the markup of a
// registry/bases/radix style (radix-nova, ...) when it is rendered inside
// style.Scope, and its unchanged new-york-v4 markup otherwise.
//
// A component returns a *Node. Its Element is the new-york-v4 element;
// Render and ElementIn (h.ElementerIn, used by h.SlotIn / h.ResolveIn for
// asChild) pick the style in scope.
package scoped

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	h "shadcn.go/libhtml"
)

// Node renders def(), or alt(st) inside a style.Scope(st, ...).
type Node struct {
	def   func() *h.Element
	alt   func(style.Style) *h.Element
	props *h.Element
}

// Props returns the user props element of a Spread node (nil otherwise),
// for setting attributes after construction.
func (n *Node) Props() *h.Element { return n.props }

// Children appends children to the props of a Spread node.
func (n *Node) Children(children ...any) *Node {
	n.props.Children(children...)
	return n
}

// New returns a Node.
func New(def func() *h.Element, alt func(style.Style) *h.Element) *Node {
	return &Node{def: def, alt: alt}
}

// Element returns the new-york-v4 element.
func (n *Node) Element() *h.Element { return n.def() }

// ElementIn returns the element in the style in scope in w, if any.
func (n *Node) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		return n.alt(st)
	}
	return n.def()
}

// Render renders the node in the style in scope, if any.
func (n *Node) Render(w io.Writer) error { return n.ElementIn(w).Render(w) }

// Spread is {...props} on a base element: h.Spread(base, props) outside a
// style, h.Spread(alt(st), props) in style st (classes merged with lib.Cn).
func Spread(base *h.Element, props *h.Element, alt func(style.Style) *h.Element) *Node {
	n := New(
		func() *h.Element { return h.Spread(base, props, lib.Cn) },
		func(st style.Style) *h.Element { return h.Spread(alt(st), props, lib.Cn) },
	)
	n.props = props
	return n
}

// Slot returns a function building <tag data-slot=slot class=site(st)>, with
// extra attrs: the alt of a Spread for a site file.Cn("...", style.Arg).
func Slot(tag, slot string, site *style.Site, attrs ...any) func(style.Style) *h.Element {
	return func(st style.Style) *h.Element {
		return h.El(tag, append([]any{h.Data("slot", slot), h.Class(site.Class(st, ""))}, attrs...)...)
	}
}
