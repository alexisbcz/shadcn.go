// Package datatable ports the upstream data-table-demo example: a table of
// payments with filtering, sorting, column visibility, row selection and
// pagination (TanStack Table upstream; its row models are ported in
// data-table.js), built on table, checkbox, input, button and
// dropdown-menu.
package datatable

import (
	_ "embed"
	"fmt"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/checkbox"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/table"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed data-table.js
var script string

func init() {
	runtime.Module("example-data-table", script)
	registry.Register(registry.Example{Name: "data-table-demo", Component: "data-table", Render: registry.Static(Demo)})
}

// Payment is a row of the demo.
type Payment struct {
	ID     string
	Amount float64
	Status string
	Email  string
}

// Payments are the demo's rows.
var Payments = []Payment{
	{"m5gr84i9", 316, "success", "ken99@example.com"},
	{"3u1reuv4", 242, "success", "Abe45@example.com"},
	{"derv1ws0", 837, "processing", "Monserrat44@example.com"},
	{"5kma53ae", 874, "success", "Silas22@example.com"},
	{"bhqecj4p", 721, "failed", "carmella@example.com"},
}

// usd formats an amount like Intl.NumberFormat("en-US", { style:
// "currency", currency: "USD" }).
func usd(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	intPart, frac := s[:len(s)-3], s[len(s)-3:]
	neg := false
	if intPart[0] == '-' {
		neg, intPart = true, intPart[1:]
	}
	out := ""
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out += ","
		}
		out += string(c)
	}
	if neg {
		return "-$" + out + frac
	}
	return "$" + out + frac
}

// rowActions is a row's actions menu.
func rowActions() h.Node {
	return dropdownmenu.DropdownMenu(
		dropdownmenu.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, h.Class("h-8 w-8 p-0"),
				h.Span(h.Class("sr-only"), "Open menu"),
				lucide.Ellipsis(),
			),
		),
		dropdownmenu.Content(dropdownmenu.AlignEnd,
			dropdownmenu.Label("Actions"),
			dropdownmenu.Item("Copy payment ID"),
			dropdownmenu.Separator(),
			dropdownmenu.Item("View customer"),
			dropdownmenu.Item("View payment details"),
		),
	)
}

// row renders a payment's row (cells in column order: select, status,
// email, amount, actions).
func row(p Payment) h.Node {
	return table.Row(h.Data("state", "false"),
		table.Cell(checkbox.Checkbox(h.Aria("label", "Select row"))),
		table.Cell(h.Div(h.Class("capitalize"), p.Status)),
		table.Cell(h.Div(h.Class("lowercase"), p.Email)),
		table.Cell(h.Div(h.Class("text-right font-medium"), usd(p.Amount))),
		table.Cell(rowActions()),
	)
}

// Demo is data-table-demo.
func Demo() h.Node {
	body := table.Body()
	ids := make([]string, 0, len(Payments))
	for _, p := range Payments {
		body.Children(row(p))
		ids = append(ids, p.ID)
	}
	columns := dropdownmenu.Content(dropdownmenu.AlignEnd)
	for _, id := range []string{"status", "email", "amount"} {
		columns.Children(dropdownmenu.CheckboxItem(h.Class("capitalize"), dropdownmenu.Checked(true), id))
	}
	return h.Group(
		h.Div(h.Class("w-full"),
			h.Div(h.Class("flex items-center py-4"),
				input.Input(h.Placeholder("Filter emails..."), h.Value(""), h.Class("max-w-sm")),
				dropdownmenu.DropdownMenu(
					dropdownmenu.Trigger(lib.AsChild,
						button.Button(button.VariantOutline, h.Class("ml-auto"), "Columns ", lucide.ChevronDown()),
					),
					columns,
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
		// The client's data (the rows' payment ids, for "Copy payment ID")
		// and the "No results." row.
		h.Template(h.Data("parity-ignore", ""), h.Data("data-table", ""), h.Data("ids", strings.Join(ids, " ")),
			table.Row(table.Cell(h.Attr("colspan", "5"), h.Class("h-24 text-center"), "No results.")),
			h.Div(h.Data("proto", "unchecked"), checkbox.Checkbox(h.Aria("label", "Select all"))),
			h.Div(h.Data("proto", "checked"), checkbox.Checkbox(h.Aria("label", "Select all"), checkbox.Checked(true))),
			h.Div(h.Data("proto", "indeterminate"), checkbox.Checkbox(h.Aria("label", "Select all"), checkbox.Indeterminate)),
		),
	)
}
