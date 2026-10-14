package item_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/separator"
	h "shadcn.go/libhtml"
)

func TestItem(t *testing.T) {
	got := h.MustRenderString(item.Item(item.VariantOutline, item.SizeSm, "x"))
	for _, s := range []string{`data-slot="item"`, `data-variant="outline"`, `data-size="sm"`, "border-border", "py-3"} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in %s", s, got)
		}
	}
	got = h.MustRenderString(item.Item(lib.AsChild, h.A(h.Href("#"), "link")))
	if !strings.HasPrefix(got, "<a") || !strings.Contains(got, `data-slot="item"`) || !strings.Contains(got, `href="#"`) {
		t.Fatalf("asChild: %s", got)
	}
}

func TestMediaAndParts(t *testing.T) {
	got := h.MustRenderString(item.Media(item.MediaImage))
	if !strings.Contains(got, `data-variant="image"`) || !strings.Contains(got, "overflow-hidden") {
		t.Fatalf("media: %s", got)
	}
	got = h.MustRenderString(item.Group(item.Separator()))
	if !strings.Contains(got, `role="list"`) || !strings.Contains(got, `data-slot="item-separator"`) || !strings.Contains(got, "my-0") {
		t.Fatalf("group: %s", got)
	}
	if got := h.MustRenderString(item.Description("<x>")); !strings.HasPrefix(got, "<p") || !strings.Contains(got, "&lt;x&gt;") {
		t.Fatalf("description: %s", got)
	}
}

func TestSeparatorOrientation(t *testing.T) {
	got := h.MustRenderString(item.Separator())
	if !strings.Contains(got, `data-orientation="horizontal"`) || !strings.Contains(got, `data-slot="item-separator"`) || !strings.Contains(got, "my-0") {
		t.Fatalf("separator: %s", got)
	}
	got = h.MustRenderString(item.Separator(separator.Vertical, h.Class("my-2")))
	if !strings.Contains(got, `data-orientation="vertical"`) || strings.Contains(got, "my-0") || !strings.Contains(got, "my-2") {
		t.Fatalf("vertical separator: %s", got)
	}
}
