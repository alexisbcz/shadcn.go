package mdx

import (
	"encoding/json"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/field"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/pagination"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// directoryEntry is one entry of ui/apps/v4/registry/directory.json.
type directoryEntry struct {
	Name        string `json:"name"`
	Homepage    string `json:"homepage"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

var (
	directoryOnce sync.Once
	directory     []directoryEntry
)

// loadDirectory reads the registry directory sorted like
// createRegistryDirectoryView without health data (/r/registries.json is
// not part of the build): by lower-cased name.
func loadDirectory(root string) []directoryEntry {
	directoryOnce.Do(func() {
		b, err := os.ReadFile(filepath.Join(root, "ui/apps/v4/registry/directory.json"))
		if err != nil {
			return
		}
		_ = json.Unmarshal(b, &directory)
		sort.SliceStable(directory, func(i, j int) bool {
			return strings.ToLower(directory[i].Name) < strings.ToLower(directory[j].Name)
		})
	})
	return directory
}

// DirectoryList: ui/apps/v4/components/directory-list.tsx. The page and
// the search query come from the URL (?q=, ?page=) and are rendered on the
// server; docs-mdx.js reloads the page with the new query as the search
// input changes (debounced like nuqs, 250 ms). The "Add" buttons open a
// dialog with the registry add command.
func DirectoryList(env *Env, _ *content.Context, _ content.Props, _ []h.Node) (h.Node, error) {
	all := loadDirectory(env.Root)
	query, page := "", 1
	if env.Request != nil {
		q := env.Request.URL.Query()
		query = q.Get("q")
		if p, err := strconv.Atoi(q.Get("page")); err == nil {
			page = p
		}
	}
	norm := func(s string) string {
		return strings.NewReplacer(" ", "", "@", "").Replace(strings.ToLower(s))
	}
	list := all
	if nq := norm(query); nq != "" {
		list = nil
		for _, r := range all {
			if strings.Contains(norm(r.Name), nq) || strings.Contains(norm(r.Description), nq) {
				list = append(list, r)
			}
		}
	}
	const pageSize = 10
	totalPages := max(1, int(math.Ceil(float64(len(list))/pageSize)))
	page = max(1, min(page, totalPages))
	start, end := (page-1)*pageSize, min(page*pageSize, len(list))

	group := item.Group(h.Class("my-8"))
	for i, r := range list[start:end] {
		group.Children(directoryItem(r))
		if i < end-start-1 {
			group.Children(item.Separator(h.Class("my-1")))
		}
	}

	pathname := "/docs/directory"
	if env.Request != nil {
		pathname = env.Request.URL.Path
	}
	children := []any{directorySearch(query, len(list)), group}
	if totalPages > 1 {
		children = append(children, directoryPagination(pathname, query, page, totalPages))
	}
	return h.Div(h.Data("not-typeset", "true"), h.Class("mt-6"), h.Data("docs-directory", ""),
		style.Scope(style.Nova, children...)), nil
}

func directorySearch(query string, count int) h.Node {
	label := "registries"
	if count == 1 {
		label = "registry"
	}
	disabled := "false"
	if query == "" {
		disabled = "true"
	}
	return field.Field(inputgroup.InputGroup(
		inputgroup.Addon(lucide.Search()),
		inputgroup.Input(h.Class("h-full"), h.Attr("placeholder", "Search"), h.Attr("value", query),
			h.Data("docs-directory-search", "")),
		inputgroup.Addon(h.Span(h.Class("text-muted-foreground tabular-nums sm:text-xs"),
			strconv.Itoa(count), " ", label)).Align(inputgroup.AlignInlineEnd),
		inputgroup.Addon(
			inputgroup.Button(lucide.X()).Size(inputgroup.ButtonSizeIconXs).
				Attr("aria-label", "Clear").Attr("title", "Clear").Data("docs-directory-clear", ""),
		).Align(inputgroup.AlignInlineEnd).Data("disabled", disabled).Class("data-[disabled=true]:hidden"),
	))
}

func homepageURL(homepage string) string {
	u, err := url.Parse(homepage)
	if err != nil {
		return homepage
	}
	// URLSearchParams.set order (url.Values.Encode would sort).
	utm := "utm_source=ui.shadcn.com&utm_medium=referral&utm_campaign=directory"
	if u.RawQuery != "" {
		u.RawQuery += "&" + utm
	} else {
		u.RawQuery = utm
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u.String()
}

func directoryItem(r directoryEntry) h.Node {
	content := item.Content(
		item.Title(h.A(h.Href(homepageURL(r.Homepage)), h.Attr("target", "_blank"),
			h.Attr("rel", "noopener noreferrer external"), h.Class("group flex items-center gap-1"),
			r.Name, " ", iconArrowUpRight(h.Class("size-4 opacity-0 group-hover:opacity-100")))),
	)
	if r.Description != "" {
		content.Children(item.Description(h.Class("text-pretty"), r.Description))
	}
	return item.Item(h.Class("group/item relative gap-6 px-0"),
		// The logos are SVG markup from the repository's directory.json
		// (upstream inserts them with dangerouslySetInnerHTML too).
		item.Media(item.MediaImage, h.Class("grayscale *:[svg]:size-8 *:[svg]:fill-foreground"), h.UnsafeHTML(r.Logo)),
		content,
		item.Actions(h.Class("relative z-10 hidden self-start sm:flex"), directoryAddButton(r.Name)),
		item.Footer(h.Class("justify-start pl-16 sm:hidden [&_button]:min-h-11"),
			button.Button("View ", iconArrowUpRight()).Size(button.SizeSm).Variant(button.VariantOutline).Data("size", ""),
			directoryAddButton(r.Name),
		),
	)
}

// directoryAddButton is DirectoryAddButton with its dialog
// (DirectoryAddProvider's desktop dialog).
func directoryAddButton(name string) h.Node {
	commands := [][2]string{
		{"pnpm", "pnpm dlx shadcn@latest registry add " + name},
		{"npm", "npx shadcn@latest registry add " + name},
		{"yarn", "yarn dlx shadcn@latest registry add " + name},
		{"bun", "bunx --bun shadcn@latest registry add " + name},
	}
	list := tabs.List(h.Class("h-auto *:data-[slot=tabs-trigger]:pt-0"))
	var panels []any
	for _, c := range commands {
		list.Children(tabs.Trigger(c[0]).Value(c[0]))
		panels = append(panels, tabs.Content(
			h.Div(h.Class("bg-surface px-3 py-3 text-surface-foreground"),
				h.Div(h.Class("no-scrollbar overflow-x-auto"),
					h.Code(h.Class("font-mono text-sm whitespace-nowrap"), c[1]))),
		).Value(c[0]).Class("mt-0"))
	}
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild,
			button.Button("Add ", tabler.IconPlus()).Size(button.SizeSm).Variant(button.VariantOutline).Class("relative z-10").
				// Base UI buttons carry no data-size: the item group's
				// has-data-[size=sm]:gap-2.5 must not match them.
				Data("size", "")),
		dialog.Content(h.Class("animate-none! sm:max-w-md"),
			dialog.Header(
				dialog.Title("Add Registry"),
				dialog.Description("Run this command to add ", name, " to your project."),
			),
			tabs.Tabs(
				h.Div(h.Class("flex items-center gap-2 border-b p-1.5"),
					list,
					button.Button(tabler.IconCopy(), h.Span(h.Class("sr-only"), "Copy command")).
						Size(button.SizeIconSm).Variant(button.VariantGhost).Class("ml-auto").
						Data("slot", "copy-button").Data("copied", "false"),
				),
				h.Group(panels...),
			).Value("pnpm").Class("gap-0 overflow-hidden rounded-xl border").Data("docs-config", "packageManager"),
			dialog.Footer(
				dialog.Close(lib.AsChild, button.Button("Done").Size(button.SizeSm).Variant(button.VariantOutline)),
			),
		),
	)
}

func directoryPageHref(pathname, query string, page int) string {
	v := url.Values{}
	if query != "" {
		v.Set("q", query)
	}
	if page > 1 {
		v.Set("page", strconv.Itoa(page))
	}
	if s := v.Encode(); s != "" {
		return pathname + "?" + s
	}
	return pathname
}

// directoryPageNumbers is getPageNumbers (0 is an ellipsis).
func directoryPageNumbers(current, total int) []int {
	if total <= 7 {
		out := make([]int, total)
		for i := range out {
			out[i] = i + 1
		}
		return out
	}
	pages := []int{1}
	if current > 4 {
		pages = append(pages, 0)
	} else if current >= 4 {
		pages = append(pages, 2)
	}
	for i := max(2, current-1); i <= min(total-1, current+1); i++ {
		pages = append(pages, i)
	}
	if current < total-3 {
		pages = append(pages, 0)
	} else if current <= total-3 {
		pages = append(pages, total-1)
	}
	return append(pages, total)
}

func directoryPagination(pathname, query string, page, total int) h.Node {
	edge := func(label string, target int, disabled bool, prev bool) h.Node {
		kids := []any{button.SizeDefault, h.Attr("aria-label", "Go to "+map[bool]string{true: "previous", false: "next"}[prev]+" page")}
		text := h.Span(h.Class("hidden sm:block"), label)
		if prev {
			kids = append(kids, h.Class("pl-1.5!"), tabler.IconChevronLeft(h.Class("cn-rtl-flip size-4")), text)
		} else {
			kids = append(kids, h.Class("pr-1.5!"), text, tabler.IconChevronRight(h.Class("cn-rtl-flip size-4")))
		}
		l := pagination.Link(kids...).Class("max-sm:min-h-11 max-sm:min-w-11")
		if disabled {
			l.Class("pointer-events-none opacity-50").Attr("aria-disabled", "true").Attr("tabindex", "-1")
		} else {
			l.Class("cursor-pointer").Attr("href", directoryPageHref(pathname, query, target))
		}
		return pagination.Item(l)
	}
	content := pagination.Content(edge("Previous", page-1, page <= 1, true))
	for _, p := range directoryPageNumbers(page, total) {
		if p == 0 {
			content.Children(pagination.Item(pagination.Ellipsis()))
			continue
		}
		l := pagination.Link(strconv.Itoa(p)).Attr("href", directoryPageHref(pathname, query, p)).
			Class("cursor-pointer max-sm:min-h-11 max-sm:min-w-11")
		if p == page {
			l.Active(true)
		} else {
			l.Active(false)
		}
		content.Children(pagination.Item(l))
	}
	content.Children(edge("Next", page+1, page >= total, false))
	return pagination.Pagination(content)
}

func init() { Register("DirectoryList", DirectoryList) }
