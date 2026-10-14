package message

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/attachment"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/message"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the message page (ui/apps/v4/examples/radix/message-*.tsx).
// Upstream's examples import the radix-rhea components, so they render in
// style.Rhea.
func init() {
	docs := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "message", Render: registry.Static(func() h.Node {
			return style.Scope(style.Rhea, fn())
		})}
	}
	registry.Register(
		docs("message-actions", DocsActions),
		docs("message-attachment", DocsAttachment),
		docs("message-avatar", DocsAvatar),
		docs("message-demo", DocsDemo),
		docs("message-group", DocsGroup),
		docs("message-header-footer", DocsHeaderFooter),
		docs("message-markdown", DocsMarkdown),
	)
}

func col(gap string, children ...any) h.Node {
	return h.Div(append([]any{h.Class("flex w-full max-w-sm flex-col " + gap + " py-12")}, children...)...)
}

func docsAvatar(src, alt, fallback string) h.Node {
	return message.Avatar(avatar.Avatar(avatar.Image(h.Src(src), h.Alt(alt)), avatar.Fallback(fallback)))
}

func ghostButton(size button.Size, label string, icon h.Node) h.Node {
	return button.Button(button.VariantGhost, size, h.Aria("label", label), h.Title(label), icon)
}

func DocsActions() h.Node {
	return col("gap-8",
		message.Message(message.Content(
			bubble.Bubble(bubble.VariantMuted, bubble.Content("The install failure is coming from the workspace package.")),
			message.Footer(
				ghostButton(button.SizeIcon, "Copy", lucide.Copy()),
				ghostButton(button.SizeIcon, "Like", lucide.ThumbsUp()),
				ghostButton(button.SizeIcon, "Dislike", lucide.ThumbsDown()),
			),
		)),
		message.Message(message.AlignEnd, message.Content(
			bubble.Bubble(bubble.Content("Okay drop me a link. Taking a look...")),
			message.Footer(h.Class("gap-2"),
				h.Span(h.Class("font-normal text-destructive"), "Failed to send"),
				button.Button(button.VariantGhost, button.SizeIconXs, h.Title("Retry"), h.Aria("label", "Retry"), lucide.RefreshCcw()),
			),
		)),
	)
}

func DocsAttachment() h.Node {
	return col("gap-8",
		message.Message(message.AlignEnd, message.Content(
			attachment.Attachment(attachment.Vertical,
				attachment.Media(attachment.MediaImage,
					h.Img(h.Src("https://images.unsplash.com/photo-1497366754035-f200968a6e72?w=900&auto=format&fit=crop&q=80"), h.Alt("Workspace")),
				),
			),
			bubble.Bubble(bubble.Content("Here's the image. Can you add it to the PDF? Use it for the cover page.")),
		)),
		message.Message(message.Content(
			bubble.Bubble(bubble.VariantMuted, bubble.Content("Done. Here's the PDF with the image added as the cover page.")),
			attachment.Attachment(
				attachment.Media(lucide.FileText()),
				attachment.Content(attachment.Title("sales-dashboard.pdf"), attachment.Description("PDF · 2.4 MB")),
				attachment.Actions(
					attachment.Action(h.Type("button"), h.Title("Download"), h.Aria("label", "Download"), button.SizeIconSm, button.VariantSecondary, lucide.Download()),
				),
			),
		)),
		message.Message(message.AlignEnd, message.Content(
			bubble.Bubble(bubble.Content("Thanks. Looks good.")),
		)),
	)
}

