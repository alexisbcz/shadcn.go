package chart

import (
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Radial bar charts (charts/chart-radial-*.tsx). The recharts elements are
// built with chart.El (RadialBarChart, RadialBar, PolarGrid,
// PolarRadiusAxis, Label), rendered by the chart-radial runtime module.

func init() {
	ex := func(name string, fn func(r *http.Request) h.Node) registry.Example {
		return registry.Example{Name: name, Component: "chart", Render: fn}
	}
	registry.Register(
		ex("chart-radial-grid", RadialGrid),
		ex("chart-radial-label", RadialLabel),
		ex("chart-radial-shape", RadialShape),
		ex("chart-radial-simple", RadialSimple),
		ex("chart-radial-stacked", RadialStacked),
		ex("chart-radial-text", RadialText),
	)
}

var dataRadialBrowsers = []chart.Row{
	{"browser": "chrome", "visitors": 275, "fill": "var(--color-chrome)"},
	{"browser": "safari", "visitors": 200, "fill": "var(--color-safari)"},
	{"browser": "firefox", "visitors": 187, "fill": "var(--color-firefox)"},
	{"browser": "edge", "visitors": 173, "fill": "var(--color-edge)"},
	{"browser": "other", "visitors": 90, "fill": "var(--color-other)"},
}

var radialSafariConfig = chart.Config{
	{Key: "visitors", Label: "Visitors"},
	{Key: "safari", Label: "Safari", Color: "var(--chart-2)"},
}

// radialCard is the card of the radial charts; contentClass is the
// CardContent's class.
func radialCard(title, contentClass string, c h.Node) h.Node {
	return card.Card(h.Class("flex flex-col"),
		card.Header(h.Class("items-center pb-0"), card.Title(title), card.Description("January - June 2024")),
		card.Content(h.Class(contentClass), c),
		card.Footer(h.Class("flex-col gap-2 text-sm"),
			h.Div(h.Class("flex items-center gap-2 leading-none font-medium"), "Trending up by 5.2% this month ", lucide.TrendingUp(h.Class("h-4 w-4"))),
			h.Div(h.Class("leading-none text-muted-foreground"), sixMonths),
		),
	)
}

func radialContainer(r *http.Request, config chart.Config, class string, root *chart.Node) h.Node {
	return chart.Container(config, nonce(r), h.Class(class), root)
}

func radialTooltip(content ...any) *chart.Node {
	return chart.Tooltip(chart.Cursor(false), chart.TooltipContent(content...))
}

// centerLabel is PolarRadiusAxis (no tick, line or axis line) with a Label
// whose content renders a <text> at the center with the given lines.
func centerLabel(textProps chart.Props, lines ...chart.Props) *chart.Node {
	ls := make([]any, len(lines))
	for i, l := range lines {
		ls[i] = map[string]any(l)
	}
	return chart.El("PolarRadiusAxis", chart.P("tick", false), chart.TickLine(false), chart.AxisLine(false),
		chart.El("Label", chart.P("content", chart.Fn{Name: "centerText", Args: []any{map[string]any(textProps), ls}})),
	)
}

func RadialSimple(r *http.Request) h.Node {
	return radialCard("Radial Chart", "flex-1 pb-0", radialContainer(r, browsers, "mx-auto aspect-square max-h-[250px]",
		chart.El("RadialBarChart", chart.Data(dataRadialBrowsers), chart.P("innerRadius", 30), chart.P("outerRadius", 110),
			radialTooltip(chart.HideLabel(), chart.NameKey("browser")),
			chart.El("RadialBar", chart.DataKey("visitors"), chart.P("background", true)),
		)))
}

func RadialGrid(r *http.Request) h.Node {
	return radialCard("Radial Chart - Grid", "flex-1 pb-0", radialContainer(r, browsers, "mx-auto aspect-square max-h-[250px]",
		chart.El("RadialBarChart", chart.Data(dataRadialBrowsers), chart.P("innerRadius", 30), chart.P("outerRadius", 100),
			radialTooltip(chart.HideLabel(), chart.NameKey("browser")),
			chart.El("PolarGrid", chart.P("gridType", "circle")),
			chart.El("RadialBar", chart.DataKey("visitors")),
		)))
}

func RadialLabel(r *http.Request) h.Node {
	return radialCard("Radial Chart - Label", "flex-1 pb-0", radialContainer(r, browsers, "mx-auto aspect-square max-h-[250px]",
		chart.El("RadialBarChart", chart.Data(dataRadialBrowsers), chart.P("startAngle", -90), chart.P("endAngle", 380), chart.P("innerRadius", 30), chart.P("outerRadius", 110),
			radialTooltip(chart.HideLabel(), chart.NameKey("browser")),
			chart.El("RadialBar", chart.DataKey("visitors"), chart.P("background", true),
				chart.LabelList(chart.Position("insideStart"), chart.DataKey("browser"), chart.ClassName("fill-white capitalize mix-blend-luminosity"), chart.FontSize(11)),
			),
		)))
}

func RadialShape(r *http.Request) h.Node {
	return radialCard("Radial Chart - Shape", "flex-1 pb-0", radialContainer(r, radialSafariConfig, "mx-auto aspect-square max-h-[250px]",
		chart.El("RadialBarChart", chart.Data([]chart.Row{{"browser": "safari", "visitors": 1260, "fill": "var(--color-safari)"}}),
			chart.P("endAngle", 100), chart.P("innerRadius", 65), chart.P("outerRadius", 95),
			chart.El("PolarGrid", chart.P("gridType", "circle"), chart.P("radialLines", false), chart.Stroke("none"),
				chart.ClassName("first:fill-muted last:fill-background"), chart.P("polarRadius", []float64{86, 74})),
			chart.El("RadialBar", chart.DataKey("visitors"), chart.P("background", true)),
			centerLabel(chart.Props{"textAnchor": "middle", "dominantBaseline": "middle"},
				chart.Props{"dy": 0, "className": "fill-foreground text-4xl font-bold", "text": "1,260"},
				chart.Props{"dy": 24, "className": "fill-muted-foreground", "text": "Visitors"}),
		)))
}

func RadialText(r *http.Request) h.Node {
	return radialCard("Radial Chart - Text", "flex-1 pb-0", radialContainer(r, radialSafariConfig, "mx-auto aspect-square max-h-[250px]",
		chart.El("RadialBarChart", chart.Data([]chart.Row{{"browser": "safari", "visitors": 200, "fill": "var(--color-safari)"}}),
			chart.P("startAngle", 0), chart.P("endAngle", 250), chart.P("outerRadius", 90), chart.P("innerRadius", 80),
			chart.El("PolarGrid", chart.P("gridType", "circle"), chart.P("radialLines", false), chart.Stroke("none"),
				chart.ClassName("first:fill-muted last:fill-background"), chart.P("polarRadius", []float64{90, 80})),
			chart.El("RadialBar", chart.DataKey("visitors"), chart.P("background", true), chart.P("cornerRadius", 10)),
			centerLabel(chart.Props{"textAnchor": "middle", "dominantBaseline": "middle"},
				chart.Props{"dy": 0, "className": "fill-foreground text-4xl font-bold", "text": "200"},
				chart.Props{"dy": 24, "className": "fill-muted-foreground", "text": "Visitors"}),
		)))
}

func RadialStacked(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
	}
	return radialCard("Radial Chart - Stacked", "flex flex-1 items-center pb-0", radialContainer(r, config, "mx-auto aspect-square w-full max-w-[250px]",
		chart.El("RadialBarChart", chart.Data([]chart.Row{{"month": "january", "mobile": 570, "desktop": 1260}}),
			chart.P("endAngle", 180), chart.P("innerRadius", 80), chart.P("outerRadius", 110),
			chart.El("RadialBar", chart.DataKey("mobile"), chart.Fill("var(--color-mobile)"), chart.StackID("a"), chart.P("cornerRadius", 5), chart.ClassName("stroke-transparent stroke-2")),
			chart.El("RadialBar", chart.DataKey("desktop"), chart.StackID("a"), chart.P("cornerRadius", 5), chart.Fill("var(--color-desktop)"), chart.ClassName("stroke-transparent stroke-2")),
			radialTooltip(chart.HideLabel()),
			centerLabel(chart.Props{"textAnchor": "middle"},
				chart.Props{"dy": -16, "className": "fill-foreground text-2xl font-bold", "text": "1,830"},
				chart.Props{"dy": 4, "className": "fill-muted-foreground", "text": "Visitors"}),
		)))
}
