package combobox

import (
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	cb "shadcn.go/components/ui/combobox"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the combobox page
// (ui/apps/v4/examples/radix/combobox-*.tsx), rendered in radix-nova.
// combobox-basic and combobox-multiple are the fixtures of the same name.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "combobox", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("combobox-demo", Basic),
		nova("combobox-auto-highlight", func() h.Node { return novaFrameworks(cb.AutoHighlight) }),
		nova("combobox-clear", func() h.Node { return novaFrameworks(cb.DefaultValue("Next.js"), cb.ShowClear(true)) }),
		nova("combobox-disabled", func() h.Node { return novaFrameworks(cb.Disabled) }),
		nova("combobox-invalid", func() h.Node { return novaFrameworks(h.Aria("invalid", "true")) }),
		nova("combobox-groups", NovaGroups),
		nova("combobox-input-group", NovaInputGroup),
		nova("combobox-popup", NovaPopup),
		nova("combobox-custom", NovaCustom),
	)
}

// novaFrameworks is the frameworks combobox; options go to the Combobox
// (AutoHighlight, DefaultValue) or its Input (the others).
func novaFrameworks(opts ...any) h.Node {
	var root, input []any
	for _, o := range opts {
		switch o.(type) {
		case cb.AutoHighlightFlag, cb.DefaultValue:
			root = append(root, o)
		default:
			input = append(input, o)
		}
	}
	return cb.Combobox(append(root,
		cb.Input(append([]any{h.Placeholder("Select a framework")}, input...)...),
		cb.Content(frameworkItems()...),
	)...)
}

func NovaGroups() h.Node {
	return cb.Combobox(
		cb.Input(h.Placeholder("Select a timezone")),
		cb.Content(cb.Empty("No timezones found."), novaTimezoneList(true)),
	)
}

func NovaInputGroup() h.Node {
	return cb.Combobox(
		cb.Input(h.Placeholder("Select a timezone"), inputgroup.Addon(lucide.Globe())),
		cb.Content(cb.AlignOffset(-28), h.Class("w-60"), cb.Empty("No timezones found."), novaTimezoneList(false)),
	)
}

var novaCountries = []country{
	{"", "", "Select country", ""},
	{"ar", "argentina", "Argentina", "South America"},
	{"au", "australia", "Australia", "Oceania"},
	{"br", "brazil", "Brazil", "South America"},
	{"ca", "canada", "Canada", "North America"},
	{"cn", "china", "China", "Asia"},
	{"co", "colombia", "Colombia", "South America"},
	{"eg", "egypt", "Egypt", "Africa"},
	{"fr", "france", "France", "Europe"},
	{"de", "germany", "Germany", "Europe"},
	{"it", "italy", "Italy", "Europe"},
	{"jp", "japan", "Japan", "Asia"},
	{"ke", "kenya", "Kenya", "Africa"},
	{"mx", "mexico", "Mexico", "North America"},
	{"nz", "new-zealand", "New Zealand", "Oceania"},
	{"ng", "nigeria", "Nigeria", "Africa"},
	{"za", "south-africa", "South Africa", "Africa"},
	{"kr", "south-korea", "South Korea", "Asia"},
	{"gb", "united-kingdom", "United Kingdom", "Europe"},
	{"us", "united-states", "United States", "North America"},
}

func NovaPopup() h.Node {
	return cb.Combobox(cb.DefaultValue(""),
		cb.Trigger(cb.Render(button.Button(button.VariantOutline, h.Class("w-64 justify-between font-normal"))), cb.Value()),
		cb.Content(
			cb.Input(cb.ShowTrigger(false), h.Placeholder("Search")),
			cb.Empty("No items found."),
			cb.List(h.Map(novaCountries, func(c country) h.Node {
				return cb.Item(cb.ItemValue(c.value), cb.ItemLabel(c.label), c.label)
			})),
		),
	)
}

func NovaCustom() h.Node {
	return cb.Combobox(
		cb.Input(h.Placeholder("Search countries...")),
		cb.Content(
			cb.Empty("No countries found."),
			cb.List(h.Map(novaCountries[1:], func(c country) h.Node {
				return cb.Item(cb.ItemValue(c.label), cb.ItemLabel(c.label),
					item.Item(item.Size("xs"), h.Class("p-0"),
						item.Content(
							item.Title(h.Class("whitespace-nowrap"), c.label),
							item.Description(c.continent, h.Comment(""), " (", h.Comment(""), c.code, h.Comment(""), ")"),
						),
					),
				)
			})),
		),
	)
}

var novaTimezones = []timezoneGroup{
	{"Americas", []string{"(GMT-5) New York", "(GMT-8) Los Angeles", "(GMT-6) Chicago", "(GMT-5) Toronto", "(GMT-8) Vancouver", "(GMT-3) São Paulo"}},
	{"Europe", []string{"(GMT+0) London", "(GMT+1) Paris", "(GMT+1) Berlin", "(GMT+1) Rome", "(GMT+1) Madrid", "(GMT+1) Amsterdam"}},
	{"Asia/Pacific", []string{"(GMT+9) Tokyo", "(GMT+8) Shanghai", "(GMT+8) Singapore", "(GMT+4) Dubai", "(GMT+11) Sydney", "(GMT+9) Seoul"}},
}

// novaTimezoneList is timezoneList with the docs' timezones.
func novaTimezoneList(separators bool) h.Node {
	return cb.List(
		h.Data("example-combobox-groups", strconv.Itoa(len(novaTimezones))),
		h.Map(novaTimezones, func(g timezoneGroup) h.Node {
			var sep any
			if separators {
				sep = cb.Separator()
			}
			return cb.Group(
				cb.Label(g.value),
				cb.Collection(h.Map(g.items, func(it string) h.Node {
					return cb.Item(cb.ItemValue(it), it)
				})),
				sep,
			)
		}),
	)
}
