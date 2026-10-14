package textarea

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the textarea page
// (ui/apps/v4/examples/radix/textarea-*.tsx), in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "textarea", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	const ph = "Type your message here."
	registry.Register(
		nova("textarea-button", func() h.Node {
			return h.Div(h.Class("grid w-full gap-2"), textarea.Textarea(h.Placeholder(ph)), button.Button("Send message"))
		}),
		nova("textarea-demo", func() h.Node { return textarea.Textarea(h.Placeholder(ph)) }),
		nova("textarea-disabled", func() h.Node {
			return field.Field(h.Data("disabled", "true"),
				field.Label(h.For("textarea-disabled"), "Message"),
				textarea.Textarea(h.ID("textarea-disabled"), h.Placeholder(ph), h.Bool("disabled")),
			)
		}),
		nova("textarea-field", func() h.Node {
			return field.Field(
				field.Label(h.For("textarea-message"), "Message"),
				field.Description("Enter your message below."),
				textarea.Textarea(h.ID("textarea-message"), h.Placeholder(ph)),
			)
		}),
		nova("textarea-invalid", func() h.Node {
			return field.Field(h.Data("invalid", "true"),
				field.Label(h.For("textarea-invalid"), "Message"),
				textarea.Textarea(h.ID("textarea-invalid"), h.Placeholder(ph), h.Aria("invalid", "true")),
				field.Description("Please enter a valid message."),
			)
		}),
		nova("textarea-rtl", func() h.Node {
			// text-start: the ui-rtl FieldDescription's text-left → text-start
			// (field does not map its classes under style.ScopeRTL).
			return style.ScopeRTL(field.Field(h.Class("w-full max-w-xs"), h.Dir("rtl"),
				field.Label(h.For("feedback"), h.Dir("rtl"), "التعليقات"),
				textarea.Textarea(h.ID("feedback"), h.Placeholder("تعليقاتك تساعدنا على التحسين..."), h.Dir("rtl"), h.Attr("rows", "4")),
				field.Description(h.Dir("rtl"), h.Class("text-start"), "شاركنا أفكارك حول خدمتنا."),
			))
		}),
	)
}
