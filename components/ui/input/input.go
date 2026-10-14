// Package input is the shadcn/ui Input: a form input field.
//
//	input.Input(h.Type("email"), h.Placeholder("Email"))
package input

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	h "shadcn.go/libhtml"
)

var base = h.Input(h.Data("slot", "input"), h.Class(lib.Cn(
	"h-9 w-full min-w-0 rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none selection:bg-primary selection:text-primary-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30",
	"focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
	"aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40",
)))

// Input returns an <input>. Pass its type, placeholder, etc. as attributes.
func Input(children ...any) *scoped.Node {
	return scoped.Spread(base, h.Input(children...), scoped.Slot("input", "input", baseInput))
}

// registry/bases/radix/ui/input.tsx, used inside a style.Scope.
var baseInput = style.Cn("cn-input w-full min-w-0 outline-none file:inline-flex file:border-0 file:bg-transparent file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50", style.Arg)
