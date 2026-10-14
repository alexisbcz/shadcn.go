package selectui_test

import (
	"regexp"
	"strings"
	"testing"

	selectui "shadcn.go/components/ui/select"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestSelectPlaceholder(t *testing.T) {
	html := render(t, selectui.Select(
		selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Pick"))),
		selectui.Content(
			selectui.Group(selectui.Label("Fruits"), selectui.Item(selectui.ItemValue("apple"), "Apple")),
			selectui.Separator(),
			selectui.Item(selectui.ItemValue("b"), selectui.Disabled, selectui.TextValue("Bee"), "B"),
		),
	))
	id := regexp.MustCompile(`aria-controls="([^"]+)"`).FindStringSubmatch(html)
	if id == nil {
		t.Fatal("no aria-controls: " + html)
	}
	for _, want := range []string{
		`role="combobox"`, `aria-expanded="false"`, `aria-autocomplete="none"`, `dir="ltr"`, `data-state="closed"`,
		`data-placeholder=""`, `data-size="default"`, `w-[180px]`,
		`data-slot="select-value" style="pointer-events: none;">Pick</span>`,
		`<template data-parity-ignore="" data-radix-select="" data-content-id="` + id[1] + `" data-position="item-aligned"`,
		`role="listbox" id="` + id[1] + `"`, `data-radix-select-viewport=""`,
		`role="group" aria-labelledby="`, `data-slot="select-label"`, `data-slot="select-separator"`,
		`data-radix-select-item="apple"`, `data-state="unchecked"`, `data-radix-collection-item=""`,
		`aria-disabled="true" data-disabled=""`, `data-radix-select-text-value="Bee"`,
		`data-radix-select-scroll-up=""`, `data-radix-select-scroll-down=""`,
		`<select data-radix-select-bubble="` + id[1] + `" aria-hidden="true" tabindex="-1"`,
		`<option value=""></option><option value="apple">Apple</option><option value="b" disabled="">B</option>`,
		`<span data-radix-select-placeholder="">Pick</span></template>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s\n%s", want, html)
		}
	}
	if w := "w-fit"; strings.Contains(html, w) {
		t.Errorf("class w-fit should be merged away")
	}
	if strings.Contains(html, `data-highlighted`) {
		t.Errorf("unexpected highlight")
	}
}

func TestSelectDefaultValue(t *testing.T) {
	html := render(t, selectui.Select(selectui.DefaultValue("b"), selectui.Name("fruit"), selectui.Required, selectui.Form("f"),
		selectui.Trigger(selectui.SizeSm, selectui.Value(selectui.Placeholder("Pick"))),
		selectui.Content(selectui.PositionPopper, selectui.SideTop,
			selectui.Item(selectui.ItemValue("a"), "Alpha"),
			selectui.Item(selectui.ItemValue("b"), "<Beta>"),
		),
	))
	for _, want := range []string{
		`aria-required="true"`, `data-size="sm"`,
		`data-slot="select-value" style="pointer-events: none;">&lt;Beta&gt;</span>`,
		`data-value="b"`, `data-position="popper"`, `data-radix-popper-content-wrapper=""`, `data-side="top"`,
		`--radix-select-trigger-width: var(--radix-popper-anchor-width)`,
		`data-[side=bottom]:translate-y-1`, `min-w-[var(--radix-select-trigger-width)]`,
		`data-state="checked"`, `name="fruit"`, `required=""`, `form="f"`,
		`<option value="b" selected="">&lt;Beta&gt;</option>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s\n%s", want, html)
		}
	}
	if strings.Contains(html, `data-placeholder`) {
		t.Error("placeholder flag with a value")
	}
	if strings.Contains(html, `<option value=""></option>`) {
		t.Error("empty option with a value")
	}
}

func TestSelectDisabled(t *testing.T) {
	html := render(t, selectui.Select(selectui.Disabled, selectui.Trigger(selectui.Value()), selectui.Content()))
	for _, want := range []string{`disabled="" data-disabled=""`, `data-disabled=""`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s\n%s", want, html)
		}
	}
	html = render(t, selectui.Select(selectui.Trigger(h.Disabled(), selectui.Value()), selectui.Content()))
	if !strings.Contains(html, `disabled="" data-disabled=""`) {
		t.Errorf("trigger disabled attribute not handled\n%s", html)
	}
}

func TestSelectItemNeedsValue(t *testing.T) {
	var b strings.Builder
	if err := selectui.Item("x").Render(&b); err == nil {
		t.Fatal("expected an error for an item without value")
	}
}

func TestSelectValueChildren(t *testing.T) {
	html := render(t, selectui.Select(selectui.DefaultValue("a"),
		selectui.Trigger(selectui.Value(h.Class("ignored"), "Custom")),
		selectui.Content(selectui.Item(selectui.ItemValue("a"), "Alpha")),
	))
	if !strings.Contains(html, `data-radix-select-value-children="" data-slot="select-value" style="pointer-events: none;">Custom</span>`) {
		t.Errorf("value children not rendered\n%s", html)
	}
	if strings.Contains(html, "ignored") {
		t.Error("SelectValue must drop className like Radix")
	}
}
