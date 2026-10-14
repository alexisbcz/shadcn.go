package site

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/drawer"
	selectui "shadcn.go/components/ui/select"
	h "shadcn.go/libhtml"
)

// The typeset docs (components/docs-panel.tsx, get-code-drawer.tsx): the
// install steps and the agent prompt, generated from the params and the
// framework. The body of the docs is also served alone at
// /typeset/docs-content, which typeset.js fetches when the params or the
// framework change, so the code is generated in one place.

// typesetDocsPanel is TypesetDocsPanel.
func typesetDocsPanel(r *http.Request, p typesetParams) h.Node {
	return h.Div(h.Class("hidden w-88 flex-col items-start gap-3 xl:flex 2xl:w-104"),
		h.Div(h.Class("isolate flex max-h-[calc(100vh-16rem)] w-full flex-col overflow-hidden rounded-2xl bg-background ring-1 ring-foreground/10"),
			typesetDocsContent(r, p),
		),
		h.A(h.Data("slot", "button"), h.Href("/docs/typeset"), h.Role("button"), h.TabIndex(0),
			h.Class("group/button inline-flex shrink-0 items-center justify-center rounded-[min(var(--radius-md),12px)] border border-transparent bg-clip-padding text-[0.8rem] font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-3.5 underline-offset-4 hover:underline h-7 gap-1 px-2.5 in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 text-muted-foreground"),
			"Read the docs"),
	)
}

// typesetGetCodeDrawer is TypesetGetCodeDrawer: a drawer from the right,
// or from the bottom on mobile (useIsMobile), so two drawers whose
// triggers show at their breakpoints.
func typesetGetCodeDrawer(r *http.Request, p typesetParams) h.Node {
	one := func(dir drawer.Direction, mobile bool, visibility string) h.Node {
		return drawer.Drawer(dir,
			drawer.Trigger(lib.AsChild,
				h.Button(h.Type("button"), h.Data("slot", "button"), h.Class(typesetFooterButtonClass+" xl:hidden "+visibility), "Get Code"),
			),
			drawer.Content(h.Data("mobile", boolString(mobile)), h.Class("data-[vaul-drawer-direction=bottom]:max-h-[85svh] [&>.bg-muted.rounded-full]:hidden"),
				drawer.Header(h.Class("pb-0"),
					drawer.Title("Get Code"),
					drawer.Description("Install typeset with the values you picked."),
				),
				typesetDocsContent(r, p),
				drawer.Footer(
					drawer.Close(lib.AsChild,
						h.Button(h.Type("button"), h.Data("slot", "button"), h.Class(typesetOutlineButtonClass), "Done"),
					),
				),
			),
		)
	}
	return h.Group(
		one(drawer.DirectionBottom, true, "md:hidden"),
		one(drawer.DirectionRight, false, "max-md:hidden"),
	)
}

const typesetOutlineButtonClass = "group/button inline-flex shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 h-8 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2"

const typesetSmallOutlineButtonClass = "group/button inline-flex shrink-0 items-center justify-center rounded-[min(var(--radius-md),12px)] border border-border bg-background bg-clip-padding text-[0.8rem] font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-3.5 hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 h-7 gap-1 px-2.5 in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 w-fit"

const typesetTabTriggerClass = "relative inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-1.5 py-0.5 text-sm font-medium whitespace-nowrap text-foreground/60 transition-all group-data-vertical/tabs:w-full group-data-vertical/tabs:justify-start hover:text-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50 has-data-[icon=inline-end]:pr-1 has-data-[icon=inline-start]:pl-1 aria-disabled:pointer-events-none aria-disabled:opacity-50 dark:text-muted-foreground dark:hover:text-foreground group-data-[variant=default]/tabs-list:data-active:shadow-sm group-data-[variant=line]/tabs-list:data-active:shadow-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 group-data-[variant=line]/tabs-list:bg-transparent group-data-[variant=line]/tabs-list:data-active:bg-transparent dark:group-data-[variant=line]/tabs-list:data-active:border-transparent dark:group-data-[variant=line]/tabs-list:data-active:bg-transparent data-active:bg-background data-active:text-foreground dark:data-active:border-input dark:data-active:bg-input/30 dark:data-active:text-foreground after:absolute after:bg-foreground after:opacity-0 after:transition-opacity group-data-horizontal/tabs:after:inset-x-0 group-data-horizontal/tabs:after:bottom-[-5px] group-data-horizontal/tabs:after:h-0.5 group-data-vertical/tabs:after:inset-y-0 group-data-vertical/tabs:after:-right-1 group-data-vertical/tabs:after:w-0.5 group-data-[variant=line]/tabs-list:data-active:after:opacity-100"

