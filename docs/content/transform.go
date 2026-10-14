package content

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"shadcn.go/docs/mdx"
)

// The remark stage of upstream's pipeline, in its order: MDX's
// remarkMarkAndUnravel, then fumadocs-mdx's plugins (remarkInclude,
// remarkGfm, remarkHeading, remarkImage, remarkCodeTab, remarkNpm,
// remarkStructure, remarkPostprocess). Plugins that would change nothing
// for the docs content are checks that fail on input they would act on.

// unravel is @mdx-js/mdx's remarkMarkAndUnravel: a paragraph holding only
// JSX elements and expressions (and whitespace) is replaced by them, as
// flow content.
func unravel(n *mdx.Node) {
	for i := 0; i < len(n.Children); i++ {
		c := n.Children[i]
		if c.Type == "paragraph" {
			all, oneOrMore := true, false
			for _, k := range c.Children {
				switch {
				case k.Type == "mdxJsxTextElement" || k.Type == "mdxTextExpression":
					oneOrMore = true
				case k.Type == "text" && strings.Trim(k.Value, " \t\n\f\r") == "":
				default:
					all = false
				}
				if !all {
					break
				}
			}
			if all && oneOrMore {
				var repl []*mdx.Node
				for _, k := range c.Children {
					switch k.Type {
					case "mdxJsxTextElement":
						k.Type = "mdxJsxFlowElement"
					case "mdxTextExpression":
						k.Type = "mdxFlowExpression"
					}
					if k.Type == "text" && k.Value != "" && strings.Trim(k.Value, "\t\r\n ") == "" {
						continue
					}
					repl = append(repl, k)
				}
				n.Children = append(n.Children[:i], append(repl, n.Children[i+1:]...)...)
				i--
				continue
			}
		}
		unravel(c)
	}
}

var customIDRe = regexp.MustCompile(`\s*\[#([\s\S]+?)]\s*$`)

// headingIDs is fumadocs' remarkHeading (generateToc: false): every heading
// gets an id, from a trailing "[#id]" or github-slugger.
func headingIDs(root *mdx.Node) {
	s := &Slugger{}
	var visit func(n *mdx.Node)
	visit = func(n *mdx.Node) {
		if n.Type == "heading" {
			if n.Data == nil {
				n.Data = map[string]any{}
			}
			id := ""
			if len(n.Children) > 0 {
				last := n.Children[len(n.Children)-1]
				if last.Type == "text" {
					if m := customIDRe.FindStringSubmatchIndex(last.Value); m != nil && m[3] > m[2] {
						id = last.Value[m[2]:m[3]]
						last.Value = last.Value[:m[0]]
					}
				}
			}
			if id == "" {
				id = s.Slug(flatten(n))
			}
			n.Data["id"] = id
			return
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(root)
}

// flatten is fumadocs' flattenNode: the concatenated values of a node's
// descendants (text, inline code, expressions).
func flatten(n *mdx.Node) string {
	if n.Children != nil || isParent(n.Type) {
		var b strings.Builder
		for _, c := range n.Children {
			b.WriteString(flatten(c))
		}
		return b.String()
	}
	return n.Value
}

func isParent(t string) bool {
	switch t {
	case "text", "inlineCode", "code", "mdxFlowExpression", "mdxTextExpression", "mdxjsEsm", "image", "thematicBreak", "break", "footnoteReference":
		return false
	}
	return true
}

// ImageSize is the intrinsic size of an image.
type ImageSize struct{ Width, Height int }

// ExternalImageSizes holds the sizes of the remote images the docs use.
// fumadocs' remarkImage downloads each external image at build time to read
// its size; these are the sizes the reference site was built with (its
// <img> tags carry them; docs/oracle serves fixtures of that size, since the
// URLs are not reliably fetchable).
var ExternalImageSizes = map[string]ImageSize{
	"https://v0.app/chat-static/button.svg": {Width: 99, Height: 32},
}

// images is fumadocs' remarkImage for the docs' images: external images get
// their size (width/height attributes).
func images(root *mdx.Node, sizes map[string]ImageSize) error {
	var err error
	var visit func(n *mdx.Node)
	visit = func(n *mdx.Node) {
		if n.Type == "image" {
			src, e := url.PathUnescape(n.URL)
			if e != nil {
				err = e
				return
			}
			if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
				err = fmt.Errorf("content: local image %q: images from files are not supported", n.URL)
				return
			}
			u, e := url.Parse(src)
			if e != nil || u.String() != src {
				err = fmt.Errorf("content: image URL %q needs WHATWG URL normalization, which is not supported", src)
				return
			}
			size, ok := sizes[src]
			if !ok {
				if strings.HasSuffix(n.URL, ".svg") {
					return // remarkImage ignores failures for SVG images
				}
				err = fmt.Errorf("content: unknown size for external image %q (add it to ExternalImageSizes)", src)
				return
			}
			if n.Data == nil {
				n.Data = map[string]any{}
			}
			n.Data["src"] = src
			n.Data["width"] = fmt.Sprint(size.Width)
			n.Data["height"] = fmt.Sprint(size.Height)
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(root)
	return err
}

// checkUnsupported fails on input that fumadocs' remarkInclude,
// remarkCodeTab or remarkNpm would transform: the docs use none of it.
func checkUnsupported(root *mdx.Node) error {
	var err error
	var visit func(n *mdx.Node)
	visit = func(n *mdx.Node) {
		if err != nil {
			return
		}
		switch n.Type {
		case "code":
			if n.Lang == "npm" || n.Lang == "package-install" {
				err = fmt.Errorf("content: %s code blocks (remarkNpm) are not supported", n.Lang)
			}
			if regexp.MustCompile(`(^|\s)tab(-group)?=`).MatchString(n.Meta) {
				err = fmt.Errorf("content: code block tabs (remarkCodeTab) are not supported")
			}
		case "mdxJsxFlowElement", "mdxJsxTextElement":
			if n.Name == "include" {
				err = fmt.Errorf("content: <include> (remarkInclude) is not supported")
			}
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(root)
	return err
}
