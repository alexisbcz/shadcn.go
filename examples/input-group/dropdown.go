package inputgroup

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "input-group", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("input-group-demo", Demo),
		ex("input-group-dropdown", Dropdown),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(
			ig.Input(h.Placeholder("Search...")),
			ig.Addon(lucide.Search()),
			ig.Addon(ig.AlignInlineEnd, "12 results"),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("example.com"), h.Class("pl-1!")),
			ig.Addon(ig.Text("https://")),
			ig.Addon(ig.AlignInlineEnd,
				tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild, ig.Button(h.Class("rounded-full"), ig.ButtonSizeIconXs, tabler.IconInfoCircle())),
					tooltip.Content("This is content in a tooltip."),
				),
			),
		),
		ig.InputGroup(
			ig.Textarea(h.Placeholder("Ask, Search or Chat...")),
			ig.Addon(ig.AlignBlockEnd,
				ig.Button(button.VariantOutline, h.Class("rounded-full"), ig.ButtonSizeIconXs, tabler.IconPlus()),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild, ig.Button(button.VariantGhost, "Auto")),
					dm.Content(dm.SideTop, dm.AlignStart, h.Class("[--radius:0.95rem]"),
						dm.Item("Auto"), dm.Item("Agent"), dm.Item("Manual"),
					),
				),
				ig.Text(h.Class("ml-auto"), "52% used"),
				separator.Separator(separator.Vertical, h.Class("h-4!")),
				ig.Button(button.VariantDefault, h.Class("rounded-full"), ig.ButtonSizeIconXs, h.Disabled(),
					lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send"),
				),
			),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("@shadcn")),
			ig.Addon(ig.AlignInlineEnd,
				h.Div(h.Class("flex size-4 items-center justify-center rounded-full bg-primary text-primary-foreground"),
					tabler.IconCheck(h.Class("size-3")),
				),
			),
		),
	)
}

func Dropdown() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-4"),
		ig.InputGroup(
			ig.Input(h.Placeholder("Enter file name")),
			ig.Addon(ig.AlignInlineEnd,
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild, ig.Button(button.VariantGhost, h.Aria("label", "More"), ig.ButtonSizeIconXs, lucide.Ellipsis())),
					dm.Content(dm.AlignEnd,
						dm.Item("Settings"), dm.Item("Copy path"), dm.Item("Open location"),
					),
				),
			),
		),
		ig.InputGroup(h.Class("[--radius:1rem]"),
			ig.Input(h.Placeholder("Enter search query")),
			ig.Addon(ig.AlignInlineEnd,
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild, ig.Button(button.VariantGhost, h.Class("pr-1.5! text-xs"),
						"Search In... ", lucide.ChevronDown(h.Class("size-3")))),
					dm.Content(dm.AlignEnd, h.Class("[--radius:0.95rem]"),
						dm.Item("Documentation"), dm.Item("Blog Posts"), dm.Item("Changelog"),
					),
				),
			),
		),
	)
}
