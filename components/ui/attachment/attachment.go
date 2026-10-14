// Package attachment is the shadcn/ui Attachment: a file attached to a
// message or a prompt, with its upload state, media, details and actions.
//
//	attachment.Attachment(attachment.StateUploading,
//		attachment.Media(spinner.Spinner()),
//		attachment.Content(attachment.Title("design.zip"), attachment.Description("64%")),
//		attachment.Actions(attachment.Action(h.Aria("label", "Cancel"), lucide.X())),
//	)
package attachment

import (
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// State is the attachment's upload state.
type State string

// States.
const (
	StateIdle       State = "idle"
	StateUploading  State = "uploading"
	StateProcessing State = "processing"
	StateError      State = "error"
	StateDone       State = "done"
)

// Size selects the attachment's size.
type Size string

// Sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
	SizeXs      Size = "xs"
)

// Orientation lays the attachment out in a row or a column.
type Orientation string

// Orientations.
const (
	Horizontal Orientation = "horizontal"
	Vertical   Orientation = "vertical"
)

// MediaVariant selects how Media renders its content.
type MediaVariant string

// Media variants.
const (
	MediaIcon  MediaVariant = "icon"
	MediaImage MediaVariant = "image"
)

var (
	StateValues       = []State{StateIdle, StateUploading, StateProcessing, StateError, StateDone}
	SizeValues        = []Size{SizeDefault, SizeSm, SizeXs}
	OrientationValues = []Orientation{Horizontal, Vertical}
)

var attachmentVariants = cva.New(
	"group/attachment relative flex w-fit max-w-full min-w-0 shrink-0 flex-wrap rounded-xl border bg-card text-card-foreground transition-colors focus-within:ring-1 focus-within:ring-ring/50 has-[>a,>button]:hover:bg-muted/50 data-[state=error]:border-destructive/30 data-[state=idle]:border-dashed",
	cva.Config{
		Variants: cva.Variants{
			"size": {
				"default": "gap-2 text-sm has-data-[slot=attachment-content]:px-2.5 has-data-[slot=attachment-content]:py-2 has-data-[slot=attachment-media]:p-2",
				"sm":      "gap-2.5 text-xs has-data-[slot=attachment-content]:px-2 has-data-[slot=attachment-content]:py-1.5 has-data-[slot=attachment-media]:p-1.5",
				"xs":      "gap-1.5 rounded-lg text-xs has-data-[slot=attachment-content]:px-1.5 has-data-[slot=attachment-content]:py-1 has-data-[slot=attachment-media]:p-1",
			},
			"orientation": {
				"horizontal": "min-w-40 items-center",
				"vertical":   "w-24 flex-col has-data-[slot=attachment-content]:w-30",
			},
		},
		Order: []string{"size", "orientation"},
	},
)

var mediaVariants = cva.New(
	"relative flex aspect-square w-10 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-muted text-foreground group-data-[orientation=vertical]/attachment:w-full group-data-[size=sm]/attachment:w-8 group-data-[size=xs]/attachment:w-7 group-data-[size=xs]/attachment:rounded-md group-data-[state=error]/attachment:bg-destructive/10 group-data-[state=error]/attachment:text-destructive group-data-[orientation=vertical]/attachment:*:data-[slot=spinner]:size-6! [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 group-data-[orientation=vertical]/attachment:[&_svg:not([class*='size-'])]:size-6 group-data-[size=xs]/attachment:[&_svg:not([class*='size-'])]:size-3.5",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"icon":  "",
				"image": "opacity-60 group-data-[state=done]/attachment:opacity-100 group-data-[state=idle]/attachment:opacity-100 *:[img]:aspect-square *:[img]:w-full *:[img]:object-cover",
			},
		},
		Defaults: cva.Selection{"variant": "icon"},
		Order:    []string{"variant"},
	},
)

// Variants returns the classes of an attachment, like
// cn(attachmentVariants({ size, orientation }), className).
func Variants(s Size, o Orientation, className ...string) string {
	return lib.Cn(append([]string{attachmentVariants.Classes(cva.Selection{"size": string(s), "orientation": string(o)})}, className...)...)
}

// MediaVariants returns the classes of Media, like
// cn(attachmentMediaVariants({ variant }), className).
func MediaVariants(v MediaVariant, className ...string) string {
	return lib.Cn(append([]string{mediaVariants.Classes(cva.Selection{"variant": string(v)})}, className...)...)
}

// Builder is returned by Attachment.
type Builder struct {
	h.Builder[Builder]
	state       State
	size        Size
	orientation Orientation
}

// Attachment returns an attachment card. Children may include a State
// (default done), a Size and an Orientation (default horizontal).
func Attachment(children ...any) *Builder {
	b := &Builder{state: StateDone, size: SizeDefault, orientation: Horizontal}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes State, Size and Orientation.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case State:
		if v == "" {
			v = StateDone
		}
		b.state = v
	case Size:
		if v == "" {
			v = SizeDefault
		}
		b.size = v
	case Orientation:
		if v == "" {
			v = Horizontal
		}
		b.orientation = v
	default:
		return false
	}
	return true
}

