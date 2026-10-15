package selectex

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// select-autocomplete (ui/apps/v4/examples/aria/select-autocomplete.tsx):
// the aria-nova Select with a React Aria Autocomplete in its popover (a
// search field filtering the options). The aria Select is not ported as a
// component: the example renders its markup, and autocomplete.js drives it
// like React Aria (popover below the trigger, 4px off, 12px from the
// viewport edge; typing filters with "contains", ignoring case and
// accents; arrows move the focused option; Enter or a click picks it).

//go:embed autocomplete.js
var autocompleteScript string

func init() {
	runtime.Module("example-select-autocomplete", autocompleteScript)
	registry.Register(registry.Example{Name: "aria--select-autocomplete", Component: "select", Render: registry.Static(AriaAutocomplete)})
}

var countries = []struct{ value, label string }{
	{"argentina", "Argentina"}, {"australia", "Australia"}, {"brazil", "Brazil"}, {"canada", "Canada"},
	{"china", "China"}, {"colombia", "Colombia"}, {"egypt", "Egypt"}, {"france", "France"},
	{"germany", "Germany"}, {"italy", "Italy"}, {"japan", "Japan"}, {"kenya", "Kenya"},
	{"mexico", "Mexico"}, {"new-zealand", "New Zealand"}, {"nigeria", "Nigeria"},
	{"south-africa", "South Africa"}, {"south-korea", "South Korea"},
	{"united-kingdom", "United Kingdom"}, {"united-states", "United States"},
}

const (
	acTriggerClass    = "flex w-full items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent py-2 pr-2 pl-2.5 text-sm whitespace-nowrap transition-colors outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 data-placeholder:text-muted-foreground data-[size=default]:h-8 data-[size=sm]:h-7 data-[size=sm]:rounded-[min(var(--radius-md),10px)] *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-1.5 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
	acContentClass    = "cn-menu-target cn-menu-translucent relative isolate z-50 w-(--trigger-width) min-w-36 origin-(--trigger-anchor-point) overflow-hidden rounded-lg bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 data-entering:animate-in data-entering:fade-in-0 data-entering:zoom-in-95 data-exiting:animate-out data-exiting:fade-out-0 data-exiting:zoom-out-95 data-[placement=bottom]:slide-in-from-top-2 data-[placement=left]:slide-in-from-right-2 data-[placement=right]:slide-in-from-left-2 data-[placement=top]:slide-in-from-bottom-2 **:data-[slot$=-item]:data-focused:bg-foreground/10"
	acInputGroupClass = "group/input-group relative flex h-8 w-full min-w-0 items-center rounded-lg border border-input transition-colors outline-none in-data-[slot=combobox-content]:focus-within:border-inherit in-data-[slot=combobox-content]:focus-within:ring-0 has-disabled:bg-input/50 has-disabled:opacity-50 has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-3 has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50 has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-3 has-[[data-slot][aria-invalid=true]]:ring-destructive/20 has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>textarea]:h-auto dark:bg-input/30 dark:has-disabled:bg-input/80 dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40 has-[>[data-align=block-end]]:[&>input]:pt-3 has-[>[data-align=block-start]]:[&>input]:pb-3 has-[>[data-align=inline-end]]:[&>input]:pr-1.5 has-[>[data-align=inline-start]]:[&>input]:pl-1.5"
	acInputClass      = "h-8 w-full min-w-0 border-input px-2.5 py-1 text-base transition-colors outline-none file:inline-flex file:h-6 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 md:text-sm dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 flex-1 rounded-none border-0 bg-transparent shadow-none ring-0 focus-visible:ring-0 disabled:bg-transparent aria-invalid:ring-0 dark:bg-transparent dark:disabled:bg-transparent [&::-webkit-search-cancel-button]:hidden"
	acAddonClass      = "flex h-auto cursor-text items-center justify-center gap-2 py-1.5 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [&>kbd]:rounded-[calc(var(--radius)-5px)] [&>svg:not([class*='size-'])]:size-4 order-first pl-2 has-[>button]:ml-[-0.3rem] has-[>kbd]:ml-[-0.15rem]"
	acItemClass       = "relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground data-focused:bg-accent data-focused:text-accent-foreground data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2"
)

func AriaAutocomplete() h.Node {
	var items []any
	for _, c := range countries {
		items = append(items, h.Div(h.Data("slot", "select-item"), h.Class(acItemClass), h.Attr("role", "option"),
			h.Attr("aria-selected", "false"), h.Data("key", c.value), h.Data("selection-mode", "single"),
			h.Span(h.Class("flex flex-1 shrink-0 gap-2 whitespace-nowrap"), c.label),
			h.Span(h.Class("pointer-events-none absolute right-2 flex size-4 items-center justify-center")),
		))
	}
	return h.Div(h.Data("slot", "select"), h.Class("w-full max-w-48"), h.Data("select-autocomplete", ""),
		h.Button(h.Data("slot", "select-trigger"), h.Data("size", "default"), h.Class(acTriggerClass), h.Type("button"),
			h.Attr("aria-haspopup", "listbox"), h.Attr("aria-expanded", "false"),
			h.Span(h.Data("slot", "select-value"), h.Class("flex flex-1 text-left data-placeholder:text-muted-foreground"), h.Data("placeholder", "true"), "Select country"),
			lucide.ChevronDown(h.Class("pointer-events-none size-4 text-muted-foreground")),
		),
		h.El("template", h.Data("select-autocomplete-popover", ""),
			h.Div(h.Data("slot", "select-content"), h.Attr("dir", "ltr"), h.Class(acContentClass), h.Data("trigger", "Select"),
				h.Attr("role", "dialog"), h.Attr("tabindex", "-1"), h.Data("placement", "bottom"),
				h.Div(h.Data("slot", "select-input-wrapper"), h.Class("p-1 pb-0"), h.Data("empty", "true"),
					h.Div(h.Data("slot", "input-group"), h.Class(acInputGroupClass), h.Attr("role", "group"),
						h.Input(h.Attr("aria-autocomplete", "list"), h.Attr("autocomplete", "off"), h.Attr("autocorrect", "off"),
							h.Attr("spellcheck", "false"), h.Attr("enterkeyhint", "go"), h.Data("slot", "select-input"), h.Class(acInputClass),
							h.Type("search"), h.Value("")),
						h.Div(h.Attr("role", "group"), h.Data("slot", "input-group-addon"), h.Data("align", "inline-start"), h.Class(acAddonClass),
							lucide.Search(h.Class("size-4 shrink-0 opacity-50"))),
					),
				),
				h.Div(h.Data("slot", "select-list"), h.Class("group/select-list max-h-[inherit] overflow-x-hidden overflow-y-auto p-0 outline-hidden"),
					h.Attr("aria-label", "Suggestions"), h.Attr("role", "listbox"), h.Attr("aria-orientation", "vertical"),
					h.Data("layout", "stack"), h.Data("orientation", "vertical"),
					h.El("section", h.Data("slot", "select-group"), h.Class("scroll-my-1 p-1"), h.Attr("role", "group"), h.Group(items...)),
					h.Div(h.Attr("role", "option"), h.Style("display: contents;"),
						h.Div(h.Data("slot", "select-empty"), h.Class("hidden w-full justify-center py-2 text-center text-sm text-muted-foreground group-data-empty/select-list:flex"), "No items found.")),
				),
			),
		),
		h.El("template", h.Data("select-autocomplete-check", ""), lucide.Check()),
	)
}
