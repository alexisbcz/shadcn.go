// Package popper is the Go side of Radix's Popper primitive
// (@radix-ui/react-popper 1.2.8): the markup of PopperContent's wrapper and
// PopperArrow, and the positioning options every popper-based component
// (Popover, Tooltip, HoverCard, Select, menus, ...) accepts on its content.
// The behavior lives in components/ui/runtime/js/popper.js.
//
// A component renders its content element as Radix does inside
// PopperContent, then wraps it:
//
//	var o popper.Options                  // defaults: bottom, center, 0, ...
//	o.Accept(popper.SideTop)              // from the content's children
//	wrapper := popper.Wrap(o, contentEl)  // <div data-radix-popper-content-wrapper>
//
// and creates the popper in its JS module once the wrapper is in the DOM:
//
//	import { createPopper } from "./popper.js"
//	const popper = createPopper(wrapper, { anchor: trigger })
//
// The options travel in a marker attribute that popper.js removes.
package popper

import (
	"encoding/json"
	"strconv"
	"strings"

	h "shadcn.go/libhtml"
)

// Side is the preferred side of the anchor to render against.
type Side string

// Sides.
const (
	SideTop    Side = "top"
	SideRight  Side = "right"
	SideBottom Side = "bottom"
	SideLeft   Side = "left"
)

// Align is the preferred alignment against the anchor.
type Align string

// Alignments.
const (
	AlignStart  Align = "start"
	AlignCenter Align = "center"
	AlignEnd    Align = "end"
)

// SideOffset is the distance in pixels from the anchor.
type SideOffset float64

// AlignOffset is an offset in pixels from the "start" or "end" alignment.
type AlignOffset float64

// ArrowPadding is the padding between the arrow and the content edges.
type ArrowPadding float64

// AvoidCollisions sets whether the content flips and shifts to stay in
// the boundary (default true).
type AvoidCollisions bool

// CollisionPadding is the distance in pixels from the boundary edges
// where collision detection should occur, on every side.
type CollisionPadding float64

// CollisionPaddingSides is CollisionPadding per side (unset sides are 0).
type CollisionPaddingSides struct{ Top, Right, Bottom, Left float64 }

// CollisionBoundary lists CSS selectors of the elements used as the
// collision boundary (default: the viewport).
type CollisionBoundary []string

// Sticky is the sticky behavior on the align axis: "partial" keeps the
// content in the boundary as long as the trigger is at least partially in
// it, "always" keeps it in the boundary regardless.
type Sticky string

// Sticky behaviors.
const (
	StickyPartial Sticky = "partial"
	StickyAlways  Sticky = "always"
)

// HideWhenDetached hides the content when the anchor is fully occluded.
type HideWhenDetached bool

// UpdatePositionStrategy selects how often the position is updated:
// "optimized" (on scroll, resize and layout changes) or "always" (every
// animation frame).
type UpdatePositionStrategy string

// Update position strategies.
const (
	UpdateOptimized UpdatePositionStrategy = "optimized"
	UpdateAlways    UpdatePositionStrategy = "always"
)

// Options are PopperContent's props.
type Options struct {
	Side                   Side
	SideOffset             float64
	Align                  Align
	AlignOffset            float64
	ArrowPadding           float64
	AvoidCollisions        bool
	CollisionBoundary      []string
	CollisionPadding       any // float64 or CollisionPaddingSides
	Sticky                 Sticky
	HideWhenDetached       bool
	UpdatePositionStrategy UpdatePositionStrategy
	Dir                    string
}

// Defaults returns Radix's defaults.
func Defaults() Options {
	return Options{Side: SideBottom, Align: AlignCenter, AvoidCollisions: true, Sticky: StickyPartial, UpdatePositionStrategy: UpdateOptimized}
}

// Accept consumes a positioning option passed as a child. It reports
// whether c was one.
func (o *Options) Accept(c any) bool {
	switch v := c.(type) {
	case Side:
		o.Side = v
	case Align:
		o.Align = v
	case SideOffset:
		o.SideOffset = float64(v)
	case AlignOffset:
		o.AlignOffset = float64(v)
	case ArrowPadding:
		o.ArrowPadding = float64(v)
	case AvoidCollisions:
		o.AvoidCollisions = bool(v)
	case CollisionPadding:
		o.CollisionPadding = float64(v)
	case CollisionPaddingSides:
		o.CollisionPadding = v
	case CollisionBoundary:
		o.CollisionBoundary = []string(v)
	case Sticky:
		o.Sticky = v
	case HideWhenDetached:
		o.HideWhenDetached = bool(v)
	case UpdatePositionStrategy:
		o.UpdatePositionStrategy = v
	default:
		return false
	}
	return true
}

func (o Options) normalized() Options {
	if o.Side == "" {
		o.Side = SideBottom
	}
	if o.Align == "" {
		o.Align = AlignCenter
	}
	if o.Sticky == "" {
		o.Sticky = StickyPartial
	}
	if o.UpdatePositionStrategy == "" {
		o.UpdatePositionStrategy = UpdateOptimized
	}
	return o
}

