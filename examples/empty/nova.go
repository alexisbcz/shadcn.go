package empty

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/empty"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the empty page (ui/apps/v4/examples/radix/empty-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "empty", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("empty-background", NovaBackground),
		nova("empty-card", NovaCard),
		nova("empty-demo", NovaDemo),
		nova("empty-rtl", NovaRtl),
	)
}

func NovaBackground() h.Node {
	return empty.Empty(h.Class("h-full bg-muted/30"),
		empty.Header(
			empty.Media(empty.MediaIcon, tabler.IconBell()),
			empty.Title("No Notifications"),
			empty.Description(h.Class("max-w-xs text-pretty"), "You're all caught up. New notifications will appear here."),
		),
		empty.Content(button.Button(button.VariantOutline, lucide.RefreshCcw(), "Refresh")),
	)
}

func NovaCard() h.Node {
	return empty.Empty(
		empty.Header(
			empty.Media(empty.MediaIcon, lucide.Folder()),
			empty.Title("No projects yet"),
			empty.Description("You haven't created any projects yet. Get started by creating your first project."),
		),
		empty.Content(
			h.Div(h.Class("flex gap-2"),
				button.Button(lib.AsChild, h.A(h.Href("#"), "Create project")),
				button.Button(button.VariantOutline, "Import project"),
			),
			button.Button(button.VariantLink, lib.AsChild, h.Class("text-muted-foreground"),
				h.A(h.Href("#"), "Learn more ", lucide.ArrowUpRight())),
		),
	)
}

func novaDemo(title, description, create, imp, learn string, icon h.Node, attrs ...any) h.Node {
	return empty.Empty(append(attrs,
		empty.Header(
			empty.Media(empty.MediaIcon, tabler.IconFolderCode()),
			empty.Title(title),
			empty.Description(description),
		),
		empty.Content(h.Class("flex-row justify-center gap-2"),
			button.Button(create),
			button.Button(button.VariantOutline, imp),
		),
		button.Button(button.VariantLink, lib.AsChild, h.Class("text-muted-foreground"), button.SizeSm,
			h.A(h.Href("#"), learn, icon)),
	)...)
}

func NovaDemo() h.Node {
	return novaDemo("No Projects Yet", "You haven't created any projects yet. Get started by creating your first project.",
		"Create Project", "Import Project", "Learn More ", lucide.ArrowUpRight())
}

func NovaRtl() h.Node {
	return style.ScopeRTL(novaDemo("لا توجد مشاريع بعد", "لم تقم بإنشاء أي مشاريع بعد. ابدأ بإنشاء مشروعك الأول.",
		"إنشاء مشروع", "استيراد مشروع", "تعرف على المزيد ",
		lucide.ArrowUpRight(h.Class("rtl:rotate-270"), h.Data("icon", "inline-end")), h.Dir("rtl")))
}
