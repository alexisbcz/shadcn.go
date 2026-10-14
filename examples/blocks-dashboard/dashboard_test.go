package blocksdashboard

import (
	"net/http/httptest"
	"strings"
	"testing"

	h "shadcn.go/libhtml"
)

func TestPageRenders(t *testing.T) {
	var b strings.Builder
	if err := Page(httptest.NewRequest("GET", "/", nil)).Render(&b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if n := strings.Count(out, `data-row-id=`); n != 10+len(items) {
		t.Fatalf("rows: got %d, want %d (first page plus every row in the template)", n, 10+len(items))
	}
	for _, want := range []string{"Page 1 of 7", "0 of 68 row(s) selected.", "Total Visitors", "Acme Inc."} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
	_ = h.Group
}
