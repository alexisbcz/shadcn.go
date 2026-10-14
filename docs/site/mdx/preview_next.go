package mdx

import (
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	h "shadcn.go/libhtml"
)

// The Go examples render next/image and next/font the way the component
// harness's shims do (a plain <img>, an empty font class). On the docs site
// the reference runs the real next/image and next/font, so a docs preview
// whose upstream source imports them is rewritten here: <img data-nimg>
// gets the optimizer URLs (/_next/image, served by the docs server) and the
// fill styles, and the wrapper of an example using a Google font gets that
// font's class.

// previewFontClass maps an example to the next/font class its wrapper
// carries in the reference build (docs/site/static/fonts/fonts.css).
var previewFontClass = map[string]string{
	"calendar-hijri": "vazirmatn_33215607-module__wVaeya__className",
}

var nextUsage sync.Map // example name -> bool

// usesNext reports whether the upstream radix example imports next/image or
// next/font.
func usesNext(root, name string) bool {
	if v, ok := nextUsage.Load(name); ok {
		return v.(bool)
	}
	src, _ := os.ReadFile(filepath.Join(root, "ui/apps/v4/examples/radix", name+".tsx"))
	uses := strings.Contains(string(src), `"next/image"`) || strings.Contains(string(src), `"next/font/google"`)
	nextUsage.Store(name, uses)
	return uses
}

// nextPreview wraps the preview of example name when its upstream source
// uses next/image or next/font.
func nextPreview(env *Env, name string, preview h.Node) h.Node {
	name = strings.TrimPrefix(name, "nova--")
	if preview == nil || env == nil || !usesNext(env.Root, name) {
		return preview
	}
	s, err := h.RenderString(preview)
	if err != nil {
		return preview
	}
	if font := previewFontClass[name]; font != "" {
		s = strings.Replace(s, `data-example-empty-class=""`, `class="`+font+`"`, 1)
	}
	s = imgTag.ReplaceAllStringFunc(s, nextImg)
	return h.UnsafeHTML(s)
}

var (
	imgTag  = regexp.MustCompile(`<img\b[^>]*\bdata-nimg="[^"]*"[^>]*>`)
	imgAttr = regexp.MustCompile(`([^\s=<>/]+)(?:="([^"]*)")?`)
)

// The images config of the reference build (next.config.mjs defaults).
var (
	nextDeviceSizes = []int{640, 750, 828, 1080, 1200, 1920, 2048, 3840}
	nextAllSizes    = []int{16, 32, 48, 64, 96, 128, 256, 384, 640, 750, 828, 1080, 1200, 1920, 2048, 3840}
)

// encodeURIComponent is JavaScript's encodeURIComponent.
func encodeURIComponent(s string) string {
	var b strings.Builder
	for _, c := range []byte(s) {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.IndexByte("-_.!~*'()", c) >= 0 {
			b.WriteByte(c)
			continue
		}
		b.WriteString("%" + strings.ToUpper(strconv.FormatUint(uint64(c)|0x100, 16)[1:]))
	}
	return b.String()
}

func nextImageURL(src string, w int) string {
	return "/_next/image?url=" + encodeURIComponent(src) + "&w=" + strconv.Itoa(w) + "&q=75"
}

// nextImg rewrites one <img data-nimg> as next/image's getImgProps renders
// it (default loader, quality 75, no sizes prop).
func nextImg(tag string) string {
	type kv struct{ k, v string }
	var attrs []kv
	get := func(k string) string {
		for _, a := range attrs {
			if a.k == k {
				return a.v
			}
		}
		return ""
	}
	set := func(k, v string) {
		for i := range attrs {
			if attrs[i].k == k {
				attrs[i].v = v
				return
			}
		}
		attrs = append(attrs, kv{k, v})
	}
	for _, m := range imgAttr.FindAllStringSubmatch(strings.TrimSuffix(strings.TrimPrefix(tag, "<img"), ">"), -1) {
		attrs = append(attrs, kv{m[1], html.UnescapeString(m[2])})
	}
	src := get("src")
	if src == "" || strings.HasPrefix(src, "/_next/") || strings.HasSuffix(src, ".svg") || strings.HasPrefix(src, "data:") {
		return tag
	}
	if get("data-nimg") == "fill" {
		var set1 []string
		for _, w := range nextDeviceSizes {
			set1 = append(set1, nextImageURL(src, w)+" "+strconv.Itoa(w)+"w")
		}
		set("sizes", "100vw")
		set("srcset", strings.Join(set1, ", "))
		set("src", nextImageURL(src, nextDeviceSizes[len(nextDeviceSizes)-1]))
		set("style", "position:absolute;height:100%;width:100%;left:0;top:0;right:0;bottom:0;color:transparent")
	} else {
		w, err := strconv.Atoi(get("width"))
		if err != nil || w <= 0 {
			return tag
		}
		pick := func(x int) int {
			for _, s := range nextAllSizes {
				if s >= x {
					return s
				}
			}
			return nextAllSizes[len(nextAllSizes)-1]
		}
		a, b := pick(w), pick(2*w)
		srcset := nextImageURL(src, a) + " 1x"
		if b != a {
			srcset += ", " + nextImageURL(src, b) + " 2x"
		}
		set("srcset", srcset)
		set("src", nextImageURL(src, b))
	}
	var out strings.Builder
	out.WriteString("<img")
	for _, a := range attrs {
		out.WriteString(" " + a.k + `="` + html.EscapeString(a.v) + `"`)
	}
	out.WriteString(">")
	return out.String()
}
