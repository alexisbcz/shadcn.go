package avatar_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/avatar"
	h "shadcn.go/libhtml"
)

func TestAvatar(t *testing.T) {
	got := h.MustRenderString(avatar.Avatar(avatar.SizeLg,
		avatar.Image(h.Src("https://example.com/a.png"), h.Alt("a")),
		avatar.Fallback("CN"),
	))
	for _, w := range []string{`data-slot="avatar"`, `data-size="lg"`, `<template data-slot="avatar-image-template"><img`, `src="https://example.com/a.png"`, `data-slot="avatar-image"`, `aspect-square size-full`, `data-slot="avatar-fallback"`, `>CN</span>`} {
		if !strings.Contains(got, w) {
			t.Errorf("missing %s in %s", w, got)
		}
	}
	got = h.MustRenderString(avatar.Fallback(avatar.DelayMs(600), "X"))
	if !strings.HasPrefix(got, `<template data-slot="avatar-fallback-template" data-delay="600"><span`) {
		t.Errorf("delayed fallback: %s", got)
	}
	if got := h.MustRenderString(avatar.Image(h.Src("javascript:alert(1)"))); strings.Contains(got, "javascript:") {
		t.Errorf("unsafe src rendered: %s", got)
	}
	for _, w := range []string{
		h.MustRenderString(avatar.Badge()), h.MustRenderString(avatar.Group()), h.MustRenderString(avatar.GroupCount("+3")),
	} {
		if !strings.Contains(w, "data-slot=\"avatar-") {
			t.Errorf("missing slot: %s", w)
		}
	}
	if got := h.MustRenderString(avatar.Group(h.Class("-space-x-4"))); strings.Contains(got, "-space-x-2") {
		t.Errorf("className must override spacing: %s", got)
	}
}
