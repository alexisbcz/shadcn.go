package dropdownmenu_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	dm "shadcn.go/components/ui/dropdown-menu"
	h "shadcn.go/libhtml"
)

func TestTriggerAndContentShareIDs(t *testing.T) {
	html := h.MustRenderString(dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		dm.Content(dm.Item("A")),
	))
	tid := regexp.MustCompile(`<button[^>]* id="([^"]+)"`).FindStringSubmatch(html)
	if tid == nil {
		t.Fatal(html)
	}
	for _, want := range []string{
		`data-slot="dropdown-menu-trigger"`, `aria-haspopup="menu"`, `aria-expanded="false"`, `data-state="closed"`,
		`data-radix-dropdown-menu=""`, `data-trigger-id="` + tid[1] + `"`, `aria-labelledby="` + tid[1] + `"`,
		`data-slot="dropdown-menu-item"`, `&#34;sideOffset&#34;:4`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if strings.Contains(html, "aria-controls") {
		t.Error("closed trigger must not render aria-controls")
	}
}

func TestOptions(t *testing.T) {
	html := h.MustRenderString(dm.DropdownMenu(dm.Modal(false), dm.DefaultOpen, dm.Dir("rtl"),
		dm.Trigger(dm.Disabled(true), "Open"),
		dm.Content(dm.Sub(dm.SubTrigger("More"), dm.SubContent(dm.Item("x")))),
	))
	for _, want := range []string{`data-modal="false"`, `data-default-open=""`, `aria-expanded="true"`, `aria-controls=`,
		`data-disabled=""`, `disabled`, `data-dir="rtl"`, `data-side="left"`, `lucide-chevron-right ml-auto size-4`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
}
