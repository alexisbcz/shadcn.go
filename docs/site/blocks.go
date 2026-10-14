package site

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/components/ui/resizable"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/components/ui/tabs"
	togglegroup "shadcn.go/components/ui/toggle-group"
	"shadcn.go/docs/content"
	"shadcn.go/docs/site/mdx"
	h "shadcn.go/libhtml"
)

// The blocks gallery: app/(app)/blocks/{layout,page}.tsx,
// app/(app)/blocks/[...categories]/page.tsx, components/block-display.tsx
// and components/block-viewer.tsx. Block data (description, meta, files)
// comes from the registry build output, public/r/styles/new-york-v4.

const blocksTitle = "Building Blocks for the Web"
const blocksDescription = "Clean, modern building blocks. Copy and paste into your apps. Works with all React frameworks. Open Source. Free forever."

// blocksStyle is getActiveStyle(): the first legacy style.
const blocksStyle = "new-york-v4"

var featuredBlocks = []string{"dashboard-01", "sidebar-07", "sidebar-03", "login-03", "login-04"}

// blockCategories is lib/categories.ts.
var blockCategories = []struct {
	name, slug string
	hidden     bool
}{
	{"Sidebar", "sidebar", false},
	{"Dashboard", "dashboard", true},
	{"Authentication", "authentication", true},
	{"Login", "login", false},
	{"Signup", "signup", false},
}

type registryFile struct {
	Path    string `json:"path"`
	Target  string `json:"target"`
	Content string `json:"content"`
}

type registryItem struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Categories  []string       `json:"categories"`
	Files       []registryFile `json:"files"`
	Meta        map[string]any `json:"meta"`
}

var registryItems sync.Map // name -> *registryItem

func (s *Site) registryItem(name string) (*registryItem, bool) {
	if v, ok := registryItems.Load(name); ok {
		return v.(*registryItem), true
	}
	b, err := os.ReadFile(filepath.Join(s.App, "public/r/styles", blocksStyle, name+".json"))
	if err != nil {
		return nil, false
	}
	var it registryItem
	if json.Unmarshal(b, &it) != nil {
		return nil, false
	}
	registryItems.Store(name, &it)
	return &it, true
}

