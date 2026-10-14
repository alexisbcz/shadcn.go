package blockssidebarb

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/components/ui/dialog"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	switchui "shadcn.go/components/ui/switch"
	h "shadcn.go/libhtml"
)

func repeat(n int, fn func(i int) h.Node) h.Node {
	out := make([]any, n)
	for i := range n {
		out[i] = fn(i)
	}
	return h.Group(out...)
}

func tile(class string) func(int) h.Node {
	return func(int) h.Node { return h.Div(h.Class(class)) }
}

// threeTiles is the common page body: three video tiles and a tall panel.
func threeTiles() h.Node {
	return h.Div(h.Class("flex flex-1 flex-col gap-4 p-4"),
		h.Div(h.Class("grid auto-rows-min gap-4 md:grid-cols-3"),
			repeat(3, tile("aspect-video rounded-xl bg-muted/50")),
		),
		h.Div(h.Class("min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min")),
	)
}

// --- sidebar-09 ------------------------------------------------------------------------

type mail struct{ name, email, subject, date, teaser string }

var mails = []mail{
	{"William Smith", "williamsmith@example.com", "Meeting Tomorrow", "09:34 AM", "Hi team, just a reminder about our meeting tomorrow at 10 AM.\nPlease come prepared with your project updates."},
	{"Alice Smith", "alicesmith@example.com", "Re: Project Update", "Yesterday", "Thanks for the update. The progress looks great so far.\nLet's schedule a call to discuss the next steps."},
	{"Bob Johnson", "bobjohnson@example.com", "Weekend Plans", "2 days ago", "Hey everyone! I'm thinking of organizing a team outing this weekend.\nWould you be interested in a hiking trip or a beach day?"},
	{"Emily Davis", "emilydavis@example.com", "Re: Question about Budget", "2 days ago", "I've reviewed the budget numbers you sent over.\nCan we set up a quick call to discuss some potential adjustments?"},
	{"Michael Wilson", "michaelwilson@example.com", "Important Announcement", "1 week ago", "Please join us for an all-hands meeting this Friday at 3 PM.\nWe have some exciting news to share about the company's future."},
	{"Sarah Brown", "sarahbrown@example.com", "Re: Feedback on Proposal", "1 week ago", "Thank you for sending over the proposal. I've reviewed it and have some thoughts.\nCould we schedule a meeting to discuss my feedback in detail?"},
	{"David Lee", "davidlee@example.com", "New Project Idea", "1 week ago", "I've been brainstorming and came up with an interesting project concept.\nDo you have time this week to discuss its potential impact and feasibility?"},
	{"Olivia Wilson", "oliviawilson@example.com", "Vacation Plans", "1 week ago", "Just a heads up that I'll be taking a two-week vacation next month.\nI'll make sure all my projects are up to date before I leave."},
	{"James Martin", "jamesmartin@example.com", "Re: Conference Registration", "1 week ago", "I've completed the registration for the upcoming tech conference.\nLet me know if you need any additional information from my end."},
	{"Sophia White", "sophiawhite@example.com", "Team Dinner", "1 week ago", "To celebrate our recent project success, I'd like to organize a team dinner.\nAre you available next Friday evening? Please let me know your preferences."},
}

var mailNav = []navItem{
	{"Inbox", lucide.Inbox, true, ""},
	{"Drafts", lucide.File, false, ""},
	{"Sent", lucide.Send, false, ""},
	{"Junk", lucide.ArchiveX, false, ""},
	{"Trash", lucide.Trash2, false, ""},
}

