// Package chart ports the upstream chart examples (examples/chart-*.tsx)
// and the area, bar and line charts (charts/chart-{area,bar,line}-*.tsx).
package chart

import (
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

func init() {
	ex := func(name string, fn func(r *http.Request) h.Node) registry.Example {
		return registry.Example{Name: name, Component: "chart", Render: fn}
	}
	registry.Register(
		ex("chart-bar-demo", BarDemo),
		ex("chart-bar-demo-axis", BarDemoAxis),
		ex("chart-bar-demo-grid", BarDemoGrid),
		ex("chart-bar-demo-legend", BarDemoLegend),
		ex("chart-bar-demo-tooltip", BarDemoTooltip),
		ex("chart-tooltip-demo", func(*http.Request) h.Node { return TooltipDemo() }),
		ex("chart-area-axes", AreaAxes),
		ex("chart-area-default", AreaDefault),
		ex("chart-area-gradient", AreaGradient),
		ex("chart-area-icons", AreaIcons),
		ex("chart-area-legend", AreaLegend),
		ex("chart-area-linear", AreaLinear),
		ex("chart-area-stacked", AreaStacked),
		ex("chart-area-stacked-expand", AreaStackedExpand),
		ex("chart-area-step", AreaStep),
		ex("chart-bar-active", BarActive),
		ex("chart-bar-default", BarDefault),
		ex("chart-bar-horizontal", BarHorizontal),
		ex("chart-bar-interactive", BarInteractive),
		ex("chart-bar-label", BarLabel),
		ex("chart-bar-label-custom", BarLabelCustom),
		ex("chart-bar-mixed", BarMixed),
		ex("chart-bar-multiple", BarMultiple),
		ex("chart-bar-negative", BarNegative),
		ex("chart-bar-stacked", BarStacked),
		ex("chart-line-default", LineDefault),
		ex("chart-line-dots", LineDots),
		ex("chart-line-dots-colors", LineDotsColors),
		ex("chart-line-dots-custom", LineDotsCustom),
		ex("chart-line-interactive", LineInteractive),
		ex("chart-line-label", LineLabel),
		ex("chart-line-label-custom", LineLabelCustom),
		ex("chart-line-linear", LineLinear),
		ex("chart-line-multiple", LineMultiple),
		ex("chart-line-step", LineStep),
	)
}

func nonce(r *http.Request) chart.Nonce { return chart.Nonce(secure.Nonce(r)) }

// --- shared configs -------------------------------------------------------------

var (
	desktopMobileBlue = chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "#2563eb"},
		{Key: "mobile", Label: "Mobile", Color: "#60a5fa"},
	}
	desktopOnly = chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
	}
	desktopMobile = chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
	}
	browsers = chart.Config{
		{Key: "visitors", Label: "Visitors"},
		{Key: "chrome", Label: "Chrome", Color: "var(--chart-1)"},
		{Key: "safari", Label: "Safari", Color: "var(--chart-2)"},
		{Key: "firefox", Label: "Firefox", Color: "var(--chart-3)"},
		{Key: "edge", Label: "Edge", Color: "var(--chart-4)"},
		{Key: "other", Label: "Other", Color: "var(--chart-5)"},
	}
	browsersColored = chart.Config{
		{Key: "visitors", Label: "Visitors", Color: "var(--chart-2)"},
		{Key: "chrome", Label: "Chrome", Color: "var(--chart-1)"},
		{Key: "safari", Label: "Safari", Color: "var(--chart-2)"},
		{Key: "firefox", Label: "Firefox", Color: "var(--chart-3)"},
		{Key: "edge", Label: "Edge", Color: "var(--chart-4)"},
		{Key: "other", Label: "Other", Color: "var(--chart-5)"},
	}
)

// --- card frames ----------------------------------------------------------------

func trending() h.Node {
	return h.Group("Trending up by 5.2% this month ", lucide.TrendingUp(h.Class("h-4 w-4")))
}

// areaCard is the card of the area charts (and line-multiple): footer
// with the trend and a period.
func areaCard(title, description, period string, c h.Node) h.Node {
	return card.Card(
		card.Header(card.Title(title), card.Description(description)),
		card.Content(c),
		card.Footer(
			h.Div(h.Class("flex w-full items-start gap-2 text-sm"),
				h.Div(h.Class("grid gap-2"),
					h.Div(h.Class("flex items-center gap-2 leading-none font-medium"), trending()),
					h.Div(h.Class("flex items-center gap-2 leading-none text-muted-foreground"), period),
				),
			),
		),
	)
}

