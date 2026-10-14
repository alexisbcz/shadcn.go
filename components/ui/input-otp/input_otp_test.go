package inputotp_test

import (
	"strings"
	"testing"

	inputotp "shadcn.go/components/ui/input-otp"
	h "shadcn.go/libhtml"
)

func TestRender(t *testing.T) {
	got := h.MustRenderString(inputotp.InputOTP(inputotp.MaxLength(4), inputotp.Pattern(inputotp.RegexpOnlyDigits), h.Name("code"),
		inputotp.Group(inputotp.Slot(0), inputotp.Slot(1)),
		inputotp.Separator(),
		inputotp.Group(inputotp.Slot(2), inputotp.Slot(3)),
	))
	for _, s := range []string{
		`<noscript></noscript><div class="flex items-center gap-2 has-disabled:opacity-50" data-input-otp-container="true" style="position: relative; cursor: text; user-select: none; pointer-events: none;">`,
		`<div data-slot="input-otp-separator" role="separator"><svg`,
		`maxlength="4"`, `pattern="^\d+$"`, `name="code"`, `inputmode="numeric"`, `autocomplete="one-time-code"`,
		`data-input-otp-placeholder-shown="true"`, `value=""`, `class="disabled:cursor-not-allowed"`,
	} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in\n%s", s, got)
		}
	}
	if n := strings.Count(got, `data-slot="input-otp-slot" data-active="false"`); n != 4 {
		t.Errorf("%d slots", n)
	}
}

func TestValueFillsSlots(t *testing.T) {
	got := h.MustRenderString(inputotp.InputOTP(inputotp.MaxLength(3), h.Value("12345"), h.Disabled(),
		inputotp.Group(inputotp.Slot(0), inputotp.Slot(1), inputotp.Slot(2)),
	))
	if !strings.Contains(got, `value="123"`) || strings.Contains(got, "placeholder-shown") || !strings.Contains(got, "cursor: default") {
		t.Fatal(got)
	}
	for _, c := range []string{">1</div>", ">2</div>", ">3</div>"} {
		if !strings.Contains(got, c) {
			t.Errorf("missing %s in %s", c, got)
		}
	}
}

func TestEscaping(t *testing.T) {
	got := h.MustRenderString(inputotp.InputOTP(inputotp.MaxLength(2), h.Value(`<"`), inputotp.Group(inputotp.Slot(0), inputotp.Slot(1))))
	if strings.Contains(got, `<"`) || strings.Contains(got, `"<`) {
		t.Fatal(got)
	}
}
