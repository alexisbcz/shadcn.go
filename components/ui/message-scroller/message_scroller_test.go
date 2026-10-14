package messagescroller

import (
	"strings"
	"testing"

	h "shadcn.go/libhtml"
)

func TestMarkup(t *testing.T) {
	var b strings.Builder
	err := Provider(AutoScroll(true), PositionStart,
		MessageScroller(
			Viewport(Content(Item(MessageID("m1"), ScrollAnchor(true), "hi"))),
			Button(DirectionStart),
		),
	).Render(&b)
	if err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{`data-message-id="m1"`, `data-scroll-anchor="true"`, `data-message-scroller-spacer=""`, `role="log"`, `data-active="false"`, `inert`, `tabindex="-1"`, `Scroll to start`, `&#34;autoScroll&#34;:true`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(out, `data-pending-scroll=""`) {
		t.Error("start position must not hide the viewport")
	}
	if strings.Count(out, `data-message-scroller="`) != 4 {
		t.Errorf("parts not sharing one provider id: %s", out)
	}
	_ = h.Div
}
