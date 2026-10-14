package inputgroup

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/card"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/kbd"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/components/ui/tooltip"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the input-group page
// (ui/apps/v4/examples/radix/input-group-*.tsx) whose content differs from
// the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "input-group", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("input-group-demo", NovaDemo),
		nova("input-group-spinner", func() h.Node { return spinnerEx(false) }),
		nova("input-group-basic", NovaBasic),
		nova("input-group-block-end", NovaBlockEnd),
		nova("input-group-block-start", NovaBlockStart),
		nova("input-group-in-card", NovaInCard),
		nova("input-group-inline-end", NovaInlineEnd),
		nova("input-group-inline-start", NovaInlineStart),
		nova("input-group-kbd", NovaKbd),
		nova("input-group-textarea-examples", NovaTextareaExamples),
		nova("input-group-with-addons", NovaWithAddons),
		nova("input-group-with-buttons", NovaWithButtons),
		nova("input-group-with-kbd", NovaWithKbd),
		nova("input-group-with-tooltip", NovaWithTooltip),
	)
}

func NovaDemo() h.Node {
	return ig.InputGroup(h.Class("max-w-xs"),
		ig.Input(h.Placeholder("Search...")),
		ig.Addon(lucide.Search()),
		ig.Addon(ig.AlignInlineEnd, "12 results"),
	)
}

func NovaBasic() h.Node {
	return field.Group(
		field.Field(
			field.Label(h.For("input-default-01"), "Default (No Input Group)"),
			input.Input(h.Placeholder("Placeholder"), h.ID("input-default-01")),
		),
		field.Field(
			field.Label(h.For("input-group-02"), "Input Group"),
			ig.InputGroup(ig.Input(h.ID("input-group-02"), h.Placeholder("Placeholder"))),
		),
		field.Field(h.Data("disabled", "true"),
			field.Label(h.For("input-disabled-03"), "Disabled"),
			ig.InputGroup(ig.Input(h.ID("input-disabled-03"), h.Placeholder("This field is disabled"), h.Disabled())),
		),
		field.Field(h.Data("invalid", "true"),
			field.Label(h.For("input-invalid-04"), "Invalid"),
			ig.InputGroup(ig.Input(h.ID("input-invalid-04"), h.Placeholder("This field is invalid"), h.Aria("invalid", "true"))),
		),
	)
}

func NovaBlockEnd() h.Node {
	return field.Group(h.Class("max-w-sm"),
		field.Field(
			field.Label(h.For("block-end-input"), "Input"),
			ig.InputGroup(h.Class("h-auto"),
				ig.Input(h.ID("block-end-input"), h.Placeholder("Enter amount")),
				ig.Addon(ig.AlignBlockEnd, ig.Text("USD")),
			),
			field.Description("Footer positioned below the input."),
		),
		field.Field(
			field.Label(h.For("block-end-textarea"), "Textarea"),
			ig.InputGroup(
				ig.Textarea(h.ID("block-end-textarea"), h.Placeholder("Write a comment...")),
				ig.Addon(ig.AlignBlockEnd,
					ig.Text("0/280"),
					ig.Button(button.VariantDefault, ig.ButtonSizeSm, h.Class("ml-auto"), "Post"),
				),
			),
			field.Description("Footer positioned below the textarea."),
		),
	)
}

func NovaBlockStart() h.Node {
	return field.Group(h.Class("max-w-sm"),
		field.Field(
			field.Label(h.For("block-start-input"), "Input"),
			ig.InputGroup(h.Class("h-auto"),
				ig.Input(h.ID("block-start-input"), h.Placeholder("Enter your name")),
				ig.Addon(ig.AlignBlockStart, ig.Text("Full Name")),
			),
			field.Description("Header positioned above the input."),
		),
		field.Field(
			field.Label(h.For("block-start-textarea"), "Textarea"),
			ig.InputGroup(
				ig.Textarea(h.ID("block-start-textarea"), h.Placeholder("console.log('Hello, world!');"), h.Class("font-mono text-sm")),
				ig.Addon(ig.AlignBlockStart,
					lucide.FileCode(h.Class("text-muted-foreground")),
					ig.Text(h.Class("font-mono"), "script.js"),
					ig.Button(ig.ButtonSizeIconXs, h.Class("ml-auto"), lucide.Copy(), h.Span(h.Class("sr-only"), "Copy")),
				),
			),
			field.Description("Header positioned above the textarea."),
		),
	)
}

