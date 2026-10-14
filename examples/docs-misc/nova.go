// Package docsmisc ports docs previews (ui/apps/v4/examples/radix/*.tsx)
// whose names match no component: muted-item-group, outline-item-group,
// file-upload-list and radio-fields, in radix-nova.
package docsmisc

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/progress"
	rg "shadcn.go/components/ui/radio-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	nova := func(name, component string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: component, Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("muted-item-group", "item", MutedItemGroup),
		nova("outline-item-group", "item", OutlineItemGroup),
		nova("file-upload-list", "item", FileUploadList),
		nova("radio-fields", "radio-group", RadioFields),
	)
}

var ordinals = []string{"First", "Second", "Third"}

func MutedItemGroup() h.Node {
	g := item.Group()
	for i, o := range ordinals {
		g.Children(item.Item(item.VariantMuted,
			item.Content(item.Title("Item "+string(rune('1'+i))), item.Description(o+" item in muted group.")),
			item.Actions(button.Button(button.VariantOutline, button.SizeSm, "Action")),
		))
	}
	return g
}

func OutlineItemGroup() h.Node {
	g := item.Group()
	for i, o := range ordinals {
		g.Children(item.Item(item.VariantOutline,
			item.Media(item.MediaIcon, lucide.Inbox()),
			item.Content(item.Title("Item "+string(rune('1'+i))), item.Description(o+" item with icon.")),
		))
	}
	return g
}

func FileUploadList() h.Node {
	files := []struct {
		name     string
		progress float64
		left     string
	}{
		{"document.pdf", 45, "2m 30s"},
		{"presentation.pptx", 78, "45s"},
		{"spreadsheet.xlsx", 12, "5m 12s"},
		{"image.jpg", 100, "Complete"},
	}
	g := item.Group()
	for _, f := range files {
		g.Children(item.Item(item.SizeXs, h.Class("px-0"),
			item.Media(item.MediaIcon, lucide.File(h.Class("size-5"))),
			item.Content(h.Class("inline-block truncate"), item.Title(h.Class("inline"), f.name)),
			item.Content(progress.Progress(progress.Value(f.progress), h.Class("w-32"))),
			item.Actions(h.Class("w-16 justify-end"), h.Span(h.Class("text-sm text-muted-foreground"), f.left)),
		))
	}
	return g
}

func radio(value, id, label string, labelAttrs []any, attrs ...any) h.Node {
	return field.Field(append(attrs, field.Horizontal,
		rg.Item(rg.Value(value), h.ID(id)),
		field.Label(append(append([]any{h.For(id)}, labelAttrs...), label)...),
	)...)
}

func RadioFields() h.Node {
	normal := []any{h.Class("font-normal")}
	contentItem := func(value, id, label string) h.Node {
		return field.Field(field.Horizontal,
			rg.Item(rg.Value(value), h.ID(id)),
			field.Content(
				field.Label(h.For(id), label),
				field.Description("Enable Touch ID to quickly unlock your device."),
			),
		)
	}
	titleItem := func(value, id, title string) h.Node {
		return field.Label(h.For(id),
			field.Field(field.Horizontal,
				rg.Item(rg.Value(value), h.ID(id)),
				field.Content(
					field.Title(title),
					field.Description("Enable Touch ID to quickly unlock your device."),
				),
			),
		)
	}
	const long = "Enable Touch ID and Face ID to make it even faster to unlock your device. This is a long label to test the layout."
	invalid := func(value, id, label string) h.Node {
		return field.Field(h.Data("invalid", "true"), field.Horizontal,
			rg.Item(rg.Value(value), h.ID(id), h.Aria("invalid", "true")),
			field.Label(h.For(id), label),
		)
	}
	disabled := func(value, id, label string) h.Node {
		return field.Field(h.Data("disabled", "true"), field.Horizontal,
			rg.Item(rg.Value(value), h.ID(id), h.Disabled()),
			field.Label(h.For(id), label),
		)
	}
	return field.Group(
		field.Set(
			field.Legend(field.LegendLabel, "Subscription Plan"),
			rg.RadioGroup(rg.Value("free"),
				radio("free", "radio-free", "Free Plan", normal),
				radio("pro", "radio-pro", "Pro Plan", normal),
				radio("enterprise", "radio-enterprise", "Enterprise", normal),
			),
		),
		field.Set(
			field.Legend(field.LegendLabel, "Battery Level"),
			field.Description("Choose your preferred battery level."),
			rg.RadioGroup(
				radio("high", "battery-high", "High", nil),
				radio("medium", "battery-medium", "Medium", nil),
				radio("low", "battery-low", "Low", nil),
			),
		),
		rg.RadioGroup(h.Class("gap-6"),
			contentItem("option1", "radio-content-1", "Enable Touch ID"),
			contentItem("option2", "radio-content-2", long),
		),
		rg.RadioGroup(h.Class("gap-3"),
			titleItem("title1", "radio-title-1", "Enable Touch ID"),
			titleItem("title2", "radio-title-2", long),
		),
		field.Set(
			field.Legend(field.LegendLabel, "Invalid Radio Group"),
			rg.RadioGroup(
				invalid("invalid1", "radio-invalid-1", "Invalid Option 1"),
				invalid("invalid2", "radio-invalid-2", "Invalid Option 2"),
			),
		),
		field.Set(
			field.Legend(field.LegendLabel, "Disabled Radio Group"),
			rg.RadioGroup(rg.Disabled(true),
				disabled("disabled1", "radio-disabled-1", "Disabled Option 1"),
				disabled("disabled2", "radio-disabled-2", "Disabled Option 2"),
			),
		),
	)
}
