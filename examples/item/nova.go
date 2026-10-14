package item

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/item"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the item page (ui/apps/v4/examples/radix/item-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "item", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("item-dropdown", NovaDropdown),
		nova("item-group", NovaGroup),
		nova("item-size", NovaSize),
		nova("item-variant", NovaVariant),
	)
}

var novaPeople = []struct{ username, avatar, email string }{
	{"shadcn", "https://github.com/shadcn.png", "shadcn@vercel.com"},
	{"maxleiter", "https://github.com/maxleiter.png", "maxleiter@vercel.com"},
	{"evilrabbit", "https://github.com/evilrabbit.png", "evilrabbit@vercel.com"},
}

func NovaDropdown() h.Node {
	var items []any
	for _, p := range novaPeople {
		items = append(items, dm.Item(
			item.Item(item.SizeXs, h.Class("w-full p-2"),
				item.Media(
					avatar.Avatar(h.Class("size-[--spacing(6.5)]"),
						avatar.Image(h.Src(p.avatar), h.Class("grayscale")),
						avatar.Fallback(p.username[:1]),
					),
				),
				item.Content(h.Class("gap-0"), item.Title(p.username), item.Description(h.Class("leading-none"), p.email)),
			),
		))
	}
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Select ", lucide.ChevronDown())),
		dm.Content(h.Class("w-48"), dm.AlignEnd, dm.Group(items...)),
	)
}

func NovaGroup() h.Node {
	group := item.Group(h.Class("max-w-sm"))
	for _, p := range novaPeople {
		group.Children(item.Item(item.VariantOutline,
			item.Media(avatar.Avatar(
				avatar.Image(h.Src(p.avatar), h.Class("grayscale")),
				avatar.Fallback(p.username[:1]),
			)),
			item.Content(h.Class("gap-1"), item.Title(p.username), item.Description(p.email)),
			item.Actions(button.Button(button.VariantGhost, button.SizeIcon, h.Class("rounded-full"), lucide.Plus())),
		))
	}
	return group
}

func novaIconItem(variant item.Variant, size item.Size, title, description string) h.Node {
	return item.Item(variant, size,
		item.Media(item.MediaIcon, lucide.Inbox()),
		item.Content(item.Title(title), item.Description(description)),
	)
}

func NovaSize() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"),
		novaIconItem(item.VariantOutline, item.SizeDefault, "Default Size", "The standard size for most use cases."),
		novaIconItem(item.VariantOutline, item.SizeSm, "Small Size", "A compact size for dense layouts."),
		novaIconItem(item.VariantOutline, item.SizeXs, "Extra Small Size", "The most compact size available."),
	)
}

func NovaVariant() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"),
		novaIconItem(item.VariantDefault, item.SizeDefault, "Default Variant", "Transparent background with no border."),
		novaIconItem(item.VariantOutline, item.SizeDefault, "Outline Variant", "Outlined style with a visible border."),
		novaIconItem(item.VariantMuted, item.SizeDefault, "Muted Variant", "Muted background for secondary content."),
	)
}
