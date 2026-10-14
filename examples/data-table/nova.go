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
	registry.Register(
		registry.Example{Name: "nova--data-table-demo", Component: "data-table", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, NovaDemo())
		})},
		registry.Example{Name: "nova--data-table-rtl", Component: "data-table", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, NovaRtl())
		})},
	)
}

// dtText is the text of the data-table previews (data-table-rtl
// translates it).
type dtText struct {
	filter, columns, status, email, amount, actions, copyID, viewCustomer, viewDetails string
	selectAll, selectRow, openMenu, noResults, of, selected, previous, next            string
	statuses                                                                           map[string]string
	amountFormat                                                                       func(float64) string
	rtl                                                                                bool
}

var dtEnglish = dtText{
	filter: "Filter emails...", columns: "Columns", status: "Status", email: "Email", amount: "Amount",
	actions: "Actions", copyID: "Copy payment ID", viewCustomer: "View customer", viewDetails: "View payment details",
	selectAll: "Select all", selectRow: "Select row", openMenu: "Open menu", noResults: "No results.",
	of: " of", selected: " row(s) selected.", previous: "Previous", next: "Next",
	amountFormat: usd,
}

var dtArabic = dtText{
	filter: "تصفية البريد الإلكتروني...", columns: "الأعمدة", status: "الحالة", email: "البريد الإلكتروني", amount: "المبلغ",
	actions: "الإجراءات", copyID: "نسخ معرف الدفع", viewCustomer: "عرض العميل", viewDetails: "عرض تفاصيل الدفع",
	selectAll: "تحديد الكل", selectRow: "تحديد الصف", openMenu: "فتح القائمة", noResults: "لا توجد نتائج.",
	of: " من", selected: " صف(وف) محدد.", previous: "السابق", next: "التالي",
	statuses:     map[string]string{"success": "ناجح", "processing": "قيد المعالجة", "failed": "فشل", "pending": "قيد الانتظار"},
	amountFormat: arSAUSD,
	rtl:          true,
}

// arSAUSD is Intl.NumberFormat("ar-SA", { style: "currency", currency:
// "USD" }): RLM, Arabic-Indic digits, Arabic decimal separator, NBSP, US$.
func arSAUSD(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	var b strings.Builder
	b.WriteString("\u200f")
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			b.WriteRune('٠' + (c - '0'))
		case c == '.':
			b.WriteRune('٫')
		default:
			b.WriteRune(c)
		}
	}
	b.WriteString("\u00a0US$")
	return b.String()
}

func novaRowActions(t dtText) h.Node {
	content := []any{dm.AlignEnd, h.Class("w-40"),
		dm.Group(
			dm.Label(t.actions),
			dm.Item(t.copyID),
		),
	}
	if t.rtl {
		content = append(content, h.Data("lang", "ar"), dm.Separator())
	}
	content = append(content, dm.Group(
		dm.Item(t.viewCustomer),
		dm.Item(t.viewDetails),
	))
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild,
			button.Button(button.VariantGhost, button.SizeIconXs,
				h.Span(h.Class("sr-only"), t.openMenu),
				lucide.Ellipsis(),
			),
		),
		dm.Content(content...),
	)
}

func NovaDemo() h.Node { return novaTable(dtEnglish) }

// NovaRtl is data-table-rtl: the ui-rtl components with Arabic text (no
// dir attribute, as upstream).
func NovaRtl() h.Node { return style.ScopeRTL(novaTable(dtArabic)) }

func novaTable(t dtText) h.Node {
	pick := func(en, rtl string) string {
		if t.rtl {
			return rtl
		}
		return en
	}
	body := table.Body()
	ids := make([]string, 0, len(Payments))
	for _, p := range Payments {
		status := p.Status
		if t.statuses != nil {
			status = t.statuses[p.Status]
		}
		body.Children(table.Row(h.Data("state", "false"),
			table.Cell(checkbox.Checkbox(h.Aria("label", t.selectRow))),
			table.Cell(h.Div(h.Class("capitalize"), status)),
			table.Cell(h.Div(h.Class("lowercase"), p.Email)),
			table.Cell(h.Div(h.Class(pick("text-right", "text-start")+" font-medium"), t.amountFormat(p.Amount))),
			table.Cell(novaRowActions(t)),
		))
		ids = append(ids, p.ID)
	}
	group := dm.Group()
	for _, id := range []string{"status", "email", "amount"} {
		group.Children(dm.CheckboxItem(h.Class("capitalize"), dm.Checked(true), id))
	}
	columnsContent := []any{dm.AlignEnd, group}
	if t.rtl {
		columnsContent = []any{dm.AlignStart, h.Data("lang", "ar"), group}
	}
	return h.Group(
		h.Div(h.Class("w-full"),
			h.Div(h.Class(pick("flex items-center py-4", "flex items-center gap-2 py-4")),
				input.Input(h.Placeholder(t.filter), h.Value(""), h.Class("max-w-sm")),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild,
						button.Button(button.VariantOutline, h.Class(pick("ml-auto", "ms-auto")), t.columns+" ", lucide.ChevronDown()),
					),
					dm.Content(columnsContent...),
				),
			),
			h.Div(h.Class("overflow-hidden rounded-md border"),
				table.Table(
					table.Header(table.Row(
						table.Head(checkbox.Checkbox(h.Aria("label", t.selectAll))),
						table.Head(t.status),
						table.Head(button.Button(button.VariantGhost, t.email, lucide.ArrowUpDown())),
						table.Head(h.Div(h.Class(pick("text-right", "text-start")), t.amount)),
						table.Head(),
					)),
					body,
				),
			),
			h.Div(h.Class(pick("flex items-center justify-end space-x-2 py-4", "flex items-center justify-end gap-2 py-4")),
				h.Div(h.Class("flex-1 text-sm text-muted-foreground"),
					"0", h.Comment(""), t.of, h.Comment(""), " ", h.Comment(""), fmt.Sprint(len(Payments)), h.Comment(""), t.selected,
				),
				h.Div(h.Class(pick("space-x-2", "flex gap-2")),
					button.Button(button.VariantOutline, button.SizeSm, h.Disabled(), t.previous),
					button.Button(button.VariantOutline, button.SizeSm, h.Disabled(), t.next),
				),
			),
		),
		h.Template(h.Data("parity-ignore", ""), h.Data("data-table", ""), h.Data("ids", strings.Join(ids, " ")),
			table.Row(table.Cell(h.Attr("colspan", "5"), h.Class("h-24 text-center"), t.noResults)),
			h.Div(h.Data("proto", "unchecked"), checkbox.Checkbox(h.Aria("label", t.selectAll))),
			h.Div(h.Data("proto", "checked"), checkbox.Checkbox(h.Aria("label", t.selectAll), checkbox.Checked(true))),
			h.Div(h.Data("proto", "indeterminate"), checkbox.Checkbox(h.Aria("label", t.selectAll), checkbox.Indeterminate)),
		),
	)
}
