// Package blocksdashboard ports the upstream dashboard-01 block
// (registry/new-york-v4/blocks/dashboard-01): an inset sidebar, a site
// header, section cards, an interactive area chart and a data table with
// tabs, row selection, column visibility, pagination and drag-and-drop
// reordering (dnd-kit and TanStack Table upstream, ported in
// dashboard.js).
package blocksdashboard

import (
	_ "embed"
	"net/http"
	"strconv"
	"strings"

	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/drawer"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/runtime"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/components/ui/table"
	"shadcn.go/components/ui/tabs"
	togglegroup "shadcn.go/components/ui/toggle-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

//go:embed dashboard.js
var script string

func init() {
	runtime.Module("example-blocks-dashboard", script)
	registry.Register(registry.Example{
		Name:      "dashboard-01",
		Component: "sidebar",
		// No Toaster: the reference mounts sonner's only for examples whose
		// own file imports sonner (page.tsx does not), so the cells' save
		// toasts are queued without being shown, as there.
		Render: Page,
	})
}

type item struct {
	id                                           int
	header, typ, status, target, limit, reviewer string
}

func nonce(r *http.Request) chart.Nonce { return chart.Nonce(secure.Nonce(r)) }

// Page is the block's page.tsx.
func Page(r *http.Request) h.Node {
	return sidebar.Provider(
		h.Style("--sidebar-width: calc(var(--spacing) * 72); --header-height: calc(var(--spacing) * 12)"),
		appSidebar(false),
		inset(r, false),
	)
}

// inset is the page's SidebarInset; example selects the variant of
// app/(app)/examples/dashboard.
func inset(r *http.Request, example bool) h.Node {
	return sidebar.Inset(
		siteHeader(example),
		h.Div(h.Class("flex flex-1 flex-col"),
			h.Div(h.Class("@container/main flex flex-1 flex-col gap-2"),
				h.Div(h.Class("flex flex-col gap-4 py-4 md:gap-6 md:py-6"),
					sectionCards(example),
					h.Div(h.Class("px-4 lg:px-6"), chartAreaInteractive(r, example)),
					dataTable(r, items),
				),
			),
		),
	)
}

// --- sidebar -------------------------------------------------------------------------

type navItem struct {
	title string
	icon  func(...any) *h.Element
}

func appSidebar(example bool) h.Node {
	navMain := []navItem{
		{"Dashboard", tabler.IconDashboard},
		{"Lifecycle", tabler.IconListDetails},
		{"Analytics", tabler.IconChartBar},
		{"Projects", tabler.IconFolder},
		{"Team", tabler.IconUsers},
	}
	documents := []navItem{
		{"Data Library", tabler.IconDatabase},
		{"Reports", tabler.IconReport},
		{"Word Assistant", tabler.IconFileWord},
	}
	navSecondary := []navItem{
		{"Settings", tabler.IconSettings},
		{"Get Help", tabler.IconHelp},
		{"Search", tabler.IconSearch},
	}
	props := []any{sidebar.CollapsibleOffcanvas, sidebar.VariantInset}
	headerClass := ""
	if example {
		props = []any{sidebar.CollapsibleNone, h.Class("h-auto border-r")}
		headerClass = "border-b"
	}
	return sidebar.Sidebar(h.Group(props...),
		sidebar.Header(h.Class(headerClass),
			sidebar.Menu(sidebar.MenuItem(
				sidebar.MenuButton(lib.AsChild, h.Class("data-[slot=sidebar-menu-button]:p-1.5!"),
					h.A(h.Href("#"),
						tabler.IconInnerShadowTop(h.Class("size-5!")),
						h.Span(h.Class("text-base font-semibold"), "Acme Inc."),
					),
				),
			)),
		),
		sidebar.Content(
			navMainGroup(navMain, example),
			navDocuments(documents),
			navSecondaryGroup(navSecondary, h.Class("mt-auto")),
		),
		sidebar.Footer(navUser()),
	)
}

func navMainGroup(items []navItem, example bool) h.Node {
	if example {
		return sidebar.Group(
			sidebar.GroupContent(
				sidebar.GroupLabel("Home"),
				sidebar.Menu(h.Map(items, func(it navItem) h.Node {
					return sidebar.MenuItem(sidebar.MenuButton(sidebar.Tooltip(it.title), it.icon(), h.Span(it.title)))
				})),
			),
		)
	}
	return sidebar.Group(
		sidebar.GroupContent(h.Class("flex flex-col gap-2"),
			sidebar.Menu(
				sidebar.MenuItem(h.Class("flex items-center gap-2"),
					sidebar.MenuButton(sidebar.Tooltip("Quick Create"),
						h.Class("min-w-8 bg-primary text-primary-foreground duration-200 ease-linear hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"),
						tabler.IconCirclePlusFilled(),
						h.Span("Quick Create"),
					),
					button.Button(button.SizeIcon, button.VariantOutline, h.Class("size-8 group-data-[collapsible=icon]:opacity-0"),
						tabler.IconMail(),
						h.Span(h.Class("sr-only"), "Inbox"),
					),
				),
			),
			sidebar.Menu(h.Map(items, func(it navItem) h.Node {
				return sidebar.MenuItem(sidebar.MenuButton(sidebar.Tooltip(it.title), it.icon(), h.Span(it.title)))
			})),
		),
	)
}

func navDocuments(items []navItem) h.Node {
	return sidebar.Group(h.Class("group-data-[collapsible=icon]:hidden"),
		sidebar.GroupLabel("Documents"),
		sidebar.Menu(
			h.Map(items, func(it navItem) h.Node {
				return sidebar.MenuItem(
					sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), it.icon(), h.Span(it.title))),
					dropdownmenu.DropdownMenu(
						dropdownmenu.Trigger(lib.AsChild,
							sidebar.MenuAction(sidebar.ShowOnHover(true), h.Class("rounded-sm data-[state=open]:bg-accent"),
								tabler.IconDots(),
								h.Span(h.Class("sr-only"), "More"),
							),
						),
						dropdownmenu.Content(h.Class("w-24 rounded-lg"), dropdownmenu.SideRight, dropdownmenu.AlignStart,
							dropdownmenu.Item(tabler.IconFolder(), h.Span("Open")),
							dropdownmenu.Item(tabler.IconShare3(), h.Span("Share")),
							dropdownmenu.Separator(),
							dropdownmenu.Item(dropdownmenu.VariantDestructive, tabler.IconTrash(), h.Span("Delete")),
						),
					),
				)
			}),
			sidebar.MenuItem(
				sidebar.MenuButton(h.Class("text-sidebar-foreground/70"),
					tabler.IconDots(h.Class("text-sidebar-foreground/70")),
					h.Span("More"),
				),
			),
		),
	)
}

