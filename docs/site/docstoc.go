package site

import (
	"net/http"
	"strconv"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// TOC is the TOC column's list: components/docs-toc.tsx (variant "list")
// inside its scroll container. static/docs.js sets data-active (scroll-spy,
// IntersectionObserver with rootMargin "0% 0% -80% 0%", as useActiveItem).
func (s *Site) TOC(page *content.Page, reg *content.Registry) (h.Node, error) {
	if len(page.TOC) == 0 {
		return nil, nil
	}
	var items []any
	for _, it := range page.TOC {
		title, err := content.RenderTOCTitle(page, it, reg)
		if err != nil {
			return nil, err
		}
		items = append(items, h.A(h.Href(it.URL),
			h.Class("text-[0.8rem] text-muted-foreground no-underline transition-colors hover:text-foreground data-[active=true]:font-medium data-[active=true]:text-foreground data-[depth=3]:pl-4 data-[depth=4]:pl-6"),
			h.Data("active", "false"), h.Data("depth", strconv.Itoa(it.Depth)), title))
	}
	return h.Div(h.Class("flex scroll-fade scrollbar-none flex-col gap-8 overflow-y-auto px-8"),
		h.Div(h.Class("flex flex-col gap-2 p-4 pt-0 text-sm"), h.Data("docs-toc", ""),
			h.P(h.Class("h-6 bg-background text-xs font-medium text-muted-foreground"), "On This Page"),
			h.Group(items...),
		),
	), nil
}

// OpenInV0Cta is components/open-in-v0-cta.tsx (the Vercel card under the
// TOC). Its button is base-nova's (Base UI adds type and tabindex).
func (s *Site) OpenInV0Cta(r *http.Request) h.Node {
	return h.Div(h.Class("group relative flex flex-col gap-2 rounded-2xl bg-surface p-6 text-sm text-surface-foreground"),
		h.Div(h.Class("text-base leading-tight font-semibold text-balance group-hover:underline"), "Deploy your shadcn/ui app on Vercel"),
		h.Div(h.Class("text-muted-foreground"), "Trusted by OpenAI, Sonos, Adobe, and more."),
		h.Div(h.Class("text-muted-foreground"), "Vercel provides tools and infrastructure to deploy apps and features at scale."),
		style.Scope(style.Nova, button.Button(button.Variant("outline"), button.Size("sm"), h.Class("mt-2 w-fit"),
			h.Type("button"), h.Attr("tabindex", "0"), "Deploy Now")),
		h.A(h.Href("https://vercel.com/new?utm_source=shadcn_site&utm_medium=web&utm_campaign=docs_cta_deploy_now_callout"),
			h.Target("_blank"), h.Rel("noreferrer"), h.Class("absolute inset-0"),
			h.Span(h.Class("sr-only"), "Deploy to Vercel")),
	)
}
