// Package buttongroup ports the upstream button-group-* examples.
package buttongroup

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/input"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// voice.js holds button-group-input-group's voice mode state.
//
//go:embed voice.js
var voiceScript string

func init() {
	runtime.Module("example-button-group-voice", voiceScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "button-group", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("button-group-input", Input),
		ex("button-group-nested", Nested),
		ex("button-group-orientation", Orientation),
		ex("button-group-separator", Separator),
		ex("button-group-size", Size),
		ex("button-group-split", Split),
		ex("button-group-popover", Popover),
		ex("button-group-input-group", InputGroup),
	)
}

func Input() h.Node {
	return buttongroup.ButtonGroup(
		input.Input(h.Placeholder("Search...")),
		button.Button(button.VariantOutline, h.Aria("label", "Search"), lucide.Search()),
	)
}

func Nested() h.Node {
	num := func(s string) h.Node { return button.Button(button.VariantOutline, button.SizeSm, s) }
	return buttongroup.ButtonGroup(
		buttongroup.ButtonGroup(num("1"), num("2"), num("3"), num("4"), num("5")),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", "Previous"), lucide.ArrowLeft()),
			button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", "Next"), lucide.ArrowRight()),
		),
	)
}

func Orientation() h.Node {
	return buttongroup.ButtonGroup(buttongroup.Vertical, h.Aria("label", "Media controls"), h.Class("h-fit"),
		button.Button(button.VariantOutline, button.SizeIcon, lucide.Plus()),
		button.Button(button.VariantOutline, button.SizeIcon, lucide.Minus()),
	)
}

func Separator() h.Node {
	return buttongroup.ButtonGroup(
		button.Button(button.VariantSecondary, button.SizeSm, "Copy"),
		buttongroup.Separator(),
		button.Button(button.VariantSecondary, button.SizeSm, "Paste"),
	)
}

func Size() h.Node {
	row := func(size, icon button.Size, label string) h.Node {
		return buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, size, label),
			button.Button(button.VariantOutline, size, "Button"),
			button.Button(button.VariantOutline, size, "Group"),
			button.Button(button.VariantOutline, icon, lucide.Plus()),
		)
	}
	return h.Div(h.Class("flex flex-col items-start gap-8"),
		row(button.SizeSm, button.SizeIconSm, "Small"),
		row(button.SizeDefault, button.SizeIcon, "Default"),
		row(button.SizeLg, button.SizeIconLg, "Large"),
	)
}

func Split() h.Node {
	return buttongroup.ButtonGroup(
		button.Button(button.VariantSecondary, "Button"),
		buttongroup.Separator(),
		button.Button(button.SizeIcon, button.VariantSecondary, tabler.IconPlus()),
	)
}

func Popover() h.Node {
	return buttongroup.ButtonGroup(
		button.Button(button.VariantOutline, lucide.Bot(), " Copilot"),
		popover.Popover(
			// The Go popover renders its closed content as a <template> in
			// place; putting it before the trigger keeps the trigger the
			// group's :last-child, as it is in React (where Popover renders
			// nothing of its own).
			popover.Content(popover.AlignEnd, h.Class("rounded-xl p-0 text-sm"),
				h.Div(h.Class("px-4 py-3"), h.Div(h.Class("text-sm font-medium"), "Agent Tasks")),
				separator.Separator(),
				h.Div(h.Class("p-4 text-sm *:[p:not(:last-child)]:mb-2"),
					textarea.Textarea(h.Placeholder("Describe your task in natural language."), h.Class("mb-4 resize-none")),
					h.P(h.Class("font-medium"), "Start a new task with Copilot"),
					h.P(h.Class("text-muted-foreground"), "Describe your task in natural language. Copilot will work in the background and open a pull request for your review."),
				),
			),
			popover.Trigger(lib.AsChild,
				button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Open Popover"), lucide.ChevronDown()),
			),
		),
	)
}

func InputGroup() h.Node {
	return buttongroup.ButtonGroup(h.Class("[--radius:9999rem]"),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeIcon, lucide.Plus()),
		),
		buttongroup.ButtonGroup(
			ig.InputGroup(
				ig.Input(h.Placeholder("Send a message...")),
				ig.Addon(ig.AlignInlineEnd,
					tooltip.Tooltip(
						tooltip.Trigger(lib.AsChild,
							ig.Button(ig.ButtonSizeIconXs, h.Data("active", "false"), h.Aria("pressed", "false"),
								h.Class("data-[active=true]:bg-orange-100 data-[active=true]:text-orange-700 dark:data-[active=true]:bg-orange-800 dark:data-[active=true]:text-orange-100"),
								lucide.AudioLines(),
							),
						),
						tooltip.Content("Voice Mode"),
					),
				),
			),
		),
	)
}