func navSecondaryGroup(items []navItem, props ...any) h.Node {
	return sidebar.Group(append(props,
		sidebar.GroupContent(
			sidebar.Menu(h.Map(items, func(it navItem) h.Node {
				return sidebar.MenuItem(sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), it.icon(), h.Span(it.title))))
			})),
		),
	)...)
}

func userBlock(avatarClass string) []any {
	return []any{
		avatar.Avatar(h.Class(avatarClass),
			avatar.Image(h.Src("/avatars/shadcn.jpg"), h.Alt("shadcn")),
			avatar.Fallback(h.Class("rounded-lg"), "CN"),
		),
		h.Div(h.Class("grid flex-1 text-left text-sm leading-tight"),
			h.Span(h.Class("truncate font-medium"), "shadcn"),
			h.Span(h.Class("truncate text-xs text-muted-foreground"), "m@example.com"),
		),
	}
}

func navUser() h.Node {
	return sidebar.Menu(sidebar.MenuItem(
		dropdownmenu.DropdownMenu(
			dropdownmenu.Trigger(lib.AsChild,
				sidebar.MenuButton(sidebar.ButtonSizeLg,
					h.Class("data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"),
					userBlock("h-8 w-8 rounded-lg grayscale"),
					tabler.IconDotsVertical(h.Class("ml-auto size-4")),
				),
			),
			dropdownmenu.Content(h.Class("w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"),
				dropdownmenu.SideRight, dropdownmenu.AlignEnd, dropdownmenu.SideOffset(4),
				dropdownmenu.Label(h.Class("p-0 font-normal"),
					h.Div(append([]any{h.Class("flex items-center gap-2 px-1 py-1.5 text-left text-sm")}, userBlock("h-8 w-8 rounded-lg")...)...),
				),
				dropdownmenu.Separator(),
				dropdownmenu.Group(
					dropdownmenu.Item(tabler.IconUserCircle(), "Account"),
					dropdownmenu.Item(tabler.IconCreditCard(), "Billing"),
					dropdownmenu.Item(tabler.IconNotification(), "Notifications"),
				),
				dropdownmenu.Separator(),
				dropdownmenu.Item(tabler.IconLogout(), "Log out"),
			),
		),
	))
}

