package chart

import (
	"net/http"

	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	selectui "shadcn.go/components/ui/select"
	h "shadcn.go/libhtml"
)

// AreaInteractive is charts/chart-area-interactive: the time range select
// filters the data. The initial range ("90d") keeps every entry; on change,
// interactive.js filters the data with upstream's date arithmetic and
// updates the chart through its el.chart API.
func AreaInteractive(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "visitors", Label: "Visitors"},
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
	}
	day := chart.DateFormat("en-US", map[string]string{"month": "short", "day": "numeric"})
	return card.Card(h.Class("pt-0"),
		card.Header(h.Class("flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row"),
			h.Div(h.Class("grid flex-1 gap-1"),
				card.Title("Area Chart - Interactive"),
				card.Description("Showing total visitors for the last 3 months"),
			),
			selectui.Select(selectui.DefaultValue("90d"),
				selectui.Trigger(h.Class("hidden w-[160px] rounded-lg sm:ml-auto sm:flex"), h.Aria("label", "Select a value"),
					selectui.Value(selectui.Placeholder("Last 3 months")),
				),
				selectui.Content(h.Class("rounded-xl"),
					selectui.Item(selectui.ItemValue("90d"), h.Class("rounded-lg"), "Last 3 months"),
					selectui.Item(selectui.ItemValue("30d"), h.Class("rounded-lg"), "Last 30 days"),
					selectui.Item(selectui.ItemValue("7d"), h.Class("rounded-lg"), "Last 7 days"),
				),
			),
		),
		card.Content(h.Class("px-2 pt-4 sm:px-6 sm:pt-6"),
			// React's useId gives the Select its id first, so the chart's
			// generated id is the second one: _r_1_.
			chart.Container(config, nonce(r), chart.ID("_r_1_"), h.Class("aspect-auto h-[250px] w-full"),
				chart.AreaChart(chart.Data(dataAreaInteractive),
					chart.Defs(gradient("fillDesktop", "desktop"), gradient("fillMobile", "mobile")),
					noVertical(),
					chart.XAxis(chart.DataKey("date"), chart.TickLine(false), chart.AxisLine(false), chart.TickMargin(8), chart.MinTickGap(32), chart.TickFormatter(day)),
					chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.LabelFormatter(day), chart.Indicator("dot"))),
					stackedAreas("url(#fillMobile)", "url(#fillDesktop)", false),
					chart.Legend(chart.LegendContent()),
				),
			),
		),
	)
}
