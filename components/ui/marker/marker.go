// Package marker is the shadcn/ui Marker: a small inline note in a
// conversation (a status, an event, a separator).
//
//	marker.Marker(marker.VariantSeparator, marker.Content("Today"))
//	marker.Marker(marker.Icon(lucide.GitBranch()), marker.Content("Switched branch"))
package marker

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Variant selects the marker's look.
type Variant string

// Marker variants.
const (
	VariantDefault   Variant = "default"
	VariantSeparator Variant = "separator"
	VariantBorder    Variant = "border"
)

// VariantValues lists every variant, in the order shadcn declares them.
var VariantValues = []Variant{VariantDefault, VariantSeparator, VariantBorder}

var markerVariants = cva.New(
	"group/marker relative flex min-h-4 w-full items-center gap-2 text-left text-sm text-muted-foreground [&_svg:not([class*='size-'])]:size-4 [a]:underline [a]:underline-offset-3 [a]:hover:text-foreground",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default":   "",
				"separator": "before:mr-1 before:h-px before:min-w-0 before:flex-1 before:bg-border after:ml-1 after:h-px after:min-w-0 after:flex-1 after:bg-border",
				"border":    "border-b border-border pb-2",
			},
		},
		Order: []string{"variant"},
	},
)

// Variants returns the classes of a marker, like
// cn(markerVariants({ variant, className })).
func Variants(v Variant, className ...string) string {
	return lib.Cn(markerVariants.Classes(cva.Selection{"variant": string(v)}, className...))
}

// Builder is returned by Marker.
type Builder struct {
	h.Builder[Builder]
	variant Variant
	asChild bool
}

// Marker returns a marker <div>. Children may include a Variant and
// lib.AsChild (to render a link or a button instead).
func Marker(children ...any) *Builder {
	b := &Builder{variant: VariantDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Variant and lib.AsChild.
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

// AsChild makes the marker render its single element child.
func (b *Builder) AsChild() *Builder { b.asChild = true; return b }

// Element resolves the marker.
func (b *Builder) Element() *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Div(h.Data("slot", "marker"), h.Data("variant", string(b.variant)), h.Class(Variants(b.variant, className)))
	el := h.Spread(base, props, nil)
	if !b.asChild {
		return el
	}
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("marker: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := h.Slot(el.Clone().SetChildren(), kids[0], nil)
	if err != nil {
		return h.Errored(fmt.Errorf("marker: AsChild: %w", err))
	}
	return out
}

// Render renders the marker.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

var (
	iconBase    = h.Span(h.Data("slot", "marker-icon"), h.Aria("hidden", "true"), h.Class(lib.Cn("size-4 shrink-0 [&_svg:not([class*='size-'])]:size-4")))
	contentBase = h.Span(h.Data("slot", "marker-content"), h.Class(lib.Cn("min-w-0 wrap-break-word group-data-[variant=separator]/marker:flex-none group-data-[variant=separator]/marker:text-center *:[a]:underline *:[a]:underline-offset-3 *:[a]:hover:text-foreground")))
)

// Icon holds the marker's icon (aria-hidden).
func Icon(children ...any) *h.Element { return h.Spread(iconBase, h.Span(children...), lib.Cn) }

// Content holds the marker's text.
func Content(children ...any) *h.Element { return h.Spread(contentBase, h.Span(children...), lib.Cn) }