// Sidebar09 is sidebar-09: collapsible nested sidebars (blocks_sidebar_b.js
// switches the active folder and shuffles the mails like upstream).
func Sidebar09() h.Node {
	app := sidebar.Sidebar(sidebar.CollapsibleIcon, h.Class("overflow-hidden *:data-[sidebar=sidebar]:flex-row"),
		h.Data("example-sidebar-09", ""),
		sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("w-[calc(var(--sidebar-width-icon)+1px)]! border-r"),
			sidebar.Header(
				sidebar.Menu(
					sidebar.MenuItem(
						sidebar.MenuButton(sidebar.ButtonSizeLg, lib.AsChild, h.Class("md:h-8 md:p-0"),
							h.A(h.Href("#"),
								h.Div(h.Class("flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"),
									lucide.Command(h.Class("size-4")),
								),
								h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
									h.Span(h.Class("truncate font-medium"), "Acme Inc"),
									h.Span(h.Class("truncate text-xs"), "Enterprise"),
								),
							),
						),
					),
				),
			),
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupContent(h.Class("px-1.5 md:px-0"),
						sidebar.Menu(h.Map(mailNav, func(it navItem) h.Node {
							return sidebar.MenuItem(
								sidebar.MenuButton(
									sidebar.Tooltip(sidebar.TooltipHidden(false), it.title),
									sidebar.IsActive(it.active),
									h.Class("px-2.5 md:px-2"),
									h.Data("mail-nav", it.title),
									it.icon(), h.Span(it.title),
								),
							)
						})),
					),
				),
			),
			sidebar.Footer(navUser(shadcnUser, true, dm.AlignEnd)),
		),
		sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("hidden flex-1 md:flex"),
			sidebar.Header(h.Class("gap-3.5 border-b p-4"),
				h.Div(h.Class("flex w-full items-center justify-between"),
					h.Div(h.Class("text-base font-medium text-foreground"), h.Data("mail-title", ""), "Inbox"),
					label.Label(h.Class("flex items-center gap-2 text-sm"),
						h.Span("Unreads"),
						switchui.Switch(h.Class("shadow-none")),
					),
				),
				sidebar.Input(h.Placeholder("Type to search...")),
			),
			sidebar.Content(
				sidebar.Group(h.Class("px-0"),
					sidebar.GroupContent(h.Data("mail-list", ""),
						h.Map(mails, func(m mail) h.Node {
							return h.A(h.Href("#"), h.Data("mail", m.email),
								h.Class("flex flex-col items-start gap-2 border-b p-4 text-sm leading-tight whitespace-nowrap last:border-b-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"),
								h.Div(h.Class("flex w-full items-center gap-2"),
									h.Span(m.name), " ",
									h.Span(h.Class("ml-auto text-xs"), m.date),
								),
								h.Span(h.Class("font-medium"), m.subject),
								h.Span(h.Class("line-clamp-2 w-[260px] text-xs whitespace-break-spaces"), m.teaser),
							)
						}),
					),
				),
			),
		),
	)
	return sidebar.Provider(h.Style("--sidebar-width: 350px"),
		app,
		sidebar.Inset(
			h.Header(h.Class("sticky top-0 flex shrink-0 items-center gap-2 border-b bg-background p-4"),
				sidebar.Trigger(h.Class("-ml-1")),
				vsep(),
				breadcrumb.Breadcrumb(breadcrumb.List(
					crumbLink("All Inboxes", true),
					crumbSep(),
					crumbPage("Inbox"),
				)),
			),
			h.Div(h.Class("flex flex-1 flex-col gap-4 p-4"),
				repeat(24, tile("aspect-video h-12 w-full rounded-lg bg-muted/50")),
			),
		),
	)
}

// --- sidebar-10 ------------------------------------------------------------------------

type action struct {
	label string
	icon  icon
}

var actions = [][]action{
	{{"Customize Page", lucide.Settings2}, {"Turn into wiki", lucide.FileText}},
	{{"Copy Link", lucide.Link}, {"Duplicate", lucide.Copy}, {"Move to", lucide.CornerUpRight}, {"Move to Trash", lucide.Trash2}},
	{{"Undo", lucide.CornerUpLeft}, {"View analytics", lucide.ChartLine}, {"Version History", lucide.GalleryVerticalEnd}, {"Show delete pages", lucide.Trash}, {"Notifications", lucide.Bell}},
	{{"Import", lucide.ArrowUp}, {"Export", lucide.ArrowDown}},
}

