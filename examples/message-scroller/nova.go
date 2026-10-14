package messagescroller

import (
	"fmt"
	"regexp"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/empty"
	hovercard "shadcn.go/components/ui/hover-card"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/message"
	ms "shadcn.go/components/ui/message-scroller"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/slider"
	"shadcn.go/components/ui/tabs"
	tg "shadcn.go/components/ui/toggle-group"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the message-scroller page
// (ui/apps/v4/examples/radix/message-scroller-*.tsx), in their initial
// state (upstream simulates chats with the AI SDK; the Go ports render the
// opening transcript and the controls). Upstream's examples import the
// radix-rhea components, so they render in style.Rhea.
func init() {
	docs := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "message-scroller", Render: registry.Static(func() h.Node {
			return style.Scope(style.Rhea, fn())
		})}
	}
	registry.Register(
		docs("message-scroller-anchoring", DocsAnchoring),
		docs("message-scroller-animation", DocsAnimation),
		docs("message-scroller-commands", DocsCommands),
		docs("message-scroller-demo", DocsDemo),
		docs("message-scroller-group-chat", DocsGroupChat),
		docs("message-scroller-load-history", DocsLoadHistory),
		docs("message-scroller-opening-position", DocsOpeningPosition),
		docs("message-scroller-previous-context", DocsPreviousContext),
		docs("message-scroller-scrollable", DocsScrollable),
		docs("message-scroller-streaming", DocsStreaming),
		docs("message-scroller-visibility", DocsVisibility),
	)
}

type chatMessage struct {
	id, role, text string
}

var paragraphSplit = regexp.MustCompile(`\n\s*\n`)

// paragraphs is text.split(/\n\s*\n/).map(trim).filter(Boolean) as <p>s.
func paragraphs(text string) []any {
	var out []any
	for _, p := range paragraphSplit.Split(text, -1) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, h.P(h.Class("whitespace-pre-wrap"), p))
		}
	}
	return out
}

// item is a transcript row: a MessageScrollerItem holding a Message with
// one bubble of paragraphs (MessageAnimated, and the examples' inline
// rows).
func item(m chatMessage, anchor bool) h.Node {
	user := m.role == "user"
	align, variant := message.AlignStart, bubble.VariantGhost
	if user {
		align, variant = message.AlignEnd, bubble.VariantMuted
	}
	return ms.Item(ms.MessageID(m.id), ms.ScrollAnchor(anchor),
		message.Message(align, message.Content(
			bubble.Bubble(variant, bubble.Content(append([]any{h.Class("space-y-2")}, paragraphs(m.text)...)...)),
		)),
	)
}

func transcript(msgs []chatMessage, anchor func(chatMessage) bool) []any {
	out := make([]any, len(msgs))
	for i, m := range msgs {
		out[i] = item(m, anchor(m))
	}
	return out
}

func userAnchor(m chatMessage) bool { return m.role == "user" }

// scroller is MessageScroller > Viewport > Content (p-(--card-spacing)) +
// Button.
func scroller(contentAttrs []any, rows ...any) h.Node {
	return ms.MessageScroller(
		ms.Viewport(ms.Content(append(append([]any{h.Class("p-(--card-spacing)")}, contentAttrs...), rows...)...)),
		ms.Button(),
	)
}

func resetButton(label string, disabled bool, withTooltip bool) h.Node {
	btn := button.Button(h.Type("button"), button.VariantOutline, button.SizeIcon, h.Aria("label", label), h.If(disabled, h.Bool("disabled")), lucide.RotateCw())
	if !withTooltip {
		return btn
	}
	return tooltip.Tooltip(
		tooltip.Trigger(lib.AsChild, btn),
		tooltip.Content(h.P("Reset")),
	)
}

func emptyState(title, description string) h.Node {
	return empty.Empty(h.Class("h-full"),
		empty.Header(
			empty.Media(empty.MediaIcon, lucide.MessageCircleDashed()),
			empty.Title(title),
			empty.Description(description),
		),
	)
}

func caption(class, text string) h.Node {
	return h.Div(h.Class(class), text)
}

