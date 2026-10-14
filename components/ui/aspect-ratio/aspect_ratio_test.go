package aspectratio_test

import (
	"strings"
	"testing"

	aspectratio "shadcn.go/components/ui/aspect-ratio"
	h "shadcn.go/libhtml"
)

func TestAspectRatio(t *testing.T) {
	got := h.MustRenderString(aspectratio.AspectRatio(aspectratio.Ratio(16.0/9), h.Class("rounded-lg  bg-muted"), h.Style("color:red"), "x"))
	for _, want := range []string{
		`style="position:relative;width:100%;padding-bottom:56.25%"`,
		`data-radix-aspect-ratio-wrapper=""`,
		`data-slot="aspect-ratio"`,
		`style="color:red;position:absolute;top:0;right:0;bottom:0;left:0"`,
		`>x</div></div>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in\n%s", want, got)
		}
	}
	if got := h.MustRenderString(aspectratio.AspectRatio()); !strings.Contains(got, "padding-bottom:100%") {
		t.Errorf("default ratio: %s", got)
	}
}
