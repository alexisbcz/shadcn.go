// Package button is the shadcn/ui Button: a button, or a component that
// looks like a button.
//
//	button.Button(button.VariantOutline, button.SizeSm, "Save")
//	button.Button(lib.AsChild, h.A(h.Href("/login"), "Login"))
package button

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Variant selects the button's look.
type Variant string

// Button variants.
const (
	VariantDefault     Variant = "default"
	VariantDestructive Variant = "destructive"
	VariantOutline     Variant = "outline"
	VariantSecondary   Variant = "secondary"
	VariantGhost       Variant = "ghost"
	VariantLink        Variant = "link"
)

// Size selects the button's size.
type Size string

// Button sizes.
const (
	SizeDefault Size = "default"
	SizeXs      Size = "xs"
	SizeSm      Size = "sm"
	SizeLg      Size = "lg"
	SizeIcon    Size = "icon"
	SizeIconXs  Size = "icon-xs"
	SizeIconSm  Size = "icon-sm"
	SizeIconLg  Size = "icon-lg"
)

// VariantValues and SizeValues list every option, in the order shadcn
// declares them.
var (
	VariantValues = []Variant{VariantDefault, VariantDestructive, VariantOutline, VariantSecondary, VariantGhost, VariantLink}
	SizeValues    = []Size{SizeDefault, SizeXs, SizeSm, SizeLg, SizeIcon, SizeIconXs, SizeIconSm, SizeIconLg}
)

var buttonVariants = cva.New(
	"inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default":     "bg-primary text-primary-foreground hover:bg-primary/90",
				"destructive": "bg-destructive text-white hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40",
				"outline":     "border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50",
				"secondary":   "bg-secondary text-secondary-foreground hover:bg-secondary/80",
				"ghost":       "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50",
				"link":        "text-primary underline-offset-4 hover:underline",
			},
			"size": {
				"default": "h-9 px-4 py-2 has-[>svg]:px-3",
				"xs":      "h-6 gap-1 rounded-md px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3",
				"sm":      "h-8 gap-1.5 rounded-md px-3 has-[>svg]:px-2.5",
				"lg":      "h-10 rounded-md px-6 has-[>svg]:px-4",
				"icon":    "size-9",
				"icon-xs": "size-6 rounded-md [&_svg:not([class*='size-'])]:size-3",
				"icon-sm": "size-8",
				"icon-lg": "size-10",
			},
		},
		Defaults: cva.Selection{"variant": "default", "size": "default"},
		Order:    []string{"variant", "size"},
	},
)

// Variants returns the classes of a button with the given variant and
// size, plus className, conflicts resolved. Use it to style other elements
// as buttons, like buttonVariants() in shadcn. Empty values select the
// defaults.
func Variants(variant Variant, size Size, className ...string) string {
	sel := cva.Selection{"variant": string(variant), "size": string(size)}
	return lib.Cn(buttonVariants.Classes(sel, className...))
}

// Builder is returned by Button.
type Builder struct {
	h.Builder[Builder]
	variant Variant
	size    Size
	asChild bool
}

// Button returns a shadcn button. Children may include a Variant, a Size
// and lib.AsChild, besides anything a libhtml element accepts:
//
//	button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Submit"), icon)
func Button(children ...any) *Builder {
	b := &Builder{variant: VariantDefault, size: SizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Variant, Size and lib.AsChild children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Variant:
		b.Variant(v)
	case Size:
		b.Size(v)
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

// Size sets the size. An empty value selects the default.
func (b *Builder) Size(s Size) *Builder {
	if s == "" {
		s = SizeDefault
	}
	b.size = s
	return b
}

// AsChild renders the button's single element child in place of a <button>,
// with the button's classes and attributes merged in. See lib.AsChild.
func (b *Builder) AsChild() *Builder { b.asChild = true; return b }

// Button-specific attributes.
func (b *Builder) Type(v string) *Builder    { b.Elem().Type(v); return b }
func (b *Builder) Name(v string) *Builder    { b.Elem().Name(v); return b }
func (b *Builder) Value(v string) *Builder   { b.Elem().Value(v); return b }
func (b *Builder) Form(v string) *Builder    { b.Elem().Form(v); return b }
func (b *Builder) Disabled(on bool) *Builder { b.Elem().Disabled(on); return b }
func (b *Builder) AutoFocus(on bool) *Builder {
	b.Elem().AutoFocus(on)
	return b
}

// Element resolves the button to the element it renders.
func (b *Builder) Element() *h.Element {
	// User attributes override the defaults; classes are merged with the
	// user's last, exactly like cn(buttonVariants({ variant, size, className })).
	btn, err := h.Slot(buttonDefaults(b.variant, b.size), b.Elem(), lib.Cn)
	if err != nil {
		return h.Errored(err)
	}
	if !b.asChild {
		return btn
	}
	kids := btn.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("button: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := h.Slot(btn.Clone().SetChildren(), kids[0], lib.Cn)
	if err != nil {
		return h.Errored(fmt.Errorf("button: AsChild: %w", err))
	}
	return out
}

// Render renders the button.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// buttonDefaults returns the element carrying a variant and size's classes
// and data attributes. Slot only reads it, so the known combinations are
// built once and shared.
func buttonDefaults(v Variant, s Size) *h.Element {
	if el, ok := buttonDefaultsCache[[2]string{string(v), string(s)}]; ok {
		return el
	}
	return newButtonDefaults(v, s)
}

func newButtonDefaults(v Variant, s Size) *h.Element {
	return h.Button(
		h.Class(Variants(v, s)),
		h.Data("slot", "button"),
		h.Data("variant", string(v)),
		h.Data("size", string(s)),
	)
}

var buttonDefaultsCache = func() map[[2]string]*h.Element {
	m := make(map[[2]string]*h.Element)
	for _, v := range VariantValues {
		for _, s := range SizeValues {
			m[[2]string{string(v), string(s)}] = newButtonDefaults(v, s)
		}
	}
	return m
}()
