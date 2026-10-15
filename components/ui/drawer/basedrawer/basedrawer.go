// Package basedrawer is the Drawer of the Base UI styles
// (styles/base-nova/ui/drawer.tsx, base-rhea and aria-nova): Base UI's
// Drawer, whose look is driven by data attributes (data-open,
// data-starting-style, data-swipe-direction, ...) and CSS variables
// (--nested-drawers, --drawer-swipe-movement-y, ...) that basedrawer.js
// maintains like Base UI does: open and close transitions, nested drawer
// stacking, swipe to dismiss, snap points, non-modal drawers.
//
//	basedrawer.Drawer(basedrawer.SwipeDirection("right"),
//		basedrawer.Trigger(button.Button(button.VariantOutline, "Open")),
//		basedrawer.Content(
//			basedrawer.Header(basedrawer.Title("Title"), basedrawer.Description("...")),
//			basedrawer.Footer(basedrawer.Close(button.Button("Close"))),
//		),
//	)
//
// The classes follow the style in scope: style.Rhea gives base-rhea's,
// anything else base-nova's (aria-nova's are the same).
package basedrawer

import (
	_ "embed"
	"encoding/json"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/aschild"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed basedrawer.js
var script string

func init() { runtime.Module("basedrawer", script) }

// config is the root's state, read by basedrawer.js from the portal
// template.
type config struct {
	SwipeDirection          string `json:"swipeDirection"`
	MobileSwipeDirection    string `json:"mobileSwipeDirection,omitempty"`
	ShowSwipeHandle         bool   `json:"showSwipeHandle,omitempty"`
	MobileShowSwipeHandle   bool   `json:"mobileShowSwipeHandle,omitempty"`
	Modal                   bool   `json:"modal"`
	DisablePointerDismissal bool   `json:"disablePointerDismissal,omitempty"`
	SnapPoints              []any  `json:"snapPoints,omitempty"`
}

// Option configures a Drawer; pass options as children of Drawer.
type Option func(*config)

// SwipeDirection is the edge the drawer is dismissed towards: "down"
// (default), "up", "left" or "right".
func SwipeDirection(d string) Option { return func(c *config) { c.SwipeDirection = d } }

// MobileSwipeDirection overrides SwipeDirection below the md breakpoint
// (the examples' useIsMobile()).
func MobileSwipeDirection(d string) Option { return func(c *config) { c.MobileSwipeDirection = d } }

// SnapPoints are heights the drawer rests at: "31rem", "200px", or
// fractions of the viewport height (1 = full height).
func SnapPoints(points ...any) Option { return func(c *config) { c.SnapPoints = points } }

var (
	// ShowSwipeHandle renders the swipe handle.
	ShowSwipeHandle Option = func(c *config) { c.ShowSwipeHandle = true }
	// MobileShowSwipeHandle renders the swipe handle below the md
	// breakpoint only.
	MobileShowSwipeHandle Option = func(c *config) { c.MobileShowSwipeHandle = true }
	// NonModal: no backdrop, the page stays interactive (modal={false}).
	NonModal Option = func(c *config) { c.Modal = false }
	// DisablePointerDismissal keeps the drawer open on outside presses.
	DisablePointerDismissal Option = func(c *config) { c.DisablePointerDismissal = true }
)

// RootBuilder is returned by Drawer.
type RootBuilder struct {
	cfg      config
	children []any
}

// Drawer is the root: it renders its children (the Trigger and the
// Content).
func Drawer(children ...any) *RootBuilder {
	b := &RootBuilder{cfg: config{SwipeDirection: "down", Modal: true}}
	for _, c := range children {
		if o, ok := c.(Option); ok {
			o(&b.cfg)
			continue
		}
		b.children = append(b.children, c)
	}
	return b
}

// Render renders the children, the Content with the root's state.
func (b *RootBuilder) Render(w io.Writer) error {
	out := make([]any, 0, len(b.children))
	for _, c := range b.children {
		if cb, ok := c.(*ContentBuilder); ok {
			cb.cfg = &b.cfg
		}
		out = append(out, c)
	}
	return h.Group(out...).Render(w)
}

// slotted renders a child with the part's props (Base UI's render prop).
type slotted struct {
	slot  *h.Element
	child any
}

func (s slotted) Render(w io.Writer) error {
	el, err := aschild.SlotIn(w, s.slot, s.child)
	if err != nil {
		return err
	}
	return el.Render(w)
}

// Trigger renders child (a button) as the drawer's trigger.
func Trigger(child any) h.Node {
	return slotted{h.Button(h.Data("slot", "drawer-trigger"), h.Attr("aria-haspopup", "dialog"), h.Attr("aria-expanded", "false"), h.Data("base-drawer-trigger", "")), child}
}

// Close renders child (a button) as a button closing the drawer.
func Close(child any) h.Node {
	return slotted{h.Button(h.Data("slot", "drawer-close"), h.Data("base-drawer-close", "")), child}
}

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	cfg      *config
	children []any
}

