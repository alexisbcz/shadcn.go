package site

import (
	"net/http"
	"slices"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sheet"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/docs/content"
	"shadcn.go/docs/site/mdx"
	h "shadcn.go/libhtml"
)

// The charts gallery: app/(app)/charts/{layout.tsx,charts.tsx,[type]/page.tsx},
// components/{chart-display,chart-toolbar,chart-iframe,chart-code-viewer,
// chart-copy-button,charts-nav}.tsx.

const chartsTitle = "Beautiful Charts & Graphs"
const chartsDescription = "A collection of ready-to-use chart components built with Recharts. From basic charts to rich data displays, copy and paste into your apps."

type chartEntry struct {
	id        string
	fullWidth bool
}

// chartList is app/(app)/charts/charts.tsx.
var chartList = map[string][]chartEntry{
	"area":    chartEntries("chart-area-interactive!", "chart-area-default", "chart-area-linear", "chart-area-step", "chart-area-legend", "chart-area-stacked", "chart-area-stacked-expand", "chart-area-icons", "chart-area-gradient", "chart-area-axes"),
	"bar":     chartEntries("chart-bar-interactive!", "chart-bar-default", "chart-bar-horizontal", "chart-bar-multiple", "chart-bar-stacked", "chart-bar-label", "chart-bar-label-custom", "chart-bar-mixed", "chart-bar-active", "chart-bar-negative"),
	"line":    chartEntries("chart-line-interactive!", "chart-line-default", "chart-line-linear", "chart-line-step", "chart-line-multiple", "chart-line-dots", "chart-line-dots-custom", "chart-line-dots-colors", "chart-line-label", "chart-line-label-custom"),
	"pie":     chartEntries("chart-pie-simple", "chart-pie-separator-none", "chart-pie-label", "chart-pie-label-custom", "chart-pie-label-list", "chart-pie-legend", "chart-pie-donut", "chart-pie-donut-active", "chart-pie-donut-text", "chart-pie-stacked", "chart-pie-interactive"),
	"radar":   chartEntries("chart-radar-default", "chart-radar-dots", "chart-radar-lines-only", "chart-radar-label-custom", "chart-radar-grid-custom", "chart-radar-grid-none", "chart-radar-grid-circle", "chart-radar-grid-circle-no-lines", "chart-radar-grid-circle-fill", "chart-radar-grid-fill", "chart-radar-multiple", "chart-radar-legend", "chart-radar-icons", "chart-radar-radius"),
	"radial":  chartEntries("chart-radial-simple", "chart-radial-label", "chart-radial-grid", "chart-radial-text", "chart-radial-shape", "chart-radial-stacked"),
	"tooltip": chartEntries("chart-tooltip-default", "chart-tooltip-indicator-line", "chart-tooltip-indicator-none", "chart-tooltip-label-custom", "chart-tooltip-label-formatter", "chart-tooltip-label-none", "chart-tooltip-formatter", "chart-tooltip-icons", "chart-tooltip-advanced"),
}

var chartTypes = []string{"area", "bar", "line", "pie", "radar", "radial", "tooltip"}

// chartEntries builds a list; a trailing "!" marks fullWidth.
func chartEntries(ids ...string) []chartEntry {
	out := make([]chartEntry, len(ids))
	for i, id := range ids {
		name, full := strings.CutSuffix(id, "!")
		out[i] = chartEntry{id: name, fullWidth: full}
	}
	return out
}

