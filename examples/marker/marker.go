// Package marker ports the marker parity fixtures (upstream ships Marker
// without new-york-v4 examples).
package marker

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "marker", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("marker-variants", Variants),
		ex("marker-shimmer", Shimmer),
	)
}

// Variants is the parity fixture marker-variants.
func Variants() h.Node {
	return h.Div(h.Class("grid w-full max-w-3xl grid-cols-2 gap-x-10 gap-y-6 py-8"),
		h.Div(h.Class("flex flex-col gap-6"),
			marker.Marker(marker.Content("A default marker for inline notes.")),
			marker.Marker(marker.VariantDefault, marker.Icon(lucide.GitBranch()), marker.Content("Switched to a new branch")),
			marker.Marker(marker.VariantSeparator, marker.Content("Conversation compacted")),
			marker.Marker(marker.VariantSeparator, marker.Icon(lucide.Search()), marker.Content("Explored 4 files")),
			marker.Marker(marker.VariantBorder, marker.Content("A border marker for row boundaries.")),
			marker.Marker(marker.VariantBorder, marker.Icon(lucide.FileText()), marker.Content("Opened implementation notes")),
			marker.Marker(h.Class("flex-col"), marker.Icon(lucide.BookOpenCheck()), marker.Content("Syncing completed")),
		),
		h.Div(h.Class("flex flex-col gap-6"),
			marker.Marker(h.Role("status"), marker.Icon(spinner.Spinner()), marker.Content("Compacting conversation")),
			marker.Marker(marker.VariantSeparator, h.Role("status"), marker.Icon(spinner.Spinner()), marker.Content("Running tests")),
			marker.Marker(lib.AsChild,
				h.A(h.Href("#links-and-buttons"),
					marker.Icon(lucide.GitBranch()),
					marker.Content("View the pull request"),
				),
			),
			marker.Marker(lib.AsChild, marker.VariantBorder,
				h.Button(h.Type("button"), h.Class("transition-colors hover:text-foreground"),
					marker.Icon(h.Class("text-foreground"), lucide.RotateCcw()),
					marker.Content("Revert this change"),
				),
			),
			marker.Marker(marker.Content("Read the ", h.A(h.Href("#docs"), "documentation"), " for details.")),
			marker.Marker(marker.VariantSeparator, marker.Content(h.Class("font-medium text-foreground"), "Today")),
			marker.Marker(marker.VariantSeparator, h.Class("text-xs"), marker.Content("Worked for 42s")),
		),
	)
}

// Shimmer is the parity fixture marker-shimmer.
func Shimmer() h.Node {
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-8 py-12"),
		marker.Marker(h.Role("status"), marker.Content(h.Class("shimmer"), "Thinking...")),
		marker.Marker(marker.VariantSeparator, h.Role("status"), marker.Content(h.Class("shimmer"), "Reading 4 files")),
	)
}
