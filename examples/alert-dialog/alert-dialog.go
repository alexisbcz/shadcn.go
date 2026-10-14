// Package alertdialog ports the upstream alert-dialog-* examples.
package alertdialog

import (
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "alert-dialog", Render: withThemeScript(fn)}
	}
	registry.Register(
		ex("alert-dialog-demo", Demo),
		ex("alert-dialog-small", Small),
	)
}

// withThemeScript puts an empty nonced <script> first in the preview, where
// the upstream docs' next-themes ThemeProvider renders its theme script.
// It matters to modal dialogs: hideOthers (aria-hidden) keeps <script>
// elements and their ancestors visible, so the preview itself is kept and
// its other children are hidden, exactly as upstream.
func withThemeScript(fn func() h.Node) func(*http.Request) h.Node {
	return func(r *http.Request) h.Node { return h.Group(secure.Script(r), fn()) }
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
