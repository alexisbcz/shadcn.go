package chart

import (
	_ "embed"

	"shadcn.go/components/ui/runtime"
)

// radar.js ports recharts' RadarChart, Radar, PolarGrid, PolarAngleAxis and
// PolarRadiusAxis; it registers the RadarChart type with chart.js and
// exports the polar grid and axes renderers for other polar charts.
//
//go:embed radar.js
var radarScript string

func init() { runtime.Module("chart-radar", radarScript) }

// RadarChart is recharts' RadarChart root.
func RadarChart(args ...any) *Node { return El("RadarChart", args...) }

// Radar is recharts' Radar.
func Radar(args ...any) *Node { return El("Radar", args...) }

// PolarGrid is recharts' PolarGrid.
func PolarGrid(args ...any) *Node { return El("PolarGrid", args...) }

// PolarAngleAxis is recharts' PolarAngleAxis.
func PolarAngleAxis(args ...any) *Node { return El("PolarAngleAxis", args...) }

// PolarRadiusAxis is recharts' PolarRadiusAxis.
func PolarRadiusAxis(args ...any) *Node { return El("PolarRadiusAxis", args...) }

// GridType is PolarGrid's gridType: "polygon" (default) or "circle".
func GridType(v string) Prop { return Prop{"gridType", v} }

// RadialLines is PolarGrid's radialLines.
func RadialLines(v bool) Prop { return Prop{"radialLines", v} }

// PolarRadius is PolarGrid's polarRadius (the concentric grid radii).
func PolarRadius(v ...float64) Prop { return Prop{"polarRadius", v} }

// PolarAngles is PolarGrid's polarAngles (the radial line angles).
func PolarAngles(v ...float64) Prop { return Prop{"polarAngles", v} }

// Angle is PolarRadiusAxis's angle (and Text's angle elsewhere).
func Angle(v float64) Prop { return Prop{"angle", v} }

// Orientation is an axis orientation ("middle", "left", "right", "outer", "inner", ...).
func Orientation(v string) Prop { return Prop{"orientation", v} }

// Tick is an axis tick prop: a bool, an object of Text props, or a Fn
// rendering the tick.
func Tick(v any) Prop { return Prop{"tick", v} }

// RadarStackedTick is chart-radar-label-custom's PolarAngleAxis tick: for
// the tick at index, a <text> (13px, medium; 10px higher for the first
// tick) reading "<rows[index][a]>/<rows[index][b]>" over rows[index][label]
// (12px, muted).
func RadarStackedTick(rows []Row, a, b, label string) Fn {
	return Fn{Name: "radarStackedTick", Args: []any{rows, a, b, label}}
}
