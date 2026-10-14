// Package label is the shadcn/ui Label (Radix Label): an accessible label
// associated with a control.
//
//	label.Label(h.For("email"), "Email")
//
// Like Radix, double-clicking the label does not select its text (see
// label.js).
package label

import (
	_ "embed"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed label.js
var script string

func init() { runtime.Module("label", script) }

var base = h.Label(h.Data("slot", "label"), h.Class(lib.Cn(
	"flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
)))

// Label returns a <label>.
func Label(children ...any) *h.Element { return h.Spread(base, h.Label(children...), lib.Cn) }
