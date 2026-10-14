// Package dialog is the shadcn/ui Dialog: a window overlaid on the primary
// window, rendering the content underneath inert.
//
//	dialog.Dialog(
//		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
//		dialog.Content(
//			dialog.Header(dialog.Title("Edit profile"), dialog.Description("...")),
//			...,
//			dialog.Footer(dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel"))),
//		),
//	)
//
// Options are passed as children of Dialog: dialog.Modal(false),
// dialog.DefaultOpen. Behavior and markup come from the Radix Dialog
// primitive (see package primitive).
package dialog

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
	// DefaultOpen opens the dialog when the page loads.
	DefaultOpen = primitive.DefaultOpen
)

// Modal sets whether the dialog is modal (default true).
func Modal(on bool) primitive.Option { return primitive.Modal(on) }

// ShowCloseButton, passed to Content (default true) or Footer (default
// false), adds a close button.
type ShowCloseButton bool

var (
	overlayClass = lib.Cn("fixed inset-0 z-50 bg-black/50 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:animate-in data-[state=open]:fade-in-0")
	contentClass = lib.Cn("fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border bg-background p-6 shadow-lg duration-200 outline-none data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 sm:max-w-lg")
	// The close button's className is a literal upstream (no cn).
	closeClass  = ("absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4")
	headerClass = lib.Cn("flex flex-col gap-2 text-center sm:text-left")
	footerClass = lib.Cn("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end")
	titleClass  = lib.Cn("text-lg leading-none font-semibold")
	descClass   = lib.Cn("text-sm text-muted-foreground")
)

// The classes of registry/bases/radix/ui/dialog.tsx, used when a style is
// in scope (style.Scope).
var (
	baseFile        = style.NewFile("ui/dialog.tsx")
	baseOverlay     = baseFile.Cn("cn-dialog-overlay fixed inset-0 isolate z-50", style.Arg)
	baseContent     = baseFile.Cn("cn-dialog-content fixed top-1/2 left-1/2 z-50 w-full -translate-x-1/2 -translate-y-1/2 outline-none", style.Arg)
	baseClose       = baseFile.Attr("cn-dialog-close")
	baseHeader      = baseFile.Cn("cn-dialog-header flex flex-col", style.Arg)
	baseFooter      = baseFile.Cn("cn-dialog-footer flex flex-col-reverse gap-2 sm:flex-row sm:justify-end", style.Arg)
	baseTitle       = baseFile.Cn("cn-dialog-title cn-font-heading", style.Arg)
	baseDescription = baseFile.Cn("cn-dialog-description", style.Arg)
)

// styleOf returns the style in scope when rendering to w ("" outside a
// style.Scope).
func styleOf(w io.Writer) style.Style {
	st, _ := style.FromWriter(w)
	return st
}

// classFor returns new-york-v4's class def, or the base site resolved in st.
func classFor(st style.Style, def string, site *style.Site) string {
	if st != "" {
		return site.Class(st, "")
	}
	return def
}

// Dialog is the root. It renders no element of its own.
func Dialog(children ...any) *primitive.RootBuilder { return primitive.Root(children...) }

// Trigger is the button that opens the dialog. Accepts lib.AsChild.
func Trigger(children ...any) *primitive.TriggerBuilder {
	return primitive.Trigger(append([]any{h.Data("slot", "dialog-trigger")}, children...)...)
}

// Portal mounts its children at the end of <body> while the dialog is open.
func Portal(children ...any) *primitive.PortalBuilder { return primitive.Portal(children...) }

// Close is a button that closes the dialog. Accepts lib.AsChild.
func Close(children ...any) *primitive.CloseBuilder {
	return primitive.Close(append([]any{h.Data("slot", "dialog-close")}, children...)...)
}

// --- Overlay -------------------------------------------------------------------------------

// OverlayBuilder is returned by Overlay.
type OverlayBuilder struct {
	h.Builder[OverlayBuilder]
	primitive.Scoped
}

