package mdx

import (
	"strconv"

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
	return img.Style("color:transparent").Attr("src", props.String("src")), nil
}

func init() { Register("Image", Image) }
