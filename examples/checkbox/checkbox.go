// Package checkbox ports the upstream checkbox-* examples.
package checkbox

import (
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "checkbox", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("checkbox-demo", Demo),
		ex("checkbox-disabled", Disabled),
		ex("checkbox-with-text", WithText),
		ex("checkbox-states", States),
		ex("checkbox-form", Form),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex flex-col gap-6"),
		h.Div(h.Class("flex items-center gap-3"),
			checkbox.Checkbox(h.ID("terms")),
			label.Label(h.For("terms"), "Accept terms and conditions"),
		),
		h.Div(h.Class("flex items-start gap-3"),
			checkbox.Checkbox(h.ID("terms-2"), checkbox.Checked(true)),
			h.Div(h.Class("grid gap-2"),
				label.Label(h.For("terms-2"), "Accept terms and conditions"),
				h.P(h.Class("text-sm text-muted-foreground"), "By clicking this checkbox, you agree to the terms and conditions."),
			),
		),
		h.Div(h.Class("flex items-start gap-3"),
			checkbox.Checkbox(h.ID("toggle"), h.Disabled()),
			label.Label(h.For("toggle"), "Enable notifications"),
		),
		label.Label(h.Class("flex items-start gap-3 rounded-lg border p-3 hover:bg-accent/50 has-[[aria-checked=true]]:border-blue-600 has-[[aria-checked=true]]:bg-blue-50 dark:has-[[aria-checked=true]]:border-blue-900 dark:has-[[aria-checked=true]]:bg-blue-950"),
			checkbox.Checkbox(h.ID("toggle-2"), checkbox.Checked(true),
				h.Class("data-[state=checked]:border-blue-600 data-[state=checked]:bg-blue-600 data-[state=checked]:text-white dark:data-[state=checked]:border-blue-700 dark:data-[state=checked]:bg-blue-700"),
			),
			h.Div(h.Class("grid gap-1.5 font-normal"),
				h.P(h.Class("text-sm leading-none font-medium"), "Enable notifications"),
				h.P(h.Class("text-sm text-muted-foreground"), "You can enable or disable notifications at any time."),
			),
		),
	)
}

func Disabled() h.Node {
	return h.Div(h.Class("flex items-center space-x-2"),
		checkbox.Checkbox(h.ID("terms2"), h.Disabled()),
		h.Label(h.For("terms2"), h.Class("text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"), "Accept terms and conditions"),
	)
}

func WithText() h.Node {
	return h.Div(h.Class("items-top flex gap-2"),
		checkbox.Checkbox(h.ID("terms1")),
		h.Div(h.Class("grid gap-1.5 leading-none"),
			h.Label(h.For("terms1"), h.Class("text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"), "Accept terms and conditions"),
			h.P(h.Class("text-sm text-muted-foreground"), "You agree to our Terms of Service and Privacy Policy."),
		),
	)
}

// States ports the checkbox-states fixture (parity/fixtures/checkbox-states.tsx):
// the Radix Checkbox props the upstream examples do not use: an
// indeterminate checkbox, checkboxes in a form (hidden bubble inputs) with
// name, value, required and a reset button, a checked disabled checkbox
// and an aria-invalid one.
func States() h.Node {
	row := func(id, text string, cb *checkbox.Builder) h.Node {
		return h.Div(h.Class("flex items-center gap-3"), cb.ID(id), label.Label(h.For(id), text))
	}
	return h.Div(h.Class("flex flex-col gap-6"),
		row("mixed", "Select all", checkbox.Checkbox(checkbox.Indeterminate)),
		h.Form(h.Class("flex flex-col gap-3"),
			row("news", "Newsletter", checkbox.Checkbox(checkbox.Name("news"), checkbox.Value("yes"), checkbox.Required(true), checkbox.Checked(true))),
			row("partial", "Partial", checkbox.Checkbox(checkbox.Name("partial"), checkbox.Indeterminate)),
			row("offers", "Offers", checkbox.Checkbox(checkbox.Name("offers"))),
			h.Button(h.Type("reset"), h.Class("w-fit rounded-md border px-2 text-sm"), "Reset"),
		),
		row("locked", "Locked", checkbox.Checkbox(checkbox.Checked(true), h.Disabled())),
		row("invalid", "Invalid", checkbox.Checkbox(h.Aria("invalid", "true"))),
	)
}

// Form ports the checkbox-form fixture (parity/fixtures/checkbox-form.tsx):
// checkboxes tied to a form through the form attribute (outside the form),
// and a disabled checkbox inside a form.
func Form() h.Node {
	row := func(id, text string, cb *checkbox.Builder) h.Node {
		return h.Div(h.Class("flex items-center gap-3"), cb.ID(id), label.Label(h.For(id), text))
	}
	return h.Div(h.Class("flex flex-col gap-4"),
		h.Form(h.ID("checkbox-form"), h.Class("flex gap-3"),
			h.Button(h.Type("reset"), h.Class("w-fit rounded-md border px-2 text-sm"), "Reset"),
		),
		row("remote", "Remote", checkbox.Checkbox(h.FormAttr("checkbox-form"), h.Name("remote"), h.Value("yes"), checkbox.Checked(true))),
		row("remote-mixed", "Remote mixed", checkbox.Checkbox(checkbox.Form("checkbox-form"), checkbox.Name("remote-mixed"), checkbox.Indeterminate, h.Required())),
		h.Form(h.Class("flex items-center gap-3"),
			checkbox.Checkbox(checkbox.Name("locked"), h.Disabled(), checkbox.Checked(true), h.ID("locked")),
			label.Label(h.For("locked"), "Locked"),
		),
	)
}
