package separator_test

import (
	"testing"

	"shadcn.go/components/ui/separator"
	h "shadcn.go/libhtml"
)

const cls = "shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px"

func TestSeparator(t *testing.T) {
	for _, tc := range []struct {
		node h.Node
		want string
	}{
		{separator.Separator(), `<div class="` + cls + `" data-orientation="horizontal" role="none" data-slot="separator"></div>`},
		{separator.Separator(separator.Vertical, h.Class("my-4")), `<div class="` + cls + ` my-4" data-orientation="vertical" role="none" data-slot="separator"></div>`},
		{separator.Separator(separator.Vertical, separator.Decorative(false)), `<div class="` + cls + `" data-orientation="vertical" aria-orientation="vertical" role="separator" data-slot="separator"></div>`},
		{separator.Separator(separator.Decorative(false)), `<div class="` + cls + `" data-orientation="horizontal" role="separator" data-slot="separator"></div>`},
		{separator.Separator(separator.Orientation("diagonal")), `<div class="` + cls + `" data-orientation="horizontal" role="none" data-slot="separator"></div>`},
	} {
		if got := h.MustRenderString(tc.node); got != tc.want {
			t.Errorf("\n got %s\nwant %s", got, tc.want)
		}
	}
}
