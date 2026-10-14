// Package command is the shadcn/ui Command: a fast, composable command
// menu (cmdk 1.1.1) with fuzzy filtering, keyboard navigation and groups.
//
//	command.Command(h.Class("rounded-lg border shadow-md"),
//		command.Input(h.Placeholder("Type a command or search...")),
//		command.List(
//			command.Empty("No results found."),
//			command.Group(command.Heading("Suggestions"),
//				command.Item(lucide.Calendar(), h.Span("Calendar")),
//				command.Item(command.Disabled(true), lucide.Calculator(), h.Span("Calculator")),
//			),
//			command.Separator(),
//			command.Group(command.Heading("Settings"),
//				command.Item(lucide.User(), h.Span("Profile"), command.Shortcut("⌘P")),
//			),
//		),
//	)
//
// Dialog wraps a Command in a shadcn Dialog (CommandDialog upstream).
//
// Markup is what cmdk renders, including the first enabled item selected
// on load. Behavior (command.js) ports cmdk: typing filters and sorts items
// with cmdk's command-score, ArrowUp/ArrowDown (Ctrl+N/J/P/K), Home/End,
// Meta+Arrow (first/last), Alt+Arrow (previous/next group) move the
// selection, pointer movement selects, Enter or a click selects an item:
// the item dispatches "shadcn:select" (detail.value). Selection changes
// dispatch "shadcn:value-change" on the root (detail.value) and search
// changes "shadcn:search-change" (detail.search).
//
// Options that cmdk takes as props and does not render (Loop, ShouldFilter,
// VimBindings, DisablePointerSelection, item Keywords, ForceMount,
// AlwaysRender) are carried by data-cmdk-* attributes, rendered only when
// they differ from cmdk's defaults.
package command

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html"
	"io"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/dialog/primitive"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed command.js
var script string

func init() { runtime.Module("command", script) }

var (
	rootClass      = lib.Cn("flex h-full w-full flex-col overflow-hidden rounded-md bg-popover text-popover-foreground")
	dialogCmdClass = "**:data-[slot=command-input-wrapper]:h-12 [&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-muted-foreground [&_[cmdk-group]]:px-2 [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-input]]:h-12 [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-item]_svg]:h-5 [&_[cmdk-item]_svg]:w-5"
	inputWrapClass = "flex h-9 items-center gap-2 border-b px-3"
	inputClass     = lib.Cn("flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-hidden placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50")
	listClass      = lib.Cn("max-h-[300px] scroll-py-1 overflow-x-hidden overflow-y-auto")
	emptyClass     = "py-6 text-center text-sm"
	groupClass     = lib.Cn("overflow-hidden p-1 text-foreground [&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1.5 [&_[cmdk-group-heading]]:text-xs [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-muted-foreground")
	separatorClass = lib.Cn("-mx-1 h-px bg-border")
	itemClass      = lib.Cn("relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50 data-[selected=true]:bg-accent data-[selected=true]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground")
	shortcutClass  = lib.Cn("ml-auto text-xs tracking-widest text-muted-foreground")
)

// cmdk's visually hidden <label> style (srOnlyStyles).
const srOnlyStyle = "position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0, 0, 0, 0);white-space:nowrap;border-width:0"

// --- options -----------------------------------------------------------------------------------

// Label is the accessible label: on Command, the text of cmdk's hidden
// <label> for the input; on List, its aria-label (default "Suggestions").
type Label string

// DefaultValue selects the item with this value on load (default: the
// first enabled item).
type DefaultValue string

// Loop makes arrow keys wrap around the list ends.
type Loop bool

// ShouldFilter(false) turns off cmdk's filtering and sorting.
type ShouldFilter bool

// VimBindings(false) turns off Ctrl+N/J/P/K navigation.
type VimBindings bool

// DisablePointerSelection stops pointer movement from selecting items.
type DisablePointerSelection bool

// Value sets an item's or group's value (default: its text content, or
// the group's heading). Values are matched by the filter.
type Value string

// Keywords are extra words an item matches.
type Keywords []string

// Disabled disables an item: it cannot be selected.
type Disabled bool

// ForceMount keeps an item or group rendered whatever the search.
type ForceMount bool

// AlwaysRender keeps a separator rendered while searching.
type AlwaysRender bool

// Heading is a group's heading.
type Heading string

// ShowCloseButton sets whether Dialog shows the corner close button
// (default true).
type ShowCloseButton bool

// Title and Description are Dialog's screen-reader title and description
// (default "Command Palette" and "Search for a command to run...").
type (
	Title       string
	Description string
)

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// --- context -----------------------------------------------------------------------------------

type ctx struct {
	inputID, listID, labelID string
	value                    string // selected value ("" = none)
}

