package attachment

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

func TestAttachment(t *testing.T) {
	out := render(t, Attachment(Content(Title("a.pdf"))))
	for _, want := range []string{`data-state="done"`, `data-size="default"`, `data-orientation="horizontal"`, "min-w-40"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
	if out := render(t, Trigger()); !strings.Contains(out, `type="button"`) {
		t.Errorf("trigger type: %s", out)
	}
	if out := render(t, Trigger(lib.AsChild, h.A(h.Href("#x")))); strings.Contains(out, `type=`) || !strings.HasPrefix(out, "<a") {
		t.Errorf("trigger asChild: %s", out)
	}
	if out := render(t, Trigger(lib.AsChild, h.Type("submit"), h.A(h.Href("#x")))); strings.Contains(out, `type=`) {
		t.Errorf("trigger asChild drops type: %s", out)
	}
	if out := render(t, Trigger(h.Type("submit"))); !strings.Contains(out, `type="submit"`) {
		t.Errorf("trigger type override: %s", out)
	}
	if out := render(t, Action(h.Aria("label", "x"))); !strings.Contains(out, `data-slot="attachment-action"`) || !strings.Contains(out, `data-size="icon-xs"`) {
		t.Errorf("action: %s", out)
	}
}
