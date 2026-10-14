package mdx

import (
	"shadcn.go/components/ui/accordion"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// Accordion, AccordionItem, AccordionTrigger, AccordionContent: the
// new-york-v4 accordion (Accordion gets className "not-typeset" from
// package content). Props: type, collapsible, defaultValue, value,
// disabled, className.
func Accordion(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	b := accordion.Accordion(nodes(children)...).Class(className(props))
	if props.String("type") == "multiple" {
		b.Type(accordion.TypeMultiple)
	}
	if v, _ := props.Get("collapsible"); v == true {
		b.Collapsible(true)
	}
	if v, _ := props.Get("disabled"); v == true {
		b.Disabled(true)
	}
	for _, name := range []string{"value", "defaultValue"} {
		switch v, _ := props.Get(name); v := v.(type) {
		case string:
			b.Value(v)
		case []any:
			var vs []string
			for _, x := range v {
				if s, ok := x.(string); ok {
					vs = append(vs, s)
				}
			}
			b.Value(vs...)
		}
	}
	return b, nil
}

func AccordionItem(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	b := accordion.Item(nodes(children)...).Value(props.String("value")).Class(className(props))
	if v, _ := props.Get("disabled"); v == true {
		b.Disabled(true)
	}
	return b, nil
}

func AccordionTrigger(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return accordion.Trigger(nodes(children)...).Class(className(props)), nil
}

func AccordionContent(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
	return accordion.Content(nodes(children)...).Class(className(props)), nil
}

func init() {
	Register("Accordion", Accordion)
	Register("AccordionItem", AccordionItem)
	Register("AccordionTrigger", AccordionTrigger)
	Register("AccordionContent", AccordionContent)
}
