package buttongroup

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/runtime"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// currency.js holds button-group-select's currency state: the trigger
// shows the selected value.
//
//go:embed currency.js
var currencyScript string

func init() {
	runtime.Module("example-button-group-currency", currencyScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "button-group", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("button-group-demo", Demo),
		ex("button-group-dropdown", Dropdown),
		ex("button-group-select", Select),
	)
}

func Demo() h.Node { return demo("w-52") }

func demo(contentClass string) h.Node {
	return buttongroup.ButtonGroup(
		buttongroup.ButtonGroup(h.Class("hidden sm:flex"),
			button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Go Back"), lucide.ArrowLeft()),
		),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, "Archive"),
			button.Button(button.VariantOutline, "Report"),
		),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, "Snooze"),
			dm.DropdownMenu(
				dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "More Options"), lucide.Ellipsis())),
				dm.Content(dm.AlignEnd, h.Class(contentClass),
					dm.Group(
						dm.Item(lucide.MailCheck(), "Mark as Read"),
						dm.Item(lucide.Archive(), "Archive"),
					),
					dm.Separator(),
					dm.Group(
						dm.Item(lucide.Clock(), "Snooze"),
						dm.Item(lucide.CalendarPlus(), "Add to Calendar"),
						dm.Item(lucide.ListFilter(), "Add to List"),
						dm.Sub(
							dm.SubTrigger(lucide.Tag(), "Label As..."),
							dm.SubContent(
								dm.RadioGroup(dm.Value("personal"),
									dm.RadioItem(dm.Value("personal"), "Personal"),
									dm.RadioItem(dm.Value("work"), "Work"),
									dm.RadioItem(dm.Value("other"), "Other"),
								),
							),
						),
					),
					dm.Separator(),
					dm.Group(
						dm.Item(dm.VariantDestructive, lucide.Trash2(), "Trash"),
					),
				),
			),
		),
	)
}

func Dropdown() h.Node { return dropdown("[--radius:1rem]") }

func dropdown(contentClass string) h.Node {
	return buttongroup.ButtonGroup(
		button.Button(button.VariantOutline, "Follow"),
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("pl-2!"), lucide.ChevronDown())),
			dm.Content(dm.AlignEnd, h.Class(contentClass),
				dm.Group(
					dm.Item(lucide.VolumeOff(), "Mute Conversation"),
					dm.Item(lucide.Check(), "Mark as Read"),
					dm.Item(lucide.TriangleAlert(), "Report Conversation"),
					dm.Item(lucide.UserRoundX(), "Block User"),
					dm.Item(lucide.Share(), "Share Conversation"),
					dm.Item(lucide.Copy(), "Copy Conversation"),
				),
				dm.Separator(),
				dm.Group(
					dm.Item(dm.VariantDestructive, lucide.Trash(), "Delete Conversation"),
				),
			),
		),
	)
}

func Select() h.Node {
	currencies := []struct{ value, label string }{
		{"$", "US Dollar"},
		{"€", "Euro"},
		{"£", "British Pound"},
	}
	items := []any{h.Class("min-w-24")}
	for _, c := range currencies {
		items = append(items, selectui.Item(selectui.ItemValue(c.value),
			c.value, " ", h.Span(h.Class("text-muted-foreground"), c.label),
		))
	}
	return buttongroup.ButtonGroup(
		buttongroup.ButtonGroup(
			selectui.Select(selectui.DefaultValue("$"),
				selectui.Trigger(h.Class("font-mono"), h.Data("example-currency", ""), "$"),
				selectui.Content(items...),
			),
			input.Input(h.Placeholder("10.00"), h.Attr("pattern", "[0-9]*")),
		),
		buttongroup.ButtonGroup(
			button.Button(h.Aria("label", "Send"), button.SizeIcon, button.VariantOutline, lucide.ArrowRight()),
		),
	)
}
