// Package pagination is the shadcn/ui Pagination: page navigation with
// next and previous links.
//
//	pagination.Pagination(pagination.Content(
//		pagination.Item(pagination.Previous(h.Href("#"))),
//		pagination.Item(pagination.Link(h.Href("#"), "1")),
//		pagination.Item(pagination.Link(h.Href("#"), pagination.IsActive, "2")),
//		pagination.Item(pagination.Ellipsis()),
//		pagination.Item(pagination.Next(h.Href("#"))),
//	))
package pagination

import (
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	h "shadcn.go/libhtml"
)

// ActiveFlag is the type of IsActive.
type ActiveFlag struct{}

// IsActive, passed as a child of Link, marks the current page.
var IsActive = ActiveFlag{}

var (
	navBase     = h.Nav(h.Role("navigation"), h.Aria("label", "pagination"), h.Data("slot", "pagination"), h.Class(lib.Cn("mx-auto flex w-full justify-center")))
	contentBase = h.Ul(h.Data("slot", "pagination-content"), h.Class(lib.Cn("flex flex-row items-center gap-1")))
	itemBase    = h.Li(h.Data("slot", "pagination-item"))
)

// Pagination is the root <nav>.
func Pagination(children ...any) *h.Element {
	return h.Spread(navBase, h.Nav(children...), lib.Cn)
}

// Content is the <ul> of items.
func Content(children ...any) *h.Element {
	return h.Spread(contentBase, h.Ul(children...), lib.Cn)
}

// Item is one <li>.
func Item(children ...any) *h.Element {
	return h.Spread(itemBase, h.Li(children...), lib.Cn)
}

// Builder is returned by Link, Previous and Next.
type Builder struct {
	h.Builder[Builder]
	active bool
	// activeSet records that the active state was given explicitly:
	// like React's data-active={isActive}, an explicit false renders
	// data-active="false".
	activeSet bool
	size      button.Size
	kind      int
}

const (
	kindLink = iota
	kindPrevious
	kindNext
)

// Link is a page link styled as a button: ghost, or outline with
// aria-current="page" when given IsActive. Children may include a
// button.Size (icon by default).
func Link(children ...any) *Builder {
	b := &Builder{size: button.SizeIcon}
	b.Builder = h.NewBuilder(b, "a")
	return b.Children(children...)
}

// Previous links to the previous page.
func Previous(children ...any) *Builder {
	b := &Builder{size: button.SizeDefault, kind: kindPrevious}
	b.Builder = h.NewBuilder(b, "a")
	return b.Children(children...)
}

// Next links to the next page.
func Next(children ...any) *Builder {
	b := &Builder{size: button.SizeDefault, kind: kindNext}
	b.Builder = h.NewBuilder(b, "a")
	return b.Children(children...)
}

// AcceptChild consumes IsActive and button.Size.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case ActiveFlag:
		b.active, b.activeSet = true, true
	case button.Size:
		b.size = v
		if v == "" {
			b.size = button.SizeIcon
		}
	default:
		return false
	}
	return true
}

// Active sets whether the link is the current page.
func (b *Builder) Active(on bool) *Builder { b.active, b.activeSet = on, true; return b }

// Size sets the link's button size.
func (b *Builder) Size(s button.Size) *Builder { b.AcceptChild(s); return b }

// Element resolves the component to its <a>.
func (b *Builder) Element() *h.Element {
	props := b.Elem().Clone()
	switch b.kind {
	case kindPrevious:
		// <PaginationLink aria-label size="default" className {...props}>
		// children</PaginationLink>: the explicit children win.
		className := lib.Cn("gap-1 px-2.5 sm:pl-2.5", props.ClassName())
		props.SetClass(className)
		props.SetChildren(lucide.ChevronLeft(), h.Span(h.Class("hidden sm:block"), "Previous"))
		return link(b.active, b.activeSet, b.size, h.Spread(h.A(h.Aria("label", "Go to previous page")), props, nil))
	case kindNext:
		className := lib.Cn("gap-1 px-2.5 sm:pr-2.5", props.ClassName())
		props.SetClass(className)
		props.SetChildren(h.Span(h.Class("hidden sm:block"), "Next"), lucide.ChevronRight())
		return link(b.active, b.activeSet, b.size, h.Spread(h.A(h.Aria("label", "Go to next page")), props, nil))
	}
	return link(b.active, b.activeSet, b.size, props)
}

func link(active, activeSet bool, size button.Size, props *h.Element) *h.Element {
	variant := button.VariantGhost
	base := h.A(h.Data("slot", "pagination-link"))
	if active {
		variant = button.VariantOutline
		base.Aria("current", "page").Data("active", "true")
	} else if activeSet {
		base.Data("active", "false")
	}
	className := props.ClassName()
	props = props.Clone()
	props.SetClass("")
	base.Class(button.Variants(variant, size, className))
	return h.Spread(base, props, lib.Cn)
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// Ellipsis stands for skipped pages. Like upstream, it always renders its
// icon and "More pages" label: child nodes are ignored, attributes apply.
func Ellipsis(children ...any) *h.Element {
	props := h.Span(children...).SetChildren()
	base := h.Span(h.Aria("hidden", "true"), h.Data("slot", "pagination-ellipsis"), h.Class(lib.Cn("flex size-9 items-center justify-center")),
		lucide.Ellipsis(h.Class("size-4")),
		h.Span(h.Class("sr-only"), "More pages"),
	)
	return h.Spread(base, props, lib.Cn)
}
