package tabler

import h "shadcn.go/libhtml"

// icon builds the <svg> exactly as @tabler/icons-react does for the given
// type ("outline" or "filled").
func icon(kind, name, shapes string, children []any) *h.Element {
	el := h.Svg(
		h.Xmlns("http://www.w3.org/2000/svg"),
		h.Width("24"), h.Height("24"),
		h.ViewBox("0 0 24 24"),
	)
	if kind == "filled" {
		el.Children(h.Fill("currentColor"), h.Stroke("none"))
	} else {
		el.Children(h.Fill("none"), h.Stroke("currentColor"), h.StrokeWidth("2"), h.StrokeLinecap("round"), h.StrokeLinejoin("round"))
	}
	// Generated from the tabler package at build time; never user input.
	return el.Children(h.Class("tabler-icon tabler-icon-"+name), h.UnsafeHTML(shapes), h.Group(children...))
}