func navActions() h.Node {
	return h.Div(h.Class("flex items-center gap-2 text-sm"),
		h.Div(h.Class("hidden font-medium text-muted-foreground md:inline-block"), "Edit Oct 08"),
		button.Button(button.VariantGhost, button.SizeIcon, h.Class("h-7 w-7"), lucide.Star()),
		popover.Popover(popover.DefaultOpen,
			popover.Trigger(lib.AsChild,
				button.Button(button.VariantGhost, button.SizeIcon, h.Class("h-7 w-7 data-[state=open]:bg-accent"), lucide.Ellipsis()),
			),
			popover.Content(h.Class("w-56 overflow-hidden rounded-lg p-0"), popover.AlignEnd,
				sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("bg-transparent"),
					sidebar.Content(h.Map(actions, func(group []action) h.Node {
						return sidebar.Group(h.Class("border-b last:border-none"),
							sidebar.GroupContent(h.Class("gap-0"),
								sidebar.Menu(h.Map(group, func(a action) h.Node {
									return sidebar.MenuItem(sidebar.MenuButton(a.icon(), " ", h.Span(a.label)))
								})),
							),
						)
					})),
				),
			),
		),
	)
}

// Sidebar10 is sidebar-10.
func Sidebar10() h.Node {
	return sidebar.Provider(
		leftSidebar(),
		sidebar.Inset(
			h.Header(h.Class("flex h-14 shrink-0 items-center gap-2"),
				h.Div(h.Class("flex flex-1 items-center gap-2 px-3"),
					sidebar.Trigger(),
					vsep(),
					breadcrumb.Breadcrumb(breadcrumb.List(
						crumbPage("Project Management & Task Tracking", h.Class("line-clamp-1")),
					)),
				),
				h.Div(h.Class("ml-auto px-3"), navActions()),
			),
			h.Div(h.Class("flex flex-1 flex-col gap-4 px-4 py-10"),
				h.Div(h.Class("mx-auto h-24 w-full max-w-3xl rounded-xl bg-muted/50")),
				h.Div(h.Class("mx-auto h-full w-full max-w-3xl rounded-xl bg-muted/50")),
			),
		),
	)
}

// --- sidebar-11 ------------------------------------------------------------------------

type treeItem struct {
	name  string
	items []treeItem
}

func leaf(name string) treeItem                   { return treeItem{name: name} }
func dir(name string, items ...treeItem) treeItem { return treeItem{name: name, items: items} }

var fileTree = []treeItem{
	// Upstream nests page.tsx, layout.tsx and blog inside api.
	dir("app", dir("api", dir("hello", leaf("route.ts")), leaf("page.tsx"), leaf("layout.tsx"), dir("blog", leaf("page.tsx")))),
	dir("components", dir("ui", leaf("button.tsx"), leaf("card.tsx")), leaf("header.tsx"), leaf("footer.tsx")),
	dir("lib", leaf("util.ts")),
	dir("public", leaf("favicon.ico"), leaf("vercel.svg")),
	leaf(".eslintrc.json"),
	leaf(".gitignore"),
	leaf("next.config.js"),
	leaf("tailwind.config.js"),
	leaf("package.json"),
	leaf("README.md"),
}

func tree(it treeItem) h.Node {
	if len(it.items) == 0 {
		return sidebar.MenuButton(sidebar.IsActive(it.name == "button.tsx"), h.Class("data-[active=true]:bg-transparent"),
			lucide.File(), it.name,
		)
	}
	return sidebar.MenuItem(
		collapsible.Collapsible(
			h.Class("group/collapsible [&[data-state=open]>button>svg:first-child]:rotate-90"),
			collapsible.Open(it.name == "components" || it.name == "ui"),
			collapsible.Trigger(lib.AsChild,
				sidebar.MenuButton(lucide.ChevronRight(h.Class("transition-transform")), lucide.Folder(), it.name),
			),
			collapsible.Content(
				sidebar.MenuSub(h.Map(it.items, tree)),
			),
		),
	)
}

