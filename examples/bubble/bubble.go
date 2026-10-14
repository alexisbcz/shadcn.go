// Package bubble ports the bubble parity fixtures (upstream ships Bubble
// without new-york-v4 examples).
package bubble

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "bubble", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("bubble-variants", Variants),
		ex("bubble-interactive", Interactive),
	)
}

// Variants is the parity fixture bubble-variants.
func Variants() h.Node {
	return h.Div(h.Class("grid w-full max-w-3xl grid-cols-2 gap-x-8 gap-y-8 py-8"),
		h.Div(h.Class("flex flex-col gap-8"),
			h.Map(bubble.VariantValues, func(v bubble.Variant) h.Node {
				return bubble.Bubble(v, bubble.Content("The ", h.Comment(""), string(v), h.Comment(""), " bubble, at the start."))
			}),
			bubble.Bubble(bubble.Content("Default variant and alignment.")),
		),
		h.Div(h.Class("flex flex-col gap-8"),
			h.Map(bubble.VariantValues, func(v bubble.Variant) h.Node {
				return bubble.Bubble(v, bubble.AlignEnd, bubble.Content("The ", h.Comment(""), string(v), h.Comment(""), " bubble, at the end, with a longer text that wraps onto a second line."))
			}),
		),
		bubble.Group(h.Class("col-span-2 gap-4"),
			bubble.Bubble(bubble.VariantMuted,
				bubble.Content("Reactions default to the bottom end."),
				bubble.Reactions(h.Role("img"), h.Aria("label", "Reaction: thumbs up"), h.Span("👍")),
			),
			bubble.Bubble(bubble.VariantSecondary, bubble.AlignEnd,
				bubble.Content("Bottom start reactions."),
				bubble.Reactions(bubble.AlignStart, h.Span("👀"), h.Span("🚀"), h.Span("+2")),
			),
			bubble.Bubble(bubble.VariantTinted, h.Class("mt-4"),
				bubble.Content("Top start reactions."),
				bubble.Reactions(bubble.SideTop, bubble.AlignStart, h.Span("🎉")),
			),
			bubble.Bubble(bubble.VariantOutline, bubble.AlignEnd, h.Class("mt-4"),
				bubble.Content(h.Class("font-medium"), "Top end reactions."),
				bubble.Reactions(bubble.SideTop, bubble.AlignEnd, h.Class("bg-accent"), h.Span("🔥"), h.Span("👏")),
			),
		),
	)
}

// Interactive is the parity fixture bubble-interactive.
func Interactive() h.Node {
	return h.Div(h.Class("flex w-full max-w-xl flex-col gap-4 py-8"),
		h.Map(bubble.VariantValues, func(v bubble.Variant) h.Node {
			return h.Div(h.Class("flex justify-between gap-4"),
				bubble.Bubble(v, bubble.Content(lib.AsChild, h.Button(string(v), h.Comment(""), " button"))),
				bubble.Bubble(v, bubble.AlignEnd, bubble.Content(lib.AsChild, h.A(h.Href("#bubble"), string(v), h.Comment(""), " link"))),
			)
		}),
		bubble.Bubble(bubble.VariantDestructive, h.Class("mb-6"),
			bubble.Content("Are you sure I can run this command?"),
			bubble.Reactions(button.Button(button.VariantGhost, button.SizeXs, "Yes, run it")),
		),
	)
}
