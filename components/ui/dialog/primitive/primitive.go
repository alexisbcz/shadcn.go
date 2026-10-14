// Package primitive is Radix UI's Dialog primitive (@radix-ui/react-dialog
// 1.1.15, with the AlertDialog variant of @radix-ui/react-alert-dialog
// 1.1.15) for shadcn.go. The shadcn Dialog, AlertDialog and Sheet are built
// on it; use those unless you need the unstyled parts.
//
//	primitive.Root(
//		primitive.Trigger(lib.AsChild, button.Button("Open")),
//		primitive.Portal(
//			primitive.Overlay(h.Class("fixed inset-0 bg-black/50")),
//			primitive.Content(primitive.Title("Title"), primitive.Description("..."), primitive.Close("Close")),
//		),
//	)
//
// Markup is what Radix renders. Root renders no element; it shares ids and
// options with its parts (Radix's DialogContext) when it renders. Closed
// content is not in the DOM (Radix Presence + Portal): Portal renders an
// inert <template> that the behavior (dialog.js) clones into <body> when
// the dialog opens and removes after the exit animation. Content and
// Overlay outside a Portal get their own template and open in place.
//
// Behavior (dialog.js, built on the shared runtime primitives): trigger
// toggles; Close closes; Escape, pointer down outside (not for alert
// dialogs) and, for non-modal dialogs, focus outside dismiss; modal dialogs
// lock page scroll, hide the rest of the page from assistive technology,
// disable outside pointer events, trap focus (Tab loops) and add focus
// guards; focus moves into the content on open (the Cancel button for
// alert dialogs) and back to the trigger on close. The dialog dispatches
// "shadcn:open-change" (cancelable, detail.open) on its template; scripts
// can drive it with the module's getDialog (see dialog.js).
package primitive

import (
	_ "embed"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed dialog.js
var script string

func init() { runtime.Module("dialog", script) }

// Marker attributes the behavior reads (and removes) on the cloned content.
const (
	markerClose   = "data-radix-dialog-close"
	markerCancel  = "data-radix-dialog-cancel"
	markerOverlay = "data-radix-dialog-overlay"
)

// Scope is the state a Root shares with its parts (Radix's DialogContext).
type Scope struct {
	ContentID     string
	TitleID       string
	DescriptionID string
	Modal         bool
	DefaultOpen   bool
	// Alert selects AlertDialog semantics: role="alertdialog", always
	// modal, no dismissal on outside interaction, focus on Cancel.
	Alert bool
}

// NewScope returns a scope with fresh ids.
func NewScope() *Scope {
	return &Scope{ContentID: uid.New("radix"), TitleID: uid.New("radix"), DescriptionID: uid.New("radix"), Modal: true}
}

// Scoped is embedded by every component taking part in a dialog: the Root
// hands it the dialog's scope before rendering. Wrappers around the parts
// (the shadcn components) embed it too and pass Scope() to the parts they
// build.
type Scoped struct {
	scope    *Scope
	inPortal bool
}

// SetDialogScope is called by Root.
func (s *Scoped) SetDialogScope(sc *Scope) { s.scope = sc }

// SetInPortal is called by Portal on the components among its children:
// Content and Overlay inside a Portal are mounted by it, outside one they
// mount in place.
func (s *Scoped) SetInPortal() { s.inPortal = true }

// Inherited carries a wrapper's scope and portal placement to the parts it
// builds; see Scoped.Inherit.
type Inherited struct {
	scope    *Scope
	inPortal bool
}

// Inherit returns a child that hands s's scope and portal placement to a
// part:
//
//	primitive.Content(b.Inherit(), primitive.WithProps(props))
func (s *Scoped) Inherit() Inherited { return Inherited{s.Scope(), s.inPortal} }

// Scope returns the dialog's scope. A part rendered outside any Root gets a
// scope of its own.
func (s *Scoped) Scope() *Scope {
	if s.scope == nil {
		s.scope = NewScope()
	}
	return s.scope
}

// Participant is implemented by everything that embeds Scoped.
type Participant interface{ SetDialogScope(*Scope) }

// Option configures a Root; pass options as children.
type Option func(*Scope)

// Modal sets whether the dialog is modal (default true). Non-modal
// dialogs have no overlay, no scroll lock, do not hide or block the page,
// and close when focus or a pointer goes outside.
func Modal(on bool) Option { return func(s *Scope) { s.Modal = on } }

// DefaultOpen opens the dialog when the page loads.
var DefaultOpen Option = func(s *Scope) { s.DefaultOpen = true }

// Alert selects AlertDialog semantics (see Scope.Alert).
var Alert Option = func(s *Scope) { s.Alert = true; s.Modal = true }

// RootBuilder is returned by Root.
type RootBuilder struct {
	h.Builder[RootBuilder]
	opts []Option
}

// Root is Dialog.Root. It renders its children only, after sharing the
// dialog's scope with every part among them (through elements, fragments
// and component builders; not through h.Lazy or NodeFunc nodes). Nested
// Roots provide their own scope. Attributes passed to Root are ignored,
// like props React's Root does not render.
func Root(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes Options.
func (b *RootBuilder) AcceptChild(c any) bool {
	if o, ok := c.(Option); ok {
		b.opts = append(b.opts, o)
		return true
	}
	return false
}

// Options adds options.
func (b *RootBuilder) Options(opts ...Option) *RootBuilder {
	b.opts = append(b.opts, opts...)
	return b
}

// Render renders the children.
func (b *RootBuilder) Render(w io.Writer) error {
	sc := NewScope()
	for _, o := range b.opts {
		o(sc)
	}
	if sc.Alert {
		sc.Modal = true
	}
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case Participant:
			c.SetDialogScope(sc)
		}
		return true
	})
	return h.Group(kids...).Render(w)
}

