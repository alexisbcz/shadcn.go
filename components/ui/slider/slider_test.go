package slider_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/slider"
	h "shadcn.go/libhtml"
)

func render(children ...any) string { return h.MustRenderString(slider.Slider(children...)) }

func TestDemo(t *testing.T) {
	got := render(slider.DefaultValue(50))
	for _, w := range []string{`dir="ltr"`, `aria-disabled="false"`, `data-orientation="horizontal"`, `--radix-slider-thumb-transform:translateX(-50%)`,
		`style="left:0%;right:50%"`, `left:calc(50% + 0px)`, `aria-valuenow="50"`, `tabindex="0"`, `data-radix-collection-item=""`, `<input style="display:none" value="50">`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	if strings.Contains(got, "aria-label") {
		t.Errorf("one thumb has no label: %s", got)
	}
}

func TestRangeAndLabels(t *testing.T) {
	got := render(slider.DefaultValue(25, 75), slider.Name("p"))
	for _, w := range []string{`aria-label="Minimum"`, `aria-label="Maximum"`, `style="left:25%;right:25%"`, `left:calc(25% + 4px)`, `left:calc(75% + -4px)`, `name="p[]"`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	if got := render(slider.DefaultValue(1, 2, 3)); !strings.Contains(got, `aria-label="Value 3 of 3"`) {
		t.Errorf("three thumbs: %s", got)
	}
}

func TestNoDefaultValue(t *testing.T) {
	got := render()
	if strings.Count(got, `role="slider"`) != 2 || !strings.Contains(got, `style="display:none"`) || strings.Count(got, "aria-valuenow") != 1 {
		t.Errorf("shadcn renders two thumbs for Radix's [min]: %s", got)
	}
}

func TestOrientationsAndOptions(t *testing.T) {
	got := render(slider.DefaultValue(30), slider.OrientationVertical, slider.Disabled(true), slider.Step(0.5), slider.Inverted(true))
	for _, w := range []string{`data-orientation="vertical"`, `translateY(50%)`, `style="top:0%;bottom:70%"`, `top:calc(30% + -12.8px)`, `aria-disabled="true"`, `data-disabled=""`, `data-step="0.5"`, `data-inverted=""`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	if strings.Contains(got, `dir=`) || strings.Contains(got, `tabindex`) {
		t.Errorf("vertical has no dir, disabled thumbs no tabindex: %s", got)
	}
	if got := render(slider.DefaultValue(30), h.Dir("rtl")); !strings.Contains(got, `dir="rtl"`) || !strings.Contains(got, `style="right:0%;left:70%"`) {
		t.Errorf("rtl: %s", got)
	}
}

func TestInvalid(t *testing.T) {
	if _, err := h.RenderString(slider.Slider(slider.Step(0))); err == nil {
		t.Error("zero step must fail")
	}
	if _, err := h.RenderString(slider.Slider(slider.Min(5), slider.Max(5))); err == nil {
		t.Error("empty range must fail")
	}
}

func TestFormAttributesAreProps(t *testing.T) {
	got := render(slider.DefaultValue(20, 60), h.Name("price"), h.FormAttr("f"), h.Disabled())
	typed := render(slider.DefaultValue(20, 60), slider.Name("price"), slider.Form("f"), slider.Disabled(true))
	if got != typed {
		t.Errorf("attributes and typed options differ:\n%s\n%s", got, typed)
	}
	root, _, _ := strings.Cut(got, ">")
	for _, bad := range []string{`name=`, `form=`, ` disabled`} {
		if strings.Contains(root, bad) {
			t.Errorf("root must not carry %s: %s", bad, root)
		}
	}
	for _, w := range []string{`aria-disabled="true"`, `name="price[]"`, `form="f"`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
}