// barCard is the card of the bar and line charts.
func barCard(title string, c h.Node) h.Node {
	return card.Card(
		card.Header(card.Title(title), card.Description("January - June 2024")),
		card.Content(c),
		card.Footer(h.Class("flex-col items-start gap-2 text-sm"),
			h.Div(h.Class("flex gap-2 leading-none font-medium"), trending()),
			h.Div(h.Class("leading-none text-muted-foreground"), "Showing total visitors for the last 6 months"),
		),
	)
}

const sixMonths = "Showing total visitors for the last 6 months"

func monthAxis(margin float64) *chart.Node {
	return chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.AxisLine(false), chart.TickMargin(margin), chart.TickFormatter(chart.Slice(0, 3)))
}

func barMonthAxis() *chart.Node {
	return chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false), chart.TickFormatter(chart.Slice(0, 3)))
}

func noVertical() *chart.Node { return chart.CartesianGrid(chart.Vertical(false)) }

func hideLabelTooltip() *chart.Node {
	return chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.HideLabel()))
}

var edgeMargin = chart.Margin{Left: 12, Right: 12}

// --- examples/chart-bar-demo* -----------------------------------------------------------

func demoBars() []*chart.Node {
	return []*chart.Node{
		chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.Radius(4)),
		chart.Bar(chart.DataKey("mobile"), chart.Fill("var(--color-mobile)"), chart.Radius(4)),
	}
}

func BarDemo(r *http.Request) h.Node {
	return chart.Container(desktopMobileBlue, nonce(r), h.Class("min-h-[200px] w-full"),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarDemo), demoBars()),
	)
}

func BarDemoGrid(r *http.Request) h.Node {
	return chart.Container(desktopMobileBlue, nonce(r), h.Class("min-h-[200px] w-full"),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarDemo), noVertical(), demoBars()),
	)
}

func BarDemoAxis(r *http.Request) h.Node {
	return chart.Container(desktopMobileBlue, nonce(r), h.Class("min-h-[200px] w-full"),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarDemo), noVertical(), barMonthAxis(), demoBars()),
	)
}

func BarDemoTooltip(r *http.Request) h.Node {
	return chart.Container(desktopMobileBlue, nonce(r), h.Class("min-h-[200px] w-full"),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarDemo), noVertical(), barMonthAxis(),
			chart.Tooltip(chart.TooltipContent()),
			demoBars(),
		),
	)
}

func BarDemoLegend(r *http.Request) h.Node {
	return chart.Container(desktopMobileBlue, nonce(r), h.Class("min-h-[200px] w-full"),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarDemo), noVertical(), barMonthAxis(),
			chart.Tooltip(chart.TooltipContent()),
			chart.Legend(chart.LegendContent()),
			demoBars(),
		),
	)
}

// --- area charts ---------------------------------------------------------------------------

func AreaDefault(r *http.Request) h.Node {
	return areaCard("Area Chart", sixMonths, "January - June 2024",
		chart.Container(desktopOnly, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaDefault), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
				chart.Area(chart.DataKey("desktop"), chart.Type("natural"), chart.Fill("var(--color-desktop)"), chart.FillOpacity(0.4), chart.Stroke("var(--color-desktop)")),
			),
		))
}

func stackedAreas(fillMobile, fillDesktop string, opacity bool) []*chart.Node {
	area := func(key, fill string) *chart.Node {
		return chart.Area(chart.DataKey(key), chart.Type("natural"), chart.Fill(fill),
			func() any {
				if opacity {
					return chart.FillOpacity(0.4)
				}
				return nil
			}(),
			chart.Stroke("var(--color-"+key+")"), chart.StackID("a"))
	}
	return []*chart.Node{area("mobile", fillMobile), area("desktop", fillDesktop)}
}

func AreaAxes(r *http.Request) h.Node {
	return areaCard("Area Chart - Axes", sixMonths, "January - June 2024",
		chart.Container(desktopMobile, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaAxes), chart.Margin{Left: -20, Right: 12},
				noVertical(), monthAxis(8),
				chart.YAxis(chart.TickLine(false), chart.AxisLine(false), chart.TickMargin(8), chart.TickCount(3)),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent()),
				stackedAreas("var(--color-mobile)", "var(--color-desktop)", true),
			),
		))
}

func gradient(id, key string) *chart.Node {
	return chart.El("linearGradient", chart.P("id", id), chart.P("x1", "0"), chart.P("y1", "0"), chart.P("x2", "0"), chart.P("y2", "1"),
		chart.El("stop", chart.P("offset", "5%"), chart.P("stopColor", "var(--color-"+key+")"), chart.P("stopOpacity", 0.8)),
		chart.El("stop", chart.P("offset", "95%"), chart.P("stopColor", "var(--color-"+key+")"), chart.P("stopOpacity", 0.1)),
	)
}

func AreaGradient(r *http.Request) h.Node {
	return areaCard("Area Chart - Gradient", sixMonths, "January - June 2024",
		chart.Container(desktopMobile, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaGradient), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent()),
				chart.Defs(gradient("fillDesktop", "desktop"), gradient("fillMobile", "mobile")),
				stackedAreas("url(#fillMobile)", "url(#fillDesktop)", true),
			),
		))
}

func AreaIcons(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)", Icon: lucide.TrendingDown},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)", Icon: lucide.TrendingUp},
	}
	return areaCard("Area Chart - Icons", sixMonths, "January - June 2024",
		chart.Container(config, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaIcons), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
				stackedAreas("var(--color-mobile)", "var(--color-desktop)", true),
				chart.Legend(chart.LegendContent()),
			),
		))
}

func AreaLegend(r *http.Request) h.Node {
	return areaCard("Area Chart - Legend", sixMonths, "January - June 2024",
		chart.Container(desktopMobile, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaLegend), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
				stackedAreas("var(--color-mobile)", "var(--color-desktop)", true),
				chart.Legend(chart.LegendContent()),
			),
		))
}

func AreaLinear(r *http.Request) h.Node {
	return areaCard("Area Chart - Linear", sixMonths, "January - June 2024",
		chart.Container(desktopOnly, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaLinear), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("dot"), chart.HideLabel())),
				chart.Area(chart.DataKey("desktop"), chart.Type("linear"), chart.Fill("var(--color-desktop)"), chart.FillOpacity(0.4), chart.Stroke("var(--color-desktop)")),
			),
		))
}

func AreaStacked(r *http.Request) h.Node {
	return areaCard("Area Chart - Stacked", sixMonths, "January - June 2024",
		chart.Container(desktopMobile, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaStacked), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("dot"))),
				stackedAreas("var(--color-mobile)", "var(--color-desktop)", true),
			),
		))
}

func AreaStackedExpand(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
		{Key: "other", Label: "Other", Color: "var(--chart-3)"},
	}
	area := func(key string, opacity float64) *chart.Node {
		return chart.Area(chart.DataKey(key), chart.Type("natural"), chart.Fill("var(--color-"+key+")"), chart.FillOpacity(opacity), chart.Stroke("var(--color-"+key+")"), chart.StackID("a"))
	}
	return areaCard("Area Chart - Stacked Expanded", "Showing total visitors for the last 6months", "January - June 2024",
		chart.Container(config, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaStackedExpand), chart.Margin{Left: 12, Right: 12, Top: 12}, chart.StackOffset("expand"),
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
				area("other", 0.1), area("mobile", 0.4), area("desktop", 0.4),
			),
		))
}

func AreaStep(r *http.Request) h.Node {
	config := chart.Config{{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)", Icon: lucide.Activity}}
	return areaCard("Area Chart - Step", sixMonths, "January - June 2024",
		chart.Container(config, nonce(r),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(dataAreaStep), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.HideLabel())),
				chart.Area(chart.DataKey("desktop"), chart.Type("step"), chart.Fill("var(--color-desktop)"), chart.FillOpacity(0.4), chart.Stroke("var(--color-desktop)")),
			),
		))
}

// --- bar charts ---------------------------------------------------------------------------

func BarDefault(r *http.Request) h.Node {
	return barCard("Bar Chart", chart.Container(desktopOnly, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarDefault),
			noVertical(), barMonthAxis(), hideLabelTooltip(),
			chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.Radius(8)),
		),
	))
}

func BarActive(r *http.Request) h.Node {
	return barCard("Bar Chart - Active", chart.Container(browsers, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarActive),
			noVertical(),
			chart.XAxis(chart.DataKey("browser"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false), chart.TickFormatter(chart.ConfigLabel())),
			hideLabelTooltip(),
			chart.Bar(chart.DataKey("visitors"), chart.StrokeWidth(2), chart.Radius(8),
				chart.Shape(chart.ActiveIndexShape(2, chart.Props{"fillOpacity": 0.8, "stroke": chart.Field("fill"), "strokeDasharray": 4, "strokeDashoffset": 4})),
			),
		),
	))
}

func BarHorizontal(r *http.Request) h.Node {
	return barCard("Bar Chart - Horizontal", chart.Container(desktopOnly, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarHorizontal), chart.Layout("vertical"), chart.Margin{Left: -20},
			chart.XAxis(chart.Type("number"), chart.DataKey("desktop"), chart.Hide()),
			chart.YAxis(chart.DataKey("month"), chart.Type("category"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false), chart.TickFormatter(chart.Slice(0, 3))),
			hideLabelTooltip(),
			chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.Radius(5)),
		),
	))
}

