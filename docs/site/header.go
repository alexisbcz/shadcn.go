package site

import (
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// Extension point: the site chrome around every page of the (app) group
// (app/(app)/layout.tsx): SiteHeader (components/site-header.tsx: MobileNav,
// MainNav, CommandMenu, GitHubLink, ModeSwitcher, the "New" button),
// SiteFooter (components/site-footer.tsx) and MobileNav
// (components/mobile-nav.tsx). The current renderings are rough
// placeholders.

// NavItem is an entry of siteConfig.navItems.
type NavItem struct{ Href, Label string }

// NavItems is siteConfig.navItems (lib/config.ts).
var NavItems = []NavItem{
	{"/", "Home"},
	{"/docs/installation", "Docs"},
	{"/docs/components", "Components"},
	{"/blocks", "Blocks"},
	{"/charts/area", "Charts"},
	{"/docs/directory", "Directory"},
	{"/typeset", "Typeset"},
	{"/create", "Create"},
}

// AppLayout is app/(app)/layout.tsx.
func (s *Site) AppLayout(r *http.Request, children ...any) h.Node {
	return h.Div(h.Data("slot", "layout"),
		h.Class("group/layout relative z-10 flex min-h-svh flex-col bg-background has-data-[slot=designer]:h-svh has-data-[slot=designer]:overflow-hidden"),
		s.SiteHeader(r),
		h.Main(h.Class("flex min-h-0 flex-1 flex-col"), h.Group(children...)),
		s.SiteFooter(r),
	)
}

// SiteHeader is components/site-header.tsx. TODO: MobileNav popover,
// MainNav active state and "new" badges, GitHub star count, separators,
// DesignerActions, exact markup.
func (s *Site) SiteHeader(r *http.Request) h.Node {
	return h.Header(h.Class("sticky top-0 z-50 w-full bg-background"),
		h.Div(h.Class("container-wrapper px-6 group-has-data-[slot=designer]/layout:max-w-none 3xl:fixed:px-0"),
			h.Div(h.Class("flex h-(--header-height) items-center **:data-[slot=separator]:h-4! group-has-data-[slot=designer]/layout:fixed:max-w-none 3xl:fixed:container"),
				s.MobileNav(r),
				s.MainNav(r),
				h.Div(h.Class("ml-auto flex items-center gap-2 md:flex-1 md:justify-end"),
					h.Div(h.Class("hidden w-full flex-1 md:flex md:w-auto md:flex-none"), s.CommandMenu(r)),
					button.Button(button.Variant("ghost"), button.Size("sm"), h.Class("h-8 shadow-none"),
						h.A(h.Href("https://github.com/shadcn-ui/ui"), h.Target("_blank"), h.Rel("noreferrer"), lucide.Github(), h.Span(h.Class("sr-only"), "GitHub")),
					).AsChild(),
					button.Button(button.Variant("ghost"), button.Size("icon"), h.Class("group/toggle extend-touch-target size-8"),
						h.Data("site-action", "toggle-theme"), h.TitleAttr("Toggle theme"),
						lucide.SunMoon(), h.Span(h.Class("sr-only"), "Toggle theme"),
					),
					style.Scope(style.Nova, button.Button(button.Size("sm"), h.Class("h-[31px] rounded-lg"),
						h.A(h.Href("/create"), lucide.Plus(), "New"),
					).AsChild()),
				),
			),
		),
	)
}

// MainNav is components/main-nav.tsx (desktop navigation).
func (s *Site) MainNav(r *http.Request) h.Node {
	return h.Nav(h.Class("items-center gap-0 hidden lg:flex"),
		h.Map(NavItems, func(it NavItem) h.Node {
			return button.Button(button.Variant("ghost"), button.Size("sm"),
				h.A(h.Href(it.Href), h.Data("active", "false"), it.Label),
			).AsChild()
		}),
	)
}

// MobileNav is components/mobile-nav.tsx. TODO: the popover with the nav
// items and the page tree.
func (s *Site) MobileNav(r *http.Request) h.Node {
	return button.Button(button.Variant("ghost"),
		h.Class("extend-touch-target h-8 touch-manipulation items-center justify-start gap-2.5 p-0! hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 active:bg-transparent dark:hover:bg-transparent flex lg:hidden"),
		h.Div(h.Class("relative flex h-8 w-4 items-center justify-center"),
			h.Div(h.Class("relative size-4"),
				h.Span(h.Class("absolute left-0 block h-0.5 w-4 bg-foreground transition-all duration-100 top-1")),
				h.Span(h.Class("absolute left-0 block h-0.5 w-4 bg-foreground transition-all duration-100 top-2.5")),
			),
			h.Span(h.Class("sr-only"), "Toggle Menu"),
		),
		h.Span(h.Class("flex h-8 items-center text-lg leading-none font-medium"), "Menu"),
	)
}

// SiteFooter is components/site-footer.tsx.
func (s *Site) SiteFooter(r *http.Request) h.Node {
	link := func(href, text string) h.Node {
		return h.A(h.Href(href), h.Target("_blank"), h.Rel("noreferrer"), h.Class("font-medium underline underline-offset-4"), text)
	}
	return h.Footer(h.Class("group-has-[.docs-nav]/body:pb-20 group-has-[.section-soft]/body:bg-surface/40 group-has-[[data-slot=designer]]/body:hidden group-has-[[data-slot=docs]]/body:hidden group-has-[.docs-nav]/body:sm:pb-0 dark:bg-transparent dark:group-has-[.section-soft]/body:bg-surface/40 3xl:fixed:bg-transparent"),
		h.Div(h.Class("container-wrapper px-4 xl:px-6"),
			h.Div(h.Class("flex h-(--footer-height) items-center justify-between"),
				h.Div(h.Class("w-full px-1 text-center text-xs leading-loose text-muted-foreground sm:text-sm"),
					"Built by ", link("https://twitter.com/shadcn", "shadcn"), " at ",
					link("https://vercel.com/new?utm_source=shadcn_site&utm_medium=web&utm_campaign=docs_cta_deploy_now_callout", "Vercel"),
					". The source code is available on ", link("https://github.com/shadcn-ui/ui", "GitHub"), ".",
				),
			),
		),
	)
}
