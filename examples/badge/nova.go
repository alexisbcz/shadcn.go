package badge

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the badge page (ui/apps/v4/examples/radix/badge-*.tsx),
// in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "badge", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("badge-colors", NovaColors),
		nova("badge-demo", NovaDemo),
		nova("badge-icon", NovaIcon),
		nova("badge-link", NovaLink),
		nova("badge-rtl", NovaRtl),
		nova("badge-spinner", NovaSpinner),
		nova("badge-variants", NovaVariants),
	)
}

func NovaColors() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-2"),
		badge.Badge(h.Class("bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300"), "Blue"),
		badge.Badge(h.Class("bg-green-50 text-green-700 dark:bg-green-950 dark:text-green-300"), "Green"),
		badge.Badge(h.Class("bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300"), "Sky"),
		badge.Badge(h.Class("bg-purple-50 text-purple-700 dark:bg-purple-950 dark:text-purple-300"), "Purple"),
		badge.Badge(h.Class("bg-red-50 text-red-700 dark:bg-red-950 dark:text-red-300"), "Red"),
	)
}

func NovaDemo() h.Node {
	return h.Div(h.Class("flex w-full flex-wrap justify-center gap-2"),
		badge.Badge("Badge"),
		badge.Badge(badge.VariantSecondary, "Secondary"),
		badge.Badge(badge.VariantDestructive, "Destructive"),
		badge.Badge(badge.VariantOutline, "Outline"),
	)
}

func NovaIcon() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-2"),
		badge.Badge(badge.VariantSecondary, lucide.BadgeCheck(h.Data("icon", "inline-start")), "Verified"),
		badge.Badge(badge.VariantOutline, "Bookmark", lucide.Bookmark(h.Data("icon", "inline-end"))),
	)
}

func NovaLink() h.Node {
	return badge.Badge(lib.AsChild, h.A(h.Href("#link"), "Open Link ", lucide.ArrowUpRight(h.Data("icon", "inline-end"))))
}

func NovaRtl() h.Node {
	return style.ScopeRTL(h.Div(h.Class("flex w-full flex-wrap justify-center gap-2"), h.Dir("rtl"),
		badge.Badge("شارة"),
		badge.Badge(badge.VariantSecondary, "ثانوي"),
		badge.Badge(badge.VariantDestructive, "مدمر"),
		badge.Badge(badge.VariantOutline, "مخطط"),
		badge.Badge(badge.VariantSecondary, lucide.BadgeCheck(h.Data("icon", "inline-start")), "متحقق"),
		badge.Badge(badge.VariantOutline, "إشارة مرجعية", lucide.Bookmark(h.Data("icon", "inline-end"))),
	))
}

func NovaSpinner() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-2"),
		badge.Badge(badge.VariantDestructive, spinner.Spinner(h.Data("icon", "inline-start")), "Deleting"),
		badge.Badge(badge.VariantSecondary, "Generating", spinner.Spinner(h.Data("icon", "inline-end"))),
	)
}

func NovaVariants() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-2"),
		badge.Badge("Default"),
		badge.Badge(badge.VariantSecondary, "Secondary"),
		badge.Badge(badge.VariantDestructive, "Destructive"),
		badge.Badge(badge.VariantOutline, "Outline"),
		badge.Badge(badge.VariantGhost, "Ghost"),
	)
}
