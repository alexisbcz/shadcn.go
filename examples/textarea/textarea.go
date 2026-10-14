// Package textarea ports the upstream textarea-* examples.
package textarea

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "textarea", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("textarea-demo", Demo),
		ex("textarea-disabled", Disabled),
		ex("textarea-with-button", WithButton),
		ex("textarea-with-label", WithLabel),
		ex("textarea-with-text", WithText),
	)
}

const placeholder = "Type your message here."

func Demo() h.Node { return textarea.Textarea(h.Placeholder(placeholder)) }

func Disabled() h.Node { return textarea.Textarea(h.Placeholder(placeholder), h.Disabled()) }

func WithButton() h.Node {
	return h.Div(h.Class("grid w-full gap-2"),
		textarea.Textarea(h.Placeholder(placeholder)),
		button.Button("Send message"),
	)
}

func WithLabel() h.Node {
	return h.Div(h.Class("grid w-full gap-3"),
		label.Label(h.For("message"), "Your message"),
		textarea.Textarea(h.Placeholder(placeholder), h.ID("message")),
	)
}

func WithText() h.Node {
	return h.Div(h.Class("grid w-full gap-3"),
		label.Label(h.For("message-2"), "Your Message"),
		textarea.Textarea(h.Placeholder(placeholder), h.ID("message-2")),
		h.P(h.Class("text-sm text-muted-foreground"), "Your message will be copied to the support team."),
	)
}
