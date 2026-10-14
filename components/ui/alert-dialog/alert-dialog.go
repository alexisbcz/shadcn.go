// Package alertdialog is the shadcn/ui AlertDialog: a modal dialog that
// interrupts the user with important content and expects a response.
//
//	alertdialog.AlertDialog(
//		alertdialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Show Dialog")),
//		alertdialog.Content(
//			alertdialog.Header(alertdialog.Title("Are you absolutely sure?"), alertdialog.Description("...")),
//			alertdialog.Footer(alertdialog.Cancel("Cancel"), alertdialog.Action("Continue")),
//		),
//	)
//
// Unlike Dialog it is always modal, is not dismissed by clicking outside,
// and focuses Cancel when it opens. Built on the Radix Dialog primitive
// (package dialog/primitive) with AlertDialog semantics.
package alertdialog

import (
	"io"
	"slices"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog/primitive"
	"shadcn.go/components/ui/internal/rtl"
	h "shadcn.go/libhtml"
)

// DefaultOpen, passed to AlertDialog, opens it when the page loads.
var DefaultOpen = primitive.DefaultOpen

// Size selects the content's size.
type Size string

// Content sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
)

var (
	overlayClass = lib.Cn("fixed inset-0 z-50 bg-black/50 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:animate-in data-[state=open]:fade-in-0")
	contentClass = lib.Cn("group/alert-dialog-content fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border bg-background p-6 shadow-lg duration-200 data-[size=sm]:max-w-xs data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[size=default]:sm:max-w-lg")
	headerClass  = lib.Cn("grid grid-rows-[auto_1fr] place-items-center gap-1.5 text-center has-data-[slot=alert-dialog-media]:grid-rows-[auto_auto_1fr] has-data-[slot=alert-dialog-media]:gap-x-6 sm:group-data-[size=default]/alert-dialog-content:place-items-start sm:group-data-[size=default]/alert-dialog-content:text-left sm:group-data-[size=default]/alert-dialog-content:has-data-[slot=alert-dialog-media]:grid-rows-[auto_1fr]")
	footerClass  = lib.Cn("flex flex-col-reverse gap-2 group-data-[size=sm]/alert-dialog-content:grid group-data-[size=sm]/alert-dialog-content:grid-cols-2 sm:flex-row sm:justify-end")
	titleClass   = lib.Cn("text-lg font-semibold sm:group-data-[size=default]/alert-dialog-content:group-has-data-[slot=alert-dialog-media]/alert-dialog-content:col-start-2")
	descClass    = lib.Cn("text-sm text-muted-foreground")
	mediaClass   = lib.Cn("mb-2 inline-flex size-16 items-center justify-center rounded-md bg-muted sm:group-data-[size=default]/alert-dialog-content:row-span-2 *:[svg:not([class*='size-'])]:size-8")
)

// The classes of registry/bases/radix/ui/alert-dialog.tsx, used when a
// style is in scope (style.Scope).
var (
	baseFile        = style.NewFile("ui/alert-dialog.tsx")
	baseOverlay     = baseFile.Cn("cn-alert-dialog-overlay fixed inset-0 z-50", style.Arg)
	baseContent     = baseFile.Cn("cn-alert-dialog-content group/alert-dialog-content fixed top-1/2 left-1/2 z-50 grid w-full -translate-x-1/2 -translate-y-1/2 outline-none", style.Arg)
	baseHeader      = baseFile.Cn("cn-alert-dialog-header", style.Arg)
	baseFooter      = baseFile.Cn("cn-alert-dialog-footer flex flex-col-reverse gap-2 group-data-[size=sm]/alert-dialog-content:grid group-data-[size=sm]/alert-dialog-content:grid-cols-2 sm:flex-row sm:justify-end", style.Arg)
	baseMedia       = baseFile.Cn("cn-alert-dialog-media", style.Arg)
	baseTitle       = baseFile.Cn("cn-alert-dialog-title cn-font-heading", style.Arg)
	baseDescription = baseFile.Cn("cn-alert-dialog-description", style.Arg)
	baseAction      = baseFile.Cn("cn-alert-dialog-action", style.Arg)
	baseCancel      = baseFile.Cn("cn-alert-dialog-cancel", style.Arg)
)

