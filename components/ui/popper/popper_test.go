package popper_test

import (
	"encoding/json"
	"strings"
	"testing"

	"shadcn.go/components/ui/popper"
	h "shadcn.go/libhtml"
)

func TestDefaultsJSON(t *testing.T) {
	var m map[string]any
	if err := json.Unmarshal([]byte(popper.Defaults().JSON()), &m); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"side": "bottom", "align": "center", "sideOffset": 0.0, "alignOffset": 0.0, "arrowPadding": 0.0,
		"avoidCollisions": true, "sticky": "partial", "hideWhenDetached": false, "updatePositionStrategy": "optimized"}
	for k, v := range want {
		if m[k] != v {
			t.Errorf("%s = %v, want %v", k, m[k], v)
		}
	}
	if _, ok := m["collisionPadding"]; ok {
		t.Error("collisionPadding must be omitted by default")
	}
}

func TestAccept(t *testing.T) {
	o := popper.Defaults()
	for _, c := range []any{popper.SideLeft, popper.AlignEnd, popper.SideOffset(8), popper.AlignOffset(-2), popper.AvoidCollisions(false),
		popper.CollisionPaddingSides{Top: 4}, popper.StickyAlways, popper.HideWhenDetached(true), popper.CollisionBoundary{"#box"}, popper.UpdateAlways} {
		if !o.Accept(c) {
			t.Fatalf("Accept(%T) = false", c)
		}
	}
	if o.Accept("text") {
		t.Fatal("Accept(string) = true")
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(o.JSON()), &m)
	if m["side"] != "left" || m["align"] != "end" || m["sideOffset"] != 8.0 || m["alignOffset"] != -2.0 || m["avoidCollisions"] != false ||
		m["sticky"] != "always" || m["hideWhenDetached"] != true || m["updatePositionStrategy"] != "always" {
		t.Fatalf("%v", m)
	}
	if cp := m["collisionPadding"].(map[string]any); cp["top"] != 4.0 || cp["left"] != 0.0 {
		t.Fatalf("collisionPadding %v", cp)
	}
	if b := m["collisionBoundary"].([]any); len(b) != 1 || b[0] != "#box" {
		t.Fatalf("collisionBoundary %v", b)
	}
}

func TestWrap(t *testing.T) {
	o := popper.Defaults()
	o.Side = popper.SideTop
	html := h.MustRenderString(popper.Wrap(o, h.Div(h.Style("color: red"), "x")))
	for _, want := range []string{
		`data-radix-popper-content-wrapper=""`,
		`style="position: fixed; left: 0px; top: 0px; transform: translate(0, -200%); min-width: max-content;"`,
		`data-side="top"`, `data-align="center"`,
		`style="color: red; animation: none;"`,
		`data-radix-popper="{`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
}

func TestArrow(t *testing.T) {
	html := h.MustRenderString(popper.Arrow(popper.SideRight, h.El("svg", h.Class("fill-x"), h.Attr("width", "20"))))
	for _, want := range []string{
		`data-radix-popper-arrow=""`,
		`left: 0px; transform-origin: 0 0; transform: translateY(50%) rotate(90deg) translateX(-50%); visibility: hidden;`,
		`width="20"`, `height="5"`, `viewBox="0 0 30 10"`, `preserveAspectRatio="none"`, `style="display: block;"`,
		`<polygon points="0,0 30,0 15,10"></polygon>`, `class="fill-x"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
}

func TestContentVars(t *testing.T) {
	if got := popper.ContentVars("tooltip"); !strings.HasPrefix(got, "--radix-tooltip-content-transform-origin: var(--radix-popper-transform-origin);") || !strings.Contains(got, "--radix-tooltip-trigger-height: var(--radix-popper-anchor-height)") {
		t.Fatal(got)
	}
}
