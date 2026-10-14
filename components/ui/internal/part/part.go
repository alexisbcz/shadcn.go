// Package part is a component part whose element depends on the style in
// scope (style.Scope): its function builds the element for a style ("" for
// new-york-v4) from the caller's props.
package part

import (
	"io"

	"shadcn.go/cn"
	"shadcn.go/components/style"
	h "shadcn.go/libhtml"
)

// Part is returned by New.
type Part struct {
	h.Builder[Part]
	fn func(st style.Style, rtl bool, props *h.Element) *h.Element
}

// New returns a part collecting props and children on a <tag> element; fn
// builds the element to render for a style (and its ui-rtl variant when rtl
// is true, inside style.ScopeRTL) from them.
func New(tag string, fn func(st style.Style, rtl bool, props *h.Element) *h.Element, children ...any) *Part {
	p := &Part{fn: fn}
	p.Builder = h.NewBuilder(p, tag)
	return p.Children(children...)
}

// Element is the new-york-v4 element.
func (p *Part) Element() *h.Element { return p.fn("", false, p.Elem()) }

// ElementIn is the element in the style of an enclosing style.Scope.
func (p *Part) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return p.fn(st, RTL(w), p.Elem())
}

// Render renders ElementIn.
func (p *Part) Render(w io.Writer) error { return p.ElementIn(w).Render(w) }

// Class is site.Class(st, values...), in the ui-rtl variant of the
// component when rtl is true (inside style.ScopeRTL): the site's literal
// classes mapped by style.RTLClass, the values (the caller's className,
// cva results of other files) unchanged, as in upstream's
// styles/<style>/ui-rtl sources.
func Class(rtl bool, site *style.Site, st style.Style, values ...string) string {
	if !rtl {
		return site.Class(st, values...)
	}
	args := site.Args(st)
	out := make([]string, 0, len(args))
	i := 0
	for _, a := range args {
		if s, ok := a.(string); ok {
			out = append(out, style.RTLClass(s))
			continue
		}
		out = append(out, values[i])
		i++
	}
	return cn.Cn(out...)
}

// RTL reports whether w renders inside style.ScopeRTL.
func RTL(w io.Writer) bool { return w != nil && style.RTLFromWriter(w) }
