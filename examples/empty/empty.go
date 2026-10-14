// Package empty ports the upstream empty-* examples.
package empty

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/empty"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/kbd"
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
		ex("empty-avatar", Avatar),
		ex("empty-avatar-group", AvatarGroup),
		ex("empty-input-group", InputGroup),
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

func Avatar() h.Node {
	return empty.Empty(
		empty.Header(
			empty.Media(empty.MediaDefault,
				avatar.Avatar(h.Class("size-12"),
					avatar.Image(h.Src("https://github.com/shadcn.png"), h.Class("grayscale")),
					avatar.Fallback("LR"),
				),
			),
			empty.Title("User Offline"),
			empty.Description("This user is currently offline. You can leave a message to notify them or try again later."),
		),
		empty.Content(button.Button(button.SizeSm, "Leave Message")),
	)
}

func AvatarGroup() h.Node {
	av := func(src, alt, fb string) h.Node {
		return avatar.Avatar(avatar.Image(h.Src(src), h.Alt(alt)), avatar.Fallback(fb))
	}
	return empty.Empty(
		empty.Header(
			empty.Media(
				h.Div(h.Class("flex -space-x-2 *:data-[slot=avatar]:size-12 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:ring-background *:data-[slot=avatar]:grayscale"),
					av("https://github.com/shadcn.png", "@shadcn", "CN"),
					av("https://github.com/maxleiter.png", "@maxleiter", "LR"),
					av("https://github.com/evilrabbit.png", "@evilrabbit", "ER"),
				),
			),
			empty.Title("No Team Members"),
			empty.Description("Invite your team to collaborate on this project."),
		),
		empty.Content(button.Button(button.SizeSm, lucide.Plus(), "Invite Members")),
	)
}

func InputGroup() h.Node {
	return empty.Empty(
		empty.Header(
			empty.Title("404 - Not Found"),
			empty.Description("The page you're looking for doesn't exist. Try searching for what you need below."),
		),
		empty.Content(
			ig.InputGroup(h.Class("sm:w-3/4"),
				ig.Input(h.Placeholder("Try searching for pages...")),
				ig.Addon(lucide.Search()),
				ig.Addon(ig.AlignInlineEnd, kbd.Kbd("/")),
			),
			empty.Description("Need help? ", h.A(h.Href("#"), "Contact support")),
		),
	)
}