func NovaInCard() h.Node {
	return card.Card(h.Class("w-full"),
		card.Header(
			card.Title("Card with Input Group"),
			card.Description("This is a card with an input group."),
		),
		card.Content(
			field.Group(
				field.Field(
					field.Label(h.For("email-input"), "Email Address"),
					ig.InputGroup(
						ig.Input(h.ID("email-input"), h.Type("email"), h.Placeholder("you@example.com")),
						ig.Addon(ig.AlignInlineEnd, lucide.Mail()),
					),
				),
				field.Field(
					field.Label(h.For("website-input"), "Website URL"),
					ig.InputGroup(
						ig.Addon(ig.Text("https://")),
						ig.Input(h.ID("website-input"), h.Placeholder("example.com")),
						ig.Addon(ig.AlignInlineEnd, lucide.ExternalLink()),
					),
				),
				field.Field(
					field.Label(h.For("feedback-textarea"), "Feedback & Comments"),
					ig.InputGroup(
						ig.Textarea(h.ID("feedback-textarea"), h.Placeholder("Share your thoughts..."), h.Class("min-h-[100px]")),
						ig.Addon(ig.AlignBlockEnd, ig.Text("0/500 characters")),
					),
				),
			),
		),
		card.Footer(h.Class("justify-end gap-2"),
			button.Button(button.VariantOutline, "Cancel"),
			button.Button("Submit"),
		),
	)
}

func NovaInlineEnd() h.Node {
	return field.Field(h.Class("max-w-sm"),
		field.Label(h.For("inline-end-input"), "Input"),
		ig.InputGroup(
			ig.Input(h.ID("inline-end-input"), h.Type("password"), h.Placeholder("Enter password")),
			ig.Addon(ig.AlignInlineEnd, lucide.EyeOff()),
		),
		field.Description("Icon positioned at the end."),
	)
}

func NovaInlineStart() h.Node {
	return field.Field(h.Class("max-w-sm"),
		field.Label(h.For("inline-start-input"), "Input"),
		ig.InputGroup(
			ig.Input(h.ID("inline-start-input"), h.Placeholder("Search...")),
			ig.Addon(ig.AlignInlineStart, lucide.Search(h.Class("text-muted-foreground"))),
		),
		field.Description("Icon positioned at the start."),
	)
}

func NovaKbd() h.Node {
	return ig.InputGroup(h.Class("max-w-sm"),
		ig.Input(h.Placeholder("Search...")),
		ig.Addon(lucide.Search(h.Class("text-muted-foreground"))),
		ig.Addon(ig.AlignInlineEnd, kbd.Kbd("⌘K")),
	)
}

