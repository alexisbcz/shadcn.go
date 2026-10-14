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
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
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
func Table(children ...any) *scoped.Node {
	props := h.Table(children...)
	return scoped.New(
		func() *h.Element { return containerBase.Clone().Children(h.Spread(tableBase, props, lib.Cn)) },
		func(st style.Style) *h.Element {
			return h.Div(h.Data("slot", "table-container"), h.Class(baseFile.Attr("cn-table-container").Class(st)),
				h.Spread(scoped.Slot("table", "table", baseTable)(st), props, lib.Cn))
		},
	)
}

func part(base *h.Element, props *h.Element, slot string, site *style.Site) *scoped.Node {
	return scoped.Spread(base, props, scoped.Slot(props.Tag(), slot, site))
}

// Header returns a <thead>.
func Header(children ...any) *scoped.Node {
	return part(headerBase, h.Thead(children...), "table-header", baseHeader)
}

// Body returns a <tbody>.
func Body(children ...any) *scoped.Node {
	return part(bodyBase, h.Tbody(children...), "table-body", baseBody)
}

// Footer returns a <tfoot>.
func Footer(children ...any) *scoped.Node {
	return part(footerBase, h.Tfoot(children...), "table-footer", baseFooter)
}

// Row returns a <tr>.
func Row(children ...any) *scoped.Node { return part(rowBase, h.Tr(children...), "table-row", baseRow) }

// Head returns a <th>.
func Head(children ...any) *scoped.Node {
	return part(headBase, h.Th(children...), "table-head", baseHead)
}

// Cell returns a <td>.
func Cell(children ...any) *scoped.Node {
	return part(cellBase, h.Td(children...), "table-cell", baseCell)
}

// Caption returns a <caption>.
func Caption(children ...any) *scoped.Node {
	return part(captionBase, h.Caption(children...), "table-caption", baseCaption)
}

// registry/bases/radix/ui/table.tsx, used inside a style.Scope.
var (
	baseFile    = style.NewFile("ui/table.tsx")
	baseTable   = baseFile.Cn("cn-table", style.Arg)
	baseHeader  = baseFile.Cn("cn-table-header", style.Arg)
	baseBody    = baseFile.Cn("cn-table-body", style.Arg)
	baseFooter  = baseFile.Cn("cn-table-footer", style.Arg)
	baseRow     = baseFile.Cn("cn-table-row has-aria-expanded:bg-muted/50", style.Arg)
	baseHead    = baseFile.Cn("cn-table-head", style.Arg)
	baseCell    = baseFile.Cn("cn-table-cell", style.Arg)
	baseCaption = baseFile.Cn("cn-table-caption", style.Arg)
)
