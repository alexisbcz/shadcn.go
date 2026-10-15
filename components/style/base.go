package style

import (
	"io"

	h "shadcn.go/libhtml"
)

// Base names the primitive library a preview mirrors ("radix", "base",
// "aria"): the shadcn styles are shared, but a few components lay out
// their parts differently per library (React Aria's breadcrumb items hold
// their separator, for one).
type Base string

// Bases.
const (
	BaseRadix Base = "radix"
	BaseUI    Base = "base"
	BaseAria  Base = "aria"
)

type baseKey struct{}

// WithBase renders children as the given library lays them out.
func WithBase(b Base, children ...any) h.Node { return h.WithValue(baseKey{}, b, children...) }

// BaseFromWriter returns the library in scope (BaseRadix by default).
func BaseFromWriter(w io.Writer) Base {
	if v, ok := h.Lookup(w, baseKey{}); ok {
		return v.(Base)
	}
	return BaseRadix
}
