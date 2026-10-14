package form

import (
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/field"
	radiogroup "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	switchui "shadcn.go/components/ui/switch"
	h "shadcn.go/libhtml"
)

// Switch is form-{rhf,tanstack,formisch}-switch.
func Switch(l Lib) h.Node {
	name := "form-" + string(l) + "-switch"
	sw := switchui.Switch(h.ID(name+"-twoFactor"), ariaInvalid())
	if l != formisch {
		sw.Name("twoFactor")
	}
	return with(name, card.Card(h.Class("w-full sm:max-w-md"),
		header("Security Settings", "Manage your account security preferences."),
		card.Content(
			formTag(l, name,
				field.Group(
					field.Field(field.Horizontal, dataInvalid(),
						field.Content(
							field.Label(h.For(name+"-twoFactor"), "Multi-factor authentication"),
							field.Description("Enable multi-factor authentication to secure your account."),
						),
						sw,
					),
				),
			),
		),
		footer(name, "Save"),
	))
}

type plan struct{ id, title, description string }

var plans = []plan{
	{"starter", "Starter (100K tokens/month)", "For everyday use with basic features."},
	{"pro", "Pro (1M tokens/month)", "For advanced AI usage with more features."},
	{"enterprise", "Enterprise (Unlimited tokens)", "For large teams and heavy usage."},
}

// RadioGroup is form-{rhf,tanstack,formisch}-radiogroup.
func RadioGroup(l Lib) h.Node {
	name := "form-" + string(l) + "-radiogroup"
	group := radiogroup.RadioGroup(radiogroup.Value(""))
	if l != formisch {
		group.Name("plan")
	}
	if l != tanstack {
		group.Children(ariaInvalid())
	}
	for _, p := range plans {
		group.Children(field.Label(h.For(name+"-"+p.id),
			field.Field(field.Horizontal, dataInvalid(),
				field.Content(
					field.Title(p.title),
					field.Description(p.description),
				),
				radiogroup.Item(radiogroup.Value(p.id), h.ID(name+"-"+p.id), ariaInvalid()),
			),
		))
	}
	set := field.Set(
		field.Legend("Plan"),
		field.Description("You can upgrade or downgrade your plan at any time."),
		group,
	)
	if l != tanstack {
		set = field.Set(dataInvalid(),
			field.Legend("Plan"),
			field.Description("You can upgrade or downgrade your plan at any time."),
			group,
		)
	}
	return with(name, card.Card(h.Class("w-full sm:max-w-md"),
		header("Subscription Plan", "See pricing and features for each plan."),
		card.Content(formTag(l, name, field.Group(set))),
		footer(name, "Save"),
	))
}

var notificationTasks = []struct{ id, label string }{
	{"push", "Push notifications"},
	{"email", "Email notifications"},
}

// Checkbox is form-{rhf,tanstack,formisch}-checkbox.
func Checkbox(l Lib) h.Node {
	name := "form-" + string(l) + "-checkbox"
	responses := checkbox.Checkbox(h.ID(name+"-responses"), checkbox.Checked(true), checkbox.Disabled(true))
	if l != formisch {
		responses.Name("responses")
	}
	responsesField := field.Field(field.Horizontal,
		responses,
		field.Label(h.For(name+"-responses"), h.Class("font-normal"), "Push notifications"),
	)
	responsesSet := field.Set()
	if l == tanstack {
		responsesField.Children(dataInvalid())
	} else {
		responsesSet.Children(dataInvalid())
	}
	responsesSet.Children(
		field.Legend(field.LegendLabel, "Responses"),
		field.Description("Get notified for requests that take time, like research or image generation."),
		field.Group(h.Data("slot", "checkbox-group"), responsesField),
	)
	group := field.Group(h.Data("slot", "checkbox-group"))
	for _, t := range notificationTasks {
		cb := checkbox.Checkbox(h.ID(name+"-"+t.id), ariaInvalid())
		if l != formisch {
			cb.Name("tasks")
		}
		group.Children(field.Field(field.Horizontal, dataInvalid(),
			cb,
			field.Label(h.For(name+"-"+t.id), h.Class("font-normal"), t.label),
		))
	}
	return with(name, card.Card(h.Class("w-full sm:max-w-md"),
		header("Notifications", "Manage your notification preferences."),
		card.Content(formTag(l, name, field.Group(
			h.Div(responsesSet),
			field.Separator(),
			field.Group(field.Set(dataInvalid(),
				field.Legend(field.LegendLabel, "Tasks"),
				field.Description("Get notified when tasks you've created have updates."),
				group,
			)),
		))),
		footer(name, "Save"),
	))
}

var spokenLanguages = []struct{ label, value string }{
	{"English", "en"},
	{"Spanish", "es"},
	{"French", "fr"},
	{"German", "de"},
	{"Italian", "it"},
	{"Chinese", "zh"},
	{"Japanese", "ja"},
}

// Select is form-{rhf,tanstack,formisch}-select.
func Select(l Lib) h.Node {
	name := "form-" + string(l) + "-select"
	content := selectui.Content(selectui.PositionItemAligned,
		selectui.Item(selectui.ItemValue("auto"), "Auto"),
		selectui.Separator(),
	)
	for _, lang := range spokenLanguages {
		content.Children(selectui.Item(selectui.ItemValue(lang.value), lang.label))
	}
	sel := selectui.Select(selectui.DefaultValue(""),
		selectui.Trigger(h.ID(name+"-language"), ariaInvalid(), h.Class("min-w-[120px]"),
			selectui.Value(selectui.Placeholder("Select")),
		),
		content,
	)
	if l != formisch {
		sel.Name("language")
	}
	return with(name, card.Card(h.Class("w-full sm:max-w-lg"),
		header("Language Preferences", "Select your preferred spoken language."),
		card.Content(formTag(l, name, field.Group(
			field.Field(field.Responsive, dataInvalid(),
				field.Content(
					field.Label(h.For(name+"-language"), "Spoken Language"),
					field.Description("For best results, select the language you speak."),
				),
				sel,
			),
		))),
		footer(name, "Save"),
	))
}
