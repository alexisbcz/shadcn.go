// Package textarea is the shadcn/ui Textarea: a multi-line text field that
// grows with its content (field-sizing: content).
//
//	textarea.Textarea(h.Placeholder("Type your message here."))
package textarea

import (
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

var base = h.Textarea(h.Data("slot", "textarea"), h.Class(lib.Cn(
	"flex field-sizing-content min-h-16 w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-xs transition-[color,box-shadow] outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:aria-invalid:ring-destructive/40",
)))

// Textarea returns a <textarea>.
func Textarea(children ...any) *h.Element { return h.Spread(base, h.Textarea(children...), lib.Cn) }
