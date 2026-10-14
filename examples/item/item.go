// Package item ports the upstream item-* examples.
package item

import (
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/item"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "item", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("item-demo", Demo),
		ex("item-header", Header),
		ex("item-icon", Icon),
		ex("item-image", Image),
		ex("item-link", Link),
		ex("item-size", Size),
		ex("item-variant", Variant),
		ex("item-avatar", Avatar),
		ex("item-group", Group),
	)
}

// nextImage renders the markup of next/image for an image with fixed
// dimensions (as the reference's next/image shim does).
func nextImage(src, alt string, width, height int, class string) *h.Element {
	return h.Img(
		h.Alt(alt),
		h.Loading("lazy"),
		h.Width(strconv.Itoa(width)),
		h.Height(strconv.Itoa(height)),
		h.Attr("decoding", "async"),
		h.Data("nimg", "1"),
		h.Style("color:transparent"),
		h.Src(src),
		h.Class(class),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"),
		item.Item(item.VariantOutline,
			item.Content(
				item.Title("Basic Item"),
				item.Description("A simple item with title and description."),
			),
			item.Actions(button.Button(button.VariantOutline, button.SizeSm, "Action")),
		),
		item.Item(item.VariantOutline, item.SizeSm, lib.AsChild,
			h.A(h.Href("#"),
				item.Media(lucide.BadgeCheck(h.Class("size-5"))),
				item.Content(item.Title("Your profile has been verified.")),
				item.Actions(lucide.ChevronRight(h.Class("size-4"))),
			),
		),
	)
}

func Size() h.Node { return Demo() }

func Header() h.Node {
	models := []struct{ name, description, image string }{
		{"v0-1.5-sm", "Everyday tasks and UI generation.", "https://images.unsplash.com/photo-1650804068570-7fb2e3dbf888?q=80&w=640&auto=format&fit=crop"},
		{"v0-1.5-lg", "Advanced thinking or reasoning.", "https://images.unsplash.com/photo-1610280777472-54133d004c8c?q=80&w=640&auto=format&fit=crop"},
		{"v0-2.0-mini", "Open Source model for everyone.", "https://images.unsplash.com/photo-1602146057681-08560aee8cde?q=80&w=640&auto=format&fit=crop"},
	}
	group := item.Group(h.Class("grid grid-cols-3 gap-4"))
	for _, m := range models {
		group.Children(item.Item(item.VariantOutline,
			item.Header(nextImage(m.image, m.name, 128, 128, "aspect-square w-full rounded-sm object-cover")),
			item.Content(item.Title(m.name), item.Description(m.description)),
		))
	}
	return h.Div(h.Class("flex w-full max-w-xl flex-col gap-6"), group)
}

func Icon() h.Node {
	return h.Div(h.Class("flex w-full max-w-lg flex-col gap-6"),
		item.Item(item.VariantOutline,
			item.Media(item.MediaIcon, lucide.ShieldAlert()),
			item.Content(
				item.Title("Security Alert"),
				item.Description("New login detected from unknown device."),
			),
			item.Actions(button.Button(button.SizeSm, button.VariantOutline, "Review")),
		),
	)
}

func Image() h.Node {
	music := []struct{ title, artist, album, duration string }{
		{"Midnight City Lights", "Neon Dreams", "Electric Nights", "3:45"},
		{"Coffee Shop Conversations", "The Morning Brew", "Urban Stories", "4:05"},
		{"Digital Rain", "Cyber Symphony", "Binary Beats", "3:30"},
	}
	group := item.Group(h.Class("gap-4"))
	for _, s := range music {
		group.Children(item.Item(item.VariantOutline, lib.AsChild, h.Role("listitem"),
			h.A(h.Href("#"),
				item.Media(item.MediaImage,
					nextImage("https://avatar.vercel.sh/"+s.title, s.title, 32, 32, "object-cover grayscale"),
				),
				item.Content(
					// React renders separate text nodes ({title} -{" "}); comments keep
					// them apart when the HTML is parsed.
					item.Title(h.Class("line-clamp-1"), s.title, h.Comment(""), " -", h.Comment(""), " ", h.Span(h.Class("text-muted-foreground"), s.album)),
					item.Description(s.artist),
				),
				item.Content(h.Class("flex-none text-center"), item.Description(s.duration)),
			),
		))
	}
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"), group)
}

