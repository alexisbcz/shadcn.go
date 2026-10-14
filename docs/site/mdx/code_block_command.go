package mdx

import (
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// CodeBlockCommand: ui/apps/v4/components/code-block-command.tsx, the
// package manager tabs (pnpm, npm, yarn, bun) that replace the <code> of
// install commands. The selected manager is useConfig's packageManager
// (localStorage "config", default pnpm): docs-mdx.js restores and
// persists it, and keeps every command block of the page in sync.
func CodeBlockCommand(_ *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	managers := []string{"pnpm", "npm", "yarn", "bun"}
	list := tabs.List().Class("rounded-none bg-transparent p-0")
	panels := h.Div(h.Class("no-scrollbar overflow-x-auto"))
	for _, m := range managers {
		list.Children(tabs.Trigger(m).Value(m).
			Class("h-7 border border-transparent pt-0.5 shadow-none! data-[state=active]:border-input data-[state=active]:bg-background!"))
		panels.Children(tabs.Content(
			h.Pre(h.Code(h.Class("relative font-mono text-sm leading-none"), h.Data("language", "bash"), props.String("__"+m+"__"))),
		).Value(m).Class("mt-0 px-4 py-3.5"))
	}
	return h.Div(h.Class("overflow-x-auto"),
		tabs.Tabs(
			h.Div(h.Class("flex items-center gap-2 border-b border-border/50 px-3 py-1"),
				h.Div(h.Class("flex size-4 items-center justify-center rounded-[1px] bg-foreground opacity-70"),
					iconTerminal(h.Class("size-3 text-code"))),
				list,
			),
			panels,
		).Value("pnpm").Class("gap-0").Data("docs-config", "packageManager"),
		button.Button(h.Span(h.Class("sr-only"), "Copy"), tabler.IconCopy()).
			Data("slot", "copy-button").Size(button.SizeIcon).Variant(button.VariantGhost).
			Class("absolute top-2 right-2 z-10 size-7 opacity-70 hover:opacity-100 focus-visible:opacity-100"),
	), nil
}

func init() { Register("CodeBlockCommand", CodeBlockCommand) }
