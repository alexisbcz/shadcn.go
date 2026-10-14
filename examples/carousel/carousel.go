// Package carousel ports the upstream carousel-* examples.
package carousel

import (
	_ "embed"
	"fmt"

	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/carousel"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed api.js
var apiScript string

func init() {
	runtime.Module("example-carousel-api", apiScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "carousel", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("carousel-api", API),
		ex("carousel-demo", Demo),
		ex("carousel-orientation", Orientation),
		ex("carousel-plugin", Plugin),
		ex("carousel-size", Size),
		ex("carousel-spacing", Spacing),
	)
}

// slides returns five numbered card slides.
func slides(itemClass, contentClass, textClass string, padded bool) []any {
	out := make([]any, 5)
	for i := range out {
		c := card.Card(card.Content(h.Class(contentClass),
			h.Span(h.Class(textClass), fmt.Sprint(i+1)),
		))
		var inner h.Node = c
		if padded {
			inner = h.Div(h.Class("p-1"), c)
		}
		out[i] = carousel.Item(h.Class(itemClass), inner)
	}
	return out
}

const square = "flex aspect-square items-center justify-center p-6"

// API shows "Slide n of 5" from the Embla API (api.js updates it).
func API() h.Node {
	return h.Div(h.Class("mx-auto max-w-xs"), h.Data("example-carousel-api", ""),
		carousel.Carousel(h.Class("w-full max-w-xs"),
			carousel.Content(slides("", square, "text-4xl font-semibold", false)...),
			carousel.Previous(),
			carousel.Next(),
		),
		// React's first render, before the effect reads the API.
		h.Div(h.Class("py-2 text-center text-sm text-muted-foreground"), "Slide ", "0", " of ", "0"),
	)
}

func Demo() h.Node {
	return carousel.Carousel(h.Class("w-full max-w-xs"),
		carousel.Content(slides("", square, "text-4xl font-semibold", true)...),
		carousel.Previous(),
		carousel.Next(),
	)
}

func Orientation() h.Node {
	return carousel.Carousel(carousel.Opts{Align: "start"}, carousel.OrientationVertical, h.Class("w-full max-w-xs"),
		carousel.Content(h.Class("-mt-1 h-[200px]"),
			h.Group(slides("pt-1 md:basis-1/2", "flex items-center justify-center p-6", "text-3xl font-semibold", true)...),
		),
		carousel.Previous(),
		carousel.Next(),
	)
}

func Plugin() h.Node {
	return carousel.Carousel(carousel.Autoplay{Delay: 2000, StopOnHover: true}, h.Class("w-full max-w-xs"),
		carousel.Content(slides("", square, "text-4xl font-semibold", true)...),
		carousel.Previous(),
		carousel.Next(),
	)
}

func Size() h.Node {
	return carousel.Carousel(carousel.Opts{Align: "start"}, h.Class("w-full max-w-sm"),
		carousel.Content(slides("md:basis-1/2 lg:basis-1/3", square, "text-3xl font-semibold", true)...),
		carousel.Previous(),
		carousel.Next(),
	)
}

func Spacing() h.Node {
	return carousel.Carousel(h.Class("w-full max-w-sm"),
		carousel.Content(h.Class("-ml-1"),
			h.Group(slides("pl-1 md:basis-1/2 lg:basis-1/3", square, "text-2xl font-semibold", true)...),
		),
		carousel.Previous(),
		carousel.Next(),
	)
}
