package field

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "field", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("field-demo", Demo),
		ex("field-select", Select),
	)
}

func selectItems(values ...string) []any {
	out := make([]any, len(values))
	for i, v := range values {
		out[i] = selectui.Item(selectui.ItemValue(v), v)
	}
	return out
}

func Demo() h.Node {
	return h.Div(h.Class("w-full max-w-md"),
		h.Form(
			field.Group(
				field.Set(
					field.Legend("Payment Method"),
					field.Description("All transactions are secure and encrypted"),
					field.Group(
						field.Field(
							field.Label(h.For("checkout-7j9-card-name-43j"), "Name on Card"),
							input.Input(h.ID("checkout-7j9-card-name-43j"), h.Placeholder("Evil Rabbit"), h.Required()),
						),
						field.Field(
							field.Label(h.For("checkout-7j9-card-number-uw1"), "Card Number"),
							input.Input(h.ID("checkout-7j9-card-number-uw1"), h.Placeholder("1234 5678 9012 3456"), h.Required()),
							field.Description("Enter your 16-digit card number"),
						),
						h.Div(h.Class("grid grid-cols-3 gap-4"),
							field.Field(
								field.Label(h.For("checkout-exp-month-ts6"), "Month"),
								selectui.Select(selectui.DefaultValue(""),
									selectui.Trigger(h.ID("checkout-exp-month-ts6"), selectui.Value(selectui.Placeholder("MM"))),
									selectui.Content(selectItems("01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12")...),
								),
							),
							field.Field(
								field.Label(h.For("checkout-7j9-exp-year-f59"), "Year"),
								selectui.Select(selectui.DefaultValue(""),
									selectui.Trigger(h.ID("checkout-7j9-exp-year-f59"), selectui.Value(selectui.Placeholder("YYYY"))),
									selectui.Content(selectItems("2024", "2025", "2026", "2027", "2028", "2029")...),
								),
							),
							field.Field(
								field.Label(h.For("checkout-7j9-cvv"), "CVV"),
								input.Input(h.ID("checkout-7j9-cvv"), h.Placeholder("123"), h.Required()),
							),
						),
					),
				),
				field.Separator(),
				field.Set(
					field.Legend("Billing Address"),
					field.Description("The billing address associated with your payment method"),
					field.Group(
						field.Field(field.Horizontal,
							checkbox.Checkbox(h.ID("checkout-7j9-same-as-shipping-wgm"), checkbox.Checked(true)),
							field.Label(h.For("checkout-7j9-same-as-shipping-wgm"), h.Class("font-normal"), "Same as shipping address"),
						),
					),
				),
				field.Set(
					field.Group(
						field.Field(
							field.Label(h.For("checkout-7j9-optional-comments"), "Comments"),
							textarea.Textarea(h.ID("checkout-7j9-optional-comments"), h.Placeholder("Add any additional comments"), h.Class("resize-none")),
						),
					),
				),
				field.Field(field.Horizontal,
					button.Button(h.Type("submit"), "Submit"),
					button.Button(button.VariantOutline, h.Type("button"), "Cancel"),
				),
			),
		),
	)
}

func Select() h.Node {
	type opt struct{ value, label string }
	opts := []opt{
		{"engineering", "Engineering"}, {"design", "Design"}, {"marketing", "Marketing"},
		{"sales", "Sales"}, {"support", "Customer Support"}, {"hr", "Human Resources"},
		{"finance", "Finance"}, {"operations", "Operations"},
	}
	return h.Div(h.Class("w-full max-w-md"),
		field.Field(
			field.Label("Department"),
			selectui.Select(
				selectui.Trigger(selectui.Value(selectui.Placeholder("Choose department"))),
				selectui.Content(h.Map(opts, func(o opt) h.Node {
					return selectui.Item(selectui.ItemValue(o.value), o.label)
				})),
			),
			field.Description("Select your department or area of work."),
		),
	)
}
