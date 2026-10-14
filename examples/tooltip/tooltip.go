// Package tooltip ports the upstream tooltip-* examples and the tooltip
// parity fixtures.
package tooltip

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	hovercard "shadcn.go/components/ui/hover-card"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "tooltip", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("tooltip-demo", Demo),
		ex("tooltip-provider-demo", ProviderDemo),
		ex("tooltip-sides", Sides),
		ex("tooltip-delays", Delays),
		ex("tooltip-default-open", DefaultOpen),
	)
}

func Demo() h.Node {
	return tooltip.Tooltip(
		tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Hover")),
		tooltip.Content(h.P("Add to library")),
	)
}

// ProviderDemo is the parity fixture tooltip-provider-demo.
func ProviderDemo() h.Node { return tooltip.Provider(Demo()) }

// Sides is the parity fixture tooltip-sides.
func Sides() h.Node {
	type variant struct {
		label string
		props []any
	}
	variants := []variant{
		{"top", nil},
		{"right", []any{tooltip.SideRight}},
		{"bottom", []any{tooltip.SideBottom, tooltip.SideOffset(8)}},
		{"left", []any{tooltip.SideLeft, tooltip.SideOffset(4)}},
		{"top start", []any{tooltip.AlignStart}},
		{"bottom end", []any{tooltip.SideBottom, tooltip.AlignEnd, tooltip.AlignOffset(6)}},
	}
	return tooltip.Provider(
		h.Div(h.Class("grid grid-cols-3 gap-x-16 gap-y-10"),
			h.Map(variants, func(v variant) h.Node {
				return tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, v.label)),
					tooltip.Content(append(v.props, h.P("Tooltip "+v.label))...),
				)
			}),
		),
		h.Div(h.Class("fixed top-2 left-2"),
			tooltip.Tooltip(
				tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, "edge")),
				tooltip.Content("Flipped below"),
			),
		),
		h.Div(h.Class("fixed right-2 bottom-2"),
			tooltip.Tooltip(
				tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, "labelled")),
				tooltip.Content(h.Aria("label", "Accessible label"), "Visible text"),
			),
		),
	)
}

// Delays is the parity fixture tooltip-delays.
func Delays() h.Node {
	trigger := func(label string) h.Node {
		return tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("w-32"), label))
	}
	return tooltip.Provider(tooltip.DelayDuration(600), tooltip.SkipDelayDuration(400),
		h.Div(h.Class("fixed top-[200px] left-[200px] flex gap-8"),
			tooltip.Tooltip(trigger("A"), tooltip.Content("Tooltip A")),
			tooltip.Tooltip(trigger("B"), tooltip.Content(tooltip.SideBottom, "Tooltip B")),
			tooltip.Tooltip(tooltip.DelayDuration(0), trigger("C"), tooltip.Content("Tooltip C")),
			tooltip.Tooltip(tooltip.DisableHoverableContent(true), trigger("D"), tooltip.Content("Tooltip D")),
		),
	)
}

// DefaultOpen is the parity fixture tooltip-default-open.
func DefaultOpen() h.Node {
	return h.Div(h.Class("flex gap-24"),
		tooltip.Provider(
			tooltip.Tooltip(tooltip.DefaultOpen,
				tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Tooltip")),
				tooltip.Content(tooltip.SideRight, "Open on load"),
			),
		),
		hovercard.HoverCard(hovercard.DefaultOpen,
			hovercard.Trigger(lib.AsChild, button.Button(button.VariantLink, "Hover card")),
			hovercard.Content("Open on load too."),
		),
	)
}
