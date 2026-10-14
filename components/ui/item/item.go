// Package item is the shadcn/ui Item: a flexible container to display
// content (media, title, description, actions) in lists.
//
//	item.Item(item.VariantOutline,
//		item.Content(item.Title("Basic Item"), item.Description("A simple item.")),
//		item.Actions(button.Button(button.VariantOutline, button.SizeSm, "Action")),
//	)
package item

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/styled"
	"shadcn.go/components/ui/separator"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Variant selects an Item's look.
type Variant string

// Item variants.
const (
	VariantDefault Variant = "default"
	VariantOutline Variant = "outline"
	VariantMuted   Variant = "muted"
)

// Size selects an Item's size.
type Size string

// Item sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
)

// MediaVariant selects the look of Media.
type MediaVariant string

// Media variants.
const (
	MediaDefault MediaVariant = "default"
	MediaIcon    MediaVariant = "icon"
	MediaImage   MediaVariant = "image"
)

var itemVariants = cva.New(
	"group/item flex flex-wrap items-center rounded-md border border-transparent text-sm transition-colors duration-100 outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 [a]:transition-colors [a]:hover:bg-accent/50",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default": "bg-transparent",
				"outline": "border-border",
				"muted":   "bg-muted/50",
			},
			"size": {
				"default": "gap-4 p-4",
				"sm":      "gap-2.5 px-4 py-3",
			},
		},
		Defaults: cva.Selection{"variant": "default", "size": "default"},
		Order:    []string{"variant", "size"},
	},
)

var mediaVariants = cva.New(
	"flex shrink-0 items-center justify-center gap-2 group-has-[[data-slot=item-description]]/item:translate-y-0.5 group-has-[[data-slot=item-description]]/item:self-start [&_svg]:pointer-events-none",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default": "bg-transparent",
				"icon":    "size-8 rounded-sm border bg-muted [&_svg:not([class*='size-'])]:size-4",
				"image":   "size-10 overflow-hidden rounded-sm [&_img]:size-full [&_img]:object-cover",
			},
		},
		Defaults: cva.Selection{"variant": "default"},
	},
)

// Variants returns the classes of an Item, like itemVariants() upstream.
func Variants(v Variant, s Size, className ...string) string {
	return lib.Cn(itemVariants.Classes(cva.Selection{"variant": string(v), "size": string(s)}, className...))
}

// MediaVariants returns the classes of Media, like itemMediaVariants().
func MediaVariants(v MediaVariant, className ...string) string {
	return lib.Cn(mediaVariants.Classes(cva.Selection{"variant": string(v)}, className...))
}

func slot(tag, name, class string) *h.Element {
	return h.El(tag, h.Data("slot", name), h.Class(lib.Cn(class)))
}

var (
	groupBase       = h.Div(h.Role("list"), h.Data("slot", "item-group"), h.Class(lib.Cn("group/item-group flex flex-col")))
	contentBase     = slot("div", "item-content", "flex flex-1 flex-col gap-1 [&+[data-slot=item-content]]:flex-none")
	titleBase       = slot("div", "item-title", "flex w-fit items-center gap-2 text-sm leading-snug font-medium")
	descriptionBase = h.P(h.Data("slot", "item-description"), h.Class(lib.Cn(
		"line-clamp-2 text-sm leading-normal font-normal text-balance text-muted-foreground",
		"[&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary")))
	actionsBase = slot("div", "item-actions", "flex items-center gap-2")
	headerBase  = slot("div", "item-header", "flex basis-full items-center justify-between gap-2")
	footerBase  = slot("div", "item-footer", "flex basis-full items-center justify-between gap-2")
)

func spread(base *h.Element, children []any) *h.Element {
	return h.Spread(base, h.El(base.Tag(), children...), lib.Cn)
}

// Group is a list of items (upstream ItemGroup, role="list").
func Group(children ...any) *styled.Node {
	return styledSpread(groupBase, children, styled.Slot("div", "item-group", sGroup, h.Role("list")))
}

// Separator separates items of a Group (upstream ItemSeparator). It is
// horizontal unless given a separator.Orientation, which, like upstream's
// {...props}, overrides the default.
func Separator(children ...any) *styled.Node {
	orientation := separator.Horizontal
	rest := make([]any, 0, len(children))
	for _, c := range children {
		if o, ok := c.(separator.Orientation); ok {
			orientation = o
			continue
		}
		rest = append(rest, c)
	}
	props := h.Div(rest...)
	className := props.ClassName()
	alt := props.Clone()
	props.SetClass(lib.Cn("my-0", className))
	def := separator.Separator(h.Data("slot", "item-separator"), orientation, spreadProp(props))
	return styled.New(def, func(st style.Style) h.Node {
		p := alt.Clone().SetClass(sSeparator.Class(st, className))
		return separator.Separator(h.Data("slot", "item-separator"), orientation, spreadProp(p))
	})
}

