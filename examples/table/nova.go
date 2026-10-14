package table

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/table"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the table page (ui/apps/v4/examples/radix/table-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "table", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("table-actions", NovaActions),
		nova("table-footer", func() h.Node {
			return novaInvoices(invoices[:3], func(s string) string { return s }, "A list of your recent invoices.",
				[4]string{"Invoice", "Status", "Method", "Amount"}, "Total")
		}),
		nova("table-rtl", NovaRtl),
	)
}

func novaInvoices(rows []invoice, t func(string) string, caption string, heads [4]string, total string, attrs ...any) h.Node {
	return table.Table(append(attrs,
		table.Caption(caption),
		table.Header(table.Row(
			table.Head(h.Class("w-[100px]"), heads[0]),
			table.Head(heads[1]),
			table.Head(heads[2]),
			table.Head(h.Class("text-right"), heads[3]),
		)),
		table.Body(h.Map(rows, func(inv invoice) h.Node {
			return table.Row(
				table.Cell(h.Class("font-medium"), inv.invoice),
				table.Cell(t(inv.paymentStatus)),
				table.Cell(t(inv.paymentMethod)),
				table.Cell(h.Class("text-right"), inv.totalAmount),
			)
		})),
		table.Footer(table.Row(
			table.Cell(h.Attr("colspan", "3"), total),
			table.Cell(h.Class("text-right"), "$2,500.00"),
		)),
	)...)
}

func NovaRtl() h.Node {
	ar := map[string]string{
		"Paid": "مدفوع", "Pending": "قيد الانتظار", "Unpaid": "غير مدفوع",
		"Credit Card": "بطاقة ائتمانية", "PayPal": "PayPal", "Bank Transfer": "تحويل بنكي",
	}
	return style.ScopeRTL(novaInvoices(invoices, func(s string) string { return ar[s] }, "قائمة بفواتيرك الأخيرة.",
		[4]string{"الفاتورة", "الحالة", "الطريقة", "المبلغ"}, "المجموع", h.Dir("rtl")))
}

func NovaActions() h.Node {
	row := func(product, price string) h.Node {
		return table.Row(
			table.Cell(h.Class("font-medium"), product),
			table.Cell(price),
			table.Cell(h.Class("text-right"),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIcon, h.Class("size-8"),
						lucide.Ellipsis(), h.Span(h.Class("sr-only"), "Open menu"))),
					dm.Content(dm.AlignEnd,
						dm.Item("Edit"),
						dm.Item("Duplicate"),
						dm.Separator(),
						dm.Item(dm.VariantDestructive, "Delete"),
					),
				),
			),
		)
	}
	return table.Table(
		table.Header(table.Row(
			table.Head("Product"),
			table.Head("Price"),
			table.Head(h.Class("text-right"), "Actions"),
		)),
		table.Body(
			row("Wireless Mouse", "$29.99"),
			row("Mechanical Keyboard", "$129.99"),
			row("USB-C Hub", "$49.99"),
		),
	)
}
