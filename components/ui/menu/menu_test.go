package menu_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/menu"
	"shadcn.go/components/ui/popper"
	h "shadcn.go/libhtml"
)

var kind = &menu.Kind{
	Name:            "test-menu",
	Content:         "content-class",
	SubContent:      "sub-class",
	Item:            "item-class px-2",
	CheckboxItem:    "cb-class",
	RadioItem:       "radio-class",
	Label:           "label-class",
	Separator:       "sep-class",
	Shortcut:        "sc-class",
	SubTrigger:      "st-class",
	ContentDefaults: func(o *popper.Options) { o.SideOffset = 4 },
	ContentIDs:      true,
}

// root shares a scope with its children, like a component root.
type root struct {
	h.Builder[root]
	sc *menu.Scope
}

func newRoot(children ...any) *root {
	r := &root{sc: menu.NewScope(kind)}
	r.Builder = h.NewBuilder(r, "template")
	return r.Children(children...)
}

func (r *root) render(t *testing.T) string {
	t.Helper()
	kids := r.Elem().ChildNodes()
	menu.Share(kids, r.sc)
	s, err := h.RenderString(h.Group(kids...))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func has(t *testing.T, html string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(html, w) {
			t.Errorf("missing %q in\n%s", w, html)
		}
	}
}

func TestContentTemplate(t *testing.T) {
	r := newRoot(menu.Content(h.Class("w-56"), popper.AlignStart, h.Style("color: red"), menu.Loop(true), menu.Item("A")))
	html := r.render(t)
	has(t, html,
		`<template data-parity-ignore=""`,
		`data-content-id="`+r.sc.ContentID+`"`,
		`data-trigger-id="`+r.sc.TriggerID+`"`,
		`data-loop=""`,
		`role="menu"`, `aria-orientation="vertical"`, `data-state="open"`, `data-radix-menu-content=""`,
		`dir="ltr"`, `id="`+r.sc.ContentID+`"`, `aria-labelledby="`+r.sc.TriggerID+`"`,
		`data-slot="test-menu-content"`, `class="content-class w-56"`, `tabindex="-1"`,
		`style="outline: none; color: red; --radix-test-menu-content-transform-origin: var(--radix-popper-transform-origin);`,
		`&#34;sideOffset&#34;:4`, `&#34;align&#34;:&#34;start&#34;`,
		`data-side="bottom"`, `data-align="start"`,
	)
}

func TestItems(t *testing.T) {
	html := newRoot(
		menu.Item(menu.Inset(true), menu.VariantDestructive, menu.Disabled(true), h.Class("px-4"), "Delete"),
		menu.Item(menu.TextValue("Zebra"), "Z"),
		menu.Item(lib.AsChild, h.A(h.Href("/docs"), "Docs")),
	).render(t)
	has(t, html,
		`role="menuitem" aria-disabled="true" data-disabled="" tabindex="-1" data-orientation="vertical" data-radix-collection-item=""`,
		`class="item-class px-4"`,
		`data-slot="test-menu-item" data-inset="true" data-variant="destructive"`,
		`data-radix-menu-text-value="Zebra"`,
		`<a class="item-class px-2" href="/docs" role="menuitem"`,
	)
	if strings.Contains(html, `"><a`) {
		t.Error("asChild item rendered a wrapper element")
	}
}

func TestCheckboxAndRadio(t *testing.T) {
	html := newRoot(
		menu.CheckboxItem(menu.Checked(true), menu.Indicator("✓"), "On"),
		menu.CheckboxItem(menu.Indeterminate, menu.Controlled, menu.Indicator("✓"), "Mixed"),
		menu.RadioGroup(menu.Value("b"),
			menu.RadioItem(menu.Value("a"), menu.Indicator("•"), "A"),
			menu.RadioItem(menu.Value("b"), menu.Indicator("•"), "B"),
		),
	).render(t)
	has(t, html,
		`role="menuitemcheckbox"`, `aria-checked="true"`, `data-state="checked"`,
		`aria-checked="mixed"`, `data-state="indeterminate"`, `data-radix-menu-controlled=""`,
		`data-radix-menu-indicator="1" data-state="checked"`,
		`role="group" data-slot="test-menu-radio-group" data-radix-menu-radio-group="`,
		`data-radix-menu-radio-value="b"`,
	)
	// Radio item checked state follows the group value.
	re := regexp.MustCompile(`role="menuitemradio"[^>]*aria-checked="(true|false)"[^>]*data-radix-menu-radio-value="(a|b)"`)
	m := re.FindAllStringSubmatch(html, -1)
	if len(m) != 2 || m[0][1] != "false" || m[1][1] != "true" {
		t.Fatalf("radio states %v", m)
	}
}

func TestSub(t *testing.T) {
	html := newRoot(menu.Content(menu.Sub(
		menu.SubTrigger(menu.Inset(true), "More"),
		menu.Portal(menu.SubContent(menu.Item("Inner"))),
	))).render(t)
	trig := regexp.MustCompile(`id="([^"]+)" aria-haspopup="menu" aria-expanded="false" aria-controls="([^"]+)" data-state="closed"`).FindStringSubmatch(html)
	if trig == nil {
		t.Fatalf("sub trigger markup: %s", html)
	}
	has(t, html,
		`<template data-parity-ignore="" data-radix-menu-sub="" data-content-id="`+trig[2]+`" data-trigger-id="`+trig[1]+`" data-portal=""`,
		`id="`+trig[2]+`" aria-labelledby="`+trig[1]+`"`,
		`data-slot="test-menu-sub-content"`, `data-side="right"`, `data-align="start"`,
	)
}

func TestSubRTL(t *testing.T) {
	r := newRoot(menu.Content(menu.Sub(menu.SubTrigger("More"), menu.SubContent(menu.Item("x")))))
	r.sc.Dir = "rtl"
	html := r.render(t)
	has(t, html, `dir="rtl"`, `data-side="left"`, `data-dir="rtl"`)
}

func TestEscaping(t *testing.T) {
	html := newRoot(menu.Item(menu.TextValue(`"><script>`), `<b>x</b>`)).render(t)
	if strings.Contains(html, "<script>") || strings.Contains(html, "<b>") {
		t.Fatal(html)
	}
}
