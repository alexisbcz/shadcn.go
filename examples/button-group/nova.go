package buttongroup

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/button-group"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/field"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the button-group page
// (ui/apps/v4/examples/radix/button-group-*.tsx) whose content differs from
// the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "button-group", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("button-group-demo", func() h.Node { return demo("w-40") }),
		nova("button-group-dropdown", func() h.Node { return dropdown("w-44") }),
		nova("button-group-nested", NovaNested),
		nova("button-group-popover", NovaPopover),
	)
}

func NovaNested() h.Node {
	return buttongroup.ButtonGroup(
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeIcon, lucide.Plus()),
		),
		buttongroup.ButtonGroup(
			ig.InputGroup(
				ig.Input(h.Placeholder("Send a message...")),
				tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild, ig.Addon(ig.AlignInlineEnd, lucide.AudioLines())),
					tooltip.Content("Voice Mode"),
				),
			),
		),
	)
}

func NovaPopover() h.Node {
	return buttongroup.ButtonGroup(
		button.Button(button.VariantOutline, lucide.Bot(), " Copilot"),
		popover.Popover(
			// Content first: see Popover.
			popover.Content(popover.AlignEnd, h.Class("rounded-xl text-sm"),
				popover.Header(
					popover.Title("Start a new task with Copilot"),
					popover.Description("Describe your task in natural language."),
				),
				field.Field(
					field.Label(h.For("task"), h.Class("sr-only"), "Task Description"),
					textarea.Textarea(h.ID("task"), h.Placeholder("I need to..."), h.Class("resize-none")),
					field.Description("Copilot will open a pull request for review."),
				),
			),
			popover.Trigger(lib.AsChild,
				button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Open Popover"), lucide.ChevronDown()),
			),
		),
	)
}

func init() {
	registry.Register(registry.Example{Name: "nova--button-group-rtl", Component: "button-group", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, style.ScopeRTL(NovaRtl()))
	})})
}

// NovaRtl is button-group-rtl in Arabic (the docs' default language for
// it).
func NovaRtl() h.Node {
	return h.Div(h.Dir("rtl"),
		buttongroup.ButtonGroup(
			buttongroup.ButtonGroup(h.Class("hidden sm:flex"),
				button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Go Back"), lucide.ArrowLeft(h.Class("rtl:rotate-180"))),
			),
			buttongroup.ButtonGroup(
				button.Button(button.VariantOutline, "أرشفة"),
				button.Button(button.VariantOutline, "تقرير"),
			),
			buttongroup.ButtonGroup(
				button.Button(button.VariantOutline, "تأجيل"),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "More Options"), lucide.Ellipsis())),
					dm.Content(dm.AlignEnd, h.Data("lang", "ar"), h.Class("w-40"),
						dm.Group(
							dm.Item(lucide.MailCheck(), "وضع علامة كمقروء"),
							dm.Item(lucide.Archive(), "أرشفة"),
						),
						dm.Separator(),
						dm.Group(
							dm.Item(lucide.Clock(), "تأجيل"),
							dm.Item(lucide.CalendarPlus(), "إضافة إلى التقويم"),
							dm.Item(lucide.ListFilter(), "إضافة إلى القائمة"),
							dm.Sub(
								dm.SubTrigger(lucide.Tag(), "تصنيف كـ..."),
								dm.SubContent(h.Data("lang", "ar"),
									dm.RadioGroup(dm.Value("personal"),
										dm.RadioItem(dm.Value("personal"), "شخصي"),
										dm.RadioItem(dm.Value("work"), "عمل"),
										dm.RadioItem(dm.Value("other"), "آخر"),
									),
								),
							),
						),
						dm.Separator(),
						dm.Group(dm.Item(dm.VariantDestructive, lucide.Trash2(), "سلة المهملات")),
					),
				),
			),
		),
	)
}