// --- header and cards ----------------------------------------------------------------

func siteHeader(example bool) h.Node {
	if example {
		return h.Header(h.Class("sticky top-0 z-10 flex h-(--header-height) shrink-0 items-center gap-2 border-b bg-background/90 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)"),
			h.Div(h.Class("flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6"),
				h.H1(h.Class("text-base font-medium"), "Documents"),
				h.Div(h.Class("ml-auto flex items-center gap-2"),
					button.Button(button.SizeSm, h.Class("hidden h-7 sm:flex"),
						tabler.IconCirclePlusFilled(),
						h.Span("Quick Create"),
					),
				),
			),
		)
	}
	return h.Header(h.Class("flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)"),
		h.Div(h.Class("flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6"),
			sidebar.Trigger(h.Class("-ml-1")),
			separator.Separator(separator.Vertical, h.Class("mx-2 data-[orientation=vertical]:h-4")),
			h.H1(h.Class("text-base font-medium"), "Documents"),
			h.Div(h.Class("ml-auto flex items-center gap-2"),
				button.Button(button.VariantGhost, lib.AsChild, button.SizeSm, h.Class("hidden sm:flex"),
					h.A(h.Href("https://github.com/shadcn-ui/ui/tree/main/apps/v4/app/(examples)/dashboard"),
						h.Rel("noopener noreferrer"), h.Target("_blank"), h.Class("dark:text-foreground"),
						"GitHub",
					),
				),
			),
		),
	)
}

type stat struct {
	description, title, badge string
	up                        bool
	footer, muted             string
}

func sectionCards(example bool) h.Node {
	class := "grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card *:data-[slot=card]:shadow-xs lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4 dark:*:data-[slot=card]:bg-card"
	if example {
		class = "grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:shadow-xs lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4"
	}
	stats := []stat{
		{"Total Revenue", "$1,250.00", "+12.5%", true, "Trending up this month", "Visitors for the last 6 months"},
		{"New Customers", "1,234", "-20%", false, "Down 20% this period", "Acquisition needs attention"},
		{"Active Accounts", "45,678", "+12.5%", true, "Strong user retention", "Engagement exceed targets"},
		{"Growth Rate", "4.5%", "+4.5%", true, "Steady performance increase", "Meets growth projections"},
	}
	return h.Div(h.Class(class),
		h.Map(stats, func(s stat) h.Node {
			icon := tabler.IconTrendingUp
			if !s.up {
				icon = tabler.IconTrendingDown
			}
			return card.Card(h.Class("@container/card"),
				card.Header(
					card.Description(s.description),
					card.Title(h.Class("text-2xl font-semibold tabular-nums @[250px]/card:text-3xl"), s.title),
					card.Action(badge.Badge(badge.VariantOutline, icon(), s.badge)),
				),
				card.Footer(h.Class("flex-col items-start gap-1.5 text-sm"),
					h.Div(h.Class("line-clamp-1 flex gap-2 font-medium"), s.footer, " ", icon(h.Class("size-4"))),
					h.Div(h.Class("text-muted-foreground"), s.muted),
				),
			)
		}),
	)
}

// --- interactive chart ---------------------------------------------------------------

