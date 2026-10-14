package contextmenu

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	cm "shadcn.go/components/ui/context-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the context-menu page
// (ui/apps/v4/examples/radix/context-menu-*.tsx), rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "context-menu", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("context-menu-basic", NovaBasic),
		nova("context-menu-checkboxes", NovaCheckboxes),
		nova("context-menu-demo", NovaDemo),
		nova("context-menu-destructive", NovaDestructive),
		nova("context-menu-groups", NovaGroups),
		nova("context-menu-icons", NovaIcons),
		nova("context-menu-radio", NovaRadio),
		nova("context-menu-rtl", NovaRtl),
		nova("context-menu-shortcuts", NovaShortcuts),
		nova("context-menu-sides", NovaSides),
		nova("context-menu-submenu", NovaSubmenu),
	)
}

const novaTriggerClass = "flex aspect-video w-full max-w-xs items-center justify-center rounded-xl border border-dashed text-sm"

func novaTrigger(fine, coarse string) h.Node {
	return cm.Trigger(h.Class(novaTriggerClass),
		h.Span(h.Class("hidden pointer-fine:inline-block"), fine),
		h.Span(h.Class("hidden pointer-coarse:inline-block"), coarse),
	)
}

func here() h.Node { return novaTrigger("Right click here", "Long press here") }

func NovaBasic() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(cm.Group(
			cm.Item("Back"),
			cm.Item(cm.Disabled(true), "Forward"),
			cm.Item("Reload"),
		)),
	)
}

// Radix's CheckboxItem has no defaultChecked: the items are unchecked and,
// without onCheckedChange, never change.
func NovaCheckboxes() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(cm.Group(
			cm.CheckboxItem(cm.Controlled, "Show Bookmarks Bar"),
			cm.CheckboxItem(cm.Controlled, "Show Full URLs"),
			cm.CheckboxItem(cm.Controlled, "Show Developer Tools"),
		)),
	)
}

func NovaDemo() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(h.Class("w-48"),
			cm.Group(
				cm.Item("Back", cm.Shortcut("⌘[")),
				cm.Item(cm.Disabled(true), "Forward", cm.Shortcut("⌘]")),
				cm.Item("Reload", cm.Shortcut("⌘R")),
				cm.Sub(
					cm.SubTrigger("More Tools"),
					cm.SubContent(h.Class("w-44"),
						cm.Group(
							cm.Item("Save Page..."),
							cm.Item("Create Shortcut..."),
							cm.Item("Name Window..."),
						),
						cm.Separator(),
						cm.Group(cm.Item("Developer Tools")),
						cm.Separator(),
						cm.Group(cm.Item(cm.VariantDestructive, "Delete")),
					),
				),
			),
			cm.Separator(),
			cm.Group(
				cm.CheckboxItem(cm.Checked(true), cm.Controlled, "Show Bookmarks"),
				cm.CheckboxItem(cm.Controlled, "Show Full URLs"),
			),
			cm.Separator(),
			cm.Group(
				cm.RadioGroup(cm.Value("pedro"), cm.Controlled,
					cm.Label("People"),
					cm.RadioItem(cm.Value("pedro"), "Pedro Duarte"),
					cm.RadioItem(cm.Value("colm"), "Colm Tuite"),
				),
			),
		),
	)
}

func NovaDestructive() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(
			cm.Group(
				cm.Item(lucide.Pencil(), "Edit"),
				cm.Item(lucide.Share(), "Share"),
			),
			cm.Separator(),
			cm.Group(cm.Item(cm.VariantDestructive, lucide.Trash(), "Delete")),
		),
	)
}

func NovaGroups() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(
			cm.Group(
				cm.Label("File"),
				cm.Item("New File", cm.Shortcut("⌘N")),
				cm.Item("Open File", cm.Shortcut("⌘O")),
				cm.Item("Save", cm.Shortcut("⌘S")),
			),
			cm.Separator(),
			cm.Group(
				cm.Label("Edit"),
				cm.Item("Undo", cm.Shortcut("⌘Z")),
				cm.Item("Redo", cm.Shortcut("⇧⌘Z")),
			),
			cm.Separator(),
			cm.Group(
				cm.Item("Cut", cm.Shortcut("⌘X")),
				cm.Item("Copy", cm.Shortcut("⌘C")),
				cm.Item("Paste", cm.Shortcut("⌘V")),
			),
			cm.Separator(),
			cm.Group(cm.Item(cm.VariantDestructive, "Delete", cm.Shortcut("⌫"))),
		),
	)
}

