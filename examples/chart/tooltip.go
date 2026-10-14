package chart

import (
	_ "embed"
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// tooltip.js registers the tooltip formatters of chart-tooltip-formatter
// and chart-tooltip-advanced (inline JSX closures upstream).
//
//go:embed tooltip.js
var tooltipScript string

func init() {
	runtime.Module("chart-example-tooltip", tooltipScript)
	ex := func(name string, fn func(r *http.Request) h.Node) registry.Example {
		return registry.Example{Name: name, Component: "chart", Render: fn}
	}
	registry.Register(
		ex("chart-tooltip-advanced", TooltipAdvanced),
		ex("chart-tooltip-default", TooltipDefault),
		ex("chart-tooltip-formatter", TooltipFormatter),
		ex("chart-tooltip-icons", TooltipIcons),
		ex("chart-tooltip-indicator-line", TooltipIndicatorLine),
		ex("chart-tooltip-indicator-none", TooltipIndicatorNone),
		ex("chart-tooltip-label-custom", TooltipLabelCustom),
		ex("chart-tooltip-label-formatter", TooltipLabelFormatter),
		ex("chart-tooltip-label-none", TooltipLabelNone),
	)
}

var dataActivities = []chart.Row{
	{"date": "2024-07-15", "running": 450, "swimming": 300},
	{"date": "2024-07-16", "running": 380, "swimming": 420},
	{"date": "2024-07-17", "running": 520, "swimming": 120},
	{"date": "2024-07-18", "running": 140, "swimming": 550},
	{"date": "2024-07-19", "running": 600, "swimming": 350},
	{"date": "2024-07-20", "running": 480, "swimming": 400},
}

var activitiesConfig = chart.Config{
	{Key: "running", Label: "Running", Color: "var(--chart-1)"},
	{Key: "swimming", Label: "Swimming", Color: "var(--chart-2)"},
}

// tooltipChart is the card and stacked bar chart every tooltip example
// shares; content is the ChartTooltipContent.
func tooltipChart(r *http.Request, title, description string, config chart.Config, content *chart.Node) h.Node {
	return card.Card(
		card.Header(card.Title(title), card.Description(description)),
		card.Content(
			chart.Container(config, nonce(r),
				chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataActivities),
					chart.XAxis(chart.DataKey("date"), chart.TickLine(false), chart.TickMargin(10), chart.AxisLine(false),
						chart.TickFormatter(chart.DateFormat("en-US", map[string]string{"weekday": "short"}))),
					chart.Bar(chart.DataKey("running"), chart.StackID("a"), chart.Fill("var(--color-running)"), chart.Radii(0, 0, 4, 4)),
					chart.Bar(chart.DataKey("swimming"), chart.StackID("a"), chart.Fill("var(--color-swimming)"), chart.Radii(4, 4, 0, 0)),
					chart.Tooltip(content, chart.Cursor(false), chart.P("defaultIndex", 1)),
				),
			),
		),
	)
}

func TooltipDefault(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - Default", "Default tooltip with ChartTooltipContent.", activitiesConfig, chart.TooltipContent())
}

func TooltipIndicatorLine(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - Line Indicator", "Tooltip with line indicator.", activitiesConfig, chart.TooltipContent(chart.Indicator("line")))
}

func TooltipIndicatorNone(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - No Indicator", "Tooltip with no indicator.", activitiesConfig, chart.TooltipContent(chart.HideIndicator()))
}

func TooltipLabelNone(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - No Label", "Tooltip with no label.", activitiesConfig, chart.TooltipContent(chart.HideIndicator(), chart.HideLabel()))
}

func TooltipLabelCustom(r *http.Request) h.Node {
	config := append(chart.Config{{Key: "activities", Label: "Activities"}}, activitiesConfig...)
	return tooltipChart(r, "Tooltip - Custom label", "Tooltip with custom label from chartConfig.", config, chart.TooltipContent(chart.LabelKey("activities"), chart.Indicator("line")))
}

func TooltipLabelFormatter(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - Label Formatter", "Tooltip with label formatter.", activitiesConfig,
		chart.TooltipContent(chart.LabelFormatter(chart.DateFormat("en-US", map[string]string{"day": "numeric", "month": "long", "year": "numeric"}))))
}

func TooltipFormatter(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - Formatter", "Tooltip with custom formatter .", activitiesConfig,
		chart.TooltipContent(chart.HideLabel(), chart.Formatter(chart.Fn{Name: "kcalFormatter"})))
}

func TooltipIcons(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "running", Label: "Running", Color: "var(--chart-1)", Icon: lucide.Footprints},
		{Key: "swimming", Label: "Swimming", Color: "var(--chart-2)", Icon: lucide.Waves},
	}
	return tooltipChart(r, "Tooltip - Icons", "Tooltip with icons.", config, chart.TooltipContent(chart.HideLabel()))
}

func TooltipAdvanced(r *http.Request) h.Node {
	return tooltipChart(r, "Tooltip - Advanced", "Tooltip with custom formatter and total.", activitiesConfig,
		chart.TooltipContent(chart.HideLabel(), chart.ClassName("w-[180px]"), chart.Formatter(chart.Fn{Name: "kcalAdvanced"})))
}
