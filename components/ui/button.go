package ui

import (
	"fmt"
	"io"

	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// ButtonVariant selects the button's look.
type ButtonVariant string

// Button variants.
const (
	ButtonVariantDefault     ButtonVariant = "default"
	ButtonVariantDestructive ButtonVariant = "destructive"
	ButtonVariantOutline     ButtonVariant = "outline"
	ButtonVariantSecondary   ButtonVariant = "secondary"
	ButtonVariantGhost       ButtonVariant = "ghost"
	ButtonVariantLink        ButtonVariant = "link"
)

// ButtonSize selects the button's size.
type ButtonSize string

// Button sizes.
const (
	ButtonSizeDefault ButtonSize = "default"
	ButtonSizeXs      ButtonSize = "xs"
	ButtonSizeSm      ButtonSize = "sm"
	ButtonSizeLg      ButtonSize = "lg"
	ButtonSizeIcon    ButtonSize = "icon"
	ButtonSizeIconXs  ButtonSize = "icon-xs"
	ButtonSizeIconSm  ButtonSize = "icon-sm"
	ButtonSizeIconLg  ButtonSize = "icon-lg"
)

// ButtonVariantValues and ButtonSizeValues list every option, in the order
// shadcn declares them.
var (
	ButtonVariantValues = []ButtonVariant{ButtonVariantDefault, ButtonVariantDestructive, ButtonVariantOutline, ButtonVariantSecondary, ButtonVariantGhost, ButtonVariantLink}
	ButtonSizeValues    = []ButtonSize{ButtonSizeDefault, ButtonSizeXs, ButtonSizeSm, ButtonSizeLg, ButtonSizeIcon, ButtonSizeIconXs, ButtonSizeIconSm, ButtonSizeIconLg}
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

// ButtonVariants returns the classes of a button with the given variant and
// size, plus className, conflicts resolved. Use it to style other elements
// as buttons, like buttonVariants() in shadcn. Empty values select the
// defaults.
func ButtonVariants(variant ButtonVariant, size ButtonSize, className ...string) string {
	sel := cva.Selection{"variant": string(variant), "size": string(size)}
	return Cn(buttonVariants.Classes(sel, className...))
}

// ButtonBuilder is returned by Button.
type ButtonBuilder struct {
	h.Builder[ButtonBuilder]
	variant ButtonVariant
	size    ButtonSize
	asChild bool
}

// Button returns a shadcn button. Children may include a ButtonVariant, a
// ButtonSize and AsChild, besides anything a libhtml element accepts:
//
//	ui.Button(ui.ButtonVariantOutline, ui.ButtonSizeIcon, h.Aria("label", "Submit"), icon)
func Button(children ...any) *ButtonBuilder {
	b := &ButtonBuilder{variant: ButtonVariantDefault, size: ButtonSizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes ButtonVariant, ButtonSize and AsChild children.
func (b *ButtonBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case ButtonVariant:
		b.Variant(v)
	case ButtonSize:
		b.Size(v)
	case asChildFlag:
		b.asChild = true
	default:
		return false
	}
	return true
}

// Variant sets the variant. An empty value selects the default.
func (b *ButtonBuilder) Variant(v ButtonVariant) *ButtonBuilder {
	if v == "" {
		v = ButtonVariantDefault
	}
	b.variant = v
	return b
}

// Size sets the size. An empty value selects the default.
func (b *ButtonBuilder) Size(s ButtonSize) *ButtonBuilder {
	if s == "" {
		s = ButtonSizeDefault
	}
	b.size = s
	return b
}

// AsChild renders the button's single element child in place of a <button>,
// with the button's classes and attributes merged in. See ui.AsChild.
func (b *ButtonBuilder) AsChild() *ButtonBuilder { b.asChild = true; return b }

// Button-specific attributes.
func (b *ButtonBuilder) Type(v string) *ButtonBuilder    { b.Elem().Type(v); return b }
func (b *ButtonBuilder) Name(v string) *ButtonBuilder    { b.Elem().Name(v); return b }
func (b *ButtonBuilder) Value(v string) *ButtonBuilder   { b.Elem().Value(v); return b }
func (b *ButtonBuilder) Form(v string) *ButtonBuilder    { b.Elem().Form(v); return b }
func (b *ButtonBuilder) Disabled(on bool) *ButtonBuilder { b.Elem().Disabled(on); return b }
func (b *ButtonBuilder) AutoFocus(on bool) *ButtonBuilder {
	b.Elem().AutoFocus(on)
	return b
}

// Element resolves the button to the element it renders.
func (b *ButtonBuilder) Element() *h.Element {
	// User attributes override the defaults; classes are merged with the
	// user's last, exactly like cn(buttonVariants({ variant, size, className })).
	btn, err := h.Slot(buttonDefaults(b.variant, b.size), b.Elem(), Cn)
	if err != nil {
		return h.Errored(err)
	}
	if !b.asChild {
		return btn
	}
	kids := btn.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("ui.Button: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := h.Slot(btn.Clone().SetChildren(), kids[0], Cn)
	if err != nil {
		return h.Errored(fmt.Errorf("ui.Button: AsChild: %w", err))
	}
	return out
}

// Render renders the button.
func (b *ButtonBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// buttonDefaults returns the element carrying a variant and size's classes
// and data attributes. Slot only reads it, so the known combinations are
// built once and shared.
func buttonDefaults(v ButtonVariant, s ButtonSize) *h.Element {
	if el, ok := buttonDefaultsCache[[2]string{string(v), string(s)}]; ok {
		return el
	}
	return newButtonDefaults(v, s)
}

func newButtonDefaults(v ButtonVariant, s ButtonSize) *h.Element {
	return h.Button(
		h.Class(ButtonVariants(v, s)),
		h.Data("slot", "button"),
		h.Data("variant", string(v)),
		h.Data("size", string(s)),
	)
}

var buttonDefaultsCache = func() map[[2]string]*h.Element {
	m := make(map[[2]string]*h.Element)
	for _, v := range ButtonVariantValues {
		for _, s := range ButtonSizeValues {
			m[[2]string{string(v), string(s)}] = newButtonDefaults(v, s)
		}
	}
	return m
}()