// rtlFor maps the component's own classes to logical ones when the user
// gave the element dir="rtl" (styles/<style>/ui-rtl upstream).
func rtlFor(user *h.Element, class string) string {
	if d, ok := user.GetAttr("dir"); ok && d == "rtl" {
		return rtl.Classes(class)
	}
	return class
}

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

// AlertDialog is the root. It renders no element of its own.
func AlertDialog(children ...any) *primitive.RootBuilder {
	return primitive.Root(append([]any{primitive.Alert}, children...)...)
}

// Trigger is the button that opens the dialog. Accepts lib.AsChild.
func Trigger(children ...any) *primitive.TriggerBuilder {
	return primitive.Trigger(append([]any{h.Data("slot", "alert-dialog-trigger")}, children...)...)
}

// Portal mounts its children at the end of <body> while the dialog is open.
func Portal(children ...any) *primitive.PortalBuilder { return primitive.Portal(children...) }

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
	props := h.Spread(h.Div(h.Data("slot", "alert-dialog-overlay"), h.Class(classFor(st, overlayClass, baseOverlay))), b.Elem(), lib.Cn)
	return primitive.Overlay(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the overlay.
func (b *OverlayBuilder) Render(w io.Writer) error { return render(b.element(styleOf(w)), w) }

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	primitive.Scoped
	size Size
}

// Content is the dialog window, portalled with its Overlay. Accepts a Size.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{size: SizeDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Size and *primitive.Scope.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Size); ok {
		b.Size(v)
		return true
	}
	return acceptScope(&b.Scoped, c)
}

// Size sets the size. An empty value selects the default.
func (b *ContentBuilder) Size(s Size) *ContentBuilder {
	if s == "" {
		s = SizeDefault
	}
	b.size = s
	return b
}

// Element resolves the portal holding the overlay and the content.
func (b *ContentBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *ContentBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *ContentBuilder) element(st style.Style) *h.Element {
	sc := b.Scope()
	props := h.Spread(h.Div(h.Data("slot", "alert-dialog-content"), h.Data("size", string(b.size)), h.Class(classFor(st, contentClass, baseContent))), b.Elem(), lib.Cn)
	if st != "" {
		props.SetClass(lib.Cn(rtlFor(b.Elem(), baseContent.Class(st, "")), b.Elem().ClassName()))
	}
	return primitive.Portal(sc,
		Overlay(sc),
		primitive.Content(sc, primitive.WithProps(props)),
	).Element()
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// --- Header, Footer, Media ---------------------------------------------------------------------

// Header lays out the media, title and description.
func Header(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "alert-dialog-header"), h.Class(classFor(styleOf(w), headerClass, baseHeader))), h.Div(children...), lib.Cn).Render(w)
	})
}

// Footer lays out the actions.
func Footer(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "alert-dialog-footer"), h.Class(classFor(styleOf(w), footerClass, baseFooter))), h.Div(children...), lib.Cn).Render(w)
	})
}

// Media holds an icon or image above (or beside) the title.
func Media(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		return h.Spread(h.Div(h.Data("slot", "alert-dialog-media"), h.Class(classFor(styleOf(w), mediaClass, baseMedia))), h.Div(children...), lib.Cn).Render(w)
	})
}

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
	props := h.Spread(h.H2(h.Data("slot", "alert-dialog-title"), h.Class(classFor(st, titleClass, baseTitle))), b.Elem(), lib.Cn)
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
	props := h.Spread(h.P(h.Data("slot", "alert-dialog-description"), h.Class(classFor(st, descClass, baseDescription))), b.Elem(), lib.Cn)
	return primitive.Description(b.Inherit(), primitive.WithProps(props)).Element()
}

