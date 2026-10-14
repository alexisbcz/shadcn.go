// Package alert is the shadcn/ui Alert: a callout for user attention.
//
//	alert.Alert(alert.VariantDestructive,
//	    lucide.CircleAlert(),
//	    alert.Title("Error"),
//	    alert.Description("Your session has expired."),
//	)
package alert

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Variant selects the alert's look.
type Variant string

// Alert variants.
const (
	VariantDefault     Variant = "default"
	VariantDestructive Variant = "destructive"
)

// VariantValues lists every variant, in the order shadcn declares them.
var VariantValues = []Variant{VariantDefault, VariantDestructive}

var alertVariants = cva.New(
	"relative grid w-full grid-cols-[0_1fr] items-start gap-y-0.5 rounded-lg border px-4 py-3 text-sm has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] has-[>svg]:gap-x-3 [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default":     "bg-card text-card-foreground",
				"destructive": "bg-card text-destructive *:data-[slot=alert-description]:text-destructive/90 [&>svg]:text-current",
			},
		},
		Defaults: cva.Selection{"variant": "default"},
		Order:    []string{"variant"},
	},
)

// Variants returns the classes of an alert with the given variant plus
// className, conflicts resolved, like cn(alertVariants({ variant }), className).
func Variants(variant Variant, className ...string) string {
	return lib.Cn(append([]string{alertVariants.Classes(cva.Selection{"variant": string(variant)})}, className...)...)
}

// Builder is returned by Alert.
type Builder struct {
	h.Builder[Builder]
	variant Variant
}

// Alert returns a <div role="alert">. Children may include a Variant.
func Alert(children ...any) *Builder {
	b := &Builder{variant: VariantDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Variant children.
func (b *Builder) AcceptChild(c any) bool {
	if v, ok := c.(Variant); ok {
		b.Variant(v)
		return true
	}
	return false
}

// Variant sets the variant. An empty value selects the default.
func (b *Builder) Variant(v Variant) *Builder {
	if v == "" {
		v = VariantDefault
	}
	b.variant = v
	return b
}

// Element resolves the alert to its <div>.
func (b *Builder) Element() *h.Element {
	base := h.Div(h.Data("slot", "alert"), h.Role("alert"), h.Class(Variants(b.variant)))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the alert.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

var (
	titleBase       = h.Div(h.Data("slot", "alert-title"), h.Class(lib.Cn("col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight")))
	descriptionBase = h.Div(h.Data("slot", "alert-description"), h.Class(lib.Cn("col-start-2 grid justify-items-start gap-1 text-sm text-muted-foreground [&_p]:leading-relaxed")))
)

// Title returns the alert's title <div>.
func Title(children ...any) *h.Element { return h.Spread(titleBase, h.Div(children...), lib.Cn) }

// Description returns the alert's description <div>.
func Description(children ...any) *h.Element {
	return h.Spread(descriptionBase, h.Div(children...), lib.Cn)
}