func NovaTextareaExamples() h.Node {
	return field.Group(
		field.Field(
			field.Label(h.For("textarea-header-footer-12"), "Default Textarea (No Input Group)"),
			textarea.Textarea(h.ID("textarea-header-footer-12"), h.Placeholder("Enter your text here...")),
		),
		field.Field(
			field.Label(h.For("textarea-header-footer-13"), "Input Group"),
			ig.InputGroup(ig.Textarea(h.ID("textarea-header-footer-13"), h.Placeholder("Enter your text here..."))),
			field.Description("This is a description of the input group."),
		),
		field.Field(h.Data("invalid", "true"),
			field.Label(h.For("textarea-header-footer-14"), "Invalid"),
			ig.InputGroup(ig.Textarea(h.ID("textarea-header-footer-14"), h.Placeholder("Enter your text here..."), h.Aria("invalid", "true"))),
			field.Description("This is a description of the input group."),
		),
		field.Field(h.Data("disabled", "true"),
			field.Label(h.For("textarea-header-footer-15"), "Disabled"),
			ig.InputGroup(ig.Textarea(h.ID("textarea-header-footer-15"), h.Placeholder("Enter your text here..."), h.Disabled())),
			field.Description("This is a description of the input group."),
		),
		field.Field(
			field.Label(h.For("prompt-31"), "Addon (block-start)"),
			ig.InputGroup(
				ig.Textarea(h.ID("prompt-31")),
				ig.Addon(ig.AlignBlockStart,
					ig.Text("Ask, Search or Chat..."),
					lucide.Info(h.Class("ml-auto text-muted-foreground")),
				),
			),
			field.Description("This is a description of the input group."),
		),
		field.Field(
			field.Label(h.For("textarea-header-footer-30"), "Addon (block-end)"),
			ig.InputGroup(
				ig.Textarea(h.ID("textarea-header-footer-30"), h.Placeholder("Enter your text here...")),
				ig.Addon(ig.AlignBlockEnd,
					ig.Text("0/280 characters"),
					ig.Button(button.VariantDefault, ig.ButtonSizeIconXs, h.Class("ml-auto rounded-full"),
						lucide.ArrowUp(), h.Span(h.Class("sr-only"), "Send")),
				),
			),
		),
		field.Field(
			field.Label(h.For("textarea-comment-31"), "Addon (Buttons)"),
			ig.InputGroup(
				ig.Textarea(h.ID("textarea-comment-31"), h.Placeholder("Share your thoughts..."), h.Class("min-h-[120px]")),
				ig.Addon(ig.AlignBlockEnd,
					ig.Button(button.VariantGhost, h.Class("ml-auto"), ig.ButtonSizeSm, "Cancel"),
					ig.Button(button.VariantDefault, ig.ButtonSizeSm, "Post Comment"),
				),
			),
		),
		field.Field(
			field.Label(h.For("textarea-code-32"), "Code Editor"),
			ig.InputGroup(
				ig.Textarea(h.ID("textarea-code-32"), h.Placeholder("console.log('Hello, world!');"), h.Class("min-h-[300px] py-3")),
				ig.Addon(ig.AlignBlockStart, h.Class("border-b"),
					ig.Text(h.Class("font-mono font-medium"), lucide.Code(), "script.js"),
					ig.Button(ig.ButtonSizeIconXs, h.Class("ml-auto"), lucide.RefreshCw()),
					ig.Button(ig.ButtonSizeIconXs, button.VariantGhost, lucide.Copy()),
				),
				ig.Addon(ig.AlignBlockEnd, h.Class("border-t"),
					ig.Text("Line 1, Column 1"),
					ig.Text(h.Class("ml-auto"), "JavaScript"),
				),
			),
		),
	)
}

func NovaWithAddons() h.Node {
	return field.Group(
		field.Field(
			field.Label(h.For("input-icon-left-05"), "Addon (inline-start)"),
			ig.InputGroup(
				ig.Input(h.ID("input-icon-left-05")),
				ig.Addon(lucide.Search(h.Class("text-muted-foreground"))),
			),
		),
		field.Field(
			field.Label(h.For("input-icon-right-07"), "Addon (inline-end)"),
			ig.InputGroup(
				ig.Input(h.ID("input-icon-right-07")),
				ig.Addon(ig.AlignInlineEnd, lucide.EyeOff()),
			),
		),
		field.Field(
			field.Label(h.For("input-icon-both-09"), "Addon (inline-start and inline-end)"),
			ig.InputGroup(
				ig.Input(h.ID("input-icon-both-09")),
				ig.Addon(lucide.Mic(h.Class("text-muted-foreground"))),
				ig.Addon(ig.AlignInlineEnd, lucide.Radio(h.Class("animate-pulse text-red-500"))),
			),
		),
		field.Field(
			field.Label(h.For("input-addon-20"), "Addon (block-start)"),
			ig.InputGroup(h.Class("h-auto"),
				ig.Input(h.ID("input-addon-20")),
				ig.Addon(ig.AlignBlockStart,
					ig.Text("First Name"),
					lucide.Info(h.Class("ml-auto text-muted-foreground")),
				),
			),
		),
		field.Field(
			field.Label(h.For("input-addon-21"), "Addon (block-end)"),
			ig.InputGroup(h.Class("h-auto"),
				ig.Input(h.ID("input-addon-21")),
				ig.Addon(ig.AlignBlockEnd,
					ig.Text("20/240 characters"),
					lucide.Info(h.Class("ml-auto text-muted-foreground")),
				),
			),
		),
		field.Field(
			field.Label(h.For("input-icon-both-10"), "Multiple Icons"),
			ig.InputGroup(
				ig.Input(h.ID("input-icon-both-10")),
				ig.Addon(ig.AlignInlineEnd,
					lucide.Star(),
					ig.Button(ig.ButtonSizeIconXs, lucide.Copy()),
				),
				ig.Addon(lucide.Radio(h.Class("animate-pulse text-red-500"))),
			),
		),
		field.Field(
			field.Label(h.For("input-description-10"), "Description"),
			ig.InputGroup(
				ig.Input(h.ID("input-description-10")),
				ig.Addon(ig.AlignInlineEnd, lucide.Info()),
			),
			field.Description("This is a description of the input group."),
		),
		field.Field(
			field.Label(h.For("input-label-10"), "Label"),
			ig.InputGroup(
				ig.Addon(field.Label(h.For("input-label-10"), "Label")),
				ig.Input(h.ID("input-label-10")),
			),
			ig.InputGroup(
				ig.Input(h.ID("input-optional-12"), h.Aria("label", "Optional")),
				ig.Addon(ig.AlignInlineEnd, ig.Text("(optional)")),
			),
		),
	)
}

