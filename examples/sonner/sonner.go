// Package sonner ports the upstream sonner-* examples.
package sonner

import (
	_ "embed"
	"net/http"

	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/sonner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// examples.js is the examples' onClick handlers: they call sonner's toast()
// like the React examples do.
//
//go:embed examples.js
var examplesScript string

func init() {
	runtime.Module("example-sonner", examplesScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "sonner", Render: registry.Static(fn), Extra: toaster}
	}
	registry.Register(
		ex("sonner-demo", Demo),
		ex("sonner-types", Types),
		ex("sonner-options", Options),
	)
}

// toaster is the <Toaster /> the upstream docs mount in their layout. Its
// theme comes from next-themes' useTheme(), which reports "light" there
// (the provider forces the page theme but `theme` stays the default
// "light" when system themes are disabled), so sonner's own theme is light
// in both page themes; the normal colors follow the page through shadcn's
// --normal-* variables.
func toaster(*http.Request) h.Node { return sonner.Toaster(sonner.ThemeLight) }

func Demo() h.Node {
	return button.Button(button.VariantOutline, "Show Toast")
}

func Types() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-2"),
		button.Button(button.VariantOutline, "Default"),
		button.Button(button.VariantOutline, "Success"),
		button.Button(button.VariantOutline, "Info"),
		button.Button(button.VariantOutline, "Warning"),
		button.Button(button.VariantOutline, "Error"),
		button.Button(button.VariantOutline, "Promise"),
	)
}

// Options is the parity fixture sonner-options.
func Options() h.Node {
	return h.Div(h.Class("flex flex-wrap gap-2"),
		button.Button(button.VariantOutline, "Rich"),
		button.Button(button.VariantOutline, "Confirm"),
		button.Button(button.VariantOutline, "Top"),
		button.Button(button.VariantOutline, "Dismiss all"),
		button.Button(button.VariantOutline, "Reject"),
		button.Button(button.VariantOutline, "Sticky"),
	)
}
