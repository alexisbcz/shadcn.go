// Package carousel is the shadcn/ui Carousel: a carousel with motion and
// swipe built on Embla.
//
//	carousel.Carousel(carousel.Opts{Align: "start"}, h.Class("w-full max-w-xs"),
//	    carousel.Content(
//	        carousel.Item(...),
//	        carousel.Item(...),
//	    ),
//	    carousel.Previous(),
//	    carousel.Next(),
//	)
//
// Root children may include an Orientation, Opts (Embla's options) and an
// Autoplay plugin (embla-carousel-autoplay's options). Previous and Next
// accept a button.Variant and button.Size like the upstream Button.
//
// Behavior (carousel.js): Embla Carousel 8.5.2 itself (vendored unmodified in
// runtime/js/embla-carousel.js) runs on the viewport, exactly like
// useEmblaCarousel; ArrowLeft/ArrowRight on the root scroll to the previous
// or next snap; the Previous/Next buttons are disabled when Embla cannot
// scroll that way. Events dispatched on the root: "shadcn:carousel-init"
// (detail: { api }, the Embla API) once Embla is running, and
// "shadcn:carousel-select" (detail: { api, index, count }) whenever the
// selected snap changes. From JavaScript, import { onCarouselApi } from
// "/_shadcn/carousel.js" to get the Embla API of a carousel element.
package carousel

import (
	_ "embed"
	"encoding/json"
	"io"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed carousel.js
var script string

func init() { runtime.Module("carousel", script) }

// Orientation of the carousel (default horizontal).
type Orientation string

// Orientations.
const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// Opts are Embla's options (the opts prop). Zero values keep Embla's
// defaults. The axis follows the Orientation.
type Opts struct {
	Align          string `json:"align,omitempty"` // "start", "center" (default) or "end"
	Loop           bool   `json:"loop,omitempty"`
	DragFree       bool   `json:"dragFree,omitempty"`
	SkipSnaps      bool   `json:"skipSnaps,omitempty"`
	StartIndex     int    `json:"startIndex,omitempty"`
	Duration       int    `json:"duration,omitempty"`
	SlidesToScroll any    `json:"slidesToScroll,omitempty"` // a number or "auto"
	ContainScroll  any    `json:"containScroll,omitempty"`  // "trimSnaps" (default), "keepSnaps" or false
	Direction      string `json:"direction,omitempty"`      // "ltr" or "rtl"
	DisableDrag    bool   `json:"-"`                        // watchDrag: false
}

// MarshalJSON adds watchDrag when dragging is disabled.
func (o Opts) MarshalJSON() ([]byte, error) {
	type plain Opts
	m := map[string]any{}
	raw, err := json.Marshal(plain(o))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if o.ContainScroll == false {
		m["containScroll"] = false
	}
	if o.DisableDrag {
		m["watchDrag"] = false
	}
	return json.Marshal(m)
}

// Autoplay adds the embla-carousel-autoplay plugin.
type Autoplay struct {
	// Delay between slides in milliseconds (default 4000).
	Delay int
	// KeepPlayingOnInteraction sets stopOnInteraction: false (upstream's
	// default is true: a drag stops autoplay for good).
	KeepPlayingOnInteraction bool
	// StopOnMouseEnter, StopOnLastSnap and Jump mirror the plugin options.
	StopOnMouseEnter bool
	StopOnLastSnap   bool
	Jump             bool
	// StopOnHover wires the root's mouseenter to the plugin's stop() and its
	// mouseleave to reset(), as the upstream carousel-plugin example does with
	// onMouseEnter={plugin.stop} onMouseLeave={plugin.reset}.
	StopOnHover bool
}

func (a Autoplay) options() map[string]any {
	m := map[string]any{}
	if a.Delay > 0 {
		m["delay"] = a.Delay
	}
	if a.KeepPlayingOnInteraction {
		m["stopOnInteraction"] = false
	}
	if a.StopOnMouseEnter {
		m["stopOnMouseEnter"] = true
	}
	if a.StopOnLastSnap {
		m["stopOnLastSnap"] = true
	}
	if a.Jump {
		m["jump"] = true
	}
	return m
}

type context struct{ orientation Orientation }

// Builder is returned by Carousel.
type Builder struct {
	h.Builder[Builder]
	orientation Orientation
	opts        *Opts
	autoplay    *Autoplay
}

// Carousel returns the carousel's root <div role="region">.
func Carousel(children ...any) *Builder {
	b := &Builder{orientation: OrientationHorizontal}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Orientation, Opts and Autoplay children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Orientation:
		b.orientation = v
	case Opts:
		b.opts = &v
	case Autoplay:
		b.autoplay = &v
	default:
		return false
	}
	return true
}

// Orientation sets the orientation.
func (b *Builder) Orientation(o Orientation) *Builder { b.orientation = o; return b }

// Opts sets Embla's options.
func (b *Builder) Opts(o Opts) *Builder { b.opts = &o; return b }

// Autoplay adds the autoplay plugin.
func (b *Builder) Autoplay(a Autoplay) *Builder { b.autoplay = &a; return b }