func NovaWithButtons() h.Node {
	group := func(id string, align ig.Align, btn h.Node) h.Node {
		return ig.InputGroup(ig.Input(h.ID(id)), ig.Addon(align, btn))
	}
	return field.Group(
		field.Field(
			field.Label(h.For("input-button-13"), "Button"),
			group("input-button-13", ig.AlignInlineStart, ig.Button("Default")),
			group("input-button-14", ig.AlignInlineStart, ig.Button(button.VariantOutline, "Outline")),
			group("input-button-15", ig.AlignInlineStart, ig.Button(button.VariantSecondary, "Secondary")),
			group("input-button-16", ig.AlignInlineEnd, ig.Button(button.VariantSecondary, "Button")),
			group("input-button-17", ig.AlignInlineEnd, ig.Button(ig.ButtonSizeIconXs, lucide.Copy())),
			group("input-button-18", ig.AlignInlineEnd, ig.Button(button.VariantSecondary, ig.ButtonSizeIconXs, lucide.Trash())),
		),
	)
}

func NovaWithKbd() h.Node {
	info := func(id, placeholder string) h.Node {
		return ig.InputGroup(
			ig.Input(h.ID(id), h.Placeholder(placeholder)),
			ig.Addon(ig.AlignInlineEnd, lucide.Info()),
		)
	}
	return field.Group(
		field.Field(
			field.Label(h.For("input-kbd-22"), "Input Group with Kbd"),
			ig.InputGroup(ig.Input(h.ID("input-kbd-22")), ig.Addon(kbd.Kbd("⌘K"))),
			ig.InputGroup(ig.Input(h.ID("input-kbd-23")), ig.Addon(ig.AlignInlineEnd, kbd.Kbd("⌘K"))),
			ig.InputGroup(
				ig.Input(h.ID("input-search-apps-24"), h.Placeholder("Search for Apps...")),
				ig.Addon(ig.AlignInlineEnd, "Ask AI"),
				ig.Addon(ig.AlignInlineEnd, kbd.Kbd("Tab")),
			),
			ig.InputGroup(
				ig.Input(h.ID("input-search-type-25"), h.Placeholder("Type to search...")),
				ig.Addon(ig.AlignInlineStart, lucide.Sparkles()),
				ig.Addon(ig.AlignInlineEnd, kbd.Group(kbd.Kbd("Ctrl"), kbd.Kbd("C"))),
			),
		),
		field.Field(
			field.Label(h.For("input-username-26"), "Username"),
			ig.InputGroup(
				ig.Input(h.ID("input-username-26"), h.Value("shadcn")),
				ig.Addon(ig.AlignInlineEnd,
					h.Div(h.Class("flex size-4 items-center justify-center rounded-full bg-green-500 dark:bg-green-800"),
						lucide.Check(h.Class("size-3 text-white")),
					),
				),
			),
			field.Description(h.Class("text-green-700"), "This username is available."),
		),
		ig.InputGroup(
			ig.Input(h.ID("input-search-docs-27"), h.Placeholder("Search documentation...")),
			ig.Addon(lucide.Search()),
			ig.Addon(ig.AlignInlineEnd, "12 results"),
		),
		ig.InputGroup(h.Data("disabled", "true"),
			ig.Input(h.ID("input-search-disabled-28"), h.Placeholder("Search documentation..."), h.Disabled()),
			ig.Addon(lucide.Search()),
			ig.Addon(ig.AlignInlineEnd, "Disabled"),
		),
		field.Group(h.Class("grid grid-cols-2 gap-4"),
			field.Field(
				field.Label(h.For("input-group-11"), "First Name"),
				info("input-group-11", "First Name"),
			),
			field.Field(
				field.Label(h.For("input-group-12"), "Last Name"),
				info("input-group-12", "Last Name"),
			),
		),
		field.Field(h.Data("disabled", "true"),
			field.Label(h.For("input-group-29"), `Loading ("data-disabled="true")`),
			ig.InputGroup(
				ig.Input(h.ID("input-group-29"), h.Disabled(), h.Value("shadcn")),
				ig.Addon(ig.AlignInlineEnd, spinner.Spinner()),
			),
			field.Description("This is a description of the input group."),
		),
	)
}

