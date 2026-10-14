package slider

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/slider"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/slider-*.tsx) whose content
// differs from the new-york-v4 example of the same name, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "slider", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("slider-controlled", NovaControlled),
		nova("slider-demo", NovaDemo),
		nova("slider-disabled", NovaDisabled),
		nova("slider-multiple", NovaMultiple),
		nova("slider-range", NovaRange),
		nova("slider-vertical", NovaVertical),
		nova("slider-rtl", NovaRtl),
	)
}

func NovaControlled() h.Node {
	return h.Div(h.Class("mx-auto grid w-full max-w-xs gap-3"),
		h.Div(h.Class("flex items-center justify-between gap-2"),
			label.Label(h.For("slider-demo-temperature"), "Temperature"),
			h.Span(h.Class("text-sm text-muted-foreground"), "0.3, 0.7"),
		),
		slider.Slider(h.ID("slider-demo-temperature"), slider.DefaultValue(0.3, 0.7), slider.Min(0), slider.Max(1), slider.Step(0.1)),
	)
}

func NovaDemo() h.Node {
	return slider.Slider(slider.DefaultValue(75), slider.Max(100), slider.Step(1), h.Class("mx-auto w-full max-w-xs"))
}

func NovaDisabled() h.Node {
	return slider.Slider(slider.DefaultValue(50), slider.Max(100), slider.Step(1), slider.Disabled(true), h.Class("mx-auto w-full max-w-xs"))
}

func NovaMultiple() h.Node {
	return slider.Slider(slider.DefaultValue(10, 20, 70), slider.Max(100), slider.Step(10), h.Class("mx-auto w-full max-w-xs"))
}

func NovaRange() h.Node {
	return slider.Slider(slider.DefaultValue(25, 50), slider.Max(100), slider.Step(5), h.Class("mx-auto w-full max-w-xs"))
}

func NovaVertical() h.Node {
	return h.Div(h.Class("mx-auto flex w-full max-w-xs items-center justify-center gap-6"),
		slider.Slider(slider.DefaultValue(50), slider.Max(100), slider.Step(1), slider.OrientationVertical, h.Class("h-40")),
		slider.Slider(slider.DefaultValue(25), slider.Max(100), slider.Step(1), slider.OrientationVertical, h.Class("h-40")),
	)
}

func novaRtl() h.Node {
	return slider.Slider(slider.DefaultValue(75), slider.Max(100), slider.Step(1), h.Class("mx-auto w-full max-w-xs"), h.Dir("rtl"))
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
