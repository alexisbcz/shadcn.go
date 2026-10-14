package aspectratio

import (
	"shadcn.go/components/style"
	aspectratio "shadcn.go/components/ui/aspect-ratio"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the aspect-ratio page
// (ui/apps/v4/examples/radix/aspect-ratio-*.tsx), in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "aspect-ratio", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("aspect-ratio-demo", func() h.Node { return novaRatio("max-w-sm", 16.0/9, "w-full ") }),
		nova("aspect-ratio-portrait", func() h.Node { return novaRatio("max-w-[10rem]", 9.0/16, "") }),
		nova("aspect-ratio-square", func() h.Node { return novaRatio("max-w-[12rem]", 1, "") }),
		nova("aspect-ratio-rtl", NovaRtl),
	)
}

// novaImage is next/image with fill, as the reference renders it.
func novaImage(class string) h.Node {
	return h.Img(
		h.Alt("Photo"),
		h.Loading("lazy"),
		h.Attr("decoding", "async"),
		h.Data("nimg", "fill"),
		h.Class(class+"rounded-lg object-cover grayscale dark:brightness-20"),
		h.Style("color:transparent"),
		h.Src("https://avatar.vercel.sh/shadcn1"),
	)
}

func novaBox(ratio float64, imgClass string) h.Node {
	return aspectratio.AspectRatio(aspectratio.Ratio(ratio), h.Class("rounded-lg bg-muted"), novaImage(imgClass))
}

func novaRatio(width string, ratio float64, imgClass string) h.Node {
	return h.Div(h.Class("w-full "+width), novaBox(ratio, imgClass))
}

func NovaRtl() h.Node {
	return h.Figure(h.Class("w-full max-w-sm"), h.Dir("rtl"),
		novaBox(16.0/9, "w-full "),
		h.Figcaption(h.Class("mt-2 text-center text-sm text-muted-foreground"), "منظر طبيعي جميل"),
	)
}
