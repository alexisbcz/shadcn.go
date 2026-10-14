package site

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"shadcn.go/components/ui/sonner"
	"strings"

	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

// Site constants (ui/apps/v4/lib/config.ts).
const (
	SiteName        = "shadcn/ui"
	SiteURL         = "https://ui.shadcn.com"
	SiteDescription = "Composable, accessible components with thoughtful defaults. Build your own component library with code you can customize, extend, and make your own."
	ThemeColorLight = "#ffffff"
	ThemeColorDark  = "#09090b"
)

// Meta is a page's metadata (Next.js generateMetadata).
type Meta struct {
	// Title is the page title; the document title is "<Title> - shadcn/ui"
	// ("shadcn/ui" when empty).
	Title       string
	Description string
	// Canonical is the canonical path ("/docs/components/radix/button").
	Canonical string
	// Article marks docs pages (og:type article, og image from the title).
	Article bool
	// PlainTitle drops the " - shadcn/ui" suffix: Next.js applies the root
	// title template only below layouts that do not set a plain title
	// themselves (the blocks layout does).
	PlainTitle bool
}

// fontClasses are the next/font variable classes the root layout puts on
// <html> (lib/fonts.ts); static/fonts/fonts.css, vendored from the
// reference build, defines them.
const fontClasses = "geist_798ff10f-module__HGtJha__variable geist_4e2d599d-module__VB9uPW__variable geist_mono_3f173115-module__-5hxnG__variable noto_sans_arabic_43e4c719-module__-KZrLG__variable noto_sans_hebrew_1a1cc4a4-module__bpro2a__variable"

// preloadFonts are the font files the reference preloads.
var preloadFonts = []string{
	"14a9d8930cac1a04-s.p.0uri4t8jwwgpt.woff2",
	"64c5f24155ea0bd5-s.p.08hbbknvbgrcb.woff2",
	"9d1f839a878dd09a-s.p.2ojwk_d41qe8t.woff2",
	"caa3a2e1cccd8315-s.p.0wgildi0cnwt9.woff2",
	"f1d25eabcf1db66d-s.p.21se6xk6bxqs6.woff2",
}

// headScripts are the inline scripts of the root layout.
type headScripts struct {
	sidebarScroll h.UnsafeJS // DOCS_SIDEBAR_SCROLL_RESTORE_SCRIPT
}

// themeColorScript sets the theme-color meta in dark mode (app/layout.tsx).
const themeColorScript = h.UnsafeJS(`
              try {
                if (localStorage.theme === 'dark' || ((!('theme' in localStorage) || localStorage.theme === 'system') && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
                  document.querySelector('meta[name="theme-color"]').setAttribute('content', '` + ThemeColorDark + `')
                }
              } catch (_) {}
            `)

// themeScript is next-themes' inline script (ThemeProvider: attribute
// "class", storage key "theme", default "system", enableSystem,
// enableColorScheme), verbatim from the reference.
const themeScript = h.UnsafeJS(`((a,b,c,d,e,f,g,h)=>{let i=document.documentElement,j=["light","dark"];function k(b){var c;(Array.isArray(a)?a:[a]).forEach(a=>{let c="class"===a,d=c&&f?e.map(a=>f[a]||a):e;c?(i.classList.remove(...d),i.classList.add(f&&f[b]?f[b]:b)):i.setAttribute(a,b)}),c=b,h&&j.includes(c)&&(i.style.colorScheme=c)}if(d)k(d);else try{let a=localStorage.getItem(b)||c,d=g&&"system"===a?window.matchMedia("(prefers-color-scheme: dark)").matches?"dark":"light":a;k(d)}catch(a){}})("class","theme","system",null,["light","dark"],null,true,true)`)

var templateLiteral = regexp.MustCompile("(?s)DOCS_SIDEBAR_SCROLL_RESTORE_SCRIPT = `(.*?)`")

// loadHeadScripts reads the scripts the root layout inlines from the
// upstream sources.
func loadHeadScripts(app string) (headScripts, error) {
	b, err := os.ReadFile(filepath.Join(app, "lib", "docs-sidebar-scroll.ts"))
	if err != nil {
		return headScripts{}, err
	}
	m := templateLiteral.FindSubmatch(b)
	if m == nil {
		return headScripts{}, fmt.Errorf("site: DOCS_SIDEBAR_SCROLL_RESTORE_SCRIPT not found")
	}
	js := strings.ReplaceAll(string(m[1]), "${DOCS_SIDEBAR_SCROLL_STORAGE_KEY}", "docs-sidebar-scroll")
	// Upstream navigates between docs pages on the client: the sidebar
	// element persists and keeps its scroll position, then the restore logic
	// only scrolls the new active item into view if needed. Here every
	// navigation loads a page, so on in-site navigations the sidebar first
	// takes the stored position of the previous page, as if it had
	// persisted. (Upstream's script runs as is otherwise.)
	const anchor = "    var active = getActiveItem(container)"
	if !strings.Contains(js, anchor) {
		return headScripts{}, fmt.Errorf("site: sidebar restore script changed upstream; update the persistence patch")
	}
	js = strings.Replace(js, anchor, `    if (stored && typeof stored.scrollTop === "number" && !container.dataset.docsScrollPersisted && document.referrer.indexOf(location.origin + "/docs") === 0) {
      container.dataset.docsScrollPersisted = "true"
      container.scrollTop = stored.scrollTop
    }
`+anchor, 1)
	return headScripts{sidebarScroll: h.UnsafeJS(js)}, nil
}

// Document is the root layout (app/layout.tsx): <html> with the font
// classes, <head> with the metadata, fonts, stylesheet and scripts, and
// <body> with the theme script before the content.
func (s *Site) Document(r *http.Request, m Meta, body ...any) h.Node {
	title := SiteName
	if m.Title != "" {
		title = m.Title + " - " + SiteName
		if m.PlainTitle {
			title = m.Title
		}
	}
	desc := m.Description
	if desc == "" {
		desc = SiteDescription
	}
	head := []any{
		h.Meta(h.Charset("utf-8")),
		h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
	}
	for _, f := range preloadFonts {
		head = append(head, h.Link(h.Rel("preload"), h.Href("/fonts/"+f), h.Attr("as", "font"), h.CrossOrigin(""), h.Type("font/woff2")))
	}
	head = append(head,
		h.Link(h.Rel("stylesheet"), h.Href("/fonts/fonts.css")),
		h.Link(h.Rel("stylesheet"), h.Href(cssPath+"?v="+s.css.version())),
		runtime.Scripts(r),
		h.Script(h.Type("module"), h.Src(staticPrefix+"site.js")),
		h.Script(s.scripts.sidebarScroll),
		h.Script(themeColorScript),
		h.Meta(h.Name("theme-color"), h.Content(ThemeColorLight)),
		h.Title(title),
		h.Meta(h.Name("description"), h.Content(desc)),
		h.Meta(h.Name("author"), h.Content("shadcn")),
		h.Link(h.Rel("manifest"), h.Href(SiteURL+"/site.webmanifest")),
		h.Meta(h.Name("keywords"), h.Content("Next.js,React,Tailwind CSS,Components,shadcn")),
		h.Meta(h.Name("creator"), h.Content("shadcn")),
	)
	if m.Canonical != "" {
		head = append(head, h.Link(h.Rel("canonical"), h.Href(SiteURL+m.Canonical)))
	}
	ogTitle := m.Title
	if ogTitle == "" {
		ogTitle = SiteName
	}
	image := SiteURL + "/opengraph-image.png"
	if m.Article {
		image = SiteURL + "/og?title=" + encodeURIComponent(m.Title) + "&description=" + encodeURIComponent(desc)
	}
	ogType := "website"
	if m.Article {
		ogType = "article"
	}
	head = append(head,
		h.Meta(h.Attr("property", "og:title"), h.Content(ogTitle)),
		h.Meta(h.Attr("property", "og:description"), h.Content(desc)),
		h.Meta(h.Attr("property", "og:image"), h.Content(image)),
		h.Meta(h.Attr("property", "og:type"), h.Content(ogType)),
		h.Meta(h.Name("twitter:card"), h.Content("summary_large_image")),
		h.Meta(h.Name("twitter:creator"), h.Content("@shadcn")),
		h.Meta(h.Name("twitter:title"), h.Content(ogTitle)),
		h.Meta(h.Name("twitter:description"), h.Content(desc)),
		h.Meta(h.Name("twitter:image"), h.Content(image)),
		h.Link(h.Rel("shortcut icon"), h.Href("/favicon-16x16.png")),
		h.Link(h.Rel("icon"), h.Href("/favicon.ico")),
		h.Link(h.Rel("apple-touch-icon"), h.Href("/apple-touch-icon.png")),
		h.Link(h.Rel("alternate"), h.Type("application/rss+xml"), h.Href(SiteURL+"/rss.xml")),
	)
	return h.Doctype(h.Html(h.Lang("en"),
		h.Class(fontClasses+" [--header-height:calc(var(--spacing)*14)] lg:[--header-height:calc(var(--spacing)*16)]"),
		h.Head(head...),
		h.Body(
			// theme-default: added on the client by ActiveThemeProvider
			// (components/active-theme.tsx) for the default theme.
			h.Class("group/body overscroll-none antialiased theme-default [--footer-height:calc(var(--spacing)*14)] xl:[--footer-height:calc(var(--spacing)*24)]"),
			h.Script(themeScript),
			h.Group(body...),
			// <Toaster position="top-center" /> (registry/bases/radix/ui/sonner):
			// next-themes reports "light" with system themes disabled, so
			// sonner's own theme stays light; colors follow the page through
			// shadcn's --normal-* variables.
			sonner.Toaster(sonner.ThemeLight, sonner.PositionTopCenter),
		),
	))
}

// encodeURIComponent is JavaScript's encodeURIComponent.
func encodeURIComponent(s string) string {
	e := url.QueryEscape(s)
	e = strings.ReplaceAll(e, "+", "%20")
	for _, c := range []string{"!", "'", "(", ")", "*"} {
		e = strings.ReplaceAll(e, url.QueryEscape(c), c)
	}
	return e
}
