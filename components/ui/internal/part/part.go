// Package part is a component part whose element depends on the style in
// scope (style.Scope): its function builds the element for a style ("" for
// new-york-v4) from the caller's props.
package part

import (
	"io"

	"shadcn.go/components/style"
	h "shadcn.go/libhtml"
)

// Part is returned by New.
type Part struct {
	h.Builder[Part]
	fn func(st style.Style, props *h.Element) *h.Element
}

// New returns a part collecting props and children on a <tag> element; fn
// builds the element to render for a style from them.
func New(tag string, fn func(st style.Style, props *h.Element) *h.Element, children ...any) *Part {
	p := &Part{fn: fn}
	p.Builder = h.NewBuilder(p, tag)
	return p.Children(children...)
}

// Element is the new-york-v4 element.
func (p *Part) Element() *h.Element { return p.fn("", p.Elem()) }

// ElementIn is the element in the style of an enclosing style.Scope.
func (p *Part) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return p.fn(st, p.Elem())
}

// Render renders ElementIn.
func (p *Part) Render(w io.Writer) error { return p.ElementIn(w).Render(w) }
