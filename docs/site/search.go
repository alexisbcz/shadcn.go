package site

import (
	"encoding/json"
	"net/http"

	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/kbd"
	h "shadcn.go/libhtml"
)

// Extension point: search (components/command-menu.tsx): the trigger in the
// header, the command dialog (pages of the tree, colors, blocks), and its
// index. The trigger is a rough placeholder; the dialog is not built yet.

// searchIndexPath serves the search index (the docs pages) as JSON.
const searchIndexPath = "/_site/search.json"

// CommandMenu is the search trigger of the header. TODO: the dialog (⌘K).
func (s *Site) CommandMenu(r *http.Request) h.Node {
	return button.Button(button.Variant("secondary"),
		h.Class("relative h-8 w-full justify-start bg-surface pl-3 font-medium text-foreground shadow-none sm:pr-12 md:w-48 lg:w-56 xl:w-64 dark:bg-card"),
		h.Data("site-action", "open-search"),
		h.Span(h.Class("hidden lg:inline-flex"), "Search documentation..."),
		h.Span(h.Class("inline-flex lg:hidden"), "Search..."),
		h.Div(h.Class("absolute top-1.5 right-1.5 hidden gap-1 sm:flex"),
			kbd.Group(kbd.Kbd(h.Class("border"), "⌘"), kbd.Kbd(h.Class("border"), "K")),
		),
	)
}

type searchEntry struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

func (s *Site) serveSearchIndex(w http.ResponseWriter, r *http.Request) {
	var out []searchEntry
	for _, p := range s.Tree.Pages() {
		e := searchEntry{Title: p.Name, URL: p.URL}
		if d, ok := s.Docs[p.URL]; ok {
			e.Description = d.Frontmatter.Description
		}
		out = append(out, e)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
