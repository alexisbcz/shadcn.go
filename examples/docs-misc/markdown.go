package docsmisc

import (
	"shadcn.go/components/style"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// markdown-demo renders upstream's Markdown (components/markdown.tsx:
// streamdown with its code plugin, shiki github-light/dark) for a fixed
// document; the tree below is streamdown's output for it.
func init() {
	registry.Register(registry.Example{Name: "nova--markdown-demo", Component: "markdown", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, MarkdownDemo())
	})})
}

func sd(name string) h.Prop { return h.Data("streamdown", name) }

func inlineCode(s string) h.Node {
	return h.Code(h.Class("rounded bg-muted px-1.5 py-0.5 font-mono text-sm"), sd("inline-code"), s)
}

func strong(s string) h.Node { return h.Span(h.Class("font-semibold"), sd("strong"), s) }

func li(kids ...any) h.Node {
	return h.Li(append([]any{h.Class("py-1 [&>p]:inline"), sd("list-item")}, kids...)...)
}

func tr(cell func(...any) *h.Element, class, slot string, a, b any) h.Node {
	return h.Tr(h.Class("border-border"), sd("table-row"),
		cell(h.Class(class), sd(slot), a),
		cell(h.Class(class), sd(slot), b))
}

const (
	kw    = "#D73A49|#F97583"
	plain = "#24292E|#E1E4E8"
)

func line(toks ...[2]string) h.Node {
	kids := []any{h.Class("block before:content-[counter(line)] before:inline-block before:[counter-increment:line] before:w-6 before:mr-4 before:text-[13px] before:text-right before:text-muted-foreground/50 before:font-mono before:select-none")}
	for _, t := range toks {
		light, dark := t[1][:7], t[1][8:]
		kids = append(kids, h.Span(h.Class("text-[var(--sdm-c,inherit)] dark:text-[var(--shiki-dark,var(--sdm-c,inherit))]"),
			h.Style("--sdm-c: "+light+"; --shiki-dark: "+dark+";"), t[0]))
	}
	return h.Span(kids...)
}

func MarkdownDemo() h.Node {
	const (
		fn  = "#6F42C1|#B392F0"
		par = "#E36209|#FFAB70"
		typ = "#005CC5|#79B8FF"
		tag = "#22863A|#85E89D"
	)
	return h.Div(h.Class("max-w-md"),
		h.Div(h.Class("space-y-4 whitespace-normal [&>*:first-child]:mt-0 [&>*:last-child]:mb-0 cn-markdown w-full min-w-0 overflow-hidden"),
			h.H2(h.Class("mt-6 mb-2 font-semibold text-2xl"), sd("heading-2"), "Getting started"),
			h.P("Markdown lets you write formatted text with a simple syntax."),
			h.H3(h.Class("mt-6 mb-2 font-semibold text-xl"), sd("heading-3"), "Features"),
			h.Ul(h.Class("list-inside list-disc whitespace-normal [li_&]:pl-6"), sd("unordered-list"),
				li(strong("Bold"), " and ", h.Em("italic"), " text"),
				li(h.Button(h.Class("wrap-anywhere appearance-none text-left font-medium text-primary underline"),
					h.Data("incomplete", "false"), sd("link"), h.Type("button"), "Links"), " and ", inlineCode("inline code")),
				li("Ordered and unordered lists"),
				li("Tables, blockquotes, and code blocks"),
			),
			h.Div(h.Class("my-4 flex flex-col gap-2 rounded-lg border border-border bg-sidebar p-2"), sd("table-wrapper"),
				h.Div(h.Class("border-collapse overflow-x-auto overflow-y-auto rounded-md border border-border bg-background"),
					h.Table(h.Class("w-full divide-y divide-border"), sd("table"),
						h.Thead(h.Class("bg-muted/80"), sd("table-header"),
							tr(h.Th, "whitespace-nowrap px-4 py-2 text-left font-semibold text-sm", "table-header-cell", "Syntax", "Result")),
						h.Tbody(h.Class("divide-y divide-border"), sd("table-body"),
							tr(h.Td, "px-4 py-2 text-sm", "table-cell", inlineCode("**bold**"), strong("bold")),
							tr(h.Td, "px-4 py-2 text-sm", "table-cell", inlineCode("*italic*"), h.Em("italic")),
							tr(h.Td, "px-4 py-2 text-sm", "table-cell", inlineCode("`code`"), inlineCode("code")),
						),
					),
				),
			),
			h.P("How about a quote? How does this look?"),
			h.Blockquote(h.Class("my-4 border-muted-foreground/30 border-l-4 pl-4 text-muted-foreground italic"), sd("blockquote"),
				h.P("The best way to learn markdown is to write it. Lorem ipsum dolor sit amet consectetur adipisicing elit. Quisquam, quos.")),
			h.P("Now let's try a code block. It should have line numbers, a copy button and syntax highlighting."),
			h.Div(h.Class("my-4 flex w-full flex-col gap-2 rounded-xl border border-border bg-sidebar p-2"), h.Data("language", "tsx"), sd("code-block"),
				h.Style("content-visibility: auto; contain-intrinsic-size: auto 200px;"),
				h.Div(h.Class("flex h-8 items-center text-muted-foreground text-xs"), h.Data("language", "tsx"), sd("code-block-header"),
					h.Span(h.Class("ml-1 font-mono lowercase"), "tsx")),
				h.Div(h.Class("language-tsx overflow-x-auto rounded-md border border-border bg-background p-4 text-sm"), h.Data("language", "tsx"), sd("code-block-body"),
					h.Pre(h.Class("language-tsx bg-[var(--sdm-bg,inherit] dark:bg-[var(--shiki-dark-bg,var(--sdm-bg,inherit)]"),
						h.Style("--sdm-bg: transparent; --sdm-fg: inherit;"),
						h.Code(h.Class("[counter-increment:line_0] [counter-reset:line]"),
							line([2]string{"export", kw}, [2]string{" ", plain}, [2]string{"function", kw}, [2]string{" ", plain},
								[2]string{"Greeting", fn}, [2]string{"({ ", plain}, [2]string{"name", par}, [2]string{" }", plain},
								[2]string{":", kw}, [2]string{" { ", plain}, [2]string{"name", par}, [2]string{":", kw},
								[2]string{" ", plain}, [2]string{"string", typ}, [2]string{" }) {", plain}),
							line([2]string{"  ", plain}, [2]string{"return", kw}, [2]string{" <", plain}, [2]string{"p", tag},
								[2]string{">Hello, {name}!</", plain}, [2]string{"p", tag}, [2]string{">", plain}),
							line([2]string{"}", plain}),
						),
					),
				),
			),
			h.P("If you need more plugins, you can pass them to the ", inlineCode("Markdown"), " component."),
		),
	)
}
