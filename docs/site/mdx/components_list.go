package mdx

import (
	"regexp"
	"strings"

	"shadcn.go/docs/content"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// pagesNew is PAGES_NEW of ui/apps/v4/lib/docs.ts.
var pagesNew = map[string]bool{
	"/docs/changelog":                       true,
	"/docs/changelog/2026-08-questionnaire": true,
	"/docs/registry/dynamic-search":         true,
	"/docs/components/radix/questionnaire":  true,
	"/docs/components/base/questionnaire":   true,
	"/docs/components/aria/questionnaire":   true,
	"/docs/react/questionnaire":             true,
	"/docs/registry/health":                 true,
}

// ComponentsList: ui/apps/v4/components/components-list.tsx. Props (from
// package content): componentsFolder (the "components" folder), currentBase,
// variant ("all" or "new").
func ComponentsList(env *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	if env.Tree == nil {
		return h.Group(), nil
	}
	base := props.String("currentBase")
	if base == "" {
		base = "base"
	}
	variant := props.String("variant")
	if variant == "" {
		variant = "all"
	}
	var pages []*pagetree.Node
	if folder := env.Tree.Find("components/" + base); folder != nil {
		for _, c := range folder.Children {
			if c.Kind == pagetree.KindPage {
				pages = append(pages, c)
			}
		}
	}
	grid := h.Div(h.Data("not-typeset", "true"),
		h.Class("mt-8 grid grid-cols-2 gap-4 md:grid-cols-3 md:gap-x-8 lg:gap-x-16 lg:gap-y-6 xl:gap-x-20"))
	n := 0
	for _, p := range pages {
		isNew := pagesNew[p.URL]
		if variant != "all" && !isNew {
			continue
		}
		n++
		link := h.A(h.Class("inline-flex items-center gap-2 text-lg font-medium underline-offset-4 hover:underline md:text-base"),
			p.Name)
		if variant == "all" && isNew {
			link.Children(h.Span(h.Class("sr-only"), "New"),
				h.Span(h.Aria("hidden", "true"), h.Class("flex size-2 rounded-full bg-blue-500")))
		}
		grid.Children(link.Href(p.URL))
	}
	if n == 0 {
		return h.Group(), nil
	}
	return grid, nil
}

func init() { Register("ComponentsList", ComponentsList) }

var (
	reComponentsListNew = regexp.MustCompile(`<ComponentsList\s+variant=["']new["']\s*/>`)
	reComponentsListAll = regexp.MustCompile(`<ComponentsList\s*/>`)
)

// ReplaceComponentsList is replaceComponentsList of ui/apps/v4/lib/llm.ts:
// the page source handed to "Copy Page" lists the components as Markdown
// links in place of <ComponentsList />. description returns a docs page's
// description by URL. Links are absoluteUrl(...), which prefixes
// NEXT_PUBLIC_APP_URL; the reference build leaves it unset ("undefined").
func ReplaceComponentsList(tree *pagetree.Node, src string, description func(url string) string) string {
	list := func(variant string) string {
		if tree == nil {
			return ""
		}
		folder := tree.Find("components/base")
		if folder == nil {
			return ""
		}
		var lines []string
		for _, c := range folder.Children {
			if c.Kind != pagetree.KindPage || (variant == "new" && !pagesNew[c.URL]) {
				continue
			}
			line := "- [" + c.Name + "](undefined" + strings.Replace(c.URL, "/base/", "/", 1) + ")"
			if d := strings.TrimSpace(description(c.URL)); d != "" {
				line += ": " + d
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n")
	}
	src = reComponentsListNew.ReplaceAllLiteralString(src, list("new"))
	return reComponentsListAll.ReplaceAllLiteralString(src, list("all"))
}
