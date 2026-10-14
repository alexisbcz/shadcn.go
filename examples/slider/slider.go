// Package slider ports the upstream slider-* examples.
package slider

import (
	"shadcn.go/components/ui/slider"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "slider", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("slider-demo", Demo),
		ex("slider-states", States),
		ex("slider-vertical", Vertical),
		ex("slider-form", Form),
	)
}

func Demo() h.Node {
	return slider.Slider(slider.DefaultValue(50), slider.Max(100), slider.Step(1), h.Class("w-[60%]"))
}

// States ports the slider-states fixture (parity/fixtures/slider-states.tsx):
// the Radix Slider props the upstream example does not use: a range (two
// thumbs) with minStepsBetweenThumbs, three thumbs, a fractional step with
// custom min and max, a coarse step, inverted, right to left, disabled, no
// default value (shadcn's hidden second thumb), and a slider in a form
// (hidden bubble inputs).
func States() h.Node {
	return h.Div(h.Class("flex w-[60%] flex-col gap-6"),
		slider.Slider(slider.DefaultValue(25, 75), slider.MinStepsBetweenThumbs(10)),
		slider.Slider(slider.DefaultValue(10, 40, 80)),
		slider.Slider(slider.DefaultValue(0.5), slider.Min(-1), slider.Max(2), slider.Step(0.1)),
		slider.Slider(slider.DefaultValue(20), slider.Step(10)),
		slider.Slider(slider.DefaultValue(30), slider.Inverted(true)),
		slider.Slider(slider.DefaultValue(30), h.Dir("rtl")),
		slider.Slider(slider.DefaultValue(60), slider.Disabled(true)),
		slider.Slider(),
		h.Form(
			slider.Slider(slider.DefaultValue(20, 60), slider.Name("price")),
		),
	)
}

// Vertical ports the slider-vertical fixture
// (parity/fixtures/slider-vertical.tsx): vertical sliders, plain, range,
// inverted and disabled.
func Vertical() h.Node {
	return h.Div(h.Class("flex h-48 gap-10"),
		slider.Slider(slider.DefaultValue(50), slider.OrientationVertical),
		slider.Slider(slider.DefaultValue(20, 70), slider.OrientationVertical),
		slider.Slider(slider.DefaultValue(30), slider.OrientationVertical, slider.Inverted(true)),
		slider.Slider(slider.DefaultValue(40), slider.OrientationVertical, slider.Disabled(true)),
	)
}

// Form ports the slider-form fixture (parity/fixtures/slider-form.tsx):
// sliders tied to a form through the form attribute (outside the form),
// right to left combined with inverted, and a right-to-left range.
func Form() h.Node {
	return h.Div(h.Class("flex w-[60%] flex-col gap-6"),
		h.Form(h.ID("slider-form")),
		slider.Slider(slider.DefaultValue(40), h.Name("volume"), h.FormAttr("slider-form")),
		slider.Slider(slider.DefaultValue(20, 70), slider.Name("price"), slider.Form("slider-form")),
		slider.Slider(slider.DefaultValue(30), h.Dir("rtl"), slider.Inverted(true)),
		slider.Slider(slider.DefaultValue(10, 90), h.Dir("rtl")),
	)
}
