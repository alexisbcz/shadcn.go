package sidebar

import (
	"fmt"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/collapsible"
	dm "shadcn.go/components/ui/dropdown-menu"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the sidebar page (ui/apps/v4/examples/radix/sidebar-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "sidebar", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(nova("sidebar-demo", NovaDemo))
}

type novaNavItem struct {
	title  string
	icon   func(...any) *h.Element
	active bool
	items  []string
}

var novaNavMain = []novaNavItem{
	{"Playground", lucide.SquareTerminal, true, []string{"History", "Starred", "Settings"}},
	{"Models", lucide.Bot, false, []string{"Genesis", "Explorer", "Quantum"}},
	{"Documentation", lucide.BookOpen, false, []string{"Introduction", "Get Started", "Tutorials", "Changelog"}},
	{"Settings", lucide.Settings2, false, []string{"General", "Team", "Billing", "Limits"}},
}

var novaTeams = []struct {
	name string
	logo func(...any) *h.Element
	plan string
}{
	{"Acme Inc", lucide.GalleryVerticalEnd, "Enterprise"},
	{"Acme Corp.", lucide.AudioWaveform, "Startup"},
	{"Evil Corp.", lucide.Command, "Free"},
}

var novaProjects = []struct {
	name string
	icon func(...any) *h.Element
}{
	{"Design Engineering", lucide.Frame},
	{"Sales & Marketing", lucide.ChartPie},
	{"Travel", lucide.Map},
}

// NovaDemo is the radix sidebar-demo: sidebar-07's app sidebar with an
// empty inset.
func NovaDemo() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(sidebar.CollapsibleIcon,
			sidebar.Header(novaTeamSwitcher()),
			sidebar.Content(novaNavMainGroup(), novaNavProjects()),
			sidebar.Footer(novaNavUser()),
			sidebar.Rail(),
		),
		sidebar.Inset(
			h.Header(h.Class("flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"),
				h.Div(h.Class("flex items-center gap-2 px-4"), sidebar.Trigger(h.Class("-ml-1"))),
			),
		),
	)
}

func novaTeamSwitcher() h.Node {
	active := novaTeams[0]
	var teams []any
	teams = append(teams, dm.Label(h.Class("text-xs text-muted-foreground"), "Teams"))
	for i, t := range novaTeams {
		teams = append(teams, dm.Item(h.Class("gap-2 p-2"),
			h.Div(h.Class("flex size-6 items-center justify-center rounded-md border"), t.logo(h.Class("size-3.5 shrink-0"))),
			t.name,
			dm.Shortcut("⌘", h.Comment(""), fmt.Sprint(i+1)),
		))
	}
	return sidebar.Menu(sidebar.MenuItem(dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			sidebar.MenuButton(sidebar.ButtonSizeLg,
				h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
				h.Div(h.Class("flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"), active.logo(h.Class("size-4"))),
				h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
					h.Span(h.Class("truncate font-medium"), active.name),
					h.Span(h.Class("truncate text-xs"), active.plan),
				),
				lucide.ChevronsUpDown(h.Class("ml-auto")),
			),
		),
		dm.Content(h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"), dm.AlignStart, dm.SideRight, dm.SideOffset(4),
			dm.Group(teams...),
			dm.Separator(),
			dm.Group(dm.Item(h.Class("gap-2 p-2"),
				h.Div(h.Class("flex size-6 items-center justify-center rounded-md border bg-transparent"), lucide.Plus(h.Class("size-4"))),
				h.Div(h.Class("font-medium text-muted-foreground"), "Add team"),
			)),
		),
	)))
}

func novaNavMainGroup() h.Node {
	return sidebar.Group(
		sidebar.GroupLabel("Platform"),
		sidebar.Menu(h.Map(novaNavMain, func(it novaNavItem) h.Node {
			return collapsible.Collapsible(lib.AsChild, collapsible.Open(it.active), h.Class("group/collapsible"),
				sidebar.MenuItem(h.Data("slot", "collapsible"),
					collapsible.Trigger(lib.AsChild,
						sidebar.MenuButton(sidebar.Tooltip(it.title),
							it.icon(), h.Span(it.title),
							lucide.ChevronRight(h.Class("ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90")),
						),
					),
					collapsible.Content(sidebar.MenuSub(h.Map(it.items, func(t string) h.Node {
						return sidebar.MenuSubItem(sidebar.MenuSubButton(lib.AsChild, h.A(h.Href("#"), h.Span(t))))
					}))),
				),
			)
		})),
	)
}