// Overlay is the translucent layer behind the content.
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
	props := h.Spread(h.Div(h.Data("slot", "dialog-overlay"), h.Class(classFor(st, overlayClass, baseOverlay))), b.Elem(), lib.Cn)
	return primitive.Overlay(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the overlay.
func (b *OverlayBuilder) Render(w io.Writer) error { return render(b.element(styleOf(w)), w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	primitive.Scoped
	showClose bool
}

// Content is the dialog window, portalled with its Overlay. It has a close
// button in its top right corner unless ShowCloseButton(false) is passed.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{showClose: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes ShowCloseButton and *primitive.Scope.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if v, ok := c.(ShowCloseButton); ok {
		b.showClose = bool(v)
		return true
	}
	return acceptScope(&b.Scoped, c)
}

// ShowCloseButton sets whether the corner close button is shown.
func (b *ContentBuilder) ShowCloseButton(on bool) *ContentBuilder { b.showClose = on; return b }

// Element resolves the portal holding the overlay and the content.
func (b *ContentBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *ContentBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *ContentBuilder) element(st style.Style) *h.Element {
	sc := b.Scope()
	props := h.Spread(h.Div(h.Data("slot", "dialog-content"), h.Class(classFor(st, contentClass, baseContent))), b.Elem(), lib.Cn)
	if b.showClose && st != "" {
		// <DialogPrimitive.Close data-slot="dialog-close" asChild><Button variant="ghost" className="cn-dialog-close" size="icon-sm">
		props.Children(primitive.Close(sc, h.Data("slot", "dialog-close"), lib.AsChild,
			button.Button(button.VariantGhost, button.SizeIconSm, h.Class(baseClose.Class(st)),
				lucide.X(),
				h.Span(h.Class("sr-only"), "Close"),
			)))
	} else if b.showClose {
		props.Children(primitive.Close(sc, h.Data("slot", "dialog-close"), h.Class(closeClass),
			lucide.X(),
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

// --- Header, Footer --------------------------------------------------------------------------

// Header lays out the title and description.
func Header(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "dialog-header"), h.Class(classFor(styleOf(w), headerClass, baseHeader))), h.Div(children...), lib.Cn).Render(w)
	})
}

// FooterBuilder is returned by Footer.
type FooterBuilder struct {
	h.Builder[FooterBuilder]
	primitive.Scoped
	showClose bool
}

// Footer lays out the actions. ShowCloseButton(true) appends an outline
// "Close" button.
func Footer(children ...any) *FooterBuilder {
	b := &FooterBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes ShowCloseButton and *primitive.Scope.
func (b *FooterBuilder) AcceptChild(c any) bool {
	if v, ok := c.(ShowCloseButton); ok {
		b.showClose = bool(v)
		return true
	}
	return acceptScope(&b.Scoped, c)
}

// ShowCloseButton sets whether the "Close" button is appended.
func (b *FooterBuilder) ShowCloseButton(on bool) *FooterBuilder { b.showClose = on; return b }

// Element resolves the footer.
func (b *FooterBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *FooterBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *FooterBuilder) element(st style.Style) *h.Element {
	el := h.Spread(h.Div(h.Data("slot", "dialog-footer"), h.Class(classFor(st, footerClass, baseFooter))), b.Elem(), lib.Cn)
	if b.showClose {
		// <DialogPrimitive.Close asChild><Button variant="outline">Close</Button></DialogPrimitive.Close>
		el.Children(primitive.Close(b.Scope(), lib.AsChild, button.Button(button.VariantOutline, "Close")))
	}
	return el
}

// Render renders the footer.
func (b *FooterBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// --- Title, Description ------------------------------------------------------------------------

// TitleBuilder is returned by Title.
type TitleBuilder struct {
	h.Builder[TitleBuilder]
	primitive.Scoped
}

// Title is the dialog's accessible name (<h2>).
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
	props := h.Spread(h.H2(h.Data("slot", "dialog-title"), h.Class(classFor(st, titleClass, baseTitle))), b.Elem(), lib.Cn)
	return primitive.Title(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the title.
func (b *TitleBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// DescriptionBuilder is returned by Description.
type DescriptionBuilder struct {
	h.Builder[DescriptionBuilder]
	primitive.Scoped
}

// Description is the dialog's accessible description (<p>).
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
	props := h.Spread(h.P(h.Data("slot", "dialog-description"), h.Class(classFor(st, descClass, baseDescription))), b.Elem(), lib.Cn)
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
