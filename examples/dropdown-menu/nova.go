package dropdownmenu

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the dropdown-menu page
// (ui/apps/v4/examples/radix/dropdown-menu-*.tsx), rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "dropdown-menu", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("dropdown-menu-avatar", NovaAvatar),
		nova("dropdown-menu-basic", NovaBasic),
		nova("dropdown-menu-checkboxes", NovaCheckboxes),
		nova("dropdown-menu-checkboxes-icons", NovaCheckboxesIcons),
		nova("dropdown-menu-complex", NovaComplex),
		nova("dropdown-menu-demo", NovaDemo),
		nova("dropdown-menu-destructive", NovaDestructive),
		nova("dropdown-menu-icons", NovaIcons),
		nova("dropdown-menu-radio-group", NovaRadioGroup),
		nova("dropdown-menu-radio-icons", NovaRadioIcons),
		nova("dropdown-menu-rtl", NovaRtl),
		nova("dropdown-menu-shortcuts", NovaShortcuts),
		nova("dropdown-menu-submenu", NovaSubmenu),
	)
}

func open(label string) h.Node {
	return dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, label))
}

func NovaAvatar() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIcon, h.Class("rounded-full"),
			avatar.Avatar(
				avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("shadcn")),
				avatar.Fallback("LR"),
			),
		)),
		dm.Content(dm.AlignEnd,
			dm.Group(
				dm.Item(lucide.BadgeCheck(), "Account"),
				dm.Item(lucide.CreditCard(), "Billing"),
				dm.Item(lucide.Bell(), "Notifications"),
			),
			dm.Separator(),
			dm.Item(lucide.LogOut(), "Sign Out"),
		),
	)
}

func NovaBasic() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(
			dm.Group(
				dm.Label("My Account"),
				dm.Item("Profile"),
				dm.Item("Billing"),
				dm.Item("Settings"),
			),
			dm.Separator(),
			dm.Item("GitHub"),
			dm.Item("Support"),
			dm.Item(dm.Disabled(true), "API"),
		),
	)
}

func NovaCheckboxes() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(h.Class("w-40"),
			dm.Group(
				dm.Label("Appearance"),
				dm.CheckboxItem(dm.Checked(true), "Status Bar"),
				dm.CheckboxItem(dm.Disabled(true), "Activity Bar"),
				dm.CheckboxItem("Panel"),
			),
		),
	)
}

func NovaCheckboxesIcons() h.Node {
	return dm.DropdownMenu(
		open("Notifications"),
		dm.Content(h.Class("w-48"),
			dm.Group(
				dm.Label("Notification Preferences"),
				dm.CheckboxItem(dm.Checked(true), lucide.Mail(), "Email notifications"),
				dm.CheckboxItem(lucide.MessageSquare(), "SMS notifications"),
				dm.CheckboxItem(dm.Checked(true), lucide.Bell(), "Push notifications"),
			),
		),
	)
}

