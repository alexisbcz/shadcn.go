package navigationmenu_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	nm "shadcn.go/components/ui/navigation-menu"
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

func TestNavigationMenuMarkup(t *testing.T) {
	html := render(t, nm.NavigationMenu(nm.UseViewport(false),
		nm.List(h.Class("flex-wrap"),
			nm.Item(nm.ItemValue("home"),
				nm.Trigger("Home"),
				nm.Content(h.Ul(nm.Link(lib.AsChild, h.A(h.Href("/x"), "X")))),
			),
			nm.Item(nm.Link(nm.Active, h.Href("/docs"), h.Class(nm.TriggerStyle()), "Docs")),
			nm.Indicator(),
		),
	))
	base := regexp.MustCompile(`id="([^"]+)-trigger-home"`).FindStringSubmatch(html)
	if base == nil {
		t.Fatal("no trigger id: " + html)
	}
	for _, want := range []string{
		`data-radix-navigation-menu="`, `aria-label="Main"`, `data-orientation="horizontal"`, `dir="ltr"`,
		`data-slot="navigation-menu"`, `data-viewport="false"`,
		`<div style="position: relative;"><ul class="group flex flex-1 list-none items-center justify-center gap-1 flex-wrap"`, `data-slot="navigation-menu-list"`,
		`<li class="relative" data-radix-navigation-menu-item="home" data-slot="navigation-menu-item">`,
		`data-state="closed" aria-expanded="false" aria-controls="` + base[1] + `-content-home"`,
		`data-slot="navigation-menu-trigger"`, `data-radix-collection-item=""`,
		`Home <svg`, `group-data-[state=open]:rotate-180`, `aria-hidden="true"`,
		`<template data-parity-ignore="" data-radix-navigation-menu-content="home"><div class="top-0`,
		`id="` + base[1] + `-content-home" aria-labelledby="` + base[1] + `-trigger-home"`,
		`data-state="open" data-slot="navigation-menu-content"`,
		`href="/x" data-slot="navigation-menu-link"`,
		`data-active="" aria-current="page"`,
		`data-radix-navigation-menu-indicator=""`, `data-slot="navigation-menu-indicator"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "navigation-menu-viewport") {
		t.Error("viewport rendered with UseViewport(false)")
	}
}

func TestNavigationMenuViewportAndOptions(t *testing.T) {
	html := render(t, nm.NavigationMenu(nm.DefaultValue("a"), nm.DelayDuration(0), nm.Dir("rtl"),
		nm.List(nm.Item(nm.ItemValue("a"), nm.Trigger(nm.Disabled, "A"), nm.Content("panel"))),
	))
	for _, want := range []string{
		`&#34;delayDuration&#34;:0`, `&#34;value&#34;:&#34;a&#34;`, `dir="rtl"`, `data-viewport="true"`,
		`disabled="" data-disabled="" data-state="open" aria-expanded="true"`,
		`<div class="absolute top-full left-0 isolate z-50 flex justify-center"><template data-parity-ignore="" data-radix-navigation-menu-viewport="">`,
		`data-slot="navigation-menu-viewport"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	// Panels in a viewport carry no data-state.
	if regexp.MustCompile(`data-state="open" data-slot="navigation-menu-content"`).MatchString(html) {
		t.Error("viewport panel has data-state")
	}
}

func TestTriggerStyle(t *testing.T) {
	if got := nm.TriggerStyle("h-10"); !strings.Contains(got, "h-10") || strings.Contains(got, "h-9") {
		t.Errorf("TriggerStyle merge: %s", got)
	}
}
