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
func (b *OverlayBuilder) Element() *h.Element {
	props := h.Spread(h.Div(h.Data("slot", "sheet-overlay"), h.Class(overlayClass)), b.Elem(), lib.Cn)
	return primitive.Overlay(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the overlay.
func (b *OverlayBuilder) Render(w io.Writer) error { return render(b.Element(), w) }

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
func (b *ContentBuilder) Element() *h.Element {
	sc := b.Scope()
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
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Header, Footer ----------------------------------------------------------------------------

// Header lays out the title and description.
func Header(children ...any) *h.Element {
	return h.Spread(h.Div(h.Data("slot", "sheet-header"), h.Class(headerClass)), h.Div(children...), lib.Cn)
}

// Footer lays out the actions at the bottom of the sheet.
func Footer(children ...any) *h.Element {
	return h.Spread(h.Div(h.Data("slot", "sheet-footer"), h.Class(footerClass)), h.Div(children...), lib.Cn)
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
func (b *TitleBuilder) Element() *h.Element {
	props := h.Spread(h.H2(h.Data("slot", "sheet-title"), h.Class(titleClass)), b.Elem(), lib.Cn)
	return primitive.Title(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the title.
func (b *TitleBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

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
func (b *DescriptionBuilder) Element() *h.Element {
	props := h.Spread(h.P(h.Data("slot", "sheet-description"), h.Class(descClass)), b.Elem(), lib.Cn)
	return primitive.Description(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the description.
func (b *DescriptionBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

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
