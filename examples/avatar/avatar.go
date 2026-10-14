// Package avatar ports the upstream avatar-* examples.
package avatar

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "avatar", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("avatar-demo", Demo),
		ex("avatar-states", States),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex flex-row flex-wrap items-center gap-12"),
		avatar.Avatar(
			avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn")),
			avatar.Fallback("CN"),
		),
		avatar.Avatar(h.Class("rounded-lg"),
			avatar.Image(h.Src("https://github.com/evilrabbit.png"), h.Alt("@evilrabbit")),
			avatar.Fallback("ER"),
		),
		h.Div(h.Class("flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:ring-background *:data-[slot=avatar]:grayscale"),
			avatar.Avatar(
				avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn")),
				avatar.Fallback("CN"),
			),
			avatar.Avatar(
				avatar.Image(h.Src("https://github.com/maxleiter.png"), h.Alt("@maxleiter")),
				avatar.Fallback("LR"),
			),
			avatar.Avatar(
				avatar.Image(h.Src("https://github.com/evilrabbit.png"), h.Alt("@evilrabbit")),
				avatar.Fallback("ER"),
			),
		),
	)
}

// pixel is a 4x4 orange PNG: an image that loads without the network.
const pixel = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAAEElEQVR4nGP4WSwGRwzEcQCBwhghZfPo6AAAAABJRU5ErkJggg=="

// States ports the avatar-states fixture (parity/fixtures/avatar-states.tsx):
// what the upstream example does not use: the sm and lg sizes, an image
// that loads from a data: URL, an image that fails to decode (fallback stays), an
// image without src, fallbacks with delayMs (shown after the delay unless
// the image loaded first), and AvatarBadge, AvatarGroup and
// AvatarGroupCount.
func States() h.Node {
	return h.Div(h.Class("flex flex-col gap-6"),
		h.Div(h.Class("flex items-center gap-4"),
			avatar.Avatar(avatar.SizeSm, avatar.Image(h.Src(pixel), h.Alt("sm")), avatar.Fallback("SM")),
			avatar.Avatar(avatar.Image(h.Src(pixel), h.Alt("default")), avatar.Fallback("DF")),
			avatar.Avatar(avatar.SizeLg, avatar.Image(h.Src(pixel), h.Alt("lg")), avatar.Fallback("LG")),
		),
		h.Div(h.Class("flex items-center gap-4"),
			avatar.Avatar(avatar.SizeSm, avatar.Image(h.Src("data:image/png;base64,AAAA"), h.Alt("broken")), avatar.Fallback("BR")),
			avatar.Avatar(avatar.Image(h.Alt("no src")), avatar.Fallback("NS")),
			avatar.Avatar(avatar.SizeLg, avatar.Image(h.Src("data:image/png;base64,AAAA"), h.Alt("delayed")), avatar.Fallback(avatar.DelayMs(600), "DL")),
			avatar.Avatar(avatar.Image(h.Src(pixel), h.Alt("fast")), avatar.Fallback(avatar.DelayMs(600), "FA")),
			avatar.Avatar(avatar.Fallback("FB")),
		),
		h.Div(h.Class("flex items-center gap-4"),
			avatar.Avatar(avatar.SizeSm, avatar.Fallback("S"), avatar.Badge()),
			avatar.Avatar(avatar.Fallback("D"), avatar.Badge(lucide.Check())),
			avatar.Avatar(avatar.SizeLg, avatar.Fallback("L"), avatar.Badge(h.Class("bg-green-600"), lucide.Check())),
		),
		avatar.Group(
			avatar.Avatar(avatar.Image(h.Src(pixel), h.Alt("one")), avatar.Fallback("A")),
			avatar.Avatar(avatar.Fallback("B")),
			avatar.GroupCount("+3"),
		),
		avatar.Group(
			avatar.Avatar(avatar.SizeLg, avatar.Fallback("A")),
			avatar.Avatar(avatar.SizeLg, avatar.Fallback("B")),
			avatar.GroupCount(lucide.Plus()),
		),
	)
}
