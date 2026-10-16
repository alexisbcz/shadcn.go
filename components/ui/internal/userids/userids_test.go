// Package userids tests that an id the user gives to a part that others
// reference (a menu trigger, an accordion trigger or content) replaces the
// generated id everywhere, so that ARIA references and the behavior, which
// looks the elements up by id, keep working.
package userids

import (
	"regexp"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/accordion"
	"shadcn.go/components/ui/button"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/menubar"
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

// count counts sub as a whole attribute (preceded by a space).
func count(s, sub string) int { return strings.Count(s, " "+sub) }

func TestDropdownMenuTriggerID(t *testing.T) {
	cases := map[string]h.Node{
		"on trigger": dropdownmenu.DropdownMenu(
			dropdownmenu.Trigger(lib.AsChild, h.ID("account"), button.Button("Open")),
			dropdownmenu.Content(dropdownmenu.Item("Profile")),
		),
		"on child": dropdownmenu.DropdownMenu(
			dropdownmenu.Trigger(lib.AsChild, button.Button(h.ID("account"), "Open")),
			dropdownmenu.Content(dropdownmenu.Item("Profile")),
		),
		"no aschild": dropdownmenu.DropdownMenu(
			dropdownmenu.Trigger(h.ID("account"), "Open"),
			dropdownmenu.Content(dropdownmenu.Item("Profile")),
		),
	}
	for name, n := range cases {
		t.Run(name, func(t *testing.T) {
			s := render(t, n)
			if count(s, `id="account"`) != 1 || count(s, `data-trigger-id="account"`) != 1 || count(s, `aria-labelledby="account"`) != 1 {
				t.Errorf("trigger id not adopted:\n%s", s)
			}
		})
	}
}

func TestMenubarTriggerID(t *testing.T) {
	s := render(t, menubar.Menubar(menubar.Menu(
		menubar.Trigger(h.ID("file-menu"), "File"),
		menubar.Content(menubar.Item("New")),
	)))
	if count(s, `id="file-menu"`) != 1 || count(s, `data-trigger-id="file-menu"`) != 1 {
		t.Errorf("trigger id not adopted:\n%s", s)
	}
}

func TestAccordionIDs(t *testing.T) {
	s := render(t, accordion.Accordion(accordion.TypeSingle,
		accordion.Item(accordion.Value("a"),
			accordion.Trigger(h.ID("faq-q"), "Question"),
			accordion.Content(h.ID("faq-a"), "Answer"),
		),
	))
	for _, want := range []string{`aria-controls="faq-a"`, `aria-labelledby="faq-q"`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s:\n%s", want, s)
		}
	}
	if n := len(regexp.MustCompile(` id="faq-a"`).FindAllString(s, -1)); n != 1 {
		t.Errorf("content id rendered %d times:\n%s", n, s)
	}
}
