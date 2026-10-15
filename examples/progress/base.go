package progress

import (
	"shadcn.go/components/style"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The base and aria progress-label and progress-rtl previews
// (ui/apps/v4/examples/{base,aria}/progress-*.tsx): Base UI's and React
// Aria's progress, whose label and value sit above the track in the root.
func init() {
	for _, prefix := range []string{"base--", "aria--"} {
		registry.Register(
			registry.Example{Name: prefix + "progress-label", Component: "progress", Render: registry.Static(func() h.Node { return style.Scope(style.Nova, BaseLabel()) })},
			registry.Example{Name: prefix + "progress-rtl", Component: "progress", Render: registry.Static(func() h.Node { return style.Scope(style.Nova, BaseRtl()) })},
		)
	}
}

// baseProgress is the nova Progress with a label and a value (the value
// node as given): the root wraps label and value, then the track.
func baseProgress(value string, dir string, label, valueNode h.Node, valueClass string) h.Node {
	root := h.Div(h.Data("slot", "progress"), h.Role("progressbar"),
		h.Aria("valuemin", "0"), h.Aria("valuemax", "100"), h.Aria("valuenow", value), h.Aria("valuetext", value+"%"),
		h.Class("flex flex-wrap gap-3 w-full max-w-sm"),
		h.Span(h.Data("slot", "progress-label"), h.Class("text-sm font-medium"), label),
		h.Span(h.Data("slot", "progress-value"), h.Aria("hidden", "true"), h.Class(valueClass+" text-sm text-muted-foreground tabular-nums"), valueNode),
		h.Div(h.Data("slot", "progress-track"), h.Class("relative flex h-1 w-full items-center overflow-x-hidden rounded-full bg-muted"),
			h.Div(h.Data("slot", "progress-indicator"), h.Class("h-full bg-primary transition-all"),
				h.Style("inset-inline-start:0;height:inherit;width:"+value+"%")),
		),
	)
	if dir != "" {
		root.Children(h.Dir(dir))
	}
	return root
}

// BaseLabel is examples/base/progress-label.tsx.
func BaseLabel() h.Node {
	return baseProgress("56", "", h.Text("Upload progress"), h.Text("56%"), "ml-auto")
}

// BaseRtl is examples/base/progress-rtl.tsx.
func BaseRtl() h.Node {
	return baseProgress("56", "rtl", h.Text("تقدم الرفع"), h.Span(h.Class("ms-auto"), "٥٦%"), "ms-auto")
}
