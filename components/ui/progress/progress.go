// Package progress is the shadcn/ui Progress (Radix Progress): an indicator
// showing the completion progress of a task, typically displayed as a
// progress bar.
//
//	progress.Progress(progress.Value(66), h.Class("w-[60%]"))
//
// Children may include a Value (0 to 100: the indicator's position) and a
// Max (default 100, only reflected in aria-valuemax), besides any
// attribute. Like Radix, an invalid Max falls back to 100.
//
// To change the value on the client, import setProgress from the
// "progress.js" runtime module:
//
//	import { setProgress } from "./progress.js"
//	setProgress(el, 66) // el: the [data-slot=progress] root; null: indeterminate
package progress

import (
	_ "embed"
	"io"
	"math"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed progress.js
var script string

func init() { runtime.Module("progress", script) }

// Value is the progress value, between 0 and Max.
type Value float64

// Max is the maximum value (default 100).
type Max float64

const defaultMax = 100

// Builder is returned by Progress.
type Builder struct {
	h.Builder[Builder]
	value *float64
	max   *float64
}

// Progress returns the progress bar: a <div role="progressbar"> and its
// indicator.
func Progress(children ...any) *Builder {
	b := &Builder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Value and Max children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		b.Value(float64(v))
	case Max:
		b.Max(float64(v))
	default:
		return false
	}
	return true
}

// Value sets the value.
func (b *Builder) Value(v float64) *Builder { b.value = &v; return b }

// Max sets the maximum value.
func (b *Builder) Max(v float64) *Builder { b.max = &v; return b }

// Element resolves the progress bar to its <div>.
//
// Like upstream, the value only moves the indicator: shadcn's Progress
// takes the value prop for the indicator's transform and does not pass it
// to Radix, so the root and indicator always render Radix's indeterminate
// state (no aria-valuenow, data-state="indeterminate").
func (b *Builder) Element() *h.Element {
	max := float64(defaultMax)
	if b.max != nil && !math.IsNaN(*b.max) && *b.max > 0 {
		max = *b.max
	}
	base := h.Div(
		h.Aria("valuemax", num(max)),
		h.Aria("valuemin", "0"),
		h.Role("progressbar"),
		h.Data("state", "indeterminate"),
		h.Data("max", num(max)),
		h.Data("slot", "progress"),
	)
	// translateX(-${100 - (value || 0)}%)
	given := 0.0
	if b.value != nil && !math.IsNaN(*b.value) {
		given = *b.value
	}
	ind := h.Div(
		h.Data("state", "indeterminate"),
		h.Data("max", num(max)),
		h.Data("slot", "progress-indicator"),
		h.Class("h-full w-full flex-1 bg-primary transition-all"),
		h.Style("transform:translateX(-"+num(100-given)+"%)"),
	)
	props := b.Elem()
	return h.Spread(base, props, nil).
		SetClass(lib.Cn("relative h-2 w-full overflow-hidden rounded-full bg-primary/20", props.ClassName())).
		SetChildren(ind)
}

// Render renders the progress bar.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// num formats a number like JavaScript's Number#toString for the values
// a progress bar takes.
func num(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
