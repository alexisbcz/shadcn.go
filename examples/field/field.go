// Package field ports the upstream field-* examples.
package field

import (
	_ "embed"

	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	rg "shadcn.go/components/ui/radio-group"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/slider"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// slider.js keeps field-slider's description in sync with the slider.
//
//go:embed slider.js
var sliderScript string

func init() {
	runtime.Module("example-field-slider", sliderScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "field", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("field-choice-card", ChoiceCard),
		ex("field-fieldset", Fieldset),
		ex("field-input", Input),
		ex("field-radio", Radio),
		ex("field-responsive", Responsive),
		ex("field-textarea", Textarea),
		ex("field-checkbox", Checkbox),
		ex("field-group", Group),
		ex("field-switch", Switch),
		ex("field-slider", Slider),
	)
}

func ChoiceCard() h.Node {
	card := func(id, value, title, desc string) h.Node {
		return field.Label(h.For(id),
			field.Field(field.Horizontal,
				field.Content(field.Title(title), field.Description(desc)),
				rg.Item(rg.Value(value), h.ID(id)),
			),
		)
	}
	return h.Div(h.Class("w-full max-w-md"),
		field.Group(
			field.Set(
				field.Label(h.For("compute-environment-p8w"), "Compute Environment"),
				field.Description("Select the compute environment for your cluster."),
				rg.RadioGroup(rg.Value("kubernetes"),
					card("kubernetes-r2h", "kubernetes", "Kubernetes", "Run GPU workloads on a K8s configured cluster."),
					card("vm-z4k", "vm", "Virtual Machine", "Access a VM configured cluster to run GPU workloads."),
				),
			),
		),
	)
}

func Fieldset() h.Node {
	f := func(id, label, placeholder string) h.Node {
		return field.Field(
			field.Label(h.For(id), label),
			input.Input(h.ID(id), h.Type("text"), h.Placeholder(placeholder)),
		)
	}
	return h.Div(h.Class("w-full max-w-md space-y-6"),
		field.Set(
			field.Legend("Address Information"),
			field.Description("We need your address to deliver your order."),
			field.Group(
				f("street", "Street Address", "123 Main St"),
				h.Div(h.Class("grid grid-cols-2 gap-4"),
					f("city", "City", "New York"),
					f("zip", "Postal Code", "90502"),
				),
			),
		),
	)
}

func Input() h.Node {
	return h.Div(h.Class("w-full max-w-md"),
		field.Set(field.Group(
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
		)),
	)
}

func Radio() h.Node {
	opt := func(value, id, label string) h.Node {
		return field.Field(field.Horizontal,
			rg.Item(rg.Value(value), h.ID(id)),
			field.Label(h.For(id), h.Class("font-normal"), label),
		)
	}
	return h.Div(h.Class("w-full max-w-md"),
		field.Set(
			field.Label("Subscription Plan"),
			field.Description("Yearly and lifetime plans offer significant savings."),
			rg.RadioGroup(rg.Value("monthly"),
				opt("monthly", "plan-monthly", "Monthly ($9.99/month)"),
				opt("yearly", "plan-yearly", "Yearly ($99.99/year)"),
				opt("lifetime", "plan-lifetime", "Lifetime ($299.99)"),
			),
		),
	)
}

func Responsive() h.Node {
	return h.Div(h.Class("w-full max-w-4xl"),
		h.Form(
			field.Set(
				field.Legend("Profile"),
				field.Description("Fill in your profile information."),
				field.Separator(),
				field.Group(
					field.Field(field.Responsive,
						field.Content(
							field.Label(h.For("name"), "Name"),
							field.Description("Provide your full name for identification"),
						),
						input.Input(h.ID("name"), h.Placeholder("Evil Rabbit"), h.Required()),
					),
					field.Separator(),
					field.Field(field.Responsive,
						field.Content(
							field.Label(h.For("lastName"), "Message"),
							field.Description("You can write your message here. Keep it short, preferably under 100 characters."),
						),
						textarea.Textarea(h.ID("message"), h.Placeholder("Hello, world!"), h.Required(), h.Class("min-h-[100px] resize-none sm:min-w-[300px]")),
					),
					field.Separator(),
					field.Field(field.Responsive,
						button.Button(h.Type("submit"), "Submit"),
						button.Button(h.Type("button"), button.VariantOutline, "Cancel"),
					),
				),
			),
		),
	)
}

func Textarea() h.Node {
	return h.Div(h.Class("w-full max-w-md"),
		field.Set(field.Group(
			field.Field(
				field.Label(h.For("feedback"), "Feedback"),
				textarea.Textarea(h.ID("feedback"), h.Placeholder("Your feedback helps us improve..."), h.Attr("rows", "4")),
				field.Description("Share your thoughts about our service."),
			),
		)),
	)
}

func Checkbox() h.Node {
	opt := func(id, label string, extra ...any) h.Node {
		return field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID(id)),
			field.Label(append([]any{h.For(id), h.Class("font-normal")}, append(extra, label)...)...),
		)
	}
	return h.Div(h.Class("w-full max-w-md"),
		field.Group(
			field.Set(
				field.Legend(field.LegendLabel, "Show these items on the desktop"),
				field.Description("Select the items you want to show on the desktop."),
				field.Group(h.Class("gap-3"),
					// Upstream passes defaultChecked to this <label>; React
					// does not render it on a label.
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
		),
	)
}

func Group() h.Node {
	opt := func(id, label string, extra ...any) h.Node {
		return field.Field(field.Horizontal,
			checkbox.Checkbox(append([]any{h.ID(id)}, extra...)...),
			field.Label(h.For(id), h.Class("font-normal"), label),
		)
	}
	return h.Div(h.Class("w-full max-w-md"),
		field.Group(
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
		),
	)
}

func Switch() h.Node {
	return h.Div(h.Class("w-full max-w-md"),
		field.Field(field.Horizontal,
			field.Content(
				field.Label(h.For("2fa"), "Multi-factor authentication"),
				field.Description("Enable multi-factor authentication. If you do not have a two-factor device, you can use a one-time code sent to your email."),
			),
			switchui.Switch(h.ID("2fa")),
		),
	)
}

func Slider() h.Node {
	return h.Div(h.Class("w-full max-w-md"),
		field.Field(
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
		),
	)
}