// Content is the drawer's popup (DrawerContent: portal, backdrop, viewport,
// popup, swipe handle and content), rendered into a template mounted on
// open.
func Content(children ...any) *ContentBuilder { return &ContentBuilder{children: children} }

// Render renders the template.
func (b *ContentBuilder) Render(w io.Writer) error {
	cfg := b.cfg
	if cfg == nil {
		cfg = &config{SwipeDirection: "down", Modal: true}
	}
	st, _ := style.FromWriter(w)
	rhea := st == style.Rhea
	raw, _ := json.Marshal(cfg)
	axis := "y"
	if cfg.SwipeDirection == "left" || cfg.SwipeDirection == "right" {
		axis = "x"
	}
	snap := len(cfg.SnapPoints) > 0
	var parts []any
	if cfg.Modal {
		overlay := h.Div(h.Data("slot", "drawer-overlay"), h.Attr("role", "presentation"), h.Class(overlayClass(rhea)))
		if snap {
			overlay.Data("snap-points", "")
		}
		parts = append(parts, overlay)
	}
	popup := h.Div(h.Data("slot", "drawer-popup"), h.Attr("role", "dialog"), h.Attr("tabindex", "-1"),
		h.Data("swipe-axis", axis), h.Data("swipe-direction", cfg.SwipeDirection), h.Class(popupClass(rhea)))
	if snap {
		popup.Data("snap-points", "")
	}
	if cfg.ShowSwipeHandle || cfg.MobileShowSwipeHandle {
		handle := h.Div(h.Data("slot", "drawer-swipe-handle"), h.Attr("aria-hidden", "true"), h.Class(handleClass(rhea)))
		if !cfg.ShowSwipeHandle {
			handle.Data("mobile-only", "")
		}
		popup.Children(handle)
	}
	popup.Children(h.Div(h.Data("drawer-content", ""), h.Data("slot", "drawer-content"), h.Class(contentClass), h.Group(b.children...)))
	modal := "false"
	if cfg.Modal {
		modal = "true"
	}
	parts = append(parts, h.Div(h.Data("slot", "drawer-viewport"), h.Attr("role", "presentation"), h.Data("modal", modal),
		h.Class("pointer-events-none fixed inset-0 z-50 select-none data-[modal=true]:pointer-events-auto"), popup))
	return h.El("template", h.Data("base-drawer", string(raw)), h.Group(parts...)).Render(w)
}

// Header is DrawerHeader.
func Header(children ...any) h.Node {
	return part("drawer-header", func(w io.Writer) string {
		if st, _ := style.FromWriter(w); st == style.Rhea {
			return "flex shrink-0 flex-col gap-0.5 p-4 pb-0 group-data-[swipe-axis=y]/drawer-popup:text-center md:gap-1.5 md:text-left"
		}
		return "flex shrink-0 flex-col gap-0.5 p-4 pb-0 group-data-[swipe-axis=y]/drawer-popup:text-center md:gap-0.5 md:text-left"
	}, h.Div, children)
}

// Footer is DrawerFooter.
func Footer(children ...any) h.Node {
	return part("drawer-footer", func(io.Writer) string { return "mt-auto flex shrink-0 flex-col gap-2 p-4 pt-0" }, h.Div, children)
}

// Title is DrawerTitle.
func Title(children ...any) h.Node {
	return part("drawer-title", func(io.Writer) string { return "cn-font-heading text-base font-medium text-foreground" }, h.H2, children)
}

// Description is DrawerDescription.
func Description(children ...any) h.Node {
	return part("drawer-description", func(io.Writer) string { return "text-sm text-balance text-muted-foreground" }, h.P, children)
}

type partNode struct {
	slot     string
	class    func(io.Writer) string
	tag      func(...any) *h.Element
	children []any
}

func part(slot string, class func(io.Writer) string, tag func(...any) *h.Element, children []any) h.Node {
	return partNode{slot, class, tag, children}
}

func (p partNode) Render(w io.Writer) error {
	el := p.tag(h.Data("slot", p.slot))
	var user []string
	var rest []any
	for _, c := range p.children {
		if cl, ok := c.(className); ok {
			user = append(user, string(cl))
			continue
		}
		rest = append(rest, c)
	}
	el.Class(lib.Cn(append([]string{p.class(w)}, user...)...))
	el.Children(rest...)
	return el.Render(w)
}

// className is a class added to a part (merged with cn).
type className string

// Class adds classes to Header, Footer, Title or Description.
func Class(c string) any { return className(c) }

func overlayClass(rhea bool) string {
	if rhea {
		return "fixed inset-0 z-50 min-h-dvh bg-black/30 opacity-[max(var(--drawer-overlay-min-opacity,0),calc(1-var(--drawer-swipe-progress)))] transition-opacity duration-450 ease-[cubic-bezier(0.32,0.72,0,1)] select-none data-ending-style:pointer-events-none data-ending-style:opacity-0 data-ending-style:duration-[calc(var(--drawer-swipe-strength)*400ms)] data-snap-points:[--drawer-overlay-min-opacity:0.5] data-starting-style:opacity-0 data-swiping:duration-0 supports-backdrop-filter:backdrop-blur-sm supports-[-webkit-touch-callout:none]:absolute"
	}
	return "fixed inset-0 z-50 min-h-dvh bg-black/10 opacity-[max(var(--drawer-overlay-min-opacity,0),calc(1-var(--drawer-swipe-progress)))] transition-opacity duration-450 ease-[cubic-bezier(0.32,0.72,0,1)] select-none data-ending-style:pointer-events-none data-ending-style:opacity-0 data-ending-style:duration-[calc(var(--drawer-swipe-strength)*400ms)] data-snap-points:[--drawer-overlay-min-opacity:0.5] data-starting-style:opacity-0 data-swiping:duration-0 supports-backdrop-filter:backdrop-blur-xs supports-[-webkit-touch-callout:none]:absolute"
}

func handleClass(rhea bool) string {
	const common = "relative z-10 flex shrink-0 cursor-grab transition-opacity duration-200 group-data-nested-drawer-open/drawer-popup:opacity-0 group-data-nested-drawer-swiping/drawer-popup:opacity-100 group-data-[swipe-axis=x]/drawer-popup:h-full group-data-[swipe-axis=x]/drawer-popup:w-3 group-data-[swipe-axis=x]/drawer-popup:items-center group-data-[swipe-axis=y]/drawer-popup:h-3 group-data-[swipe-axis=y]/drawer-popup:w-full group-data-[swipe-axis=y]/drawer-popup:justify-center group-data-[swipe-direction=down]/drawer-popup:items-end group-data-[swipe-direction=left]/drawer-popup:order-last group-data-[swipe-direction=left]/drawer-popup:justify-start group-data-[swipe-direction=right]/drawer-popup:justify-end group-data-[swipe-direction=up]/drawer-popup:order-last group-data-[swipe-direction=up]/drawer-popup:items-start after:block after:shrink-0 after:rounded-full after:bg-muted "
	if rhea {
		return common + "group-data-[swipe-axis=x]/drawer-popup:after:h-[100px] group-data-[swipe-axis=x]/drawer-popup:after:w-1.5 group-data-[swipe-axis=y]/drawer-popup:after:h-1.5 group-data-[swipe-axis=y]/drawer-popup:after:w-[100px] active:cursor-grabbing"
	}
	return common + "group-data-[swipe-axis=x]/drawer-popup:after:h-24 group-data-[swipe-axis=x]/drawer-popup:after:w-1 group-data-[swipe-axis=y]/drawer-popup:after:h-1 group-data-[swipe-axis=y]/drawer-popup:after:w-24 active:cursor-grabbing"
}

func popupClass(rhea bool) string {
	base := "group/drawer-popup pointer-events-auto fixed z-50 m-(--drawer-inset,0px) flex h-(--drawer-content-height) max-h-(--drawer-content-max-height,none) min-h-0 w-(--drawer-content-width,auto) transform-[translate3d(var(--translate-x,0px),var(--translate-y,0px),0)_scale(var(--stack-scale))] flex-col bg-popover text-sm text-popover-foreground transition-[transform,height,opacity,filter] duration-450 ease-[cubic-bezier(0.22,1,0.36,1)] will-change-transform outline-none select-none [interpolate-size:allow-keywords] data-[swipe-direction=down]:rounded-t-xl data-[swipe-direction=down]:border-t data-[swipe-direction=left]:rounded-r-xl data-[swipe-direction=left]:border-r data-[swipe-direction=right]:rounded-l-xl data-[swipe-direction=right]:border-l data-[swipe-direction=up]:rounded-b-xl data-[swipe-direction=up]:border-b"
	if rhea {
		base = "group/drawer-popup pointer-events-auto fixed z-50 m-(--drawer-inset,0px) flex h-(--drawer-content-height) max-h-(--drawer-content-max-height,none) min-h-0 w-(--drawer-content-width,auto) transform-[translate3d(var(--translate-x,0px),var(--translate-y,0px),0)_scale(var(--stack-scale))] flex-col rounded-[min(var(--radius-4xl),24px)] border border-popover bg-popover text-sm text-popover-foreground shadow-xl transition-[transform,height,opacity,filter] duration-450 ease-[cubic-bezier(0.22,1,0.36,1)] will-change-transform outline-none select-none [--drawer-bleed-background:transparent] [--drawer-inset:--spacing(2)] [--drawer-stacked-shadow:0_-20px_25px_-5px_rgb(0_0_0/0.1),0_-8px_10px_-6px_rgb(0_0_0/0.1)] [interpolate-size:allow-keywords] data-[swipe-direction=down]:data-nested-drawer-open:shadow-(--drawer-stacked-shadow) dark:border-border"
	}
	return lib.Cn(base,
		"data-nested-drawer-open:overflow-hidden data-nested-drawer-open:brightness-95",
		"after:pointer-events-none after:absolute after:bg-(--drawer-bleed-background,var(--color-popover)) data-[swipe-axis=x]:after:inset-y-0 data-[swipe-axis=x]:after:w-(--bleed) data-[swipe-axis=y]:after:inset-x-0 data-[swipe-axis=y]:after:h-(--bleed) data-[swipe-direction=down]:after:top-full data-[swipe-direction=left]:after:right-full data-[swipe-direction=right]:after:left-full data-[swipe-direction=up]:after:bottom-full",
		"[--drawer-content-height:var(--drawer-height,auto)] data-[swipe-axis=x]:[--drawer-content-width:75%] data-[swipe-axis=y]:[--drawer-content-max-height:calc(100dvh-6rem)] data-[swipe-axis=y]:data-snap-points:[--drawer-content-height:100dvh] data-[swipe-axis=x]:sm:[--drawer-content-width:24rem]",
		"[--bleed:3rem] [--peek:1rem] [--stack-height:var(--drawer-frontmost-height,var(--drawer-height,0px))] [--stack-peek-offset:max(0px,calc((var(--nested-drawers)-var(--stack-progress))*var(--peek)))] [--stack-progress:clamp(0,var(--drawer-swipe-progress),1)] [--stack-scale-base:max(0,calc(1-(var(--nested-drawers)*var(--stack-step))))] [--stack-scale:clamp(0,calc(var(--stack-scale-base)+(var(--stack-step)*var(--stack-progress))),1)] [--stack-shrink:calc(1-var(--stack-scale))] [--stack-step:0.05]",
		"data-ending-style:transform-(--closed-transform) data-ending-style:opacity-[0.9999] data-ending-style:duration-[calc(var(--drawer-swipe-strength)*400ms)] data-nested-drawer-swiping:duration-0 data-ending-style:data-nested-drawer-swiping:duration-[calc(var(--drawer-swipe-strength)*400ms)] data-starting-style:transform-(--closed-transform) data-swiping:duration-0 data-ending-style:data-swiping:duration-[calc(var(--drawer-swipe-strength)*400ms)]",
		"data-[swipe-axis=y]:inset-x-0 data-[swipe-axis=y]:data-nested-drawer-open:h-(--stack-height)",
		"data-[swipe-axis=x]:inset-y-0 data-[swipe-axis=x]:flex-row",
		"data-[swipe-direction=down]:bottom-0 data-[swipe-direction=down]:origin-bottom data-[swipe-direction=down]:[--closed-transform:translate3d(0,calc(100%+var(--drawer-inset,0px)+2px),0)] data-[swipe-direction=down]:[--translate-y:calc(var(--drawer-snap-point-offset,0px)+var(--drawer-swipe-movement-y)-var(--stack-peek-offset)-(var(--stack-shrink)*var(--stack-height)))]",
		"data-[swipe-direction=up]:top-0 data-[swipe-direction=up]:origin-top data-[swipe-direction=up]:[--closed-transform:translate3d(0,calc(-100%-var(--drawer-inset,0px)-2px),0)] data-[swipe-direction=up]:[--translate-y:calc(var(--drawer-snap-point-offset,0px)+var(--drawer-swipe-movement-y)+var(--stack-peek-offset)+(var(--stack-shrink)*var(--stack-height)))]",
		"data-[swipe-direction=left]:left-0 data-[swipe-direction=left]:origin-left data-[swipe-direction=left]:[--closed-transform:translate3d(calc(-100%-var(--drawer-inset,0px)-2px),0,0)] data-[swipe-direction=left]:[--translate-x:calc(var(--drawer-swipe-movement-x)+var(--stack-peek-offset)+(var(--stack-shrink)*100%))]",
		"data-[swipe-direction=right]:right-0 data-[swipe-direction=right]:origin-right data-[swipe-direction=right]:[--closed-transform:translate3d(calc(100%+var(--drawer-inset,0px)+2px),0,0)] data-[swipe-direction=right]:[--translate-x:calc(var(--drawer-swipe-movement-x)-var(--stack-peek-offset)-(var(--stack-shrink)*100%))]",
	)
}

const contentClass = "flex min-h-0 flex-1 flex-col overflow-hidden overscroll-contain rounded-[inherit] transition-opacity duration-300 ease-[cubic-bezier(0.45,1.005,0,1.005)] select-text group-data-nested-drawer-open/drawer-popup:opacity-0 group-data-nested-drawer-swiping/drawer-popup:opacity-100 group-data-swiping/drawer-popup:select-none"
