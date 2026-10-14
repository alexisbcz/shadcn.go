// Package spinner ports the upstream spinner-* examples that need no other
// components yet.
package spinner

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "spinner", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("spinner-basic", Basic),
		ex("spinner-size", Size),
		ex("spinner-color", Color),
		ex("spinner-custom", Custom),
		ex("spinner-button", Button),
	)
}

func Basic() h.Node {
	return h.Div(h.Class("flex flex-col items-center justify-center gap-8"), spinner.Spinner())
}

func Size() h.Node {
	return h.Div(h.Class("flex items-center gap-6"),
		spinner.Spinner(h.Class("size-3")),
		spinner.Spinner(h.Class("size-4")),
		spinner.Spinner(h.Class("size-6")),
		spinner.Spinner(h.Class("size-8")),
	)
}

func Color() h.Node {
	return h.Div(h.Class("flex items-center gap-6"),
		h.Map([]string{"text-red-500", "text-green-500", "text-blue-500", "text-yellow-500", "text-purple-500"}, func(c string) h.Node {
			return spinner.Spinner(h.Class("size-6 " + c))
		}),
	)
}

// Custom mirrors the example's local Spinner built on LoaderIcon.
func Custom() h.Node {
	custom := func(class string) *h.Element {
		return h.Spread(lucide.Loader(h.Role("status"), h.Aria("label", "Loading"), h.Class("size-4 animate-spin")), h.Svg(h.Class(class)), lib.Cn)
	}
	return h.Div(h.Class("flex items-center gap-4"), custom(""))
}

func Button() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-4"),
		button.Button(h.Disabled(), button.SizeSm, spinner.Spinner(), "Loading..."),
		button.Button(button.VariantOutline, h.Disabled(), button.SizeSm, spinner.Spinner(), "Please wait"),
		button.Button(button.VariantSecondary, h.Disabled(), button.SizeSm, spinner.Spinner(), "Processing"),
	)
}
