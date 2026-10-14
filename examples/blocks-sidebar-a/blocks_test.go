package blockssidebara

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

// <Collapsible asChild><SidebarMenuItem> renders one <li> carrying the
// collapsible's props.
func TestCollapsibleMenuItem(t *testing.T) {
	s := render(t, Sidebar07())
	for _, want := range []string{
		`<li class="group/menu-item relative group/collapsible" data-slot="collapsible" data-sidebar="menu-item" data-state="open">`,
		`data-sidebar-collapsible="icon"`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(s, `data-collapsible="icon"`) {
		t.Error("an expanded page must not match group-has-data-[collapsible=icon]")
	}
}

func TestAllRender(t *testing.T) {
	for _, fn := range []func() h.Node{Sidebar01, Sidebar02, Sidebar03, Sidebar04, Sidebar05, Sidebar06, Sidebar07, Sidebar08} {
		render(t, fn())
	}
}
