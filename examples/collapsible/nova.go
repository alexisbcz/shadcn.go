package collapsible

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/collapsible-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "collapsible", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("collapsible-basic", NovaBasic),
		nova("collapsible-demo", NovaDemo),
		nova("collapsible-file-tree", NovaFileTree),
		nova("collapsible-settings", NovaSettings),
		nova("collapsible-rtl", NovaRtl),
	)
}

func NovaBasic() h.Node {
	return card.Card(h.Class("mx-auto w-full max-w-sm"),
		card.Content(
			collapsible.Collapsible(h.Class("rounded-md data-[state=open]:bg-muted"),
				collapsible.Trigger(lib.AsChild,
					button.Button(button.VariantGhost, h.Class("group w-full"),
						"Product details",
						lucide.ChevronDown(h.Class("ml-auto group-data-[state=open]:rotate-180")),
					),
				),
				collapsible.Content(h.Class("flex flex-col items-start gap-2 p-2.5 pt-0 text-sm"),
					h.Div("This panel can be expanded or collapsed to reveal additional content."),
					button.Button(button.SizeXs, "Learn More"),
				),
			),
		),
	)
}

type orderText struct{ orderNumber, status, shipped, shippingAddress, address, items, itemsDescription string }

func novaOrder(t orderText, attrs ...any) h.Node {
	return collapsible.Collapsible(append([]any{h.Class("flex w-[350px] flex-col gap-2")}, append(attrs,
		h.Div(h.Class("flex items-center justify-between gap-4 px-4"),
			h.H4(h.Class("text-sm font-semibold"), t.orderNumber),
			collapsible.Trigger(lib.AsChild,
				button.Button(button.VariantGhost, button.SizeIcon, h.Class("size-8"),
					lucide.ChevronsUpDown(),
					h.Span(h.Class("sr-only"), "Toggle details"),
				),
			),
		),
		h.Div(h.Class("flex items-center justify-between rounded-md border px-4 py-2 text-sm"),
			h.Span(h.Class("text-muted-foreground"), t.status),
			h.Span(h.Class("font-medium"), t.shipped),
		),
		collapsible.Content(h.Class("flex flex-col gap-2"),
			h.Div(h.Class("rounded-md border px-4 py-2 text-sm"),
				h.P(h.Class("font-medium"), t.shippingAddress),
				h.P(h.Class("text-muted-foreground"), t.address),
			),
			h.Div(h.Class("rounded-md border px-4 py-2 text-sm"),
				h.P(h.Class("font-medium"), t.items),
				h.P(h.Class("text-muted-foreground"), t.itemsDescription),
			),
		),
	)...)...)
}

func NovaDemo() h.Node {
	return novaOrder(orderText{"Order #4189", "Status", "Shipped", "Shipping address", "100 Market St, San Francisco", "Items", "2x Studio Headphones"})
}

func NovaRtl() h.Node {
	return novaOrder(orderText{"الطلب #4189", "الحالة", "تم الشحن", "عنوان الشحن", "100 Market St, San Francisco", "العناصر", "2x سماعات الاستوديو"}, h.Dir("rtl"))
}

type fileItem struct {
	name  string
	items []fileItem
}

func dir(name string, items ...fileItem) fileItem { return fileItem{name: name, items: items} }

func file(name string) fileItem { return fileItem{name: name} }

func NovaFileTree() h.Node {
	tree := []fileItem{
		dir("components",
			dir("ui", file("button.tsx"), file("card.tsx"), file("dialog.tsx"), file("input.tsx"), file("select.tsx"), file("table.tsx")),
			file("login-form.tsx"), file("register-form.tsx"),
		),
		dir("lib", file("utils.ts"), file("cn.ts"), file("api.ts")),
		dir("hooks", file("use-media-query.ts"), file("use-debounce.ts"), file("use-local-storage.ts")),
		dir("types", file("index.d.ts"), file("api.d.ts")),
		dir("public", file("favicon.ico"), file("logo.svg"), file("images")),
		file("app.tsx"), file("layout.tsx"), file("globals.css"), file("package.json"), file("tsconfig.json"), file("README.md"), file(".gitignore"),
	}
	var render func(it fileItem) h.Node
	render = func(it fileItem) h.Node {
		if it.items != nil {
			var kids []any
			for _, c := range it.items {
				kids = append(kids, render(c))
			}
			return collapsible.Collapsible(
				collapsible.Trigger(lib.AsChild,
					button.Button(button.VariantGhost, button.SizeSm, h.Class("group w-full justify-start transition-none hover:bg-accent hover:text-accent-foreground"),
						lucide.ChevronRight(h.Class("transition-transform group-data-[state=open]:rotate-90")),
						lucide.Folder(),
						it.name,
					),
				),
				collapsible.Content(h.Class("mt-1 ml-5 style-lyra:ml-4"),
					h.Div(append([]any{h.Class("flex flex-col gap-1")}, kids...)...),
				),
			)
		}
		return button.Button(button.VariantLink, button.SizeSm, h.Class("w-full justify-start gap-2 text-foreground"),
			lucide.File(),
			h.Span(it.name),
		)
	}
	var items []any
	for _, it := range tree {
		items = append(items, render(it))
	}
	return card.Card(h.Class("mx-auto w-full max-w-[16rem] gap-2"), card.SizeSm,
		card.Header(
			tabs.Tabs(tabs.Value("explorer"),
				tabs.List(h.Class("w-full"),
					tabs.Trigger(tabs.Value("explorer"), "Explorer"),
					tabs.Trigger(tabs.Value("settings"), "Outline"),
				),
			),
		),
		card.Content(h.Div(append([]any{h.Class("flex flex-col gap-1")}, items...)...)),
	)
}

func NovaSettings() h.Node {
	radius := func(label, labelFor string) h.Node {
		return field.Field(
			field.Label(h.For(labelFor), h.Class("sr-only"), label),
			input.Input(h.ID("radius"), h.Placeholder("0"), h.Value("0")),
		)
	}
	return card.Card(h.Class("mx-auto w-full max-w-xs"), card.SizeSm,
		card.Header(
			card.Title("Radius"),
			card.Description("Set the corner radius of the element."),
		),
		card.Content(
			collapsible.Collapsible(h.Class("flex items-start gap-2"),
				field.Group(h.Class("grid w-full grid-cols-2 gap-2"),
					radius("Radius X", "radius-x"),
					radius("Radius Y", "radius-y"),
					collapsible.Content(h.Class("col-span-full grid grid-cols-subgrid gap-2"),
						radius("Radius X", "radius-x"),
						radius("Radius Y", "radius-y"),
					),
				),
				collapsible.Trigger(lib.AsChild,
					// {isOpen ? <MinimizeIcon /> : <MaximizeIcon />}: both
					// icons, shown by the trigger's data-state.
					button.Button(button.VariantOutline, button.SizeIcon,
						lucide.Maximize(h.Class("[[data-state=open]>&]:hidden")),
						lucide.Minimize(h.Class("[[data-state=closed]>&]:hidden")),
					),
				),
			),
		),
	)
}
