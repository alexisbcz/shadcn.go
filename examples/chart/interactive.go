package chart

import (
	_ "embed"
	"net/http"
	"strconv"

	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

// interactive.js switches the series of the interactive examples: the
// state upstream keeps in React (activeChart) lives in the buttons'
// data-active, and the chart is updated through its el.chart API.
//
//go:embed interactive.js
var interactiveScript string

func init() { runtime.Module("chart-example-interactive", interactiveScript) }

var interactiveConfig = func(desktop, mobile string) chart.Config {
	return chart.Config{
		{Key: "views", Label: "Page Views"},
		{Key: "desktop", Label: "Desktop", Color: desktop},
		{Key: "mobile", Label: "Mobile", Color: mobile},
	}
}

// localeString formats n like Number.prototype.toLocaleString("en-US").
func localeString(n int) string {
	s := strconv.Itoa(n)
	out := make([]byte, 0, len(s)+len(s)/3)
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, s[i])
	}
	return string(out)
}

func total(rows []chart.Row, key string) int {
	t := 0
	for _, r := range rows {
		t += r[key].(int)
	}
	return t
}

func seriesButtons(config chart.Config, buttonClass string) h.Node {
	buttons := make([]any, 0, 2)
	for _, key := range []string{"desktop", "mobile"} {
		item, _ := config.Get(key)
		buttons = append(buttons, h.Button(
			h.Data("active", strconv.FormatBool(key == "desktop")),
			h.Class(buttonClass),
			h.Span(h.Class("text-xs text-muted-foreground"), item.Label),
			h.Span(h.Class("text-lg leading-none font-bold sm:text-3xl"), localeString(total(dataBarInteractive, key))),
		))
	}
	return h.Div(h.Class("flex"), buttons)
}

func dateAxis() *chart.Node {
	return chart.XAxis(chart.DataKey("date"), chart.TickLine(false), chart.AxisLine(false), chart.TickMargin(8), chart.MinTickGap(32),
		chart.TickFormatter(chart.DateFormat("en-US", map[string]string{"month": "short", "day": "numeric"})))
}

func dateTooltip(cursor bool) *chart.Node {
	content := chart.TooltipContent(chart.ClassName("w-[150px]"), chart.NameKey("views"),
		chart.LabelFormatter(chart.DateFormat("en-US", map[string]string{"month": "short", "day": "numeric", "year": "numeric"})))
	if cursor {
		return chart.Tooltip(content)
	}
	return chart.Tooltip(chart.Cursor(false), content)
}

func BarInteractive(r *http.Request) h.Node {
	config := interactiveConfig("var(--chart-2)", "var(--chart-1)")
	return card.Card(h.Class("py-0"),
		card.Header(h.Class("flex flex-col items-stretch border-b p-0! sm:flex-row"),
			h.Div(h.Class("flex flex-1 flex-col justify-center gap-1 px-6 pt-4 pb-3 sm:py-0!"),
				card.Title("Bar Chart - Interactive"),
				card.Description("Showing total visitors for the last 3 months"),
			),
			seriesButtons(config, "relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l data-[active=true]:bg-muted/50 sm:border-t-0 sm:border-l sm:px-8 sm:py-6"),
		),
		card.Content(h.Class("px-2 sm:p-6"),
			chart.Container(config, nonce(r), h.Class("aspect-auto h-[250px] w-full"),
				chart.BarChart(chart.AccessibilityLayer(), chart.Data(dataBarInteractive), edgeMargin,
					noVertical(), dateAxis(), dateTooltip(true),
					chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)")),
				),
			),
		),
	)
}

func LineInteractive(r *http.Request) h.Node {
	config := interactiveConfig("var(--chart-1)", "var(--chart-2)")
	return card.Card(h.Class("py-4 sm:py-0"), card.Header(h.Class("flex flex-col items-stretch border-b p-0! sm:flex-row"),
		h.Div(h.Class("flex flex-1 flex-col justify-center gap-1 px-6 pb-3 sm:pb-0"),
			card.Title("Line Chart - Interactive"),
			card.Description("Showing total visitors for the last 3 months"),
		),
		seriesButtons(config, "flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l data-[active=true]:bg-muted/50 sm:border-t-0 sm:border-l sm:px-8 sm:py-6"),
	),
		card.Content(h.Class("px-2 sm:p-6"),
			chart.Container(config, nonce(r), h.Class("aspect-auto h-[250px] w-full"),
				chart.LineChart(chart.AccessibilityLayer(), chart.Data(dataLineInteractive), edgeMargin,
					noVertical(), dateAxis(), dateTooltip(true),
					chart.Line(chart.DataKey("desktop"), chart.Type("monotone"), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2), chart.Dot(false)),
				),
			),
		),
	)
}
