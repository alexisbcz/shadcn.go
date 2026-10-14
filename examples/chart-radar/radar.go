// Package chartradar ports the upstream radar charts
// (charts/chart-radar-*.tsx).
package chartradar

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
		ex("chart-radar-default", Default),
		ex("chart-radar-dots", Dots),
		ex("chart-radar-grid-circle", GridCircle),
		ex("chart-radar-grid-circle-fill", GridCircleFill),
		ex("chart-radar-grid-circle-no-lines", GridCircleNoLines),
		ex("chart-radar-grid-custom", GridCustom),
		ex("chart-radar-grid-fill", GridFill),
		ex("chart-radar-grid-none", GridNone),
		ex("chart-radar-icons", Icons),
		ex("chart-radar-label-custom", LabelCustom),
		ex("chart-radar-legend", Legend),
		ex("chart-radar-lines-only", LinesOnly),
		ex("chart-radar-multiple", Multiple),
		ex("chart-radar-radius", RadiusAxis),
	)
}

func nonce(r *http.Request) chart.Nonce { return chart.Nonce(secure.Nonce(r)) }

var (
	desktopOnly = chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
	}
	desktopMobile = chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
	}
)

func desktop(values ...float64) []chart.Row {
	months := []string{"January", "February", "March", "April", "May", "June"}
	rows := make([]chart.Row, len(values))
	for i, v := range values {
		rows[i] = chart.Row{"month": months[i], "desktop": v}
	}
	return rows
}

func desktopAndMobile(values ...[2]float64) []chart.Row {
	months := []string{"January", "February", "March", "April", "May", "June"}
	rows := make([]chart.Row, len(values))
	for i, v := range values {
		rows[i] = chart.Row{"month": months[i], "desktop": v[0], "mobile": v[1]}
	}
	return rows
}

var (
	dataDefault    = desktop(186, 305, 237, 273, 209, 214)
	dataFill       = desktop(186, 285, 237, 203, 209, 264)
	dataNoLines    = desktop(186, 305, 237, 203, 209, 214)
	dataMultiple   = desktopAndMobile([2]float64{186, 80}, [2]float64{305, 200}, [2]float64{237, 120}, [2]float64{73, 190}, [2]float64{209, 130}, [2]float64{214, 140})
	dataLinesOnly  = desktopAndMobile([2]float64{186, 160}, [2]float64{185, 170}, [2]float64{207, 180}, [2]float64{173, 160}, [2]float64{160, 190}, [2]float64{174, 204})
	dotProps       = chart.Dot(chart.Props{"r": 4, "fillOpacity": 1})
	tooltipDefault = chart.Tooltip(chart.Cursor(false), chart.TooltipContent())
	tooltipNoLabel = chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.HideLabel()))
	tooltipLine    = chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line")))
)

// radarCard is the card of the radar charts. compact is the default frame
// (header pb-4, content pb-0); the legend charts use the other one.
func radarCard(title string, compactHeader, compactContent bool, c h.Node) h.Node {
	headerClass := "items-center"
	if compactHeader {
		headerClass = "items-center pb-4"
	}
	contentClass := ""
	footerClass := "flex-col gap-2 pt-4 text-sm"
	if compactContent {
		contentClass = "pb-0"
		footerClass = "flex-col gap-2 text-sm"
	}
	return card.Card(
		card.Header(h.Class(headerClass),
			card.Title(title),
			card.Description("Showing total visitors for the last 6 months"),
		),
		card.Content(h.Class(contentClass), c),
		card.Footer(h.Class(footerClass),
			h.Div(h.Class("flex items-center gap-2 leading-none font-medium"),
				"Trending up by 5.2% this month ", lucide.TrendingUp(h.Class("h-4 w-4")),
			),
			h.Div(h.Class("flex items-center gap-2 leading-none text-muted-foreground"), "January - June 2024"),
		),
	)
}

func container(r *http.Request, config chart.Config, c *chart.Node) h.Node {
	return chart.Container(config, nonce(r), h.Class("mx-auto aspect-square max-h-[250px]"), c)
}

func month() *chart.Node { return chart.PolarAngleAxis(chart.DataKey("month")) }

func desktopRadar(fillOpacity float64, extra ...any) *chart.Node {
	return chart.Radar(append([]any{chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.FillOpacity(fillOpacity)}, extra...)...)
}

func mobileRadar() *chart.Node {
	return chart.Radar(chart.DataKey("mobile"), chart.Fill("var(--color-mobile)"))
}

func Default(r *http.Request) h.Node {
	return radarCard("Radar Chart", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataDefault), tooltipDefault, month(), chart.PolarGrid(), desktopRadar(0.6))))
}

func Dots(r *http.Request) h.Node {
	return radarCard("Radar Chart - Dots", false, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataDefault), tooltipDefault, month(), chart.PolarGrid(), desktopRadar(0.6, dotProps))))
}

func GridCircle(r *http.Request) h.Node {
	return radarCard("Radar Chart - Grid Circle", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataDefault), tooltipNoLabel, chart.PolarGrid(chart.GridType("circle")), month(), desktopRadar(0.6, dotProps))))
}

func GridCircleFill(r *http.Request) h.Node {
	return radarCard("Radar Chart - Grid Circle Filled", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataFill), tooltipDefault,
			chart.PolarGrid(chart.ClassName("fill-(--color-desktop) opacity-20"), chart.GridType("circle")),
			month(), desktopRadar(0.5))))
}

func GridCircleNoLines(r *http.Request) h.Node {
	return radarCard("Radar Chart - Grid Circle - No lines", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataNoLines), tooltipNoLabel,
			chart.PolarGrid(chart.GridType("circle"), chart.RadialLines(false)),
			month(), desktopRadar(0.6, dotProps))))
}

func GridCustom(r *http.Request) h.Node {
	return radarCard("Radar Chart - Grid Custom", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataDefault), tooltipNoLabel,
			chart.PolarGrid(chart.RadialLines(false), chart.PolarRadius(90), chart.StrokeWidth(1)),
			month(), desktopRadar(0.6))))
}

func GridFill(r *http.Request) h.Node {
	return radarCard("Radar Chart - Grid Filled", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataFill), tooltipNoLabel,
			chart.PolarGrid(chart.ClassName("fill-(--color-desktop) opacity-20")),
			month(), desktopRadar(0.5))))
}

func GridNone(r *http.Request) h.Node {
	return radarCard("Radar Chart - Grid None", true, true, container(r, desktopOnly,
		chart.RadarChart(chart.Data(dataDefault), tooltipNoLabel, month(), desktopRadar(0.6, dotProps))))
}

func legendChart(r *http.Request, config chart.Config) h.Node {
	return container(r, config,
		chart.RadarChart(chart.Data(dataMultiple), chart.Margin{Top: -40, Bottom: -10},
			tooltipLine, month(), chart.PolarGrid(), desktopRadar(0.6), mobileRadar(),
			chart.Legend(chart.LegendContent(chart.ClassName("mt-8"))),
		))
}

func Icons(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)", Icon: lucide.ArrowDownFromLine},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)", Icon: lucide.ArrowUpFromLine},
	}
	return radarCard("Radar Chart - Icons", true, false, legendChart(r, config))
}

func Legend(r *http.Request) h.Node {
	return radarCard("Radar Chart - Legend", false, false, legendChart(r, desktopMobile))
}

func LabelCustom(r *http.Request) h.Node {
	return radarCard("Radar Chart - Custom Label", true, true, container(r, desktopMobile,
		chart.RadarChart(chart.Data(dataMultiple), chart.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
			tooltipLine,
			chart.PolarAngleAxis(chart.DataKey("month"), chart.Tick(chart.RadarStackedTick(dataMultiple, "desktop", "mobile", "month"))),
			chart.PolarGrid(), desktopRadar(0.6), mobileRadar(),
		)))
}

func LinesOnly(r *http.Request) h.Node {
	return radarCard("Radar Chart - Lines Only", true, true, container(r, desktopMobile,
		chart.RadarChart(chart.Data(dataLinesOnly), tooltipLine, month(), chart.PolarGrid(chart.RadialLines(false)),
			chart.Radar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), chart.FillOpacity(0), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2)),
			chart.Radar(chart.DataKey("mobile"), chart.Fill("var(--color-mobile)"), chart.FillOpacity(0), chart.Stroke("var(--color-mobile)"), chart.StrokeWidth(2)),
		)))
}

func Multiple(r *http.Request) h.Node {
	return radarCard("Radar Chart - Multiple", true, true, container(r, desktopMobile,
		chart.RadarChart(chart.Data(dataMultiple), tooltipLine, month(), chart.PolarGrid(), desktopRadar(0.6), mobileRadar())))
}

func RadiusAxis(r *http.Request) h.Node {
	return radarCard("Radar Chart - Radius Axis", true, true, container(r, desktopMobile,
		chart.RadarChart(chart.Data(dataMultiple),
			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"), chart.LabelKey("month"))),
			chart.PolarGrid(), desktopRadar(0.6), mobileRadar(),
			chart.PolarRadiusAxis(chart.Angle(60), chart.Stroke("var(--foreground)"), chart.Orientation("middle"), chart.AxisLine(false)),
		)))
}
