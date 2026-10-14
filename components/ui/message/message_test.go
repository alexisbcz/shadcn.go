package message

import (
	"strings"
	"testing"

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

func TestMessage(t *testing.T) {
	out := render(t, Message("hi"))
	for _, want := range []string{`data-slot="message"`, `data-align="start"`, "group/message"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
	if out := render(t, Message(AlignEnd, h.Class("mt-2"))); !strings.Contains(out, `data-align="end"`) || !strings.Contains(out, "mt-2") {
		t.Errorf("align/class: %s", out)
	}
	for name, n := range map[string]h.Node{
		"message-group":   Group(),
		"message-avatar":  Avatar(),
		"message-content": Content(),
		"message-header":  Header(),
		"message-footer":  Footer(),
	} {
		if out := render(t, n); !strings.Contains(out, `data-slot="`+name+`"`) {
			t.Errorf("%s: %s", name, out)
		}
	}
	// className merges through cn: px-0 replaces px-3.
	if out := render(t, Header(h.Class("px-0"))); strings.Contains(out, "px-3 ") || !strings.Contains(out, "px-0") {
		t.Errorf("header class merge: %s", out)
	}
	// Text is escaped.
	if out := render(t, Content("<b>")); !strings.Contains(out, "&lt;b&gt;") {
		t.Errorf("escaping: %s", out)
	}
}