// chatFooter is the input footer of demo, streaming and previous-context:
// the next scripted message and the send controls.
func chatFooter(next string, extra ...any) h.Node {
	addon := []any{ig.AlignBlockEnd, h.Class("pt-1"),
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild, ig.Button(h.Aria("label", "Add files"), h.Type("button"), ig.ButtonSizeIconSm, button.VariantOutline, lucide.Plus())),
			dm.Content(dm.AlignStart, dm.SideTop, h.Class("w-44"),
				dm.Item(lucide.Paperclip(), "Add Photos & Files"),
				dm.Separator(),
				dm.Item(lucide.Image(), "Create Image"),
				dm.Item(lucide.Telescope(), "Deep Research"),
				dm.Item(lucide.Globe(), "Web Search"),
			),
		),
	}
	addon = append(addon, extra...)
	addon = append(addon, ig.Button(h.Type("submit"), button.VariantDefault, ig.ButtonSizeIconSm, h.Class("ml-auto"),
		lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send")))
	return card.Footer(h.Class("flex-col gap-2"),
		h.Form(h.Class("w-full"),
			ig.InputGroup(
				h.Div(h.Class("h-14 w-full px-3 py-2.5"),
					h.Span(h.Class("line-clamp-2 opacity-60 data-[status=ready]:opacity-100"), h.Data("status", "ready"), next),
				),
				ig.Addon(addon...),
			),
		),
	)
}

const (
	scrollChat1 = "I'm building a chat for our app and the scroll behavior is driving me nuts. Every time the AI streams a reply, the whole thread jumps around."
	scrollChat2 = "That's the classic streaming scroll problem. Wrap your message list in `MessageScroller` and turn on `autoScroll` — the viewport pins to the bottom as tokens arrive, so users always see the latest text land in place.\n\nThe important part: it only auto-scrolls while the reader is already at the bottom. The moment they scroll up to read something earlier, auto-scroll backs off and their position is preserved. You get smooth streaming without fighting the user's intent."
	scrollChat3 = "Okay, but when someone sends a new message the view still feels jarring — like the whole conversation reloads from the top."
)

func chatCard(header h.Node, content h.Node, footer h.Node) h.Node {
	return card.Card(h.Class("mx-auto h-140 w-full max-w-sm gap-0"), header, content, footer)
}

func DocsDemo() h.Node {
	return ms.Provider(
		h.Div(h.Class("relative flex flex-col gap-4"),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("New Chat"),
					card.Description("How can I help you today?"),
					card.Action(resetButton("Reset conversation", false, true)),
				),
				card.Content(h.Class("flex-1 overflow-hidden p-0"),
					emptyState("Morning, shadcn!", "What are we working on today? Press send to start a new conversation"),
				),
				chatFooter(scrollChat1),
			),
			caption("px-0.5 text-center text-xs text-muted-foreground", "Demo is read only. Press send to send messages."),
		),
	)
}

func DocsStreaming() h.Node {
	return ms.Provider(ms.AutoScroll(true),
		h.Div(h.Class("relative flex flex-col gap-4"),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("Streaming Messages"),
					card.Description("Auto-scroll follows the live edge of the conversation."),
					card.Action(resetButton("Reset stream", true, true)),
				),
				card.Content(h.Class("flex-1 overflow-hidden p-0"),
					emptyState("Ready to Stream", "Press send to stream a scripted launch summary."),
				),
				chatFooter(scrollChat1),
			),
			caption("px-0.5 text-center text-xs text-muted-foreground", "Streaming is simulated. `autoScroll` is enabled."),
		),
	)
}

func DocsPreviousContext() h.Node {
	initial := []chatMessage{{"msg-1", "user", scrollChat1}, {"msg-2", "assistant", scrollChat2}}
	return ms.Provider(ms.ScrollMargin(24), ms.ScrollPreviousItemPeek(64),
		h.Div(h.Class("relative flex flex-col gap-4"),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("Keeping Context Visible"),
					card.Description("New turns keep part of the previous reply in view."),
					card.Action(resetButton("Reset context example", false, true)),
				),
				card.Content(h.Class("flex-1 overflow-hidden p-0"),
					scroller([]any{h.Aria("busy", "false")}, transcript(initial, userAnchor)...),
				),
				chatFooter(scrollChat3,
					h.Div(h.Class("flex w-28 items-center gap-2"),
						h.Span(h.Class("text-xs text-muted-foreground tabular-nums"), "64", h.Comment(""), "px"),
						slider.Slider(h.Aria("label", "Previous context peek"), slider.DefaultValue(64), slider.Min(64), slider.Max(128), slider.Step(1)),
					),
				),
			),
			caption("px-0.5 text-center text-xs text-muted-foreground", "Adjust the slider and send. Observe the previous message peak"),
		),
	)
}