// typesetDocsContent is TypesetDocsContent: Base UI tabs (rendered here as
// plain markup; typeset.js switches them) and the framework select.
func typesetDocsContent(r *http.Request, p typesetParams) h.Node {
	tab := func(value, label string, active bool) h.Node {
		return h.Button(h.Type("button"), h.Role("tab"), h.Data("slot", "tabs-trigger"), h.Data("orientation", "horizontal"),
			h.Data("typeset-tab", value), h.If(active, h.Data("active", "")), h.Aria("selected", boolString(active)),
			h.TabIndex(map[bool]int{true: 0, false: -1}[active]), h.Class(typesetTabTriggerClass), label)
	}
	var frameworks []any
	for _, t := range typesetData.Templates {
		frameworks = append(frameworks, selectui.Item(selectui.ItemValue(t.Value), t.Title))
	}
	return h.Div(h.Data("slot", "tabs"), h.Data("orientation", "horizontal"), h.Data("typeset-tabs", ""),
		h.Class("group/tabs flex data-horizontal:flex-col min-h-0 flex-1 flex-col gap-0"),
		h.Div(h.Class("flex items-center justify-between gap-2 border-b px-4 py-3"),
			h.Div(h.Role("tablist"), h.Data("slot", "tabs-list"), h.Data("variant", "default"), h.Data("orientation", "horizontal"),
				h.Class("group/tabs-list inline-flex w-fit items-center justify-center rounded-lg p-[3px] text-muted-foreground group-data-horizontal/tabs:h-8 group-data-vertical/tabs:h-fit group-data-vertical/tabs:flex-col data-[variant=line]:rounded-none bg-muted"),
				tab("docs", "Docs", true),
				tab("prompt", "Prompt", false),
			),
			style.Scope(style.Nova,
				h.Div(h.Class("contents"), h.Data("typeset-framework", ""),
					selectui.Select(selectui.DefaultValue("next"),
						selectui.Trigger(selectui.SizeSm, selectui.Value()),
						selectui.Content(selectui.PositionPopper, selectui.AlignEnd, selectui.Group(frameworks...)),
					),
				),
			),
		),
		h.Div(h.Data("typeset-docs-body", ""), h.Class("min-h-0 flex-1 scroll-fade scrollbar-none overflow-y-auto p-4 md:p-6"),
			typesetDocsBody(r, p, "next"),
		),
	)
}

// TypesetDocsContent serves /typeset/docs-content: the docs body for the
// query's params and framework.
func (s *Site) TypesetDocsContent(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	framework := q.Get("framework")
	ok := false
	for _, t := range typesetData.Templates {
		ok = ok || t.Value == framework
	}
	if !ok {
		framework = "next"
	}
	_ = h.WriteHTTP(w, http.StatusOK, typesetDocsBody(r, parseTypesetParams(q), framework))
}

