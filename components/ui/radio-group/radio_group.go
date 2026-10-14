// Package radiogroup is the shadcn/ui RadioGroup (Radix RadioGroup): a set
// of checkable buttons, known as radio buttons, where no more than one can
// be checked at a time.
//
//	radiogroup.RadioGroup(radiogroup.Value("comfortable"),
//	    h.Div(h.Class("flex items-center gap-3"),
//	        radiogroup.Item(radiogroup.Value("default"), h.ID("r1")),
//	        label.Label(h.For("r1"), "Default"),
//	    ),
//	    ...
//	)
//
// Root children may include Value (the initially checked item), Name,
// Required, Disabled, Orientation and Loop; items take a Value (default
// "on", like Radix). Items may sit anywhere under the root.
//
// Behavior (radio-group.js): Radix's roving focus (arrow keys move focus
// and check the focused item, loop by default), Space checks, Enter does
// not, and the cancelable "shadcn:value-change" event (detail: { value })
// is dispatched on the group before every change.
//
// Each item is followed by Radix's hidden <input type="radio"> so the
// group works in native forms; like Radix, it is removed when the item is
// not inside a <form> (and has no form attribute).
//
// Right to left (Radix's dir prop): pass h.Dir("rtl") to the root; it
// replaces the default dir="ltr" and the arrow keys follow it.
package radiogroup

import (
	_ "embed"
	"fmt"
	"io"
	"shadcn.go/components/style"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed radio_group.js
var script string

func init() { runtime.Module("radio-group", script) }

// Value is an item's value; on the root, the initially checked value.
type Value string

// Name is the form field name of the group's hidden inputs.
type Name string

// Required marks the group as required.
type Required bool

// Disabled, passed to the root, disables every item.
type Disabled bool

// Orientation sets the arrow keys that move focus ("horizontal" or
// "vertical"); unset, both work.
type Orientation string

// Orientations.
const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// Loop sets whether keyboard navigation wraps around (default true).
type Loop bool

type context struct {
	value    *string
	name     string
	required bool
	disabled bool
	orient   Orientation
	style    style.Style
}

// Builder is returned by RadioGroup.
type Builder struct {
	h.Builder[Builder]
	ctx  context
	loop bool
}

// RadioGroup returns the group's root <div role="radiogroup">.
func RadioGroup(children ...any) *Builder {
	b := &Builder{loop: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the group's typed options.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		s := string(v)
		b.ctx.value = &s
	case Name:
		b.ctx.name = string(v)
	case Required:
		b.ctx.required = bool(v)
	case Disabled:
		b.ctx.disabled = bool(v)
	case Orientation:
		b.ctx.orient = v
	case Loop:
		b.loop = bool(v)
	default:
		return false
	}
	return true
}

// Value sets the initially checked value.
func (b *Builder) Value(v string) *Builder { b.ctx.value = &v; return b }

// Name sets the form field name.
func (b *Builder) Name(v string) *Builder { b.ctx.name = v; return b }

// Required marks the group as required.
func (b *Builder) Required(on bool) *Builder { b.ctx.required = on; return b }

// Disabled disables every item.
func (b *Builder) Disabled(on bool) *Builder { b.ctx.disabled = on; return b }

// Orientation sets the navigation orientation.
func (b *Builder) Orientation(o Orientation) *Builder { b.ctx.orient = o; return b }

// Loop sets whether keyboard navigation wraps around.
func (b *Builder) Loop(on bool) *Builder { b.loop = on; return b }

// Element resolves the group to its <div>.
func (b *Builder) Element() *h.Element { return b.element("") }

// element renders in style st ("" for new-york-v4).
func (b *Builder) element(st style.Style) *h.Element {
	ctx := b.ctx
	ctx.style = st
	focusable := 0
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ItemBuilder:
			c.ctx = &ctx
			if !c.disabled() {
				focusable++
			}
		case *Builder:
			return false
		}
		return true
	})
	base := h.Div(
		h.TabIndex(h.IfElse(focusable == 0, -1, 0)),
		h.Style("outline:none"),
		h.Role("radiogroup"),
		h.Aria("required", strconv.FormatBool(ctx.required)),
		h.Dir("ltr"),
		h.Data("slot", "radio-group"),
	)
	if ctx.orient != "" {
		base.Data("orientation", string(ctx.orient))
		base.Aria("orientation", string(ctx.orient))
	}
	if ctx.disabled {
		base.Data("disabled", "")
	}
	if !b.loop {
		base.Data("loop", "false")
	}
	props := b.Elem()
	if st != "" {
		return h.Spread(base, props, nil).SetClass(styleRoot.Class(st, props.ClassName()))
	}
	return h.Spread(base, props, nil).SetClass(lib.Cn("grid gap-3", props.ClassName()))
}

// Render renders the group.
// Inside style.Scope it renders that style.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return b.element(st)
}

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	ctx   *context
	value string
}

