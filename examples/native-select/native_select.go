// Package nativeselect ports the upstream native-select-* examples.
package nativeselect

import (
	nativeselect "shadcn.go/components/ui/native-select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "native-select", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("native-select-demo", Demo),
		ex("native-select-disabled", Disabled),
		ex("native-select-groups", Groups),
		ex("native-select-invalid", Invalid),
		ex("native-select-sm", Sm),
	)
}

func opt(value, label string) h.Node { return nativeselect.Option(h.Value(value), label) }

func Demo() h.Node {
	return nativeselect.NativeSelect(
		opt("", "Select status"),
		opt("todo", "Todo"),
		opt("in-progress", "In Progress"),
		opt("done", "Done"),
		opt("cancelled", "Cancelled"),
	)
}

func Disabled() h.Node {
	return nativeselect.NativeSelect(h.Disabled(),
		opt("", "Select priority"),
		opt("low", "Low"),
		opt("medium", "Medium"),
		opt("high", "High"),
		opt("critical", "Critical"),
	)
}

func Groups() h.Node {
	return nativeselect.NativeSelect(
		opt("", "Select department"),
		nativeselect.OptGroup(h.LabelAttr("Engineering"),
			opt("frontend", "Frontend"),
			opt("backend", "Backend"),
			opt("devops", "DevOps"),
		),
		nativeselect.OptGroup(h.LabelAttr("Sales"),
			opt("sales-rep", "Sales Rep"),
			opt("account-manager", "Account Manager"),
			opt("sales-director", "Sales Director"),
		),
		nativeselect.OptGroup(h.LabelAttr("Operations"),
			opt("support", "Customer Support"),
			opt("product-manager", "Product Manager"),
			opt("ops-manager", "Operations Manager"),
		),
	)
}

func Invalid() h.Node {
	return nativeselect.NativeSelect(h.Aria("invalid", "true"),
		opt("", "Select role"),
		opt("admin", "Admin"),
		opt("editor", "Editor"),
		opt("viewer", "Viewer"),
		opt("guest", "Guest"),
	)
}

// Sm ports the native-select-sm fixture (parity/fixtures/native-select-sm.tsx).
func Sm() h.Node {
	return nativeselect.NativeSelect(nativeselect.SizeSm, h.Class("w-40"),
		opt("", "Select size"),
		opt("s", "Small"),
		opt("m", "Medium"),
	)
}
