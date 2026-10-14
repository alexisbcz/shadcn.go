package hovercard

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	hovercard "shadcn.go/components/ui/hover-card"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the hover-card page
// (ui/apps/v4/examples/radix/hover-card-*.tsx), rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "hover-card", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("hover-card-demo", NovaDemo),
		nova("hover-card-sides", NovaSides),
		nova("hover-card-rtl", NovaRtl),
	)
}

var novaSides = []hovercard.Side{hovercard.SideLeft, hovercard.SideTop, hovercard.SideBottom, hovercard.SideRight}

func NovaDemo() h.Node {
	return hovercard.HoverCard(hovercard.OpenDelay(10), hovercard.CloseDelay(100),
		hovercard.Trigger(lib.AsChild, button.Button(button.VariantLink, "Hover Here")),
		hovercard.Content(h.Class("flex w-64 flex-col gap-0.5"),
			h.Div(h.Class("font-semibold"), "@nextjs"),
			h.Div("The React Framework – created and maintained by @vercel."),
			h.Div(h.Class("mt-1 text-xs text-muted-foreground"), "Joined December 2021"),
		),
	)
}

func NovaSides() h.Node {
	kids := []any{h.Class("flex flex-wrap justify-center gap-2")}
	for _, s := range novaSides {
		kids = append(kids, hovercard.HoverCard(hovercard.OpenDelay(100), hovercard.CloseDelay(100),
			hovercard.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("capitalize"), string(s))),
			hovercard.Content(s,
				h.Div(h.Class("flex flex-col gap-1"),
					h.H4(h.Class("font-medium"), "Hover Card"),
					h.P("This hover card appears on the "+string(s)+" side of the trigger."),
				),
			),
		))
	}
	return h.Div(kids...)
}

func NovaRtl() h.Node {
	labels := map[hovercard.Side]string{hovercard.SideLeft: "يسار", hovercard.SideTop: "أعلى", hovercard.SideBottom: "أسفل", hovercard.SideRight: "يمين"}
	kids := []any{h.Class("flex flex-wrap justify-center gap-2")}
	for _, s := range novaSides {
		kids = append(kids, hovercard.HoverCard(hovercard.OpenDelay(10), hovercard.CloseDelay(100),
			hovercard.Trigger(lib.AsChild, button.Button(button.VariantOutline, labels[s])),
			hovercard.Content(s, h.Class("flex w-64 flex-col gap-1"), h.Dir("rtl"),
				h.Div(h.Class("font-semibold"), "سماعات لاسلكية"),
				h.Div(h.Class("text-sm text-muted-foreground"), "٩٩.٩٩ $"),
			),
		))
	}
	return h.Div(kids...)
}