func DocsAnimation() h.Node {
	presets := []struct{ id, name string }{
		{"fade", "Fade"}, {"slide-up", "Slide Up"}, {"slide-side", "Slide Side"}, {"pop", "Pop"},
		{"spring-bounce", "Spring Bounce"}, {"blur-fade", "Blur Fade"}, {"scale-fade", "Scale Fade"},
	}
	var items []any
	for _, p := range presets {
		items = append(items, selectui.Item(selectui.ItemValue(p.id), p.name))
	}
	return h.Div(h.Class("relative flex flex-col gap-4"),
		chatCard(
			card.Header(h.Class("border-b"),
				card.Title("Animation"),
				card.Description("Choose how user messages are animated when they are added to the conversation."),
				card.Action(h.Class("flex items-center gap-2"), resetButton("Reset animated messages", true, false)),
			),
			card.Content(h.Class("min-h-0 flex-1 overflow-hidden p-0"),
				emptyState("No Messages Yet", "Click the button below to send the first message."),
			),
			card.Footer(h.Class("border-t"),
				selectui.Select(selectui.DefaultValue("fade"),
					selectui.Trigger(h.Aria("label", "Animation preset"), selectui.Value()),
					selectui.Content(selectui.AlignStart, selectui.SideTop, selectui.PositionPopper, selectui.Group(items...)),
				),
				button.Button(h.Type("button"), button.SizeIcon, h.Class("ml-auto"), lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send Message")),
			),
		),
		caption("mx-auto max-w-sm px-0.5 text-center text-xs text-balance text-muted-foreground", "Select an animation then click send to see it in action."),
	)
}

func DocsAnchoring() h.Node {
	return h.Div(h.Class("relative flex flex-col gap-4"),
		chatCard(
			card.Header(h.Class("border-b"),
				card.Title("Anchoring Turns"),
				card.Description("Choose which role settles near the top edge."),
				card.Action(resetButton("Reset anchored turns", true, false)),
			),
			card.Content(h.Class("min-h-0 flex-1 overflow-hidden p-0"),
				emptyState("No anchored messages yet", "Send the first message to see the selected role anchor."),
			),
			card.Footer(
				tg.ToggleGroup(tg.TypeSingle, tg.Value("user"), h.Aria("label", "Select scroll anchor role"),
					tg.Item(tg.Value("user"), h.Aria("label", "Anchor user messages"), "User"),
					tg.Item(tg.Value("assistant"), h.Aria("label", "Anchor assistant messages"), "Assistant"),
				),
				button.Button(h.Type("button"), button.SizeIcon, h.Class("ml-auto"), lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send Message")),
			),
		),
		caption("mx-auto max-w-xs px-0.5 text-center text-xs text-muted-foreground", "Toggle the anchor role, then send messages to compare where turns settle."),
	)
}

var commandMessages = []chatMessage{
	{"command-activation", "user", "We're seeing activation dip after workspace creation. Can you help me find the likely step?"},
	{"msg-1", "assistant", "The sharpest drop is between creating the workspace and inviting the first teammate.\n\nWorkspace creation is still healthy, but the invite step is where users pause. That suggests the product is asking for collaboration before the user has enough confidence in the workspace."},
	{"command-compare", "user", "What should I compare before we change the onboarding flow?"},
	{"msg-2", "assistant", "Compare three cohorts:\n\n1. Users who choose a template before inviting teammates.\n2. Users who start from a blank workspace.\n3. Users who skip invites and return within 24 hours.\n\nIf template users invite faster, the fix is probably better first-run guidance rather than a louder invite prompt."},
	{"command-experiment", "user", "Can you turn that into an experiment?"},
	{"msg-3", "assistant", "Yes. Create a variant that shows a short checklist after workspace creation:\n\n- Pick a template.\n- Add one project detail.\n- Invite a teammate when the workspace has context.\n\nMeasure first invite completion, 24-hour return rate, and whether teams create a second project."},
	{"command-risk", "user", "What's the risk if we delay the invite prompt?"},
	{"msg-4", "assistant", "The main risk is reducing team creation for accounts that already know who they want to invite.\n\nTo protect that path, keep the invite action visible in the header and only change the primary empty-state guidance. That gives confident teams a direct route without forcing uncertain users through the invite step too early."},
}

// trimmed is getTrimmedMessageText.
func trimmed(text string) string {
	if len(text) > 42 {
		return text[:39] + "..."
	}
	return text
}

func DocsCommands() h.Node {
	var jumps []any
	for _, m := range commandMessages {
		if m.role == "user" {
			jumps = append(jumps, dm.Item(h.Data("example-message-scroller-jump", m.id), h.Span(h.Class("line-clamp-1 min-w-0"), trimmed(m.text))))
		}
	}
	return ms.Provider(ms.PositionEnd,
		h.Div(h.Class("relative flex flex-col gap-4"),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("Commands"),
					card.Description("Drive the transcript from outside."),
					card.Action(dm.DropdownMenu(
						dm.Trigger(lib.AsChild, button.Button(h.Type("button"), button.VariantSecondary, "Jump to...")),
						dm.Content(append([]any{dm.AlignEnd, dm.SideBottom, h.Class("w-64"), dm.Label("Conversations")}, jumps...)...),
					)),
				),
				card.Content(h.Class("flex-1 overflow-hidden p-0"),
					scroller(nil, transcript(commandMessages, userAnchor)...),
				),
				nil,
			),
			caption("mx-auto max-w-sm px-0.5 text-center text-xs text-balance text-muted-foreground", "Use the controls to jump to any message in the conversation."),
		),
	)
}

