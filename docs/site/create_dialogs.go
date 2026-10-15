package site

import (
	_ "embed"

	"shadcn.go/components/lib"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/command"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/drawer"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	radiogroup "shadcn.go/components/ui/radio-group"
	"shadcn.go/components/ui/runtime"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/tabs"
	togglegroup "shadcn.go/components/ui/toggle-group"
	h "shadcn.go/libhtml"
)

// The designer's dialogs: Reset (reset-button.tsx), Open Preset
// (open-preset.tsx), Get Code (project-form.tsx) and the action menu
// (action-menu.tsx). create-dialogs.js opens them and keeps their state.

//go:embed static/create-dialogs.js
var createDialogsScript string

func init() { runtime.Module("create-dialogs", createDialogsScript) }

func createDialogs(ds designSystem) h.Node {
	return h.Group(resetDialog(), openPresetDialogs(), projectFormDialog(ds), actionMenuDialog(ds))
}

func resetDialog() h.Node {
	return h.Div(h.Data("create-dialog", "reset"),
		alertdialog.AlertDialog(
			alertdialog.Content(alertdialog.SizeSm,
				alertdialog.Header(
					alertdialog.Title("Reset to defaults?"),
					alertdialog.Description("This will reset all customization options to their default values."),
				),
				alertdialog.Footer(
					alertdialog.Cancel("Cancel"),
					alertdialog.Action(h.Data("create-reset-confirm", ""), "Reset"),
				),
			),
		),
	)
}

const presetExample = "b2D0wqNxT"

func presetField() h.Node {
	return field.Field(h.Data("create-preset-field", ""),
		field.Label(h.For("preset-code"), h.Class("sr-only"), "Preset code"),
		field.Content(
			input.Input(h.ID("preset-code"), h.Placeholder(presetExample+" or --preset "+presetExample),
				h.Attr("autocapitalize", "none"), h.Attr("autocorrect", "off"), h.Attr("spellcheck", "false"),
				h.Class("h-10 md:h-8"), h.Data("create-preset-input", "")),
		),
	)
}

// openPresetDialogs are OpenPreset's dialog (desktop) and drawer (mobile).
func openPresetDialogs() h.Node {
	return h.Group(
		h.Div(h.Data("create-dialog", "open-preset"),
			dialog.Dialog(
				dialog.Content(h.Class("dark"),
					h.Form(h.Data("create-preset-form", ""),
						dialog.Header(
							dialog.Title("Open Preset"),
							dialog.Description("Paste a preset code to load a saved configuration."),
						),
						h.Div(h.Class("py-4"), presetField()),
						dialog.Footer(
							dialog.Close(lib.AsChild, button.Button(button.VariantOutline, h.Type("button"), "Cancel")),
							button.Button(h.Type("submit"), h.Bool("disabled"), "Open"),
						),
					),
				),
			),
		),
		h.Div(h.Data("create-dialog", "open-preset-mobile"),
			drawer.Drawer(
				drawer.Content(h.Class("dark rounded-t-2xl!"),
					drawer.Header(
						drawer.Title(h.Class("text-xl"), "Open Preset"),
						drawer.Description("Paste a preset code to load a saved configuration."),
					),
					h.Form(h.Data("create-preset-form", ""),
						h.Div(h.Class("px-4 py-2"), presetField()),
						drawer.Footer(
							button.Button(h.Type("submit"), h.Class("h-10"), h.Bool("disabled"), "Open"),
							drawer.Close(lib.AsChild, button.Button(button.VariantOutline, h.Type("button"), h.Class("h-10"), "Cancel")),
						),
					),
				),
			),
		),
	)
}

// --- Get Code (project-form.tsx) ---

var packageManagers = []string{"pnpm", "npm", "yarn", "bun"}

func commandTabs(kind string) h.Node {
	var triggers, contents []any
	for _, pm := range packageManagers {
		triggers = append(triggers, tabs.Trigger(tabs.Value(pm), h.Class("py-0 leading-none data-[state=active]:shadow-none"), pm))
		contents = append(contents, tabs.Content(tabs.Value(pm),
			h.Div(h.Class("relative overflow-hidden border-t bg-popover p-3"),
				h.Div(h.Class("no-scrollbar overflow-x-auto"),
					h.Code(h.Class("font-mono text-sm whitespace-nowrap"), h.Data("create-command", kind+":"+pm)),
				),
			),
		))
	}
	return tabs.Tabs(tabs.Value("pnpm"), h.Class("min-w-0 gap-0 overflow-hidden rounded-xl border-0 ring-1 ring-border"), h.Data("create-package-manager", ""),
		h.Div(h.Class("flex items-center gap-2 py-1 pr-1.5 pl-1"),
			tabs.List(h.Class("bg-transparent font-mono"), h.Group(triggers...)),
			button.Button(button.SizeIconSm, button.VariantGhost, h.Class("ml-auto"), h.Data("create-copy", kind),
				crHugeiconMarkup("Copy01Icon", ""), h.Span(h.Class("sr-only"), "Copy command")),
		),
		h.Group(contents...),
	)
}