// Render renders the description.
func (b *DescriptionBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// --- Action, Cancel ----------------------------------------------------------------------------

// ButtonBuilder is returned by Action and Cancel.
type ButtonBuilder struct {
	h.Builder[ButtonBuilder]
	primitive.Scoped
	variant button.Variant
	size    button.Size
	cancel  bool
}

// Action is the button confirming the action; it closes the dialog.
// Accepts a button.Variant (default "default") and button.Size.
func Action(children ...any) *ButtonBuilder {
	b := &ButtonBuilder{variant: button.VariantDefault, size: button.SizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// Cancel is the button dismissing the dialog; it has focus when the dialog
// opens. Accepts a button.Variant (default "outline") and button.Size.
func Cancel(children ...any) *ButtonBuilder {
	b := &ButtonBuilder{variant: button.VariantOutline, size: button.SizeDefault, cancel: true}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes button.Variant, button.Size and *primitive.Scope.
func (b *ButtonBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case button.Variant:
		b.Variant(v)
	case button.Size:
		b.Size(v)
	default:
		return acceptScope(&b.Scoped, c)
	}
	return true
}

// Variant sets the button variant.
func (b *ButtonBuilder) Variant(v button.Variant) *ButtonBuilder {
	if v == "" {
		v = button.VariantDefault
	}
	b.variant = v
	return b
}

// Size sets the button size.
func (b *ButtonBuilder) Size(s button.Size) *ButtonBuilder {
	if s == "" {
		s = button.SizeDefault
	}
	b.size = s
	return b
}

// Element resolves the button:
// <Button variant size asChild><AlertDialogPrimitive.Action|Cancel data-slot className={cn(className)} {...props} /></Button>
func (b *ButtonBuilder) Element() *h.Element { return b.element("") }

// ElementIn resolves the element in the style in scope when rendering to w.
func (b *ButtonBuilder) ElementIn(w io.Writer) *h.Element { return b.element(styleOf(w)) }

func (b *ButtonBuilder) element(st style.Style) *h.Element {
	slot, part, site := "alert-dialog-action", primitive.Close, baseAction
	if b.cancel {
		slot, part, site = "alert-dialog-cancel", primitive.Cancel, baseCancel
	}
	props := h.Spread(h.Button(h.Data("slot", slot)), b.Elem(), lib.Cn)
	if st != "" {
		props.SetClass(site.Class(st, props.ClassName()))
	} else if props.ClassName() != "" {
		props.SetClass(lib.Cn(props.ClassName()))
	}
	child := part(b.Scope(), primitive.WithProps(props)).Element()
	// <Button asChild> renders a Radix Slot: the Button's props are merged
	// into the child (the child's attributes win) and the classes are
	// concatenated, Button's first, without cn.
	defaults := h.Button(
		h.Data("slot", "button"),
		h.Data("variant", string(b.variant)),
		h.Data("size", string(b.size)),
		h.Class(buttonClasses(st, b.variant, b.size)),
	)
	out, err := h.Slot(defaults, child, joinClasses)
	if err != nil {
		return h.Errored(err)
	}
	return out
}

// joinClasses is Radix Slot's className merge: [slot, child].filter(Boolean).join(" ").
func joinClasses(classes ...string) string {
	return strings.Join(slices.DeleteFunc(classes, func(c string) bool { return c == "" }), " ")
}

// Render renders the button.
func (b *ButtonBuilder) Render(w io.Writer) error { return b.element(styleOf(w)).Render(w) }

// buttonClasses is buttonVariants({ variant, size }) in style st.
func buttonClasses(st style.Style, v button.Variant, s button.Size) string {
	if st != "" {
		return button.StyleVariants(st, v, s)
	}
	return button.Variants(v, s)
}

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
