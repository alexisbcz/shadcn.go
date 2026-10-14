// Package dropdownmenu ports the upstream dropdown-menu-* examples, the
// mode-toggle example and the dropdown-menu parity fixtures.
package dropdownmenu

import (
	_ "embed"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// dialog.js opens dropdown-menu-dialog's dialogs from the menu items'
// onSelect (React state of the example).
//
//go:embed dialog.js
var dialogScript string

func init() {
	runtime.Module("example-dropdown-menu-dialog", dialogScript)
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "dropdown-menu", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("dropdown-menu-checkboxes", Checkboxes),
		ex("dropdown-menu-demo", Demo),
		ex("dropdown-menu-dialog", Dialog),
		ex("dropdown-menu-radio-group", RadioGroup),
		ex("mode-toggle", ModeToggle),
		ex("dropdown-menu-nested", Nested),
		ex("dropdown-menu-sides", Sides),
		ex("dropdown-menu-rtl", Rtl),
		ex("dropdown-menu-default-open", DefaultOpen),
	)
}

func Checkboxes() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		dm.Content(h.Class("w-56"),
			dm.Label("Appearance"),
			dm.Separator(),
			dm.CheckboxItem(dm.Checked(true), "Status Bar"),
			dm.CheckboxItem(dm.Checked(false), dm.Disabled(true), "Activity Bar"),
			dm.CheckboxItem(dm.Checked(false), "Panel"),
		),
	)
}

func Demo() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		dm.Content(h.Class("w-56"), dm.AlignStart,
			dm.Label("My Account"),
			dm.Group(
				dm.Item("Profile", dm.Shortcut("⇧⌘P")),
				dm.Item("Billing", dm.Shortcut("⌘B")),
				dm.Item("Settings", dm.Shortcut("⌘S")),
				dm.Item("Keyboard shortcuts", dm.Shortcut("⌘K")),
			),
			dm.Separator(),
			dm.Group(
				dm.Item("Team"),
				dm.Sub(
					dm.SubTrigger("Invite users"),
					dm.Portal(
						dm.SubContent(
							dm.Item("Email"),
							dm.Item("Message"),
							dm.Separator(),
							dm.Item("More..."),
						),
					),
				),
				dm.Item("New Team", dm.Shortcut("⌘+T")),
			),
			dm.Separator(),
			dm.Item("GitHub"),
			dm.Item("Support"),
			dm.Item(dm.Disabled(true), "API"),
			dm.Separator(),
			dm.Item("Log out", dm.Shortcut("⇧⌘Q")),
		),
	)
}

func Dialog() h.Node {
	return h.Group(
		dm.DropdownMenu(dm.Modal(false),
			dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Aria("label", "Open menu"), button.SizeIconSm, lucide.Ellipsis())),
			dm.Content(h.Class("w-40"), dm.AlignEnd,
				dm.Label("File Actions"),
				dm.Group(
					dm.Item("New File..."),
					dm.Item("Share..."),
					dm.Item(dm.Disabled(true), "Download"),
				),
			),
		),
		dialog.Dialog(
			dialog.Content(h.Class("sm:max-w-[425px]"),
				dialog.Header(
					dialog.Title("Create New File"),
					dialog.Description("Provide a name for your new file. Click create when you're done."),
				),
				field.Group(h.Class("pb-3"),
					field.Field(
						field.Label(h.For("filename"), "File Name"),
						input.Input(h.ID("filename"), h.Name("filename"), h.Placeholder("document.txt")),
					),
				),
				dialog.Footer(
					dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
					button.Button(h.Type("submit"), "Create"),
				),
			),
		),
		dialog.Dialog(
			dialog.Content(h.Class("sm:max-w-[425px]"),
				dialog.Header(
					dialog.Title("Share File"),
					dialog.Description("Anyone with the link will be able to view this file."),
				),
				field.Group(h.Class("py-3"),
					field.Field(
						label.Label(h.For("email"), "Email Address"),
						input.Input(h.ID("email"), h.Name("email"), h.Type("email"), h.Placeholder("shadcn@vercel.com"), h.Attr("autocomplete", "off")),
					),
					field.Field(
						field.Label(h.For("message"), "Message (Optional)"),
						textarea.Textarea(h.ID("message"), h.Name("message"), h.Placeholder("Check out this file")),
					),
				),
				dialog.Footer(
					dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
					button.Button(h.Type("submit"), "Send Invite"),
				),
			),
		),
	)
}

func RadioGroup() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		dm.Content(h.Class("w-56"),
			dm.Label("Panel Position"),
			dm.Separator(),
			dm.RadioGroup(dm.Value("bottom"),
				dm.RadioItem(dm.Value("top"), "Top"),
				dm.RadioItem(dm.Value("bottom"), "Bottom"),
				dm.RadioItem(dm.Value("right"), "Right"),
			),
		),
	)
}

