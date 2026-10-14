// Package checkbox is the shadcn/ui Checkbox (Radix Checkbox): a control
// that allows the user to toggle between checked and not checked, with an
// optional indeterminate state.
//
//	checkbox.Checkbox(h.ID("terms"))
//	checkbox.Checkbox(h.ID("terms-2"), checkbox.Checked(true))
//	checkbox.Checkbox(checkbox.Indeterminate)
//	checkbox.Checkbox(h.Disabled())
//
// Children may include Checked (Radix defaultChecked), Indeterminate, Name,
// Value (default "on"), Required, Disabled and Form, besides any attribute.
// Like the Radix props, the name, value, required and form attributes
// (h.Name, h.Value, h.Required, h.FormAttr) are equivalent to the typed
// options: they go to the hidden input (value also to the button, required
// as aria-required), not to the button.
//
// Behavior (checkbox.js): a click (or Space) toggles the state, an
// indeterminate checkbox becomes checked, Enter does nothing, the indicator
// mounts while checked or indeterminate (Radix Presence), and form resets
// restore the initial state. The cancelable "shadcn:checked-change" event
// (detail: { checked: true | false | "indeterminate" }) is dispatched before
// every change.
//
// The button is followed by Radix's hidden <input type="checkbox"> so the
// checkbox works in native forms; like Radix, it is removed when the
// checkbox is not inside a <form> (and has no form attribute).
package checkbox

import (
	_ "embed"
	"io"
	"shadcn.go/components/style"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed checkbox.js
var script string

func init() { runtime.Module("checkbox", script) }

// Checked sets the initial checked state (Radix defaultChecked).
type Checked bool

// IndeterminateFlag is the type of Indeterminate.
type IndeterminateFlag struct{}

// Indeterminate makes the checkbox initially indeterminate
// (defaultChecked="indeterminate").
var Indeterminate = IndeterminateFlag{}

// Name is the form field name of the hidden input.
type Name string

// Value is the form value of the hidden input (default "on").
type Value string

// Required marks the checkbox as required.
type Required bool

// Disabled disables the checkbox.
type Disabled bool

// Form is the id of the form the hidden input belongs to (Radix form).
type Form string

type state int

const (
	unchecked state = iota
	checked
	indeterminate
)

// Builder is returned by Checkbox.
type Builder struct {
	h.Builder[Builder]
	state    state
	name     string
	value    string
	required *bool
	disabled bool
	form     string
}

// Checkbox returns the checkbox <button> (followed by its hidden form
// <input> when rendered).
func Checkbox(children ...any) *Builder {
	b := &Builder{value: "on"}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes the checkbox's typed options.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Checked:
		b.Checked(bool(v))
	case IndeterminateFlag:
		b.Indeterminate()
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
func (b *Builder) Checked(on bool) *Builder {
	b.state = unchecked
	if on {
		b.state = checked
	}
	return b
}

// Indeterminate makes the checkbox initially indeterminate.
func (b *Builder) Indeterminate() *Builder {
	b.state = indeterminate
	return b
}

// Name sets the form field name.
func (b *Builder) Name(v string) *Builder { b.name = v; return b }

// Value sets the form value.
func (b *Builder) Value(v string) *Builder { b.value = v; return b }

// Required marks the checkbox as required (aria-required).
func (b *Builder) Required(on bool) *Builder { b.required = &on; return b }

// Disabled disables the checkbox.
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

const rootClass = "peer size-4 shrink-0 rounded-[4px] border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground dark:bg-input/30 dark:aria-invalid:ring-destructive/40 dark:data-[state=checked]:bg-primary"

// Element resolves the checkbox to its <button> (without the hidden input;
// Render writes both).
func (b *Builder) Element() *h.Element { return b.element("") }

// styleFile is registry/bases/radix/ui/checkbox.tsx (style.Scope).
var styleFile = style.NewFile("ui/checkbox.tsx")

var (
	styleRoot      = styleFile.Cn("cn-checkbox peer relative shrink-0 outline-none after:absolute after:-inset-x-3 after:-inset-y-2 disabled:cursor-not-allowed disabled:opacity-50", style.Arg)
	styleIndicator = styleFile.Attr("cn-checkbox-indicator grid place-content-center text-current transition-none")
)

func (b *Builder) element(st style.Style) *h.Element {
	props, o := b.resolve()
	disabled := o.disabled
	base := h.Button(
		h.Type("button"),
		h.Role("checkbox"),
		h.Aria("checked", ariaChecked(b.state)),
		h.Data("state", stateName(b.state)),
		h.Value(o.value),
		h.Data("slot", "checkbox"),
	)
	if o.required != nil {
		base.Aria("required", strconv.FormatBool(*o.required))
	}
	if disabled {
		base.Data("disabled", "").Disabled(true)
	}
	out := h.Spread(base, props, nil).SetClass(lib.Cn(rootClass, props.ClassName()))
	ind := indicator(b.state, disabled)
	if st != "" {
		out.SetClass(styleRoot.Class(st, props.ClassName()))
		ind = h.Span(
			h.Data("state", stateName(b.state)),
			h.Data("slot", "checkbox-indicator"),
			h.Class(styleIndicator.Class(st)),
			h.Style("pointer-events:none"),
			lucide.Check(),
		).When(disabled, func(e *h.Element) { e.Data("disabled", "") })
	}
	if b.state != unchecked {
		out.Children(ind)
	} else {
		// Radix mounts the indicator when the checkbox gets checked; the
		// behavior takes its markup from this template (and removes it).
		out.Children(h.Template(h.Data("slot", "checkbox-indicator-template"), ind))
	}
	return out
}

func indicator(s state, disabled bool) *h.Element {
	el := h.Span(
		h.Data("state", stateName(s)),
		h.Data("slot", "checkbox-indicator"),
		h.Class("grid place-content-center text-current transition-none"),
		h.Style("pointer-events:none"),
		lucide.Check(h.Class("size-3.5")),
	)
	if disabled {
		el.Data("disabled", "")
	}
	return el
}

// input returns Radix's CheckboxBubbleInput as first rendered.
func (b *Builder) input() *h.Element {
	_, o := b.resolve()
	in := h.Input(
		h.Type("checkbox"),
		h.Aria("hidden", "true"),
		h.Value(o.value),
		h.TabIndex(-1),
		h.Style("position:absolute;pointer-events:none;opacity:0;margin:0;transform:translateX(-100%)"),
	)
	if b.state == checked {
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

// Render renders the checkbox's button and hidden input.
func (b *Builder) Render(w io.Writer) error {
	el := b.ElementIn(w)
	if el.Err() != nil {
		return el.Render(w)
	}
	return h.Group(el, b.input()).Render(w)
}

func stateName(s state) string {
	switch s {
	case checked:
		return "checked"
	case indeterminate:
		return "indeterminate"
	}
	return "unchecked"
}

func ariaChecked(s state) string {
	switch s {
	case checked:
		return "true"
	case indeterminate:
		return "mixed"
	}
	return "false"
}

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return b.element(st)
}
