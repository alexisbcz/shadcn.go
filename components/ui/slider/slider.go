// Package slider is the shadcn/ui Slider (Radix Slider): an input where the
// user selects a value from within a given range.
//
//	slider.Slider(slider.DefaultValue(50), slider.Max(100), slider.Step(1), h.Class("w-[60%]"))
//	slider.Slider(slider.DefaultValue(25, 75))                     // a range: two thumbs
//	slider.Slider(slider.DefaultValue(50), slider.OrientationVertical)
//
// Children may include DefaultValue (one value per thumb), Min (default 0),
// Max (default 100), Step (default 1), MinStepsBetweenThumbs, an
// Orientation, Inverted, Disabled, Name and Form, besides any attribute.
// Pass h.Dir("rtl") for right to left. Like the Radix props, the name, form
// and disabled attributes (h.Name, h.FormAttr, h.Disabled) are equivalent
// to Name, Form and Disabled: they are not rendered on the root.
//
// Like upstream, without a DefaultValue the slider renders two thumbs
// (shadcn's [min, max]) but holds Radix's default value [min]: the second
// thumb is hidden.
//
// Behavior (slider.js): pointer down on the track moves the closest thumb
// there and drags it; dragging a thumb moves it; arrow keys step
// (Shift: 10 steps), Page Up/Down step by 10, Home/End go to the minimum
// (first thumb) and maximum (last thumb). The cancelable
// "shadcn:value-change" event (detail: { value: number[] }) is dispatched
// on the root before every change, and "shadcn:value-commit" (detail:
// { value }) at the end of a drag or keyboard change.
//
// Inside a <form> each thumb is followed by Radix's hidden <input> carrying
// its value (named Name, or Name+"[]" for several thumbs); outside a form
// the behavior removes them, like Radix.
package slider

import (
	_ "embed"
	"fmt"
	"io"
	"math"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed slider.js
var script string

func init() { runtime.Module("slider", script) }

// Values is the initial value of each thumb (Radix defaultValue).
type Values []float64

// DefaultValue returns the initial value of each thumb.
func DefaultValue(v ...float64) Values { return Values(v) }

// Min is the minimum value (default 0).
type Min float64

// Max is the maximum value (default 100).
type Max float64

// Step is the stepping interval (default 1).
type Step float64

// MinStepsBetweenThumbs is the minimum number of steps between thumbs.
type MinStepsBetweenThumbs float64

// Orientation sets the slider's orientation.
type Orientation string

// Orientations.
const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// Inverted makes the slider go from max to min (right to left, or top to
// bottom).
type Inverted bool

// Disabled disables the slider.
type Disabled bool

// Name is the form field name of the hidden inputs.
type Name string

// Form is the id of the form the hidden inputs belong to.
type Form string

// Builder is returned by Slider.
type Builder struct {
	h.Builder[Builder]
	values   Values
	min, max float64
	step     float64
	minSteps float64
	orient   Orientation
	inverted bool
	disabled bool
	name     string
	form     string
}

// Slider returns the slider's root <span> with its track, range and thumbs.
func Slider(children ...any) *Builder {
	b := &Builder{min: 0, max: 100, step: 1, orient: OrientationHorizontal}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// AcceptChild consumes the slider's typed options.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Values:
		b.values = v
	case Min:
		b.min = float64(v)
	case Max:
		b.max = float64(v)
	case Step:
		b.step = float64(v)
	case MinStepsBetweenThumbs:
		b.minSteps = float64(v)
	case Orientation:
		b.Orientation(v)
	case Inverted:
		b.inverted = bool(v)
	case Disabled:
		b.disabled = bool(v)
	case Name:
		b.name = string(v)
	case Form:
		b.form = string(v)
	default:
		return false
	}
	return true
}

// DefaultValue sets the initial value of each thumb.
func (b *Builder) DefaultValue(v ...float64) *Builder { b.values = v; return b }

// Min sets the minimum value.
func (b *Builder) Min(v float64) *Builder { b.min = v; return b }

// Max sets the maximum value.
func (b *Builder) Max(v float64) *Builder { b.max = v; return b }

// Step sets the stepping interval.
func (b *Builder) Step(v float64) *Builder { b.step = v; return b }

// MinStepsBetweenThumbs sets the minimum number of steps between thumbs.
func (b *Builder) MinStepsBetweenThumbs(v float64) *Builder { b.minSteps = v; return b }

// Orientation sets the orientation. An empty value selects horizontal.
func (b *Builder) Orientation(o Orientation) *Builder {
	if o == "" {
		o = OrientationHorizontal
	}
	b.orient = o
	return b
}

// Inverted inverts the slider.
func (b *Builder) Inverted(on bool) *Builder { b.inverted = on; return b }

// Disabled disables the slider.
func (b *Builder) Disabled(on bool) *Builder { b.disabled = on; return b }

// Name sets the form field name.
func (b *Builder) Name(v string) *Builder { b.name = v; return b }

// Form sets the id of the form the hidden inputs belong to.
func (b *Builder) Form(v string) *Builder { b.form = v; return b }

const (
	rootClass  = "relative flex w-full touch-none items-center select-none data-[disabled]:opacity-50 data-[orientation=vertical]:h-full data-[orientation=vertical]:min-h-44 data-[orientation=vertical]:w-auto data-[orientation=vertical]:flex-col"
	trackClass = "relative grow overflow-hidden rounded-full bg-muted data-[orientation=horizontal]:h-1.5 data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-1.5"
	rangeClass = "absolute bg-primary data-[orientation=horizontal]:h-full data-[orientation=vertical]:w-full"
	thumbClass = "block size-4 shrink-0 rounded-full border border-primary bg-white shadow-sm ring-ring/50 transition-[color,box-shadow] hover:ring-4 focus-visible:ring-4 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
)

var (
	trackClassCn = lib.Cn(trackClass)
	rangeClassCn = lib.Cn(rangeClass)
	thumbClassCn = lib.Cn(thumbClass)
)

// defaultThumbSize is the size of a thumb with the default classes
// (size-4); the behavior measures the real one, like Radix's useSize.
const defaultThumbSize = 16

