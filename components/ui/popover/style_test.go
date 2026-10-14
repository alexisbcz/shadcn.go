package popover_test

import (
	"strings"
	"testing"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/popover"
	h "shadcn.go/libhtml"
)

func TestPopoverNova(t *testing.T) {
	html := h.MustRenderString(style.Scope(style.Nova, popover.Popover(popover.Trigger("t"),
		popover.Content(popover.Header(popover.Title("T"), popover.Description("D"))))))
	for _, want := range []string{"rounded-lg bg-popover p-2.5", "flex flex-col gap-0.5 text-sm", `class="font-medium" data-slot="popover-title"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in %s", want, html)
		}
	}
	if strings.Contains(html, "cn-popover") {
		t.Error("unresolved token")
	}
}
