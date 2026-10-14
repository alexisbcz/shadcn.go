package pagination

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/pagination"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/pagination-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "pagination", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("pagination-icons-only", NovaIconsOnly),
		nova("pagination-simple", NovaSimple),
		nova("pagination-rtl", NovaRtl),
	)
}

func NovaIconsOnly() h.Node {
	var items []any
	for _, v := range []string{"10", "25", "50", "100"} {
		items = append(items, selectui.Item(selectui.ItemValue(v), v))
	}
	return h.Div(h.Class("flex items-center justify-between gap-4"),
		field.Field(field.Horizontal, h.Class("w-fit"),
			field.Label(h.For("select-rows-per-page"), "Rows per page"),
			selectui.Select(selectui.DefaultValue("25"),
				selectui.Trigger(h.Class("w-20"), h.ID("select-rows-per-page"), selectui.Value()),
				selectui.Content(selectui.AlignStart, selectui.Group(items...)),
			),
		),
		pagination.Pagination(h.Class("mx-0 w-auto"),
			pagination.Content(
				pagination.Item(pagination.Previous(h.Href("#"))),
				pagination.Item(pagination.Next(h.Href("#"))),
			),
		),
	)
}

func NovaSimple() h.Node {
	var items []any
	for _, n := range []string{"1", "2", "3", "4", "5"} {
		link := pagination.Link(h.Href("#"), n)
		if n == "2" {
			link.Active(true)
		}
		items = append(items, pagination.Item(link))
	}
	return pagination.Pagination(pagination.Content(items...))
}

func novaRtl() h.Node {
	return pagination.Pagination(h.Dir("rtl"),
		pagination.Content(
			pagination.Item(pagination.Previous(h.Href("#"), pagination.Text("السابق"))),
			pagination.Item(pagination.Link(h.Href("#"), "١")),
			pagination.Item(pagination.Link(h.Href("#"), pagination.IsActive, "٢")),
			pagination.Item(pagination.Link(h.Href("#"), "٣")),
			pagination.Item(pagination.Ellipsis()),
			pagination.Item(pagination.Next(h.Href("#"), pagination.Text("التالي"))),
		),
	)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
