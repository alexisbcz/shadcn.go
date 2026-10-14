// Package card ports the upstream card-* examples.
package card

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(
		registry.Example{Name: "card-demo", Component: "card", Render: registry.Static(Demo)},
	)
}

func Demo() h.Node {
	return card.Card(h.Class("w-full max-w-sm"),
		card.Header(
			card.Title("Login to your account"),
			card.Description("Enter your email below to login to your account"),
			card.Action(button.Button(button.VariantLink, "Sign Up")),
		),
		card.Content(
			h.Form(
				h.Div(h.Class("flex flex-col gap-6"),
					h.Div(h.Class("grid gap-2"),
						label.Label(h.For("email"), "Email"),
						input.Input(h.ID("email"), h.Type("email"), h.Placeholder("m@example.com"), h.Required()),
					),
					h.Div(h.Class("grid gap-2"),
						h.Div(h.Class("flex items-center"),
							label.Label(h.For("password"), "Password"),
							h.A(h.Href("#"), h.Class("ml-auto inline-block text-sm underline-offset-4 hover:underline"),
								"Forgot your password?",
							),
						),
						input.Input(h.ID("password"), h.Type("password"), h.Required()),
					),
				),
			),
		),
		card.Footer(h.Class("flex-col gap-2"),
			button.Button(h.Type("submit"), h.Class("w-full"), "Login"),
			button.Button(button.VariantOutline, h.Class("w-full"), "Login with Google"),
		),
	)
}
