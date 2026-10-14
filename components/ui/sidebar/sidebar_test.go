package sidebar_test

import (
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/collapsible"
	"shadcn.go/components/ui/sidebar"
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

func contains(t *testing.T, html string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(html, w) {
			t.Errorf("missing %s\nin %s", w, html)
		}
	}
}

func TestProviderAndSidebar(t *testing.T) {
	html := render(t, sidebar.Provider(h.Style("--sidebar-width: 20rem"), h.Class("min-h-0"),
		sidebar.Sidebar(sidebar.Content("x")),
	))
	contains(t, html,
		`data-slot="sidebar-wrapper"`,
		`class="group/sidebar-wrapper flex w-full has-data-[variant=inset]:bg-sidebar min-h-0"`,
		`style="--sidebar-width: 16rem; --sidebar-width-icon: 3rem; --sidebar-width: 20rem"`,
		`data-sidebar-default-open="true"`,
		`data-state="expanded" data-collapsible="" data-variant="sidebar" data-side="left" data-slot="sidebar"`,
		`data-slot="sidebar-gap"`, `data-slot="sidebar-container"`,
		`data-sidebar="sidebar" data-slot="sidebar-inner"`,
		`<template data-parity-ignore="" data-sidebar-mobile="" data-sidebar-collapsible="offcanvas">`,
		`data-mobile="true"`, `data-sidebar-mobile-body=""`, `Displays the mobile sidebar.`,
	)
}

func TestCollapsedState(t *testing.T) {
	html := render(t, sidebar.Provider(sidebar.DefaultOpen(false),
		sidebar.Sidebar(sidebar.CollapsibleIcon, sidebar.VariantFloating, sidebar.SideRight, h.Class("extra")),
	))
	contains(t, html,
		`data-sidebar-default-open="false"`,
		`data-state="collapsed" data-collapsible="icon" data-variant="floating" data-side="right"`,
		`group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4)))]`,
		`right-0 group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)] p-2`,
		` extra"`,
	)
}

func TestCollapsibleNone(t *testing.T) {
	html := render(t, sidebar.Provider(sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("h-64"), "x")))
	contains(t, html, `<div class="flex w-(--sidebar-width) flex-col bg-sidebar text-sidebar-foreground h-64" data-slot="sidebar">x</div>`)
	if strings.Contains(html, "data-sidebar-mobile") {
		t.Error("collapsible none has no mobile sheet")
	}
}

func TestMenuButton(t *testing.T) {
	html := render(t, sidebar.MenuButton(sidebar.IsActive(true), sidebar.ButtonSizeLg, sidebar.ButtonVariantOutline, h.Class("p-0"), "Go"))
	contains(t, html, `data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="lg" data-active="true"`,
		`shadow-[0_0_0_1px_var(--sidebar-border)]`, `h-12 text-sm group-data-[collapsible=icon]:p-0! p-0"`)
	if strings.Contains(html, " p-2 ") {
		t.Error("cn must drop p-2 in favor of p-0")
	}

	link := render(t, sidebar.MenuButton(lib.AsChild, h.A(h.Href("#"), h.Class("x"), "Home")))
	contains(t, link, `<a class="peer/menu-button`, ` x"`, `data-active="false"`, `data-size="default"`)
}

func TestMenuButtonTooltip(t *testing.T) {
	expanded := render(t, sidebar.Provider(sidebar.Sidebar(sidebar.MenuButton(sidebar.Tooltip(h.Class("font-bold"), "Home"), "Home"))))
	contains(t, expanded, `data-sidebar-tooltip=""`, `data-radix-tooltip-trigger=`, `data-side="right"`, `font-bold`, ` hidden`)
	collapsed := render(t, sidebar.Provider(sidebar.DefaultOpen(false), sidebar.Sidebar(sidebar.MenuButton(sidebar.Tooltip("Home"), "Home"))))
	re := regexp.MustCompile(`<div[^>]*data-slot="tooltip-content"[^>]*>`)
	if m := re.FindString(collapsed); m == "" || strings.Contains(m, " hidden") {
		t.Errorf("collapsed tooltip content must be visible: %q", m)
	}
	// As a slot child (e.g. inside a CollapsibleTrigger with AsChild), the
	// tooltip template goes first inside the button.
	slotted := render(t, collapsible.Collapsible(collapsible.Trigger(lib.AsChild, sidebar.MenuButton(sidebar.Tooltip("Docs"), "Docs"))))
	contains(t, slotted, `data-slot="collapsible-trigger"`, `aria-expanded="false"`, `<template data-parity-ignore="" data-radix-tooltip=""`)
}

