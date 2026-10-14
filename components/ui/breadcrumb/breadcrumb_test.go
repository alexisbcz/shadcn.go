package breadcrumb_test

import (
	"strings"
	"testing"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/breadcrumb"
	h "shadcn.go/libhtml"
)

func TestLinkAsChild(t *testing.T) {
	got := h.MustRenderString(breadcrumb.Link(lib.AsChild, h.A(h.Href("/docs"), h.Class("x"), "Docs")))
	want := `<a class="transition-colors hover:text-foreground x" href="/docs" data-slot="breadcrumb-link">Docs</a>`
	if got != want {
		t.Fatalf("\n got %s\nwant %s", got, want)
	}
	if got := h.MustRenderString(breadcrumb.Link(h.Href("/"), "Home")); got != `<a class="transition-colors hover:text-foreground" data-slot="breadcrumb-link" href="/">Home</a>` {
		t.Fatal(got)
	}
	if _, err := h.RenderString(breadcrumb.Link(lib.AsChild, "text")); err == nil {
		t.Fatal("AsChild with text must fail")
	}
}

func TestSeparator(t *testing.T) {
	got := h.MustRenderString(breadcrumb.Separator())
	if !strings.Contains(got, "lucide-chevron-right") || !strings.Contains(got, `role="presentation"`) || !strings.Contains(got, `aria-hidden="true"`) {
		t.Fatal(got)
	}
	got = h.MustRenderString(breadcrumb.Separator(lucide.Slash()))
	if strings.Contains(got, "chevron") || !strings.Contains(got, "lucide-slash") {
		t.Fatal(got)
	}
}

func TestPageAndEllipsis(t *testing.T) {
	got := h.MustRenderString(breadcrumb.Page("Here"))
	want := `<span class="font-normal text-foreground" data-slot="breadcrumb-page" role="link" aria-disabled="true" aria-current="page">Here</span>`
	if got != want {
		t.Fatalf("\n got %s\nwant %s", got, want)
	}
	got = h.MustRenderString(breadcrumb.Ellipsis(h.Class("size-4")))
	if !strings.Contains(got, `class="flex items-center justify-center size-4"`) || !strings.Contains(got, `<span class="sr-only">More</span>`) {
		t.Fatal(got)
	}
}

func TestSecurity(t *testing.T) {
	got := h.MustRenderString(breadcrumb.Link(lib.AsChild, h.A(h.Href("javascript:alert(1)"), "x")))
	if strings.Contains(got, "javascript:") {
		t.Fatal(got)
	}
}

func TestEllipsisIgnoresChildren(t *testing.T) {
	// Upstream renders its own children after {...props}: the user's are dropped.
	got := h.MustRenderString(breadcrumb.Ellipsis(h.ID("e"), "extra"))
	if strings.Contains(got, "extra") || !strings.Contains(got, `id="e"`) || strings.Count(got, "More") != 1 {
		t.Fatal(got)
	}
}
