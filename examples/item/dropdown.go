package item

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/item"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(registry.Example{Name: "item-dropdown", Component: "item", Render: registry.Static(Dropdown)})
}

func Dropdown() h.Node {
	people := []struct{ username, avatar, email string }{
		{"shadcn", "https://github.com/shadcn.png", "shadcn@vercel.com"},
		{"maxleiter", "https://github.com/maxleiter.png", "maxleiter@vercel.com"},
		{"evilrabbit", "https://github.com/evilrabbit.png", "evilrabbit@vercel.com"},
	}
	var items []any
	items = append(items, h.Class("w-72 [--radius:0.65rem]"), dm.AlignEnd)
	for _, p := range people {
		items = append(items, dm.Item(h.Class("p-0"),
			item.Item(item.SizeSm, h.Class("w-full p-2"),
				item.Media(
					avatar.Avatar(h.Class("size-8"),
						avatar.Image(h.Src(p.avatar), h.Class("grayscale")),
						avatar.Fallback(p.username[:1]),
					),
				),
				item.Content(h.Class("gap-0.5"), item.Title(p.username), item.Description(p.email)),
			),
		))
	}
	return h.Div(h.Class("flex min-h-64 w-full max-w-md flex-col items-center gap-6"),
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, h.Class("w-fit"), "Select ", lucide.ChevronDown())),
			dm.Content(items...),
		),
	)
}
