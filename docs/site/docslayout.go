package site

import (
	"io/fs"
	"net/http"
	"strings"

	"shadcn.go/components/ui/button"
	"shadcn.go/docs/content"
	"shadcn.go/docs/site/mdx"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// The docs chrome. DocsLayout is app/(app)/docs/layout.tsx (SidebarProvider
// + DocsSidebar, docssidebar.go), DocsPage is
// app/(app)/docs/[[...slug]]/page.tsx: the page header (title, description,
// DocsCopyPage in docscopypage.go, prev/next buttons), the MDX
// content, the bottom pager, and the TOC column (TOC and
// OpenInV0Cta, docstoc.go). Client behaviors (TOC scroll-spy, copy page,
// sidebar scroll saving) are in static/docs.js.

// DocsPage renders /docs/[[...slug]].
func (s *Site) DocsPage(r *http.Request) (h.Node, error) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	doc, ok := s.Docs[path]
	if !ok {
		return s.docsRedirect(r)
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
	raw, _ := fs.ReadFile(s.content, doc.Path)
	main := s.DocsArticle(r, DocsArticle{
		Title:       fm.Title,
		Description: fm.Description,
		URL:         SiteURL + doc.URL,
		Raw: mdx.ReplaceComponentsList(s.Tree, string(raw), func(url string) string {
			if d := s.Docs[url]; d != nil {
				return d.Frontmatter.Description
			}
			return ""
		}),
		Prev: prev,
		Next: next,
		Body: body,
		TOC:  toc,
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
		h.Script(h.Type("module"), h.Src(staticPrefix+"docs.js")),
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
	// URL is the page's absolute URL and Raw its source, for DocsCopyPage.
	URL, Raw string
	// Prev and Next are the neighbouring pages (nil: none).
	Prev, Next *pagetree.Node
	// Actions replace the default header actions (copy page + prev/next),
	// e.g. the changelog's RSS button.
	Actions h.Node
	// Before is rendered at the top of the content.
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
		kids := []any{h.Href(p.URL), h.Span(h.Class("sr-only"), label), iconArrowRight()}
		if left {
			kids = []any{h.Href(p.URL), iconArrowLeft(), h.Span(h.Class("sr-only"), label)}
		}
		return button.Button(button.Variant("secondary"), button.Size("icon"), h.Class("extend-touch-target size-8 shadow-none md:size-7"),
			h.A(kids...)).AsChild()
	}
	actions := a.Actions
	if actions == nil {
		actions = h.Div(h.Class("docs-nav flex items-center gap-2"),
			h.Div(h.Class("hidden sm:block"), s.DocsCopyPage(r, a.URL, a.Raw)),
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
			h.A(h.Href(a.Prev.URL), iconArrowLeft(), " ", a.Prev.Name)).AsChild()
	}
	if a.Next != nil {
		nextLink = button.Button(button.Variant("secondary"), button.Size("sm"), h.Class("ml-auto shadow-none"),
			h.A(h.Href(a.Next.URL), a.Next.Name, " ", iconArrowRight())).AsChild()
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
