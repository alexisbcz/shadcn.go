package chart

import (
	_ "embed"

	"shadcn.go/components/ui/runtime"
)

// radial.js adds recharts' RadialBarChart to the runtime: build it with
// El("RadialBarChart", ...), El("RadialBar", ...), El("PolarGrid", ...),
// El("PolarRadiusAxis", ..., El("Label", P("content", Fn{Name: "centerText", ...}))).
//
//go:embed radial.js
var radialScript string

func init() { runtime.Module("chart-radial", radialScript) }
