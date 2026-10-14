// Package messagescroller is the shadcn/ui MessageScroller: a scroll
// container for chat transcripts that opens at the end (or the start, or
// the last turn), follows new messages, anchors a new turn near the top and
// shows scroll-to-end / scroll-to-start buttons while there is more to see.
//
//	messagescroller.Provider(messagescroller.AutoScroll(true),
//		messagescroller.MessageScroller(
//			messagescroller.Viewport(
//				messagescroller.Content(
//					messagescroller.Item(messagescroller.MessageID("m1"), messagescroller.ScrollAnchor(true), ...),
//					messagescroller.Item(messagescroller.MessageID("m2"), ...),
//				),
//			),
//			messagescroller.Button(),
//		),
//	)
//
// The behavior (message-scroller.js) is a port of the headless
// @shadcn/react/message-scroller primitive; scripts reach it with
// getMessageScroller(el) (scrollToEnd, scrollToStart, scrollToMessage,
// subscribeVisibility) and observe "shadcn:scrollable-change".
package messagescroller

import (
	_ "embed"
	"encoding/json"
	"io"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed message-scroller.js
var script string

func init() { runtime.Module("message-scroller", script) }

// Provider options.
type (
	// AutoScroll follows new content at the bottom while the viewport is at
	// the end.
	AutoScroll bool
	// DefaultScrollPosition is where the transcript opens.
	DefaultScrollPosition string
	// ScrollEdgeThreshold is the distance from an edge (px) that still counts
	// as at the edge (default 8).
	ScrollEdgeThreshold float64
	// ScrollPreviousItemPeek is how much of the previous item (px) stays
	// visible above a newly anchored item (default 64).
	ScrollPreviousItemPeek float64
	// ScrollMargin is the margin (px) on the aligned edge of scroll commands
	// (default 0).
	ScrollMargin float64
)

// Opening positions.
const (
	PositionStart      DefaultScrollPosition = "start"
	PositionEnd        DefaultScrollPosition = "end"
	PositionLastAnchor DefaultScrollPosition = "last-anchor"
)

// Viewport, Content, Item and Button options.
type (
	// PreserveScrollOnPrepend keeps the first visible row in place when rows
	// are prepended (Viewport; default true).
	PreserveScrollOnPrepend bool
	// SpacerClass is the class of the Content's tail spacer.
	SpacerClass string
	// MessageID identifies an Item for scrollToMessage and visibility.
	MessageID string
	// ScrollAnchor marks an Item as a turn boundary: appended, it is
	// scrolled near the top.
	ScrollAnchor bool
	// Direction is the transcript edge a Button scrolls toward.
	Direction string
	// Behavior is a Button's native scroll behavior (default "smooth").
	Behavior string
)

// Button directions.
const (
	DirectionEnd   Direction = "end"
	DirectionStart Direction = "start"
)

type config struct {
	id      string
	options map[string]any
}

type scoped struct{ cfg *config }

func (s *scoped) setConfig(c *config) { s.cfg = c }

func (s *scoped) config() *config {
	if s.cfg == nil {
		s.cfg = &config{id: uid.New("message-scroller"), options: map[string]any{}}
	}
	return s.cfg
}

type participant interface{ setConfig(*config) }

func (s *scoped) mark(el *h.Element, part string) {
	el.Attr("data-message-scroller", s.config().id)
	el.Attr("data-message-scroller-part", part)
}

// --- Provider ----------------------------------------------------------------------------------

// ProviderBuilder is returned by Provider.
type ProviderBuilder struct {
	children []any
	cfg      *config
	position DefaultScrollPosition
}

// Provider holds the scroller's options and renders no element. Children
// may include AutoScroll, DefaultScrollPosition (default end),
// ScrollEdgeThreshold, ScrollPreviousItemPeek and ScrollMargin.
func Provider(children ...any) *ProviderBuilder {
	b := &ProviderBuilder{cfg: &config{id: uid.New("message-scroller"), options: map[string]any{}}, position: PositionEnd}
	for _, c := range children {
		switch v := c.(type) {
		case AutoScroll:
			b.cfg.options["autoScroll"] = bool(v)
		case DefaultScrollPosition:
			if v != "" {
				b.position = v
			}
		case ScrollEdgeThreshold:
			b.cfg.options["scrollEdgeThreshold"] = float64(v)
		case ScrollPreviousItemPeek:
			b.cfg.options["scrollPreviousItemPeek"] = float64(v)
		case ScrollMargin:
			b.cfg.options["scrollMargin"] = float64(v)
		default:
			b.children = append(b.children, c)
		}
	}
	b.cfg.options["defaultScrollPosition"] = string(b.position)
	return b
}

// pending reports whether the opening position hides the viewport until it
// is applied (data-pending-scroll).
func (c *config) pending() bool {
	p, _ := c.options["defaultScrollPosition"].(string)
	return p == "" || p == string(PositionEnd) || p == string(PositionLastAnchor)
}

// Render renders the children with the provider's options shared.
func (b *ProviderBuilder) Render(w io.Writer) error {
	tree.Walk(b.children, func(n any) bool {
		switch c := n.(type) {
		case *ProviderBuilder:
			return false
		case participant:
			c.setConfig(b.cfg)
		}
		return true
	})
	return h.Group(b.children...).Render(w)
}

// --- MessageScroller (root) --------------------------------------------------------------------

// RootBuilder is returned by MessageScroller.
type RootBuilder struct {
	h.Builder[RootBuilder]
	scoped
}

// MessageScroller is the frame laying out the viewport and the buttons.
func MessageScroller(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the root.
func (b *RootBuilder) Element() *h.Element {
	base := h.Div(h.Data("slot", "message-scroller"), h.Class(lib.Cn("group/message-scroller relative flex size-full min-h-0 flex-col overflow-hidden")))
	el := h.Spread(base, b.Elem(), lib.Cn)
	if b.config().pending() {
		el.Attr("data-pending-scroll", "")
	}
	b.mark(el, "root")
	return el
}

// Render renders the root.
func (b *RootBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Viewport ----------------------------------------------------------------------------------

// ViewportBuilder is returned by Viewport.
type ViewportBuilder struct {
	h.Builder[ViewportBuilder]
	scoped
	preserve bool
}

// Viewport is the scrolling region (role="region", aria-label="Messages",
// tabindex 0 unless overridden). Children may include
// PreserveScrollOnPrepend.
func Viewport(children ...any) *ViewportBuilder {
	b := &ViewportBuilder{preserve: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes PreserveScrollOnPrepend.
func (b *ViewportBuilder) AcceptChild(c any) bool {
	if v, ok := c.(PreserveScrollOnPrepend); ok {
		b.preserve = bool(v)
		return true
	}
	return false
}

const viewportClass = "size-full min-h-0 min-w-0 scroll-fade-b scrollbar-thin scrollbar-gutter-stable overflow-y-auto overscroll-contain contain-content data-autoscrolling:scrollbar-none data-pending-scroll:invisible"

// Element resolves the viewport.
func (b *ViewportBuilder) Element() *h.Element {
	base := h.Div(h.Role("region"), h.Aria("label", "Messages"), h.TabIndex(0), h.Data("slot", "message-scroller-viewport"), h.Class(lib.Cn(viewportClass)))
	el := h.Spread(base, b.Elem(), lib.Cn)
	cfg := b.config()
	if cfg.pending() {
		el.Attr("data-pending-scroll", "")
	}
	opts := map[string]any{"preserveScrollOnPrepend": b.preserve}
	for k, v := range cfg.options {
		opts[k] = v
	}
	js, _ := json.Marshal(opts)
	el.Attr("data-message-scroller-options", string(js))
	b.mark(el, "viewport")
	return el
}

// Render renders the viewport.
func (b *ViewportBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Content -----------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	scoped
	spacerClass string
}

// Content is the list of rows (role="log", aria-relevant="additions"),
// followed by the hidden tail spacer used when anchoring a turn near the
// top. Children may include SpacerClass.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes SpacerClass.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if v, ok := c.(SpacerClass); ok {
		b.spacerClass = string(v)
		return true
	}
	return false
}

// Element resolves the content.
func (b *ContentBuilder) Element() *h.Element {
	base := h.Div(h.Role("log"), h.Aria("relevant", "additions"), h.Data("slot", "message-scroller-content"), h.Class(lib.Cn("flex h-max min-h-full flex-col gap-8")))
	el := h.Spread(base, b.Elem(), lib.Cn)
	spacer := h.Div(h.Aria("hidden", "true"), h.Data("message-scroller-spacer", ""), h.Hidden())
	if b.spacerClass != "" {
		spacer.Class(b.spacerClass)
	}
	el.SetChildren(append(append([]any{}, el.ChildNodes()...), spacer)...)
	b.mark(el, "content")
	return el
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Item --------------------------------------------------------------------------------------

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	messageID string
	anchor    bool
}

// Item is one row: a message, a marker, a typing indicator, ... Children
// may include MessageID and ScrollAnchor.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes MessageID and ScrollAnchor.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case MessageID:
		b.messageID = string(v)
	case ScrollAnchor:
		b.anchor = bool(v)
	default:
		return false
	}
	return true
}

// Element resolves the item.
func (b *ItemBuilder) Element() *h.Element {
	base := h.Div(h.Data("slot", "message-scroller-item"), h.Class(lib.Cn("min-w-0 shrink-0 [contain-intrinsic-size:auto_10rem] [content-visibility:auto]")))
	if b.messageID != "" {
		base.Attr("data-message-id", b.messageID)
	}
	base.Attr("data-scroll-anchor", strconv.FormatBool(b.anchor))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Button ------------------------------------------------------------------------------------

// ButtonBuilder is returned by Button.
type ButtonBuilder struct {
	h.Builder[ButtonBuilder]
	scoped
	direction Direction
	behavior  Behavior
	variant   button.Variant
	size      button.Size
}

// Button scrolls to the end (or, with DirectionStart, the start) of the
// transcript; it is hidden while there is nothing to scroll to. Children
// may include a Direction, a Behavior, a button.Variant (default
// secondary) and a button.Size (default icon-sm); without other children
// it shows an arrow and a screen-reader label.
func Button(children ...any) *ButtonBuilder {
	b := &ButtonBuilder{direction: DirectionEnd, behavior: "smooth", variant: button.VariantSecondary, size: button.SizeIconSm}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Direction, Behavior, button.Variant and button.Size.
func (b *ButtonBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Direction:
		if v == "" {
			v = DirectionEnd
		}
		b.direction = v
	case Behavior:
		if v == "" {
			v = "smooth"
		}
		b.behavior = v
	case button.Variant:
		b.variant = v
	case button.Size:
		b.size = v
	default:
		return false
	}
	return true
}

const buttonClass = "absolute inset-s-1/2 -translate-x-1/2 border-border bg-background text-foreground transition-[translate,scale,opacity] duration-200 hover:bg-muted hover:text-foreground data-[active=false]:pointer-events-none data-[active=false]:scale-95 data-[active=false]:opacity-0 data-[active=false]:duration-400 data-[active=false]:ease-[cubic-bezier(0.7,0,0.84,0)] data-[active=true]:translate-y-0 data-[active=true]:scale-100 data-[active=true]:opacity-100 data-[active=true]:ease-[cubic-bezier(0.23,1,0.32,1)] data-[direction=end]:bottom-4 data-[direction=end]:data-[active=false]:translate-y-full data-[direction=start]:top-4 data-[direction=start]:data-[active=false]:-translate-y-full rtl:translate-x-1/2 data-[direction=start]:[&_svg]:rotate-180"

// Element resolves the button.
func (b *ButtonBuilder) Element() *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	kids := props.ChildNodes()
	if len(kids) == 0 {
		label := "Scroll to end"
		if b.direction == DirectionStart {
			label = "Scroll to start"
		}
		kids = []any{lucide.ArrowDown(), h.Span(h.Class("sr-only"), label)}
	}
	userTab, hasTab := props.GetAttr("tabindex")
	props.SetChildren()
	attrs := []any{
		b.variant, b.size,
		h.Data("slot", "message-scroller-button"),
		h.Data("direction", string(b.direction)),
		h.Data("active", "false"),
		h.Type("button"),
		h.Bool("inert"),
		h.Class(lib.Cn(buttonClass, className)),
	}
	el := button.Button(append(attrs, kids...)...).Element()
	el = h.Spread(el, props, nil)
	el.Attr("tabindex", "-1")
	if hasTab {
		el.Attr("data-message-scroller-tabindex", userTab)
	}
	if b.behavior != "smooth" {
		el.Attr("data-message-scroller-behavior", string(b.behavior))
	}
	b.mark(el, "button")
	return el
}

// Render renders the button.
func (b *ButtonBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
