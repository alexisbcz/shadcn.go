package scrollarea_test

import (
	"strings"
	"testing"

	scrollarea "shadcn.go/components/ui/scroll-area"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	s, err := h.RenderString(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestDefault(t *testing.T) {
	s := render(t, scrollarea.ScrollArea(h.Class("h-72 relative"), h.P("x")))
	for _, want := range []string{
		`dir="ltr" data-slot="scroll-area"`,
		`style="position:relative;--radix-scroll-area-corner-width:0px;--radix-scroll-area-corner-height:0px"`,
		`style="overflow-x:hidden;overflow-y:scroll"`,
		`<div style="min-width:100%;display:table"><p>x</p></div>`,
		`<template data-slot="scroll-area-scrollbar-template"><div class="flex touch-none p-px transition-colors select-none h-full w-2.5 border-l border-l-transparent" data-orientation="vertical"`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
	for _, bad := range []string{"<style", "data-type", "data-scroll-hide-delay"} {
		if strings.Contains(s, bad) {
			t.Errorf("unexpected %q in %s", bad, s)
		}
	}
}

func TestOptions(t *testing.T) {
	s := render(t, scrollarea.ScrollArea(scrollarea.TypeAlways, scrollarea.ScrollHideDelay(200), scrollarea.Nonce("abc"), h.Style("color:red"),
		h.P("x"), scrollarea.ScrollBar(scrollarea.Horizontal, h.Class("h-4"))))
	for _, want := range []string{
		`data-type="always"`, `data-scroll-hide-delay="200"`,
		`<style nonce="abc">` + scrollarea.ViewportCSS + `</style>`,
		`overflow-x:scroll;overflow-y:scroll`,
		`;color:red"`,
		`class="flex touch-none p-px transition-colors select-none flex-col border-t border-t-transparent h-4" data-orientation="horizontal" data-slot="scroll-area-scrollbar"`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
}