// Sidebar11 is sidebar-11: a file tree.
func Sidebar11() h.Node {
	changes := []struct{ file, state string }{{"README.md", "M"}, {"api/hello/route.ts", "U"}, {"app/layout.tsx", "M"}}
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Changes"),
					sidebar.GroupContent(
						sidebar.Menu(h.Map(changes, func(c struct{ file, state string }) h.Node {
							return sidebar.MenuItem(
								sidebar.MenuButton(lucide.File(), c.file),
								sidebar.MenuBadge(c.state),
							)
						})),
					),
				),
				sidebar.Group(
					sidebar.GroupLabel("Files"),
					sidebar.GroupContent(sidebar.Menu(h.Map(fileTree, tree))),
				),
			),
			sidebar.Rail(),
		),
		sidebar.Inset(
			h.Header(h.Class("flex h-16 shrink-0 items-center gap-2 border-b px-4"),
				sidebar.Trigger(h.Class("-ml-1")),
				vsep(),
				breadcrumb.Breadcrumb(breadcrumb.List(
					crumbLink("components", true),
					crumbSep(),
					crumbLink("ui", true),
					crumbSep(),
					crumbPage("button.tsx"),
				)),
			),
			threeTiles(),
		),
	)
}

// --- sidebar-12 ------------------------------------------------------------------------

// Sidebar12 is sidebar-12: a calendar sidebar.
func Sidebar12() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(append(calendarSidebar("rounded-sm"), sidebar.Rail())...),
		sidebar.Inset(
			h.Header(h.Class("sticky top-0 flex h-16 shrink-0 items-center gap-2 border-b bg-background px-4"),
				sidebar.Trigger(h.Class("-ml-1")),
				vsep(),
				breadcrumb.Breadcrumb(breadcrumb.List(crumbPage("October 2024"))),
			),
			h.Div(h.Class("flex flex-1 flex-col gap-4 p-4"),
				h.Div(h.Class("grid auto-rows-min gap-4 md:grid-cols-5"),
					repeat(20, tile("aspect-square rounded-xl bg-muted/50")),
				),
			),
		),
	)
}

// --- sidebar-13 ------------------------------------------------------------------------

// Sidebar13 is sidebar-13: a settings dialog with a sidebar.
func Sidebar13() h.Node {
	nav := []navItem{
		{"Notifications", lucide.Bell, false, ""},
		{"Navigation", lucide.Menu, false, ""},
		{"Home", lucide.House, false, ""},
		{"Appearance", lucide.Paintbrush, false, ""},
		{"Messages & media", lucide.MessageCircle, true, ""},
		{"Language & region", lucide.Globe, false, ""},
		{"Accessibility", lucide.Keyboard, false, ""},
		{"Mark as read", lucide.Check, false, ""},
		{"Audio & video", lucide.Video, false, ""},
		{"Connected accounts", lucide.Link, false, ""},
		{"Privacy & visibility", lucide.Lock, false, ""},
		{"Advanced", lucide.Settings, false, ""},
	}
	return h.Div(h.Class("flex h-svh items-center justify-center"),
		dialog.Dialog(dialog.DefaultOpen,
			dialog.Trigger(lib.AsChild, button.Button(button.SizeSm, "Open Dialog")),
			dialog.Content(h.Class("overflow-hidden p-0 md:max-h-[500px] md:max-w-[700px] lg:max-w-[800px]"),
				dialog.Title(h.Class("sr-only"), "Settings"),
				dialog.Description(h.Class("sr-only"), "Customize your settings here."),
				sidebar.Provider(h.Class("items-start"),
					sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("hidden md:flex"),
						sidebar.Content(
							sidebar.Group(
								sidebar.GroupContent(navMain(nav)),
							),
						),
					),
					h.Main(h.Class("flex h-[480px] flex-1 flex-col overflow-hidden"),
						h.Header(h.Class("flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"),
							h.Div(h.Class("flex items-center gap-2 px-4"),
								breadcrumb.Breadcrumb(breadcrumb.List(
									crumbLink("Settings", true),
									crumbSep(),
									crumbPage("Messages & media"),
								)),
							),
						),
						h.Div(h.Class("flex flex-1 flex-col gap-4 overflow-y-auto p-4 pt-0"),
							repeat(10, tile("aspect-video max-w-3xl rounded-xl bg-muted/50")),
						),
					),
				),
			),
		),
	)
}

// --- sidebar-14 ------------------------------------------------------------------------

type docSection struct {
	title  string
	items  []string
	active string
}