func gradient(id, key string, top float64) *chart.Node {
	return chart.El("linearGradient", chart.P("id", id), chart.P("x1", "0"), chart.P("y1", "0"), chart.P("x2", "0"), chart.P("y2", "1"),
		chart.El("stop", chart.P("offset", "5%"), chart.P("stopColor", "var(--color-"+key+")"), chart.P("stopOpacity", top)),
		chart.El("stop", chart.P("offset", "95%"), chart.P("stopColor", "var(--color-"+key+")"), chart.P("stopOpacity", 0.1)),
	)
}

func chartAreaInteractive(r *http.Request, example bool) h.Node {
	// The example starts on the last 7 days and has no isMobile effect;
	// dashboard.js reads the initial range from data-dashboard-chart.
	initial := "90d"
	var tooltipExtra []any
	if example {
		// defaultIndex={isMobile ? -1 : 10}: the chart is hidden below md,
		// where isMobile holds.
		initial, tooltipExtra = "7d", []any{chart.P("defaultIndex", 10)}
	}
	config := chart.Config{
		{Key: "visitors", Label: "Visitors"},
		{Key: "desktop", Label: "Desktop", Color: "var(--primary)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--primary)"},
	}
	day := chart.DateFormat("en-US", map[string]string{"month": "short", "day": "numeric"})
	ranges := [][2]string{{"90d", "Last 3 months"}, {"30d", "Last 30 days"}, {"7d", "Last 7 days"}}
	area := func(key string) *chart.Node {
		return chart.Area(chart.DataKey(key), chart.Type("natural"), chart.Fill("url(#fill"+strings.ToUpper(key[:1])+key[1:]+")"),
			chart.Stroke("var(--color-"+key+")"), chart.StackID("a"))
	}
	return card.Card(h.Class("@container/card"), h.Data("dashboard-chart", h.IfElse(example, initial, "")),
		card.Header(
			card.Title("Total Visitors"),
			card.Description(
				h.Span(h.Class("hidden @[540px]/card:block"), "Total for the last 3 months"),
				h.Span(h.Class("@[540px]/card:hidden"), "Last 3 months"),
			),
			card.Action(
				togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value(initial), togglegroup.VariantOutline,
					h.Class("hidden *:data-[slot=toggle-group-item]:px-4! @[767px]/card:flex"),
					h.Map(ranges, func(rg [2]string) h.Node {
						return togglegroup.Item(togglegroup.Value(rg[0]), h.Data("value", rg[0]), rg[1])
					}),
				),
				selectui.Select(selectui.DefaultValue(initial),
					selectui.Trigger(selectui.SizeSm, h.Aria("label", "Select a value"),
						h.Class("flex w-40 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate @[767px]/card:hidden"),
						selectui.Value(selectui.Placeholder("Last 3 months")),
					),
					selectui.Content(h.Class("rounded-xl"),
						h.Map(ranges, func(rg [2]string) h.Node {
							return selectui.Item(selectui.ItemValue(rg[0]), h.Class("rounded-lg"), rg[1])
						}),
					),
				),
			),
		),
		card.Content(h.Class("px-2 pt-4 sm:px-6 sm:pt-6"),
			chart.Container(config, nonce(r), h.Class("aspect-auto h-[250px] w-full"),
				chart.AreaChart(chart.Data(visitors),
					chart.Defs(gradient("fillDesktop", "desktop", 1.0), gradient("fillMobile", "mobile", 0.8)),
					chart.CartesianGrid(chart.Vertical(false)),
					chart.XAxis(chart.DataKey("date"), chart.TickLine(false), chart.AxisLine(false), chart.TickMargin(8), chart.MinTickGap(32), chart.TickFormatter(day)),
					chart.Tooltip(append([]any{chart.Cursor(false), chart.TooltipContent(chart.LabelFormatter(day), chart.Indicator("dot"))}, tooltipExtra...)...),
					area("mobile"),
					area("desktop"),
				),
			),
		),
	)
}

// --- data table ----------------------------------------------------------------------

const cellInputClass = "h-8 w-16 border-transparent bg-transparent text-right shadow-none hover:bg-input/30 focus-visible:border focus-visible:bg-background dark:bg-transparent dark:hover:bg-input/30 dark:focus-visible:bg-input/30"

var hideable = []string{"type", "status", "target", "limit", "reviewer"}

