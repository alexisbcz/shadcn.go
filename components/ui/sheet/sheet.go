// Package sheet is the shadcn/ui Sheet: a Dialog sliding in from an edge
// of the screen, complementing the main content.
//
//	sheet.Sheet(
//		sheet.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
//		sheet.Content(sheet.SideLeft,
//			sheet.Header(sheet.Title("Edit profile"), sheet.Description("...")),
//			sheet.Footer(sheet.Close(lib.AsChild, button.Button(button.VariantOutline, "Close"))),
//		),
//	)
//
// Built on the Radix Dialog primitive (package dialog/primitive).
package sheet

import (
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog/primitive"
	h "shadcn.go/libhtml"
)

// Root options.
var (
	// DefaultOpen opens the sheet when the page loads.
	DefaultOpen = primitive.DefaultOpen
)

// Modal sets whether the sheet is modal (default true).
func Modal(on bool) primitive.Option { return primitive.Modal(on) }

// Side selects the edge the sheet slides in from.
type Side string

// Sides.
const (
	SideTop    Side = "top"
	SideRight  Side = "right"
	SideBottom Side = "bottom"
	SideLeft   Side = "left"
)

// ShowCloseButton, passed to Content (default true), toggles the corner
// close button.
type ShowCloseButton bool

const (
	contentBase = "fixed z-50 flex flex-col gap-4 bg-background shadow-lg transition ease-in-out data-[state=closed]:animate-out data-[state=closed]:duration-300 data-[state=open]:animate-in data-[state=open]:duration-500"
	// The close button's className is a literal upstream (no cn).
	closeClass = "absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-secondary"
)

var (
	overlayClass = lib.Cn("fixed inset-0 z-50 bg-black/50 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:animate-in data-[state=open]:fade-in-0")
	sideClasses  = map[Side]string{
		SideRight:  "inset-y-0 right-0 h-full w-3/4 border-l data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right sm:max-w-sm",
		SideLeft:   "inset-y-0 left-0 h-full w-3/4 border-r data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left sm:max-w-sm",
		SideTop:    "inset-x-0 top-0 h-auto border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top",
		SideBottom: "inset-x-0 bottom-0 h-auto border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom",
	}
	headerClass = lib.Cn("flex flex-col gap-1.5 p-4")
	footerClass = lib.Cn("mt-auto flex flex-col gap-2 p-4")
	titleClass  = lib.Cn("font-semibold text-foreground")
	descClass   = lib.Cn("text-sm text-muted-foreground")
)

// The classes of registry/bases/radix/ui/sheet.tsx, used when a style is
// in scope (style.Scope).
var (
	baseFile        = style.NewFile("ui/sheet.tsx")
	baseOverlay     = baseFile.Cn("cn-sheet-overlay fixed inset-0 z-50 duration-100 data-open:animate-in data-open:fade-in-0 data-closed:animate-out data-closed:fade-out-0", style.Arg)
	baseContent     = baseFile.Cn("cn-sheet-content data-open:animate-in data-open:fade-in-0 data-[side=bottom]:data-open:slide-in-from-bottom-10 data-[side=left]:data-open:slide-in-from-left-10 data-[side=right]:data-open:slide-in-from-right-10 data-[side=top]:data-open:slide-in-from-top-10 data-closed:animate-out data-closed:fade-out-0 data-[side=bottom]:data-closed:slide-out-to-bottom-10 data-[side=left]:data-closed:slide-out-to-left-10 data-[side=right]:data-closed:slide-out-to-right-10 data-[side=top]:data-closed:slide-out-to-top-10", style.Arg)
	baseClose       = baseFile.Attr("cn-sheet-close")
	baseHeader      = baseFile.Cn("cn-sheet-header flex flex-col", style.Arg)
	baseFooter      = baseFile.Cn("cn-sheet-footer mt-auto flex flex-col", style.Arg)
	baseTitle       = baseFile.Cn("cn-sheet-title cn-font-heading", style.Arg)
	baseDescription = baseFile.Cn("cn-sheet-description", style.Arg)
)

func styleOf(w io.Writer) style.Style {
	st, _ := style.FromWriter(w)
	return st
}

func classFor(st style.Style, def string, site *style.Site) string {
	if st != "" {
		return site.Class(st, "")
	}
	return def
}

// Sheet is the root. It renders no element of its own.
func Sheet(children ...any) *primitive.RootBuilder { return primitive.Root(children...) }

// Trigger is the button that opens the sheet. Accepts lib.AsChild.
func Trigger(children ...any) *primitive.TriggerBuilder {
	return primitive.Trigger(append([]any{h.Data("slot", "sheet-trigger")}, children...)...)
}

// Close is a button that closes the sheet. Accepts lib.AsChild.
func Close(children ...any) *primitive.CloseBuilder {
	return primitive.Close(append([]any{h.Data("slot", "sheet-close")}, children...)...)
}

// Portal mounts its children at the end of <body> while the sheet is open.
// (Not exported upstream; provided for custom compositions.)
func Portal(children ...any) *primitive.PortalBuilder { return primitive.Portal(children...) }

// --- Overlay -------------------------------------------------------------------------------

// OverlayBuilder is returned by Overlay.
type OverlayBuilder struct {
	h.Builder[OverlayBuilder]
	primitive.Scoped
}

