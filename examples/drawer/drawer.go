// Package drawer ports the upstream drawer-* examples.
package drawer

import (
	_ "embed"
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/drawer"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// demo.js keeps drawer-demo's goal (React state of the example) and
// updates the number and the buttons' disabled state.
//
//go:embed demo.js
var demoScript string

// media.js swaps drawer-dialog's Dialog and Drawer versions.
//
//go:embed media.js
var mediaScript string

func init() {
	runtime.Module("example-drawer-demo", demoScript)
	runtime.Module("example-drawer-dialog", mediaScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "drawer", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("drawer-demo", Demo),
		ex("drawer-dialog", DialogExample),
		ex("drawer-directions", Directions),
	)
}

var goals = []float64{400, 300, 200, 300, 200, 278, 189, 239, 300, 200, 278, 189, 349}

func Demo() h.Node {
	return drawer.Drawer(
		drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open Drawer")),
		drawer.Content(
			h.Div(h.Class("mx-auto w-full max-w-sm"),
				drawer.Header(
					drawer.Title("Move Goal"),
					drawer.Description("Set your daily activity goal."),
				),
				h.Div(h.Class("p-4 pb-0"),
					h.Div(h.Class("flex items-center justify-center space-x-2"),
						button.Button(button.VariantOutline, button.SizeIcon, h.Class("h-8 w-8 shrink-0 rounded-full"),
							lucide.Minus(),
							h.Span(h.Class("sr-only"), "Decrease"),
						),
						h.Div(h.Class("flex-1 text-center"),
							h.Div(h.Class("text-7xl font-bold tracking-tighter"), "350"),
							h.Div(h.Class("text-[0.70rem] text-muted-foreground uppercase"), "Calories/day"),
						),
						button.Button(button.VariantOutline, button.SizeIcon, h.Class("h-8 w-8 shrink-0 rounded-full"),
							lucide.Plus(),
							h.Span(h.Class("sr-only"), "Increase"),
						),
					),
					h.Div(h.Class("mt-3 h-[120px]"), barChart(352, 120, goals)),
				),
				drawer.Footer(
					button.Button("Submit"),
					drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
				),
			),
		),
	)
}

// barChart renders what recharts 3.8 renders for
// <ResponsiveContainer width="100%" height="100%"><BarChart data={data}>
// <Bar dataKey="goal" style={{ fill: "hsl(var(--foreground))", opacity: 0.9 }} />
// </BarChart></ResponsiveContainer> once the container measured width x height
// (default 5px margins, a band category axis with a 10% category gap and a
// [0, nice max] number axis, neither drawn).
func barChart(width, height float64, data []float64) h.Node {
	const margin = 5
	plotW, plotH := width-2*margin, height-2*margin
	maxV := 0.0
	for _, v := range data {
		maxV = math.Max(maxV, v)
	}
	domainMax := niceMax(maxV)
	band := plotW / float64(len(data))
	offset := band * 0.1
	barW := math.Floor(band - 2*offset) // recharts truncates the bar size
	clipID := fmt.Sprintf("recharts%d-clip", rand.Uint32())
	var bars []any
	for i, v := range data {
		x := margin + float64(i)*band + offset
		y := margin + plotH*(1-v/domainMax)
		hgt := margin + plotH - y
		d := fmt.Sprintf("M %s,%s h %s v %s h %s Z", num(x), num(y), num(barW), num(hgt), num(-barW))
		bars = append(bars, h.G(h.Class("recharts-layer recharts-bar-rectangle"),
			h.G(h.Class("recharts-layer recharts-inactive-bar"),
				h.Path(h.Attr("name", "undefined"), h.Attr("x", num(x)), h.Attr("y", num(y)), h.Attr("width", num(barW)), h.Attr("height", num(hgt)),
					h.Attr("radius", "0"), h.Class("recharts-rectangle"), h.Attr("d", d), h.Style("fill: hsl(var(--foreground)); opacity: 0.9;")),
			),
		))
	}
	layer := func(z string) *h.Element { return h.G(h.TabIndex(-1), h.Class("recharts-zIndex-layer_"+z)) }
	w, ht := num(width), num(height)
	return h.Div(h.Class("recharts-responsive-container"), h.Style("width: 100%; height: 100%; min-width: 0px;"),
		h.Div(h.Style("width: 0px; height: 0px; overflow: visible;"),
			h.Div(h.Attr("width", w), h.Attr("height", ht), h.Class("recharts-wrapper"), h.Style("position: relative; cursor: default; width: "+w+"px; height: "+ht+"px;"),
				h.Svg(h.Role("application"), h.TabIndex(0), h.Class("recharts-surface"), h.Attr("width", w), h.Attr("height", ht), h.Attr("viewBox", "0 0 "+w+" "+ht),
					h.Style("width: 100%; height: 100%; display: block;"),
					h.Title(), h.El("desc"),
					layer("-100"), layer("-50"),
					h.Defs(h.ClipPath(h.ID(clipID), h.Rect(h.Attr("x", num(margin)), h.Attr("y", num(margin)), h.Attr("height", num(plotH)), h.Attr("width", num(plotW))))),
					layer("100"), layer("200"),
					h.G(h.TabIndex(-1), h.Class("recharts-zIndex-layer_300"),
						h.G(h.Class("recharts-layer recharts-bar"), h.ID(fmt.Sprintf("recharts-bar-%d", rand.Uint32())),
							h.G(h.Class("recharts-layer recharts-bar-rectangles"),
								h.G(append([]any{h.Class("recharts-layer")}, bars...)...),
							),
						),
					),
					layer("400"), layer("500"), layer("600"), layer("1000"), layer("1100"), layer("1200"), layer("2000"),
				),
			),
		),
	)
}

// niceMax is the upper bound of recharts' nice tick domain for [0, max]
// with 5 ticks.
func niceMax(max float64) float64 {
	if max <= 0 {
		return 1
	}
	step := max / 4
	mag := math.Pow(10, math.Floor(math.Log10(step)))
	norm := step / mag
	var nice float64
	switch {
	case norm <= 1:
		nice = 1
	case norm <= 2:
		nice = 2
	case norm <= 2.5:
		nice = 2.5
	case norm <= 5:
		nice = 5
	default:
		nice = 10
	}
	return math.Ceil(max/(nice*mag)) * nice * mag
}

// num formats like recharts: rounded to 4 decimals, JavaScript notation.
func num(v float64) string {
	r := math.Round(v*1e4) / 1e4
	if r == 0 {
		r = 0
	}
	return strconv.FormatFloat(r, 'f', -1, 64)
}

// DialogExample is drawer-dialog: a Dialog on screens at least 768px wide
// (the reference viewport), a Drawer below (useMediaQuery). The Drawer
// version waits in a template; media.js swaps the two when the media query
// changes, like the example's re-render.
func DialogExample() h.Node {
	return h.Group(
		dialog.Dialog(
			dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Edit Profile")),
			dialog.Content(h.Class("sm:max-w-[425px]"),
				dialog.Header(
					dialog.Title("Edit profile"),
					dialog.Description("Make changes to your profile here. Click save when you're done."),
				),
				profileForm(""),
			),
		),
		h.Template(h.Data("parity-ignore", ""), h.Data("example-media-swap", "(min-width: 768px)"),
			drawer.Drawer(
				drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Edit Profile")),
				drawer.Content(
					drawer.Header(h.Class("text-left"),
						drawer.Title("Edit profile"),
						drawer.Description("Make changes to your profile here. Click save when you're done."),
					),
					profileForm("px-4"),
					drawer.Footer(h.Class("pt-2"),
						drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
					),
				),
			),
		),
	)
}

func profileForm(className string) h.Node {
	return h.Form(h.Class(lib.Cn("grid items-start gap-6", className)),
		h.Div(h.Class("grid gap-3"),
			label.Label(h.For("email"), "Email"),
			input.Input(h.Type("email"), h.ID("email"), h.Value("shadcn@example.com")),
		),
		h.Div(h.Class("grid gap-3"),
			label.Label(h.For("username"), "Username"),
			input.Input(h.ID("username"), h.Value("@shadcn")),
		),
		button.Button(h.Type("submit"), "Save changes"),
	)
}

// Directions is the parity fixture drawer-directions.
func Directions() h.Node {
	var items []any
	for _, d := range []drawer.Direction{drawer.DirectionTop, drawer.DirectionRight, drawer.DirectionBottom, drawer.DirectionLeft} {
		items = append(items, drawer.Drawer(d,
			drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("capitalize"), string(d))),
			drawer.Content(
				drawer.Header(
					drawer.Title("Drawer from the "+string(d)),
					drawer.Description("Drag it back to its edge to close it."),
				),
				h.Div(h.Class("px-4 text-sm"), "Content"),
				drawer.Footer(
					button.Button("Submit"),
					drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
				),
			),
		))
	}
	items = append(items, drawer.Drawer(drawer.Modal(false),
		drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Non-modal")),
		drawer.Content(
			drawer.Header(
				drawer.Title("Non-modal drawer"),
				drawer.Description("The page stays interactive."),
			),
			drawer.Footer(
				drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Close")),
			),
		),
	))
	return h.Div(append([]any{h.Class("flex flex-wrap items-center gap-2")}, items...)...)
}
