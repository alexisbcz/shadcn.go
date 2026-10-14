package mdx

import (
	"net/url"
	"strconv"
	"strings"

	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Image: next/image with the props merged by package content (className,
// src, width, height, alt, ...). Same <img> attributes as next/image, but
// the file is served as is (no /_next/image optimizer).
func Image(_ *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	num := func(name string) string {
		switch v, _ := props.Get(name); v := v.(type) {
		case float64:
			if v != v || v == 0 { // NaN or 0: next/image omits it
				return ""
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		case string:
			return v
		}
		return ""
	}
	img := h.Img(h.Alt(props.String("alt")), h.Attr("loading", "lazy"))
	if w := num("width"); w != "" {
		img.Attr("width", w)
	}
	if hh := num("height"); hh != "" {
		img.Attr("height", hh)
	}
	img.Attr("decoding", "async").Data("nimg", "1")
	if c := className(props); c != "" {
		img.Class(c)
	}
	img.Style("color:transparent")
	src := props.String("src")
	w, _ := strconv.Atoi(num("width"))
	if w <= 0 || !strings.HasPrefix(src, "/") || strings.HasSuffix(src, ".svg") {
		return img.Attr("src", src), nil
	}
	// next/image with a fixed width: the optimizer URL at 1x and 2x
	// (getImgProps, sizes of the reference build), q 75.
	opt := func(x int) string {
		for _, s := range nextImageSizes {
			if s >= x {
				x = s
				break
			}
		}
		if x > nextImageSizes[len(nextImageSizes)-1] {
			x = nextImageSizes[len(nextImageSizes)-1]
		}
		return "/_next/image?url=" + url.QueryEscape(src) + "&w=" + strconv.Itoa(x) + "&q=75"
	}
	return img.Attr("srcset", opt(w)+" 1x, "+opt(2*w)+" 2x").Attr("src", opt(2*w)), nil
}

// nextImageSizes is imageSizes + deviceSizes of the reference build
// (.next/required-server-files.json), sorted.
var nextImageSizes = []int{32, 48, 64, 96, 128, 256, 384, 640, 750, 828, 1080, 1200, 1920, 2048, 3840}

func init() { Register("Image", Image) }