type scoped struct{ c *ctx }

func (s *scoped) setCtx(c *ctx) { s.c = c }

func (s *scoped) ctx() *ctx {
	if s.c == nil {
		s.c = &ctx{inputID: uid.New("cmdk"), listID: uid.New("cmdk"), labelID: uid.New("cmdk")}
	}
	return s.c
}

type participant interface{ setCtx(*ctx) }

// --- Command -----------------------------------------------------------------------------------

// RootBuilder is returned by Command.
type RootBuilder struct {
	h.Builder[RootBuilder]
	label        string
	defaultValue *string
	loop         bool
	shouldFilter bool
	vim          bool
	noPointer    bool
	extraClass   string
}

// Command is the root: the menu's container, holding its Input and List.
// Accepts Label, DefaultValue, Loop, ShouldFilter, VimBindings and
// DisablePointerSelection.
func Command(children ...any) *RootBuilder {
	b := &RootBuilder{shouldFilter: true, vim: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the root options.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Label:
		b.label = string(v)
	case DefaultValue:
		s := string(v)
		b.defaultValue = &s
	case Loop:
		b.loop = bool(v)
	case ShouldFilter:
		b.shouldFilter = bool(v)
	case VimBindings:
		b.vim = bool(v)
	case DisablePointerSelection:
		b.noPointer = bool(v)
	default:
		return false
	}
	return true
}

// Element resolves the root.
func (b *RootBuilder) Element() *h.Element {
	c := &ctx{inputID: uid.New("cmdk"), listID: uid.New("cmdk"), labelID: uid.New("cmdk")}
	kids := b.Elem().ChildNodes()
	first, found := "", false
	tree.Walk(kids, func(n any) bool {
		switch v := n.(type) {
		case *RootBuilder:
			return false
		case *ItemBuilder:
			v.setCtx(c)
			if !found && !v.disabled {
				first, found = v.value(), true
			}
		case participant:
			v.setCtx(c)
		}
		return true
	})
	if b.defaultValue != nil {
		c.value = *b.defaultValue
	} else {
		c.value = first
	}
	el := h.Spread(h.Div(h.Data("slot", "command"), h.Class(rootClass)), b.Elem(), lib.Cn)
	if b.extraClass != "" {
		el.SetClass(lib.Cn(el.ClassName(), b.extraClass))
	}
	if _, ok := el.GetAttr("tabindex"); !ok {
		el.Attr("tabindex", "-1")
	}
	el.Attr("cmdk-root", "")
	if b.loop {
		el.Attr("data-cmdk-loop", "")
	}
	if !b.shouldFilter {
		el.Attr("data-cmdk-should-filter", "false")
	}
	if !b.vim {
		el.Attr("data-cmdk-vim-bindings", "false")
	}
	if b.noPointer {
		el.Attr("data-cmdk-disable-pointer-selection", "")
	}
	label := h.El("label", h.Attr("cmdk-label", ""), h.For(c.inputID), h.ID(c.labelID), h.Style(srOnlyStyle))
	if b.label != "" {
		label.Children(b.label)
	}
	return el.SetChildren(append([]any{label}, kids...)...)
}

