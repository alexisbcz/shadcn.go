package inputgroup_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/button"
	inputgroup "shadcn.go/components/ui/input-group"
	h "shadcn.go/libhtml"
)

func TestButton(t *testing.T) {
	got := h.MustRenderString(inputgroup.Button("x"))
	for _, s := range []string{`type="button"`, `data-size="xs"`, `data-variant="ghost"`, "h-6", "shadow-none"} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in %s", s, got)
		}
	}
	if strings.Contains(got, "h-9") {
		t.Errorf("button size classes must give way: %s", got)
	}
	got = h.MustRenderString(inputgroup.Button(inputgroup.ButtonSizeSm, button.VariantDefault, h.Type("submit")))
	if !strings.Contains(got, `data-size="sm"`) || !strings.Contains(got, `data-variant="default"`) || !strings.Contains(got, `type="submit"`) {
		t.Fatalf("custom button: %s", got)
	}
}

func TestAddonAndControls(t *testing.T) {
	got := h.MustRenderString(inputgroup.Addon(inputgroup.AlignBlockEnd, h.Class("border-t")))
	if !strings.Contains(got, `data-align="block-end"`) || !strings.Contains(got, "border-t") || !strings.Contains(got, "order-last") {
		t.Fatalf("addon: %s", got)
	}
	got = h.MustRenderString(inputgroup.Input(h.Placeholder("p")))
	if !strings.Contains(got, `data-slot="input-group-control"`) || strings.Contains(got, "shadow-xs") || strings.Contains(got, "rounded-md") {
		t.Fatalf("input: %s", got)
	}
	got = h.MustRenderString(inputgroup.Textarea())
	if !strings.Contains(got, `data-slot="input-group-control"`) || !strings.Contains(got, "py-3") || strings.Contains(got, "py-2 ") {
		t.Fatalf("textarea: %s", got)
	}
}
