package hugeicons

import (
	"strings"
	"testing"

	h "shadcn.go/libhtml"
)

// The markup HugeiconsIcon renders for <HugeiconsIcon icon={Add01Icon}
// strokeWidth={2} /> (react-dom/server), less its empty class="".
const add01 = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" color="currentColor" stroke-width="2" stroke="currentColor"><path d="M12.001 5.00003V19.002" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"></path>`

func TestIcon(t *testing.T) {
	s, err := h.RenderString(Add01Icon())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(s, add01) {
		t.Errorf("got %s", s)
	}
	s, _ = h.RenderString(Add01Icon(h.Class("size-4")))
	if !strings.Contains(s, `class="size-4"`) {
		t.Errorf("class: %s", s)
	}
	if ByName["AddIcon"] == nil || ByName["Add01Icon"] == nil {
		t.Error("ByName misses AddIcon")
	}
}
