// Package bubble is the shadcn/ui Bubble: a chat bubble, with optional
// reactions, grouped in a BubbleGroup.
//
//	bubble.Group(
//		bubble.Bubble(bubble.VariantMuted, bubble.Content("Hello!")),
//		bubble.Bubble(bubble.AlignEnd, bubble.Content("Hi"),
//			bubble.Reactions(h.Span("👍"))),
//	)
package bubble

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/styled"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Variant selects the bubble's look.
type Variant string

// Bubble variants.
const (
	VariantDefault     Variant = "default"
	VariantSecondary   Variant = "secondary"
	VariantMuted       Variant = "muted"
	VariantTinted      Variant = "tinted"
	VariantOutline     Variant = "outline"
	VariantGhost       Variant = "ghost"
	VariantDestructive Variant = "destructive"
)

// VariantValues lists every variant, in the order shadcn declares them.
var VariantValues = []Variant{VariantDefault, VariantSecondary, VariantMuted, VariantTinted, VariantOutline, VariantGhost, VariantDestructive}

// Align places a Bubble (or Reactions) at the start or the end.
type Align string

// Alignments.
const (
	AlignStart Align = "start"
	AlignEnd   Align = "end"
)

// Side places Reactions above or below the bubble.
type Side string

// Sides.
const (
	SideTop    Side = "top"
	SideBottom Side = "bottom"
)

var bubbleVariants = cva.New(
	"group/bubble relative flex w-fit max-w-[80%] min-w-0 flex-col gap-1 group-data-[align=end]/message:self-end data-[align=end]:self-end data-[variant=ghost]:max-w-full",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default":     "*:data-[slot=bubble-content]:bg-primary *:data-[slot=bubble-content]:text-primary-foreground [&>[data-slot=bubble-content]:is(button,a):hover]:bg-primary/80",
				"secondary":   "*:data-[slot=bubble-content]:bg-secondary *:data-[slot=bubble-content]:text-secondary-foreground [&>[data-slot=bubble-content]:is(button,a):hover]:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)]",
				"muted":       "*:data-[slot=bubble-content]:bg-muted [&>[data-slot=bubble-content]:is(button,a):hover]:bg-[color-mix(in_oklch,var(--muted),var(--foreground)_5%)]",
				"tinted":      "*:data-[slot=bubble-content]:bg-[oklch(from_var(--primary)_0.93_calc(c*0.4)_h)] *:data-[slot=bubble-content]:text-foreground dark:*:data-[slot=bubble-content]:bg-[oklch(from_var(--primary)_0.3_calc(c*0.4)_h)] [&>[data-slot=bubble-content]:is(button,a):hover]:bg-[oklch(from_var(--primary)_0.88_calc(c*0.5)_h)] dark:[&>[data-slot=bubble-content]:is(button,a):hover]:bg-[oklch(from_var(--primary)_0.35_calc(c*0.5)_h)]",
				"outline":     "*:data-[slot=bubble-content]:border-border *:data-[slot=bubble-content]:bg-background [&>[data-slot=bubble-content]:is(button,a):hover]:bg-muted [&>[data-slot=bubble-content]:is(button,a):hover]:text-foreground dark:[&>[data-slot=bubble-content]:is(button,a):hover]:bg-input/30",
				"ghost":       "border-none *:data-[slot=bubble-content]:rounded-none *:data-[slot=bubble-content]:bg-transparent *:data-[slot=bubble-content]:p-0 [&>[data-slot=bubble-content]:is(button,a):hover]:bg-muted [&>[data-slot=bubble-content]:is(button,a):hover]:text-foreground dark:[&>[data-slot=bubble-content]:is(button,a):hover]:bg-muted/50",
				"destructive": "*:data-[slot=bubble-content]:bg-destructive/10 *:data-[slot=bubble-content]:text-destructive dark:*:data-[slot=bubble-content]:bg-destructive/20 [&>[data-slot=bubble-content]:is(button,a):hover]:bg-destructive/20 dark:[&>[data-slot=bubble-content]:is(button,a):hover]:bg-destructive/30",
			},
		},
		Defaults: cva.Selection{"variant": "default"},
		Order:    []string{"variant"},
	},
)

var reactionsVariants = cva.New(
	"absolute z-10 flex w-fit shrink-0 items-center justify-center gap-1 rounded-full bg-muted px-1.5 py-0.5 text-sm ring-3 ring-card has-[button]:p-0",
	cva.Config{
		Variants: cva.Variants{
			"side":  {"top": "top-0 -translate-y-3/4", "bottom": "bottom-0 translate-y-3/4"},
			"align": {"start": "left-3", "end": "right-3"},
		},
		Defaults: cva.Selection{"side": "bottom", "align": "end"},
		Order:    []string{"side", "align"},
	},
)

