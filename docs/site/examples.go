package site

import (
	"net/http"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// The examples: app/(app)/examples/layout.tsx and components/examples-nav.tsx.
// Each example page lives in examples_<name>.go.

const examplesTitle = "The Foundation for your Design System"
const examplesDescription = "A set of beautifully designed components that you can customize, extend, and build on. Start here then make it your own. Open Source. Open Code."

var exampleLinks = []struct{ name, href string }{
	{"Dashboard", "/examples/dashboard"},
	{"Tasks", "/examples/tasks"},
	{"Playground", "/examples/playground"},
	{"Authentication", "/examples/authentication"},
	{"RTL", "/examples/rtl"},
}

// examplesLayout renders an example page inside the examples layout, with the
// layout's metadata.
func (s *Site) examplesLayout(r *http.Request, canonical string, children ...any) h.Node {
	return s.examplesLayoutMeta(r, Meta{Title: examplesTitle, Description: examplesDescription, Canonical: canonical}, children...)
}

// examplesLayoutMeta renders an example page with its own metadata (a plain
// title: the layout sets a plain title, so no " - shadcn/ui" suffix).
func (s *Site) examplesLayoutMeta(r *http.Request, m Meta, children ...any) h.Node {
	return s.Document(r, m,
		s.AppLayout(r,
			pageHeader(
				announcement(),
				pageHeaderHeadingClass(examplesTitle, "max-w-4xl"),
				pageHeaderDescription(examplesDescription),
				// The examples layout's Button is styles/radix-luma/ui/button.
				pageActions(style.Scope(style.Luma,
					button.Button(lib.AsChild, h.Class("h-[35px]"), h.A(h.Href("/docs/installation"), "Get Started")),
					button.Button(lib.AsChild, button.VariantSecondary, h.A(h.Href("/docs/components"), "View Components")),
				)),
			),
			pageNav(h.ID("examples"), h.Class("hidden md:flex"),
				examplesNav(r.URL.Path),
				themeSelector("mr-4 hidden md:flex"),
			),
			h.Div(h.Class("container-wrapper flex flex-1 flex-col section-soft pb-6"),
				h.Div(h.Class("container flex flex-1 scroll-mt-20 flex-col theme-container"),
					h.Div(h.Class("flex flex-col overflow-hidden rounded-lg border bg-background bg-clip-padding has-[[data-slot=rtl-components]]:overflow-visible has-[[data-slot=rtl-components]]:border-0 has-[[data-slot=rtl-components]]:bg-transparent md:flex-1 xl:rounded-xl"),
						h.Group(children...),
					),
				),
			),
		),
	)
}

func pageHeaderHeadingClass(text, class string) h.Node {
	return h.H1(h.Class(lib.Cn("leading-tighter max-w-3xl text-3xl font-semibold tracking-tight text-balance text-primary lg:leading-[1.1] lg:font-semibold xl:text-5xl xl:tracking-tighter", class)), text)
}

func examplesNav(pathname string) h.Node {
	link := func(name, href string, active bool, extra ...any) h.Node {
		return h.A(h.Href(href),
			h.Class("flex h-7 items-center justify-center gap-2 px-4 text-center text-base font-medium text-muted-foreground transition-colors hover:text-primary data-[active=true]:text-primary"),
			h.Data("active", boolString(active)),
			name, h.Group(extra...))
	}
	links := []any{link("Examples", "/", pathname == "/")}
	for _, e := range exampleLinks {
		var extra []any
		if e.name == "RTL" {
			extra = append(extra, h.Span(h.Class("flex size-2 rounded-full bg-blue-500"), h.Title("New")))
		}
		links = append(links, link(e.name, e.href, len(pathname) >= len(e.href) && pathname[:len(e.href)] == e.href, extra...))
	}
	return h.Div(h.Class("flex items-center flex-1 overflow-hidden [&>a:first-child]:text-primary"),
		h.Div(h.Data("slot", "scroll-area"), h.Class("relative max-w-[96%] md:max-w-[600px] lg:max-w-none"),
			h.Div(h.Data("slot", "scroll-area-viewport"), h.Class("size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1"), h.Style("overflow: scroll"),
				h.Div(h.Style("min-width: 100%; display: table"),
					h.Div(h.Class("flex items-center"), h.Group(links...)),
				),
			),
		),
	)
}

// themeSelector (components/theme-selector.tsx) is in examples_theme.go.
