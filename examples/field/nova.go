package field

import (
	"fmt"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	rg "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/slider"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the field page (ui/apps/v4/examples/radix/field-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "field", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("field-demo", func() h.Node { return demo(true) }),
		nova("field-checkbox", NovaCheckbox),
		nova("field-choice-card", NovaChoiceCard),
		nova("field-fieldset", NovaFieldset),
		nova("field-group", NovaGroup),
		nova("field-input", NovaInput),
		nova("field-radio", NovaRadio),
		nova("field-responsive", NovaResponsive),
		nova("field-select", NovaSelect),
		nova("field-slider", NovaSlider),
		nova("field-switch", NovaSwitch),
		nova("field-textarea", NovaTextarea),
	)
}

func NovaCheckbox() h.Node {
	opt := func(id, label string) h.Node {
		return field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID(id)),
			field.Label(h.For(id), h.Class("font-normal"), label),
		)
	}
	return field.Group(h.Class("w-full max-w-xs"),
		field.Set(
			field.Legend(field.LegendLabel, "Show these items on the desktop"),
			field.Description("Select the items you want to show on the desktop."),
			field.Group(h.Class("gap-3"),
				opt("finder-pref-9k2-hard-disks-ljj", "Hard disks"),
				opt("finder-pref-9k2-external-disks-1yg", "External disks"),
				opt("finder-pref-9k2-cds-dvds-fzt", "CDs, DVDs, and iPods"),
				opt("finder-pref-9k2-connected-servers-6l2", "Connected servers"),
			),
		),
		field.Separator(),
		field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("finder-pref-9k2-sync-folders-nep"), checkbox.Checked(true)),
			field.Content(
				field.Label(h.For("finder-pref-9k2-sync-folders-nep"), "Sync Desktop & Documents folders"),
				field.Description("Your Desktop & Documents folders are being synced with iCloud Drive. You can access them from other devices."),
			),
		),
	)
}

func NovaChoiceCard() h.Node {
	card := func(id, value, title, desc string) h.Node {
		return field.Label(h.For(id),
			field.Field(field.Horizontal,
				field.Content(field.Title(title), field.Description(desc)),
				rg.Item(rg.Value(value), h.ID(id)),
			),
		)
	}
	return field.Group(h.Class("w-full max-w-xs"),
		field.Set(
			field.Legend(field.LegendLabel, "Compute Environment"),
			field.Description("Select the compute environment for your cluster."),
			rg.RadioGroup(rg.Value("kubernetes"),
				card("kubernetes-r2h", "kubernetes", "Kubernetes", "Run GPU workloads on a K8s cluster."),
				card("vm-z4k", "vm", "Virtual Machine", "Access a cluster to run GPU workloads."),
			),
		),
	)
}

func NovaFieldset() h.Node {
	f := func(id, label, placeholder string) h.Node {
		return field.Field(
			field.Label(h.For(id), label),
			input.Input(h.ID(id), h.Type("text"), h.Placeholder(placeholder)),
		)
	}
	return field.Set(h.Class("w-full max-w-sm"),
		field.Legend("Address Information"),
		field.Description("We need your address to deliver your order."),
		field.Group(
			f("street", "Street Address", "123 Main St"),
			h.Div(h.Class("grid grid-cols-2 gap-4"),
				f("city", "City", "New York"),
				f("zip", "Postal Code", "90502"),
			),
		),
	)
}

func NovaGroup() h.Node {
	opt := func(id, label string, extra ...any) h.Node {
		return field.Field(field.Horizontal,
			checkbox.Checkbox(append([]any{h.ID(id)}, extra...)...),
			field.Label(h.For(id), h.Class("font-normal"), label),
		)
	}
	return field.Group(h.Class("w-full max-w-xs"),
		field.Set(
			field.Label("Responses"),
			field.Description("Get notified when ChatGPT responds to requests that take time, like research or image generation."),
			field.Group(h.Data("slot", "checkbox-group"),
				opt("push", "Push notifications", checkbox.Checked(true), h.Disabled()),
			),
		),
		field.Separator(),
		field.Set(
			field.Label("Tasks"),
			field.Description("Get notified when tasks you've created have updates.", h.Comment(""), " ", h.A(h.Href("#"), "Manage tasks")),
			field.Group(h.Data("slot", "checkbox-group"),
				opt("push-tasks", "Push notifications"),
				opt("email-tasks", "Email notifications"),
			),
		),
	)
}