var toc = []docSection{
	{"Getting Started", []string{"Installation", "Project Structure"}, ""},
	{"Build Your Application", []string{"Routing", "Data Fetching", "Rendering", "Caching", "Styling", "Optimizing", "Configuring", "Testing", "Authentication", "Deploying", "Upgrading", "Examples"}, "Data Fetching"},
	{"API Reference", []string{"Components", "File Conventions", "Functions", "next.config.js Options", "CLI", "Edge Runtime"}, ""},
	{"Architecture", []string{"Accessibility", "Fast Refresh", "Next.js Compiler", "Supported Browsers", "Turbopack"}, ""},
	{"Community", []string{"Contribution Guide"}, ""},
}

// Sidebar14 is sidebar-14: a table of contents on the right.
func Sidebar14() h.Node {
	return sidebar.Provider(
		sidebar.Inset(
			h.Header(h.Class("flex h-16 shrink-0 items-center gap-2 border-b px-4"),
				breadcrumb.Breadcrumb(breadcrumb.List(
					crumbLink("Build Your Application", true),
					crumbSep(),
					crumbPage("Data Fetching"),
				)),
				sidebar.Trigger(h.Class("-mr-1 ml-auto rotate-180")),
			),
			threeTiles(),
		),
		sidebar.Sidebar(sidebar.SideRight,
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Table of Contents"),
					sidebar.GroupContent(
						sidebar.Menu(h.Map(toc, func(s docSection) h.Node {
							return sidebar.MenuItem(
								sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), h.Class("font-medium"), s.title)),
								sidebar.MenuSub(h.Map(s.items, func(t string) h.Node {
									return sidebar.MenuSubItem(
										sidebar.MenuSubButton(lib.AsChild, sidebar.IsActive(t == s.active), h.A(h.Href("#"), t)),
									)
								})),
							)
						})),
					),
				),
			),
			sidebar.Rail(),
		),
	)
}

// --- sidebar-15 ------------------------------------------------------------------------

// Sidebar15 is sidebar-15: a left and a right sidebar.
func Sidebar15() h.Node {
	return sidebar.Provider(
		leftSidebar(),
		sidebar.Inset(
			h.Header(h.Class("sticky top-0 flex h-14 shrink-0 items-center gap-2 bg-background"),
				h.Div(h.Class("flex flex-1 items-center gap-2 px-3"),
					sidebar.Trigger(),
					vsep(),
					breadcrumb.Breadcrumb(breadcrumb.List(
						crumbPage("Project Management & Task Tracking", h.Class("line-clamp-1")),
					)),
				),
			),
			h.Div(h.Class("flex flex-1 flex-col gap-4 p-4"),
				h.Div(h.Class("mx-auto h-24 w-full max-w-3xl rounded-xl bg-muted/50")),
				h.Div(h.Class("mx-auto h-[100vh] w-full max-w-3xl rounded-xl bg-muted/50")),
			),
		),
		sidebar.Sidebar(append([]any{sidebar.CollapsibleNone, h.Class("sticky top-0 hidden h-svh border-l lg:flex")},
			calendarSidebar("rounded-xs")...)...),
	)
}

// --- sidebar-16 ------------------------------------------------------------------------

type navGroup struct {
	title  string
	icon   icon
	active bool
	items  []string
}

var platform = []navGroup{
	{"Playground", lucide.SquareTerminal, true, []string{"History", "Starred", "Settings"}},
	{"Models", lucide.Bot, false, []string{"Genesis", "Explorer", "Quantum"}},
	{"Documentation", lucide.BookOpen, false, []string{"Introduction", "Get Started", "Tutorials", "Changelog"}},
	{"Settings", lucide.Settings2, false, []string{"General", "Team", "Billing", "Limits"}},
}

func platformNav() h.Node {
	return sidebar.Group(
		sidebar.GroupLabel("Platform"),
		sidebar.Menu(h.Map(platform, func(g navGroup) h.Node {
			return collapsible.Collapsible(lib.AsChild, collapsible.Open(g.active),
				sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, sidebar.Tooltip(g.title),
						h.A(h.Href("#"), g.icon(), h.Span(g.title)),
					),
					collapsible.Trigger(lib.AsChild,
						sidebar.MenuAction(h.Class("data-[state=open]:rotate-90"),
							lucide.ChevronRight(), h.Span(h.Class("sr-only"), "Toggle"),
						),
					),
					collapsible.Content(
						sidebar.MenuSub(h.Map(g.items, func(t string) h.Node {
							return sidebar.MenuSubItem(
								sidebar.MenuSubButton(lib.AsChild, h.A(h.Href("#"), h.Span(t))),
							)
						})),
					),
				),
			)
		})),
	)
}

