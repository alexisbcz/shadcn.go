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
	"cmp"
	"io"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/part"
	"shadcn.go/cva"

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

// styleFile is registry/bases/radix/ui/pagination.tsx (style.Scope).
var styleFile = style.NewFile("ui/pagination.tsx")

var (
	styleRoot         = styleFile.Cn("cn-pagination mx-auto flex w-full justify-center", style.Arg)
	styleContent      = styleFile.Cn("cn-pagination-content flex items-center", style.Arg)
	styleLink         = styleFile.Cn("cn-pagination-link", style.Arg)
	stylePrevious     = styleFile.Cn("cn-pagination-previous", style.Arg)
	styleNext         = styleFile.Cn("cn-pagination-next", style.Arg)
	styleFlip         = styleFile.Attr("cn-rtl-flip")
	stylePreviousText = styleFile.Attr("cn-pagination-previous-text hidden sm:block")
	styleNextText     = styleFile.Attr("cn-pagination-next-text hidden sm:block")
	styleEllipsis     = styleFile.Cn("cn-pagination-ellipsis flex items-center justify-center", style.Arg)
)

// buttonFile is registry/bases/radix/ui/button.tsx: PaginationLink is a
// Button asChild.
var buttonFile = style.NewFile("ui/button.tsx")

var buttonVariants = buttonFile.CVA(
	"cn-button group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
	[]style.Variant{
		{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-button-variant-default"},
			{Name: "outline", Class: "cn-button-variant-outline"},
			{Name: "secondary", Class: "cn-button-variant-secondary"},
			{Name: "ghost", Class: "cn-button-variant-ghost"},
			{Name: "destructive", Class: "cn-button-variant-destructive"},
			{Name: "link", Class: "cn-button-variant-link"},
		}},
		{Name: "size", Options: []style.Option{
			{Name: "default", Class: "cn-button-size-default"},
			{Name: "xs", Class: "cn-button-size-xs"},
			{Name: "sm", Class: "cn-button-size-sm"},
			{Name: "lg", Class: "cn-button-size-lg"},
			{Name: "icon", Class: "cn-button-size-icon"},
			{Name: "icon-xs", Class: "cn-button-size-icon-xs"},
			{Name: "icon-sm", Class: "cn-button-size-icon-sm"},
			{Name: "icon-lg", Class: "cn-button-size-icon-lg"},
		}},
	},
	cva.Selection{"variant": "default", "size": "default"},
)

// Pagination is the root <nav>.
func Pagination(children ...any) *part.Part {
	return part.New("nav", func(st style.Style, props *h.Element) *h.Element {
		if st != "" {
			base := navBase.Clone().SetClass(styleRoot.Class(st, ""))
			return h.Spread(base, props, lib.Cn)
		}
		return h.Spread(navBase, props, lib.Cn)
	}, children...)
}

// Content is the <ul> of items.
func Content(children ...any) *part.Part {
	return part.New("ul", func(st style.Style, props *h.Element) *h.Element {
		if st != "" {
			base := contentBase.Clone().SetClass(styleContent.Class(st, ""))
			return h.Spread(base, props, lib.Cn)
		}
		return h.Spread(contentBase, props, lib.Cn)
	}, children...)
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
	text      string
}

// Text, passed to Previous or Next, replaces its label ("Previous",
// "Next"), like the text prop of registry/bases/radix.
type Text string

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
	case Text:
		b.text = string(v)
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
func (b *Builder) Element() *h.Element { return b.element("") }

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return b.element(st)
}

func (b *Builder) element(st style.Style) *h.Element {
	props := b.Elem().Clone()
	if st != "" {
		return b.styled(st, props)
	}
	switch b.kind {
	case kindPrevious:
		// <PaginationLink aria-label size="default" className {...props}>
		// children</PaginationLink>: the explicit children win.
		className := lib.Cn("gap-1 px-2.5 sm:pl-2.5", props.ClassName())
		props.SetClass(className)
		props.SetChildren(lucide.ChevronLeft(), h.Span(h.Class("hidden sm:block"), cmp.Or(b.text, "Previous")))
		return link(b.active, b.activeSet, b.size, h.Spread(h.A(h.Aria("label", "Go to previous page")), props, nil))
	case kindNext:
		className := lib.Cn("gap-1 px-2.5 sm:pr-2.5", props.ClassName())
		props.SetClass(className)
		props.SetChildren(h.Span(h.Class("hidden sm:block"), cmp.Or(b.text, "Next")), lucide.ChevronRight())
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
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// styled is the bases PaginationLink / Previous / Next in style st.
func (b *Builder) styled(st style.Style, props *h.Element) *h.Element {
	switch b.kind {
	case kindPrevious:
		props.SetClass(stylePrevious.Class(st, props.ClassName()))
		props.SetChildren(
			lucide.ChevronLeft(h.Data("icon", "inline-start"), h.Class(styleFlip.Class(st))),
			h.Span(h.Class(stylePreviousText.Class(st)), cmp.Or(b.text, "Previous")),
		)
		return styledLink(st, b.active, b.activeSet, b.size, h.Spread(h.A(h.Aria("label", "Go to previous page")), props, nil))
	case kindNext:
		props.SetClass(styleNext.Class(st, props.ClassName()))
		props.SetChildren(
			h.Span(h.Class(styleNextText.Class(st)), cmp.Or(b.text, "Next")),
			lucide.ChevronRight(h.Data("icon", "inline-end"), h.Class(styleFlip.Class(st))),
		)
		return styledLink(st, b.active, b.activeSet, b.size, h.Spread(h.A(h.Aria("label", "Go to next page")), props, nil))
	}
	return styledLink(st, b.active, b.activeSet, b.size, props)
}

// styledLink is <Button asChild variant size className={cn("cn-pagination-link",
// className)}><a aria-current data-slot data-active {...props} /></Button>.
func styledLink(st style.Style, active, activeSet bool, size button.Size, props *h.Element) *h.Element {
	variant := button.VariantGhost
	if active {
		variant = button.VariantOutline
	}
	cls := lib.Cn(buttonVariants.Classes(st, cva.Selection{"variant": string(variant), "size": string(size)}, styleLink.Class(st, props.ClassName())))
	base := h.A(h.Data("slot", "pagination-link"), h.Data("variant", string(variant)), h.Data("size", string(size)))
	if active {
		base.Aria("current", "page").Data("active", "true")
	} else if activeSet {
		base.Data("active", "false")
	}
	props = props.Clone()
	props.SetClass("")
	return h.Spread(base, props, nil).SetClass(cls)
}

// Ellipsis stands for skipped pages. Like upstream, it always renders its
// icon and "More pages" label: child nodes are ignored, attributes apply.
func Ellipsis(children ...any) *part.Part {
	return part.New("span", func(st style.Style, props *h.Element) *h.Element {
		props = props.Clone().SetChildren()
		if st != "" {
			base := h.Span(h.Aria("hidden", "true"), h.Data("slot", "pagination-ellipsis"), h.Class(styleEllipsis.Class(st, "")),
				lucide.Ellipsis(),
				h.Span(h.Class("sr-only"), "More pages"),
			)
			return h.Spread(base, props, lib.Cn)
		}
		return ellipsis(props)
	}, children...)
}

func ellipsis(props *h.Element) *h.Element {
	base := h.Span(h.Aria("hidden", "true"), h.Data("slot", "pagination-ellipsis"), h.Class(lib.Cn("flex size-9 items-center justify-center")),
		lucide.Ellipsis(h.Class("size-4")),
		h.Span(h.Class("sr-only"), "More pages"),
	)
	return h.Spread(base, props, lib.Cn)
}
