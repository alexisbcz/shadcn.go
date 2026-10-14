// Package form ports the upstream form-* examples: forms built with Field
// and validated on the client with react-hook-form, TanStack Form or
// Formisch (zod / valibot schemas), and the Next.js server-action forms.
//
// The markup is the examples' first render; the behavior (the libraries'
// validation semantics, the conditional error messages, toasts) is in the
// examples' JavaScript modules, built on components/ui/form's
// form-state.js. Each example renders an inert marker template (left out
// of the parity DOM comparison) naming the example, next to its card.
package form

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	uiform "shadcn.go/components/ui/form"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/sonner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed examples.js
var examplesScript string

func init() {
	runtime.Module("example-form", examplesScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "form", Render: registry.Static(fn), Extra: toaster}
	}
	registry.Register(
		ex("form-rhf-input", func() h.Node { return Input(rhf) }),
		ex("form-tanstack-input", func() h.Node { return Input(tanstack) }),
		ex("form-formisch-input", func() h.Node { return Input(formisch) }),
		ex("form-rhf-textarea", func() h.Node { return Textarea(rhf) }),
		ex("form-tanstack-textarea", func() h.Node { return Textarea(tanstack) }),
		ex("form-formisch-textarea", func() h.Node { return Textarea(formisch) }),
		ex("form-rhf-demo", func() h.Node { return Demo(rhf) }),
		ex("form-tanstack-demo", func() h.Node { return Demo(tanstack) }),
		ex("form-formisch-demo", func() h.Node { return Demo(formisch) }),
		ex("form-rhf-password", Password),
		registry.Example{Name: "form-hook-demo", Component: "form", Render: registry.Static(HookDemo)},
		ex("form-next-complex", func() h.Node { return NextComplex(NextComplexDefaults, false) }),
		ex("form-next-demo", func() h.Node { return NextDemo(NextDemoState{}, false) }),
		ex("form-rhf-array", func() h.Node { return Array(rhf) }),
		ex("form-tanstack-array", func() h.Node { return Array(tanstack) }),
		ex("form-formisch-array", func() h.Node { return Array(formisch) }),
		ex("form-rhf-complex", func() h.Node { return Complex(rhf) }),
		ex("form-tanstack-complex", func() h.Node { return Complex(tanstack) }),
		ex("form-formisch-complex", func() h.Node { return Complex(formisch) }),
		ex("form-rhf-select", func() h.Node { return Select(rhf) }),
		ex("form-tanstack-select", func() h.Node { return Select(tanstack) }),
		ex("form-formisch-select", func() h.Node { return Select(formisch) }),
		ex("form-rhf-switch", func() h.Node { return Switch(rhf) }),
		ex("form-tanstack-switch", func() h.Node { return Switch(tanstack) }),
		ex("form-formisch-switch", func() h.Node { return Switch(formisch) }),
		ex("form-rhf-radiogroup", func() h.Node { return RadioGroup(rhf) }),
		ex("form-tanstack-radiogroup", func() h.Node { return RadioGroup(tanstack) }),
		ex("form-formisch-radiogroup", func() h.Node { return RadioGroup(formisch) }),
		ex("form-rhf-checkbox", func() h.Node { return Checkbox(rhf) }),
		ex("form-tanstack-checkbox", func() h.Node { return Checkbox(tanstack) }),
		ex("form-formisch-checkbox", func() h.Node { return Checkbox(formisch) }),
	)
}

// toaster is the <Toaster /> the upstream docs mount in their layout (see
// examples/sonner: sonner's own theme is light in both page themes).
func toaster(*http.Request) h.Node { return sonner.Toaster(sonner.ThemeLight) }

// Lib is the form library an example uses upstream.
type Lib string

const (
	rhf      Lib = "rhf"
	tanstack Lib = "tanstack"
	formisch Lib = "formisch"
)

// marker is the example's inert template: the examples' module finds the
// example by it (its previous sibling is the example's root).
func marker(name string) h.Node {
	return h.Template(h.Data("parity-ignore", ""), h.Data("form-example", name))
}

// with renders an example's root followed by its marker.
func with(name string, root h.Node) h.Node { return h.Group(root, marker(name)) }

// formTag is the <form>: Formisch's <Form> adds noValidate.
func formTag(l Lib, id string, children ...any) *h.Element {
	f := h.Form(h.ID(id))
	if l == formisch {
		f.Bool("novalidate", true)
	}
	return f.Children(children...)
}

// fieldName is the name attribute of a field's control: Formisch names
// fields by their JSON path.
func fieldName(l Lib, path ...any) string {
	if l == formisch {
		b, _ := json.Marshal(path)
		return string(b)
	}
	s := ""
	for i, p := range path {
		if i > 0 {
			s += "."
		}
		switch v := p.(type) {
		case string:
			s += v
		case int:
			b, _ := json.Marshal(v)
			s += string(b)
		}
	}
	return s
}

// invalid renders React's data-invalid={false} / aria-invalid={false}.
func dataInvalid() h.Prop { return h.Data("invalid", "false") }
func ariaInvalid() h.Prop { return h.Aria("invalid", "false") }

// footer is the examples' CardFooter: Reset and Submit.
func footer(formID, submit string, children ...any) *h.Element {
	return card.Footer(append(children,
		field.Field(field.Horizontal,
			button.Button(h.Type("button"), button.VariantOutline, "Reset"),
			button.Button(h.Type("submit"), h.FormAttr(formID), submit),
		),
	)...)
}

func header(title, description string, children ...any) *h.Element {
	return card.Header(append(children, card.Title(title), card.Description(description))...)
}

// The form component's client modules (form-state.js) are linked through
// this import.
var _ = uiform.Form
