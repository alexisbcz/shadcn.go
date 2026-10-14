package accordion_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/accordion"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	s, err := h.RenderString(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func item(v string, extra ...any) *accordion.ItemBuilder {
	return accordion.Item(append([]any{accordion.Value(v), accordion.Trigger(v), accordion.Content(h.Class("x"), "body "+v)}, extra...)...)
}

func TestSingle(t *testing.T) {
	s := render(t, accordion.Accordion(accordion.Value("a"), item("a"), item("b")))
	for _, want := range []string{
		`data-slot="accordion" data-orientation="vertical"`,
		`<h3 class="flex" data-orientation="vertical" data-state="open">`,
		`aria-disabled="true"`, // open and not collapsible
		`role="region"`,
		`<div class="pt-0 pb-4 x">body a</div>`,
		`--radix-accordion-content-height:var(--radix-collapsible-content-height)`,
		`transition-duration:0s;animation-name:none`,
		`data-radix-collection-item=""`,
		`class="lucide lucide-chevron-down pointer-events-none size-4`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
	if strings.Contains(s, `data-type=`) || strings.Contains(s, `data-collapsible`) {
		t.Errorf("single non-collapsible needs no options: %s", s)
	}
	if strings.Count(s, `aria-disabled`) != 1 {
		t.Errorf("only the open item's trigger is aria-disabled: %s", s)
	}
}

func TestMultipleDisabled(t *testing.T) {
	s := render(t, accordion.Accordion(accordion.TypeMultiple, accordion.Value("a"), accordion.Value("b"),
		item("a"), item("b"), item("c", accordion.Disabled(true))))
	if strings.Count(s, `aria-expanded="true"`) != 2 {
		t.Errorf("two items must be open: %s", s)
	}
	if strings.Contains(s, "aria-disabled") {
		t.Errorf("multiple accordions are collapsible: %s", s)
	}
	for _, want := range []string{`data-type="multiple"`, `data-disabled=""`, ` disabled`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
}

func TestSingleKeepsOneValue(t *testing.T) {
	s := render(t, accordion.Accordion(accordion.Collapsible(true), accordion.Value("a"), accordion.Value("b"), item("a"), item("b")))
	if strings.Count(s, `aria-expanded="true"`) != 1 {
		t.Errorf("a single accordion opens one item: %s", s)
	}
	if !strings.Contains(s, `data-collapsible=""`) {
		t.Errorf("missing data-collapsible: %s", s)
	}
}

func TestPartsOutside(t *testing.T) {
	for _, n := range []h.Node{accordion.Item(), accordion.Trigger(), accordion.Content()} {
		if _, err := h.RenderString(n); err == nil {
			t.Errorf("%T outside its parent must fail", n)
		}
	}
}
