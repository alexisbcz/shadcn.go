package toggle_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/toggle"
	h "shadcn.go/libhtml"
)

func TestToggleMarkup(t *testing.T) {
	got := h.MustRenderString(toggle.Toggle(h.Aria("label", "Bold"), "B"))
	for _, want := range []string{`<button `, `type="button"`, `aria-pressed="false"`, `data-state="off"`, `data-slot="toggle"`, `aria-label="Bold"`, `>B</button>`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, "data-disabled") || strings.Contains(got, `" disabled`) || strings.Contains(got, " disabled>") {
		t.Errorf("unexpected disabled: %s", got)
	}
}

func TestTogglePressedDisabled(t *testing.T) {
	a := h.MustRenderString(toggle.Toggle(toggle.Pressed(true), h.Disabled(), "x"))
	b := h.MustRenderString(toggle.Toggle("x").Pressed(true).Disabled(true))
	if a != b {
		t.Fatalf("\n%s\n%s", a, b)
	}
	for _, want := range []string{`aria-pressed="true"`, `data-state="on"`, `data-disabled=""`, ` disabled`} {
		if !strings.Contains(a, want) {
			t.Errorf("missing %s in %s", want, a)
		}
	}
}

func TestToggleVariants(t *testing.T) {
	cls := toggle.Variants(toggle.VariantOutline, toggle.SizeLg, "px-8")
	for _, want := range []string{"border-input", "h-10", "min-w-10", "px-8", "shadow-xs"} {
		if !strings.Contains(" "+cls+" ", " "+want+" ") {
			t.Errorf("missing %s in %s", want, cls)
		}
	}
	if strings.Contains(cls, "px-2.5") {
		t.Errorf("className did not override px: %s", cls)
	}
	if toggle.Variants("", "") != toggle.Variants(toggle.VariantDefault, toggle.SizeDefault) {
		t.Error("empty values must select defaults")
	}
	got := h.MustRenderString(toggle.Toggle(toggle.SizeSm, h.Class("rounded-full")))
	if !strings.Contains(got, `class="`+h.EscapeAttr(toggle.Variants("", toggle.SizeSm, "rounded-full"))+`"`) {
		t.Errorf("classes differ: %s", got)
	}
}

func TestToggleEscaping(t *testing.T) {
	payload := `"><script>alert(1)</script>`
	got := h.MustRenderString(toggle.Toggle(payload, h.Aria("label", payload)))
	if strings.Contains(got, "<script>") {
		t.Fatal(got)
	}
}