func DocsAvatar() h.Node {
	return col("gap-6",
		message.Message(
			docsAvatar("/avatars/03.png", "@avatar", "R"),
			message.Content(bubble.Bubble(bubble.VariantMuted, bubble.Content("The build failed during dependency installation."))),
		),
		message.Message(message.AlignEnd,
			docsAvatar("/avatars/10.png", "@avatar", "R"),
			message.Content(bubble.Bubble(bubble.Content("Can you share the exact error?"))),
		),
		message.Message(
			docsAvatar("/avatars/03.png", "@avatar", "R"),
			message.Content(bubble.Group(
				bubble.Bubble(bubble.VariantMuted, bubble.Content("Here's the error from the logs")),
				bubble.Bubble(bubble.VariantMuted, bubble.Content("Something went wrong with the build. The libraries are not installed correctly. Try running the build again.")),
			)),
		),
	)
}

func DocsDemo() h.Node {
	me := func() h.Node { return docsAvatar("/avatars/10.png", "@me", "ME") }
	rabbit := func() h.Node { return docsAvatar("/avatars/02.png", "@rabbit", "R") }
	return col("gap-6",
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
			message.Content(bubble.Group(
				bubble.Bubble(bubble.VariantMuted, bubble.Content("It's always a one-line change 😭.")),
				bubble.Bubble(bubble.VariantMuted,
					bubble.Content("Alright, let me take a look."),
					bubble.Reactions(h.Aria("label", "Reactions: thumbs up"), h.Span("👍")),
				),
			)),
		),
		marker.Marker(h.Role("status"),
			marker.Content(h.Class("shimmer"), h.Span(h.Class("font-medium"), "Oliver"), " is typing..."),
		),
	)
}

func DocsGroup() h.Node {
	return col("gap-6",
		message.Group(
			message.Message(
				message.Avatar(),
				message.Content(bubble.Bubble(bubble.VariantMuted, bubble.Content("I checked the registry addresses."))),
			),
			message.Message(
				docsAvatar("/avatars/02.png", "@avatar", "CN"),
				message.Content(bubble.Bubble(bubble.VariantMuted, bubble.Content("The component and example JSON now live under the UI registry."))),
			),
		),
	)
}

func DocsHeaderFooter() h.Node {
	return col("gap-8",
		message.Message(message.Content(
			message.Header("Olivia"),
			bubble.Bubble(bubble.VariantMuted, bubble.Content("I already checked the logs.")),
		)),
		message.Message(message.AlignEnd, message.Content(
			bubble.Bubble(bubble.Content("Send the report to the team. Ping @shadcn if you need help.")),
			message.Footer(h.Div("Read ", h.Span(h.Class("font-normal"), "Yesterday"))),
		)),
	)
}

func DocsMarkdown() h.Node {
	return col("gap-8",
		message.Message(message.AlignEnd, message.Content(
			bubble.Bubble(bubble.Content("How do I render markdown in a message?")),
		)),
		message.Message(message.Content(
			bubble.Bubble(bubble.VariantGhost, bubble.Content(
				h.Div(h.Class("space-y-4 whitespace-normal [&>*:first-child]:mt-0 [&>*:last-child]:mb-0 cn-markdown w-full min-w-0 overflow-hidden"),
					h.P("Here's how to render markdown in a message:"),
					h.Ol(h.Class("list-inside list-decimal whitespace-normal [li_&]:pl-6"), h.Data("streamdown", "ordered-list"),
						mdItem("Render assistant text through ", mdStrong("Markdown"), "."),
						mdItem("Keep user messages as plain text."),
						mdItem("Use a ", mdCode("ghost"), " bubble so the response is unframed."),
					),
				),
			)),
		)),
	)
}

// Streamdown's markdown elements (see examples/bubble).
func mdItem(children ...any) h.Node {
	return h.Li(append([]any{h.Class("py-1 [&>p]:inline"), h.Data("streamdown", "list-item")}, children...)...)
}

func mdStrong(text string) h.Node {
	return h.Span(h.Class("font-semibold"), h.Data("streamdown", "strong"), text)
}

func mdCode(text string) h.Node {
	return h.Code(h.Class("rounded bg-muted px-1.5 py-0.5 font-mono text-sm"), h.Data("streamdown", "inline-code"), text)
}
