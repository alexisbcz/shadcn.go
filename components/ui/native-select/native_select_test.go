package nativeselect_test

import (
	"strings"
	"testing"

	nativeselect "shadcn.go/components/ui/native-select"
	h "shadcn.go/libhtml"
)

func TestPropsGoToSelect(t *testing.T) {
	got := h.MustRenderString(nativeselect.NativeSelect(nativeselect.SizeSm, h.Disabled(), h.Name("status"), h.Class("w-40"),
		nativeselect.Option(h.Value(""), "Pick"),
		nativeselect.OptGroup(h.LabelAttr("G"), nativeselect.Option(h.Value("a"), "A")),
	))
	for _, s := range []string{
		`<div class="group/native-select relative w-fit has-[select:disabled]:opacity-50" data-slot="native-select-wrapper"><select`,
		`data-slot="native-select" data-size="sm"`, ` disabled`, `name="status"`, `w-40`,
		`<option class="bg-[Canvas] text-[CanvasText]" data-slot="native-select-option" value="">Pick</option>`,
		`<optgroup class="bg-[Canvas] text-[CanvasText]" data-slot="native-select-optgroup" label="G">`,
		`data-slot="native-select-icon"`, `aria-hidden="true"`,
	} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q in\n%s", s, got)
		}
	}
	if strings.Contains(got, " w-full ") {
		t.Errorf("w-full must be overridden by w-40: %s", got)
	}
}

func TestDefaults(t *testing.T) {
	a := h.MustRenderString(nativeselect.NativeSelect())
	b := h.MustRenderString(nativeselect.NativeSelect().Size(""))
	if a != b || !strings.Contains(a, `data-size="default"`) {
		t.Fatalf("%s\n%s", a, b)
	}
}
