// Package blockssidebara ports the upstream sidebar-01 to sidebar-08 blocks
// (registry/new-york-v4/blocks/sidebar-0N): full pages built on
// components/ui/sidebar.
package blockssidebara

import (
	_ "embed"
	"fmt"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/collapsible"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed switcher.js
var switcherScript string

func init() {
	runtime.Module("example-blocks-sidebar-switcher", switcherScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "sidebar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("sidebar-01", Sidebar01),
		ex("sidebar-02", Sidebar02),
		ex("sidebar-03", Sidebar03),
		ex("sidebar-04", Sidebar04),
		ex("sidebar-05", Sidebar05),
		ex("sidebar-06", Sidebar06),
		ex("sidebar-07", Sidebar07),
		ex("sidebar-08", Sidebar08),
	)
}

// --- shared page parts -------------------------------------------------------------------------

func crumbs() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		breadcrumb.Item(h.Class("hidden md:block"), breadcrumb.Link(h.Href("#"), "Build Your Application")),
		breadcrumb.Separator(h.Class("hidden md:block")),
		breadcrumb.Item(breadcrumb.Page("Data Fetching")),
	))
}

func sepData() h.Node {
	return separator.Separator(separator.Vertical, h.Class("mr-2 data-[orientation=vertical]:h-4"))
}

func sepH4() h.Node {
	return separator.Separator(separator.Vertical, h.Class("mr-2 h-4"))
}

// headerFlat is the header of sidebar-01/04/05/06: trigger, separator and
// breadcrumb directly in the header.
func headerFlat(class string) h.Node {
	return h.Header(h.Class(class), sidebar.Trigger(h.Class("-ml-1")), sepData(), crumbs())
}

// headerNested wraps the header items in a div (sidebar-07/08).
func headerNested(class string) h.Node {
	return h.Header(h.Class(class),
		h.Div(h.Class("flex items-center gap-2 px-4"), sidebar.Trigger(h.Class("-ml-1")), sepData(), crumbs()),
	)
}

func gridBody(class string) h.Node {
	box := func() h.Node { return h.Div(h.Class("aspect-video rounded-xl bg-muted/50")) }
	return h.Div(h.Class(class),
		h.Div(h.Class("grid auto-rows-min gap-4 md:grid-cols-3"), box(), box(), box()),
		h.Div(h.Class("min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min")),
	)
}

func page(app h.Node, head, body h.Node, provider ...any) h.Node {
	return sidebar.Provider(append(provider, app, sidebar.Inset(head, body))...)
}

// --- data --------------------------------------------------------------------------------------

type docItem struct {
	title  string
	active bool
}

type docSection struct {
	title string
	items []docItem
}

func items(titles ...string) []docItem {
	out := make([]docItem, len(titles))
	for i, t := range titles {
		out[i] = docItem{title: t, active: t == "Data Fetching"}
	}
	return out
}

var docs = []docSection{
	{"Getting Started", items("Installation", "Project Structure")},
	{"Build Your Application", items("Routing", "Data Fetching", "Rendering", "Caching", "Styling", "Optimizing", "Configuring", "Testing", "Authentication", "Deploying", "Upgrading", "Examples")},
	{"API Reference", items("Components", "File Conventions", "Functions", "next.config.js Options", "CLI", "Edge Runtime")},
	{"Architecture", items("Accessibility", "Fast Refresh", "Next.js Compiler", "Supported Browsers", "Turbopack")},
}

var docsCommunity = append(append([]docSection{}, docs...), docSection{"Community", items("Contribution Guide")})

var versions = []string{"1.0.1", "1.1.0-alpha", "2.0.0-beta1"}

// --- shared sidebar parts ----------------------------------------------------------------------

func logoBox(icon h.Node) h.Node {
	return h.Div(h.Class("flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"), icon)
}

