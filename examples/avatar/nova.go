package avatar

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/avatar-*.tsx) whose content
// differs from the new-york-v4 example of the same name, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "avatar", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("avatar-badge", NovaBadge),
		nova("avatar-badge-icon", NovaBadgeIcon),
		nova("avatar-basic", NovaBasic),
		nova("avatar-demo", NovaDemo),
		nova("avatar-dropdown", NovaDropdown),
		nova("avatar-group", NovaGroup),
		nova("avatar-group-count", NovaGroupCount),
		nova("avatar-group-count-icon", NovaGroupCountIcon),
		nova("avatar-size", NovaSize),
		nova("avatar-rtl", NovaRtl),
	)
}

func gh(user, fallback string, attrs ...any) *avatar.Builder {
	return avatar.Avatar(append(attrs,
		avatar.Image(h.Src("https://github.com/"+user+".png"), h.Alt("@"+user)),
		avatar.Fallback(fallback),
	)...)
}

func NovaBadge() h.Node {
	return avatar.Avatar(
		avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn")),
		avatar.Fallback("CN"),
		avatar.Badge(h.Class("bg-green-600 dark:bg-green-800")),
	)
}

func NovaBadgeIcon() h.Node {
	return avatar.Avatar(h.Class("grayscale"),
		avatar.Image(h.Src("https://github.com/pranathip.png"), h.Alt("@pranathip")),
		avatar.Fallback("PP"),
		avatar.Badge(lucide.Plus()),
	)
}

func NovaBasic() h.Node {
	return avatar.Avatar(
		avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn"), h.Class("grayscale")),
		avatar.Fallback("CN"),
	)
}

func three() []any {
	return []any{gh("shadcn", "CN"), gh("maxleiter", "LR"), gh("evilrabbit", "ER")}
}

func novaDemo(more string, attrs ...any) h.Node {
	return h.Div(append([]any{h.Class("flex flex-row flex-wrap items-center gap-6 md:gap-12")}, append(attrs,
		avatar.Avatar(
			avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn"), h.Class("grayscale")),
			avatar.Fallback("CN"),
		),
		avatar.Avatar(
			avatar.Image(h.Src("https://github.com/evilrabbit.png"), h.Alt("@evilrabbit")),
			avatar.Fallback("ER"),
			avatar.Badge(h.Class("bg-green-600 dark:bg-green-800")),
		),
		avatar.Group(append([]any{h.Class("grayscale")}, append(three(), avatar.GroupCount(more))...)...),
	)...)...)
}

func NovaDemo() h.Node { return novaDemo("+3") }

func novaRtl() h.Node { return novaDemo("+٣", h.Dir("rtl")) }

func NovaDropdown() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, button.SizeIcon, h.Class("rounded-full"),
				avatar.Avatar(
					avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("shadcn")),
					avatar.Fallback("CN"),
				),
			),
		),
		dm.Content(h.Class("w-32"),
			dm.Group(dm.Item("Profile"), dm.Item("Billing"), dm.Item("Settings")),
			dm.Separator(),
			dm.Group(dm.Item(dm.VariantDestructive, "Log out")),
		),
	)
}

func NovaGroup() h.Node {
	return avatar.Group(append([]any{h.Class("grayscale")}, three()...)...)
}

func NovaGroupCount() h.Node {
	return avatar.Group(append([]any{h.Class("grayscale")}, append(three(), avatar.GroupCount("+3"))...)...)
}

func NovaGroupCountIcon() h.Node {
	return avatar.Group(append([]any{h.Class("grayscale")}, append(three(), avatar.GroupCount(lucide.Plus()))...)...)
}

func NovaSize() h.Node {
	return h.Div(h.Class("flex flex-wrap items-center gap-2 grayscale"),
		gh("shadcn", "CN", avatar.SizeSm),
		gh("shadcn", "CN"),
		gh("shadcn", "CN", avatar.SizeLg),
	)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
