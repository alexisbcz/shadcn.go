// Package skeleton is the shadcn/ui Skeleton: a placeholder shown while
// content loads.
//
//	skeleton.Skeleton(h.Class("h-4 w-[250px]"))
package skeleton

import (
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

var base = h.Div(h.Data("slot", "skeleton"), h.Class(lib.Cn("animate-pulse rounded-md bg-accent")))

// Skeleton returns a pulsing placeholder <div>.
func Skeleton(children ...any) *h.Element { return h.Spread(base, h.Div(children...), lib.Cn) }
