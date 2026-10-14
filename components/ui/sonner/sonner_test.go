package sonner_test

import (
	"strings"
	"testing"
	"time"

	"shadcn.go/components/ui/sonner"
	h "shadcn.go/libhtml"
)

func TestToasterMarkup(t *testing.T) {
	html := h.MustRenderString(sonner.Toaster())
	for _, want := range []string{
		`<section aria-label="Notifications alt+T" tabindex="-1" aria-live="polite" aria-relevant="additions text" aria-atomic="false">`,
		`data-sonner-toaster=""`,
		`data-theme="system"`,
		`data-position="bottom-right"`,
		`data-class="toaster group"`,
		`data-sonner-icon="success"`,
		`lucide-loader-circle size-4 animate-spin`,
		`data-parity-ignore=""`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if !strings.Contains(html, h.EscapeAttr(`{"--normal-bg":"var(--popover)","--normal-text":"var(--popover-foreground)","--normal-border":"var(--border)","--border-radius":"var(--radius)"}`)) {
		t.Errorf("style: %s", html)
	}
}

func TestToasterOptions(t *testing.T) {
	html := h.MustRenderString(sonner.Toaster(sonner.ThemeDark, sonner.PositionTopCenter, sonner.RichColors, sonner.Duration(2*time.Second),
		sonner.Hotkey("altKey", "KeyN"), sonner.ContainerAriaLabel("Toasts"), h.Class("mine")))
	for _, want := range []string{`data-theme="dark"`, `data-position="top-center"`, `data-rich-colors="true"`, `data-duration="2000"`, `aria-label="Toasts alt+N"`, `data-class="mine"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
}

func TestOnClickEscapes(t *testing.T) {
	html := h.MustRenderString(h.Button(sonner.OnClick(sonner.Success(`<b>"hi"</b>`).WithDescription("d")), "x"))
	if strings.Contains(html, "<b>") {
		t.Fatal(html)
	}
	if !strings.Contains(html, "data-sonner-trigger=") {
		t.Fatal(html)
	}
}
