package form

import (
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	inputgroup "shadcn.go/components/ui/input-group"
	h "shadcn.go/libhtml"
)

// emailName is the name attribute of the email at index.
func emailName(l Lib, index int) string {
	switch l {
	case tanstack:
		return "emails[" + strconv.Itoa(index) + "].address"
	case formisch:
		return fieldName(l, "emails", index, "address")
	}
	return "emails." + strconv.Itoa(index) + ".address"
}

func emailRemove(index int) *inputgroup.Builder {
	return inputgroup.Addon(inputgroup.AlignInlineEnd,
		inputgroup.Button(h.Type("button"), button.VariantGhost, inputgroup.ButtonSizeIconXs,
			h.Aria("label", "Remove email "+strconv.Itoa(index+1)),
			lucide.X(),
		),
	)
}

// emailItem is one email of the array examples.
func emailItem(l Lib, name string, index int, removable bool) *field.Builder {
	group := inputgroup.InputGroup(
		inputgroup.Input(
			h.Name(emailName(l, index)),
			h.Value(""),
			h.ID(name+"-email-"+strconv.Itoa(index)),
			ariaInvalid(),
			h.Placeholder("name@example.com"),
			h.Type("email"),
			h.AutoComplete("email"),
		),
	)
	if removable {
		group.Children(emailRemove(index))
	}
	return field.Field(field.Horizontal, dataInvalid(), field.Content(group))
}

// Array is form-{rhf,tanstack,formisch}-array.
func Array(l Lib) h.Node {
	name := "form-" + string(l) + "-array"
	n := 2
	if l == tanstack {
		n = 1
	}
	items := field.Group(h.Class("gap-4"))
	for i := 0; i < n; i++ {
		items.Children(emailItem(l, name, i, n > 1))
	}
	items.Children(button.Button(h.Type("button"), button.VariantOutline, button.SizeSm, "Add Email Address"))
	return h.Group(
		card.Card(h.Class("w-full sm:max-w-md"),
			header("Contact Emails", "Manage your contact email addresses.", h.Class("border-b")),
			card.Content(formTag(l, name,
				field.Set(h.Class("gap-4"),
					field.Legend(field.LegendLabel, "Email Addresses"),
					field.Description("Add up to 5 email addresses where we can contact you."),
					items,
				),
			)),
			card.Footer(h.Class("border-t"),
				field.Field(field.Horizontal,
					button.Button(h.Type("button"), button.VariantOutline, "Reset"),
					button.Button(h.Type("submit"), h.FormAttr(name), "Save"),
				),
			),
		),
		// The item and remove-button prototypes the client clones.
		h.Template(h.Data("parity-ignore", ""), h.Data("form-example", name),
			h.Div(h.Data("proto", "item"), emailItem(l, name, 0, false)),
			h.Div(h.Data("proto", "remove"), emailRemove(0)),
		),
	)
}
