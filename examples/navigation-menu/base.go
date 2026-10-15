package navigationmenu

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	nm "shadcn.go/components/ui/navigation-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The base navigation-menu-demo preview
// (ui/apps/v4/examples/base/navigation-menu-demo.tsx): the radix demo plus
// a "With Icon" menu before Docs.
func init() {
	registry.Register(registry.Example{Name: "base--navigation-menu-demo", Component: "navigation-menu", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, BaseDemo())
	})})
}

// BaseDemo is examples/base/navigation-menu-demo.tsx.
func BaseDemo() h.Node {
	t := novaEnglish
	comps := make([]any, len(t.items))
	for i, c := range t.items {
		comps[i] = novaListItem(components[i].href, c[0], c[1])
	}
	iconLink := func(icon h.Node, text string) h.Node {
		return nm.Link(lib.AsChild, h.A(h.Href("#"), h.Class("flex-row items-center gap-2"), icon, text))
	}
	return nm.NavigationMenu(
		nm.List(
			nm.Item(
				nm.Trigger(t.gettingStarted),
				nm.Content(
					h.Ul(h.Class("w-96"),
						novaListItem("/docs", t.introduction, t.introductionDesc),
						novaListItem("/docs/installation", t.installation, t.installationDesc),
						novaListItem("/docs/primitives/typography", t.typography, t.typographyDesc),
					),
				),
			),
			nm.Item(h.Class("hidden md:flex"),
				nm.Trigger(t.components),
				nm.Content(
					h.Ul(append([]any{h.Class("grid w-[400px] gap-2 md:w-[500px] md:grid-cols-2 lg:w-[600px]")}, comps...)...),
				),
			),
			nm.Item(
				nm.Trigger("With Icon"),
				nm.Content(
					h.Ul(h.Class("grid w-[200px]"),
						h.Li(
							iconLink(lucide.CircleAlert(), "Backlog"),
							iconLink(lucide.CircleDashed(), "To Do"),
							iconLink(lucide.CircleCheck(), "Done"),
						),
					),
				),
			),
			nm.Item(
				nm.Link(lib.AsChild, h.Class(nm.TriggerStyleFor(style.Nova)), h.A(h.Href("/docs"), t.docs)),
			),
		),
	)
}
