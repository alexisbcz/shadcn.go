package site

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

// The create designer: /create (app/(app)/(create)/create/page.tsx and
// components/*). The page renders the design system of the search params;
// create.js keeps the designer's state (the pickers, history, locks,
// shuffle, the preview iframe sync and the dialogs), as the React hooks do.

//go:embed static/create.js
var createScript string

func init() { runtime.Module("create", createScript) }

const createDescription = "Customize everything. Pick your component library, icons, base color, theme, fonts and create your own version of shadcn/ui."

// Create serves /create.
func (s *Site) Create(r *http.Request) (h.Node, error) {
	ds := designSystemFromQuery(r.URL.Query())
	return s.Document(r, Meta{Title: "New Project", Description: createDescription, Canonical: "/create"},
		s.AppLayout(r,
			h.Div(h.Class("relative z-10 flex min-h-0 flex-1 flex-col overflow-hidden section-soft [--customizer-width:--spacing(48)] [--gap:--spacing(4)] md:[--gap:--spacing(6)] 2xl:[--customizer-width:--spacing(56)]"),
				h.Data("create", createState(ds)),
				h.Div(h.Data("slot", "designer"), h.Class("flex min-h-0 flex-1 flex-col gap-(--gap) p-(--gap) pt-[calc(var(--gap)*0.25)] md:flex-row-reverse"),
					createPreviewFrame(ds),
					style.Scope(style.Nova, crCustomizer(ds)),
				),
				style.Scope(style.Nova, welcomeDialog(), createDialogs(ds)),
				h.Script(h.Type("application/json"), h.Data("create-config", ""), h.UnsafeJS(createClientConfig())),
			),
		),
	), nil
}

// createState is the design system the page was rendered with.
func createState(ds designSystem) string {
	b, _ := json.Marshal(map[string]any{
		"base": ds.Base, "item": ds.Item, "style": ds.Style, "baseColor": ds.BaseColor, "theme": ds.Theme,
		"chartColor": ds.ChartColor, "iconLibrary": ds.IconLibrary, "font": ds.Font, "fontHeading": ds.FontHeading,
		"radius": ds.Radius, "menuAccent": ds.MenuAccent, "menuColor": ds.MenuColor, "pointer": ds.Pointer, "rtl": ds.RTL,
		"template": "next",
	})
	return string(b)
}

