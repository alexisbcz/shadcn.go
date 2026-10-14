package site

import (
	"net/http"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/separator"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// The site chrome around every page of the (app) group
// (app/(app)/layout.tsx): SiteHeader (components/site-header.tsx: MobileNav,
// MainNav, CommandMenu, GitHubLink, ModeSwitcher, the "New" button),
// SiteFooter (components/site-footer.tsx) and MobileNav
// (components/mobile-nav.tsx).

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

// GitHubStars is the star count the reference build bakes in
// (components/github-link.tsx fetches it at build time; DOCS_PLAN.md D0).
const GitHubStars = "126k"

// headerPagesNew is PAGES_NEW of lib/docs.ts.
var headerPagesNew = map[string]bool{
	"/docs/changelog":                       true,
	"/docs/changelog/2026-08-questionnaire": true,
	"/docs/registry/dynamic-search":         true,
	"/docs/components/radix/questionnaire":  true,
	"/docs/components/base/questionnaire":   true,
	"/docs/components/aria/questionnaire":   true,
	"/docs/react/questionnaire":             true,
	"/docs/registry/health":                 true,
}

// mobileSections is TOP_LEVEL_SECTIONS of components/mobile-nav.tsx.
var mobileSections = []NavItem{
	{"/docs", "Introduction"},
	{"/docs/components", "Components"},
	{"/docs/installation", "Installation"},
	{"/docs/theming", "Theming"},
	{"/docs/cli", "CLI"},
	{"/docs/rtl", "RTL"},
	{"/docs/skills", "Skills"},
	{"/docs/mcp", "MCP Server"},
	{"/docs/registry", "Registry"},
	{"/docs/forms", "Forms"},
	{"/docs/changelog", "Changelog"},
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

func vsep(class ...string) h.Node {
	return separator.Separator(separator.Vertical, h.Class(strings.Join(class, " ")))
}

// SiteHeader is components/site-header.tsx. DesignerActions renders only on
// designer pages (/create), so its hidden placeholders are kept for layout.
func (s *Site) SiteHeader(r *http.Request) h.Node {
	return h.Header(h.Class("sticky top-0 z-50 w-full bg-background"),
		h.Div(h.Class("container-wrapper px-6 group-has-data-[slot=designer]/layout:max-w-none 3xl:fixed:px-0"),
			h.Div(h.Class("flex h-(--header-height) items-center **:data-[slot=separator]:h-4! group-has-data-[slot=designer]/layout:fixed:max-w-none 3xl:fixed:container"),
				s.MobileNav(r),
				s.MainNav(r),
				h.Div(h.Class("ml-auto flex items-center gap-2 md:flex-1 md:justify-end"),
					h.Div(h.Class("hidden w-full flex-1 md:flex md:w-auto md:flex-none"), s.CommandMenu(r)),
					vsep("ml-2 hidden lg:block"),
					GitHubLink(),
					vsep(),
					ModeSwitcher(),
					h.Div(h.Class("flex items-center gap-2 group-has-data-[slot=designer]/layout:hidden"),
						vsep(),
						style.Scope(style.Nova, button.Button(button.SizeSm, h.Class("h-[31px] rounded-lg"), lib.AsChild,
							h.A(h.Href("/create"), plusSignIcon(), "New"),
						)),
					),
				),
			),
		),
	)
}

// GitHubLink is components/github-link.tsx.
func GitHubLink() h.Node {
	return button.Button(button.SizeSm, button.VariantGhost, h.Class("h-8 shadow-none"), lib.AsChild,
		h.A(h.Href("https://github.com/shadcn-ui/ui"), h.Target("_blank"), h.Rel("noreferrer"),
			h.El("svg", h.Attr("viewBox", "0 0 438.549 438.549"),
				h.Path(h.Attr("fill", "currentColor"), h.Attr("d", githubPath)),
			),
			h.Span(h.Class("w-fit text-xs text-muted-foreground tabular-nums"), GitHubStars),
		),
	)
}

// ModeSwitcher is components/mode-switcher.tsx; site.js toggles the theme.
func ModeSwitcher() h.Node {
	path := func(d string) h.Node { return h.Path(h.Attr("d", d)) }
	return button.Button(button.VariantGhost, button.SizeIcon, h.Class("group/toggle extend-touch-target size-8"),
		h.Data("site-action", "toggle-theme"),
		h.El("svg", h.Attr("xmlns", "http://www.w3.org/2000/svg"), h.Attr("width", "24"), h.Attr("height", "24"),
			h.Attr("viewBox", "0 0 24 24"), h.Attr("fill", "none"), h.Attr("stroke", "currentColor"),
			h.Attr("stroke-width", "2"), h.Attr("stroke-linecap", "round"), h.Attr("stroke-linejoin", "round"),
			h.Class("size-4.5"),
			h.Path(h.Attr("stroke", "none"), h.Attr("d", "M0 0h24v24H0z"), h.Attr("fill", "none")),
			path("M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"),
			path("M12 3l0 18"),
			path("M12 9l4.65 -4.65"),
			path("M12 14.3l7.37 -7.37"),
			path("M12 19.6l8.85 -8.85"),
		),
		h.Span(h.Class("sr-only"), "Toggle theme"),
	)
}

// plusSignIcon is hugeicons' PlusSignIcon as HugeiconsIcon renders it.
func plusSignIcon() h.Node {
	return h.El("svg", h.Attr("xmlns", "http://www.w3.org/2000/svg"), h.Attr("width", "24"), h.Attr("height", "24"),
		h.Attr("viewBox", "0 0 24 24"), h.Attr("fill", "none"), h.Attr("color", "currentColor"), h.Class(""),
		h.Path(h.Attr("d", "M12 4V20M20 12H4"), h.Attr("stroke", "currentColor"), h.Attr("stroke-linecap", "round"),
			h.Attr("stroke-linejoin", "round"), h.Attr("stroke-width", "1.5")),
	)
}

// MainNav is components/main-nav.tsx (desktop navigation).
func (s *Site) MainNav(r *http.Request) h.Node {
	return h.Nav(h.Class("items-center gap-0 hidden lg:flex"),
		h.Map(NavItems, func(it NavItem) h.Node {
			return button.Button(button.VariantGhost, button.SizeSm, h.Class("px-2.5"), lib.AsChild,
				h.A(h.Href(it.Href), h.Data("active", boolStr(r.URL.Path == it.Href)),
					h.Data("new", boolStr(headerPagesNew[it.Href])), h.Class("relative items-center"), it.Label),
			)
		}),
	)
}

// MobileNav is components/mobile-nav.tsx: the "Menu" button below lg and
// its full-screen popover with the nav items, the sections and the page
// tree. The hamburger lines turn into a cross while the trigger is open
// (data-state=open).
func (s *Site) MobileNav(r *http.Request) h.Node {
	line := "absolute left-0 block h-0.5 w-4 bg-foreground transition-all duration-100 in-data-[state=open]:top-[0.4rem]"
	link := func(href string, class string, children ...any) h.Node {
		return h.A(h.Class(lib.Cn("flex items-center gap-2 text-2xl font-medium", class)), h.Href(href), h.Group(children...))
	}
	dot := func(href string, title bool) h.Node {
		if !headerPagesNew[href] {
			return nil
		}
		return h.Span(h.Class("flex size-2 rounded-full bg-blue-500"), h.If(title, h.TitleAttr("New")))
	}
	group := func(gap, label string, links ...any) h.Node {
		return h.Div(h.Class("flex flex-col gap-4"),
			h.Div(h.Class("text-sm font-medium text-muted-foreground"), label),
			h.Div(h.Class("flex flex-col "+gap), h.Group(links...)),
		)
	}
	var items, sections, tree []any
	for _, it := range NavItems {
		items = append(items, link(it.Href, "", it.Label))
	}
	for _, it := range mobileSections {
		sections = append(sections, link(it.Href, "", it.Label, dot(it.Href, true)))
	}
	base := mobileCurrentBase(r.URL.Path)
	if s.Tree != nil {
		for _, f := range s.Tree.Children {
			if f.Kind != pagetree.KindFolder {
				continue
			}
			var links []any
			for _, p := range mobileFolderPages(f, base) {
				// changelog/meta.json lists no pages ("pages": []): fumadocs
				// gives the folder no children, the pagetree lists them all.
				if f.Path == "changelog" {
					break
				}
				links = append(links, h.A(h.Class("text-2xl font-medium flex items-center gap-2"), h.Href(p.URL), p.Name, " ", dot(p.URL, false)))
			}
			tree = append(tree, group("gap-3", f.Name, links...))
		}
	}
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantGhost,
			h.Class("extend-touch-target h-8 touch-manipulation items-center justify-start gap-2.5 p-0! hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 active:bg-transparent dark:hover:bg-transparent flex lg:hidden"),
			h.Div(h.Class("relative flex h-8 w-4 items-center justify-center"),
				h.Div(h.Class("relative size-4"),
					h.Span(h.Class(line+" top-1 in-data-[state=open]:-rotate-45")),
					h.Span(h.Class(line+" top-2.5 in-data-[state=open]:rotate-45")),
				),
				h.Span(h.Class("sr-only"), "Toggle Menu"),
			),
			h.Span(h.Class("flex h-8 items-center text-lg leading-none font-medium"), "Menu"),
		)),
		popover.Content(popover.SideBottom, popover.AlignStart, popover.AlignOffset(-16), popover.SideOffset(14),
			h.Class("no-scrollbar h-(--radix-popper-available-height) w-(--radix-popper-available-width) overflow-y-auto rounded-none border-none bg-background/90 p-0 shadow-none backdrop-blur duration-100 data-open:animate-none!"),
			h.Div(h.Class("flex flex-col gap-12 overflow-auto px-6 py-6"),
				group("gap-3", "Menu", items...),
				group("gap-3", "Sections", sections...),
				h.Div(h.Class("flex flex-col gap-8"), h.Group(tree...)),
			),
		),
	)
}

// mobileCurrentBase is getCurrentBase of lib/page-tree.ts.
func mobileCurrentBase(path string) string {
	const p = "/docs/components/"
	if rest, ok := strings.CutPrefix(path, p); ok {
		if i := strings.Index(rest, "/"); i > 0 {
			switch b := rest[:i]; b {
			case "radix", "base", "aria":
				return b
			}
		}
	}
	return "base"
}

// mobileFolderPages is getPagesFromFolder of lib/page-tree.ts.
func mobileFolderPages(folder *pagetree.Node, base string) []*pagetree.Node {
	direct := func(f *pagetree.Node) (out []*pagetree.Node) {
		for _, c := range f.Children {
			if c.Kind == pagetree.KindPage {
				out = append(out, c)
			}
		}
		return out
	}
	if folder.Path == "components" || folder.Name == "Components" {
		names := map[string]string{"radix": "Radix UI", "base": "Base UI", "aria": "React Aria"}
		for _, c := range folder.Children {
			if c.Kind == pagetree.KindFolder && (strings.HasSuffix(c.Path, "/"+base) || c.Name == names[base]) {
				return direct(c)
			}
		}
		var all []*pagetree.Node
		var walk func(*pagetree.Node)
		walk = func(f *pagetree.Node) {
			for _, c := range f.Children {
				if c.Kind == pagetree.KindPage && !strings.HasSuffix(c.URL, "/components") {
					all = append(all, c)
				} else if c.Kind == pagetree.KindFolder {
					walk(c)
				}
			}
		}
		walk(folder)
		return all
	}
	return direct(folder)
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

// githubPath is Icons.gitHub of components/icons.tsx.
const githubPath = "M409.132 114.573c-19.608-33.596-46.205-60.194-79.798-79.8-33.598-19.607-70.277-29.408-110.063-29.408-39.781 0-76.472 9.804-110.063 29.408-33.596 19.605-60.192 46.204-79.8 79.8C9.803 148.168 0 184.854 0 224.63c0 47.78 13.94 90.745 41.827 128.906 27.884 38.164 63.906 64.572 108.063 79.227 5.14.954 8.945.283 11.419-1.996 2.475-2.282 3.711-5.14 3.711-8.562 0-.571-.049-5.708-.144-15.417a2549.81 2549.81 0 01-.144-25.406l-6.567 1.136c-4.187.767-9.469 1.092-15.846 1-6.374-.089-12.991-.757-19.842-1.999-6.854-1.231-13.229-4.086-19.13-8.559-5.898-4.473-10.085-10.328-12.56-17.556l-2.855-6.57c-1.903-4.374-4.899-9.233-8.992-14.559-4.093-5.331-8.232-8.945-12.419-10.848l-1.999-1.431c-1.332-.951-2.568-2.098-3.711-3.429-1.142-1.331-1.997-2.663-2.568-3.997-.572-1.335-.098-2.43 1.427-3.289 1.525-.859 4.281-1.276 8.28-1.276l5.708.853c3.807.763 8.516 3.042 14.133 6.851 5.614 3.806 10.229 8.754 13.846 14.842 4.38 7.806 9.657 13.754 15.846 17.847 6.184 4.093 12.419 6.136 18.699 6.136 6.28 0 11.704-.476 16.274-1.423 4.565-.952 8.848-2.383 12.847-4.285 1.713-12.758 6.377-22.559 13.988-29.41-10.848-1.14-20.601-2.857-29.264-5.14-8.658-2.286-17.605-5.996-26.835-11.14-9.235-5.137-16.896-11.516-22.985-19.126-6.09-7.614-11.088-17.61-14.987-29.979-3.901-12.374-5.852-26.648-5.852-42.826 0-23.035 7.52-42.637 22.557-58.817-7.044-17.318-6.379-36.732 1.997-58.24 5.52-1.715 13.706-.428 24.554 3.853 10.85 4.283 18.794 7.952 23.84 10.994 5.046 3.041 9.089 5.618 12.135 7.708 17.705-4.947 35.976-7.421 54.818-7.421s37.117 2.474 54.823 7.421l10.849-6.849c7.419-4.57 16.18-8.758 26.262-12.565 10.088-3.805 17.802-4.853 23.134-3.138 8.562 21.509 9.325 40.922 2.279 58.24 15.036 16.18 22.559 35.787 22.559 58.817 0 16.178-1.958 30.497-5.853 42.966-3.9 12.471-8.941 22.457-15.125 29.979-6.191 7.521-13.901 13.85-23.131 18.986-9.232 5.14-18.182 8.85-26.84 11.136-8.662 2.286-18.415 4.004-29.263 5.146 9.894 8.562 14.842 22.077 14.842 40.539v60.237c0 3.422 1.19 6.279 3.572 8.562 2.379 2.279 6.136 2.95 11.276 1.995 44.163-14.653 80.185-41.062 108.068-79.226 27.88-38.161 41.825-81.126 41.825-128.906-.01-39.771-9.818-76.454-29.414-110.049z"
