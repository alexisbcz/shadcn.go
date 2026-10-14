package primitive_test

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog/primitive"
	h "shadcn.go/libhtml"
)

func attr(t *testing.T, html, re string) string {
	t.Helper()
	m := regexp.MustCompile(re).FindStringSubmatch(html)
	if m == nil {
		t.Fatalf("%s not found in %s", re, html)
	}
	return m[1]
}

func TestRootSharesScope(t *testing.T) {
	html := h.MustRenderString(primitive.Root(
		h.Div(primitive.Trigger("Open")), // through plain elements
		primitive.Portal(
			primitive.Overlay(),
			primitive.Content(
				primitive.Title("T"),
				primitive.Description("D"),
				primitive.Close("x"),
			),
		),
	))
	controls := attr(t, html, `aria-controls="([^"]+)"`)
	if got := attr(t, html, `role="dialog" id="([^"]+)"`); got != controls {
		t.Fatalf("content id %q, trigger controls %q", got, controls)
	}
	if got := attr(t, html, `data-content-id="([^"]+)"`); got != controls {
		t.Fatalf("template content id %q", got)
	}
	if attr(t, html, `aria-labelledby="([^"]+)"`) != attr(t, html, `<h2 id="([^"]+)"`) {
		t.Fatal("title id mismatch: " + html)
	}
	if attr(t, html, `aria-describedby="([^"]+)"`) != attr(t, html, `<p id="([^"]+)"`) {
		t.Fatal("description id mismatch: " + html)
	}
	if attr(t, html, `data-radix-dialog-close="([^"]+)"`) != controls {
		t.Fatal("close marker")
	}
	for _, want := range []string{
		`<button type="button" aria-haspopup="dialog" aria-expanded="false"`,
		`data-state="closed">Open</button>`,
		`<template data-parity-ignore="" data-radix-dialog="" data-content-id=`,
		`<div data-state="open" data-radix-dialog-overlay="" style="pointer-events: auto;"></div>`,
		`tabindex="-1"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in\n%s", want, html)
		}
	}
	// The trigger is outside the template, the content inside.
	if strings.Index(html, "aria-controls") > strings.Index(html, "<template") {
		t.Fatal(html)
	}
}

func TestNestedRootsHaveTheirOwnScope(t *testing.T) {
	html := h.MustRenderString(primitive.Root(
		primitive.Trigger("outer"),
		primitive.Portal(primitive.Content(
			primitive.Root(primitive.Trigger("inner"), primitive.Portal(primitive.Content("in"))),
		)),
	))
	ids := regexp.MustCompile(`aria-controls="([^"]+)"`).FindAllStringSubmatch(html, -1)
	if len(ids) != 2 || ids[0][1] == ids[1][1] {
		t.Fatalf("expected two distinct dialogs: %v\n%s", ids, html)
	}
}

func TestOptions(t *testing.T) {
	html := h.MustRenderString(primitive.Root(primitive.Modal(false), primitive.DefaultOpen,
		primitive.Trigger("t"), primitive.Portal(primitive.Overlay(), primitive.Content())))
	for _, want := range []string{`data-modal="false"`, `data-default-open=""`, `aria-expanded="true"`, `data-state="open">t</button>`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if strings.Contains(html, "data-radix-dialog-overlay") {
		t.Error("non-modal dialogs render no overlay")
	}
	html = h.MustRenderString(primitive.Root(primitive.Alert, primitive.Modal(false), primitive.Portal(primitive.Content(primitive.Cancel("c")))))
	for _, want := range []string{`role="alertdialog"`, `data-radix-dialog="alert"`, `data-radix-dialog-cancel=""`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if strings.Contains(html, `data-modal="false"`) {
		t.Error("alert dialogs are always modal")
	}
}

func TestContentOutsidePortalMountsInPlace(t *testing.T) {
	html := h.MustRenderString(primitive.Root(primitive.Content("c")))
	if !strings.Contains(html, `data-in-place=""`) {
		t.Fatal(html)
	}
	html = h.MustRenderString(primitive.Root(primitive.Portal(primitive.Content("c"))))
	if strings.Contains(html, `data-in-place`) || strings.Count(html, "<template") != 1 {
		t.Fatal(html)
	}
}

// Radix Slot with a component child: the slot's props reach the component
// as its own props, so they beat the component's defaults (React spreads
// them after its defaults) while the component's user props beat them.
func TestSlotIntoComponent(t *testing.T) {
	html := h.MustRenderString(primitive.Root(
		primitive.Trigger(h.Data("slot", "dialog-trigger"), lib.AsChild,
			button.Button(button.VariantOutline, h.Class("w-full"), "Open")),
	))
	for _, want := range []string{
		`<button class="` + h.EscapeAttr(button.Variants(button.VariantOutline, "", "w-full")) + `"`,
		` data-slot="dialog-trigger"`, ` data-variant="outline"`, ` type="button"`, ` aria-haspopup="dialog"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %s in %s", want, html)
		}
	}
	if strings.Count(html, "data-slot") != 1 {
		t.Fatal(html)
	}
	// A plain element child keeps its own props.
	html = h.MustRenderString(primitive.Close(lib.AsChild, h.A(h.Href("/x"), h.Data("slot", "mine"), h.Type("link"), "x")))
	if !strings.Contains(html, `data-slot="mine"`) || !strings.Contains(html, `type="link"`) || !strings.HasPrefix(html, "<a ") {
		t.Fatal(html)
	}
	// AsChild needs exactly one child.
	if _, err := h.RenderString(primitive.Close(lib.AsChild)); err == nil {
		t.Fatal("expected an error")
	}
}

func TestWithPropsAndUserOverrides(t *testing.T) {
	props := h.Spread(h.H2(h.Data("slot", "dialog-title"), h.Class("text-lg")), h.H2(h.Class("text-sm"), "Hi"), lib.Cn)
	html := h.MustRenderString(primitive.Title(primitive.WithProps(props), h.ID("custom")))
	if html != `<h2 class="text-sm" id="custom" data-slot="dialog-title">Hi</h2>` {
		t.Fatal(html)
	}
}

func TestEscaping(t *testing.T) {
	payload := `"><script>alert(1)</script>`
	html := h.MustRenderString(primitive.Root(
		primitive.Trigger(payload, h.Class(payload)),
		primitive.Portal(primitive.Content(primitive.Title(payload), h.Data("x", payload))),
	))
	if strings.Contains(html, "<script>") {
		t.Fatal(html)
	}
	if _, err := h.RenderString(primitive.Content(h.Attr("onclick", "alert(1)"))); !h.IsSecurityError(err) {
		t.Fatal("event handler accepted")
	}
}
