package mdx

import (
	"shadcn.go/components/ui/kbd"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Kbd: the new-york-v4 Kbd.
func Kbd(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	args := nodes(children)
	if c := className(props); c != "" {
		args = append([]any{h.Class(c)}, args...)
	}
	return kbd.Kbd(args...), nil
}

func init() { Register("Kbd", Kbd) }