const contentClass = "w-fit max-w-full min-w-0 overflow-hidden rounded-xl border border-transparent px-3 py-2 text-sm leading-relaxed wrap-break-word group-data-[align=end]/bubble:self-end [button]:text-left [button,a]:transition-colors [button,a]:outline-none [button,a]:focus-visible:border-ring [button,a]:focus-visible:ring-3 [button,a]:focus-visible:ring-ring/50"

// Variants returns the classes of a bubble with the given variant plus
// className, like cn(bubbleVariants({ variant }), className).
func Variants(v Variant, className ...string) string {
	return lib.Cn(append([]string{bubbleVariants.Classes(cva.Selection{"variant": string(v)})}, className...)...)
}

// ReactionsVariants returns the classes of Reactions, like
// cn(bubbleReactionsVariants({ side, align }), className).
func ReactionsVariants(side Side, align Align, className ...string) string {
	return lib.Cn(append([]string{reactionsVariants.Classes(cva.Selection{"side": string(side), "align": string(align)})}, className...)...)
}

var groupBase = h.Div(h.Data("slot", "bubble-group"), h.Class(lib.Cn("flex min-w-0 flex-col gap-2")))

// Group stacks bubbles (upstream BubbleGroup).
func Group(children ...any) *styled.Node {
	return styled.Spread(groupBase, h.Div(children...), styled.Slot("div", "bubble-group", sGroup))
}

// Builder is returned by Bubble.
type Builder struct {
	h.Builder[Builder]
	variant Variant
	align   Align
}

// Bubble returns a bubble. Children may include a Variant and an Align.
func Bubble(children ...any) *Builder {
	b := &Builder{variant: VariantDefault, align: AlignStart}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Variant and Align.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Variant:
		b.Variant(v)
	case Align:
		b.Align(v)
	default:
		return false
	}
	return true
}

// Variant sets the variant. An empty value selects the default.
func (b *Builder) Variant(v Variant) *Builder {
	if v == "" {
		v = VariantDefault
	}
	b.variant = v
	return b
}

// Align sets the alignment. An empty value selects the start.
func (b *Builder) Align(a Align) *Builder {
	if a == "" {
		a = AlignStart
	}
	b.align = a
	return b
}

// Element resolves the bubble.
func (b *Builder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the bubble as it renders inside w (in the style of an
// enclosing style.Scope).
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Div(
		h.Data("slot", "bubble"),
		h.Data("variant", string(b.variant)),
		h.Data("align", string(b.align)),
		h.Class(bubbleClass(styleOf(w), b.variant, className)),
	)
	return h.Spread(base, props, nil)
}

// Render renders the bubble.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	asChild bool
}

// Content is the bubble's body. With lib.AsChild it renders its single
// element child (a button or a link) instead of a <div>.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if _, ok := c.(lib.AsChildFlag); ok {
		b.asChild = true
		return true
	}
	return false
}

// AsChild makes the content render its single element child.
func (b *ContentBuilder) AsChild() *ContentBuilder { b.asChild = true; return b }

// Element resolves the content.
func (b *ContentBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the content as it renders inside w.
func (b *ContentBuilder) ElementIn(w io.Writer) *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Div(h.Data("slot", "bubble-content"), h.Class(contentClassFor(styleOf(w), className)))
	el := h.Spread(base, props, nil)
	if !b.asChild {
		return el
	}
	return slotChild(w, "bubble: Content", el)
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ReactionsBuilder is returned by Reactions.
type ReactionsBuilder struct {
	h.Builder[ReactionsBuilder]
	side  Side
	align Align
}

// Reactions pins reactions (emoji, a button) to the bubble's edge. Children
// may include a Side (default bottom) and an Align (default end).
func Reactions(children ...any) *ReactionsBuilder {
	b := &ReactionsBuilder{side: SideBottom, align: AlignEnd}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Side and Align.
func (b *ReactionsBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Side:
		if v == "" {
			v = SideBottom
		}
		b.side = v
	case Align:
		if v == "" {
			v = AlignEnd
		}
		b.align = v
	default:
		return false
	}
	return true
}

// Element resolves the reactions.
func (b *ReactionsBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the reactions as it renders inside w.
func (b *ReactionsBuilder) ElementIn(w io.Writer) *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Div(
		h.Data("slot", "bubble-reactions"),
		h.Data("align", string(b.align)),
		h.Data("side", string(b.side)),
		h.Class(reactionsClass(styleOf(w), b.side, b.align, className)),
	)
	return h.Spread(base, props, nil)
}

// Render renders the reactions.
func (b *ReactionsBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// slotChild renders el's single element child with el's attributes merged
// in, like Radix Slot (classes concatenated, slot first).
func slotChild(w io.Writer, what string, el *h.Element) *h.Element {
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("%s: AsChild needs exactly one element child, got %d children", what, len(kids)))
	}
	out, err := h.SlotIn(w, el.Clone().SetChildren(), kids[0], nil)
	if err != nil {
		return h.Errored(fmt.Errorf("%s: AsChild: %w", what, err))
	}
	return out
}
