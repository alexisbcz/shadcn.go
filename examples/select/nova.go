package selectex

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/field"
	selectui "shadcn.go/components/ui/select"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the select page (ui/apps/v4/examples/radix/select-*.tsx),
// rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "select", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("select-demo", NovaDemo),
		nova("select-scrollable", NovaScrollable),
		nova("select-disabled", NovaDisabled),
		nova("select-groups", NovaGroups),
		nova("select-invalid", NovaInvalid),
		nova("select-align-item", NovaAlignItem),
		nova("select-rtl", NovaRtl),
	)
}

var fruits = []opt{{"apple", "Apple"}, {"banana", "Banana"}, {"blueberry", "Blueberry"}, {"grapes", "Grapes"}, {"pineapple", "Pineapple"}}

func NovaDemo() h.Node {
	return selectui.Select(
		selectui.Trigger(h.Class("w-full max-w-48"), selectui.Value(selectui.Placeholder("Select a fruit"))),
		selectui.Content(group("Fruits", fruits...)),
	)
}

func NovaDisabled() h.Node {
	return selectui.Select(selectui.Disabled,
		selectui.Trigger(h.Class("w-full max-w-48"), selectui.Value(selectui.Placeholder("Select a fruit"))),
		selectui.Content(selectui.Group(
			selectui.Item(selectui.ItemValue("apple"), "Apple"),
			selectui.Item(selectui.ItemValue("banana"), "Banana"),
			selectui.Item(selectui.ItemValue("blueberry"), "Blueberry"),
			selectui.Item(selectui.ItemValue("grapes"), selectui.Disabled, "Grapes"),
			selectui.Item(selectui.ItemValue("pineapple"), "Pineapple"),
		)),
	)
}

func NovaGroups() h.Node {
	return selectui.Select(
		selectui.Trigger(h.Class("w-full max-w-48"), selectui.Value(selectui.Placeholder("Select a fruit"))),
		selectui.Content(
			group("Fruits", opt{"apple", "Apple"}, opt{"banana", "Banana"}, opt{"blueberry", "Blueberry"}),
			selectui.Separator(),
			group("Vegetables", opt{"carrot", "Carrot"}, opt{"broccoli", "Broccoli"}, opt{"spinach", "Spinach"}),
		),
	)
}

func NovaInvalid() h.Node {
	return field.Field(h.Data("invalid", "true"), h.Class("w-full max-w-48"),
		field.Label("Fruit"),
		selectui.Select(
			selectui.Trigger(h.Aria("invalid", "true"), selectui.Value(selectui.Placeholder("Select a fruit"))),
			selectui.Content(selectui.Group(items(opt{"apple", "Apple"}, opt{"banana", "Banana"}, opt{"blueberry", "Blueberry"})...)),
		),
		field.Error("Please select a fruit."),
	)
}

func NovaAlignItem() h.Node {
	return field.Group(h.Class("w-full max-w-xs"),
		field.Field(field.Horizontal,
			field.Content(
				field.Label(h.For("align-item"), "Align Item"),
				field.Description("Toggle to align the item with the trigger."),
			),
			switchui.Switch(h.ID("align-item"), switchui.Checked(true)),
		),
		field.Field(
			selectui.Select(selectui.DefaultValue("banana"),
				selectui.Trigger(selectui.Value()),
				selectui.Content(selectui.Group(items(fruits...)...)),
			),
		),
	)
}

func NovaScrollable() h.Node {
	return selectui.Select(
		selectui.Trigger(h.Class("w-full max-w-64"), selectui.Value(selectui.Placeholder("Select a timezone"))),
		selectui.Content(
			group("North America",
				opt{"est", "Eastern Standard Time"}, opt{"cst", "Central Standard Time"},
				opt{"mst", "Mountain Standard Time"}, opt{"pst", "Pacific Standard Time"},
				opt{"akst", "Alaska Standard Time"}, opt{"hst", "Hawaii Standard Time"},
			),
			group("Europe & Africa",
				opt{"gmt", "Greenwich Mean Time"}, opt{"cet", "Central European Time"},
				opt{"eet", "Eastern European Time"}, opt{"west", "Western European Summer Time"},
				opt{"cat", "Central Africa Time"}, opt{"eat", "East Africa Time"},
			),
			group("Asia",
				opt{"msk", "Moscow Time"}, opt{"ist", "India Standard Time"},
				opt{"cst_china", "China Standard Time"}, opt{"jst", "Japan Standard Time"},
				opt{"kst", "Korea Standard Time"}, opt{"ist_indonesia", "Indonesia Central Standard Time"},
			),
			group("Australia & Pacific",
				opt{"awst", "Australian Western Standard Time"}, opt{"acst", "Australian Central Standard Time"},
				opt{"aest", "Australian Eastern Standard Time"}, opt{"nzst", "New Zealand Standard Time"},
				opt{"fjt", "Fiji Time"},
			),
			group("South America",
				opt{"art", "Argentina Time"}, opt{"bot", "Bolivia Time"},
				opt{"brt", "Brasilia Time"}, opt{"clt", "Chile Standard Time"},
			),
		),
	)
}

// NovaRtl is examples/radix/select-rtl.tsx in its default language (Arabic).
func NovaRtl() h.Node {
	// Only the trigger and the content get dir="rtl" (the root's direction
	// stays ltr, so the item-aligned position uses left).
	return style.ScopeRTL(selectui.Select(
		selectui.Trigger(h.Class("w-32"), h.Dir("rtl"), selectui.Value(selectui.Placeholder("اختر فاكهة"))),
		selectui.Content(h.Dir("rtl"), h.Data("lang", "ar"),
			group("الفواكه", opt{"apple", "تفاح"}, opt{"banana", "موز"}, opt{"blueberry", "توت أزرق"}, opt{"grapes", "عنب"}, opt{"pineapple", "أناناس"}),
			selectui.Separator(),
			group("الخضروات", opt{"carrot", "جزر"}, opt{"broccoli", "بروكلي"}, opt{"spinach", "سبانخ"}),
		),
	))
}
