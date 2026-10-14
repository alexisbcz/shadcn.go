// Package scrollarea ports the upstream scroll-area-* examples and the
// port's scroll-area-* fixtures.
package scrollarea

import (
	"fmt"

	scrollarea "shadcn.go/components/ui/scroll-area"
	"shadcn.go/components/ui/separator"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "scroll-area", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("scroll-area-demo", Demo),
		ex("scroll-area-horizontal-demo", HorizontalDemo),
		ex("scroll-area-types", Types),
		ex("scroll-area-both", Both),
	)
}

func tags(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("v1.2.0-beta.%d", n-i)
	}
	return out
}

func tagList(n int) h.Node {
	return h.Div(h.Class("p-4"),
		h.H4(h.Class("mb-4 text-sm leading-none font-medium"), "Tags"),
		h.Map(tags(n), func(tag string) h.Node {
			return h.Group(h.Div(h.Class("text-sm"), tag), separator.Separator(h.Class("my-2")))
		}),
	)
}

func Demo() h.Node {
	return scrollarea.ScrollArea(h.Class("h-72 w-48 rounded-md border"), tagList(50))
}

type artwork struct{ artist, art string }

var works = []artwork{
	{"Ornella Binni", "https://images.unsplash.com/photo-1465869185982-5a1a7522cbcb?auto=format&fit=crop&w=300&q=80"},
	{"Tom Byrom", "https://images.unsplash.com/photo-1548516173-3cabfa4607e9?auto=format&fit=crop&w=300&q=80"},
	{"Vladimir Malyavko", "https://images.unsplash.com/photo-1494337480532-3725c85fd2ab?auto=format&fit=crop&w=300&q=80"},
}

// HorizontalDemo uses a plain <img> with the markup next/image renders for
// an unoptimized image with fixed dimensions (as the reference's shim does).
func HorizontalDemo() h.Node {
	return scrollarea.ScrollArea(h.Class("w-96 rounded-md border whitespace-nowrap"),
		h.Div(h.Class("flex w-max space-x-4 p-4"),
			h.Map(works, func(a artwork) h.Node {
				return h.Figure(h.Class("shrink-0"),
					h.Div(h.Class("overflow-hidden rounded-md"),
						h.Img(
							h.Alt("Photo by "+a.artist),
							h.Loading("lazy"),
							h.Width("300"),
							h.Height("400"),
							h.Attr("decoding", "async"),
							h.Data("nimg", "1"),
							h.Class("aspect-[3/4] h-fit w-fit object-cover"),
							h.Src(a.art),
							h.Style("color:transparent"),
						),
					),
					h.El("figcaption", h.Class("pt-2 text-xs text-muted-foreground"),
						"Photo by ",
						h.Span(h.Class("font-semibold text-foreground"), a.artist),
					),
				)
			}),
		),
		scrollarea.ScrollBar(scrollarea.Horizontal),
	)
}

// Types is the scroll-area-types fixture: the Radix types the examples do
// not use (scroll with a shorter hide delay, auto with and without
// overflow, always).
func Types() h.Node {
	box := func(children ...any) h.Node {
		return scrollarea.ScrollArea(append([]any{h.Class("h-40 w-40 rounded-md border")}, children...)...)
	}
	return h.Div(h.Class("grid grid-cols-2 gap-4"),
		box(scrollarea.TypeScroll, scrollarea.ScrollHideDelay(300), tagList(20)),
		box(scrollarea.TypeAuto, tagList(20)),
		box(scrollarea.TypeAuto, tagList(2)),
		box(scrollarea.TypeAlways, tagList(20)),
	)
}

func grid() h.Node {
	return h.Div(h.Class("grid w-[480px] grid-cols-6 gap-2 p-4"),
		h.Range(48, func(i int) h.Node {
			return h.Div(h.Class("flex h-12 items-center justify-center rounded-md border text-sm"), fmt.Sprint(i+1))
		}),
	)
}

// Both is the scroll-area-both fixture: both scrollbars (and the corner)
// with type="always", left to right and right to left.
func Both() h.Node {
	return h.Div(h.Class("flex gap-6"),
		scrollarea.ScrollArea(scrollarea.TypeAlways, h.Class("h-56 w-64 rounded-md border"),
			grid(),
			scrollarea.ScrollBar(scrollarea.Horizontal),
		),
		scrollarea.ScrollArea(scrollarea.TypeAlways, h.Dir("rtl"), h.Class("h-56 w-64 rounded-md border"),
			grid(),
			scrollarea.ScrollBar(scrollarea.Horizontal),
		),
	)
}