// blockNames is getAllBlockIds(["registry:block"], categories).
func (s *Site) blockNames(categories []string) []string {
	entries, _ := os.ReadDir(filepath.Join(s.App, "public/r/styles", blocksStyle))
	var names []string
	for _, e := range entries {
		name, ok := strings.CutSuffix(e.Name(), ".json")
		if !ok || strings.HasPrefix(name, "chart-") {
			continue
		}
		it, ok := s.registryItem(name)
		if !ok || it.Type != "registry:block" {
			continue
		}
		if len(categories) > 0 && !slices.ContainsFunc(it.Categories, func(c string) bool { return slices.Contains(categories, c) }) {
			continue
		}
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// Blocks serves /blocks (featured blocks).
func (s *Site) Blocks(r *http.Request) (h.Node, error) {
	var items []any
	for _, name := range featuredBlocks {
		items = append(items, s.blockDisplay(r, name))
	}
	items = append(items, h.Div(h.Class("container-wrapper"),
		h.Div(h.Class("container flex justify-center py-6"),
			button.Button(button.VariantOutline, lib.AsChild, h.A(h.Href("/blocks/sidebar"), "Browse more blocks")),
		),
	))
	return s.blocksLayout(r, h.Div(h.Class("flex flex-col gap-12 md:gap-24"), h.Group(items...))), nil
}

// BlocksCategory serves /blocks/{categories...}.
func (s *Site) BlocksCategory(r *http.Request) (h.Node, error) {
	cats := strings.Split(strings.Trim(r.PathValue("categories"), "/"), "/")
	known := false
	for _, c := range blockCategories {
		if c.slug == cats[0] {
			known = true
		}
	}
	if !known {
		return nil, errNotFound
	}
	var items []any
	for _, name := range s.blockNames(cats) {
		items = append(items, s.blockDisplay(r, name))
	}
	title := ""
	for _, c := range blockCategories {
		if c.slug == cats[0] {
			title = c.name + " Blocks"
		}
	}
	return s.blocksLayoutMeta(r, title, "/blocks/"+strings.Join(cats, "/"), h.Div(h.Class("flex flex-col gap-12 md:gap-24"), h.Group(items...))), nil
}

func (s *Site) blocksLayout(r *http.Request, children ...any) h.Node {
	return s.blocksLayoutMeta(r, "", "/blocks", children...)
}

// blocksLayoutMeta renders the layout with the page's metadata: the layout's
// title and description, overridden by a category page's title.
func (s *Site) blocksLayoutMeta(r *http.Request, title, canonical string, children ...any) h.Node {
	m := Meta{Title: blocksTitle, Description: blocksDescription, Canonical: canonical}
	if title != "" {
		m.Title, m.PlainTitle = title, true
	}
	return s.Document(r, m,
		s.AppLayout(r,
			pageHeader(
				announcement(),
				pageHeaderHeading(blocksTitle),
				pageHeaderDescription(blocksDescription),
				// The blocks layout's Button is styles/radix-luma/ui/button.
				pageActions(style.Scope(style.Luma,
					button.Button(lib.AsChild, h.Class("h-[35px]"), h.A(h.Href("#blocks"), "Browse Blocks")),
					button.Button(lib.AsChild, button.VariantSecondary, h.A(h.Href("/docs/components"), "View Components")),
				)),
			),
			pageNav(h.ID("blocks"),
				blocksNav(r.URL.Path),
				style.Scope(style.Luma, button.Button(lib.AsChild, button.VariantSecondary, button.SizeSm, h.Class("mr-7 hidden shadow-none lg:flex"),
					h.A(h.Href("/blocks/sidebar"), "Browse all blocks"))),
			),
			h.Div(h.Class("container-wrapper flex-1 section-soft md:py-12"),
				h.Div(h.Class("container"), h.Group(children...)),
			),
		),
	)
}

// --- page header (components/page-header.tsx, page-nav.tsx, announcement.tsx) ---

func pageHeader(children ...any) h.Node {
	return h.Section(h.Class("border-grid"),
		h.Div(h.Class("container-wrapper"),
			h.Div(h.Class("container flex flex-col items-center gap-2 px-6 py-8 text-center md:py-16 lg:py-20 xl:gap-4"), h.Group(children...)),
		),
	)
}

func pageHeaderHeading(text string) h.Node {
	return h.H1(h.Class("leading-tighter max-w-3xl text-3xl font-semibold tracking-tight text-balance text-primary lg:leading-[1.1] lg:font-semibold xl:text-5xl xl:tracking-tighter"), text)
}

func pageHeaderDescription(text string) h.Node {
	return h.P(h.Class("max-w-4xl text-base text-balance text-foreground sm:text-lg"), text)
}

func pageActions(children ...any) h.Node {
	return h.Div(h.Class("flex w-full items-center justify-center gap-2 pt-2 **:data-[slot=button]:shadow-none"), h.Group(children...))
}

func pageNav(children ...any) h.Node {
	el := h.Div(h.Class("container-wrapper scroll-mt-24"))
	inner := h.Div(h.Class("container flex items-center justify-between gap-4 py-4"))
	for _, c := range children {
		if p, ok := c.(h.Prop); ok {
			el.Children(p)
		} else {
			inner.Children(c)
		}
	}
	return el.Children(inner)
}

func announcement() h.Node {
	return badge.Badge(badge.VariantSecondary, lib.AsChild, h.Class("bg-muted"),
		h.A(h.Href("/docs/changelog"), "New Questionnaire component ", lucide.ArrowRight()))
}

// blocksNav is components/blocks-nav.tsx (the ScrollArea renders plain
// overflow markup on the server).
func blocksNav(pathname string) h.Node {
	link := func(name, slug string) h.Node {
		return h.A(h.Href("/blocks/"+slug),
			h.Class("flex h-7 items-center justify-center px-4 text-center text-base font-medium text-muted-foreground transition-colors hover:text-primary data-[active=true]:text-primary"),
			h.Data("active", boolString(pathname == strings.TrimSuffix("/blocks/"+slug, "/"))),
			name)
	}
	links := []any{link("Featured", "")}
	for _, c := range blockCategories {
		if !c.hidden {
			links = append(links, link(c.name, c.slug))
		}
	}
	return h.Div(h.Class("relative overflow-hidden"),
		h.Div(h.Data("slot", "scroll-area"), h.Class("relative max-w-none"),
			h.Div(h.Data("slot", "scroll-area-viewport"), h.Class("size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1"), h.Style("overflow: scroll"),
				h.Div(h.Style("min-width: 100%; display: table"),
					h.Div(h.Class("flex items-center"), h.Group(links...)),
				),
			),
		),
	)
}

func boolString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// --- block display and viewer --------------------------------------------------------

func (s *Site) blockDisplay(r *http.Request, name string) h.Node {
	it, ok := s.registryItem(name)
	if !ok || len(it.Files) == 0 {
		return nil
	}
	height := "930px"
	if v, ok := it.Meta["iframeHeight"].(string); ok && v != "" {
		height = v
	}
	containerClass, _ := it.Meta["containerClassName"].(string)
	preview, _ := mdx.ComponentPreview(s.Env(r), nil, content.Props{
		{Name: "name", Value: it.Name},
		{Name: "styleName", Value: blocksStyle},
		{Name: "hideCode", Value: true},
		{Name: "className", Value: lib.Cn("my-0 **:[.preview]:h-auto **:[.preview]:p-4 **:[.preview>.p-6]:p-0", containerClass)},
	}, nil)

	active := ""
	if len(it.Files) > 0 {
		active = it.Files[0].Target
	}
	return h.Div(h.ID(it.Name), h.Data("view", "preview"), h.Data("block-viewer", ""),
		h.Class("group/block-view-wrapper flex min-w-0 scroll-mt-24 flex-col-reverse items-stretch gap-4 overflow-hidden md:flex-col"),
		h.Style("--height: "+height),
		blockToolbar(it),
		blockView(it, strings.TrimSuffix(height, "px")),
		blockCode(it, active),
		blockMobile(it, preview),
	)
}

func blockToolbar(it *registryItem) h.Node {
	return h.Div(h.Class("hidden w-full items-center gap-2 pl-2 md:pr-6 lg:flex"),
		tabs.Tabs(tabs.Value("preview"),
			tabs.List(h.Class("grid h-8! grid-cols-2 items-center rounded-lg p-1 *:data-[slot=tabs-trigger]:h-6 *:data-[slot=tabs-trigger]:rounded-sm *:data-[slot=tabs-trigger]:px-2 *:data-[slot=tabs-trigger]:text-xs"),
				tabs.Trigger(tabs.Value("preview"), h.Data("block-view", "preview"), "Preview"),
				tabs.Trigger(tabs.Value("code"), h.Data("block-view", "code"), "Code"),
			),
		),
		separator.Separator(separator.Vertical, h.Class("mx-2 h-4!")),
		h.A(h.Href("#"+it.Name), h.Class("flex-1 text-center text-sm font-medium underline-offset-2 hover:underline md:flex-auto md:text-left"),
			strings.TrimSuffix(it.Description, ".")),
		h.Div(h.Class("ml-auto flex items-center gap-2"),
			h.Div(h.Class("h-8 items-center gap-1.5 rounded-md border p-[3px] shadow-none"),
				togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("100%"), h.Data("block-size", ""),
					h.Class("gap-1 *:data-[slot=toggle-group-item]:size-6! *:data-[slot=toggle-group-item]:rounded-sm!"),
					togglegroup.Item(togglegroup.Value("100%"), h.TitleAttr("Desktop"), lucide.Monitor()),
					togglegroup.Item(togglegroup.Value("60%"), h.TitleAttr("Tablet"), lucide.Tablet()),
					togglegroup.Item(togglegroup.Value("30%"), h.TitleAttr("Mobile"), lucide.Smartphone()),
					separator.Separator(separator.Vertical, h.Class("h-4!")),
					button.Button(button.SizeIcon, button.VariantGhost, lib.AsChild, h.Class("size-6 rounded-sm p-0"), h.TitleAttr("Open in New Tab"),
						h.A(h.Href("/view/"+blocksStyle+"/"+it.Name), h.Target("_blank"),
							h.Span(h.Class("sr-only"), "Open in New Tab"), lucide.Fullscreen())),
					separator.Separator(separator.Vertical, h.Class("h-4!")),
					button.Button(button.SizeIcon, button.VariantGhost, h.Class("size-6 rounded-sm p-0"), h.TitleAttr("Refresh Preview"), h.Data("block-refresh", ""),
						lucide.RotateCw(), h.Span(h.Class("sr-only"), "Refresh Preview")),
				),
			),
			separator.Separator(separator.Vertical, h.Class("mx-1 h-4!")),
			button.Button(button.VariantOutline, button.SizeSm, h.Class("w-fit gap-1 px-2 shadow-none"),
				h.Data("block-copy", "npx shadcn@latest add "+it.Name),
				lucide.Terminal(), h.Span("npx shadcn add "+it.Name)),
			separator.Separator(separator.Vertical, h.Class("mx-1 h-4!")),
			openInV0Button(it.Name),
		),
	)
}

// openInV0Button is components/open-in-v0-button.tsx.
func openInV0Button(name string) h.Node {
	return button.Button(button.SizeSm, lib.AsChild, h.Class("h-[1.8rem] gap-1"),
		// NEXT_PUBLIC_V0_URL and NEXT_PUBLIC_APP_URL are unset in the
		// reference build, which renders "undefined" for both.
		h.A(h.Href("undefined/chat/api/open?url=undefined/r/styles/new-york-v4/"+name+".json"), h.Target("_blank"),
			"Open in ", v0Icon()))
}

func blockView(it *registryItem, height string) h.Node {
	return h.Div(h.Class("hidden group-data-[view=code]/block-view-wrapper:hidden md:h-(--height) lg:flex"),
		h.Div(h.Class("relative grid w-full gap-4"),
			h.Div(h.Class("absolute inset-0 right-4 [background-image:radial-gradient(#d4d4d4_1px,transparent_1px)] [background-size:20px_20px] dark:[background-image:radial-gradient(#404040_1px,transparent_1px)]")),
			resizable.PanelGroup(resizable.Horizontal,
				h.Class("relative z-10 after:absolute after:inset-0 after:right-3 after:z-0 after:rounded-xl after:bg-surface/50"),
				resizable.Panel(resizable.DefaultSize("100%"), resizable.MinSize("30%"), h.Data("block-panel", ""),
					h.Class("relative aspect-[4/2.5] overflow-hidden rounded-lg border bg-background md:aspect-auto md:rounded-xl"),
					h.El("iframe", h.Src("/view/"+blocksStyle+"/"+it.Name), h.Height(height), h.Loading("lazy"),
						h.Class("relative z-20 no-scrollbar w-full bg-background")),
				),
				resizable.Handle(h.Class("relative hidden w-3 bg-transparent p-0 after:absolute after:top-1/2 after:right-0 after:h-8 after:w-[6px] after:translate-x-[-1px] after:-translate-y-1/2 after:rounded-full after:bg-border after:transition-all after:hover:h-10 md:block")),
				resizable.Panel(resizable.DefaultSize("0%"), resizable.MinSize("0%")),
			),
		),
	)
}

func blockCode(it *registryItem, active string) h.Node {
	var figures []any
	for _, f := range it.Files {
		if f.Target == "" {
			continue
		}
		lang := f.Path[strings.LastIndex(f.Path, ".")+1:]
		figures = append(figures, h.El("figure", h.Data("rehype-pretty-code-figure", ""), h.Data("block-file", f.Target),
			h.If(f.Target != active, h.Hidden()),
			h.Class("mx-0! mt-0 flex min-w-0 flex-1 flex-col rounded-xl border-none"),
			h.El("figcaption", h.Data("language", lang),
				h.Class("flex h-12 shrink-0 items-center gap-2 border-b px-4 py-2 text-code-foreground [&_svg]:size-4 [&_svg]:text-code-foreground [&_svg]:opacity-70"),
				content.LanguageIcon(lang), f.Target,
				h.Div(h.Class("ml-auto flex items-center gap-2"),
					button.Button(button.VariantGhost, button.SizeIcon, h.Class("size-7"), h.Data("block-copy-file", f.Target), lucide.Clipboard()),
				),
			),
			h.Div(h.Class("no-scrollbar overflow-y-auto"), mdx.HighlightCode(f.Content, lang)),
		))
	}
	if len(figures) == 0 {
		return nil
	}
	return h.Div(h.Class("mr-[14px] flex overflow-hidden rounded-xl border bg-code text-code-foreground group-data-[view=preview]/block-view-wrapper:hidden md:h-(--height)"),
		h.Div(h.Class("w-72"), blockFileTree(it, active)),
		h.Group(figures...),
	)
}

type fileTree struct {
	name     string
	path     string
	children []*fileTree
}

// createFileTreeForRegistryItemFiles (lib/registry.ts).
func newFileTree(files []registryFile) []*fileTree {
	var root []*fileTree
	for _, f := range files {
		p := f.Target
		if p == "" {
			p = f.Path
		}
		level := &root
		parts := strings.Split(p, "/")
		for i, part := range parts {
			isFile := i == len(parts)-1
			var node *fileTree
			for _, n := range *level {
				if n.name == part {
					node = n
				}
			}
			if node != nil {
				if isFile {
					node.path = p
				} else {
					level = &node.children
				}
				continue
			}
			node = &fileTree{name: part}
			if isFile {
				node.path = p
			}
			*level = append(*level, node)
			if !isFile {
				level = &node.children
			}
		}
	}
	return root
}

const treeButtonClass = "rounded-none pl-(--index) whitespace-nowrap hover:bg-muted-foreground/15 focus:bg-muted-foreground/15 focus-visible:bg-muted-foreground/15 active:bg-muted-foreground/15 data-[active=true]:bg-muted-foreground/15"

func blockFileTree(it *registryItem, active string) h.Node {
	var items []any
	for _, n := range newFileTree(it.Files) {
		items = append(items, treeItem(n, 1, active))
	}
	return sidebar.Provider(h.Class("flex min-h-full! flex-col border-r"),
		sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("w-full flex-1"),
			sidebar.GroupLabel(h.Class("h-12 rounded-none border-b px-4 text-sm"), "Files"),
			sidebar.Group(h.Class("p-0"),
				sidebar.GroupContent(sidebar.Menu(h.Class("translate-x-0 gap-1.5"), h.Group(items...))),
			),
		),
	)
}