// versionSwitcher is VersionSwitcher (switcher.js keeps the selected
// version).
func versionSwitcher() h.Node {
	return sidebar.Menu(sidebar.MenuItem(
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild,
				sidebar.MenuButton(sidebar.ButtonSizeLg,
					h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
					h.Data("example-blocks-sidebar-switcher", "version"),
					logoBox(lucide.GalleryVerticalEnd(h.Class("size-4"))),
					h.Div(h.Class("flex flex-col gap-0.5 leading-none"),
						h.Span(h.Class("font-medium"), "Documentation"),
						h.Span(h.Class(""), "v", h.Comment(""), versions[0]),
					),
					lucide.ChevronsUpDown(h.Class("ml-auto")),
				),
			),
			dm.Content(h.Class("w-(--radix-dropdown-menu-trigger-width)"), dm.AlignStart,
				h.Map(versions, func(v string) h.Node {
					return dm.Item("v", h.Comment(""), v, h.Comment(""), " ", h.If(v == versions[0], lucide.Check(h.Class("ml-auto"))))
				}),
			),
		),
	))
}

func searchForm() h.Node {
	return h.Form(
		sidebar.Group(h.Class("py-0"),
			sidebar.GroupContent(h.Class("relative"),
				label.Label(h.For("search"), h.Class("sr-only"), "Search"),
				sidebar.Input(h.ID("search"), h.Placeholder("Search the docs..."), h.Class("pl-8")),
				lucide.Search(h.Class("pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none")),
			),
		),
	)
}

func docLink(it docItem) h.Node {
	return sidebar.MenuItem(sidebar.MenuButton(lib.AsChild, sidebar.IsActive(it.active), h.A(h.Href("#"), it.title)))
}

func docSubMenu(its []docItem, class ...any) h.Node {
	return sidebar.MenuSub(append(class, h.Map(its, func(it docItem) h.Node {
		return sidebar.MenuSubItem(sidebar.MenuSubButton(lib.AsChild, sidebar.IsActive(it.active), h.A(h.Href("#"), it.title)))
	}))...)
}

func docsHeader() h.Node { return sidebar.Header(docsMenu()) }

func docsMenu() h.Node {
	return sidebar.Menu(sidebar.MenuItem(
		sidebar.MenuButton(sidebar.ButtonSizeLg, lib.AsChild,
			h.A(h.Href("#"),
				logoBox(lucide.GalleryVerticalEnd(h.Class("size-4"))),
				h.Div(h.Class("flex flex-col gap-0.5 leading-none"),
					h.Span(h.Class("font-medium"), "Documentation"),
					h.Span(h.Class(""), "v1.0.0"),
				),
			),
		),
	))
}

// --- sidebar-01 ------------------------------------------------------------------------------

// Sidebar01 is sidebar-01: a simple sidebar with navigation grouped by
// section.
func Sidebar01() h.Node {
	app := sidebar.Sidebar(
		sidebar.Header(versionSwitcher(), searchForm()),
		sidebar.Content(h.Map(docs, func(s docSection) h.Node {
			return sidebar.Group(
				sidebar.GroupLabel(s.title),
				sidebar.GroupContent(sidebar.Menu(h.Map(s.items, docLink))),
			)
		})),
		sidebar.Rail(),
	)
	return page(app, headerFlat("flex h-16 shrink-0 items-center gap-2 border-b px-4"), gridBody("flex flex-1 flex-col gap-4 p-4"))
}

// --- sidebar-02 ------------------------------------------------------------------------------