func Link() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-4"),
		item.Item(lib.AsChild,
			h.A(h.Href("#"),
				item.Content(
					item.Title("Visit our documentation"),
					item.Description("Learn how to get started with our components."),
				),
				item.Actions(lucide.ChevronRight(h.Class("size-4"))),
			),
		),
		item.Item(item.VariantOutline, lib.AsChild,
			h.A(h.Href("#"), h.Target("_blank"), h.Rel("noopener noreferrer"),
				item.Content(
					item.Title("External resource"),
					item.Description("Opens in a new tab with security attributes."),
				),
				item.Actions(lucide.ExternalLink(h.Class("size-4"))),
			),
		),
	)
}

func Variant() h.Node {
	v := func(variant item.Variant, title, desc string) h.Node {
		return item.Item(variant,
			item.Content(item.Title(title), item.Description(desc)),
			item.Actions(button.Button(button.VariantOutline, button.SizeSm, "Open")),
		)
	}
	return h.Div(h.Class("flex flex-col gap-6"),
		v(item.VariantDefault, "Default Variant", "Standard styling with subtle background and borders."),
		v(item.VariantOutline, "Outline Variant", "Outlined style with clear borders and transparent background."),
		v(item.VariantMuted, "Muted Variant", "Subdued appearance with muted colors for secondary content."),
	)
}

func Avatar() h.Node {
	return h.Div(h.Class("flex w-full max-w-lg flex-col gap-6"),
		item.Item(item.VariantOutline,
			item.Media(
				avatar.Avatar(h.Class("size-10"),
					avatar.Image(h.Src("https://github.com/evilrabbit.png")),
					avatar.Fallback("ER"),
				),
			),
			item.Content(item.Title("Evil Rabbit"), item.Description("Last seen 5 months ago")),
			item.Actions(
				button.Button(button.SizeIconSm, button.VariantOutline, h.Class("rounded-full"), h.Aria("label", "Invite"), lucide.Plus()),
			),
		),
		item.Item(item.VariantOutline,
			item.Media(
				h.Div(h.Class("flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:ring-background *:data-[slot=avatar]:grayscale"),
					avatar.Avatar(h.Class("hidden sm:flex"),
						avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn")),
						avatar.Fallback("CN"),
					),
					avatar.Avatar(h.Class("hidden sm:flex"),
						avatar.Image(h.Src("https://github.com/maxleiter.png"), h.Alt("@maxleiter")),
						avatar.Fallback("LR"),
					),
					avatar.Avatar(
						avatar.Image(h.Src("https://github.com/evilrabbit.png"), h.Alt("@evilrabbit")),
						avatar.Fallback("ER"),
					),
				),
			),
			item.Content(item.Title("No Team Members"), item.Description("Invite your team to collaborate on this project.")),
			item.Actions(button.Button(button.SizeSm, button.VariantOutline, "Invite")),
		),
	)
}

func Group() h.Node {
	people := []struct{ username, avatar, email string }{
		{"shadcn", "https://github.com/shadcn.png", "shadcn@vercel.com"},
		{"maxleiter", "https://github.com/maxleiter.png", "maxleiter@vercel.com"},
		{"evilrabbit", "https://github.com/evilrabbit.png", "evilrabbit@vercel.com"},
	}
	group := item.Group()
	for i, p := range people {
		group.Children(item.Item(
			item.Media(avatar.Avatar(
				avatar.Image(h.Src(p.avatar), h.Class("grayscale")),
				avatar.Fallback(p.username[:1]),
			)),
			item.Content(h.Class("gap-1"), item.Title(p.username), item.Description(p.email)),
			item.Actions(button.Button(button.VariantGhost, button.SizeIcon, h.Class("rounded-full"), lucide.Plus())),
		))
		if i != len(people)-1 {
			group.Children(item.Separator())
		}
	}
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"), group)
}
