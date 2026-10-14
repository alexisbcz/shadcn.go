// Package popover ports the upstream popover-* examples and the popover
// parity fixtures.
package popover

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/popover"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "popover", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("popover-demo", Demo),
		ex("popover-sides", Sides),
		ex("popover-collision", Collision),
		ex("popover-modal-anchor", ModalAnchor),
		ex("popover-default-open", DefaultOpen),
	)
}

func Demo() h.Node {
	field := func(id, text, value string) h.Node {
		return h.Div(h.Class("grid grid-cols-3 items-center gap-4"),
			label.Label(h.For(id), text),
			input.Input(h.ID(id), h.Value(value), h.Class("col-span-2 h-8")),
		)
	}
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open popover")),
		popover.Content(h.Class("w-80"),
			h.Div(h.Class("grid gap-4"),
				h.Div(h.Class("space-y-2"),
					h.H4(h.Class("leading-none font-medium"), "Dimensions"),
					h.P(h.Class("text-sm text-muted-foreground"), "Set the dimensions for the layer."),
				),
				h.Div(h.Class("grid gap-2"),
					field("width", "Width", "100%"),
					field("maxWidth", "Max. width", "300px"),
					field("height", "Height", "25px"),
					field("maxHeight", "Max. height", "none"),
				),
			),
		),
	)
}

// Sides is the parity fixture popover-sides.
func Sides() h.Node {
	type variant struct {
		side        popover.Side
		align       popover.Align
		sideOffset  *float64
		alignOffset *float64
	}
	f := func(v float64) *float64 { return &v }
	variants := []variant{
		{popover.SideTop, popover.AlignStart, nil, nil},
		{popover.SideTop, popover.AlignCenter, nil, nil},
		{popover.SideTop, popover.AlignEnd, nil, nil},
		{popover.SideRight, popover.AlignStart, nil, nil},
		{popover.SideRight, popover.AlignCenter, f(16), nil},
		{popover.SideRight, popover.AlignEnd, nil, nil},
		{popover.SideBottom, popover.AlignStart, nil, f(12)},
		{popover.SideBottom, popover.AlignCenter, f(0), nil},
		{popover.SideBottom, popover.AlignEnd, nil, f(-8)},
		{popover.SideLeft, popover.AlignStart, nil, f(6)},
		{popover.SideLeft, popover.AlignCenter, nil, nil},
		{popover.SideLeft, popover.AlignEnd, f(12), nil},
	}
	return h.Div(h.Class("grid grid-cols-3 gap-x-24 gap-y-6"),
		h.Map(variants, func(v variant) h.Node {
			opts := []any{h.Class("w-48"), v.side, v.align}
			if v.sideOffset != nil {
				opts = append(opts, popover.SideOffset(*v.sideOffset))
			}
			if v.alignOffset != nil {
				opts = append(opts, popover.AlignOffset(*v.alignOffset))
			}
			opts = append(opts, popover.Header(
				popover.Title(string(v.side)+" / "+string(v.align)),
				popover.Description("Placed by the popper."),
			))
			return popover.Popover(
				popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, string(v.side)+" "+string(v.align))),
				popover.Content(opts...),
			)
		}),
	)
}

// Collision is the parity fixture popover-collision.
func Collision() h.Node {
	pop := func(pos, label string, content ...any) h.Node {
		return h.Div(h.Class(pos),
			popover.Popover(
				popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, label)),
				popover.Content(content...),
			),
		)
	}
	return h.Div(
		pop("fixed top-2 left-2", "flip top", popover.SideTop, popover.AlignEnd, "Flipped to the bottom, shifted right."),
		pop("fixed right-2 bottom-2", "flip bottom", popover.CollisionPadding(24), "Flipped to the top, padded."),
		pop("fixed top-40 left-2", "flip left", popover.SideLeft, popover.CollisionPaddingSides{Top: 40, Left: 8}, "Flipped to the right."),
		pop("fixed top-2 right-2", "no avoid", popover.AvoidCollisions(false), popover.SideTop, "Overflows the viewport."),
		pop("fixed bottom-2 left-2", "sticky always", popover.SideRight, popover.AlignEnd, popover.StickyAlways, "Always inside."),
		h.Div(h.Class("h-32 w-64 overflow-auto rounded-md border p-4"), h.Data("testid", "scroller"),
			h.Div(h.Class("h-96 pt-8"),
				popover.Popover(
					popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, "detached")),
					popover.Content(popover.HideWhenDetached(true), h.Class("w-40"), popover.SideRight, "Hidden when detached."),
				),
			),
		),
	)
}

// ModalAnchor is the parity fixture popover-modal-anchor.
func ModalAnchor() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-6"),
		popover.Popover(popover.Modal(true),
			popover.Anchor(h.Class("w-72 rounded-md border border-dashed p-2 text-center text-sm"), "Anchor"),
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open modal")),
			popover.Content(h.Class("grid gap-2"), popover.SideTop,
				input.Input(h.Value("First")),
				input.Input(h.Value("Second")),
				popover.Close(lib.AsChild, button.Button(button.SizeSm, "Done")),
				popover.Arrow(h.Class("fill-popover")),
			),
		),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Anchored below")),
			popover.Anchor(h.Class("w-72 rounded-md border border-dashed p-2 text-center text-sm"), "Second anchor"),
			popover.Content(popover.AlignStart, "Positioned against the anchor after the trigger."),
		),
		button.Button(button.VariantSecondary, "Outside"),
	)
}

// DefaultOpen is the parity fixture popover-default-open.
func DefaultOpen() h.Node {
	return popover.Popover(popover.DefaultOpen,
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open on load")),
		popover.Content(popover.AlignStart,
			"Opened by defaultOpen. ", h.A(h.Href("#x"), h.Class("underline"), "Link"), " ",
			h.Button(h.Class("underline"), "Button"),
		),
	)
}
