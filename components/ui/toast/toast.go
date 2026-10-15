// Package toast is the Toast of the Base UI styles
// (styles/base-nova/ui/toast.tsx): Base UI's toast manager rendering
// stacked toasts in a viewport at the bottom right. The behavior lives in
// toast.js, which exports the manager:
//
//	import { toast } from "/_shadcn/toast.js"
//	const id = toast.add({ title: "Saved", description: "...", type: "success" })
//	toast.close(id)
//	toast.promise(fetch("/api"), { loading: "Saving…", success: "Saved.", error: "Failed." })
//
// Server-rendered buttons can show toasts declaratively with Add and
// Promise. Render one Toaster per page (the root layout's <Toaster />);
// without one, toast.js creates the viewport on first use.
package toast

import (
	_ "embed"
	"encoding/json"

	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed toast.js
var script string

func init() { runtime.Module("toast", script) }

// ViewportClass is the class of ToastViewport.
const ViewportClass = "pointer-events-none fixed inset-x-4 bottom-4 z-50 mx-auto w-auto max-w-sm outline-none sm:right-4 sm:left-auto sm:mx-0 sm:w-full"

// Toaster renders the toast portal and its (empty) viewport.
func Toaster() h.Node {
	return h.Div(h.Data("base-ui-portal", ""), h.Data("slot", "toast-portal"),
		h.Div(h.Attr("tabindex", "-1"), h.Attr("role", "region"), h.Attr("aria-live", "polite"), h.Attr("aria-atomic", "false"),
			h.Attr("aria-relevant", "additions text"), h.Attr("aria-label", "Notifications"), h.Data("slot", "toast-viewport"),
			h.Class(ViewportClass)),
	)
}

// Options are the options of toast.add.
type Options struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	// Type is "success", "info", "warning", "error" or "loading".
	Type string `json:"type,omitempty"`
	// Priority "high" announces the toast as an alert.
	Priority string `json:"priority,omitempty"`
	// Action is the label of an action button closing the toast.
	Action string `json:"action,omitempty"`
}

// Add is a prop showing a toast when the element is clicked.
func Add(o Options) h.Prop {
	b, _ := json.Marshal(o)
	return h.Data("toast-add", string(b))
}

// PromiseOptions are the messages of toast.promise; Success may contain
// "{name}", replaced by Name (the resolved value's name).
type PromiseOptions struct {
	Loading string `json:"loading"`
	Success string `json:"success"`
	Error   string `json:"error"`
	// Delay is how long the simulated work takes, in milliseconds.
	Delay int    `json:"delay"`
	Name  string `json:"name,omitempty"`
}

// Promise is a prop showing a loading toast when the element is clicked,
// turning into a success toast after Delay milliseconds.
func Promise(o PromiseOptions) h.Prop {
	b, _ := json.Marshal(o)
	return h.Data("toast-promise", string(b))
}
