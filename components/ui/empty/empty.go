// Package empty is the shadcn/ui Empty: an empty state with media, a title,
// a description and actions.
//
//	empty.Empty(
//		empty.Header(
//			empty.Media(empty.MediaIcon, tabler.IconFolderCode()),
//			empty.Title("No Projects Yet"),
//			empty.Description("Get started by creating your first project."),
//		),
//		empty.Content(button.Button("Create Project")),
//	)
package empty

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Builder is returned by every Empty sub-component.
type Builder struct {
	h.Builder[Builder]
	base    *h.Element
	variant MediaVariant
	media   bool
}

func styled(tag, slot, class string, children []any) *Builder {
	b := &Builder{base: h.El(tag, h.Data("slot", slot), h.Class(class))}
	b.Builder = h.NewBuilder(b, tag)
	return b.Children(children...)
}

// Empty is the root: a centered column, dashed when given a border.
func Empty(children ...any) *Builder {
	return styled("div", "empty", "flex min-w-0 flex-1 flex-col items-center justify-center gap-6 rounded-lg border-dashed p-6 text-center text-balance md:p-12", children)
}

// Header groups the media, title and description.
func Header(children ...any) *Builder {
	return styled("div", "empty-header", "flex max-w-sm flex-col items-center gap-2 text-center", children)
}

// Title is the empty state's title.
func Title(children ...any) *Builder {
	return styled("div", "empty-title", "text-lg font-medium tracking-tight", children)
}

// Description is the empty state's description. Like upstream, it renders
// a <div>.
func Description(children ...any) *Builder {
	return styled("div", "empty-description", "text-sm/relaxed text-muted-foreground [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary", children)
}

// Content holds the actions.
func Content(children ...any) *Builder {
	return styled("div", "empty-content", "flex w-full max-w-sm min-w-0 flex-col items-center gap-4 text-sm text-balance", children)
}

// MediaVariant selects the look of Media.
type MediaVariant string

// Media variants.
const (
	MediaDefault MediaVariant = "default"
	MediaIcon    MediaVariant = "icon"
)

var mediaVariants = cva.New(
	"mb-2 flex shrink-0 items-center justify-center [&_svg]:pointer-events-none [&_svg]:shrink-0",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default": "bg-transparent",
				"icon":    "flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted text-foreground [&_svg:not([class*='size-'])]:size-6",
			},
		},
		Defaults: cva.Selection{"variant": "default"},
	},
)

// MediaVariants returns the classes of Media for a variant plus className,
// like emptyMediaVariants() upstream.
func MediaVariants(v MediaVariant, className ...string) string {
	return lib.Cn(mediaVariants.Classes(cva.Selection{"variant": string(v)}, className...))
}

// Media holds an icon, avatar or image. Children may include a
// MediaVariant.
func Media(children ...any) *Builder {
	b := &Builder{media: true, variant: MediaDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes a MediaVariant child (Media only).
func (b *Builder) AcceptChild(c any) bool {
	if v, ok := c.(MediaVariant); ok && b.media {
		b.Variant(v)
		return true
	}
	return false
}

// Variant sets Media's variant. An empty value selects the default.
func (b *Builder) Variant(v MediaVariant) *Builder {
	if v == "" {
		v = MediaDefault
	}
	b.variant = v
	return b
}

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element {
	if b.media {
		// cn(emptyMediaVariants({ variant, className })): the user's classes
		// go through cva, so they are merged once with the variant classes.
		props := b.Elem().Clone()
		className := props.ClassName()
		props.SetClass("")
		base := h.Div(h.Data("slot", "empty-icon"), h.Data("variant", string(b.variant)), h.Class(MediaVariants(b.variant, className)))
		return h.Spread(base, props, lib.Cn)
	}
	return h.Spread(b.base, b.Elem(), lib.Cn)
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }
