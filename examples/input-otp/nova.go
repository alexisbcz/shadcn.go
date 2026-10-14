package inputotp

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	inputotp "shadcn.go/components/ui/input-otp"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the input-otp page
// (ui/apps/v4/examples/radix/input-otp-*.tsx), rendered in radix-nova.
// input-otp-controlled and -separator are the new-york-v4 examples.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "input-otp", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("input-otp-demo", NovaDemo),
		nova("input-otp-pattern", NovaPattern),
		nova("input-otp-alphanumeric", NovaAlphanumeric),
		nova("input-otp-disabled", NovaDisabled),
		nova("input-otp-four-digits", NovaFourDigits),
		nova("input-otp-invalid", NovaInvalid),
		nova("input-otp-form", NovaForm),
		nova("input-otp-rtl", NovaRtl),
	)
}

func NovaDemo() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(6), h.Value("123456"), inputotp.Group(slots(0, 6)...))
}

func NovaPattern() h.Node {
	return field.Field(h.Class("w-fit"),
		field.Label(h.For("digits-only"), "Digits Only"),
		inputotp.InputOTP(h.ID("digits-only"), inputotp.MaxLength(6), inputotp.Pattern(inputotp.RegexpOnlyDigits),
			inputotp.Group(slots(0, 6)...)),
	)
}

func NovaAlphanumeric() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(6), inputotp.Pattern(inputotp.RegexpOnlyDigitsAndChars),
		inputotp.Group(slots(0, 3)...),
		inputotp.Separator(),
		inputotp.Group(slots(3, 6)...),
	)
}

func NovaDisabled() h.Node {
	return inputotp.InputOTP(h.ID("disabled"), inputotp.MaxLength(6), h.Disabled(), h.Value("123456"),
		inputotp.Group(slots(0, 3)...),
		inputotp.Separator(),
		inputotp.Group(slots(3, 6)...),
	)
}

func NovaFourDigits() h.Node {
	return inputotp.InputOTP(inputotp.MaxLength(4), inputotp.Pattern(inputotp.RegexpOnlyDigits), inputotp.Group(slots(0, 4)...))
}

func NovaInvalid() h.Node {
	slot := func(i int) any { return inputotp.Slot(i, h.Aria("invalid", "true")) }
	return inputotp.InputOTP(inputotp.MaxLength(6), h.Value("000000"),
		inputotp.Group(slot(0), slot(1)),
		inputotp.Separator(),
		inputotp.Group(slot(2), slot(3)),
		inputotp.Separator(),
		inputotp.Group(slot(4), slot(5)),
	)
}

func NovaForm() h.Node {
	big := "*:data-[slot=input-otp-slot]:h-12 *:data-[slot=input-otp-slot]:w-11 *:data-[slot=input-otp-slot]:text-xl"
	return card.Card(h.Class("mx-auto max-w-md"),
		card.Header(
			card.Title("Verify your login"),
			card.Description("Enter the verification code we sent to your email address:", " ",
				h.Span(h.Class("font-medium"), "m@example.com"), "."),
		),
		card.Content(
			field.Field(
				h.Div(h.Class("flex items-center justify-between"),
					field.Label(h.For("otp-verification"), "Verification code"),
					button.Button(button.VariantOutline, button.SizeXs, lucide.RefreshCw(), "Resend Code"),
				),
				inputotp.InputOTP(inputotp.MaxLength(6), h.ID("otp-verification"), h.Attr("required", ""),
					inputotp.Group(h.Class(big), slots(0, 3)),
					inputotp.Separator(h.Class("mx-2")),
					inputotp.Group(h.Class(big), slots(3, 6)),
				),
				field.Description(h.A(h.Href("#"), "I no longer have access to this email address.")),
			),
		),
		card.Footer(
			field.Field(
				button.Button(h.Type("submit"), h.Class("w-full"), "Verify"),
				h.Div(h.Class("text-sm text-muted-foreground"), "Having trouble signing in?", " ",
					h.A(h.Href("#"), h.Class("underline underline-offset-4 transition-colors hover:text-primary"), "Contact support")),
			),
		),
	)
}

// NovaRtl is examples/radix/input-otp-rtl.tsx in its default language
// (Arabic).
func NovaRtl() h.Node {
	return field.Field(h.Class("mx-auto max-w-xs"),
		field.Label(h.For("input-otp-rtl"), "رمز التحقق"),
		inputotp.InputOTP(inputotp.MaxLength(6), h.Value("123456"), h.Dir("rtl"), h.ID("input-otp-rtl"),
			inputotp.Group(slots(0, 6)...)),
	)
}