func typesetDocsBody(r *http.Request, p typesetParams, framework string) h.Node {
	c := newTypesetCode(r, p, framework)
	section := func(step int, title string, children ...any) h.Node {
		return h.Section(h.Class("flex flex-col gap-2.5"),
			h.H3(h.Class("text-sm font-medium"), strconv.Itoa(step), ". ", title),
			h.Group(children...))
	}
	prose := func(children ...any) h.Node {
		return h.P(h.Class("text-sm leading-relaxed text-muted-foreground"), h.Group(children...))
	}
	code := func(text string, class ...string) h.Node {
		return h.Div(h.Class("relative rounded-lg bg-muted/60"),
			h.Pre(h.Class(lib.Cn("scrollbar-none overflow-x-auto p-3 font-mono text-sm leading-relaxed", strings.Join(class, " "))), text))
	}
	copyButton := func(action, label string) h.Node {
		return h.Button(h.Type("button"), h.Data("slot", "button"), h.Data("typeset-copy", action), h.Class(typesetSmallOutlineButtonClass),
			lucide.Copy(h.Data("icon", "inline-start")), label)
	}
	var fonts []any
	if framework == "next" {
		fonts = []any{
			prose("Load the fonts with ", h.Code(h.Class("font-mono"), "next/font"), " in your root layout:"),
			code(c.nextFontCode()),
		}
	} else {
		fonts = []any{
			typesetDocsCommand(c.fontsourceCommand()),
			prose("Then import them in your CSS file:"),
			code(c.fontsourceCSS()),
		}
	}
	return h.Group(
		h.Div(h.Role("tabpanel"), h.Data("slot", "tabs-content"), h.Data("typeset-panel", "docs"), h.TabIndex(0),
			h.Class("flex-1 text-sm outline-none flex flex-col gap-6"),
			section(1, "Create typeset.css",
				prose("Copy the stylesheet into a ", h.Code(h.Class("font-mono"), "typeset.css"), " file next to your main CSS file, then import it:"),
				copyButton("css", "Copy typeset.css"),
				code("@import \"tailwindcss\";\n@import \"./typeset.css\";"),
			),
			section(2, "Add the fonts", fonts...),
			section(3, "Create your custom typeset", code(c.preset())),
			section(4, "Wrap your content", code(c.usage())),
		),
		h.Div(h.Role("tabpanel"), h.Data("slot", "tabs-content"), h.Data("typeset-panel", "prompt"), h.TabIndex(0), h.Bool("hidden"),
			h.Class("flex-1 text-sm outline-none flex flex-col gap-2.5"),
			prose("One prompt with your picks baked in. Copy it and paste it into your coding agent."),
			code(c.prompt(), "max-h-102 scroll-fade overflow-y-auto whitespace-pre-wrap"),
			copyButton("prompt", "Copy Prompt"),
		),
	)
}

// typesetDocsCommand is DocsCommand: the install command per package
// manager (typeset.js switches them and stores the choice in the docs'
// "config").
func typesetDocsCommand(command string) h.Node {
	var tabs, cmds []any
	for _, pm := range []string{"pnpm", "npm", "yarn", "bun"} {
		tabs = append(tabs, h.Button(h.Type("button"), h.Data("active", boolString(pm == "pnpm")), h.Data("typeset-pm", pm),
			h.Class("font-mono text-xs text-muted-foreground transition-colors hover:text-foreground data-[active=true]:text-foreground"), pm))
		cmds = append(cmds, h.Pre(h.Data("typeset-pm-command", pm), h.If(pm != "pnpm", h.Bool("hidden")),
			h.Class("scrollbar-none overflow-x-auto p-3 font-mono text-xs leading-relaxed"), commandFor(pm, command)))
	}
	return h.Div(h.Class("rounded-lg bg-muted/60"),
		h.Div(h.Class("flex items-center gap-3 px-3 pt-2"), h.Group(tabs...)),
		h.Group(cmds...),
	)
}

// commandFor is getCommandForPackageManager.
func commandFor(pm, command string) string {
	switch pm {
	case "pnpm":
		return strings.Replace(strings.Replace(command, "npm install", "pnpm add", 1), "npx", "pnpm dlx", 1)
	case "yarn":
		return strings.Replace(strings.Replace(command, "npm install", "yarn add", 1), "npx", "yarn dlx", 1)
	case "bun":
		return strings.Replace(strings.Replace(command, "npm install", "bun add", 1), "npx", "bunx --bun", 1)
	}
	return command
}

// typesetCode generates the docs' code for the params.
type typesetCode struct {
	origin    string
	p         typesetParams
	framework string
	fonts     []typesetFontDefinition
}

func newTypesetCode(r *http.Request, p typesetParams, framework string) *typesetCode {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	c := &typesetCode{origin: scheme + "://" + r.Host, p: p, framework: framework}
	seen := map[string]bool{}
	for _, id := range []string{p.Body, p.Heading, p.Mono} {
		if id == "inherit" || findTypesetFont(id) == nil || seen[id] {
			continue
		}
		seen[id] = true
		for _, d := range typesetData.Definitions {
			if d.Name == id {
				c.fonts = append(c.fonts, d)
			}
		}
	}
	return c
}

func (c *typesetCode) presetName() string { return "typeset-" + c.p.Item }

func (c *typesetCode) presetCSS() string {
	heading := c.p.Heading
	if heading == "inherit" {
		heading = c.p.Body
	}
	return "." + c.presetName() + " {\n" +
		"  --typeset-font-body: var(--font-" + c.p.Body + ");\n" +
		"  --typeset-font-heading: var(--font-" + heading + ");\n" +
		"  --typeset-font-mono: var(--font-" + c.p.Mono + ");\n" +
		"  --typeset-size: " + c.p.Scale + "px;\n" +
		"  --typeset-leading: " + c.p.Leading + ";\n" +
		"  --typeset-flow: " + c.p.Flow + ";\n}"
}

