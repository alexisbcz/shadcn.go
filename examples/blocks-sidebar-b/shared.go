// Package blockssidebarb ports the upstream sidebar-09 to sidebar-16
// blocks (registry/new-york-v4/blocks/sidebar-*): full pages built on
// components/ui/sidebar.
package blockssidebarb

import (
	_ "embed"
	"fmt"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/collapsible"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed blocks_sidebar_b.js
var script string

func init() {
	runtime.Module("example-blocks-sidebar-b", script)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "sidebar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("sidebar-09", Sidebar09),
		ex("sidebar-10", Sidebar10),
		ex("sidebar-11", Sidebar11),
		ex("sidebar-12", Sidebar12),
		ex("sidebar-13", Sidebar13),
		ex("sidebar-14", Sidebar14),
		ex("sidebar-15", Sidebar15),
		ex("sidebar-16", Sidebar16),
	)
}

type icon = func(...any) *h.Element

type user struct{ name, email, avatar string }

var shadcnUser = user{"shadcn", "m@example.com", "/avatars/shadcn.jpg"}

func userAvatar(u user) h.Node {
	return avatar.Avatar(h.Class("h-8 w-8 rounded-lg"),
		avatar.Image(h.Src(u.avatar), h.Alt(u.name)),
		avatar.Fallback(h.Class("rounded-lg"), "CN"),
	)
}

func userText(u user) h.Node {
	return h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
		h.Span(h.Class("truncate font-medium"), u.name),
		h.Span(h.Class("truncate text-xs"), u.email),
	)
}

// navUser is NavUser; sidebar-09 and sidebar-16 use the compact
// (md:h-8 md:p-0, align end) variant, sidebar-12 and -15 align start.
func navUser(u user, compact bool, align dm.Align) h.Node {
	cls := "data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
	if compact {
		cls += " md:h-8 md:p-0"
	}
	return sidebar.Menu(
		sidebar.MenuItem(
			dm.DropdownMenu(
				dm.Trigger(lib.AsChild,
					sidebar.MenuButton(sidebar.ButtonSizeLg, h.Class(cls),
						userAvatar(u),
						userText(u),
						lucide.ChevronsUpDown(h.Class("ml-auto size-4")),
					),
				),
				dm.Content(h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"),
					dm.SideRight, align, dm.SideOffset(4),
					dm.Label(h.Class("p-0 font-normal"),
						h.Div(h.Class("flex items-center gap-2 px-1 py-1.5 text-left text-sm"),
							userAvatar(u),
							userText(u),
						),
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
			),
		),
	)
}

// vsep is the header's vertical separator.
func vsep() h.Node {
	return separator.Separator(separator.Vertical, h.Class("mr-2 data-[orientation=vertical]:h-4"))
}

func crumbLink(label string, hideMobile bool) h.Node {
	return breadcrumb.Item(h.If(hideMobile, h.Class("hidden md:block")), breadcrumb.Link(h.Href("#"), label))
}

func crumbSep() h.Node { return breadcrumb.Separator(h.Class("hidden md:block")) }

func crumbPage(label string, extra ...any) h.Node {
	return breadcrumb.Item(breadcrumb.Page(append(extra, label)...))
}

// --- sidebar-10 / sidebar-15 left sidebar ---------------------------------------------

type team struct {
	name string
	logo icon
	plan string
}

var teams = []team{
	{"Acme Inc", lucide.Command, "Enterprise"},
	{"Acme Corp.", lucide.AudioWaveform, "Startup"},
	{"Evil Corp.", lucide.Command, "Free"},
}

func teamSwitcher() h.Node {
	active := teams[0]
	var items []any
	for i, t := range teams {
		items = append(items, dm.Item(h.Class("gap-2 p-2"),
			h.Div(h.Class("flex size-6 items-center justify-center rounded-xs border"),
				t.logo(h.Class("size-4 shrink-0")),
			),
			t.name,
			dm.Shortcut("⌘", fmt.Sprint(i+1)),
		))
	}
	return sidebar.Menu(
		sidebar.MenuItem(
			dm.DropdownMenu(
				dm.Trigger(lib.AsChild,
					sidebar.MenuButton(h.Class("w-fit px-1.5"),
						h.Div(h.Class("flex aspect-square size-5 items-center justify-center rounded-md bg-sidebar-primary text-sidebar-primary-foreground"),
							active.logo(h.Class("size-3")),
						),
						h.Span(h.Class("truncate font-medium"), active.name),
						lucide.ChevronDown(h.Class("opacity-50")),
					),
				),
				dm.Content(append([]any{h.Class("w-64 rounded-lg"), dm.AlignStart, dm.SideBottom, dm.SideOffset(4),
					dm.Label(h.Class("text-xs text-muted-foreground"), "Teams"),
				}, append(items,
					dm.Separator(),
					dm.Item(h.Class("gap-2 p-2"),
						h.Div(h.Class("flex size-6 items-center justify-center rounded-md border bg-background"),
							lucide.Plus(h.Class("size-4")),
						),
						h.Div(h.Class("font-medium text-muted-foreground"), "Add team"),
					),
				)...)...),
			),
		),
	)
}

type navItem struct {
	title  string
	icon   icon
	active bool
	badge  string
}

var navMainItems = []navItem{
	{"Search", lucide.Search, false, ""},
	{"Ask AI", lucide.Sparkles, false, ""},
	{"Home", lucide.House, true, ""},
	{"Inbox", lucide.Inbox, false, "10"},
}

var navSecondaryItems = []navItem{
	{"Calendar", lucide.Calendar, false, ""},
	{"Settings", lucide.Settings2, false, ""},
	{"Templates", lucide.Blocks, false, ""},
	{"Trash", lucide.Trash2, false, ""},
	{"Help", lucide.MessageCircleQuestion, false, ""},
}

func navMain(items []navItem) h.Node {
	return sidebar.Menu(h.Map(items, func(it navItem) h.Node {
		return sidebar.MenuItem(
			sidebar.MenuButton(lib.AsChild, sidebar.IsActive(it.active),
				h.A(h.Href("#"), it.icon(), h.Span(it.title)),
			),
		)
	}))
}

func navSecondary(items []navItem, size sidebar.ButtonSize, props ...any) h.Node {
	return sidebar.Group(append(props,
		sidebar.GroupContent(
			sidebar.Menu(h.Map(items, func(it navItem) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, size,
						h.A(h.Href("#"), it.icon(), h.Span(it.title)),
					),
					h.If(it.badge != "", sidebar.MenuBadge(it.badge)),
				)
			})),
		),
	)...)
}

