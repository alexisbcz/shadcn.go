package command

import (
	"bytes"
	"strings"
	"testing"

	"shadcn.go/components/ui/dialog"
	h "shadcn.go/libhtml"
)

func renderHTML(t *testing.T, n h.Node) string {
	t.Helper()
	var b bytes.Buffer
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestFirstEnabledItemSelected(t *testing.T) {
	out := renderHTML(t, Command(
		Input(h.Placeholder("Search")),
		List(
			Empty("No results."),
			Group(Heading("G"),
				Item(Disabled(true), h.Span("A")),
				Item(h.Span("B"), Shortcut("⌘B")),
				Item(h.Span("C")),
			),
		),
	))
	for _, want := range []string{
		`aria-disabled="true"`,
		`aria-selected="true" data-disabled="false" data-selected="true" data-value="B⌘B"`,
		`cmdk-root=""`,
		`tabindex="-1"`,
		`role="combobox"`,
		`aria-label="Suggestions"`,
		`cmdk-list-sizer=""`,
		`data-cmdk-empty-template=""`,
		`data-value="G"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	if strings.Count(out, `data-selected="true"`) != 1 {
		t.Fatalf("want exactly one selected item:\n%s", out)
	}
	i := strings.Index(out, `data-selected="true"`)
	j := strings.Index(out, `data-value="B⌘B"`)
	k := strings.LastIndex(out[:j], "<div")
	if i < k || i > j+200 {
		t.Errorf("selected item is not B:\n%s", out)
	}
}

func TestDefaultValueAndOptions(t *testing.T) {
	out := renderHTML(t, Command(DefaultValue("c"), Loop(true), VimBindings(false), ShouldFilter(false),
		List(Item(Value(" c "), Keywords{"x"}, "Third"), Item("Other")),
	))
	for _, want := range []string{`data-value="c"`, `data-cmdk-loop=""`, `data-cmdk-vim-bindings="false"`, `data-cmdk-should-filter="false"`, `data-cmdk-keywords="[&#34;x&#34;]"`} {
		if !strings.Contains(out, want) && !strings.Contains(out, strings.ReplaceAll(want, "&#34;", "&quot;")) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	if !strings.Contains(out, `data-selected="true" data-value="c"`) && strings.Count(out, `data-selected="true"`) != 1 {
		t.Errorf("want c selected:\n%s", out)
	}
}

func TestTextContentEscaping(t *testing.T) {
	out := renderHTML(t, Command(List(Item(`<b>"x"&y</b>`))))
	if strings.Contains(out, "<b>") {
		t.Fatalf("item text not escaped:\n%s", out)
	}
	if got := textContent([]any{h.Span(`a "<b>" & c`), " d"}); got != `a "<b>" & c d` {
		t.Errorf("textContent = %q", got)
	}
}

func TestDialog(t *testing.T) {
	out := renderHTML(t, Dialog(dialog.DefaultOpen, Input(), List(Item("A"))))
	for _, want := range []string{"Command Palette", "Search for a command to run...", `data-radix-dialog=""`, `data-default-open=""`, "overflow-hidden p-0", "[&amp;_[cmdk-item]]:py-3"} {
		if !strings.Contains(out, want) && !strings.Contains(out, strings.ReplaceAll(want, "&amp;", "&")) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
}