func novaNavProjects() h.Node {
	return sidebar.Group(h.Class("group-data-[collapsible=icon]:hidden"),
		sidebar.GroupLabel("Projects"),
		sidebar.Menu(
			h.Map(novaProjects, func(p struct {
				name string
				icon func(...any) *h.Element
			}) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), p.icon(), h.Span(p.name))),
					dm.DropdownMenu(
						dm.Trigger(lib.AsChild,
							sidebar.MenuAction(sidebar.ShowOnHover(true), lucide.Ellipsis(), h.Span(h.Class("sr-only"), "More")),
						),
						dm.Content(h.Class("w-48 rounded-lg"), dm.SideRight, dm.AlignStart,
							dm.Group(
								dm.Item(lucide.Folder(h.Class("text-muted-foreground")), h.Span("View Project")),
								dm.Item(lucide.Forward(h.Class("text-muted-foreground")), h.Span("Share Project")),
							),
							dm.Separator(),
							dm.Group(dm.Item(lucide.Trash2(h.Class("text-muted-foreground")), h.Span("Delete Project"))),
						),
					),
				)
			}),
			sidebar.MenuItem(sidebar.MenuButton(h.Class("text-sidebar-foreground/70"),
				lucide.Ellipsis(h.Class("text-sidebar-foreground/70")), h.Span("More"))),
		),
	)
}

func novaUserBlock() h.Node {
	return h.Group(
		avatar.Avatar(h.Class("h-8 w-8 rounded-lg"),
			avatar.Image(h.Src("/avatars/shadcn.jpg"), h.Alt("shadcn")),
			avatar.Fallback(h.Class("rounded-lg"), "CN"),
		),
		h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
			h.Span(h.Class("truncate font-medium"), "shadcn"),
			h.Span(h.Class("truncate text-xs"), "m@example.com"),
		),
	)
}

func novaNavUser() h.Node {
	return sidebar.Menu(sidebar.MenuItem(dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			sidebar.MenuButton(sidebar.ButtonSizeLg,
				h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
				novaUserBlock(),
				lucide.ChevronsUpDown(h.Class("ml-auto size-4")),
			),
		),
		dm.Content(h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"), dm.SideRight, dm.AlignEnd, dm.SideOffset(4),
			dm.Group(dm.Label(h.Class("p-0 font-normal"),
				h.Div(h.Class("flex items-center gap-2 px-1 py-1.5 text-left text-sm"), novaUserBlock()),
			)),
			dm.Separator(),
			dm.Group(dm.Item(lucide.Sparkles(), "Upgrade to Pro")),
			dm.Separator(),
			dm.Group(
				dm.Item(lucide.BadgeCheck(), "Account"),
				dm.Item(lucide.CreditCard(), "Billing"),
				dm.Item(lucide.Bell(), "Notifications"),
			),
			dm.Separator(),
			dm.Group(dm.Item(lucide.LogOut(), "Log out")),
		),
	)))
}

func init() {
	registry.Register(registry.Example{Name: "nova--sidebar-rtl", Component: "sidebar", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, style.ScopeRTL(NovaRtl()))
	})})
}

