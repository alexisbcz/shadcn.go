// Package sidebar ports the upstream sidebar-* internal demos
// (registry/new-york-v4/internal).
package sidebar

import (
	_ "embed"
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/collapsible"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/components/ui/sonner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed controlled.js
var controlledScript string

//go:embed group_action.js
var groupActionScript string

func init() {
	runtime.Module("example-sidebar-controlled", controlledScript)
	runtime.Module("example-sidebar-group-action", groupActionScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "sidebar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("sidebar-demo", Demo),
		ex("sidebar-controlled", Controlled),
		ex("sidebar-group", GroupExample),
		ex("sidebar-group-collapsible", GroupCollapsible),
		ex("sidebar-menu", MenuExample),
		ex("sidebar-menu-badge", MenuBadge),
		ex("sidebar-menu-sub", MenuSub),
		ex("sidebar-menu-collapsible", MenuCollapsible),
		ex("sidebar-rsc", RSC),
		ex("view--sidebar-rsc", RSCResolved),
		ex("sidebar-header", HeaderExample),
		ex("sidebar-footer", FooterExample),
		ex("sidebar-menu-action", MenuAction),
		registry.Example{Name: "sidebar-group-action", Component: "sidebar", Render: registry.Static(GroupAction), Extra: toaster},
	)
}

type link struct {
	title string
	icon  func(...any) *h.Element
	badge string
}

var projects = []link{
	{"Design Engineering", lucide.Frame, "24"},
	{"Sales & Marketing", lucide.ChartPie, "12"},
	{"Travel", lucide.Map, "3"},
	{"Support", lucide.LifeBuoy, "21"},
	{"Feedback", lucide.Send, "8"},
}

func linkButton(l link, extra ...any) h.Node {
	return sidebar.MenuButton(append([]any{lib.AsChild}, append(extra,
		h.A(h.Href("#"), l.icon(), h.Span(l.title)),
	)...)...)
}

func projectMenu() h.Node {
	return sidebar.Menu(h.Map(projects, func(p link) h.Node {
		return sidebar.MenuItem(linkButton(p))
	}))
}

func header(children ...any) h.Node {
	return h.Header(append([]any{h.Class("flex h-12 items-center justify-between px-4")}, children...)...)
}

// Demo is sidebar-demo.
func Demo() h.Node {
	items := []link{
		{"Home", lucide.House, ""},
		{"Inbox", lucide.Inbox, ""},
		{"Calendar", lucide.Calendar, ""},
		{"Search", lucide.Search, ""},
		{"Settings", lucide.Settings, ""},
	}
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Application"),
					sidebar.GroupContent(
						sidebar.Menu(h.Map(items, func(l link) h.Node {
							return sidebar.MenuItem(linkButton(l))
						})),
					),
				),
			),
		),
		sidebar.Inset(header(sidebar.Trigger())),
	)
}

// Controlled is sidebar-controlled: a button outside the sidebar drives the
// provider's state (controlled.js).
func Controlled() h.Node {
	return sidebar.Provider(
		h.Data("example-sidebar-controlled", ""),
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupContent(projectMenu()),
				),
			),
		),
		sidebar.Inset(header(
			button.Button(button.SizeSm, button.VariantGhost,
				lucide.PanelLeftClose(),
				h.Span("Close", " Sidebar"),
			),
			h.Template(h.Data("parity-ignore", ""), h.Data("icon", "open"), lucide.PanelLeftOpen()),
			h.Template(h.Data("parity-ignore", ""), h.Data("icon", "close"), lucide.PanelLeftClose()),
		)),
	)
}

func helpMenu() h.Node {
	return sidebar.Menu(
		sidebar.MenuItem(sidebar.MenuButton(lucide.LifeBuoy(), "Support")),
		sidebar.MenuItem(sidebar.MenuButton(lucide.Send(), "Feedback")),
	)
}

// GroupExample is sidebar-group.
func GroupExample() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Help"),
					sidebar.GroupContent(helpMenu()),
				),
			),
		),
	)
}

