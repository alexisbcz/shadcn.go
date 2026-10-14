// Package message ports the message parity fixtures (upstream ships
// Message without new-york-v4 examples).
package message

import (
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/message"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(registry.Example{Name: "message-variants", Component: "message", Render: registry.Static(Variants)})
}

// Variants is the parity fixture message-variants.
func Variants() h.Node {
	return h.Div(h.Class("grid w-full max-w-3xl grid-cols-2 gap-8 py-6"),
		message.Group(
			message.Message(
				message.Avatar(avatar.Avatar(avatar.Fallback("R"))),
				message.Content(
					message.Header("Rabbit"),
					bubble.Bubble(bubble.VariantMuted, bubble.Content("It's 4:55 PM. On a Friday.")),
				),
			),
			message.Message(message.AlignEnd,
				message.Avatar(avatar.Avatar(avatar.Fallback("ME"))),
				message.Content(
					bubble.Bubble(bubble.Content("It's a one-line change.")),
					message.Footer("Delivered"),
				),
			),
			message.Message(message.AlignStart,
				message.Avatar(h.Class("bg-primary text-primary-foreground"), "AI"),
				message.Content(
					bubble.Group(
						bubble.Bubble(bubble.VariantMuted, bubble.Content("It's always a one-line change.")),
						bubble.Bubble(bubble.VariantMuted,
							bubble.Content("Alright, let me take a look."),
							bubble.Reactions(h.Aria("label", "Reactions: thumbs up"), h.Span("👍")),
						),
					),
					message.Footer("Seen"),
				),
			),
		),
		message.Group(h.Class("gap-6"),
			message.Message(
				message.Content(
					message.Header("Assistant"),
					bubble.Bubble(bubble.VariantGhost, bubble.Content("Ghost bubbles take the full width, and the header and footer lose their padding to line up with the text.")),
					message.Footer(h.Div("Read ", h.Span(h.Class("font-normal"), "Yesterday"))),
				),
			),
			message.Message(message.AlignEnd, h.Class("text-base"),
				message.Content(h.Class("gap-1"),
					message.Header("You"),
					bubble.Bubble(bubble.VariantSecondary, bubble.Content("Send the report to the team.")),
					message.Footer(h.Class("text-foreground"), "Sent"),
				),
			),
			marker.Marker(h.Role("status"),
				marker.Content(h.Class("shimmer"), h.Span(h.Class("font-medium"), "Oliver"), " is typing..."),
			),
		),
	)
}
