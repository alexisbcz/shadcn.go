// Package breadcrumb is the shadcn/ui Breadcrumb: the path to the current
// page as a list of links.
//
//	breadcrumb.Breadcrumb(breadcrumb.List(
//		breadcrumb.Item(breadcrumb.Link(h.Href("/"), "Home")),
//		breadcrumb.Separator(),
//		breadcrumb.Item(breadcrumb.Page("Breadcrumb")),
//	))
package breadcrumb

import (
	"fmt"
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

// Builder is returned by every Breadcrumb sub-component.
type Builder struct {
	h.Builder[Builder]
	base       *h.Element
	asChild    bool
	acceptsAs  bool
	defaultKid func() *h.Element // rendered when the user gives no children
	fixedKids  bool              // the user's children are ignored (Ellipsis)
}

func styled(base *h.Element, tag string, children []any) *Builder {
	b := &Builder{base: base}
	b.Builder = h.NewBuilder(b, tag)
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild (Link only).
func (b *Builder) AcceptChild(c any) bool {
	if _, ok := c.(lib.AsChildFlag); ok && b.acceptsAs {
		b.asChild = true
		return true
	}
	return false
}

// AsChild makes Link render its single element child (e.g. a router link)
// with the link's classes and attributes merged in.
func (b *Builder) AsChild() *Builder { b.asChild = true; return b }

// Breadcrumb is the root <nav aria-label="breadcrumb">.
func Breadcrumb(children ...any) *Builder {
	return styled(h.Nav(h.Aria("label", "breadcrumb"), h.Data("slot", "breadcrumb")), "nav", children)
}

// List is the ordered list of items.
func List(children ...any) *Builder {
	return styled(h.Ol(h.Data("slot", "breadcrumb-list"), h.Class("flex flex-wrap items-center gap-1.5 text-sm break-words text-muted-foreground sm:gap-2.5")), "ol", children)
}

// Item is one entry of the list.
func Item(children ...any) *Builder {
	return styled(h.Li(h.Data("slot", "breadcrumb-item"), h.Class("inline-flex items-center gap-1.5")), "li", children)
}

// Link is a link to an ancestor page. Pass lib.AsChild to style your own
// link element instead of rendering an <a>.
func Link(children ...any) *Builder {
	b := &Builder{base: h.A(h.Data("slot", "breadcrumb-link"), h.Class("transition-colors hover:text-foreground")), acceptsAs: true}
	b.Builder = h.NewBuilder(b, "a")
	return b.Children(children...)
}

// Page is the current page.
func Page(children ...any) *Builder {
	return styled(h.Span(h.Data("slot", "breadcrumb-page"), h.Role("link"), h.Aria("disabled", "true"), h.Aria("current", "page"), h.Class("font-normal text-foreground")), "span", children)
}

// Separator separates two items; it shows a chevron unless given children.
func Separator(children ...any) *Builder {
	b := styled(h.Li(h.Data("slot", "breadcrumb-separator"), h.Role("presentation"), h.Aria("hidden", "true"), h.Class("[&>svg]:size-3.5")), "li", children)
	b.defaultKid = func() *h.Element { return lucide.ChevronRight() }
	return b
}

// Ellipsis stands for collapsed items. Like upstream, it always renders
// its icon and "More" label: child nodes are ignored, attributes apply.
func Ellipsis(children ...any) *Builder {
	b := styled(h.Span(h.Data("slot", "breadcrumb-ellipsis"), h.Role("presentation"), h.Aria("hidden", "true"), h.Class("flex size-9 items-center justify-center"),
		lucide.Ellipsis(h.Class("size-4")),
		h.Span(h.Class("sr-only"), "More"),
	), "span", children)
	b.fixedKids = true
	return b
}

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element {
	props := b.Elem()
	if b.fixedKids {
		props = props.Clone().SetChildren()
	} else if b.defaultKid != nil && len(props.ChildNodes()) == 0 {
		props = props.Clone().SetChildren(b.defaultKid())
	}
	el := h.Spread(b.base, props, lib.Cn)
	if !b.asChild {
		return el
	}
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("breadcrumb: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	// Radix Slot: the child's props win, classes are concatenated.
	out, err := h.Slot(el.Clone().SetChildren(), kids[0], nil)
	if err != nil {
		return h.Errored(fmt.Errorf("breadcrumb: AsChild: %w", err))
	}
	return out
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }
