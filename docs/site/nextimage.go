package site

import (
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// nextImagePath is the next/image optimizer route. The Go server does not
// resize: it serves the reference build's optimizer output, vendored by
// gen/next-image.mjs under static/nextimage/<w>/<src>.webp, and falls back
// to the original public file for anything not vendored.
const nextImagePath = "/_next/image"

func (s *Site) serveNextImage(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("url")
	width, err := strconv.Atoi(r.URL.Query().Get("w"))
	if err != nil || !strings.HasPrefix(src, "/") || strings.Contains(src, "..") {
		http.Error(w, `"url" parameter is invalid`, http.StatusBadRequest)
		return
	}
	src = path.Clean(src)
	w.Header().Set("Vary", "Accept")
	w.Header().Set("Cache-Control", "public, max-age=14400, must-revalidate")
	if b, err := fs.ReadFile(static, "static/nextimage/"+strconv.Itoa(width)+src+".webp"); err == nil {
		w.Header().Set("Content-Type", "image/webp")
		_, _ = w.Write(b)
		return
	}
	http.ServeFile(w, r, path.Join(s.App, "public", src))
}