// Sidebar02 is sidebar-02: a sidebar with collapsible sections.
func Sidebar02() h.Node {
	app := sidebar.Sidebar(
		sidebar.Header(versionSwitcher(), searchForm()),
		sidebar.Content(h.Class("gap-0"), h.Map(docsCommunity, func(s docSection) h.Node {
			return collapsible.Collapsible(h.Attr("title", s.title), collapsible.Open(true), h.Class("group/collapsible"),
				sidebar.Group(
					sidebar.GroupLabel(lib.AsChild,
						h.Class("group/label text-sm text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"),
						collapsible.Trigger(s.title, " ",
							lucide.ChevronRight(h.Class("ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90")),
						),
					),
					collapsible.Content(sidebar.GroupContent(sidebar.Menu(h.Map(s.items, docLink)))),
				),
			)
		})),
		sidebar.Rail(),
	)
	var boxes []any
	for range 24 {
		boxes = append(boxes, h.Div(h.Class("aspect-video h-12 w-full rounded-lg bg-muted/50")))
	}
	head := h.Header(h.Class("sticky top-0 flex h-16 shrink-0 items-center gap-2 border-b bg-background px-4"),
		sidebar.Trigger(h.Class("-ml-1")), sepH4(), crumbs())
	return page(app, head, h.Div(append([]any{h.Class("flex flex-1 flex-col gap-4 p-4")}, boxes...)...))
}

// --- sidebar-03 / 04 -------------------------------------------------------------------------

func subMenuSidebar(variant sidebar.Variant, menuClass, subClass string, rail bool) h.Node {
	var subOpts []any
	if subClass != "" {
		subOpts = append(subOpts, h.Class(subClass))
	}
	var menuOpts []any
	if menuClass != "" {
		menuOpts = append(menuOpts, h.Class(menuClass))
	}
	return sidebar.Sidebar(sidebar.Variant(variant),
		docsHeader(),
		sidebar.Content(sidebar.Group(
			sidebar.Menu(append(menuOpts, h.Map(docsCommunity, func(s docSection) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), h.Class("font-medium"), s.title)),
					docSubMenu(s.items, subOpts...),
				)
			}))...),
		)),
		h.If(rail, sidebar.Rail()),
	)
}

// Sidebar03 is sidebar-03: a sidebar with submenus.
func Sidebar03() h.Node {
	head := h.Header(h.Class("flex h-16 shrink-0 items-center gap-2 border-b"),
		h.Div(h.Class("flex items-center gap-2 px-3"), sidebar.Trigger(), sepH4(), crumbs()),
	)
	return page(subMenuSidebar(sidebar.VariantSidebar, "", "", true), head, gridBody("flex flex-1 flex-col gap-4 p-4"))
}

// Sidebar04 is sidebar-04: a floating sidebar with submenus.
func Sidebar04() h.Node {
	return page(subMenuSidebar(sidebar.VariantFloating, "gap-2", "ml-0 border-l-0 px-1.5", false),
		headerFlat("flex h-16 shrink-0 items-center gap-2 px-4"),
		gridBody("flex flex-1 flex-col gap-4 p-4 pt-0"),
		h.Style("--sidebar-width: 19rem"),
	)
}

// --- sidebar-05 ------------------------------------------------------------------------------

// Sidebar05 is sidebar-05: a sidebar with collapsible submenus.
func Sidebar05() h.Node {
	var groups []any
	for i, s := range docsCommunity {
		groups = append(groups, collapsible.Collapsible(collapsible.Open(i == 1), h.Class("group/collapsible"),
			sidebar.MenuItem(
				collapsible.Trigger(lib.AsChild,
					sidebar.MenuButton(s.title, " ",
						lucide.Plus(h.Class("ml-auto group-data-[state=open]/collapsible:hidden")),
						lucide.Minus(h.Class("ml-auto group-data-[state=closed]/collapsible:hidden")),
					),
				),
				collapsible.Content(docSubMenu(s.items)),
			),
		))
	}
	app := sidebar.Sidebar(
		sidebar.Header(
			docsMenu(),
			searchForm(),
		),
		sidebar.Content(sidebar.Group(sidebar.Menu(groups...))),
		sidebar.Rail(),
	)
	return page(app, headerFlat("flex h-16 shrink-0 items-center gap-2 border-b px-4"), gridBody("flex flex-1 flex-col gap-4 p-4"))
}

