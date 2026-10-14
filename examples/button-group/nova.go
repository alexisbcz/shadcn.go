package buttongroup

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/field"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the button-group page
// (ui/apps/v4/examples/radix/button-group-*.tsx) whose content differs from
// the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "button-group", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("button-group-demo", func() h.Node { return demo("w-40") }),
		nova("button-group-dropdown", func() h.Node { return dropdown("w-44") }),
		nova("button-group-nested", NovaNested),
		nova("button-group-popover", NovaPopover),
	)
}

func NovaNested() h.Node {
	return buttongroup.ButtonGroup(
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeIcon, lucide.Plus()),
		),
		buttongroup.ButtonGroup(
			ig.InputGroup(
				ig.Input(h.Placeholder("Send a message...")),
				tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild, ig.Addon(ig.AlignInlineEnd, lucide.AudioLines())),
					tooltip.Content("Voice Mode"),
				),
			),
		),
	)
}

func NovaPopover() h.Node {
	return buttongroup.ButtonGroup(
		button.Button(button.VariantOutline, lucide.Bot(), " Copilot"),
		popover.Popover(
			// Content first: see Popover.
			popover.Content(popover.AlignEnd, h.Class("rounded-xl text-sm"),
				popover.Header(
					popover.Title("Start a new task with Copilot"),
					popover.Description("Describe your task in natural language."),
				),
				field.Field(
					field.Label(h.For("task"), h.Class("sr-only"), "Task Description"),
					textarea.Textarea(h.ID("task"), h.Placeholder("I need to..."), h.Class("resize-none")),
					field.Description("Copilot will open a pull request for review."),
				),
			),
			popover.Trigger(lib.AsChild,
				button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Open Popover"), lucide.ChevronDown()),
			),
		),
	)
}
