package field

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/styled"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// The registry/bases/radix/ui/field.tsx markup, resolved per style.
var file = style.NewFile("ui/field.tsx")

var (
	sSet         = file.Cn("cn-field-set flex flex-col", style.Arg)
	sLegend      = file.Cn("cn-field-legend", style.Arg)
	sGroup       = file.Cn("cn-field-group group/field-group @container/field-group flex w-full flex-col", style.Arg)
	sContent     = file.Cn("cn-field-content group/field-content flex flex-1 flex-col leading-snug", style.Arg)
	sLabel       = file.Cn("cn-field-label group/field-label peer/field-label flex w-fit", "has-[>[data-slot=field]]:w-full has-[>[data-slot=field]]:flex-col", style.Arg)
	sTitle       = file.Cn("cn-field-title flex w-fit items-center", style.Arg)
	sDescription = file.Cn("cn-field-description leading-normal font-normal group-has-data-horizontal/field:text-balance", "last:mt-0 nth-last-2:-mt-1", "[&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary", style.Arg)
	sSeparator   = file.Cn("cn-field-separator relative", style.Arg)
	sSepContent  = file.Attr("cn-field-separator-content relative mx-auto block w-fit bg-background")
	sSepLine     = file.Attr("absolute inset-0 top-1/2")
	sError       = file.Cn("cn-field-error font-normal", style.Arg)

	sFieldVariants = file.CVA("cn-field group/field flex w-full",
		[]style.Variant{{Name: "orientation", Options: []style.Option{
			{Name: "vertical", Class: "cn-field-orientation-vertical flex-col *:w-full [&>.sr-only]:w-auto"},
			{Name: "horizontal", Class: "cn-field-orientation-horizontal flex-row items-center has-[>[data-slot=field-content]]:items-start *:data-[slot=field-label]:flex-auto has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px"},
			{Name: "responsive", Class: "cn-field-orientation-responsive flex-col *:w-full @md/field-group:flex-row @md/field-group:items-center @md/field-group:*:w-auto @md/field-group:has-[>[data-slot=field-content]]:items-start @md/field-group:*:data-[slot=field-label]:flex-auto [&>.sr-only]:w-auto @md/field-group:has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px"},
		}}},
		cva.Selection{"orientation": "vertical"},
	)
)

func styledSpread(base *h.Element, children []any, alt func(style.Style) *h.Element) *styled.Node {
	return styled.Spread(base, h.El(base.Tag(), children...), alt)
}

// fieldClass is cn(fieldVariants({ orientation }), className) in style st.
func fieldClass(st style.Style, o Orientation, className string) string {
	return lib.Cn(sFieldVariants.Classes(st, cva.Selection{"orientation": string(o)}), className)
}
