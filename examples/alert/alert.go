// Package alert ports the upstream alert-* examples.
package alert

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/alert"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "alert", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("alert-demo", Demo),
		ex("alert-destructive", Destructive),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("grid w-full max-w-xl items-start gap-4"),
		alert.Alert(
			lucide.CircleCheck(),
			alert.Title("Success! Your changes have been saved"),
			alert.Description("This is an alert with icon, title and description."),
		),
		alert.Alert(
			lucide.Popcorn(),
			alert.Title("This Alert has a title and an icon. No description."),
		),
		alert.Alert(alert.VariantDestructive,
			lucide.CircleAlert(),
			alert.Title("Unable to process your payment."),
			alert.Description(
				h.P("Please verify your billing information and try again."),
				h.Ul(h.Class("list-inside list-disc text-sm"),
					h.Li("Check your card details"),
					h.Li("Ensure sufficient funds"),
					h.Li("Verify billing address"),
				),
			),
		),
	)
}

func Destructive() h.Node {
	return alert.Alert(alert.VariantDestructive,
		lucide.CircleAlert(),
		alert.Title("Error"),
		alert.Description("Your session has expired. Please log in again."),
	)
}