// Charts serves /charts/{type}.
func (s *Site) Charts(r *http.Request) (h.Node, error) {
	typ := r.PathValue("type")
	if !slices.Contains(chartTypes, typ) {
		return nil, errNotFound
	}
	title := strings.ToUpper(typ[:1]) + typ[1:] + " Charts"
	cells := make([]any, 0, 12)
	list := chartList[typ]
	for i := range 12 {
		if i < len(list) {
			if it, ok := s.registryItem(list[i].id); ok && len(it.Files) > 0 {
				cells = append(cells, chartDisplay(it, list[i].fullWidth))
				continue
			}
		}
		cells = append(cells, h.Div(h.Class("hidden aspect-square w-full rounded-lg border border-dashed xl:block")))
	}
	page := h.Div(h.Class("grid flex-1 gap-12 lg:gap-24"),
		h.H2(h.Class("sr-only"), title),
		h.Div(h.Class("grid flex-1 scroll-mt-20 items-stretch gap-10 md:grid-cols-2 md:gap-6 lg:grid-cols-3 xl:gap-10"), h.Group(cells...)),
	)
	return s.Document(r, Meta{Title: title, PlainTitle: true, Description: chartsDescription, Canonical: "/charts/" + typ},
		s.AppLayout(r,
			pageHeader(
				announcement(),
				pageHeaderHeading(chartsTitle),
				pageHeaderDescription(chartsDescription),
				// The charts layout's Button is styles/radix-luma/ui/button.
				pageActions(style.Scope(style.Luma,
					button.Button(lib.AsChild, h.Class("h-[35px]"), h.A(h.Href("#charts"), "Browse Charts")),
					button.Button(lib.AsChild, button.VariantSecondary, h.A(h.Href("/docs/components/chart"), "Documentation")),
				)),
			),
			pageNav(h.ID("charts"), chartsNav(r.URL.Path)),
			h.Div(h.Class("container-wrapper flex-1"),
				h.Div(h.Class("container pb-6"), h.Section(h.Class("theme-container"), page)),
			),
		),
	), nil
}

func chartsNav(pathname string) h.Node {
	links := []struct{ name, href string }{
		{"Area Charts", "/charts/area#charts"},
		{"Bar Charts", "/charts/bar#charts"},
		{"Line Charts", "/charts/line#charts"},
		{"Pie Charts", "/charts/pie#charts"},
		{"Radar Charts", "/charts/radar#charts"},
		{"Radial Charts", "/charts/radial#charts"},
		{"Tooltips", "/charts/tooltip#charts"},
	}
	var as []any
	for _, l := range links {
		as = append(as, h.A(h.Href(l.href), h.Data("active", boolString(strings.HasPrefix(l.href, pathname))),
			h.Class("flex h-7 shrink-0 items-center justify-center px-4 text-center text-base font-medium text-muted-foreground transition-colors hover:text-primary data-[active=true]:text-primary"),
			l.name))
	}
	return h.Div(h.Class("relative overflow-hidden"),
		h.Div(h.Data("slot", "scroll-area"), h.Class("relative max-w-[600px] lg:max-w-none"),
			h.Div(h.Data("slot", "scroll-area-viewport"), h.Class("size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1"), h.Style("overflow: scroll"),
				h.Div(h.Style("min-width: 100%; display: table"),
					h.Div(h.Class("flex items-center"), h.Group(as...)),
				),
			),
		),
	)
}

func chartDisplay(it *registryItem, fullWidth bool) h.Node {
	code := it.Files[0].Content
	return h.Div(h.Class(lib.Cn("themes-wrapper group relative flex flex-col overflow-hidden rounded-xl transition-all duration-200 ease-in-out hover:z-30", h.Ternary(fullWidth, "md:col-span-2 lg:col-span-3", ""))),
		h.Div(h.Class("relative z-20 flex items-center gap-2 justify-end px-3 py-2.5"),
			h.Div(h.Class("flex items-center gap-1.5 pl-1 text-[13px] text-muted-foreground [&>svg]:h-[0.9rem] [&>svg]:w-[0.9rem]"), chartTitle(it.Name)),
			h.Div(h.Class("ml-auto flex items-center gap-2 [&>form]:flex"),
				chartCopyButton(code, "[&_svg]-h-3 h-6 w-6 rounded-[6px] bg-transparent text-foreground shadow-none hover:bg-muted dark:text-foreground [&_svg]:w-3"),
				separator.Separator(separator.Vertical, h.Class("mx-0 hidden h-4! md:flex")),
				chartCodeViewer(it, code),
			),
		),
		h.Div(h.Class("relative z-10 overflow-hidden rounded-xl bg-background"),
			h.El("iframe", h.Src("/view/"+blocksStyle+"/"+it.Name), h.Data("chart-iframe", ""),
				h.Class("w-full border-none transition-opacity duration-300 opacity-0"),
				h.Height("460"), h.Loading("lazy"), h.TitleAttr(it.Name)),
		),
	)
}

func chartTitle(name string) h.Node {
	switch {
	case strings.Contains(name, "chart-line"):
		return h.Group(lucide.ChartLine(), " Line Chart")
	case strings.Contains(name, "chart-bar"):
		return h.Group(lucide.ChartColumnBig(), " Bar Chart")
	case strings.Contains(name, "chart-pie"):
		return h.Group(lucide.ChartPie(), " Pie Chart")
	case strings.Contains(name, "chart-area"):
		return h.Group(lucide.ChartArea(), " Area Chart")
	case strings.Contains(name, "chart-radar"):
		return h.Group(lucide.Hexagon(), " Radar Chart")
	case strings.Contains(name, "chart-radial"):
		return h.Group(lucide.Radar(), " Radial Chart")
	case strings.Contains(name, "chart-tooltip"):
		return h.Group(lucide.MousePointer2(), "Tooltip")
	}
	return h.Text(name)
}

// chartCopyButton is components/chart-copy-button.tsx.
func chartCopyButton(code, class string) h.Node {
	return tooltip.Tooltip(
		tooltip.Trigger(lib.AsChild,
			button.Button(button.SizeIcon, button.VariantGhost, h.Class(lib.Cn("[&_svg]-h-3.5 h-7 w-7 rounded-[6px] [&_svg]:w-3.5", class)),
				h.Data("copy-text", code),
				h.Span(h.Class("sr-only"), "Copy"), tabler.IconCopy()),
		),
		tooltip.Content(h.Class("bg-black text-white"), "Copy code"),
	)
}

// chartCodeViewer is the desktop (Sheet) branch of
// components/chart-code-viewer.tsx.
func chartCodeViewer(it *registryItem, code string) h.Node {
	content := h.Div(h.Class("flex min-h-0 flex-1 flex-col gap-0"),
		h.Div(h.Class("chart-wrapper hidden theme-container sm:block [&_[data-chart]]:mx-auto [&_[data-chart]]:max-h-[35vh] [&>div]:rounded-none [&>div]:border-0 [&>div]:border-b [&>div]:shadow-none")),
		h.Div(h.Class("flex min-w-0 flex-1 flex-col overflow-hidden p-4"),
			h.El("figure", h.Data("rehype-pretty-code-figure", ""), h.Class("mt-0 flex h-auto min-w-0 flex-1 flex-col overflow-hidden"),
				h.El("figcaption", h.Data("language", "tsx"),
					h.Class("flex h-12 shrink-0 items-center gap-2 border-b py-2 pr-2 pl-4 text-foreground [&>svg]:size-4 [&>svg]:text-foreground [&>svg]:opacity-70"),
					content.LanguageIcon("tsx"), it.Name,
					h.Div(h.Class("ml-auto flex items-center gap-2"),
						chartCopyButton(code, ""),
						openInV0ButtonClass(it.Name, "rounded-sm"),
					),
				),
				h.Div(h.Class("no-scrollbar overflow-y-auto"), mdx.HighlightCode(code, "tsx")),
			),
		),
	)
	return sheet.Sheet(
		sheet.Trigger(lib.AsChild,
			button.Button(button.SizeSm, button.VariantOutline, h.Class("h-6 rounded-[6px] border bg-transparent px-2 text-xs text-foreground shadow-none hover:bg-muted dark:text-foreground"), "View Code"),
		),
		sheet.Content(sheet.SideRight, h.Class("flex flex-col gap-0 border-l-0 p-0 sm:max-w-sm md:w-[700px] md:max-w-[700px] dark:border-l"),
			sheet.Header(h.Class("sr-only"), sheet.Title("Code"), sheet.Description("View the code for the chart.")),
			content,
		),
	)
}