// ModeToggle: upstream calls next-themes' setTheme, which the docs (and the
// parity reference) force to the page's theme, so selecting only closes the
// menu.
func ModeToggle() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIcon,
			lucide.Sun(h.Class("h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90")),
			lucide.Moon(h.Class("absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0")),
			h.Span(h.Class("sr-only"), "Toggle theme"),
		)),
		dm.Content(dm.AlignEnd,
			dm.Item("Light"),
			dm.Item("Dark"),
			dm.Item("System"),
		),
	)
}

// Nested is the parity fixture dropdown-menu-nested.
func Nested() h.Node {
	return dm.DropdownMenu(
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Nested")),
		dm.Content(h.Class("w-56"), dm.Loop(true),
			dm.Label(dm.Inset(true), "Format"),
			dm.CheckboxItem(dm.Indeterminate, "Bold", dm.Shortcut("⌘B")),
			dm.CheckboxItem(dm.Checked(false), "Italic"),
			dm.Separator(),
			dm.Sub(
				dm.SubTrigger(dm.Inset(true), "Level 1"),
				dm.SubContent(dm.Loop(true),
					dm.Item("Alpha"),
					dm.Sub(
						dm.SubTrigger("Level 2"),
						dm.Portal(
							dm.SubContent(
								dm.RadioGroup(dm.Value("m"),
									dm.RadioItem(dm.Value("s"), "Small"),
									dm.RadioItem(dm.Value("m"), "Medium"),
									dm.RadioItem(dm.Value("l"), dm.Disabled(true), "Large"),
								),
								dm.Separator(),
								dm.Sub(
									dm.SubTrigger("Level 3"),
									dm.SubContent(
										dm.Item("Deep item"),
										dm.Item(dm.VariantDestructive, "Delete", dm.Shortcut("⌘⌫")),
									),
								),
							),
						),
					),
					dm.Item(dm.Disabled(true), "Disabled"),
				),
			),
			dm.Sub(
				dm.SubTrigger(dm.Disabled(true), "Disabled sub"),
				dm.SubContent(dm.Item("Never shown")),
			),
			dm.Item(dm.TextValue("Zebra"), "Custom text value"),
			dm.Item(lib.AsChild, h.A(h.Href("#docs"), "Docs link")),
		),
	)
}

// Sides is the parity fixture dropdown-menu-sides.
func Sides() h.Node {
	type variant struct {
		side        dm.Side
		align       dm.Align
		nonModal    bool
		sideOffset  *float64
		alignOffset *float64
	}
	f := func(v float64) *float64 { return &v }
	variants := []variant{
		{dm.SideTop, dm.AlignStart, false, nil, nil},
		{dm.SideTop, dm.AlignCenter, true, nil, nil},
		{dm.SideTop, dm.AlignEnd, false, f(12), nil},
		{dm.SideRight, dm.AlignStart, false, nil, f(8)},
		{dm.SideRight, dm.AlignCenter, true, nil, nil},
		{dm.SideRight, dm.AlignEnd, false, nil, nil},
		{dm.SideBottom, dm.AlignStart, false, f(0), nil},
		{dm.SideBottom, dm.AlignCenter, true, nil, nil},
		{dm.SideBottom, dm.AlignEnd, false, nil, f(-6)},
		{dm.SideLeft, dm.AlignStart, false, nil, nil},
		{dm.SideLeft, dm.AlignCenter, true, f(10), nil},
		{dm.SideLeft, dm.AlignEnd, false, nil, nil},
	}
	return h.Div(h.Class("grid grid-cols-3 gap-x-24 gap-y-6"),
		h.Map(variants, func(v variant) h.Node {
			opts := []any{h.Class("w-40"), v.side, v.align}
			if v.sideOffset != nil {
				opts = append(opts, dm.SideOffset(*v.sideOffset))
			}
			if v.alignOffset != nil {
				opts = append(opts, dm.AlignOffset(*v.alignOffset))
			}
			label := string(v.side) + " " + string(v.align)
			opts = append(opts,
				dm.Label(string(v.side)+" / "+string(v.align)),
				dm.Item("First"),
				dm.Item("Second"),
			)
			return dm.DropdownMenu(dm.Modal(!v.nonModal),
				dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeSm, label)),
				dm.Content(opts...),
			)
		}),
	)
}

// Rtl is the parity fixture dropdown-menu-rtl.
func Rtl() h.Node {
	return dm.DropdownMenu(dm.Dir("rtl"),
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "RTL")),
		dm.Content(h.Class("w-48"),
			dm.Item("First"),
			dm.Sub(
				dm.SubTrigger("More"),
				dm.SubContent(
					dm.Item("Inner one"),
					dm.Item("Inner two"),
				),
			),
			dm.Separator(),
			dm.RadioGroup(dm.Value("one"),
				dm.RadioItem(dm.Value("one"), "One"),
				dm.RadioItem(dm.Value("two"), "Two"),
			),
		),
	)
}

// DefaultOpen is the parity fixture dropdown-menu-default-open.
func DefaultOpen() h.Node {
	return dm.DropdownMenu(dm.DefaultOpen,
		dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open on load")),
		dm.Content(h.Class("w-48"),
			dm.Label("Account"),
			dm.Separator(),
			dm.Item("Profile"),
			dm.Item("Settings"),
		),
	)
}
