// Package skeleton is the shadcn/ui Skeleton: a placeholder shown while
// content loads.
//
//	skeleton.Skeleton(h.Class("h-4 w-[250px]"))
package skeleton

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	h "shadcn.go/libhtml"
)

var base = h.Div(h.Data("slot", "skeleton"), h.Class(lib.Cn("animate-pulse rounded-md bg-accent")))

// Skeleton returns a pulsing placeholder <div>.
func Skeleton(children ...any) *scoped.Node {
	return scoped.Spread(base, h.Div(children...), scoped.Slot("div", "skeleton", baseSkeleton))
}

// registry/bases/radix/ui/skeleton.tsx, used inside a style.Scope.
var baseSkeleton = style.Cn("cn-skeleton animate-pulse", style.Arg)