func DocsGroupChat() h.Node {
	type groupItem struct{ id, sender, role, text string }
	items := []groupItem{
		{"group-1", "Grace", "participant", "@mary, the astrophage line keeps matching Venus energy output. Can you check my math?"},
		{"group-2", "Mary (Agent)", "assistant", "Yes. Confirmed. The curve points to a microorganism harvesting stellar energy and breeding near carbon dioxide. If @rocky agrees, this is the clue we need."},
		{"group-3", "Grace", "participant", "ping @rocky"},
	}
	var rows []any
	for _, it := range items {
		current := it.sender == "Grace"
		variant, align := bubble.VariantTinted, message.AlignStart
		if current {
			variant, align = bubble.VariantMuted, message.AlignEnd
		} else if it.role == "assistant" {
			variant = bubble.VariantGhost
		}
		var header any
		if !current {
			header = message.Header(it.sender)
		}
		rows = append(rows, ms.Item(ms.MessageID(it.id), ms.ScrollAnchor(it.id == "group-3"),
			message.Message(align, message.Content(header, bubble.Bubble(variant, bubble.Content(it.text)))),
		))
	}
	return ms.Provider(
		h.Div(h.Class("relative flex flex-col gap-4"),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("Group Chat"),
					card.Description("A group chat with several participants and an assistant. The Marker is marked as a turn."),
					card.Action(resetButton("Reset conversation", true, true)),
				),
				card.Content(h.Class("min-h-0 flex-1 p-0"),
					ms.Provider(scroller(nil, rows...)),
				),
				card.Footer(h.Class("flex flex-col items-center gap-2 border-t"),
					button.Button(h.Type("button"), h.Class("w-full"), button.VariantSecondary, "Add Rocky"),
					h.P(h.Class("text-xs text-muted-foreground"), "This will create a marker and make it the anchor"),
				),
			),
			caption("mx-auto max-w-sm px-0.5 text-center text-xs text-balance text-muted-foreground", "When a user joins, a marker is created. scrollAnchor on the marker marks it as the next turn"),
		),
	)
}