func navProjects() h.Node {
	projects := []navItem{
		{"Design Engineering", lucide.Frame, false, ""},
		{"Sales & Marketing", lucide.ChartPie, false, ""},
		{"Travel", lucide.Map, false, ""},
	}
	return sidebar.Group(h.Class("group-data-[collapsible=icon]:hidden"),
		sidebar.GroupLabel("Projects"),
		sidebar.Menu(
			h.Map(projects, func(p navItem) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), p.icon(), h.Span(p.title))),
					dm.DropdownMenu(
						dm.Trigger(lib.AsChild,
							sidebar.MenuAction(sidebar.ShowOnHover(true),
								lucide.Ellipsis(), h.Span(h.Class("sr-only"), "More"),
							),
						),
						dm.Content(h.Class("w-48"), dm.SideRight, dm.AlignStart,
							dm.Item(lucide.Folder(h.Class("text-muted-foreground")), h.Span("View Project")),
							dm.Item(lucide.Share(h.Class("text-muted-foreground")), h.Span("Share Project")),
							dm.Separator(),
							dm.Item(lucide.Trash2(h.Class("text-muted-foreground")), h.Span("Delete Project")),
						),
					),
				)
			}),
			sidebar.MenuItem(sidebar.MenuButton(lucide.Ellipsis(), h.Span("More"))),
		),
	)
}

// Sidebar16 is sidebar-16: a sticky site header above the sidebar.
func Sidebar16() h.Node {
	secondary := []navItem{
		{"Support", lucide.LifeBuoy, false, ""},
		{"Feedback", lucide.Send, false, ""},
	}
	siteHeader := h.Header(h.Class("sticky top-0 z-50 flex w-full items-center border-b bg-background"),
		h.Div(h.Class("flex h-(--header-height) w-full items-center gap-2 px-4"),
			button.Button(h.Class("h-8 w-8"), button.VariantGhost, button.SizeIcon, h.Data("example-sidebar-toggle", ""),
				lucide.PanelLeft(),
			),
			separator.Separator(separator.Vertical, h.Class("mr-2 h-4")),
			breadcrumb.Breadcrumb(h.Class("hidden sm:block"), breadcrumb.List(
				breadcrumb.Item(breadcrumb.Link(h.Href("#"), "Build Your Application")),
				breadcrumb.Separator(),
				crumbPage("Data Fetching"),
			)),
			h.Form(h.Class("w-full sm:ml-auto sm:w-auto"),
				h.Div(h.Class("relative"),
					label.Label(h.For("search"), h.Class("sr-only"), "Search"),
					sidebar.Input(h.ID("search"), h.Placeholder("Type to search..."), h.Class("h-8 pl-7")),
					lucide.Search(h.Class("pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none")),
				),
			),
		),
	)
	app := sidebar.Sidebar(h.Class("top-(--header-height) h-[calc(100svh-var(--header-height))]!"),
		sidebar.Header(
			sidebar.Menu(
				sidebar.MenuItem(
					sidebar.MenuButton(sidebar.ButtonSizeLg, lib.AsChild,
						h.A(h.Href("#"),
							h.Div(h.Class("flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"),
								lucide.Command(h.Class("size-4")),
							),
							h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
								h.Span(h.Class("truncate font-medium"), "Acme Inc"),
								h.Span(h.Class("truncate text-xs"), "Enterprise"),
							),
						),
					),
				),
			),
		),
		sidebar.Content(
			platformNav(),
			navProjects(),
			navSecondary(secondary, sidebar.ButtonSizeSm, h.Class("mt-auto")),
		),
		sidebar.Footer(navUser(shadcnUser, false, dm.AlignEnd)),
	)
	return h.Div(h.Class("[--header-height:calc(--spacing(14))]"),
		sidebar.Provider(h.Class("flex flex-col"),
			siteHeader,
			h.Div(h.Class("flex flex-1"),
				app,
				sidebar.Inset(threeTiles()),
			),
		),
	)
}
