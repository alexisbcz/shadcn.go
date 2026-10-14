package site

import (
	"net/http"

	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// Extension point: search (components/command-menu.tsx): the trigger in the
// header, the command dialog (pages of the tree, colors, blocks), and its
// index. The trigger is a rough placeholder; the dialog is not built yet.

// searchIndexPath is the route site.go registers for the search API.
const searchIndexPath = searchAPIPath

// CommandMenu is the search trigger of the header (the DialogTrigger of
// components/command-menu.tsx). TODO: the dialog (⌘K).
func (s *Site) CommandMenu(r *http.Request) h.Node {
	return button.Button(button.VariantOutline,
		h.Class("relative h-8 w-full justify-start rounded-lg border-none bg-muted pl-3 text-foreground shadow-none transition-colors hover:bg-muted/50 md:w-48 lg:w-40 xl:w-64 dark:bg-card"),
		h.Data("site-action", "open-search"),
		h.Aria("haspopup", "dialog"), h.Aria("expanded", "false"), h.Data("state", "closed"),
		h.Span(h.Class("hidden xl:inline-flex"), "Search documentation..."),
		h.Span(h.Class("inline-flex xl:hidden"), "Search..."),
	)
}

func (s *Site) serveSearchIndex(w http.ResponseWriter, r *http.Request) { s.serveSearch(w, r) }
