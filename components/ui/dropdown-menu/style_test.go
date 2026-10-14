package dropdownmenu_test

import (
	"strings"
	"testing"

	"shadcn.go/components/style"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	h "shadcn.go/libhtml"
)

func TestDropdownMenuNova(t *testing.T) {
	html := h.MustRenderString(style.Scope(style.Nova, dropdownmenu.DropdownMenu(
		dropdownmenu.Trigger("t"),
		dropdownmenu.Content(
			dropdownmenu.Item("a"),
			dropdownmenu.CheckboxItem(dropdownmenu.Checked(true), "c"),
			dropdownmenu.Sub(dropdownmenu.SubTrigger("s"), dropdownmenu.SubContent(dropdownmenu.Item("x"))),
		),
	)))
	for _, want := range []string{
		`class="cn-menu-target cn-menu-translucent z-50 max-h-(--radix-dropdown-menu-content-available-height) w-(--radix-dropdown-menu-trigger-width) min-w-32 origin-(--radix-dropdown-menu-content-transform-origin) overflow-x-hidden overflow-y-auto rounded-lg bg-popover p-1 text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:overflow-hidden data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"`,
		`class="pointer-events-none absolute right-2 flex items-center justify-center" data-slot="dropdown-menu-checkbox-item-indicator"`,
		`&#34;align&#34;:&#34;start&#34;`,
		`cn-rtl-flip ml-auto`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "absolute left-2") {
		t.Error("new-york indicator in nova")
	}
}
