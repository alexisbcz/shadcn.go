package form

import (
	"strconv"
	"unicode/utf16"

	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/spinner"
	h "shadcn.go/libhtml"
)

// NextDemoValues are form-next-demo's fields.
type NextDemoValues struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// NextDemoState is form-next-demo's FormState: the submitted values, the
// errors per field (zod's flatten().fieldErrors) and whether the action
// succeeded.
type NextDemoState struct {
	Values  NextDemoValues
	Errors  map[string][]string
	Success bool
}

// jsLen is a string's length in UTF-16 code units, like JavaScript's.
func jsLen(s string) int { return len(utf16.Encode([]rune(s))) }

// NextDemoAction is upstream's demoFormAction server action, run on the
// server: it validates the values with the example's schema.
func NextDemoAction(values NextDemoValues) NextDemoState {
	errs := map[string][]string{}
	check := func(field, v string, min, max int, minMsg, maxMsg string) {
		if jsLen(v) < min {
			errs[field] = append(errs[field], minMsg)
		}
		if jsLen(v) > max {
			errs[field] = append(errs[field], maxMsg)
		}
	}
	check("title", values.Title, 5, 32, "Bug title must be at least 5 characters.", "Bug title must be at most 32 characters.")
	check("description", values.Description, 20, 100, "Description must be at least 20 characters.", "Description must be at most 100 characters.")
	if len(errs) > 0 {
		return NextDemoState{Values: values, Errors: errs}
	}
	return NextDemoState{Success: true}
}

// NextDemo is form-next-demo: a Next.js form posting to a server action
// (React's useActionState), rendered for a state. pending renders the
// submission in flight. Upstream's reference build replaces next/form with
// a plain <form> without the action, so in the parity harness the form
// submits natively and the page reloads in its initial state; the Go page
// does the same.
func NextDemo(state NextDemoState, pending bool) h.Node {
	const name = "form-next-demo"
	invalid := func(f string) string { return strconv.FormatBool(len(state.Errors[f]) > 0) }
	errorOf := func(f string) h.Node {
		if len(state.Errors[f]) == 0 {
			return nil
		}
		return field.Error(state.Errors[f][0])
	}
	title := input.Input(
		h.ID("title"),
		h.Name("title"),
		h.Value(state.Values.Title),
		h.Aria("invalid", invalid("title")),
		h.Placeholder("Login button not working on mobile"),
		h.AutoComplete("off"),
		h.If(pending, h.Bool("disabled")),
	)
	desc := inputgroup.Textarea(
		h.ID("description"),
		h.Name("description"),
		h.Placeholder("I'm having an issue with the login button on mobile."),
		h.Attr("rows", "6"),
		h.Class("min-h-24 resize-none"),
		h.Aria("invalid", invalid("description")),
		state.Values.Description,
	)
	submit := button.Button(h.Type("submit"), h.FormAttr("bug-report-form"))
	if pending {
		title.Props().Bool("disabled", true)
		desc.Props().Bool("disabled", true)
		submit.Disabled(true)
		submit.Children(spinner.Spinner())
	}
	submit.Children("Submit")
	return with(name, card.Card(h.Class("w-full max-w-md"),
		header("Bug Report", "Help us improve by reporting bugs you encounter."),
		card.Content(
			h.Form(h.ID("bug-report-form"),
				field.Group(
					field.Field(h.Data("invalid", invalid("title")),
						field.Label(h.For("title"), "Bug Title"),
						title,
						errorOf("title"),
					),
					field.Field(h.Data("invalid", invalid("description")),
						field.Label(h.For("description"), "Description"),
						inputgroup.InputGroup(
							desc,
							inputgroup.Addon(inputgroup.AlignBlockEnd,
								inputgroup.Text(h.Class("tabular-nums"), strconv.Itoa(jsLen(state.Values.Description)), h.Comment(""), "/100 characters"),
							),
						),
						field.Description("Include steps to reproduce, expected behavior, and what actually happened."),
						errorOf("description"),
					),
				),
			),
		),
		card.Footer(field.Field(field.Horizontal, submit)),
	))
}

// NextComplexValues are form-next-complex's fields.
type NextComplexValues = complexValues

// NextComplexState is form-next-complex's FormState.
type NextComplexState struct {
	Values  NextComplexValues
	Errors  map[string][]string
	Success bool
}

// NextComplexDefaults is form-next-complex's initial state.
var NextComplexDefaults = NextComplexState{Values: NextComplexValues{Plan: "basic", BillingPeriod: "monthly", Addons: []string{}}}

// NextComplexAction is upstream's complexFormAction server action, run on
// the server: it validates the values with the example's schema (every
// message per field, like zod's flatten().fieldErrors).
func NextComplexAction(values NextComplexValues) NextComplexState {
	errs := map[string][]string{}
	add := func(f, msg string) { errs[f] = append(errs[f], msg) }
	if jsLen(values.Plan) < 1 {
		add("plan", "Please select a subscription plan")
	}
	if values.Plan != "basic" && values.Plan != "pro" {
		add("plan", "Invalid plan selection. Please choose Basic or Pro")
	}
	if jsLen(values.BillingPeriod) < 1 {
		add("billingPeriod", "Please select a billing period")
	}
	if len(values.Addons) < 1 {
		add("addons", "Please select at least one add-on")
	}
	if len(values.Addons) > 3 {
		add("addons", "You can select up to 3 add-ons")
	}
	for _, a := range values.Addons {
		known := false
		for _, x := range addons {
			if x.id == a {
				known = true
			}
		}
		if !known {
			add("addons", "You selected an invalid add-on")
			break
		}
	}
	if len(errs) > 0 {
		return NextComplexState{Values: values, Errors: errs}
	}
	return NextComplexState{Values: values, Success: true}
}

// NextComplex is form-next-complex rendered for a state (see NextDemo for
// how the parity page submits).
func NextComplex(state NextComplexState, pending bool) h.Node {
	o := complexOpts{next: true, formID: "subscription-form", withName: true,
		idOf:    func(s string) string { return s },
		values:  state.Values,
		errors:  state.Errors,
		pending: pending,
	}
	submit := button.Button(h.Type("submit"), h.FormAttr(o.formID))
	if pending {
		submit.Disabled(true)
		submit.Children(spinner.Spinner())
	}
	submit.Children("Save Preferences")
	return with("form-next-complex", card.Card(h.Class("w-full max-w-sm"),
		card.Content(h.Form(h.ID(o.formID), subscriptionForm(o))),
		card.Footer(field.Field(field.Horizontal, h.Class("justify-end"), submit)),
	))
}
