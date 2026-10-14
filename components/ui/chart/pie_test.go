package chart

import (
	"strings"
	"testing"
)

func TestStyle(t *testing.T) {
	var b strings.Builder
	if err := Style("pie-interactive", Config{{Key: "may", Label: "May", Color: "var(--chart-5)"}}, "n0nce").Render(&b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{`nonce="n0nce"`, "[data-chart=pie-interactive]", "--color-may: var(--chart-5);", ".dark [data-chart=pie-interactive]"} {
		if !strings.Contains(out, want) {
			t.Errorf("Style: missing %q in %s", want, out)
		}
	}
	if Style("pie", Config{{Key: "x", Label: "X"}}, "") != nil {
		t.Error("Style without colors: want nil")
	}
	if err := Style("bad id", nil, "").Render(&strings.Builder{}); err == nil {
		t.Error("Style with an unsafe id: want an error")
	}
}
