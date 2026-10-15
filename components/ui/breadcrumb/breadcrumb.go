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
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/part"

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
	isList bool // List: React Aria items hold their separator (style.BaseAria)
	isSep  bool // Separator
	// alt returns the base element and default child of a style
	// (registry/bases/radix/ui/breadcrumb.tsx), for style.Scope.
	alt func(st style.Style) (base *h.Element, defaultKid func() *h.Element)
}

// styleFile is registry/bases/radix/ui/breadcrumb.tsx (style.Scope).
var styleFile = style.NewFile("ui/breadcrumb.tsx")

var (
	styleRoot      = styleFile.Cn("cn-breadcrumb", style.Arg)
	styleList      = styleFile.Cn("cn-breadcrumb-list flex flex-wrap items-center wrap-break-word", style.Arg)
	styleItem      = styleFile.Cn("cn-breadcrumb-item inline-flex items-center", style.Arg)
	styleLink      = styleFile.Cn("cn-breadcrumb-link", style.Arg)
	stylePage      = styleFile.Cn("cn-breadcrumb-page", style.Arg)
	styleSeparator = styleFile.Cn("cn-breadcrumb-separator", style.Arg)
	styleSepIcon   = styleFile.Attr("cn-rtl-flip")
	styleEllipsis  = styleFile.Cn("cn-breadcrumb-ellipsis flex items-center justify-center", style.Arg)
)

// withAlt sets the part's style base: the new-york-v4 base with its class
// replaced by site's (the user's class is merged when spreading).
func (b *Builder) withAlt(site *style.Site, kid func(st style.Style) func() *h.Element) *Builder {
	b.alt = func(st style.Style) (*h.Element, func() *h.Element) {
		base := b.base.Clone().SetClass(site.Class(st, ""))
		var k func() *h.Element
		if kid != nil {
			k = kid(st)
		}
		return base, k
	}
	return b
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
	return styled(h.Nav(h.Aria("label", "breadcrumb"), h.Data("slot", "breadcrumb")), "nav", children).withAlt(styleRoot, nil)
}

// List is the ordered list of items.
func List(children ...any) *Builder {
	b := styled(h.Ol(h.Data("slot", "breadcrumb-list"), h.Class("flex flex-wrap items-center gap-1.5 text-sm break-words text-muted-foreground sm:gap-2.5")), "ol", children).withAlt(styleList, nil)
	b.isList = true
	return b
}

// Item is one entry of the list.
func Item(children ...any) *Builder {
	return styled(h.Li(h.Data("slot", "breadcrumb-item"), h.Class("inline-flex items-center gap-1.5")), "li", children).withAlt(styleItem, nil)
}

// Link is a link to an ancestor page. Pass lib.AsChild to style your own
// link element instead of rendering an <a>.
func Link(children ...any) *Builder {
	b := &Builder{base: h.A(h.Data("slot", "breadcrumb-link"), h.Class("transition-colors hover:text-foreground")), acceptsAs: true}
	b.Builder = h.NewBuilder(b, "a")
	b.withAlt(styleLink, nil)
	return b.Children(children...)
}

// Page is the current page.
func Page(children ...any) *Builder {
	return styled(h.Span(h.Data("slot", "breadcrumb-page"), h.Role("link"), h.Aria("disabled", "true"), h.Aria("current", "page"), h.Class("font-normal text-foreground")), "span", children).withAlt(stylePage, nil)
}

// Separator separates two items; it shows a chevron unless given children.
func Separator(children ...any) *Builder {
	b := styled(h.Li(h.Data("slot", "breadcrumb-separator"), h.Role("presentation"), h.Aria("hidden", "true"), h.Class("[&>svg]:size-3.5")), "li", children)
	b.defaultKid = func() *h.Element { return lucide.ChevronRight() }
	b.isSep = true
	return b.withAlt(styleSeparator, func(st style.Style) func() *h.Element {
		return func() *h.Element { return lucide.ChevronRight(h.Class(styleSepIcon.Class(st))) }
	})
}

// Ellipsis stands for collapsed items. Like upstream, it always renders
// its icon and "More" label: child nodes are ignored, attributes apply.
func Ellipsis(children ...any) *Builder {
	b := styled(h.Span(h.Data("slot", "breadcrumb-ellipsis"), h.Role("presentation"), h.Aria("hidden", "true"), h.Class("flex size-9 items-center justify-center"),
		lucide.Ellipsis(h.Class("size-4")),
		h.Span(h.Class("sr-only"), "More"),
	), "span", children)
	b.fixedKids = true
	b.alt = func(st style.Style) (*h.Element, func() *h.Element) {
		return h.Span(h.Data("slot", "breadcrumb-ellipsis"), h.Role("presentation"), h.Aria("hidden", "true"), h.Class(styleEllipsis.Class(st, "")),
			lucide.Ellipsis(),
			h.Span(h.Class("sr-only"), "More"),
		), nil
	}
	return b
}

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element { return b.element(nil, "") }

// ElementIn resolves the component as it renders inside w, in the style of
// an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return b.element(w, st)
}

func (b *Builder) element(w io.Writer, st style.Style) *h.Element {
	props := b.Elem()
	base, defaultKid := b.base, b.defaultKid
	if st != "" && b.alt != nil {
		base, defaultKid = b.alt(st)
		if part.RTL(w) {
			// The ui-rtl variant: the part's own classes mapped (the
			// default separator's cn-rtl-flip included).
			base = base.Clone().SetClass(style.RTLClass(base.ClassName()))
			if defaultKid != nil {
				k := defaultKid
				defaultKid = func() *h.Element { e := k(); return e.SetClass(style.RTLClass(e.ClassName())) }
			}
		}
	}
	if b.fixedKids {
		props = props.Clone().SetChildren()
	} else if defaultKid != nil && len(props.ChildNodes()) == 0 {
		props = props.Clone().SetChildren(defaultKid())
	}
	if b.isList && w != nil && style.BaseFromWriter(w) == style.BaseAria {
		props = ariaItems(w, props)
	}
	el := h.Spread(base, props, lib.Cn)
	if !b.asChild {
		return el
	}
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("breadcrumb: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	// Radix Slot: the child's props win, classes are concatenated.
	var out *h.Element
	var err error
	if w != nil {
		out, err = h.SlotIn(w, el.Clone().SetChildren(), kids[0], nil)
	} else {
		out, err = h.Slot(el.Clone().SetChildren(), kids[0], nil)
	}
	if err != nil {
		return h.Errored(fmt.Errorf("breadcrumb: AsChild: %w", err))
	}
	return out
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ariaItems lays out a list's children as React Aria's Breadcrumbs does
// (registry/bases/aria/ui/breadcrumb.tsx): an item holds the separator after
// it, as a <span>, instead of the separator being a sibling <li>.
func ariaItems(w io.Writer, props *h.Element) *h.Element {
	kids := props.ChildNodes()
	out := make([]any, 0, len(kids))
	var last *h.Element
	for _, k := range kids {
		kb, ok := k.(*Builder)
		if ok && kb.isSep && last != nil {
			sep := kb.ElementIn(w)
			span := h.Span(h.Data("slot", "breadcrumb-separator"), h.Role("presentation"), h.Aria("hidden", "true"),
				h.Class(sep.ClassName()), h.Group(sep.ChildNodes()...))
			last.Children(span)
			continue
		}
		if ok && !kb.isSep && !kb.isList {
			last = kb.ElementIn(w).Clone()
			out = append(out, last)
			continue
		}
		last = nil
		out = append(out, k)
	}
	return props.Clone().SetChildren(out...)
}
