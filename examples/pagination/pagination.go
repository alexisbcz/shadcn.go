// Package pagination ports the upstream pagination-* examples.
package pagination

import (
	"shadcn.go/components/ui/pagination"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(registry.Example{Name: "pagination-demo", Component: "pagination", Render: registry.Static(Demo)})
}

func Demo() h.Node {
	return pagination.Pagination(pagination.Content(
		pagination.Item(pagination.Previous(h.Href("#"))),
		pagination.Item(pagination.Link(h.Href("#"), "1")),
		pagination.Item(pagination.Link(h.Href("#"), pagination.IsActive, "2")),
		pagination.Item(pagination.Link(h.Href("#"), "3")),
		pagination.Item(pagination.Ellipsis()),
		pagination.Item(pagination.Next(h.Href("#"))),
	))
}
