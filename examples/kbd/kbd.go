// Package kbd ports the upstream kbd-* examples.
package kbd

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/kbd"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "kbd", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("kbd-button", Button),
		ex("kbd-demo", Demo),
		ex("kbd-group", Group),
	)
}

func Button() h.Node {
	return h.Div(h.Class("flex flex-wrap items-center gap-4"),
		button.Button(button.VariantOutline, button.SizeSm, h.Class("pr-2"), "Accept ", kbd.Kbd("⏎")),
		button.Button(button.VariantOutline, button.SizeSm, h.Class("pr-2"), "Cancel ", kbd.Kbd("Esc")),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-4"),
		kbd.Group(kbd.Kbd("⌘"), kbd.Kbd("⇧"), kbd.Kbd("⌥"), kbd.Kbd("⌃")),
		kbd.Group(kbd.Kbd("Ctrl"), h.Span("+"), kbd.Kbd("B")),
	)
}

func Group() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-4"),
		h.P(h.Class("text-sm text-muted-foreground"),
			"Use", " ",
			kbd.Group(kbd.Kbd("Ctrl + B"), kbd.Kbd("Ctrl + K")),
			" ", "to open the command palette",
		),
	)
}