type emojiLink struct{ name, emoji string }

var favorites = []emojiLink{
	{"Project Management & Task Tracking", "📊"},
	{"Family Recipe Collection & Meal Planning", "🍳"},
	{"Fitness Tracker & Workout Routines", "💪"},
	{"Book Notes & Reading List", "📚"},
	{"Sustainable Gardening Tips & Plant Care", "🌱"},
	{"Language Learning Progress & Resources", "🗣️"},
	{"Home Renovation Ideas & Budget Tracker", "🏠"},
	{"Personal Finance & Investment Portfolio", "💰"},
	{"Movie & TV Show Watchlist with Reviews", "🎬"},
	{"Daily Habit Tracker & Goal Setting", "✅"},
}

func navFavorites() h.Node {
	return sidebar.Group(h.Class("group-data-[collapsible=icon]:hidden"),
		sidebar.GroupLabel("Favorites"),
		sidebar.Menu(
			h.Map(favorites, func(f emojiLink) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild,
						h.A(h.Href("#"), h.TitleAttr(f.name), h.Span(f.emoji), h.Span(f.name)),
					),
					dm.DropdownMenu(
						dm.Trigger(lib.AsChild,
							sidebar.MenuAction(sidebar.ShowOnHover(true),
								lucide.Ellipsis(), h.Span(h.Class("sr-only"), "More"),
							),
						),
						dm.Content(h.Class("w-56 rounded-lg"), dm.SideRight, dm.AlignStart,
							dm.Item(lucide.StarOff(h.Class("text-muted-foreground")), h.Span("Remove from Favorites")),
							dm.Separator(),
							dm.Item(lucide.Link(h.Class("text-muted-foreground")), h.Span("Copy Link")),
							dm.Item(lucide.ArrowUpRight(h.Class("text-muted-foreground")), h.Span("Open in New Tab")),
							dm.Separator(),
							dm.Item(lucide.Trash2(h.Class("text-muted-foreground")), h.Span("Delete")),
						),
					),
				)
			}),
			sidebar.MenuItem(
				sidebar.MenuButton(h.Class("text-sidebar-foreground/70"), lucide.Ellipsis(), h.Span("More")),
			),
		),
	)
}

type workspace struct {
	name, emoji string
	pages       []emojiLink
}

var workspaces = []workspace{
	{"Personal Life Management", "🏠", []emojiLink{{"Daily Journal & Reflection", "📔"}, {"Health & Wellness Tracker", "🍏"}, {"Personal Growth & Learning Goals", "🌟"}}},
	{"Professional Development", "💼", []emojiLink{{"Career Objectives & Milestones", "🎯"}, {"Skill Acquisition & Training Log", "🧠"}, {"Networking Contacts & Events", "🤝"}}},
	{"Creative Projects", "🎨", []emojiLink{{"Writing Ideas & Story Outlines", "✍️"}, {"Art & Design Portfolio", "🖼️"}, {"Music Composition & Practice Log", "🎵"}}},
	{"Home Management", "🏡", []emojiLink{{"Household Budget & Expense Tracking", "💰"}, {"Home Maintenance Schedule & Tasks", "🔧"}, {"Family Calendar & Event Planning", "📅"}}},
	{"Travel & Adventure", "🧳", []emojiLink{{"Trip Planning & Itineraries", "🗺️"}, {"Travel Bucket List & Inspiration", "🌎"}, {"Travel Journal & Photo Gallery", "📸"}}},
}