func DocsLoadHistory() h.Node {
	history := []chatMessage{
		{"msg-1", "user", "Can you summarize the incident channel?"},
		{"msg-2", "assistant", "The first alert was a delayed export job. It started backing up around 09:42 UTC and triggered the warning once the retry queue crossed the threshold.\n\nNo customer-facing checkout paths were affected, but exports for larger workspaces were running about 12 minutes behind."},
		{"msg-3", "user", "Was checkout affected?"},
		{"msg-4", "assistant", "No checkout errors were reported. Payment authorization, order creation, and confirmation emails stayed inside their normal latency bands.\n\nThe only elevated metric was export queue depth, which maps to analytics downloads instead of checkout."},
		{"msg-5", "user", "What changed in the last deploy?"},
		{"msg-6", "assistant", "Only the export queue worker changed. The deploy moved large CSV jobs onto the shared retry policy, which made each failed attempt hold a worker slot longer than before.\n\nThe app deploy did not include checkout, pricing, or billing API changes."},
		{"msg-7", "user", "Do we need to roll back?"},
		{"msg-8", "assistant", "Not yet. Queue depth is recovering after we reduced retry concurrency, and the oldest pending job is now under five minutes old.\n\nKeep rollback ready if the queue starts climbing again, but the current trend points toward recovery."},
		{"msg-9", "user", "Keep watching for customer-visible issues."},
		{"msg-10", "assistant", "I will watch the queue and support tags for another 15 minutes. I am tracking export failures, delayed download requests, and any support thread that mentions missing reports.\n\nIf those stay quiet through the next batch window, we can close this as an internal degradation."},
	}
	rows := transcript(history[len(history)-5:], func(chatMessage) bool { return false })
	rows = append(rows, ms.Item(ms.ScrollAnchor(false), marker.Marker(marker.VariantSeparator, marker.Content("End of Conversation"))))
	return ms.Provider(
		h.Div(h.Class("relative flex flex-col gap-4"),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("Load History"),
					card.Description("Prepended messages keep your place."),
					card.Action(resetButton("Reset loaded messages", true, true)),
				),
				card.Content(h.Class("flex-1 overflow-hidden p-0"), scroller(nil, rows...)),
				card.Footer(h.Class("flex flex-col items-center gap-2 border-t"),
					button.Button(h.Type("button"), h.Class("w-full"), button.VariantSecondary, "Load History"),
					h.P(h.Class("text-xs text-muted-foreground"), "Restore earlier messages while keeping your place."),
				),
			),
			caption("mx-auto max-w-sm px-0.5 text-center text-xs text-balance text-muted-foreground", "Click Load History to load the entire conversation"),
		),
	)
}

func DocsOpeningPosition() h.Node {
	msgs := []chatMessage{
		{"open-1", "user", "This is the first message the user sent in the conversation."},
		{"open-2", "assistant", "Workspace creation rose 8%, but first invite completion only rose 2%."},
		{"open-3", "user", "This is the last message the user sent in the conversation."},
		{"open-4", "assistant", "Start with the invite step. Teams are creating workspaces but waiting to add collaborators.\n\nRecommended follow-up:\n\n1. Compare invite drop-off by account size.\n2. Check whether users who skip invites still return within 24 hours.\n3. Review the empty-state copy on the first project screen.\n4. Segment activation by template, since template users may not need invites right away.\n\nIf that pattern holds, the next experiment should make collaboration useful earlier instead of prompting for invites harder."},
	}
	var triggers []any
	for _, p := range []string{"start", "end", "last-anchor"} {
		triggers = append(triggers, tabs.Trigger(tabs.Value(p), p))
	}
	return h.Div(h.Class("relative flex flex-col gap-4"),
		chatCard(
			card.Header(h.Class("gap-1 border-b"),
				card.Title("Opening Position"),
				card.Description("Choose where a saved transcript opens."),
			),
			card.Content(h.Class("flex-1 overflow-hidden p-0"),
				ms.Provider(ms.PositionLastAnchor, scroller(nil, transcript(msgs, userAnchor)...)),
			),
			card.Footer(h.Class("flex items-center justify-center border-t"),
				tabs.Tabs(tabs.Value("last-anchor"), h.Class("w-full"),
					tabs.List(append([]any{h.Class("w-full")}, triggers...)...),
				),
			),
		),
		caption("mx-auto max-w-sm px-0.5 text-center text-xs text-muted-foreground", "Toggle the defaultScrollPosition to see where the transcript starts when you open the thread"),
	)
}

func DocsScrollable() h.Node {
	var msgs []chatMessage
	for i := 0; i < 12; i++ {
		m := chatMessage{id: fmt.Sprintf("scrollable-%d", i+1)}
		if i%2 == 0 {
			m.role, m.text = "user", fmt.Sprintf("Review scroll checkpoint %d.", i+1)
		} else {
			m.role, m.text = "assistant", fmt.Sprintf("Checkpoint %d is synced. The scrollable hook updates as the viewport moves.\n\nWhen the reader is at the first message, the footer should only point them down. Once they move into the middle of the transcript, it should explain that both directions are available.\n\nAt the latest message, the footer should switch again and only point them back up.", i+1)
		}
		msgs = append(msgs, m)
	}
	return h.Div(h.Class("mx-auto flex w-full max-w-sm flex-col gap-4"),
		card.Card(h.Class("h-140 w-full gap-0 overflow-hidden"),
			card.Header(h.Class("gap-1 border-b"),
				card.Title("Scroll Status"),
				card.Description("Where the reader can go scroll to based on current scroll position."),
			),
			ms.Provider(ms.PositionStart,
				card.Content(h.Class("flex-1 overflow-hidden p-0"),
					ms.MessageScroller(
						ms.Viewport(ms.Content(append([]any{h.Class("gap-4 p-(--card-spacing)")}, transcript(msgs, userAnchor)...)...)),
						ms.Button(),
					),
				),
				card.Footer(h.Class("justify-center border-t text-center text-sm text-muted-foreground"),
					h.Data("example-message-scroller-status", ""),
					"You are at the top. You can only scroll down.",
				),
			),
		),
		caption("px-0.5 text-center text-xs text-muted-foreground", "Scroll the transcript to see the footer update."),
	)
}

