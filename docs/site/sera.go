package site

import (
	"net/http"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/components/ui/table"
	togglegroup "shadcn.go/components/ui/toggle-group"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

// The Sera style page: app/(app)/(styles)/sera. Its components are
// styles/base-sera (rendered here with the radix components under the Sera
// style). Its style.css (the .preview and .theme-* rules) is a stylesheet
// of its own in the reference build, vendored as static/sera.css (without
// the stray `@utility font-heading`, which the browser drops).

const seraTitle = "Introducing Sera"
const seraDescription = "Minimal. Editorial. Typographic. Underline Controls and Uppercase Headings. Shaped by Print Design Principles."

// seraPreviewClass is the class of each preview section (and of the lazy
// previews' placeholders, with p-4 sm:p-6 instead).
const seraPreviewClass = "preview theme-taupe @container/preview w-full flex-1 bg-muted pt-4 font-sans ring-1 ring-foreground/5 [--gap:--spacing(4)] sm:pt-0 md:[--gap:--spacing(6)] xl:[--gap:--spacing(8)] 2xl:py-8 **:[.container]:px-(--gap)"

// Sera serves /sera.
func (s *Site) Sera(r *http.Request) (h.Node, error) {
	m := Meta{Title: seraTitle, Description: seraDescription, Canonical: ""}
	return s.Document(r, m,
		s.AppLayout(r,
			h.Link(h.Rel("stylesheet"), h.Href(staticPrefix+"sera.css")),
			h.Script(h.Type("module"), h.Src(staticPrefix+"sera.js")),
			pageHeader(
				h.H1(h.Class(lib.Cn("leading-tighter max-w-3xl text-3xl font-semibold tracking-tight text-balance text-primary lg:leading-[1.1] lg:font-semibold xl:text-5xl xl:tracking-tighter", "font-(family-name:--font-playfair-display) text-[2.875rem] tracking-tight!")), seraTitle),
				h.P(h.Class(lib.Cn("max-w-4xl text-base text-balance text-foreground sm:text-lg", "max-w-2xl text-pretty md:text-balance")), seraDescription),
				h.Div(h.Class("flex w-full items-center justify-center gap-2 pt-2 **:data-[slot=button]:shadow-none **:[.container]:justify-start"),
					style.Scope(style.Sera,
						button.Button(lib.AsChild, button.SizeSm,
							h.A(h.Href("/create?preset=b4xFeBLg4O"), "Open in shadcn/create", lucide.ArrowRight(h.Data("icon", "inline-end")))),
					),
				),
			),
			seraImagePreview(),
			h.Div(h.Class("container-wrapper hidden flex-1 flex-col section-soft px-0 md:flex md:px-2 md:py-12"),
				h.Div(h.Class("container flex flex-1 flex-col gap-10 px-0 3xl:max-w-[2000px] md:px-6"),
					style.Scope(style.Sera, seraAudienceAnalytics(r)),
					seraLazyPreview("articleDirectory", 760),
					seraLazyPreview("emptyState", 560),
					seraLazyPreview("editArticle", 980),
					seraLazyPreview("mediaLibrary", 880),
					seraLazyPreview("mediaLibraryTable", 980),
				),
			),
		),
	), nil
}

// seraImagePreview is components/image-preview.tsx (below md).
func seraImagePreview() h.Node {
	item := func(name string) h.Node {
		img := func(theme, class string) h.Node {
			u := "%2Fimages%2F" + name + "-" + theme + ".png"
			return h.Img(h.Alt(name), h.Attr("loading", "lazy"), h.Width("1440"), h.Height("900"),
				h.Attr("decoding", "async"), h.Data("nimg", "1"), h.Class(class), h.Style("color:transparent"),
				h.UnsafeAttr("srcset", "/_next/image?url="+u+"&w=1920&q=75 1x, /_next/image?url="+u+"&w=3840&q=75 2x"),
				h.UnsafeAttr("src", "/_next/image?url="+u+"&w=3840&q=75"))
		}
		return h.Div(h.Class("theme-taupe overflow-hidden bg-muted px-4 py-2 first:pt-4 last:pb-4"),
			img("light", "dark:hidden"),
			img("dark", "hidden dark:block"),
		)
	}
	return h.Div(h.Class("mt-8 flex flex-col overflow-hidden md:hidden"),
		item("sera-01"), item("sera-03"), item("sera-02"), item("sera-06"))
}

// seraLazyPreview is components/lazy-preview.tsx before the preview loads:
// its placeholder.
func seraLazyPreview(name string, minHeight int) h.Node {
	return h.Div(h.Data("sera-preview", name),
		h.El("template", style.Scope(style.Sera, seraPreviews()[name]())),
		h.Div(h.Aria("hidden", "true"),
			h.Class("preview theme-taupe @container/preview w-full flex-1 bg-muted p-4 font-sans ring-1 ring-foreground/5 [--gap:--spacing(4)] sm:p-6 md:[--gap:--spacing(6)] xl:[--gap:--spacing(8)]"),
			h.Style("min-height: "+strconv.Itoa(minHeight)+"px"),
			h.Div(h.Class("container flex flex-col gap-(--gap) py-(--gap)"),
				h.Div(h.Class("flex items-center justify-between gap-4"),
					h.Div(h.Class("flex flex-col gap-3"),
						h.Div(h.Class("h-5 w-44 bg-background/80")),
						h.Div(h.Class("h-3 w-56 max-w-full bg-background/60")),
					),
					h.Div(h.Class("hidden h-8 w-28 bg-background/70 sm:block")),
				),
				h.Div(h.Class("grid grid-cols-1 gap-(--gap) md:grid-cols-3"),
					h.Div(h.Class("min-h-48 bg-background/70 md:col-span-2")),
					h.Div(h.Class("min-h-48 bg-background/70")),
				),
			),
		),
	)
}

// --- audience-analytics ---------------------------------------------------------

func seraAudienceAnalytics(r *http.Request) h.Node {
	return h.Div(h.Class(seraPreviewClass),
		seraAnalyticsHeader(),
		separator.Separator(h.Class("hidden sm:block")),
		h.Div(h.Class("container grid grid-cols-12 gap-(--gap) py-(--gap)"),
			seraMetricsGrid(),
			h.Div(h.Class("col-span-full md:col-span-6 lg:col-span-8"), seraTrafficOverview(r)),
			seraDemographics(h.Class("col-span-full md:col-span-6 lg:col-span-4")),
			seraTopEditorial(h.Class("col-span-full")),
		),
	)
}

var seraDateRanges = [][2]string{
	{"Last 7 days", "last-7-days"},
	{"Last 30 days", "last-30-days"},
	{"This month", "this-month"},
	{"Last month", "last-month"},
}

func seraAnalyticsHeader() h.Node {
	var items []any
	for _, o := range seraDateRanges {
		items = append(items, dm.RadioItem(dm.Value(o[1]), o[0]))
	}
	return h.Header(
		h.Div(h.Class("container flex flex-col items-center justify-center gap-(--gap) py-(--gap) sm:flex-row sm:justify-between"),
			h.Div(h.Class("flex flex-col gap-2 text-center sm:text-left"),
				h.H1(h.Class("line-clamp-1 font-heading text-3xl tracking-wide uppercase md:text-3xl lg:text-4xl"), "Audience Analytics"),
				h.Div(h.Class("line-clamp-1 text-sm font-medium tracking-wider text-muted-foreground uppercase"), "Editorial Performance Dashboard"),
			),
			buttongroup.ButtonGroup(h.Class("gap-2 sm:ml-auto md:gap-4"),
				dm.DropdownMenu(h.Data("sera-date-range", ""),
					dm.Trigger(lib.AsChild,
						button.Button(button.VariantOutline, h.Class("bg-background hover:bg-background/80 data-popup-open:bg-background"),
							h.Span(h.Data("sera-date-label", ""), h.Class("contents"), "Last 30 days"), " ",
							lucide.ChevronDown(h.Data("icon", "inline-end"))),
					),
					dm.Content(dm.AlignEnd,
						dm.Group(dm.RadioGroup(dm.Value("last-30-days"), h.Group(items...))),
					),
				),
				button.Button(
					lucide.Download(h.Data("icon", "inline-start")),
					h.Span(h.Class("lg:hidden"), "Export"),
					h.Span(h.Class("hidden lg:inline"), "Export Report"),
				),
			),
		),
	)
}

func seraMetricsGrid() h.Node {
	metrics := []struct {
		label, value, comparison string
		up                       bool
	}{
		{"Total visitors", "248.5k", "12.4%", true},
		{"Unique readers", "182.1k", "8.7%", true},
		{"Avg. time on page", "3m 42s", "1.2%", false},
		{"Bounce rate", "42.8%", "3.5%", false},
	}
	var cards []any
	for _, m := range metrics {
		icon := lucide.TrendingUp(h.Class("inline-block size-2.5 text-muted-foreground"))
		if !m.up {
			icon = lucide.TrendingDown(h.Class("inline-block size-2.5 text-muted-foreground"))
		}
		cards = append(cards, card.Card(h.Class("gap-0 col-span-full md:col-span-6 lg:col-span-3"),
			card.Content(h.Class("flex flex-col gap-2"),
				card.Description(h.Class("text-xs uppercase"), m.label),
				card.Title(h.Class("text-5xl tracking-tight lowercase"), m.value),
				card.Description(icon, " ", h.Span(h.Class("text-foreground"), m.comparison), " ", h.Span("vs last period")),
			),
		))
	}
	return h.Group(cards...)
}

var seraTrafficData = []chart.Row{
	{"date": "2025-10-01", "views": 2600, "unique": 1600},
	{"date": "2025-10-04", "views": 4500, "unique": 3000},
	{"date": "2025-10-08", "views": 3500, "unique": 2500},
	{"date": "2025-10-10", "views": 6400, "unique": 4500},
	{"date": "2025-10-13", "views": 5400, "unique": 4000},
	{"date": "2025-10-15", "views": 8300, "unique": 6500},
	{"date": "2025-10-17", "views": 7400, "unique": 6000},
	{"date": "2025-10-18", "views": 9240, "unique": 7105},
	{"date": "2025-10-22", "views": 7700, "unique": 6400},
	{"date": "2025-10-26", "views": 8800, "unique": 7000},
	{"date": "2025-10-29", "views": 9800, "unique": 8400},
}

func seraTrafficOverview(r *http.Request) h.Node {
	config := chart.Config{
		{Key: "views", Label: "Views", Theme: &chart.Theme{Light: "var(--chart-5)", Dark: "var(--chart-1)"}},
		{Key: "unique", Label: "Unique", Theme: &chart.Theme{Light: "var(--chart-1)", Dark: "var(--chart-2)"}},
	}
	return card.Card(
		card.Header(
			card.Title(h.Class("text-2xl"), "Traffic Overview"),
			card.Description("Traffic for the last 30 days has increased by 12.4% compared to the previous period."),
		),
		card.Content(
			chart.Container(config, chart.Nonce(secure.Nonce(r)), h.Class("h-82 w-full"),
				chart.LineChart(chart.Data(seraTrafficData),
					chart.CartesianGrid(chart.Vertical(false), chart.P("strokeDasharray", "3 6"), chart.Stroke("var(--border)")),
					chart.XAxis(chart.DataKey("date"), chart.TickLine(false), chart.AxisLine(false), chart.P("interval", "preserveStartEnd"),
						chart.TickMargin(10), chart.TickFormatter(chart.DateFormat("en-US", map[string]string{"month": "short", "day": "numeric", "timeZone": "UTC"}))),
					chart.YAxis(chart.TickLine(false), chart.AxisLine(false), chart.P("width", 44), chart.P("domain", []int{0, 10000}),
						chart.P("ticks", []int{0, 2500, 5000, 7500, 10000}), chart.Hide()),
					chart.Tooltip(chart.TooltipContent()),
					chart.Line(chart.Type("linear"), chart.DataKey("views"), chart.Stroke("var(--color-views)"), chart.StrokeWidth(2.2),
						chart.Dot(false), chart.ActiveDot(map[string]any{"r": 3.5, "fill": "var(--color-views)"})),
					chart.Line(chart.Type("linear"), chart.DataKey("unique"), chart.Stroke("var(--color-unique)"), chart.StrokeWidth(2),
						chart.P("strokeDasharray", "4 6"), chart.Dot(false), chart.ActiveDot(false)),
					chart.El("ReferenceDot", chart.P("x", "2025-10-18"), chart.P("y", 9240), chart.P("r", 2.5),
						chart.Fill("var(--color-views)"), chart.Stroke("var(--color-views)")),
				),
			),
		),
	)
}

func seraDemographics(props ...any) h.Node {
	data := []struct {
		age string
		pct int
	}{{"18 - 24", 22}, {"25 - 34", 64}, {"35 - 44", 12}, {"45+", 5}}
	var bars []any
	for i, d := range data {
		bars = append(bars, seraProgress("sera-demographics-"+strconv.Itoa(i), d.age, d.pct))
	}
	return card.Card(h.Group(props...),
		card.Header(
			card.Title(h.Class("text-2xl"), "Demographics"),
			card.Description("Reader Profile"),
		),
		card.Content(h.Class("flex flex-col gap-10"), h.Group(bars...)),
		card.Footer(
			button.Button(button.VariantLink, h.Class("w-full"), "View all source ", lucide.MoveRight(h.Data("icon", "inline-end"))),
		),
	)
}

// seraProgress is the base-sera Progress with a label and a value (Base
// UI's progress: the value is formatted as a percentage).
func seraProgress(id, label string, value int) h.Node {
	v := strconv.Itoa(value)
	return h.Div(h.Aria("label", label), h.Aria("labelledby", id), h.Aria("valuemax", "100"), h.Aria("valuemin", "0"),
		h.Aria("valuenow", v), h.Aria("valuetext", v+"%"), h.Class("flex flex-wrap gap-3"), h.Data("progressing", ""),
		h.Data("slot", "progress"), h.Role("progressbar"),
		h.Span(h.Class("text-xs font-semibold tracking-wide uppercase"), h.Data("progressing", ""), h.Data("slot", "progress-label"), h.ID(id), h.Role("presentation"), label),
		h.Span(h.Aria("hidden", "true"), h.Class("ml-auto text-sm text-muted-foreground tabular-nums"), h.Data("progressing", ""), h.Data("slot", "progress-value"), v+"%"),
		h.Div(h.Class("relative flex h-0.5 w-full items-center overflow-x-hidden rounded-none bg-muted"), h.Data("progressing", ""), h.Data("slot", "progress-track"),
			h.Div(h.Class("h-full bg-primary transition-all"), h.Data("progressing", ""), h.Data("slot", "progress-indicator"),
				h.Style("height: inherit; inset-inline-start: 0px; width: "+v+"%")),
		),
	)
}

var seraEditorialRows = [][5]string{
	{"The New Vanguard of Minimalist Architecture", "Elena Rostova", "Oct 12", "45.2k", "04:15"},
	{"Autumn Sartorial Code: Deconstructed Classics", "Julian Vance", "Oct 05", "38.9k", "03:42"},
	{"Interview: Director Sofia Coppola on The Aesthetics of Isolation", "Marcus Trent", "Sep 28", "31.4k", "06:20"},
	{"Sourcing Ceramics from Kyoto's Oldest Kilns", "Sarah Lin", "Oct 18", "22.1k", "02:55"},
	{"Field Notes from Copenhagen Design Week", "Noah Bennett", "Oct 21", "19.7k", "03:18"},
	{"A Studio Visit with Milan's Most Elusive Lighting Designer", "Claire Duval", "Oct 09", "17.4k", "04:02"},
	{"Collecting the New Avant-Garde in Contemporary Furniture", "Tommy Rhodes", "Sep 30", "15.9k", "03:36"},
	{"Inside Lisbon's Quiet Culinary Renaissance", "Amara Iqbal", "Oct 14", "14.2k", "05:08"},
	{"Why Slow Interiors Are Defining the Next Luxury Wave", "Henry Vale", "Oct 03", "12.7k", "03:11"},
	{"The Return of Print: Independent Magazine Covers to Watch", "Mina Okafor", "Sep 26", "11.3k", "02:49"},
}

// seraTopEditorial is top-editorial.tsx: five rows, the rest revealed by
// "Load more content" after two seconds (static/sera.js).
func seraTopEditorial(props ...any) h.Node {
	var rows []any
	for i, row := range seraEditorialRows {
		cells := []any{
			table.Cell(h.Class("translate-y-1 align-text-top"), strconv.Itoa(i+1)),
			table.Cell(
				h.Div(h.Class("flex flex-col gap-2"),
					h.P(h.Class("font-heading text-xl tracking-tight text-foreground"), row[0]),
					h.P(h.Class("text-xs font-semibold tracking-wide text-muted-foreground uppercase"), "By "+row[1]),
				),
			),
			table.Cell(row[2]),
			table.Cell(row[3]),
			table.Cell(row[4]),
			table.Cell(h.Class("text-right"),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild,
						button.Button(button.VariantGhost, button.SizeIconXs, h.Aria("label", "Open actions for "+row[0]), lucide.Ellipsis())),
					dm.Content(dm.AlignEnd,
						dm.Item("Edit"),
						dm.Item("Publish"),
						dm.Item(dm.VariantDestructive, "Delete"),
					),
				),
			),
		}
		// Until the hidden rows show, the last shown row is the last one.
		cells = append(cells, h.Class("[&:not(:has(~tr:not([hidden])))]:border-0"))
		if i >= 5 {
			cells = append(cells, h.Bool("hidden"), h.Data("sera-more", ""))
		}
		rows = append(rows, table.Row(cells...))
	}
	group := togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("views"), togglegroup.VariantOutline,
		h.Aria("label", "Top editorials metric selector"), h.Class("w-full sm:ml-auto sm:w-fit"))
	for _, m := range [][2]string{{"views", "VIEWS"}, {"time", "TIME"}, {"shares", "SHARES"}} {
		group.Children(togglegroup.Item(togglegroup.Value(m[0]), h.Class("flex-1"), m[1]))
	}
	return card.Card(h.Group(props...),
		card.Header(
			h.Div(h.Class("flex flex-col gap-(--gap) sm:flex-row"),
				h.Div(h.Class("flex flex-col gap-1.5"),
					card.Title(h.Class("text-2xl"), "Top Editorials"),
					card.Description("Ranked by engagement"),
				),
				group,
			),
		),
		card.Content(h.Class("flex-1 **:data-[slot=table-container]:no-scrollbar **:data-[slot=table-container]:overflow-y-hidden"),
			table.Table(
				table.Header(table.Row(
					table.Head("#"), table.Head("Title"), table.Head("Published"), table.Head("Page Views"), table.Head("Read Time"),
					table.Head(h.Class("text-right"), "Actions"),
				)),
				table.Body(h.Group(rows...)),
			),
		),
		card.Footer(h.Class("justify-center"),
			button.Button(h.Type("button"), button.VariantOutline, h.Data("sera-load-more", ""),
				"Load more content ", lucide.ArrowDown(h.Data("icon", "inline-end"))),
			h.El("template", h.Data("sera-spinner", ""), spinner.Spinner(h.Data("slot", "spinner"), h.Data("icon", "inline-end"))),
		),
	)
}
