// Package alertdialog ports the upstream alert-dialog-* examples.
package alertdialog

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "alert-dialog", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("alert-dialog-demo", Demo),
		ex("alert-dialog-small", Small),
	)
}

func Demo() h.Node {
	return alertdialog.AlertDialog(
		alertdialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Show Dialog")),
		alertdialog.Content(
			alertdialog.Header(
				alertdialog.Title("Are you absolutely sure?"),
				alertdialog.Description("This action cannot be undone. This will permanently delete your account and remove your data from our servers."),
			),
			alertdialog.Footer(
				alertdialog.Cancel("Cancel"),
				alertdialog.Action("Continue"),
			),
		),
	)
}

// Small is the parity fixture alert-dialog-small.
func Small() h.Node {
	return h.Div(h.Class("flex gap-2"),
		alertdialog.AlertDialog(
			alertdialog.Trigger(lib.AsChild, button.Button(button.VariantDestructive, "Delete (sm)")),
			alertdialog.Content(alertdialog.SizeSm,
				alertdialog.Header(
					alertdialog.Media(lucide.Trash2()),
					alertdialog.Title("Delete project?"),
					alertdialog.Description("This permanently deletes the project and its files."),
				),
				alertdialog.Footer(
					alertdialog.Cancel("Cancel"),
					alertdialog.Action(button.VariantDestructive, h.Class("rounded-full"), "Delete"),
				),
			),
		),
		alertdialog.AlertDialog(
			alertdialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Delete (media)")),
			alertdialog.Content(
				alertdialog.Header(
					alertdialog.Media(lucide.Trash2()),
					alertdialog.Title("Delete project?"),
					alertdialog.Description("This permanently deletes the project and its files."),
				),
				alertdialog.Footer(
					alertdialog.Cancel(button.SizeSm, "Cancel"),
					alertdialog.Action(button.VariantDestructive, button.SizeSm, "Delete"),
				),
			),
		),
	)
}
