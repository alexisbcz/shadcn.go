// Package badge ports the upstream badge-* examples.
package badge

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/badge"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "badge", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("badge-demo", Demo),
		ex("badge-destructive", Destructive),
		ex("badge-link", Link),
		ex("badge-outline", Outline),
		ex("badge-secondary", Secondary),
	)
}

func Demo() h.Node {
	const counter = "h-5 min-w-5 rounded-full px-1 font-mono tabular-nums"
	return h.Div(h.Class("flex flex-col items-center gap-2"),
		h.Div(h.Class("flex w-full flex-wrap gap-2"),
			badge.Badge("Badge"),
			badge.Badge(badge.VariantSecondary, "Secondary"),
			badge.Badge(badge.VariantDestructive, "Destructive"),
			badge.Badge(badge.VariantOutline, "Outline"),
		),
		h.Div(h.Class("flex w-full flex-wrap gap-2"),
			badge.Badge(badge.VariantSecondary, h.Class("bg-blue-500 text-white dark:bg-blue-600"),
				lucide.BadgeCheck(),
				"Verified",
			),
			badge.Badge(h.Class(counter), "8"),
			badge.Badge(h.Class(counter), badge.VariantDestructive, "99"),
			badge.Badge(h.Class(counter), badge.VariantOutline, "20+"),
		),
	)
}

func Destructive() h.Node { return badge.Badge(badge.VariantDestructive, "Destructive") }
func Outline() h.Node     { return badge.Badge(badge.VariantOutline, "Outline") }
func Secondary() h.Node   { return badge.Badge(badge.VariantSecondary, "Secondary") }

// Link ports the badge-link fixture (parity/fixtures/badge-link.tsx):
// badges rendered as links with AsChild, plus aria-invalid.
func Link() h.Node {
	link := func(v badge.Variant, children ...any) h.Node {
		return badge.Badge(v, lib.AsChild, h.A(append([]any{h.Href("#")}, children...)...))
	}
	return h.Div(h.Class("flex flex-wrap items-center gap-2"),
		link(badge.VariantDefault, "Default"),
		link(badge.VariantSecondary, "Secondary"),
		link(badge.VariantDestructive, "Destructive"),
		link(badge.VariantOutline, "Outline ", lucide.ArrowUpRight()),
		link(badge.VariantGhost, "Ghost"),
		link(badge.VariantLink, "Link"),
		badge.Badge(h.Aria("invalid", "true"), "Invalid"),
	)
}
