package sonner

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/toast"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the base toast page (ui/apps/v4/examples/base/toast-*.tsx):
// Base UI toasts (package toast), in base-nova. They live here, next to
// the sonner examples, as the toast component's examples.
func init() {
	reg := func(name string, fn func() h.Node) {
		registry.Register(registry.Example{Name: "base--" + name, Component: "toast", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})})
	}
	reg("toast-demo", ToastDemo)
	reg("toast-types", ToastTypes)
	reg("toast-promise", ToastPromise)
}

func ToastDemo() h.Node {
	return button.Button(button.VariantOutline, toast.Add(toast.Options{
		Title:       "Event created",
		Description: "Sunday, December 3 at 9:00 AM",
		Action:      "Undo",
	}), "Show Toast")
}

func ToastTypes() h.Node {
	b := func(label string, o toast.Options) h.Node {
		return button.Button(button.VariantOutline, toast.Add(o), label)
	}
	return h.Div(h.Class("flex flex-wrap gap-2"),
		b("Default", toast.Options{Description: "Event has been created."}),
		b("Success", toast.Options{Type: "success", Description: "Event has been created."}),
		b("Info", toast.Options{Type: "info", Description: "Arrive 10 minutes before the event."}),
		b("Warning", toast.Options{Type: "warning", Description: "The event cannot start before 8:00 AM."}),
		b("Error", toast.Options{Type: "error", Description: "The event could not be created.", Priority: "high"}),
	)
}

func ToastPromise() h.Node {
	return button.Button(button.VariantOutline, toast.Promise(toast.PromiseOptions{
		Loading: "Creating event…",
		Success: "{name} created.",
		Error:   "Could not create event.",
		Delay:   2000,
		Name:    "Event",
	}), "Create Event")
}