// GroupCollapsible is sidebar-group-collapsible.
func GroupCollapsible() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				collapsible.Collapsible(collapsible.Open(true), h.Class("group/collapsible"),
					sidebar.Group(
						sidebar.GroupLabel(lib.AsChild, h.Class("text-sm hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"),
							collapsible.Trigger(
								"Help",
								lucide.ChevronDown(h.Class("ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180")),
							),
						),
						collapsible.Content(
							sidebar.GroupContent(helpMenu()),
						),
					),
				),
			),
		),
	)
}

// MenuExample is sidebar-menu.
func MenuExample() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupContent(projectMenu()),
				),
			),
		),
	)
}

// MenuBadge is sidebar-menu-badge.
func MenuBadge() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupContent(
						sidebar.Menu(h.Map(projects, func(p link) h.Node {
							return sidebar.MenuItem(
								linkButton(p, h.Class("group-has-[[data-state=open]]/menu-item:bg-sidebar-accent")),
								sidebar.MenuBadge(p.badge),
							)
						})),
					),
				),
			),
		),
	)
}

type section struct {
	title string
	items []string
}

var docs = []section{
	{"Getting Started", []string{"Installation", "Project Structure"}},
	{"Build Your Application", []string{"Routing", "Data Fetching", "Rendering", "Caching", "Styling", "Optimizing", "Configuring", "Testing", "Authentication", "Deploying", "Upgrading", "Examples"}},
	{"API Reference", []string{"Components", "File Conventions", "Functions", "next.config.js Options", "CLI", "Edge Runtime"}},
	{"Architecture", []string{"Accessibility", "Fast Refresh", "Next.js Compiler", "Supported Browsers", "Turbopack"}},
}

func subMenu(items []string) h.Node {
	return sidebar.MenuSub(h.Map(items, func(t string) h.Node {
		return sidebar.MenuSubItem(sidebar.MenuSubButton(lib.AsChild, h.A(h.Href("#"), h.Span(t))))
	}))
}

// MenuSub is sidebar-menu-sub.
func MenuSub() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupContent(
						sidebar.Menu(h.Map(docs, func(s section) h.Node {
							return sidebar.MenuItem(
								sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), h.Span(s.title))),
								subMenu(s.items),
							)
						})),
					),
				),
			),
		),
	)
}

// MenuCollapsible is sidebar-menu-collapsible.
func MenuCollapsible() h.Node {
	var items []any
	for i, s := range docs {
		items = append(items, collapsible.Collapsible(h.Class("group/collapsible"), collapsible.Open(i == 0),
			sidebar.MenuItem(
				collapsible.Trigger(lib.AsChild,
					sidebar.MenuButton(
						h.Span(s.title),
						lucide.ChevronRight(h.Class("ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90")),
					),
				),
				collapsible.Content(subMenu(s.items)),
			),
		))
	}
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupContent(sidebar.Menu(items...)),
				),
			),
		),
	)
}

// RSC is sidebar-rsc: the Suspense fallback (NavProjectsSkeleton). The
// async NavProjects never resolves in a client-rendered React tree, so
// the fallback is what upstream shows.
func RSC() h.Node {
	var items []any
	for range 5 {
		items = append(items, sidebar.MenuItem(sidebar.MenuSkeleton(sidebar.ShowIcon(true))))
	}
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupContent(sidebar.Menu(items...)),
				),
			),
		),
	)
}

// RSCResolved is sidebar-rsc as the docs site serves it: there the server
// component resolves (at build time), so NavProjects replaces the skeleton.
// The docs /view route uses it ("view--" names override the example there).
func RSCResolved() h.Node {
	projects := []struct {
		name string
		icon func(...any) *h.Element
	}{
		{"Design Engineering", lucide.Frame},
		{"Sales & Marketing", lucide.ChartPie},
		{"Travel", lucide.Map},
		{"Support", lucide.LifeBuoy},
		{"Feedback", lucide.Send},
	}
	var items []any
	for _, p := range projects {
		items = append(items, sidebar.MenuItem(
			sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), p.icon(), h.Span(p.name))),
		))
	}
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupContent(sidebar.Menu(items...)),
				),
			),
		),
	)
}