func NovaIcons() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(
			cm.Group(
				cm.Item(lucide.Copy(), "Copy"),
				cm.Item(lucide.Scissors(), "Cut"),
				cm.Item(lucide.ClipboardPaste(), "Paste"),
			),
			cm.Separator(),
			cm.Group(cm.Item(cm.VariantDestructive, lucide.Trash(), "Delete")),
		),
	)
}

func NovaRadio() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(
			cm.Group(
				cm.Label("People"),
				cm.RadioGroup(cm.Value("pedro"),
					cm.RadioItem(cm.Value("pedro"), "Pedro Duarte"),
					cm.RadioItem(cm.Value("colm"), "Colm Tuite"),
				),
			),
			cm.Separator(),
			cm.Group(
				cm.Label("Theme"),
				cm.RadioGroup(cm.Value("light"),
					cm.RadioItem(cm.Value("light"), "Light"),
					cm.RadioItem(cm.Value("dark"), "Dark"),
					cm.RadioItem(cm.Value("system"), "System"),
				),
			),
		),
	)
}

func NovaRtl() h.Node {
	lang := h.Data("lang", "ar")
	return cm.ContextMenu(novaTrigger("انقر بزر الماوس الأيمن هنا", "اضغط مطولاً هنا"),
		cm.Content(h.Class("w-48"), lang,
			cm.Group(
				cm.Sub(
					cm.SubTrigger("التنقل"),
					cm.SubContent(h.Class("w-44"), lang,
						cm.Group(
							cm.Item(lucide.ArrowLeft(), "رجوع", cm.Shortcut("⌘[")),
							cm.Item(cm.Disabled(true), lucide.ArrowRight(), "تقدم", cm.Shortcut("⌘]")),
							cm.Item(lucide.RotateCw(), "إعادة تحميل", cm.Shortcut("⌘R")),
						),
					),
				),
				cm.Sub(
					cm.SubTrigger("المزيد من الأدوات"),
					cm.SubContent(h.Class("w-44"), lang,
						cm.Group(
							cm.Item("حفظ الصفحة..."),
							cm.Item("إنشاء اختصار..."),
							cm.Item("تسمية النافذة..."),
						),
						cm.Separator(),
						cm.Group(cm.Item("أدوات المطور")),
						cm.Separator(),
						cm.Group(cm.Item(cm.VariantDestructive, "حذف")),
					),
				),
			),
			cm.Separator(),
			cm.Group(
				cm.CheckboxItem(cm.Checked(true), cm.Controlled, "إظهار الإشارات المرجعية"),
				cm.CheckboxItem(cm.Controlled, "إظهار عناوين URL الكاملة"),
			),
			cm.Separator(),
			cm.Group(
				cm.RadioGroup(cm.Value("pedro"),
					cm.Label("الأشخاص"),
					cm.RadioItem(cm.Value("pedro"), "Pedro Duarte"),
					cm.RadioItem(cm.Value("colm"), "Colm Tuite"),
				),
			),
		),
	)
}

func NovaShortcuts() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(
			cm.Group(
				cm.Item("Back", cm.Shortcut("⌘[")),
				cm.Item(cm.Disabled(true), "Forward", cm.Shortcut("⌘]")),
				cm.Item("Reload", cm.Shortcut("⌘R")),
			),
			cm.Separator(),
			cm.Group(
				cm.Item("Save", cm.Shortcut("⌘S")),
				cm.Item("Save As...", cm.Shortcut("⇧⌘S")),
			),
		),
	)
}

// NovaSides: ContextMenu forces its content's side (right), so the side
// props have no effect.
func NovaSides() h.Node {
	kids := []any{h.Class("grid w-full max-w-sm grid-cols-2 gap-4")}
	for _, s := range []string{"top", "right", "bottom", "left"} {
		kids = append(kids, cm.ContextMenu(
			novaTrigger("Right click ("+s+")", "Long press ("+s+")"),
			cm.Content(cm.Group(
				cm.Item("Back"),
				cm.Item("Forward"),
				cm.Item("Reload"),
			)),
		))
	}
	return h.Div(kids...)
}

func NovaSubmenu() h.Node {
	return cm.ContextMenu(here(),
		cm.Content(
			cm.Group(
				cm.Item("Copy", cm.Shortcut("⌘C")),
				cm.Item("Cut", cm.Shortcut("⌘X")),
			),
			cm.Sub(
				cm.SubTrigger("More Tools"),
				cm.SubContent(
					cm.Group(
						cm.Item("Save Page..."),
						cm.Item("Create Shortcut..."),
						cm.Item("Name Window..."),
					),
					cm.Separator(),
					cm.Group(cm.Item("Developer Tools")),
					cm.Separator(),
					cm.Group(cm.Item(cm.VariantDestructive, "Delete")),
				),
			),
		),
	)
}
