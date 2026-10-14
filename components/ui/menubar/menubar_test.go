package menubar_test

import (
	"strings"
	"testing"

	mb "shadcn.go/components/ui/menubar"
	h "shadcn.go/libhtml"
)

func TestMarkup(t *testing.T) {
	html := h.MustRenderString(mb.Menubar(mb.Loop(false), mb.Dir("rtl"), mb.DefaultValue("file"),
		mb.Menu(mb.Value("file"), mb.Trigger("File"), mb.Content(mb.Item("New"))),
		mb.Menu(mb.Trigger(mb.Disabled(true), "Edit"), mb.Content(mb.Item("Undo"))),
	))
	for _, want := range []string{
		`role="menubar"`, `data-slot="menubar"`, `tabindex="0"`, `data-orientation="horizontal"`,
		`data-radix-menubar-loop="false"`, `data-radix-menubar-dir="rtl"`,
		`role="menuitem"`, `data-slot="menubar-trigger"`, `aria-expanded="true"`, `data-value="file"`, `data-default-open=""`,
		`data-radix-menubar-content=""`, `&#34;alignOffset&#34;:-4`, `&#34;sideOffset&#34;:8`, `data-modal="false"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
}