func dragHandle() h.Node {
	return button.Button(button.VariantGhost, button.SizeIcon,
		h.Role("button"), h.TabIndex(0), h.Aria("disabled", "false"),
		h.Aria("roledescription", "sortable"), h.Aria("describedby", "DndDescribedBy-0"),
		h.Data("drag-handle", ""),
		h.Class("size-7 text-muted-foreground hover:bg-transparent"),
		tabler.IconGripVertical(h.Class("size-3 text-muted-foreground")),
		h.Span(h.Class("sr-only"), "Drag to reorder"),
	)
}

func editableCell(it item, field, value string) h.Node {
	id := strconv.Itoa(it.id) + "-" + field
	title := strings.ToUpper(field[:1]) + field[1:]
	return h.Form(h.Data("saving", "Saving "+it.header),
		label.Label(h.For(id), h.Class("sr-only"), title),
		input.Input(h.Class(cellInputClass), h.Value(value), h.ID(id)),
	)
}

func reviewerCell(it item) []any {
	if it.reviewer != "Assign reviewer" {
		return []any{it.reviewer}
	}
	id := strconv.Itoa(it.id) + "-reviewer"
	return []any{
		label.Label(h.For(id), h.Class("sr-only"), "Reviewer"),
		selectui.Select(
			selectui.Trigger(selectui.SizeSm, h.ID(id), h.Class("w-38 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate"),
				selectui.Value(selectui.Placeholder("Assign reviewer")),
			),
			selectui.Content(selectui.AlignEnd,
				selectui.Item(selectui.ItemValue("Eddie Lake"), "Eddie Lake"),
				selectui.Item(selectui.ItemValue("Jamik Tashpulatov"), "Jamik Tashpulatov"),
			),
		),
	}
}

func rowActions() h.Node {
	return dropdownmenu.DropdownMenu(
		dropdownmenu.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, button.SizeIcon, h.Class("flex size-8 text-muted-foreground data-[state=open]:bg-muted"),
				tabler.IconDotsVertical(),
				h.Span(h.Class("sr-only"), "Open menu"),
			),
		),
		dropdownmenu.Content(dropdownmenu.AlignEnd, h.Class("w-32"),
			dropdownmenu.Item("Edit"),
			dropdownmenu.Item("Make a copy"),
			dropdownmenu.Item("Favorite"),
			dropdownmenu.Separator(),
			dropdownmenu.Item(dropdownmenu.VariantDestructive, "Delete"),
		),
	)
}

func statusBadge(status string) h.Node {
	var icon h.Node
	if status == "Done" {
		icon = tabler.IconCircleCheckFilled(h.Class("fill-green-500 dark:fill-green-400"))
	} else {
		icon = tabler.IconLoader()
	}
	return badge.Badge(badge.VariantOutline, h.Class("px-1.5 text-muted-foreground"), icon, status)
}

func row(r *http.Request, it item) h.Node {
	return table.Row(
		h.Data("state", "false"), h.Data("dragging", "false"), h.Data("row-id", strconv.Itoa(it.id)),
		h.Class("relative z-0 data-[dragging=true]:z-10 data-[dragging=true]:opacity-80"),
		table.Cell(h.Data("col", "drag"), dragHandle()),
		table.Cell(h.Data("col", "select"), h.Div(h.Class("flex items-center justify-center"), checkbox.Checkbox(h.Aria("label", "Select row")))),
		table.Cell(h.Data("col", "header"), cellViewer(r, it)),
		table.Cell(h.Data("col", "type"), h.Div(h.Class("w-32"), badge.Badge(badge.VariantOutline, h.Class("px-1.5 text-muted-foreground"), it.typ))),
		table.Cell(h.Data("col", "status"), statusBadge(it.status)),
		table.Cell(h.Data("col", "target"), editableCell(it, "target", it.target)),
		table.Cell(h.Data("col", "limit"), editableCell(it, "limit", it.limit)),
		table.Cell(append([]any{h.Data("col", "reviewer")}, reviewerCell(it)...)...),
		table.Cell(h.Data("col", "actions"), rowActions()),
	)
}

