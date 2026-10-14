package site

import (
	"embed"
	"io/fs"
	"net/http"
)

// static holds the vendored fonts (static/fonts: the woff2 files and the
// next/font CSS of the reference build, with URLs rewritten to /fonts/) and
// the site's own scripts (static/*.js).
//
//go:embed static
var static embed.FS

// staticPrefix is the URL prefix of the site scripts.
const staticPrefix = "/_site/"

func fontsHandler() http.Handler {
	sub, _ := fs.Sub(static, "static/fonts")
	return immutable(http.FileServer(http.FS(sub)))
}

func staticHandler() http.Handler {
	sub, _ := fs.Sub(static, "static")
	return http.FileServer(http.FS(sub))
}

func immutable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		next.ServeHTTP(w, r)
	})
}
