package alert_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/alert"
	h "shadcn.go/libhtml"
)

func TestAlert(t *testing.T) {
	got := h.MustRenderString(alert.Alert(alert.VariantDestructive, alert.Title("T"), alert.Description("D")))
	if !strings.Contains(got, `data-slot="alert" role="alert"`) {
		t.Fatalf("got %s", got)
	}
	// cn() dedupes [&>svg]:text-current, repeated by the destructive variant.
	if n := strings.Count(got, "[&amp;&gt;svg]:text-current"); n != 1 {
		t.Errorf("[&>svg]:text-current appears %d times: %s", n, got)
	}
	for _, s := range []string{`data-slot="alert-title"`, `data-slot="alert-description"`, "text-destructive"} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %s", s)
		}
	}
}
