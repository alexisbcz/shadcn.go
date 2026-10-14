package datatable

import (
	"fmt"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/checkbox"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/table"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The data-table docs preview (ui/apps/v4/examples/radix/data-table-demo.tsx),
// in radix-nova: the same table as data-table-demo, with the radix
// examples' menus (groups, icon-xs trigger). data-table.js drives both.
func init() {
	registry.Register(registry.Example{Name: "nova--data-table-demo", Component: "data-table", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, NovaDemo())
	})})
}

func novaRowActions() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, button.SizeIconXs,
				h.Span(h.Class("sr-only"), "Open menu"),
				lucide.Ellipsis(),
			),
		),
		dm.Content(dm.AlignEnd, h.Class("w-40"),
			dm.Group(
				dm.Label("Actions"),
				dm.Item("Copy payment ID"),
			),
			dm.Group(
				dm.Item("View customer"),
				dm.Item("View payment details"),
			),
		),
	)
}

func NovaDemo() h.Node {
	body := table.Body()
	ids := make([]string, 0, len(Payments))
	for _, p := range Payments {
		body.Children(table.Row(h.Data("state", "false"),
			table.Cell(checkbox.Checkbox(h.Aria("label", "Select row"))),
			table.Cell(h.Div(h.Class("capitalize"), p.Status)),
			table.Cell(h.Div(h.Class("lowercase"), p.Email)),
			table.Cell(h.Div(h.Class("text-right font-medium"), usd(p.Amount))),
			table.Cell(novaRowActions()),
		))
		ids = append(ids, p.ID)
	}
	group := dm.Group()
	for _, id := range []string{"status", "email", "amount"} {
		group.Children(dm.CheckboxItem(h.Class("capitalize"), dm.Checked(true), id))
	}
	return h.Group(
		h.Div(h.Class("w-full"),
			h.Div(h.Class("flex items-center py-4"),
				input.Input(h.Placeholder("Filter emails..."), h.Value(""), h.Class("max-w-sm")),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild,
						button.Button(button.VariantOutline, h.Class("ml-auto"), "Columns ", lucide.ChevronDown()),
					),
					dm.Content(dm.AlignEnd, group),
				),
			),
			h.Div(h.Class("overflow-hidden rounded-md border"),
				table.Table(
					table.Header(table.Row(
						table.Head(checkbox.Checkbox(h.Aria("label", "Select all"))),
						table.Head("Status"),
						table.Head(button.Button(button.VariantGhost, "Email", lucide.ArrowUpDown())),
						table.Head(h.Div(h.Class("text-right"), "Amount")),
						table.Head(),
					)),
					body,
				),
			),
			h.Div(h.Class("flex items-center justify-end space-x-2 py-4"),
				h.Div(h.Class("flex-1 text-sm text-muted-foreground"),
					"0", h.Comment(""), " of", h.Comment(""), " ", h.Comment(""), fmt.Sprint(len(Payments)), h.Comment(""), " row(s) selected.",
				),
				h.Div(h.Class("space-x-2"),
					button.Button(button.VariantOutline, button.SizeSm, h.Disabled(), "Previous"),
					button.Button(button.VariantOutline, button.SizeSm, h.Disabled(), "Next"),
				),
			),
		),
		h.Template(h.Data("parity-ignore", ""), h.Data("data-table", ""), h.Data("ids", strings.Join(ids, " ")),
			table.Row(table.Cell(h.Attr("colspan", "5"), h.Class("h-24 text-center"), "No results.")),
			h.Div(h.Data("proto", "unchecked"), checkbox.Checkbox(h.Aria("label", "Select all"))),
			h.Div(h.Data("proto", "checked"), checkbox.Checkbox(h.Aria("label", "Select all"), checkbox.Checked(true))),
			h.Div(h.Data("proto", "indeterminate"), checkbox.Checkbox(h.Aria("label", "Select all"), checkbox.Indeterminate)),
		),
	)
}