// Overlay is the translucent layer behind the content. (Not exported
// upstream; Content includes it.)
func Overlay(children ...any) *OverlayBuilder {
	b := &OverlayBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes *primitive.Scope.
func (b *OverlayBuilder) AcceptChild(c any) bool { return acceptScope(&b.Scoped, c) }

// Element resolves the overlay.
func (b *OverlayBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *OverlayBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *OverlayBuilder) element(st style.Style) *h.Element {
	props := h.Spread(h.Div(h.Data("slot", "sheet-overlay"), h.Class(classFor(st, overlayClass, baseOverlay))), b.Elem(), lib.Cn)
	return primitive.Overlay(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the overlay.
func (b *OverlayBuilder) Render(w io.Writer) error { return render(b.element(styleOf(w)), w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	primitive.Scoped
	side      Side
	showClose bool
}

// Content is the sheet panel, portalled with its Overlay. Accepts a Side
// (default SideRight) and ShowCloseButton.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{side: SideRight, showClose: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Side, ShowCloseButton and *primitive.Scope.
func (b *ContentBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Side:
		b.Side(v)
	case ShowCloseButton:
		b.showClose = bool(v)
	default:
		return acceptScope(&b.Scoped, c)
	}
	return true
}

// Side sets the edge. An empty value selects SideRight.
func (b *ContentBuilder) Side(s Side) *ContentBuilder {
	if s == "" {
		s = SideRight
	}
	b.side = s
	return b
}

// ShowCloseButton sets whether the corner close button is shown.
func (b *ContentBuilder) ShowCloseButton(on bool) *ContentBuilder { b.showClose = on; return b }

// Classes returns the content's classes for a side plus className, like
// the cn(...) call of SheetContent.
func Classes(side Side, className ...string) string {
	return lib.Cn(append([]string{contentBase, sideClasses[side]}, className...)...)
}

// Element resolves the portal holding the overlay and the panel.
func (b *ContentBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *ContentBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *ContentBuilder) element(st style.Style) *h.Element {
	sc := b.Scope()
	if st != "" {
		props := h.Spread(h.Div(h.Data("slot", "sheet-content"), h.Data("side", string(b.side))), b.Elem(), nil)
		props.SetClass(baseContent.Class(st, b.Elem().ClassName()))
		if b.showClose {
			props.Children(primitive.Close(sc, h.Data("slot", "sheet-close"), lib.AsChild,
				button.Button(button.VariantGhost, button.SizeIconSm, h.Class(baseClose.Class(st)),
					lucide.X(),
					h.Span(h.Class("sr-only"), "Close"),
				)))
		}
		return primitive.Portal(sc,
			Overlay(sc),
			primitive.Content(sc, primitive.WithProps(props)),
		).Element()
	}
	props := h.Spread(h.Div(h.Data("slot", "sheet-content")), b.Elem(), nil)
	props.SetClass(Classes(b.side, b.Elem().ClassName()))
	if b.showClose {
		props.Children(primitive.Close(sc, h.Class(closeClass),
			lucide.X(h.Class("size-4")),
			h.Span(h.Class("sr-only"), "Close"),
		))
	}
	return primitive.Portal(sc,
		Overlay(sc),
		primitive.Content(sc, primitive.WithProps(props)),
	).Element()
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// --- Header, Footer ----------------------------------------------------------------------------

// Header lays out the title and description.
func Header(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "sheet-header"), h.Class(classFor(styleOf(w), headerClass, baseHeader))), h.Div(children...), lib.Cn).Render(w)
	})
}

// Footer lays out the actions at the bottom of the sheet.
func Footer(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "sheet-footer"), h.Class(classFor(styleOf(w), footerClass, baseFooter))), h.Div(children...), lib.Cn).Render(w)
	})
}

// --- Title, Description ------------------------------------------------------------------------

// TitleBuilder is returned by Title.
type TitleBuilder struct {
	h.Builder[TitleBuilder]
	primitive.Scoped
}

// Title is the sheet's accessible name (<h2>).
func Title(children ...any) *TitleBuilder {
	b := &TitleBuilder{}
	b.Builder = h.NewBuilder(b, "h2")
	return b.Children(children...)
}

// AcceptChild consumes *primitive.Scope.
func (b *TitleBuilder) AcceptChild(c any) bool { return acceptScope(&b.Scoped, c) }

// Element resolves the title.
func (b *TitleBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *TitleBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *TitleBuilder) element(st style.Style) *h.Element {
	props := h.Spread(h.H2(h.Data("slot", "sheet-title"), h.Class(classFor(st, titleClass, baseTitle))), b.Elem(), lib.Cn)
	return primitive.Title(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the title.
func (b *TitleBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// DescriptionBuilder is returned by Description.
type DescriptionBuilder struct {
	h.Builder[DescriptionBuilder]
	primitive.Scoped
}

// Description is the sheet's accessible description (<p>).
func Description(children ...any) *DescriptionBuilder {
	b := &DescriptionBuilder{}
	b.Builder = h.NewBuilder(b, "p")
	return b.Children(children...)
}

// AcceptChild consumes *primitive.Scope.
func (b *DescriptionBuilder) AcceptChild(c any) bool { return acceptScope(&b.Scoped, c) }

// Element resolves the description.
func (b *DescriptionBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *DescriptionBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *DescriptionBuilder) element(st style.Style) *h.Element {
	props := h.Spread(h.P(h.Data("slot", "sheet-description"), h.Class(classFor(st, descClass, baseDescription))), b.Elem(), lib.Cn)
	return primitive.Description(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the description.
func (b *DescriptionBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// --- helpers -----------------------------------------------------------------------------------

func acceptScope(s *primitive.Scoped, c any) bool {
	if sc, ok := c.(*primitive.Scope); ok {
		s.SetDialogScope(sc)
		return true
	}
	return false
}

func render(el *h.Element, w io.Writer) error {
	if el == nil {
		return nil
	}
	return el.Render(w)
}