// --- sidebar-06 ------------------------------------------------------------------------------

func optInForm() h.Node {
	return card.Card(h.Class("gap-2 py-4 shadow-none"),
		card.Header(h.Class("px-4"),
			card.Title(h.Class("text-sm"), "Subscribe to our newsletter"),
			card.Description("Opt-in to receive updates and news about the sidebar."),
		),
		card.Content(h.Class("px-4"),
			h.Form(h.Div(h.Class("grid gap-2.5"),
				sidebar.Input(h.Type("email"), h.Placeholder("Email")),
				button.Button(h.Class("w-full bg-sidebar-primary text-sidebar-primary-foreground shadow-none"), button.SizeSm, "Subscribe"),
			)),
		),
	)
}

// Sidebar06 is sidebar-06: a sidebar with submenus as dropdowns.
func Sidebar06() h.Node {
	nav := sidebar.Group(sidebar.Menu(h.Map(docs, func(s docSection) h.Node {
		return dm.DropdownMenu(sidebar.MenuItem(
			dm.Trigger(lib.AsChild,
				sidebar.MenuButton(h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
					s.title, " ", lucide.Ellipsis(h.Class("ml-auto")),
				),
			),
			dm.Content(dm.SideRight, dm.AlignStart, h.Class("min-w-56 rounded-lg"),
				h.Map(s.items, func(it docItem) h.Node {
					return dm.Item(lib.AsChild, h.A(h.Href("#"), it.title))
				}),
			),
		))
	})))
	app := sidebar.Sidebar(
		docsHeader(),
		sidebar.Content(nav),
		sidebar.Footer(h.Div(h.Class("p-1"), optInForm())),
		sidebar.Rail(),
	)
	return page(app, headerFlat("flex h-16 shrink-0 items-center gap-2 border-b px-4"), gridBody("flex flex-1 flex-col gap-4 p-4"))
}

// --- sidebar-07 / 08 -------------------------------------------------------------------------

type navItem struct {
	title  string
	icon   func(...any) *h.Element
	active bool
	items  []string
}

var platform = []navItem{
	{"Playground", lucide.SquareTerminal, true, []string{"History", "Starred", "Settings"}},
	{"Models", lucide.Bot, false, []string{"Genesis", "Explorer", "Quantum"}},
	{"Documentation", lucide.BookOpen, false, []string{"Introduction", "Get Started", "Tutorials", "Changelog"}},
	{"Settings", lucide.Settings2, false, []string{"General", "Team", "Billing", "Limits"}},
}

type project struct {
	name string
	icon func(...any) *h.Element
}

var projects = []project{
	{"Design Engineering", lucide.Frame},
	{"Sales & Marketing", lucide.ChartPie},
	{"Travel", lucide.Map},
}

type team struct {
	name string
	logo func(...any) *h.Element
	plan string
}

var teams = []team{
	{"Acme Inc", lucide.GalleryVerticalEnd, "Enterprise"},
	{"Acme Corp.", lucide.AudioWaveform, "Startup"},
	{"Evil Corp.", lucide.Command, "Free"},
}

func subItems(its []string) h.Node {
	return sidebar.MenuSub(h.Map(its, func(t string) h.Node {
		return sidebar.MenuSubItem(sidebar.MenuSubButton(lib.AsChild, h.A(h.Href("#"), h.Span(t))))
	}))
}

// collapsibleItem is <Collapsible asChild><SidebarMenuItem>: the menu
// item receives the collapsible's props, its data-slot included.
func collapsibleItem(open bool, class string, children ...any) h.Node {
	opts := []any{lib.AsChild, collapsible.Open(open)}
	if class != "" {
		opts = append(opts, h.Class(class))
	}
	return collapsible.Collapsible(append(opts, sidebar.MenuItem(append([]any{h.Data("slot", "collapsible")}, children...)...))...)
}

