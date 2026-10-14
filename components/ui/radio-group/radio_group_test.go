package radiogroup_test

import (
	"strings"
	"testing"

	rg "shadcn.go/components/ui/radio-group"
	h "shadcn.go/libhtml"
)

func TestRadioGroup(t *testing.T) {
	got := h.MustRenderString(rg.RadioGroup(rg.Value("b"), rg.Name("density"),
		h.Div(rg.Item(rg.Value("a"), h.ID("r1"))),
		h.Div(rg.Item(rg.Value("b"), h.ID("r2"))),
	))
	for _, want := range []string{
		`role="radiogroup"`, `aria-required="false"`, `dir="ltr"`, `tabindex="0"`, `style="outline:none"`,
		`data-slot="radio-group"`, `class="grid gap-3"`,
		`value="a"`, `value="b"`, `name="density"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Count(got, `aria-checked="true"`) != 1 || strings.Count(got, `data-slot="radio-group-indicator"`) != 2 {
		t.Errorf("checked state: %s", got)
	}
	if strings.Count(got, `<template data-slot="radio-group-indicator-template">`) != 1 {
		t.Errorf("unchecked item carries the indicator template: %s", got)
	}
	if strings.Count(got, `<input `) != 2 || strings.Count(got, ` checked`) != 1 {
		t.Errorf("bubble inputs: %s", got)
	}
}

func TestRadioGroupDisabled(t *testing.T) {
	got := h.MustRenderString(rg.RadioGroup(rg.Disabled(true), rg.Item(rg.Value("a"))))
	if !strings.Contains(got, `tabindex="-1" style="outline:none"`) {
		t.Errorf("no focusable item: group not a tab stop: %s", got)
	}
	if !strings.Contains(got, `data-disabled=""`) || !strings.Contains(got, " disabled") {
		t.Errorf("disabled: %s", got)
	}
	if _, err := h.RenderString(rg.Item()); err == nil {
		t.Error("item outside a group must fail")
	}
}

func TestRadioGroupOptions(t *testing.T) {
	got := h.MustRenderString(rg.RadioGroup(h.Dir("rtl"), rg.Loop(false), rg.OrientationHorizontal, rg.Item(rg.Value("a"))))
	for _, want := range []string{`dir="rtl"`, `data-loop="false"`, `aria-orientation="horizontal"`, `data-orientation="horizontal"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, `dir="ltr"`) {
		t.Errorf("caller's dir must replace the default: %s", got)
	}
}
