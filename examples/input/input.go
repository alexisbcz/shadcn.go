// Package input ports the upstream input-* examples.
package input

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "input", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("input-demo", Demo),
		ex("input-disabled", Disabled),
		ex("input-file", File),
		ex("input-with-button", WithButton),
		ex("input-with-label", WithLabel),
		ex("input-with-text", WithText),
	)
}

func Demo() h.Node { return input.Input(h.Type("email"), h.Placeholder("Email")) }

func Disabled() h.Node {
	return input.Input(h.Disabled(), h.Type("email"), h.Placeholder("Email"))
}

func File() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm items-center gap-3"),
		label.Label(h.For("picture"), "Picture"),
		input.Input(h.ID("picture"), h.Type("file")),
	)
}

func WithButton() h.Node {
	return h.Div(h.Class("flex w-full max-w-sm items-center gap-2"),
		input.Input(h.Type("email"), h.Placeholder("Email")),
		button.Button(h.Type("submit"), button.VariantOutline, "Subscribe"),
	)
}

func WithLabel() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm items-center gap-3"),
		label.Label(h.For("email"), "Email"),
		input.Input(h.Type("email"), h.ID("email"), h.Placeholder("Email")),
	)
}

func WithText() h.Node {
	return h.Div(h.Class("grid w-full max-w-sm items-center gap-3"),
		label.Label(h.For("email-2"), "Email"),
		input.Input(h.Type("email"), h.ID("email-2"), h.Placeholder("Email")),
		h.P(h.Class("text-sm text-muted-foreground"), "Enter your email address."),
	)
}
