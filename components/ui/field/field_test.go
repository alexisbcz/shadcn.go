package field_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/field"
	h "shadcn.go/libhtml"
)

func TestField(t *testing.T) {
	got := h.MustRenderString(field.Field(field.Horizontal, h.Data("invalid", "true"), "x"))
	for _, s := range []string{`role="group"`, `data-slot="field"`, `data-orientation="horizontal"`, "flex-row", `data-invalid="true"`} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in %s", s, got)
		}
	}
	if !strings.Contains(h.MustRenderString(field.Field()), `data-orientation="vertical"`) {
		t.Error("default orientation must be vertical")
	}
}

func TestLabelAndLegend(t *testing.T) {
	got := h.MustRenderString(field.Label(h.For("a"), h.Class("font-normal"), "A"))
	if !strings.Contains(got, `data-slot="field-label"`) || !strings.Contains(got, "font-normal") || strings.Contains(got, "font-medium") || strings.Contains(got, "leading-none") {
		t.Fatalf("label: %s", got)
	}
	got = h.MustRenderString(field.Legend(field.LegendLabel, "L"))
	if !strings.Contains(got, `data-variant="label"`) || !strings.HasPrefix(got, "<legend") {
		t.Fatalf("legend: %s", got)
	}
}

func TestSeparator(t *testing.T) {
	got := h.MustRenderString(field.Separator())
	if !strings.Contains(got, `data-content="false"`) || strings.Contains(got, "field-separator-content") {
		t.Fatalf("empty separator: %s", got)
	}
	got = h.MustRenderString(field.Separator("Or"))
	if !strings.Contains(got, `data-content="true"`) || !strings.Contains(got, `data-slot="field-separator-content">Or</span>`) {
		t.Fatalf("separator with content: %s", got)
	}
}

func TestErrors(t *testing.T) {
	if got := h.MustRenderString(field.Error()); got != "" {
		t.Fatalf("empty Error must render nothing: %q", got)
	}
	if got := h.MustRenderString(field.Errors(nil)); got != "" {
		t.Fatalf("no errors must render nothing: %q", got)
	}
	got := h.MustRenderString(field.Errors([]string{"bad", "bad"}))
	if !strings.Contains(got, `role="alert"`) || !strings.Contains(got, ">bad</div>") {
		t.Fatalf("single error: %s", got)
	}
	got = h.MustRenderString(field.Errors([]string{"a", "<b>", "a"}))
	if strings.Count(got, "<li>") != 2 || !strings.Contains(got, "&lt;b&gt;") {
		t.Fatalf("error list: %s", got)
	}
}
