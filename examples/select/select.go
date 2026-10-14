// Package selectex ports the upstream select-* examples and the select
// parity fixtures.
package selectex

import (
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "select", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("select-demo", Demo),
		ex("select-scrollable", Scrollable),
	)
}

type opt struct{ value, label string }

func items(opts ...opt) []any {
	out := make([]any, len(opts))
	for i, o := range opts {
		out[i] = selectui.Item(selectui.ItemValue(o.value), o.label)
	}
	return out
}

func group(label string, opts ...opt) *selectui.GroupBuilder {
	return selectui.Group(append([]any{selectui.Label(label)}, items(opts...)...)...)
}

func Demo() h.Node {
	return selectui.Select(
		selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Select a fruit"))),
		selectui.Content(
			group("Fruits",
				opt{"apple", "Apple"}, opt{"banana", "Banana"}, opt{"blueberry", "Blueberry"},
				opt{"grapes", "Grapes"}, opt{"pineapple", "Pineapple"},
			),
		),
	)
}

func Scrollable() h.Node {
	return selectui.Select(
		selectui.Trigger(h.Class("w-[280px]"), selectui.Value(selectui.Placeholder("Select a timezone"))),
		selectui.Content(
			group("North America",
				opt{"est", "Eastern Standard Time (EST)"}, opt{"cst", "Central Standard Time (CST)"},
				opt{"mst", "Mountain Standard Time (MST)"}, opt{"pst", "Pacific Standard Time (PST)"},
				opt{"akst", "Alaska Standard Time (AKST)"}, opt{"hst", "Hawaii Standard Time (HST)"},
			),
			group("Europe & Africa",
				opt{"gmt", "Greenwich Mean Time (GMT)"}, opt{"cet", "Central European Time (CET)"},
				opt{"eet", "Eastern European Time (EET)"}, opt{"west", "Western European Summer Time (WEST)"},
				opt{"cat", "Central Africa Time (CAT)"}, opt{"eat", "East Africa Time (EAT)"},
			),
			group("Asia",
				opt{"msk", "Moscow Time (MSK)"}, opt{"ist", "India Standard Time (IST)"},
				opt{"cst_china", "China Standard Time (CST)"}, opt{"jst", "Japan Standard Time (JST)"},
				opt{"kst", "Korea Standard Time (KST)"}, opt{"ist_indonesia", "Indonesia Central Standard Time (WITA)"},
			),
			group("Australia & Pacific",
				opt{"awst", "Australian Western Standard Time (AWST)"}, opt{"acst", "Australian Central Standard Time (ACST)"},
				opt{"aest", "Australian Eastern Standard Time (AEST)"}, opt{"nzst", "New Zealand Standard Time (NZST)"},
				opt{"fjt", "Fiji Time (FJT)"},
			),
			group("South America",
				opt{"art", "Argentina Time (ART)"}, opt{"bot", "Bolivia Time (BOT)"},
				opt{"brt", "Brasilia Time (BRT)"}, opt{"clt", "Chile Standard Time (CLT)"},
			),
		),
	)
}
