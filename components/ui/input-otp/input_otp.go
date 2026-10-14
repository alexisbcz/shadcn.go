// Package inputotp is the shadcn/ui InputOTP: an accessible one-time
// password input split into slots, a port of the input-otp library
// (v1.4.2): one transparent <input> over the visual slots, with the same
// selection handling, fake caret and password-manager badge handling.
//
//	inputotp.InputOTP(inputotp.MaxLength(6),
//		inputotp.Group(inputotp.Slot(0), inputotp.Slot(1), inputotp.Slot(2)),
//		inputotp.Separator(),
//		inputotp.Group(inputotp.Slot(3), inputotp.Slot(4), inputotp.Slot(5)),
//	)
//
// Attributes given to InputOTP (h.Class, h.Name, h.Value, h.Disabled,
// h.Aria, ...) go to the <input>, like upstream; ContainerClass styles the
// container. Slots are matched to indexes in document order (Slot(0) must
// come first, then Slot(1), ...).
//
// Events (bubbling from the input): "shadcn:value-change" with
// detail.value after every change and "shadcn:complete" with detail.value
// when the last slot gets filled.
package inputotp

import (
	_ "embed"
	"io"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed input-otp.js
var script string

func init() { runtime.Module("input-otp", script) }

// The patterns input-otp exports.
const (
	RegexpOnlyDigits         = `^\d+$`
	RegexpOnlyChars          = `^[a-zA-Z]+$`
	RegexpOnlyDigitsAndChars = `^[a-zA-Z0-9]+$`
)

// MaxLength is the number of slots (required).
type MaxLength int

// Pattern is a regular expression every value must match (both a JS RegExp
// and an HTML pattern), such as RegexpOnlyDigits.
type Pattern string

// ContainerClass adds classes to the container <div>.
type ContainerClass string

// InputMode overrides the input's inputmode (default "numeric").
type InputMode string

// TextAlign sets where new characters are typed: "left" (default),
// "center" or "right".
type TextAlign string

// Builder is returned by InputOTP.
type Builder struct {
	h.Builder[Builder]
	maxLength      int
	pattern        Pattern
	containerClass string
	inputMode      string
	textAlign      string
	st             style.Style // the enclosing style.Scope, "" for new-york-v4
}

// The classes of registry/bases/radix/ui/input-otp.tsx, used inside a
// style.Scope (see package style).
var (
	file        = style.NewFile("ui/input-otp.tsx")
	stContainer = file.Cn("cn-input-otp flex items-center has-disabled:opacity-50", style.Arg)
	stInput     = file.Cn("cn-input-otp-input disabled:cursor-not-allowed", style.Arg)
	stGroup     = file.Cn("cn-input-otp-group flex items-center", style.Arg)
	stSlot      = file.Cn("cn-input-otp-slot relative flex items-center justify-center data-[active=true]:z-10", style.Arg)
	stSeparator = file.Attr("cn-input-otp-separator flex items-center")
)

// InputOTP returns the input. Children may include MaxLength, Pattern,
// ContainerClass, InputMode and TextAlign
// values; attributes apply to the <input>, nodes render inside the
// container.
func InputOTP(children ...any) *Builder {
	b := &Builder{inputMode: "numeric", textAlign: "left"}
	b.Builder = h.NewBuilder(b, "input")
	return b.Children(children...)
}

// AcceptChild consumes the typed options.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case MaxLength:
		b.maxLength = int(v)
	case Pattern:
		b.pattern = v
	case ContainerClass:
		b.containerClass = lib.Cn(b.containerClass, string(v))
	case InputMode:
		b.inputMode = string(v)
	case TextAlign:
		b.textAlign = string(v)
	default:
		return false
	}
	return true
}

// Builder-style setters.
func (b *Builder) MaxLength(n int) *Builder         { b.maxLength = n; return b }
func (b *Builder) Pattern(p string) *Builder        { b.pattern = Pattern(p); return b }
func (b *Builder) ContainerClass(c string) *Builder { return b.Children(ContainerClass(c)) }
func (b *Builder) Name(v string) *Builder           { b.Elem().Name(v); return b }
func (b *Builder) Value(v string) *Builder          { b.Elem().Value(v); return b }
func (b *Builder) Disabled(on bool) *Builder        { b.Elem().Disabled(on); return b }

const inputStyle = "position: absolute; inset: 0px; width: 100%; height: 100%; display: flex; opacity: 1; color: transparent; pointer-events: all; background: transparent; caret-color: transparent; border: 0px solid transparent; outline: transparent solid 0px; box-shadow: none; line-height: 1; letter-spacing: -0.5em; font-size: var(--root-height); font-family: monospace; font-variant-numeric: tabular-nums;"

