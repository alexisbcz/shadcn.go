// Package chartpie ports the upstream pie charts (charts/chart-pie-*.tsx).
package chartpie

import (
	_ "embed"
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/components/ui/popper"
	"shadcn.go/components/ui/runtime"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

// interactive.js keeps chart-pie-interactive's active month (React state
// upstream) in the select and updates the chart's shape and label.
//
//go:embed interactive.js
var interactiveScript string

func init() {
	runtime.Module("chart-pie-example-interactive", interactiveScript)
	ex := func(name string, fn func(r *http.Request) h.Node) registry.Example {
		return registry.Example{Name: name, Component: "chart", Render: fn}
	}
	registry.Register(
		ex("chart-pie-donut", Donut),
		ex("chart-pie-donut-active", DonutActive),
		ex("chart-pie-donut-text", DonutText),
		ex("chart-pie-interactive", Interactive),
		ex("chart-pie-label", Label),
		ex("chart-pie-label-custom", LabelCustom),
		ex("chart-pie-label-list", LabelList),
		ex("chart-pie-legend", Legend),
		ex("chart-pie-separator-none", SeparatorNone),
		ex("chart-pie-simple", Simple),
		ex("chart-pie-stacked", Stacked),
	)
}

func nonce(r *http.Request) chart.Nonce { return chart.Nonce(secure.Nonce(r)) }

var browsers = chart.Config{
	{Key: "visitors", Label: "Visitors"},
	{Key: "chrome", Label: "Chrome", Color: "var(--chart-1)"},
	{Key: "safari", Label: "Safari", Color: "var(--chart-2)"},
	{Key: "firefox", Label: "Firefox", Color: "var(--chart-3)"},
	{Key: "edge", Label: "Edge", Color: "var(--chart-4)"},
	{Key: "other", Label: "Other", Color: "var(--chart-5)"},
}

var months = chart.Config{
	{Key: "visitors", Label: "Visitors"},
	{Key: "desktop", Label: "Desktop"},
	{Key: "mobile", Label: "Mobile"},
	{Key: "january", Label: "January", Color: "var(--chart-1)"},
	{Key: "february", Label: "February", Color: "var(--chart-2)"},
	{Key: "march", Label: "March", Color: "var(--chart-3)"},
	{Key: "april", Label: "April", Color: "var(--chart-4)"},
	{Key: "may", Label: "May", Color: "var(--chart-5)"},
}

func browserData(values ...int) []chart.Row {
	keys := []string{"chrome", "safari", "firefox", "edge", "other"}
	rows := make([]chart.Row, len(keys))
	for i, k := range keys {
		rows[i] = chart.Row{"browser": k, "visitors": values[i], "fill": "var(--color-" + k + ")"}
	}
	return rows
}

func monthData(key string, values ...int) []chart.Row {
	keys := []string{"january", "february", "march", "april", "may"}
	rows := make([]chart.Row, len(keys))
	for i, k := range keys {
		rows[i] = chart.Row{"month": k, key: values[i], "fill": "var(--color-" + k + ")"}
	}
	return rows
}

var (
	defaultBrowsers = func() []chart.Row { return browserData(275, 200, 187, 173, 90) }
	desktopData     = func() []chart.Row { return monthData("desktop", 186, 305, 237, 173, 209) }
	mobileData      = func() []chart.Row { return monthData("mobile", 80, 200, 120, 190, 130) }
)

// pieCard is the card of the pie charts: centered header, the chart, and
// the trend footer (unless footer is false).
func pieCard(title string, footer bool, c h.Node) h.Node {
	var f h.Node
	if footer {
		f = card.Footer(h.Class("flex-col gap-2 text-sm"),
			h.Div(h.Class("flex items-center gap-2 leading-none font-medium"),
				"Trending up by 5.2% this month ", lucide.TrendingUp(h.Class("h-4 w-4"))),
			h.Div(h.Class("leading-none text-muted-foreground"), "Showing total visitors for the last 6 months"),
		)
	}
	return card.Card(h.Class("flex flex-col"),
		card.Header(h.Class("items-center pb-0"),
			card.Title(title),
			card.Description("January - June 2024"),
		),
		card.Content(h.Class("flex-1 pb-0"), c),
		f,
	)
}

func container(r *http.Request, config chart.Config, class string, root *chart.Node) h.Node {
	return chart.Container(config, nonce(r), h.Class(class), root)
}

func hideLabelTooltip(cursor bool) *chart.Node {
	if cursor {
		return chart.Tooltip(chart.TooltipContent(chart.HideLabel()))
	}
	return chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.HideLabel()))
}

const square = "mx-auto aspect-square max-h-[250px]"

func Simple(r *http.Request) h.Node {
	return pieCard("Pie Chart", true, container(r, browsers, square,
		chart.PieChart(
			hideLabelTooltip(false),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"), chart.NameKey("browser")),
		)))
}

func Donut(r *http.Request) h.Node {
	return pieCard("Pie Chart - Donut", true, container(r, browsers, square,
		chart.PieChart(
			hideLabelTooltip(false),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"), chart.NameKey("browser"), chart.InnerRadius(60)),
		)))
}

func DonutActive(r *http.Request) h.Node {
	return pieCard("Pie Chart - Donut Active", true, container(r, browsers, square,
		chart.PieChart(
			hideLabelTooltip(false),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"), chart.NameKey("browser"),
				chart.InnerRadius(60), chart.StrokeWidth(5), chart.Shape(chart.ActiveSector(0, 10))),
		)))
}

