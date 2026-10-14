// Package menubar ports the upstream menubar-* examples and the menubar
// parity fixtures.
package menubar

import (
	mb "shadcn.go/components/ui/menubar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "menubar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("menubar-demo", Demo),
		ex("menubar-options", Options),
	)
}

// Demo: upstream passes checked and value without change handlers, so the
// state is fixed (Controlled).
func Demo() h.Node {
	return mb.Menubar(
		mb.Menu(
			mb.Trigger("File"),
			mb.Content(
				mb.Item("New Tab ", mb.Shortcut("⌘T")),
				mb.Item("New Window ", mb.Shortcut("⌘N")),
				mb.Item(mb.Disabled(true), "New Incognito Window"),
				mb.Separator(),
				mb.Sub(
					mb.SubTrigger("Share"),
					mb.SubContent(
						mb.Item("Email link"),
						mb.Item("Messages"),
						mb.Item("Notes"),
					),
				),
				mb.Separator(),
				mb.Item("Print... ", mb.Shortcut("⌘P")),
			),
		),
		mb.Menu(
			mb.Trigger("Edit"),
			mb.Content(
				mb.Item("Undo ", mb.Shortcut("⌘Z")),
				mb.Item("Redo ", mb.Shortcut("⇧⌘Z")),
				mb.Separator(),
				mb.Sub(
					mb.SubTrigger("Find"),
					mb.SubContent(
						mb.Item("Search the web"),
						mb.Separator(),
						mb.Item("Find..."),
						mb.Item("Find Next"),
						mb.Item("Find Previous"),
					),
				),
				mb.Separator(),
				mb.Item("Cut"),
				mb.Item("Copy"),
				mb.Item("Paste"),
			),
		),
		mb.Menu(
			mb.Trigger("View"),
			mb.Content(
				mb.CheckboxItem(mb.Controlled, "Always Show Bookmarks Bar"),
				mb.CheckboxItem(mb.Checked(true), mb.Controlled, "Always Show Full URLs"),
				mb.Separator(),
				mb.Item(mb.Inset(true), "Reload ", mb.Shortcut("⌘R")),
				mb.Item(mb.Disabled(true), mb.Inset(true), "Force Reload ", mb.Shortcut("⇧⌘R")),
				mb.Separator(),
				mb.Item(mb.Inset(true), "Toggle Fullscreen"),
				mb.Separator(),
				mb.Item(mb.Inset(true), "Hide Sidebar"),
			),
		),
		mb.Menu(
			mb.Trigger("Profiles"),
			mb.Content(
				mb.RadioGroup(mb.Value("benoit"), mb.Controlled,
					mb.RadioItem(mb.Value("andy"), "Andy"),
					mb.RadioItem(mb.Value("benoit"), "Benoit"),
					mb.RadioItem(mb.Value("Luis"), "Luis"),
				),
				mb.Separator(),
				mb.Item(mb.Inset(true), "Edit..."),
				mb.Separator(),
				mb.Item(mb.Inset(true), "Add Profile..."),
			),
		),
	)
}

// Options is the parity fixture menubar-options.
func Options() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-24"),
		mb.Menubar(mb.Loop(false), mb.DefaultValue("edit"),
			mb.Menu(mb.Value("file"),
				mb.Trigger("File"),
				mb.Content(
					mb.Item("New"),
					mb.Sub(
						mb.SubTrigger("Export"),
						mb.Portal(
							mb.SubContent(
								mb.Item("PDF"),
								mb.Sub(
									mb.SubTrigger("Image"),
									mb.SubContent(
										mb.Item("PNG"),
										mb.Item("SVG"),
									),
								),
							),
						),
					),
				),
			),
			mb.Menu(mb.Value("edit"),
				mb.Trigger("Edit"),
				mb.Content(
					mb.Item("Undo"),
					mb.Separator(),
					mb.Item(mb.Disabled(true), "Redo"),
				),
			),
			mb.Menu(mb.Value("locked"),
				mb.Trigger(mb.Disabled(true), "Locked"),
				mb.Content(mb.Item("Never shown")),
			),
			mb.Menu(mb.Value("help"),
				mb.Trigger("Help"),
				mb.Content(mb.Item("About")),
			),
		),
		mb.Menubar(mb.Dir("rtl"),
			mb.Menu(
				mb.Trigger("Right"),
				mb.Content(
					mb.Item("One"),
					mb.Sub(
						mb.SubTrigger("Sub"),
						mb.SubContent(mb.Item("Inner")),
					),
				),
			),
			mb.Menu(
				mb.Trigger("Left"),
				mb.Content(mb.Item("Two")),
			),
		),
	)
}
