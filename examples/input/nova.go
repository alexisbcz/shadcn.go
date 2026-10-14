package input

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the input page (ui/apps/v4/examples/radix/input-*.tsx),
// in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "input", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("input-badge", NovaBadge),
		nova("input-basic", NovaBasic),
		nova("input-button-group", NovaButtonGroup),
		nova("input-demo", NovaDemo),
		nova("input-disabled", NovaDisabled),
		nova("input-field", NovaField),
		nova("input-fieldgroup", NovaFieldgroup),
		nova("input-file", NovaFile),
		nova("input-form", NovaForm),
		nova("input-grid", NovaGrid),
		nova("input-inline", NovaInline),
		nova("input-input-group", NovaInputGroup),
		nova("input-invalid", NovaInvalid),
		nova("input-required", NovaRequired),
		nova("input-rtl", NovaRtl),
	)
}

func NovaBadge() h.Node {
	return field.Field(
		field.Label(h.For("input-badge"), "Webhook URL", " ", badge.Badge(badge.VariantSecondary, h.Class("ml-auto"), "Beta")),
		input.Input(h.ID("input-badge"), h.Type("url"), h.Placeholder("https://api.example.com/webhook")),
	)
}

func NovaBasic() h.Node { return input.Input(h.Placeholder("Enter text")) }

func NovaButtonGroup() h.Node {
	return field.Field(
		field.Label(h.For("input-button-group"), "Search"),
		buttongroup.ButtonGroup(
			input.Input(h.ID("input-button-group"), h.Placeholder("Type to search...")),
			button.Button(button.VariantOutline, "Search"),
		),
	)
}

func described(id, label, description string, inputAttrs []any, attrs ...any) h.Node {
	return describedClass(id, label, description, "", inputAttrs, attrs...)
}

func describedClass(id, label, description, descriptionClass string, inputAttrs []any, attrs ...any) h.Node {
	return field.Field(append(attrs,
		field.Label(h.For(id), label),
		input.Input(append([]any{h.ID(id)}, inputAttrs...)...),
		field.Description(h.If(descriptionClass != "", h.Class(descriptionClass)), description),
	)...)
}

func NovaDemo() h.Node {
	return described("input-demo-api-key", "API Key", "Your API key is encrypted and stored securely.",
		[]any{h.Type("password"), h.Placeholder("sk-...")})
}

func NovaDisabled() h.Node {
	return described("input-demo-disabled", "Email", "This field is currently disabled.",
		[]any{h.Type("email"), h.Placeholder("Email"), h.Bool("disabled")}, h.Data("disabled", "true"))
}

func NovaField() h.Node {
	return described("input-field-username", "Username", "Choose a unique username for your account.",
		[]any{h.Type("text"), h.Placeholder("Enter your username")})
}

func NovaFieldgroup() h.Node {
	return field.Group(
		field.Field(
			field.Label(h.For("fieldgroup-name"), "Name"),
			input.Input(h.ID("fieldgroup-name"), h.Placeholder("Jordan Lee")),
		),
		field.Field(
			field.Label(h.For("fieldgroup-email"), "Email"),
			input.Input(h.ID("fieldgroup-email"), h.Type("email"), h.Placeholder("name@example.com")),
			field.Description("We'll send updates to this address."),
		),
		field.Field(field.Horizontal,
			button.Button(h.Type("reset"), button.VariantOutline, "Reset"),
			button.Button(h.Type("submit"), "Submit"),
		),
	)
}

func NovaFile() h.Node {
	return described("picture", "Picture", "Select a picture to upload.", []any{h.Type("file")})
}

func NovaForm() h.Node {
	return h.Form(h.Class("w-full max-w-sm"),
		field.Group(
			field.Field(
				field.Label(h.For("form-name"), "Name"),
				input.Input(h.ID("form-name"), h.Type("text"), h.Placeholder("Evil Rabbit"), h.Required()),
			),
			field.Field(
				field.Label(h.For("form-email"), "Email"),
				input.Input(h.ID("form-email"), h.Type("email"), h.Placeholder("john@example.com")),
				field.Description("We'll never share your email with anyone."),
			),
			h.Div(h.Class("grid grid-cols-2 gap-4"),
				field.Field(
					field.Label(h.For("form-phone"), "Phone"),
					input.Input(h.ID("form-phone"), h.Type("tel"), h.Placeholder("+1 (555) 123-4567")),
				),
				field.Field(
					field.Label(h.For("form-country"), "Country"),
					selectui.Select(selectui.DefaultValue("us"),
						selectui.Trigger(h.ID("form-country"), selectui.Value()),
						selectui.Content(
							selectui.Item(selectui.ItemValue("us"), "United States"),
							selectui.Item(selectui.ItemValue("uk"), "United Kingdom"),
							selectui.Item(selectui.ItemValue("ca"), "Canada"),
						),
					),
				),
			),
			field.Field(
				field.Label(h.For("form-address"), "Address"),
				input.Input(h.ID("form-address"), h.Type("text"), h.Placeholder("123 Main St")),
			),
			field.Field(field.Horizontal,
				button.Button(h.Type("button"), button.VariantOutline, "Cancel"),
				button.Button(h.Type("submit"), "Submit"),
			),
		),
	)
}

func NovaGrid() h.Node {
	return field.Group(h.Class("grid max-w-sm grid-cols-2"),
		field.Field(field.Label(h.For("first-name"), "First Name"), input.Input(h.ID("first-name"), h.Placeholder("Jordan"))),
		field.Field(field.Label(h.For("last-name"), "Last Name"), input.Input(h.ID("last-name"), h.Placeholder("Lee"))),
	)
}

func NovaInline() h.Node {
	return field.Field(field.Horizontal,
		input.Input(h.Type("search"), h.Placeholder("Search...")),
		button.Button("Search"),
	)
}

func NovaInputGroup() h.Node {
	return field.Field(
		field.Label(h.For("input-group-url"), "Website URL"),
		inputgroup.InputGroup(
			inputgroup.Input(h.ID("input-group-url"), h.Placeholder("example.com")),
			inputgroup.Addon(inputgroup.Text("https://")),
			inputgroup.Addon(inputgroup.AlignInlineEnd, lucide.Info()),
		),
	)
}

func NovaInvalid() h.Node {
	return described("input-invalid", "Invalid Input", "This field contains validation errors.",
		[]any{h.Placeholder("Error"), h.Aria("invalid", "true")}, h.Data("invalid", "true"))
}

func NovaRequired() h.Node {
	return field.Field(
		field.Label(h.For("input-required"), "Required Field ", h.Span(h.Class("text-destructive"), "*")),
		input.Input(h.ID("input-required"), h.Placeholder("This field is required"), h.Required()),
		field.Description("This field must be filled out."),
	)
}

func NovaRtl() h.Node {
	// text-start: the ui-rtl FieldDescription's text-left → text-start
	// (field does not map its classes under style.ScopeRTL).
	return style.ScopeRTL(describedClass("input-rtl-api-key", "مفتاح API", "مفتاح API الخاص بك مشفر ومخزن بأمان.", "text-start",
		[]any{h.Type("password"), h.Placeholder("sk-..."), h.Dir("rtl")}, h.Dir("rtl")))
}
