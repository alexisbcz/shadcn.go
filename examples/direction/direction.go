// Package direction ports the direction parity fixture (upstream ships
// DirectionProvider without new-york-v4 examples).
package direction

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/direction"
	"shadcn.go/components/ui/label"
	rg "shadcn.go/components/ui/radio-group"
	"shadcn.go/components/ui/slider"
	"shadcn.go/components/ui/tabs"
	tg "shadcn.go/components/ui/toggle-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(registry.Example{Name: "direction-provider", Component: "direction", Render: registry.Static(Provider)})
}

func items() []any {
	return []any{
		tg.Item(tg.Value("bold"), h.Aria("label", "Toggle bold"), lucide.Bold()),
		tg.Item(tg.Value("italic"), h.Aria("label", "Toggle italic"), lucide.Italic()),
		tg.Item(tg.Value("underline"), h.Aria("label", "Toggle underline"), lucide.Underline()),
	}
}

// Provider is the parity fixture direction-provider.
func Provider() h.Node {
	radio := func(value, id, text string) h.Node {
		return h.Div(h.Class("flex items-center gap-2"), rg.Item(rg.Value(value), h.ID(id)), label.Label(h.For(id), text))
	}
	return direction.Provider(direction.RTL,
		h.Div(h.Class("flex w-full max-w-md flex-col gap-6"),
			tabs.Tabs(tabs.Value("one"),
				tabs.List(
					tabs.Trigger(tabs.Value("one"), "One"),
					tabs.Trigger(tabs.Value("two"), "Two"),
					tabs.Trigger(tabs.Value("three"), "Three"),
				),
				tabs.Content(tabs.Value("one"), "First panel."),
				tabs.Content(tabs.Value("two"), "Second panel."),
				tabs.Content(tabs.Value("three"), "Third panel."),
			),
			rg.RadioGroup(rg.Value("a"), h.Class("flex"),
				radio("a", "dir-a", "First"),
				radio("b", "dir-b", "Second"),
			),
			tg.ToggleGroup(tg.TypeMultiple, tg.VariantOutline, items()),
			slider.Slider(slider.DefaultValue(30), slider.Max(100), slider.Step(1)),
			tg.ToggleGroup(tg.TypeSingle, h.Dir("ltr"), items()),
			direction.Provider(direction.LTR,
				slider.Slider(slider.DefaultValue(30), slider.Max(100), slider.Step(1)),
			),
		),
	)
}
