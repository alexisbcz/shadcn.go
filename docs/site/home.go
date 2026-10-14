package site

import (
	"encoding/json"
	"net/http"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Home renders / (app/(app)/(root)/page.tsx): the page header, then the
// dashboard screenshot on mobile and the cards demo (home_cards.go) from md
// up.
func (s *Site) Home(r *http.Request) (h.Node, error) {
	const title = "The Foundation for your Design System"
	ld, _ := json.Marshal(struct {
		Context       string   `json:"@context"`
		Type          string   `json:"@type"`
		ID            string   `json:"@id"`
		URL           string   `json:"url"`
		Name          string   `json:"name"`
		AlternateName []string `json:"alternateName"`
		Description   string   `json:"description"`
		InLanguage    string   `json:"inLanguage"`
		SameAs        []string `json:"sameAs"`
	}{"https://schema.org", "WebSite", SiteURL + "/#website", SiteURL, SiteName,
		[]string{"shadcn", "ui.shadcn.com"}, SiteDescription, "en-US",
		[]string{"https://github.com/shadcn-ui/ui", "https://twitter.com/shadcn"}})
	image := func(theme, class string) h.Node {
		src := "/_next/image?url=%2Fimages%2Ffull-" + theme + ".png&w=3840&q=75"
		return h.Img(h.Alt("Dashboard"), h.Width("2560"), h.Height("2764"), h.Attr("decoding", "async"),
			h.Data("nimg", "1"), h.Class(class), h.Style("color:transparent"), h.Attr("srcset", src+" 1x"), h.Src(src))
	}
	return s.Document(r, Meta{Title: SiteName + " - " + title, PlainTitle: true, Description: SiteDescription, Canonical: "/"},
		s.AppLayout(r,
			h.Div(h.Class("flex flex-1 flex-col"),
				h.Script(h.Type("application/ld+json"), h.UnsafeJS(strings.ReplaceAll(string(ld), "<", `<`))),
				pageHeaderClass("md:**:[.container]:pb-8 lg:**:[.container]:pb-12",
					announcement(),
					pageHeaderHeadingClass(title, "max-w-4xl"),
					pageHeaderDescription(SiteDescription),
					// The root page's Button is styles/radix-luma/ui/button.
					pageActions(style.Scope(style.Luma,
						button.Button(lib.AsChild, h.Class("h-[35px]"), h.A(h.Href("/docs/installation"), "Get Started")),
						button.Button(lib.AsChild, button.VariantSecondary, h.A(h.Href("/docs/components"), "View Components")),
					)),
				),
				h.Div(h.Class("container-wrapper flex-1 p-0"),
					h.Div(h.Class("container overflow-hidden md:px-0 lg:max-w-none"),
						h.Section(h.Class("-mx-4 w-[140vw] overflow-hidden md:hidden"),
							image("light", "block h-auto w-full dark:hidden"),
							image("dark", "hidden h-auto w-full dark:block"),
						),
						h.Section(h.Class("hidden md:block"), cardsDemo(r)),
					),
				),
			),
		),
	), nil
}

// pageHeaderClass is pageHeader with a className.
func pageHeaderClass(class string, children ...any) h.Node {
	return h.Section(h.Class(lib.Cn("border-grid", class)),
		h.Div(h.Class("container-wrapper"),
			h.Div(h.Class("container flex flex-col items-center gap-2 px-6 py-8 text-center md:py-16 lg:py-20 xl:gap-4"), h.Group(children...)),
		),
	)
}

// messageScrollerDemo is examples/radix/message-scroller-demo (radix-rhea).
func messageScrollerDemo(r *http.Request) h.Node {
	if ex, ok := registry.Get("nova--message-scroller-demo"); ok {
		return ex.Render(r)
	}
	if ex, ok := registry.Get("message-scroller-demo"); ok {
		return style.Scope(style.Rhea, ex.Render(r))
	}
	return h.Group()
}
