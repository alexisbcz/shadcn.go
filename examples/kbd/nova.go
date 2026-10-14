package kbd

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/kbd"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the kbd page (ui/apps/v4/examples/radix/kbd-*.tsx) whose
// content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "kbd", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("kbd-button", NovaButton),
		nova("kbd-input-group", NovaInputGroup),
		nova("kbd-rtl", NovaRtl),
		nova("kbd-tooltip", NovaTooltip),
	)
}

func NovaButton() h.Node {
	return button.Button(button.VariantOutline, "Accept", " ", kbd.Kbd(h.Data("icon", "inline-end"), h.Class("translate-x-0.5"), "⏎"))
}

func NovaInputGroup() h.Node {
	return h.Div(h.Class("flex w-full max-w-xs flex-col gap-6"),
		ig.InputGroup(
			ig.Input(h.Placeholder("Search...")),
			ig.Addon(lucide.Search()),
			ig.Addon(ig.AlignInlineEnd, kbd.Kbd("⌘"), kbd.Kbd("K")),
		),
	)
}

func NovaRtl() h.Node {
	return style.ScopeRTL(h.Div(h.Class("flex flex-col items-center gap-4"), h.Dir("rtl"),
		kbd.Group(kbd.Kbd("⌘"), kbd.Kbd("⇧"), kbd.Kbd("⌥"), kbd.Kbd("⌃")),
		kbd.Group(kbd.Kbd("Ctrl"), h.Span("+"), kbd.Kbd("B")),
	))
}

func NovaTooltip() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-4"),
		buttongroup.ButtonGroup(
			tooltip.Tooltip(
				tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Save")),
				tooltip.Content("Save Changes ", kbd.Kbd("S")),
			),
			tooltip.Tooltip(
				tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Print")),
				tooltip.Content("Print Document", " ", kbd.Group(kbd.Kbd("Ctrl"), kbd.Kbd("P"))),
			),
		),
	)
}
