// Package empty ports the upstream empty-* examples.
//
// Not ported yet (they need components owned by other groups):
// empty-avatar and empty-avatar-group (avatar), empty-input-group
// (input-group).
package empty

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/empty"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "empty", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("empty-background", Background),
		ex("empty-demo", Demo),
		ex("empty-icon", Icon),
		ex("empty-outline", Outline),
	)
}

func Background() h.Node {
	return empty.Empty(h.Class("h-full bg-gradient-to-b from-muted/50 from-30% to-background"),
		empty.Header(
			empty.Media(empty.MediaIcon, tabler.IconBell()),
			empty.Title("No Notifications"),
			empty.Description("You're all caught up. New notifications will appear here."),
		),
		empty.Content(
			button.Button(button.VariantOutline, button.SizeSm, lucide.RefreshCcw(), "Refresh"),
		),
	)
}

func Demo() h.Node {
	return empty.Empty(
		empty.Header(
			empty.Media(empty.MediaIcon, tabler.IconFolderCode()),
			empty.Title("No Projects Yet"),
			empty.Description("You haven't created any projects yet. Get started by creating your first project."),
		),
		empty.Content(
			h.Div(h.Class("flex gap-2"),
				button.Button("Create Project"),
				button.Button(button.VariantOutline, "Import Project"),
			),
		),
		button.Button(button.VariantLink, lib.AsChild, h.Class("text-muted-foreground"), button.SizeSm,
			h.A(h.Href("#"), "Learn More ", lucide.ArrowUpRight()),
		),
	)
}

func Icon() h.Node {
	item := func(icon h.Node, title, desc string) h.Node {
		return empty.Empty(empty.Header(
			empty.Media(empty.MediaIcon, icon),
			empty.Title(title),
			empty.Description(desc),
		))
	}
	return h.Div(h.Class("grid gap-8 md:grid-cols-2"),
		item(tabler.IconInbox(), "No messages", "Your inbox is empty. New messages will appear here."),
		item(tabler.IconStar(), "No favorites", "Items you mark as favorites will appear here."),
		item(tabler.IconHeart(), "No likes yet", "Content you like will be saved here for easy access."),
		item(tabler.IconBookmark(), "No bookmarks", "Save interesting content by bookmarking it."),
	)
}

func Outline() h.Node {
	return empty.Empty(h.Class("border border-dashed"),
		empty.Header(
			empty.Media(empty.MediaIcon, tabler.IconCloud()),
			empty.Title("Cloud Storage Empty"),
			empty.Description("Upload files to your cloud storage to access them anywhere."),
		),
		empty.Content(button.Button(button.VariantOutline, button.SizeSm, "Upload Files")),
	)
}
