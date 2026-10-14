package checkbox

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/table"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/checkbox-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "checkbox", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("checkbox-basic", NovaBasic),
		nova("checkbox-demo", NovaDemo),
		nova("checkbox-description", NovaDescription),
		nova("checkbox-disabled", NovaDisabled),
		nova("checkbox-group", NovaGroup),
		nova("checkbox-invalid", NovaInvalid),
		nova("checkbox-table", NovaTable),
		nova("checkbox-rtl", NovaRtl),
	)
}

func NovaBasic() h.Node {
	return field.Group(h.Class("mx-auto w-56"),
		field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("terms-checkbox-basic"), h.Name("terms-checkbox-basic")),
			field.Label(h.For("terms-checkbox-basic"), "Accept terms and conditions"),
		),
	)
}

type demoText struct {
	acceptTerms, acceptTermsDescription, enableNotifications, enableNotificationsDescription string
}

func novaDemo(t demoText, suffix string, attrs ...any) h.Node {
	return field.Group(append(attrs,
		field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("terms-checkbox"+suffix), h.Name("terms-checkbox")),
			label.Label(h.For("terms-checkbox"+suffix), t.acceptTerms),
		),
		field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("terms-checkbox-2"+suffix), h.Name("terms-checkbox-2"), checkbox.Checked(true)),
			field.Content(
				field.Label(h.For("terms-checkbox-2"+suffix), t.acceptTerms),
				field.Description(t.acceptTermsDescription),
			),
		),
		field.Field(field.Horizontal, h.Data("disabled", "true"),
			checkbox.Checkbox(h.ID("toggle-checkbox"+suffix), h.Name("toggle-checkbox"), h.Disabled()),
			field.Label(h.For("toggle-checkbox"+suffix), t.enableNotifications),
		),
		field.Label(
			field.Field(field.Horizontal,
				checkbox.Checkbox(h.ID("toggle-checkbox-2"), h.Name("toggle-checkbox-2")),
				field.Content(
					field.Title(t.enableNotifications),
					field.Description(t.enableNotificationsDescription),
				),
			),
		),
	)...)
}

func NovaDemo() h.Node {
	return novaDemo(demoText{
		"Accept terms and conditions",
		"By clicking this checkbox, you agree to the terms.",
		"Enable notifications",
		"You can enable or disable notifications at any time.",
	}, "", h.Class("max-w-sm"))
}

func NovaRtl() h.Node {
	return novaDemo(demoText{
		"قبول الشروط والأحكام",
		"بالنقر على هذا المربع، فإنك توافق على الشروط.",
		"تفعيل الإشعارات",
		"يمكنك تفعيل أو إلغاء تفعيل الإشعارات في أي وقت.",
	}, "-rtl", h.Class("max-w-sm"), h.Dir("rtl"))
}

func NovaDescription() h.Node {
	return field.Group(h.Class("mx-auto w-72"),
		field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("terms-checkbox-desc"), h.Name("terms-checkbox-desc"), checkbox.Checked(true)),
			field.Content(
				field.Label(h.For("terms-checkbox-desc"), "Accept terms and conditions"),
				field.Description("By clicking this checkbox, you agree to the terms and conditions."),
			),
		),
	)
}

func NovaDisabled() h.Node {
	return field.Group(h.Class("mx-auto w-56"),
		field.Field(field.Horizontal, h.Data("disabled", "true"),
			checkbox.Checkbox(h.ID("toggle-checkbox-disabled"), h.Name("toggle-checkbox-disabled"), h.Disabled()),
			field.Label(h.For("toggle-checkbox-disabled"), "Enable notifications"),
		),
	)
}

func NovaGroup() h.Node {
	item := func(id, text string, checked bool) h.Node {
		return field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID(id), h.Name(id), checkbox.Checked(checked)),
			field.Label(h.For(id), h.Class("font-normal"), text),
		)
	}
	return field.Set(
		field.Legend(field.LegendLabel, "Show these items on the desktop:"),
		field.Description("Select the items you want to show on the desktop."),
		field.Group(h.Class("gap-3"),
			item("finder-pref-9k2-hard-disks-ljj-checkbox", "Hard disks", true),
			item("finder-pref-9k2-external-disks-1yg-checkbox", "External disks", true),
			item("finder-pref-9k2-cds-dvds-fzt-checkbox", "CDs, DVDs, and iPods", false),
			item("finder-pref-9k2-connected-servers-6l2-checkbox", "Connected servers", false),
		),
	)
}

func NovaInvalid() h.Node {
	return field.Group(h.Class("mx-auto w-56"),
		field.Field(field.Horizontal, h.Data("invalid", "true"),
			checkbox.Checkbox(h.ID("terms-checkbox-invalid"), h.Name("terms-checkbox-invalid"), h.Aria("invalid", "true")),
			field.Label(h.For("terms-checkbox-invalid"), "Accept terms and conditions"),
		),
	)
}

func NovaTable() h.Node {
	rows := [][4]string{
		{"1", "Sarah Chen", "sarah.chen@example.com", "Admin"},
		{"2", "Marcus Rodriguez", "marcus.rodriguez@example.com", "User"},
		{"3", "Priya Patel", "priya.patel@example.com", "User"},
		{"4", "David Kim", "david.kim@example.com", "Editor"},
	}
	var body []any
	for _, r := range rows {
		selected := r[0] == "1"
		row := table.Row(
			table.Cell(checkbox.Checkbox(h.ID("row-"+r[0]+"-checkbox"), h.Name("row-"+r[0]+"-checkbox"), checkbox.Checked(selected))),
			table.Cell(h.Class("font-medium"), r[1]),
			table.Cell(r[2]),
			table.Cell(r[3]),
		)
		if selected {
			row = table.Row(h.Data("state", "selected"),
				table.Cell(checkbox.Checkbox(h.ID("row-"+r[0]+"-checkbox"), h.Name("row-"+r[0]+"-checkbox"), checkbox.Checked(true))),
				table.Cell(h.Class("font-medium"), r[1]),
				table.Cell(r[2]),
				table.Cell(r[3]),
			)
		}
		body = append(body, row)
	}
	return table.Table(
		table.Header(table.Row(
			table.Head(h.Class("w-8"), checkbox.Checkbox(h.ID("select-all-checkbox"), h.Name("select-all-checkbox"))),
			table.Head("Name"),
			table.Head("Email"),
			table.Head("Role"),
		)),
		table.Body(body...),
	)
}
