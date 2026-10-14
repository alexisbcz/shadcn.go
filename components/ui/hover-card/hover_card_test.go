package hovercard_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	hovercard "shadcn.go/components/ui/hover-card"
	h "shadcn.go/libhtml"
)

func TestHoverCardMarkup(t *testing.T) {
	html := h.MustRenderString(hovercard.HoverCard(hovercard.OpenDelay(100), hovercard.CloseDelay(50),
		hovercard.Trigger(lib.AsChild, button.Button(button.VariantLink, "@nextjs")),
		hovercard.Content(h.Class("w-80"), hovercard.SideRight, "x", hovercard.Arrow()),
	))
	for _, want := range []string{
		`<button`, `data-slot="hover-card-trigger"`, `data-state="closed"`, `data-radix-hover-card-trigger="`,
		`data-open-delay="100" data-close-delay="50"`,
		`data-slot="hover-card-content"`, `data-state="open"`, `data-side="right"`, `&#34;sideOffset&#34;:4`,
		`--radix-hover-card-content-transform-origin: var(--radix-popper-transform-origin)`,
		`data-radix-popper-arrow=""`, `left: 0px; transform-origin: 0 0;`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, " w-64") {
		t.Error("w-80 must replace w-64")
	}
}

func TestHoverCardDefaults(t *testing.T) {
	html := h.MustRenderString(hovercard.HoverCard(hovercard.DefaultOpen, hovercard.Trigger(h.Href("#x"), "link"), hovercard.Content("c")))
	for _, want := range []string{`<a data-state="open"`, `href="#x"`, `data-open-delay="700" data-close-delay="300"`, `data-default-open=""`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
}
