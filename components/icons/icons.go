// Package icons provides Lucide icons (https://lucide.dev, ISC license) as
// libhtml elements. Only the icons used so far are included; add more by
// copying their shapes from lucide-static.
//
//	icons.ArrowUp(h.Class("size-6"))
package icons

import (
	h "shadcn.go/libhtml"
)

// Lucide returns the standard 24x24 Lucide <svg> wrapper for the named icon.
// children may add classes and attributes as well as the shapes.
func Lucide(name string, children ...any) *h.Element {
	return h.Svg(
		h.Class("lucide lucide-"+name),
		h.Xmlns("http://www.w3.org/2000/svg"),
		h.Width("24"), h.Height("24"),
		h.ViewBox("0 0 24 24"),
		h.Fill("none"),
		h.Stroke("currentColor"),
		h.StrokeWidth("2"),
		h.StrokeLinecap("round"),
		h.StrokeLinejoin("round"),
		h.Aria("hidden", "true"),
		h.Group(children...),
	)
}

func path(d string) *h.Element { return h.Path(h.D(d)) }

// ArrowUp is lucide's arrow-up.
func ArrowUp(children ...any) *h.Element {
	return Lucide("arrow-up", path("m5 12 7-7 7 7"), path("M12 19V5"), h.Group(children...))
}

// ArrowUpRight is lucide's arrow-up-right.
func ArrowUpRight(children ...any) *h.Element {
	return Lucide("arrow-up-right", path("M7 7h10v10"), path("M7 17 17 7"), h.Group(children...))
}

// CircleFadingArrowUp is lucide's circle-fading-arrow-up.
func CircleFadingArrowUp(children ...any) *h.Element {
	return Lucide("circle-fading-arrow-up",
		path("M12 2a10 10 0 0 1 7.38 16.75"),
		path("m16 12-4-4-4 4"),
		path("M12 16V8"),
		path("M2.5 8.875a10 10 0 0 0-.5 3"),
		path("M2.83 16a10 10 0 0 0 2.43 3.4"),
		path("M4.636 5.235a10 10 0 0 1 .891-.857"),
		path("M8.644 21.42a10 10 0 0 0 7.631-.38"),
		h.Group(children...),
	)
}

// GitBranch is lucide's git-branch.
func GitBranch(children ...any) *h.Element {
	return Lucide("git-branch",
		h.Line(h.Attr("x1", "6"), h.Attr("x2", "6"), h.Attr("y1", "3"), h.Attr("y2", "15")),
		h.Circle(h.Attr("cx", "18"), h.Attr("cy", "6"), h.Attr("r", "3")),
		h.Circle(h.Attr("cx", "6"), h.Attr("cy", "18"), h.Attr("r", "3")),
		path("M18 9a9 9 0 0 1-9 9"),
		h.Group(children...),
	)
}

// LoaderCircle is lucide's loader-circle (Loader2Icon in lucide-react).
func LoaderCircle(children ...any) *h.Element {
	return Lucide("loader-circle", path("M21 12a9 9 0 1 1-6.219-8.56"), h.Group(children...))
}
