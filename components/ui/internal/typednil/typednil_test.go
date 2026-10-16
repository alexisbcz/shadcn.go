// Package typednil tests that typed nils, such as the nil *Builder that
// h.If(false, component) returns, render as nothing wherever a child goes.
package typednil

import (
	"testing"

	"shadcn.go/components/ui/accordion"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/card"
	h "shadcn.go/libhtml"
)

func TestTypedNilChildren(t *testing.T) {
	var nilElement *h.Element
	cases := map[string]h.Node{
		"element child":   h.Div(h.If(false, badge.Badge("New"))),
		"fragment":        h.Group(h.If(false, badge.Badge("New"))),
		"component child": card.Card(h.If(false, badge.Badge("New"))),
		"accordion": accordion.Accordion(accordion.TypeSingle,
			accordion.Item(accordion.Value("a"), nilElement,
				accordion.Trigger("A"), accordion.Content("a", h.If(false, badge.Badge("New"))),
			),
		),
		"iff": h.Div(h.Iff(false, func() *badge.Builder { return badge.Badge("New") })),
	}
	for name, n := range cases {
		t.Run(name, func(t *testing.T) {
			s, err := h.RenderString(n)
			if err != nil {
				t.Fatal(err)
			}
			if want := "New"; len(s) > 0 && contains(s, want) {
				t.Errorf("rendered the hidden badge: %s", s)
			}
		})
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
