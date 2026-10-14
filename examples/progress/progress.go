// Package progress ports the upstream progress-* examples.
package progress

import (
	_ "embed"

	"shadcn.go/components/ui/progress"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// demo.js sets progress-demo's value to 66 after 500ms, like the example's
// React effect.
//
//go:embed demo.js
var demoScript string

func init() {
	runtime.Module("example-progress-demo", demoScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "progress", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("progress-demo", Demo),
		ex("progress-states", States),
	)
}

func Demo() h.Node {
	return progress.Progress(progress.Value(13), h.Class("w-[60%]"), h.Data("example", "progress-demo"))
}

// States ports the progress-states fixture (parity/fixtures/progress-states.tsx):
// the Radix Progress props the upstream example does not use: no value
// (indeterminate), a complete bar, a custom max, a fractional value and
// an out-of-range value.
func States() h.Node {
	return h.Div(h.Class("flex w-[60%] flex-col gap-4"),
		progress.Progress(),
		progress.Progress(progress.Value(100)),
		progress.Progress(progress.Value(30), progress.Max(40)),
		progress.Progress(progress.Value(33.3)),
		progress.Progress(progress.Value(0)),
		progress.Progress(progress.Value(40), progress.Max(40), h.Class("h-4 bg-muted")),
	)
}
