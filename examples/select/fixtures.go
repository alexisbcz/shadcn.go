package selectex

import (
	"fmt"

	"shadcn.go/components/icons/lucide"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "select", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("select-popper", Popper),
		ex("select-states", States),
		ex("select-form", FormFixture),
		ex("select-default-open", DefaultOpenFixture),
	)
}

// Popper is the parity fixture select-popper.
func Popper() h.Node {
	many := make([]any, 30)
	for i := range many {
		many[i] = selectui.Item(selectui.ItemValue(fmt.Sprintf("item-%d", i+1)), fmt.Sprintf("Item %d", i+1))
	}
	return h.Div(h.Class("flex items-start gap-6"),
		selectui.Select(selectui.DefaultValue("banana"),
			selectui.Trigger(selectui.SizeSm, h.Class("w-[160px]"), selectui.Value(selectui.Placeholder("Fruit"))),
			selectui.Content(selectui.PositionPopper,
				selectui.Group(
					selectui.Label("Fruits"),
					selectui.Item(selectui.ItemValue("apple"), "Apple"),
					selectui.Item(selectui.ItemValue("banana"), "Banana"),
					selectui.Item(selectui.ItemValue("blueberry"), selectui.Disabled, "Blueberry"),
				),
				selectui.Separator(),
				selectui.Group(
					selectui.Label("Vegetables"),
					selectui.Item(selectui.ItemValue("carrot"), "Carrot"),
					selectui.Item(selectui.ItemValue("potato"), selectui.TextValue("Spud"), "Potato"),
				),
			),
		),
		selectui.Select(
			selectui.Trigger(h.Class("w-[160px]"), selectui.Value(selectui.Placeholder("Top, end"))),
			selectui.Content(selectui.PositionPopper, selectui.SideTop, selectui.AlignEnd, selectui.SideOffset(8),
				selectui.Item(selectui.ItemValue("one"), "One"),
				selectui.Item(selectui.ItemValue("two"), "Two"),
				selectui.Item(selectui.ItemValue("three"), "Three"),
			),
		),
		selectui.Select(
			selectui.Trigger(h.Class("w-[160px]"), selectui.Value(selectui.Placeholder("Many"))),
			selectui.Content(append([]any{selectui.PositionPopper, selectui.SideRight, selectui.AlignOffset(4)}, many...)...),
		),
	)
}

// States is the parity fixture select-states.
func States() h.Node {
	return h.Div(h.Class("grid grid-cols-2 gap-4"),
		selectui.Select(selectui.Disabled, selectui.DefaultValue("apple"),
			selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Disabled"))),
			selectui.Content(
				selectui.Item(selectui.ItemValue("apple"), "Apple"),
				selectui.Item(selectui.ItemValue("cherry"), "Cherry"),
			),
		),
		selectui.Select(
			selectui.Trigger(h.Class("w-[180px]"), selectui.Disabled, selectui.Value(selectui.Placeholder("Disabled trigger"))),
			selectui.Content(selectui.Item(selectui.ItemValue("apple"), "Apple")),
		),
		selectui.Select(selectui.DefaultValue("cherry"),
			selectui.Trigger(h.Class("w-[180px]"), h.Aria("invalid", "true"), selectui.Value(selectui.Placeholder("With icons"))),
			selectui.Content(
				selectui.Item(selectui.ItemValue("apple"), lucide.Apple(), "Apple"),
				selectui.Item(selectui.ItemValue("cherry"), lucide.Cherry(), "Cherry"),
			),
		),
		selectui.Select(selectui.DefaultValue("b"),
			selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Custom"), "Custom value")),
			selectui.Content(
				selectui.Item(selectui.ItemValue("a"), "Alpha"),
				selectui.Item(selectui.ItemValue("b"), "Beta"),
			),
		),
		selectui.Select(selectui.Dir("rtl"),
			selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Right to left"))),
			selectui.Content(
				selectui.Item(selectui.ItemValue("one"), "One"),
				selectui.Item(selectui.ItemValue("two"), "Two"),
			),
		),
		selectui.Select(
			selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Disabled first"))),
			selectui.Content(
				selectui.Item(selectui.ItemValue("one"), selectui.Disabled, "One"),
				selectui.Item(selectui.ItemValue("two"), "Two"),
				selectui.Item(selectui.ItemValue("three"), selectui.Disabled, "Three"),
				selectui.Item(selectui.ItemValue("four"), "Four"),
			),
		),
	)
}

// FormFixture is the parity fixture select-form.
func FormFixture() h.Node {
	return h.Div(h.Class("flex flex-col gap-4"),
		h.Form(h.ID("select-form"), h.Class("flex gap-3"),
			selectui.Select(selectui.Name("fruit"), selectui.Required, selectui.AutoComplete("off"),
				selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Fruit"))),
				selectui.Content(
					selectui.Item(selectui.ItemValue("apple"), "Apple"),
					selectui.Item(selectui.ItemValue("banana"), "Banana"),
					selectui.Item(selectui.ItemValue("cherry"), selectui.Disabled, "Cherry"),
				),
			),
			selectui.Select(selectui.Name("size"), selectui.DefaultValue("m"), selectui.Disabled,
				selectui.Trigger(h.Class("w-[120px]"), selectui.Value()),
				selectui.Content(
					selectui.Item(selectui.ItemValue("s"), "Small"),
					selectui.Item(selectui.ItemValue("m"), "Medium"),
				),
			),
		),
		selectui.Select(selectui.Name("color"), selectui.Form("select-form"), selectui.DefaultValue("red"),
			selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Color"))),
			selectui.Content(
				selectui.Item(selectui.ItemValue("red"), "Red"),
				selectui.Item(selectui.ItemValue("green"), "Green"),
			),
		),
	)
}

// DefaultOpenFixture is the parity fixture select-default-open.
func DefaultOpenFixture() h.Node {
	return selectui.Select(selectui.DefaultOpen, selectui.DefaultValue("two"),
		selectui.Trigger(h.Class("w-[180px]"), selectui.Value(selectui.Placeholder("Pick one"))),
		selectui.Content(
			selectui.Item(selectui.ItemValue("one"), "One"),
			selectui.Item(selectui.ItemValue("two"), "Two"),
			selectui.Item(selectui.ItemValue("three"), "Three"),
		),
	)
}
