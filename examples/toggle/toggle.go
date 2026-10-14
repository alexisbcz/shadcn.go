// Package toggle ports the upstream toggle-* examples.
package toggle

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/toggle"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "toggle", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("toggle-demo", Demo),
		ex("toggle-disabled", Disabled),
		ex("toggle-lg", Lg),
		ex("toggle-outline", Outline),
		ex("toggle-sm", Sm),
		ex("toggle-with-text", WithText),
		ex("toggle-states", States),
	)
}

func Demo() h.Node {
	return toggle.Toggle(
		h.Aria("label", "Toggle bookmark"),
		toggle.SizeSm,
		toggle.VariantOutline,
		h.Class("data-[state=on]:bg-transparent data-[state=on]:*:[svg]:fill-blue-500 data-[state=on]:*:[svg]:stroke-blue-500"),
		lucide.Bookmark(),
		"Bookmark",
	)
}

func Disabled() h.Node {
	return toggle.Toggle(h.Aria("label", "Toggle italic"), h.Disabled(), lucide.Underline(h.Class("h-4 w-4")))
}

func Lg() h.Node {
	return toggle.Toggle(toggle.SizeLg, h.Aria("label", "Toggle italic"), lucide.Italic())
}

func Outline() h.Node {
	return toggle.Toggle(toggle.VariantOutline, h.Aria("label", "Toggle italic"), lucide.Italic())
}

func Sm() h.Node {
	return toggle.Toggle(toggle.SizeSm, h.Aria("label", "Toggle italic"), lucide.Italic())
}

func WithText() h.Node {
	return toggle.Toggle(h.Aria("label", "Toggle italic"), lucide.Italic(), "Italic")
}

// States ports the toggle-states fixture (parity/fixtures/toggle-states.tsx):
// initially pressed toggles (default and outline) and a pressed disabled one.
func States() h.Node {
	return h.Div(h.Class("flex items-center gap-2"),
		toggle.Toggle(toggle.Pressed(true), h.Aria("label", "Toggle bold"), lucide.Bold()),
		toggle.Toggle(toggle.VariantOutline, toggle.Pressed(true), h.Aria("label", "Toggle italic"), lucide.Italic(), "Italic"),
		toggle.Toggle(toggle.Pressed(true), h.Disabled(), h.Aria("label", "Toggle locked"), lucide.Bold()),
	)
}
