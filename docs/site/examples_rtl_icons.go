package site

import h "shadcn.go/libhtml"

// Tabler icons the RTL example uses that components/icons/tabler does not
// generate (it only scans the registry), rendered as @tabler/icons-react 3.34.1.

func rtlTablerIcon(name, shapes string, children ...any) *h.Element {
	el := h.Svg(h.Xmlns("http://www.w3.org/2000/svg"), h.Width("24"), h.Height("24"), h.ViewBox("0 0 24 24"),
		h.Fill("none"), h.Stroke("currentColor"), h.StrokeWidth("2"), h.StrokeLinecap("round"), h.StrokeLinejoin("round"))
	// Generated from the tabler package; never user input.
	return el.Children(h.Class("tabler-icon tabler-icon-"+name), h.UnsafeHTML(shapes), h.Group(children...))
}

func rtlIconMinus(children ...any) *h.Element {
	return rtlTablerIcon("minus", "<path d=\"M5 12l14 0\"></path>", children...)
}

func rtlIconApps(children ...any) *h.Element {
	return rtlTablerIcon("apps", "<path d=\"M4 4m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z\"></path><path d=\"M4 14m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z\"></path><path d=\"M14 14m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z\"></path><path d=\"M14 7l6 0\"></path><path d=\"M17 4l0 6\"></path>", children...)
}

func rtlIconArrowUp(children ...any) *h.Element {
	return rtlTablerIcon("arrow-up", "<path d=\"M12 5l0 14\"></path><path d=\"M18 11l-6 -6\"></path><path d=\"M6 11l6 -6\"></path>", children...)
}

func rtlIconAt(children ...any) *h.Element {
	return rtlTablerIcon("at", "<path d=\"M12 12m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0\"></path><path d=\"M16 12v1.5a2.5 2.5 0 0 0 5 0v-1.5a9 9 0 1 0 -5.5 8.28\"></path>", children...)
}

func rtlIconBook(children ...any) *h.Element {
	return rtlTablerIcon("book", "<path d=\"M3 19a9 9 0 0 1 9 0a9 9 0 0 1 9 0\"></path><path d=\"M3 6a9 9 0 0 1 9 0a9 9 0 0 1 9 0\"></path><path d=\"M3 6l0 13\"></path><path d=\"M12 6l0 13\"></path><path d=\"M21 6l0 13\"></path>", children...)
}

func rtlIconCircleDashedPlus(children ...any) *h.Element {
	return rtlTablerIcon("circle-dashed-plus", "<path d=\"M8.56 3.69a9 9 0 0 0 -2.92 1.95\"></path><path d=\"M3.69 8.56a9 9 0 0 0 -.69 3.44\"></path><path d=\"M3.69 15.44a9 9 0 0 0 1.95 2.92\"></path><path d=\"M8.56 20.31a9 9 0 0 0 3.44 .69\"></path><path d=\"M15.44 20.31a9 9 0 0 0 2.92 -1.95\"></path><path d=\"M20.31 15.44a9 9 0 0 0 .69 -3.44\"></path><path d=\"M20.31 8.56a9 9 0 0 0 -1.95 -2.92\"></path><path d=\"M15.44 3.69a9 9 0 0 0 -3.44 -.69\"></path><path d=\"M9 12h6\"></path><path d=\"M12 9v6\"></path>", children...)
}

func rtlIconPaperclip(children ...any) *h.Element {
	return rtlTablerIcon("paperclip", "<path d=\"M15 7l-6.5 6.5a1.5 1.5 0 0 0 3 3l6.5 -6.5a3 3 0 0 0 -6 -6l-6.5 6.5a4.5 4.5 0 0 0 9 9l6.5 -6.5\"></path>", children...)
}

func rtlIconWorld(children ...any) *h.Element {
	return rtlTablerIcon("world", "<path d=\"M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0\"></path><path d=\"M3.6 9h16.8\"></path><path d=\"M3.6 15h16.8\"></path><path d=\"M11.5 3a17 17 0 0 0 0 18\"></path><path d=\"M12.5 3a17 17 0 0 1 0 18\"></path>", children...)
}

func rtlIconX(children ...any) *h.Element {
	return rtlTablerIcon("x", "<path d=\"M18 6l-12 12\"></path><path d=\"M6 6l12 12\"></path>", children...)
}
