// Package aspectratio ports the upstream aspect-ratio-* examples.
package aspectratio

import (
	aspectratio "shadcn.go/components/ui/aspect-ratio"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(
		registry.Example{Name: "aspect-ratio-demo", Component: "aspect-ratio", Render: registry.Static(Demo)},
	)
}

// Demo uses a plain <img> with the markup next/image renders for a fill
// image (as the reference's next/image shim does).
func Demo() h.Node {
	return aspectratio.AspectRatio(aspectratio.Ratio(16.0/9), h.Class("rounded-lg bg-muted"),
		h.Img(
			h.Alt("Photo by Drew Beamer"),
			h.Loading("lazy"),
			h.Attr("decoding", "async"),
			h.Data("nimg", "fill"),
			h.Style("color:transparent"),
			h.Src("https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=800&dpr=2&q=80"),
			h.Class("h-full w-full rounded-lg object-cover dark:brightness-[0.2] dark:grayscale"),
		),
	)
}
