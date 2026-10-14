// Package messagescroller ports the message-scroller parity fixtures
// (upstream ships MessageScroller without new-york-v4 examples).
package messagescroller

import (
	_ "embed"
	"fmt"
	"slices"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/bubble"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/message"
	ms "shadcn.go/components/ui/message-scroller"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// dynamic.js keeps message-scroller-dynamic's rows and status like the
// fixture's React state.
//
//go:embed dynamic.js
var dynamicScript string

func init() {
	runtime.Module("example-message-scroller-dynamic", dynamicScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "message-scroller", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("message-scroller-demo", Demo),
		ex("message-scroller-positions", Positions),
		ex("message-scroller-dynamic", Dynamic),
	)
}

// Demo is the parity fixture message-scroller-demo.
func Demo() h.Node {
	items := make([]any, 0, 12)
	for i := 0; i < 12; i++ {
		user := i%2 == 0
		text := fmt.Sprintf("Message %d", i+1)
		if i%3 == 0 {
			text = fmt.Sprintf("Message %d: a longer message that wraps onto a few lines so the transcript is taller than its viewport and needs to scroll.", i+1)
		}
		align, variant := message.AlignStart, bubble.VariantMuted
		if user {
			align, variant = message.AlignEnd, bubble.VariantDefault
		}
		items = append(items, ms.Item(ms.MessageID(fmt.Sprintf("m%d", i+1)), ms.ScrollAnchor(user),
			message.Message(align, message.Content(bubble.Bubble(variant, bubble.Content(text)))),
		))
	}
	return h.Div(h.Class("h-[420px] w-full max-w-md rounded-xl border"),
		ms.Provider(
			ms.MessageScroller(
				ms.Viewport(ms.Content(h.Class("p-4"), items)),
				ms.Button(ms.DirectionStart),
				ms.Button(),
			),
		),
	)
}

func rows(anchors ...int) []any {
	out := make([]any, 0, 10)
	for row := 1; row <= 10; row++ {
		anchor := slices.Contains(anchors, row)
		var body h.Node
		if anchor {
			body = marker.Marker(marker.VariantSeparator, marker.Content("Turn ", h.Comment(""), fmt.Sprint(row)))
		} else {
			body = bubble.Bubble(bubble.VariantSecondary, bubble.Content("Row ", h.Comment(""), fmt.Sprint(row), h.Comment(""), " of the transcript."))
		}
		out = append(out, ms.Item(ms.MessageID(fmt.Sprintf("row-%d", row)), ms.ScrollAnchor(anchor), body))
	}
	return out
}

// Positions is the parity fixture message-scroller-positions.
func Positions() h.Node {
	box := func(children ...any) h.Node {
		return h.Div(append([]any{h.Class("h-64 rounded-lg border")}, children...)...)
	}
	return h.Div(h.Class("grid w-full max-w-4xl grid-cols-2 gap-4"),
		box(ms.Provider(ms.PositionStart,
			ms.MessageScroller(
				ms.Viewport(h.Aria("label", "Start"), ms.Content(h.Class("gap-3 p-3"), rows())),
				ms.Button(button.VariantOutline, button.SizeSm, ms.Behavior("auto"), lucide.ChevronsDown(), "Latest"),
			),
		)),
		box(ms.Provider(ms.PositionLastAnchor, ms.ScrollPreviousItemPeek(24), ms.ScrollMargin(4),
			ms.MessageScroller(
				ms.Viewport(ms.Content(h.Class("gap-3 p-3"), ms.SpacerClass("bg-muted"), rows(3, 6))),
				ms.Button(ms.DirectionStart),
				ms.Button(),
			),
		)),
		box(ms.Provider(ms.AutoScroll(true), ms.ScrollEdgeThreshold(16),
			ms.MessageScroller(h.Class("bg-muted/30"),
				ms.Viewport(h.TabIndex(-1), ms.PreserveScrollOnPrepend(false),
					ms.Content(h.Class("gap-3 p-3"), h.Role("list"),
						h.Div(h.Class("text-xs text-muted-foreground"), "A row without messageId."),
						rows(9),
					),
				),
				ms.Button(button.VariantDefault),
			),
		)),
		box(ms.Provider(
			ms.MessageScroller(
				ms.Viewport(ms.Content(h.Class("p-3"))),
				ms.Button(),
			),
		)),
	)
}

type dynRow struct {
	id         int
	user, long bool
}

func (r dynRow) text() string {
	if r.long {
		return fmt.Sprintf("Reply %d: a longer reply that wraps onto a few lines so new turns have to scroll into place.", r.id)
	}
	if r.user {
		return fmt.Sprintf("Question %d", r.id)
	}
	return fmt.Sprintf("Reply %d", r.id)
}

func (r dynRow) node() h.Node {
	variant, align := bubble.VariantMuted, bubble.AlignStart
	if r.user {
		variant, align = bubble.VariantDefault, bubble.AlignEnd
	}
	return ms.Item(ms.MessageID(fmt.Sprintf("r%d", r.id)), ms.ScrollAnchor(r.user),
		bubble.Bubble(variant, align, bubble.Content(r.text())),
	)
}

// Dynamic is the parity fixture message-scroller-dynamic; dynamic.js
// plays the fixture's state.
func Dynamic() h.Node {
	rows := make([]any, 0, 6)
	for i := 0; i < 6; i++ {
		rows = append(rows, dynRow{id: i + 1, user: i%3 == 0, long: i%2 == 1}.node())
	}
	action := func(name, label string, variant button.Variant) h.Node {
		return button.Button(button.SizeSm, variant, h.Data("example-ms-action", name), label)
	}
	return ms.Provider(ms.AutoScroll(true),
		h.Div(h.Class("flex w-full max-w-xl flex-col gap-3"), h.Data("example-ms-dynamic", ""),
			h.Div(h.Class("flex flex-wrap gap-2"),
				action("reply", "Add reply", button.VariantDefault),
				action("turn", "Add turn", button.VariantDefault),
				action("prepend", "Prepend", button.VariantDefault),
				action("jump2", "Jump to 2", button.VariantOutline),
				action("jump4", "Jump to 4", button.VariantOutline),
				action("start", "Start", button.VariantOutline),
				action("end", "End", button.VariantOutline),
			),
			h.P(h.Class("text-xs text-muted-foreground"), h.Data("example-ms-status", ""),
				"Current: ", h.Comment(""), "none", h.Comment(""), " · Visible: ", h.Comment(""), " · Scrollable: ", h.Comment(""), "-", h.Comment(""), " ", h.Comment(""), "-"),
			h.Div(h.Class("h-72 rounded-lg border"),
				ms.MessageScroller(
					ms.Viewport(ms.Content(h.Class("gap-4 p-4"), rows)),
					ms.Button(),
				),
			),
			h.Template(h.Data("parity-ignore", ""), h.Data("example-ms-rows", ""),
				dynRow{id: 0, user: true}.node(),
				dynRow{id: 0}.node(),
			),
		),
	)
}
