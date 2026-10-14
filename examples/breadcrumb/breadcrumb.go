// Package breadcrumb ports the upstream breadcrumb-* examples.
//
// Not ported yet (they need components owned by other groups):
// breadcrumb-demo and breadcrumb-dropdown (dropdown-menu),
// breadcrumb-responsive (dropdown-menu, drawer).
package breadcrumb

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/breadcrumb"
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
