// Package message is the shadcn/ui Message: one entry of a conversation,
// with an avatar, a header, content (bubbles, attachments) and a footer.
//
//	message.Message(message.AlignEnd,
//		message.Avatar(avatar.Avatar(avatar.Fallback("ME"))),
//		message.Content(bubble.Bubble(bubble.Content("Hi!")), message.Footer("Delivered")),
//	)
package message

import (
	"io"

	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

// Align puts the message on the start (others) or the end (you) side.
type Align string

// Alignments.
const (
	AlignStart Align = "start"
	AlignEnd   Align = "end"
)

const messageClass = "group/message relative flex w-full min-w-0 gap-2 text-sm data-[align=end]:flex-row-reverse"

var (
	groupBase   = h.Div(h.Data("slot", "message-group"), h.Class(lib.Cn("flex min-w-0 flex-col gap-2")))
	avatarBase  = h.Div(h.Data("slot", "message-avatar"), h.Class(lib.Cn("flex w-fit min-w-8 shrink-0 items-center justify-center self-end overflow-hidden rounded-full bg-muted group-has-data-[slot=message-footer]/message:-translate-y-8")))
	contentBase = h.Div(h.Data("slot", "message-content"), h.Class(lib.Cn("flex w-full min-w-0 flex-col gap-2.5 wrap-break-word group-data-[align=end]/message:*:data-slot:self-end")))
	headerBase  = h.Div(h.Data("slot", "message-header"), h.Class(lib.Cn("flex max-w-full min-w-0 items-center px-3 text-xs font-medium text-muted-foreground group-has-data-[variant=ghost]/message:px-0")))
	footerBase  = h.Div(h.Data("slot", "message-footer"), h.Class(lib.Cn("flex max-w-full min-w-0 items-center px-3 text-xs font-medium text-muted-foreground group-has-data-[variant=ghost]/message:px-0 group-data-[align=end]/message:justify-end")))
)

func spread(base *h.Element, children []any) *h.Element {
	return h.Spread(base, h.Div(children...), lib.Cn)
}

// Group stacks messages (upstream MessageGroup).
func Group(children ...any) *h.Element { return spread(groupBase, children) }

// Avatar holds the author's avatar, aligned with the message's last line
// (lifted above the footer when there is one).
func Avatar(children ...any) *h.Element { return spread(avatarBase, children) }

// Content holds the message's bubbles, attachments, header and footer.
func Content(children ...any) *h.Element { return spread(contentBase, children) }

// Header is a small line above the content (the author's name).
func Header(children ...any) *h.Element { return spread(headerBase, children) }

// Footer is a small line below the content (status, actions).
func Footer(children ...any) *h.Element { return spread(footerBase, children) }

// Builder is returned by Message.
type Builder struct {
	h.Builder[Builder]
	align Align
}

// Message returns a message row. Children may include an Align (default
// start).
func Message(children ...any) *Builder {
	b := &Builder{align: AlignStart}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Align.
func (b *Builder) AcceptChild(c any) bool {
	if a, ok := c.(Align); ok {
		b.Align(a)
		return true
	}
	return false
}

// Align sets the side. An empty value selects the start.
func (b *Builder) Align(a Align) *Builder {
	if a == "" {
		a = AlignStart
	}
	b.align = a
	return b
}

// Element resolves the message.
func (b *Builder) Element() *h.Element {
	base := h.Div(h.Data("slot", "message"), h.Data("align", string(b.align)), h.Class(lib.Cn(messageClass)))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the message.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }
