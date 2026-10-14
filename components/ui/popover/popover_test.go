package popover_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/popover"
	h "shadcn.go/libhtml"
)

func TestPopoverMarkup(t *testing.T) {
	html := h.MustRenderString(popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		popover.Anchor("anchor"),
		popover.Content(h.Class("w-80"), popover.SideTop, popover.AlignStart, popover.SideOffset(8),
			popover.Header(popover.Title("T"), popover.Description("D")),
			popover.Close("x"),
			popover.Arrow(),
		),
	))
	id := regexp.MustCompile(`aria-controls="([^"]+)"`).FindStringSubmatch(html)
	if id == nil {
		t.Fatal("no aria-controls: " + html)
	}
	for _, want := range []string{
		`data-slot="popover-trigger"`, `aria-haspopup="dialog"`, `aria-expanded="false"`, `data-state="closed"`, `type="button"`,
		`data-slot="popover-anchor"`, `data-radix-popover-anchor="` + id[1] + `"`,
		`<template data-parity-ignore="" data-radix-popover="" data-content-id="` + id[1] + `"`,
		`id="` + id[1] + `"`, `role="dialog"`, `tabindex="-1"`, `data-slot="popover-content"`,
		`data-side="top"`, `data-align="start"`, `&#34;sideOffset&#34;:8`,
		`--radix-popover-content-transform-origin: var(--radix-popper-transform-origin)`,
		`data-slot="popover-header"`, `data-slot="popover-title"`, `data-slot="popover-description"`,
		`data-radix-popover-close="` + id[1] + `"`,
		`data-radix-popper-arrow=""`, `bottom: 0px;`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	// className replaces w-72 through cn; the trigger is the Button itself.
	if strings.Contains(html, " w-72") {
		t.Error("w-80 must replace w-72")
	}
	if strings.Contains(html, `data-slot="button"`) {
		t.Error("the trigger's data-slot must replace the button's")
	}
}

func TestPopoverOptions(t *testing.T) {
	html := h.MustRenderString(popover.Popover(popover.Modal(true), popover.DefaultOpen, popover.Trigger("t"), popover.Content("c")))
	for _, want := range []string{`data-modal=""`, `data-default-open=""`, `aria-expanded="true"`, `data-state="open"`, `&#34;sideOffset&#34;:4`, `data-align="center"`, `data-side="bottom"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
}

func TestPopoverEscapes(t *testing.T) {
	html := h.MustRenderString(popover.Popover(popover.Trigger("<b>"), popover.Content(`"><script>`)))
	if strings.Contains(html, "<b>") || strings.Contains(html, "<script>") {
		t.Fatal(html)
	}
}
