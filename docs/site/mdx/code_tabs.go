package mdx

import (
	"shadcn.go/components/ui/tabs"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// CodeTabs: ui/apps/v4/components/code-tabs.tsx. Tabs whose value is the
// remembered installation type (useConfig: "cli" by default).
// TODO: persist the choice like useConfig (localStorage "config").
func CodeTabs(_ *Env, _ *content.Context, _ content.Props, children []h.Node) (h.Node, error) {
	return tabs.Tabs(nodes(children)...).Value("cli").Class("relative mt-6 w-full *:data-[slot=tabs-list]:gap-6"), nil
}

func init() { Register("CodeTabs", CodeTabs) }