func treeItem(n *fileTree, index int, active string) h.Node {
	if n.children == nil {
		indent := float64(index) * 1.3
		if index == 2 {
			indent = float64(index) * 1.2
		}
		return sidebar.MenuItem(
			sidebar.MenuButton(h.Class(treeButtonClass), h.Data("index", itoa(index)), h.Data("block-select-file", n.path),
				h.Style("--index: "+remString(indent)),
				lucide.ChevronRight(h.Class("invisible")), lucide.File(h.Class("h-4 w-4")), n.name,
			).IsActive(n.path == active),
		)
	}
	indent := float64(index) * 1.2
	if index == 1 {
		indent = 1
	}
	var subs []any
	for _, c := range n.children {
		subs = append(subs, treeItem(c, index+1, active))
	}
	return sidebar.MenuItem(
		collapsible.Collapsible(h.Class("group/collapsible [&[data-state=open]>button>svg:first-child]:rotate-90"),
			collapsible.Trigger(lib.AsChild,
				sidebar.MenuButton(h.Class(treeButtonClass), h.Style("--index: "+remString(indent)),
					lucide.ChevronRight(h.Class("transition-transform")), lucide.Folder(), n.name),
			),
			collapsible.Content(sidebar.MenuSub(h.Class("m-0 w-full translate-x-0 border-none p-0"), h.Group(subs...))),
		).Open(true),
	)
}

