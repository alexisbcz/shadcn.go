package empty

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// registry/bases/radix/ui/empty.tsx, used inside a style.Scope.
var (
	baseFile          = style.NewFile("ui/empty.tsx")
	baseMediaVariants = baseFile.CVA("cn-empty-media flex shrink-0 items-center justify-center [&_svg]:pointer-events-none [&_svg]:shrink-0",
		[]style.Variant{{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-empty-media-default"},
			{Name: "icon", Class: "cn-empty-media-icon"},
		}}},
		cva.Selection{"variant": "default"})
	baseSites = map[string]*style.Site{
		"empty":             baseFile.Cn("cn-empty flex w-full min-w-0 flex-1 flex-col items-center justify-center text-center text-balance", style.Arg),
		"empty-header":      baseFile.Cn("cn-empty-header flex max-w-sm flex-col items-center", style.Arg),
		"empty-title":       baseFile.Cn("cn-empty-title cn-font-heading", style.Arg),
		"empty-description": baseFile.Cn("cn-empty-description text-muted-foreground [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary", style.Arg),
		"empty-content":     baseFile.Cn("cn-empty-content flex w-full max-w-sm min-w-0 flex-col items-center text-balance", style.Arg),
	}
)

// styledElement resolves the component in style st.
func (b *Builder) styledElement(st style.Style) *h.Element {
	if b.media {
		props := b.Elem().Clone()
		className := props.ClassName()
		props.SetClass("")
		class := lib.Cn(baseMediaVariants.Classes(st, cva.Selection{"variant": string(b.variant)}, className))
		base := h.Div(h.Data("slot", "empty-icon"), h.Data("variant", string(b.variant)), h.Class(class))
		return h.Spread(base, props, lib.Cn)
	}
	slot, _ := b.base.GetAttr("data-slot")
	base := h.El(b.base.Tag(), h.Data("slot", slot), h.Class(baseSites[slot].Class(st, "")))
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the component, in the style in scope (style.Scope) if any.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn resolves the element as it renders inside w: in the style in
// scope (style.Scope), if any.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		return b.styledElement(st)
	}
	return b.Element()
}
