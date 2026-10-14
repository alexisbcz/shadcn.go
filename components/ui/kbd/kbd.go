// Package kbd is the shadcn/ui Kbd: displays keyboard input.
//
//	kbd.Group(kbd.Kbd("Ctrl"), kbd.Kbd("B"))
package kbd

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
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
func Kbd(children ...any) *scoped.Node {
	return scoped.Spread(kbdBase, h.Kbd(children...), scoped.Slot("kbd", "kbd", baseKbd))
}

// registry/bases/radix/ui/kbd.tsx, used inside a style.Scope.
var (
	baseFile  = style.NewFile("ui/kbd.tsx")
	baseKbd   = baseFile.Cn("cn-kbd pointer-events-none inline-flex items-center justify-center select-none", style.Arg)
	baseGroup = baseFile.Cn("cn-kbd-group inline-flex items-center", style.Arg)
)

// Group returns a <kbd> grouping several keys (upstream KbdGroup, which
// renders a <kbd> too).
func Group(children ...any) *scoped.Node {
	return scoped.Spread(groupBase, h.Kbd(children...), scoped.Slot("kbd", "kbd-group", baseGroup))
}