// previewQuery is serializeDesignSystemSearchParams: the non-default params.
func (ds designSystem) previewQuery() string {
	def := defaultDesignSystem()
	q := url.Values{}
	if ds.presetCode() != "b0" {
		q.Set("preset", ds.presetCode())
	}
	for _, k := range []string{"iconLibrary", "style", "theme", "chartColor", "font", "fontHeading", "baseColor", "menuAccent", "menuColor", "radius"} {
		if v := *ds.field(k); v != *def.field(k) {
			q.Set(k, v)
		}
	}
	if ds.Pointer {
		q.Set("pointer", "true")
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// --- header (components/designer-actions.tsx, v0-button.tsx) ---

// designerActions is DesignerActions: rendered in the site header on the
// create page only (/typeset carries its own).
func designerActions(r *http.Request) h.Node {
	if r.URL.Path != "/create" {
		return nil
	}
	ds := designSystemFromQuery(r.URL.Query())
	return style.Scope(style.Nova,
		h.Div(h.Class("hidden items-center gap-2 group-has-data-[slot=designer]/layout:md:flex"),
			vsep(),
			v0Button(ds, button.VariantOutline),
			projectFormTrigger(""),
		),
		h.Div(h.Class("hidden items-center gap-2 group-has-data-[slot=designer]/layout:flex group-has-data-[slot=designer]/layout:md:hidden"),
			vsep(),
			v0Button(ds, button.VariantDefault),
		),
	)
}

const v0Logo = `<svg viewBox="0 0 40 20" fill="none" xmlns="http://www.w3.org/2000/svg" class="size-5" data-icon="inline-end"><path d="M23.3919 0H32.9188C36.7819 0 39.9136 3.13165 39.9136 6.99475V16.0805H36.0006V6.99475C36.0006 6.90167 35.9969 6.80925 35.9898 6.71766L26.4628 16.079C26.4949 16.08 26.5272 16.0805 26.5595 16.0805H36.0006V19.7762H26.5595C22.6964 19.7762 19.4788 16.6139 19.4788 12.7508V3.68923H23.3919V12.7508C23.3919 12.9253 23.4054 13.0977 23.4316 13.2668L33.1682 3.6995C33.0861 3.6927 33.003 3.68923 32.9188 3.68923H23.3919V0Z" fill="currentColor"></path><path d="M13.7688 19.0956L0 3.68759H5.53933L13.6231 12.7337V3.68759H17.7535V17.5746C17.7535 19.6705 15.1654 20.6584 13.7688 19.0956Z" fill="currentColor"></path></svg>`

// v0Button is V0Button: "Open in v0" (outline on desktop, default on
// mobile). NEXT_PUBLIC_APP_URL and NEXT_PUBLIC_V0_URL are unset in the
// reference build, as in a local build.
func v0Button(ds designSystem, variant button.Variant) h.Node {
	return button.Button(variant, lib.AsChild, h.Class("h-[31px] gap-1 rounded-lg"),
		h.A(h.UnsafeAttr("href", v0Href(ds)), h.Target("_blank"), h.Role("link"), h.Data("create-v0", ""),
			h.Span("Open in"), h.UnsafeHTML(v0Logo)),
	)
}

func v0Href(ds designSystem) string {
	q := "preset=" + ds.presetCode() + "&base=" + ds.Base
	return "undefined/chat/api/open?url=undefined/init/v0?" + q + "&title=New " + ds.Base + "-" + ds.Style + " project"
}

// projectFormTrigger is ProjectForm's "Get Code" trigger.
func projectFormTrigger(class string) h.Node {
	return button.Button(h.Class(class), h.Data("create-open", "project-form"), "Get Code")
}

// --- preview (components/preview.tsx, preview-switcher.tsx) ---

func createPreviewFrame(ds designSystem) h.Node {
	item := func(label, value string) h.Node {
		return button.Button(button.VariantGhost, button.SizeSm, h.Data("active", boolString(ds.Item == value)),
			h.Class("h-7 min-w-8 cursor-pointer rounded-lg px-2.5 text-xs font-medium text-muted-foreground transition-colors hover:text-foreground data-[active=true]:bg-accent data-[active=true]:text-accent-foreground"),
			h.Data("create-item", value), label)
	}
	switcher := h.Div(h.Class("dark absolute right-3 bottom-3 z-20 flex items-center gap-1 rounded-xl bg-card/90 p-1 shadow-xl backdrop-blur-xl"),
		h.Data("create-preview-switcher", ""),
		item("01", "preview-02"), item("02", "preview"),
	)
	if ds.Item != "preview" && !strings.HasPrefix(ds.Item, "preview-0") {
		switcher = nil
	}
	return h.Div(h.Class("relative flex flex-1 flex-col justify-center overflow-hidden rounded-2xl ring ring-foreground/10 md:ring-muted dark:ring-foreground/10"),
		h.Div(h.Class("relative z-0 mx-auto flex w-full flex-1 flex-col overflow-hidden"),
			h.Div(h.Class("absolute inset-0 bg-muted dark:bg-muted/30")),
			h.Iframe(h.UnsafeAttr("src", "/preview/"+ds.Base+"/"+ds.Item+ds.previewQuery()), h.Class("z-10 size-full flex-1"), h.Attr("title", "Preview"),
				h.Data("create-iframe", "")),
		),
		style.Scope(style.Luma, switcher),
	)
}

// --- welcome dialog (components/welcome-dialog.tsx) ---

const shadcnLogo = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" class="size-12"><rect width="256" height="256" fill="none"></rect><line x1="208" y1="128" x2="128" y2="208" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></line><line x1="192" y1="40" x2="40" y2="192" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></line></svg>`

// welcomeDialog opens on the first visit (create.js reads localStorage).
func welcomeDialog() h.Node {
	return h.Div(h.Data("create-dialog", "welcome"),
		dialog.Dialog(
			dialog.Content(dialog.ShowCloseButton(false), baseUILayer,
				h.Class("dialog-ring max-w-92 min-w-0 gap-0 overflow-hidden rounded-xl p-0 sm:max-w-sm dark:bg-neutral-900"),
				h.Div(h.Class("flex aspect-[2/1.2] w-full items-center justify-center rounded-t-xl bg-neutral-950 text-center text-neutral-100 sm:aspect-2/1"),
					h.Div(h.Class("font-mono text-2xl font-bold"), h.UnsafeHTML(shadcnLogo)),
				),
				dialog.Header(h.Class("gap-1 p-4"),
					dialog.Title(h.Class("text-left text-base"), "Build your own shadcn/ui"),
					dialog.Description(h.Class("text-left leading-relaxed text-foreground"),
						"Customize everything from the ground up. Pick your component library, font, color scheme, and more."),
					dialog.Description(h.Class("mt-2 text-left leading-relaxed font-medium text-foreground"),
						"Available for all major React frameworks."),
				),
				dialog.Footer(h.Class("m-0"),
					dialog.Close(lib.AsChild, button.Button(h.Class("w-full"), "Get Started")),
				),
			),
		),
	)
}

// --- crCustomizer (components/crCustomizer.tsx) ---

func crCustomizer(ds designSystem) h.Node {
	return card.Card(card.SizeSm, h.Data("create-customizer", ""),
		h.Class("dark top-24 right-12 isolate z-10 max-h-full min-h-0 w-full self-start rounded-2xl bg-card/90 backdrop-blur-xl md:w-(--customizer-width)"),
		card.Header(h.Class("hidden items-center justify-between gap-2 border-b group-data-reversed/layout:flex-row-reverse md:flex"),
			crMainMenu(),
		),
		card.Content(h.Class("no-scrollbar min-h-0 flex-1 overflow-x-auto overflow-y-hidden md:overflow-y-auto"),
			field.Group(h.Class("flex-row gap-2.5 py-px **:data-[slot=field-separator]:-mx-4 **:data-[slot=field-separator]:w-auto md:flex-col md:gap-3.25"),
				crPicker(crStylePicker(ds)),
				field.Separator(h.Class("hidden md:block")),
				crPicker(crColorPicker(ds, "baseColor", "Base Color")),
				crPicker(crColorPicker(ds, "theme", "Theme")),
				crPicker(crColorPicker(ds, "chartColor", "Chart Color")),
				field.Separator(h.Class("hidden md:block")),
				crPicker(crFontPicker(ds, "fontHeading", "Heading")),
				crPicker(crFontPicker(ds, "font", "Font")),
				field.Separator(h.Class("hidden md:block")),
				crPicker(crIconLibraryPicker(ds)),
				crPicker(crRadiusPicker(ds)),
				field.Separator(h.Class("hidden md:block")),
				crPicker(crMenuColorPicker(ds)),
				crPicker(crMenuAccentPicker(ds)),
				// BasePicker renders on mobile only (useIsMobile).
				h.Div(h.Class("contents md:hidden"), crPicker(crBasePicker(ds))),
			),
		),
		card.Footer(h.Class("flex min-w-0 gap-2 md:flex-col md:rounded-b-none md:**:[button,a]:w-full"),
			crCopyPresetButton(ds),
			button.Button(button.VariantOutline, h.Data("create-open", "open-preset"),
				h.Class(crFooterButtonClass+" max-w-20 min-w-0 flex-1 sm:max-w-none md:flex-none"),
				h.Span(h.Class("md:hidden"), "Open"), h.Span(h.Class("hidden md:inline"), "Open Preset"),
			),
			button.Button(button.VariantOutline, h.Data("create-action", "shuffle"),
				h.Class(crFooterButtonClass+" max-w-20 min-w-0 flex-1 sm:max-w-none md:flex-none"),
				h.Span(h.Class("w-full truncate text-center font-medium"), "Shuffle"),
			),
		),
		card.Footer(h.Class("-mt-3 hidden min-w-0 gap-2 md:flex md:flex-col md:**:[button,a]:w-full"),
			projectFormTrigger(""),
		),
		h.Template(h.Data("create-lock-icon", "locked"), crHugeiconMarkup("SquareLock01Icon", "size-5 text-foreground")),
		h.Template(h.Data("create-lock-icon", "unlocked"), crHugeiconMarkup("SquareUnlock01Icon", "size-5 text-foreground")),
	)
}

const crFooterButtonClass = "touch-manipulation bg-transparent! px-2! py-0! text-sm! transition-none select-none hover:bg-muted! pointer-coarse:h-10!"

// crCopyPresetButton is CopyPreset.
func crCopyPresetButton(ds designSystem) h.Node {
	label := "--preset " + ds.presetCode()
	return button.Button(button.VariantOutline, h.Attr("title", label), h.Data("create-action", "copy-preset"),
		h.Class(crFooterButtonClass+" min-w-0 flex-1 md:flex-none"),
		h.Span(h.Class("block min-w-0 truncate"), h.Data("create-preset-label", ""), label),
	)
}

// crMainMenu is MainMenu: the trigger of the designer's menu.
func crMainMenu() h.Node {
	trigger := h.Button(h.Type("button"), h.Aria("haspopup", "menu"), h.Aria("expanded", "false"), h.Data("slot", "dropdown-menu-trigger"),
		h.Data("create-menu-trigger", "main"),
		h.Class(lib.Cn(crPickerTriggerClass, "flex items-center justify-between gap-2 rounded-lg px-1.75 ring-1 ring-foreground/10 focus-visible:ring-1")),
		h.Span(h.Class("font-medium"), "Menu"),
		crHugeiconMarkup("Menu09Icon", "size-5"),
	)
	item := func(action, label, shortcut string) h.Node {
		return h.Div(h.Role("menuitem"), h.TabIndex(-1), h.Data("slot", "dropdown-menu-item"), h.Data("variant", "default"),
			h.Data("create-menu-action", action), h.Class(crPickerItemClass),
			label, h.Span(h.Data("slot", "dropdown-menu-shortcut"), h.Class(crPickerShortcutClass), h.Data("create-shortcut", shortcut), shortcut))
	}
	popup := h.Div(h.Role("menu"), h.Data("slot", "dropdown-menu-content"), h.Class(crPickerContentClass),
		h.Div(h.Role("group"), h.Data("slot", "dropdown-menu-group"),
			item("navigate", "Navigate...", "Ctrl+P"),
			item("open-preset", "Open Preset...", "O"),
			item("shuffle", "Shuffle", "R"),
			item("toggle-theme", "Light/Dark", "D"),
		),
		crPickerSeparator(),
		h.Div(h.Role("group"), h.Data("slot", "dropdown-menu-group"),
			item("undo", "Undo", "Ctrl+Z"),
			item("redo", "Redo", "Ctrl+Shift+Z"),
			crPickerSeparator(),
			item("reset", "Reset", "⇧R"),
		),
	)
	return h.Group(trigger, h.Template(h.Data("create-menu", "main"), h.Data("side", "right"), h.Data("align-offset", "-8"), popup))
}

// --- pickers (components/picker.tsx and the *-picker.tsx files) ---

const (
	crPickerTriggerClass  = "relative w-36 shrink-0 touch-manipulation rounded-xl p-3 ring-1 ring-foreground/10 select-none hover:bg-muted focus-visible:ring-foreground/50 focus-visible:outline-none disabled:opacity-50 data-popup-open:bg-muted md:w-full md:rounded-lg md:px-2.5 md:py-2"
	crPickerContentClass  = "cn-menu-target z-50 no-scrollbar max-h-(--available-height) w-[calc(var(--available-width)-(--spacing(6)))] min-w-32 origin-(--transform-origin) translate-y-2 overflow-x-hidden overflow-y-auto rounded-xl border-0 bg-neutral-950/80 p-1.5 text-neutral-100 ring-1 ring-neutral-950/80 backdrop-blur-xl outline-none md:w-52 dark:bg-neutral-800/90 dark:ring-neutral-700/50 data-closed:overflow-hidden"
	crPickerItemClass     = "group/dropdown-menu-item relative flex cursor-default items-center gap-2 rounded-lg px-2 py-1.5 text-sm font-medium outline-hidden select-none **:text-neutral-100 focus:bg-neutral-600 focus:text-neutral-100 focus:**:text-neutral-100 data-inset:pl-8 dark:focus:bg-neutral-700/80 pointer-coarse:gap-3 pointer-coarse:py-2.5 pointer-coarse:pl-3 pointer-coarse:text-base data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
	crPickerRadioClass    = "relative flex cursor-default items-center gap-2 rounded-lg py-1.5 pr-8 pl-2 text-sm font-medium outline-hidden select-none **:text-neutral-100 focus:bg-neutral-600 focus:text-neutral-100 focus:**:text-neutral-100 data-inset:pl-8 dark:focus:bg-neutral-700/80 pointer-coarse:gap-3 pointer-coarse:py-2.5 pointer-coarse:pl-3 pointer-coarse:text-base data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
	crPickerLabelClass    = "px-2 py-1.5 text-xs font-medium text-neutral-400 data-inset:pl-8"
	crPickerShortcutClass = "ml-auto text-xs tracking-widest text-neutral-400! group-focus/dropdown-menu-item:text-neutral-100"
	crIndicatorClass      = "pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 items-center justify-center select-none md:right-2.5"
	crColorDotClass       = "pointer-events-none absolute top-1/2 right-4 size-4 -translate-y-1/2 rounded-full bg-(--color) select-none md:right-2.5"
)

type crPickerItem struct {
	value, title     string
	disabled, hidden bool
}

type crPickerGroup struct {
	label string
	param string // the radio group's param, when it differs (menu)
	items []crPickerItem
}

type crPickerSpec struct {
	param, label  string
	wrapClass     string
	valueClass    string
	value         string // the checked item
	valueTitle    string
	indicator     func(value string) h.Node
	groups        []crPickerGroup
	contentClass  string
	disabled      bool
	lock          bool
	indicatorKeys []string // values to pre-render indicators for
}

func crPickerSeparator() h.Node {
	return h.Div(h.Role("separator"), h.Aria("orientation", "horizontal"), h.Data("slot", "dropdown-menu-separator"), h.Class("-mx-1.5 my-1.5 h-px bg-neutral-600 dark:bg-neutral-700"))
}

func crCheckIndicator() h.Node {
	return h.Span(h.Class("pointer-events-none absolute right-2 flex items-center justify-center"), h.Data("slot", "dropdown-menu-radio-item-indicator"),
		h.Span(h.Data("create-check", ""), h.Aria("hidden", "true"),
			ip("CheckIcon|IconCheck|Tick02Icon|CheckIcon|RiCheckLine", "size-4 pointer-coarse:size-5")),
	)
}

// crPicker renders a crPicker: the trigger with its value and indicator, the
// lock button, and the menu in a template create.js opens.
func crPicker(p crPickerSpec) h.Node {
	valueClass := "text-sm font-medium text-foreground"
	if p.valueClass != "" {
		valueClass = p.valueClass
	}
	var groups []any
	for i, g := range p.groups {
		if i > 0 {
			groups = append(groups, crPickerSeparator())
		}
		param := p.param
		if g.param != "" {
			param = g.param
		}
		var items []any
		if g.label != "" {
			items = append(items, h.Div(h.Role("presentation"), h.Data("slot", "dropdown-menu-label"), h.Class(crPickerLabelClass), g.label))
		}
		for _, it := range g.items {
			el := h.Div(h.Role("menuitemradio"), h.TabIndex(-1), h.Data("slot", "dropdown-menu-radio-item"),
				h.Data("value", it.value), h.Data("title", it.title), h.Class(crPickerRadioClass),
				crCheckIndicator(), it.title)
			if it.disabled {
				el.Children(h.Data("disabled", ""), h.Aria("disabled", "true"))
			}
			if it.hidden {
				el.Children(h.Bool("hidden"))
			}
			items = append(items, el)
		}
		groups = append(groups, h.Div(h.Role("group"), h.Data("slot", "dropdown-menu-radio-group"), h.Data("create-radio-group", param),
			h.Div(h.Role("group"), h.Data("slot", "dropdown-menu-group"), h.Group(items...))))
	}
	// Each item's indicator, for create.js to swap in.
	var indicators []any
	if p.indicator != nil {
		keys := p.indicatorKeys
		if keys == nil {
			for _, g := range p.groups {
				for _, it := range g.items {
					keys = append(keys, it.value)
				}
			}
		}
		for _, k := range keys {
			indicators = append(indicators, h.Template(h.Data("create-indicator", k), p.indicator(k)))
		}
	}
	trigger := h.Button(h.Type("button"), h.Aria("haspopup", "menu"), h.Aria("expanded", "false"), h.Data("slot", "dropdown-menu-trigger"),
		h.Data("create-menu-trigger", p.param), h.Class(crPickerTriggerClass),
		h.Div(h.Class("flex flex-col justify-start text-left"),
			h.Div(h.Class("text-xs text-muted-foreground"), p.label),
			h.Div(h.Class(valueClass), h.Data("create-picker-value", ""), p.valueTitle),
		),
		h.Div(h.Class("contents"), h.Data("create-picker-indicator", ""), crIndicatorOf(p)),
	)
	if p.disabled {
		trigger.Children(h.Bool("disabled"))
	}
	wrap := "group/picker relative"
	if p.wrapClass != "" {
		wrap = p.wrapClass
	}
	contentClass := crPickerContentClass
	if p.contentClass != "" {
		contentClass = lib.Cn(crPickerContentClass, p.contentClass)
	}
	return h.Div(h.Class(wrap), h.Data("create-picker", p.param),
		trigger,
		h.Template(h.Data("create-menu", p.param), h.Data("side", "right"),
			h.Div(h.Role("menu"), h.Data("slot", "dropdown-menu-content"), h.Class(contentClass), h.Group(groups...))),
		h.Group(indicators...),
		h.If(p.lock, crLockButton(p.param)),
	)
}

func crIndicatorOf(p crPickerSpec) h.Node {
	if p.indicator == nil {
		return nil
	}
	return p.indicator(p.value)
}

// crLockButton is LockButton.
func crLockButton(param string) h.Node {
	return h.Button(h.Type("button"), h.Attr("title", "Lock"), h.Aria("label", "Lock"), h.Data("locked", "false"),
		h.Data("create-lock", param),
		h.Class("flex size-4 cursor-pointer items-center justify-center rounded opacity-0 ring-foreground/60 transition-opacity outline-none group-focus-within/picker:opacity-100 group-hover/picker:opacity-100 focus:opacity-100 focus-visible:ring-1 data-[locked=true]:opacity-100 pointer-coarse:hidden absolute top-1/2 right-8 -translate-y-1/2"),
		crHugeiconMarkup("SquareUnlock01Icon", "size-5 text-foreground"),
	)
}

// crHugeiconMarkup renders one of the create components' hugeicons.
func crHugeiconMarkup(name, class string) h.Node {
	svg := createDataOnce().Hugeicons[name]
	return h.UnsafeHTML(strings.Replace(svg, "__CLS__", class, 1))
}

func crStylePicker(ds designSystem) crPickerSpec {
	d := createDataOnce()
	p := crPickerSpec{param: "style", label: "Style", value: ds.Style, lock: true}
	var items []crPickerItem
	icons := map[string]string{}
	for _, s := range d.Styles {
		items = append(items, crPickerItem{value: s.Name, title: s.Title})
		icons[s.Name] = s.Icon
		if s.Name == ds.Style {
			p.valueTitle = s.Title
		}
	}
	p.groups = []crPickerGroup{{items: items}}
	p.indicator = func(v string) h.Node {
		return h.Div(h.Class(crIndicatorClass), h.UnsafeHTML(icons[v]))
	}
	return p
}

// crColorPicker is BaseColorPicker, ThemePicker and ChartColorPicker.
func crColorPicker(ds designSystem, param, label string) crPickerSpec {
	d := createDataOnce()
	p := crPickerSpec{param: param, label: label, value: *ds.field(param), lock: true}
	dot := func(c *createColor, key string) h.Node {
		if c == nil {
			return nil
		}
		return h.Div(h.Class(crColorDotClass), h.Style("--color: "+c.CSSVars["dark"].get(key)))
	}
	if param == "baseColor" {
		var items []crPickerItem
		for _, c := range d.BaseColors {
			items = append(items, crPickerItem{value: c.Name, title: c.Title})
		}
		p.groups = []crPickerGroup{{items: items}}
		if c := d.baseColor(p.value); c != nil {
			p.valueTitle = c.Title
		}
		p.indicator = func(v string) h.Node { return dot(d.baseColor(v), "muted-foreground") }
		return p
	}
	// Themes (and chart colors) of the base color: the base colors, then
	// the colored themes.
	// Every base color theme is rendered; the ones of other base colors
	// stay hidden (create.js follows the base color).
	var base, colored []crPickerItem
	for _, t := range d.Themes {
		if d.isBaseColor(t.Name) {
			base = append(base, crPickerItem{value: t.Name, title: t.Title, hidden: t.Name != ds.BaseColor})
		} else {
			colored = append(colored, crPickerItem{value: t.Name, title: t.Title})
		}
	}
	p.groups = []crPickerGroup{{items: base}, {items: colored}}
	p.contentClass = "max-h-92"
	if t := d.theme(p.value); t != nil {
		p.valueTitle = t.Title
	}
	p.indicator = func(v string) h.Node {
		key := "primary"
		if d.isBaseColor(v) {
			key = "muted-foreground"
		}
		return dot(d.theme(v), key)
	}
	// Every theme, for base color changes.
	for _, t := range d.Themes {
		p.indicatorKeys = append(p.indicatorKeys, t.Name)
	}
	return p
}

func crFontPicker(ds designSystem, param, label string) crPickerSpec {
	d := createDataOnce()
	value := *ds.field(param)
	p := crPickerSpec{param: param, label: label, value: value, lock: true, contentClass: "max-h-96",
		valueClass: "line-clamp-1 max-w-[80%] truncate text-sm font-medium text-foreground"}
	body, _ := d.font(ds.Font)
	if param == "fontHeading" {
		p.groups = append(p.groups, crPickerGroup{items: []crPickerItem{{value: "inherit", title: body.Title}}})
	}
	// FONTS (lib/fonts.ts) grouped by type, in order.
	order := []string{"geist", "inter", "noto-sans", "nunito-sans", "figtree", "roboto", "raleway", "dm-sans", "public-sans", "outfit", "oxanium", "manrope", "space-grotesk", "montserrat", "ibm-plex-sans", "source-sans-3", "instrument-sans", "geist-mono", "jetbrains-mono", "noto-serif", "roboto-slab", "merriweather", "lora", "playfair-display", "eb-garamond", "instrument-serif"}
	byType := map[string]int{}
	for _, name := range order {
		f, _ := d.font(name)
		i, ok := byType[f.Type]
		if !ok {
			i = len(p.groups)
			byType[f.Type] = i
			p.groups = append(p.groups, crPickerGroup{label: strings.ToUpper(f.Type[:1]) + f.Type[1:]})
		}
		p.groups[i].items = append(p.groups[i].items, crPickerItem{value: f.Name, title: f.Title})
	}
	if f, ok := d.font(value); ok {
		p.valueTitle = f.Title
	} else {
		p.valueTitle = body.Title
	}
	p.indicator = func(v string) h.Node {
		f, ok := d.font(v)
		if !ok {
			f = body
		}
		return h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 items-center justify-center text-base text-foreground select-none md:right-2.5"),
			h.UnsafeAttr("style", "font-family: "+f.fontFamily()+";"), "Aa")
	}
	return p
}

func crIconLibraryPicker(ds designSystem) crPickerSpec {
	d := createDataOnce()
	p := crPickerSpec{param: "iconLibrary", label: "Icon Library", value: ds.IconLibrary, lock: true}
	var items []crPickerItem
	for _, l := range d.IconLibraries {
		items = append(items, crPickerItem{value: l.Name, title: l.Title})
		if l.Name == ds.IconLibrary {
			p.valueTitle = l.Title
		}
	}
	p.groups = []crPickerGroup{{items: items}}
	p.indicator = func(v string) h.Node {
		return h.Div(h.Class(crIndicatorClass+" text-base text-foreground *:[svg]:text-foreground!"), h.UnsafeHTML(d.IconLibraryLogos[v]))
	}
	return p
}

const crRadiusIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" class="text-foreground"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 20v-5C4 8.925 8.925 4 15 4h5"></path></svg>`

func crRadiusPicker(ds designSystem) crPickerSpec {
	d := createDataOnce()
	locked := ds.Style == "lyra" || ds.Style == "sera"
	value := ds.Radius
	if locked {
		value = "none"
	}
	p := crPickerSpec{param: "radius", label: "Radius", value: value, lock: true, disabled: locked}
	var others []crPickerItem
	for _, r := range d.Radii {
		if r.Name == value {
			p.valueTitle = r.Label
		}
		if r.Name != "default" {
			others = append(others, crPickerItem{value: r.Name, title: r.Label, disabled: ds.Style == "rhea" && r.Name == "large"})
		}
	}
	p.groups = []crPickerGroup{{items: []crPickerItem{{value: "default", title: "Default"}}}, {items: others}}
	p.indicator = func(string) h.Node {
		return h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 rotate-90 items-center justify-center text-base text-foreground select-none md:right-2.5"), h.UnsafeHTML(crRadiusIcon))
	}
	p.indicatorKeys = []string{}
	return p
}

var crMenuColorLabels = map[string]string{
	"default": "Default / Solid", "default-translucent": "Default / Translucent",
	"inverted": "Inverted / Solid", "inverted-translucent": "Inverted / Translucent",
}

func crMenuColorPicker(ds designSystem) crPickerSpec {
	color, surface := "default", "solid"
	if strings.HasPrefix(ds.MenuColor, "inverted") {
		color = "inverted"
	}
	if strings.HasSuffix(ds.MenuColor, "translucent") {
		surface = "translucent"
	}
	p := crPickerSpec{param: "menuColor", label: "Menu", value: color + " " + surface, valueTitle: crMenuColorLabels[ds.MenuColor], lock: true,
		valueClass: "line-clamp-1 max-w-[80%] truncate text-sm font-medium text-foreground"}
	p.groups = []crPickerGroup{
		{label: "Color", param: "menuColor:color", items: []crPickerItem{{value: "default", title: "Default"}, {value: "inverted", title: "Inverted"}}},
		{label: "Appearance", param: "menuColor:surface", items: []crPickerItem{{value: "solid", title: "Solid"}, {value: "translucent", title: "Translucent"}}},
	}
	p.indicator = func(string) h.Node {
		return h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 items-center justify-center text-base text-foreground select-none md:right-2.5"),
			crHugeiconMarkup("Menu02Icon", "size-4"))
	}
	p.indicatorKeys = []string{}
	return p
}

func crMenuAccentPicker(ds designSystem) crPickerSpec {
	translucent := strings.HasSuffix(ds.MenuColor, "translucent")
	p := crPickerSpec{param: "menuAccent", label: "Menu Accent", value: ds.MenuAccent, lock: true, wrapClass: "group/picker relative pr-3 md:pr-0"}
	p.groups = []crPickerGroup{{items: []crPickerItem{{value: "subtle", title: "Subtle"}, {value: "bold", title: "Bold", disabled: translucent}}}}
	p.valueTitle = map[string]string{"subtle": "Subtle", "bold": "Bold"}[ds.MenuAccent]
	p.indicator = func(v string) h.Node {
		path := func(d string) string {
			return `<path d="` + d + `" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-accent="` + v + `" class="fill-muted-foreground/30 data-[accent=bold]:fill-foreground"></path>`
		}
		return h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 flex size-4 -translate-y-1/2 items-center justify-center text-base text-foreground select-none md:right-2.5"),
			h.UnsafeHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="128" height="128" viewBox="0 0 24 24" fill="none" class="size-4 text-foreground">`+
				path("M19 12.1294L12.9388 18.207C11.1557 19.9949 10.2641 20.8889 9.16993 20.9877C8.98904 21.0041 8.80705 21.0041 8.62616 20.9877C7.53195 20.8889 6.64039 19.9949 4.85726 18.207L2.83687 16.1811C1.72104 15.0622 1.72104 13.2482 2.83687 12.1294M19 12.1294L10.9184 4.02587M19 12.1294H2.83687M10.9184 4.02587L2.83687 12.1294M10.9184 4.02587L8.89805 2")+
				path("M22 20C22 21.1046 21.1046 22 20 22C18.8954 22 18 21.1046 18 20C18 18.8954 20 17 20 17C20 17 22 18.8954 22 20Z")+`</svg>`))
	}
	return p
}

