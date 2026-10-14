package progress

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/progress"
	"shadcn.go/components/ui/slider"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/progress-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "progress", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("progress-controlled", NovaControlled),
		nova("progress-label", NovaLabel),
		nova("progress-rtl", NovaRtl),
	)
}

func NovaControlled() h.Node {
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-4"),
		progress.Progress(progress.Value(50)),
		slider.Slider(slider.DefaultValue(50), slider.Min(0), slider.Max(100), slider.Step(1)),
	)
}

func NovaLabel() h.Node {
	return field.Field(h.Class("w-full max-w-sm"),
		field.Label(h.For("progress-upload"),
			h.Span("Upload progress"),
			h.Span(h.Class("ml-auto"), "66%"),
		),
		progress.Progress(progress.Value(66), h.ID("progress-upload")),
	)
}

func NovaRtl() h.Node {
	return field.Field(h.Class("w-full max-w-sm"), h.Dir("rtl"),
		field.Label(h.For("progress-upload"),
			h.Span("تقدم الرفع"),
			h.Span(h.Class("ms-auto"), "٦٦", h.Comment(""), "%"),
		),
		progress.Progress(progress.Value(66), h.ID("progress-upload"), h.Class("rtl:rotate-180")),
	)
}
