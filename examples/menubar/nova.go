package menubar

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	mb "shadcn.go/components/ui/menubar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the menubar page (ui/apps/v4/examples/radix/menubar-*.tsx),
// rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "menubar", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("menubar-checkbox", NovaCheckbox),
		nova("menubar-demo", NovaDemo),
		nova("menubar-icons", NovaIcons),
		nova("menubar-radio", NovaRadio),
		nova("menubar-rtl", NovaRtl),
		nova("menubar-submenu", NovaSubmenu),
	)
}

var inset = mb.Inset(true)

func NovaCheckbox() h.Node {
	return mb.Menubar(h.Class("w-72"),
		mb.Menu(
			mb.Trigger("View"),
			mb.Content(h.Class("w-64"),
				mb.CheckboxItem(mb.Controlled, "Always Show Bookmarks Bar"),
				mb.CheckboxItem(mb.Checked(true), mb.Controlled, "Always Show Full URLs"),
				mb.Separator(),
				mb.Item(inset, "Reload ", mb.Shortcut("⌘R")),
				mb.Item(mb.Disabled(true), inset, "Force Reload ", mb.Shortcut("⇧⌘R")),
			),
		),
		mb.Menu(
			mb.Trigger("Format"),
			mb.Content(
				mb.CheckboxItem(mb.Checked(true), mb.Controlled, "Strikethrough"),
				mb.CheckboxItem(mb.Controlled, "Code"),
				mb.CheckboxItem(mb.Controlled, "Superscript"),
			),
		),
	)
}

// novaDemo is menubar-demo and, translated, menubar-rtl.
func novaDemo(t map[string]string, rtl bool) h.Node {
	var lang any
	var opts []any
	if rtl {
		lang = h.Data("lang", "ar")
		opts = append(opts, mb.Dir("rtl"))
	}
	profiles := []any{mb.Value("benoit")}
	if !rtl {
		profiles = append(profiles, mb.Controlled)
	}
	return mb.Menubar(append(opts, h.Class("w-72"),
		mb.Menu(
			mb.Trigger(t["file"]),
			mb.Content(lang,
				mb.Group(
					mb.Item(t["newTab"]+" ", mb.Shortcut("⌘T")),
					mb.Item(t["newWindow"]+" ", mb.Shortcut("⌘N")),
					mb.Item(mb.Disabled(true), t["newIncognitoWindow"]),
				),
				mb.Separator(),
				mb.Group(
					mb.Sub(
						mb.SubTrigger(t["share"]),
						mb.SubContent(lang,
							mb.Group(
								mb.Item(t["emailLink"]),
								mb.Item(t["messages"]),
								mb.Item(t["notes"]),
							),
						),
					),
				),
				mb.Separator(),
				mb.Group(mb.Item(t["print"]+" ", mb.Shortcut("⌘P"))),
			),
		),
		mb.Menu(
			mb.Trigger(t["edit"]),
			mb.Content(lang,
				mb.Group(
					mb.Item(t["undo"]+" ", mb.Shortcut("⌘Z")),
					mb.Item(t["redo"]+" ", mb.Shortcut("⇧⌘Z")),
				),
				mb.Separator(),
				mb.Group(
					mb.Sub(
						mb.SubTrigger(t["find"]),
						mb.SubContent(lang,
							mb.Group(mb.Item(t["searchTheWeb"])),
							mb.Separator(),
							mb.Group(
								mb.Item(t["findItem"]),
								mb.Item(t["findNext"]),
								mb.Item(t["findPrevious"]),
							),
						),
					),
				),
				mb.Separator(),
				mb.Group(
					mb.Item(t["cut"]),
					mb.Item(t["copy"]),
					mb.Item(t["paste"]),
				),
			),
		),
		mb.Menu(
			mb.Trigger(t["view"]),
			mb.Content(h.Class("w-44"), lang,
				mb.Group(
					mb.CheckboxItem(mb.Controlled, t["bookmarksBar"]),
					mb.CheckboxItem(mb.Checked(true), mb.Controlled, t["fullUrls"]),
				),
				mb.Separator(),
				mb.Group(
					mb.Item(inset, t["reload"]+" ", mb.Shortcut("⌘R")),
					mb.Item(mb.Disabled(true), inset, t["forceReload"]+" ", mb.Shortcut("⇧⌘R")),
				),
				mb.Separator(),
				mb.Group(mb.Item(inset, t["toggleFullscreen"])),
				mb.Separator(),
				mb.Group(mb.Item(inset, t["hideSidebar"])),
			),
		),
		mb.Menu(
			mb.Trigger(t["profiles"]),
			mb.Content(lang,
				mb.RadioGroup(append(profiles,
					mb.RadioItem(mb.Value("andy"), t["andy"]),
					mb.RadioItem(mb.Value("benoit"), t["benoit"]),
					mb.RadioItem(mb.Value("Luis"), t["luis"]),
				)...),
				mb.Separator(),
				mb.Group(mb.Item(inset, t["editProfile"])),
				mb.Separator(),
				mb.Group(mb.Item(inset, t["addProfile"])),
			),
		),
	)...)
}

