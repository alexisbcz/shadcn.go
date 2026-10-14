// Package content turns the shadcn/ui docs' MDX (ui/apps/v4/content/docs)
// into pages, exactly as upstream's fumadocs-mdx pipeline does: frontmatter,
// heading ids, table of contents, the HTML tree (with code highlighted by
// package highlight), and rendering to libhtml nodes with the element
// mapping of ui/apps/v4/mdx-components.tsx. Custom JSX components (Callout,
// ComponentPreview, ...) are resolved through a Registry.
package content

import (
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"shadcn.go/docs/hast"
	"shadcn.go/docs/highlight"
	"shadcn.go/docs/mdx"
)

// Page is one docs page.
type Page struct {
	// Path is the file path relative to the content root, with slashes
	// ("components/radix/button.mdx").
	Path string
	// Slugs and URL are fumadocs' (baseUrl /docs): folder groups such as
	// "(root)" are dropped, index pages take their folder's URL.
	Slugs []string
	URL   string

	Frontmatter *Frontmatter
	// Tree is the page's HTML tree after the rehype plugins (rehype-toc,
	// rehype-pretty-code): what MDX compiles to JSX.
	Tree []hast.Node
	// TOC is the table of contents (fumadocs' rehype-toc).
	TOC []TOCItem
	// Imports are the page's ESM imports (icons used in JSX).
	Imports []hast.Import
}

// TOCItem is one table of contents entry. Title holds the heading's
// children (rendered without the docs' element mapping).
type TOCItem struct {
	Depth int
	URL   string
	Title []hast.Node
}

// Options configures parsing.
type Options struct {
	// Highlighter highlights code blocks; nil uses highlight.Default().
	Highlighter *highlight.Highlighter
	// ImageSizes gives the sizes of external images; nil uses
	// ExternalImageSizes.
	ImageSizes map[string]ImageSize
}

// ParsePage parses one MDX file. relPath is its path relative to the
// content root (it determines the URL).
func ParsePage(relPath string, src string, opts Options) (*Page, error) {
	matter, body := SplitFrontmatter(src)
	fm, err := ParseFrontmatter(matter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", relPath, err)
	}
	root, err := mdx.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", relPath, err)
	}
	tree, toc, err := Transform(root, fm, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", relPath, err)
	}
	p := &Page{Path: relPath, Frontmatter: fm, Tree: tree, TOC: toc}
	for _, n := range tree {
		if e, ok := n.(*hast.MdxEsm); ok {
			p.Imports = append(p.Imports, e.Imports...)
		}
	}
	p.Slugs, err = Slugs(relPath)
	if err != nil {
		return nil, err
	}
	p.URL = "/docs"
	if len(p.Slugs) > 0 {
		p.URL += "/" + strings.Join(p.Slugs, "/")
	}
	return p, nil
}

// Transform runs the remark and rehype stages on a parsed tree.
func Transform(root *mdx.Node, fm *Frontmatter, opts Options) ([]hast.Node, []TOCItem, error) {
	if err := checkUnsupported(root); err != nil {
		return nil, nil, err
	}
	unravel(root)
	headingIDs(root)
	sizes := opts.ImageSizes
	if sizes == nil {
		sizes = ExternalImageSizes
	}
	if err := images(root, sizes); err != nil {
		return nil, nil, err
	}
	if fm != nil && fm.Title == "" {
		// remarkPostprocess: the first h1 becomes the title.
		for _, c := range root.Children {
			if c.Type == "heading" && c.Depth == 1 {
				fm.Title = flatten(c)
				break
			}
		}
	}
	tree, err := toHast(root)
	if err != nil {
		return nil, nil, err
	}
	toc := rehypeToc(&tree)
	h := opts.Highlighter
	if h == nil {
		if h, err = highlight.Default(); err != nil {
			return nil, nil, err
		}
	}
	if err := prettyCode(tree, h); err != nil {
		return nil, nil, err
	}
	return tree, toc, nil
}

var headingTag = regexp.MustCompile(`^h[1-6]$`)

