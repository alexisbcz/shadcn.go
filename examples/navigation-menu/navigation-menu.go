// Package navigationmenu ports the upstream navigation-menu-* examples and
// the navigation-menu parity fixtures.
package navigationmenu

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	nm "shadcn.go/components/ui/navigation-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "navigation-menu", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("navigation-menu-demo", Demo),
	)
}

var components = []struct{ title, href, description string }{
	{"Alert Dialog", "/docs/primitives/alert-dialog", "A modal dialog that interrupts the user with important content and expects a response."},
	{"Hover Card", "/docs/primitives/hover-card", "For sighted users to preview content available behind a link."},
	{"Progress", "/docs/primitives/progress", "Displays an indicator showing the completion progress of a task, typically displayed as a progress bar."},
	{"Scroll-area", "/docs/primitives/scroll-area", "Visually or semantically separates content."},
	{"Tabs", "/docs/primitives/tabs", "A set of layered sections of content—known as tab panels—that are displayed one at a time."},
	{"Tooltip", "/docs/primitives/tooltip", "A popup that displays information related to an element when the element receives keyboard focus or the mouse hovers over it."},
}

func listItem(href, title string, children ...any) h.Node {
	return h.Li(
		nm.Link(lib.AsChild, h.A(h.Href(href),
			h.Div(h.Class("text-sm leading-none font-medium"), title),
			h.P(append([]any{h.Class("line-clamp-2 text-sm leading-snug text-muted-foreground")}, children...)...),
		)),
	)
}

func link(href string, children ...any) h.Node {
	return nm.Link(lib.AsChild, h.A(append([]any{h.Href(href)}, children...)...))
}

// Demo is navigation-menu-demo. Upstream passes viewport={useIsMobile()},
// which is false at the parity viewport's width (and on the first render).
func Demo() h.Node {
	comps := make([]any, len(components))
	for i, c := range components {
		comps[i] = listItem(c.href, c.title, c.description)
	}
	return nm.NavigationMenu(nm.UseViewport(false),
		nm.List(h.Class("flex-wrap"),
			nm.Item(
				nm.Trigger("Home"),
				nm.Content(
					h.Ul(h.Class("grid gap-2 md:w-[400px] lg:w-[500px] lg:grid-cols-[.75fr_1fr]"),
						h.Li(h.Class("row-span-3"),
							nm.Link(lib.AsChild, h.A(
								h.Class("flex h-full w-full flex-col justify-end rounded-md bg-linear-to-b from-muted/50 to-muted p-4 no-underline outline-hidden transition-all duration-200 select-none focus:shadow-md md:p-6"),
								h.Href("/"),
								h.Div(h.Class("mb-2 text-lg font-medium sm:mt-4"), "shadcn/ui"),
								h.P(h.Class("text-sm leading-tight text-muted-foreground"), "Beautifully designed components built with Tailwind CSS."),
							)),
						),
						listItem("/docs", "Introduction", "Re-usable components built using Radix UI and Tailwind CSS."),
						listItem("/docs/installation", "Installation", "How to install dependencies and structure your app."),
						listItem("/docs/primitives/typography", "Typography", "Styles for headings, paragraphs, lists...etc"),
					),
				),
			),
			nm.Item(
				nm.Trigger("Components"),
				nm.Content(
					h.Ul(append([]any{h.Class("grid gap-2 sm:w-[400px] md:w-[500px] md:grid-cols-2 lg:w-[600px]")}, comps...)...),
				),
			),
			nm.Item(
				nm.Link(lib.AsChild, h.Class(nm.TriggerStyle()), h.A(h.Href("/docs"), "Docs")),
			),
			nm.Item(h.Class("hidden md:block"),
				nm.Trigger("List"),
				nm.Content(
					h.Ul(h.Class("grid w-[300px] gap-4"),
						h.Li(
							link("#", h.Div(h.Class("font-medium"), "Components"), h.Div(h.Class("text-muted-foreground"), "Browse all components in the library.")),
							link("#", h.Div(h.Class("font-medium"), "Documentation"), h.Div(h.Class("text-muted-foreground"), "Learn how to use the library.")),
							link("#", h.Div(h.Class("font-medium"), "Blog"), h.Div(h.Class("text-muted-foreground"), "Read our latest blog posts.")),
						),
					),
				),
			),
			nm.Item(h.Class("hidden md:block"),
				nm.Trigger("Simple"),
				nm.Content(
					h.Ul(h.Class("grid w-[200px] gap-4"),
						h.Li(
							link("#", "Components"),
							link("#", "Documentation"),
							link("#", "Blocks"),
						),
					),
				),
			),
			nm.Item(h.Class("hidden md:block"),
				nm.Trigger("With Icon"),
				nm.Content(
					h.Ul(h.Class("grid w-[200px] gap-4"),
						h.Li(
							link("#", h.Class("flex-row items-center gap-2"), lucide.CircleHelp(), "Backlog"),
							link("#", h.Class("flex-row items-center gap-2"), lucide.Circle(), "To Do"),
							link("#", h.Class("flex-row items-center gap-2"), lucide.CircleCheck(), "Done"),
						),
					),
				),
			),
		),
	)
}
