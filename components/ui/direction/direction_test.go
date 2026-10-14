package direction

import (
	"io"
	"strings"
	"testing"

	rg "shadcn.go/components/ui/radio-group"
	h "shadcn.go/libhtml"
)

func TestProvider(t *testing.T) {
	var b strings.Builder
	err := Provider(RTL,
		rg.RadioGroup(rg.Item(rg.Value("a"))),
		rg.RadioGroup(h.Dir("ltr"), rg.Item(rg.Value("b"))),
		Provider(LTR, rg.RadioGroup(rg.Item(rg.Value("c")))),
	).Render(&b)
	if err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if n := strings.Count(out, `dir="rtl"`); n != 1 {
		t.Errorf("want 1 rtl group, got %d: %s", n, out)
	}
	if n := strings.Count(out, `dir="ltr"`); n != 2 {
		t.Errorf("want 2 ltr groups, got %d: %s", n, out)
	}
}

type receiver struct{ dir string }

func (r *receiver) SetDirectionContext(d string) { r.dir = d }
func (r *receiver) Render(w io.Writer) error     { _, err := io.WriteString(w, "["+r.dir+"]"); return err }

func TestReceiver(t *testing.T) {
	outer, inner := &receiver{}, &receiver{}
	var b strings.Builder
	if err := Provider(RTL, h.Div(outer), Provider(LTR, inner)).Render(&b); err != nil {
		t.Fatal(err)
	}
	if outer.dir != "rtl" || inner.dir != "ltr" {
		t.Errorf("receivers got %q and %q", outer.dir, inner.dir)
	}
}
