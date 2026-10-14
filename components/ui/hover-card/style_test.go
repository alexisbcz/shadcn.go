package hovercard_test

import (
	"strings"
	"testing"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/hover-card"
	h "shadcn.go/libhtml"
)

func TestHoverCardNova(t *testing.T) {
	html := h.MustRenderString(style.Scope(style.Nova, hovercard.HoverCard(hovercard.Trigger("t"), hovercard.Content("c"))))
	want := `class="z-50 w-64 origin-(--radix-hover-card-content-transform-origin) rounded-lg bg-popover p-2.5 text-sm text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-hidden duration-100 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"`
	if !strings.Contains(html, want) {
		t.Errorf("missing nova class in %s", html)
	}
}
