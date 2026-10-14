// Package togglegroup ports the upstream toggle-group-* examples.
package togglegroup

import (
	"shadcn.go/components/icons/lucide"
	tg "shadcn.go/components/ui/toggle-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "toggle-group", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("toggle-group-demo", Demo),
		ex("toggle-group-disabled", Disabled),
		ex("toggle-group-lg", Lg),
		ex("toggle-group-outline", Outline),
		ex("toggle-group-single", Single),
		ex("toggle-group-sm", Sm),
		ex("toggle-group-spacing", Spacing),
		ex("toggle-group-behavior", Behavior),
	)
}

// textItems are the bold/italic/strikethrough items every example shares.
func textItems() []any {
	return []any{
		tg.Item(tg.Value("bold"), h.Aria("label", "Toggle bold"), lucide.Bold(h.Class("h-4 w-4"))),
		tg.Item(tg.Value("italic"), h.Aria("label", "Toggle italic"), lucide.Italic(h.Class("h-4 w-4"))),
		tg.Item(tg.Value("strikethrough"), h.Aria("label", "Toggle strikethrough"), lucide.Underline(h.Class("h-4 w-4"))),
	}
}

func Demo() h.Node {
	return tg.ToggleGroup(tg.VariantOutline, tg.TypeMultiple, textItems())
}

func Disabled() h.Node {
	return tg.ToggleGroup(tg.TypeMultiple, tg.Disabled(true), textItems())
}

func Lg() h.Node {
	return tg.ToggleGroup(tg.TypeMultiple, tg.SizeLg, textItems())
}

func Outline() h.Node {
	return tg.ToggleGroup(tg.TypeMultiple, tg.VariantOutline, textItems())
}

func Single() h.Node {
	return tg.ToggleGroup(tg.TypeSingle, textItems())
}

func Sm() h.Node {
	return tg.ToggleGroup(tg.TypeSingle, tg.SizeSm, textItems())
}

func Spacing() h.Node {
	return tg.ToggleGroup(tg.TypeMultiple, tg.VariantOutline, tg.Spacing(2), tg.SizeSm,
		tg.Item(tg.Value("star"), h.Aria("label", "Toggle star"),
			h.Class("data-[state=on]:bg-transparent data-[state=on]:*:[svg]:fill-yellow-500 data-[state=on]:*:[svg]:stroke-yellow-500"),
			lucide.Star(), "Star"),
		tg.Item(tg.Value("heart"), h.Aria("label", "Toggle heart"),
			h.Class("data-[state=on]:bg-transparent data-[state=on]:*:[svg]:fill-red-500 data-[state=on]:*:[svg]:stroke-red-500"),
			lucide.Heart(), "Heart"),
		tg.Item(tg.Value("bookmark"), h.Aria("label", "Toggle bookmark"),
			h.Class("data-[state=on]:bg-transparent data-[state=on]:*:[svg]:fill-blue-500 data-[state=on]:*:[svg]:stroke-blue-500"),
			lucide.Bookmark(), "Bookmark"),
	)
}

// Behavior ports the toggle-group-behavior fixture
// (parity/fixtures/toggle-group-behavior.tsx): an initial value, vertical
// orientation, Loop(false), a disabled item, dir="rtl" and RovingFocus(false).
func Behavior() h.Node {
	return h.Div(h.Class("flex flex-col gap-6"),
		tg.ToggleGroup(tg.TypeMultiple, tg.VariantOutline, tg.Value("italic"), tg.OrientationVertical, tg.Loop(false), h.Class("flex-col"),
			tg.Item(tg.Value("bold"), h.Aria("label", "Toggle bold"), lucide.Bold(h.Class("h-4 w-4"))),
			tg.Item(tg.Value("italic"), h.Aria("label", "Toggle italic"), lucide.Italic(h.Class("h-4 w-4"))),
			tg.Item(tg.Value("underline"), h.Aria("label", "Toggle underline"), h.Disabled(), lucide.Underline(h.Class("h-4 w-4"))),
		),
		tg.ToggleGroup(tg.TypeSingle, tg.Value("center"), h.Dir("rtl"), tg.Spacing(1),
			tg.Item(tg.Value("left"), "Left"),
			tg.Item(tg.Value("center"), "Center"),
			tg.Item(tg.Value("right"), "Right"),
		),
		tg.ToggleGroup(tg.TypeSingle, tg.RovingFocus(false),
			tg.Item(tg.Value("a"), "A"),
			tg.Item(tg.Value("b"), "B"),
		),
	)
}