// NovaWithTooltip is rendered by the docs without its country prop (an
// empty dropdown label).
func NovaWithTooltip() h.Node {
	return field.Group(
		field.Field(
			field.Label(h.For("input-tooltip-20"), "Tooltip"),
			ig.InputGroup(
				ig.Input(h.ID("input-tooltip-20")),
				ig.Addon(ig.AlignInlineEnd,
					tooltip.Tooltip(
						tooltip.Trigger(lib.AsChild, ig.Button(h.Class("rounded-full"), ig.ButtonSizeIconXs, lucide.Info())),
						tooltip.Content("This is content in a tooltip."),
					),
				),
			),
			field.Description("This is a description of the input group."),
		),
		field.Field(
			field.Label(h.For("input-dropdown-21"), "Dropdown"),
			ig.InputGroup(
				ig.Input(h.ID("input-dropdown-21")),
				ig.Addon(
					dm.DropdownMenu(
						dm.Trigger(lib.AsChild, ig.Button(h.Class("text-muted-foreground tabular-nums"), " ", lucide.ChevronDown())),
						dm.Content(dm.AlignStart, h.Class("min-w-16"), dm.SideOffset(10), dm.AlignOffset(-8),
							dm.Item("+1"), dm.Item("+44"), dm.Item("+46"),
						),
					),
				),
			),
			field.Description("This is a description of the input group."),
		),
		field.Field(
			field.Label(h.For("input-secure-19"), "Popover"),
			ig.InputGroup(
				popover.Popover(
					popover.Trigger(lib.AsChild, ig.Addon(ig.Button(button.VariantSecondary, ig.ButtonSizeIconXs, lucide.Info()))),
					popover.Content(popover.AlignStart,
						popover.Header(
							popover.Title("Your connection is not secure."),
							popover.Description("You should not enter any sensitive information on this site."),
						),
					),
				),
				ig.Addon(h.Class("pl-1 text-muted-foreground"), "https://"),
				ig.Input(h.ID("input-secure-19")),
				ig.Addon(ig.AlignInlineEnd, ig.Button(ig.ButtonSizeIconXs, lucide.Star())),
			),
			field.Description("This is a description of the input group."),
		),
		field.Field(
			field.Label(h.For("url"), "Button Group"),
			buttongroup.ButtonGroup(
				buttongroup.Text("https://"),
				ig.InputGroup(
					ig.Input(h.ID("url")),
					ig.Addon(ig.AlignInlineEnd, lucide.Info()),
				),
				buttongroup.Text(".com"),
			),
			field.Description("This is a description of the input group."),
		),
	)
}

func init() {
	registry.Register(registry.Example{Name: "nova--input-group-rtl", Component: "input-group", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, style.ScopeRTL(NovaRtl()))
	})})
}

// NovaRtl is input-group-rtl in Arabic (the docs' default language for
// it).
func NovaRtl() h.Node {
	rtl := h.Dir("rtl")
	return h.Div(h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(h.Class("max-w-xs"), rtl,
			ig.Input(h.Placeholder("بحث...")),
			ig.Addon(lucide.Search()),
			ig.Addon(ig.AlignInlineEnd, "١٢ نتيجة"),
		),
		ig.InputGroup(rtl,
			ig.Input(h.Placeholder("جاري البحث...")),
			ig.Addon(ig.AlignInlineEnd, spinner.Spinner()),
		),
		ig.InputGroup(rtl,
			ig.Input(h.Placeholder("جاري حفظ التغييرات...")),
			ig.Addon(ig.AlignInlineEnd, ig.Text("جاري الحفظ..."), spinner.Spinner()),
		),
		field.Group(h.Class("max-w-sm"), rtl,
			field.Field(
				field.Label(h.For("rtl-textarea"), "منطقة النص"),
				ig.InputGroup(rtl,
					ig.Textarea(h.ID("rtl-textarea"), h.Placeholder("اكتب تعليقًا...")),
					ig.Addon(ig.AlignBlockEnd,
						ig.Text("٠/٢٨٠"),
						ig.Button(button.VariantDefault, ig.ButtonSizeSm, h.Class("ml-auto"), "نشر"),
					),
				),
				field.Description("تذييل موضع أسفل منطقة النص."),
			),
		),
	)
}
