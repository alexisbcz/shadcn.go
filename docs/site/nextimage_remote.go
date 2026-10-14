package site

import (
	"crypto/sha1"
	"encoding/hex"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
)

// serveRemoteNextImage serves /_next/image for a remote url (the docs
// previews' next/image, see mdx/preview_next.go): the reference optimizer's
// output vendored by gen/next-image-remote.mjs under
// static/nextimage/remote/<w>/<sha1(url)>.webp, or a redirect to the
// original image when it is not vendored.
func (s *Site) serveRemoteNextImage(w http.ResponseWriter, r *http.Request, src string, width int) {
	sum := sha1.Sum([]byte(src))
	if b, err := fs.ReadFile(static, "static/nextimage/remote/"+strconv.Itoa(width)+"/"+hex.EncodeToString(sum[:])+".webp"); err == nil {
		w.Header().Set("Vary", "Accept")
		w.Header().Set("Cache-Control", "public, max-age=14400, must-revalidate")
		w.Header().Set("Content-Type", "image/webp")
		_, _ = w.Write(b)
		return
	}
	if u, err := url.Parse(src); err != nil || u.Scheme != "https" || !nextImageRemoteHosts[u.Hostname()] {
		http.Error(w, `"url" parameter is not allowed`, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, src, http.StatusFound)
}

// nextImageRemoteHosts is images.remotePatterns of ui/apps/v4/next.config.mjs.
var nextImageRemoteHosts = map[string]bool{
	"avatars.githubusercontent.com": true,
	"images.unsplash.com":           true,
	"avatar.vercel.sh":              true,
}