// Render renders the root.
func (b *RootBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Dialog ------------------------------------------------------------------------------------

// DialogBuilder is returned by Dialog.
type DialogBuilder struct {
	h.Builder[DialogBuilder]
	title, description string
	showClose          bool
	opts               []any
}

// Dialog is CommandDialog: a Command inside a Dialog, with a visually
// hidden title and description. Its children are the Command's. Accepts
// Title, Description, ShowCloseButton, the dialog's root options
// (dialog.DefaultOpen, dialog.Modal) and, as h.Class, classes for the
// dialog content.
func Dialog(children ...any) *DialogBuilder {
	b := &DialogBuilder{title: "Command Palette", description: "Search for a command to run...", showClose: true}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Title, Description, ShowCloseButton and dialog
// options.
func (b *DialogBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Title:
		b.title = string(v)
	case Description:
		b.description = string(v)
	case ShowCloseButton:
		b.showClose = bool(v)
	case primitive.Option:
		b.opts = append(b.opts, v)
	default:
		return false
	}
	return true
}

// Render renders the dialog: its sr-only header in place, and the dialog
// content (holding the Command) in the dialog's portal template.
func (b *DialogBuilder) Render(w io.Writer) error {
	cmd := Command(b.Elem().ChildNodes()...)
	cmd.extraClass = dialogCmdClass
	content := dialog.Content(dialog.ShowCloseButton(b.showClose), h.Class(lib.Cn("overflow-hidden p-0", b.Elem().ClassName())), cmd)
	root := dialog.Dialog(append(b.opts,
		dialog.Header(h.Class("sr-only"),
			dialog.Title(b.title),
			dialog.Description(b.description),
		),
		content,
	)...)
	return root.Render(w)
}

// --- Input -------------------------------------------------------------------------------------

// InputBuilder is returned by Input.
type InputBuilder struct {
	h.Builder[InputBuilder]
	scoped
}

// Input is the search field, with a search icon.
func Input(children ...any) *InputBuilder {
	b := &InputBuilder{}
	b.Builder = h.NewBuilder(b, "input")
	return b.Children(children...)
}

// Element resolves the input and its wrapper.
func (b *InputBuilder) Element() *h.Element {
	c := b.ctx()
	in := h.Spread(h.Input(h.Data("slot", "command-input"), h.Class(inputClass)), b.Elem(), lib.Cn)
	in.Attr("cmdk-input", "")
	in.Attr("autocomplete", "off")
	in.Attr("autocorrect", "off")
	in.Attr("spellcheck", "false")
	in.Attr("aria-autocomplete", "list")
	in.Attr("role", "combobox")
	in.Attr("aria-expanded", "true")
	in.Attr("aria-controls", c.listID)
	in.Attr("aria-labelledby", c.labelID)
	in.Attr("id", c.inputID)
	in.Attr("type", "text")
	in.Attr("value", "")
	return h.Div(h.Data("slot", "command-input-wrapper"), h.Class(inputWrapClass),
		lucide.Search(h.Class("size-4 shrink-0 opacity-50")),
		in,
	)
}

// Render renders the input.
func (b *InputBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- List --------------------------------------------------------------------------------------

// ListBuilder is returned by List.
type ListBuilder struct {
	h.Builder[ListBuilder]
	scoped
	label string
}

// List holds the items, groups, separators and the Empty message.
// Accepts Label (its aria-label, default "Suggestions").
func List(children ...any) *ListBuilder {
	b := &ListBuilder{label: "Suggestions"}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Label.
func (b *ListBuilder) AcceptChild(c any) bool {
	if v, ok := c.(Label); ok {
		b.label = string(v)
		return true
	}
	return false
}

// Element resolves the list.
func (b *ListBuilder) Element() *h.Element {
	c := b.ctx()
	el := h.Spread(h.Div(h.Data("slot", "command-list"), h.Class(listClass)), b.Elem(), lib.Cn)
	kids := el.ChildNodes()
	el.Attr("cmdk-list", "")
	el.Attr("role", "listbox")
	el.Attr("tabindex", "-1")
	el.Attr("aria-label", b.label)
	el.Attr("id", c.listID)
	return el.SetChildren(h.Div(h.Attr("cmdk-list-sizer", ""), h.Group(kids...)))
}

// Render renders the list.
func (b *ListBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Empty -------------------------------------------------------------------------------------

// EmptyBuilder is returned by Empty.
type EmptyBuilder struct {
	h.Builder[EmptyBuilder]
	scoped
}

// Empty is shown when no item matches the search.
func Empty(children ...any) *EmptyBuilder {
	b := &EmptyBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the empty message. It is not rendered while items
// match: the behavior mounts it from an inert template.
func (b *EmptyBuilder) Element() *h.Element {
	el := h.Spread(h.Div(h.Data("slot", "command-empty"), h.Class(emptyClass)), b.Elem(), lib.Cn)
	el.Attr("cmdk-empty", "")
	el.Attr("role", "presentation")
	return h.Template(h.Data("parity-ignore", ""), h.Attr("data-cmdk-empty-template", ""), el)
}

// Render renders the empty message.
func (b *EmptyBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Group -------------------------------------------------------------------------------------

// GroupBuilder is returned by Group.
type GroupBuilder struct {
	h.Builder[GroupBuilder]
	scoped
	heading    string
	value      *string
	forceMount bool
}

// Group groups items under a Heading. Accepts Heading, Value and
// ForceMount.
func Group(children ...any) *GroupBuilder {
	b := &GroupBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Heading, Value and ForceMount.
func (b *GroupBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Heading:
		b.heading = string(v)
	case Value:
		s := string(v)
		b.value = &s
	case ForceMount:
		b.forceMount = bool(v)
	default:
		return false
	}
	return true
}

// Element resolves the group.
func (b *GroupBuilder) Element() *h.Element {
	el := h.Spread(h.Div(h.Data("slot", "command-group"), h.Class(groupClass)), b.Elem(), lib.Cn)
	kids := el.ChildNodes()
	el.Attr("cmdk-group", "")
	el.Attr("role", "presentation")
	if b.forceMount {
		el.Attr("data-cmdk-force-mount", "")
	}
	value := "undefined" // cmdk sets data-value={undefined} without value or heading
	if b.value != nil {
		value = strings.TrimSpace(*b.value)
	} else if b.heading != "" {
		value = strings.TrimSpace(b.heading)
	}
	el.Attr("data-value", value)
	items := h.Div(h.Attr("cmdk-group-items", ""), h.Role("group"))
	var out []any
	if b.heading != "" {
		hid := uid.New("cmdk")
		out = append(out, h.Div(h.Attr("cmdk-group-heading", ""), h.Aria("hidden", "true"), h.ID(hid), b.heading))
		items.Attr("aria-labelledby", hid)
	}
	items.Children(kids...)
	return el.SetChildren(append(out, items)...)
}

// Render renders the group.
func (b *GroupBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Separator ---------------------------------------------------------------------------------

// SeparatorBuilder is returned by Separator.
type SeparatorBuilder struct {
	h.Builder[SeparatorBuilder]
	always bool
}

// Separator separates groups. It is hidden while searching unless
// AlwaysRender(true) is passed.
func Separator(children ...any) *SeparatorBuilder {
	b := &SeparatorBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes AlwaysRender.
func (b *SeparatorBuilder) AcceptChild(c any) bool {
	if v, ok := c.(AlwaysRender); ok {
		b.always = bool(v)
		return true
	}
	return false
}

// Element resolves the separator.
func (b *SeparatorBuilder) Element() *h.Element {
	el := h.Spread(h.Div(h.Data("slot", "command-separator"), h.Class(separatorClass)), b.Elem(), lib.Cn)
	el.Attr("cmdk-separator", "")
	el.Attr("role", "separator")
	if b.always {
		el.Attr("data-cmdk-always-render", "")
	}
	return el
}

// Render renders the separator.
func (b *SeparatorBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Item --------------------------------------------------------------------------------------

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	scoped
	val        *string
	keywords   []string
	disabled   bool
	forceMount bool
}

// Item is a selectable option. Accepts Value, Keywords, Disabled and
// ForceMount.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes Value, Keywords, Disabled and ForceMount.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		s := string(v)
		b.val = &s
	case Keywords:
		b.keywords = append(b.keywords, v...)
	case Disabled:
		b.disabled = bool(v)
	case ForceMount:
		b.forceMount = bool(v)
	default:
		return false
	}
	return true
}

// value is cmdk's item value: the Value option, else the text content,
// trimmed.
func (b *ItemBuilder) value() string {
	if b.val != nil {
		return strings.TrimSpace(*b.val)
	}
	return strings.TrimSpace(textContent(b.Elem().ChildNodes()))
}

// Element resolves the item.
func (b *ItemBuilder) Element() *h.Element {
	c := b.ctx()
	v := b.value()
	selected := c.value != "" && c.value == v
	el := h.Spread(h.Div(h.Data("slot", "command-item"), h.Class(itemClass)), b.Elem(), lib.Cn)
	el.Attr("id", uid.New("cmdk"))
	el.Attr("cmdk-item", "")
	el.Attr("role", "option")
	el.Attr("aria-disabled", boolStr(b.disabled))
	el.Attr("aria-selected", boolStr(selected))
	el.Attr("data-disabled", boolStr(b.disabled))
	el.Attr("data-selected", boolStr(selected))
	el.Attr("data-value", v)
	if len(b.keywords) > 0 {
		kw, _ := json.Marshal(b.keywords)
		el.Attr("data-cmdk-keywords", string(kw))
	}
	if b.forceMount {
		el.Attr("data-cmdk-force-mount", "")
	}
	return el
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

// --- Shortcut ----------------------------------------------------------------------------------

// Shortcut shows an item's keyboard shortcut.
func Shortcut(children ...any) *h.Element {
	return h.Spread(h.Span(h.Data("slot", "command-shortcut"), h.Class(shortcutClass)), h.Span(children...), lib.Cn)
}

// --- helpers -----------------------------------------------------------------------------------

// textContent returns the text the nodes render, like the DOM's
// textContent.
func textContent(nodes []any) string {
	var buf bytes.Buffer
	if err := h.Group(nodes...).Render(&buf); err != nil {
		return ""
	}
	src := buf.String()
	var out strings.Builder
	inTag := false
	var quote byte
	for i := 0; i < len(src); i++ {
		ch := src[i]
		switch {
		case inTag && quote != 0:
			if ch == quote {
				quote = 0
			}
		case inTag:
			if ch == '"' || ch == '\'' {
				quote = ch
			} else if ch == '>' {
				inTag = false
			}
		case ch == '<':
			inTag = true
		default:
			out.WriteByte(ch)
		}
	}
	return html.UnescapeString(out.String())
}
