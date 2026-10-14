package spinner

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the spinner page (ui/apps/v4/examples/radix/spinner-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "spinner", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("spinner-badge", NovaBadge),
		nova("spinner-button", NovaButton),
		nova("spinner-rtl", NovaRtl),
	)
}

func start() h.Node { return spinner.Spinner(h.Data("icon", "inline-start")) }

func NovaBadge() h.Node {
	return h.Div(h.Class("flex items-center gap-4 [--radius:1.2rem]"),
		badge.Badge(start(), "Syncing"),
		badge.Badge(badge.VariantSecondary, start(), "Updating"),
		badge.Badge(badge.VariantOutline, start(), "Processing"),
	)
}

func NovaButton() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-4"),
		button.Button(h.Bool("disabled"), button.SizeSm, start(), "Loading..."),
		button.Button(button.VariantOutline, h.Bool("disabled"), button.SizeSm, start(), "Please wait"),
		button.Button(button.VariantSecondary, h.Bool("disabled"), button.SizeSm, start(), "Processing"),
	)
}

func NovaRtl() h.Node {
	return style.ScopeRTL(h.Div(h.Class("flex w-full max-w-xs flex-col gap-4 [--radius:1rem]"), h.Dir("rtl"),
		item.Item(item.VariantMuted, h.Dir("rtl"),
			item.Media(spinner.Spinner()),
			item.Content(item.Title(h.Class("line-clamp-1"), "جاري معالجة الدفع...")),
			item.Content(h.Class("flex-none justify-end"), h.Span(h.Class("text-sm tabular-nums"), "١٠٠.٠٠ دولار")),
		),
	))
}
