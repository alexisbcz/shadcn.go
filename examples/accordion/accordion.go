// Package accordion ports the upstream accordion-* examples and the port's
// accordion-* fixtures.
package accordion

import (
	"shadcn.go/components/ui/accordion"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "accordion", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("accordion-demo", Demo),
		ex("accordion-multiple", Multiple),
		ex("accordion-states", States),
		ex("accordion-nested", Nested),
	)
}

func Demo() h.Node {
	content := func(a, b string) *accordion.ContentBuilder {
		return accordion.Content(h.Class("flex flex-col gap-4 text-balance"), h.P(a), h.P(b))
	}
	return accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), h.Class("w-full"), accordion.Value("item-1"),
		accordion.Item(accordion.Value("item-1"),
			accordion.Trigger("Product Information"),
			content(
				"Our flagship product combines cutting-edge technology with sleek design. Built with premium materials, it offers unparalleled performance and reliability.",
				"Key features include advanced processing capabilities, and an intuitive user interface designed for both beginners and experts.",
			),
		),
		accordion.Item(accordion.Value("item-2"),
			accordion.Trigger("Shipping Details"),
			content(
				"We offer worldwide shipping through trusted courier partners. Standard delivery takes 3-5 business days, while express shipping ensures delivery within 1-2 business days.",
				"All orders are carefully packaged and fully insured. Track your shipment in real-time through our dedicated tracking portal.",
			),
		),
		accordion.Item(accordion.Value("item-3"),
			accordion.Trigger("Return Policy"),
			content(
				"We stand behind our products with a comprehensive 30-day return policy. If you're not completely satisfied, simply return the item in its original condition.",
				"Our hassle-free return process includes free return shipping and full refunds processed within 48 hours of receiving the returned item.",
			),
		),
	)
}

// Multiple is the accordion-multiple fixture: type="multiple" with two
// items open by default and a disabled item.
func Multiple() h.Node {
	return accordion.Accordion(accordion.TypeMultiple, h.Class("w-full max-w-md"), accordion.Value("a"), accordion.Value("c"),
		accordion.Item(accordion.Value("a"),
			accordion.Trigger("Is it accessible?"),
			accordion.Content("Yes. It adheres to the WAI-ARIA design pattern."),
		),
		accordion.Item(accordion.Value("b"),
			accordion.Trigger("Is it styled?"),
			accordion.Content("Yes. It comes with default styles that match the other components."),
		),
		accordion.Item(accordion.Value("c"),
			accordion.Trigger("Is it animated?"),
			accordion.Content("Yes. It's animated by default, but you can disable it if you prefer."),
		),
		accordion.Item(accordion.Value("d"), accordion.Disabled(true),
			accordion.Trigger("Is it disabled?"),
			accordion.Content("This item cannot be opened."),
		),
		accordion.Item(accordion.Value("e"),
			accordion.Trigger("Can it nest?"),
			accordion.Content("Yes."),
		),
	)
}

// States is the accordion-states fixture: a single, non-collapsible
// accordion with a disabled item, a horizontal one in a right-to-left
// context, and a fully disabled one.
func States() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"),
		accordion.Accordion(accordion.TypeSingle, accordion.Value("one"),
			accordion.Item(accordion.Value("one"),
				accordion.Trigger("One"),
				accordion.Content("First content."),
			),
			accordion.Item(accordion.Value("two"), accordion.Disabled(true),
				accordion.Trigger("Two (disabled)"),
				accordion.Content("Second content."),
			),
			accordion.Item(accordion.Value("three"),
				accordion.Trigger("Three"),
				accordion.Content("Third content."),
			),
		),
		h.Div(h.Dir("rtl"),
			accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), accordion.OrientationHorizontal, accordion.DirRTL, h.Class("flex gap-4"),
				accordion.Item(accordion.Value("x"),
					accordion.Trigger("X"),
					accordion.Content("X content."),
				),
				accordion.Item(accordion.Value("y"),
					accordion.Trigger("Y"),
					accordion.Content("Y content."),
				),
				accordion.Item(accordion.Value("z"),
					accordion.Trigger("Z"),
					accordion.Content("Z content."),
				),
			),
		),
		accordion.Accordion(accordion.TypeSingle, accordion.Disabled(true), accordion.Value("p"),
			accordion.Item(accordion.Value("p"),
				accordion.Trigger("Disabled open"),
				accordion.Content("Disabled accordion content."),
			),
			accordion.Item(accordion.Value("q"),
				accordion.Trigger("Disabled closed"),
				accordion.Content("Hidden."),
			),
		),
	)
}

// Nested is the accordion-nested fixture: a horizontal accordion read left
// to right, a vertical one with dir="rtl", and an accordion nested in
// another one's content.
func Nested() h.Node {
	return h.Div(h.Class("flex w-full max-w-md flex-col gap-6"),
		accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), accordion.OrientationHorizontal, h.Class("flex gap-4"),
			accordion.Item(accordion.Value("l"),
				accordion.Trigger("L"),
				accordion.Content("L content."),
			),
			accordion.Item(accordion.Value("m"),
				accordion.Trigger("M"),
				accordion.Content("M content."),
			),
			accordion.Item(accordion.Value("n"),
				accordion.Trigger("N"),
				accordion.Content("N content."),
			),
		),
		accordion.Accordion(accordion.TypeMultiple, accordion.DirRTL,
			accordion.Item(accordion.Value("outer-1"),
				accordion.Trigger("Outer one"),
				accordion.Content(
					accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("inner-1"),
						accordion.Item(accordion.Value("inner-1"),
							accordion.Trigger("Inner one"),
							accordion.Content("Inner one content."),
						),
						accordion.Item(accordion.Value("inner-2"),
							accordion.Trigger("Inner two"),
							accordion.Content("Inner two content."),
						),
					),
				),
			),
			accordion.Item(accordion.Value("outer-2"),
				accordion.Trigger("Outer two"),
				accordion.Content("Outer two content."),
			),
		),
	)
}
