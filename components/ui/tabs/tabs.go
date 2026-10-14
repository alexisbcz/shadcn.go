// Package tabs is the shadcn/ui Tabs (Radix Tabs): a set of layered
// sections of content, known as tab panels, displayed one at a time.
//
//	tabs.Tabs(tabs.Value("account"),
//	    tabs.List(
//	        tabs.Trigger(tabs.Value("account"), "Account"),
//	        tabs.Trigger(tabs.Value("password"), "Password"),
//	    ),
//	    tabs.Content(tabs.Value("account"), ...),
//	    tabs.Content(tabs.Value("password"), ...),
//	)
//
// Root children may include Value (the initially selected tab, Radix
// defaultValue), Orientation and ActivationMode; List takes a ListVariant
// and Loop; Trigger and Content take the Value they belong to.
//
// Behavior (tabs.js): Radix's roving focus on the list (arrow keys per
// orientation, Home/End, loop), automatic activation on focus (or manual:
// Enter/Space), selection on mousedown, and panels mounted only while
// selected (Radix Presence: exit animations run before unmounting). A
// cancelable "shadcn:value-change" event (detail: { value }) is dispatched
// on the root before every change.
//
// Right to left (Radix's dir prop): pass h.Dir("rtl") to the root; it
// replaces the default dir="ltr" and the arrow keys follow it.
//
// Inactive panels render empty and hidden, like Radix; their content is
// carried in a <template> the behavior mounts from (a fresh copy on every
// activation, as React remounts it).
package tabs

import (
	_ "embed"
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

//go:embed tabs.js
var script string

func init() { runtime.Module("tabs", script) }

// Value is a trigger's or a panel's value; on the root, the initially
// selected value.
type Value string

// Orientation of the tabs (default horizontal).
type Orientation string

// Orientations.
const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// ActivationMode selects whether focusing a trigger selects it
// (automatic, the default) or only Enter/Space/click do (manual).
type ActivationMode string

// Activation modes.
const (
	ActivationAutomatic ActivationMode = "automatic"
	ActivationManual    ActivationMode = "manual"
)

// ListVariant selects the list's look.
type ListVariant string

// List variants.
const (
	ListVariantDefault ListVariant = "default"
	ListVariantLine    ListVariant = "line"
)

// Loop sets whether keyboard navigation in the list wraps (default true).
type Loop bool

type context struct {
	baseID      string
	value       string
	orientation Orientation
}

func (c *context) triggerID(v string) string { return c.baseID + "-trigger-" + v }
func (c *context) contentID(v string) string { return c.baseID + "-content-" + v }

// Builder is returned by Tabs.
type Builder struct {
	h.Builder[Builder]
	value       string
	orientation Orientation
	activation  ActivationMode
}

// Tabs returns the tabs' root <div>.
func Tabs(children ...any) *Builder {
	b := &Builder{orientation: OrientationHorizontal}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Value, Orientation and ActivationMode children.
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		b.value = string(v)
	case Orientation:
		b.orientation = v
	case ActivationMode:
		b.activation = v
	default:
		return false
	}
	return true
}

// Value sets the initially selected tab.
func (b *Builder) Value(v string) *Builder { b.value = v; return b }

// Orientation sets the orientation.
func (b *Builder) Orientation(o Orientation) *Builder { b.orientation = o; return b }

// ActivationMode sets the activation mode.
func (b *Builder) ActivationMode(m ActivationMode) *Builder { b.activation = m; return b }

// Element resolves the tabs to the root <div>.
func (b *Builder) Element() *h.Element {
	orientation := b.orientation
	if orientation == "" {
		orientation = OrientationHorizontal
	}
	ctx := &context{baseID: uid.New("radix"), value: b.value, orientation: orientation}
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ListBuilder:
			c.ctx = ctx
		case *TriggerBuilder:
			c.ctx = ctx
		case *ContentBuilder:
			c.ctx = ctx
			return false // a panel's content may hold other tabs
		case *Builder:
			return false
		}
		return true
	})
	base := h.Div(
		h.Dir("ltr"),
		h.Data("orientation", string(orientation)),
		h.Data("slot", "tabs"),
	)
	if b.activation == ActivationManual {
		base.Data("activation-mode", "manual")
	}
	props := b.Elem()
	return h.Spread(base, props, nil).SetClass(lib.Cn("group/tabs flex gap-2 data-[orientation=horizontal]:flex-col", props.ClassName()))
}

// Render renders the tabs.
func (b *Builder) Render(w io.Writer) error { return b.Element().Render(w) }

var listVariants = cva.New(
	"group/tabs-list inline-flex w-fit items-center justify-center rounded-lg p-[3px] text-muted-foreground group-data-[orientation=horizontal]/tabs:h-9 group-data-[orientation=vertical]/tabs:h-fit group-data-[orientation=vertical]/tabs:flex-col data-[variant=line]:rounded-none",
	cva.Config{
		Variants: cva.Variants{
			"variant": {
				"default": "bg-muted",
				"line":    "gap-1 bg-transparent",
			},
		},
		Defaults: cva.Selection{"variant": "default"},
	},
)

// ListVariants returns the list's classes for a variant, plus className,
// conflicts resolved (tabsListVariants upstream).
func ListVariants(v ListVariant, className ...string) string {
	return lib.Cn(listVariants.Classes(cva.Selection{"variant": string(v)}, className...))
}

// ListBuilder is returned by List.
type ListBuilder struct {
	h.Builder[ListBuilder]
	ctx     *context
	variant ListVariant
	loop    bool
}

// List returns the tab list <div role="tablist">.
func List(children ...any) *ListBuilder {
	b := &ListBuilder{variant: ListVariantDefault, loop: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes ListVariant and Loop children.
func (b *ListBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case ListVariant:
		b.variant = v
	case Loop:
		b.loop = bool(v)
	default:
		return false
	}
	return true
}

// Variant sets the list's variant.
func (b *ListBuilder) Variant(v ListVariant) *ListBuilder { b.variant = v; return b }

// Loop sets whether keyboard navigation wraps.
func (b *ListBuilder) Loop(on bool) *ListBuilder { b.loop = on; return b }

// Element resolves the list to its <div>.
func (b *ListBuilder) Element() *h.Element {
	if b.ctx == nil {
		return h.Errored(fmt.Errorf("tabs: List must be inside Tabs"))
	}
	variant := b.variant
	if variant == "" {
		variant = ListVariantDefault
	}
	focusable := 0
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		if t, ok := n.(*TriggerBuilder); ok && !t.Elem().Has("disabled") {
			focusable++
		}
		_, nested := n.(*Builder)
		return !nested
	})
	o := string(b.ctx.orientation)
	base := h.Div(
		h.TabIndex(h.IfElse(focusable == 0, -1, 0)),
		h.Data("orientation", o),
		h.Style("outline:none"),
		h.Role("tablist"),
		h.Aria("orientation", o),
		h.Data("slot", "tabs-list"),
		h.Data("variant", string(variant)),
	)
	if !b.loop {
		base.Data("loop", "false")
	}
	props := b.Elem()
	return h.Spread(base, props, nil).SetClass(lib.Cn(listVariants.Classes(cva.Selection{"variant": string(variant)}), props.ClassName()))
}

// Render renders the list.
func (b *ListBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// TriggerBuilder is returned by Trigger.
type TriggerBuilder struct {
	h.Builder[TriggerBuilder]
	ctx   *context
	value string
}

// Trigger returns a tab <button role="tab">. Children may include its
// Value.
func Trigger(children ...any) *TriggerBuilder {
	b := &TriggerBuilder{}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// AcceptChild consumes Value children.
func (b *TriggerBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Value); ok {
		b.value = string(v)
		return true
	}
	return false
}

// Value sets the trigger's value.
func (b *TriggerBuilder) Value(v string) *TriggerBuilder { b.value = v; return b }

// Disabled disables the trigger.
func (b *TriggerBuilder) Disabled(on bool) *TriggerBuilder { b.Elem().Disabled(on); return b }

const triggerClass = "relative inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap text-foreground/60 transition-all group-data-[orientation=vertical]/tabs:w-full group-data-[orientation=vertical]/tabs:justify-start hover:text-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50 group-data-[variant=default]/tabs-list:data-[state=active]:shadow-sm group-data-[variant=line]/tabs-list:data-[state=active]:shadow-none dark:text-muted-foreground dark:hover:text-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
const triggerLineClass = "group-data-[variant=line]/tabs-list:bg-transparent group-data-[variant=line]/tabs-list:data-[state=active]:bg-transparent dark:group-data-[variant=line]/tabs-list:data-[state=active]:border-transparent dark:group-data-[variant=line]/tabs-list:data-[state=active]:bg-transparent"
const triggerActiveClass = "data-[state=active]:bg-background data-[state=active]:text-foreground dark:data-[state=active]:border-input dark:data-[state=active]:bg-input/30 dark:data-[state=active]:text-foreground"
const triggerAfterClass = "after:absolute after:bg-foreground after:opacity-0 after:transition-opacity group-data-[orientation=horizontal]/tabs:after:inset-x-0 group-data-[orientation=horizontal]/tabs:after:bottom-[-5px] group-data-[orientation=horizontal]/tabs:after:h-0.5 group-data-[orientation=vertical]/tabs:after:inset-y-0 group-data-[orientation=vertical]/tabs:after:-right-1 group-data-[orientation=vertical]/tabs:after:w-0.5 group-data-[variant=line]/tabs-list:data-[state=active]:after:opacity-100"

// Element resolves the trigger to its <button>.
func (b *TriggerBuilder) Element() *h.Element {
	if b.ctx == nil {
		return h.Errored(fmt.Errorf("tabs: Trigger must be inside Tabs"))
	}
	props := b.Elem()
	selected := b.value == b.ctx.value
	disabled := props.Has("disabled")
	base := h.Button(
		h.TabIndex(-1),
		h.Data("orientation", string(b.ctx.orientation)),
		h.Data("radix-collection-item", ""),
		h.Type("button"),
		h.Role("tab"),
		h.Aria("selected", h.IfElse(selected, "true", "false")),
		h.Aria("controls", b.ctx.contentID(b.value)),
		h.Data("state", h.IfElse(selected, "active", "inactive")),
		h.ID(b.ctx.triggerID(b.value)),
		h.Data("slot", "tabs-trigger"),
	)
	if disabled {
		base.Data("disabled", "")
	}
	return h.Spread(base, props, nil).SetClass(lib.Cn(triggerClass, triggerLineClass, triggerActiveClass, triggerAfterClass, props.ClassName()))
}

// Render renders the trigger.
func (b *TriggerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// ContentBuilder is returned by Content.
type ContentBuilder struct {
	h.Builder[ContentBuilder]
	ctx   *context
	value string
}

// Content returns a tab panel <div role="tabpanel">. Children may include
// its Value.
func Content(children ...any) *ContentBuilder {
	b := &ContentBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Value children.
func (b *ContentBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Value); ok {
		b.value = string(v)
		return true
	}
	return false
}

// Value sets the panel's value.
func (b *ContentBuilder) Value(v string) *ContentBuilder { b.value = v; return b }

// Element resolves the panel to its <div>.
func (b *ContentBuilder) Element() *h.Element {
	if b.ctx == nil {
		return h.Errored(fmt.Errorf("tabs: Content must be inside Tabs"))
	}
	props := b.Elem()
	selected := b.value == b.ctx.value
	base := h.Div(
		h.Data("state", h.IfElse(selected, "active", "inactive")),
		h.Data("orientation", string(b.ctx.orientation)),
		h.Role("tabpanel"),
		h.Aria("labelledby", b.ctx.triggerID(b.value)),
		h.ID(b.ctx.contentID(b.value)),
		h.TabIndex(0),
		h.Data("slot", "tabs-content"),
	)
	// Children are rendered only while present; the template carries them
	// for the behavior to mount on activation.
	userProps := props.Clone().SetChildren()
	style, hasStyle := props.GetAttr("style")
	if selected {
		// Radix: animationDuration "0s" on the panel selected at mount, so
		// it does not animate in.
		if hasStyle && style != "" {
			style += ";"
		}
		userProps.Attr("style", style+"animation-duration:0s")
	} else {
		base.Hidden(true)
	}
	out := h.Spread(base, userProps, nil).SetClass(lib.Cn("flex-1 outline-none", props.ClassName()))
	kids := props.ChildNodes()
	if selected {
		out.Children(kids...)
	}
	out.Children(h.Template(append([]any{h.Data("slot", "tabs-content-template")}, kids...)...))
	return out
}

// Render renders the panel.
func (b *ContentBuilder) Render(w io.Writer) error { return b.Element().Render(w) }
