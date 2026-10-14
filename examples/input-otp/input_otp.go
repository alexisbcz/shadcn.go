// Package inputotp ports the upstream input-otp-* examples.
package inputotp

import (
	_ "embed"

	inputotp "shadcn.go/components/ui/input-otp"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// controlled.js updates input-otp-controlled's text from the input's
// shadcn:value-change events, like the example's React state.
//
//go:embed controlled.js
var controlledScript string

func init() {
	runtime.Module("example-input-otp-controlled", controlledScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "input-otp", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("input-otp-controlled", Controlled),
		ex("input-otp-demo", Demo),
		ex("input-otp-pattern", Pattern),
		ex("input-otp-separator", Separator),
		ex("input-otp-disabled", Disabled),
	)
}

func slots(from, to int) []any {
	var out []any
	for i := from; i < to; i++ {
		out = append(out, inputotp.Slot(i))
	}
	return out
}

func Controlled() h.Node {
	return h.Div(h.Class("space-y-2"),
		inputotp.InputOTP(inputotp.MaxLength(6), inputotp.Group(slots(0, 6)...)),
		h.Div(h.Class("text-center text-sm"), "Enter your one-time password."),
	)
}

func Demo() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(6),
		inputotp.Group(slots(0, 3)...),
		inputotp.Separator(),
		inputotp.Group(slots(3, 6)...),
	)
}

func Pattern() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(6), inputotp.Pattern(inputotp.RegexpOnlyDigitsAndChars),
		inputotp.Group(slots(0, 6)...),
	)
}

func Separator() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(6),
		inputotp.Group(slots(0, 2)...),
		inputotp.Separator(),
		inputotp.Group(slots(2, 4)...),
		inputotp.Separator(),
		inputotp.Group(slots(4, 6)...),
	)
}

// Disabled ports the input-otp-disabled fixture
// (parity/fixtures/input-otp-disabled.tsx).
func Disabled() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(4), h.Disabled(),
		inputotp.Group(slots(0, 2)...),
		inputotp.Separator(),
		inputotp.Group(slots(2, 4)...),
	)
}