func navMain07() h.Node {
	return sidebar.Group(
		sidebar.GroupLabel("Platform"),
		sidebar.Menu(h.Map(platform, func(it navItem) h.Node {
			return collapsibleItem(it.active, "group/collapsible",
				collapsible.Trigger(lib.AsChild,
					sidebar.MenuButton(sidebar.Tooltip(it.title),
						it.icon(), h.Span(it.title),
						lucide.ChevronRight(h.Class("ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90")),
					),
				),
				collapsible.Content(subItems(it.items)),
			)
		})),
	)
}

func navMain08() h.Node {
	return sidebar.Group(
		sidebar.GroupLabel("Platform"),
		sidebar.Menu(h.Map(platform, func(it navItem) h.Node {
			return collapsibleItem(it.active, "",
				sidebar.MenuButton(lib.AsChild, sidebar.Tooltip(it.title),
					h.A(h.Href("#"), it.icon(), h.Span(it.title)),
				),
				collapsible.Trigger(lib.AsChild,
					sidebar.MenuAction(h.Class("data-[state=open]:rotate-90"),
						lucide.ChevronRight(), h.Span(h.Class("sr-only"), "Toggle"),
					),
				),
				collapsible.Content(subItems(it.items)),
			)
		})),
	)
}

func navProjects(v08 bool) h.Node {
	contentClass, share, moreBtn, moreIcon := "w-48 rounded-lg", lucide.Forward, []any{h.Class("text-sidebar-foreground/70")}, []any{h.Class("text-sidebar-foreground/70")}
	if v08 {
		contentClass, share, moreBtn, moreIcon = "w-48", lucide.Share, nil, nil
	}
	return sidebar.Group(h.Class("group-data-[collapsible=icon]:hidden"),
		sidebar.GroupLabel("Projects"),
		sidebar.Menu(
			h.Map(projects, func(p project) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), p.icon(), h.Span(p.name))),
					dm.DropdownMenu(
						dm.Trigger(lib.AsChild,
							sidebar.MenuAction(sidebar.ShowOnHover(true), lucide.Ellipsis(), h.Span(h.Class("sr-only"), "More")),
						),
						dm.Content(h.Class(contentClass), dm.SideRight, dm.AlignStart,
							dm.Item(lucide.Folder(h.Class("text-muted-foreground")), h.Span("View Project")),
							dm.Item(share(h.Class("text-muted-foreground")), h.Span("Share Project")),
							dm.Separator(),
							dm.Item(lucide.Trash2(h.Class("text-muted-foreground")), h.Span("Delete Project")),
						),
					),
				)
			}),
			sidebar.MenuItem(sidebar.MenuButton(append(moreBtn, lucide.Ellipsis(moreIcon...), h.Span("More"))...)),
		),
	)
}

func userBlock() h.Node {
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

func navUser() h.Node {
	return sidebar.Menu(sidebar.MenuItem(dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			sidebar.MenuButton(sidebar.ButtonSizeLg,
				h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
				userBlock(),
				lucide.ChevronsUpDown(h.Class("ml-auto size-4")),
			),
		),
		dm.Content(h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"), dm.SideRight, dm.AlignEnd, dm.SideOffset(4),
			dm.Label(h.Class("p-0 font-normal"),
				h.Div(h.Class("flex items-center gap-2 px-1 py-1.5 text-left text-sm"), userBlock()),
			),
			dm.Separator(),
			dm.Group(dm.Item(lucide.Sparkles(), "Upgrade to Pro")),
			dm.Separator(),
			dm.Group(
				dm.Item(lucide.BadgeCheck(), "Account"),
				dm.Item(lucide.CreditCard(), "Billing"),
				dm.Item(lucide.Bell(), "Notifications"),
			),
			dm.Separator(),
			dm.Item(lucide.LogOut(), "Log out"),
		),
	)))
}

