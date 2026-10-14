package dialog

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the dialog page (ui/apps/v4/examples/radix/dialog-*.tsx)
// whose content differs from the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "dialog", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("dialog-demo", NovaDemo),
		nova("dialog-close-button", NovaCloseButton),
		nova("dialog-no-close-button", NovaNoCloseButton),
		nova("dialog-scrollable-content", NovaScrollableContent),
		nova("dialog-sticky-footer", NovaStickyFooter),
		nova("dialog-rtl", NovaRTL),
	)
}

// NovaProfileForm is the edit-profile dialog of dialog-demo and dialog-rtl.
func NovaProfileForm(t map[string]string, contentAttrs ...any) h.Node {
	return dialog.Dialog(
		h.Form(
			dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, t["openDialog"])),
			dialog.Content(append([]any{h.Class("sm:max-w-sm")}, append(contentAttrs,
				dialog.Header(
					dialog.Title(t["editProfile"]),
					dialog.Description(t["description"]),
				),
				field.Group(
					field.Field(
						label.Label(h.For("name-1"), t["name"]),
						input.Input(h.ID("name-1"), h.Name("name"), h.Value("Pedro Duarte")),
					),
					field.Field(
						label.Label(h.For("username-1"), t["username"]),
						input.Input(h.ID("username-1"), h.Name("username"), h.Value("@peduarte")),
					),
				),
				dialog.Footer(
					dialog.Close(lib.AsChild, button.Button(button.VariantOutline, t["cancel"])),
					button.Button(h.Type("submit"), t["saveChanges"]),
				),
			)...)...,
			),
		),
	)
}

func NovaDemo() h.Node {
	return NovaProfileForm(map[string]string{
		"openDialog":  "Open Dialog",
		"editProfile": "Edit profile",
		"description": "Make changes to your profile here. Click save when you're done.",
		"name":        "Name",
		"username":    "Username",
		"cancel":      "Cancel",
		"saveChanges": "Save changes",
	})
}

func NovaRTL() h.Node {
	return NovaProfileForm(map[string]string{
		"openDialog":  "فتح الحوار",
		"editProfile": "تعديل الملف الشخصي",
		"description": "قم بإجراء تغييرات على ملفك الشخصي هنا. انقر فوق حفظ عند الانتهاء.",
		"name":        "الاسم",
		"username":    "اسم المستخدم",
		"cancel":      "إلغاء",
		"saveChanges": "حفظ التغييرات",
	}, h.Dir("rtl"), h.Data("lang", "ar"))
}

func NovaCloseButton() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Share")),
		dialog.Content(h.Class("sm:max-w-md"),
			dialog.Header(
				dialog.Title("Share link"),
				dialog.Description("Anyone who has this link will be able to view this."),
			),
			h.Div(h.Class("flex items-center gap-2"),
				h.Div(h.Class("grid flex-1 gap-2"),
					label.Label(h.For("link"), h.Class("sr-only"), "Link"),
					input.Input(h.ID("link"), h.Value("https://ui.shadcn.com/docs/installation"), h.ReadOnly()),
				),
			),
			dialog.Footer(h.Class("sm:justify-start"),
				dialog.Close(lib.AsChild, button.Button(h.Type("button"), "Close")),
			),
		),
	)
}

func NovaNoCloseButton() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "No Close Button")),
		dialog.Content(dialog.ShowCloseButton(false),
			dialog.Header(
				dialog.Title("No Close Button"),
				dialog.Description("This dialog doesn't have a close button in the top-right corner."),
			),
		),
	)
}

const lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func novaLorem() h.Node {
	ps := []any{h.Class("-mx-4 no-scrollbar max-h-[50vh] overflow-y-auto px-4")}
	for range 10 {
		ps = append(ps, h.P(h.Class("mb-4 leading-normal"), lorem))
	}
	return h.Div(ps...)
}

func NovaScrollableContent() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Scrollable Content")),
		dialog.Content(
			dialog.Header(
				dialog.Title("Scrollable Content"),
				dialog.Description("This is a dialog with scrollable content."),
			),
			novaLorem(),
		),
	)
}

func NovaStickyFooter() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Sticky Footer")),
		dialog.Content(
			dialog.Header(
				dialog.Title("Sticky Footer"),
				dialog.Description("This dialog has a sticky footer that stays visible while the content scrolls."),
			),
			novaLorem(),
			dialog.Footer(
				dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Close")),
			),
		),
	)
}
