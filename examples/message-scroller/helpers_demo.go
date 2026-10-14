package messagescroller

import (
	_ "embed"
	"encoding/json"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	dm "shadcn.go/components/ui/dropdown-menu"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/message"
	ms "shadcn.go/components/ui/message-scroller"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The @shadcn/helpers docs previews (ui/apps/v4/examples/base/
// ai-sdk-helper-demo.tsx and tanstack-ai-helper-demo.tsx, styleName
// base-rhea): a scripted chat (createChat) played through useChat. Both
// render the same card; helpers_demo.js plays the script: the user turn
// appears (slide-up), then after the scripted sleep the assistant reply
// streams word by word every 20ms (splitTextDeltas, delayMs 20), reasoning
// first when the turn has some.
//
//go:embed helpers_demo.js
var helpersScript string

type helperPart struct {
	Type  string `json:"type"` // "text", "reasoning" or "sleep"
	Text  string `json:"text,omitempty"`
	Sleep int    `json:"sleep,omitempty"`
}

type helperTurn struct {
	Role  string       `json:"role"`
	Sleep int          `json:"sleep"` // sleep before the turn streams
	Parts []helperPart `json:"parts"`
}

var helperScript = []helperTurn{
	{Role: "user", Parts: []helperPart{{Type: "text", Text: scrollChat1}}},
	{Role: "assistant", Sleep: 1000, Parts: []helperPart{
		{Type: "reasoning", Text: "They are describing a streaming transcript that keeps taking control of the viewport. I should explain when auto-scroll follows and when it stops."},
		{Type: "sleep", Sleep: 1000},
		{Type: "text", Text: scrollChat2},
	}},
	{Role: "user", Parts: []helperPart{{Type: "text", Text: scrollChat3}}},
	{Role: "assistant", Sleep: 1000, Parts: []helperPart{{Type: "text", Text: "MessageScrollerItem fixes that with turn anchoring. Set `scrollAnchor` on the turn that should settle near the top instead of blindly snapping to the document bottom.\n\nIt also leaves a small peek of the previous exchange visible above the anchor, so context isn't lost. The reply starts in view without that disorienting jump you get from a plain overflow container."}}},
	{Role: "user", Parts: []helperPart{{Type: "text", Text: "And if they've scrolled up to re-read an older answer? I don't want to yank them back down."}}},
	{Role: "assistant", Sleep: 1000, Parts: []helperPart{{Type: "text", Text: "You won't. Auto-scroll only runs when the viewport is already pinned to the bottom, so scrolling up is a deliberate opt-out — their place in the thread stays put even as new tokens keep arriving below.\n\nWhen there is content they haven't seen yet, `MessageScrollerButton` appears at the bottom of the viewport. One tap jumps them back to the newest message and re-engages auto-scroll. Same pattern as Slack or iMessage: quiet when you're caught up, helpful when you're not."}}},
	{Role: "user", Parts: []helperPart{{Type: "text", Text: "Last one — does this work with assistive tech?"}}},
	{Role: "assistant", Sleep: 1000, Parts: []helperPart{{Type: "text", Text: "`MessageScrollerContent` sets `role=\"log\"` and `aria-relevant=\"additions\"` by default, so screen readers announce new messages as they stream in.\n\nThe scroll button is a real `<button>` with an sr-only label, and it's removed from the tab order when you're already at the bottom — no ghost focus stops."}}},
}

func init() {
	runtime.Module("example-helpers-chat", helpersScript)
	for _, name := range []string{"ai-sdk-helper-demo", "tanstack-ai-helper-demo"} {
		registry.Register(registry.Example{Name: "nova--" + name, Component: "message-scroller", Render: registry.Static(func() h.Node {
			return style.Scope(style.Rhea, HelperChatDemo())
		})})
	}
}

// HelperChatDemo is the card of the helpers demos, in its initial state
// (no messages, the first user turn queued in the textarea).
func HelperChatDemo() h.Node {
	script, _ := json.Marshal(helperScript)
	return ms.Provider(
		h.Div(h.Class("relative flex flex-col gap-4"), h.Data("example-helpers-chat", string(script)),
			chatCard(
				card.Header(h.Class("gap-1 border-b"),
					card.Title("New Chat"),
					card.Description("How can I help you today?"),
					card.Action(resetButton("Reset conversation", false, true)),
				),
				card.Content(h.Class("flex-1 overflow-hidden p-0"),
					emptyState("Morning, shadcn!", "What are we working on today? Press send to start a new conversation"),
				),
				card.Footer(h.Class("flex-col gap-2"),
					h.Form(h.Class("w-full"),
						ig.InputGroup(
							ig.Textarea(h.Aria("label", "Next predefined message"),
								h.Class("h-14 min-h-14 overflow-hidden px-3 py-2.5 opacity-60 data-[status=ready]:opacity-100"),
								h.Data("status", "ready"),
								h.Attr("placeholder", "No messages queued. Reset the conversation."),
								h.Bool("readonly"),
								scrollChat1,
							),
							ig.Addon(ig.AlignBlockEnd, h.Class("pt-1"),
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
								ig.Button(h.Type("submit"), button.VariantDefault, ig.ButtonSizeIconSm, h.Class("ml-auto"),
									lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send")),
							),
						),
					),
				),
			),
			caption("px-0.5 text-center text-xs text-muted-foreground", "Demo is read only. Press send to send messages."),
			// Prototypes cloned by helpers_demo.js.
			h.El("template").Data("example-helpers-chat-proto", "").Children(
				ms.MessageScroller(
					ms.Viewport(ms.Content(h.Class("p-(--card-spacing)"), h.Aria("busy", "false"))),
					ms.Button(),
				),
				ms.Item(ms.ScrollAnchor(true), h.Data("role", "user"),
					message.Message(message.AlignEnd, message.Content(
						bubble.Bubble(bubble.VariantMuted, bubble.Content(h.Class("space-y-2"))),
					)),
				),
				ms.Item(ms.ScrollAnchor(false), h.Data("role", "assistant"),
					message.Message(message.AlignStart, message.Content(
						bubble.Bubble(bubble.VariantGhost, bubble.Content(h.Class("space-y-2"))),
					)),
				),
				h.Div(h.Class("w-full border-l-2 border-muted-foreground/30 pl-3 text-muted-foreground"), h.Data("role", "reasoning"),
					h.Div(h.Class("mb-1 flex items-center gap-1.5 text-xs font-medium"), lucide.Brain(h.Class("size-3.5")), "Reasoning"),
					h.Div(h.Class("space-y-1.5 text-sm")),
				),
			),
		),
	)
}
