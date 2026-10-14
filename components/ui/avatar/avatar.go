// Package avatar is the shadcn/ui Avatar (Radix Avatar): an image element
// with a fallback for representing the user.
//
//	avatar.Avatar(
//	    avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn")),
//	    avatar.Fallback("CN"),
//	)
//
// Like Radix, the image is rendered only once it has loaded, and the
// fallback only while it has not (and, with DelayMs, once the delay has
// passed). The server renders the fallback (or nothing, with a delay); the
// behavior (avatar.js) loads the image and swaps them. The image markup
// travels in a <template data-slot="avatar-image-template">, so nothing is
// fetched before the behavior runs.
//
// The "shadcn:loading-status-change" event (detail: { status: "loading" |
// "loaded" | "error" }) is dispatched on the image's avatar when its image's
// loading status changes (Radix onLoadingStatusChange).
package avatar

import (
	_ "embed"
	"fmt"
	"io"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/part"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed avatar.js
var script string

func init() { runtime.Module("avatar", script) }

// Size selects the avatar's size.
type Size string

// Avatar sizes.
const (
	SizeDefault Size = "default"
	SizeSm      Size = "sm"
	SizeLg      Size = "lg"
)

// DelayMs delays rendering the fallback (Radix AvatarFallback delayMs), to
// avoid flashing it for images that load fast.
type DelayMs int

// Builder is returned by Avatar.
type Builder struct {
	h.Builder[Builder]
	size Size
}

// Avatar returns the avatar's root <span>. Children may include a Size.
func Avatar(children ...any) *Builder {
	b := &Builder{size: SizeDefault}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// AcceptChild consumes Size children.
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

// Element resolves the avatar to its <span>.
func (b *Builder) Element() *h.Element { return b.element("") }

// element renders in style st ("" for new-york-v4).
func (b *Builder) element(st style.Style) *h.Element {
	base := h.Span(h.Data("slot", "avatar"), h.Data("size", string(b.size)))
	props := b.Elem()
	if st != "" {
		return h.Spread(base, props, nil).SetClass(styleRoot.Class(st, props.ClassName()))
	}
	return h.Spread(base, props, nil).SetClass(lib.Cn("group/avatar relative flex size-8 shrink-0 overflow-hidden rounded-full select-none data-[size=lg]:size-10 data-[size=sm]:size-6", props.ClassName()))
}

// Render renders the avatar.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	return b.element(st)
}

// ImageBuilder is returned by Image.
type ImageBuilder struct {
	h.Builder[ImageBuilder]
}

// Image returns the avatar's image: an <img> that is shown once it has
// loaded. Pass h.Src and h.Alt (and h.Attr("referrerpolicy", ...) or
// h.CrossOrigin, which also apply to the loading request).
func Image(children ...any) *ImageBuilder {
	b := &ImageBuilder{}
	b.Builder = h.NewBuilder(b, "img")
	return b.Children(children...)
}

// Element resolves the image to the <img> it shows once loaded.
func (b *ImageBuilder) Element() *h.Element { return b.element("") }

func (b *ImageBuilder) element(st style.Style) *h.Element {
	props := b.Elem()
	if st != "" {
		return h.Spread(h.Img(h.Data("slot", "avatar-image")), props, nil).SetClass(styleImage.Class(st, props.ClassName()))
	}
	return h.Spread(h.Img(h.Data("slot", "avatar-image")), props, nil).SetClass(lib.Cn("aspect-square size-full", props.ClassName()))
}

// Render writes the image's template: Radix renders nothing until the
// image has loaded.
func (b *ImageBuilder) Render(w io.Writer) error {
	st, _ := style.FromWriter(w)
	el := b.element(st)
	if el.Err() != nil {
		return el.Render(w)
	}
	return h.Template(h.Data("slot", "avatar-image-template"), el).Render(w)
}

// FallbackBuilder is returned by Fallback.
type FallbackBuilder struct {
	h.Builder[FallbackBuilder]
	delay *int
}

// Fallback returns the avatar's fallback <span>, shown while the image is
// not loaded. Children may include a DelayMs.
func Fallback(children ...any) *FallbackBuilder {
	b := &FallbackBuilder{}
	b.Builder = h.NewBuilder(b, "span")
	return b.Children(children...)
}

// AcceptChild consumes DelayMs children.
func (b *FallbackBuilder) AcceptChild(c any) bool {
	if d, ok := c.(DelayMs); ok {
		b.DelayMs(int(d))
		return true
	}
	return false
}

// DelayMs delays rendering the fallback.
func (b *FallbackBuilder) DelayMs(ms int) *FallbackBuilder { b.delay = &ms; return b }

// Element resolves the fallback to its <span>.
func (b *FallbackBuilder) Element() *h.Element { return b.element("") }

func (b *FallbackBuilder) element(st style.Style) *h.Element {
	props := b.Elem()
	if st != "" {
		return h.Spread(h.Span(h.Data("slot", "avatar-fallback")), props, nil).SetClass(styleFallback.Class(st, props.ClassName()))
	}
	return h.Spread(h.Span(h.Data("slot", "avatar-fallback")), props, nil).SetClass(lib.Cn("flex size-full items-center justify-center rounded-full bg-muted text-sm text-muted-foreground group-data-[size=sm]/avatar:text-xs", props.ClassName()))
}

// Render renders the fallback, or, with a delay, its template.
func (b *FallbackBuilder) Render(w io.Writer) error {
	st, _ := style.FromWriter(w)
	el := b.element(st)
	if el.Err() != nil || b.delay == nil {
		return el.Render(w)
	}
	if *b.delay < 0 {
		return h.Errored(fmt.Errorf("avatar: negative DelayMs %d", *b.delay)).Render(w)
	}
	return h.Template(h.Data("slot", "avatar-fallback-template"), h.Data("delay", strconv.Itoa(*b.delay)), el).Render(w)
}

// styleFile is registry/bases/radix/ui/avatar.tsx (style.Scope).
var styleFile = style.NewFile("ui/avatar.tsx")

var (
	styleRoot     = styleFile.Cn("cn-avatar group/avatar relative flex shrink-0 select-none after:absolute after:inset-0 after:border after:border-border after:mix-blend-darken dark:after:mix-blend-lighten", style.Arg)
	styleImage    = styleFile.Cn("cn-avatar-image aspect-square size-full object-cover", style.Arg)
	styleFallback = styleFile.Cn("cn-avatar-fallback flex size-full items-center justify-center text-sm group-data-[size=sm]/avatar:text-xs", style.Arg)
	styleBadge    = styleFile.Cn(
		"cn-avatar-badge absolute right-0 bottom-0 z-10 inline-flex items-center justify-center rounded-full bg-blend-color ring-2 select-none",
		"group-data-[size=sm]/avatar:size-2 group-data-[size=sm]/avatar:[&>svg]:hidden",
		"group-data-[size=default]/avatar:size-2.5 group-data-[size=default]/avatar:[&>svg]:size-2",
		"group-data-[size=lg]/avatar:size-3 group-data-[size=lg]/avatar:[&>svg]:size-2",
		style.Arg,
	)
	styleGroup      = styleFile.Cn("cn-avatar-group group/avatar-group flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:ring-background", style.Arg)
	styleGroupCount = styleFile.Cn("cn-avatar-group-count relative flex shrink-0 items-center justify-center ring-2 ring-background", "", style.Arg)
)

// Badge returns a <span> badge in the avatar's bottom-right corner.
func Badge(children ...any) *part.Part {
	return part.New("span", func(st style.Style, props *h.Element) *h.Element {
		e := h.Spread(h.Span(h.Data("slot", "avatar-badge")), props, nil)
		if st != "" {
			return e.SetClass(styleBadge.Class(st, props.ClassName()))
		}
		return e.SetClass(lib.Cn(
			"absolute right-0 bottom-0 z-10 inline-flex items-center justify-center rounded-full bg-primary text-primary-foreground ring-2 ring-background select-none",
			"group-data-[size=sm]/avatar:size-2 group-data-[size=sm]/avatar:[&>svg]:hidden",
			"group-data-[size=default]/avatar:size-2.5 group-data-[size=default]/avatar:[&>svg]:size-2",
			"group-data-[size=lg]/avatar:size-3 group-data-[size=lg]/avatar:[&>svg]:size-2",
			props.ClassName(),
		))
	}, children...)
}

// Group returns a <div> stacking avatars.
func Group(children ...any) *part.Part {
	return part.New("div", func(st style.Style, props *h.Element) *h.Element {
		e := h.Spread(h.Div(h.Data("slot", "avatar-group")), props, nil)
		if st != "" {
			return e.SetClass(styleGroup.Class(st, props.ClassName()))
		}
		return e.SetClass(lib.Cn("group/avatar-group flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:ring-background", props.ClassName()))
	}, children...)
}

// GroupCount returns a <div> showing the number of avatars not shown in a
// Group.
func GroupCount(children ...any) *part.Part {
	return part.New("div", func(st style.Style, props *h.Element) *h.Element {
		e := h.Spread(h.Div(h.Data("slot", "avatar-group-count")), props, nil)
		if st != "" {
			return e.SetClass(styleGroupCount.Class(st, props.ClassName()))
		}
		return e.SetClass(lib.Cn("relative flex size-8 shrink-0 items-center justify-center rounded-full bg-muted text-sm text-muted-foreground ring-2 ring-background group-has-data-[size=lg]/avatar-group:size-10 group-has-data-[size=sm]/avatar-group:size-6 [&>svg]:size-4 group-has-data-[size=lg]/avatar-group:[&>svg]:size-5 group-has-data-[size=sm]/avatar-group:[&>svg]:size-3", props.ClassName()))
	}, children...)
}