func dataTable(r *http.Request, data []item) h.Node {
	views := [][2]string{{"outline", "Outline"}, {"past-performance", "Past Performance"}, {"key-personnel", "Key Personnel"}, {"focus-documents", "Focus Documents"}}
	columns := dropdownmenu.Content(dropdownmenu.AlignEnd, h.Class("w-56"))
	for _, c := range hideable {
		columns.Children(dropdownmenu.CheckboxItem(h.Class("capitalize"), dropdownmenu.Checked(true), c))
	}
	body := table.Body(h.Class("**:data-[slot=table-cell]:first:w-8"))
	for _, it := range data[:min(10, len(data))] {
		body.Children(row(r, it))
	}
	all := make([]any, 0, len(data))
	ids := make([]string, 0, len(data))
	for _, it := range data {
		all = append(all, row(r, it))
		ids = append(ids, strconv.Itoa(it.id))
	}
	pages := (len(data) + 9) / 10
	placeholder := func(v string) h.Node {
		return tabs.Content(tabs.Value(v), h.Class("flex flex-col px-4 lg:px-6"),
			h.Div(h.Class("aspect-video w-full flex-1 rounded-lg border border-dashed")),
		)
	}
	head := func(col string, children ...any) h.Node {
		return table.Head(append([]any{h.Data("col", col)}, children...)...)
	}
	return tabs.Tabs(tabs.Value("outline"), h.Class("w-full flex-col justify-start gap-6"), h.Data("dashboard-table", ""),
		h.Div(h.Class("flex items-center justify-between px-4 lg:px-6"),
			label.Label(h.For("view-selector"), h.Class("sr-only"), "View"),
			selectui.Select(selectui.DefaultValue("outline"),
				selectui.Trigger(selectui.SizeSm, h.ID("view-selector"), h.Class("flex w-fit @4xl/main:hidden"),
					selectui.Value(selectui.Placeholder("Select a view")),
				),
				selectui.Content(h.Map(views, func(v [2]string) h.Node {
					return selectui.Item(selectui.ItemValue(v[0]), v[1])
				})),
			),
			tabs.List(h.Class("hidden **:data-[slot=badge]:size-5 **:data-[slot=badge]:rounded-full **:data-[slot=badge]:bg-muted-foreground/30 **:data-[slot=badge]:px-1 @4xl/main:flex"),
				tabs.Trigger(tabs.Value("outline"), "Outline"),
				tabs.Trigger(tabs.Value("past-performance"), "Past Performance ", badge.Badge(badge.VariantSecondary, "3")),
				tabs.Trigger(tabs.Value("key-personnel"), "Key Personnel ", badge.Badge(badge.VariantSecondary, "2")),
				tabs.Trigger(tabs.Value("focus-documents"), "Focus Documents"),
			),
			h.Div(h.Class("flex items-center gap-2"),
				dropdownmenu.DropdownMenu(
					dropdownmenu.Trigger(lib.AsChild,
						button.Button(button.VariantOutline, button.SizeSm, h.Data("columns-trigger", ""),
							tabler.IconLayoutColumns(),
							h.Span(h.Class("hidden lg:inline"), "Customize Columns"),
							h.Span(h.Class("lg:hidden"), "Columns"),
							tabler.IconChevronDown(),
						),
					),
					columns,
				),
				button.Button(button.VariantOutline, button.SizeSm,
					tabler.IconPlus(),
					h.Span(h.Class("hidden lg:inline"), "Add Section"),
				),
			),
		),
		tabs.Content(tabs.Value("outline"), h.Class("relative flex flex-col gap-4 overflow-auto px-4 lg:px-6"),
			h.Div(h.Class("overflow-hidden rounded-lg border"),
				table.Table(
					table.Header(h.Class("sticky top-0 z-10 bg-muted"),
						table.Row(
							head("drag"),
							head("select", h.Div(h.Class("flex items-center justify-center"), checkbox.Checkbox(h.Aria("label", "Select all")))),
							head("header", "Header"),
							head("type", "Section Type"),
							head("status", "Status"),
							head("target", h.Div(h.Class("w-full text-right"), "Target")),
							head("limit", h.Div(h.Class("w-full text-right"), "Limit")),
							head("reviewer", "Reviewer"),
							head("actions"),
						),
					),
					body,
				),
			),
			h.Div(h.Class("flex items-center justify-between px-4"),
				h.Div(h.Class("hidden flex-1 text-sm text-muted-foreground lg:flex"), h.Data("selected-text", ""),
					"0 of "+strconv.Itoa(len(data))+" row(s) selected."),
				h.Div(h.Class("flex w-full items-center gap-8 lg:w-fit"),
					h.Div(h.Class("hidden items-center gap-2 lg:flex"),
						label.Label(h.For("rows-per-page"), h.Class("text-sm font-medium"), "Rows per page"),
						selectui.Select(selectui.DefaultValue("10"),
							selectui.Trigger(selectui.SizeSm, h.Class("w-20"), h.ID("rows-per-page"),
								selectui.Value(selectui.Placeholder("10")),
							),
							selectui.Content(selectui.SideTop, h.Map([]string{"10", "20", "30", "40", "50"}, func(s string) h.Node {
								return selectui.Item(selectui.ItemValue(s), s)
							})),
						),
					),
					h.Div(h.Class("flex w-fit items-center justify-center text-sm font-medium"), h.Data("page-text", ""),
						"Page 1 of "+strconv.Itoa(pages)),
					h.Div(h.Class("ml-auto flex items-center gap-2 lg:ml-0"),
						button.Button(button.VariantOutline, h.Class("hidden h-8 w-8 p-0 lg:flex"), h.Data("page", "first"), h.Disabled(),
							h.Span(h.Class("sr-only"), "Go to first page"), tabler.IconChevronsLeft()),
						button.Button(button.VariantOutline, button.SizeIcon, h.Class("size-8"), h.Data("page", "previous"), h.Disabled(),
							h.Span(h.Class("sr-only"), "Go to previous page"), tabler.IconChevronLeft()),
						button.Button(button.VariantOutline, button.SizeIcon, h.Class("size-8"), h.Data("page", "next"),
							h.Span(h.Class("sr-only"), "Go to next page"), tabler.IconChevronRight()),
						button.Button(button.VariantOutline, button.SizeIcon, h.Class("hidden size-8 lg:flex"), h.Data("page", "last"),
							h.Span(h.Class("sr-only"), "Go to last page"), tabler.IconChevronsRight()),
					),
				),
			),
			// dnd-kit's screen reader instructions and live region.
			h.Div(h.ID("DndDescribedBy-0"), h.Style("display: none;"),
				"To pick up a draggable item, press the space bar. While dragging, use the arrow keys to move the item. Press space bar again to drop the item in its new position, or press escape to cancel."),
			h.Div(h.ID("DndLiveRegion-0"), h.Role("status"), h.Aria("live", "assertive"), h.Aria("atomic", "true"),
				h.Style("position: fixed; top: 0px; left: 0px; width: 1px; height: 1px; margin: -1px; border: 0px; padding: 0px; overflow: hidden; clip: rect(0px, 0px, 0px, 0px); clip-path: inset(100%); white-space: nowrap;")),
			// Every row, cloned by dashboard.js when it enters the page.
			h.Template(h.Data("parity-ignore", ""), h.Data("rows", strings.Join(ids, " ")), all,
				h.Div(h.Data("proto", "checked"), checkbox.Checkbox(h.Aria("label", "Select all"), checkbox.Checked(true))),
				h.Div(h.Data("proto", "indeterminate"), checkbox.Checkbox(h.Aria("label", "Select all"), checkbox.Indeterminate)),
			),
		),
		placeholder("past-performance"),
		placeholder("key-personnel"),
		placeholder("focus-documents"),
	)
}

