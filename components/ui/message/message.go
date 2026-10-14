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
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/styled"
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

func spread(base *h.Element, children []any, site *style.Site) *styled.Node {
	slot, _ := base.GetAttr("data-slot")
	return styled.Spread(base, h.Div(children...), styled.Slot("div", slot, site))
}

// Group stacks messages (upstream MessageGroup).
func Group(children ...any) *styled.Node { return spread(groupBase, children, sGroup) }

// Avatar holds the author's avatar, aligned with the message's last line
// (lifted above the footer when there is one).
func Avatar(children ...any) *styled.Node { return spread(avatarBase, children, sAvatar) }

// Content holds the message's bubbles, attachments, header and footer.
func Content(children ...any) *styled.Node { return spread(contentBase, children, sContent) }

// Header is a small line above the content (the author's name).
func Header(children ...any) *styled.Node { return spread(headerBase, children, sHeader) }

// Footer is a small line below the content (status, actions).
func Footer(children ...any) *styled.Node { return spread(footerBase, children, sFooter) }

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
func (b *Builder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the message as it renders inside w (in the style of
// an enclosing style.Scope).
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	class := lib.Cn(messageClass)
	if w != nil {
		if st, ok := style.FromWriter(w); ok {
			class = sMessage.Class(st, "")
		}
	}
	base := h.Div(h.Data("slot", "message"), h.Data("align", string(b.align)), h.Class(class))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the message.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }
