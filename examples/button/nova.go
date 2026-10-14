package button

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the button page (ui/apps/v4/examples/radix/button-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "button", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("button-aschild", AsChild),
		nova("button-rtl", NovaRtl),
		nova("button-size", NovaSize),
		nova("button-spinner", NovaSpinner),
	)
}

func NovaSize() h.Node {
	pair := func(s button.Size, label string, icon button.Size) h.Node {
		return h.Div(h.Class("flex items-start gap-2"),
			button.Button(s, button.VariantOutline, label),
			button.Button(icon, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
		)
	}
	return h.Div(h.Class("flex flex-col items-start gap-8 sm:flex-row"),
		pair(button.SizeXs, "Extra Small", button.SizeIconXs),
		pair(button.SizeSm, "Small", button.SizeIconSm),
		pair(button.SizeDefault, "Default", button.SizeIcon),
		pair(button.SizeLg, "Large", button.SizeIconLg),
	)
}

func NovaSpinner() h.Node {
	return h.Div(h.Class("flex gap-2"),
		button.Button(button.VariantOutline, h.Bool("disabled"), spinner.Spinner(h.Data("icon", "inline-start")), "Generating"),
		button.Button(button.VariantSecondary, h.Bool("disabled"), "Downloading", spinner.Spinner(h.Data("icon", "inline-start"))),
	)
}

func NovaRtl() h.Node {
	return style.ScopeRTL(h.Div(h.Class("flex flex-wrap items-center gap-2 md:flex-row"), h.Dir("rtl"),
		button.Button(button.VariantOutline, "زر"),
		button.Button(button.VariantDestructive, "حذف"),
		button.Button(button.VariantOutline, "إرسال", " ", lucide.ArrowRight(h.Class("rtl:rotate-180"), h.Data("icon", "inline-end"))),
		button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Add"), lucide.Plus()),
		button.Button(button.VariantSecondary, h.Bool("disabled"), spinner.Spinner(h.Data("icon", "inline-start")), " ", "جاري التحميل"),
	))
}