// --- row drawer ----------------------------------------------------------------------

func field(id, title string, control h.Node) h.Node {
	return h.Div(h.Class("flex flex-col gap-3"), label.Label(h.For(id), title), control)
}

func fieldSelect(id, value, placeholder string, options ...string) h.Node {
	return selectui.Select(selectui.DefaultValue(value),
		selectui.Trigger(h.ID(id), h.Class("w-full"), selectui.Value(selectui.Placeholder(placeholder))),
		selectui.Content(h.Map(options, func(o string) h.Node { return selectui.Item(selectui.ItemValue(o), o) })),
	)
}

var sixMonths = []chart.Row{
	{"month": "January", "desktop": 186, "mobile": 80},
	{"month": "February", "desktop": 305, "mobile": 200},
	{"month": "March", "desktop": 237, "mobile": 120},
	{"month": "April", "desktop": 73, "mobile": 190},
	{"month": "May", "desktop": 209, "mobile": 130},
	{"month": "June", "desktop": 214, "mobile": 140},
}

func cellViewer(r *http.Request, it item) h.Node {
	config := chart.Config{
		{Key: "desktop", Label: "Desktop", Color: "var(--primary)"},
		{Key: "mobile", Label: "Mobile", Color: "var(--primary)"},
	}
	area := func(key string, opacity float64) *chart.Node {
		return chart.Area(chart.DataKey(key), chart.Type("natural"), chart.Fill("var(--color-"+key+")"), chart.FillOpacity(opacity),
			chart.Stroke("var(--color-"+key+")"), chart.StackID("a"))
	}
	return drawer.Drawer(drawer.DirectionRight,
		drawer.Trigger(lib.AsChild,
			button.Button(button.VariantLink, h.Class("w-fit px-0 text-left text-foreground"), it.header),
		),
		drawer.Content(
			drawer.Header(h.Class("gap-1"),
				drawer.Title(it.header),
				drawer.Description("Showing total visitors for the last 6 months"),
			),
			h.Div(h.Class("flex flex-col gap-4 overflow-y-auto px-4 text-sm"),
				chart.Container(config, nonce(r),
					chart.AreaChart(chart.AccessibilityLayer(), chart.Data(sixMonths), chart.Margin{Left: 0, Right: 10},
						chart.CartesianGrid(chart.Vertical(false)),
						chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.AxisLine(false), chart.TickMargin(8), chart.TickFormatter(chart.Slice(0, 3)), chart.Hide()),
						chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("dot"))),
						area("mobile", 0.6),
						area("desktop", 0.4),
					),
				),
				separator.Separator(),
				h.Div(h.Class("grid gap-2"),
					h.Div(h.Class("flex gap-2 leading-none font-medium"),
						"Trending up by 5.2% this month ", tabler.IconTrendingUp(h.Class("size-4")),
					),
					h.Div(h.Class("text-muted-foreground"),
						"Showing total visitors for the last 6 months. This is just some random text to test the layout. It spans multiple lines and should wrap around."),
				),
				separator.Separator(),
				h.Form(h.Class("flex flex-col gap-4"),
					field("header", "Header", input.Input(h.ID("header"), h.Value(it.header))),
					h.Div(h.Class("grid grid-cols-2 gap-4"),
						field("type", "Type", fieldSelect("type", it.typ, "Select a type",
							"Table of Contents", "Executive Summary", "Technical Approach", "Design", "Capabilities", "Focus Documents", "Narrative", "Cover Page")),
						field("status", "Status", fieldSelect("status", it.status, "Select a status", "Done", "In Progress", "Not Started")),
					),
					h.Div(h.Class("grid grid-cols-2 gap-4"),
						field("target", "Target", input.Input(h.ID("target"), h.Value(it.target))),
						field("limit", "Limit", input.Input(h.ID("limit"), h.Value(it.limit))),
					),
					field("reviewer", "Reviewer", fieldSelect("reviewer", it.reviewer, "Select a reviewer", "Eddie Lake", "Jamik Tashpulatov", "Emily Whalen")),
				),
			),
			drawer.Footer(
				button.Button("Submit"),
				drawer.Close(lib.AsChild, button.Button(button.VariantOutline, "Done")),
			),
		),
	)
}

// ExamplePage is the sidebar of app/(app)/examples/dashboard/page.tsx: the
// dashboard-01 page with a static (collapsible="none") sidebar, a sticky
// header, plainer cards and a chart that starts on the last 7 days.
func ExamplePage(r *http.Request, props ...any) h.Node {
	return sidebar.Provider(
		h.Style("--sidebar-width: calc(var(--spacing) * 64); --header-height: calc(var(--spacing) * 12 + 1px)"),
		h.Group(props...),
		appSidebar(true),
		inset(r, true),
	)
}
