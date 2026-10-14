package separator

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/separator"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the separator page
// (ui/apps/v4/examples/radix/separator-*.tsx), in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "separator", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("separator-demo", func() h.Node {
			return novaDemo("shadcn/ui", "The Foundation for your Design System",
				"A set of beautifully designed components that you can customize, extend, and build on.")
		}),
		nova("separator-list", NovaList),
		nova("separator-menu", NovaMenu),
		nova("separator-rtl", func() h.Node {
			return style.ScopeRTL(novaDemo("shadcn/ui", "الأساس لنظام التصميم الخاص بك",
				"مجموعة من المكونات المصممة بشكل جميل يمكنك تخصيصها وتوسيعها والبناء عليها.", h.Dir("rtl")))
		}),
		nova("separator-vertical", NovaVertical),
	)
}

func novaDemo(title, subtitle, description string, attrs ...any) h.Node {
	return h.Div(append([]any{h.Class("flex max-w-sm flex-col gap-4 text-sm")}, append(attrs,
		h.Div(h.Class("flex flex-col gap-1.5"),
			h.Div(h.Class("leading-none font-medium"), title),
			h.Div(h.Class("text-muted-foreground"), subtitle),
		),
		separator.Separator(),
		h.Div(description),
	)...)...)
}

func NovaList() h.Node {
	row := func(n string) h.Node {
		return h.Dl(h.Class("flex items-center justify-between"),
			h.Dt("Item "+n), h.Dd(h.Class("text-muted-foreground"), "Value "+n))
	}
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-2 text-sm"),
		row("1"), separator.Separator(), row("2"), separator.Separator(), row("3"))
}

func NovaMenu() h.Node {
	col := func(class, title, sub string) h.Node {
		return h.Div(h.Class(class),
			h.Span(h.Class("font-medium"), title),
			h.Span(h.Class("text-xs text-muted-foreground"), sub))
	}
	return h.Div(h.Class("flex items-center gap-2 text-sm md:gap-4"),
		col("flex flex-col gap-1", "Settings", "Manage preferences"),
		separator.Separator(separator.Vertical),
		col("flex flex-col gap-1", "Account", "Profile & security"),
		separator.Separator(separator.Vertical, h.Class("hidden md:block")),
		col("hidden flex-col gap-1 md:flex", "Help", "Support & docs"),
	)
}

func NovaVertical() h.Node {
	return h.Div(h.Class("flex h-5 items-center gap-4 text-sm"),
		h.Div("Blog"), separator.Separator(separator.Vertical),
		h.Div("Docs"), separator.Separator(separator.Vertical),
		h.Div("Source"))
}
