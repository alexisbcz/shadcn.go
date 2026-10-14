// Package nativeselect is the shadcn/ui NativeSelect: a styled native
// <select>.
//
//	nativeselect.NativeSelect(
//		nativeselect.Option(h.Value(""), "Select status"),
//		nativeselect.OptGroup(h.LabelAttr("Work"), nativeselect.Option(h.Value("todo"), "Todo")),
//	)
//
// Attributes and classes given to NativeSelect go to the <select>, like
// upstream, not to its wrapper.
package nativeselect

import (
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	h "shadcn.go/libhtml"
)

// Size selects the select's height.
type Size string

// Sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
)

// Builder is returned by NativeSelect.
type Builder struct {
	h.Builder[Builder]
	size Size
}

// NativeSelect returns the select. Children may include a Size.
func NativeSelect(children ...any) *Builder {
	b := &Builder{size: SizeDefault}
	b.Builder = h.NewBuilder(b, "select")
	return b.Children(children...)
}

// AcceptChild consumes a Size child.
func (b *Builder) AcceptChild(c any) bool {
	if s, ok := c.(Size); ok {
		b.Size(s)
		return true
	}
	return false
}

// Size sets the size. An empty value selects the default.
func (b *Builder) Size(s Size) *Builder {
	if s == "" {
		s = SizeDefault
	}
	b.size = s
	return b
}

// Select-specific attributes.
func (b *Builder) Name(v string) *Builder     { b.Elem().Name(v); return b }
func (b *Builder) Form(v string) *Builder     { b.Elem().Form(v); return b }
func (b *Builder) Disabled(on bool) *Builder  { b.Elem().Disabled(on); return b }
func (b *Builder) Required(on bool) *Builder  { b.Elem().Bool("required", on); return b }
func (b *Builder) Multiple(on bool) *Builder  { b.Elem().Bool("multiple", on); return b }
func (b *Builder) AutoFocus(on bool) *Builder { b.Elem().AutoFocus(on); return b }

const selectClass = "h-9 w-full min-w-0 appearance-none rounded-md border border-input bg-transparent px-3 py-2 pr-9 text-sm shadow-xs transition-[color,box-shadow] outline-none selection:bg-primary selection:text-primary-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed data-[size=sm]:h-8 data-[size=sm]:py-1 dark:bg-input/30 dark:hover:bg-input/50"

// Element resolves the component to its wrapper <div>.
func (b *Builder) Element() *h.Element {
	props := b.Elem()
	sel := h.Spread(h.Select(h.Data("slot", "native-select"), h.Data("size", string(b.size))), props.Clone().SetClass(""), nil)
	sel.SetClass(lib.Cn(
		selectClass,
		"focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
		"aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40",
		props.ClassName(),
	))
	return h.Div(
		h.Class("group/native-select relative w-fit has-[select:disabled]:opacity-50"),
		h.Data("slot", "native-select-wrapper"),
		sel,
		lucide.ChevronDown(
			h.Class("pointer-events-none absolute top-1/2 right-3.5 size-4 -translate-y-1/2 text-muted-foreground opacity-50 select-none"),
			h.Aria("hidden", "true"),
			h.Data("slot", "native-select-icon"),
		),
	)
}

// OptionBuilder is returned by Option and OptGroup.
type OptionBuilder struct {
	h.Builder[OptionBuilder]
	slot string
}

func option(tag, slot string, children []any) *OptionBuilder {
	b := &OptionBuilder{slot: slot}
	b.Builder = h.NewBuilder(b, tag)
	return b.Children(children...)
}

// Option is an <option>.
func Option(children ...any) *OptionBuilder {
	return option("option", "native-select-option", children)
}

// OptGroup is an <optgroup>; set its label with h.LabelAttr.
func OptGroup(children ...any) *OptionBuilder {
	return option("optgroup", "native-select-optgroup", children)
}

// Element resolves the option to its element.
func (b *OptionBuilder) Element() *h.Element {
	base := h.El(b.Elem().Tag(), h.Data("slot", b.slot), h.Class("bg-[Canvas] text-[CanvasText]"))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the option.
func (b *OptionBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
