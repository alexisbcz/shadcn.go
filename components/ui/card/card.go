// Package card is the shadcn/ui Card: a container with a header, content
// and footer.
//
//	card.Card(
//	    card.Header(card.Title("Title"), card.Description("..."), card.Action(btn)),
//	    card.Content(...),
//	    card.Footer(...),
//	)
package card

import (
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

func slot(name, class string) *h.Element {
	return h.Div(h.Data("slot", name), h.Class(lib.Cn(class)))
}

var (
	cardBase        = slot("card", "flex flex-col gap-6 rounded-xl border bg-card py-6 text-card-foreground shadow-sm")
	headerBase      = slot("card-header", "@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6")
	titleBase       = slot("card-title", "leading-none font-semibold")
	descriptionBase = slot("card-description", "text-sm text-muted-foreground")
	actionBase      = slot("card-action", "col-start-2 row-span-2 row-start-1 self-start justify-self-end")
	contentBase     = slot("card-content", "px-6")
	footerBase      = slot("card-footer", "flex items-center px-6 [.border-t]:pt-6")
)

func spread(base *h.Element, children []any) *h.Element {
	return h.Spread(base, h.Div(children...), lib.Cn)
}

// Card returns the card's root <div>.
func Card(children ...any) *h.Element { return spread(cardBase, children) }

// Header returns the card's header <div>.
func Header(children ...any) *h.Element { return spread(headerBase, children) }

// Title returns the card's title <div>.
func Title(children ...any) *h.Element { return spread(titleBase, children) }

// Description returns the card's description <div>.
func Description(children ...any) *h.Element { return spread(descriptionBase, children) }

// Action returns the card's action <div>, placed at the header's end.
func Action(children ...any) *h.Element { return spread(actionBase, children) }

// Content returns the card's content <div>.
func Content(children ...any) *h.Element { return spread(contentBase, children) }

// Footer returns the card's footer <div>.
func Footer(children ...any) *h.Element { return spread(footerBase, children) }
