package chart

import (
	_ "embed"
	"fmt"

	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

// pie.js ports recharts' PieChart, Pie, Sector and polar labels; it
// registers the PieChart type with chart.js.
//
//go:embed pie.js
var pieScript string

func init() { runtime.Module("chart-pie", pieScript) }

// PieChart is recharts' PieChart root.
func PieChart(args ...any) *Node { return El("PieChart", args...) }

// Pie is recharts' Pie. Its data is given with Data; Cell, Label and
// LabelList children are supported.
func Pie(args ...any) *Node { return El("Pie", args...) }

// Label is recharts' Label (inside a Pie, positioned on the polar view
// box); give it a Content function or a Value.
func Label(args ...any) *Node { return El("Label", args...) }

// Polar props.
func InnerRadius(v any) Prop  { return Prop{"innerRadius", v} }
func OuterRadius(v any) Prop  { return Prop{"outerRadius", v} }
func LabelLine(v any) Prop    { return Prop{"labelLine", v} }
func Content(f Fn) Prop       { return Prop{"content", f} }
func Value(v any) Prop        { return Prop{"value", v} }
func PaddingAngle(v any) Prop { return Prop{"paddingAngle", v} }

// LabelProp is a graphical item's label prop (true, or an object of Text
// props, or a Fn rendering the label).
func LabelProp(v any) Prop { return Prop{"label", v} }

// ActiveSector is a Pie shape: the sector at index is drawn grow pixels
// further out; ({index, outerRadius, ...props}) => index === i ?
// <Sector {...props} outerRadius={outerRadius + grow}/> : <Sector {...props} outerRadius={outerRadius}/>.
func ActiveSector(index int, grow float64) Fn {
	return Fn{Name: "pieActiveSector", Args: []any{index, grow}}
}

// ActiveSectorRing is chart-pie-interactive's shape: the sector at index is
// drawn 10 pixels further out with a ring from outerRadius+12 to
// outerRadius+25 around it, in a <g>.
func ActiveSectorRing(index int) Fn { return Fn{Name: "pieActiveSectorRing", Args: []any{index}} }

// PieCenterText is a Label content rendering value (large) and caption at the
// view box center, as the donut examples do.
func PieCenterText(value any, caption string) Fn {
	return Fn{Name: "pieCenterText", Args: []any{value, caption}}
}

// PieValueLabel is a Pie label rendering payload[field] as a <text> with
// the label's position and the given fill.
func PieValueLabel(field, fill string) Fn {
	return Fn{Name: "pieValueLabel", Args: []any{field, fill}}
}

// PayloadConfigLabel is the labelFormatter (_, payload) => chartConfig[payload?.[0].dataKey].label.
func PayloadConfigLabel() Fn { return Fn{Name: "piePayloadConfigLabel"} }

// Style is ChartStyle: the config's --color-* variables for elements with
// data-chart=id (rendered by an element other than the container, as
// chart-pie-interactive does on its card).
func Style(id string, config Config, nonce Nonce) h.Node {
	if !safeKey.MatchString(id) {
		return h.Errored(fmt.Errorf("chart: id %q: only letters, digits, - and _ are allowed", id))
	}
	css, err := config.css(id)
	if err != nil {
		return h.Errored(err)
	}
	if css == "" {
		return nil
	}
	st := h.StyleEl(h.UnsafeCSS(css))
	if nonce != "" {
		st.Nonce(string(nonce))
	}
	return st
}
