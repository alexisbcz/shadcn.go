package mdx

import (
	"shadcn.go/components/ui/button"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Button: the new-york-v4 Button with className "not-typeset" merged by
// package content. Props: variant, size, asChild, className.
// TODO: other props.
func Button(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	b := button.Button(nodes(children)...).Class(className(props))
	if v := props.String("variant"); v != "" {
		b.Variant(button.Variant(v))
	}
	if s := props.String("size"); s != "" {
		b.Size(button.Size(s))
	}
	if v, _ := props.Get("asChild"); v == true {
		b.AsChild()
	}
	return b, nil
}

func init() { Register("Button", Button) }
