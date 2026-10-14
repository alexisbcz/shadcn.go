package buttongroup_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/label"
	h "shadcn.go/libhtml"
)

func TestOrientation(t *testing.T) {
	got := h.MustRenderString(buttongroup.ButtonGroup("x"))
	if strings.Contains(got, "data-orientation") || !strings.Contains(got, `role="group"`) || !strings.Contains(got, "rounded-l-none") {
		t.Fatalf("default group: %s", got)
	}
	got = h.MustRenderString(buttongroup.ButtonGroup(buttongroup.Vertical, h.Class("h-fit")))
	if !strings.Contains(got, `data-orientation="vertical"`) || !strings.Contains(got, "flex-col") || !strings.Contains(got, "h-fit") {
		t.Fatalf("vertical group: %s", got)
	}
}

func TestSeparator(t *testing.T) {
	got := h.MustRenderString(buttongroup.Separator(h.Class("bg-red-500")))
	for _, s := range []string{`data-slot="button-group-separator"`, `data-orientation="vertical"`, `role="none"`, "bg-red-500", "m-0!"} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in %s", s, got)
		}
	}
	if strings.Contains(got, "bg-input") || strings.Contains(got, "bg-border") {
		t.Errorf("bg classes must be overridden: %s", got)
	}
}

func TestTextAsChild(t *testing.T) {
	got := h.MustRenderString(buttongroup.Text(lib.AsChild, label.Label(h.For("url"), "https://")))
	if !strings.HasPrefix(got, "<label") || strings.Contains(got, "leading-none") || !strings.Contains(got, "bg-muted") {
		t.Fatalf("Text asChild label: %s", got)
	}
	got = h.MustRenderString(buttongroup.Text(lib.AsChild, h.Span(h.Class("px-2"), "a")))
	if !strings.Contains(got, "px-4") || !strings.Contains(got, "px-2") {
		t.Fatalf("plain child classes are concatenated like Radix Slot: %s", got)
	}
	if _, err := h.RenderString(buttongroup.Text(lib.AsChild, "a", "b")); err == nil {
		t.Fatal("AsChild with two children must fail")
	}
}
