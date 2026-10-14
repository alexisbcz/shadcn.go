package alertdialog_test

import (
	"strings"
	"testing"

	"shadcn.go/components/icons/lucide"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

func TestAlertDialogMarkup(t *testing.T) {
	html := h.MustRenderString(alertdialog.AlertDialog(
		alertdialog.Trigger("Open"),
		alertdialog.Content(alertdialog.SizeSm,
			alertdialog.Header(alertdialog.Media(lucide.Trash2()), alertdialog.Title("T"), alertdialog.Description("D")),
			alertdialog.Footer(alertdialog.Cancel("No"), alertdialog.Action(button.VariantDestructive, "Yes")),
		),
	))
	for _, want := range []string{
		`data-radix-dialog="alert"`, `role="alertdialog"`, `data-size="sm"`,
		`data-slot="alert-dialog-trigger"`, `data-slot="alert-dialog-media"`,
		`data-slot="alert-dialog-cancel"`, `data-radix-dialog-cancel=""`, `data-variant="outline"`,
		`data-slot="alert-dialog-action"`, `data-variant="destructive"`,
		`class="` + h.EscapeAttr(button.Variants(button.VariantDestructive, "")) + `"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if strings.Contains(html, `data-slot="button"`) {
		t.Error("Action/Cancel keep their own data-slot")
	}
}

// AlertDialogAction is <Button asChild><Action className={cn(className)} /></Button>:
// Radix Slot concatenates the classes without merging them.
func TestActionClassConcatenated(t *testing.T) {
	html := h.MustRenderString(alertdialog.AlertDialog(alertdialog.Content(
		alertdialog.Action(h.Class("bg-red-500"), "Yes"),
	)))
	want := `class="` + h.EscapeAttr(button.Variants(button.VariantDefault, button.SizeDefault)+" bg-red-500") + `"`
	if !strings.Contains(html, want) {
		t.Errorf("missing %s in %s", want, html)
	}
	if !strings.Contains(html, `data-slot="alert-dialog-action"`) || !strings.Contains(html, `data-variant="default"`) {
		t.Errorf("bad attributes: %s", html)
	}
}
