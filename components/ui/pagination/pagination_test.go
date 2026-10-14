package pagination_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/pagination"
	h "shadcn.go/libhtml"
)

func TestLink(t *testing.T) {
	got := h.MustRenderString(pagination.Link(h.Href("#"), "1"))
	if strings.Contains(got, "aria-current") || strings.Contains(got, "data-active") || !strings.Contains(got, "size-9") || !strings.Contains(got, "hover:bg-accent") {
		t.Fatalf("link: %s", got)
	}
	got = h.MustRenderString(pagination.Link(pagination.IsActive, "2"))
	if !strings.Contains(got, `aria-current="page"`) || !strings.Contains(got, `data-active="true"`) || !strings.Contains(got, "border") {
		t.Fatalf("active link: %s", got)
	}
	// isActive={false}: React renders data-active="false", no aria-current.
	got = h.MustRenderString(pagination.Link("3").Active(false))
	if strings.Contains(got, "aria-current") || !strings.Contains(got, `data-active="false"`) {
		t.Fatalf("inactive link: %s", got)
	}
}

func TestPreviousNext(t *testing.T) {
	got := h.MustRenderString(pagination.Previous(h.Href("/p/1"), "ignored"))
	if !strings.Contains(got, `aria-label="Go to previous page"`) || strings.Contains(got, "ignored") || !strings.Contains(got, "Previous") || !strings.Contains(got, "sm:pl-2.5") {
		t.Fatalf("previous: %s", got)
	}
	got = h.MustRenderString(pagination.Next(h.Aria("label", "Next")))
	if !strings.Contains(got, `aria-label="Next"`) || strings.Contains(got, "Go to next page") {
		t.Fatalf("next: %s", got)
	}
	got = h.MustRenderString(pagination.Ellipsis())
	if !strings.Contains(got, `aria-hidden="true"`) || !strings.Contains(got, "More pages") {
		t.Fatalf("ellipsis: %s", got)
	}
}
