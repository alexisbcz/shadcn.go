// Package switchui ports the upstream switch-* examples.
package switchui

import (
	"shadcn.go/components/ui/label"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "switch", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("switch-demo", Demo),
		ex("switch-states", States),
		ex("switch-form", Form),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex items-center space-x-2"),
		switchui.Switch(h.ID("airplane-mode")),
		label.Label(h.For("airplane-mode"), "Airplane Mode"),
	)
}

// States ports the switch-states fixture (parity/fixtures/switch-states.tsx):
// the props the upstream example does not use: the sm size, defaultChecked,
// disabled (checked and unchecked), aria-invalid, and switches in a form
// (hidden bubble inputs) with name, value, required and a reset button.
func States() h.Node {
	row := func(id, text string, sw *switchui.Builder) h.Node {
		return h.Div(h.Class("flex items-center gap-2"), sw.ID(id), label.Label(h.For(id), text))
	}
	return h.Div(h.Class("flex flex-col gap-4"),
		row("small", "Small", switchui.Switch(switchui.SizeSm)),
		row("small-on", "Small on", switchui.Switch(switchui.SizeSm, switchui.Checked(true))),
		row("on", "On", switchui.Switch(switchui.Checked(true))),
		row("off-disabled", "Disabled", switchui.Switch(h.Disabled())),
		row("on-disabled", "Disabled on", switchui.Switch(switchui.Checked(true), h.Disabled())),
		row("invalid", "Invalid", switchui.Switch(h.Aria("invalid", "true"))),
		h.Form(h.Class("flex flex-col gap-3"),
			row("wifi", "Wi-Fi", switchui.Switch(switchui.Name("wifi"), switchui.Checked(true))),
			row("bt", "Bluetooth", switchui.Switch(switchui.Name("bt"), switchui.Value("yes"), switchui.Required(true))),
			h.Button(h.Type("reset"), h.Class("w-fit rounded-md border px-2 text-sm"), "Reset"),
		),
	)
}

// Form ports the switch-form fixture (parity/fixtures/switch-form.tsx):
// switches tied to a form through the form attribute (outside the form),
// and a disabled switch inside a form.
func Form() h.Node {
	row := func(id, text string, sw *switchui.Builder) h.Node {
		return h.Div(h.Class("flex items-center gap-2"), sw.ID(id), label.Label(h.For(id), text))
	}
	return h.Div(h.Class("flex flex-col gap-4"),
		h.Form(h.ID("switch-form"), h.Class("flex gap-3"),
			h.Button(h.Type("reset"), h.Class("w-fit rounded-md border px-2 text-sm"), "Reset"),
		),
		row("remote", "Remote", switchui.Switch(h.FormAttr("switch-form"), h.Name("remote"), h.Value("yes"), switchui.Checked(true))),
		row("remote-off", "Remote off", switchui.Switch(switchui.Form("switch-form"), switchui.Name("remote-off"), h.Required())),
		h.Form(h.Class("flex items-center gap-2"),
			switchui.Switch(switchui.Name("locked"), h.Disabled(), switchui.Checked(true), h.ID("locked")),
			label.Label(h.For("locked"), "Locked"),
		),
	)
}
