// Package badge is the shadcn/ui Badge: a badge or a component that looks
// like a badge.
//
//	badge.Badge(badge.VariantSecondary, "Secondary")
//	badge.Badge(lib.AsChild, h.A(h.Href("#"), "Link"))
package badge

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Variant selects the badge's look.
type Variant string

// Badge variants.
const (
	VariantDefault     Variant = "default"
	VariantSecondary   Variant = "secondary"
	VariantDestructive Variant = "destructive"
	VariantOutline     Variant = "outline"
	VariantGhost       Variant = "ghost"
	VariantLink        Variant = "link"
)

// VariantValues lists every variant, in the order shadcn declares them.
var VariantValues = []Variant{VariantDefault, VariantSecondary, VariantDestructive, VariantOutline, VariantGhost, VariantLink}

var badgeVariants = cva.New(
	"inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-full border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default":     "bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
				"secondary":   "bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
				"destructive": "bg-destructive text-white focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40 [a&]:hover:bg-destructive/90",
				"outline":     "border-border text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
				"ghost":       "[a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
				"link":        "text-primary underline-offset-4 [a&]:hover:underline",
			},
		},
		Defaults: cva.Selection{"variant": "default"},
		Order:    []string{"variant"},
	},
)

// Variants returns the classes of a badge with the given variant plus
// className, conflicts resolved, like cn(badgeVariants({ variant }), className).
// Use it to style other elements as badges.
func Variants(variant Variant, className ...string) string {
	return lib.Cn(append([]string{badgeVariants.Classes(cva.Selection{"variant": string(variant)})}, className...)...)
}

// Builder is returned by Badge.
type Builder struct {
	h.Builder[Builder]
	variant Variant
	asChild bool
}

// Badge returns a <span> badge. Children may include a Variant and
// lib.AsChild.
func Badge(children ...any) *Builder {
	b := &Builder{variant: VariantDefault}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// AcceptChild consumes Variant and lib.AsChild children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Variant:
		b.Variant(v)
	case lib.AsChildFlag:
		b.asChild = true
	default:
		return false
	}
	return true
}

// Variant sets the variant. An empty value selects the default.
func (b *Builder) Variant(v Variant) *Builder {
	if v == "" {
		v = VariantDefault
	}
	b.variant = v
	return b
}

// AsChild renders the badge's single element child in place of a <span>.
func (b *Builder) AsChild() *Builder { b.asChild = true; return b }

// Element resolves the badge to the element it renders.
func (b *Builder) Element() *h.Element {
	base := h.Span(h.Data("slot", "badge"), h.Data("variant", string(b.variant)), h.Class(Variants(b.variant)))
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
	out, err := h.Slot(el.Clone().SetChildren(), kids[0], lib.Cn)
	if err != nil {
		return h.Errored(fmt.Errorf("badge: AsChild: %w", err))
	}
	return out
}

// Render renders the badge.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }
