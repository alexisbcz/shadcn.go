// Package inputgroup ports the upstream input-group-* examples.
package inputgroup

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// autosize.js ports react-textarea-autosize (8.5.9), which
// input-group-custom uses.
//
//go:embed autosize.js
var autosizeScript string

// button.js holds input-group-button's state (copy and favorite buttons).
//
//go:embed button.js
var buttonScript string

func init() {
	runtime.Module("example-textarea-autosize", autosizeScript)
	runtime.Module("example-input-group-button", buttonScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "input-group", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("input-group-button-group", ButtonGroup),
		ex("input-group-custom", Custom),
		ex("input-group-icon", Icon),
		ex("input-group-spinner", Spinner),
		ex("input-group-text", Text),
		ex("input-group-textarea", Textarea),
		ex("input-group-button", Button),
		ex("input-group-label", Label),
		ex("input-group-tooltip", Tooltip),
	)
}

func ButtonGroup() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		buttongroup.ButtonGroup(
			buttongroup.Text(lib.AsChild, label.Label(h.For("url"), "https://")),
			ig.InputGroup(
				ig.Input(h.ID("url")),
				ig.Addon(ig.AlignInlineEnd, lucide.Link2()),
			),
			buttongroup.Text(".com"),
		),
	)
}

func Custom() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(
			h.Textarea(h.Data("slot", "input-group-control"),
				h.Class("flex field-sizing-content min-h-16 w-full resize-none rounded-md bg-transparent px-3 py-2.5 text-base transition-[color,box-shadow] outline-none md:text-sm"),
				h.Placeholder("Autoresize textarea..."),
			),
			ig.Addon(ig.AlignBlockEnd,
				ig.Button(h.Class("ml-auto"), ig.ButtonSizeSm, button.VariantDefault, "Submit"),
			),
		),
	)
}

func Icon() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(
			ig.Input(h.Placeholder("Search...")),
			ig.Addon(lucide.Search()),
		),
		ig.InputGroup(
			ig.Input(h.Type("email"), h.Placeholder("Enter your email")),
			ig.Addon(lucide.Mail()),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("Card number")),
			ig.Addon(lucide.CreditCard()),
			ig.Addon(ig.AlignInlineEnd, lucide.Check()),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("Card number")),
			ig.Addon(ig.AlignInlineEnd, lucide.Star(), lucide.Info()),
		),
	)
}

func Spinner() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-4"),
		ig.InputGroup(h.Data("disabled", "true"),
			ig.Input(h.Placeholder("Searching..."), h.Disabled()),
			ig.Addon(ig.AlignInlineEnd, spinner.Spinner()),
		),
		ig.InputGroup(h.Data("disabled", "true"),
			ig.Input(h.Placeholder("Processing..."), h.Disabled()),
			ig.Addon(spinner.Spinner()),
		),
		ig.InputGroup(h.Data("disabled", "true"),
			ig.Input(h.Placeholder("Saving changes..."), h.Disabled()),
			ig.Addon(ig.AlignInlineEnd, ig.Text("Saving..."), spinner.Spinner()),
		),
		ig.InputGroup(h.Data("disabled", "true"),
			ig.Input(h.Placeholder("Refreshing data..."), h.Disabled()),
			ig.Addon(lucide.Loader(h.Class("animate-spin"))),
			ig.Addon(ig.AlignInlineEnd, ig.Text(h.Class("text-muted-foreground"), "Please wait...")),
		),
	)
}

func Text() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(
			ig.Addon(ig.Text("$")),
			ig.Input(h.Placeholder("0.00")),
			ig.Addon(ig.AlignInlineEnd, ig.Text("USD")),
		),
		ig.InputGroup(
			ig.Addon(ig.Text("https://")),
			ig.Input(h.Placeholder("example.com"), h.Class("pl-0.5!")),
			ig.Addon(ig.AlignInlineEnd, ig.Text(".com")),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("Enter your username")),
			ig.Addon(ig.AlignInlineEnd, ig.Text("@company.com")),
		),
		ig.InputGroup(
			ig.Textarea(h.Placeholder("Enter your message")),
			ig.Addon(ig.AlignBlockEnd, ig.Text(h.Class("text-xs text-muted-foreground"), "120 characters left")),
		),
	)
}

