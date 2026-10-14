// Package toggle is the shadcn/ui Toggle (Radix Toggle): a two-state button
// that can be either on or off.
//
//	toggle.Toggle(toggle.VariantOutline, toggle.SizeSm, h.Aria("label", "Toggle bold"), lucide.Bold())
//	toggle.Toggle(toggle.Pressed(true), "On")
//
// Clicking toggles aria-pressed and data-state ("on"/"off") and dispatches
// a cancelable "shadcn:pressed-change" event (detail: { pressed }) on the
// button; cancelling it keeps the current state (see toggle.js).
package toggle

import (
	_ "embed"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

//go:embed toggle.js
var script string

func init() { runtime.Module("toggle", script) }

// Variant selects the toggle's look.
type Variant string

// Toggle variants.
const (
	VariantDefault Variant = "default"
	VariantOutline Variant = "outline"
)

// Size selects the toggle's size.
type Size string

// Toggle sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
	SizeLg      Size = "lg"
)

// Pressed, passed as a child, sets the initial pressed state (Radix
// defaultPressed / pressed).
type Pressed bool

var toggleVariants = cva.New(
	"inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default": "bg-transparent",
				"outline": "border border-input bg-transparent shadow-xs hover:bg-accent hover:text-accent-foreground",
			},
			"size": {
				"default": "h-9 min-w-9 px-2",
				"sm":      "h-8 min-w-8 px-1.5",
				"lg":      "h-10 min-w-10 px-2.5",
			},
		},
		Defaults: cva.Selection{"variant": "default", "size": "default"},
		Order:    []string{"variant", "size"},
	},
)

// VariantClasses returns toggleVariants({ variant, size }) followed by
// className, unmerged (cva output). Empty values select the defaults.
func VariantClasses(variant Variant, size Size, className ...string) string {
	return toggleVariants.Classes(cva.Selection{"variant": string(variant), "size": string(size)}, className...)
}

// Variants returns the classes of a toggle with the given variant and size,
// plus className, conflicts resolved: cn(toggleVariants({ variant, size,
// className })).
func Variants(variant Variant, size Size, className ...string) string {
	return lib.Cn(VariantClasses(variant, size, className...))
}

// Builder is returned by Toggle.
type Builder struct {
	h.Builder[Builder]
	variant Variant
	size    Size
	pressed bool
}

// Toggle returns a toggle <button>. Children may include a Variant, a Size
// and Pressed, besides anything a libhtml element accepts.
func Toggle(children ...any) *Builder {
	b := &Builder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Variant, Size and Pressed children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Variant:
		b.variant = v
	case Size:
		b.size = v
	case Pressed:
		b.pressed = bool(v)
	default:
		return false
	}
	return true
}

// Variant sets the variant. An empty value selects the default.
func (b *Builder) Variant(v Variant) *Builder { b.variant = v; return b }

// Size sets the size. An empty value selects the default.
func (b *Builder) Size(s Size) *Builder { b.size = s; return b }

// Pressed sets the initial pressed state.
func (b *Builder) Pressed(on bool) *Builder { b.pressed = on; return b }

// Disabled disables the toggle.
func (b *Builder) Disabled(on bool) *Builder { b.Elem().Disabled(on); return b }

// Element resolves the toggle to its <button>.
func (b *Builder) Element() *h.Element {
	props := b.Elem()
	cls := Variants(b.variant, b.size, props.ClassName())
	base := Root(b.pressed, props.Has("disabled"))
	base.Data("slot", "toggle")
	out := h.Spread(base, props, nil)
	return out.SetClass(cls)
}

// Render renders the toggle.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// Root returns the attributes Radix's Toggle.Root renders before the
// caller's props: type="button", aria-pressed, data-state and
// data-disabled. Toggle group items build on it.
func Root(pressed, disabled bool) *h.Element {
	return h.Button(
		h.Type("button"),
		h.Aria("pressed", boolString(pressed)),
		h.Data("state", State(pressed)),
	).When(disabled, func(e *h.Element) { e.Data("disabled", "") })
}

// State returns the data-state value for a pressed state: "on" or "off".
func State(pressed bool) string {
	if pressed {
		return "on"
	}
	return "off"
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
