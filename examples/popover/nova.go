package popover

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/popover"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the popover page (ui/apps/v4/examples/radix/popover-*.tsx),
// rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "popover", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("popover-demo", Demo),
		nova("popover-basic", NovaBasic),
		nova("popover-alignments", NovaAlignments),
		nova("popover-form", NovaForm),
		nova("popover-rtl", NovaRtl),
	)
}

func NovaAlignments() h.Node {
	one := func(label string, align popover.Align, text string) h.Node {
		return popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, label)),
			popover.Content(align, h.Class("w-40"), text),
		)
	}
	return h.Div(h.Class("flex gap-6"),
		one("Start", popover.AlignStart, "Aligned to start"),
		one("Center", popover.AlignCenter, "Aligned to center"),
		one("End", popover.AlignEnd, "Aligned to end"),
	)
}

func NovaBasic() h.Node {
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open Popover")),
		popover.Content(popover.AlignStart,
			popover.Header(
				popover.Title("Dimensions"),
				popover.Description("Set the dimensions for the layer."),
			),
		),
	)
}

func NovaForm() h.Node {
	row := func(id, text, value string) h.Node {
		return field.Field(field.Horizontal,
			field.Label(h.For(id), h.Class("w-1/2"), text),
			input.Input(h.ID(id), h.Value(value)),
		)
	}
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open Popover")),
		popover.Content(h.Class("w-64"), popover.AlignStart,
			popover.Header(
				popover.Title("Dimensions"),
				popover.Description("Set the dimensions for the layer."),
			),
			field.Group(h.Class("gap-4"),
				row("width", "Width", "100%"),
				row("height", "Height", "25px"),
			),
		),
	)
}

func NovaRtl() h.Node {
	sides := []struct {
		side  popover.Side
		label string
	}{{popover.SideLeft, "يسار"}, {popover.SideTop, "أعلى"}, {popover.SideBottom, "أسفل"}, {popover.SideRight, "يمين"}}
	var kids []any
	for _, s := range sides {
		kids = append(kids, popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, s.label)),
			popover.Content(s.side, h.Dir("rtl"),
				popover.Header(
					popover.Title("الأبعاد"),
					popover.Description("تعيين الأبعاد للطبقة."),
				),
			),
		))
	}
	return h.Div(append([]any{h.Class("flex flex-wrap justify-center gap-2")}, kids...)...)
}
