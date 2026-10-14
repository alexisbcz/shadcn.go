// Package resizable is the shadcn/ui Resizable: panel groups resizable with
// the mouse, touch or keyboard, a port of react-resizable-panels (v4.5.8).
//
//	resizable.PanelGroup(resizable.Horizontal, h.Class("max-w-md rounded-lg border"),
//		resizable.Panel(resizable.DefaultSize("25%"), sidebar),
//		resizable.Handle(resizable.WithHandle),
//		resizable.Panel(resizable.DefaultSize("75%"), content),
//	)
//
// Sizes (DefaultSize, MinSize, MaxSize, CollapsedSize) accept the units
// react-resizable-panels accepts: "%" (also a bare number), "px", "rem",
// "em", "vh" and "vw".
//
// Like react-resizable-panels, a Panel's h.Class and h.Style apply to its
// inner <div>; its other attributes to the outer one.
//
// Ids: like the React components, groups, panels and handles get ids
// (rendered as id and data-testid) unless given one with h.ID. They are
// generated like React's useId ("_r_0_", "_r_1_", ...) in render order
// within each outermost group; the client renumbers duplicates when a page
// renders several groups.
//
// Events, dispatched on the group element: "shadcn:layout-change" with
// detail.layout ({panelId: percentage}) whenever the layout changes, and
// "shadcn:layout-changed" when a change is complete (not during a drag).
// The group element exposes getLayout() and setLayout(layout) as
// el.resizableGroup; a panel exposes collapse(), expand(), getSize(),
// isCollapsed() and resize(size) as el.resizablePanel.
package resizable

import (
	_ "embed"
	"encoding/json"
	"io"
	"math"
	"strconv"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed resizable.js
var script string

func init() { runtime.Module("resizable", script) }

// Orientation is the direction panels are laid out in.
type Orientation string

// Orientations.
const (
	Horizontal Orientation = "horizontal"
	Vertical   Orientation = "vertical"
)

// Disabled, given to PanelGroup, disables resizing.
type Disabled bool

// DisableCursor, given to PanelGroup, keeps the cursor unchanged while
// hovering or dragging handles.
type DisableCursor bool

// DefaultLayout, given to PanelGroup, sets the initial layout: panel id ->
// percentage. It is used only if it names every panel.
type DefaultLayout map[string]float64

// ResizeTargetMinimumSize, given to PanelGroup, sets the minimum size in
// pixels of the area where a handle can be grabbed, for fine (mouse) and
// coarse (touch) pointers. Default: 10 and 20.
type ResizeTargetMinimumSize struct{ Fine, Coarse float64 }

// DefaultSize, given to Panel, sets its initial size.
type DefaultSize string

// MinSize, given to Panel, sets its minimum size (default "0%").
type MinSize string

// MaxSize, given to Panel, sets its maximum size (default "100%").
type MaxSize string

// CollapsedSize, given to Panel, sets its size when collapsed (default
// "0%").
type CollapsedSize string

// Collapsible, given to Panel, lets it collapse below its minimum size.
type Collapsible bool

// WithHandleFlag is the type of WithHandle.
type WithHandleFlag struct{}

// WithHandle, given to Handle, shows a grip.
var WithHandle = WithHandleFlag{}

// --- group ---------------------------------------------------------------------

// GroupBuilder is returned by PanelGroup.
type GroupBuilder struct {
	h.Builder[GroupBuilder]
	orientation   Orientation
	disabled      bool
	disableCursor bool
	defaultLayout DefaultLayout
	minimumSize   *ResizeTargetMinimumSize
}

// PanelGroup returns a group of panels separated by handles. Children may
// include an Orientation (default Horizontal), Disabled, DisableCursor,
// DefaultLayout and ResizeTargetMinimumSize.
func PanelGroup(children ...any) *GroupBuilder {
	b := &GroupBuilder{orientation: Horizontal}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the group options.
func (b *GroupBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Orientation:
		b.Orientation(v)
	case Disabled:
		b.disabled = bool(v)
	case DisableCursor:
		b.disableCursor = bool(v)
	case DefaultLayout:
		b.defaultLayout = v
	case ResizeTargetMinimumSize:
		b.minimumSize = &v
	default:
		return false
	}
	return true
}

// Orientation sets the orientation. An empty value selects Horizontal.
func (b *GroupBuilder) Orientation(o Orientation) *GroupBuilder {
	if o == "" {
		o = Horizontal
	}
	b.orientation = o
	return b
}

// Element resolves the group.
func (b *GroupBuilder) Element() *h.Element {
	return b.element(&idScope{})
}

// Render renders the group.
func (b *GroupBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

type idScope struct{ n int }

// next returns the next id in React's useId format.
func (s *idScope) next() string {
	id := "_r_" + strconv.FormatInt(int64(s.n), 32) + "_"
	s.n++
	return id
}

func (b *GroupBuilder) element(ids *idScope) *h.Element {
	props := b.Elem()
	id := ids.next()
	if v, ok := props.GetAttr("id"); ok {
		id = v
	}
	dir := "row"
	if b.orientation == Vertical {
		dir = "column"
	}
	userStyle, _ := props.GetAttr("style")
	// Resolve the children in render order: the group's id comes first,
	// then its descendants', like React's useId.
	kids := resolve(props.ChildNodes(), ids)

	// Initial layout, computed like the client does when it mounts.
	var panels []*PanelBuilder
	for _, k := range kids {
		if v, ok := k.(*resolvedPanel); ok {
			panels = append(panels, v.b)
		}
	}
	layout, constraints, ok := initialLayout(panels, b.defaultLayout)

	el := h.Spread(h.Div(), props.Clone().SetChildren(), nil)
	el.SetClass(lib.Cn("flex h-full w-full aria-[orientation=vertical]:flex-col", props.ClassName()))
	el.Data("slot", "resizable-panel-group")
	el.Data("group", "true")
	el.Data("testid", id)
	el.ID(id)
	el.Style(joinStyle("height: 100%; width: 100%; overflow: hidden;", userStyle, "display: flex; flex-direction: "+dir+"; flex-wrap: nowrap;"))
	if cfg := b.config(); cfg != "" {
		el.Data("resizable", cfg)
	}
	pi := 0
	for _, k := range kids {
		switch v := k.(type) {
		case *resolvedPanel:
			grow := "1"
			if ok {
				grow = formatNumber(layout[pi])
			}
			pi++
			el.Children(v.b.render(v.id, grow, v.kids))
		case *resolvedHandle:
			var aria *handleAria
			if ok {
				aria = separatorAria(kids, v, layout, constraints)
			}
			el.Children(v.b.render(v.id, b.orientation, aria))
		default:
			el.Children(k)
		}
	}
	return el
}

func (b *GroupBuilder) config() string {
	cfg := map[string]any{}
	if b.disabled {
		cfg["disabled"] = true
	}
	if b.disableCursor {
		cfg["disableCursor"] = true
	}
	if b.defaultLayout != nil {
		cfg["defaultLayout"] = b.defaultLayout
	}
	if b.minimumSize != nil {
		cfg["resizeTargetMinimumSize"] = map[string]float64{"fine": b.minimumSize.Fine, "coarse": b.minimumSize.Coarse}
	}
	if len(cfg) == 0 {
		return ""
	}
	out, _ := json.Marshal(cfg)
	return string(out)
}

type resolvedPanel struct {
	b    *PanelBuilder
	id   string
	kids []any
}

type resolvedHandle struct {
	b  *HandleBuilder
	id string
}

func (*resolvedPanel) Render(io.Writer) error  { return nil }
func (*resolvedHandle) Render(io.Writer) error { return nil }

// resolve assigns ids in render order to the group's direct panels and
// handles and, depth first, to nested groups.
func resolve(nodes []any, ids *idScope) []any {
	var out []any
	for _, n := range nodes {
		switch v := n.(type) {
		case *PanelBuilder:
			id := ids.next()
			if u, ok := v.Elem().GetAttr("id"); ok {
				id = u
			}
			out = append(out, &resolvedPanel{b: v, id: id, kids: resolveNested(v.Elem().ChildNodes(), ids)})
		case *HandleBuilder:
			id := ids.next()
			if u, ok := v.Elem().GetAttr("id"); ok {
				id = u
			}
			out = append(out, &resolvedHandle{b: v, id: id})
		default:
			out = append(out, resolveNested([]any{n}, ids)...)
		}
	}
	return out
}

// resolveNested renders nested groups within the given nodes with the
// shared id scope.
func resolveNested(nodes []any, ids *idScope) []any {
	out := make([]any, 0, len(nodes))
	for _, n := range nodes {
		switch v := n.(type) {
		case *GroupBuilder:
			out = append(out, v.element(ids))
		case *h.Element:
			c := v.Clone().SetChildren(resolveNested(v.ChildNodes(), ids)...)
			out = append(out, c)
		default:
			out = append(out, n)
		}
	}
	return out
}

// --- panel -----------------------------------------------------------------------

// PanelBuilder is returned by Panel.
type PanelBuilder struct {
	h.Builder[PanelBuilder]
	defaultSize   string
	minSize       string
	maxSize       string
	collapsedSize string
	collapsible   bool
}

// Panel returns a panel. Children may include DefaultSize, MinSize,
// MaxSize, CollapsedSize and Collapsible.
func Panel(children ...any) *PanelBuilder {
	b := &PanelBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the panel options.
func (b *PanelBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case DefaultSize:
		b.defaultSize = string(v)
	case MinSize:
		b.minSize = string(v)
	case MaxSize:
		b.maxSize = string(v)
	case CollapsedSize:
		b.collapsedSize = string(v)
	case Collapsible:
		b.collapsible = bool(v)
	default:
		return false
	}
	return true
}

// Render renders the panel on its own (outside a PanelGroup it cannot
// resize).
func (b *PanelBuilder) Render(w io.Writer) error {
	ids := &idScope{}
	id := ids.next()
	return b.render(id, "1", resolveNested(b.Elem().ChildNodes(), ids)).Render(w)
}

const panelStyle = "min-height: 0px; max-height: 100%; height: auto; min-width: 0px; max-width: 100%; width: auto; border: none; border-width: 0px; padding: 0px; margin: 0px; display: flex; flex-basis: 0px; flex-shrink: 1; overflow: hidden;"

func (b *PanelBuilder) render(id, grow string, kids []any) *h.Element {
	props := b.Elem()
	userStyle, _ := props.GetAttr("style")
	outerProps := props.Clone().SetChildren().SetClass("").Remove("style")
	el := h.Spread(h.Div(h.Data("slot", "resizable-panel")), outerProps, nil)
	el.Data("panel", "true")
	el.Data("testid", id)
	el.ID(id)
	el.Style(panelStyle + " flex-grow: " + grow + ";")
	if cfg := b.config(); cfg != "" {
		el.Data("resizable", cfg)
	}
	inner := h.Div()
	if c := props.ClassName(); c != "" {
		inner.Class(c)
	}
	inner.Style(joinStyle("max-height: 100%; max-width: 100%; flex-grow: 1;", userStyle))
	inner.Children(kids...)
	return el.Children(inner)
}

func (b *PanelBuilder) config() string {
	cfg := map[string]any{}
	if b.defaultSize != "" {
		cfg["defaultSize"] = b.defaultSize
	}
	if b.minSize != "" {
		cfg["minSize"] = b.minSize
	}
	if b.maxSize != "" {
		cfg["maxSize"] = b.maxSize
	}
	if b.collapsedSize != "" {
		cfg["collapsedSize"] = b.collapsedSize
	}
	if b.collapsible {
		cfg["collapsible"] = true
	}
	if len(cfg) == 0 {
		return ""
	}
	out, _ := json.Marshal(cfg)
	return string(out)
}

// --- handle ----------------------------------------------------------------------

// HandleBuilder is returned by Handle.
type HandleBuilder struct {
	h.Builder[HandleBuilder]
	withHandle bool
}

// Handle returns the separator between two panels. Pass WithHandle to show
// a grip.
func Handle(children ...any) *HandleBuilder {
	b := &HandleBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes WithHandle.
func (b *HandleBuilder) AcceptChild(c any) bool {
	if _, ok := c.(WithHandleFlag); ok {
		b.withHandle = true
		return true
	}
	return false
}

// Render renders the handle on its own.
func (b *HandleBuilder) Render(w io.Writer) error {
	return b.render((&idScope{}).next(), Horizontal, nil).Render(w)
}

const handleClass = "relative flex w-px items-center justify-center bg-border after:absolute after:inset-y-0 after:left-1/2 after:w-1 after:-translate-x-1/2 focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:outline-hidden aria-[orientation=horizontal]:h-px aria-[orientation=horizontal]:w-full aria-[orientation=horizontal]:after:left-0 aria-[orientation=horizontal]:after:h-1 aria-[orientation=horizontal]:after:w-full aria-[orientation=horizontal]:after:translate-x-0 aria-[orientation=horizontal]:after:-translate-y-1/2 [&[aria-orientation=horizontal]>div]:rotate-90"

type handleAria struct {
	controls       string
	min, max, now  float64
	hasMin, hasMax bool
}

func (b *HandleBuilder) render(id string, groupOrientation Orientation, aria *handleAria) *h.Element {
	props := b.Elem()
	userStyle, _ := props.GetAttr("style")
	el := h.Spread(h.Div(h.Data("slot", "resizable-handle")), props.Clone().SetChildren().SetClass(""), nil)
	if aria != nil {
		el.Aria("controls", aria.controls)
		if aria.hasMax {
			el.Aria("valuemax", formatNumber(aria.max))
		}
		if aria.hasMin {
			el.Aria("valuemin", formatNumber(aria.min))
		}
		el.Aria("valuenow", formatNumber(aria.now))
	}
	orientation := "vertical"
	if groupOrientation == Vertical {
		orientation = "horizontal"
	}
	el.Aria("orientation", orientation)
	el.SetClass(lib.Cn(handleClass, props.ClassName()))
	el.Data("separator", "inactive")
	el.Data("testid", id)
	el.ID(id)
	el.Role("separator")
	el.Style(joinStyle("flex-basis: auto;", userStyle, "flex-grow: 0; flex-shrink: 0;"))
	el.TabIndex(0)
	if b.withHandle {
		el.Children(h.Div(h.Class("z-10 flex h-4 w-3 items-center justify-center rounded-xs border bg-border"),
			lucide.GripVertical(h.Class("size-2.5")),
		))
	}
	el.Children(props.ChildNodes()...)
	return el
}

func joinStyle(parts ...string) string {
	var b strings.Builder
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strings.TrimSuffix(p, ";"))
		b.WriteByte(';')
	}
	return b.String()
}

// formatNumber formats like JavaScript's Number#toString for the values
// layouts hold (at most 3 decimals).
func formatNumber(f float64) string {
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