// Item returns a radio <button> of the group, followed by its hidden
// form <input>. Children may include a Value.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{value: "on"}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Value children.
func (b *ItemBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Value); ok {
		b.value = string(v)
		return true
	}
	return false
}

// Value sets the item's value.
func (b *ItemBuilder) Value(v string) *ItemBuilder { b.value = v; return b }

// Disabled disables the item.
func (b *ItemBuilder) Disabled(on bool) *ItemBuilder { b.Elem().Disabled(on); return b }

func (b *ItemBuilder) disabled() bool {
	return (b.ctx != nil && b.ctx.disabled) || b.Elem().Has("disabled")
}

func (b *ItemBuilder) checked() bool {
	return b.ctx != nil && b.ctx.value != nil && *b.ctx.value == b.value
}

const itemClass = "aspect-square size-4 shrink-0 rounded-full border border-input text-primary shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:bg-input/30 dark:aria-invalid:ring-destructive/40"

// Element resolves the item to its <button> (without the hidden input;
// Render writes both).
func (b *ItemBuilder) Element() *h.Element {
	if b.ctx == nil {
		return h.Errored(fmt.Errorf("radiogroup: Item must be inside a RadioGroup"))
	}
	checked, disabled := b.checked(), b.disabled()
	base := h.Button(
		h.Type("button"),
		h.Role("radio"),
		h.Aria("checked", strconv.FormatBool(checked)),
		h.Data("state", state(checked)),
		h.Value(b.value),
		h.TabIndex(-1),
		h.Data("radix-collection-item", ""),
		h.Data("slot", "radio-group-item"),
	)
	if b.ctx.orient != "" {
		base.Data("orientation", string(b.ctx.orient))
	}
	if disabled {
		base.Data("disabled", "").Disabled(true)
	}
	props := b.Elem()
	out := h.Spread(base, props, nil).SetClass(lib.Cn(itemClass, props.ClassName()))
	ind := indicator(checked, disabled)
	if st := b.ctx.style; st != "" {
		out.SetClass(styleItem.Class(st, props.ClassName()))
		ind = h.Span(
			h.Data("state", state(checked)),
			h.Data("slot", "radio-group-indicator"),
			h.Class(styleIndicator.Class(st)),
			h.Span(h.Class(styleIndicatorIcon.Class(st))),
		).When(disabled, func(e *h.Element) { e.Data("disabled", "") })
	}
	if checked {
		out.Children(ind)
	} else {
		// Radix mounts the indicator when the item gets checked; the
		// behavior takes its markup from this template (and removes it).
		out.Children(h.Template(h.Data("slot", "radio-group-indicator-template"), ind))
	}
	return out
}

// styleFile is registry/bases/radix/ui/radio-group.tsx (style.Scope).
var styleFile = style.NewFile("ui/radio-group.tsx")

var (
	styleRoot          = styleFile.Cn("cn-radio-group w-full", style.Arg)
	styleItem          = styleFile.Cn("cn-radio-group-item group/radio-group-item peer relative aspect-square shrink-0 border outline-none after:absolute after:-inset-x-3 after:-inset-y-2 disabled:cursor-not-allowed disabled:opacity-50", style.Arg)
	styleIndicator     = styleFile.Attr("cn-radio-group-indicator")
	styleIndicatorIcon = styleFile.Attr("cn-radio-group-indicator-icon")
)

func indicator(checked, disabled bool) *h.Element {
	el := h.Span(
		h.Data("state", state(checked)),
		h.Data("slot", "radio-group-indicator"),
		h.Class("relative flex items-center justify-center"),
		lucide.Circle(h.Class("absolute top-1/2 left-1/2 size-2 -translate-x-1/2 -translate-y-1/2 fill-primary")),
	)
	if disabled {
		el.Data("disabled", "")
	}
	return el
}

// input returns Radix's RadioBubbleInput as first rendered.
func (b *ItemBuilder) input() *h.Element {
	in := h.Input(
		h.Type("radio"),
		h.Aria("hidden", "true"),
		h.Value(b.value),
		h.TabIndex(-1),
		h.Style("transform:translateX(-100%);position:absolute;pointer-events:none;opacity:0;margin:0"),
	)
	if b.checked() {
		in.Checked(true)
	}
	if b.ctx.name != "" {
		in.Name(b.ctx.name)
	}
	if b.ctx.required {
		in.Required(true)
	}
	if b.disabled() {
		in.Disabled(true)
	}
	if f, ok := b.Elem().GetAttr("form"); ok {
		in.Form(f)
	}
	return in
}

// Render renders the item's button and hidden input.
func (b *ItemBuilder) Render(w io.Writer) error {
	el := b.Element()
	if el.Err() != nil {
		return el.Render(w)
	}
	return h.Group(el, b.input()).Render(w)
}

func state(checked bool) string {
	if checked {
		return "checked"
	}
	return "unchecked"
}