func TestParts(t *testing.T) {
	cases := map[string]h.Node{
		`data-sidebar="rail" data-slot="sidebar-rail" aria-label="Toggle Sidebar" tabindex="-1" title="Toggle Sidebar"`: sidebar.Rail(),
		`data-sidebar="trigger"`:                                 sidebar.Trigger(),
		`data-slot="sidebar-trigger"`:                            sidebar.Trigger(h.Class("-ml-1")),
		`size-7 -ml-1"`:                                          sidebar.Trigger(h.Class("-ml-1")),
		`data-slot="sidebar-inset"`:                              sidebar.Inset(),
		`data-slot="sidebar-input" data-sidebar="input"`:         sidebar.Input(h.Placeholder("Search")),
		`h-8 w-full bg-background shadow-none`:                   sidebar.Input(),
		`data-slot="sidebar-separator"`:                          sidebar.Separator(),
		`mx-2 w-auto bg-sidebar-border`:                          sidebar.Separator(),
		`data-slot="sidebar-menu" data-sidebar="menu"`:           sidebar.Menu(),
		`data-slot="sidebar-menu-item" data-sidebar="menu-item"`: sidebar.MenuItem(),
		`md:opacity-0`:                                           sidebar.MenuAction(sidebar.ShowOnHover(true)),
		`data-slot="sidebar-menu-sub-button" data-sidebar="menu-sub-button" data-size="sm" data-active="true"`: sidebar.MenuSubButton(sidebar.SubButtonSizeSm, sidebar.IsActive(true)),
		`text-xs group-data-[collapsible=icon]:hidden"`:                                                        sidebar.MenuSubButton(sidebar.SubButtonSizeSm),
		`data-slot="sidebar-group-label" data-sidebar="group-label">x</span>`:                                  sidebar.GroupLabel(lib.AsChild, h.Span("x")),
	}
	for want, n := range cases {
		contains(t, render(t, n), want)
	}
}

func TestMenuSkeleton(t *testing.T) {
	re := regexp.MustCompile(`--skeleton-width: (\d+)%`)
	for range 50 {
		html := render(t, sidebar.MenuSkeleton(sidebar.ShowIcon(true)))
		contains(t, html, `data-sidebar="menu-skeleton-icon"`, `class="animate-pulse bg-accent size-4 rounded-md"`)
		m := re.FindStringSubmatch(html)
		if m == nil {
			t.Fatal(html)
		}
		if w, _ := strconv.Atoi(m[1]); w < 50 || w > 89 {
			t.Fatalf("width %d%% out of [50, 90)", w)
		}
	}
	if strings.Contains(render(t, sidebar.MenuSkeleton()), "menu-skeleton-icon") {
		t.Error("no icon without ShowIcon")
	}
}

func TestEscaping(t *testing.T) {
	html := render(t, sidebar.MenuButton(sidebar.Tooltip(`<b>"x"</b>`), `<script>`))
	if strings.Contains(html, "<script>") || strings.Contains(html, "<b>") {
		t.Error(html)
	}
}

func TestCookieState(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if _, ok := sidebar.CookieState(r); ok {
		t.Error("no cookie")
	}
	r.Header.Set("Cookie", "sidebar_state=false")
	if open, ok := sidebar.CookieState(r); !ok || open {
		t.Error("collapsed")
	}
	r.Header.Set("Cookie", "sidebar_state=true")
	if open, ok := sidebar.CookieState(r); !ok || !open {
		t.Error("expanded")
	}
}
