package site

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/docs/content"
	"shadcn.go/docs/site/mdx"
	h "shadcn.go/libhtml"
)

// latestChangelogPages is NUMBER_OF_LATEST_PAGES.
const latestChangelogPages = 5

// ChangelogPages returns the changelog entries (lib/changelog.ts
// getChangelogPages): pages under changelog/, newest first.
func (s *Site) ChangelogPages() []*Doc {
	var out []*Doc
	for _, d := range s.Docs {
		if len(d.Slugs) > 1 && d.Slugs[0] == "changelog" {
			out = append(out, d)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].Frontmatter, out[j].Frontmatter
		if a.HasDate != b.HasDate {
			return a.HasDate
		}
		if !a.Date.Equal(b.Date) {
			return a.Date.After(b.Date)
		}
		return out[i].Path < out[j].Path
	})
	return out
}

// Changelog renders /docs/changelog (app/(app)/docs/changelog/page.tsx).
func (s *Site) Changelog(r *http.Request) (h.Node, error) {
	pages := s.ChangelogPages()
	latest, older := pages, []*Doc(nil)
	if len(pages) > latestChangelogPages {
		latest, older = pages[:latestChangelogPages], pages[latestChangelogPages:]
	}
	reg := mdx.Registry(s.Env(r))
	var articles, tocLinks []any
	for _, d := range latest {
		p, err := d.Page()
		if err != nil {
			return nil, err
		}
		body, err := content.Render(p, reg)
		if err != nil {
			return nil, err
		}
		articles = append(articles, h.El("article", h.Class("mb-12 border-b pb-12"),
			h.H2(h.Class("font-heading text-xl font-semibold tracking-tight"), d.Frontmatter.Title),
			h.Div(h.Class("typeset mt-6 *:first:mt-0"), body),
		))
		tocLinks = append(tocLinks, h.A(h.Href(d.URL), h.Class("text-[0.8rem] text-muted-foreground no-underline transition-colors hover:text-foreground"), d.Frontmatter.Title))
	}
	if len(older) > 0 {
		articles = append(articles, h.Div(h.ID("more-updates"), h.Class("mb-24 scroll-mt-24"),
			h.H2(h.Class("mb-6 font-heading text-xl font-semibold tracking-tight"), "More Updates"),
			h.Div(h.Class("grid auto-rows-fr gap-3 sm:grid-cols-2"),
				h.Map(older, func(d *Doc) h.Node {
					date, title, _ := strings.Cut(d.Frontmatter.Title, " - ")
					return h.A(h.Href(d.URL), h.Class("flex w-full flex-col rounded-2xl bg-surface px-4 py-3 text-surface-foreground transition-colors hover:bg-surface/80"),
						h.Span(h.Class("text-xs text-muted-foreground"), date),
						h.Span(h.Class("text-sm font-medium"), title),
					)
				}),
			),
		))
		tocLinks = append(tocLinks, h.A(h.Href("#more-updates"), h.Class("text-[0.8rem] text-muted-foreground no-underline transition-colors hover:text-foreground"), "More Updates"))
	}
	rss := style.Scope(style.Nova, button.Button(button.Variant("secondary"), button.Size("sm"),
		h.A(h.Href("/rss.xml"), h.Target("_blank"), h.Rel("noopener noreferrer"), lucide.Rss(), "RSS")).AsChild())
	main := h.Div(h.Data("slot", "docs"), h.Class("flex scroll-mt-24 items-stretch pb-8 text-[1.05rem] sm:text-[15px] xl:w-full"),
		h.Div(h.Class("flex min-w-0 flex-1 flex-col"),
			h.Div(h.Class("h-(--top-spacing) shrink-0")),
			h.Div(h.Class("mx-auto flex w-full max-w-160 min-w-0 flex-1 flex-col gap-6 px-4 py-6 text-foreground md:px-0 lg:py-8 dark:text-foreground"),
				h.Div(h.Class("flex flex-col gap-2"),
					h.Div(h.Class("flex items-center justify-between"),
						h.H1(h.Class("scroll-m-24 text-4xl font-semibold tracking-tight sm:text-3xl"), "Changelog"),
						rss,
					),
					h.P(h.Class("text-[1.05rem] text-muted-foreground sm:text-base sm:text-balance md:max-w-[80%]"), "Latest updates and announcements."),
				),
				h.Div(h.Class("w-full flex-1 pb-16 sm:pb-0"), h.Group(articles...)),
			),
		),
		h.Div(h.Class("sticky top-[calc(var(--header-height)+1px)] z-30 ml-auto hidden h-[90svh] w-72 flex-col gap-4 overflow-hidden overscroll-none pb-8 lg:flex"),
			h.Div(h.Class("h-(--top-spacing) shrink-0")),
			h.Div(h.Class("no-scrollbar flex flex-col gap-8 overflow-y-auto px-8"),
				h.Div(h.Class("flex flex-col gap-2 p-4 pt-0 text-sm"),
					h.P(h.Class("sticky top-0 h-6 bg-background text-xs font-medium text-muted-foreground"), "On This Page"),
					h.Group(tocLinks...),
				),
			),
			h.Div(h.Class("hidden flex-1 flex-col gap-6 px-6 xl:flex"), s.OpenInV0Cta(r)),
		),
	)
	return s.Document(r, Meta{Title: "Changelog", Description: "Latest updates and announcements.", Article: true},
		s.AppLayout(r, s.DocsLayout(r, "/docs/changelog", main)),
	), nil
}

// serveRSS serves /rss.xml (app/rss.xml/route.ts).
func (s *Site) serveRSS(w http.ResponseWriter, r *http.Request) {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>` + SiteName + ` Changelog</title>
    <link>` + SiteURL + `</link>
    <description>` + SiteDescription + `</description>
    <language>en-us</language>
    <atom:link href="` + SiteURL + `/rss.xml" rel="self" type="application/rss+xml"/>
`)
	for i, d := range s.ChangelogPages() {
		if i > 0 {
			b.WriteString("\n")
		}
		date := time.Now()
		if d.Frontmatter.HasDate {
			date = d.Frontmatter.Date
		}
		link := SiteURL + d.URL
		b.WriteString(`    <item>
      <title><![CDATA[` + d.Frontmatter.Title + `]]></title>
      <link>` + link + `</link>
      <guid>` + link + `</guid>
      <description><![CDATA[` + d.Frontmatter.Description + `]]></description>
      <pubDate>` + date.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") + `</pubDate>
    </item>`)
	}
	b.WriteString("\n  </channel>\n</rss>")
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	_, _ = w.Write([]byte(b.String()))
}
