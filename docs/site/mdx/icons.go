package mdx

import h "shadcn.go/libhtml"

// Tabler icons (@tabler/icons-react 3.34.1) the docs components use that
// the shared tabler package does not generate. Same markup as tabler's
// outline icons; the shapes are constants.
func tablerIcon(name, shapes string, children ...any) *h.Element {
	return h.Svg(h.Xmlns("http://www.w3.org/2000/svg"), h.Width("24"), h.Height("24"), h.ViewBox("0 0 24 24"),
		h.Fill("none"), h.Stroke("currentColor"), h.StrokeWidth("2"), h.StrokeLinecap("round"), h.StrokeLinejoin("round"),
		h.Class("tabler-icon tabler-icon-"+name), h.UnsafeHTML(shapes), h.Group(children...))
}

func iconAlertCircle(children ...any) *h.Element {
	return tablerIcon("alert-circle", `<path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path><path d="M12 8v4"></path><path d="M12 16h.01"></path>`, children...)
}

func iconTerminal(children ...any) *h.Element {
	return tablerIcon("terminal", `<path d="M5 7l5 5l-5 5"></path><path d="M12 19l7 0"></path>`, children...)
}

func iconArrowUpRight(children ...any) *h.Element {
	return tablerIcon("arrow-up-right", `<path d="M17 7l-10 10"></path><path d="M8 7l9 0l0 9"></path>`, children...)
}