func NovaComplex() h.Node {
	return dm.DropdownMenu(
		open("Complex Menu"),
		dm.Content(h.Class("w-44"),
			dm.Group(
				dm.Label("File"),
				dm.Item(lucide.File(), "New File", dm.Shortcut("⌘N")),
				dm.Item(lucide.Folder(), "New Folder", dm.Shortcut("⇧⌘N")),
				dm.Sub(
					dm.SubTrigger(lucide.FolderOpen(), "Open Recent"),
					dm.Portal(dm.SubContent(
						dm.Group(
							dm.Label("Recent Projects"),
							dm.Item(lucide.FileCode(), "Project Alpha"),
							dm.Item(lucide.FileCode(), "Project Beta"),
							dm.Sub(
								dm.SubTrigger(lucide.Ellipsis(), "More Projects"),
								dm.Portal(dm.SubContent(
									dm.Item(lucide.FileCode(), "Project Gamma"),
									dm.Item(lucide.FileCode(), "Project Delta"),
								)),
							),
						),
						dm.Separator(),
						dm.Group(dm.Item(lucide.FolderSearch(), "Browse...")),
					)),
				),
				dm.Separator(),
				dm.Item(lucide.Save(), "Save", dm.Shortcut("⌘S")),
				dm.Item(lucide.Download(), "Export", dm.Shortcut("⇧⌘E")),
			),
			dm.Separator(),
			dm.Group(
				dm.Label("View"),
				dm.CheckboxItem(dm.Checked(true), lucide.Eye(), "Show Sidebar"),
				dm.CheckboxItem(lucide.PanelsTopLeft(), "Show Status Bar"),
				dm.Sub(
					dm.SubTrigger(lucide.Palette(), "Theme"),
					dm.Portal(dm.SubContent(
						dm.Group(
							dm.Label("Appearance"),
							dm.RadioGroup(dm.Value("light"),
								dm.RadioItem(dm.Value("light"), lucide.Sun(), "Light"),
								dm.RadioItem(dm.Value("dark"), lucide.Moon(), "Dark"),
								dm.RadioItem(dm.Value("system"), lucide.Monitor(), "System"),
							),
						),
					)),
				),
			),
			dm.Separator(),
			dm.Group(
				dm.Label("Account"),
				dm.Item(lucide.User(), "Profile", dm.Shortcut("⇧⌘P")),
				dm.Item(lucide.CreditCard(), "Billing"),
				dm.Sub(
					dm.SubTrigger(lucide.Settings(), "Settings"),
					dm.Portal(dm.SubContent(
						dm.Group(
							dm.Label("Preferences"),
							dm.Item(lucide.Keyboard(), "Keyboard Shortcuts"),
							dm.Item(lucide.Languages(), "Language"),
							dm.Sub(
								dm.SubTrigger(lucide.Bell(), "Notifications"),
								dm.Portal(dm.SubContent(
									dm.Group(
										dm.Label("Notification Types"),
										dm.CheckboxItem(dm.Checked(true), lucide.Bell(), "Push Notifications"),
										dm.CheckboxItem(dm.Checked(true), lucide.Mail(), "Email Notifications"),
									),
								)),
							),
						),
						dm.Separator(),
						dm.Group(dm.Item(lucide.Shield(), "Privacy & Security")),
					)),
				),
			),
			dm.Separator(),
			dm.Group(
				dm.Item(lucide.CircleHelp(), "Help & Support"),
				dm.Item(lucide.FileText(), "Documentation"),
			),
			dm.Separator(),
			dm.Group(
				dm.Item(dm.VariantDestructive, lucide.LogOut(), "Sign Out", dm.Shortcut("⇧⌘Q")),
			),
		),
	)
}

func NovaDemo() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(h.Class("w-40"), dm.AlignStart,
			dm.Group(
				dm.Label("My Account"),
				dm.Item("Profile", dm.Shortcut("⇧⌘P")),
				dm.Item("Billing", dm.Shortcut("⌘B")),
				dm.Item("Settings", dm.Shortcut("⌘S")),
			),
			dm.Separator(),
			dm.Group(
				dm.Item("Team"),
				dm.Sub(
					dm.SubTrigger("Invite users"),
					dm.Portal(dm.SubContent(
						dm.Item("Email"),
						dm.Item("Message"),
						dm.Separator(),
						dm.Item("More..."),
					)),
				),
				dm.Item("New Team", dm.Shortcut("⌘+T")),
			),
			dm.Separator(),
			dm.Group(
				dm.Item("GitHub"),
				dm.Item("Support"),
				dm.Item(dm.Disabled(true), "API"),
			),
			dm.Separator(),
			dm.Group(dm.Item("Log out", dm.Shortcut("⇧⌘Q"))),
		),
	)
}

func NovaDestructive() h.Node {
	return dm.DropdownMenu(
		open("Actions"),
		dm.Content(
			dm.Group(
				dm.Item(lucide.Pencil(), "Edit"),
				dm.Item(lucide.Share(), "Share"),
			),
			dm.Separator(),
			dm.Group(dm.Item(dm.VariantDestructive, lucide.Trash(), "Delete")),
		),
	)
}

func NovaIcons() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(
			dm.Item(lucide.User(), "Profile"),
			dm.Item(lucide.CreditCard(), "Billing"),
			dm.Item(lucide.Settings(), "Settings"),
			dm.Separator(),
			dm.Item(dm.VariantDestructive, lucide.LogOut(), "Log out"),
		),
	)
}

