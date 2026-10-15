package table

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/table"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The aria table previews (ui/apps/v4/examples/aria/table-*.tsx): React
// Aria's TableHeader renders its own header row, without the row classes
// (no hover background), and the invoices tables put their caption in a
// figcaption after the table.
func init() {
	aria := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "aria--" + name, Component: "table", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		aria("table-demo", func() h.Node {
			return ariaInvoices(invoices, func(s string) string { return s }, "A list of your recent invoices.",
				[4]string{"Invoice", "Status", "Method", "Amount"}, "Total")
		}),
		aria("table-footer", func() h.Node {
			return ariaInvoices(invoices[:3], func(s string) string { return s }, "A list of your recent invoices.",
				[4]string{"Invoice", "Status", "Method", "Amount"}, "Total")
		}),
		aria("table-rtl", AriaRtl),
	)
}

func ariaInvoices(rows []invoice, t func(string) string, caption string, heads [4]string, total string, attrs ...any) h.Node {
	return h.El("figure", h.Class("w-full"),
		table.Table(append(attrs,
			table.Header(h.Tr(
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
		)...),
		h.El("figcaption", h.Data("slot", "table-caption"), h.Class("mt-4 text-center text-sm text-muted-foreground"), caption),
	)
}

// AriaRtl is examples/aria/table-rtl.tsx.
func AriaRtl() h.Node {
	ar := map[string]string{
		"Paid": "مدفوع", "Pending": "قيد الانتظار", "Unpaid": "غير مدفوع",
		"Credit Card": "بطاقة ائتمانية", "PayPal": "PayPal", "Bank Transfer": "تحويل بنكي",
	}
	return style.ScopeRTL(ariaInvoices(invoices, func(s string) string { return ar[s] }, "قائمة بفواتيرك الأخيرة.",
		[4]string{"الفاتورة", "الحالة", "الطريقة", "المبلغ"}, "المجموع", h.Dir("rtl")))
}
