package radiogroup

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/label"
	rg "shadcn.go/components/ui/radio-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/radio-group-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "radio-group", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("radio-group-demo", NovaDemo),
		nova("radio-group-choice-card", NovaChoiceCard),
		nova("radio-group-description", NovaDescription),
		nova("radio-group-disabled", NovaDisabled),
		nova("radio-group-fieldset", NovaFieldset),
		nova("radio-group-invalid", NovaInvalid),
		nova("radio-group-rtl", NovaRtl),
	)
}

func NovaDemo() h.Node {
	item := func(value, id, text string) h.Node {
		return h.Div(h.Class("flex items-center gap-3"),
			rg.Item(rg.Value(value), h.ID(id)),
			label.Label(h.For(id), text),
		)
	}
	return rg.RadioGroup(rg.Value("comfortable"), h.Class("w-fit"),
		item("default", "r1", "Default"),
		item("comfortable", "r2", "Comfortable"),
		item("compact", "r3", "Compact"),
	)
}

func NovaChoiceCard() h.Node {
	card := func(value, id, title, desc string) h.Node {
		return field.Label(h.For(id),
			field.Field(field.Horizontal,
				field.Content(field.Title(title), field.Description(desc)),
				rg.Item(rg.Value(value), h.ID(id)),
			),
		)
	}
	return rg.RadioGroup(rg.Value("plus"), h.Class("max-w-sm"),
		card("plus", "plus-plan", "Plus", "For individuals and small teams."),
		card("pro", "pro-plan", "Pro", "For growing businesses."),
		card("enterprise", "enterprise-plan", "Enterprise", "For large teams and enterprises."),
	)
}

// described renders the radio-group-description items, also used (with
// dir) by the RTL preview.
func described(dir string, items [][4]string, attrs ...any) h.Node {
	var kids []any
	for _, it := range items {
		var d []any
		if dir != "" {
			d = []any{h.Dir(dir)}
		}
		kids = append(kids, field.Field(field.Horizontal,
			rg.Item(append([]any{rg.Value(it[0]), h.ID(it[1])}, d...)...),
			field.Content(
				field.Label(append([]any{h.For(it[1])}, append(d, it[2])...)...),
				field.Description(append(d, it[3])...),
			),
		))
	}
	return rg.RadioGroup(append(append([]any{rg.Value("comfortable"), h.Class("w-fit")}, attrs...), kids...)...)
}

func NovaDescription() h.Node {
	return described("", [][4]string{
		{"default", "desc-r1", "Default", "Standard spacing for most use cases."},
		{"comfortable", "desc-r2", "Comfortable", "More space between elements."},
		{"compact", "desc-r3", "Compact", "Minimal spacing for dense layouts."},
	})
}

func NovaRtl() h.Node {
	return described("rtl", [][4]string{
		{"default", "r1-rtl", "افتراضي", "تباعد قياسي لمعظم حالات الاستخدام."},
		{"comfortable", "r2-rtl", "مريح", "مساحة أكبر بين العناصر."},
		{"compact", "r3-rtl", "مضغوط", "تباعد أدنى للتخطيطات الكثيفة."},
	}, h.Dir("rtl"))
}

func NovaDisabled() h.Node {
	return rg.RadioGroup(rg.Value("option2"), h.Class("w-fit"),
		field.Field(field.Horizontal, h.Data("disabled", "true"),
			rg.Item(rg.Value("option1"), h.ID("disabled-1"), h.Disabled()),
			field.Label(h.For("disabled-1"), h.Class("font-normal"), "Disabled"),
		),
		field.Field(field.Horizontal,
			rg.Item(rg.Value("option2"), h.ID("disabled-2")),
			field.Label(h.For("disabled-2"), h.Class("font-normal"), "Option 2"),
		),
		field.Field(field.Horizontal,
			rg.Item(rg.Value("option3"), h.ID("disabled-3")),
			field.Label(h.For("disabled-3"), h.Class("font-normal"), "Option 3"),
		),
	)
}

func plainItems(invalid bool, items [][3]string) []any {
	var out []any
	for _, it := range items {
		f := []any{field.Horizontal}
		item := []any{rg.Value(it[0]), h.ID(it[1])}
		if invalid {
			f = append(f, h.Data("invalid", "true"))
			item = append(item, h.Aria("invalid", "true"))
		}
		out = append(out, field.Field(append(f,
			rg.Item(item...),
			field.Label(h.For(it[1]), h.Class("font-normal"), it[2]),
		)...))
	}
	return out
}

func NovaFieldset() h.Node {
	return field.Set(h.Class("w-full max-w-xs"),
		field.Legend(field.LegendLabel, "Subscription Plan"),
		field.Description("Yearly and lifetime plans offer significant savings."),
		rg.RadioGroup(append([]any{rg.Value("monthly")}, plainItems(false, [][3]string{
			{"monthly", "plan-monthly", "Monthly ($9.99/month)"},
			{"yearly", "plan-yearly", "Yearly ($99.99/year)"},
			{"lifetime", "plan-lifetime", "Lifetime ($299.99)"},
		})...)...),
	)
}

func NovaInvalid() h.Node {
	return field.Set(h.Class("w-full max-w-xs"),
		field.Legend(field.LegendLabel, "Notification Preferences"),
		field.Description("Choose how you want to receive notifications."),
		rg.RadioGroup(append([]any{rg.Value("email")}, plainItems(true, [][3]string{
			{"email", "invalid-email", "Email only"},
			{"sms", "invalid-sms", "SMS only"},
			{"both", "invalid-both", "Both Email & SMS"},
		})...)...),
	)
}
