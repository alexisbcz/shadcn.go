package form

import (
	"shadcn.go/components/ui/button"
	uiform "shadcn.go/components/ui/form"
	"shadcn.go/components/ui/input"
	h "shadcn.go/libhtml"
)

// HookDemo is the form-hook-demo parity fixture: the Form component
// (form.tsx) validated with react-hook-form's semantics.
func HookDemo() h.Node {
	return with("form-hook-demo", uiform.Form(
		h.Form(h.Class("w-2/3 space-y-6"),
			uiform.Field("username", uiform.Item(
				uiform.Label("Username"),
				uiform.Control(input.Input(h.Placeholder("shadcn"), h.Name("username"), h.Value(""))),
				uiform.Description("This is your public display name."),
				uiform.Message(),
			)),
			uiform.Field("bio", uiform.Item(
				uiform.Label("Bio"),
				uiform.Control(input.Input(h.Placeholder("Tell us about you"), h.Name("bio"), h.Value(""))),
				uiform.Message("Optional, up to 20 characters."),
			)),
			button.Button(h.Type("submit"), "Submit"),
		),
	))
}
