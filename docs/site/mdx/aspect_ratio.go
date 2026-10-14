package mdx

import (
	aspectratio "shadcn.go/components/ui/aspect-ratio"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// AspectRatio: the new-york-v4 aspect-ratio. Props: ratio, className.
func AspectRatio(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	args := nodes(children)
	if r, ok := props.Get("ratio"); ok {
		if f, ok := r.(float64); ok && f > 0 {
			args = append(args, aspectratio.Ratio(f))
		}
	}
	if c := className(props); c != "" {
		args = append(args, h.Class(c))
	}
	return aspectratio.AspectRatio(args...), nil
}

func init() { Register("AspectRatio", AspectRatio) }
