package checkbox_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/checkbox"
	h "shadcn.go/libhtml"
)

func has(t *testing.T, got string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in\n%s", w, got)
		}
	}
}

func TestUnchecked(t *testing.T) {
	got := h.MustRenderString(checkbox.Checkbox(h.ID("terms")))
	has(t, got, `role="checkbox"`, `aria-checked="false"`, `data-state="unchecked"`, `value="on"`, `type="button"`,
		`<template data-slot="checkbox-indicator-template">`, `<input type="checkbox" aria-hidden="true"`, `tabindex="-1"`)
	if strings.Contains(got, "aria-required") || strings.Contains(got, "data-disabled") || strings.Contains(got, " checked") {
		t.Errorf("unexpected attribute: %s", got)
	}
}

func TestCheckedAndIndeterminate(t *testing.T) {
	got := h.MustRenderString(checkbox.Checkbox(checkbox.Checked(true), checkbox.Name("n"), checkbox.Value("yes"), checkbox.Required(true)))
	has(t, got, `aria-checked="true"`, `data-state="checked"`, `data-slot="checkbox-indicator"`, `aria-required="true"`, ` checked`, `name="n"`, `value="yes"`, ` required`)
	if strings.Contains(got, "<template") {
		t.Errorf("checked checkbox must render its indicator: %s", got)
	}
	got = h.MustRenderString(checkbox.Checkbox(checkbox.Indeterminate))
	has(t, got, `aria-checked="mixed"`, `data-state="indeterminate"`, `data-slot="checkbox-indicator"`)
	if strings.Contains(got, " checked") {
		t.Errorf("indeterminate must not check the bubble input: %s", got)
	}
}

func TestDisabledAndBuilder(t *testing.T) {
	a := h.MustRenderString(checkbox.Checkbox(h.Disabled(), checkbox.Checked(true)))
	b := h.MustRenderString(checkbox.Checkbox().Disabled(true).Checked(true))
	has(t, a, `data-disabled=""`, ` disabled`)
	has(t, b, `data-disabled=""`, ` disabled`)
	if strings.Count(a, `data-disabled=""`) != 2 {
		t.Errorf("button and indicator must both carry data-disabled: %s", a)
	}
}

func TestClassMerge(t *testing.T) {
	got := h.MustRenderString(checkbox.Checkbox(h.Class("rounded-full <x>")))
	has(t, got, "rounded-full", "&lt;x&gt;")
	if strings.Contains(got, "rounded-[4px]") {
		t.Errorf("className must override rounded: %s", got)
	}
}

func TestFormAttributesGoToInput(t *testing.T) {
	got := h.MustRenderString(checkbox.Checkbox(h.FormAttr("f"), h.Name("n"), h.Value("yes"), h.Required()))
	typed := h.MustRenderString(checkbox.Checkbox(checkbox.Form("f"), checkbox.Name("n"), checkbox.Value("yes"), checkbox.Required(true)))
	if got != typed {
		t.Errorf("attributes and typed options differ:\n%s\n%s", got, typed)
	}
	button, input, _ := strings.Cut(got, "<input")
	has(t, button, `value="yes"`, `aria-required="true"`)
	has(t, input, `form="f"`, `name="n"`, `value="yes"`, ` required`)
	for _, bad := range []string{`form=`, `name=`, ` required`} {
		if strings.Contains(button, bad) {
			t.Errorf("button must not carry %s: %s", bad, button)
		}
	}
}
