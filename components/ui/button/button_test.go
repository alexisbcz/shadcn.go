package button_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// TestButtonVariantsMatchShadcn compares every variant x size (with a few
// className overrides) against cn(buttonVariants(...)) computed by the real
// class-variance-authority and tailwind-merge (see scripts/jsref).
func TestButtonVariantsMatchShadcn(t *testing.T) {
	data, err := os.ReadFile("testdata/button_ref.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct{ Variant, Size, ClassName, Want string }
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no reference cases")
	}
	for _, c := range cases {
		got := button.Variants(button.Variant(c.Variant), button.Size(c.Size), c.ClassName)
		if got != c.Want {
			t.Errorf("variant=%s size=%s className=%q\n got %s\nwant %s", c.Variant, c.Size, c.ClassName, got, c.Want)
		}
		// The rendered component must carry the same classes.
		html := h.MustRenderString(button.Button(button.Variant(c.Variant), button.Size(c.Size), h.Class(c.ClassName)))
		if !strings.Contains(html, `class="`+h.EscapeAttr(c.Want)+`"`) {
			t.Errorf("rendered button classes differ for %s/%s/%q: %s", c.Variant, c.Size, c.ClassName, html)
		}
	}
}

func TestButtonRender(t *testing.T) {
	got := h.MustRenderString(button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Submit"), h.Type("submit"), "Go"))
	want := `<button class="` + h.EscapeAttr(button.Variants(button.VariantOutline, button.SizeIcon)) +
		`" aria-label="Submit" type="submit" data-slot="button" data-variant="outline" data-size="icon">Go</button>`
	if got != want {
		t.Fatalf("\n got %s\nwant %s", got, want)
	}
}

func TestButtonBuilderEqualsChildren(t *testing.T) {
	a := h.MustRenderString(button.Button(button.VariantGhost, button.SizeSm, h.Class("w-full"), h.Disabled(), "Save"))
	b := h.MustRenderString(button.Button("Save").Variant(button.VariantGhost).Size(button.SizeSm).Class("w-full").Disabled(true))
	if a != b {
		t.Fatalf("\n%s\n%s", a, b)
	}
}

func TestButtonDefaultsAndOverrides(t *testing.T) {
	got := h.MustRenderString(button.Button("x"))
	if !strings.Contains(got, `data-variant="default" data-size="default"`) {
		t.Fatal(got)
	}
	// User props override the component's data attributes, like {...props}
	// spread after them in React.
	got = h.MustRenderString(button.Button(h.Data("slot", "custom"), "x"))
	if !strings.Contains(got, `data-slot="custom"`) || strings.Count(got, "data-slot") != 1 {
		t.Fatal(got)
	}
	// Empty values fall back to defaults.
	got = h.MustRenderString(button.Button("x").Variant("").Size(""))
	if !strings.Contains(got, `data-variant="default" data-size="default"`) {
		t.Fatal(got)
	}
}

func TestButtonAsChild(t *testing.T) {
	for _, b := range []*button.Builder{
		button.Button(lib.AsChild, h.A(h.Href("/login"), h.Class("rounded-full"), "Login")),
		button.Button(h.A(h.Href("/login"), h.Class("rounded-full"), "Login")).AsChild(),
	} {
		got := h.MustRenderString(b)
		want := `<a class="` + h.EscapeAttr(button.Variants("", "", "rounded-full")) +
			`" href="/login" data-slot="button" data-variant="default" data-size="default">Login</a>`
		if got != want {
			t.Fatalf("\n got %s\nwant %s", got, want)
		}
	}
	// A component child works too.
	got := h.MustRenderString(button.Button(lib.AsChild, button.VariantLink, button.Button("inner")))
	if !strings.HasPrefix(got, "<button") || strings.Count(got, "<button") != 1 {
		t.Fatal(got)
	}
	for _, bad := range []*button.Builder{
		button.Button(lib.AsChild),
		button.Button(lib.AsChild, "text"),
		button.Button(lib.AsChild, h.A(), h.A()),
	} {
		if _, err := h.RenderString(bad); err == nil {
			t.Error("expected AsChild error")
		}
	}
}

func TestButtonSecurity(t *testing.T) {
	payload := `"><script>alert(1)</script>`
	got := h.MustRenderString(button.Button(payload, h.Class(payload), h.Value(payload)))
	if strings.Contains(got, "<script>") {
		t.Fatal(got)
	}
	if _, err := h.RenderString(button.Button(h.Attr("onclick", "alert(1)"))); !h.IsSecurityError(err) {
		t.Fatal("event handler from plain string accepted")
	}
	got = h.MustRenderString(button.Button(lib.AsChild, h.A(h.Href("javascript:alert(1)"), "x")))
	if strings.Contains(got, "javascript:") {
		t.Fatal(got)
	}
}