func gridOption(id, logo, title string, contentClass, titleClass, logoClass string) h.Node {
	return field.Label(h.For(id), h.Class("block w-full"),
		field.Field(field.Horizontal, h.Class("w-full rounded-md transition-colors duration-150 hover:bg-neutral-700/45"),
			field.Content(h.Class("flex flex-row items-center gap-2 "+contentClass),
				h.Div(h.Class(logoClass), h.UnsafeHTML(logo)),
				field.Title(h.Class(titleClass), title),
			),
			radiogroup.Item(radiogroup.Value(id[len("template-"):]), h.ID(id), h.Class("sr-only absolute")),
		),
	)
}

func projectFormDialog(ds designSystem) h.Node {
	d := createDataOnce()
	var templates []any
	for _, t := range d.Templates {
		templates = append(templates, gridOption("template-"+t.Value, t.Logo, t.Title, "px-2.5 py-1.5", "", "size-4 text-neutral-100 [&_svg]:size-4 *:[svg]:text-neutral-100!"))
	}
	var bases []any
	for _, b := range d.Bases {
		bases = append(bases, field.Label(h.For("base-"+b.Name), h.Class("block w-full"),
			field.Field(field.Horizontal, h.Class("w-full rounded-md transition-colors duration-150 hover:bg-neutral-700/45"),
				field.Content(h.Class("flex flex-row items-center gap-2 py-1"),
					h.Div(h.Class("size-4 shrink-0 text-neutral-100 [&_svg]:size-4 *:[svg]:text-neutral-100!"), h.UnsafeHTML(b.Logo)),
					field.Title(h.Class("whitespace-nowrap"), b.Title),
				),
				radiogroup.Item(radiogroup.Value(b.Name), h.ID("base-"+b.Name), h.Class("sr-only absolute")),
			),
		))
	}
	option := func(id, label string, icon h.Node) h.Node {
		return field.Field(field.Horizontal,
			field.Label(h.For(id), icon, label),
			switchui.Switch(h.ID(id), h.Data("create-option", id)),
		)
	}
	applyModes := []struct{ value, title, description string }{
		{"full", "Full preset", "Everything from the preset, including components, theme, and fonts."},
		{"theme", "Theme only", "Theme tokens only, like colors, radii, and shadows. Components stay as they are."},
		{"font", "Fonts only", "Only preset fonts for body and headings. Components stay as they are."},
	}
	var modes []any
	for _, m := range applyModes {
		modes = append(modes, field.Label(h.For("apply-"+m.value),
			field.Field(field.Horizontal,
				radiogroup.Item(radiogroup.Value(m.value), h.ID("apply-"+m.value)),
				field.Content(field.Title(m.title), field.Description(m.description)),
			),
		))
	}
	copyButton := func(kind, label string) h.Node {
		return button.Button(h.Class("h-9 w-full"), h.Data("create-copy", kind), h.Data("create-copy-label", label), label)
	}
	return h.Div(h.Data("create-dialog", "project-form"),
		dialog.Dialog(
			dialog.Content(h.Class("dark top-[64px] no-scrollbar flex max-h-[calc(100svh-2rem)] translate-y-0 flex-col rounded-2xl p-0 shadow-xl **:data-[slot=dialog-close]:top-4.5 **:data-[slot=dialog-close]:right-4 **:data-[slot=field-separator]:h-2 sm:max-w-md"),
				h.Div(h.Class("flex min-w-0 flex-1 flex-col gap-0 overflow-hidden rounded-2xl"),
					dialog.Header(h.Class("border-b px-6 py-5"),
						togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("new-project"), togglegroup.Spacing(2), h.Aria("label", "Project type"),
							h.Class("**:data-[slot=toggle-group-item]:data-pressed:bg-neutral-700/70 **:data-[slot=toggle-group-item]:data-[state=on]:bg-neutral-700/70"),
							h.Data("create-project-tabs", ""),
							togglegroup.Item(togglegroup.Value("new-project"), "New Project"),
							togglegroup.Item(togglegroup.Value("existing-project"), "Existing Project"),
							togglegroup.Item(togglegroup.Value("theme"), "Theme"),
						),
					),
					h.Div(h.Class("no-scrollbar overflow-y-auto"), h.Data("create-project-tab", "new-project"),
						field.Group(h.Class("px-6 py-4"),
							field.Field(h.Class("gap-3"),
								field.Label("Template"),
								radiogroup.RadioGroup(radiogroup.Value("next"), h.Class("grid grid-cols-2 gap-2"), h.Data("create-template", ""), h.Group(templates...)),
							),
							field.Separator(h.Class("-mx-6")),
							field.Field(
								field.Label("Base"),
								radiogroup.RadioGroup(radiogroup.Value(ds.Base), h.Aria("label", "Base"), h.Class("grid grid-cols-1 gap-2 sm:grid-cols-3"), h.Data("create-base", ""), h.Group(bases...)),
							),
							field.Separator(h.Class("-mx-6")),
							field.Set(
								field.Legend(field.LegendLabel, h.Class("sr-only"), "Options"),
								option("pointer", "Use pointer on buttons", crHugeiconMarkup("HandPointingRight04Icon", "size-4 -rotate-90")),
								field.Separator(h.Class("-mx-6")),
								option("monorepo", "Create a monorepo", h.Span(h.Class("size-4 text-neutral-100 [&_svg]:size-4 [&_svg]:fill-current"), h.UnsafeHTML(d.TurborepoLogo))),
								field.Separator(h.Class("-mx-6")),
								option("rtl", "Enable RTL support", crHugeiconMarkup("Globe02Icon", "size-4")),
							),
						),
						dialog.Footer(h.Class("m-0 min-w-0 p-6"),
							h.Div(h.Class("flex w-full min-w-0 flex-col gap-3"), commandTabs("command"), copyButton("command", "Copy Command")),
						),
					),
					h.Div(h.Class("no-scrollbar overflow-y-auto"), h.Data("create-project-tab", "existing-project"), h.Bool("hidden"),
						field.Group(h.Class("px-6 py-4"),
							field.Set(h.Class("gap-3"),
								field.Legend(field.LegendLabel, "Apply Preset"),
								field.Description("Pick which parts of the preset to apply."),
								radiogroup.RadioGroup(radiogroup.Value("full"), h.Aria("label", "Apply"), h.Data("create-apply-mode", ""), h.Group(modes...)),
							),
						),
						dialog.Footer(h.Class("m-0 min-w-0 p-6"),
							h.Div(h.Class("flex w-full min-w-0 flex-col gap-3"), commandTabs("apply"), copyButton("apply", "Copy Command")),
						),
					),
					h.Div(h.Class("no-scrollbar overflow-y-auto"), h.Data("create-project-tab", "theme"), h.Bool("hidden"),
						field.Group(h.Class("min-w-0 px-6 py-4"),
							field.Set(h.Class("min-w-0 gap-3"),
								field.Legend(field.LegendLabel, "Theme Tokens"),
								field.Description("Copy the CSS variables for this preset."),
								h.Div(h.Class("w-full min-w-0 overflow-hidden rounded-xl border-0 ring-1 ring-border"),
									h.Div(h.Class("flex items-center gap-2 py-1 pr-1.5 pl-3"),
										h.Div(h.Class("min-w-0 truncate font-mono text-sm text-muted-foreground"), "globals.css"),
										button.Button(button.SizeIconSm, button.VariantGhost, h.Class("ml-auto"), h.Data("create-copy", "theme"),
											crHugeiconMarkup("Copy01Icon", ""), h.Span(h.Class("sr-only"), "Copy theme")),
									),
									h.Div(h.Class("relative no-scrollbar max-h-[45svh] overflow-auto border-t bg-popover p-3"),
										h.Pre(h.Class("min-w-max font-mono leading-normal whitespace-pre"), h.Code(h.Data("create-theme-css", ""))),
									),
								),
							),
						),
						dialog.Footer(h.Class("m-0 min-w-0 p-6"), copyButton("theme", "Copy Theme")),
					),
				),
			),
		),
		h.Template(h.Data("create-copy-icon", "copy"), crHugeiconMarkup("Copy01Icon", "")),
		h.Template(h.Data("create-copy-icon", "tick"), crHugeiconMarkup("Tick02Icon", "")),
	)
}

// --- action menu (action-menu.tsx) ---

// actionMenuDialog is the Navigate... command palette. Of the registry
// items of the base, the Go site renders the preview block only (the
// component examples are not ported).
func actionMenuDialog(ds designSystem) h.Node {
	item := func(name, label string) h.Node {
		return command.Item(command.Value(label), h.Data("checked", boolString(ds.Item == name)), h.Class("px-2"), h.Data("create-navigate", name), label)
	}
	return h.Div(h.Data("create-dialog", "action-menu"),
		command.Dialog(h.Class("animate-none!"),
			command.Command(command.Loop(true),
				command.Input(h.Placeholder("Search")),
				command.List(
					command.Empty("No items found."),
					// getAllItems drops names ending in a digit (preview-02).
					command.Group(item("preview", "Preview")),
				),
			),
		),
	)
}
