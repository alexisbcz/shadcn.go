package lucide

import h "shadcn.go/libhtml"

// icon builds the <svg> exactly as lucide-react's Icon component does:
// default attributes, class "lucide lucide-<name>" plus the caller's classes,
// then the shapes followed by any child nodes. Children may also add
// attributes and classes (h.Class, h.Aria, ...).
func icon(name, shapes string, children []any) *h.Element {
	return h.Svg(
		h.Xmlns("http://www.w3.org/2000/svg"),
		h.Width("24"), h.Height("24"),
		h.ViewBox("0 0 24 24"),
		h.Fill("none"),
		h.Stroke("currentColor"),
		h.StrokeWidth("2"),
		h.StrokeLinecap("round"),
		h.StrokeLinejoin("round"),
		h.Class("lucide lucide-"+name),
		// Generated from the lucide package at build time; never user input.
		h.UnsafeHTML(shapes),
		h.Group(children...),
	)
}
