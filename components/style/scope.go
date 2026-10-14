package style

import (
	"io"

	h "shadcn.go/libhtml"
)

type scopeKey struct{}

// Scope renders children in style s: style-aware components below it render
// that style's markup and classes. Pages can mix styles, e.g. the docs site
// chrome in new-york-v4 with radix-nova previews.
func Scope(s Style, children ...any) h.Node { return h.WithValue(scopeKey{}, s, children...) }

// FromWriter returns the style set by an enclosing Scope, if any, for use in
// a component's Render method (w is the writer libhtml passes to it).
func FromWriter(w io.Writer) (Style, bool) {
	v, ok := h.Lookup(w, scopeKey{})
	if !ok {
		return "", false
	}
	s, ok := v.(Style)
	return s, ok
}
