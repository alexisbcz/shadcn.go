package buttongroup

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The registry/bases/radix/ui/button-group.tsx classes, resolved per style.
var file = style.NewFile("ui/button-group.tsx")

var (
	// As in the new-york-v4 variants, a child followed only by a <template>
	// (a Go dropdown menu's or tooltip's content) is the visual last child.
	sVariants = file.CVA(
		"cn-button-group group/button-group flex w-fit items-stretch *:focus-visible:relative *:focus-visible:z-10 [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1",
		[]style.Variant{{Name: "orientation", Options: []style.Option{
			{Name: "horizontal", Class: "cn-button-group-orientation-horizontal [&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child):not(:has(+template:last-child))]:rounded-r-none"},
			{Name: "vertical", Class: "cn-button-group-orientation-vertical flex-col [&>*:not(:first-child)]:rounded-t-none [&>*:not(:first-child)]:border-t-0 [&>*:not(:last-child):not(:has(+template:last-child))]:rounded-b-none"},
		}}},
		cva.Selection{"orientation": "horizontal"},
	)
	sText      = file.Cn("cn-button-group-text flex items-center [&_svg]:pointer-events-none", style.Arg)
	sSeparator = file.Cn("cn-button-group-separator relative self-stretch data-horizontal:mx-px data-horizontal:w-auto data-vertical:my-px data-vertical:h-auto", style.Arg)
)

func groupClass(st style.Style, o Orientation, className string) string {
	if st == "" {
		return Variants(o, className)
	}
	return lib.Cn(sVariants.Classes(st, cva.Selection{"orientation": string(o)}), className)
}

func textClassFor(st style.Style) string {
	if st == "" {
		return lib.Cn(textClass)
	}
	return sText.Class(st, "")
}

func separatorClassFor(st style.Style, className string) string {
	if st == "" {
		return lib.Cn(separatorClass, className)
	}
	return sSeparator.Class(st, className)
}
