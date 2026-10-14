// Package kbd is the shadcn/ui Kbd: displays keyboard input.
//
//	kbd.Group(kbd.Kbd("Ctrl"), kbd.Kbd("B"))
package kbd

import (
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

var (
	kbdBase = h.Kbd(h.Data("slot", "kbd"), h.Class(lib.Cn(
		"pointer-events-none inline-flex h-5 w-fit min-w-5 items-center justify-center gap-1 rounded-sm bg-muted px-1 font-sans text-xs font-medium text-muted-foreground select-none",
		"[&_svg:not([class*='size-'])]:size-3",
		"[[data-slot=tooltip-content]_&]:bg-background/20 [[data-slot=tooltip-content]_&]:text-background dark:[[data-slot=tooltip-content]_&]:bg-background/10",
	)))
	groupBase = h.Kbd(h.Data("slot", "kbd-group"), h.Class(lib.Cn("inline-flex items-center gap-1")))
)

// Kbd returns a <kbd> key.
func Kbd(children ...any) *h.Element { return h.Spread(kbdBase, h.Kbd(children...), lib.Cn) }

// Group returns a <kbd> grouping several keys (upstream KbdGroup, which
// renders a <kbd> too).
func Group(children ...any) *h.Element { return h.Spread(groupBase, h.Kbd(children...), lib.Cn) }
