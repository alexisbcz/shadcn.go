package bubble

import (
	"strings"
	"testing"

	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestBubble(t *testing.T) {
	out := render(t, Bubble(VariantMuted, AlignEnd, h.Class("mt-4"), Content("<hi>")))
	for _, want := range []string{`data-slot="bubble"`, `data-variant="muted"`, `data-align="end"`, `mt-4`, `&lt;hi&gt;`, `data-slot="bubble-content"`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
	out = render(t, Content(lib.AsChild, h.Button("x")))
	if !strings.HasPrefix(out, "<button") || !strings.Contains(out, `data-slot="bubble-content"`) {
		t.Errorf("asChild: %s", out)
	}
	out = render(t, Reactions(SideTop, h.Span("👍")))
	if !strings.Contains(out, `data-side="top"`) || !strings.Contains(out, `data-align="end"`) || !strings.Contains(out, "right-3") {
		t.Errorf("reactions: %s", out)
	}
}
