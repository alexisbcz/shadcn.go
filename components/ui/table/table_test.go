package table_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/table"
	h "shadcn.go/libhtml"
)

func TestTable(t *testing.T) {
	got := h.MustRenderString(table.Table(h.Class("x"), table.Body(table.Row(table.Cell(h.Attr("colspan", "3"), "a")))))
	want := `<div class="relative w-full overflow-x-auto" data-slot="table-container"><table class="w-full caption-bottom text-sm x" data-slot="table">`
	if !strings.HasPrefix(got, want) {
		t.Fatalf("\n got %s\nwant prefix %s", got, want)
	}
	if !strings.Contains(got, `<td class="p-2 align-middle whitespace-nowrap`) || !strings.Contains(got, `colspan="3"`) {
		t.Errorf("cell: %s", got)
	}
}
