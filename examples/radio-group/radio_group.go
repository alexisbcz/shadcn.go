// Package radiogroup ports the upstream radio-group-* examples.
package radiogroup

import (
	"shadcn.go/components/ui/label"
	rg "shadcn.go/components/ui/radio-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	registry.Register(registry.Example{Name: "radio-group-demo", Component: "radio-group", Render: registry.Static(Demo)})
	registry.Register(registry.Example{Name: "radio-group-behavior", Component: "radio-group", Render: registry.Static(Behavior)})
}

func Demo() h.Node {
	item := func(value, id, text string) h.Node {
		return h.Div(h.Class("flex items-center gap-3"),
			rg.Item(rg.Value(value), h.ID(id)),
			label.Label(h.For(id), text),
		)
	}
	return rg.RadioGroup(rg.Value("comfortable"),
		item("default", "r1", "Default"),
		item("comfortable", "r2", "Comfortable"),
		item("compact", "r3", "Compact"),
	)
}

// Behavior ports the radio-group-behavior fixture
// (parity/fixtures/radio-group-behavior.tsx): a group in a form (hidden
// bubble inputs) with Name, Required, horizontal orientation and a disabled
// item; a group without value, dir="rtl" and Loop(false); a disabled group.
func Behavior() h.Node {
	item := func(value, id, text string, attrs ...any) h.Node {
		return h.Div(h.Class("flex items-center gap-2"),
			rg.Item(append([]any{rg.Value(value), h.ID(id)}, attrs...)...),
			label.Label(h.For(id), text),
		)
	}
	return h.Div(h.Class("flex flex-col gap-6"),
		h.Form(
			rg.RadioGroup(rg.Name("plan"), rg.Value("pro"), rg.OrientationHorizontal, rg.Required(true), h.Class("flex"),
				item("free", "plan-free", "Free"),
				item("pro", "plan-pro", "Pro"),
				item("team", "plan-team", "Team", h.Disabled()),
			),
		),
		rg.RadioGroup(h.Dir("rtl"), rg.Loop(false),
			item("a", "rtl-a", "First"),
			item("b", "rtl-b", "Second"),
		),
		rg.RadioGroup(rg.Value("x"), rg.Disabled(true),
			item("x", "dis-x", "Locked"),
		),
	)
}
