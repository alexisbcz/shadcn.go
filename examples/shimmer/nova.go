// Package shimmer ports the docs previews of the shimmer utility
// (ui/apps/v4/examples/radix/shimmer-*.tsx, shadcn/tailwind.css).
package shimmer

import (
	_ "embed"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/marker"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed replay.js
var replayScript string

func init() {
	runtime.Module("example-shimmer-replay", replayScript)
	nova := func(name string, st style.Style, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "shimmer", Render: registry.Static(func() h.Node {
			return style.Scope(st, fn())
		})}
	}
	registry.Register(
		nova("shimmer-angle", style.Nova, func() h.Node { return pair("shimmer", "shimmer-angle-45") }),
		nova("shimmer-color", style.Nova, NovaColor),
		nova("shimmer-demo", style.Nova, func() h.Node {
			return h.P(h.Class("shimmer text-sm text-muted-foreground"), generating)
		}),
		nova("shimmer-duration", style.Nova, func() h.Node { return pair("shimmer", "shimmer-duration-1000") }),
		// The marker and once previews import styles/radix-rhea.
		nova("shimmer-marker", style.Rhea, NovaMarker),
		nova("shimmer-none", style.Nova, func() h.Node {
			return h.Div(h.Class("flex flex-col items-center gap-3 text-sm text-muted-foreground"),
				h.P(h.Class("shimmer md:shimmer-none"), generating),
				h.P(h.Class("font-mono text-xs"), "shimmer md:shimmer-none"),
			)
		}),
		nova("shimmer-once", style.Rhea, NovaOnce),
		nova("shimmer-rtl", style.Nova, NovaRtl),
		nova("shimmer-spread", style.Nova, func() h.Node { return pair("shimmer-spread-4", "shimmer-spread-24") }),
	)
}

const generating = "Generating response…"

func cell(class, label string, attrs ...any) h.Node {
	return h.Div(h.Class("flex flex-col gap-3"),
		h.P(append(attrs, h.Class(class))...),
		h.P(h.Class("font-mono text-xs"), label),
	)
}

func grid(cells ...any) h.Node {
	return h.Div(append([]any{h.Class("mx-auto grid w-full max-w-lg gap-6 text-center text-sm text-muted-foreground sm:grid-cols-2")}, cells...)...)
}

// pair shows the default shimmer (or a) next to b, each with its label.
func pair(a, b string) h.Node {
	ca := a
	if a != "shimmer" {
		ca = "shimmer " + a
	}
	return grid(cell(ca, a, generating), cell("shimmer "+b, b, generating))
}

func NovaColor() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-2 text-sm text-muted-foreground"),
		h.P(h.Class("shimmer shimmer-color-blue-500/60"), generating),
		h.P(h.Class("shimmer shimmer-color-[#378ADD]"), generating),
	)
}

func NovaMarker() h.Node {
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-4"),
		marker.Marker(h.Role("status"),
			marker.Icon(spinner.Spinner()),
			marker.Content(h.Class("shimmer"), "Thinking..."),
		),
		marker.Marker(marker.VariantSeparator, h.Role("status"),
			marker.Content(h.Class("shimmer"), "Reading 4 files"),
		),
	)
}

func NovaOnce() h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-4"), h.Data("example-shimmer-replay", ""),
		h.P(h.Class("shimmer text-sm text-muted-foreground shimmer-duration-1100 shimmer-once"), generating),
		button.Button(button.VariantOutline, button.SizeSm, "Replay"),
	)
}

func NovaRtl() h.Node {
	return grid(
		cell("shimmer", `dir="ltr"`, h.Dir("ltr"), generating),
		cell("shimmer", `dir="rtl"`, h.Dir("rtl"), "جارٍ إنشاء الرد…"),
	)
}