// Element resolves the component. It renders, like input-otp, an (empty)
// <noscript> followed by the container, so use it as a fragment.
func (b *Builder) Element() *h.Element {
	props := b.Elem()
	value, _ := props.GetAttr("value")
	runes := []rune(value)
	if b.maxLength > 0 && len(runes) > b.maxLength {
		runes = runes[:b.maxLength]
	}
	value = string(runes)
	_, disabled := props.GetAttr("disabled")

	// The user's attributes, without the children (they go to the container).
	inputProps := props.Clone().SetChildren()
	input := h.Spread(h.Input(h.AutoComplete("one-time-code")), inputProps, nil)
	input.Data("slot", "input-otp")
	if b.st != "" {
		input.SetClass(lib.Cn(stInput.Class(b.st, ""), props.ClassName()))
		if _, ok := input.GetAttr("spellcheck"); !ok {
			input.Attr("spellcheck", "false")
		}
	} else {
		input.SetClass(lib.Cn("disabled:cursor-not-allowed", props.ClassName()))
	}
	input.Data("input-otp", "true")
	if value == "" {
		input.Data("input-otp-placeholder-shown", "true")
	}
	input.Attr("inputmode", b.inputMode)
	if b.pattern != "" {
		input.Attr("pattern", string(b.pattern))
	}
	input.Attr("maxlength", strconv.Itoa(b.maxLength))
	input.Value(value)
	input.Style(inputStyle + " text-align: " + b.textAlign + ";")

	cursor := "text"
	if disabled {
		cursor = "default"
	}
	kids := props.ChildNodes()
	fill(kids, &slotContext{value: runes})
	containerClass := lib.Cn("flex items-center gap-2 has-disabled:opacity-50", b.containerClass)
	if b.st != "" {
		// Upstream's transform leaves containerClassName alone: the token
		// stays and the style's gap never applies.
		containerClass = lib.Cn("cn-input-otp flex items-center has-disabled:opacity-50", b.containerClass)
	}
	return h.Div(
		h.Data("input-otp-container", "true"),
		h.Class(containerClass),
		h.Style("position: relative; cursor: "+cursor+"; user-select: none; pointer-events: none;"),
		h.Group(kids...),
		h.Div(h.Style("position: absolute; inset: 0px; pointer-events: none;"), input),
	)
}

// Fragment returns the <noscript> input-otp emits, then the container.
// The container resolves when rendered, in the enclosing style.Scope.
func (b *Builder) Fragment() h.Fragment {
	return h.Group(h.El("noscript"), h.NodeFunc(func(w io.Writer) error {
		b.st, _ = style.FromWriter(w)
		return b.Element().Render(w)
	}))
}

// Render renders Fragment.
func (b *Builder) Render(w io.Writer) error { return b.Fragment().Render(w) }

type slotContext struct{ value []rune }

// fill hands the value to the slots among nodes.
func fill(nodes []any, ctx *slotContext) {
	for _, n := range nodes {
		switch v := n.(type) {
		case *SlotBuilder:
			v.ctx = ctx
		case *GroupBuilder:
			fill(v.Elem().ChildNodes(), ctx)
		case *h.Element:
			fill(v.ChildNodes(), ctx)
		case h.Fragment:
			fill(v, ctx)
		case []any:
			fill(v, ctx)
		}
	}
}

// GroupBuilder is returned by Group and Separator.
type GroupBuilder struct {
	h.Builder[GroupBuilder]
	separator bool
	st        style.Style
}

// Group groups slots visually.
func Group(children ...any) *GroupBuilder {
	b := &GroupBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Separator separates groups with a minus icon.
func Separator(children ...any) *GroupBuilder {
	b := &GroupBuilder{separator: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the group or separator.
func (b *GroupBuilder) Element() *h.Element {
	if b.separator {
		if b.st != "" {
			user := b.Elem().Clone().SetChildren()
			el := h.Spread(h.Div(h.Data("slot", "input-otp-separator"), h.Class(stSeparator.Class(b.st)), h.Role("separator"), lucide.Minus()), user, nil)
			if c := user.ClassName(); c != "" {
				el.SetClass(c) // {...props} replaces className
			}
			return el
		}
		return h.Spread(h.Div(h.Data("slot", "input-otp-separator"), h.Role("separator"), lucide.Minus()), b.Elem().Clone().SetChildren(), lib.Cn)
	}
	if b.st != "" {
		return h.Spread(h.Div(h.Data("slot", "input-otp-group"), h.Class(stGroup.Class(b.st, ""))), b.Elem(), lib.Cn)
	}
	return h.Spread(h.Div(h.Data("slot", "input-otp-group"), h.Class("flex items-center")), b.Elem(), lib.Cn)
}

// Render renders the group or separator.
func (b *GroupBuilder) Render(w io.Writer) error {
	b.st, _ = style.FromWriter(w)
	return b.Element().Render(w)
}

// SlotBuilder is returned by Slot.
type SlotBuilder struct {
	h.Builder[SlotBuilder]
	index int
	ctx   *slotContext
	st    style.Style
}

// Slot shows the character at index.
func Slot(index int, children ...any) *SlotBuilder {
	b := &SlotBuilder{index: index}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

const slotClass = "relative flex h-9 w-9 items-center justify-center border-y border-r border-input text-sm shadow-xs transition-all outline-none first:rounded-l-md first:border-l last:rounded-r-md aria-invalid:border-destructive data-[active=true]:z-10 data-[active=true]:border-ring data-[active=true]:ring-[3px] data-[active=true]:ring-ring/50 data-[active=true]:aria-invalid:border-destructive data-[active=true]:aria-invalid:ring-destructive/20 dark:bg-input/30 dark:data-[active=true]:aria-invalid:ring-destructive/40"

// Element resolves the slot. Server-side it is never active (the input is
// not focused before scripts run).
func (b *SlotBuilder) Element() *h.Element {
	cls := slotClass
	if b.st != "" {
		cls = stSlot.Class(b.st, "")
	}
	el := h.Spread(h.Div(h.Data("slot", "input-otp-slot"), h.Data("active", "false"), h.Class(cls)), b.Elem().Clone().SetChildren(), lib.Cn)
	if b.ctx != nil && b.index >= 0 && b.index < len(b.ctx.value) {
		el.Children(string(b.ctx.value[b.index]))
	}
	return el
}

// Render renders the slot.
func (b *SlotBuilder) Render(w io.Writer) error {
	b.st, _ = style.FromWriter(w)
	return b.Element().Render(w)
}
