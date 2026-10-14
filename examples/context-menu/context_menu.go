// Package contextmenu ports the upstream context-menu-* examples and the
// context-menu parity fixtures.
package contextmenu

import (
	cm "shadcn.go/components/ui/context-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "context-menu", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("context-menu-demo", Demo),
		ex("context-menu-options", Options),
	)
}

// Demo: upstream passes checked and value without change handlers, so the
// state is fixed (Controlled).
func Demo() h.Node {
	return cm.ContextMenu(
		cm.Trigger(h.Class("flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"), "Right click here"),
		cm.Content(h.Class("w-52"),
			cm.Item(cm.Inset(true), "Back", cm.Shortcut("⌘[")),
			cm.Item(cm.Inset(true), cm.Disabled(true), "Forward", cm.Shortcut("⌘]")),
			cm.Item(cm.Inset(true), "Reload", cm.Shortcut("⌘R")),
			cm.Sub(
				cm.SubTrigger(cm.Inset(true), "More Tools"),
				cm.SubContent(h.Class("w-44"),
					cm.Item("Save Page..."),
					cm.Item("Create Shortcut..."),
					cm.Item("Name Window..."),
					cm.Separator(),
					cm.Item("Developer Tools"),
					cm.Separator(),
					cm.Item(cm.VariantDestructive, "Delete"),
				),
			),
			cm.Separator(),
			cm.CheckboxItem(cm.Checked(true), cm.Controlled, "Show Bookmarks"),
			cm.CheckboxItem(cm.Controlled, "Show Full URLs"),
			cm.Separator(),
			cm.RadioGroup(cm.Value("pedro"), cm.Controlled,
				cm.Label(cm.Inset(true), "People"),
				cm.RadioItem(cm.Value("pedro"), "Pedro Duarte"),
				cm.RadioItem(cm.Value("colm"), "Colm Tuite"),
			),
		),
	)
}

// Options is the parity fixture context-menu-options.
func Options() h.Node {
	return h.Div(h.Class("flex flex-col gap-6"),
		cm.ContextMenu(cm.Modal(false),
			cm.Trigger(h.Class("flex h-[120px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"), "Non-modal area"),
			cm.Content(h.Class("w-48"), cm.Loop(true),
				cm.Label("View"),
				cm.CheckboxItem(cm.Checked(false), "Show grid"),
				cm.Separator(),
				cm.Sub(
					cm.SubTrigger("Arrange"),
					cm.Portal(
						cm.SubContent(h.Class("w-40"),
							cm.Item("Bring forward"),
							cm.Sub(
								cm.SubTrigger("Align"),
								cm.SubContent(
									cm.Item("Left"),
									cm.Item("Right"),
								),
							),
						),
					),
				),
				cm.Item(cm.Disabled(true), "Disabled"),
			),
		),
		cm.ContextMenu(
			cm.Trigger(cm.Disabled(true), h.Class("flex h-[80px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"), "Disabled area"),
			cm.Content(cm.Item("Never shown")),
		),
	)
}