func BarLabel(r *http.Request) h.Node {
	return barCard("Bar Chart - Label", chart.Container(desktopOnly, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarLabel), chart.Margin{Top: 20},
			noVertical(), barMonthAxis(), hideLabelTooltip(),
			chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.Radius(8),
				chart.LabelList(chart.Position("top"), chart.Offset(12), chart.ClassName("fill-foreground"), chart.FontSize(12)),
			),
		),
	))
}

func BarLabelCustom(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-2)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
		{Key: "label", Color: "var(--background)"},
	}
	return barCard("Bar Chart - Custom Label", chart.Container(config, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarLabelCustom), chart.Layout("vertical"), chart.Margin{Right: 16},
			chart.CartesianGrid(chart.Horizontal(false)),
			chart.YAxis(chart.DataKey("month"), chart.Type("category"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false), chart.TickFormatter(chart.Slice(0, 3)), chart.Hide()),
			chart.XAxis(chart.DataKey("desktop"), chart.Type("number"), chart.Hide()),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
			chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.Radius(4),
				chart.LabelList(chart.DataKey("month"), chart.Position("insideLeft"), chart.Offset(8), chart.ClassName("fill-(--color-label)"), chart.FontSize(12)),
				chart.LabelList(chart.DataKey("desktop"), chart.Position("right"), chart.Offset(8), chart.ClassName("fill-foreground"), chart.FontSize(12)),
			),
		),
	))
}

func BarMixed(r *http.Request) h.Node {
	return barCard("Bar Chart - Mixed", chart.Container(browsers, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarMixed), chart.Layout("vertical"), chart.Margin{},
			chart.YAxis(chart.DataKey("browser"), chart.Type("category"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false), chart.TickFormatter(chart.ConfigLabel())),
			chart.XAxis(chart.DataKey("visitors"), chart.Type("number"), chart.Hide()),
			hideLabelTooltip(),
			chart.Bar(chart.DataKey("visitors"), chart.Radius(5)),
		),
	))
}

func BarMultiple(r *http.Request) h.Node {
	return barCard("Bar Chart - Multiple", chart.Container(desktopMobile, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarMultiple),
			noVertical(), barMonthAxis(),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("dashed"))),
			chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.Radius(4)),
			chart.Bar(chart.DataKey("mobile"), chart.Fill("var(--color-mobile)"), chart.Radius(4)),
		),
	))
}

func BarNegative(r *http.Request) h.Node {
	config := chart.Config{{Key: "visitors", Label: "Visitors"}}
	cells := make([]*chart.Node, 0, len(dataBarNegative))
	for _, row := range dataBarNegative {
		fill := "var(--chart-2)"
		if row["visitors"].(int) > 0 {
			fill = "var(--chart-1)"
		}
		cells = append(cells, chart.Cell(chart.Fill(fill)))
	}
	return barCard("Bar Chart - Negative", chart.Container(config, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarNegative),
			noVertical(),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.HideLabel(), chart.HideIndicator())),
			chart.Bar(chart.DataKey("visitors"),
				chart.LabelList(chart.Position("top"), chart.DataKey("month"), chart.P("fillOpacity", 1)),
				cells,
			),
		),
	))
}

func BarStacked(r *http.Request) h.Node {
	return barCard("Bar Chart - Stacked + Legend", chart.Container(desktopMobile, nonce(r),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarStacked),
			noVertical(), barMonthAxis(),
			chart.Tooltip(chart.TooltipContent(chart.HideLabel())),
			chart.Legend(chart.LegendContent()),
			chart.Bar(chart.DataKey("desktop"), chart.StackID("a"), chart.Fill("var(--color-desktop)"), chart.Radii(0, 0, 4, 4)),
			chart.Bar(chart.DataKey("mobile"), chart.StackID("a"), chart.Fill("var(--color-mobile)"), chart.Radii(4, 4, 0, 0)),
		),
	))
}

// --- line charts -------------------------------------------------------------------------

func lineCard(title string, config chart.Config, r *http.Request, root *chart.Node) h.Node {
	return barCard(title, chart.Container(config, nonce(r), root))
}

func LineDefault(r *http.Request) h.Node {
	return lineCard("Line Chart", desktopOnly, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineDefault), edgeMargin,
			noVertical(), monthAxis(8), hideLabelTooltip(),
			chart.Line(chart.DataKey("desktop"), chart.Type("natural"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2), chart.Dot(false)),
		))
}

func LineLinear(r *http.Request) h.Node {
	return lineCard("Line Chart - Linear", desktopOnly, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineLinear), edgeMargin,
			noVertical(), monthAxis(8), hideLabelTooltip(),
			chart.Line(chart.DataKey("desktop"), chart.Type("linear"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2), chart.Dot(false)),
		))
}

func LineStep(r *http.Request) h.Node {
	return lineCard("Line Chart - Step", desktopOnly, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineStep), edgeMargin,
			noVertical(), monthAxis(8), hideLabelTooltip(),
			chart.Line(chart.DataKey("desktop"), chart.Type("step"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2), chart.Dot(false)),
		))
}

func LineDots(r *http.Request) h.Node {
	return lineCard("Line Chart - Dots", desktopMobile, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineDots), edgeMargin,
			noVertical(), monthAxis(8), hideLabelTooltip(),
			chart.Line(chart.DataKey("desktop"), chart.Type("natural"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2),
				chart.Dot(chart.Props{"fill": "var(--color-desktop)"}), chart.ActiveDot(chart.Props{"r": 6})),
		))
}

func LineDotsColors(r *http.Request) h.Node {
	return lineCard("Line Chart - Dots Colors", browsersColored, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineDotsColors), chart.Margin{Top: 24, Left: 24, Right: 24},
			noVertical(),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"), chart.NameKey("visitors"), chart.HideLabel())),
			chart.Line(chart.DataKey("visitors"), chart.Type("natural"), chart.Stroke("var(--color-visitors)"), chart.StrokeWidth(2),
				chart.Dot(chart.DotShape(chart.Props{"r": 5, "fill": chart.Field("fill"), "stroke": chart.Field("fill")}))),
		))
}

func LineDotsCustom(r *http.Request) h.Node {
	return lineCard("Line Chart - Custom Dots", desktopMobile, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineDotsCustom), edgeMargin,
			noVertical(), monthAxis(8), hideLabelTooltip(),
			chart.Line(chart.DataKey("desktop"), chart.Type("natural"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2),
				chart.Dot(chart.IconDot(lucide.GitCommitVertical, 24, chart.Props{"fill": "hsl(var(--background))", "stroke": "var(--color-desktop)"}))),
		))
}

func LineLabel(r *http.Request) h.Node {
	return lineCard("Line Chart - Label", desktopMobile, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineLabel), chart.Margin{Top: 20, Left: 12, Right: 12},
			noVertical(), monthAxis(8),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
			chart.Line(chart.DataKey("desktop"), chart.Type("natural"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2),
				chart.Dot(chart.Props{"fill": "var(--color-desktop)"}), chart.ActiveDot(chart.Props{"r": 6}),
				chart.LabelList(chart.Position("top"), chart.Offset(12), chart.ClassName("fill-foreground"), chart.FontSize(12)),
			),
		))
}

func LineLabelCustom(r *http.Request) h.Node {
	return lineCard("Line Chart - Custom Label", browsersColored, r,
		chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineLabelCustom), chart.Margin{Top: 24, Left: 24, Right: 24},
			noVertical(),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"), chart.NameKey("visitors"), chart.HideLabel())),
			chart.Line(chart.DataKey("visitors"), chart.Type("natural"), chart.Stroke("var(--color-visitors)"), chart.StrokeWidth(2),
				chart.Dot(chart.Props{"fill": "var(--color-visitors)"}), chart.ActiveDot(chart.Props{"r": 6}),
				chart.LabelList(chart.Position("top"), chart.Offset(12), chart.ClassName("fill-foreground"), chart.FontSize(12), chart.DataKey("browser"), chart.Formatter(chart.ConfigLabel())),
			),
		))
}

func LineMultiple(r *http.Request) h.Node {
	line := func(key string) *chart.Node {
		return chart.Line(chart.DataKey(key), chart.Type("monotone"), chart.Stroke("var(--color-"+key+")"), chart.StrokeWidth(2), chart.Dot(false))
	}
	return card.Card(
		card.Header(card.Title("Line Chart - Multiple"), card.Description("January - June 2024")),
		card.Content(chart.Container(desktopMobile, nonce(r),
			chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineMultiple), edgeMargin,
				noVertical(), monthAxis(8),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent()),
				line("desktop"), line("mobile"),
			),
		)),
		card.Footer(
			h.Div(h.Class("flex w-full items-start gap-2 text-sm"),
				h.Div(h.Class("grid gap-2"),
					h.Div(h.Class("flex items-center gap-2 leading-none font-medium"), trending()),
					h.Div(h.Class("flex items-center gap-2 leading-none text-muted-foreground"), sixMonths),
				),
			),
		),
	)
}