func DocsVisibility() h.Node {
	msgs := []chatMessage{
		{"vis-brief", "user", "Review the incident handoff and tell me what to read first."},
		{"msg-1", "assistant", "Start with the summary and the impact section. The regression affected the upload queue, but the recovery path completed for every queued job."},
		{"vis-impact", "user", "What was the customer impact?"},
		{"msg-2", "assistant", "Impact was limited to delayed processing.\n\nNo records were dropped, and the reconciliation worker confirmed each retry batch. Support saw confusion from two customers, but there were no checkout or billing errors."},
		{"vis-actions", "user", "What actions are open?"},
		{"msg-3", "assistant", "Keep the retry window enabled until the next deploy, then add a queue-depth alert as the long-term fix.\n\nThe alert should fire on sustained queue growth, not a single short spike."},
		{"vis-checklist", "user", "Give me the follow-up checklist."},
		{"msg-4", "assistant", "After that, compare the queue recovery graph with the deploy timeline so the handoff shows exactly when processing returned to baseline. That makes it easier for support and engineering to answer the same customer questions without re-reading the whole incident thread.\n\nI would also add a short owner note beside each follow-up item. The checklist is small, but ownership keeps the retry-window decision, alert tuning, and support macro from drifting into separate follow-up conversations.\n\nKeep the retry window enabled until the next deploy, then add a queue-depth alert as the long-term fix.\n\nThe alert should fire on sustained queue growth, not a single short spike."},
	}
	var bars, links []any
	for _, m := range msgs {
		if m.role != "user" {
			continue
		}
		current := m.id == "vis-checklist"
		bars = append(bars, h.Span(h.Data("current", fmt.Sprint(current)), h.Class("h-0.5 w-4 rounded-full bg-muted-foreground/40 data-[current=true]:bg-foreground")))
		links = append(links, h.Button(h.Type("button"), h.If(current, h.Aria("current", "location")),
			h.Data("example-message-scroller-jump", m.id),
			h.Class("flex min-h-7 items-center rounded-xl px-2 py-1.5 text-left text-sm transition-colors outline-none hover:bg-accent hover:text-accent-foreground focus-visible:bg-accent focus-visible:text-accent-foreground aria-current:bg-accent aria-current:text-accent-foreground"),
			h.Span(h.Class("line-clamp-1 min-w-0"), trimmed(m.text)),
		))
	}
	return ms.Provider(ms.ScrollMargin(12),
		h.Div(h.Class("relative flex flex-col gap-4"),
			h.Div(h.Class("relative mx-auto w-full max-w-sm"),
				card.Card(h.Class("h-140 w-full gap-0"),
					card.Header(h.Class("gap-1 border-b"),
						card.Title("Transcript Outline"),
						card.Description("Track the current anchored turn."),
					),
					card.Content(h.Class("flex-1 overflow-hidden p-0"),
						scroller(nil, transcript(msgs, userAnchor)...),
					),
				),
				h.Div(h.Class("absolute top-1/2 -right-12 -translate-y-1/2"),
					hovercard.HoverCard(hovercard.OpenDelay(0), hovercard.CloseDelay(0),
						hovercard.Trigger(lib.AsChild,
							h.Button(h.Type("button"), h.Aria("label", "Open transcript outline"),
								h.Class("flex h-9 w-9 flex-col items-center justify-center gap-1 rounded-md transition-colors outline-none focus-visible:ring-3 focus-visible:ring-ring/50"),
								bars,
							),
						),
						hovercard.Content(hovercard.AlignCenter, hovercard.SideLeft, hovercard.SideOffset(-28),
							h.Class("flex w-64 flex-col gap-1 rounded-2xl p-1"),
							links,
						),
					),
				),
			),
			caption("mx-auto max-w-sm px-0.5 text-center text-xs text-muted-foreground", "Open the outline to jump between anchored turns as you read."),
		),
	)
}
