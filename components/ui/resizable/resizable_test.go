package resizable

import (
	"reflect"
	"strings"
	"testing"

	h "shadcn.go/libhtml"
)

func TestIdsInRenderOrder(t *testing.T) {
	got := h.MustRenderString(PanelGroup(
		Panel(DefaultSize("50%"), "One"),
		Handle(),
		Panel(DefaultSize("50%"),
			PanelGroup(Vertical,
				Panel(DefaultSize("25%"), "Two"),
				Handle(WithHandle),
				Panel(DefaultSize("75%"), "Three"),
			),
		),
	))
	var ids []string
	for _, part := range strings.Split(got, ` id="`)[1:] {
		ids = append(ids, part[:strings.Index(part, `"`)])
	}
	want := []string{"_r_0_", "_r_1_", "_r_2_", "_r_3_", "_r_4_", "_r_5_", "_r_6_", "_r_7_"}
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("ids %v", ids)
	}
	for _, s := range []string{
		`aria-controls="_r_1_" aria-valuemax="100" aria-valuemin="0" aria-valuenow="50" aria-orientation="vertical"`,
		`aria-controls="_r_5_" aria-valuemax="100" aria-valuemin="0" aria-valuenow="25" aria-orientation="horizontal"`,
		`flex-grow: 75;`, `flex-direction: column;`, `lucide-grip-vertical`,
		`data-resizable="{&#34;defaultSize&#34;:&#34;50%&#34;}"`,
	} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in\n%s", s, got)
		}
	}
}

func TestLayout(t *testing.T) {
	cons := []constraint{{maxSize: 100}, {maxSize: 100, minSize: 30}, {maxSize: 100}}
	l := validateLayout([]float64{50, 10, 40}, cons)
	if !reflect.DeepEqual(l, []float64{30, 30, 40}) {
		t.Fatal(l)
	}
	l = adjustLayoutByDelta(-100, []float64{30, 30, 40}, cons, [2]int{0, 1}, []float64{30, 30, 40})
	if !reflect.DeepEqual(l, []float64{0, 60, 40}) {
		t.Fatal(l)
	}
	if got := round3(1.23456); got != 1.235 {
		t.Fatal(got)
	}
	if formatNumber(50) != "50" || formatNumber(58.929) != "58.929" {
		t.Fatal(formatNumber(58.929))
	}
	layout, _, ok := initialLayout([]*PanelBuilder{Panel(), Panel(DefaultSize("40")), Panel()}, nil)
	if !ok || !reflect.DeepEqual(layout, []float64{30, 40, 30}) {
		t.Fatal(layout, ok)
	}
	if _, _, ok := initialLayout([]*PanelBuilder{Panel(DefaultSize("200px"))}, nil); ok {
		t.Fatal("px sizes need the client")
	}
}

func TestPanelPropsPlacement(t *testing.T) {
	got := h.MustRenderString(PanelGroup(Panel(h.Class("p-4"), h.Style("color: red"), h.Data("x", "1"), "c")))
	if !strings.Contains(got, `<div style="max-height: 100%; max-width: 100%; flex-grow: 1; color: red;" class="p-4">c</div>`) && !strings.Contains(got, `<div class="p-4" style="max-height: 100%; max-width: 100%; flex-grow: 1; color: red;">c</div>`) {
		t.Fatal(got)
	}
	if !strings.Contains(got, `data-x="1"`) {
		t.Fatal(got)
	}
}
