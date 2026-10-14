package carousel

import (
	"fmt"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/carousel"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the carousel page
// (ui/apps/v4/examples/radix/carousel-*.tsx), rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "carousel", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("carousel-api", NovaAPI),
		nova("carousel-demo", NovaDemo),
		nova("carousel-orientation", NovaOrientation),
		nova("carousel-plugin", NovaPlugin),
		nova("carousel-size", NovaSize),
		nova("carousel-spacing", NovaSpacing),
		nova("carousel-multiple", NovaMultiple),
	)
}

// novaSlides is slides with a class on the card.
func novaSlides(itemClass, cardClass, contentClass, textClass string, padded bool) []any {
	out := make([]any, 5)
	for i := range out {
		c := card.Card(h.If(cardClass != "", h.Class(cardClass)), card.Content(h.Class(contentClass),
			h.Span(h.Class(textClass), fmt.Sprint(i+1)),
		))
		var inner h.Node = c
		if padded {
			inner = h.Div(h.Class("p-1"), c)
		}
		out[i] = carousel.Item(h.If(itemClass != "", h.Class(itemClass)), inner)
	}
	return out
}

func NovaAPI() h.Node {
	return h.Div(h.Class("mx-auto max-w-[10rem] sm:max-w-xs"), h.Data("example-carousel-api", ""),
		carousel.Carousel(h.Class("w-full max-w-xs"),
			carousel.Content(novaSlides("", "m-px", square, "text-4xl font-semibold", false)...),
			carousel.Previous(),
			carousel.Next(),
		),
		h.Div(h.Class("py-2 text-center text-sm text-muted-foreground"), "Slide ", "0", " of ", "0"),
	)
}

func NovaDemo() h.Node {
	return carousel.Carousel(h.Class("w-full max-w-[12rem] sm:max-w-xs"),
		carousel.Content(novaSlides("", "", square, "text-4xl font-semibold", true)...),
		carousel.Previous(),
		carousel.Next(),
	)
}

func NovaOrientation() h.Node {
	return carousel.Carousel(carousel.Opts{Align: "start"}, carousel.OrientationVertical, h.Class("w-full max-w-xs"),
		carousel.Content(h.Class("-mt-1 h-[270px]"),
			h.Group(novaSlides("basis-1/2 pt-1", "", "flex items-center justify-center p-6", "text-3xl font-semibold", true)...),
		),
		carousel.Previous(),
		carousel.Next(),
	)
}

func NovaPlugin() h.Node {
	return carousel.Carousel(carousel.Autoplay{Delay: 2000, StopOnHover: true}, h.Class("w-full max-w-[10rem] sm:max-w-xs"),
		carousel.Content(novaSlides("", "", square, "text-4xl font-semibold", true)...),
		carousel.Previous(),
		carousel.Next(),
	)
}

func NovaSize() h.Node {
	return carousel.Carousel(carousel.Opts{Align: "start"}, h.Class("w-full max-w-[12rem] sm:max-w-xs md:max-w-sm"),
		carousel.Content(novaSlides("basis-1/2 lg:basis-1/3", "", square, "text-3xl font-semibold", true)...),
		carousel.Previous(),
		carousel.Next(),
	)
}

func NovaSpacing() h.Node {
	return carousel.Carousel(h.Class("w-full max-w-[12rem] sm:max-w-xs md:max-w-sm"),
		carousel.Content(h.Class("-ml-1"),
			h.Group(novaSlides("basis-1/2 pl-1 lg:basis-1/3", "", square, "text-2xl font-semibold", true)...),
		),
		carousel.Previous(),
		carousel.Next(),
	)
}

func NovaMultiple() h.Node {
	return carousel.Carousel(h.Class("mx-auto max-w-xs sm:max-w-sm"), carousel.Opts{Align: "start"},
		carousel.Content(novaSlides("sm:basis-1/2 lg:basis-1/3", "", square, "text-3xl font-semibold", true)...),
		carousel.Previous(h.Class("hidden sm:inline-flex")),
		carousel.Next(h.Class("hidden sm:inline-flex")),
	)
}
