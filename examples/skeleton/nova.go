package skeleton

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/skeleton"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the skeleton page
// (ui/apps/v4/examples/radix/skeleton-*.tsx), in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "skeleton", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("skeleton-avatar", NovaAvatar),
		nova("skeleton-card", NovaCard),
		nova("skeleton-demo", func() h.Node { return novaDemo() }),
		nova("skeleton-form", NovaForm),
		nova("skeleton-rtl", func() h.Node { return style.ScopeRTL(novaDemo(h.Dir("rtl"))) }),
		nova("skeleton-table", NovaTable),
		nova("skeleton-text", NovaText),
	)
}

func sk(class string) h.Node { return skeleton.Skeleton(h.Class(class)) }

func NovaAvatar() h.Node {
	return h.Div(h.Class("flex w-fit items-center gap-4"),
		sk("size-10 shrink-0 rounded-full"),
		h.Div(h.Class("grid gap-2"), sk("h-4 w-[150px]"), sk("h-4 w-[100px]")),
	)
}

func NovaCard() h.Node {
	return card.Card(h.Class("w-full max-w-xs"),
		card.Header(sk("h-4 w-2/3"), sk("h-4 w-1/2")),
		card.Content(sk("aspect-video w-full")),
	)
}

func novaDemo(attrs ...any) h.Node {
	return h.Div(append([]any{h.Class("flex items-center gap-4")}, append(attrs,
		sk("h-12 w-12 rounded-full"),
		h.Div(h.Class("space-y-2"), sk("h-4 w-[250px]"), sk("h-4 w-[200px]")),
	)...)...)
}

func NovaForm() h.Node {
	return h.Div(h.Class("flex w-full max-w-xs flex-col gap-7"),
		h.Div(h.Class("flex flex-col gap-3"), sk("h-4 w-20"), sk("h-8 w-full")),
		h.Div(h.Class("flex flex-col gap-3"), sk("h-4 w-24"), sk("h-8 w-full")),
		sk("h-8 w-24"),
	)
}

func NovaTable() h.Node {
	rows := []any{h.Class("flex w-full max-w-sm flex-col gap-2")}
	for range 5 {
		rows = append(rows, h.Div(h.Class("flex gap-4"), sk("h-4 flex-1"), sk("h-4 w-24"), sk("h-4 w-20")))
	}
	return h.Div(rows...)
}

func NovaText() h.Node {
	return h.Div(h.Class("flex w-full max-w-xs flex-col gap-2"), sk("h-4 w-full"), sk("h-4 w-full"), sk("h-4 w-3/4"))
}
