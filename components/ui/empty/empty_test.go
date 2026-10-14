package empty_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/empty"
	h "shadcn.go/libhtml"
)

func TestMediaVariants(t *testing.T) {
	got := h.MustRenderString(empty.Media(empty.MediaIcon, h.Class("size-12"), "x"))
	want := `<div class="mb-2 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 flex shrink-0 items-center justify-center rounded-lg bg-muted text-foreground [&amp;_svg:not([class*=&#39;size-&#39;])]:size-6 size-12" data-slot="empty-icon" data-variant="icon">x</div>`
	if got != want {
		t.Fatalf("\n got %s\nwant %s", got, want)
	}
	got = h.MustRenderString(empty.Media())
	if !strings.Contains(got, `data-variant="default"`) || !strings.Contains(got, "bg-transparent") {
		t.Fatal(got)
	}
	if empty.MediaVariants("") != empty.MediaVariants(empty.MediaDefault) {
		t.Fatal("empty variant must select the default")
	}
}

func TestSlotsAndClassMerge(t *testing.T) {
	got := h.MustRenderString(empty.Empty(h.Class("p-2 border"), empty.Title("T"), empty.Description("D")))
	for _, s := range []string{`data-slot="empty"`, `data-slot="empty-title"`, `data-slot="empty-description"`, "p-2", "md:p-12"} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in %s", s, got)
		}
	}
	if strings.Contains(got, " p-6 ") {
		t.Errorf("p-6 must be overridden by p-2: %s", got)
	}
}

func TestEscaping(t *testing.T) {
	got := h.MustRenderString(empty.Title(`<script>alert(1)</script>`))
	if strings.Contains(got, "<script>") {
		t.Fatal(got)
	}
}
