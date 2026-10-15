package hugeicons

import h "shadcn.go/libhtml"

// icon builds the <svg> exactly as @hugeicons/react's HugeiconsIcon does
// with strokeWidth={2}: default attributes, color, the caller's classes,
// stroke width and stroke, then the shapes followed by any child nodes.
// (Without classes React renders class=""; libhtml omits the empty
// attribute, which renders the same.) Children may also add attributes and
// classes (h.Class, h.Aria, ...).
//
//	hugeicons.Search01Icon(h.Class("size-4"))
func icon(shapes string, children []any) *h.Element {
	return h.Svg(
		h.Xmlns("http://www.w3.org/2000/svg"),
		h.Width("24"), h.Height("24"),
		h.ViewBox("0 0 24 24"),
		h.Fill("none"),
		h.Attr("color", "currentColor"),
		h.StrokeWidth("2"),
		h.Stroke("currentColor"),
		// Generated from @hugeicons/core-free-icons at build time; never
		// user input.
		h.UnsafeHTML(shapes),
		h.Group(children...),
	)
}
