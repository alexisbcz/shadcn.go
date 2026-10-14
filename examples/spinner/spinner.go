// Package spinner ports the upstream spinner-* examples.
package spinner

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/empty"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/progress"
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
		ex("spinner-badge", Badge),
		ex("spinner-demo", Demo),
		ex("spinner-empty", Empty),
		ex("spinner-input-group", InputGroup),
		ex("spinner-item", Item),
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

func Badge() h.Node {
	return h.Div(h.Class("flex items-center gap-4 [--radius:1.2rem]"),
		badge.Badge(spinner.Spinner(), "Syncing"),
		badge.Badge(badge.VariantSecondary, spinner.Spinner(), "Updating"),
		badge.Badge(badge.VariantOutline, spinner.Spinner(), "Processing"),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex w-full max-w-xs flex-col gap-4 [--radius:1rem]"),
		item.Item(item.VariantMuted,
			item.Media(spinner.Spinner()),
			item.Content(item.Title(h.Class("line-clamp-1"), "Processing payment...")),
			item.Content(h.Class("flex-none justify-end"), h.Span(h.Class("text-sm tabular-nums"), "$100.00")),
		),
	)
}

func Empty() h.Node {
	return empty.Empty(h.Class("w-full"),
		empty.Header(
			empty.Media(empty.MediaIcon, spinner.Spinner()),
			empty.Title("Processing your request"),
			empty.Description("Please wait while we process your request. Do not refresh the page."),
		),
		empty.Content(button.Button(button.VariantOutline, button.SizeSm, "Cancel")),
	)
}

func InputGroup() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-4"),
		ig.InputGroup(
			ig.Input(h.Placeholder("Send a message..."), h.Disabled()),
			ig.Addon(ig.AlignInlineEnd, spinner.Spinner()),
		),
		ig.InputGroup(
			ig.Textarea(h.Placeholder("Send a message..."), h.Disabled()),
			ig.Addon(ig.AlignBlockEnd,
				spinner.Spinner(), " Validating...",
				ig.Button(h.Class("ml-auto"), button.VariantDefault, lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send")),
			),
		),
	)
}

func Item() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-4 [--radius:1rem]"),
		item.Item(item.VariantOutline,
			item.Media(item.MediaIcon, spinner.Spinner()),
			item.Content(item.Title("Downloading..."), item.Description("129 MB / 1000 MB")),
			item.Actions(h.Class("hidden sm:flex"), button.Button(button.VariantOutline, button.SizeSm, "Cancel")),
			item.Footer(progress.Progress(progress.Value(75))),
		),
	)
}
