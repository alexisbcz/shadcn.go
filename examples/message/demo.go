package message

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/attachment"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/message"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "message", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("message-demo", Demo),
		ex("message-parts", Parts),
	)
}

// pixel is a 4x4 orange PNG: an image that loads without the network.
const pixel = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAAEElEQVR4nGP4WSwGRwzEcQCBwhghZfPo6AAAAABJRU5ErkJggg=="

// Demo is the parity fixture message-demo.
func Demo() h.Node {
	me := func() h.Node {
		return message.Avatar(avatar.Avatar(avatar.Image(h.Src(pixel), h.Alt("@me")), avatar.Fallback("ME")))
	}
	rabbit := func() h.Node { return message.Avatar(avatar.Avatar(avatar.Fallback("R"))) }
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-6 py-12"),
		message.Message(message.AlignEnd, me(),
			message.Content(bubble.Bubble(bubble.Content("Deploying to prod real quick."))),
		),
		message.Message(rabbit(),
			message.Content(bubble.Bubble(bubble.VariantMuted, bubble.Content("It's 4:55 PM. On a Friday."))),
		),
		message.Message(message.AlignEnd, me(),
			message.Content(
				bubble.Bubble(bubble.Content("It's a one-line change.")),
				message.Footer("Delivered"),
			),
		),
		message.Message(rabbit(),
			message.Content(
				bubble.Group(
					bubble.Bubble(bubble.VariantMuted, bubble.Content("It's always a one-line change 😭.")),
					bubble.Bubble(bubble.VariantMuted,
						bubble.Content("Alright, let me take a look."),
						bubble.Reactions(h.Aria("label", "Reactions: thumbs up"), h.Span("👍")),
					),
				),
			),
		),
		marker.Marker(h.Role("status"),
			marker.Content(h.Class("shimmer"), h.Span(h.Class("font-medium"), "Oliver"), " is typing..."),
		),
	)
}

// Parts is the parity fixture message-parts.
func Parts() h.Node {
	action := func(label string, size button.Size, icon h.Node) h.Node {
		return button.Button(button.VariantGhost, size, h.Aria("label", label), h.Attr("title", label), icon)
	}
	return h.Div(h.Class("grid w-full max-w-3xl grid-cols-2 gap-10 py-8"),
		h.Div(h.Class("flex flex-col gap-8"),
			message.Message(
				message.Content(
					message.Header("Olivia"),
					bubble.Bubble(bubble.VariantMuted, bubble.Content("I already checked the logs.")),
				),
			),
			message.Message(message.AlignEnd,
				message.Content(
					message.Header("You"),
					bubble.Bubble(bubble.Content("Send the report to the team. Ping @shadcn if you need help.")),
					message.Footer(h.Div("Read ", h.Span(h.Class("font-normal"), "Yesterday"))),
				),
			),
			message.Message(
				message.Content(
					message.Header("Assistant"),
					bubble.Bubble(bubble.VariantGhost, bubble.Content("The install failure is coming from the workspace package.")),
					message.Footer(
						action("Copy", button.SizeIcon, lucide.Copy()),
						action("Like", button.SizeIcon, lucide.ThumbsUp()),
						action("Dislike", button.SizeIcon, lucide.ThumbsDown()),
					),
				),
			),
		),
		h.Div(h.Class("flex flex-col gap-8"),
			message.Message(message.AlignEnd,
				message.Content(
					bubble.Bubble(bubble.Content("Okay drop me a link. Taking a look...")),
					message.Footer(h.Class("gap-2"),
						h.Span(h.Class("font-normal text-destructive"), "Failed to send"),
						action("Retry", button.SizeIconXs, lucide.RefreshCcw()),
					),
				),
			),
			message.Message(message.AlignEnd,
				message.Avatar(h.Class("size-8 text-xs font-medium"), "ME"),
				message.Content(
					attachment.Attachment(
						attachment.Media(lucide.FileText()),
						attachment.Content(attachment.Title("quarterly-report.pdf"), attachment.Description("PDF · 2.4 MB")),
					),
					bubble.Bubble(bubble.Content("Here is the report.")),
				),
			),
			message.Group(h.Class("gap-1"),
				message.Message(
					message.Avatar(h.Class("size-8 text-xs font-medium"), "AI"),
					message.Content(bubble.Bubble(bubble.VariantSecondary, bubble.Content("First message in a group."))),
				),
				message.Message(h.Class("pl-10"),
					message.Content(h.Class("gap-1"),
						bubble.Bubble(bubble.VariantSecondary, bubble.Content("A second message without an avatar, and a long text that wraps: averyveryverylongwordwithoutanybreakopportunityatall.")),
					),
				),
			),
		),
	)
}