func blockMobile(it *registryItem, preview h.Node) h.Node {
	var body h.Node
	if m, _ := it.Meta["mobile"].(string); m == "component" {
		body = preview
	} else {
		img := func(theme, class string) h.Node {
			u := "%2Fr%2Fstyles%2Fnew-york%2F" + it.Name + "-" + theme + ".png"
			return h.Img(h.Alt(it.Name), h.Data("block", it.Name), h.Loading("lazy"), h.Width("1440"), h.Height("900"),
				h.Attr("decoding", "async"), h.Data("nimg", "1"), h.Class(class), h.Style("color:transparent"),
				h.UnsafeAttr("srcset", "/_next/image?url="+u+"&w=1920&q=75 1x, /_next/image?url="+u+"&w=3840&q=75 2x"),
				h.UnsafeAttr("src", "/_next/image?url="+u+"&w=3840&q=75"))
		}
		body = h.Div(h.Class("overflow-hidden rounded-xl border"),
			img("light", "object-cover dark:hidden"),
			img("dark", "hidden object-cover dark:block"),
		)
	}
	return h.Div(h.Class("flex flex-col gap-2 lg:hidden"),
		h.Div(h.Class("flex items-center gap-2 px-2"),
			h.Div(h.Class("line-clamp-1 text-sm font-medium"), it.Description),
			h.Div(h.Class("ml-auto shrink-0 font-mono text-xs text-muted-foreground"), it.Name),
		),
		body,
	)
}
