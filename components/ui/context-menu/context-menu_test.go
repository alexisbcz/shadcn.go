package contextmenu_test

import (
	"strings"
	"testing"

	cm "shadcn.go/components/ui/context-menu"
	h "shadcn.go/libhtml"
)

func TestMarkup(t *testing.T) {
	html := h.MustRenderString(cm.ContextMenu(
		cm.Trigger(h.Class("area"), "Right click"),
		cm.Content(cm.Label(cm.Inset(true), "L"), cm.Item("A")),
	))
	for _, want := range []string{
		`data-state="closed" data-slot="context-menu-trigger" data-radix-context-menu-trigger="`, `class="area"`,
		`data-radix-context-menu=""`, `&#34;side&#34;:&#34;right&#34;`, `&#34;sideOffset&#34;:2`, `&#34;align&#34;:&#34;start&#34;`,
		`text-foreground`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if strings.Contains(html, "aria-labelledby") || strings.Contains(html, "style=\"-webkit") {
		t.Error("context menu content has no labelledby and the trigger no style")
	}
}