// Element resolves the slider to its root <span>.
func (b *Builder) Element() *h.Element {
	if b.step <= 0 || math.IsNaN(b.step) {
		return h.Errored(fmt.Errorf("slider: Step must be positive, got %v", b.step))
	}
	if !(b.max > b.min) {
		return h.Errored(fmt.Errorf("slider: Max (%v) must be greater than Min (%v)", b.max, b.min))
	}
	props := b.Elem().Clone()
	name, form, disabled := b.name, b.form, b.disabled
	if v, ok := props.GetAttr("name"); ok {
		name = v
		props.Remove("name")
	}
	if v, ok := props.GetAttr("form"); ok {
		form = v
		props.Remove("form")
	}
	if props.Has("disabled") {
		disabled = true
		props.Remove("disabled")
	}
	horizontal := b.orient != OrientationVertical
	orient := string(b.orient)

	// Edges and direction (Radix SliderHorizontal / SliderVertical).
	var start, end, cssVar string
	dir := 1.0
	if horizontal {
		d, _ := props.GetAttr("dir")
		if d == "" {
			d = "ltr"
		}
		fromLeft := (d == "ltr") != b.inverted
		start, end = "left", "right"
		if !fromLeft {
			start, end, dir = "right", "left", -1
		}
		cssVar = "translateX(-50%)"
	} else {
		start, end = "bottom", "top"
		if b.inverted {
			start, end, dir = "top", "bottom", -1
		}
		cssVar = "translateY(50%)"
	}

	// shadcn renders one thumb per value of value ?? defaultValue ?? [min, max];
	// Radix holds defaultValue ?? [min].
	values := []float64(b.values)
	thumbs := len(values)
	if b.values == nil {
		values = []float64{b.min}
		thumbs = 2
	}

	base := h.Span(
		h.Aria("disabled", strconv.FormatBool(disabled)),
		h.Data("orientation", orient),
		h.Data("slot", "slider"),
	)
	if horizontal {
		base.Attr("dir", "ltr")
	}
	if disabled {
		base.Data("disabled", "")
	}
	// Props with no DOM trace, for the behavior (which removes them).
	if b.step != 1 {
		base.Data("step", num(b.step))
	}
	if b.minSteps != 0 {
		base.Data("min-steps-between-thumbs", num(b.minSteps))
	}
	if b.inverted {
		base.Data("inverted", "")
	}
	root := h.Spread(base, props, nil).SetClass(lib.Cn(rootClass, props.ClassName()))
	style, _ := props.GetAttr("style")
	if style != "" {
		style += ";"
	}
	root.Style(style + "--radix-slider-thumb-transform:" + cssVar)

	// Track and range.
	pcts := make([]float64, len(values))
	for i, v := range values {
		pcts[i] = percentage(v, b.min, b.max)
	}
	offsetStart := 0.0
	if len(values) > 1 {
		offsetStart = minOf(pcts)
	}
	rangeStyle := start + ":" + num(offsetStart) + "%"
	if len(values) > 0 {
		// With no values Radix computes 100 - Math.max() = Infinity, an
		// invalid length the browser drops.
		rangeStyle += ";" + end + ":" + num(100-maxOf(pcts)) + "%"
	}
	rng := h.Span(h.Data("orientation", orient), h.Data("slot", "slider-range"), h.Class(rangeClassCn), h.Style(rangeStyle))
	track := h.Span(h.Data("orientation", orient), h.Data("slot", "slider-track"), h.Class(trackClassCn), rng)
	if disabled {
		rng.Data("disabled", "")
		track.Data("disabled", "")
	}

	children := []any{track}
	for i := 0; i < thumbs; i++ {
		var value *float64
		if i < len(values) {
			value = &values[i]
		}
		percent, size := 0.0, 0.0
		if value != nil {
			percent, size = percentage(*value, b.min, b.max), defaultThumbSize
		}
		offset := 0.0
		if size != 0 {
			offset = thumbInBoundsOffset(size, percent, dir)
		}
		thumb := h.Span(
			h.Role("slider"),
			h.Aria("valuemin", num(b.min)),
			h.Aria("valuemax", num(b.max)),
			h.Aria("orientation", orient),
			h.Data("orientation", orient),
			h.Data("slot", "slider-thumb"),
			h.Class(thumbClassCn),
			h.Data("radix-collection-item", ""),
		)
		if label := thumbLabel(i, len(values)); label != "" {
			thumb.Aria("label", label)
		}
		if value != nil {
			thumb.Aria("valuenow", num(*value))
			// React rendered style={display: "none"} first (before the
			// thumb knew its index), then removed the declaration.
			thumb.Style("")
		} else {
			thumb.Style("display:none")
		}
		if disabled {
			thumb.Data("disabled", "")
		} else {
			thumb.TabIndex(0)
		}
		wrapper := h.Span(
			h.Style("transform:var(--radix-slider-thumb-transform);position:absolute;"+start+":calc("+num(percent)+"% + "+num(offset)+"px)"),
			thumb,
		)
		// Radix's bubble input; the behavior removes it outside forms.
		in := h.Input(h.Style("display:none"))
		if name != "" {
			inputName := name
			if len(values) > 1 {
				inputName += "[]"
			}
			in.Name(inputName)
		}
		if form != "" {
			in.Form(form)
		}
		if value != nil {
			in.Value(num(*value))
		}
		wrapper.Children(in)
		children = append(children, wrapper)
	}
	return root.SetChildren(children...)
}

// Render renders the slider.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

func percentage(value, min, max float64) float64 {
	maxSteps := max - min
	percentPerStep := 100 / maxSteps
	p := percentPerStep * (value - min)
	return math.Max(0, math.Min(100, p))
}

func thumbInBoundsOffset(width, left, direction float64) float64 {
	halfWidth := width / 2
	offset := linearScale([2]float64{0, 50}, [2]float64{0, halfWidth})
	return (halfWidth - offset(left)*direction) * direction
}

func linearScale(input, output [2]float64) func(float64) float64 {
	return func(v float64) float64 {
		if input[0] == input[1] || output[0] == output[1] {
			return output[0]
		}
		ratio := (output[1] - output[0]) / (input[1] - input[0])
		return output[0] + ratio*(v-input[0])
	}
}

func thumbLabel(index, total int) string {
	switch {
	case total > 2:
		return fmt.Sprintf("Value %d of %d", index+1, total)
	case total == 2:
		return [2]string{"Minimum", "Maximum"}[index]
	}
	return ""
}

func minOf(v []float64) float64 {
	m := v[0]
	for _, x := range v[1:] {
		m = math.Min(m, x)
	}
	return m
}

func maxOf(v []float64) float64 {
	m := v[0]
	for _, x := range v[1:] {
		m = math.Max(m, x)
	}
	return m
}

// num formats a number like JavaScript's Number#toString does for the
// values a slider takes.
func num(f float64) string {
	if f == 0 {
		return "0" // also -0, like JavaScript
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
