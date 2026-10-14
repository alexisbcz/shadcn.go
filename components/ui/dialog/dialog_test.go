package dialog_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	h "shadcn.go/libhtml"
)

func TestDialogMarkup(t *testing.T) {
	html := h.MustRenderString(dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		dialog.Content(h.Class("sm:max-w-md"),
			dialog.Header(dialog.Title("T"), dialog.Description("D")),
			dialog.Footer(dialog.ShowCloseButton(true), dialog.Close(lib.AsChild, button.Button("Cancel"))),
		),
	))
	for _, want := range []string{
		`data-slot="dialog-trigger"`,
		`data-slot="dialog-overlay"`,
		`class="` + lib.Cn("fixed inset-0 z-50 bg-black/50 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:animate-in data-[state=open]:fade-in-0") + `"`,
		`data-slot="dialog-content"`,
		// className overrides sm:max-w-lg through cn.
		`sm:max-w-md"`,
		`data-slot="dialog-header"`, `data-slot="dialog-title"`, `data-slot="dialog-description"`, `data-slot="dialog-footer"`,
		`<span class="sr-only">Close</span>`,
		`class="lucide lucide-x"`,
		`>Close</button>`, // the footer's close button
	} {
		if !strings.Contains(h.EscapeAttr(html), h.EscapeAttr(want)) && !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "sm:max-w-lg") {
		t.Error("className must replace sm:max-w-lg")
	}
	// Two close buttons inside the content (corner + Cancel) plus the footer one.
	if n := strings.Count(html, "data-radix-dialog-close="); n != 3 {
		t.Errorf("%d close markers", n)
	}
}

func TestShowCloseButtonFalse(t *testing.T) {
	html := h.MustRenderString(dialog.Dialog(dialog.Content(dialog.ShowCloseButton(false), "x")))
	if strings.Contains(html, "sr-only") {
		t.Fatal(html)
	}
	html = h.MustRenderString(dialog.Dialog(dialog.Content("x").ShowCloseButton(false)))
	if strings.Contains(html, "sr-only") {
		t.Fatal(html)
	}
}
