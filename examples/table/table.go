// Package table ports the upstream table-* examples.
package table

import (
	"shadcn.go/components/ui/table"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(
		registry.Example{Name: "table-demo", Component: "table", Render: registry.Static(Demo)},
	)
}

type invoice struct{ invoice, paymentStatus, totalAmount, paymentMethod string }

var invoices = []invoice{
	{"INV001", "Paid", "$250.00", "Credit Card"},
	{"INV002", "Pending", "$150.00", "PayPal"},
	{"INV003", "Unpaid", "$350.00", "Bank Transfer"},
	{"INV004", "Paid", "$450.00", "Credit Card"},
	{"INV005", "Paid", "$550.00", "PayPal"},
	{"INV006", "Pending", "$200.00", "Bank Transfer"},
	{"INV007", "Unpaid", "$300.00", "Credit Card"},
}

func Demo() h.Node {
	return table.Table(
		table.Caption("A list of your recent invoices."),
		table.Header(
			table.Row(
				table.Head(h.Class("w-[100px]"), "Invoice"),
				table.Head("Status"),
				table.Head("Method"),
				table.Head(h.Class("text-right"), "Amount"),
			),
		),
		table.Body(
			h.Map(invoices, func(inv invoice) h.Node {
				return table.Row(
					table.Cell(h.Class("font-medium"), inv.invoice),
					table.Cell(inv.paymentStatus),
					table.Cell(inv.paymentMethod),
					table.Cell(h.Class("text-right"), inv.totalAmount),
				)
			}),
		),
		table.Footer(
			table.Row(
				table.Cell(h.Attr("colspan", "3"), "Total"),
				table.Cell(h.Class("text-right"), "$2,500.00"),
			),
		),
	)
}
