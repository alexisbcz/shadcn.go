package card

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the base and aria card pages that differ from the radix
// ones (ui/apps/v4/examples/{base,aria}/card-*.tsx).
func init() {
	small := func() h.Node { return style.Scope(style.Nova, BaseSmall()) }
	registry.Register(
		registry.Example{Name: "base--card-small", Component: "card", Render: registry.Static(small)},
		registry.Example{Name: "aria--card-small", Component: "card", Render: registry.Static(small)},
	)
}

// BaseSmall is examples/base/card-small.tsx.
func BaseSmall() h.Node {
	item := func(text string) h.Node {
		return h.Li(h.Class("flex gap-2"),
			lucide.ChevronRight(h.Class("mt-0.5 size-4 shrink-0 text-muted-foreground")),
			h.Span(text),
		)
	}
	return card.Card(card.SizeSm, h.Class("mx-auto w-full max-w-xs"),
		card.Header(
			card.Title("Scheduled reports"),
			card.Description("Weekly snapshots. No more manual exports."),
		),
		card.Content(h.Ul(h.Class("grid gap-2 py-2 text-sm"),
			item("Choose a schedule (daily, or weekly)."),
			item("Send to channels or specific teammates."),
			item("Include charts, tables, and key metrics."),
		)),
		card.Footer(h.Class("flex-col gap-2"),
			button.Button(button.SizeSm, h.Class("w-full"), "Set up scheduled reports"),
			button.Button(button.VariantOutline, button.SizeSm, h.Class("w-full"), "See what's new"),
		),
	)
}
