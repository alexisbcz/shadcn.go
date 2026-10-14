// Package scrollarea is the shadcn/ui ScrollArea (Radix ScrollArea):
// augments native scroll functionality for custom, cross-browser styling.
//
//	scrollarea.ScrollArea(h.Class("h-72 w-48 rounded-md border"), content...)
//	scrollarea.ScrollArea(h.Class("w-96 whitespace-nowrap"), content,
//	    scrollarea.ScrollBar(scrollarea.Horizontal))
//
// ScrollArea always has a vertical scrollbar; add a horizontal ScrollBar
// among its children for horizontal scrolling, as upstream. Root children
// may include a Type (when the scrollbars show: hover, the default, scroll,
// auto or always), ScrollHideDelay (ms before hiding, default 600) and
// Nonce.
//
// Behavior (scroll-area.js): the native scrollbars of the viewport are
// hidden and replaced with the Radix scrollbars, mounted (and unmounted)
// per Type: while the pointer is over the area (hover), while scrolling
// (scroll), while the content overflows (auto) or always; each has a thumb
// sized and positioned from the viewport's scroll position, dragged with
// the pointer, and the track scrolls on click and wheel. When both
// scrollbars show, a corner fills the gap.
//
// Radix renders a <style> hiding the native scrollbars inside the area.
// Under a nonce-based Content-Security-Policy a <style> only applies with
// the request's nonce: pass it with Nonce(secure.Nonce(r)) to render the
// <style> server-side; without it the behavior inserts the <style> (with
// the page's script nonce) when it mounts.
package scrollarea

import (
	_ "embed"
	"fmt"
	"io"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/part"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed scroll-area.js
var script string

func init() { runtime.Module("scroll-area", script) }

// Type selects when the scrollbars are visible.
type Type string

// Types, as in Radix: hover (default) shows the scrollbars while the
// pointer is over the area and the content overflows; scroll while
// scrolling; auto while the content overflows; always always.
const (
	TypeHover  Type = "hover"
	TypeScroll Type = "scroll"
	TypeAuto   Type = "auto"
	TypeAlways Type = "always"
)

// ScrollHideDelay is the delay in milliseconds before hover and scroll
// scrollbars hide (default 600).
type ScrollHideDelay int

// Nonce is the request's CSP nonce, for the <style> Radix renders.
type Nonce string

// Orientation of a ScrollBar.
type Orientation string

// Orientations.
const (
	Vertical   Orientation = "vertical"
	Horizontal Orientation = "horizontal"
)

// ViewportCSS is the stylesheet Radix renders with every viewport.
const ViewportCSS = `[data-radix-scroll-area-viewport]{scrollbar-width:none;-ms-overflow-style:none;-webkit-overflow-scrolling:touch;}[data-radix-scroll-area-viewport]::-webkit-scrollbar{display:none}`

// Builder is returned by ScrollArea.
type Builder struct {
	h.Builder[Builder]
	typ       Type
	hideDelay int
	nonce     string
}

// ScrollArea returns the scroll area: a root <div> holding the viewport
// (with the children), the vertical scrollbar and the corner.
func ScrollArea(children ...any) *Builder {
	b := &Builder{typ: TypeHover, hideDelay: 600}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Type, ScrollHideDelay and Nonce children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Type:
		b.typ = v
	case ScrollHideDelay:
		b.hideDelay = int(v)
	case Nonce:
		b.nonce = string(v)
	default:
		return false
	}
	return true
}

// Type sets when the scrollbars are visible.
func (b *Builder) Type(t Type) *Builder { b.typ = t; return b }

// ScrollHideDelay sets the hide delay in milliseconds.
func (b *Builder) ScrollHideDelay(ms int) *Builder { b.hideDelay = ms; return b }

// Nonce sets the CSP nonce of the <style> element.
func (b *Builder) Nonce(n string) *Builder { b.nonce = n; return b }

// styleFile is registry/bases/radix/ui/scroll-area.tsx (style.Scope).
var styleFile = style.NewFile("ui/scroll-area.tsx")

var (
	styleRoot      = styleFile.Cn("cn-scroll-area relative", style.Arg)
	styleViewport  = styleFile.Attr("cn-scroll-area-viewport size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1")
	styleScrollbar = styleFile.Cn("cn-scroll-area-scrollbar flex touch-none p-px transition-colors select-none", style.Arg)
	styleThumb     = styleFile.Attr("cn-scroll-area-thumb relative flex-1 bg-border")
)

const viewportClass = "size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1"

// Element resolves the scroll area to its root <div>.
func (b *Builder) Element() *h.Element { return b.element("") }

