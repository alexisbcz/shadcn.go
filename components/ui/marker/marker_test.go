package marker

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

func TestMarker(t *testing.T) {
	var b strings.Builder
	if err := Marker(VariantSeparator, lib.AsChild, h.A(h.Href("#x"), Content("t"))).Render(&b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.HasPrefix(out, "<a") || !strings.Contains(out, `data-variant="separator"`) || !strings.Contains(out, "before:h-px") {
		t.Errorf("marker: %s", out)
	}
}