func NovaInput() h.Node {
	return field.Set(h.Class("w-full max-w-xs"),
		field.Group(
			field.Field(
				field.Label(h.For("username"), "Username"),
				input.Input(h.ID("username"), h.Type("text"), h.Placeholder("Max Leiter")),
				field.Description("Choose a unique username for your account."),
			),
			field.Field(
				field.Label(h.For("password"), "Password"),
				field.Description("Must be at least 8 characters long."),
				input.Input(h.ID("password"), h.Type("password"), h.Placeholder("••••••••")),
			),
		),
	)
}

func NovaRadio() h.Node {
	opt := func(value, id, label string) h.Node {
		return field.Field(field.Horizontal,
			rg.Item(rg.Value(value), h.ID(id)),
			field.Label(h.For(id), h.Class("font-normal"), label),
		)
	}
	return field.Set(h.Class("w-full max-w-xs"),
		field.Legend(field.LegendLabel, "Subscription Plan"),
		field.Description("Yearly and lifetime plans offer significant savings."),
		rg.RadioGroup(rg.Value("monthly"),
			opt("monthly", "plan-monthly", "Monthly ($9.99/month)"),
			opt("yearly", "plan-yearly", "Yearly ($99.99/year)"),
			opt("lifetime", "plan-lifetime", "Lifetime ($299.99)"),
		),
	)
}

func NovaResponsive() h.Node {
	return h.Div(h.Class("w-full max-w-lg"),
		h.Form(
			field.Set(
				field.Legend("Profile"),
				field.Description("Fill in your profile information."),
				field.Group(
					field.Field(field.Responsive,
						field.Content(
							field.Label(h.For("name"), "Name"),
							field.Description("Provide your full name for identification"),
						),
						input.Input(h.ID("name"), h.Placeholder("Evil Rabbit"), h.Required()),
					),
					field.Field(field.Responsive,
						button.Button(h.Type("submit"), "Submit"),
						button.Button(h.Type("button"), button.VariantOutline, "Cancel"),
					),
				),
			),
		),
	)
}

func NovaSelect() h.Node {
	type opt struct{ value, label string }
	opts := []opt{
		{"engineering", "Engineering"}, {"design", "Design"}, {"marketing", "Marketing"},
		{"sales", "Sales"}, {"support", "Customer Support"}, {"hr", "Human Resources"},
		{"finance", "Finance"}, {"operations", "Operations"},
	}
	return field.Field(h.Class("w-full max-w-xs"),
		field.Label("Department"),
		selectui.Select(
			selectui.Trigger(selectui.Value(selectui.Placeholder("Choose department"))),
			selectui.Content(selectui.Group(h.Map(opts, func(o opt) h.Node {
				return selectui.Item(selectui.ItemValue(o.value), o.label)
			}))),
		),
		field.Description("Select your department or area of work."),
	)
}

func NovaSlider() h.Node {
	return field.Field(h.Class("w-full max-w-xs"),
		field.Title("Price Range"),
		field.Description(
			"Set your budget range ($",
			h.Span(h.Class("font-medium tabular-nums"), "200"),
			" -", h.Comment(""), " ",
			h.Span(h.Class("font-medium tabular-nums"), "800"),
			").",
		),
		slider.Slider(slider.DefaultValue(200, 800), slider.Max(1000), slider.Min(0), slider.Step(10),
			h.Class("mt-2 w-full"), h.Aria("label", "Price Range")),
	)
}

func NovaSwitch() h.Node {
	return field.Field(field.Horizontal, h.Class("w-fit"),
		field.Label(h.For("2fa"), "Multi-factor authentication"),
		switchui.Switch(h.ID("2fa")),
	)
}

func NovaTextarea() h.Node {
	return field.Set(h.Class("w-full max-w-xs"),
		field.Group(
			field.Field(
				field.Label(h.For("feedback"), "Feedback"),
				textarea.Textarea(h.ID("feedback"), h.Placeholder("Your feedback helps us improve..."), h.Attr("rows", "4")),
				field.Description("Share your thoughts about our service."),
			),
		),
	)
}

func init() {
	registry.Register(registry.Example{Name: "nova--field-rtl", Component: "field", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, style.ScopeRTL(NovaRtl()))
	})})
}

// NovaRtl is field-rtl in Arabic (the docs' default language for it).
func NovaRtl() h.Node {
	items := func(labels ...string) []any {
		out := make([]any, len(labels))
		for i, l := range labels {
			out[i] = selectui.Item(selectui.ItemValue(l), l)
		}
		return out
	}
	months := []any{selectui.Item(selectui.ItemValue("MM"), "ش.ش")}
	for i, ar := range []string{"٠١", "٠٢", "٠٣", "٠٤", "٠٥", "٠٦", "٠٧", "٠٨", "٠٩", "١٠", "١١", "١٢"} {
		months = append(months, selectui.Item(selectui.ItemValue(fmt.Sprintf("%02d", i+1)), ar))
	}
	return h.Div(h.Class("w-full max-w-md py-6"), h.Dir("rtl"),
		h.Form(
			field.Group(
				field.Set(
					field.Legend("طريقة الدفع"),
					field.Description("جميع المعاملات آمنة ومشفرة"),
					field.Group(
						field.Field(
							field.Label(h.For("checkout-7j9-card-name-43j-rtl"), "الاسم على البطاقة"),
							input.Input(h.ID("checkout-7j9-card-name-43j-rtl"), h.Placeholder("Evil Rabbit"), h.Required()),
						),
						field.Field(
							field.Label(h.For("checkout-7j9-card-number-uw1-rtl"), "رقم البطاقة"),
							input.Input(h.ID("checkout-7j9-card-number-uw1-rtl"), h.Placeholder("1234 5678 9012 3456"), h.Required()),
							field.Description("أدخل رقم البطاقة المكون من 16 رقمًا"),
						),
						h.Div(h.Class("grid grid-cols-3 gap-4"),
							field.Field(
								field.Label(h.For("checkout-exp-month-ts6-rtl"), "الشهر"),
								selectui.Select(
									selectui.Trigger(h.ID("checkout-exp-month-ts6-rtl"), selectui.Value()),
									selectui.Content(h.Dir("rtl"), selectui.Group(months...)),
								),
							),
							field.Field(
								field.Label(h.For("checkout-7j9-exp-year-f59-rtl"), "السنة"),
								selectui.Select(
									selectui.Trigger(h.ID("checkout-7j9-exp-year-f59-rtl"), selectui.Value()),
									selectui.Content(h.Dir("rtl"), selectui.Group(items("YYYY", "2024", "2025", "2026", "2027", "2028", "2029")...)),
								),
							),
							field.Field(
								field.Label(h.For("checkout-7j9-cvv-rtl"), "CVV"),
								input.Input(h.ID("checkout-7j9-cvv-rtl"), h.Placeholder("123"), h.Required()),
							),
						),
					),
				),
				field.Separator(),
				field.Set(
					field.Legend("عنوان الفوترة"),
					field.Description("عنوان الفوترة المرتبط بطريقة الدفع الخاصة بك"),
					field.Group(
						field.Field(field.Horizontal,
							checkbox.Checkbox(h.ID("checkout-7j9-same-as-shipping-wgm-rtl"), checkbox.Checked(true)),
							field.Label(h.For("checkout-7j9-same-as-shipping-wgm-rtl"), h.Class("font-normal"), "نفس عنوان الشحن"),
						),
					),
				),
				field.Set(
					field.Group(
						field.Field(
							field.Label(h.For("checkout-7j9-optional-comments-rtl"), "تعليقات"),
							textarea.Textarea(h.ID("checkout-7j9-optional-comments-rtl"), h.Placeholder("أضف أي تعليقات إضافية"), h.Class("resize-none")),
						),
					),
				),
				field.Field(field.Horizontal,
					button.Button(h.Type("submit"), "إرسال"),
					button.Button(button.VariantOutline, h.Type("button"), "إلغاء"),
				),
			),
		),
	)
}