func DonutText(r *http.Request) h.Node {
	data := browserData(275, 200, 287, 173, 190)
	total := 0
	for _, row := range data {
		total += row["visitors"].(int)
	}
	return pieCard("Pie Chart - Donut with Text", true, container(r, browsers, square,
		chart.PieChart(
			hideLabelTooltip(false),
			chart.Pie(chart.Data(data), chart.DataKey("visitors"), chart.NameKey("browser"),
				chart.InnerRadius(60), chart.StrokeWidth(5),
				chart.Label(chart.Content(chart.PieCenterText(total, "Visitors"))),
			),
		)))
}

func Label(r *http.Request) h.Node {
	return pieCard("Pie Chart - Label", true, container(r, browsers, square+" pb-0 [&_.recharts-pie-label-text]:fill-foreground",
		chart.PieChart(
			hideLabelTooltip(true),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"), chart.LabelProp(true), chart.NameKey("browser")),
		)))
}

func LabelCustom(r *http.Request) h.Node {
	return pieCard("Pie Chart - Custom Label", true, container(r, browsers, square+" px-0",
		chart.PieChart(
			chart.Tooltip(chart.TooltipContent(chart.NameKey("visitors"), chart.HideLabel())),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"), chart.LabelLine(false),
				chart.LabelProp(chart.PieValueLabel("visitors", "var(--foreground)")), chart.NameKey("browser")),
		)))
}

func LabelList(r *http.Request) h.Node {
	return pieCard("Pie Chart - Label List", true, container(r, browsers, square+" [&_.recharts-text]:fill-background",
		chart.PieChart(
			chart.Tooltip(chart.TooltipContent(chart.NameKey("visitors"), chart.HideLabel())),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"),
				chart.LabelList(chart.DataKey("browser"), chart.ClassName("fill-background"), chart.Stroke("none"), chart.FontSize(12), chart.Formatter(chart.ConfigLabel())),
			),
		)))
}

func Legend(r *http.Request) h.Node {
	return pieCard("Pie Chart - Legend", false, container(r, browsers, "mx-auto aspect-square max-h-[300px]",
		chart.PieChart(
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors")),
			chart.Legend(chart.LegendContent(chart.NameKey("browser"), chart.ClassName("-translate-y-2 flex-wrap gap-2 *:basis-1/4 *:justify-center"))),
		)))
}

func SeparatorNone(r *http.Request) h.Node {
	return pieCard("Pie Chart - Separator None", true, container(r, browsers, square,
		chart.PieChart(
			hideLabelTooltip(false),
			chart.Pie(chart.Data(defaultBrowsers()), chart.DataKey("visitors"), chart.NameKey("browser"), chart.Stroke("0")),
		)))
}

func Stacked(r *http.Request) h.Node {
	return pieCard("Pie Chart - Stacked", true, container(r, months, square,
		chart.PieChart(
			chart.Tooltip(chart.TooltipContent(chart.LabelKey("visitors"), chart.NameKey("month"), chart.Indicator("line"), chart.LabelFormatter(chart.PayloadConfigLabel()))),
			chart.Pie(chart.Data(desktopData()), chart.DataKey("desktop"), chart.OuterRadius(60)),
			chart.Pie(chart.Data(mobileData()), chart.DataKey("mobile"), chart.InnerRadius(70), chart.OuterRadius(90)),
		)))
}

// Interactive is charts/chart-pie-interactive: the month select sets the
// active sector and the center label (interactive.js).
func Interactive(r *http.Request) h.Node {
	const id = "pie-interactive"
	data := desktopData()
	items := make([]any, 0, len(data))
	for _, row := range data {
		key := row["month"].(string)
		cfg, _ := months.Get(key)
		items = append(items, selectui.Item(selectui.ItemValue(key), h.Class("rounded-lg [&_span]:flex"),
			h.Div(h.Class("flex items-center gap-2 text-xs"),
				h.Span(h.Class("flex h-3 w-3 shrink-0 rounded-xs"), h.Style("background-color:var(--color-"+key+")")),
				cfg.Label,
			),
		))
	}
	return card.Card(h.Data("chart", id), h.Class("flex flex-col"),
		chart.Style(id, months, nonce(r)),
		card.Header(h.Class("flex-row items-start space-y-0 pb-0"),
			h.Div(h.Class("grid gap-1"),
				card.Title("Pie Chart - Interactive"),
				card.Description("January - June 2024"),
			),
			selectui.Select(selectui.DefaultValue("january"),
				selectui.Trigger(h.Class("ml-auto h-7 w-[130px] rounded-lg pl-2.5"), h.Aria("label", "Select a value"),
					selectui.Value(selectui.Placeholder("Select month")),
				),
				selectui.Content(popper.AlignEnd, h.Class("rounded-xl"), items),
			),
		),
		card.Content(h.Class("flex flex-1 justify-center pb-0"),
			chart.Container(months, nonce(r), chart.ID(id), h.Class("mx-auto aspect-square w-full max-w-[300px]"),
				chart.PieChart(
					hideLabelTooltip(false),
					chart.Pie(chart.Data(data), chart.DataKey("desktop"), chart.NameKey("month"),
						chart.InnerRadius(60), chart.StrokeWidth(5), chart.Shape(chart.ActiveSectorRing(0)),
						chart.Label(chart.Content(chart.PieCenterText(data[0]["desktop"], "Visitors"))),
					),
				),
			),
		),
	)
}