// toaster is the shadcn <Toaster /> the upstream docs mount in their
// layout (see examples/sonner).
func toaster(*http.Request) h.Node { return sonner.Toaster(sonner.ThemeLight) }

// plainToaster is sonner's own <Toaster position="bottom-left"
// toastOptions={{ className: "ml-[160px]" }} /> (not shadcn's wrapper: no
// className, style or icons).
func plainToaster() h.Node {
	el := sonner.Toaster(sonner.ThemeLight, sonner.SonnerIcons, sonner.PositionBottomLeft,
		sonner.ToastOptions(map[string]any{"className": "ml-[160px]"}),
	).Element()
	for _, c := range el.ChildNodes() {
		if tpl, ok := c.(*h.Element); ok && tpl.Tag() == "template" {
			tpl.Remove("data-class")
			tpl.Remove("data-style")
		}
	}
	return el
}

// GroupAction is sidebar-group-action: the action shows a toast
// (group_action.js).
func GroupAction() h.Node {
	return sidebar.Provider(
		plainToaster(),
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupAction(h.Attr("title", "Add Project"), h.Data("example-sidebar-group-action", ""),
						lucide.Plus(), " ", h.Span(h.Class("sr-only"), "Add Project"),
					),
					sidebar.GroupContent(
						sidebar.Menu(h.Map(projects[:3], func(p link) h.Node {
							return sidebar.MenuItem(linkButton(p))
						})),
					),
				),
			),
		),
	)
}

func menuDropdown(title string, chevron h.Node, items []string, content ...any) h.Node {
	return sidebar.Menu(
		sidebar.MenuItem(
			dropdownmenu.DropdownMenu(
				dropdownmenu.Trigger(lib.AsChild,
					sidebar.MenuButton(h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
						title, chevron,
					),
				),
				dropdownmenu.Content(append(content, h.Map(items, func(t string) h.Node {
					return dropdownmenu.Item(h.Span(t))
				}))...),
			),
		),
	)
}

// HeaderExample is sidebar-header.
func HeaderExample() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Header(
				menuDropdown("Select Workspace", lucide.ChevronDown(h.Class("ml-auto")), []string{"Acme Inc", "Acme Corp."},
					h.Class("w-(--radix-popper-anchor-width)")),
			),
		),
		sidebar.Inset(header(sidebar.Trigger())),
	)
}

// FooterExample is sidebar-footer.
func FooterExample() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Header(),
			sidebar.Content(),
			sidebar.Footer(
				menuDropdown("Username", lucide.ChevronUp(h.Class("ml-auto")), []string{"Account", "Billing", "Sign out"},
					dropdownmenu.SideTop, h.Class("w-(--radix-popper-anchor-width)")),
			),
		),
		sidebar.Inset(header(sidebar.Trigger())),
	)
}

// MenuAction is sidebar-menu-action.
func MenuAction() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Projects"),
					sidebar.GroupContent(
						sidebar.Menu(h.Map(projects, func(p link) h.Node {
							return sidebar.MenuItem(
								linkButton(p, h.Class("group-has-[[data-state=open]]/menu-item:bg-sidebar-accent")),
								dropdownmenu.DropdownMenu(
									dropdownmenu.Trigger(lib.AsChild,
										sidebar.MenuAction(lucide.Ellipsis(), h.Span(h.Class("sr-only"), "More")),
									),
									dropdownmenu.Content(dropdownmenu.SideRight, dropdownmenu.AlignStart,
										dropdownmenu.Item(h.Span("Edit Project")),
										dropdownmenu.Item(h.Span("Delete Project")),
									),
								),
							)
						})),
					),
				),
			),
		),
	)
}