func navWorkspaces() h.Node {
	return sidebar.Group(
		sidebar.GroupLabel("Workspaces"),
		sidebar.GroupContent(
			sidebar.Menu(
				h.Map(workspaces, func(w workspace) h.Node {
					return collapsible.Collapsible(
						sidebar.MenuItem(
							sidebar.MenuButton(lib.AsChild,
								h.A(h.Href("#"), h.Span(w.emoji), h.Span(w.name)),
							),
							collapsible.Trigger(lib.AsChild,
								sidebar.MenuAction(sidebar.ShowOnHover(true),
									h.Class("left-2 bg-sidebar-accent text-sidebar-accent-foreground data-[state=open]:rotate-90"),
									lucide.ChevronRight(),
								),
							),
							sidebar.MenuAction(sidebar.ShowOnHover(true), lucide.Plus()),
							collapsible.Content(
								sidebar.MenuSub(h.Map(w.pages, func(p emojiLink) h.Node {
									return sidebar.MenuSubItem(
										sidebar.MenuSubButton(lib.AsChild,
											h.A(h.Href("#"), h.Span(p.emoji), h.Span(p.name)),
										),
									)
								})),
							),
						),
					)
				}),
				sidebar.MenuItem(
					sidebar.MenuButton(h.Class("text-sidebar-foreground/70"), lucide.Ellipsis(), h.Span("More")),
				),
			),
		),
	)
}

// leftSidebar is sidebar-10's AppSidebar and sidebar-15's SidebarLeft.
func leftSidebar() h.Node {
	return sidebar.Sidebar(h.Class("border-r-0"),
		sidebar.Header(
			teamSwitcher(),
			navMain(navMainItems),
		),
		sidebar.Content(
			navFavorites(),
			navWorkspaces(),
			navSecondary(navSecondaryItems, sidebar.ButtonSizeDefault, h.Class("mt-auto")),
		),
		sidebar.Rail(),
	)
}

// --- sidebar-12 / sidebar-15 calendars -----------------------------------------------

type calendarGroup struct {
	name  string
	items []string
}

var calendarGroups = []calendarGroup{
	{"My Calendars", []string{"Personal", "Work", "Family"}},
	{"Favorites", []string{"Holidays", "Birthdays"}},
	{"Other", []string{"Travel", "Reminders", "Deadlines"}},
}

func datePicker() h.Node {
	return sidebar.Group(h.Class("px-0"),
		sidebar.GroupContent(
			calendar.Calendar(h.Class("[&_[role=gridcell]]:w-[33px] [&_[role=gridcell].bg-accent]:bg-sidebar-primary [&_[role=gridcell].bg-accent]:text-sidebar-primary-foreground")),
		),
	)
}

// calendars is Calendars; rounded is the check box radius (rounded-sm in
// sidebar-12, rounded-xs in sidebar-15).
func calendars(rounded string) h.Node {
	var out []any
	for i, c := range calendarGroups {
		out = append(out,
			sidebar.Group(h.Class("py-0"),
				collapsible.Collapsible(collapsible.Open(i == 0), h.Class("group/collapsible"),
					sidebar.GroupLabel(lib.AsChild,
						h.Class("group/label w-full text-sm text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"),
						collapsible.Trigger(
							c.name, " ",
							lucide.ChevronRight(h.Class("ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90")),
						),
					),
					collapsible.Content(
						sidebar.GroupContent(
							sidebar.Menu(h.MapIdx(c.items, func(j int, item string) h.Node {
								return sidebar.MenuItem(
									sidebar.MenuButton(
										h.Div(h.Data("active", fmt.Sprint(j < 2)),
											h.Class("group/calendar-item flex aspect-square size-4 shrink-0 items-center justify-center "+rounded+" border border-sidebar-border text-sidebar-primary-foreground data-[active=true]:border-sidebar-primary data-[active=true]:bg-sidebar-primary"),
											lucide.Check(h.Class("hidden size-3 group-data-[active=true]/calendar-item:block")),
										),
										item,
									),
								)
							})),
						),
					),
				),
			),
			sidebar.Separator(h.Class("mx-0")),
		)
	}
	return h.Group(out...)
}

// calendarSidebar is the body of sidebar-12's AppSidebar and sidebar-15's
// SidebarRight.
func calendarSidebar(rounded string) []any {
	return []any{
		sidebar.Header(h.Class("h-16 border-b border-sidebar-border"),
			navUser(shadcnUser, false, dm.AlignStart),
		),
		sidebar.Content(
			datePicker(),
			sidebar.Separator(h.Class("mx-0")),
			calendars(rounded),
		),
		sidebar.Footer(
			sidebar.Menu(
				sidebar.MenuItem(
					sidebar.MenuButton(lucide.Plus(), h.Span("New Calendar")),
				),
			),
		),
	}
}
