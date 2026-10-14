// Package command ports the upstream command-* examples.
package command

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/command"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// dialog.js toggles command-dialog's palette with ⌘J / Ctrl+J, like the
// upstream example's keydown listener.
//
//go:embed dialog.js
var dialogScript string

func init() {
	runtime.Module("example-command-dialog", dialogScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "command", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("command-demo", Demo),
		ex("command-dialog", DialogDemo),
	)
}

func Demo() h.Node {
	return command.Command(h.Class("rounded-lg border shadow-md md:min-w-[450px]"),
		command.Input(h.Placeholder("Type a command or search...")),
		command.List(
			command.Empty("No results found."),
			command.Group(command.Heading("Suggestions"),
				command.Item(lucide.Calendar(), h.Span("Calendar")),
				command.Item(lucide.Smile(), h.Span("Search Emoji")),
				command.Item(command.Disabled(true), lucide.Calculator(), h.Span("Calculator")),
			),
			command.Separator(),
			command.Group(command.Heading("Settings"),
				command.Item(lucide.User(), h.Span("Profile"), command.Shortcut("⌘P")),
				command.Item(lucide.CreditCard(), h.Span("Billing"), command.Shortcut("⌘B")),
				command.Item(lucide.Settings(), h.Span("Settings"), command.Shortcut("⌘S")),
			),
		),
	)
}

func DialogDemo() h.Node {
	return h.Group(
		h.P(h.Class("text-sm text-muted-foreground"),
			"Press ",
			h.Kbd(h.Class("pointer-events-none inline-flex h-5 items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100 select-none"),
				h.Span(h.Class("text-xs"), "⌘"), "J",
			),
		),
		command.Dialog(
			command.Input(h.Placeholder("Type a command or search...")),
			command.List(
				command.Empty("No results found."),
				command.Group(command.Heading("Suggestions"),
					command.Item(lucide.Calendar(), h.Span("Calendar")),
					command.Item(lucide.Smile(), h.Span("Search Emoji")),
					command.Item(lucide.Calculator(), h.Span("Calculator")),
				),
				command.Separator(),
				command.Group(command.Heading("Settings"),
					command.Item(lucide.User(), h.Span("Profile"), command.Shortcut("⌘P")),
					command.Item(lucide.CreditCard(), h.Span("Billing"), command.Shortcut("⌘B")),
					command.Item(lucide.Settings(), h.Span("Settings"), command.Shortcut("⌘S")),
				),
			),
		),
	)
}
