// Package table is the shadcn/ui Table: a responsive table.
//
//	table.Table(
//	    table.Caption("A list of your recent invoices."),
//	    table.Header(table.Row(table.Head("Invoice"))),
//	    table.Body(table.Row(table.Cell("INV001"))),
//	)
package table

import (
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

func slot(el *h.Element, name, class string) *h.Element {
	return el.Children(h.Data("slot", name), h.Class(lib.Cn(class)))
}

var (
	containerBase = h.Div(h.Data("slot", "table-container"), h.Class("relative w-full overflow-x-auto"))
	tableBase     = slot(h.Table(), "table", "w-full caption-bottom text-sm")
	headerBase    = slot(h.Thead(), "table-header", "[&_tr]:border-b")
	bodyBase      = slot(h.Tbody(), "table-body", "[&_tr:last-child]:border-0")
	footerBase    = slot(h.Tfoot(), "table-footer", "border-t bg-muted/50 font-medium [&>tr]:last:border-b-0")
	rowBase       = slot(h.Tr(), "table-row", "border-b transition-colors hover:bg-muted/50 has-aria-expanded:bg-muted/50 data-[state=selected]:bg-muted")
	headBase      = slot(h.Th(), "table-head", "h-10 px-2 text-left align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]")
	cellBase      = slot(h.Td(), "table-cell", "p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]")
	captionBase   = slot(h.Caption(), "table-caption", "mt-4 text-sm text-muted-foreground")
)

// Table returns a <table> inside a horizontally scrolling container <div>.
// Attributes apply to the <table>.
func Table(children ...any) *h.Element {
	return containerBase.Clone().Children(h.Spread(tableBase, h.Table(children...), lib.Cn))
}

// Header returns a <thead>.
func Header(children ...any) *h.Element { return h.Spread(headerBase, h.Thead(children...), lib.Cn) }

// Body returns a <tbody>.
func Body(children ...any) *h.Element { return h.Spread(bodyBase, h.Tbody(children...), lib.Cn) }

// Footer returns a <tfoot>.
func Footer(children ...any) *h.Element { return h.Spread(footerBase, h.Tfoot(children...), lib.Cn) }

// Row returns a <tr>.
func Row(children ...any) *h.Element { return h.Spread(rowBase, h.Tr(children...), lib.Cn) }

// Head returns a <th>.
func Head(children ...any) *h.Element { return h.Spread(headBase, h.Th(children...), lib.Cn) }

// Cell returns a <td>.
func Cell(children ...any) *h.Element { return h.Spread(cellBase, h.Td(children...), lib.Cn) }

// Caption returns a <caption>.
func Caption(children ...any) *h.Element {
	return h.Spread(captionBase, h.Caption(children...), lib.Cn)
}
