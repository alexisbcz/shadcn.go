package bubble

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the bubble page (ui/apps/v4/examples/radix/bubble-*.tsx).
// Upstream's examples import the radix-rhea components, so they render in
// style.Rhea.
// collapsible.js swaps bubble-collapsible's text when it opens.
//
//go:embed collapsible.js
var collapsibleScript string

func init() {
	runtime.Module("example-bubble-collapsible", collapsibleScript)
	docs := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "bubble", Render: registry.Static(func() h.Node {
			return style.Scope(style.Rhea, fn())
		})}
	}
	registry.Register(
		docs("bubble-alignment", DocsAlignment),
		docs("bubble-collapsible", DocsCollapsible),
		docs("bubble-demo", DocsDemo),
		docs("bubble-group-demo", DocsGroupDemo),
		docs("bubble-link-button", DocsLinkButton),
		docs("bubble-markdown", DocsMarkdown),
		docs("bubble-popover", DocsPopover),
		docs("bubble-reactions", DocsReactions),
		docs("bubble-tooltip", DocsTooltip),
		docs("bubble-variants", DocsVariants),
	)
}

func col(gap string, children ...any) h.Node {
	return h.Div(append([]any{h.Class("flex w-full max-w-sm flex-col " + gap + " py-12")}, children...)...)
}

func reactions(label string, emojis ...string) h.Node {
	kids := []any{h.Role("img"), h.Aria("label", label)}
	for _, e := range emojis {
		kids = append(kids, h.Span(e))
	}
	return bubble.Reactions(kids...)
}

func DocsAlignment() h.Node {
	return col("gap-8",
		bubble.Bubble(bubble.VariantMuted, bubble.Content("This bubble is aligned to the start. This is the default alignment.")),
		bubble.Bubble(bubble.AlignEnd, bubble.Content("This bubble is aligned to the end. Use this for user messages.")),
	)
}

const collapsibleText = "The accessibility review found two focus states that were visually too subtle in dark mode.\n\n" +
	"I checked the dialog, menu, and drawer paths because each one renders focusable controls inside a layered surface.\n\n" +
	"The dialog and drawer are fine. The menu needs the hover and focus tokens split so keyboard focus stays visible when the pointer is not involved.\n\n" +
	"I also recommend keeping the change in the style file instead of the primitive so the other themes can choose their own focus treatment later."

func DocsCollapsible() h.Node {
	preview := collapsibleText[:180] + "..."
	return col("gap-8",
		bubble.Bubble(bubble.VariantMuted, bubble.Content("How can I help you today?")),
		bubble.Bubble(bubble.VariantMuted, bubble.AlignEnd,
			bubble.Content(h.Class("whitespace-pre-line"),
				collapsible.Collapsible(
					h.Div(h.Data("example-bubble-collapsible", ""), h.Data("full", collapsibleText), preview),
					collapsible.Trigger(lib.AsChild,
						button.Button(button.VariantLink, h.Class("gap-1 p-0 text-muted-foreground"),
							"Show more",
							lucide.ChevronDown(h.Data("icon", "inline-end"), h.Class("group-data-open/button:rotate-180")),
						),
					),
				),
			),
		),
	)
}

func DocsDemo() h.Node {
	return col("gap-8",
		bubble.Bubble(bubble.AlignEnd, bubble.Content("Hey there! what's up?")),
		bubble.Group(
			bubble.Bubble(bubble.VariantMuted, bubble.Content("Hey! Want to see chat bubbles?")),
			bubble.Bubble(bubble.VariantMuted,
				bubble.Content("I can group messages, switch sides, and keep the whole thread easy to scan."),
				reactions("Reaction: thumbs up", "👍"),
			),
		),
		bubble.Bubble(bubble.AlignEnd, bubble.Content("Sure. Hit me with your best demo.")),
		bubble.Bubble(bubble.VariantMuted,
			bubble.Content("Yes. You are reading a demo that is demoing itself. Very meta. Very on-brand."),
			reactions("Reactions: thumbs up, fire, eyes, and 2 more", "👍", "🔥", "👀", "+2"),
		),
	)
}

func DocsGroupDemo() h.Node {
	return col("gap-8",
		bubble.Bubble(bubble.VariantMuted, bubble.Content("Can you tell me what's the issue?")),
		bubble.Group(
			bubble.Bubble(bubble.AlignEnd, bubble.Content("You tell me!")),
			bubble.Bubble(bubble.AlignEnd, bubble.Content("It worked yesterday. You broke it!")),
			bubble.Bubble(bubble.AlignEnd,
				bubble.Content("Find the bug and fix it."),
				bubble.Reactions(h.Aria("label", "Reactions: eyes"), bubble.AlignStart, h.Span("👀")),
			),
		),
		bubble.Bubble(bubble.VariantMuted, bubble.Content("Want me to diff yesterday's you against today's you? It's a bit embarrassing.")),
	)
}

func DocsLinkButton() h.Node {
	btn := func(text string) h.Node {
		return bubble.Bubble(bubble.VariantTinted, bubble.AlignEnd, bubble.Content(lib.AsChild, h.Button(text)))
	}
	return col("gap-8",
		bubble.Bubble(bubble.VariantMuted, bubble.Content("How can I help you today?")),
		bubble.Group(
			btn("I forgot my password"),
			btn("I need help with my subscription"),
			btn("Something else. Talk to a human."),
		),
	)
}

