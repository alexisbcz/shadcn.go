// Package tree lets a compound component pass context to its parts, the way
// React context reaches every descendant: the root walks its children
// (through plain elements, fragments and other components' builders) and
// hands state to the parts it finds before anything renders.
//
//	func (b *RootBuilder) Element() *h.Element {
//		tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
//			switch c := n.(type) {
//			case *ItemBuilder:
//				c.ctx = ctx
//			case *RootBuilder:
//				return false // a nested root provides its own context
//			}
//			return true
//		})
//		...
//	}
//
// Parts reached only through a lazily computed node (h.Lazy, a NodeFunc)
// are not visited.
package tree

import h "shadcn.go/libhtml"

// Container is implemented by builders embedding h.Builder.
type Container interface {
	Elem() *h.Element
}

// Walk calls visit for every node under children, depth first, in
// document order. visit returns false to skip a node's descendants.
func Walk(children []any, visit func(n any) bool) {
	for _, c := range children {
		walk(c, visit)
	}
}

func walk(c any, visit func(n any) bool) {
	if h.IsNil(c) {
		return
	}
	switch v := c.(type) {
	case nil:
		return
	case h.Fragment:
		Walk(v, visit)
		return
	case []any:
		Walk(v, visit)
		return
	case h.Scoped:
		// WithValue / style.Scope: no element of its own.
		Walk(v.Children(), visit)
		return
	}
	if !visit(c) {
		return
	}
	switch v := c.(type) {
	case *h.Element:
		Walk(v.ChildNodes(), visit)
	case Container:
		if el := v.Elem(); el != nil {
			Walk(el.ChildNodes(), visit)
		}
	}
}