var novaEn = map[string]string{
	"file": "File", "newTab": "New Tab", "newWindow": "New Window", "newIncognitoWindow": "New Incognito Window",
	"share": "Share", "emailLink": "Email link", "messages": "Messages", "notes": "Notes", "print": "Print...",
	"edit": "Edit", "undo": "Undo", "redo": "Redo", "find": "Find", "searchTheWeb": "Search the web",
	"findItem": "Find...", "findNext": "Find Next", "findPrevious": "Find Previous", "cut": "Cut", "copy": "Copy",
	"paste": "Paste", "view": "View", "bookmarksBar": "Bookmarks Bar", "fullUrls": "Full URLs", "reload": "Reload",
	"forceReload": "Force Reload", "toggleFullscreen": "Toggle Fullscreen", "hideSidebar": "Hide Sidebar",
	"profiles": "Profiles", "andy": "Andy", "benoit": "Benoit", "luis": "Luis", "editProfile": "Edit...",
	"addProfile": "Add Profile...",
}

var novaAr = map[string]string{
	"file": "ملف", "newTab": "علامة تبويب جديدة", "newWindow": "نافذة جديدة", "newIncognitoWindow": "نافذة التصفح المتخفي الجديدة",
	"share": "مشاركة", "emailLink": "رابط البريد الإلكتروني", "messages": "الرسائل", "notes": "الملاحظات", "print": "طباعة...",
	"edit": "تعديل", "undo": "تراجع", "redo": "إعادة", "find": "بحث", "searchTheWeb": "البحث على الويب",
	"findItem": "بحث...", "findNext": "البحث التالي", "findPrevious": "البحث السابق", "cut": "قص", "copy": "نسخ",
	"paste": "لصق", "view": "عرض", "bookmarksBar": "شريط الإشارات المرجعية", "fullUrls": "عناوين URL الكاملة", "reload": "إعادة تحميل",
	"forceReload": "إعادة تحميل قسري", "toggleFullscreen": "تبديل وضع ملء الشاشة", "hideSidebar": "إخفاء الشريط الجانبي",
	"profiles": "الملفات الشخصية", "andy": "Andy", "benoit": "Benoit", "luis": "Luis", "editProfile": "تعديل...",
	"addProfile": "إضافة ملف شخصي...",
}

func NovaDemo() h.Node { return novaDemo(novaEn, false) }

func NovaRtl() h.Node { return novaDemo(novaAr, true) }

func NovaIcons() h.Node {
	return mb.Menubar(h.Class("w-72"),
		mb.Menu(
			mb.Trigger("File"),
			mb.Content(
				mb.Item(lucide.File(), "New File ", mb.Shortcut("⌘N")),
				mb.Item(lucide.Folder(), "Open Folder"),
				mb.Separator(),
				mb.Item(lucide.Save(), "Save ", mb.Shortcut("⌘S")),
			),
		),
		mb.Menu(
			mb.Trigger("More"),
			mb.Content(
				mb.Group(
					mb.Item(lucide.Settings(), "Settings"),
					mb.Item(lucide.CircleHelp(), "Help"),
					mb.Separator(),
					mb.Item(mb.VariantDestructive, lucide.Trash(), "Delete"),
				),
			),
		),
	)
}

func NovaRadio() h.Node {
	return mb.Menubar(h.Class("w-72"),
		mb.Menu(
			mb.Trigger("Profiles"),
			mb.Content(
				mb.RadioGroup(mb.Value("benoit"),
					mb.RadioItem(mb.Value("andy"), "Andy"),
					mb.RadioItem(mb.Value("benoit"), "Benoit"),
					mb.RadioItem(mb.Value("luis"), "Luis"),
				),
				mb.Separator(),
				mb.Item(inset, "Edit..."),
				mb.Item(inset, "Add Profile..."),
			),
		),
		mb.Menu(
			mb.Trigger("Theme"),
			mb.Content(
				mb.RadioGroup(mb.Value("system"),
					mb.RadioItem(mb.Value("light"), "Light"),
					mb.RadioItem(mb.Value("dark"), "Dark"),
					mb.RadioItem(mb.Value("system"), "System"),
				),
			),
		),
	)
}

func NovaSubmenu() h.Node {
	return mb.Menubar(h.Class("w-72"),
		mb.Menu(
			mb.Trigger("File"),
			mb.Content(
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
	)
}
