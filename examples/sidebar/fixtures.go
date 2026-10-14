package sidebar

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "sidebar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("sidebar-icon", Icon),
		ex("sidebar-right", Right),
		ex("sidebar-none", None),
	)
}

// Icon is the sidebar-icon fixture: collapsible icon, floating variant,
// rail, tooltips and the less used parts.
func Icon() h.Node {
	return sidebar.Provider(
		sidebar.Sidebar(sidebar.CollapsibleIcon, sidebar.VariantFloating,
			sidebar.Header(
				sidebar.Menu(
					sidebar.MenuItem(
						sidebar.MenuButton(sidebar.ButtonSizeLg, sidebar.ButtonVariantOutline, sidebar.Tooltip("Acme Inc"),
							lucide.Frame(), h.Span("Acme Inc"), lucide.ChevronsUpDown(h.Class("ml-auto")),
						),
					),
				),
				sidebar.Input(h.Placeholder("Search...")),
			),
			sidebar.Separator(),
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Platform"),
					sidebar.GroupAction(h.Attr("title", "Add"), lucide.Plus(), h.Span(h.Class("sr-only"), "Add")),
					sidebar.GroupContent(
						sidebar.Menu(
							sidebar.MenuItem(
								sidebar.MenuButton(sidebar.IsActive(true), sidebar.Tooltip("Home"), lucide.House(), h.Span("Home")),
								sidebar.MenuBadge("3"),
							),
							sidebar.MenuItem(
								sidebar.MenuButton(sidebar.Tooltip(h.Class("font-bold"), "Inbox"), lucide.Inbox(), h.Span("Inbox")),
								sidebar.MenuAction(sidebar.ShowOnHover(true), lucide.Ellipsis(), h.Span(h.Class("sr-only"), "More")),
							),
							sidebar.MenuItem(
								sidebar.MenuButton(sidebar.ButtonSizeSm, h.Bool("disabled"), sidebar.Tooltip("Settings"), lucide.Settings(), h.Span("Settings")),
								sidebar.MenuSub(
									sidebar.MenuSubItem(sidebar.MenuSubButton(h.Href("#"), sidebar.IsActive(true), h.Span("General"))),
									sidebar.MenuSubItem(sidebar.MenuSubButton(h.Href("#"), sidebar.SubButtonSizeSm, h.Span("Team"))),
								),
							),
						),
					),
				),
			),
			sidebar.Footer(
				sidebar.Menu(
					sidebar.MenuItem(
						sidebar.MenuButton(sidebar.Tooltip("Account"), lucide.Settings(), h.Span("Account")),
					),
				),
			),
			sidebar.Rail(),
		),
		sidebar.Inset(
			h.Header(h.Class("flex h-12 items-center gap-2 border-b px-4"),
				sidebar.Trigger(h.Class("-ml-1")),
				h.Span(h.Class("text-sm"), "Dashboard"),
			),
		),
	)
}

// Right is the sidebar-right fixture: right side, inset variant, rail,
// starting collapsed.
func Right() h.Node {
	return sidebar.Provider(sidebar.DefaultOpen(false),
		sidebar.Inset(
			h.Header(h.Class("flex h-12 items-center justify-end px-4"), sidebar.Trigger()),
			h.Div(h.Class("p-4 text-sm"), "Main content"),
		),
		sidebar.Sidebar(sidebar.SideRight, sidebar.VariantInset,
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupLabel("Mail"),
					sidebar.GroupContent(
						sidebar.Menu(
							sidebar.MenuItem(sidebar.MenuButton(lucide.House(), h.Span("Home"))),
							sidebar.MenuItem(sidebar.MenuButton(lucide.Inbox(), h.Span("Inbox"))),
						),
					),
				),
			),
			sidebar.Rail(),
		),
	)
}

// None is the sidebar-none fixture: a sidebar that never collapses.
func None() h.Node {
	return sidebar.Provider(h.Class("min-h-0 w-auto"),
		sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("h-64 rounded-lg border"),
			sidebar.Content(
				sidebar.Group(
					sidebar.GroupContent(
						sidebar.Menu(
							sidebar.MenuItem(sidebar.MenuButton(lucide.House(), h.Span("Home"))),
							sidebar.MenuItem(sidebar.MenuButton(lucide.Inbox(), h.Span("Inbox"))),
						),
					),
				),
			),
		),
	)
}
