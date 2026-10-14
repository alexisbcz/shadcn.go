package button

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestStyleScope(t *testing.T) {
	plain := render(t, Button("Save"))
	nova := render(t, style.Scope(style.Nova, Button("Save")))
	if plain == nova {
		t.Fatal("style.Scope did not change the rendering")
	}
	if !strings.Contains(nova, "rounded-lg") || !strings.Contains(plain, "rounded-md") {
		t.Errorf("unexpected classes:\nplain %s\nnova  %s", plain, nova)
	}
	// The ui-rtl variant maps physical classes to logical ones.
	rtl := render(t, style.Scope(style.Nova, style.ScopeRTL(Button("Save"))))
	if !strings.Contains(rtl, "has-data-[icon=inline-end]:pe-2") || strings.Contains(rtl, ":pr-2") {
		t.Errorf("rtl classes: %s", rtl)
	}
	// asChild keeps the style for the slotted child.
	link := render(t, style.Scope(style.Nova, Button(lib.AsChild, h.A(h.Href("/x"), "Go"))))
	if !strings.HasPrefix(link, "<a") || !strings.Contains(link, "rounded-lg") {
		t.Errorf("asChild in scope: %s", link)
	}
}