// Element resolves the attachment.
func (b *Builder) Element() *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Div(
		h.Data("slot", "attachment"),
		h.Data("state", string(b.state)),
		h.Data("size", string(b.size)),
		h.Data("orientation", string(b.orientation)),
		h.Class(Variants(b.size, b.orientation, className)),
	)
	return h.Spread(base, props, nil)
}

// Render renders the attachment.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

// MediaBuilder is returned by Media.
type MediaBuilder struct {
	h.Builder[MediaBuilder]
	variant MediaVariant
}

// Media holds an icon, a spinner or an image (MediaImage).
func Media(children ...any) *MediaBuilder {
	b := &MediaBuilder{variant: MediaIcon}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes MediaVariant.
func (b *MediaBuilder) AcceptChild(c any) bool {
	if v, ok := c.(MediaVariant); ok {
		if v == "" {
			v = MediaIcon
		}
		b.variant = v
		return true
	}
	return false
}

// Element resolves the media.
func (b *MediaBuilder) Element() *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Div(h.Data("slot", "attachment-media"), h.Data("variant", string(b.variant)), h.Class(MediaVariants(b.variant, className)))
	return h.Spread(base, props, nil)
}

// Render renders the media.
func (b *MediaBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

var (
	contentBase     = h.Div(h.Data("slot", "attachment-content"), h.Class(lib.Cn("max-w-full min-w-0 flex-1 leading-tight group-data-[orientation=vertical]/attachment:px-1")))
	titleBase       = h.Span(h.Data("slot", "attachment-title"), h.Class(lib.Cn("block max-w-full min-w-0 truncate font-medium group-data-[state=processing]/attachment:shimmer group-data-[state=uploading]/attachment:shimmer")))
	descriptionBase = h.Span(h.Data("slot", "attachment-description"), h.Class(lib.Cn(
		"mt-0.5 block min-w-0 truncate text-xs text-muted-foreground group-data-[state=error]/attachment:text-destructive/80",
		"max-w-full")))
	actionsBase = h.Div(h.Data("slot", "attachment-actions"), h.Class(lib.Cn("relative z-20 flex shrink-0 items-center group-data-[orientation=vertical]/attachment:absolute group-data-[orientation=vertical]/attachment:top-3 group-data-[orientation=vertical]/attachment:right-3 group-data-[orientation=vertical]/attachment:gap-1")))
	groupBase   = h.Div(h.Data("slot", "attachment-group"), h.Class(lib.Cn("flex min-w-0 scroll-fade-x snap-x snap-mandatory scroll-px-1 scrollbar-none gap-3 overflow-x-auto overscroll-x-contain py-1 *:data-[slot=attachment]:flex-none *:data-[slot=attachment]:snap-start")))
)

func spread(base *h.Element, children []any) *h.Element {
	return h.Spread(base, h.El(base.Tag(), children...), lib.Cn)
}

// Content holds the Title and Description.
func Content(children ...any) *h.Element { return spread(contentBase, children) }

// Title is the file name (a truncated <span>, shimmering while uploading
// or processing).
func Title(children ...any) *h.Element { return spread(titleBase, children) }

// Description is the file's details or status (a truncated <span>).
func Description(children ...any) *h.Element { return spread(descriptionBase, children) }

// Actions holds the Action buttons, above the Trigger.
func Actions(children ...any) *h.Element { return spread(actionsBase, children) }

// Group lays attachments out in a horizontally scrolling, snapping row.
func Group(children ...any) *h.Element { return spread(groupBase, children) }

// Action is a ghost icon-xs button (upstream AttachmentAction). A
// button.Variant or button.Size among children overrides the defaults.
func Action(children ...any) *button.Builder {
	return button.Button(append([]any{h.Data("slot", "attachment-action"), button.VariantGhost, button.SizeIconXs}, children...)...)
}

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	asChild bool
}

// Trigger is a button covering the whole attachment (to open or preview
// it). With lib.AsChild it renders its single element child (a link).
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes lib.AsChild.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	if _, ok := c.(lib.AsChildFlag); ok {
		b.asChild = true
		return true
	}
	return false
}

// AsChild makes the trigger render its single element child.
func (b *TriggerBuilder) AsChild() *TriggerBuilder { b.asChild = true; return b }

// Type sets the button type (default "button").
func (b *TriggerBuilder) Type(v string) *TriggerBuilder { b.Elem().Type(v); return b }

// Element resolves the trigger.
func (b *TriggerBuilder) Element() *h.Element {
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	base := h.Button(h.Data("slot", "attachment-trigger"), h.Class(lib.Cn("absolute inset-0 z-10 outline-none", className)))
	if !b.asChild {
		if _, ok := props.GetAttr("type"); !ok {
			base.Attr("type", "button")
		}
	} else {
		// type={asChild ? undefined : (type ?? "button")}: upstream takes
		// type out of the props, so the slotted child never receives it.
		props.Remove("type")
	}
	el := h.Spread(base, props, nil)
	if !b.asChild {
		return el
	}
	kids := el.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("attachment: Trigger AsChild needs exactly one element child, got %d children", len(kids)))
	}
	out, err := h.Slot(el.Clone().SetChildren(), kids[0], nil)
	if err != nil {
		return h.Errored(fmt.Errorf("attachment: Trigger AsChild: %w", err))
	}
	return out
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
