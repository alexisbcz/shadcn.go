package site

import (
	"net/http"
	"strconv"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/button"
	"shadcn.go/docs/content"
	"shadcn.go/docs/site/mdx"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// Extension point: the docs chrome. DocsLayout is app/(app)/docs/layout.tsx
// (SidebarProvider + DocsSidebar), DocsPage is
// app/(app)/docs/[[...slug]]/page.tsx (page header with title, copy-page
// menu and prev/next buttons, the MDX content, the bottom pager, the TOC
// column with OpenInV0Cta). The sidebar, TOC and pager are rough
// placeholders.

// DocsPage renders /docs/[[...slug]].
func (s *Site) DocsPage(r *http.Request) (h.Node, error) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	doc, ok := s.Docs[path]
	if !ok {
		return nil, errNotFound
	}
	page, err := doc.Page()
	if err != nil {
		return nil, err
	}
	env := s.Env(r)
	body, err := content.Render(page, mdx.Registry(env))
	if err != nil {
		return nil, err
	}
	fm := doc.Frontmatter
	var prev, next *pagetree.Node
	if len(doc.Slugs) == 0 || doc.Slugs[0] != "changelog" {
		prev, next = s.Tree.Neighbours(doc.URL)
	}
	toc, err := s.TOC(page, mdx.Registry(env))
	if err != nil {
		return nil, err
	}
	var baseSwitcher h.Node
	if len(doc.Slugs) >= 3 && doc.Slugs[0] == "components" {
		baseSwitcher = s.DocsBaseSwitcher(r, doc.Slugs[1], doc.Slugs[2])
	}
	main := s.DocsArticle(r, DocsArticle{
		Title:       fm.Title,
		Description: fm.Description,
		Prev:        prev,
		Next:        next,
		Before:      baseSwitcher,
		Body:        body,
		TOC:         toc,
	})
	return s.Document(r, Meta{Title: fm.Title, Description: fm.Description, Canonical: doc.URL, Article: true},
		s.AppLayout(r, s.DocsLayout(r, doc.URL, main)),
	), nil
}

// DocsLayout is app/(app)/docs/layout.tsx: the sidebar and the page.
func (s *Site) DocsLayout(r *http.Request, current string, children ...any) h.Node {
	return h.Div(h.Class("container-wrapper flex flex-1 flex-col px-2"),
		h.Div(h.Data("slot", "sidebar-wrapper"),
			h.Style("--sidebar-width:calc(var(--spacing) * 72);--sidebar-width-icon:3rem"),
			h.Class("group/sidebar-wrapper flex w-full has-data-[variant=inset]:bg-sidebar min-h-min flex-1 items-start px-0 [--top-spacing:0] lg:grid lg:grid-cols-[var(--sidebar-width)_minmax(0,1fr)] lg:[--top-spacing:calc(var(--spacing)*4)] 3xl:fixed:container 3xl:fixed:px-3"),
			s.DocsSidebar(r, current),
			h.Div(h.Class("h-full w-full"), h.Group(children...)),
		),
	)
}

// topLevelSections is TOP_LEVEL_SECTIONS of components/docs-sidebar.tsx.
var topLevelSections = []NavItem{
	{"/docs", "Introduction"},
	{"/docs/components", "Components"},
	{"/docs/installation", "Installation"},
	{"/docs/theming", "Theming"},
	{"/docs/cli", "CLI"},
	{"/docs/typeset", "Typeset"},
	{"/docs/skills", "Skills"},
	{"/docs/registry", "Registry"},
	{"/docs/changelog", "Changelog"},
}

var excludedSections = map[string]bool{"installation": true, "dark-mode": true, "changelog": true, "rtl": true}

// DocsSidebar is components/docs-sidebar.tsx. TODO: the new-york-v4
// Sidebar markup, base-dependent component folders, PAGES_NEW badges,
// scroll restoration.
func (s *Site) DocsSidebar(r *http.Request, current string) h.Node {
	item := func(href, label string, active bool) h.Node {
		return h.Li(h.Data("slot", "sidebar-menu-item"), h.Class("group/menu-item relative"),
			h.A(h.Href(href), h.Data("slot", "sidebar-menu-button"), h.Data("active", boolStr(active)),
				h.Class("peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring transition-[width,height,padding] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground h-[30px] text-[0.8rem] font-medium 3xl:fixed:w-full 3xl:fixed:max-w-48 relative w-fit overflow-visible border border-transparent data-[active=true]:border-accent"),
				label),
		)
	}
	group := func(label string, items ...any) h.Node {
		return h.Div(h.Data("slot", "sidebar-group"), h.Class("relative flex w-full min-w-0 flex-col p-2"),
			h.Div(h.Data("slot", "sidebar-group-label"), h.Class("flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium text-muted-foreground"), label),
			h.Div(h.Data("slot", "sidebar-group-content"), h.Class("w-full text-sm"),
				h.Ul(h.Data("slot", "sidebar-menu"), h.Class("flex w-full min-w-0 flex-col gap-0.5"), h.Group(items...)),
			),
		)
	}
	var groups []any
	var sections []any
	for _, sec := range topLevelSections {
		active := current == sec.Href || (sec.Href != "/docs" && strings.HasPrefix(current, sec.Href))
		sections = append(sections, item(sec.Href, sec.Label, active))
	}
	groups = append(groups, h.Div(h.Class("pt-12"), group("Sections", sections...)))
	for _, f := range s.Tree.Children {
		if f.Kind != pagetree.KindFolder || excludedSections[f.Path] {
			continue
		}
		folder := f
		if f.Path == "components" {
			base := "radix"
			if parts := strings.Split(strings.TrimPrefix(current, "/docs/components/"), "/"); len(parts) > 1 && strings.HasPrefix(current, "/docs/components/") {
				base = parts[0]
			}
			if b := f.Find("components/" + base); b != nil {
				folder = b
				folder.Name = "Components"
			}
		}
		var items []any
		for _, p := range folder.Children {
			if p.Kind != pagetree.KindPage || p.URL == "/docs" || p.URL == "/docs/changelog" {
				continue
			}
			items = append(items, item(p.URL, p.Name, p.URL == current))
		}
		if len(items) > 0 {
			groups = append(groups, group(folder.Name, items...))
		}
	}
	return h.Div(h.Data("slot", "sidebar"),
		h.Class("text-sidebar-foreground w-(--sidebar-width) flex-col sticky top-[calc(var(--header-height)+0.6rem)] z-30 hidden h-[calc(100svh-10rem)] overflow-hidden overscroll-none bg-transparent [--sidebar-menu-width:--spacing(56)] lg:flex"),
		h.Div(h.Class("absolute top-12 right-2 bottom-0 hidden h-full w-px bg-[linear-gradient(to_bottom,transparent_0%,var(--border)_10%,var(--border)_90%,transparent_100%)] lg:flex")),
		h.Div(h.Data("slot", "sidebar-content"), h.Data("docs-sidebar-content", ""),
			h.Class("flex min-h-0 flex-1 flex-col gap-2 overflow-auto w-(--sidebar-menu-width) scroll-fade scrollbar-none overflow-x-hidden pl-2.5"),
			h.Group(groups...),
		),
	)
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// DocsArticle is the content column and TOC column of a docs page.
type DocsArticle struct {
	Title, Description string
	// Prev and Next are the neighbouring pages (nil: none).
	Prev, Next *pagetree.Node
	// Actions replace the default header actions (copy page + prev/next),
	// e.g. the changelog's RSS button.
	Actions h.Node
	// Before is rendered at the top of the content (DocsBaseSwitcher).
	Before h.Node
	Body   h.Node
	TOC    h.Node
}

// DocsArticle renders the page layout of app/(app)/docs/[[...slug]]/page.tsx.
func (s *Site) DocsArticle(r *http.Request, a DocsArticle) h.Node {
	navButton := func(p *pagetree.Node, label string, left bool) h.Node {
		if p == nil {
			return nil
		}
		icon := h.Node(lucide.ArrowRight())
		kids := []any{h.Span(h.Class("sr-only"), label), icon}
		if left {
			kids = []any{lucide.ArrowLeft(), h.Span(h.Class("sr-only"), label)}
		}
		return button.Button(button.Variant("secondary"), button.Size("icon"), h.Class("extend-touch-target size-8 shadow-none md:size-7"),
			h.A(append([]any{h.Href(p.URL)}, kids...)...)).AsChild()
	}
	actions := a.Actions
	if actions == nil {
		actions = h.Div(h.Class("docs-nav flex items-center gap-2"),
			h.Div(h.Class("hidden sm:block"), s.DocsCopyPage(r)),
			h.Div(h.Class("ml-auto flex gap-2"), navButton(a.Prev, "Previous", true), navButton(a.Next, "Next", false)),
		)
	}
	var desc h.Node
	if a.Description != "" {
		desc = h.P(h.Class("text-[1.05rem] text-muted-foreground sm:text-base sm:text-balance md:max-w-[80%]"), a.Description)
	}
	var prevLink, nextLink h.Node
	if a.Prev != nil {
		prevLink = button.Button(button.Variant("secondary"), button.Size("sm"), h.Class("shadow-none"),
			h.A(h.Href(a.Prev.URL), lucide.ArrowLeft(), " ", a.Prev.Name)).AsChild()
	}
	if a.Next != nil {
		nextLink = button.Button(button.Variant("secondary"), button.Size("sm"), h.Class("ml-auto shadow-none"),
			h.A(h.Href(a.Next.URL), a.Next.Name, " ", lucide.ArrowRight())).AsChild()
	}
	pager := h.Div(h.Class("hidden h-16 w-full items-center gap-2 px-4 sm:flex sm:px-0"), prevLink, nextLink)
	return h.Div(h.Data("slot", "docs"), h.Class("flex scroll-mt-24 items-stretch pb-8 text-[1.05rem] sm:text-[15px] xl:w-full"),
		h.Div(h.Class("flex min-w-0 flex-1 flex-col"),
			h.Div(h.Class("h-(--top-spacing) shrink-0")),
			h.Div(h.Class("mx-auto flex w-full max-w-160 min-w-0 flex-1 flex-col gap-6 px-4 py-6 text-foreground md:px-0 lg:py-8 dark:text-foreground"),
				h.Div(h.Class("flex flex-col gap-2"),
					h.Div(h.Class("flex flex-col gap-2"),
						h.Div(h.Class("flex items-center justify-between md:items-start"),
							h.H1(h.Class("scroll-m-24 text-3xl font-semibold tracking-tight sm:text-3xl"), a.Title),
							actions,
						),
						desc,
					),
				),
				h.Div(h.Class("typeset w-full flex-1 pb-16 *:data-[slot=alert]:first:mt-0 sm:pb-0"), a.Before, a.Body),
				pager,
			),
		),
		h.Div(h.Class("sticky top-[calc(var(--header-height)+1px)] z-30 ml-auto hidden h-[90svh] w-(--sidebar-width) flex-col gap-4 overflow-hidden overscroll-none pb-8 xl:flex"),
			h.Div(h.Class("h-(--top-spacing) shrink-0")),
			a.TOC,
			h.Div(h.Class("hidden flex-1 flex-col gap-6 px-6 xl:flex"), s.OpenInV0Cta(r)),
		),
	)
}

// TOC is components/docs-toc.tsx inside its column. TODO: scroll-spy
// (data-active), exact markup.
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
		cls := "text-[0.8rem] text-muted-foreground no-underline transition-colors hover:text-foreground data-[active=true]:text-foreground"
		switch it.Depth {
		case 3:
			cls += " pl-4"
		case 4:
			cls += " pl-6"
		}
		items = append(items, h.A(h.Href(it.URL), h.Data("active", "false"), h.Data("depth", itoa(it.Depth)), h.Class(cls), title))
	}
	return h.Div(h.Class("flex scroll-fade scrollbar-none flex-col gap-8 overflow-y-auto px-8"),
		h.Div(h.Class("flex flex-col gap-2 p-4 pt-0 text-sm"),
			h.P(h.Class("sticky top-0 h-6 bg-background text-xs font-medium text-muted-foreground"), "On This Page"),
			h.Group(items...),
		),
	), nil
}

func itoa(i int) string { return strconv.Itoa(i) }

// DocsCopyPage is components/docs-copy-page.tsx. TODO: the button group
// with the dropdown (copy markdown, open in v0/ChatGPT/Claude).
func (s *Site) DocsCopyPage(r *http.Request) h.Node {
	return button.Button(button.Variant("secondary"), button.Size("sm"), h.Class("h-8 shadow-none md:h-7 md:text-[0.8rem]"),
		lucide.Copy(), "Copy Page")
}

// DocsBaseSwitcher is components/docs-base-switcher.tsx (Radix UI / Base
// UI / React Aria tabs above component pages). TODO.
func (s *Site) DocsBaseSwitcher(r *http.Request, base, component string) h.Node {
	return nil
}

// OpenInV0Cta is components/open-in-v0-cta.tsx. TODO.
func (s *Site) OpenInV0Cta(r *http.Request) h.Node {
	return nil
}