// Content holds the title and description.
func Content(children ...any) *styled.Node {
	return styledSpread(contentBase, children, styled.Slot("div", "item-content", sContent))
}

// Title is the item's title.
func Title(children ...any) *styled.Node {
	return styledSpread(titleBase, children, styled.Slot("div", "item-title", sTitle))
}

// Description is the item's description (a <p>, clamped to two lines).
func Description(children ...any) *styled.Node {
	return styledSpread(descriptionBase, children, styled.Slot("p", "item-description", sDescription))
}

// Actions holds buttons or other controls.
func Actions(children ...any) *styled.Node {
	return styledSpread(actionsBase, children, styled.Slot("div", "item-actions", sActions))
}

// Header spans the item's full width above its content.
func Header(children ...any) *styled.Node {
	return styledSpread(headerBase, children, styled.Slot("div", "item-header", sHeader))
}

// Footer spans the item's full width below its content.
func Footer(children ...any) *styled.Node {
	return styledSpread(footerBase, children, styled.Slot("div", "item-footer", sFooter))
}

// Builder is returned by Item and Media.
type Builder struct {
	h.Builder[Builder]
	media        bool
	variant      Variant
	size         Size
	mediaVariant MediaVariant
	asChild      bool
}

// Item returns an item <div>. Children may include a Variant, a Size and
// lib.AsChild (to render e.g. an <a> instead).
func Item(children ...any) *Builder {
	b := &Builder{variant: VariantDefault, size: SizeDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Media holds an icon, avatar or image. Children may include a
// MediaVariant.
func Media(children ...any) *Builder {
	b := &Builder{media: true, mediaVariant: MediaDefault}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Variant, Size and lib.AsChild (Item) and
// MediaVariant (Media).
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Variant:
		if b.media {
			return false
		}
		b.Variant(v)
	case Size:
		if b.media {
			return false
		}
		b.Size(v)
	case lib.AsChildFlag:
		if b.media {
			return false
		}
		b.asChild = true
	case MediaVariant:
		if !b.media {
			return false
		}
		if v == "" {
			v = MediaDefault
		}
		b.mediaVariant = v
	default:
		return false
	}
	return true
}

// Variant sets the Item's variant.
func (b *Builder) Variant(v Variant) *Builder {
	if v == "" {
		v = VariantDefault
	}
	b.variant = v
	return b
}

// Size sets the Item's size.
func (b *Builder) Size(s Size) *Builder {
	if s == "" {
		s = SizeDefault
	}
	b.size = s
	return b
}

// AsChild makes the Item render its single element child.
func (b *Builder) AsChild() *Builder { b.asChild = true; return b }

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the component to its element as it renders inside w
// (in the style of an enclosing style.Scope).
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	if b.media {
		class := MediaVariants(b.mediaVariant, className)
		if st != "" {
			class = lib.Cn(sMediaVariants.Classes(st, cva.Selection{"variant": string(b.mediaVariant)}, className))
		}
		base := h.Div(h.Data("slot", "item-media"), h.Data("variant", string(b.mediaVariant)), h.Class(class))
		return h.Spread(base, props, lib.Cn)
	}
	class := Variants(b.variant, b.size, className)
	if st != "" {
		class = lib.Cn(sItemVariants.Classes(st, cva.Selection{"variant": string(b.variant), "size": string(b.size)}, className))
	}
	base := h.Div(h.Data("slot", "item"), h.Data("variant", string(b.variant)), h.Data("size", string(b.size)), h.Class(class))
	el := h.Spread(base, props, lib.Cn)
	if !b.asChild {
		return el
	}
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("item: AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := h.SlotIn(w, el.Clone().SetChildren(), kids[0], slotClassMerger(kids[0]))
	if err != nil {
		return h.Errored(fmt.Errorf("item: AsChild: %w", err))
	}
	return out
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error {
	return b.ElementIn(w).Render(w)
}

// slotClassMerger returns how Radix Slot's merged className ends up on the
// rendered child: concatenated for a plain element; for a shadcn component
// (it carries a data-slot), which runs cn(its base classes, merged), the
// Slot's classes win conflicts.
func slotClassMerger(child any) func(...string) string {
	var el *h.Element
	switch c := child.(type) {
	case *h.Element:
		el = c
	case h.Elementer:
		el = c.Element()
	}
	if el == nil || !el.Has("data-slot") {
		return nil
	}
	return func(classes ...string) string { return lib.Cn(classes[1], classes[0]) }
}

// spreadProp applies props' attributes, classes and children to the
// element it is passed to, like {...props}.
func spreadProp(props *h.Element) h.Prop {
	return func(e *h.Element) { *e = *h.Spread(e, props, nil) }
}