// Props carries props prepared by a wrapper component; see WithProps.
type Props struct{ el *h.Element }

// WithProps passes a prepared element's attributes, classes and children to
// a part, as if they had been passed one by one before the part's other
// children. Wrappers use it to forward props they merged themselves:
//
//	props := h.Spread(h.H2(h.Class("text-lg")), b.Elem(), lib.Cn)
//	primitive.Title(b.Scope(), primitive.WithProps(props))
func WithProps(el *h.Element) Props { return Props{el} }

// part holds what every part shares.
type part struct {
	Scoped
	asChild bool
	props   *h.Element
	w       io.Writer // set by Render: resolves the asChild child in it
}

func (p *part) accept(c any) bool {
	switch v := c.(type) {
	case lib.AsChildFlag:
		p.asChild = true
	case *Scope:
		p.SetDialogScope(v)
	case Inherited:
		p.SetDialogScope(v.scope)
		p.inPortal = p.inPortal || v.inPortal
	case Props:
		if p.props == nil {
			p.props = v.el
		} else {
			p.props = h.Spread(p.props, v.el, nil)
		}
	default:
		return false
	}
	return true
}

// resolve spreads the user's props over defaults and applies asChild.
func (p *part) resolve(name string, defaults, user *h.Element) *h.Element {
	if p.props != nil {
		user = h.Spread(p.props, user, nil)
	}
	props := h.Spread(defaults, user, nil)
	if p.asChild {
		return asChild(name, props, p.w)
	}
	return props
}

// renderIn records the writer the part renders into, for asChild: the
// slotted child is resolved in it (see SlotIn), keeping the style in scope.
func (p *part) renderIn(w io.Writer) { p.w = w }

func state(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}

// --- Trigger -------------------------------------------------------------------------------

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	part
}

// Trigger is Dialog.Trigger: a button that toggles the dialog. Accepts
// lib.AsChild.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *TriggerBuilder) AcceptChild(c any) bool { return b.accept(c) }

// AsChild renders the single child in place of the <button>.
func (b *TriggerBuilder) AsChild() *TriggerBuilder { b.asChild = true; return b }

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	sc := b.Scope()
	d := h.Button(
		h.Type("button"),
		h.Aria("haspopup", "dialog"),
		h.Aria("expanded", map[bool]string{true: "true", false: "false"}[sc.DefaultOpen]),
		h.Aria("controls", sc.ContentID),
		h.Data("state", state(sc.DefaultOpen)),
	)
	return b.resolve("Trigger", d, b.Elem())
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error {
	b.renderIn(w)
	return b.Element().Render(w)
}

// --- Portal --------------------------------------------------------------------------------

// PortalBuilder is returned by Portal.
type PortalBuilder struct {
	h.Builder[PortalBuilder]
	Scoped
}

// Portal is Dialog.Portal: its children are mounted at the end of <body>
// while the dialog is open (each with its own presence, so each plays its
// exit animation before leaving). It renders the inert <template> they
// are cloned from.
func Portal(children ...any) *PortalBuilder {
	b := &PortalBuilder{}
	b.Builder = h.NewBuilder(b, "template")
	return b.Children(children...)
}

// AcceptChild consumes *Scope.
func (b *PortalBuilder) AcceptChild(c any) bool {
	if sc, ok := c.(*Scope); ok {
		b.SetDialogScope(sc)
		return true
	}
	return false
}

// Element resolves the portal's template.
func (b *PortalBuilder) Element() *h.Element {
	sc := b.Scope()
	kids := b.Elem().ChildNodes()
	tree.Walk(kids, func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case interface{ SetInPortal() }:
			c.SetInPortal()
		}
		return true
	})
	return template(sc, false, kids...)
}

