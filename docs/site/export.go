package site

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Export writes the site as static files to dir, for a static host such as
// GitHub Pages serving it under basePath ("/shadcn.go", or "" at the root).
//
// It crawls the routes from "/" and every docs page, following the
// root-relative URLs of the HTML, CSS and JavaScript it gets. A page is
// written as <path>.html (and <path>/index.html when other pages live
// under it); query strings are dropped from links to pages and folded into
// the file name of other resources (/_next/image?url=...). Every
// root-relative URL is then prefixed with basePath.
//
// What needs the server does not work on the static site: full-text
// search (/api/search) and the query parameters of /create and /typeset.
func (s *Site) Export(dir, basePath string) error {
	basePath = strings.TrimSuffix(basePath, "/")
	e := &exporter{
		s: s, h: s.Handler(), dir: dir, base: basePath,
		seen: map[string]bool{}, mapped: map[string]string{}, text: map[string][]byte{},
		pages: map[string]bool{},
	}
	e.segs = map[string]bool{}
	for _, seg := range routeSegments {
		e.segs[seg] = true
	}
	if entries, err := os.ReadDir(filepath.Join(s.App, "public")); err == nil {
		for _, en := range entries {
			e.segs[en.Name()] = true
		}
	}
	var segs []string
	for seg := range e.segs {
		segs = append(segs, regexp.QuoteMeta(seg))
	}
	sort.Slice(segs, func(i, j int) bool { return len(segs[i]) > len(segs[j]) })
	// A root-relative URL after a quote, a parenthesis (CSS url()), white
	// space or a comma (srcset) or "=" (meta refresh), whose first segment
	// the site serves.
	e.urlRe = regexp.MustCompile("([\"'`(\\s,=])(/(?:" + strings.Join(segs, "|") + ")(?:[/?#][^\"'`()\\s,<>]*)?)([\"'`)\\s,])")

	e.enqueue("/")
	for _, u := range []string{"/charts", "/colors", "/sera", "/typeset", "/create", "/blocks",
		"/examples/dashboard", "/examples/tasks", "/examples/playground", "/examples/authentication", "/examples/rtl"} {
		e.enqueue(u)
	}
	for _, it := range NavItems {
		e.enqueue(it.Href)
	}
	for u := range s.Docs {
		e.enqueue(u)
		if rest, ok := strings.CutPrefix(u, "/docs/components/radix/"); ok {
			e.enqueue("/docs/components/base/" + rest)
			e.enqueue("/docs/components/aria/" + rest)
		}
	}
	for len(e.queue) > 0 {
		u := e.queue[0]
		e.queue = e.queue[1:]
		if err := e.fetch(u); err != nil {
			return err
		}
	}

	// 404.html: GitHub Pages serves it for unknown paths.
	rec := httptest.NewRecorder()
	s.notFound(rec, httptest.NewRequest("GET", "/404", nil))
	e.text["/404.html"] = rec.Body.Bytes()

	for p, b := range e.text {
		if err := e.write(p, e.rewrite(b)); err != nil {
			return err
		}
	}
	for p := range e.pages {
		if p == "/" {
			continue
		}
		// A page with pages under it is also written as its folder's
		// index, whichever of the two the host serves.
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(p))); err == nil {
			if err := e.write(p+"/index.html", e.rewrite(e.text[p+".html"])); err != nil {
				return err
			}
		}
	}
	// Folders starting with "_" (/_shadcn, /_site, /_next) are not served
	// by GitHub Pages' Jekyll build.
	return os.WriteFile(filepath.Join(dir, ".nojekyll"), nil, 0o644)
}

// routeSegments are the first path segments of the routes of Handler
// (besides the public directory's entries).
var routeSegments = []string{
	"docs", "_shadcn", "_site", "assets", "fonts", "_next", "api", "view",
	"blocks", "charts", "colors", "examples", "sera", "typeset", "typeset.css",
	"create", "preview", "rss.xml",
}

type exporter struct {
	s     *Site
	h     http.Handler
	dir   string
	base  string
	segs  map[string]bool
	urlRe *regexp.Regexp

	queue []string
	seen  map[string]bool
	// mapped is the URL a crawled URL with a query is written at.
	mapped map[string]string
	// text holds the text files by output path, rewritten at the end.
	text map[string][]byte
	// pages are the paths of the HTML pages.
	pages map[string]bool
}

func (e *exporter) enqueue(u string) {
	u, _, _ = strings.Cut(u, "#")
	if !e.seen[u] {
		e.seen[u] = true
		e.queue = append(e.queue, u)
	}
}

// staticQuery reports whether the query of a URL under p is only cache
// busting (the response does not depend on it).
func staticQuery(p string) bool {
	for _, pre := range []string{"/assets/", "/_site/", "/_shadcn/", "/fonts/"} {
		if strings.HasPrefix(p, pre) {
			return true
		}
	}
	return false
}

func (e *exporter) fetch(u string) error {
	p, query, _ := strings.Cut(u, "?")
	if query != "" && staticQuery(p) {
		e.mapped[u] = p
		e.enqueue(p)
		return nil
	}
	if strings.HasPrefix(p, "/api/") {
		return nil
	}
	rec := httptest.NewRecorder()
	e.h.ServeHTTP(rec, httptest.NewRequest("GET", u, nil))
	if rec.Code != http.StatusOK {
		if query == "" {
			e.s.logf("export: %s: %d", u, rec.Code)
		}
		return nil
	}
	ct := rec.Header().Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)
	body := rec.Body.Bytes()
	isHTML := mt == "text/html"
	if query != "" {
		if isHTML {
			// Links to a page with a query go to the page itself.
			e.mapped[u] = p
			e.enqueue(p)
			return nil
		}
		sum := sha256.Sum256([]byte(query))
		ext := path.Ext(p)
		if exts, _ := mime.ExtensionsByType(mt); len(exts) > 0 && ext == "" {
			ext = exts[0]
		}
		mp := strings.TrimSuffix(p, "/") + "/" + hex.EncodeToString(sum[:8]) + ext
		e.mapped[u] = mp
		p = mp
	}
	out := p
	if isHTML {
		e.pages[p] = true
		if p == "/" {
			out = "/index.html"
		} else {
			out = strings.TrimSuffix(p, "/") + ".html"
		}
	}
	if isHTML || mt == "text/css" || mt == "text/javascript" || mt == "application/javascript" {
		e.text[out] = body
		e.discover(string(body), isHTML)
		if !isHTML {
			// Relative imports of the ES modules.
			for _, m := range relImportRe.FindAllStringSubmatch(string(body), -1) {
				e.enqueue(path.Join(path.Dir(p), m[1]))
			}
		}
		return nil
	}
	return e.write(out, body)
}

var relImportRe = regexp.MustCompile(`(?:from|import)\s*\(?\s*["'](\.\.?/[^"']+)["']`)

// discover enqueues the root-relative URLs of a text file.
func (e *exporter) discover(s string, isHTML bool) {
	for _, m := range e.urlRe.FindAllStringSubmatch(s, -1) {
		u := m[2]
		if strings.Contains(u, "${") {
			continue
		}
		if isHTML {
			u = html.UnescapeString(u)
		}
		e.enqueue(u)
	}
}

// rewrite points the root-relative URLs of a text file at their exported
// files, under the base path.
func (e *exporter) rewrite(b []byte) []byte {
	out := e.urlRe.ReplaceAllStringFunc(string(b), func(m string) string {
		sm := e.urlRe.FindStringSubmatch(m)
		u := sm[2]
		if mp, ok := e.mapped[html.UnescapeString(u)]; ok {
			u = mp
		}
		return sm[1] + e.base + u + sm[3]
	})
	if e.base != "" {
		// The home page.
		out = strings.ReplaceAll(out, `href="/"`, `href="`+e.base+`/"`)
	}
	return []byte(out)
}

func (e *exporter) write(p string, b []byte) error {
	f := filepath.Join(e.dir, filepath.FromSlash(strings.TrimPrefix(p, "/")))
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(f, b, 0o644); err != nil {
		return fmt.Errorf("export: %w", err)
	}
	return nil
}
