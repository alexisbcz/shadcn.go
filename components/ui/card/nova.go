package card

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	h "shadcn.go/libhtml"
)

// registry/bases/radix/ui/card.tsx, used inside a style.Scope.
var (
	baseFile        = style.NewFile("ui/card.tsx")
	baseCard        = baseFile.Cn("cn-card group/card flex flex-col", style.Arg)
	baseHeader      = baseFile.Cn("cn-card-header group/card-header @container/card-header grid auto-rows-min items-start has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto]", style.Arg)
	baseTitle       = baseFile.Cn("cn-card-title cn-font-heading", style.Arg)
	baseDescription = baseFile.Cn("cn-card-description", style.Arg)
	baseAction      = baseFile.Cn("cn-card-action col-start-2 row-span-2 row-start-1 self-start justify-self-end", style.Arg)
	baseContent     = baseFile.Cn("cn-card-content", style.Arg)
	baseFooter      = baseFile.Cn("cn-card-footer flex items-center", style.Arg)
)

// Size selects the card's size in the bases styles (data-size).
type Size string

// Card sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
)

func spread(base *h.Element, slot string, site *style.Site, children []any) *scoped.Node {
	props := h.Div(children...)
	return scoped.Spread(base, props, scoped.Slot("div", slot, site))
}
