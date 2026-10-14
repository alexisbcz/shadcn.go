package mdx

import (
	"shadcn.go/components/ui/alert"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Alert, AlertTitle, AlertDescription: the new-york-v4 alert.
func Alert(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	b := alert.Alert(nodes(children)...).Class(className(props))
	if v := props.String("variant"); v != "" {
		b.Variant(alert.Variant(v))
	}
	return b, nil
}

func AlertTitle(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return alert.Title(append([]any{h.Class(className(props))}, nodes(children)...)...), nil
}

func AlertDescription(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return alert.Description(append([]any{h.Class(className(props))}, nodes(children)...)...), nil
}

func init() {
	Register("Alert", Alert)
	Register("AlertTitle", AlertTitle)
	Register("AlertDescription", AlertDescription)
}