// element renders in style st ("" for new-york-v4).
func (b *Builder) element(st style.Style) *h.Element {
	props := b.Elem()
	kids := props.ChildNodes()
	horizontal := false
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *ScrollBarBuilder:
			if c.orientation == Horizontal {
				horizontal = true
			}
			return false
		case *Builder:
			return false // a nested area has its own scrollbars
		}
		return true
	})

	typ := b.typ
	switch typ {
	case TypeHover, TypeScroll, TypeAuto, TypeAlways:
	default:
		typ = TypeHover
	}
	base := h.Div(h.Dir("ltr"), h.Data("slot", "scroll-area"))
	// Radix props with no DOM trace, read (and removed) by scroll-area.js.
	if typ != TypeHover {
		base.Data("type", string(typ))
	}
	if b.hideDelay != 600 {
		base.Data("scroll-hide-delay", strconv.Itoa(b.hideDelay))
	}
	css := "position:relative;--radix-scroll-area-corner-width:0px;--radix-scroll-area-corner-height:0px"
	if us, ok := props.GetAttr("style"); ok && us != "" {
		css += ";" + us
	}
	user := props.Clone().SetChildren()
	user.Attr("style", css)
	root := h.Spread(base, user, nil).SetClass(lib.Cn("relative", props.ClassName()))
	vpClass := viewportClass
	if st != "" {
		root.SetClass(styleRoot.Class(st, props.ClassName()))
		vpClass = styleViewport.Class(st)
	}

	if b.nonce != "" {
		root.Children(h.StyleEl(h.Nonce(b.nonce), h.UnsafeCSS(ViewportCSS)))
	}
	overflowX := "hidden"
	if horizontal {
		overflowX = "scroll"
	}
	viewport := h.Div(
		h.Data("radix-scroll-area-viewport", ""),
		h.Data("slot", "scroll-area-viewport"),
		h.Class(vpClass),
		h.Style("overflow-x:"+overflowX+";overflow-y:scroll"),
		h.Div(h.Style("min-width:100%;display:table"), kids),
	)
	return root.Children(viewport, ScrollBar())
}

// Render renders the scroll area.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return b.element(st)
}

// ScrollBarBuilder is returned by ScrollBar.
type ScrollBarBuilder struct {
	h.Builder[ScrollBarBuilder]
	rtl         bool // inside style.ScopeRTL, set by ElementIn
	orientation Orientation
}

// ScrollBar returns a scrollbar of its ScrollArea (vertical unless given
// Horizontal). Radix mounts it only while visible: it renders as a
// <template> the behavior mounts it from.
func ScrollBar(children ...any) *ScrollBarBuilder {
	b := &ScrollBarBuilder{orientation: Vertical}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Orientation children.
func (b *ScrollBarBuilder) AcceptChild(c any) bool {
	if o, ok := c.(Orientation); ok {
		b.Orientation(o)
		return true
	}
	return false
}

// Orientation sets the scrollbar's orientation.
func (b *ScrollBarBuilder) Orientation(o Orientation) *ScrollBarBuilder {
	if o != Horizontal {
		o = Vertical
	}
	b.orientation = o
	return b
}

// Element resolves the scrollbar to its <template>.
func (b *ScrollBarBuilder) Element() *h.Element { return b.element("") }

// element renders in style st ("" for new-york-v4).
func (b *ScrollBarBuilder) element(st style.Style) *h.Element {
	o := b.orientation
	if o != Horizontal {
		o = Vertical
	}
	props := b.Elem()
	if len(props.ChildNodes()) > 0 {
		return h.Errored(fmt.Errorf("scrollarea: ScrollBar takes no content"))
	}
	class := lib.Cn(
		"flex touch-none p-px transition-colors select-none",
		h.IfElse(o == Vertical, "h-full w-2.5 border-l border-l-transparent", ""),
		h.IfElse(o == Horizontal, "h-2.5 flex-col border-t border-t-transparent", ""),
		props.ClassName(),
	)
	base := h.Div(h.Data("orientation", string(o)), h.Data("slot", "scroll-area-scrollbar"))
	bar := h.Spread(base, props, nil).SetClass(class)
	thumbClass := "relative flex-1 rounded-full bg-border"
	if st != "" {
		bar.SetClass(part.Class(b.rtl, styleScrollbar, st, props.ClassName()))
		thumbClass = styleThumb.Class(st)
	}
	bar.Children(h.Div(h.Data("slot", "scroll-area-thumb"), h.Class(thumbClass)))
	return h.Template(h.Data("slot", "scroll-area-scrollbar-template"), bar)
}

// Render renders the scrollbar's template.
func (b *ScrollBarBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *ScrollBarBuilder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	b.rtl = part.RTL(w)
	return b.element(st)
}