// rehypeToc is fumadocs' rehype-toc: headings with an id become TOC
// entries; "[!toc]" excludes a heading, "[toc]" keeps it only in the TOC.
func rehypeToc(tree *[]hast.Node) []TOCItem {
	var items []TOCItem
	var visit func(list *[]hast.Node)
	visit = func(list *[]hast.Node) {
		for i := 0; i < len(*list); i++ {
			switch e := (*list)[i].(type) {
			case *hast.Element:
				if headingTag.MatchString(e.Tag) && len(e.Children) > 0 {
					id, ok := e.Get("id")
					if !ok {
						continue
					}
					tocOnly := false
					if last, ok := e.Children[len(e.Children)-1].(*hast.Text); ok {
						if v, ok := handleTag(last.Value, "[!toc]"); ok {
							last.Value = v
							continue
						}
						if v, ok := handleTag(last.Value, "[toc]"); ok {
							tocOnly = true
							last.Value = v
						}
					}
					var title []hast.Node
					for _, c := range e.Children {
						title = append(title, hast.CloneNode(c))
					}
					items = append(items, TOCItem{Depth: int(e.Tag[1] - '0'), URL: "#" + id, Title: title})
					if tocOnly {
						*list = append((*list)[:i], (*list)[i+1:]...)
						i--
					}
					continue
				}
				visit(&e.Children)
			case *hast.MdxJsx:
				visit(&e.Children)
			}
		}
	}
	visit(tree)
	return items
}

func handleTag(value, tag string) (string, bool) {
	i := strings.Index(value, tag)
	if i < 0 {
		return "", false
	}
	return strings.TrimRight(value[:i], " \t\n\r\f\v\u00a0\ufeff") + value[i+len(tag):], true
}

// prettyCode is rehype-pretty-code over the whole tree: every <pre><code>
// is highlighted (inline code is left alone: the docs never use its
// {:lang} syntax, which is rejected).
func prettyCode(nodes []hast.Node, h *highlight.Highlighter) error {
	for _, n := range nodes {
		switch e := n.(type) {
		case *hast.Element:
			if highlight.IsBlockCode(e) {
				if err := h.Apply(e); err != nil {
					return err
				}
				continue
			}
			if e.Tag == "code" && len(e.Children) == 1 {
				if t, ok := e.Children[0].(*hast.Text); ok && inlineLangRe.MatchString(t.Value) {
					return fmt.Errorf("content: inline code language syntax {:lang} is not supported: %q", t.Value)
				}
			}
			if err := prettyCode(e.Children, h); err != nil {
				return err
			}
		case *hast.MdxJsx:
			if err := prettyCode(e.Children, h); err != nil {
				return err
			}
		}
	}
	return nil
}

var inlineLangRe = regexp.MustCompile(`\{:[a-zA-Z.-]+\}$`)

var groupRe = regexp.MustCompile(`^\(.+\)$`)

// Slugs is fumadocs' getSlugs for a page path.
func Slugs(relPath string) ([]string, error) {
	dir, file := path.Split(relPath)
	name := strings.TrimSuffix(file, path.Ext(file))
	var slugs []string
	for _, seg := range strings.Split(strings.Trim(dir, "/"), "/") {
		if seg != "" && !groupRe.MatchString(seg) {
			slugs = append(slugs, encodeURI(seg))
		}
	}
	if groupRe.MatchString(name) {
		return nil, fmt.Errorf("content: cannot use folder group in file names: %s", relPath)
	}
	if name != "index" {
		slugs = append(slugs, encodeURI(name))
	}
	return slugs, nil
}

func encodeURI(s string) string {
	u := url.URL{Path: s}
	return u.EscapedPath()
}

// Load parses every .mdx file under root in fsys, sorted by path.
func Load(fsys fs.FS, root string, opts Options) ([]*Page, error) {
	var paths []string
	err := fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".mdx") {
			paths = append(paths, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	// Pages are independent: parse them in parallel (the highlighter runs
	// several JavaScript runtimes).
	pages := make([]*Page, len(paths))
	errs := make([]error, len(paths))
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.GOMAXPROCS(0))
	for i, p := range paths {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			b, err := fs.ReadFile(fsys, p)
			if err != nil {
				errs[i] = err
				return
			}
			rel := strings.TrimPrefix(strings.TrimPrefix(p, root), "/")
			pages[i], errs[i] = ParsePage(rel, string(b), opts)
		}(i, p)
	}
	wg.Wait()
	taken := map[string]bool{}
	for i, pg := range pages {
		if errs[i] != nil {
			return nil, errs[i]
		}
		key := strings.Join(pg.Slugs, "/")
		if taken[key] {
			return nil, fmt.Errorf("content: duplicated slugs %q", key)
		}
		taken[key] = true
	}
	return pages, nil
}
