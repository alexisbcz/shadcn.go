package site

import h "shadcn.go/libhtml"

// Tabler icons (@tabler/icons-react 3.34.1) the docs chrome uses that the
// generated tabler package does not have yet. Same DOM as React renders
// (class "tabler-icon tabler-icon-<name> <className>").

func tablerIcon(name, shapes string, className ...string) *h.Element {
	cls := "tabler-icon tabler-icon-" + name + " "
	if len(className) > 0 {
		cls += className[0]
	}
	return h.Svg(
		h.Xmlns("http://www.w3.org/2000/svg"),
		h.Width("24"), h.Height("24"),
		h.ViewBox("0 0 24 24"),
		h.Fill("none"), h.Stroke("currentColor"), h.StrokeWidth("2"), h.StrokeLinecap("round"), h.StrokeLinejoin("round"),
		h.Class(cls),
		h.UnsafeHTML(shapes),
	)
}

func iconArrowLeft() h.Node {
	return tablerIcon("arrow-left", `<path d="M5 12l14 0"></path><path d="M5 12l6 6"></path><path d="M5 12l6 -6"></path>`)
}

func iconArrowRight() h.Node {
	return tablerIcon("arrow-right", `<path d="M5 12l14 0"></path><path d="M13 18l6 -6"></path><path d="M13 6l6 6"></path>`)
}

func iconCopy() h.Node {
	return tablerIcon("copy", `<path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path><path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>`)
}

func iconCheck(className ...string) h.Node {
	return tablerIcon("check", `<path d="M5 12l5 5l10 -10"></path>`, className...)
}

func iconChevronDown(className string) h.Node {
	return tablerIcon("chevron-down", `<path d="M6 9l6 6l6 -6"></path>`, className)
}
