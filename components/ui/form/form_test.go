package form

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/input"
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

func TestItemWithoutError(t *testing.T) {
	out := render(t, Item(
		Label("Username"),
		Control(input.Input(h.Name("username"))),
		Description("Public name."),
		Message(),
	).ID("x"))
	for _, want := range []string{
		`data-slot="form-item"`,
		`data-slot="form-label"`, `data-error="false"`, `for="x-form-item"`,
		`data-slot="form-control"`, `id="x-form-item"`, `aria-describedby="x-form-item-description"`, `aria-invalid="false"`,
		`id="x-form-item-description"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
	if strings.Contains(out, "form-message") {
		t.Errorf("empty message rendered: %s", out)
	}
	if strings.Contains(out, `data-slot="input"`) {
		t.Errorf("control keeps the input's data-slot: %s", out)
	}
}

func TestItemWithError(t *testing.T) {
	out := render(t, Item(Error("Too <short>"),
		Label("Username"),
		Control(input.Input()),
		Message("hint"),
	).ID("x"))
	for _, want := range []string{
		`data-error="true"`, `aria-invalid="true"`,
		`aria-describedby="x-form-item-description x-form-item-message"`,
		`id="x-form-item-message"`, `Too &lt;short&gt;`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
	if strings.Contains(out, "hint") {
		t.Errorf("message children shown with an error: %s", out)
	}
}

func TestMessageChildrenWithoutError(t *testing.T) {
	out := render(t, Item(Message("hint")).ID("x"))
	if !strings.Contains(out, `<p class="text-sm text-destructive" data-slot="form-message" id="x-form-item-message">hint</p>`) &&
		!strings.Contains(out, ">hint</p>") {
		t.Errorf("message: %s", out)
	}
}

func TestPartOutsideItem(t *testing.T) {
	if _, err := h.RenderString(Label("x")); err == nil {
		t.Error("expected an error for a part outside an Item")
	}
}
