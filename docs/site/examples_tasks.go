package site

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/command"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/table"
	h "shadcn.go/libhtml"
)

// The tasks example: app/(app)/examples/tasks (a TanStack Table of tasks
// with a toolbar, faceted filters, column visibility, sorting, selection
// and pagination). The server renders the first page; examples-tasks.js
// keeps the table state on the client, with every row prerendered.

//go:embed data/examples/tasks/tasks.json
var tasksJSON []byte

//go:embed static/examples-tasks.js
var tasksScript string

func init() { runtime.Module("example-tasks", tasksScript) }

type task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Label    string `json:"label"`
	Priority string `json:"priority"`
}

type taskOption struct {
	value, label string
	icon         func(...any) *h.Element
}

var taskLabels = []taskOption{
	{"bug", "Bug", nil},
	{"feature", "Feature", nil},
	{"documentation", "Documentation", nil},
}

var taskStatuses = []taskOption{
	{"backlog", "Backlog", lucide.CircleHelp},
	{"todo", "Todo", lucide.Circle},
	{"in progress", "In Progress", lucide.Timer},
	{"done", "Done", lucide.CircleCheckBig},
	{"canceled", "Canceled", lucide.CircleOff},
}

var taskPriorities = []taskOption{
	{"low", "Low", lucide.ArrowDown},
	{"medium", "Medium", lucide.ArrowRight},
	{"high", "High", lucide.ArrowUp},
}

const tasksPageSize = 25

func findTaskOption(opts []taskOption, v string) *taskOption {
	for i := range opts {
		if opts[i].value == v {
			return &opts[i]
		}
	}
	return nil
}

// ExampleTasks serves /examples/tasks.
func (s *Site) ExampleTasks(r *http.Request) (h.Node, error) {
	var tasks []task
	if err := json.Unmarshal(tasksJSON, &tasks); err != nil {
		return nil, err
	}
	img := func(theme, class string) h.Node {
		u := "%2Fexamples%2Ftasks-" + theme + ".png"
		return h.Img(h.Alt("Playground"), h.Attr("loading", "lazy"), h.Width("1280"), h.Height("998"),
			h.Attr("decoding", "async"), h.Data("nimg", "1"), h.Class(class), h.Style("color:transparent"),
			h.UnsafeAttr("srcset", "/_next/image?url="+u+"&w=1920&q=75 1x, /_next/image?url="+u+"&w=3840&q=75 2x"),
			h.UnsafeAttr("src", "/_next/image?url="+u+"&w=3840&q=75"))
	}
	return s.examplesLayoutMeta(r, Meta{Title: "Tasks", PlainTitle: true, Description: "A task and issue tracker build using Tanstack Table.", Canonical: "/examples/tasks"},
		h.Div(h.Class("md:hidden"),
			img("light", "block dark:hidden"),
			img("dark", "hidden dark:block"),
		),
		h.Div(h.Class("hidden h-full flex-1 flex-col gap-8 p-8 md:flex"),
			h.Div(h.Class("flex items-center justify-between gap-2"),
				h.Div(h.Class("flex flex-col gap-1"),
					h.H2(h.Class("text-2xl font-semibold tracking-tight"), "Welcome back!"),
					h.P(h.Class("text-muted-foreground"), "Here's a list of your tasks for this month."),
				),
				h.Div(h.Class("flex items-center gap-2"), tasksUserNav()),
			),
			tasksDataTable(tasks),
		),
	), nil
}

// tasksUserNav is components/user-nav.tsx.
func tasksUserNav() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, h.Class("relative h-8 w-8 rounded-full"),
				avatar.Avatar(h.Class("h-9 w-9"),
					avatar.Image(h.Src("/avatars/03.png"), h.Alt("@shadcn")),
					avatar.Fallback("SC"),
				),
			),
		),
		dm.Content(h.Class("w-56"), dm.AlignEnd,
			dm.Label(h.Class("font-normal"),
				h.Div(h.Class("flex flex-col space-y-1"),
					h.P(h.Class("text-sm leading-none font-medium"), "shadcn"),
					h.P(h.Class("text-xs leading-none text-muted-foreground"), "m@example.com"),
				),
			),
			dm.Separator(),
			dm.Group(
				dm.Item("Profile", dm.Shortcut("⇧⌘P")),
				dm.Item("Billing", dm.Shortcut("⌘B")),
				dm.Item("Settings", dm.Shortcut("⌘S")),
				dm.Item("New Team"),
			),
			dm.Separator(),
			dm.Item("Log out", dm.Shortcut("⇧⌘Q")),
		),
	)
}

// tasksColumnHeader is components/data-table-column-header.tsx.
func tasksColumnHeader(title string, sortable bool) h.Node {
	if !sortable {
		return h.Div(h.Class(""), title)
	}
	return h.Div(h.Class("flex items-center gap-2"),
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild,
				button.Button(button.VariantGhost, button.SizeSm, h.Class("-ml-3 h-8 data-[state=open]:bg-accent"),
					h.Span(title), lucide.ChevronsUpDown(),
				),
			),
			dm.Content(dm.AlignStart,
				dm.Item(lucide.ArrowUp(), "Asc"),
				dm.Item(lucide.ArrowDown(), "Desc"),
				dm.Separator(),
				dm.Item(lucide.EyeOff(), "Hide"),
			),
		),
	)
}

// tasksFacetedFilter is components/data-table-faceted-filter.tsx (nothing
// selected; examples-tasks.js renders the selection).
func tasksFacetedFilter(column, title string, options []taskOption, counts map[string]int) h.Node {
	items := make([]any, 0, len(options))
	for _, o := range options {
		var count any
		if n := counts[o.value]; n > 0 {
			count = h.Span(h.Class("ml-auto flex size-4 items-center justify-center font-mono text-xs text-muted-foreground"), strconv.Itoa(n))
		}
		items = append(items, command.Item(h.Data("tasks-option", o.value),
			h.Div(h.Class("flex size-4 items-center justify-center rounded-[4px] border border-input [&_svg]:invisible"),
				lucide.Check(h.Class("size-3.5 text-primary-foreground")),
			),
			o.icon(h.Class("size-4 text-muted-foreground")),
			h.Span(o.label),
			count,
		))
	}
	return popover.Popover(
		popover.Trigger(lib.AsChild,
			button.Button(button.VariantOutline, button.SizeSm, h.Class("h-8 border-dashed"), h.Data("tasks-filter", column),
				lucide.CirclePlus(), title,
			),
		),
		popover.Content(h.Class("w-[200px] p-0"), popover.AlignStart,
			command.Command(
				command.Input(h.Placeholder(title)),
				command.List(
					command.Empty("No results found."),
					command.Group(h.Group(items...)),
				),
			),
		),
	)
}

// tasksRowActions is components/data-table-row-actions.tsx.
func tasksRowActions(t task) h.Node {
	radios := make([]any, 0, len(taskLabels))
	for _, l := range taskLabels {
		radios = append(radios, dm.RadioItem(dm.Value(l.value), l.label))
	}
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, button.SizeIcon, h.Class("size-8 data-[state=open]:bg-muted"),
				lucide.Ellipsis(), h.Span(h.Class("sr-only"), "Open menu"),
			),
		),
		dm.Content(dm.AlignEnd, h.Class("w-[160px]"),
			dm.Item("Edit"),
			dm.Item("Make a copy"),
			dm.Item("Favorite"),
			dm.Separator(),
			dm.Sub(
				dm.SubTrigger("Labels"),
				dm.SubContent(dm.RadioGroup(dm.Value(t.Label), h.Group(radios...))),
			),
			dm.Separator(),
			dm.Item(dm.VariantDestructive, "Delete", dm.Shortcut("⌘⌫")),
		),
	)
}

