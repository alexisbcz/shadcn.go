package togglegroup_test

import (
	"strings"
	"testing"

	tg "shadcn.go/components/ui/toggle-group"
	h "shadcn.go/libhtml"
)

func items() []any {
	return []any{
		tg.Item(tg.Value("a"), "A"),
		h.Span(tg.Item(tg.Value("b"), "B")), // context reaches nested items
		tg.Item(tg.Value("c"), h.Disabled(), "C"),
	}
}

func TestMultiple(t *testing.T) {
	got := h.MustRenderString(tg.ToggleGroup(tg.TypeMultiple, tg.VariantOutline, tg.Value("a"), tg.Value("b"), items()))
	for _, want := range []string{
		`role="group"`, `dir="ltr"`, `tabindex="0"`, `style="outline:none;--gap:0"`, `data-slot="toggle-group"`,
		`data-variant="outline"`, `data-spacing="0"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Count(got, `aria-pressed="true"`) != 2 || strings.Count(got, `aria-pressed="false"`) != 1 {
		t.Errorf("pressed state: %s", got)
	}
	if strings.Count(got, `data-radix-collection-item=""`) != 3 || strings.Count(got, `tabindex="-1"`) != 3 {
		t.Errorf("roving items: %s", got)
	}
	if !strings.Contains(got, `data-disabled="" disabled`) && !strings.Contains(got, `disabled data-disabled=""`) && strings.Count(got, "data-disabled") != 1 {
		t.Errorf("disabled item: %s", got)
	}
	if strings.Contains(got, `value="`) {
		t.Errorf("Radix does not render item values: %s", got)
	}
}

func TestSingle(t *testing.T) {
	got := h.MustRenderString(tg.ToggleGroup(tg.TypeSingle, tg.Value("b"), items()))
	if strings.Contains(got, "aria-pressed") {
		t.Errorf("single items use aria-checked: %s", got)
	}
	if strings.Count(got, `role="radio"`) != 3 || strings.Count(got, `aria-checked="true"`) != 1 {
		t.Errorf("single: %s", got)
	}
	if strings.Contains(got, "data-variant") || strings.Contains(got, "data-size") {
		t.Errorf("no variant/size given, no attributes: %s", got)
	}
}

func TestDisabledGroup(t *testing.T) {
	got := h.MustRenderString(tg.ToggleGroup(tg.TypeMultiple, tg.Disabled(true), items()))
	if !strings.Contains(got, `role="group" dir="ltr" tabindex="-1"`) && !strings.Contains(got, `tabindex="-1" `) {
		t.Errorf("group without focusable items is not a tab stop: %s", got)
	}
	if strings.Count(got, "data-disabled") != 3 {
		t.Errorf("all items disabled: %s", got)
	}
}

func TestSpacingAndErrors(t *testing.T) {
	got := h.MustRenderString(tg.ToggleGroup(tg.TypeMultiple, tg.Spacing(2), items()))
	if !strings.Contains(got, `--gap:2`) || strings.Count(got, `data-spacing="2"`) != 4 {
		t.Errorf("spacing: %s", got)
	}
	if _, err := h.RenderString(tg.ToggleGroup(items())); err == nil {
		t.Error("missing Type must fail")
	}
	if _, err := h.RenderString(tg.Item(tg.Value("x"))); err == nil {
		t.Error("item outside a group must fail")
	}
}

func TestToggleGroupOptions(t *testing.T) {
	got := h.MustRenderString(tg.ToggleGroup(tg.TypeSingle, h.Dir("rtl"), tg.Loop(false),
		tg.Item(tg.Value("a"), "A")))
	for _, want := range []string{`dir="rtl"`, `data-loop="false"`, `role="radio"`, `aria-checked="false"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, `dir="ltr"`) || strings.Contains(got, `aria-pressed`) {
		t.Errorf("unexpected attribute: %s", got)
	}
	got = h.MustRenderString(tg.ToggleGroup(tg.TypeSingle, tg.RovingFocus(false),
		tg.Item(tg.Value("a"), "A")))
	if !strings.Contains(got, `data-roving-focus="false"`) || strings.Contains(got, "tabindex") || strings.Contains(got, "data-radix-collection-item") {
		t.Errorf("rovingFocus false: %s", got)
	}
}
