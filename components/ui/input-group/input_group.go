// Package inputgroup is the shadcn/ui InputGroup: an input or textarea with
// addons (icons, text, buttons, kbd) inside the same border.
//
//	inputgroup.InputGroup(
//		inputgroup.Input(h.Placeholder("Search...")),
//		inputgroup.Addon(lucide.Search()),
//		inputgroup.Addon(inputgroup.AlignInlineEnd, "12 results"),
//	)
//
// Clicking an addon (outside its buttons) focuses the group's input, like
// upstream (see input-group.js).
package inputgroup

import (
	_ "embed"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

//go:embed input-group.js
var script string

func init() { runtime.Module("input-group", script) }

// Align places an Addon.
type Align string

// Addon alignments.
const (
	AlignInlineStart Align = "inline-start"
	AlignInlineEnd   Align = "inline-end"
	AlignBlockStart  Align = "block-start"
	AlignBlockEnd    Align = "block-end"
)

// ButtonSize selects the size of a Button.
type ButtonSize string

// Button sizes.
const (
	ButtonSizeXs     ButtonSize = "xs"
	ButtonSizeSm     ButtonSize = "sm"
	ButtonSizeIconXs ButtonSize = "icon-xs"
	ButtonSizeIconSm ButtonSize = "icon-sm"
)

var groupClass = lib.Cn(
	"group/input-group relative flex w-full items-center rounded-md border border-input shadow-xs transition-[color,box-shadow] outline-none dark:bg-input/30",
	"h-9 min-w-0 has-[>textarea]:h-auto",
	"has-[>[data-align=inline-start]]:[&>input]:pl-2",
	"has-[>[data-align=inline-end]]:[&>input]:pr-2",
	"has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>[data-align=block-start]]:[&>input]:pb-3",
	"has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-end]]:[&>input]:pt-3",
	"has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-[3px] has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50",
	"has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-destructive/20 dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40",
)

var addonVariants = cva.New(
	"flex h-auto cursor-text items-center justify-center gap-2 py-1.5 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [&>kbd]:rounded-[calc(var(--radius)-5px)] [&>svg:not([class*='size-'])]:size-4",
	cva.Config{
		Variants: cva.Variants{
			"align": {
				"inline-start": "order-first pl-3 has-[>button]:ml-[-0.45rem] has-[>kbd]:ml-[-0.35rem]",
				"inline-end":   "order-last pr-3 has-[>button]:mr-[-0.45rem] has-[>kbd]:mr-[-0.35rem]",
				"block-start":  "order-first w-full justify-start px-3 pt-3 group-has-[>input]/input-group:pt-2.5 [.border-b]:pb-3",
				"block-end":    "order-last w-full justify-start px-3 pb-3 group-has-[>input]/input-group:pb-2.5 [.border-t]:pt-3",
			},
		},
		Defaults: cva.Selection{"align": "inline-start"},
	},
)

var buttonVariants = cva.New(
	"flex items-center gap-2 text-sm shadow-none",
	cva.Config{
		Variants: cva.Variants{
			"size": {
				"xs":      "h-6 gap-1 rounded-[calc(var(--radius)-5px)] px-2 has-[>svg]:px-2 [&>svg:not([class*='size-'])]:size-3.5",
				"sm":      "h-8 gap-1.5 rounded-md px-2.5 has-[>svg]:px-2.5",
				"icon-xs": "size-6 rounded-[calc(var(--radius)-5px)] p-0 has-[>svg]:p-0",
				"icon-sm": "size-8 p-0 has-[>svg]:p-0",
			},
		},
		Defaults: cva.Selection{"size": "xs"},
	},
)

// AddonVariants returns the classes of an Addon, like
// inputGroupAddonVariants() upstream.
func AddonVariants(a Align, className ...string) string {
	return lib.Cn(addonVariants.Classes(cva.Selection{"align": string(a)}, className...))
}

// ButtonVariants returns the classes InputGroupButton adds to a Button.
func ButtonVariants(s ButtonSize, className ...string) string {
	return lib.Cn(buttonVariants.Classes(cva.Selection{"size": string(s)}, className...))
}

var (
	groupBase = h.Div(h.Data("slot", "input-group"), h.Role("group"), h.Class(groupClass))
	textBase  = h.Span(h.Class(lib.Cn("flex items-center gap-2 text-sm text-muted-foreground [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4")))
	inputCls  = "flex-1 rounded-none border-0 bg-transparent shadow-none focus-visible:ring-0 dark:bg-transparent"
	areaCls   = "flex-1 resize-none rounded-none border-0 bg-transparent py-3 shadow-none focus-visible:ring-0 dark:bg-transparent"
)

// InputGroup is the group's bordered container.
func InputGroup(children ...any) *h.Element {
	return h.Spread(groupBase, h.Div(children...), lib.Cn)
}

// Text is muted inline text inside an Addon (upstream InputGroupText).
func Text(children ...any) *h.Element {
	return h.Spread(textBase, h.Span(children...), lib.Cn)
}

// Input is the group's <input> (upstream InputGroupInput, an Input).
func Input(children ...any) *h.Element {
	props := h.Input(children...)
	props.SetClass(lib.Cn(inputCls, props.ClassName()))
	return input.Input(h.Data("slot", "input-group-control"), spreadProp(props))
}

// Textarea is the group's <textarea> (upstream InputGroupTextarea).
func Textarea(children ...any) *h.Element {
	props := h.Textarea(children...)
	props.SetClass(lib.Cn(areaCls, props.ClassName()))
	return textarea.Textarea(h.Data("slot", "input-group-control"), spreadProp(props))
}

// Builder is returned by Addon and Button.
type Builder struct {
	h.Builder[Builder]
	isButton bool
	align    Align
	size     ButtonSize
	variant  button.Variant
}

// Addon holds icons, text, buttons or kbd next to the control. Children
// may include an Align (inline-start by default).
func Addon(children ...any) *Builder {
	b := &Builder{align: AlignInlineStart}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Button is a compact button for an Addon (upstream InputGroupButton): a
// ghost button.Button of size xs by default. Children may include a
// ButtonSize and a button.Variant.
func Button(children ...any) *Builder {
	b := &Builder{isButton: true, size: ButtonSizeXs, variant: button.VariantGhost}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Align (Addon), ButtonSize and button.Variant
// (Button).
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Align:
		if b.isButton {
			return false
		}
		b.Align(v)
	case ButtonSize:
		if !b.isButton {
			return false
		}
		b.Size(v)
	case button.Variant:
		if !b.isButton {
			return false
		}
		b.Variant(v)
	default:
		return false
	}
	return true
}

// Align sets an Addon's alignment.
func (b *Builder) Align(a Align) *Builder {
	if a == "" {
		a = AlignInlineStart
	}
	b.align = a
	return b
}

// Size sets a Button's size.
func (b *Builder) Size(s ButtonSize) *Builder {
	if s == "" {
		s = ButtonSizeXs
	}
	b.size = s
	return b
}

// Variant sets a Button's variant.
func (b *Builder) Variant(v button.Variant) *Builder {
	if v == "" {
		v = button.VariantGhost
	}
	b.variant = v
	return b
}

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	if b.isButton {
		props.SetClass(ButtonVariants(b.size, className))
		return button.Button(b.variant, h.Type("button"), h.Data("size", string(b.size)), spreadProp(props)).Element()
	}
	props.SetClass("")
	base := h.Div(h.Role("group"), h.Data("slot", "input-group-addon"), h.Data("align", string(b.align)), h.Class(lib.Cn(AddonVariants(b.align), className)))
	return h.Spread(base, props, lib.Cn)
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// spreadProp applies props' attributes, classes and children to the
// element it is passed to, like {...props}.
func spreadProp(props *h.Element) h.Prop {
	return func(e *h.Element) { *e = *h.Spread(e, props, nil) }
}
