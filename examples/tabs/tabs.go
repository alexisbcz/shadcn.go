// Package tabs ports the upstream tabs-* examples.
package tabs

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(registry.Example{Name: "tabs-demo", Component: "tabs", Render: registry.Static(Demo)})
	registry.Register(registry.Example{Name: "tabs-behavior", Component: "tabs", Render: registry.Static(Behavior)})
}

func Demo() h.Node {
	field := func(id, text string, attrs ...any) h.Node {
		return h.Div(h.Class("grid gap-3"),
			label.Label(h.For(id), text),
			input.Input(append([]any{h.ID(id)}, attrs...)...),
		)
	}
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-6"),
		tabs.Tabs(tabs.Value("account"),
			tabs.List(
				tabs.Trigger(tabs.Value("account"), "Account"),
				tabs.Trigger(tabs.Value("password"), "Password"),
			),
			tabs.Content(tabs.Value("account"),
				card.Card(
					card.Header(
						card.Title("Account"),
						card.Description("Make changes to your account here. Click save when you're done."),
					),
					card.Content(h.Class("grid gap-6"),
						field("tabs-demo-name", "Name", h.Value("Pedro Duarte")),
						field("tabs-demo-username", "Username", h.Value("@peduarte")),
					),
					card.Footer(button.Button("Save changes")),
				),
			),
			tabs.Content(tabs.Value("password"),
				card.Card(
					card.Header(
						card.Title("Password"),
						card.Description("Change your password here. After saving, you'll be logged out."),
					),
					card.Content(h.Class("grid gap-6"),
						field("tabs-demo-current", "Current password", h.Type("password")),
						field("tabs-demo-new", "New password", h.Type("password")),
					),
					card.Footer(button.Button("Save password")),
				),
			),
		),
	)
}

// Behavior ports the tabs-behavior fixture (parity/fixtures/tabs-behavior.tsx):
// vertical orientation, manual activation, a disabled trigger, the line
// list variant, Loop(false) and dir="rtl".
func Behavior() h.Node {
	return h.Div(h.Class("flex flex-col gap-6"),
		tabs.Tabs(tabs.Value("alpha"), tabs.OrientationVertical, tabs.ActivationManual,
			tabs.List(
				tabs.Trigger(tabs.Value("alpha"), "Alpha"),
				tabs.Trigger(tabs.Value("beta"), h.Disabled(), "Beta"),
				tabs.Trigger(tabs.Value("gamma"), "Gamma"),
			),
			tabs.Content(tabs.Value("alpha"), "Alpha panel"),
			tabs.Content(tabs.Value("beta"), "Beta panel"),
			tabs.Content(tabs.Value("gamma"), "Gamma panel"),
		),
		tabs.Tabs(tabs.Value("one"), h.Dir("rtl"),
			tabs.List(tabs.ListVariantLine, tabs.Loop(false),
				tabs.Trigger(tabs.Value("one"), "One"),
				tabs.Trigger(tabs.Value("two"), "Two"),
				tabs.Trigger(tabs.Value("three"), "Three"),
			),
			tabs.Content(tabs.Value("one"), "First panel"),
			tabs.Content(tabs.Value("two"), "Second panel"),
			tabs.Content(tabs.Value("three"), "Third panel"),
		),
	)
}
