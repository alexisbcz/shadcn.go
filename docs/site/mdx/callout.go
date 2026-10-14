package mdx

import (
	"shadcn.go/cn"
	"shadcn.go/components/ui/alert"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Callout: ui/apps/v4/components/callout.tsx (new-york-v4 Alert). Props:
// title, icon (an element), variant, className.
func Callout(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	variant := props.String("variant")
	if variant == "" {
		variant = "default"
	}
	var kids []any
	if icon, ok := props.Get("icon"); ok {
		if n, ok := icon.(h.Node); ok {
			kids = append(kids, n)
		}
	}
	if t := props.String("title"); t != "" {
		kids = append(kids, alert.Title(t))
	}
	kids = append(kids, alert.Description(append([]any{h.Class("text-card-foreground/80")}, nodes(children)...)...))
	return alert.Alert(kids...).Data("variant", variant).
		Class(cn.Cn("not-typeset mt-6 w-auto rounded-2xl border-surface bg-surface text-surface-foreground md:-mx-1 **:[code]:border", className(props))), nil
}

func init() { Register("Callout", Callout) }
