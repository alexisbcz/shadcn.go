package site

import (
	"net/http"
	"regexp"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// topLevelSections is TOP_LEVEL_SECTIONS of components/docs-sidebar.tsx.
var topLevelSections = []NavItem{
	{"/docs", "Introduction"},
	{"/docs/components", "Components"},
	{"/docs/installation", "Installation"},
	{"/docs/theming", "Theming"},
	{"/docs/cli", "CLI"},
	{"/docs/typeset", "Typeset"},
	{"/docs/skills", "Skills"},
	{"/docs/registry", "Registry"},
	{"/docs/changelog", "Changelog"},
}

// excludedSections and excludedPages are EXCLUDED_SECTIONS and
// EXCLUDED_PAGES of components/docs-sidebar.tsx.
var (
	excludedSections = map[string]bool{"installation": true, "dark-mode": true, "changelog": true, "rtl": true}
	excludedPages    = map[string]bool{"/docs": true, "/docs/changelog": true, "/docs/rtl": true, "/docs/new": true}
)

// PagesNew is PAGES_NEW of lib/docs.ts: pages with a blue "new" dot.
var PagesNew = map[string]bool{
	"/docs/changelog":                       true,
	"/docs/changelog/2026-08-questionnaire": true,
	"/docs/registry/dynamic-search":         true,
	"/docs/components/radix/questionnaire":  true,
	"/docs/components/base/questionnaire":   true,
	"/docs/components/aria/questionnaire":   true,
	"/docs/react/questionnaire":             true,
	"/docs/registry/health":                 true,
}

var baseRe = regexp.MustCompile(`/docs/components/(radix|base|aria)/`)

// CurrentBase is getCurrentBase of lib/page-tree.ts: the component base of
// a docs path, "radix" by default (the only base the Go port has).
func CurrentBase(pathname string) string {
	if m := baseRe.FindStringSubmatch(pathname); m != nil {
		return m[1]
	}
	return "radix"
}

// PagesFromFolder is getPagesFromFolder of lib/page-tree.ts: a folder's
// pages; for the components folder, the pages of the current base.
func PagesFromFolder(folder *pagetree.Node, currentBase string) []*pagetree.Node {
	if folder.Path == "components" || folder.Name == "Components" {
		names := map[string]string{"radix": "Radix UI", "base": "Base UI", "aria": "React Aria"}
		for _, c := range folder.Children {
			if c.Kind == pagetree.KindFolder && (pathBase(c.Path) == currentBase || c.Name == names[currentBase]) {
				return directPages(c)
			}
		}
		var all []*pagetree.Node
		var walk func(*pagetree.Node)
		walk = func(f *pagetree.Node) {
			for _, c := range f.Children {
				switch c.Kind {
				case pagetree.KindPage:
					if !strings.HasSuffix(c.URL, "/components") {
						all = append(all, c)
					}
				case pagetree.KindFolder:
					walk(c)
				}
			}
		}
		walk(folder)
		return all
	}
	return directPages(folder)
}

func pathBase(p string) string {
	if i := strings.LastIndex(p, "/"); i >= 0 {
		return p[i+1:]
	}
	return p
}

func directPages(f *pagetree.Node) []*pagetree.Node {
	var out []*pagetree.Node
	for _, c := range f.Children {
		if c.Kind == pagetree.KindPage {
			out = append(out, c)
		}
	}
	return out
}

const sidebarMenuButtonClass = "relative h-[30px] w-fit overflow-visible border border-transparent text-[0.8rem] font-medium after:absolute after:inset-x-0 after:-inset-y-1 after:z-0 after:rounded-md data-[active=true]:border-accent data-[active=true]:bg-accent 3xl:fixed:w-full 3xl:fixed:max-w-48"

// DocsSidebar is components/docs-sidebar.tsx. Scroll restoration is the
// root layout's inline script plus static/docs.js (saving the position).
func (s *Site) DocsSidebar(r *http.Request, current string) h.Node {
	item := func(href, name string, active bool) h.Node {
		var dot h.Node
		if PagesNew[href] {
			dot = h.Span(h.Class("flex size-2 rounded-full bg-blue-500"), h.TitleAttr("New"))
		}
		return sidebar.MenuItem(
			sidebar.MenuButton(lib.AsChild, sidebar.IsActive(active), h.Class(sidebarMenuButtonClass),
				h.A(h.Href(href),
					h.Span(h.Class("absolute inset-0 flex w-(--sidebar-menu-width) bg-transparent")),
					name, dot,
				),
			),
		)
	}
	label := func(name string) h.Node {
		return sidebar.GroupLabel(h.Class("font-medium text-muted-foreground"), name)
	}

	var sections []any
	for _, sec := range topLevelSections {
		active := current == sec.Href
		if sec.Href != "/docs" {
			active = strings.HasPrefix(current, sec.Href)
		}
		sections = append(sections, item(sec.Href, sec.Label, active))
	}
	groups := []any{
		sidebar.Group(h.Class("pt-12"),
			label("Sections"),
			sidebar.GroupContent(sidebar.Menu(sections...)),
		),
	}

	base := CurrentBase(current)
	for _, n := range s.Tree.Children {
		id := n.Path
		if n.Kind == pagetree.KindPage {
			id = strings.TrimSuffix(n.Path, ".mdx")
		}
		if excludedSections[id] {
			continue
		}
		var menu h.Node
		if n.Kind == pagetree.KindFolder {
			var items []any
			for _, p := range PagesFromFolder(n, base) {
				if excludedPages[p.URL] {
					continue
				}
				items = append(items, item(p.URL, p.Name, p.URL == current))
			}
			menu = sidebar.Menu(h.Class("gap-0.5"), h.Group(items...))
		}
		groups = append(groups, sidebar.Group(label(n.Name), sidebar.GroupContent(menu)))
	}

	return sidebar.Sidebar(sidebar.CollapsibleNone,
		h.Class("sticky top-[calc(var(--header-height)+0.6rem)] z-30 hidden h-[calc(100svh-10rem)] overflow-hidden overscroll-none bg-transparent [--sidebar-menu-width:--spacing(56)] lg:flex"),
		h.Div(h.Class("absolute top-12 right-2 bottom-0 hidden h-full w-px bg-[linear-gradient(to_bottom,transparent_0%,var(--border)_10%,var(--border)_90%,transparent_100%)] lg:flex")),
		sidebar.Content(h.Data("docs-sidebar-content", ""),
			h.Class("w-(--sidebar-menu-width) scroll-fade scrollbar-none overflow-x-hidden pl-2.5"),
			h.Group(groups...),
		),
	)
}