func crBasePicker(ds designSystem) crPickerSpec {
	d := createDataOnce()
	p := crPickerSpec{param: "base", label: "Base", value: ds.Base}
	var items []crPickerItem
	logos := map[string]string{}
	for _, b := range d.Bases {
		items = append(items, crPickerItem{value: b.Name, title: b.Title})
		logos[b.Name] = b.Logo
		if b.Name == ds.Base {
			p.valueTitle = b.Title
		}
	}
	p.groups = []crPickerGroup{{items: items}}
	p.indicator = func(v string) h.Node {
		return h.Div(h.Class("pointer-events-none absolute top-1/2 right-4 size-4 -translate-y-1/2 text-foreground select-none md:right-2.5 *:[svg]:size-4 *:[svg]:text-foreground!"), h.UnsafeHTML(logos[v]))
	}
	return p
}

// createClientConfig is the data create.js needs: the theme variables and
// font families (to build the preview's theme as DesignSystemProvider
// does), the preset fields, and the shuffle's value lists.
func createClientConfig() string {
	d := createDataOnce()
	fonts := map[string]map[string]string{}
	for _, f := range d.Fonts {
		fonts[f.Name] = map[string]string{"title": f.Title, "family": f.fontFamily(), "type": f.Type}
	}
	b, _ := json.Marshal(map[string]any{
		"baseColors": d.BaseColors, "themes": d.Themes, "radii": d.Radii, "fonts": fonts,
		"styles": d.Styles, "iconLibraries": d.IconLibraries, "bases": d.Bases, "presets": d.Presets,
		"menuColorLabels": crMenuColorLabels,
	})
	return string(b)
}
