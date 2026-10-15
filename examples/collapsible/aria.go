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
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The aria collapsible-settings preview
// (ui/apps/v4/examples/aria/collapsible-settings.tsx): React Aria's
// DisclosurePanel wraps the extra fields in a div and, once hydrated, hides
// itself with hidden="until-found", so the closed panel still takes a row
// (and a gap) of the field group's grid.
func init() {
	registry.Register(registry.Example{Name: "aria--collapsible-settings", Component: "collapsible", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, AriaSettings())
	})})
}

// AriaSettings is examples/aria/collapsible-settings.tsx.
func AriaSettings() h.Node {
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
					collapsible.Content(h.Data("hidden-until-found", ""),
						h.Div(h.Class("col-span-full grid grid-cols-subgrid gap-2"),
							radius("Radius X", "radius-x"),
							radius("Radius Y", "radius-y"),
						),
					),
				),
				collapsible.Trigger(lib.AsChild,
					button.Button(button.VariantOutline, button.SizeIcon,
						lucide.Maximize(h.Class("[[data-state=open]>&]:hidden")),
						lucide.Minimize(h.Class("[[data-state=closed]>&]:hidden")),
					),
				),
			),
		),
	)
}
