package mdx

import (
	"shadcn.go/cn"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/separator"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// CodeCollapsibleWrapper: ui/apps/v4/components/code-collapsible-wrapper.tsx.
// A Radix Collapsible (closed) whose content is force-mounted, with two
// triggers ("Expand"/"Collapse"); docs-mdx.js toggles it.
func CodeCollapsibleWrapper(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return codeCollapsibleWrapper(className(props), nodes(children)...), nil
}

func codeCollapsibleWrapper(class string, children ...any) h.Node {
	id := newID("collapsible")
	trigger := func(el *h.Element) *h.Element {
		return el.Attr("type", "button").Aria("controls", id).Aria("expanded", "false").
			Data("state", "closed").Data("slot", "collapsible-trigger")
	}
	return h.Div(h.Data("state", "closed"), h.Data("slot", "collapsible"), h.Data("docs-collapsible", ""),
		h.Class(cn.Cn("group/collapsible relative md:-mx-1", class)),
		trigger(h.Div(h.Class("absolute top-1.5 right-9 z-10 flex items-center"),
			button.Button("Expand").Variant(button.VariantGhost).Size(button.SizeSm).
				Class("h-7 rounded-md px-2 text-muted-foreground").Data("docs-toggle-label", ""),
			separator.Separator(separator.Vertical).Class("mx-1.5 h-4!"),
		)),
		h.Div(h.Data("state", "closed"), h.ID(id), h.Data("slot", "collapsible-content"),
			h.Class("relative mt-6 overflow-hidden data-[state=closed]:max-h-64 data-[state=closed]:[content-visibility:auto] [&>figure]:mt-0 [&>figure]:md:mx-0!"),
			h.Group(children...),
		),
		trigger(h.Button(h.Class("absolute inset-x-0 -bottom-2 flex h-20 items-center justify-center rounded-b-lg bg-gradient-to-b from-code/70 to-code text-sm text-muted-foreground group-data-[state=open]/collapsible:hidden"),
			h.Data("docs-toggle-label", ""), "Expand")),
	)
}

func init() { Register("CodeCollapsibleWrapper", CodeCollapsibleWrapper) }
