// Package combobox ports the upstream combobox-* examples (a Command in a
// Popover, or in a Drawer on small screens) and the parity fixtures of the
// Base UI Combobox (components/ui/combobox).
package combobox

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/command"
	"shadcn.go/components/ui/drawer"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// combobox.js keeps the examples' React state: the picked value is shown
// by the trigger and marked in the list the popover mounts next time, and
// picking closes the popover.
//
//go:embed combobox.js
var script string

func init() {
	runtime.Module("example-combobox", script)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "combobox", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("combobox-demo", Demo),
		ex("combobox-popover", PopoverExample),
		ex("combobox-responsive", Responsive),
		ex("combobox-dropdown-menu", DropdownMenu),
	)
	registerFixtures(ex)
}

type option struct{ value, label string }

var frameworks = []option{
	{"next.js", "Next.js"},
	{"sveltekit", "SvelteKit"},
	{"nuxt.js", "Nuxt.js"},
	{"remix", "Remix"},
	{"astro", "Astro"},
}

var statuses = []option{
	{"backlog", "Backlog"},
	{"todo", "Todo"},
	{"in progress", "In Progress"},
	{"done", "Done"},
	{"canceled", "Canceled"},
}

// options serializes the value/label pairs for combobox.js.
func options(list []option) string {
	s := ""
	for i, o := range list {
		if i > 0 {
			s += "\n"
		}
		s += o.value + "\t" + o.label
	}
	return s
}

func Demo() h.Node {
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline,
			h.Role("combobox"),
			h.Aria("expanded", "false"),
			h.Class("w-[200px] justify-between"),
			h.Data("example-combobox", "demo"),
			h.Data("example-combobox-options", options(frameworks)),
			"Select framework...",
			lucide.ChevronsUpDown(h.Class("opacity-50")),
		)),
		popover.Content(h.Class("w-[200px] p-0"),
			command.Command(
				command.Input(h.Placeholder("Search framework..."), h.Class("h-9")),
				command.List(
					command.Empty("No framework found."),
					command.Group(
						h.Map(frameworks, func(f option) h.Node {
							return command.Item(command.Value(f.value),
								f.label,
								lucide.Check(h.Class(lib.Cn("ml-auto", "opacity-0"))),
							)
						}),
					),
				),
			),
		),
	)
}

func PopoverExample() h.Node {
	// The content's template follows the row: inside it, it would be the
	// last child and take space-x-4's margin off the trigger.
	return popover.Popover(
		h.Div(h.Class("flex items-center space-x-4"),
			h.P(h.Class("text-sm text-muted-foreground"), "Status"),
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline,
				h.Class("w-[150px] justify-start"),
				h.Data("example-combobox", "status"),
				h.Data("example-combobox-options", options(statuses)),
				"+ Set status",
			)),
		),
		popover.Content(h.Class("p-0"), popover.SideRight, popover.AlignStart,
			command.Command(
				command.Input(h.Placeholder("Change status...")),
				statusItems(),
			),
		),
	)
}

func statusItems() h.Node {
	return command.List(
		command.Empty("No results found."),
		command.Group(
			h.Map(statuses, func(s option) h.Node {
				return command.Item(command.Value(s.value), s.label)
			}),
		),
	)
}

func statusList() h.Node {
	return command.Command(
		command.Input(h.Placeholder("Filter status...")),
		statusItems(),
	)
}

// Responsive is combobox-responsive: a Popover on screens at least 768px
// wide, a Drawer otherwise. The server renders the Popover version; the
// Drawer version waits in a template that combobox.js swaps in when the
// media query stops matching.
func Responsive() h.Node {
	trigger := func() h.Node {
		return button.Button(button.VariantOutline,
			h.Class("w-[150px] justify-start"),
			h.Data("example-combobox", "status"),
			h.Data("example-combobox-options", options(statuses)),
			"+ Set status",
		)
	}
	return h.Group(
		popover.Popover(
			popover.Trigger(lib.AsChild, trigger()),
			popover.Content(h.Class("w-[200px] p-0"), popover.AlignStart, statusList()),
		),
		h.Template(h.Data("parity-ignore", ""), h.Data("example-combobox-media", "(min-width: 768px)"),
			drawer.Drawer(
				drawer.Trigger(lib.AsChild, trigger()),
				drawer.Content(
					h.Div(h.Class("mt-4 border-t"), statusList()),
				),
			),
		),
	)
}

var labels = []string{"feature", "bug", "enhancement", "documentation", "design", "question", "maintenance"}

// DropdownMenu is combobox-dropdown-menu: a Command filtering labels in a
// dropdown submenu; picking one sets the label and closes the menu.
func DropdownMenu() h.Node {
	return h.Div(h.Class("flex w-full flex-col items-start justify-between rounded-md border px-4 py-3 sm:flex-row sm:items-center"),
		h.P(h.Class("text-sm leading-none font-medium"),
			h.Span(h.Class("mr-2 rounded-lg bg-primary px-2 py-1 text-xs text-primary-foreground"), h.Data("example-combobox-label", ""), "feature"),
			h.Span(h.Class("text-muted-foreground"), "Create a new project"),
		),
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeSm, lucide.Ellipsis())),
			dm.Content(dm.AlignEnd, h.Class("w-[200px]"),
				dm.Label("Actions"),
				dm.Group(
					dm.Item("Assign to..."),
					dm.Item("Set due date..."),
					dm.Separator(),
					dm.Sub(
						dm.SubTrigger("Apply label"),
						dm.SubContent(h.Class("p-0"),
							command.Command(
								command.Input(h.Placeholder("Filter label..."), h.Data("example-combobox-autofocus", ""), h.Class("h-9")),
								command.List(
									command.Empty("No label found."),
									command.Group(
										h.Map(labels, func(l string) h.Node {
											return command.Item(command.Value(l), l)
										}),
									),
								),
							),
						),
					),
					dm.Separator(),
					dm.Item(h.Class("text-red-600"), "Delete", dm.Shortcut("⌘⌫")),
				),
			),
		),
	)
}
