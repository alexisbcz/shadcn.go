// Package hovercard ports the upstream hover-card-* examples and the
// hover card parity fixtures.
package hovercard

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/button"
	hovercard "shadcn.go/components/ui/hover-card"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "hover-card", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("hover-card-demo", Demo),
		ex("hover-card-sides", Sides),
	)
}

func Demo() h.Node {
	return hovercard.HoverCard(
		hovercard.Trigger(lib.AsChild, button.Button(button.VariantLink, "@nextjs")),
		hovercard.Content(h.Class("w-80"),
			h.Div(h.Class("flex justify-between gap-4"),
				avatar.Avatar(
					avatar.Image(h.Src("https://github.com/vercel.png")),
					avatar.Fallback("VC"),
				),
				h.Div(h.Class("space-y-1"),
					h.H4(h.Class("text-sm font-semibold"), "@nextjs"),
					h.P(h.Class("text-sm"), "The React Framework – created and maintained by @vercel."),
					h.Div(h.Class("text-xs text-muted-foreground"), "Joined December 2021"),
				),
			),
		),
	)
}

// Sides is the parity fixture hover-card-sides.
func Sides() h.Node {
	type variant struct {
		side  hovercard.Side
		align hovercard.Align
	}
	variants := []variant{
		{hovercard.SideTop, hovercard.AlignStart},
		{hovercard.SideRight, hovercard.AlignCenter},
		{hovercard.SideBottom, hovercard.AlignEnd},
		{hovercard.SideLeft, hovercard.AlignCenter},
	}
	return h.Div(h.Class("flex flex-col items-center gap-8"),
		h.Div(h.Class("flex gap-24"),
			h.Map(variants, func(v variant) h.Node {
				return hovercard.HoverCard(hovercard.OpenDelay(100), hovercard.CloseDelay(100),
					hovercard.Trigger(lib.AsChild, button.Button(button.VariantOutline, string(v.side))),
					hovercard.Content(v.side, v.align, h.Class("w-48"),
						h.P(h.Class("text-sm"), "Card on the "+string(v.side)+"."),
						h.A(h.Href("#more"), h.Class("text-sm underline"), "More"),
						hovercard.Arrow(h.Class("fill-popover")),
					),
				)
			}),
		),
		hovercard.HoverCard(hovercard.OpenDelay(100), hovercard.CloseDelay(100),
			hovercard.Trigger(h.Href("#profile"), h.Class("text-sm underline"), "Plain link"),
			hovercard.Content("Default placement."),
		),
	)
}
