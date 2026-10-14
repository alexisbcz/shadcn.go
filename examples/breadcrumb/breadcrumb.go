// Package breadcrumb ports the upstream breadcrumb-* examples.
package breadcrumb

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/drawer"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "breadcrumb", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("breadcrumb-ellipsis", Ellipsis),
		ex("breadcrumb-link", Link),
		ex("breadcrumb-separator", Separator),
		ex("breadcrumb-demo", Demo),
		ex("breadcrumb-dropdown", Dropdown),
		ex("breadcrumb-responsive", Responsive),
	)
}

// link is BreadcrumbLink asChild around next/link, which renders an <a>.
func link(href, label string) h.Node {
	return breadcrumb.Item(breadcrumb.Link(lib.AsChild, h.A(h.Href(href), label)))
}

func Ellipsis() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("/", "Home"),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Ellipsis()),
		breadcrumb.Separator(),
		link("/docs/components", "Components"),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

func Link() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("/", "Home"),
		breadcrumb.Separator(),
		link("/components", "Components"),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

func Separator() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("/", "Home"),
		breadcrumb.Separator(lucide.Slash()),
		link("/components", "Components"),
		breadcrumb.Separator(lucide.Slash()),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

func menuItems() []any {
	return []any{dm.Item("Documentation"), dm.Item("Themes"), dm.Item("GitHub")}
}

func Demo() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("/", "Home"),
		breadcrumb.Separator(),
		breadcrumb.Item(dm.DropdownMenu(
			dm.Trigger(h.Class("flex items-center gap-1"),
				breadcrumb.Ellipsis(h.Class("size-4")),
				h.Span(h.Class("sr-only"), "Toggle menu"),
			),
			dm.Content(append([]any{dm.AlignStart}, menuItems()...)...),
		)),
		breadcrumb.Separator(),
		link("/docs/components", "Components"),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

func Dropdown() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("/", "Home"),
		breadcrumb.Separator(lucide.Slash()),
		breadcrumb.Item(dm.DropdownMenu(
			dm.Trigger(h.Class("flex items-center gap-1 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-3.5"),
				"Components", lucide.ChevronDown(),
			),
			dm.Content(append([]any{dm.AlignStart}, menuItems()...)...),
		)),
		breadcrumb.Separator(lucide.Slash()),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

// Responsive is breadcrumb-responsive: a DropdownMenu on screens at least
// 768px wide (the reference viewport), a Drawer below (useMediaQuery). The
// Drawer version waits in a template that examples/drawer's media.js swaps
// in when the media query stops matching.
func Responsive() h.Node {
	middle := []string{"Documentation", "Build Your Application"}
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("#", "Home"),
		breadcrumb.Separator(),
		breadcrumb.Item(
			dm.DropdownMenu(
				dm.Trigger(h.Class("flex items-center gap-1"), h.Aria("label", "Toggle menu"),
					breadcrumb.Ellipsis(h.Class("size-4")),
				),
				dm.Content(dm.AlignStart, h.Map(middle, func(l string) h.Node {
					return dm.Item(h.A(h.Href("#"), l))
				})),
			),
			h.Template(h.Data("parity-ignore", ""), h.Data("example-media-swap", "(min-width: 768px)"),
				drawer.Drawer(
					drawer.Trigger(h.Aria("label", "Toggle Menu"), breadcrumb.Ellipsis(h.Class("h-4 w-4"))),
					drawer.Content(
						drawer.Header(h.Class("text-left"),
							drawer.Title("Navigate to"),
							drawer.Description("Select a page to navigate to."),
						),
						h.Div(h.Class("grid gap-1 px-4"), h.Map(middle, func(l string) h.Node {
							return h.A(h.Href("#"), h.Class("py-1 text-sm"), l)
						})),
						drawer.Footer(h.Class("pt-4"),
							drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Close")),
						),
					),
				),
			),
		),
		breadcrumb.Separator(),
		breadcrumb.Item(
			breadcrumb.Link(lib.AsChild, h.Class("max-w-20 truncate md:max-w-none"), h.A(h.Href("#"), "Data Fetching")),
			breadcrumb.Separator(),
		),
		breadcrumb.Item(breadcrumb.Page(h.Class("max-w-20 truncate md:max-w-none"), "Caching and Revalidating")),
	))
}