// Render renders the template.
func (b *PortalBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// template returns the <template> the behavior mounts from. inPlace makes
// it mount where the template is instead of in <body>.
func template(sc *Scope, inPlace bool, children ...any) *h.Element {
	kind := ""
	if sc.Alert {
		kind = "alert"
	}
	return h.Template(
		h.Data("parity-ignore", ""),
		h.Data("radix-dialog", kind),
		h.Data("content-id", sc.ContentID),
		h.If(!sc.Modal, h.Data("modal", "false")),
		h.If(sc.DefaultOpen, h.Data("default-open", "")),
		h.If(inPlace, h.Data("in-place", "")),
		h.Group(children...),
	)
}

// --- Overlay -------------------------------------------------------------------------------

// OverlayBuilder is returned by Overlay.
type OverlayBuilder struct {
	h.Builder[OverlayBuilder]
	part
}

// Overlay is Dialog.Overlay: the layer covering the page behind modal
// content (it holds the scroll lock). Non-modal dialogs render none.
func Overlay(children ...any) *OverlayBuilder {
	b := &OverlayBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *OverlayBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the overlay as it is when open.
func (b *OverlayBuilder) Element() *h.Element {
	sc := b.Scope()
	if !sc.Modal {
		return nil
	}
	d := h.Div(h.Data("state", "open"), h.Attr(markerOverlay, ""))
	el := b.resolve("Overlay", d, b.Elem())
	// style={{ pointerEvents: "auto", ...props.style }}
	st, _ := el.GetAttr("style")
	el.Attr("style", "pointer-events: auto;"+st)
	if !b.inPortal {
		return template(sc, true, el)
	}
	return el
}

// Render renders the overlay.
func (b *OverlayBuilder) Render(w io.Writer) error {
	el := b.Element()
	if el == nil {
		return nil
	}
	return el.Render(w)
}

// --- Content -------------------------------------------------------------------------------

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	part
}

// Content is Dialog.Content: the dialog itself (role="dialog", or
// "alertdialog" in an alert dialog), labelled by its Title and described
// by its Description.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *ContentBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the content as it is when open.
func (b *ContentBuilder) Element() *h.Element {
	sc := b.Scope()
	role := "dialog"
	if sc.Alert {
		role = "alertdialog"
	}
	d := h.Div(
		h.Role(role),
		h.ID(sc.ContentID),
		h.Aria("describedby", sc.DescriptionID),
		h.Aria("labelledby", sc.TitleID),
		h.Data("state", "open"),
		h.TabIndex(-1),
	)
	el := b.resolve("Content", d, b.Elem())
	if !b.inPortal {
		return template(sc, true, el)
	}
	return el
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Title, Description ------------------------------------------------------------------------

// TitleBuilder is returned by Title.
type TitleBuilder struct {
	h.Builder[TitleBuilder]
	part
}

// Title is Dialog.Title: an <h2> labelling the content.
func Title(children ...any) *TitleBuilder {
	b := &TitleBuilder{}
	b.Builder = h.NewBuilder(b, "h2")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *TitleBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the title.
func (b *TitleBuilder) Element() *h.Element {
	return b.resolve("Title", h.H2(h.ID(b.Scope().TitleID)), b.Elem())
}

// Render renders the title.
func (b *TitleBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// DescriptionBuilder is returned by Description.
type DescriptionBuilder struct {
	h.Builder[DescriptionBuilder]
	part
}

// Description is Dialog.Description: a <p> describing the content.
func Description(children ...any) *DescriptionBuilder {
	b := &DescriptionBuilder{}
	b.Builder = h.NewBuilder(b, "p")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *DescriptionBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the description.
func (b *DescriptionBuilder) Element() *h.Element {
	return b.resolve("Description", h.P(h.ID(b.Scope().DescriptionID)), b.Elem())
}

// Render renders the description.
func (b *DescriptionBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Close, Cancel -------------------------------------------------------------------------

// CloseBuilder is returned by Close and Cancel.
type CloseBuilder struct {
	h.Builder[CloseBuilder]
	part
	cancel bool
}

// Close is Dialog.Close (and AlertDialog.Action): a button that closes the
// dialog. Accepts lib.AsChild.
func Close(children ...any) *CloseBuilder {
	b := &CloseBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// Cancel is AlertDialog.Cancel: a Close button that also receives focus
// when the alert dialog opens.
func Cancel(children ...any) *CloseBuilder {
	b := Close(children...)
	b.cancel = true
	return b
}

// AcceptChild consumes lib.AsChild and *Scope.
func (b *CloseBuilder) AcceptChild(c any) bool { return b.accept(c) }

// AsChild renders the single child in place of the <button>.
func (b *CloseBuilder) AsChild() *CloseBuilder { b.asChild = true; return b }

// Element resolves the button.
func (b *CloseBuilder) Element() *h.Element {
	d := h.Button(h.Type("button"), h.Attr(markerClose, b.Scope().ContentID), h.If(b.cancel, h.Attr(markerCancel, "")))
	return b.resolve(map[bool]string{false: "Close", true: "Cancel"}[b.cancel], d, b.Elem())
}

// Render renders the button.
func (b *CloseBuilder) Render(w io.Writer) error {
	b.renderIn(w)
	return b.Element().Render(w)
}