func (b *Builder) config() string {
	orientation := b.orientation
	if orientation == "" {
		orientation = OrientationHorizontal
	}
	cfg := map[string]any{"axis": "x"}
	if orientation == OrientationVertical {
		cfg["axis"] = "y"
	}
	if b.opts != nil {
		cfg["opts"] = b.opts
	}
	if b.autoplay != nil {
		cfg["autoplay"] = b.autoplay.options()
		if b.autoplay.StopOnHover {
			cfg["stopOnHover"] = true
		}
	}
	out, _ := json.Marshal(cfg)
	return string(out)
}

// Element resolves the carousel to its root <div>.
func (b *Builder) Element() *h.Element {
	orientation := b.orientation
	if orientation == "" {
		orientation = OrientationHorizontal
	}
	ctx := &context{orientation: orientation}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ContentBuilder:
			c.ctx = ctx
		case *ItemBuilder:
			c.ctx = ctx
		case *PreviousBuilder:
			c.ctx = ctx
		case *NextBuilder:
			c.ctx = ctx
		case *Builder:
			return false
		}
		return true
	})
	base := h.Div(
		h.Class("relative"),
		h.Role("region"),
		h.Aria("roledescription", "carousel"),
		h.Data("slot", "carousel"),
		// Read and removed by carousel.js on mount.
		h.Data("carousel", b.config()),
	)
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the carousel.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

func orientationOf(c *context) Orientation {
	if c == nil {
		return OrientationHorizontal
	}
	return c.orientation
}

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	ctx *context
}

// Content returns the viewport <div> holding the slide container; its
// children are the slides.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the content to the viewport <div>.
func (b *ContentBuilder) Element() *h.Element {
	cls := "-ml-4"
	if orientationOf(b.ctx) == OrientationVertical {
		cls = "-mt-4 flex-col"
	}
	props := b.Elem()
	inner := h.Spread(h.Div(), props, nil).SetClass(lib.Cn("flex", cls, props.ClassName()))
	return h.Div(h.Class("overflow-hidden"), h.Data("slot", "carousel-content"), inner)
}

// Render renders the content.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	ctx *context
}

// Item returns a slide <div role="group">.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the item to its <div>.
func (b *ItemBuilder) Element() *h.Element {
	cls := "pl-4"
	if orientationOf(b.ctx) == OrientationVertical {
		cls = "pt-4"
	}
	base := h.Div(h.Role("group"), h.Aria("roledescription", "slide"), h.Data("slot", "carousel-item"))
	props := b.Elem()
	return h.Spread(base, props, nil).SetClass(lib.Cn("min-w-0 shrink-0 grow-0 basis-full", cls, props.ClassName()))
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// navButton renders Previous and Next.
func navButton(ctx *context, props *h.Element, variant button.Variant, size button.Size, slot, horizontal, vertical string, icon h.Node, label string) *h.Element {
	cls := horizontal
	if orientationOf(ctx) == OrientationVertical {
		cls = vertical
	}
	base := button.Button(variant, size,
		h.Data("slot", slot),
		h.Class(lib.Cn("absolute size-8 rounded-full", cls, props.ClassName())),
		h.Disabled(), // canScroll* is false on React's first render
		icon,
		h.Span(h.Class("sr-only"), label),
	).Element()
	userProps := props.Clone().SetClass("")
	out := h.Spread(base, userProps.SetChildren(), nil)
	if kids := props.ChildNodes(); len(kids) > 0 {
		out.SetChildren(kids...) // children passed in props replace the icon and label
	}
	return out
}

// PreviousBuilder is returned by Previous.
type PreviousBuilder struct {
	h.Builder[PreviousBuilder]
	ctx     *context
	variant button.Variant
	size    button.Size
}

// Previous returns the button that scrolls to the previous slide.
func Previous(children ...any) *PreviousBuilder {
	b := &PreviousBuilder{variant: button.VariantOutline, size: button.SizeIcon}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes button.Variant and button.Size children.
func (b *PreviousBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case button.Variant:
		b.variant = v
	case button.Size:
		b.size = v
	default:
		return false
	}
	return true
}

// Element resolves the button.
func (b *PreviousBuilder) Element() *h.Element {
	return navButton(b.ctx, b.Elem(), b.variant, b.size, "carousel-previous",
		"top-1/2 -left-12 -translate-y-1/2", "-top-12 left-1/2 -translate-x-1/2 rotate-90",
		lucide.ArrowLeft(), "Previous slide")
}

// Render renders the button.
func (b *PreviousBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// NextBuilder is returned by Next.
type NextBuilder struct {
	h.Builder[NextBuilder]
	ctx     *context
	variant button.Variant
	size    button.Size
}

// Next returns the button that scrolls to the next slide.
func Next(children ...any) *NextBuilder {
	b := &NextBuilder{variant: button.VariantOutline, size: button.SizeIcon}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes button.Variant and button.Size children.
func (b *NextBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case button.Variant:
		b.variant = v
	case button.Size:
		b.size = v
	default:
		return false
	}
	return true
}

// Element resolves the button.
func (b *NextBuilder) Element() *h.Element {
	return navButton(b.ctx, b.Elem(), b.variant, b.size, "carousel-next",
		"top-1/2 -right-12 -translate-y-1/2", "-bottom-12 left-1/2 -translate-x-1/2 rotate-90",
		lucide.ArrowRight(), "Next slide")
}

// Render renders the button.
func (b *NextBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
