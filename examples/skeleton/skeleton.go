// Package skeleton ports the upstream skeleton-* examples.
package skeleton

import (
	"shadcn.go/components/ui/skeleton"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "skeleton", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("skeleton-card", Card),
		ex("skeleton-demo", Demo),
	)
}

func Card() h.Node {
	return h.Div(h.Class("flex flex-col space-y-3"),
		skeleton.Skeleton(h.Class("h-[125px] w-[250px] rounded-xl")),
		h.Div(h.Class("space-y-2"),
			skeleton.Skeleton(h.Class("h-4 w-[250px]")),
			skeleton.Skeleton(h.Class("h-4 w-[200px]")),
		),
	)
}

func Demo() h.Node {
	return h.Div(h.Class("flex items-center space-x-4"),
		skeleton.Skeleton(h.Class("h-12 w-12 rounded-full")),
		h.Div(h.Class("space-y-2"),
			skeleton.Skeleton(h.Class("h-4 w-[250px]")),
			skeleton.Skeleton(h.Class("h-4 w-[200px]")),
		),
	)
}
