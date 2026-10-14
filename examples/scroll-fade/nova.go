// Package scrollfade ports the docs previews of the scroll-fade utility
// (ui/apps/v4/examples/radix/scroll-fade-*.tsx, shadcn/tailwind.css).
package scrollfade

import (
	"strconv"

	"shadcn.go/components/style"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "scroll-fade", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("scroll-fade-demo", func() h.Node { return frame("h-72 scroll-fade scrollbar-none overflow-y-auto", numbered(12)) }),
		nova("scroll-fade-edge", NovaEdge),
		nova("scroll-fade-horizontal", func() h.Node {
			return frame("scroll-fade-x scrollbar-none overflow-x-auto", tags("rounded-lg px-3", horizontalTags))
		}),
		nova("scroll-fade-none", func() h.Node {
			return column("w-full ",
				labeled("h-48 scroll-fade scrollbar-none overflow-y-auto", "scroll-fade", numbered(8)),
				labeled("h-48 scroll-fade scrollbar-none overflow-y-auto scroll-fade-none", "scroll-fade scroll-fade-none", numbered(8)),
			)
		}),
		nova("scroll-fade-overflow", func() h.Node { return frame("scroll-fade scrollbar-none overflow-y-auto", numbered(3)) }),
		nova("scroll-fade-rtl", func() h.Node {
			return frame("scroll-fade-x scrollbar-none overflow-x-auto", tags("rounded-lg px-3", rtlTags), h.Dir("rtl"))
		}),
		nova("scroll-fade-size", func() h.Node {
			return h.Div(h.Class("mx-auto flex w-full max-w-xs flex-col gap-6"),
				labeled("h-48 scroll-fade scrollbar-none overflow-y-auto scroll-fade-4", "scroll-fade-4", numbered(8)),
				labeled("h-48 scroll-fade scrollbar-none overflow-y-auto scroll-fade-24", "scroll-fade-24", numbered(8)),
			)
		}),
	)
}

var (
	horizontalTags = []string{"Design", "Engineering", "Marketing", "Product", "Research", "Sales", "Support", "Operations", "Finance", "Legal", "People", "Security"}
	rtlTags        = []string{"تصميم", "هندسة", "تسويق", "منتج", "أبحاث", "مبيعات", "دعم", "عمليات", "مالية", "قانوني"}
	edgeItems      = []string{"Inbox triage", "Design review", "API contract", "QA pass", "Launch notes", "Metrics follow-up"}
	edgeTags       = []string{"Design", "Engineering", "Marketing", "Product", "Research", "Sales", "Support", "Operations"}
)

func list(items []string) h.Node {
	kids := []any{h.Class("flex flex-col gap-1.5 p-1.5")}
	for _, it := range items {
		kids = append(kids, h.Div(h.Class("rounded-lg bg-muted px-3 py-2.5 text-sm"), it))
	}
	return h.Div(kids...)
}

func numbered(n int) h.Node {
	items := make([]string, n)
	for i := range items {
		items[i] = "Item " + strconv.Itoa(i+1)
	}
	return list(items)
}

func tags(class string, names []string) h.Node {
	kids := []any{h.Class("flex w-max gap-1.5 p-1.5")}
	for _, t := range names {
		kids = append(kids, h.Div(h.Class("shrink-0 "+class+" bg-muted py-2.5 text-sm"), t))
	}
	return h.Div(kids...)
}

func frame(scroller string, content h.Node, attrs ...any) h.Node {
	return h.Div(append([]any{h.Class("mx-auto w-full max-w-xs overflow-hidden rounded-2xl border")}, append(attrs,
		h.Div(h.Class(scroller), content))...)...)
}

func labeled(scroller, label string, content h.Node) h.Node {
	return h.Div(h.Class("flex flex-col gap-3"),
		h.Div(h.Class("overflow-hidden rounded-2xl border"), h.Div(h.Class(scroller), content)),
		h.P(h.Class("text-center font-mono text-xs text-muted-foreground"), label),
	)
}

func column(width string, kids ...any) h.Node {
	return h.Div(append([]any{h.Class("mx-auto flex " + width + "max-w-xs min-w-0 flex-col gap-6")}, kids...)...)
}

func NovaEdge() h.Node {
	return column("",
		labeled("h-36 scroll-fade-t scrollbar-none overflow-y-auto", "scroll-fade-t", list(edgeItems)),
		labeled("h-36 scroll-fade-b scrollbar-none overflow-y-auto", "scroll-fade-b", list(edgeItems)),
		labeled("scroll-fade-s scrollbar-none overflow-x-auto", "scroll-fade-s", tags("rounded-xl px-4", edgeTags)),
		labeled("scroll-fade-e scrollbar-none overflow-x-auto", "scroll-fade-e", tags("rounded-xl px-4", edgeTags)),
	)
}
