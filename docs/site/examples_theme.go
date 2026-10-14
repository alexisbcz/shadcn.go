package site

import (
	_ "embed"
	"encoding/json"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/drawer"
	"shadcn.go/components/ui/label"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/tabs"
	h "shadcn.go/libhtml"
)

// The examples layout's theme selector: components/theme-selector.tsx,
// CopyCodeButton and CustomizerCode of components/theme-customizer.tsx, and
// the active theme of components/active-theme.tsx (examples-theme.js). The
// themes are lib/themes.ts over registry/_legacy-base-colors.ts, extracted
// by gen/themes.mjs.

//go:embed data/examples/theme/themes.json
var themesJSON []byte

type themeVars struct {
	Light [][2]string `json:"light"`
	Dark  [][2]string `json:"dark"`
}

type exampleTheme struct {
	Name  string     `json:"name"`
	Label string     `json:"label"`
	HSL   themeVars  `json:"hsl"`
	OKLCH *themeVars `json:"oklch"`
}

var exampleThemes = func() []exampleTheme {
	var t []exampleTheme
	if err := json.Unmarshal(themesJSON, &t); err != nil {
		panic(err)
	}
	return t
}()

// defaultExampleTheme is the selector's value for the "default" active
// theme.
const defaultExampleTheme = "neutral"

// themeSelector is components/theme-selector.tsx.
func themeSelector(class string) h.Node {
	var items []any
	for _, t := range exampleThemes {
		items = append(items, selectui.Item(selectui.ItemValue(t.Name), h.Class("data-[state=checked]:opacity-50"), t.Label))
	}
	// The selector is styles/base-nova; the radix select of the nova
	// style renders the closest.
	return h.Div(h.Class(lib.Cn("flex items-center gap-2", class)), h.Data("theme-selector", ""),
		style.Scope(style.Nova,
			label.Label(h.For("theme-selector"), h.Class("sr-only"), "Theme"),
			selectui.Select(selectui.DefaultValue(defaultExampleTheme),
				selectui.Trigger(h.ID("theme-selector"), h.Class("w-36"), selectui.Value(selectui.Placeholder("Select a theme"))),
				selectui.Content(selectui.Group(selectui.Label("Theme"), h.Group(items...))),
			),
		),
		copyCodeButton(button.VariantSecondary, button.SizeIconSm, "icon-sm", "rounded-lg border bg-transparent"),
		h.Script(h.Type("application/json"), h.Data("themes", ""), h.UnsafeJS(strings.ReplaceAll(string(themesJSON), "<", "\\u003c"))),
		h.Script(h.Type("module"), h.Src(staticPrefix+"examples-theme.js")),
	)
}

// copyCodeButton is CopyCodeButton: a drawer below sm, a dialog above.
func copyCodeButton(variant button.Variant, size button.Size, sizeName, class string) h.Node {
	name := defaultExampleTheme
	title := func() h.Node { return h.Span(h.Data("theme-name", ""), name) }
	return h.Group(
		drawer.Drawer(
			drawer.Trigger(lib.AsChild, button.Button(variant, size, h.Class(lib.Cn("sm:hidden", class)), "Copy Code")),
			drawer.Content(h.Class("h-auto rounded-t-2xl"),
				drawer.Header(
					drawer.Title(h.Class("capitalize"), title()),
					drawer.Description("Copy and paste the following code into your CSS file."),
				),
				customizerCode(name),
			),
		),
		dialog.Dialog(
			dialog.Trigger(lib.AsChild, button.Button(variant, size, h.Data("size", sizeName), h.Class(lib.Cn("group/button hidden sm:flex", class)),
				iconCopy(),
				h.Span(h.Class("group-data-[size=icon-sm]/button:sr-only"), "Copy Code"),
			)),
			dialog.Content(h.Class("rounded-2xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-2xl dark:bg-neutral-800 dark:ring-neutral-900"),
				dialog.Header(
					dialog.Title(h.Class("capitalize"), title()),
					dialog.Description("Copy and paste the following code into your CSS file."),
				),
				customizerCode(name),
			),
		),
	)
}

