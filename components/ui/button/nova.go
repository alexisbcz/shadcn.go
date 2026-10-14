package button

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// registry/bases/radix/ui/button.tsx, used inside a style.Scope.
var baseFile = style.NewFile("ui/button.tsx")

var baseVariants = baseFile.CVA(
	"cn-button group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
	[]style.Variant{
		{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-button-variant-default"},
			{Name: "outline", Class: "cn-button-variant-outline"},
			{Name: "secondary", Class: "cn-button-variant-secondary"},
			{Name: "ghost", Class: "cn-button-variant-ghost"},
			{Name: "destructive", Class: "cn-button-variant-destructive"},
			{Name: "link", Class: "cn-button-variant-link"},
		}},
		{Name: "size", Options: []style.Option{
			{Name: "default", Class: "cn-button-size-default"},
			{Name: "xs", Class: "cn-button-size-xs"},
			{Name: "sm", Class: "cn-button-size-sm"},
			{Name: "lg", Class: "cn-button-size-lg"},
			{Name: "icon", Class: "cn-button-size-icon"},
			{Name: "icon-xs", Class: "cn-button-size-icon-xs"},
			{Name: "icon-sm", Class: "cn-button-size-icon-sm"},
			{Name: "icon-lg", Class: "cn-button-size-icon-lg"},
		}},
	},
	cva.Selection{"variant": "default", "size": "default"},
)

// StyleVariants is Variants for style st (buttonVariants of
// styles/radix-<st>/ui/button.tsx).
func StyleVariants(st style.Style, variant Variant, size Size, className ...string) string {
	if variant == "" {
		variant = VariantDefault
	}
	if size == "" {
		size = SizeDefault
	}
	sel := cva.Selection{"variant": string(variant), "size": string(size)}
	return lib.Cn(baseVariants.Classes(st, sel, className...))
}

// ElementIn resolves the button as it renders inside w: in the style in
// scope (style.Scope), if any.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, ok := style.FromWriter(w)
	if !ok {
		return b.element(buttonDefaults(b.variant, b.size), w)
	}
	return b.element(h.Button(
		h.Class(StyleVariants(st, b.variant, b.size)),
		h.Data("slot", "button"),
		h.Data("variant", string(b.variant)),
		h.Data("size", string(b.size)),
	), w)
}

// Render renders the button, in the style in scope (style.Scope) if any.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }
