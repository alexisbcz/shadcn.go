// Package spinner is the shadcn/ui Spinner: an indicator that shows that
// something is loading.
//
//	spinner.Spinner(h.Class("size-8"))
package spinner

import (
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

// Builder is returned by Spinner.
type Builder struct {
	h.Builder[Builder]
}

// Spinner returns a spinning Loader2 icon with role="status".
func Spinner(children ...any) *Builder {
	b := &Builder{}
	b.Builder = h.NewBuilder(b, "svg")
	return b.Children(children...)
}

// Element resolves the spinner to its <svg>.
func (b *Builder) Element() *h.Element {
	base := lucide.LoaderCircle(h.Role("status"), h.Aria("label", "Loading"), h.Class("size-4 animate-spin"))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the spinner.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }
