package chart

import (
	"net/http"
	"strconv"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the chart page (ui/apps/v4/examples/radix/chart-*.tsx).
// Most are the new-york-v4 examples under another name; the radix-nova chart
// has the same classes as new-york-v4.
func init() {
	nova := func(name string, fn func(r *http.Request) h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "chart", Render: func(r *http.Request) h.Node {
			return style.Scope(style.Nova, fn(r))
		}}
	}
	registry.Register(
		nova("chart-example", BarDemo),
		nova("chart-example-axis", BarDemoAxis),
		nova("chart-example-grid", BarDemoGrid),
		nova("chart-example-legend", BarDemoLegend),
		nova("chart-example-tooltip", BarDemoTooltip),
		nova("chart-tooltip", func(*http.Request) h.Node { return TooltipDemo() }),
		nova("chart-rtl", NovaChartRtl),
		// chart-demo imports the new-york-v4 card and chart: no style scope.
		registry.Example{Name: "nova--chart-demo", Component: "chart", Render: NovaChartDemo},
	)
}

// NovaChartDemo is examples/radix/chart-demo.tsx: chart-bar-interactive
// with the April data only and a bottom padding on the card.
func NovaChartDemo(r *http.Request) h.Node {
	config := interactiveConfig("var(--chart-2)", "var(--chart-1)")
	data := dataBarInteractive[:30]
	buttons := make([]any, 0, 2)
	for _, key := range []string{"desktop", "mobile"} {
		item, _ := config.Get(key)
		buttons = append(buttons, h.Button(
			h.Data("active", strconv.FormatBool(key == "desktop")),
			h.Class("relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l data-[active=true]:bg-muted/50 sm:border-t-0 sm:border-l sm:px-8 sm:py-6"),
			h.Span(h.Class("text-xs text-muted-foreground"), item.Label),
			h.Span(h.Class("text-lg leading-none font-bold sm:text-3xl"), localeString(total(data, key))),
		))
	}
	return card.Card(h.Class("py-0 pb-4"),
		card.Header(h.Class("flex flex-col items-stretch border-b p-0! sm:flex-row"),
			h.Div(h.Class("flex flex-1 flex-col justify-center gap-1 px-6 pt-4 pb-3 sm:py-0!"),
				card.Title("Bar Chart - Interactive"),
				card.Description("Showing total visitors for the last 3 months"),
			),
			h.Div(h.Class("flex"), buttons),
		),
		card.Content(h.Class("px-2 sm:p-6"),
			chart.Container(config, nonce(r), h.Class("aspect-auto h-[250px] w-full"),
				chart.BarChart(chart.AccessibilityLayer(), chart.Data(data), edgeMargin,
					noVertical(), dateAxis(), dateTooltip(true),
					chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)")),
				),
			),
		),
	)
}

// rtlData is the chart-rtl data with the months already translated to
// Arabic (the preview's default language): the axis shows their first three
// letters and the tooltip the full name, as t[value] does upstream.
var rtlData = []chart.Row{
	{"month": "يناير", "desktop": 186, "mobile": 80},
	{"month": "فبراير", "desktop": 305, "mobile": 200},
	{"month": "مارس", "desktop": 237, "mobile": 120},
	{"month": "أبريل", "desktop": 73, "mobile": 190},
	{"month": "مايو", "desktop": 209, "mobile": 130},
	{"month": "يونيو", "desktop": 214, "mobile": 140},
}

// NovaChartRtl is examples/radix/chart-rtl.tsx in its default language
// (Arabic, right to left).
func NovaChartRtl(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "سطح المكتب", Color: "var(--chart-2)"},
		{Key: "mobile", Label: "الجوال", Color: "var(--chart-1)"},
	}
	return chart.Container(config, nonce(r), h.Class("min-h-[200px] w-full"),
		chart.BarChart(chart.AccessibilityLayer(), chart.Data(rtlData),
			chart.CartesianGrid(chart.Vertical(false), chart.P("orientation", "right")),
			chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false),
				chart.TickFormatter(chart.Slice(0, 3)), chart.P("reversed", true)),
			chart.Tooltip(chart.TooltipContent(chart.LabelClassName("w-32"))),
			chart.Legend(chart.LegendContent()),
			demoBars(),
		),
	)
}
