// Package switchui is the shadcn/ui Switch (Radix Switch): a control that
// allows the user to toggle between checked and not checked.
//
//	switchui.Switch(h.ID("airplane-mode"))
//	switchui.Switch(switchui.Checked(true), switchui.SizeSm)
//
// Children may include Checked (Radix defaultChecked), a Size, Name, Value
// (default "on"), Required, Disabled and Form, besides any attribute. Like
// the Radix props, the name, value, required and form attributes (h.Name,
// h.Value, h.Required, h.FormAttr) are equivalent to the typed options:
// they go to the hidden input (value also to the button, required as
// aria-required), not to the button.
//
// Behavior (switch.js): a click (or Space, Enter natively) toggles the
// switch. The cancelable "shadcn:checked-change" event (detail: { checked })
// is dispatched before every change.
//
// The button is followed by Radix's hidden <input type="checkbox"> so the
// switch works in native forms; like Radix, it is removed when the switch is
// not inside a <form> (and has no form attribute).
package switchui

import (
	_ "embed"
	"io"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed switch.js
var script string

func init() { runtime.Module("switch", script) }

// Checked sets the initial checked state (Radix defaultChecked).
type Checked bool

// Size selects the switch's size.
type Size string

// Switch sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
)

// Name is the form field name of the hidden input.
type Name string

// Value is the form value of the hidden input (default "on").
type Value string

// Required marks the switch as required.
type Required bool

// Disabled disables the switch.
type Disabled bool

// Form is the id of the form the hidden input belongs to (Radix form).
type Form string

// Builder is returned by Switch.
type Builder struct {
	h.Builder[Builder]
	checked  bool
	size     Size
	name     string
	value    string
	required *bool
	disabled bool
	form     string
}

// Switch returns the switch <button> with its thumb (followed by its hidden
// form <input> when rendered).
func Switch(children ...any) *Builder {
	b := &Builder{value: "on", size: SizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes the switch's typed options.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Checked:
		b.checked = bool(v)
	case Size:
		b.Size(v)
	case Name:
		b.name = string(v)
	case Value:
		b.value = string(v)
	case Required:
		b.Required(bool(v))
	case Disabled:
		b.disabled = bool(v)
	case Form:
		b.form = string(v)
	default:
		return false
	}
	return true
}

// Checked sets the initial checked state.
func (b *Builder) Checked(on bool) *Builder { b.checked = on; return b }

// Size sets the size. An empty value selects the default.
func (b *Builder) Size(s Size) *Builder {
	if s == "" {
		s = SizeDefault
	}
	b.size = s
	return b
}

// Name sets the form field name.
func (b *Builder) Name(v string) *Builder { b.name = v; return b }

// Value sets the form value.
func (b *Builder) Value(v string) *Builder { b.value = v; return b }

// Required marks the switch as required (aria-required).
func (b *Builder) Required(on bool) *Builder { b.required = &on; return b }

// Disabled disables the switch.
func (b *Builder) Disabled(on bool) *Builder { b.disabled = on; return b }

// Form sets the id of the form the hidden input belongs to.
func (b *Builder) Form(v string) *Builder { b.form = v; return b }

// options are the Radix props, from the typed options and the equivalent
// attributes.
type options struct {
	name, value, form string
	required          *bool
	disabled          bool
}

// resolve returns the attributes for the button (without the ones Radix
// takes as props) and the props.
func (b *Builder) resolve() (*h.Element, options) {
	props := b.Elem().Clone()
	o := options{name: b.name, value: b.value, form: b.form, required: b.required, disabled: b.disabled}
	if v, ok := props.GetAttr("name"); ok {
		o.name = v
		props.Remove("name")
	}
	if v, ok := props.GetAttr("value"); ok {
		o.value = v
		props.Remove("value")
	}
	if v, ok := props.GetAttr("form"); ok {
		o.form = v
		props.Remove("form")
	}
	if props.Has("required") {
		on := true
		o.required = &on
		props.Remove("required")
	}
	if props.Has("disabled") {
		o.disabled = true
		props.Remove("disabled")
	}
	return props, o
}

const (
	rootClass  = "peer group/switch inline-flex shrink-0 items-center rounded-full border border-transparent shadow-xs transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 data-[size=default]:h-[1.15rem] data-[size=default]:w-8 data-[size=sm]:h-3.5 data-[size=sm]:w-6 data-[state=checked]:bg-primary data-[state=unchecked]:bg-input dark:data-[state=unchecked]:bg-input/80"
	thumbClass = "pointer-events-none block rounded-full bg-background ring-0 transition-transform group-data-[size=default]/switch:size-4 group-data-[size=sm]/switch:size-3 data-[state=checked]:translate-x-[calc(100%-2px)] data-[state=unchecked]:translate-x-0 dark:data-[state=checked]:bg-primary-foreground dark:data-[state=unchecked]:bg-foreground"
)

var thumbClassCn = lib.Cn(thumbClass)

// Element resolves the switch to its <button> (without the hidden input;
// Render writes both).
func (b *Builder) Element() *h.Element {
	props, o := b.resolve()
	disabled := o.disabled
	base := h.Button(
		h.Type("button"),
		h.Role("switch"),
		h.Aria("checked", strconv.FormatBool(b.checked)),
		h.Data("state", state(b.checked)),
		h.Value(o.value),
		h.Data("slot", "switch"),
		h.Data("size", string(b.size)),
	)
	if o.required != nil {
		base.Aria("required", strconv.FormatBool(*o.required))
	}
	if disabled {
		base.Data("disabled", "").Disabled(true)
	}
	out := h.Spread(base, props, nil).SetClass(lib.Cn(rootClass, props.ClassName()))
	thumb := h.Span(
		h.Data("state", state(b.checked)),
		h.Data("slot", "switch-thumb"),
		h.Class(thumbClassCn),
	)
	if disabled {
		thumb.Data("disabled", "")
	}
	return out.SetChildren(thumb)
}

// input returns Radix's SwitchBubbleInput as first rendered.
func (b *Builder) input() *h.Element {
	_, o := b.resolve()
	in := h.Input(
		h.Type("checkbox"),
		h.Aria("hidden", "true"),
		h.Value(o.value),
		h.TabIndex(-1),
		h.Style("transform:translateX(-100%);position:absolute;pointer-events:none;opacity:0;margin:0"),
	)
	if b.checked {
		in.Checked(true)
	}
	if o.required != nil && *o.required {
		in.Required(true)
	}
	if o.disabled {
		in.Disabled(true)
	}
	if o.name != "" {
		in.Name(o.name)
	}
	if o.form != "" {
		in.Form(o.form)
	}
	return in
}

// Render renders the switch's button and hidden input.
func (b *Builder) Render(w io.Writer) error {
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
