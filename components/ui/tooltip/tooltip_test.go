package tooltip_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/tooltip"
	h "shadcn.go/libhtml"
)

func TestTooltipMarkup(t *testing.T) {
	html := h.MustRenderString(tooltip.Tooltip(
		tooltip.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Hover")),
		tooltip.Content(h.P("Add to library")),
	))
	id := regexp.MustCompile(`data-radix-tooltip-trigger="([^"]+)"`).FindStringSubmatch(html)
	if id == nil {
		t.Fatal(html)
	}
	for _, want := range []string{
		`data-slot="tooltip-trigger"`, `data-state="closed"`,
		`data-radix-tooltip="" data-content-id="` + id[1] + `" data-provider-id="default" data-delay-duration="0" data-skip-delay-duration="300"`,
		`data-slot="tooltip-content"`, `data-side="top"`, `&#34;sideOffset&#34;:0`,
		`class="z-50 size-2.5 translate-y-[calc(-50%_-_2px)] rotate-45 rounded-[2px] bg-foreground fill-foreground"`,
		`<span id="` + id[1] + `" role="tooltip"`, `data-radix-tooltip-copy=""`,
		`--radix-tooltip-trigger-width: var(--radix-popper-anchor-width)`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "aria-describedby") {
		t.Error("closed trigger must not be described")
	}
}

func TestProviderAndOptions(t *testing.T) {
	html := h.MustRenderString(tooltip.Provider(tooltip.DelayDuration(500), tooltip.SkipDelayDuration(100),
		tooltip.Tooltip(tooltip.Trigger("a"), tooltip.Content("x")),
		tooltip.Tooltip(tooltip.DelayDuration(0), tooltip.DisableHoverableContent(true), tooltip.DefaultOpen, tooltip.Trigger("b"),
			tooltip.Content(h.Aria("label", "Label"), "y")),
	))
	if n := strings.Count(html, `data-provider-id="tooltip-provider-`); n != 2 {
		t.Errorf("%d tooltips in the provider", n)
	}
	for _, want := range []string{`data-delay-duration="500" data-skip-delay-duration="100"`, `data-delay-duration="0"`, `data-disable-hoverable-content=""`,
		`data-default-open=""`, `data-state="instant-open"`, `aria-describedby=`, `role="tooltip" style="position: absolute; border: 0px;`, `>Label</span>`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, `aria-label="Label"`) {
		t.Error("aria-label must move to the visually hidden copy")
	}
}
