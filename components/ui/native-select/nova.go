package nativeselect

import (
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/scoped"
	h "shadcn.go/libhtml"
)

// registry/bases/radix/ui/native-select.tsx, used inside a style.Scope.
var (
	baseFile    = style.NewFile("ui/native-select.tsx")
	baseWrapper = baseFile.Cn("cn-native-select-wrapper group/native-select relative w-fit has-[select:disabled]:opacity-50", style.Arg)
	baseSelect  = baseFile.Attr("cn-native-select outline-none disabled:pointer-events-none disabled:cursor-not-allowed")
	baseIcon    = baseFile.Attr("cn-native-select-icon pointer-events-none absolute select-none")
)

// styledElement resolves the component in style st. Unlike new-york-v4,
// the className goes to the wrapper and the other props to the <select>.
func (b *Builder) styledElement(st style.Style, w io.Writer) *h.Element {
	props := b.Elem()
	sel := h.Spread(h.Select(h.Data("slot", "native-select"), h.Data("size", string(b.size))), props.Clone().SetClass(""), nil)
	sel.SetClass(scoped.Class(w, baseSelect.Class(st)))
	return h.Div(
		h.Class(baseWrapper.Class(st, props.ClassName())), // className is the user's: not mapped
		h.Data("slot", "native-select-wrapper"),
		h.Data("size", string(b.size)),
		sel,
		lucide.ChevronDown(
			h.Class(scoped.Class(w, baseIcon.Class(st))),
			h.Aria("hidden", "true"),
			h.Data("slot", "native-select-icon"),
		),
	)
}

// Render renders the component, in the style in scope (style.Scope) if any.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn resolves the element as it renders inside w: in the style in
// scope (style.Scope), if any.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	if st, ok := style.FromWriter(w); ok {
		return b.styledElement(st, w)
	}
	return b.Element()
}
