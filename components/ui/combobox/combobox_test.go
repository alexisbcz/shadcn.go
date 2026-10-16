package combobox

import (
	"strings"
	"testing"

	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	s, err := h.RenderString(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The root renders no element; attributes given to it wrap the parts.
func TestRootAttributes(t *testing.T) {
	plain := render(t, Combobox(Input(), Content(List(Item(ItemValue("apple"), "Apple")))))
	if !strings.HasPrefix(plain, `<div class="group/input-group`) {
		t.Errorf("root without attributes rendered an element: %.80s", plain)
	}
	s := render(t, Combobox(h.Class("w-64"), h.ID("fruit"), Input(), Content(List(Item(ItemValue("apple"), "Apple")))))
	if !strings.HasPrefix(s, `<div class="w-64" id="fruit">`) || !strings.HasSuffix(s, "</div>") {
		t.Errorf("attributes not on a wrapping div: %.120s", s)
	}
}