func NovaRadioGroup() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(h.Class("w-32"),
			dm.Group(
				dm.Label("Panel Position"),
				dm.RadioGroup(dm.Value("bottom"),
					dm.RadioItem(dm.Value("top"), "Top"),
					dm.RadioItem(dm.Value("bottom"), "Bottom"),
					dm.RadioItem(dm.Value("right"), "Right"),
				),
			),
		),
	)
}

func NovaRadioIcons() h.Node {
	return dm.DropdownMenu(
		open("Payment Method"),
		dm.Content(h.Class("min-w-56"),
			dm.Group(
				dm.Label("Select Payment Method"),
				dm.RadioGroup(dm.Value("card"),
					dm.RadioItem(dm.Value("card"), lucide.CreditCard(), "Credit Card"),
					dm.RadioItem(dm.Value("paypal"), lucide.Wallet(), "PayPal"),
					dm.RadioItem(dm.Value("bank"), lucide.Building2(), "Bank Transfer"),
				),
			),
		),
	)
}

func NovaRtl() h.Node {
	lang := h.Data("lang", "ar")
	return dm.DropdownMenu(
		open("افتح القائمة"),
		dm.Content(h.Class("w-36"), lang,
			dm.Group(
				dm.Sub(
					dm.SubTrigger("الحساب"),
					dm.Portal(dm.SubContent(lang,
						dm.Group(
							dm.Item(lucide.User(), "الملف الشخصي"),
							dm.Item(lucide.CreditCard(), "الفوترة"),
							dm.Item(lucide.Settings(), "الإعدادات"),
						),
					)),
				),
			),
			dm.Separator(),
			dm.Group(
				dm.Label("الفريق"),
				dm.Item("الفريق"),
				dm.Sub(
					dm.SubTrigger("دعوة المستخدمين"),
					dm.Portal(dm.SubContent(lang,
						dm.Item("البريد الإلكتروني"),
						dm.Item("رسالة"),
						dm.Sub(
							dm.SubTrigger("المزيد"),
							dm.Portal(dm.SubContent(lang,
								dm.Item("تقويم"),
								dm.Item("دردشة"),
								dm.Separator(),
								dm.Item("خطاف ويب"),
							)),
						),
						dm.Separator(),
						dm.Item("متقدم..."),
					)),
				),
				dm.Item("فريق جديد", dm.Shortcut("⌘+T")),
			),
			dm.Separator(),
			dm.Group(
				dm.Label("عرض"),
				dm.CheckboxItem(dm.Checked(true), "شريط الحالة"),
				dm.CheckboxItem("شريط النشاط"),
				dm.CheckboxItem("اللوحة"),
			),
			dm.Separator(),
			dm.Group(
				dm.Label("الموضع"),
				dm.RadioGroup(dm.Value("bottom"),
					dm.RadioItem(dm.Value("top"), "أعلى"),
					dm.RadioItem(dm.Value("bottom"), "أسفل"),
					dm.RadioItem(dm.Value("right"), "يمين"),
					dm.RadioItem(dm.Value("left"), "يسار"),
				),
			),
			dm.Separator(),
			dm.Group(dm.Item(dm.VariantDestructive, "تسجيل الخروج")),
		),
	)
}

func NovaShortcuts() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(
			dm.Group(
				dm.Label("My Account"),
				dm.Item("Profile", dm.Shortcut("⇧⌘P")),
				dm.Item("Billing", dm.Shortcut("⌘B")),
				dm.Item("Settings", dm.Shortcut("⌘S")),
			),
			dm.Separator(),
			dm.Item("Log out", dm.Shortcut("⇧⌘Q")),
		),
	)
}

func NovaSubmenu() h.Node {
	return dm.DropdownMenu(
		open("Open"),
		dm.Content(
			dm.Group(
				dm.Item("Team"),
				dm.Sub(
					dm.SubTrigger("Invite users"),
					dm.Portal(dm.SubContent(
						dm.Item("Email"),
						dm.Item("Message"),
						dm.Sub(
							dm.SubTrigger("More options"),
							dm.Portal(dm.SubContent(
								dm.Item("Calendly"),
								dm.Item("Slack"),
								dm.Separator(),
								dm.Item("Webhook"),
							)),
						),
						dm.Separator(),
						dm.Item("Advanced..."),
					)),
				),
				dm.Item("New Team", dm.Shortcut("⌘+T")),
			),
		),
	)
}