// JSON returns the options as popper.js reads them.
func (o Options) JSON() string {
	o = o.normalized()
	m := map[string]any{
		"side":                   string(o.Side),
		"sideOffset":             o.SideOffset,
		"align":                  string(o.Align),
		"alignOffset":            o.AlignOffset,
		"arrowPadding":           o.ArrowPadding,
		"avoidCollisions":        o.AvoidCollisions,
		"sticky":                 string(o.Sticky),
		"hideWhenDetached":       o.HideWhenDetached,
		"updatePositionStrategy": string(o.UpdatePositionStrategy),
	}
	switch p := o.CollisionPadding.(type) {
	case float64:
		m["collisionPadding"] = p
	case CollisionPaddingSides:
		m["collisionPadding"] = map[string]float64{"top": p.Top, "right": p.Right, "bottom": p.Bottom, "left": p.Left}
	}
	if len(o.CollisionBoundary) > 0 {
		m["collisionBoundary"] = o.CollisionBoundary
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// Wrap returns PopperContent's markup before the first position is known:
// the wrapper div holding content, which gets data-side, data-align and
// `animation: none` (appended to its style). content is modified.
func Wrap(o Options, content *h.Element) *h.Element {
	o = o.normalized()
	content.Data("side", string(o.Side)).Data("align", string(o.Align))
	AppendStyle(content, "animation: none")
	return h.Div(
		h.Data("radix-popper-content-wrapper", ""),
		h.Style("position: fixed; left: 0px; top: 0px; transform: translate(0, -200%); min-width: max-content;"),
		h.If(o.Dir != "", h.Dir(o.Dir)),
		h.Attr("data-radix-popper", o.JSON()),
		content,
	)
}

// AppendStyle appends declarations to el's style attribute.
func AppendStyle(el *h.Element, decl string) {
	st, _ := el.GetAttr("style")
	st = strings.TrimSpace(st)
	if st != "" && !strings.HasSuffix(st, ";") {
		st += ";"
	}
	if st != "" {
		st += " "
	}
	el.Attr("style", st+decl+";")
}

// ContentVars returns the declarations a component's content sets to
// re-namespace the popper's CSS variables, e.g. for prefix "popover":
// --radix-popover-content-transform-origin: var(--radix-popper-transform-origin); ...
func ContentVars(prefix string) string {
	p := "--radix-" + prefix
	return p + "-content-transform-origin: var(--radix-popper-transform-origin); " +
		p + "-content-available-width: var(--radix-popper-available-width); " +
		p + "-content-available-height: var(--radix-popper-available-height); " +
		p + "-trigger-width: var(--radix-popper-anchor-width); " +
		p + "-trigger-height: var(--radix-popper-anchor-height)"
}

// arrowWrapperStyle is PopperArrow's span style before positioning, for
// the given side.
func arrowWrapperStyle(side Side) string {
	base := map[Side]string{SideTop: "bottom", SideRight: "left", SideBottom: "top", SideLeft: "right"}[side]
	origin := map[Side]string{SideTop: "", SideRight: "0 0", SideBottom: "center 0", SideLeft: "100% 0"}[side]
	transform := map[Side]string{
		SideTop:    "translateY(100%)",
		SideRight:  "translateY(50%) rotate(90deg) translateX(-50%)",
		SideBottom: "rotate(180deg)",
		SideLeft:   "translateY(50%) rotate(-90deg) translateX(50%)",
	}[side]
	s := "position: absolute; " + base + ": 0px; "
	if origin != "" {
		s += "transform-origin: " + origin + "; "
	}
	return s + "transform: " + transform + "; visibility: hidden;"
}

// Arrow returns PopperArrow's markup (the measuring span and the Arrow svg)
// for content placed on side. props are the user's props for the svg
// (attributes, class, children replacing the polygon when asChild-like
// children are given); width and height default to 10 and 5.
func Arrow(side Side, props *h.Element) *h.Element {
	if side == "" {
		side = SideBottom
	}
	w, hgt := "10", "5"
	if props != nil {
		if v, ok := props.GetAttr("width"); ok {
			w = v
		}
		if v, ok := props.GetAttr("height"); ok {
			hgt = v
		}
	}
	svg := h.El("svg")
	if props != nil {
		svg = h.Spread(svg, props, nil)
	}
	svg.Attr("width", w).Attr("height", hgt).Attr("viewBox", "0 0 30 10").Attr("preserveAspectRatio", "none")
	AppendStyle(svg, "display: block")
	if len(svg.ChildNodes()) == 0 {
		svg.Children(h.El("polygon", h.Attr("points", "0,0 30,0 15,10")))
	}
	return h.Span(h.Attr("data-radix-popper-arrow", ""), h.Style(arrowWrapperStyle(side)), svg)
}

// FormatNumber formats a float like JavaScript's Number#toString for the
// values options use.
func FormatNumber(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
