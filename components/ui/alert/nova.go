package alert

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// registry/bases/radix/ui/alert.tsx, used inside a style.Scope.
var (
	baseFile     = style.NewFile("ui/alert.tsx")
	baseVariants = baseFile.CVA("cn-alert group/alert relative w-full",
		[]style.Variant{{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-alert-variant-default"},
			{Name: "destructive", Class: "cn-alert-variant-destructive"},
		}}},
		cva.Selection{"variant": "default"})
	baseTitle       = baseFile.Cn("cn-alert-title [&_a]:underline [&_a]:underline-offset-3 [&_a]:hover:text-foreground", style.Arg)
	baseDescription = baseFile.Cn("cn-alert-description [&_a]:underline [&_a]:underline-offset-3 [&_a]:hover:text-foreground", style.Arg)
	baseAction      = baseFile.Cn("cn-alert-action", style.Arg)
)

// styledElement resolves the alert in style st.
func (b *Builder) styledElement(st style.Style) *h.Element {
	class := lib.Cn(baseVariants.Classes(st, cva.Selection{"variant": string(b.variant)}))
	base := h.Div(h.Data("slot", "alert"), h.Role("alert"), h.Class(class))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the alert, in the style in scope (style.Scope) if any.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn resolves the element as it renders inside w: in the style in
// scope (style.Scope), if any.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		return b.styledElement(st)
	}
	return b.Element()
}

// Title returns the alert's title <div>.
func Title(children ...any) *scoped.Node {
	return scoped.Spread(titleBase, h.Div(children...), scoped.Slot("div", "alert-title", baseTitle))
}

// Description returns the alert's description <div>.
func Description(children ...any) *scoped.Node {
	return scoped.Spread(descriptionBase, h.Div(children...), scoped.Slot("div", "alert-description", baseDescription))
}

// Action returns the alert's action <div> (bases styles only; outside a
// style it renders the same slot without classes).
func Action(children ...any) *scoped.Node {
	return scoped.Spread(h.Div(h.Data("slot", "alert-action")), h.Div(children...), scoped.Slot("div", "alert-action", baseAction))
}
