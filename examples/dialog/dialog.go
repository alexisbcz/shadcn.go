// Package dialog ports the upstream dialog-* examples.
package dialog

import (
	"fmt"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "dialog", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("dialog-close-button", CloseButton),
		ex("dialog-demo", Demo),
		ex("dialog-non-modal", NonModal),
		ex("dialog-nested", Nested),
		ex("dialog-scrollable", Scrollable),
		ex("dialog-default-open", DefaultOpen),
	)
}

func CloseButton() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Share")),
		dialog.Content(h.Class("sm:max-w-md"),
			dialog.Header(
				dialog.Title("Share link"),
				dialog.Description("Anyone who has this link will be able to view this."),
			),
			h.Div(h.Class("flex items-center gap-2"),
				h.Div(h.Class("grid flex-1 gap-2"),
					label.Label(h.For("link"), h.Class("sr-only"), "Link"),
					input.Input(h.ID("link"), h.Value("https://ui.shadcn.com/docs/installation"), h.ReadOnly()),
				),
			),
			dialog.Footer(h.Class("sm:justify-start"),
				dialog.Close(lib.AsChild, button.Button(h.Type("button"), button.VariantSecondary, "Close")),
			),
		),
	)
}

func Demo() h.Node {
	return dialog.Dialog(
		h.Form(
			dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open Dialog")),
			dialog.Content(h.Class("sm:max-w-[425px]"),
				dialog.Header(
					dialog.Title("Edit profile"),
					dialog.Description("Make changes to your profile here. Click save when you're done."),
				),
				h.Div(h.Class("grid gap-4"),
					h.Div(h.Class("grid gap-3"),
						label.Label(h.For("name-1"), "Name"),
						input.Input(h.ID("name-1"), h.Name("name"), h.Value("Pedro Duarte")),
					),
					h.Div(h.Class("grid gap-3"),
						label.Label(h.For("username-1"), "Username"),
						input.Input(h.ID("username-1"), h.Name("username"), h.Value("@peduarte")),
					),
				),
				dialog.Footer(
					dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
					button.Button(h.Type("submit"), "Save changes"),
				),
			),
		),
	)
}

// NonModal is the parity fixture dialog-non-modal.
func NonModal() h.Node {
	return h.Div(h.Class("flex items-center gap-2"),
		dialog.Dialog(dialog.Modal(false),
			dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Non-modal")),
			dialog.Content(dialog.ShowCloseButton(false),
				dialog.Header(
					dialog.Title("Non-modal dialog"),
					dialog.Description("The page behind stays interactive."),
				),
				dialog.Footer(dialog.ShowCloseButton(true),
					button.Button("Save"),
				),
			),
		),
		button.Button(button.VariantSecondary, "Outside"),
	)
}

// Nested is the parity fixture dialog-nested.
func Nested() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open outer")),
		dialog.Content(
			dialog.Header(
				dialog.Title("Outer dialog"),
				dialog.Description("Open another dialog on top of this one."),
			),
			dialog.Dialog(
				dialog.Trigger(lib.AsChild, button.Button("Open inner")),
				dialog.Content(h.Class("sm:max-w-sm"),
					dialog.Header(
						dialog.Title("Inner dialog"),
						dialog.Description("Escape closes this one first."),
					),
					dialog.Footer(dialog.ShowCloseButton(true)),
				),
			),
		),
	)
}

// Scrollable is the parity fixture dialog-scrollable.
func Scrollable() h.Node {
	return h.Div(h.Class("flex h-[1400px] flex-col items-center gap-4 pt-10"),
		dialog.Dialog(
			dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Terms")),
			dialog.Content(h.Class("sm:max-w-lg"),
				dialog.Header(
					dialog.Title("Terms of service"),
					dialog.Description("Scroll to read everything."),
				),
				h.Div(h.Class("-mx-6 max-h-[240px] overflow-y-auto px-6 text-sm"),
					h.Range(12, func(i int) h.Node {
						return h.P(h.Class("mb-4 leading-normal"),
							fmt.Sprintf("Section %d. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", i+1))
					}),
				),
				dialog.Footer(
					dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Decline")),
					button.Button("Accept"),
				),
			),
		),
		h.P(h.Class("text-sm text-muted-foreground"), "The page scrolls."),
	)
}

// DefaultOpen is the parity fixture dialog-default-open.
func DefaultOpen() h.Node {
	return dialog.Dialog(dialog.DefaultOpen, dialog.Modal(false),
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open again")),
		dialog.Content(h.Class("sm:max-w-sm"),
			dialog.Header(
				dialog.Title("Welcome"),
				dialog.Description("This dialog is open when the page loads."),
			),
			dialog.Footer(
				dialog.Close(lib.AsChild, button.Button("Got it")),
			),
		),
	)
}