// teamSwitcher is TeamSwitcher (switcher.js keeps the active team).
func teamSwitcher() h.Node {
	active := teams[0]
	var rows []any
	rows = append(rows, dm.Label(h.Class("text-xs text-muted-foreground"), "Teams"))
	for i, t := range teams {
		rows = append(rows, dm.Item(h.Class("gap-2 p-2"),
			h.Div(h.Class("flex size-6 items-center justify-center rounded-md border"), t.logo(h.Class("size-3.5 shrink-0"))),
			t.name,
			dm.Shortcut("⌘", h.Comment(""), fmt.Sprint(i+1)),
		))
	}
	rows = append(rows, dm.Separator(), dm.Item(h.Class("gap-2 p-2"),
		h.Div(h.Class("flex size-6 items-center justify-center rounded-md border bg-transparent"), lucide.Plus(h.Class("size-4"))),
		h.Div(h.Class("font-medium text-muted-foreground"), "Add team"),
	))
	var logos []any
	for i, t := range teams {
		logos = append(logos, h.Template(h.Data("parity-ignore", ""), h.Data("team-logo", fmt.Sprint(i)), h.Data("plan", t.plan), t.logo(h.Class("size-4"))))
	}
	return sidebar.Menu(sidebar.MenuItem(
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild,
				sidebar.MenuButton(sidebar.ButtonSizeLg,
					h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
					h.Data("example-blocks-sidebar-switcher", "team"),
					logoBox(active.logo(h.Class("size-4"))),
					h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
						h.Span(h.Class("truncate font-medium"), active.name),
						h.Span(h.Class("truncate text-xs"), active.plan),
					),
					lucide.ChevronsUpDown(h.Class("ml-auto")),
				),
			),
			dm.Content(append([]any{h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"), dm.AlignStart, dm.SideRight, dm.SideOffset(4)}, rows...)...),
		),
		logos,
	))
}

// Sidebar07 is sidebar-07: a sidebar that collapses to icons.
func Sidebar07() h.Node {
	app := sidebar.Sidebar(sidebar.CollapsibleIcon,
		sidebar.Header(teamSwitcher()),
		sidebar.Content(navMain07(), navProjects(false)),
		sidebar.Footer(navUser()),
		sidebar.Rail(),
	)
	return page(app,
		headerNested("flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"),
		gridBody("flex flex-1 flex-col gap-4 p-4 pt-0"))
}

// Sidebar08 is sidebar-08: an inset sidebar with secondary navigation.
func Sidebar08() h.Node {
	secondary := sidebar.Group(h.Class("mt-auto"), sidebar.GroupContent(sidebar.Menu(
		sidebar.MenuItem(sidebar.MenuButton(lib.AsChild, sidebar.ButtonSizeSm, h.A(h.Href("#"), lucide.LifeBuoy(), h.Span("Support")))),
		sidebar.MenuItem(sidebar.MenuButton(lib.AsChild, sidebar.ButtonSizeSm, h.A(h.Href("#"), lucide.Send(), h.Span("Feedback")))),
	)))
	app := sidebar.Sidebar(sidebar.VariantInset,
		sidebar.Header(sidebar.Menu(sidebar.MenuItem(
			sidebar.MenuButton(sidebar.ButtonSizeLg, lib.AsChild,
				h.A(h.Href("#"),
					logoBox(lucide.Command(h.Class("size-4"))),
					h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
						h.Span(h.Class("truncate font-medium"), "Acme Inc"),
						h.Span(h.Class("truncate text-xs"), "Enterprise"),
					),
				),
			),
		))),
		sidebar.Content(navMain08(), navProjects(true), secondary),
		sidebar.Footer(navUser()),
	)
	return page(app, headerNested("flex h-16 shrink-0 items-center gap-2"), gridBody("flex flex-1 flex-col gap-4 p-4 pt-0"))
}
