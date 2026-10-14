package form

import (
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/textarea"
	h "shadcn.go/libhtml"
)

// Input is form-{rhf,tanstack,formisch}-input.
func Input(l Lib) h.Node {
	name := "form-" + string(l) + "-input"
	return with(name, card.Card(h.Class("w-full sm:max-w-md"),
		header("Profile Settings", "Update your profile information below."),
		card.Content(
			formTag(l, name,
				field.Group(
					field.Field(dataInvalid(),
						field.Label(h.For(name+"-username"), "Username"),
						input.Input(
							h.Name(fieldName(l, "username")),
							h.Value(""),
							h.ID(name+"-username"),
							ariaInvalid(),
							h.Placeholder("shadcn"),
							h.AutoComplete("username"),
						),
						field.Description("This is your public display name. Must be between 3 and 10 characters. Must only contain letters, numbers, and underscores."),
					),
				),
			),
		),
		footer(name, "Save"),
	))
}

// Textarea is form-{rhf,tanstack,formisch}-textarea.
func Textarea(l Lib) h.Node {
	name := "form-" + string(l) + "-textarea"
	return with(name, card.Card(h.Class("w-full sm:max-w-md"),
		header("Personalization", "Customize your experience by telling us more about yourself."),
		card.Content(
			formTag(l, name,
				field.Group(
					field.Field(dataInvalid(),
						field.Label(h.For(name+"-about"), "More about you"),
						textarea.Textarea(
							h.Name(fieldName(l, "about")),
							h.ID(name+"-about"),
							ariaInvalid(),
							h.Placeholder("I'm a software engineer..."),
							h.Class("min-h-[120px]"),
						),
						field.Description("Tell us more about yourself. This will be used to help us personalize your experience."),
					),
				),
			),
		),
		footer(name, "Save"),
	))
}

// Demo is form-{rhf,tanstack,formisch}-demo, the bug report form.
func Demo(l Lib) h.Node {
	name := "form-" + string(l) + "-demo"
	formID := name
	titleID, descID := name+"-title", name+"-description"
	if l == tanstack {
		formID = "bug-report-form"
		titleID, descID = "title", "description"
	}
	return with(name, card.Card(h.Class("w-full sm:max-w-md"),
		header("Bug Report", "Help us improve by reporting bugs you encounter."),
		card.Content(
			formTag(l, formID,
				field.Group(
					field.Field(dataInvalid(),
						field.Label(h.For(titleID), "Bug Title"),
						input.Input(
							h.Name(fieldName(l, "title")),
							h.Value(""),
							h.ID(titleID),
							ariaInvalid(),
							h.Placeholder("Login button not working on mobile"),
							h.AutoComplete("off"),
						),
					),
					field.Field(dataInvalid(),
						field.Label(h.For(descID), "Description"),
						inputgroup.InputGroup(
							inputgroup.Textarea(
								h.Name(fieldName(l, "description")),
								h.ID(descID),
								h.Placeholder("I'm having an issue with the login button on mobile."),
								h.Attr("rows", "6"),
								h.Class("min-h-24 resize-none"),
								ariaInvalid(),
							),
							inputgroup.Addon(inputgroup.AlignBlockEnd,
								inputgroup.Text(h.Class("tabular-nums"), "0", h.Comment(""), "/100 characters"),
							),
						),
						field.Description("Include steps to reproduce, expected behavior, and what actually happened."),
					),
				),
			),
		),
		footer(formID, "Submit"),
	))
}
