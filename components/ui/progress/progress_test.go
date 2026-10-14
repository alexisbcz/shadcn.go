package progress_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/progress"
	h "shadcn.go/libhtml"
)

func TestProgress(t *testing.T) {
	got := h.MustRenderString(progress.Progress(progress.Value(33.3)))
	for _, w := range []string{`role="progressbar"`, `aria-valuemax="100"`, `aria-valuemin="0"`, `data-state="indeterminate"`, `data-max="100"`, `transform:translateX(-66.7%)`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	// Like upstream, the value never reaches Radix.
	if strings.Contains(got, "aria-valuenow") || strings.Contains(got, "data-value") {
		t.Errorf("unexpected value attributes: %s", got)
	}
	if got := h.MustRenderString(progress.Progress()); !strings.Contains(got, "translateX(-100%)") {
		t.Errorf("no value: %s", got)
	}
	got = h.MustRenderString(progress.Progress(progress.Max(40)))
	if !strings.Contains(got, `aria-valuemax="40"`) {
		t.Errorf("max: %s", got)
	}
	if got := h.MustRenderString(progress.Progress(progress.Max(-1))); !strings.Contains(got, `aria-valuemax="100"`) {
		t.Errorf("invalid max must fall back to 100: %s", got)
	}
}
