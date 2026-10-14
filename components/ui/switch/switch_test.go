package switchui_test

import (
	"strings"
	"testing"

	switchui "shadcn.go/components/ui/switch"
	h "shadcn.go/libhtml"
)

func TestSwitch(t *testing.T) {
	got := h.MustRenderString(switchui.Switch(h.ID("a")))
	for _, w := range []string{`role="switch"`, `aria-checked="false"`, `data-state="unchecked"`, `data-size="default"`, `value="on"`, `data-slot="switch-thumb"`, `<input type="checkbox" aria-hidden="true"`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	got = h.MustRenderString(switchui.Switch(switchui.SizeSm, switchui.Checked(true), h.Disabled(), switchui.Name("x")))
	for _, w := range []string{`data-size="sm"`, `aria-checked="true"`, `data-state="checked"`, `data-disabled=""`, ` checked`, `name="x"`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	if strings.Count(got, `data-state="checked"`) != 2 {
		t.Errorf("thumb must follow state: %s", got)
	}
	if h.MustRenderString(switchui.Switch().Size("")) != h.MustRenderString(switchui.Switch()) {
		t.Error("empty size must select the default")
	}
}

func TestFormAttributesGoToInput(t *testing.T) {
	got := h.MustRenderString(switchui.Switch(h.FormAttr("f"), h.Name("n"), h.Value("yes"), h.Required()))
	typed := h.MustRenderString(switchui.Switch(switchui.Form("f"), switchui.Name("n"), switchui.Value("yes"), switchui.Required(true)))
	if got != typed {
		t.Errorf("attributes and typed options differ:\n%s\n%s", got, typed)
	}
	button, input, _ := strings.Cut(got, "<input")
	for _, w := range []string{`value="yes"`, `aria-required="true"`} {
		if !strings.Contains(button, w) {
			t.Errorf("button misses %s: %s", w, button)
		}
	}
	for _, w := range []string{`form="f"`, `name="n"`, `value="yes"`, ` required`} {
		if !strings.Contains(input, w) {
			t.Errorf("input misses %s: %s", w, input)
		}
	}
	for _, bad := range []string{`form=`, `name=`, ` required`} {
		if strings.Contains(button, bad) {
			t.Errorf("button must not carry %s: %s", bad, button)
		}
	}
}
