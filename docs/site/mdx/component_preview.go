package mdx

import (
	"strings"

	"shadcn.go/components/style"
	"shadcn.go/docs/content"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// ComponentPreview: ui/apps/v4/components/component-preview.tsx and
// component-preview-tabs.tsx. Props: name, styleName ("radix-nova"), type,
// className, previewClassName, align, hideCode, chromeLessOnMobile,
// direction, caption.
//
// Placeholder: the Go example in a bordered preview box, rendered in
// radix-nova (the explicit "nova--<name>" registration if any, otherwise
// the example of the same name in style.Scope(style.Nova), like
// cmd/parity). TODO: the code panel (ComponentSource of the example), the
// block type (iframe + screenshots), caption, alignment.
func ComponentPreview(env *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	name := props.String("name")
	var preview h.Node
	if ex, ok := LookupExample(name); ok {
		preview = ex(env)
	} else {
		preview = h.P(h.Class("text-sm text-muted-foreground"), "Example ", h.Code(name), " is not ported yet.")
	}
	align := map[string]string{"start": "items-start", "end": "items-end"}[props.String("align")]
	if align == "" {
		align = "items-center"
	}
	return h.Div(h.Data("slot", "component-preview"), h.Data("name", name),
		h.Class("group relative mt-4 mb-12 flex flex-col overflow-hidden rounded-xl border "+className(props)),
		h.Data("not-typeset", ""),
		h.Div(h.Data("slot", "preview"),
			h.Class("preview flex h-[450px] w-full justify-center p-10 "+align+" "+props.String("previewClassName")),
			preview,
		),
	), nil
}

// LookupExample finds the Go port of a docs example by its upstream name
// (ui/apps/v4/examples/radix/<name>.tsx) and returns a function rendering
// it in radix-nova.
func LookupExample(name string) (func(env *Env) h.Node, bool) {
	if ex, ok := registry.Get("nova--" + name); ok {
		return func(env *Env) h.Node { return ex.Render(env.Request) }, true
	}
	ex, ok := registry.Get(strings.TrimPrefix(name, "nova--"))
	if !ok {
		return nil, false
	}
	return func(env *Env) h.Node { return style.Scope(style.Nova, ex.Render(env.Request)) }, true
}

func init() { Register("ComponentPreview", ComponentPreview) }