// NovaRtl is sidebar-rtl in Arabic (the docs' default language for it).
func NovaRtl() h.Node {
	nav := []novaNavItem{
		{"ملعب", lucide.SquareTerminal, true, []string{"السجل", "المميز", "الإعدادات"}},
		{"النماذج", lucide.Bot, false, []string{"جينيسيس", "إكسبلورر", "كوانتوم"}},
		{"التوثيق", lucide.BookOpen, false, []string{"مقدمة", "ابدأ", "الدروس", "سجل التغييرات"}},
		{"الإعدادات", lucide.Settings2, false, []string{"عام", "الفريق", "الفوترة", "الحدود"}},
	}
	projects := []struct {
		name string
		icon func(...any) *h.Element
	}{
		{"هندسة التصميم", lucide.Frame},
		{"المبيعات والتسويق", lucide.ChartPie},
		{"السفر", lucide.Map},
	}
	userBlock := func(avatarClass ...any) h.Node {
		return h.Group(
			avatar.Avatar(append(avatarClass,
				avatar.Image(h.Src("/avatars/shadcn.jpg"), h.Alt("shadcn")),
				avatar.Fallback(h.Class("rounded-lg"), "CN"),
			)...),
			h.Div(h.Class("grid flex-1 text-start text-sm leading-tight"),
				h.Span(h.Class("truncate font-medium"), "shadcn"),
				h.Span(h.Class("truncate text-xs"), "m@example.com"),
			),
		)
	}
	app := sidebar.Provider(
		sidebar.Sidebar(sidebar.CollapsibleIcon, sidebar.SideRight,
			sidebar.Header(sidebar.Menu(sidebar.MenuItem(
				sidebar.MenuButton(sidebar.ButtonSizeLg, lib.AsChild,
					h.A(h.Href("#"),
						h.Div(h.Class("flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"), lucide.GalleryVerticalEnd(h.Class("size-4"))),
						h.Div(h.Class("flex flex-col gap-0.5 leading-none"),
							h.Span(h.Class("font-medium"), "شركة أكمي"),
							h.Span(h.Class(""), "المؤسسة"),
						),
					),
				),
			))),
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("المنصة"),
					sidebar.Menu(h.Map(nav, func(it novaNavItem) h.Node {
						return collapsible.Collapsible(lib.AsChild, collapsible.Open(it.active), h.Class("group/collapsible"),
							sidebar.MenuItem(h.Data("slot", "collapsible"),
								collapsible.Trigger(lib.AsChild,
									sidebar.MenuButton(sidebar.Tooltip(it.title),
										it.icon(), h.Span(it.title),
										lucide.ChevronRight(h.Class("ms-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90 rtl:rotate-180 rtl:group-data-[state=open]/collapsible:rotate-90")),
									),
								),
								collapsible.Content(sidebar.MenuSub(h.Map(it.items, func(t string) h.Node {
									return sidebar.MenuSubItem(sidebar.MenuSubButton(lib.AsChild, h.A(h.Href("#"), h.Span(t))))
								}))),
							),
						)
					})),
				),
				sidebar.Group(h.Class("group-data-[collapsible=icon]:hidden"),
					sidebar.GroupLabel("المشاريع"),
					sidebar.Menu(
						h.Map(projects, func(p struct {
							name string
							icon func(...any) *h.Element
						}) h.Node {
							return sidebar.MenuItem(
								sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), p.icon(), h.Span(p.name))),
								dm.DropdownMenu(
									dm.Trigger(lib.AsChild,
										sidebar.MenuAction(sidebar.ShowOnHover(true), lucide.Ellipsis(), h.Span(h.Class("sr-only"), "المزيد")),
									),
									dm.Content(h.Class("w-48 rounded-lg"), dm.SideRight, dm.AlignStart,
										dm.Group(
											dm.Item(lucide.Folder(h.Class("text-muted-foreground")), h.Span("عرض المشروع")),
											dm.Item(lucide.Forward(h.Class("text-muted-foreground")), h.Span("مشاركة المشروع")),
										),
										dm.Separator(),
										dm.Group(dm.Item(lucide.Trash2(h.Class("text-muted-foreground")), h.Span("حذف المشروع"))),
									),
								),
							)
						}),
						sidebar.MenuItem(sidebar.MenuButton(h.Class("text-sidebar-foreground/70"),
							lucide.Ellipsis(h.Class("text-sidebar-foreground/70")), h.Span("المزيد"))),
					),
				),
			),
			sidebar.Footer(sidebar.Menu(sidebar.MenuItem(dm.DropdownMenu(
				dm.Trigger(lib.AsChild,
					sidebar.MenuButton(sidebar.ButtonSizeLg,
						h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
						userBlock(h.Class("rounded-lg")),
						lucide.ChevronsUpDown(h.Class("ms-auto size-4")),
					),
				),
				dm.Content(h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"), dm.SideRight, dm.AlignEnd, dm.SideOffset(4),
					dm.Group(dm.Label(h.Class("p-0 font-normal"),
						h.Div(h.Class("flex items-center gap-2 px-1 py-1.5 text-left text-sm"), userBlock()),
					)),
					dm.Separator(),
					dm.Group(dm.Item(lucide.Sparkles(), "ترقية إلى Pro")),
					dm.Separator(),
					dm.Group(
						dm.Item(lucide.BadgeCheck(), "الحساب"),
						dm.Item(lucide.CreditCard(), "الفوترة"),
						dm.Item(lucide.Bell(), "الإشعارات"),
					),
					dm.Separator(),
					dm.Group(dm.Item(lucide.LogOut(), "تسجيل الخروج")),
				),
			)))),
			sidebar.Rail(),
		),
		sidebar.Inset(
			h.Header(h.Class("flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"), h.Dir("rtl"),
				h.Div(h.Class("flex items-center gap-2 px-4"), sidebar.Trigger(h.Class("-ml-1"))),
			),
		),
	)
	return h.Div(h.Class("relative"), h.Dir("rtl"), languageSelector("ar"), app)
}

// languageSelector is the docs' LanguageSelector (a base-ui Select upstream;
// the radix Select renders the same closed trigger).
func languageSelector(value string) h.Node {
	return selectui.Select(selectui.DefaultValue(value),
		selectui.Trigger(selectui.SizeSm, h.Class("absolute top-4 right-4 z-10 w-36 rtl:right-auto rtl:left-4"), h.Dir("ltr"), h.Data("name", "language-selector"), selectui.Value()),
		selectui.Content(h.Dir("ltr"), h.Class("data-open:animate-none data-closed:animate-none"),
			selectui.Group(
				selectui.Item(selectui.ItemValue("en"), "English"),
				selectui.Item(selectui.ItemValue("ar"), "Arabic (العربية)"),
				selectui.Item(selectui.ItemValue("he"), "Hebrew (עברית)"),
			),
		),
	)
}
