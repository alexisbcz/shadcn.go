package marker

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the marker page (ui/apps/v4/examples/radix/marker-*.tsx).
// Upstream's examples import the radix-rhea components, so they render in
// style.Rhea.
func init() {
	docs := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "marker", Render: registry.Static(func() h.Node {
			return style.Scope(style.Rhea, fn())
		})}
	}
	registry.Register(
		docs("marker-border", DocsBorder),
		docs("marker-demo", DocsDemo),
		docs("marker-icon", DocsIcon),
		docs("marker-link-button", DocsLinkButton),
		docs("marker-separator", DocsSeparator),
		docs("marker-shimmer", Shimmer),
		docs("marker-status", DocsStatus),
		docs("marker-variants", DocsVariants),
	)
}

func col(gap string, children ...any) h.Node {
	return h.Div(append([]any{h.Class("flex w-full max-w-sm flex-col " + gap + " py-12")}, children...)...)
}

func DocsBorder() h.Node {
	return col("gap-3",
		marker.Marker(marker.VariantBorder, marker.Icon(lucide.GitBranch()), marker.Content("Switched to release-candidate")),
		marker.Marker(marker.VariantBorder, marker.Icon(lucide.Search()), marker.Content("Reviewed 8 related files")),
		marker.Marker(marker.VariantBorder, marker.Icon(lucide.FileText()), marker.Content("Opened implementation notes")),
	)
}

func DocsDemo() h.Node {
	return col("gap-8",
		marker.Marker(marker.Icon(lucide.GitBranch()), marker.Content("Switched to a new branch")),
		marker.Marker(h.Role("status"), marker.Icon(spinner.Spinner()), marker.Content(h.Class("shimmer"), "Thinking...")),
		marker.Marker(marker.VariantSeparator, marker.Content("Conversation compacted")),
		marker.Marker(marker.Icon(lucide.Search()), marker.Content("Explored 4 files")),
	)
}

func DocsIcon() h.Node {
	return col("gap-12",
		marker.Marker(marker.Icon(lucide.GitBranch()), marker.Content("Switched to a new branch")),
		marker.Marker(marker.VariantSeparator, marker.Icon(lucide.Search()), marker.Content("Explored 4 files")),
		marker.Marker(h.Class("flex-col"), marker.Icon(lucide.BookOpenCheck()), marker.Content("Syncing completed")),
	)
}

func DocsLinkButton() h.Node {
	return col("gap-8",
		marker.Marker(lib.AsChild,
			h.A(h.Href("#links-and-buttons"), marker.Icon(lucide.GitBranch()), marker.Content("View the pull request")),
		),
		marker.Marker(lib.AsChild,
			h.Button(h.Type("button"), h.Class("transition-colors hover:text-foreground"),
				marker.Icon(lucide.RotateCcw()),
				marker.Content("Revert this change"),
			),
		),
	)
}

func DocsSeparator() h.Node {
	return col("gap-8",
		marker.Marker(marker.VariantSeparator, marker.Content("Today")),
		marker.Marker(marker.VariantSeparator, marker.Content("Worked for 42s")),
		marker.Marker(marker.VariantSeparator, marker.Content("Conversation compacted")),
	)
}

func DocsStatus() h.Node {
	return col("gap-8",
		marker.Marker(h.Role("status"), marker.Icon(spinner.Spinner()), marker.Content("Compacting conversation")),
		marker.Marker(marker.VariantSeparator, h.Role("status"), marker.Icon(spinner.Spinner()), marker.Content("Running tests")),
	)
}

func DocsVariants() h.Node {
	return col("gap-8",
		marker.Marker(marker.Content("A default marker for inline notes.")),
		marker.Marker(marker.VariantSeparator, marker.Content("A separator marker")),
		marker.Marker(marker.VariantBorder, marker.Content("A border marker for row boundaries.")),
	)
}
