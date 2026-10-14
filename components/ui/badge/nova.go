package badge

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// registry/bases/radix/ui/badge.tsx, used inside a style.Scope.
var baseFile = style.NewFile("ui/badge.tsx")

var baseVariants = baseFile.CVA(
	"cn-badge group/badge inline-flex w-fit shrink-0 items-center justify-center overflow-hidden whitespace-nowrap focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none",
	[]style.Variant{{Name: "variant", Options: []style.Option{
		{Name: "default", Class: "cn-badge-variant-default"},
		{Name: "secondary", Class: "cn-badge-variant-secondary"},
		{Name: "destructive", Class: "cn-badge-variant-destructive"},
		{Name: "outline", Class: "cn-badge-variant-outline"},
		{Name: "ghost", Class: "cn-badge-variant-ghost"},
		{Name: "link", Class: "cn-badge-variant-link"},
	}}},
	cva.Selection{"variant": "default"},
)

// StyleVariants is Variants for style st.
func StyleVariants(st style.Style, variant Variant, className ...string) string {
	if variant == "" {
		variant = VariantDefault
	}
	return lib.Cn(append([]string{baseVariants.Classes(st, cva.Selection{"variant": string(variant)})}, className...)...)
}

// ElementIn resolves the badge as it renders inside w: in the style in
// scope (style.Scope), if any.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		return b.element(scoped.Class(w, StyleVariants(st, b.variant)), w)
	}
	return b.element(Variants(b.variant), w)
}

// Element resolves the badge to the element it renders.
func (b *Builder) Element() *h.Element { return b.element(Variants(b.variant), nil) }

func (b *Builder) element(class string, w io.Writer) *h.Element {
	base := h.Span(h.Data("slot", "badge"), h.Data("variant", string(b.variant)), h.Class(class))
	el, err := h.Slot(base, b.Elem(), lib.Cn)
	if err != nil {
		return h.Errored(err)
	}
	if !b.asChild {
		return el
	}
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("badge: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := h.SlotIn(w, el.Clone().SetChildren(), kids[0], lib.Cn)
	if err != nil {
		return h.Errored(fmt.Errorf("badge: AsChild: %w", err))
	}
	return out
}

// Render renders the badge, in the style in scope (style.Scope) if any.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }
