package nativeselect

import (
	"shadcn.go/components/style"
	nativeselect "shadcn.go/components/ui/native-select"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the native-select page
// (ui/apps/v4/examples/radix/native-select-*.tsx) whose content differs from
// the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "native-select", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	fruits := func(attrs ...any) h.Node {
		return nativeselect.NativeSelect(append(attrs,
			opt("apple", "Apple"), opt("banana", "Banana"), opt("blueberry", "Blueberry"))...)
	}
	registry.Register(
		nova("native-select-disabled", func() h.Node { return fruits(h.Disabled(), opt("", "Disabled")) }),
		nova("native-select-invalid", func() h.Node { return fruits(h.Aria("invalid", "true"), opt("", "Error state")) }),
		nova("native-select-rtl", func() h.Node {
			return style.ScopeRTL(nativeselect.NativeSelect(h.Dir("rtl"),
				opt("", "اختر الحالة"),
				opt("todo", "مهام"),
				opt("in-progress", "قيد التنفيذ"),
				opt("done", "منجز"),
				opt("cancelled", "ملغي"),
			))
		}),
	)
}
