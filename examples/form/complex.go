package form

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/field"
	radiogroup "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	switchui "shadcn.go/components/ui/switch"
	h "shadcn.go/libhtml"
)

var addons = []struct{ id, title, description string }{
	{"analytics", "Analytics", "Advanced analytics and reporting"},
	{"backup", "Backup", "Automated daily backups"},
	{"support", "Priority Support", "24/7 premium customer support"},
}

var complexPlans = []struct{ id, title, description string }{
	{"basic", "Basic", "For individuals and small teams"},
	{"pro", "Pro", "For businesses with higher demands"},
}

// complexValues are the subscription form's values.
type complexValues struct {
	Plan               string
	BillingPeriod      string
	Addons             []string
	EmailNotifications bool
}

// complexOpts describes the variants of the subscription form.
type complexOpts struct {
	lib      Lib
	next     bool // form-next-complex (uncontrolled, React's useActionState)
	formID   string
	idOf     func(string) string // element ids by field / option
	values   complexValues
	errors   map[string][]string
	pending  bool
	withName bool
}

func (o complexOpts) invalid(f string) h.Prop {
	if len(o.errors[f]) > 0 {
		return h.Data("invalid", "true")
	}
	return dataInvalid()
}

func (o complexOpts) ariaInvalid(f string) h.Prop {
	if len(o.errors[f]) > 0 {
		return h.Aria("invalid", "true")
	}
	return ariaInvalid()
}

func (o complexOpts) errorOf(f string) h.Node {
	if len(o.errors[f]) == 0 {
		return nil
	}
	if o.next {
		return field.Error(o.errors[f][0])
	}
	return field.Errors(o.errors[f])
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// subscriptionForm renders form-{rhf,tanstack,formisch,next}-complex's
// fields. TanStack's and Next's variants mark the plan's items invalid and
// not the set, and leave the select trigger without aria-invalid.
func subscriptionForm(o complexOpts) h.Node {
	tanstackLike := o.lib == tanstack || o.next
	group := radiogroup.RadioGroup(radiogroup.Value(o.values.Plan))
	if o.withName {
		group.Name("plan")
	}
	if !tanstackLike || o.next {
		group.Children(o.ariaInvalid("plan"))
	}
	if o.pending {
		group.Disabled(true)
	}
	for _, p := range complexPlans {
		f := field.Field(field.Horizontal,
			field.Content(field.Title(p.title), field.Description(p.description)),
		)
		item := radiogroup.Item(radiogroup.Value(p.id), h.ID(o.idOf(p.id)))
		if tanstackLike && !o.next {
			f.Children(o.invalid("plan"))
			item.Children(o.ariaInvalid("plan"))
		}
		f.Children(item)
		group.Children(field.Label(h.For(o.idOf(p.id)), f))
	}
	planSet := field.Set()
	if !tanstackLike || o.next {
		planSet.Children(o.invalid("plan"))
	}
	legend := field.Legend(field.LegendLabel, "Subscription Plan")
	if tanstackLike {
		legend = field.Legend("Subscription Plan")
	}
	planSet.Children(legend, field.Description("Choose your subscription plan."), group, o.errorOf("plan"))

	trigger := selectui.Trigger(h.ID(o.idOf("billingPeriod")), selectui.Value(selectui.Placeholder("Select")))
	if !tanstackLike {
		trigger.Children(o.ariaInvalid("billingPeriod"))
	}
	if o.pending {
		trigger.Disabled(true)
	}
	sel := selectui.Select(selectui.DefaultValue(o.values.BillingPeriod),
		trigger,
		selectui.Content(
			selectui.Item(selectui.ItemValue("monthly"), "Monthly"),
			selectui.Item(selectui.ItemValue("yearly"), "Yearly"),
		),
	)
	if o.withName {
		sel.Name("billingPeriod")
	}

	boxes := field.Group(h.Data("slot", "checkbox-group"))
	for _, a := range addons {
		cb := checkbox.Checkbox(h.ID(o.idOf(a.id)), o.ariaInvalid("addons"), checkbox.Checked(contains(o.values.Addons, a.id)))
		if o.withName {
			cb.Name("addons")
		}
		if o.next {
			cb.Value(a.id)
		}
		if o.pending {
			cb.Disabled(true)
		}
		boxes.Children(field.Field(field.Horizontal, o.invalid("addons"),
			cb,
			field.Content(
				field.Label(h.For(o.idOf(a.id)), a.title),
				field.Description(a.description),
			),
		))
	}

	sw := switchui.Switch(h.ID(o.idOf("emailNotifications")), o.ariaInvalid("emailNotifications"), switchui.Checked(o.values.EmailNotifications))
	if o.withName {
		sw.Name("emailNotifications")
	}
	if o.pending {
		sw.Disabled(true)
	}
	notifications := field.Field(field.Horizontal,
		field.Content(
			field.Label(h.For(o.idOf("emailNotifications")), "Email Notifications"),
			field.Description("Receive email updates about your subscription"),
		),
		sw,
		o.errorOf("emailNotifications"),
	)
	if !o.next {
		notifications.Children(o.invalid("emailNotifications"))
	}

	return field.Group(
		planSet,
		field.Separator(),
		field.Field(o.invalid("billingPeriod"),
			field.Label(h.For(o.idOf("billingPeriod")), "Billing Period"),
			sel,
			field.Description("Choose how often you want to be billed."),
			o.errorOf("billingPeriod"),
		),
		field.Separator(),
		field.Set(
			field.Legend("Add-ons"),
			field.Description("Select additional features you'd like to include."),
			boxes,
			o.errorOf("addons"),
		),
		field.Separator(),
		notifications,
	)
}

// Complex is form-{rhf,tanstack,formisch}-complex.
func Complex(l Lib) h.Node {
	name := "form-" + string(l) + "-complex"
	o := complexOpts{lib: l, formID: name, withName: l != formisch,
		idOf:   func(s string) string { return name + "-" + s },
		values: complexValues{Plan: "basic"},
	}
	if l == tanstack {
		o.formID = "subscription-form"
		o.idOf = func(s string) string { return s }
		o.values.BillingPeriod = "monthly"
		return with(name, card.Card(h.Class("w-full max-w-sm"),
			card.Content(formTag(l, o.formID, subscriptionForm(o))),
			card.Footer(field.Field(field.Horizontal, h.Class("justify-end"),
				button.Button(h.Type("submit"), h.FormAttr(o.formID), "Save Preferences"),
			)),
		))
	}
	return with(name, card.Card(h.Class("w-full max-w-sm"),
		header("You're almost there!", "Choose your subscription plan and billing period.", h.Class("border-b")),
		card.Content(formTag(l, o.formID, subscriptionForm(o))),
		card.Footer(h.Class("border-t"),
			field.Field(
				button.Button(h.Type("submit"), h.FormAttr(o.formID), "Save Preferences"),
				button.Button(h.Type("button"), button.VariantOutline, "Reset"),
			),
		),
	))
}
