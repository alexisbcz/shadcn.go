package site

import (
	"net/http"

	blocksdashboard "shadcn.go/examples/blocks-dashboard"
	h "shadcn.go/libhtml"
)

// ExampleDashboard serves /examples/dashboard (app/(app)/examples/dashboard):
// screenshots below md, the dashboard-01 variant above.
func (s *Site) ExampleDashboard(r *http.Request) (h.Node, error) {
	img := func(theme, class string) h.Node {
		u := "%2Fexamples%2Fdashboard-" + theme + ".png"
		return h.Img(h.Alt("Authentication"), h.Width("1280"), h.Height("843"),
			h.Attr("decoding", "async"), h.Data("nimg", "1"), h.Class(class), h.Style("color:transparent"),
			h.UnsafeAttr("srcset", "/_next/image?url="+u+"&w=1920&q=75 1x, /_next/image?url="+u+"&w=3840&q=75 2x"),
			h.UnsafeAttr("src", "/_next/image?url="+u+"&w=3840&q=75"))
	}
	return s.examplesLayout(r, "/examples/dashboard",
		h.Div(h.Class("md:hidden"),
			img("light", "block dark:hidden"),
			img("dark", "hidden dark:block"),
		),
		blocksdashboard.ExamplePage(r, h.Class("hidden md:flex")),
	), nil
}
