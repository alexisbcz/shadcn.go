// Package aspectratio is the shadcn/ui AspectRatio (Radix AspectRatio):
// displays content within a desired ratio.
//
//	aspectratio.AspectRatio(aspectratio.Ratio(16.0/9), h.Class("rounded-lg bg-muted"), img)
package aspectratio

import (
	"io"
	"strconv"
	"strings"

	h "shadcn.go/libhtml"
)

// Ratio, passed as a child, sets the width / height ratio (default 1).
type Ratio float64

// Builder is returned by AspectRatio.
type Builder struct {
	h.Builder[Builder]
	ratio float64
}

// AspectRatio returns the ratio box: a relative wrapper <div> sized with the
// padding-bottom trick, holding an absolutely positioned <div> that carries
// the attributes and children. Children may include a Ratio.
func AspectRatio(children ...any) *Builder {
	b := &Builder{ratio: 1}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Ratio children.
func (b *Builder) AcceptChild(c any) bool {
	if r, ok := c.(Ratio); ok {
		b.ratio = float64(r)
		return true
	}
	return false
}

// Ratio sets the width / height ratio.
func (b *Builder) Ratio(r float64) *Builder { b.ratio = r; return b }

// jsNumber formats a float like JavaScript's Number#toString for the
// values a ratio produces.
func jsNumber(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// Element resolves the aspect ratio to its wrapper <div>.
func (b *Builder) Element() *h.Element {
	props := b.Elem().Clone()
	style, _ := props.GetAttr("style")
	style = strings.TrimSuffix(strings.TrimSpace(style), ";")
	if style != "" {
		style += ";"
	}
	props.Attr("style", style+"position:absolute;top:0;right:0;bottom:0;left:0")
	inner := h.Spread(h.Div(h.Data("slot", "aspect-ratio")), props, nil) // className is passed through, not cn()ed
	return h.Div(
		h.Style("position:relative;width:100%;padding-bottom:"+jsNumber(100/b.ratio)+"%"),
		h.Data("radix-aspect-ratio-wrapper", ""),
		inner,
	)
}

// Render renders the aspect ratio.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }
