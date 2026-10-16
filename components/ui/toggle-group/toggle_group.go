// Package togglegroup is the shadcn/ui ToggleGroup (Radix ToggleGroup): a
// set of two-state buttons that can be toggled on or off.
//
//	togglegroup.ToggleGroup(togglegroup.TypeMultiple, togglegroup.VariantOutline,
//	    togglegroup.Item(togglegroup.Value("bold"), h.Aria("label", "Toggle bold"), lucide.Bold()),
//	    togglegroup.Item(togglegroup.Value("italic"), h.Aria("label", "Toggle italic"), lucide.Italic()),
//	)
//
// Root children may include a Type (required, like Radix), a Variant, a
// Size, Spacing, Value (the initially pressed items' values, repeatable for
// TypeMultiple), Disabled, Orientation, Loop, RovingFocus and Required. Items take a
// Value, a Variant and a Size (the group's win, like upstream).
//
// Behavior (toggle-group.js): Radix's roving focus (arrow keys, Home/End,
// loop), single or multiple selection; a cancelable "shadcn:value-change"
// event is dispatched on the group before every change, with detail
// { item, pressed, pressedItems } (Radix does not render item values, so the
// detail carries the item elements; cancelling keeps the current state),
// and "shadcn:value-changed", with the same detail, after the items' state
// has changed.
//
// Right to left (Radix's dir prop): pass h.Dir("rtl") to the root; it
// replaces the default dir="ltr" and the arrow keys follow it.
package togglegroup

import (
	"cmp"
	_ "embed"
	"fmt"
	"io"
	"shadcn.go/components/ui/internal/part"
	"strconv"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/toggle"
	h "shadcn.go/libhtml"
)

//go:embed toggle_group.js
var script string

func init() { runtime.Module("toggle-group", script) }

// Variant and Size are the toggle's.
type (
	Variant = toggle.Variant
	Size    = toggle.Size
)

// Variants and sizes, as in package toggle.
const (
	VariantDefault = toggle.VariantDefault
	VariantOutline = toggle.VariantOutline
	SizeDefault    = toggle.SizeDefault
	SizeSm         = toggle.SizeSm
	SizeLg         = toggle.SizeLg
)

// Type selects single or multiple selection.
type Type string

// Selection types.
const (
	TypeSingle   Type = "single"
	TypeMultiple Type = "multiple"
)

// Value is an item's value; on the root, an initially pressed value.
type Value string

// Spacing is the gap between items, in Tailwind spacing units (default 0:
// items are joined).
type Spacing int

// Disabled, passed to the root, disables every item.
type Disabled bool

// Orientation sets the arrow keys that move focus ("horizontal" or
// "vertical"); unset, both work.
type Orientation string

// Orientations.
const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// Loop sets whether keyboard navigation wraps around (default true).
type Loop bool

// RovingFocus sets whether arrow keys move focus among items (default
// true); without it every item is a tab stop.
type RovingFocus bool

// Required keeps one item pressed in a TypeSingle group: pressing the
// pressed item no longer releases it (Radix lets it, leaving no value).
type Required bool

type context struct {
	typ         Type
	variant     Variant
	size        Size
	spacing     int
	spacingSet  bool
	values      []string
	disabled    bool
	orientation Orientation
	rovingFocus bool
	style       style.Style
	rtl         bool
}

// Builder is returned by ToggleGroup.
type Builder struct {
	h.Builder[Builder]
	rtl      bool // inside style.ScopeRTL, set by ElementIn
	ctx      context
	loop     bool
	required bool
}

// ToggleGroup returns the group's root <div>.
func ToggleGroup(children ...any) *Builder {
	b := &Builder{ctx: context{rovingFocus: true}, loop: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the group's typed options.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Type:
		b.ctx.typ = v
	case Variant:
		b.ctx.variant = v
	case Size:
		b.ctx.size = v
	case Spacing:
		b.ctx.spacing = int(v)
		b.ctx.spacingSet = true
	case Value:
		b.ctx.values = append(b.ctx.values, string(v))
	case Disabled:
		b.ctx.disabled = bool(v)
	case Orientation:
		b.ctx.orientation = v
	case Loop:
		b.loop = bool(v)
	case RovingFocus:
		b.ctx.rovingFocus = bool(v)
	case Required:
		b.required = bool(v)
	default:
		return false
	}
	return true
}

// Type sets single or multiple selection.
func (b *Builder) Type(t Type) *Builder { b.ctx.typ = t; return b }

// Variant sets the items' variant.
func (b *Builder) Variant(v Variant) *Builder { b.ctx.variant = v; return b }

// Size sets the items' size.
func (b *Builder) Size(s Size) *Builder { b.ctx.size = s; return b }

// Spacing sets the gap between items.
func (b *Builder) Spacing(n int) *Builder { b.ctx.spacing, b.ctx.spacingSet = n, true; return b }

// Value sets the initially pressed values.
func (b *Builder) Value(values ...string) *Builder { b.ctx.values = values; return b }

// Disabled disables every item.
func (b *Builder) Disabled(on bool) *Builder { b.ctx.disabled = on; return b }

// Orientation sets the navigation orientation.
func (b *Builder) Orientation(o Orientation) *Builder { b.ctx.orientation = o; return b }

// Loop sets whether keyboard navigation wraps around.
func (b *Builder) Loop(on bool) *Builder { b.loop = on; return b }

// Required keeps one item pressed in a TypeSingle group.
func (b *Builder) Required(on bool) *Builder { b.required = on; return b }

// RovingFocus sets whether arrow keys move focus among items.
func (b *Builder) RovingFocus(on bool) *Builder { b.ctx.rovingFocus = on; return b }

const rootClass = "group/toggle-group flex w-fit items-center gap-[--spacing(var(--gap))] rounded-md data-[spacing=default]:data-[variant=outline]:shadow-xs"

// styleFile is registry/bases/radix/ui/toggle-group.tsx (style.Scope).
var styleFile = style.NewFile("ui/toggle-group.tsx")

var (
	styleRoot = styleFile.Cn("cn-toggle-group group/toggle-group flex w-fit flex-row items-center gap-[--spacing(var(--gap))] data-vertical:flex-col data-vertical:items-stretch", style.Arg)
	styleItem = styleFile.Cn("cn-toggle-group-item shrink-0 focus:z-10 focus-visible:z-10 group-data-horizontal/toggle-group:data-[spacing=0]:data-[variant=outline]:border-l-0 group-data-vertical/toggle-group:data-[spacing=0]:data-[variant=outline]:border-t-0 group-data-horizontal/toggle-group:data-[spacing=0]:data-[variant=outline]:first:border-l group-data-vertical/toggle-group:data-[spacing=0]:data-[variant=outline]:first:border-t", style.Arg, style.Arg)
)

// Element resolves the group to its <div>.
func (b *Builder) Element() *h.Element { return b.element("") }

// element renders the group in style st ("" for new-york-v4). In a style,
// as in registry/bases/radix: spacing defaults to 2, data-orientation is
// always rendered (default "horizontal") and the orientation is not given
// to Radix (arrow keys work in both directions).
func (b *Builder) element(st style.Style) *h.Element {
	if b.ctx.typ != TypeSingle && b.ctx.typ != TypeMultiple {
		return h.Errored(fmt.Errorf("togglegroup: missing Type (TypeSingle or TypeMultiple)"))
	}
	ctx := b.ctx
	ctx.style = st
	ctx.rtl = b.rtl
	orientation := ctx.orientation
	if st != "" {
		if !ctx.spacingSet {
			ctx.spacing = 2
		}
		ctx.orientation = ""
	}
	if ctx.typ == TypeSingle && len(ctx.values) > 1 {
		ctx.values = ctx.values[:1]
	}
	focusable := 0
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ItemBuilder:
			c.ctx = &ctx
			if !c.disabled() {
				focusable++
			}
		case *Builder:
			return false
		}
		return true
	})
	base := h.Div(h.Role("group"), h.Dir("ltr"))
	style := "--gap:" + strconv.Itoa(ctx.spacing)
	if ctx.rovingFocus {
		base.TabIndex(h.IfElse(focusable == 0, -1, 0))
		if ctx.orientation != "" {
			base.Data("orientation", string(ctx.orientation))
		}
		style = "outline:none;" + style
		if !b.loop {
			base.Data("loop", "false")
		}
	} else {
		base.Data("roving-focus", "false")
	}
	if b.required {
		base.Data("required", "true")
	}
	base.Data("slot", "toggle-group")
	if ctx.variant != "" {
		base.Data("variant", string(ctx.variant))
	}
	if ctx.size != "" {
		base.Data("size", string(ctx.size))
	}
	base.Data("spacing", strconv.Itoa(ctx.spacing))
	if st != "" {
		if orientation == "" {
			orientation = OrientationHorizontal
		}
		base.Data("orientation", string(orientation))
	}
	props := b.Elem()
	if s, ok := props.GetAttr("style"); ok {
		// Like {...props} after style={{"--gap": spacing}}: the caller's
		// style replaces it; Radix's outline:none stays.
		props = props.Clone()
		props.Attr("style", strings.TrimSuffix(style, "--gap:"+strconv.Itoa(ctx.spacing))+s)
	} else {
		base.Style(style)
	}
	out := h.Spread(base, props, nil)
	if st != "" {
		return out.SetClass(styleRoot.Class(st, props.ClassName()))
	}
	return out.SetClass(lib.Cn(rootClass, props.ClassName()))
}

// Render renders the group, in the style of an enclosing style.Scope.
func (b *Builder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ElementIn is Element in the style of an enclosing style.Scope.
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	b.rtl = part.RTL(w)
	return b.element(st)
}

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	ctx     *context
	value   string
	variant Variant
	size    Size
}

// Item returns a toggle <button> of the group. Children may include a
// Value, a Variant and a Size.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Value, Variant and Size children.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		b.value = string(v)
	case Variant:
		b.variant = v
	case Size:
		b.size = v
	default:
		return false
	}
	return true
}

// Value sets the item's value.
func (b *ItemBuilder) Value(v string) *ItemBuilder { b.value = v; return b }

// Variant sets the item's variant (the group's wins).
func (b *ItemBuilder) Variant(v Variant) *ItemBuilder { b.variant = v; return b }

// Size sets the item's size (the group's wins).
func (b *ItemBuilder) Size(s Size) *ItemBuilder { b.size = s; return b }

// Disabled disables the item.
func (b *ItemBuilder) Disabled(on bool) *ItemBuilder { b.Elem().Disabled(on); return b }

func (b *ItemBuilder) disabled() bool {
	return (b.ctx != nil && b.ctx.disabled) || b.Elem().Has("disabled")
}

const itemClass = "w-auto min-w-0 shrink-0 px-3 focus:z-10 focus-visible:z-10"
const itemSpacingClass = "data-[spacing=0]:rounded-none data-[spacing=0]:shadow-none data-[spacing=0]:first:rounded-l-md data-[spacing=0]:last:rounded-r-md data-[spacing=0]:data-[variant=outline]:border-l-0 data-[spacing=0]:data-[variant=outline]:first:border-l"

// Element resolves the item to its <button>.
func (b *ItemBuilder) Element() *h.Element {
	ctx := b.ctx
	if ctx == nil {
		return h.Errored(fmt.Errorf("togglegroup: Item must be inside a ToggleGroup"))
	}
	variant := ctx.variant
	if variant == "" {
		variant = b.variant
	}
	size := ctx.size
	if size == "" {
		size = b.size
	}
	if ctx.style != "" {
		// variant = "default", size = "default" in the bases item.
		variant = cmp.Or(variant, VariantDefault)
		size = cmp.Or(size, SizeDefault)
	}
	pressed := false
	for _, v := range ctx.values {
		if v == b.value {
			pressed = true
		}
	}
	disabled := b.disabled()
	base := toggle.Root(pressed, disabled)
	if ctx.typ == TypeSingle {
		base.Remove("aria-pressed")
		base.Role("radio").Aria("checked", strconv.FormatBool(pressed))
	}
	if ctx.rovingFocus {
		base.TabIndex(-1)
		if ctx.orientation != "" {
			base.Data("orientation", string(ctx.orientation))
		}
		base.Data("radix-collection-item", "")
	}
	if disabled {
		base.Disabled(true)
	}
	base.Data("slot", "toggle-group-item")
	if variant != "" {
		base.Data("variant", string(variant))
	}
	if size != "" {
		base.Data("size", string(size))
	}
	base.Data("spacing", strconv.Itoa(ctx.spacing))
	props := b.Elem()
	out := h.Spread(base, props, nil)
	if ctx.style != "" {
		return out.SetClass(part.Class(ctx.rtl, styleItem, ctx.style, toggle.StyleVariantClassesDir(ctx.style, ctx.rtl, variant, size), props.ClassName()))
	}
	return out.SetClass(lib.Cn(toggle.VariantClasses(variant, size), itemClass, itemSpacingClass, props.ClassName()))
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
