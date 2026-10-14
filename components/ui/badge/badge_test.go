package badge_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/badge"
	h "shadcn.go/libhtml"
)

func TestBadge(t *testing.T) {
	got := h.MustRenderString(badge.Badge(badge.VariantOutline, "<b>"))
	if !strings.HasPrefix(got, `<span class="`) || !strings.Contains(got, `data-slot="badge" data-variant="outline"`) || !strings.HasSuffix(got, `">&lt;b&gt;</span>`) {
		t.Fatalf("got %s", got)
	}
	if !strings.Contains(got, "border-border text-foreground") {
		t.Errorf("outline classes missing: %s", got)
	}
	got = h.MustRenderString(badge.Badge(lib.AsChild, h.A(h.Href("#"), h.Class("px-4"), "Link")))
	if !strings.HasPrefix(got, `<a `) || !strings.Contains(got, `href="#"`) || !strings.Contains(got, `data-variant="default"`) || !strings.Contains(got, " px-4") || strings.Contains(got, " px-2 ") {
		t.Errorf("asChild: %s", got)
	}
	for _, v := range badge.VariantValues {
		if c := badge.Variants(v); c == "" {
			t.Errorf("no classes for %s", v)
		}
	}
}
