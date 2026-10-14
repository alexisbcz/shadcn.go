package mdx

import (
	"shadcn.go/components/ui/tabs"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Tabs, TabsList, TabsTrigger, TabsContent: the new-york-v4 tabs
// (registry/new-york-v4/ui/tabs), with the className defaults of
// mdx-components.tsx already merged by package content. The Go tabs wire
// triggers and panels together at render time, so the parts can be built
// independently. TODO: extra props (orientation, ...).
func Tabs(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	b := tabs.Tabs(nodes(children)...).Class(className(props))
	if v := props.String("defaultValue"); v != "" {
		b.Value(v)
	}
	return b, nil
}

func TabsList(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return tabs.List(nodes(children)...).Class(className(props)), nil
}

func TabsTrigger(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return tabs.Trigger(nodes(children)...).Value(props.String("value")).Class(className(props)), nil
}

func TabsContent(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return tabs.Content(nodes(children)...).Value(props.String("value")).Class(className(props)), nil
}

func init() {
	Register("Tabs", Tabs)
	Register("TabsList", TabsList)
	Register("TabsTrigger", TabsTrigger)
	Register("TabsContent", TabsContent)
}
