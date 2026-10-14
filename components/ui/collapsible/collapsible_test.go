package collapsible_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/collapsible"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	s, err := h.RenderString(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestClosed(t *testing.T) {
	s := render(t, collapsible.Collapsible(h.Class("flex"),
		collapsible.Trigger("Toggle"),
		collapsible.Content(h.Class("c"), h.P("Body")),
	))
	for _, want := range []string{
		`<div class="flex" data-state="closed" data-slot="collapsible">`,
		`type="button"`, `aria-expanded="false"`,
		`data-slot="collapsible-content"`, `hidden`, `style=""`,
		`<template data-slot="collapsible-content-template"><p>Body</p></template>`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
	if strings.Count(s, "<p>Body</p>") != 1 {
		t.Errorf("closed content must only carry its children in the template: %s", s)
	}
	// aria-controls points at the content.
	i := strings.Index(s, `aria-controls="`) + len(`aria-controls="`)
	id := s[i : i+strings.IndexByte(s[i:], '"')]
	if !strings.Contains(s, `id="`+id+`"`) {
		t.Errorf("aria-controls %q does not match the content id: %s", id, s)
	}
}

func TestOpenDisabledAsChild(t *testing.T) {
	s := render(t, collapsible.Collapsible(collapsible.Open(true), collapsible.Disabled(true),
		collapsible.Trigger(lib.AsChild, button.Button(button.VariantGhost, "T")),
		collapsible.Content(h.Style("color:red"), h.P("Body")),
	))
	for _, want := range []string{
		`data-state="open" data-disabled="" data-slot="collapsible"`,
		`data-slot="collapsible-trigger"`, `data-variant="ghost"`, `aria-expanded="true"`, `disabled`,
		`style="color:red;transition-duration:0s;animation-name:none"`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
	if strings.Contains(s, `data-slot="button"`) {
		t.Errorf("the trigger's data-slot must win over the button's: %s", s)
	}
	if strings.Count(s, "<p>Body</p>") != 2 {
		t.Errorf("open content renders its children and the template: %s", s)
	}
}

func TestPartsOutsideRoot(t *testing.T) {
	for _, n := range []h.Node{collapsible.Trigger(), collapsible.Content()} {
		if _, err := h.RenderString(n); err == nil {
			t.Errorf("%T outside Collapsible must fail", n)
		}
	}
}
