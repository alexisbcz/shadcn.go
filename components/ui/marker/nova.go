package marker

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The registry/bases/radix/ui/marker.tsx classes, resolved per style. The
// parts render them inside style.Scope; outside, the new-york-v4 classes.
var file = style.NewFile("ui/marker.tsx")

var (
	sMarker = file.CVA("cn-marker group/marker relative flex w-full items-center",
		[]style.Variant{{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-marker-variant-default"},
			{Name: "separator", Class: "cn-marker-variant-separator"},
			{Name: "border", Class: "cn-marker-variant-border"},
		}}},
		cva.Selection{"variant": "default"},
	)
	sIcon    = file.Cn("cn-marker-icon shrink-0", style.Arg)
	sContent = file.Cn("cn-marker-content min-w-0 wrap-break-word", style.Arg)
)

func styleOf(w io.Writer) (style.Style, bool) {
	if w == nil {
		return "", false
	}
	return style.FromWriter(w)
}

func markerClass(w io.Writer, v Variant, className string) string {
	if st, ok := styleOf(w); ok {
		return lib.Cn(sMarker.Classes(st, cva.Selection{"variant": string(v)}, className))
	}
	return Variants(v, className)
}
