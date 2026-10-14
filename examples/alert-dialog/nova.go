package alertdialog

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	ad "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the alert-dialog page
// (ui/apps/v4/examples/radix/alert-dialog-*.tsx) whose content differs from
// the new-york-v4 examples, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "alert-dialog", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("alert-dialog-demo", NovaDemo),
		nova("alert-dialog-basic", NovaBasic),
		nova("alert-dialog-destructive", NovaDestructive),
		nova("alert-dialog-media", NovaMedia),
		nova("alert-dialog-small", NovaSmall),
		nova("alert-dialog-small-media", NovaSmallMedia),
		nova("alert-dialog-rtl", NovaRTL),
	)
}

// novaAlert is the common shape of the previews: a trigger button and a
// content with an optional media, a title, a description and two actions.
func novaAlert(trigger *button.Builder, content []any, media any, title, description any, cancel, action any) h.Node {
	return novaAlertHeader(trigger, content, nil, media, title, description, cancel, action)
}

func novaAlertHeader(trigger *button.Builder, content, header []any, media any, title, description any, cancel, action any) h.Node {
	if media != nil {
		header = append(header, media)
	}
	header = append(header, ad.Title(title), ad.Description(description))
	return ad.AlertDialog(
		ad.Trigger(lib.AsChild, trigger),
		ad.Content(append(content,
			ad.Header(header...),
			ad.Footer(cancel, action),
		)...),
	)
}

func NovaDemo() h.Node {
	return novaAlert(button.Button(button.VariantOutline, "Show Dialog"), nil, nil,
		"Are you absolutely sure?",
		"This action cannot be undone. This will permanently delete your account from our servers.",
		ad.Cancel("Cancel"), ad.Action("Continue"))
}

func NovaBasic() h.Node {
	return novaAlert(button.Button(button.VariantOutline, "Show Dialog"), nil, nil,
		"Are you absolutely sure?",
		"This action cannot be undone. This will permanently delete your account and remove your data from our servers.",
		ad.Cancel("Cancel"), ad.Action("Continue"))
}

func NovaDestructive() h.Node {
	return novaAlert(button.Button(button.VariantDestructive, "Delete Chat"), []any{ad.SizeSm},
		ad.Media(h.Class("bg-destructive/10 text-destructive dark:bg-destructive/20 dark:text-destructive"), lucide.Trash2()),
		"Delete chat?",
		h.Group("This will permanently delete this chat conversation. View ", h.A(h.Href("#"), "Settings"), " delete any memories saved during this chat."),
		ad.Cancel(button.VariantOutline, "Cancel"), ad.Action(button.VariantDestructive, "Delete"))
}

func NovaMedia() h.Node {
	return novaAlert(button.Button(button.VariantOutline, "Share Project"), nil,
		ad.Media(lucide.CircleFadingPlus()),
		"Share this project?",
		"Anyone with the link will be able to view and edit this project.",
		ad.Cancel("Cancel"), ad.Action("Share"))
}

func NovaSmall() h.Node {
	return novaAlert(button.Button(button.VariantOutline, "Show Dialog"), []any{ad.SizeSm}, nil,
		"Allow accessory to connect?",
		"Do you want to allow the USB accessory to connect to this device?",
		ad.Cancel("Don't allow"), ad.Action("Allow"))
}

func NovaSmallMedia() h.Node {
	return novaAlert(button.Button(button.VariantOutline, "Show Dialog"), []any{ad.SizeSm},
		ad.Media(lucide.Bluetooth()),
		"Allow accessory to connect?",
		"Do you want to allow the USB accessory to connect to this device?",
		ad.Cancel("Don't allow"), ad.Action("Allow"))
}

// NovaRTL is alert-dialog-rtl in its default language, Arabic.
func NovaRTL() h.Node {
	rtl := []any{h.Dir("rtl"), h.Data("lang", "ar")}
	// styles/radix-nova/ui-rtl/alert-dialog.tsx: text-start for text-left
	// (the Go components have no RTL transform).
	header := []any{h.Class("sm:group-data-[size=default]/alert-dialog-content:text-start")}
	return h.Div(h.Class("flex gap-4"), h.Dir("rtl"),
		novaAlertHeader(button.Button(button.VariantOutline, "إظهار الحوار"), rtl, header, nil,
			"هل أنت متأكد تمامًا؟",
			"لا يمكن التراجع عن هذا الإجراء. سيؤدي هذا إلى حذف حسابك نهائيًا من خوادمنا.",
			ad.Cancel("إلغاء"), ad.Action("متابعة")),
		novaAlertHeader(button.Button(button.VariantOutline, "إظهار الحوار (صغير)"), append([]any{ad.SizeSm}, rtl...), header,
			ad.Media(lucide.Bluetooth()),
			"السماح للملحق بالاتصال؟",
			"هل تريد السماح لملحق USB بالاتصال بهذا الجهاز؟",
			ad.Cancel("عدم السماح"), ad.Action("السماح")),
	)
}
