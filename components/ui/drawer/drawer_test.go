package drawer_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/drawer"
	h "shadcn.go/libhtml"
)

func TestDrawerMarkup(t *testing.T) {
	html := h.MustRenderString(drawer.Drawer(
		drawer.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		drawer.Content(
			drawer.Header(drawer.Title("T"), drawer.Description("D")),
			drawer.Footer(drawer.Close(lib.AsChild, button.Button("Cancel"))),
		),
	))
	for _, want := range []string{
		`data-slot="drawer-trigger"`,
		`aria-haspopup="dialog"`,
		`data-vaul-root=""`,
		`data-direction="bottom"`,
		`data-slot="drawer-overlay"`,
		`data-vaul-overlay=""`,
		`data-vaul-snap-points-overlay="true"`,
		`data-vaul-drawer-direction="bottom"`,
		`data-vaul-drawer=""`,
		`data-vaul-animate="true"`,
		`data-slot="drawer-content"`,
		`role="dialog"`,
		`data-slot="drawer-header"`,
		`data-slot="drawer-title"`,
		`data-slot="drawer-description"`,
		`data-slot="drawer-footer"`,
		`data-slot="drawer-close"`,
		`mx-auto mt-4 hidden h-2 w-[100px]`,
	} {
		if !strings.Contains(html, h.EscapeAttr(want)) && !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	// The trigger controls the content.
	i := strings.Index(html, `aria-controls="`)
	id := html[i+len(`aria-controls="`):]
	id = id[:strings.IndexByte(id, '"')]
	if !strings.Contains(html, `id="`+id+`"`) || !strings.Contains(html, `data-content-id="`+id+`"`) {
		t.Errorf("content id %q not wired", id)
	}
}

func TestDrawerOptions(t *testing.T) {
	html := h.MustRenderString(drawer.Drawer(drawer.DirectionRight, drawer.Modal(false), drawer.Dismissible(false), drawer.DefaultOpen,
		drawer.Content(h.Class("extra"), "x")))
	for _, want := range []string{`data-direction="right"`, `data-modal="false"`, `data-dismissible="false"`, `data-default-open=""`, `data-vaul-drawer-direction="right"`, `data-vaul-animate="false"`, `extra`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	// Non-modal drawers render no overlay.
	if strings.Contains(html, "drawer-overlay") {
		t.Error("non-modal drawer must not render an overlay")
	}
}

func TestSnapPoints(t *testing.T) {
	html := h.MustRenderString(drawer.Drawer(drawer.SnapPoints(0.5, "355px", 1), drawer.Content("x")))
	if !strings.Contains(html, `data-snap-points="[&#34;0.5&#34;,&#34;355px&#34;,&#34;1&#34;]"`) {
		t.Error(html)
	}
	if !strings.Contains(html, `data-vaul-snap-points="true"`) {
		t.Error("snap points attribute")
	}
}

func TestClassesMerge(t *testing.T) {
	if c := drawer.Classes("z-40"); strings.Contains(c, "z-50") || !strings.Contains(c, "z-40") {
		t.Error(c)
	}
}
