package item

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/styled"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// The registry/bases/radix/ui/item.tsx markup, resolved per style.
var file = style.NewFile("ui/item.tsx")

var (
	sGroup       = file.Cn("cn-item-group group/item-group flex w-full flex-col", style.Arg)
	sSeparator   = file.Cn("cn-item-separator", style.Arg)
	sContent     = file.Cn("cn-item-content flex flex-1 flex-col [&+[data-slot=item-content]]:flex-none", style.Arg)
	sTitle       = file.Cn("cn-item-title line-clamp-1 flex w-fit items-center", style.Arg)
	sDescription = file.Cn("cn-item-description line-clamp-2 font-normal [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary", style.Arg)
	sActions     = file.Cn("cn-item-actions flex items-center", style.Arg)
	sHeader      = file.Cn("cn-item-header flex basis-full items-center justify-between", style.Arg)
	sFooter      = file.Cn("cn-item-footer flex basis-full items-center justify-between", style.Arg)

	sItemVariants = file.CVA(
		"cn-item group/item flex w-full flex-wrap items-center transition-colors duration-100 outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 [a]:transition-colors",
		[]style.Variant{
			{Name: "variant", Options: []style.Option{
				{Name: "default", Class: "cn-item-variant-default"},
				{Name: "outline", Class: "cn-item-variant-outline"},
				{Name: "muted", Class: "cn-item-variant-muted"},
			}},
			{Name: "size", Options: []style.Option{
				{Name: "default", Class: "cn-item-size-default"},
				{Name: "sm", Class: "cn-item-size-sm"},
				{Name: "xs", Class: "cn-item-size-xs"},
			}},
		},
		cva.Selection{"variant": "default", "size": "default"},
	)
	sMediaVariants = file.CVA(
		"cn-item-media flex shrink-0 items-center justify-center [&_svg]:pointer-events-none",
		[]style.Variant{
			{Name: "variant", Options: []style.Option{
				{Name: "default", Class: "cn-item-media-variant-default"},
				{Name: "icon", Class: "cn-item-media-variant-icon"},
				{Name: "image", Class: "cn-item-media-variant-image"},
			}},
		},
		cva.Selection{"variant": "default"},
	)
)

// SizeXs is the extra-small size of the radix bases (styles only; it
// renders as the default size in new-york-v4).
const SizeXs Size = "xs"

func styledSpread(base *h.Element, children []any, alt func(style.Style) *h.Element) *styled.Node {
	return styled.Spread(base, h.El(base.Tag(), children...), alt)
}