// tasksRow renders a task's row (cells: select, id, title, status,
// priority, actions).
func tasksRow(t task) h.Node {
	var labelBadge any
	if l := findTaskOption(taskLabels, t.Label); l != nil {
		labelBadge = badge.Badge(badge.VariantOutline, l.label)
	}
	var status, priority any
	if o := findTaskOption(taskStatuses, t.Status); o != nil {
		status = h.Div(h.Class("flex w-[100px] items-center gap-2"), o.icon(h.Class("size-4 text-muted-foreground")), h.Span(o.label))
	}
	if o := findTaskOption(taskPriorities, t.Priority); o != nil {
		priority = h.Div(h.Class("flex items-center gap-2"), o.icon(h.Class("size-4 text-muted-foreground")), h.Span(o.label))
	}
	return table.Row(h.Data("state", "false"),
		table.Cell(checkbox.Checkbox(h.Aria("label", "Select row"), h.Class("translate-y-[2px]"))),
		table.Cell(h.Div(h.Class("w-[80px]"), t.ID)),
		table.Cell(h.Div(h.Class("flex gap-2"), labelBadge, h.Span(h.Class("max-w-[500px] truncate font-medium"), t.Title))),
		table.Cell(status),
		table.Cell(priority),
		table.Cell(tasksRowActions(t)),
	)
}

