// Package site is the Go docs site: the routes of ui/apps/v4 (the Next.js
// app of ui.shadcn.com) rendered with the Go components and package
// content. cmd/docs serves it.
//
// Layout of the package (extension points, one concern per file):
//
//	site.go        Site: content loading, routes, static files
//	shell.go       the HTML document (app/layout.tsx): <head>, fonts, theme
//	css.go         the stylesheet (styles.css built by package tailwind)
//	header.go      site header, footer, mobile nav (components/site-*.tsx)
//	docslayout.go  docs sidebar, page header, TOC, pager (app/(app)/docs/*)
//	search.go      search: command menu (components/command-menu.tsx)
//	changelog.go   /docs/changelog
//	home.go        /
//	mdx/           the custom MDX components, one file per component
//	pagetree/      the docs navigation tree (fumadocs' pageTree)
package site

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"shadcn.go/components/ui/runtime"
	"shadcn.go/docs/content"
	"shadcn.go/docs/site/mdx"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// Options configures a Site.
type Options struct {
	// Root is the repository root (holding ui/apps/v4).
	Root string
	// Preload parses every page at startup instead of on first request.
	Preload bool
	// Logf logs progress (nil: log.Printf).
	Logf func(format string, args ...any)
}

// Site serves the docs.
type Site struct {
	Root string
	// App is Root/ui/apps/v4.
	App string
	// Tree is the docs page tree.
	Tree *pagetree.Node
	// Docs are the docs pages by URL ("/docs", "/docs/components/radix/button").
	Docs map[string]*Doc

	css     *stylesheet
	logf    func(string, ...any)
	scripts headScripts
}

// Doc is one docs page: its metadata is read at startup, the page itself is
// parsed (Markdown, highlighting) on first use.
type Doc struct {
	Path        string // relative to content/docs ("components/radix/button.mdx")
	URL         string
	Slugs       []string
	Frontmatter *content.Frontmatter

	site *Site
	once sync.Once
	page *content.Page
	err  error
}

// Page returns the parsed page.
func (d *Doc) Page() (*content.Page, error) {
	d.once.Do(func() {
		b, err := os.ReadFile(filepath.Join(d.site.contentDir(), filepath.FromSlash(d.Path)))
		if err != nil {
			d.err = err
			return
		}
		d.page, d.err = content.ParsePage(d.Path, string(b), content.Options{})
	})
	return d.page, d.err
}

func (s *Site) contentDir() string { return filepath.Join(s.App, "content", "docs") }

// New loads the content and starts building the stylesheet.
func New(opts Options) (*Site, error) {
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return nil, err
	}
	s := &Site{Root: root, App: filepath.Join(root, "ui", "apps", "v4"), Docs: map[string]*Doc{}, logf: opts.Logf}
	if s.logf == nil {
		s.logf = log.Printf
	}
	if s.scripts, err = loadHeadScripts(s.App); err != nil {
		return nil, err
	}
	s.css = newStylesheet(root, s.logf)

	fsys := os.DirFS(s.contentDir())
	var infos []pagetree.PageInfo
	err = fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".mdx") {
			return err
		}
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		matter, _ := content.SplitFrontmatter(string(b))
		fm, err := content.ParseFrontmatter(matter)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		slugs, err := content.Slugs(p)
		if err != nil {
			return err
		}
		url := "/docs"
		if len(slugs) > 0 {
			url += "/" + strings.Join(slugs, "/")
		}
		doc := &Doc{Path: p, URL: url, Slugs: slugs, Frontmatter: fm, site: s}
		s.Docs[url] = doc
		meta := map[string]any{}
		for _, kv := range fm.Values {
			meta[kv.Key] = kv.Value
		}
		infos = append(infos, pagetree.PageInfo{Path: p, URL: url, Title: fm.Title, Meta: meta})
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.Tree = pagetree.Build(fsys, ".", infos)
	if opts.Preload {
		var wg sync.WaitGroup
		sem := make(chan struct{}, 8)
		for _, d := range s.Docs {
			wg.Add(1)
			go func(d *Doc) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				if _, err := d.Page(); err != nil {
					s.logf("docs: %v", err)
				}
			}(d)
		}
		wg.Wait()
	}
	return s, nil
}

// Env returns the environment of MDX components for a request.
func (s *Site) Env(r *http.Request) *mdx.Env {
	return &mdx.Env{Request: r, Tree: s.Tree, Root: s.Root}
}

// Handler returns the site's routes.
func (s *Site) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+cssPath, s.serveCSS)
	mux.Handle("GET /fonts/", http.StripPrefix("/fonts/", fontsHandler()))
	mux.Handle("GET "+staticPrefix, http.StripPrefix(staticPrefix, staticHandler()))
	mux.Handle("GET "+runtime.Prefix, runtime.Handler())
	mux.HandleFunc("GET /{$}", s.page(s.Home))
	mux.HandleFunc("GET /docs/changelog", s.page(s.Changelog))
	mux.HandleFunc("GET /docs", s.page(s.DocsPage))
	mux.HandleFunc("GET /docs/{slug...}", s.page(s.DocsPage))
	mux.HandleFunc("GET "+searchIndexPath, s.serveSearchIndex)
	mux.HandleFunc("GET "+nextImagePath, s.serveNextImage)
	mux.HandleFunc("GET /rss.xml", s.serveRSS)
	mux.HandleFunc("GET /view/{style}/{name}", s.page(s.View))
	mux.HandleFunc("GET /blocks", s.page(s.Blocks))
	mux.HandleFunc("GET /blocks/{categories...}", s.page(s.BlocksCategory))
	mux.HandleFunc("GET /charts/{type}", s.page(s.Charts))
	mux.HandleFunc("GET /colors", s.page(s.Colors))
	mux.HandleFunc("GET /examples/dashboard", s.page(s.ExampleDashboard))
	mux.HandleFunc("GET /examples/tasks", s.page(s.ExampleTasks))
	mux.HandleFunc("GET /examples/playground", s.page(s.ExamplePlayground))
	mux.HandleFunc("GET /examples/authentication", s.page(s.ExampleAuthentication))
	mux.HandleFunc("GET /examples/rtl", s.page(s.ExampleRTL))
	// Everything else: the app's public directory (favicons, images, /r/
	// registry files, llms.txt, ...), as Next.js serves it.
	public := http.FileServer(http.Dir(filepath.Join(s.App, "public")))
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(filepath.Join(s.App, "public", filepath.FromSlash(filepath.Clean("/"+r.URL.Path)))); err != nil {
			s.notFound(w, r)
			return
		}
		public.ServeHTTP(w, r)
	}))
	return mux
}

// errNotFound is returned by page functions for unknown pages.
var errNotFound = fmt.Errorf("not found")

// page adapts a page function (which returns the whole document) to a
// handler.
func (s *Site) page(fn func(r *http.Request) (h.Node, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n, err := fn(r)
		if err == errNotFound {
			s.notFound(w, r)
			return
		}
		if err != nil {
			s.logf("%s: %v", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.WriteHTTP(w, http.StatusOK, n); err != nil {
			s.logf("%s: %v", r.URL.Path, err)
		}
	}
}

// notFound renders the 404 page: the app has no not-found.tsx, so Next.js
// renders its default one inside the root layout.
func (s *Site) notFound(w http.ResponseWriter, r *http.Request) {
	const font = `system-ui,"Segoe UI",Roboto,Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji"`
	// The document title stays the root one: Next.js hoists the error page's
	// <title> before the metadata's, which wins.
	doc := s.Document(r, Meta{},
		h.Div(h.UnsafeAttr("style", "font-family:"+font+";height:100vh;text-align:center;display:flex;flex-direction:column;align-items:center;justify-content:center"),
			h.Div(
				h.StyleEl(h.UnsafeCSS("body{color:#000;background:#fff;margin:0}.next-error-h1{border-right:1px solid rgba(0,0,0,.3)}@media (prefers-color-scheme:dark){body{color:#fff;background:#000}.next-error-h1{border-right:1px solid rgba(255,255,255,.3)}}")),
				h.H1(h.Class("next-error-h1"), h.Style("display:inline-block;margin:0 20px 0 0;padding:0 23px 0 0;font-size:24px;font-weight:500;vertical-align:top;line-height:49px"), "404"),
				h.Div(h.Style("display:inline-block"),
					h.H2(h.Style("font-size:14px;font-weight:400;line-height:49px;margin:0"), "This page could not be found."),
				),
			),
		))
	if err := h.WriteHTTP(w, http.StatusNotFound, doc); err != nil {
		s.logf("%s: %v", r.URL.Path, err)
	}
}
