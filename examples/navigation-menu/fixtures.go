package navigationmenu

import (
	nm "shadcn.go/components/ui/navigation-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "navigation-menu", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("navigation-menu-viewport", ViewportFixture),
		ex("navigation-menu-options", OptionsFixture),
		ex("navigation-menu-vertical", VerticalFixture),
	)
}

func plainLink(href, text string, extra ...any) h.Node {
	return h.Li(nm.Link(append([]any{h.Href(href)}, append(extra, text)...)...))
}

// ViewportFixture is the parity fixture navigation-menu-viewport.
func ViewportFixture() h.Node {
	return nm.NavigationMenu(nm.DelayDuration(0),
		nm.List(
			nm.Item(
				nm.Trigger("Getting started"),
				nm.Content(h.Ul(h.Class("grid w-[400px] gap-2"),
					plainLink("#intro", "Introduction", nm.Active),
					plainLink("#install", "Installation"),
				)),
			),
			nm.Item(
				nm.Trigger("Components"),
				nm.Content(h.Ul(h.Class("grid w-[260px] gap-2"),
					plainLink("#button", "Button"),
					plainLink("#card", "Card"),
					plainLink("#dialog", "Dialog"),
				)),
			),
			nm.Item(
				nm.Trigger(nm.Disabled, "Disabled"),
				nm.Content(h.Ul(h.Class("w-[200px]"), plainLink("#never", "Never"))),
			),
			nm.Item(
				nm.Link(h.Href("#docs"), h.Class(nm.TriggerStyle()), "Docs"),
			),
			nm.Indicator(),
		),
	)
}

// OptionsFixture is the parity fixture navigation-menu-options.
func OptionsFixture() h.Node {
	return nm.NavigationMenu(nm.UseViewport(false), nm.DefaultValue("two"), nm.DelayDuration(0), nm.SkipDelayDuration(0), nm.Dir("rtl"),
		nm.List(
			nm.Item(nm.ItemValue("one"),
				nm.Trigger("One"),
				nm.Content(h.Ul(h.Class("grid w-[200px] gap-2"), plainLink("#a", "Alpha"), plainLink("#b", "Beta"))),
			),
			nm.Item(nm.ItemValue("two"),
				nm.Trigger("Two"),
				nm.Content(h.Ul(h.Class("grid w-[240px] gap-2"), plainLink("#c", "Gamma"), plainLink("#d", "Delta"))),
			),
			nm.Item(nm.ItemValue("three"),
				nm.Trigger("Three"),
				nm.Content(h.Ul(h.Class("grid w-[200px] gap-2"), plainLink("#e", "Epsilon"))),
			),
		),
	)
}

// VerticalFixture is the parity fixture navigation-menu-vertical.
func VerticalFixture() h.Node {
	contentClass := h.Class("group-data-[viewport=false]/navigation-menu:top-0 group-data-[viewport=false]/navigation-menu:mt-0 md:left-full md:ml-2")
	return nm.NavigationMenu(nm.OrientationVertical, nm.UseViewport(false), nm.DelayDuration(0),
		nm.List(h.Class("flex-col items-start"),
			nm.Item(nm.ItemValue("one"),
				nm.Trigger("One"),
				nm.Content(contentClass, h.Ul(h.Class("grid w-[200px] gap-2"), plainLink("#a", "Alpha"), plainLink("#b", "Beta"))),
			),
			nm.Item(nm.ItemValue("two"),
				nm.Trigger("Two"),
				nm.Content(contentClass, h.Ul(h.Class("grid w-[240px] gap-2"), plainLink("#c", "Gamma"), plainLink("#d", "Delta"), plainLink("#e", "Epsilon"))),
			),
			nm.Item(nm.ItemValue("three"),
				nm.Link(h.Href("#three"), "Three"),
			),
		),
	)
}
