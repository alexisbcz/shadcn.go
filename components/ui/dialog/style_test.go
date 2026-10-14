package dialog

import (
	"strings"
	"testing"

	"shadcn.go/components/style"
	h "shadcn.go/libhtml"
)

func TestNovaScope(t *testing.T) {
	d := func() h.Node {
		return Dialog(DefaultOpen, Content(h.Dir("rtl"), Header(Title("T")), Footer(ShowCloseButton(true))))
	}
	ny := h.MustRenderString(d())
	nova := h.MustRenderString(style.Scope(style.Nova, d()))
	if !strings.Contains(ny, "rounded-lg border bg-background p-6") || strings.Contains(ny, "ring-foreground/10") {
		t.Errorf("new-york-v4 content classes changed:\n%s", ny)
	}
	for _, want := range []string{
		"bg-popover",            // cn-dialog-content
		"start-1/2",             // ui-rtl mapping of left-1/2 under dir=rtl
		`data-variant="ghost"`,  // the close button is a Button
		"end-2",                 // ui-rtl cn-dialog-close
		"cn-font-heading",       // the title
		"rounded-b-xl border-t", // the footer
	} {
		if !strings.Contains(nova, want) {
			t.Errorf("nova rendering lacks %q:\n%s", want, nova)
		}
	}
}
