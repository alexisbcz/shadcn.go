package switchui

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/field"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/switch-*.tsx) whose content
// differs from the new-york-v4 example of the same name, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "switch", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("switch-choice-card", NovaChoiceCard),
		nova("switch-description", NovaDescription),
		nova("switch-disabled", NovaDisabled),
		nova("switch-invalid", NovaInvalid),
		nova("switch-sizes", NovaSizes),
		nova("switch-rtl", NovaRtl),
	)
}

func NovaChoiceCard() h.Node {
	card := func(id, title, desc string, checked bool) h.Node {
		return field.Label(h.For(id),
			field.Field(field.Horizontal,
				field.Content(
					field.Title(title),
					field.Description(desc),
				),
				switchui.Switch(h.ID(id), switchui.Checked(checked)),
			),
		)
	}
	return field.Group(h.Class("w-full max-w-sm"),
		card("switch-share", "Share across devices", "Focus is shared across devices, and turns off when you leave the app.", false),
		card("switch-notifications", "Enable notifications", "Receive notifications when focus mode is enabled or disabled.", true),
	)
}

func NovaDescription() h.Node {
	return field.Field(field.Horizontal, h.Class("max-w-sm"),
		field.Content(
			field.Label(h.For("switch-focus-mode"), "Share across devices"),
			field.Description("Focus is shared across devices, and turns off when you leave the app."),
		),
		switchui.Switch(h.ID("switch-focus-mode")),
	)
}

func NovaDisabled() h.Node {
	return field.Field(field.Horizontal, h.Data("disabled", "true"), h.Class("w-fit"),
		switchui.Switch(h.ID("switch-disabled-unchecked"), h.Disabled()),
		field.Label(h.For("switch-disabled-unchecked"), "Disabled"),
	)
}

func NovaInvalid() h.Node {
	return field.Field(field.Horizontal, h.Class("max-w-sm"), h.Data("invalid", "true"),
		field.Content(
			field.Label(h.For("switch-terms"), "Accept terms and conditions"),
			field.Description("You must accept the terms and conditions to continue."),
		),
		switchui.Switch(h.ID("switch-terms"), h.Aria("invalid", "true")),
	)
}

func NovaSizes() h.Node {
	return field.Group(h.Class("w-full max-w-[10rem]"),
		field.Field(field.Horizontal,
			switchui.Switch(h.ID("switch-size-sm"), switchui.SizeSm),
			field.Label(h.For("switch-size-sm"), "Small"),
		),
		field.Field(field.Horizontal,
			switchui.Switch(h.ID("switch-size-default"), switchui.SizeDefault),
			field.Label(h.For("switch-size-default"), "Default"),
		),
	)
}

func novaRtl() h.Node {
	return field.Field(field.Horizontal, h.Class("max-w-sm"), h.Dir("rtl"),
		field.Content(
			field.Label(h.For("switch-focus-mode-rtl"), h.Dir("rtl"), "المشاركة عبر الأجهزة"),
			// field does not map its classes under style.ScopeRTL yet:
			// text-start is the ui-rtl FieldDescription's text-left.
			field.Description(h.Dir("rtl"), h.Class("text-start"), "يتم مشاركة التركيز عبر الأجهزة، ويتم إيقاف تشغيله عند مغادرة التطبيق."),
		),
		switchui.Switch(h.ID("switch-focus-mode-rtl"), h.Dir("rtl")),
	)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
