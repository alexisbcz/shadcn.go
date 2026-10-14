// Package textarea is the shadcn/ui Textarea: a multi-line text field that
// grows with its content (field-sizing: content).
//
//	textarea.Textarea(h.Placeholder("Type your message here."))
package textarea

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	h "shadcn.go/libhtml"
)

var base = h.Textarea(h.Data("slot", "textarea"), h.Class(lib.Cn(
	"flex field-sizing-content min-h-16 w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-xs transition-[color,box-shadow] outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:aria-invalid:ring-destructive/40",
)))

// Textarea returns a <textarea>.
func Textarea(children ...any) *scoped.Node {
	return scoped.Spread(base, h.Textarea(children...), scoped.Slot("textarea", "textarea", baseTextarea))
}

// registry/bases/radix/ui/textarea.tsx, used inside a style.Scope.
var baseTextarea = style.Cn("cn-textarea flex field-sizing-content min-h-16 w-full outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50", style.Arg)