func DocsMarkdown() h.Node {
	return col("gap-8",
		bubble.Bubble(bubble.AlignEnd, bubble.VariantMuted, bubble.Content(
			markdown(h.P("Hello! Are you actually ", mdStrong("thinking"), "?")),
		)),
		bubble.Bubble(bubble.VariantGhost, bubble.Content(ghostMarkdown())),
	)
}

func DocsPopover() h.Node {
	return col("gap-4",
		bubble.Bubble(bubble.AlignEnd, bubble.Content("Run the build script.")),
		bubble.Bubble(bubble.VariantDestructive,
			bubble.Content("Failed to run the command."),
			bubble.Reactions(
				popover.Popover(
					popover.Trigger(lib.AsChild,
						button.Button(button.VariantGhost, button.SizeIconXs, h.Aria("label", "Show error details"), h.Class("aria-expanded:text-destructive"), lucide.Info()),
					),
					popover.Content(
						popover.Header(
							popover.Title(h.Class("text-sm"), "Command failed with exit code 1"),
							popover.Description(h.Class("text-sm"), "ENOENT: no such file or directory, open pnpm-lock.yaml"),
						),
					),
				),
			),
		),
	)
}

func DocsReactions() h.Node {
	return col("gap-12",
		bubble.Bubble(bubble.VariantMuted, bubble.AlignEnd,
			bubble.Content("I don't need tests, I know my code works."),
			bubble.Reactions(bubble.AlignStart, h.Role("img"), h.Aria("label", "Reactions: thumbs up, surprised"), h.Span("👍"), h.Span("😮")),
		),
		bubble.Bubble(bubble.VariantMuted,
			bubble.Content("Bold. Fine I'll add some tests. I'll let you know when they're done."),
			reactions("Reactions: eyes, rocket, and 2 more", "👀", "🚀", "+2"),
		),
		bubble.Bubble(bubble.VariantDefault, bubble.AlignEnd,
			bubble.Content("Tests passed on the first try. All 142 of them. Looking good!"),
			bubble.Reactions(bubble.SideTop, bubble.AlignStart, h.Role("img"), h.Aria("label", "Reactions: party popper, clapping hands"), h.Span("🎉"), h.Span("👏")),
		),
		bubble.Bubble(bubble.VariantDestructive,
			bubble.Content("Are you sure I can run this command?"),
			bubble.Reactions(button.Button(button.VariantGhost, button.SizeXs, "Yes, run it")),
		),
	)
}

func DocsTooltip() h.Node {
	return col("gap-4",
		bubble.Bubble(bubble.VariantSecondary, bubble.Content("Did you remove the stale route?")),
		bubble.Bubble(bubble.AlignEnd,
			bubble.Content("Yes, removed it from the registry."),
			bubble.Reactions(h.Class("p-0"),
				tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIconXs, lucide.Check())),
					tooltip.Content("Read on Jan 5, 2026 at 4:32 PM"),
				),
			),
		),
	)
}

func DocsVariants() h.Node {
	return col("gap-12",
		bubble.Bubble(bubble.Content("This is the default primary bubble.")),
		bubble.Bubble(bubble.VariantSecondary, bubble.AlignEnd, bubble.Content("This is the secondary variant.")),
		bubble.Bubble(bubble.VariantMuted,
			bubble.Content("This one is muted. It uses a lower emphasis color for the chat bubble."),
			reactions("Reaction: thumbs up", "👍"),
		),
		bubble.Bubble(bubble.VariantTinted, bubble.AlignEnd, bubble.Content("This one is tinted. The tint is a softer color derived from the primary color.")),
		bubble.Bubble(bubble.VariantOutline, bubble.Content("We can also use an outlined variant.")),
		bubble.Bubble(bubble.VariantDestructive, bubble.AlignEnd,
			bubble.Content("Or a destructive variant with a reaction."),
			reactions("Reaction: fire", "🔥"),
		),
		bubble.Bubble(bubble.VariantGhost, bubble.Content(ghostMarkdown())),
	)
}

// markdown is the output of the docs' Markdown (Streamdown): one <p> per
// line, strong and inline code as Streamdown renders them.
func markdown(children ...any) h.Node {
	return h.Div(append([]any{h.Class("space-y-4 whitespace-normal [&>*:first-child]:mt-0 [&>*:last-child]:mb-0 cn-markdown w-full min-w-0 overflow-hidden")}, children...)...)
}

func mdStrong(text string) h.Node {
	return h.Span(h.Class("font-semibold"), h.Data("streamdown", "strong"), text)
}

func mdCode(text string) h.Node {
	return h.Code(h.Class("rounded bg-muted px-1.5 py-0.5 font-mono text-sm"), h.Data("streamdown", "inline-code"), text)
}

func ghostMarkdown() h.Node {
	return markdown(
		h.P("Ghost bubbles work for assistant text, ", mdStrong("markdown"), ", and other content that should not be framed."),
		h.P("This is perfect for assistant messages that should not have a frame and can take the full width of the container. You can also render ", mdCode("code"), " in it."),
		h.P("Ghost bubbles are full width and can take the full width of the container."),
	)
}