func (c *typesetCode) preset() string {
	return "/* globals.css */\n@import \"tailwindcss\";\n@import \"./typeset.css\";\n\n" + c.presetCSS()
}

func (c *typesetCode) usage() string {
	return "<div className=\"typeset " + c.presetName() + " max-w-[" + typesetMeasureWidth(c.p.Measure) + "]\">\n  {content}\n</div>"
}

var camelRe = regexp.MustCompile(`-(\w)`)

func toCamelCase(s string) string {
	return camelRe.ReplaceAllStringFunc(s, func(m string) string { return strings.ToUpper(m[1:]) })
}

func jsonString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// nextFontCode is getNextFontCode.
func (c *typesetCode) nextFontCode() string {
	imports := []string{}
	for _, f := range c.fonts {
		imports = append(imports, f.Import)
	}
	lines := []string{"// app/layout.tsx", "import { " + strings.Join(imports, ", ") + " } from \"next/font/google\"", ""}
	vars := []string{}
	for _, f := range c.fonts {
		block := []string{
			"const " + toCamelCase(f.Name) + " = " + f.Import + "({",
			"  subsets: " + jsonString(f.Subsets) + ",",
		}
		if f.Weight != nil {
			block = append(block, "  weight: "+jsonString(f.Weight)+",")
		}
		block = append(block, "  variable: \"--font-"+f.Name+"\",", "})", "")
		lines = append(lines, strings.Join(block, "\n"))
		vars = append(vars, "${"+toCamelCase(f.Name)+".variable}")
	}
	lines = append(lines, "<html className={`"+strings.Join(vars, " ")+"`}>")
	return strings.Join(lines, "\n")
}

// fontsourceCommand is getFontsourceCommand.
func (c *typesetCode) fontsourceCommand() string {
	deps := []string{}
	for _, f := range c.fonts {
		deps = append(deps, f.Dependency)
	}
	return "npm install " + strings.Join(deps, " ")
}

// fontsourceCSS is getFontsourceCss.
func (c *typesetCode) fontsourceCSS() string {
	lines := []string{"/* globals.css */"}
	for _, f := range c.fonts {
		lines = append(lines, "@import \""+f.Dependency+"\";")
	}
	lines = append(lines, "", ":root {")
	for _, f := range c.fonts {
		lines = append(lines, "  --font-"+f.Name+": "+f.Family+";")
	}
	return strings.Join(append(lines, "}"), "\n")
}

func (c *typesetCode) prompt() string {
	var fontStep string
	if c.framework == "next" {
		fontStep = "Load the fonts in the root layout and update the HTML element:\n\n" + c.nextFontCode()
	} else {
		fontStep = "Install the fonts:\n\n" + c.fontsourceCommand() + "\n\nThen import them in the main CSS file:\n\n" + c.fontsourceCSS()
	}
	return `Install shadcn/typeset in this project.

Typeset is a single stylesheet that styles rendered markdown: wrap the output in a ` + "`typeset`" + ` container and everything inside (headings, lists, tables, code, blockquotes, math) is styled. Everything outside is untouched.

1. Download ` + c.origin + `/typeset.css and save it as typeset.css next to the project's main CSS file (where Tailwind is imported). If the file already exists, replace it with the downloaded copy.

2. Import it in the main CSS file, after the Tailwind import:

@import "./typeset.css";

3. ` + fontStep + `

4. Add this preset to the main CSS file, after the typeset import. If a class named .` + c.presetName() + ` already exists, update its values in place. Leave any other typeset-* presets untouched: they are separate surfaces:

` + c.presetCSS() + `

5. Do not apply the class anywhere yet. Search the project for surfaces that render markdown or rich content: react-markdown, Streamdown, or MDX components, dangerouslySetInnerHTML with parsed markdown, prose classes, CMS content renderers. Present the candidates you find as a short list and ask the user which surface should use typeset. Then wrap only the surface they pick:

` + c.usage() + `

If the picked surface already has its own typography (a prose class, styled markdown components), list those styles and let the user decide what to remove before wrapping.

Notes:

- To exclude an embedded component from typeset styles, add the not-typeset class or the data-not-typeset attribute to it.
- Verify on the surface the user picked: headings, lists, tables, and code inside the container should be styled with no classes on the content itself.
- Docs: ` + c.origin + `/docs/typeset`
}
