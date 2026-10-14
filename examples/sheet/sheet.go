// Package sheet ports the upstream sheet-* examples.
package sheet

import (
	"net/http"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/sheet"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "sheet", Render: withThemeScript(fn)}
	}
	registry.Register(
		ex("sheet-demo", Demo),
		ex("sheet-side", SideExample),
	)
}

// withThemeScript puts an empty nonced <script> first in the preview, where
// the upstream docs' next-themes ThemeProvider renders its theme script.
// It matters to modal dialogs: hideOthers (aria-hidden) keeps <script>
// elements and their ancestors visible, so the preview itself is kept and
// its other children are hidden, exactly as upstream.
func withThemeScript(fn func() h.Node) func(*http.Request) h.Node {
	return func(r *http.Request) h.Node { return h.Group(secure.Script(r), fn()) }
}

func Demo() h.Node {
	return sheet.Sheet(
		sheet.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
		sheet.Content(
			sheet.Header(
				sheet.Title("Edit profile"),
				sheet.Description("Make changes to your profile here. Click save when you're done."),
			),
			h.Div(h.Class("grid flex-1 auto-rows-min gap-6 px-4"),
				h.Div(h.Class("grid gap-3"),
					label.Label(h.For("sheet-demo-name"), "Name"),
					input.Input(h.ID("sheet-demo-name"), h.Value("Pedro Duarte")),
				),
				h.Div(h.Class("grid gap-3"),
					label.Label(h.For("sheet-demo-username"), "Username"),
					input.Input(h.ID("sheet-demo-username"), h.Value("@peduarte")),
				),
			),
			sheet.Footer(
				button.Button(h.Type("submit"), "Save changes"),
				sheet.Close(lib.AsChild, button.Button(button.VariantOutline, "Close")),
			),
		),
	)
}

func SideExample() h.Node {
	sides := []sheet.Side{sheet.SideTop, sheet.SideRight, sheet.SideBottom, sheet.SideLeft}
	return h.Div(h.Class("grid grid-cols-2 gap-2"),
		h.Map(sides, func(side sheet.Side) h.Node {
			return sheet.Sheet(
				sheet.Trigger(lib.AsChild, button.Button(button.VariantOutline, string(side))),
				sheet.Content(side,
					sheet.Header(
						sheet.Title("Edit profile"),
						sheet.Description("Make changes to your profile here. Click save when you're done."),
					),
					h.Div(h.Class("grid gap-4 py-4"),
						h.Div(h.Class("grid grid-cols-4 items-center gap-4"),
							label.Label(h.For("name"), h.Class("text-right"), "Name"),
							input.Input(h.ID("name"), h.Value("Pedro Duarte"), h.Class("col-span-3")),
						),
						h.Div(h.Class("grid grid-cols-4 items-center gap-4"),
							label.Label(h.For("username"), h.Class("text-right"), "Username"),
							input.Input(h.ID("username"), h.Value("@peduarte"), h.Class("col-span-3")),
						),
					),
					sheet.Footer(
						sheet.Close(lib.AsChild, button.Button(h.Type("submit"), "Save changes")),
					),
				),
			)
		}),
	)
}
