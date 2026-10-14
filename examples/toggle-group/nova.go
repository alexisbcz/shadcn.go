package togglegroup

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/field"
	tg "shadcn.go/components/ui/toggle-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/toggle-group-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "toggle-group", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("toggle-group-outline", NovaOutline),
		nova("toggle-group-sizes", NovaSizes),
		nova("toggle-group-spacing", NovaSpacing),
		nova("toggle-group-vertical", NovaVertical),
		nova("toggle-group-font-weight-selector", NovaFontWeightSelector),
		nova("toggle-group-rtl", NovaRtl),
	)
}

func NovaOutline() h.Node {
	return tg.ToggleGroup(tg.VariantOutline, tg.TypeSingle, tg.Value("all"),
		tg.Item(tg.Value("all"), h.Aria("label", "Toggle all"), "All"),
		tg.Item(tg.Value("missed"), h.Aria("label", "Toggle missed"), "Missed"),
	)
}

func positions() []any {
	var out []any
	for _, p := range [][2]string{{"top", "Top"}, {"bottom", "Bottom"}, {"left", "Left"}, {"right", "Right"}} {
		out = append(out, tg.Item(tg.Value(p[0]), h.Aria("label", "Toggle "+p[0]), p[1]))
	}
	return out
}

func NovaSizes() h.Node {
	return h.Div(h.Class("flex flex-col gap-4"),
		tg.ToggleGroup(tg.TypeSingle, tg.SizeSm, tg.Value("top"), tg.VariantOutline, positions()),
		tg.ToggleGroup(tg.TypeSingle, tg.Value("top"), tg.VariantOutline, positions()),
	)
}

func NovaSpacing() h.Node {
	return tg.ToggleGroup(tg.TypeSingle, tg.SizeSm, tg.Value("top"), tg.VariantOutline, tg.Spacing(2), positions())
}

func NovaVertical() h.Node {
	return tg.ToggleGroup(tg.TypeMultiple, tg.OrientationVertical, tg.Spacing(1), tg.Value("bold"), tg.Value("italic"),
		tg.Item(tg.Value("bold"), h.Aria("label", "Toggle bold"), lucide.Bold()),
		tg.Item(tg.Value("italic"), h.Aria("label", "Toggle italic"), lucide.Italic()),
		tg.Item(tg.Value("underline"), h.Aria("label", "Toggle underline"), lucide.Underline()),
	)
}

func NovaFontWeightSelector() h.Node {
	item := func(value, label, weight string) h.Node {
		return tg.Item(tg.Value(value), h.Aria("label", label), h.Class("flex size-16 flex-col items-center justify-center rounded-xl"),
			h.Span(h.Class("text-2xl leading-none "+weight), "Aa"),
			h.Span(h.Class("text-xs text-muted-foreground"), label),
		)
	}
	return field.Field(
		field.Label("Font Weight"),
		tg.ToggleGroup(tg.TypeSingle, tg.Value("normal"), tg.VariantOutline, tg.Spacing(2), tg.SizeLg,
			item("light", "Light", "font-light"),
			item("normal", "Normal", "font-normal"),
			item("medium", "Medium", "font-medium"),
			item("bold", "Bold", "font-bold"),
		),
		field.Description(
			// React renders separate text nodes: split them with comments.
			"Use", h.Comment(""), " ",
			h.Code(h.Class("rounded-md bg-muted px-1 py-0.5 font-mono"), "font-", h.Comment(""), "normal"),
			" ", h.Comment(""), "to set the font weight.",
		),
	)
}

func novaRtl() h.Node {
	return tg.ToggleGroup(tg.VariantOutline, tg.TypeSingle, tg.Value("list"),
		tg.Item(tg.Value("list"), h.Aria("label", "قائمة"), "قائمة"),
		tg.Item(tg.Value("grid"), h.Aria("label", "شبكة"), "شبكة"),
		tg.Item(tg.Value("cards"), h.Aria("label", "بطاقات"), "بطاقات"),
	)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
