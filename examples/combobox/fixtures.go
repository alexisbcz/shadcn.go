package combobox

import (
	_ "embed"
	"strconv"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/button"
	cb "shadcn.go/components/ui/combobox"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// fixtures.js mirrors combobox-groups' render function: a group's
// separator shows while its index among the rendered groups is below the
// number of groups minus one.
//
//go:embed fixtures.js
var fixturesScript string

func registerFixtures(ex func(string, func() h.Node) registry.Example) {
	runtime.Module("example-combobox-fixtures", fixturesScript)
	registry.Register(
		ex("combobox-basic", Basic),
		ex("combobox-states", States),
		ex("combobox-groups", Groups),
		ex("combobox-multiple", MultipleChips),
		ex("combobox-popup", PopupInput),
		ex("combobox-custom", Custom),
	)
}

var frameworkNames = []string{"Next.js", "SvelteKit", "Nuxt.js", "Remix", "Astro"}

func frameworkItems() []any {
	return []any{
		cb.Empty("No items found."),
		cb.List(h.Map(frameworkNames, func(f string) h.Node {
			return cb.Item(cb.ItemValue(f), f)
		})),
	}
}

// Basic is the parity fixture combobox-basic.
func Basic() h.Node {
	return cb.Combobox(
		cb.Input(h.Placeholder("Select a framework")),
		cb.Content(frameworkItems()...),
	)
}

// States is the parity fixture combobox-states.
func States() h.Node {
	return h.Div(h.Class("flex w-72 flex-col gap-4"),
		cb.Combobox(cb.AutoHighlight,
			cb.Input(h.Placeholder("Auto highlight")),
			cb.Content(frameworkItems()...),
		),
		cb.Combobox(cb.DefaultValue("Next.js"),
			cb.Input(h.Placeholder("With clear"), cb.ShowClear(true)),
			cb.Content(frameworkItems()...),
		),
		cb.Combobox(
			cb.Input(h.Placeholder("Disabled"), cb.Disabled),
			cb.Content(frameworkItems()...),
		),
		cb.Combobox(
			cb.Input(h.Placeholder("Invalid"), h.Aria("invalid", "true")),
			cb.Content(frameworkItems()...),
		),
	)
}

type timezoneGroup struct {
	value string
	items []string
}

var timezones = []timezoneGroup{
	{"Americas", []string{"(GMT-5) New York", "(GMT-8) Los Angeles", "(GMT-6) Chicago", "(GMT-5) Toronto"}},
	{"Europe", []string{"(GMT+0) London", "(GMT+1) Paris", "(GMT+1) Berlin", "(GMT+1) Rome"}},
	{"Asia/Pacific", []string{"(GMT+9) Tokyo", "(GMT+8) Shanghai", "(GMT+8) Singapore", "(GMT+4) Dubai"}},
}

func timezoneList(separators bool) h.Node {
	return cb.List(
		h.Data("example-combobox-groups", strconv.Itoa(len(timezones))),
		h.Map(timezones, func(g timezoneGroup) h.Node {
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

// Groups is the parity fixture combobox-groups.
func Groups() h.Node {
	return h.Div(h.Class("flex w-72 flex-col gap-4"),
		cb.Combobox(
			cb.Input(h.Placeholder("Select a timezone")),
			cb.Content(
				cb.Empty("No timezones found."),
				timezoneList(true),
			),
		),
		cb.Combobox(
			cb.Input(h.Placeholder("With addon"),
				inputgroup.Addon(lucide.Globe()),
			),
			cb.Content(cb.AlignOffset(-28), h.Class("w-60"),
				cb.Empty("No timezones found."),
				timezoneList(false),
			),
		),
	)
}

// MultipleChips is the parity fixture combobox-multiple.
func MultipleChips() h.Node {
	anchor := cb.UseAnchor()
	return cb.Combobox(cb.Multiple, cb.AutoHighlight, cb.DefaultValue("Next.js"),
		cb.Chips(anchor, h.Class("w-full max-w-xs"),
			cb.Value(cb.Chip(), cb.ChipsInput()),
		),
		cb.Content(append([]any{anchor}, frameworkItems()...)...),
	)
}

type country struct{ code, value, label, continent string }

var popupCountries = []country{
	{"", "", "Select country", ""},
	{"ar", "argentina", "Argentina", "South America"},
	{"au", "australia", "Australia", "Oceania"},
	{"br", "brazil", "Brazil", "South America"},
	{"ca", "canada", "Canada", "North America"},
	{"cn", "china", "China", "Asia"},
	{"fr", "france", "France", "Europe"},
	{"de", "germany", "Germany", "Europe"},
	{"jp", "japan", "Japan", "Asia"},
}

// PopupInput is the parity fixture combobox-popup.
func PopupInput() h.Node {
	return cb.Combobox(cb.DefaultValue(""),
		cb.Trigger(cb.Render(button.Button(button.VariantOutline, h.Class("w-64 justify-between font-normal"))), cb.Value()),
		cb.Content(
			cb.Input(cb.ShowTrigger(false), h.Placeholder("Search")),
			cb.Empty("No items found."),
			cb.List(h.Map(popupCountries, func(c country) h.Node {
				return cb.Item(cb.ItemValue(c.value), cb.ItemLabel(c.label), c.label)
			})),
		),
	)
}

var customCountries = []country{
	{"ar", "argentina", "Argentina", "South America"},
	{"au", "australia", "Australia", "Oceania"},
	{"br", "brazil", "Brazil", "South America"},
	{"ca", "canada", "Canada", "North America"},
	{"cn", "china", "China", "Asia"},
	{"eg", "egypt", "Egypt", "Africa"},
	{"fr", "france", "France", "Europe"},
}

// Custom is the parity fixture combobox-custom: object items whose form
// value is their label (itemToStringValue).
func Custom() h.Node {
	return h.Form(
		cb.Combobox(cb.Name("country"),
			cb.Input(h.Placeholder("Search countries...")),
			cb.Content(
				cb.Empty("No countries found."),
				cb.List(h.Map(customCountries, func(c country) h.Node {
					return cb.Item(cb.ItemValue(c.label), cb.ItemLabel(c.label),
						item.Item(item.SizeSm, h.Class("p-0"),
							item.Content(
								item.Title(h.Class("whitespace-nowrap"), c.label),
								item.Description(c.continent, h.Comment(""), " (", h.Comment(""), c.code, h.Comment(""), ")"),
							),
						),
					)
				})),
			),
		),
	)
}
