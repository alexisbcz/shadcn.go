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

// Base is an entry of BASES (registry/bases.ts).
type Base struct{ Name, Title, Logo string }

// Bases is BASES of registry/bases.ts, in order.
var Bases = []Base{
	{"base", "Base UI", "<svg width='17' height='24' viewBox='0 0 17 24'><path fill='currentColor' d='M9.5001 7.01537C9.2245 6.99837 9 7.22385 9 7.49999V23C13.4183 23 17 19.4183 17 15C17 10.7497 13.6854 7.27351 9.5001 7.01537Z'></path><path fill='currentColor'   d='M8 9.8V12V23C3.58172 23 0 19.0601 0 14.2V12V1C4.41828 1 8 4.93989 8 9.8Z'></path></svg>"},
	{"aria", "React Aria", "<svg xmlns='http://www.w3.org/2000/svg' viewBox='200 206 800 790' fill='none'><path d='M720.67 205.995C867.583 205.995 986.679 325.091 986.68 472.003C986.68 590.753 908.865 691.325 801.446 725.521L979.312 948.055C994.438 966.98 980.963 995 956.736 995H795.612C778.743 995 762.715 987.629 751.734 974.823L697.365 911.421L493.126 653.39C457.134 607.918 489.518 540.979 547.511 540.977L720.67 540.971C758.758 540.971 789.635 510.091 789.635 472.003C789.634 433.915 758.758 403.038 720.67 403.038H429.939C404.955 403.038 388.623 391.886 373.994 373.623L277.349 252.966C262.194 234.045 275.664 205.996 299.905 205.995H720.67Z M396.605 720.706C407.798 705.406 430.443 704.843 442.381 719.568L503.816 797.018H502.786L535.569 838.934C548.074 854.358 549.943 877.191 538.047 893.09L476.638 972.545C465.692 986.707 448.803 995 430.903 995H242.276C218.18 995 204.665 967.248 219.523 948.278L337.992 797.018H337.923L396.605 720.706Z' fill='currentColor' /></svg>"},
	{"radix", "Radix UI", "<svg role='img' viewBox='0 0 24 24' xmlns='http://www.w3.org/2000/svg'><title>Radix UI</title><path fill='currentColor' d='M11.52 24a7.68 7.68 0 0 1-7.68-7.68 7.68 7.68 0 0 1 7.68-7.68V24Zm0-24v7.68H3.84V0h7.68Zm4.8 7.68a3.84 3.84 0 1 1 0-7.68 3.84 3.84 0 0 1 0 7.68Z'/></svg>"},
}

// DocsBaseSwitcher is components/docs-base-switcher.tsx: the Base UI /
// React Aria / Radix UI tabs above a component page (className "mb-4").
func (s *Site) DocsBaseSwitcher(r *http.Request, base, component string) h.Node {
	var links []any
	var logo string
	for _, b := range Bases {
		if b.Name == base {
			logo = b.Logo
		}
		href := "/docs/components/" + b.Name + "/" + component
		if _, ok := s.Docs[href]; !ok {
			continue
		}
		links = append(links, h.A(h.Href(href), h.Data("active", strconv.FormatBool(b.Name == base)),
			h.Class("relative inline-flex items-center justify-center gap-1 pt-1 pb-0.5 text-base font-medium text-muted-foreground transition-colors after:absolute after:inset-x-0 after:bottom-[-4px] after:h-0.5 after:bg-foreground after:opacity-0 after:transition-opacity hover:text-foreground data-[active=true]:text-foreground data-[active=true]:after:opacity-100"),
			b.Title))
	}
	var logoNode h.Node
	if logo != "" {
		logoNode = h.Div(h.Class("ml-auto size-4 shrink-0 text-muted-foreground opacity-80 [&_svg]:size-4"), h.UnsafeHTML(logo))
	}
	return h.Div(h.Class("not-typeset inline-flex w-full items-center gap-6 mb-4"), h.Group(links...), logoNode)
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