func findExampleTheme(name string) exampleTheme {
	for _, t := range exampleThemes {
		if t.Name == name {
			return t
		}
	}
	return exampleThemes[0]
}

// customizerCode is CustomizerCode; examples-theme.js re-renders the code
// lines (same markup) when the active theme changes.
func customizerCode(name string) h.Node {
	t := findExampleTheme(name)
	figure := func(format string, lines ...any) h.Node {
		return h.Figure(h.Data("rehype-pretty-code-figure", ""), h.Class("mx-0! mt-0 rounded-xl"),
			h.Figcaption(h.Class("flex items-center gap-2 text-code-foreground [&_svg]:size-4 [&_svg]:text-code-foreground [&_svg]:opacity-70"),
				h.Data("rehype-pretty-code-title", ""), h.Data("language", "css"), h.Data("theme", "github-dark github-light-default"),
				h.Svg(h.Xmlns("http://www.w3.org/2000/svg"), h.ViewBox("0 0 24 24"), h.Class("fill-foreground"), h.Path(h.D(cssIconPath))),
				"app/globals.css",
			),
			h.Pre(h.Class("no-scrollbar max-h-[300px] min-w-0 overflow-x-auto px-4 py-3.5 outline-none has-[[data-highlighted-line]]:px-0 has-[[data-line-numbers]]:px-0 has-[[data-slot=tabs]]:p-0 md:max-h-[450px]"),
				button.Button(h.Data("slot", "copy-button"), button.SizeIcon, button.VariantGhost, h.Data("theme-copy", format),
					h.Class("absolute top-3 right-2 z-10 size-7 bg-code text-code-foreground shadow-none hover:opacity-100 focus-visible:opacity-100"),
					h.Span(h.Class("sr-only"), "Copy"), iconCopy(),
				),
				h.Code(h.Data("line-numbers", ""), h.Data("language", "css"), h.Data("theme-code", format), h.Group(lines...)),
			),
		)
	}
	return tabs.Tabs(tabs.Value("v4-oklch"), h.Class("min-w-0 px-4 pb-4 md:p-0"),
		tabs.List(
			tabs.Trigger(tabs.Value("v4-oklch"), "OKLCH"),
			tabs.Trigger(tabs.Value("v4-hsl"), "HSL"),
			tabs.Trigger(tabs.Value("v3"), "Tailwind v3"),
		),
		tabs.Content(tabs.Value("v4-oklch"), figure("v4-oklch", themeCodeLines(t, "v4-oklch")...)),
		tabs.Content(tabs.Value("v4-hsl"), figure("v4-hsl", themeCodeLines(t, "v4-hsl")...)),
		tabs.Content(tabs.Value("v3"), figure("v3", themeCodeLines(t, "v3")...)),
	)
}

const nbsp = " "

// themeCodeLines renders the code lines of a tab (mirrored by lines() in
// examples-theme.js).
func themeCodeLines(t exampleTheme, format string) []any {
	var out []any
	line := func(class string, children ...any) {
		out = append(out, h.Span(h.Data("line", ""), h.Class(class), h.Group(children...)))
	}
	ind := func(color string) h.Node {
		return h.Span(h.Class("inline-block size-3 border border-border/50"), h.Style("background-color:"+color))
	}
	const fg = "line text-code-foreground"
	switch format {
	case "v4-oklch", "v4-hsl":
		vars := t.HSL
		if format == "v4-oklch" && t.OKLCH != nil {
			vars = *t.OKLCH
		}
		entry := func(kv [2]string) {
			if format == "v4-oklch" {
				line(fg, strings.Repeat(nbsp, 3)+"--"+kv[0]+": ", ind(kv[1]), " "+kv[1]+";")
			} else {
				line(fg, strings.Repeat(nbsp, 3)+"--"+kv[0]+": ", ind("hsl("+kv[1]+")"), " hsl("+kv[1]+");")
			}
		}
		line(fg, nbsp+":root {")
		line(fg, strings.Repeat(nbsp, 3)+"--radius: 0.65rem;")
		for _, kv := range vars.Light {
			entry(kv)
		}
		line(fg, nbsp+"}")
		line(fg, nbsp)
		line(fg, nbsp+".dark {")
		for _, kv := range vars.Dark {
			entry(kv)
		}
		line(fg, nbsp+"}")
	case "v3":
		get := func(vars [][2]string, k string) string {
			for _, kv := range vars {
				if kv[0] == k {
					return kv[1]
				}
			}
			return ""
		}
		v := func(vars [][2]string, k string) {
			val := get(vars, k)
			line("line", strings.Repeat(nbsp, 4)+"--"+k+": ", ind("hsl("+val+")"), " "+val+";")
		}
		block := func(sel string, vars [][2]string, radius bool) {
			line("line", nbsp+nbsp+sel+" {")
			v(vars, "background")
			v(vars, "foreground")
			for _, p := range []string{"card", "popover", "primary", "secondary", "muted", "accent", "destructive"} {
				v(vars, p)
				v(vars, p+"-foreground")
			}
			v(vars, "border")
			v(vars, "input")
			v(vars, "ring")
			if radius {
				line("line", strings.Repeat(nbsp, 4)+"--radius: 0.5rem;")
			}
			for _, p := range []string{"chart-1", "chart-2", "chart-3", "chart-4", "chart-5"} {
				v(vars, p)
			}
			line("line", nbsp+nbsp+"}")
		}
		line("line", "@layer base {")
		block(":root", t.HSL.Light, true)
		line("line", nbsp)
		block(".dark", t.HSL.Dark, false)
		line("line", "}")
	}
	return out
}

// cssIconPath is Icons.css (components/icons.tsx).
const cssIconPath = "M0 0v20.16A3.84 3.84 0 0 0 3.84 24h16.32A3.84 3.84 0 0 0 24 20.16V3.84A3.84 3.84 0 0 0 20.16 0Zm14.256 13.08c1.56 0 2.28 1.08 2.304 2.64h-1.608c.024-.288-.048-.6-.144-.84-.096-.192-.288-.264-.552-.264-.456 0-.696.264-.696.84-.024.576.288.888.768 1.08.72.288 1.608.744 1.92 1.296q.432.648.432 1.656c0 1.608-.912 2.592-2.496 2.592-1.656 0-2.4-1.032-2.424-2.688h1.68c0 .792.264 1.176.792 1.176.264 0 .456-.072.552-.24.192-.312.24-1.176-.048-1.512-.312-.408-.912-.6-1.32-.816q-.828-.396-1.224-.936c-.24-.36-.36-.888-.36-1.536 0-1.44.936-2.472 2.424-2.448m5.4 0c1.584 0 2.304 1.08 2.328 2.64h-1.608c0-.288-.048-.6-.168-.84-.096-.192-.264-.264-.528-.264-.48 0-.72.264-.72.84s.288.888.792 1.08c.696.288 1.608.744 1.92 1.296.264.432.408.984.408 1.656.024 1.608-.888 2.592-2.472 2.592-1.68 0-2.424-1.056-2.448-2.688h1.68c0 .744.264 1.176.792 1.176.264 0 .456-.072.552-.24.216-.312.264-1.176-.048-1.512-.288-.408-.888-.6-1.32-.816-.552-.264-.96-.576-1.2-.936s-.36-.888-.36-1.536c-.024-1.44.912-2.472 2.4-2.448m-11.031.018c.711-.006 1.419.198 1.839.63.432.432.672 1.128.648 1.992H9.336c.024-.456-.096-.792-.432-.96-.312-.144-.768-.048-.888.24-.12.264-.192.576-.168.864v3.504c0 .744.264 1.128.768 1.128a.65.65 0 0 0 .552-.264c.168-.24.192-.552.168-.84h1.776c.096 1.632-.984 2.712-2.568 2.688-1.536 0-2.496-.864-2.472-2.472v-4.032c0-.816.24-1.44.696-1.848.432-.408 1.146-.624 1.857-.63"
