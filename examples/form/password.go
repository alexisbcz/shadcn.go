package form

import (
	"encoding/json"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/progress"
	h "shadcn.go/libhtml"
)

var passwordRequirements = []string{
	"At least 8 characters",
	"One lowercase letter",
	"One uppercase letter",
	"One number",
	"One special character",
}

// progressClass is the strength meter's class for a strength color, like
// cn() inside upstream's Progress.
func progressClass(color string) string {
	return lib.Cn("relative h-2 w-full overflow-hidden rounded-full bg-primary/20", color)
}

func iconClass(icon *h.Element) string { return icon.ClassName() }

// Password is form-rhf-password.
func Password() h.Node {
	const name = "form-rhf-password"
	list := h.Div(h.Class("space-y-1.5"))
	for _, r := range passwordRequirements {
		list.Children(h.Div(h.Class("flex items-center gap-2 text-sm"),
			lucide.Check(h.Class("size-4 text-muted-foreground")),
			h.Span(h.Class("text-muted-foreground"), r),
		))
	}
	// The class names the client switches between, computed here like
	// upstream's cn() and lucide's className merge.
	classes := map[string]string{
		"bg-neutral-200":   progressClass("bg-neutral-200"),
		"bg-red-500":       progressClass("bg-red-500"),
		"bg-yellow-500":    progressClass("bg-yellow-500"),
		"bg-green-500":     progressClass("bg-green-500"),
		"addonMet":         iconClass(lucide.Check(h.Class("text-green-500"))),
		"addonUnmet":       iconClass(lucide.Check(h.Class("text-muted-foreground"))),
		"requirementMet":   iconClass(lucide.Check(h.Class("size-4 text-green-500"))),
		"requirementUnmet": iconClass(lucide.Check(h.Class("size-4 text-muted-foreground"))),
	}
	payload, _ := json.Marshal(classes)
	return h.Group(
		card.Card(h.Class("w-full sm:max-w-md"),
			header("Create Password", "Choose a strong password to secure your account.", h.Class("border-b")),
			card.Content(formTag(rhf, name, field.Group(
				field.Field(dataInvalid(),
					field.Label(h.For(name+"-input"), "Password"),
					inputgroup.InputGroup(
						inputgroup.Input(
							h.Name("password"),
							h.Value(""),
							h.ID(name+"-input"),
							h.Type("password"),
							h.Placeholder("Enter your password"),
							ariaInvalid(),
							h.AutoComplete("new-password"),
						),
						inputgroup.Addon(inputgroup.AlignInlineEnd, lucide.Check(h.Class("text-muted-foreground"))),
					),
					h.Div(h.Class("space-y-2"),
						progress.Progress(progress.Value(0), h.Class("bg-neutral-200")),
						list,
					),
				),
			))),
			card.Footer(h.Class("border-t"),
				field.Field(
					button.Button(h.Type("submit"), h.FormAttr(name), "Create Password"),
					button.Button(h.Type("button"), button.VariantOutline, "Reset"),
				),
			),
		),
		h.Template(h.Data("parity-ignore", ""), h.Data("form-example", name), h.Data("classes", string(payload))),
	)
}
