package sheet

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/sheet"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the sheet page (ui/apps/v4/examples/radix/sheet-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "sheet", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("sheet-no-close-button", NovaNoCloseButton),
		nova("sheet-side", NovaSide),
		nova("sheet-rtl", NovaRTL),
	)
}

func NovaNoCloseButton() h.Node {
	return sheet.Sheet(
		sheet.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open Sheet")),
		sheet.Content(sheet.ShowCloseButton(false),
			sheet.Header(
				sheet.Title("No Close Button"),
				sheet.Description("This sheet doesn't have a close button in the top-right corner. Click outside to close."),
			),
		),
	)
}

const novaLorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func NovaSide() h.Node {
	row := []any{h.Class("flex flex-wrap gap-2")}
	for _, side := range []sheet.Side{sheet.SideTop, sheet.SideRight, sheet.SideBottom, sheet.SideLeft} {
		ps := []any{h.Class("no-scrollbar overflow-y-auto px-4")}
		for range 10 {
			ps = append(ps, h.P(h.Class("mb-2 leading-relaxed"), novaLorem))
		}
		row = append(row, sheet.Sheet(
			sheet.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("capitalize"), string(side))),
			sheet.Content(side, h.Class("data-[side=bottom]:max-h-[50vh] data-[side=top]:max-h-[50vh]"),
				sheet.Header(
					sheet.Title("Edit profile"),
					sheet.Description("Make changes to your profile here. Click save when you're done."),
				),
				h.Div(ps...),
				sheet.Footer(
					button.Button(h.Type("submit"), "Save changes"),
					sheet.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
				),
			),
		))
	}
	return h.Div(row...)
}

// NovaRTL is sheet-rtl in its default language, Arabic.
func NovaRTL() h.Node {
	return sheet.Sheet(
		sheet.Trigger(lib.AsChild, button.Button(button.VariantOutline, "فتح")),
		sheet.Content(h.Dir("rtl"), sheet.SideLeft, h.Data("lang", "ar"),
			sheet.Header(
				sheet.Title("تعديل الملف الشخصي"),
				sheet.Description("قم بإجراء تغييرات على ملفك الشخصي هنا. انقر حفظ عند الانتهاء."),
			),
			field.Group(h.Class("px-4"),
				field.Field(
					field.Label(h.For("sheet-rtl-name"), "الاسم"),
					input.Input(h.ID("sheet-rtl-name"), h.Value("Pedro Duarte")),
				),
				field.Field(
					field.Label(h.For("sheet-rtl-username"), "اسم المستخدم"),
					input.Input(h.ID("sheet-rtl-username"), h.Value("peduarte")),
				),
			),
			sheet.Footer(
				button.Button(h.Type("submit"), "حفظ التغييرات"),
				sheet.Close(lib.AsChild, button.Button(button.VariantOutline, "إغلاق")),
			),
		),
	)
}
