package inputgroup

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The registry/bases/radix/ui/input-group.tsx classes, resolved per style.
var file = style.NewFile("ui/input-group.tsx")

var (
	sGroup = file.Cn("group/input-group cn-input-group relative flex w-full min-w-0 items-center outline-none has-[>textarea]:h-auto", style.Arg)
	sAddon = file.CVA("cn-input-group-addon flex cursor-text items-center justify-center select-none",
		[]style.Variant{{Name: "align", Options: []style.Option{
			{Name: "inline-start", Class: "cn-input-group-addon-align-inline-start order-first"},
			{Name: "inline-end", Class: "cn-input-group-addon-align-inline-end order-last"},
			{Name: "block-start", Class: "cn-input-group-addon-align-block-start order-first w-full justify-start"},
			{Name: "block-end", Class: "cn-input-group-addon-align-block-end order-last w-full justify-start"},
		}}},
		cva.Selection{"align": "inline-start"},
	)
	sButton = file.CVA("cn-input-group-button flex items-center shadow-none",
		[]style.Variant{{Name: "size", Options: []style.Option{
			{Name: "xs", Class: "cn-input-group-button-size-xs"},
			{Name: "sm", Class: "cn-input-group-button-size-sm"},
			{Name: "icon-xs", Class: "cn-input-group-button-size-icon-xs"},
			{Name: "icon-sm", Class: "cn-input-group-button-size-icon-sm"},
		}}},
		cva.Selection{"size": "xs"},
	)
	sText     = file.Cn("cn-input-group-text flex items-center [&_svg]:pointer-events-none", style.Arg)
	sInput    = file.Cn("cn-input-group-input flex-1", style.Arg)
	sTextarea = file.Cn("cn-input-group-textarea flex-1 resize-none", style.Arg)
)

func addonClass(st style.Style, a Align, className string) string {
	if st == "" {
		return lib.Cn(AddonVariants(a), className)
	}
	return lib.Cn(sAddon.Classes(st, cva.Selection{"align": string(a)}), className)
}

func buttonClass(st style.Style, s ButtonSize, className string) string {
	if st == "" {
		return ButtonVariants(s, className)
	}
	return lib.Cn(sButton.Classes(st, cva.Selection{"size": string(s)}), className)
}
