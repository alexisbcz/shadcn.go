package tabs_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/ui/tabs"
	h "shadcn.go/libhtml"
)

func TestTabs(t *testing.T) {
	got := h.MustRenderString(tabs.Tabs(tabs.Value("a"),
		tabs.List(tabs.Trigger(tabs.Value("a"), "A"), tabs.Trigger(tabs.Value("b"), "B")),
		tabs.Content(tabs.Value("a"), h.P("panel a")),
		tabs.Content(tabs.Value("b"), h.P("panel b")),
	))
	for _, want := range []string{
		`dir="ltr" data-orientation="horizontal" data-slot="tabs"`,
		`role="tablist"`, `aria-orientation="horizontal"`, `data-variant="default"`,
		`aria-selected="true"`, `aria-selected="false"`,
		`style="animation-duration:0s"`, ` hidden`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	// Ids: trigger ids are referenced by the panels and vice versa.
	ids := regexp.MustCompile(`id="([^"]+)-trigger-a"`).FindStringSubmatch(got)
	if ids == nil || !strings.Contains(got, `aria-controls="`+ids[1]+`-content-a"`) || !strings.Contains(got, `aria-labelledby="`+ids[1]+`-trigger-a"`) {
		t.Errorf("id references: %s", got)
	}
	// Panel b is hidden and empty, its content in the template.
	if strings.Count(got, "panel b") != 1 || strings.Count(got, "panel a") != 2 {
		t.Errorf("panel contents: %s", got)
	}
}

func TestTabsErrors(t *testing.T) {
	for _, n := range []h.Node{tabs.List(), tabs.Trigger(), tabs.Content()} {
		if _, err := h.RenderString(n); err == nil {
			t.Errorf("%T outside Tabs must fail", n)
		}
	}
}
