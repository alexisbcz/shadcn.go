package breadcrumb

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/breadcrumb-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "breadcrumb", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("breadcrumb-basic", NovaBasic),
		nova("breadcrumb-demo", NovaDemo),
		nova("breadcrumb-dropdown", NovaDropdown),
		nova("breadcrumb-separator", NovaSeparator),
		nova("breadcrumb-rtl", NovaRtl),
	)
}

func NovaBasic() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		breadcrumb.Item(breadcrumb.Link(h.Href("#"), "Home")),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Link(h.Href("#"), "Components")),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

func novaMenuItems(labels ...string) h.Node {
	var items []any
	for _, l := range labels {
		items = append(items, dm.Item(l))
	}
	return dm.Group(items...)
}

func NovaDemo() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("#", "Home"),
		breadcrumb.Separator(),
		breadcrumb.Item(dm.DropdownMenu(
			dm.Trigger(lib.AsChild,
				button.Button(button.SizeIconSm, button.VariantGhost,
					breadcrumb.Ellipsis(),
					h.Span(h.Class("sr-only"), "Toggle menu"),
				),
			),
			dm.Content(dm.AlignStart, novaMenuItems("Documentation", "Themes", "GitHub")),
		)),
		breadcrumb.Separator(),
		link("#", "Components"),
		breadcrumb.Separator(),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

type dropdownText struct{ home, components, documentation, themes, github, breadcrumb string }

func novaDropdown(t dropdownText, rootAttrs []any, contentAttrs ...any) h.Node {
	return breadcrumb.Breadcrumb(append(rootAttrs, breadcrumb.List(
		link("/", t.home),
		breadcrumb.Separator(lucide.Dot()),
		breadcrumb.Item(dm.DropdownMenu(
			dm.Trigger(lib.AsChild,
				h.Button(h.Class("flex items-center gap-1"),
					t.components,
					lucide.ChevronDown(h.Class("size-3.5")),
				),
			),
			dm.Content(append(contentAttrs, novaMenuItems(t.documentation, t.themes, t.github))...),
		)),
		breadcrumb.Separator(lucide.Dot()),
		breadcrumb.Item(breadcrumb.Page(t.breadcrumb)),
	))...)
}

func NovaDropdown() h.Node {
	return novaDropdown(dropdownText{"Home", "Components", "Documentation", "Themes", "GitHub", "Breadcrumb"}, nil, dm.AlignStart)
}

func novaRtl() h.Node {
	return novaDropdown(dropdownText{"الرئيسية", "المكونات", "التوثيق", "السمات", "جيت هاب", "مسار التنقل"},
		[]any{h.Dir("rtl")}, h.Data("lang", "ar"))
}

func NovaSeparator() h.Node {
	return breadcrumb.Breadcrumb(breadcrumb.List(
		link("/", "Home"),
		breadcrumb.Separator(lucide.Dot()),
		link("/components", "Components"),
		breadcrumb.Separator(lucide.Dot()),
		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
	))
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
