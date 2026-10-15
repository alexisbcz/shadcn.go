package command

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/command"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the command page (ui/apps/v4/examples/radix/command-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "command", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("command-demo", NovaDemo),
		nova("command-basic", NovaBasic),
		nova("command-groups", NovaGroups),
		nova("command-scrollable", NovaScrollable),
		nova("command-shortcuts", NovaShortcuts),
		nova("command-rtl", NovaRTL),
	)
	// The aria-nova Command is react-aria's Autocomplete + Menu: no item is
	// highlighted until the pointer or the keyboard moves to one.
	aria := func(name string, fn func(...any) h.Node) registry.Example {
		return registry.Example{Name: "aria--" + name, Component: "command", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn(command.DefaultValue("")))
		})}
	}
	menu := func(groups ...novaGroup) func(...any) h.Node {
		return func(root ...any) h.Node { return novaMenu(groups, root...) }
	}
	registry.Register(
		aria("command-demo", novaDemo),
		aria("command-rtl", novaRTL),
		aria("command-groups", menu(novaSuggestions, novaSettings)),
		aria("command-shortcuts", menu(novaSettings)),
	)
}

type novaItem struct {
	icon     func(...any) *h.Element
	label    string
	shortcut string
	disabled bool
}

type novaGroup struct {
	heading string
	items   []novaItem
}

func novaList(empty string, groups []novaGroup) *command.ListBuilder {
	kids := []any{command.Empty(empty)}
	for i, g := range groups {
		if i > 0 {
			kids = append(kids, command.Separator())
		}
		items := []any{command.Heading(g.heading)}
		for _, it := range g.items {
			var parts []any
			if it.disabled {
				parts = append(parts, command.Disabled(true))
			}
			if it.icon != nil {
				parts = append(parts, it.icon(), h.Span(it.label))
			} else {
				parts = append(parts, it.label)
			}
			if it.shortcut != "" {
				parts = append(parts, command.Shortcut(it.shortcut))
			}
			items = append(items, command.Item(parts...))
		}
		kids = append(kids, command.Group(items...))
	}
	return command.List(kids...)
}

// novaMenu is an "Open Menu" button opening a CommandDialog holding a
// Command (examples/command/dialog.js binds the button).
func novaMenu(groups []novaGroup, root ...any) h.Node {
	return h.Div(h.Class("flex flex-col gap-4"),
		button.Button(button.VariantOutline, h.Class("w-fit"), h.Data("example-command-open", ""), "Open Menu"),
		command.Dialog(
			command.Command(h.Group(root...),
				command.Input(h.Placeholder("Type a command or search...")),
				novaList("No results found.", groups),
			),
		),
	)
}

var (
	novaSuggestions = novaGroup{"Suggestions", []novaItem{{icon: lucide.Calendar, label: "Calendar"}, {icon: lucide.Smile, label: "Search Emoji"}, {icon: lucide.Calculator, label: "Calculator"}}}
	novaSettings    = novaGroup{"Settings", []novaItem{{lucide.User, "Profile", "⌘P", false}, {lucide.CreditCard, "Billing", "⌘B", false}, {lucide.Settings, "Settings", "⌘S", false}}}
)

func NovaDemo() h.Node { return novaDemo() }

func novaDemo(root ...any) h.Node {
	suggestions := novaSuggestions
	suggestions.items = []novaItem{suggestions.items[0], suggestions.items[1], {icon: lucide.Calculator, label: "Calculator", disabled: true}}
	return command.Command(h.Class("max-w-sm rounded-lg border"), h.Group(root...),
		command.Input(h.Placeholder("Type a command or search...")),
		novaList("No results found.", []novaGroup{suggestions, novaSettings}),
	)
}

func NovaBasic() h.Node {
	return novaMenu([]novaGroup{{"Suggestions", []novaItem{{label: "Calendar"}, {label: "Search Emoji"}, {label: "Calculator"}}}})
}

func NovaGroups() h.Node { return novaMenu([]novaGroup{novaSuggestions, novaSettings}) }

func NovaShortcuts() h.Node { return novaMenu([]novaGroup{novaSettings}) }

func NovaScrollable() h.Node {
	return novaMenu([]novaGroup{
		{"Navigation", []novaItem{{lucide.House, "Home", "⌘H", false}, {lucide.Inbox, "Inbox", "⌘I", false}, {lucide.FileText, "Documents", "⌘D", false}, {lucide.Folder, "Folders", "⌘F", false}}},
		{"Actions", []novaItem{{lucide.Plus, "New File", "⌘N", false}, {lucide.FolderPlus, "New Folder", "⇧⌘N", false}, {lucide.Copy, "Copy", "⌘C", false}, {lucide.Scissors, "Cut", "⌘X", false}, {lucide.ClipboardPaste, "Paste", "⌘V", false}, {lucide.Trash, "Delete", "⌫", false}}},
		{"View", []novaItem{{icon: lucide.LayoutGrid, label: "Grid View"}, {icon: lucide.List, label: "List View"}, {lucide.ZoomIn, "Zoom In", "⌘+", false}, {lucide.ZoomOut, "Zoom Out", "⌘-", false}}},
		{"Account", []novaItem{{lucide.User, "Profile", "⌘P", false}, {lucide.CreditCard, "Billing", "⌘B", false}, {lucide.Settings, "Settings", "⌘S", false}, {icon: lucide.Bell, label: "Notifications"}, {icon: lucide.CircleHelp, label: "Help & Support"}}},
		{"Tools", []novaItem{{icon: lucide.Calculator, label: "Calculator"}, {icon: lucide.Calendar, label: "Calendar"}, {icon: lucide.Image, label: "Image Editor"}, {icon: lucide.Code, label: "Code Editor"}}},
	})
}

// NovaRTL is command-rtl in its default language, Arabic.
func NovaRTL() h.Node { return novaRTL() }

func novaRTL(root ...any) h.Node {
	return command.Command(h.Class("max-w-sm rounded-lg border"), h.Dir("rtl"), h.Group(root...),
		command.Input(h.Placeholder("اكتب أمرًا أو ابحث..."), h.Dir("rtl")),
		novaList("لم يتم العثور على نتائج.", []novaGroup{
			{"اقتراحات", []novaItem{{icon: lucide.Calendar, label: "التقويم"}, {icon: lucide.Smile, label: "البحث عن الرموز التعبيرية"}, {icon: lucide.Calculator, label: "الآلة الحاسبة", disabled: true}}},
			{"الإعدادات", []novaItem{{lucide.User, "الملف الشخصي", "⌘P", false}, {lucide.CreditCard, "الفوترة", "⌘B", false}, {lucide.Settings, "الإعدادات", "⌘S", false}}},
		}),
	)
}
