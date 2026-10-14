package sheet_test

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/sheet"
	h "shadcn.go/libhtml"
)

func TestSheetSides(t *testing.T) {
	for side, want := range map[sheet.Side]string{
		sheet.SideRight:  "data-[state=open]:slide-in-from-right",
		sheet.SideLeft:   "data-[state=open]:slide-in-from-left",
		sheet.SideTop:    "data-[state=open]:slide-in-from-top",
		sheet.SideBottom: "data-[state=open]:slide-in-from-bottom",
	} {
		html := h.MustRenderString(sheet.Sheet(sheet.Content(side, h.Class("p-8"), "x")))
		if !strings.Contains(html, `class="`+h.EscapeAttr(sheet.Classes(side, "p-8"))+`"`) || !strings.Contains(html, want) {
			t.Errorf("%s: %s", side, html)
		}
	}
	if sheet.Classes(sheet.SideRight) != lib.Cn(sheet.Classes(sheet.SideRight)) {
		t.Error("classes not normalized")
	}
	// Default side is right; the close button has no data-slot.
	html := h.MustRenderString(sheet.Sheet(sheet.Trigger("o"), sheet.Content("x")))
	if !strings.Contains(html, "slide-in-from-right") || !strings.Contains(html, `<svg class="lucide lucide-x size-4"`) || !strings.Contains(html, `data-slot="sheet-trigger"`) {
		t.Fatal(html)
	}
}