func Textarea() h.Node {
	return h.Div(h.Class("grid w-full max-w-md gap-4"),
		ig.InputGroup(
			ig.Textarea(h.ID("textarea-code-32"), h.Placeholder("console.log('Hello, world!');"), h.Class("min-h-[200px]")),
			ig.Addon(ig.AlignBlockEnd, h.Class("border-t"),
				ig.Text("Line 1, Column 1"),
				ig.Button(ig.ButtonSizeSm, h.Class("ml-auto"), button.VariantDefault, "Run ", tabler.IconCornerDownLeft()),
			),
			ig.Addon(ig.AlignBlockStart, h.Class("border-b"),
				ig.Text(h.Class("font-mono font-medium"), tabler.IconBrandJavascript(), "script.js"),
				ig.Button(h.Class("ml-auto"), ig.ButtonSizeIconXs, tabler.IconRefresh()),
				ig.Button(button.VariantGhost, ig.ButtonSizeIconXs, tabler.IconCopy()),
			),
		),
	)
}

func Button() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(
			ig.Input(h.Placeholder("https://x.com/shadcn"), h.ReadOnly()),
			ig.Addon(ig.AlignInlineEnd,
				ig.Button(h.Aria("label", "Copy"), h.TitleAttr("Copy"), ig.ButtonSizeIconXs, tabler.IconCopy()),
			),
		),
		ig.InputGroup(h.Class("[--radius:9999px]"),
			popover.Popover(
				popover.Trigger(lib.AsChild,
					ig.Addon(ig.Button(button.VariantSecondary, ig.ButtonSizeIconXs, tabler.IconInfoCircle())),
				),
				popover.Content(popover.AlignStart, h.Class("flex flex-col gap-1 rounded-xl text-sm"),
					h.P(h.Class("font-medium"), "Your connection is not secure."),
					h.P("You should not enter any sensitive information on this site."),
				),
			),
			ig.Addon(h.Class("pl-1.5 text-muted-foreground"), "https://"),
			ig.Input(h.ID("input-secure-19")),
			ig.Addon(ig.AlignInlineEnd,
				ig.Button(ig.ButtonSizeIconXs,
					tabler.IconStar(h.Data("favorite", "false"), h.Class("data-[favorite=true]:fill-blue-600 data-[favorite=true]:stroke-blue-600")),
				),
			),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("Type to search...")),
			ig.Addon(ig.AlignInlineEnd, ig.Button(button.VariantSecondary, "Search")),
		),
	)
}

func Label() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm gap-4"),
		ig.InputGroup(
			ig.Input(h.ID("email"), h.Placeholder("shadcn")),
			ig.Addon(label.Label(h.For("email"), "@")),
		),
		ig.InputGroup(
			ig.Input(h.ID("email-2"), h.Placeholder("shadcn@vercel.com")),
			ig.Addon(ig.AlignBlockStart,
				label.Label(h.For("email-2"), h.Class("text-foreground"), "Email"),
				tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild,
						ig.Button(button.VariantGhost, h.Aria("label", "Help"), h.Class("ml-auto rounded-full"), ig.ButtonSizeIconXs, lucide.Info()),
					),
					tooltip.Content(h.P("We'll use this to send you notifications")),
				),
			),
		),
	)
}

func Tooltip() h.Node {
	tip := func(label string, icon h.Node, text string) h.Node {
		return tooltip.Tooltip(
			tooltip.Trigger(lib.AsChild, ig.Button(button.VariantGhost, h.Aria("label", label), ig.ButtonSizeIconXs, icon)),
			tooltip.Content(h.P(text)),
		)
	}
	return h.Div(h.Class("grid w-full max-w-sm gap-4"),
		ig.InputGroup(
			ig.Input(h.Placeholder("Enter password"), h.Type("password")),
			ig.Addon(ig.AlignInlineEnd, tip("Info", lucide.Info(), "Password must be at least 8 characters")),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("Your email address")),
			ig.Addon(ig.AlignInlineEnd, tip("Help", lucide.CircleHelp(), "We'll use this to send you notifications")),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("Enter API key")),
			tooltip.Tooltip(
				tooltip.Trigger(lib.AsChild,
					ig.Addon(ig.Button(button.VariantGhost, h.Aria("label", "Help"), ig.ButtonSizeIconXs, lucide.CircleHelp())),
				),
				tooltip.Content(tooltip.SideLeft, h.P("Click for help with API keys")),
			),
		),
	)
}