// tasksDataTable is components/data-table.tsx with its toolbar and
// pagination.
func tasksDataTable(tasks []task) h.Node {
	statusCounts, priorityCounts := map[string]int{}, map[string]int{}
	for _, t := range tasks {
		statusCounts[t.Status]++
		priorityCounts[t.Priority]++
	}
	body := table.Body()
	rest := h.Tbody()
	for i, t := range tasks {
		if i < tasksPageSize {
			body.Children(tasksRow(t))
		} else {
			rest.Children(tasksRow(t))
		}
	}
	pageCount := (len(tasks) + tasksPageSize - 1) / tasksPageSize
	head := func(children ...any) h.Node { return table.Head(h.Attr("colspan", "1"), h.Group(children...)) }

	viewItems := []any{dm.Label("Toggle columns"), dm.Separator()}
	for _, id := range []string{"title", "status", "priority"} {
		viewItems = append(viewItems, dm.CheckboxItem(h.Class("capitalize"), dm.Checked(true), id))
	}
	pageSizes := []any{}
	for _, n := range []int{10, 20, 25, 30, 40, 50} {
		pageSizes = append(pageSizes, selectui.Item(selectui.ItemValue(strconv.Itoa(n)), strconv.Itoa(n)))
	}
	pageButton := func(label string, icon h.Node, class string, disabled bool) h.Node {
		b := button.Button(button.VariantOutline, button.SizeIcon, h.Class(class), h.Span(h.Class("sr-only"), label), icon)
		if disabled {
			b.Children(h.Disabled())
		}
		return b
	}
	data, _ := json.Marshal(tasks)
	return h.Div(h.Class("flex flex-col gap-4"), h.Data("tasks-table", ""),
		// Toolbar (components/data-table-toolbar.tsx).
		h.Div(h.Class("flex items-center justify-between"),
			h.Div(h.Class("flex flex-1 items-center gap-2"),
				input.Input(h.Placeholder("Filter tasks..."), h.Value(""), h.Class("h-8 w-[150px] lg:w-[250px]")),
				tasksFacetedFilter("status", "Status", taskStatuses, statusCounts),
				tasksFacetedFilter("priority", "Priority", taskPriorities, priorityCounts),
			),
			h.Div(h.Class("flex items-center gap-2"),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild,
						button.Button(button.VariantOutline, button.SizeSm, h.Class("ml-auto hidden h-8 lg:flex"), h.Data("tasks-view", ""),
							lucide.Settings2(), "View",
						),
					),
					dm.Content(dm.AlignEnd, h.Class("w-[150px]"), h.Group(viewItems...)),
				),
				button.Button(button.SizeSm, "Add Task"),
			),
		),
		h.Div(h.Class("overflow-hidden rounded-md border"),
			table.Table(
				table.Header(table.Row(
					head(checkbox.Checkbox(h.Aria("label", "Select all"), h.Class("translate-y-[2px]"))),
					head(tasksColumnHeader("Task", false)),
					head(tasksColumnHeader("Title", true)),
					head(tasksColumnHeader("Status", true)),
					head(tasksColumnHeader("Priority", true)),
					head(),
				)),
				body,
			),
		),
		// Pagination (components/data-table-pagination.tsx).
		h.Div(h.Class("flex items-center justify-between px-2"),
			h.Div(h.Class("flex-1 text-sm text-muted-foreground"), h.Data("tasks-selected", ""),
				"0", h.Comment(""), " of", h.Comment(""), " ", h.Comment(""), strconv.Itoa(len(tasks)), h.Comment(""), " row(s) selected.",
			),
			h.Div(h.Class("flex items-center space-x-6 lg:space-x-8"),
				h.Div(h.Class("flex items-center space-x-2"),
					h.P(h.Class("text-sm font-medium"), "Rows per page"),
					selectui.Select(selectui.DefaultValue(strconv.Itoa(tasksPageSize)),
						selectui.Trigger(h.Class("h-8 w-[70px]"), h.Data("tasks-page-size", ""), selectui.Value(selectui.Placeholder(strconv.Itoa(tasksPageSize)))),
						selectui.Content(selectui.SideTop, h.Group(pageSizes...)),
					),
				),
				h.Div(h.Class("flex w-[100px] items-center justify-center text-sm font-medium"), h.Data("tasks-page", ""),
					"Page ", h.Comment(""), "1", h.Comment(""), " of ", h.Comment(""), strconv.Itoa(pageCount),
				),
				h.Div(h.Class("flex items-center space-x-2"),
					pageButton("Go to first page", lucide.ChevronsLeft(), "hidden size-8 lg:flex", true),
					pageButton("Go to previous page", lucide.ChevronLeft(), "size-8", true),
					pageButton("Go to next page", lucide.ChevronRight(), "size-8", pageCount <= 1),
					pageButton("Go to last page", lucide.ChevronsRight(), "hidden size-8 lg:flex", pageCount <= 1),
				),
			),
		),
		// The client's data: the tasks, the rows past the first page, the
		// "No results." row and the pieces React renders on state changes.
		h.Template(h.Data("parity-ignore", ""), h.Data("tasks-data", string(data)),
			h.Table(rest),
			h.Table(h.Tbody(h.Data("proto", "no-results"), table.Row(table.Cell(h.Attr("colspan", "7"), h.Class("h-24 text-center"), "No results.")))),
			h.Div(h.Data("proto", "checked"), checkbox.Checkbox(h.Aria("label", "Select all"), h.Class("translate-y-[2px]"), checkbox.Checked(true))),
			h.Div(h.Data("proto", "indeterminate"), checkbox.Checkbox(h.Aria("label", "Select all"), h.Class("translate-y-[2px]"), checkbox.Indeterminate)),
			h.Div(h.Data("proto", "sort-asc"), lucide.ArrowUp()),
			h.Div(h.Data("proto", "sort-desc"), lucide.ArrowDown()),
			h.Div(h.Data("proto", "sort-none"), lucide.ChevronsUpDown()),
			h.Div(h.Data("proto", "separator"), separator.Separator(separator.Vertical, h.Class("mx-2 h-4"))),
			h.Div(h.Data("proto", "badge-count"), badge.Badge(badge.VariantSecondary, h.Class("rounded-sm px-1 font-normal lg:hidden"), "0")),
			h.Div(h.Data("proto", "badge"), badge.Badge(badge.VariantSecondary, h.Class("rounded-sm px-1 font-normal"), "0")),
			h.Div(h.Data("proto", "badges"), h.Div(h.Class("hidden gap-1 lg:flex"))),
			h.Div(h.Data("proto", "count"), h.Span(h.Class("ml-auto flex size-4 items-center justify-center font-mono text-xs text-muted-foreground"), "0")),
			h.Div(h.Data("proto", "reset"), button.Button(button.VariantGhost, button.SizeSm, "Reset", lucide.X())),
			h.Div(h.Data("proto", "clear"),
				command.Command(command.List(command.Separator(), command.Group(command.Item(h.Class("justify-center text-center"), "Clear filters")))),
			),
		),
	)
}
