// Package questionnaire is the shadcn/ui Questionnaire: a multi-step form
// asking one question at a time, with choices (single or multiple), a
// free-text answer, keyboard shortcuts, validation and skipping. The form
// is a native <form>: read the answers with FormData on submit.
//
//	questionnaire.Questionnaire(questionnaire.DefaultItem("direction"), questionnaire.ShortcutsLetters,
//		questionnaire.Progress(),
//		questionnaire.Item(questionnaire.Name("direction"), questionnaire.Required(true),
//			questionnaire.Title("What should the agent build next?"),
//			questionnaire.Choices(
//				questionnaire.Choice(questionnaire.Value("tool-calls"), "Tool call timeline"),
//				questionnaire.Choice(questionnaire.Value("approvals"), "Approval checkpoints"),
//			),
//			questionnaire.Error(),
//		),
//		questionnaire.Actions(questionnaire.Previous(), questionnaire.Skip(), questionnaire.Next(), questionnaire.Submit()),
//	)
//
// It renders registry/bases/radix/ui/questionnaire.tsx in the style of an
// enclosing style.Scope (the default style otherwise), with the state the
// @shadcn/react/questionnaire primitive has after hydration. The behavior
// (questionnaire.js) ports the primitive: navigation, validation, Skip,
// Enter / Cmd+Enter / arrows / letter or number shortcuts. Scripts observe
// "shadcn:item-change" (detail.item) and "shadcn:status-change" (on the
// item, detail.status), and drive a questionnaire with
// getQuestionnaire(form).setItem(name).
package questionnaire

import (
	_ "embed"
	"fmt"
	"io"
	"strconv"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/internal/tree"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed questionnaire.js
var script string

func init() { runtime.Module("questionnaire", script) }

// Root options.
type (
	// DefaultItem is the item shown first (the first enabled item when
	// absent or disabled).
	DefaultItem string
	// ActiveItem controls the shown item (upstream's item prop); scripts
	// change it with getQuestionnaire(form).setItem.
	ActiveItem string
	// Shortcuts enables single-key shortcuts for choices.
	Shortcuts string
)

// Shortcut modes.
const (
	ShortcutsLetters Shortcuts = "letters"
	ShortcutsNumbers Shortcuts = "numbers"
)

// ItemDef describes an item for Items (upstream's Root.items): the logical
// order, required and disabled items, and the choices that get shortcuts.
type ItemDef struct {
	Name     string
	Required bool
	Disabled bool
	Choices  []ChoiceDef
}

// ChoiceDef is a choice of an ItemDef.
type ChoiceDef struct {
	Value    string
	Disabled bool
}

// ItemsOption is returned by Items.
type ItemsOption struct{ defs []ItemDef }

// Items declares the questionnaire's items (upstream's Root.items).
func Items(defs ...ItemDef) ItemsOption { return ItemsOption{defs} }

// Item and answer options.
type (
	// Name is an Item's name: the form field of its answers.
	Name string
	// Required makes an Item require an answer (no Skip).
	Required bool
	// Multiple makes an Item's choices checkboxes.
	Multiple bool
	// Disabled disables an Item (it leaves the flow) or a Choice.
	Disabled bool
	// Value is a Choice's value.
	Value string
	// DefaultChecked selects a Choice initially.
	DefaultChecked bool
)

// RenderOption is returned by Render.
type RenderOption struct{ el any }

// Render renders a part as el (upstream's render prop): the part's
// attributes are merged into el, e.g. Title(questionnaire.Render(card.Title())).
func Render(el any) RenderOption { return RenderOption{el} }

var file = style.NewFile("ui/questionnaire.tsx")

var (
	sRoot        = file.Cn("cn-questionnaire flex w-full min-w-0 flex-col", style.Arg)
	sProgress    = file.Cn("cn-questionnaire-progress min-h-[1lh] w-fit min-w-[14ch] font-medium text-muted-foreground tabular-nums", style.Arg)
	sItem        = file.Cn("cn-questionnaire-item min-w-0 border-0 p-0 outline-none", style.Arg)
	sTitle       = file.Cn("cn-questionnaire-title cn-font-heading text-pretty", style.Arg)
	sDescription = file.Cn("cn-questionnaire-description text-pretty text-muted-foreground", style.Arg)
	sChoices     = file.Cn("cn-questionnaire-choices group/questionnaire-choices grid min-w-0", style.Arg)
	sChoice      = file.Cn("cn-questionnaire-choice group/questionnaire-choice relative flex min-h-11 cursor-pointer items-start text-start transition-colors outline-none select-none", "data-disabled:pointer-events-none data-disabled:cursor-not-allowed data-disabled:opacity-50", style.Arg)
	sChoiceInput = file.Attr("cn-questionnaire-choice-input absolute inset-0 z-10 size-full cursor-pointer opacity-0")
	sIndicator   = file.Attr("cn-questionnaire-choice-indicator pointer-events-none relative flex shrink-0 items-center justify-center border group-data-[type=radio]/questionnaire-choice:rounded-full")
	sDot         = file.Attr("cn-questionnaire-choice-indicator-dot hidden rounded-full group-data-[type=checkbox]/questionnaire-choice:hidden group-data-checked/questionnaire-choice:block")
	sCheck       = file.Attr("cn-questionnaire-choice-indicator-check hidden group-data-[type=radio]/questionnaire-choice:hidden group-data-checked/questionnaire-choice:block")
	sChoiceLabel = file.Attr("cn-questionnaire-choice-label cn-questionnaire-choice-content flex min-w-0 flex-1 flex-col leading-snug")
	sShortcut    = file.Attr("cn-questionnaire-choice-shortcut cn-questionnaire-shortcut pointer-events-none ms-auto hidden shrink-0 group-data-[shortcut]/questionnaire-choice:inline-flex")
	sChoiceDesc  = file.Cn("cn-questionnaire-choice-description", style.Arg)
	sInputWrap   = file.Attr("cn-questionnaire-input-wrapper group/questionnaire-input relative min-w-0")
	sInput       = file.Cn("cn-questionnaire-input min-h-11 w-full min-w-0 transition-[color,box-shadow,background-color] outline-none disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 sm:min-h-0", "selection:bg-primary selection:text-primary-foreground placeholder:text-muted-foreground", style.Arg)
	sError       = file.Cn("cn-questionnaire-error text-destructive", style.Arg)
	sActions     = file.Cn("cn-questionnaire-actions grid min-h-11 w-full grid-cols-[minmax(0,1fr)_auto_auto] items-center", style.Arg)
	sPrevious    = file.Cn(style.Arg, "cn-questionnaire-previous col-start-1 row-start-1 min-h-11 justify-self-start sm:min-h-0", style.Arg)
	sSkip        = file.Cn(style.Arg, "cn-questionnaire-skip col-start-2 row-start-1 min-h-11 justify-self-end sm:min-h-0", style.Arg)
	sNext        = file.Cn(style.Arg, "cn-questionnaire-next col-start-3 row-start-1 min-h-11 justify-self-end sm:min-h-0", style.Arg)
	sSubmit      = file.Cn(style.Arg, "cn-questionnaire-submit col-start-3 row-start-1 min-h-11 justify-self-end sm:min-h-0", style.Arg)
)

// --- state -------------------------------------------------------------------------------------

type rootState struct {
	shortcuts     Shortcuts
	defs          []ItemDef
	hasDefs       bool
	active        string
	current       int
	total         int
	first, last   bool
	activeItem    *ItemBuilder
	activeRequired *bool
}

type answer struct {
	choice  *ChoiceBuilder
	input   *InputBuilder
	enabled bool
}

// participant is a part that receives the root's state.
type participant interface{ setRoot(*rootState) }

type rooted struct{ root *rootState }

func (r *rooted) setRoot(s *rootState) { r.root = s }

func (r *rooted) state() *rootState {
	if r.root == nil {
		return &rootState{}
	}
	return r.root
}

func styleOf(w io.Writer) style.Style {
	if w != nil {
		if st, ok := style.FromWriter(w); ok {
			return st
		}
	}
	return style.Default()
}

func boolAttr(e *h.Element, name string, on bool) {
	if on {
		e.Attr(name, "")
	}
}

// --- Root --------------------------------------------------------------------------------------

// RootBuilder is returned by Questionnaire.
type RootBuilder struct {
	h.Builder[RootBuilder]
	defaultItem string
	activeItem  string
	shortcuts   Shortcuts
	defs        []ItemDef
	hasDefs     bool
}

// Questionnaire is the root <form>. Children may include DefaultItem,
// ActiveItem, Shortcuts and Items.
func Questionnaire(children ...any) *RootBuilder {
	b := &RootBuilder{}
	b.Builder = h.NewBuilder(b, "form")
	return b.Children(children...)
}

// AcceptChild consumes the root options.
func (b *RootBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case DefaultItem:
		b.defaultItem = string(v)
	case ActiveItem:
		b.activeItem = string(v)
	case Shortcuts:
		b.shortcuts = v
	case ItemsOption:
		b.defs, b.hasDefs = v.defs, true
	default:
		return false
	}
	return true
}

func shortcutKeys(mode Shortcuts) []string {
	var keys []string
	switch mode {
	case ShortcutsLetters:
		for c := 'A'; c <= 'Z'; c++ {
			keys = append(keys, string(c))
		}
	case ShortcutsNumbers:
		for i := 1; i <= 9; i++ {
			keys = append(keys, strconv.Itoa(i))
		}
	}
	return keys
}

// resolve hands the root state to every part and computes it, like the
// primitive after its layout effects ran.
func (b *RootBuilder) resolve() *rootState {
	s := &rootState{shortcuts: b.shortcuts, defs: b.defs, hasDefs: b.hasDefs}
	var items []*ItemBuilder
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *RootBuilder:
			return false
		case *ItemBuilder:
			items = append(items, c)
			c.collect()
		}
		if p, ok := n.(participant); ok {
			p.setRoot(s)
		}
		return true
	})
	byName := map[string]*ItemBuilder{}
	for _, it := range items {
		byName[it.name] = it
		if b.hasDefs {
			for _, d := range b.defs {
				if d.Name == it.name {
					it.def = &d
					break
				}
			}
		}
	}
	// Logical order: the enabled definitions, or the enabled rendered items.
	var logical []string
	if b.hasDefs {
		for _, d := range b.defs {
			if !d.Disabled {
				logical = append(logical, d.Name)
			}
		}
	} else {
		for _, it := range items {
			if !it.disabled {
				logical = append(logical, it.name)
			}
		}
	}
	active := b.activeItem
	if active == "" {
		active = b.defaultItem
	}
	idx := -1
	for i, n := range logical {
		if n == active {
			idx = i
		}
	}
	if idx < 0 && len(logical) > 0 {
		idx, active = 0, logical[0]
	}
	s.active = active
	s.total = len(logical)
	if idx >= 0 {
		s.current = idx + 1
		s.first = idx == 0
		s.last = idx == len(logical)-1
		s.activeItem = byName[active]
		req := false
		if it := byName[active]; it != nil {
			req = it.required
		}
		if b.hasDefs {
			for _, d := range b.defs {
				if d.Name == active {
					req = d.Required
				}
			}
		}
		s.activeRequired = &req
	}
	for _, it := range items {
		it.finish(s)
	}
	return s
}

// Element resolves the root (in the default style).
func (b *RootBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the root as it renders inside w.
func (b *RootBuilder) ElementIn(w io.Writer) *h.Element {
	s := b.resolve()
	base := h.Form(h.Data("slot", "questionnaire"), h.Bool("novalidate"), h.Class(sRoot.Class(styleOf(w), "")))
	if s.total > 0 {
		base.Attr("data-current", strconv.Itoa(s.current))
		boolAttr(base, "data-first", s.first)
		boolAttr(base, "data-last", s.last)
		base.Attr("data-total", strconv.Itoa(s.total))
	}
	if s.shortcuts != "" {
		base.Attr("data-shortcuts", string(s.shortcuts))
	}
	if b.activeItem != "" {
		base.Attr("data-questionnaire-controlled", "")
	}
	if b.hasDefs {
		base.Attr("data-questionnaire-order", strings.Join(logicalNames(b.defs), " "))
	}
	return h.Spread(base, b.Elem(), lib.Cn)
}

func logicalNames(defs []ItemDef) []string {
	var out []string
	for _, d := range defs {
		out = append(out, d.Name)
	}
	return out
}

// Render renders the root.
func (b *RootBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- Progress ----------------------------------------------------------------------------------

// ProgressBuilder is returned by Progress.
type ProgressBuilder struct {
	h.Builder[ProgressBuilder]
	rooted
}

// Progress shows "Question n of total" (role="progressbar"). Children
// replace the label.
func Progress(children ...any) *ProgressBuilder {
	b := &ProgressBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the progress.
func (b *ProgressBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the progress as it renders inside w.
func (b *ProgressBuilder) ElementIn(w io.Writer) *h.Element {
	s := b.state()
	base := h.Div(h.Aria("label", "Questionnaire progress"), h.Aria("live", "polite"), h.Role("progressbar"),
		h.Data("slot", "questionnaire-progress"), h.Class(sProgress.Class(styleOf(w), "")))
	if s.total > 0 {
		label := fmt.Sprintf("Question %d of %d", s.current, s.total)
		base.Aria("valuemax", strconv.Itoa(s.total)).Aria("valuemin", "1").Aria("valuenow", strconv.Itoa(s.current)).Aria("valuetext", label)
		base.Attr("data-current", strconv.Itoa(s.current))
		boolAttr(base, "data-first", s.first)
		boolAttr(base, "data-last", s.last)
		base.Attr("data-total", strconv.Itoa(s.total))
		if len(b.Elem().ChildNodes()) == 0 {
			base.Children(label)
		}
	}
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the progress.
func (b *ProgressBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- Item --------------------------------------------------------------------------------------

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	name                         string
	required, multiple, disabled bool
	def                          *ItemDef
	answers                      []*answer
	descIDs, errIDs              []string
	root                         *rootState
	status                       string
	shortcutByValue              map[string]string
}

// Item is one question (a <fieldset>). Children must include a Name and
// may include Required, Multiple and Disabled.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "fieldset")
	return b.Children(children...)
}

// AcceptChild consumes Name, Required, Multiple and Disabled.
func (b *ItemBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Name:
		b.name = string(v)
	case Required:
		b.required = bool(v)
	case Multiple:
		b.multiple = bool(v)
	case Disabled:
		b.disabled = bool(v)
	default:
		return false
	}
	return true
}

// collect hands the item to its parts and gives them ids.
func (b *ItemBuilder) collect() {
	b.answers, b.descIDs, b.errIDs = nil, nil, nil
	tree.Walk(b.Elem().ChildNodes(), func(n any) bool {
		switch c := n.(type) {
		case *ItemBuilder, *RootBuilder:
			return false
		case *ChoiceBuilder:
			c.item = b
			if c.id == "" {
				c.id = uid.New("questionnaire")
			}
			b.answers = append(b.answers, &answer{choice: c, enabled: !c.disabled && !b.disabled})
		case *InputBuilder:
			c.item = b
			if c.id == "" {
				c.id = uid.New("questionnaire")
			}
			b.answers = append(b.answers, &answer{input: c, enabled: !c.Elem().Has("disabled") && !b.disabled})
		case *DescriptionBuilder:
			c.item = b
			id, ok := c.Elem().GetAttr("id")
			if !ok || id == "" {
				id = uid.New("questionnaire")
				c.Elem().ID(id)
			}
			b.descIDs = append(b.descIDs, id)
		case *ErrorBuilder:
			c.item = b
			id, ok := c.Elem().GetAttr("id")
			if !ok || id == "" {
				id = uid.New("questionnaire")
				c.Elem().ID(id)
			}
			b.errIDs = append(b.errIDs, id)
		case *ChoicesBuilder:
			c.item = b
		}
		return true
	})
}

// finish computes the item's status and the shortcuts of its choices.
func (b *ItemBuilder) finish(s *rootState) {
	b.root = s
	b.status = "unanswered"
	selected := false
	for _, a := range b.answers {
		if !a.enabled {
			continue
		}
		if (a.choice != nil && a.choice.defaultChecked) || (a.input != nil && a.input.filled()) {
			selected = true
		}
	}
	if selected {
		b.status = "answered"
	}
	keys := shortcutKeys(s.shortcuts)
	b.shortcutByValue = map[string]string{}
	if s.hasDefs {
		if b.def != nil {
			i := 0
			for _, c := range b.def.Choices {
				if c.Disabled {
					continue
				}
				if i >= len(keys) {
					break
				}
				b.shortcutByValue[c.Value] = keys[i]
				i++
			}
		}
	} else {
		i := 0
		for _, a := range b.answers {
			if a.choice == nil || !a.enabled {
				continue
			}
			if i >= len(keys) {
				break
			}
			b.shortcutByValue[a.choice.value] = keys[i]
			i++
		}
	}
}

func (b *ItemBuilder) active() bool {
	return b.root != nil && !b.disabled && b.root.active == b.name
}

func (b *ItemBuilder) hasInputAnswer() bool {
	for _, a := range b.answers {
		if a.input != nil && a.enabled {
			return true
		}
	}
	return false
}

func (b *ItemBuilder) hasAnswers() bool {
	for _, a := range b.answers {
		if a.enabled {
			return true
		}
	}
	return false
}

// Element resolves the item.
func (b *ItemBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the item as it renders inside w.
func (b *ItemBuilder) ElementIn(w io.Writer) *h.Element {
	active := b.active()
	base := h.El("fieldset", h.Data("slot", "questionnaire-item"), h.Data("questionnaire-name", b.name), h.TabIndex(-1),
		h.Class(sItem.Class(styleOf(w), "")))
	user := b.Elem()
	describedBy := append([]string{}, b.descIDs...)
	if v, ok := user.GetAttr("aria-describedby"); ok && v != "" {
		describedBy = append(describedBy, v)
	}
	if len(describedBy) > 0 {
		base.Aria("describedby", strings.Join(describedBy, " "))
	}
	var keys []string
	if v, ok := user.GetAttr("aria-keyshortcuts"); ok && v != "" {
		keys = append(keys, v)
	}
	if active {
		keys = append(keys, "Meta+Enter Control+Enter")
		if b.hasAnswers() {
			keys = append(keys, "ArrowUp ArrowDown")
		}
		if !b.root.first {
			keys = append(keys, "ArrowLeft")
		}
		if !b.root.last && b.status != "unanswered" {
			keys = append(keys, "ArrowRight")
		}
	}
	if len(keys) > 0 {
		base.Aria("keyshortcuts", strings.Join(keys, " "))
	}
	if b.disabled {
		base.Bool("disabled", true)
	}
	if !active {
		base.Bool("hidden", true).Bool("inert", true)
	}
	boolAttr(base, "data-active", active)
	boolAttr(base, "data-disabled", b.disabled)
	boolAttr(base, "data-multiple", b.multiple)
	boolAttr(base, "data-required", b.required)
	base.Attr("data-status", b.status)
	props := user.Clone()
	props.Remove("aria-describedby")
	props.Remove("aria-keyshortcuts")
	return h.Spread(base, props, lib.Cn)
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- Title, Description, ChoiceDescription, Actions ------------------------------------------------

// partBuilder is a simple part: its own tag or a Render element.
type partBuilder struct {
	render any
}

func (p *partBuilder) accept(c any) bool {
	if r, ok := c.(RenderOption); ok {
		p.render = r.el
		return true
	}
	return false
}

// out renders el, or the Render element with el's attributes merged in
// (they win, like the props useRender passes to a render element; classes
// merge with the element's own first).
func (p *partBuilder) out(w io.Writer, el *h.Element) *h.Element {
	if p.render == nil {
		return el
	}
	target := h.ResolveIn(w, p.render)
	if target == nil {
		return el
	}
	kids := el.ChildNodes()
	slot := el.Clone().SetChildren()
	return h.Spread(target.Clone().SetChildren(append(append([]any{}, target.ChildNodes()...), kids...)...), slot, lib.Cn)
}

// TitleBuilder is returned by Title.
type TitleBuilder struct {
	h.Builder[TitleBuilder]
	partBuilder
}

// Title is the question (a <legend>). Accepts Render.
func Title(children ...any) *TitleBuilder {
	b := &TitleBuilder{}
	b.Builder = h.NewBuilder(b, "legend")
	return b.Children(children...)
}

// AcceptChild consumes Render.
func (b *TitleBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the title.
func (b *TitleBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the title as it renders inside w.
func (b *TitleBuilder) ElementIn(w io.Writer) *h.Element {
	base := h.El("legend", h.Data("slot", "questionnaire-title"), h.Class(sTitle.Class(styleOf(w), "")))
	return b.out(w, h.Spread(base, b.Elem(), lib.Cn))
}

// Render renders the title.
func (b *TitleBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// DescriptionBuilder is returned by Description.
type DescriptionBuilder struct {
	h.Builder[DescriptionBuilder]
	partBuilder
	item *ItemBuilder
}

// Description describes the question (a <p>, referenced by the item's
// aria-describedby). Accepts Render.
func Description(children ...any) *DescriptionBuilder {
	b := &DescriptionBuilder{}
	b.Builder = h.NewBuilder(b, "p")
	return b.Children(children...)
}

// AcceptChild consumes Render.
func (b *DescriptionBuilder) AcceptChild(c any) bool { return b.accept(c) }

// Element resolves the description.
func (b *DescriptionBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the description as it renders inside w.
func (b *DescriptionBuilder) ElementIn(w io.Writer) *h.Element {
	base := h.P(h.Data("slot", "questionnaire-description"), h.Class(sDescription.Class(styleOf(w), "")))
	return b.out(w, h.Spread(base, b.Elem(), lib.Cn))
}

// Render renders the description.
func (b *DescriptionBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ChoicesBuilder is returned by Choices.
type ChoicesBuilder struct {
	h.Builder[ChoicesBuilder]
	item *ItemBuilder
}

// Choices lays out the answers of an item.
func Choices(children ...any) *ChoicesBuilder {
	b := &ChoicesBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Element resolves the choices.
func (b *ChoicesBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the choices as it renders inside w.
func (b *ChoicesBuilder) ElementIn(w io.Writer) *h.Element {
	base := h.Div(h.Data("slot", "questionnaire-choices"), h.Class(sChoices.Class(styleOf(w), "")))
	if b.item != nil && b.item.root != nil && b.item.root.shortcuts != "" {
		base.Attr("data-shortcuts", string(b.item.root.shortcuts))
	}
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the choices.
func (b *ChoicesBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// ChoiceDescription is secondary text inside a Choice.
func ChoiceDescription(children ...any) h.Node {
	return h.NodeFunc(func(w io.Writer) error {
		base := h.Span(h.Data("slot", "questionnaire-choice-description"), h.Class(sChoiceDesc.Class(styleOf(w), "")))
		return h.Spread(base, h.Span(children...), lib.Cn).Render(w)
	})
}

// Actions lays out Previous, Skip, Next and Submit.
func Actions(children ...any) h.Node {
	props := h.Div(children...)
	return actionsNode{props}
}

type actionsNode struct{ props *h.Element }

func (a actionsNode) Elem() *h.Element { return a.props }

func (a actionsNode) Render(w io.Writer) error {
	base := h.Div(h.Data("slot", "questionnaire-actions"), h.Class(sActions.Class(styleOf(w), "")))
	return h.Spread(base, a.props, lib.Cn).Render(w)
}

// --- Choice ------------------------------------------------------------------------------------

// ChoiceBuilder is returned by Choice.
type ChoiceBuilder struct {
	h.Builder[ChoiceBuilder]
	value          string
	disabled       bool
	defaultChecked bool
	id             string
	item           *ItemBuilder
}

// Choice is an answer (a <label> wrapping a radio or checkbox). Children
// must include a Value and may include Disabled and DefaultChecked.
func Choice(children ...any) *ChoiceBuilder {
	b := &ChoiceBuilder{}
	b.Builder = h.NewBuilder(b, "label")
	return b.Children(children...)
}

// AcceptChild consumes Value, Disabled and DefaultChecked.
func (b *ChoiceBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Value:
		b.value = string(v)
	case Disabled:
		b.disabled = bool(v)
	case DefaultChecked:
		b.defaultChecked = bool(v)
	default:
		return false
	}
	return true
}

// Element resolves the choice.
func (b *ChoiceBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the choice as it renders inside w.
func (b *ChoiceBuilder) ElementIn(w io.Writer) *h.Element {
	st := styleOf(w)
	it := b.item
	if it == nil {
		it = &ItemBuilder{}
	}
	disabled := it.disabled || b.disabled
	checked := b.defaultChecked
	typ := "radio"
	if it.multiple {
		typ = "checkbox"
	}
	shortcut := it.shortcutByValue[b.value]
	invalid := false
	stateAttrs := func(e *h.Element) {
		boolAttr(e, "data-checked", checked)
		boolAttr(e, "data-disabled", disabled)
		boolAttr(e, "data-invalid", invalid)
		if shortcut != "" {
			e.Attr("data-shortcut", shortcut)
		}
		e.Attr("data-type", typ)
		boolAttr(e, "data-unchecked", !checked)
	}
	input := h.Input(h.Type(typ), h.ID(b.id), h.Name(it.name), h.Value(b.value),
		h.Data("slot", "questionnaire-choice-input"), h.Class(sChoiceInput.Class(st)))
	if keys := answerKeyShortcuts(shortcut, !disabled && checked); keys != "" {
		input.Aria("keyshortcuts", keys)
	}
	input.Bool("checked", checked)
	input.Bool("disabled", disabled)
	input.Bool("required", it.required && !it.multiple && !it.hasInputAnswer())
	stateAttrs(input)
	short := h.Span(h.Aria("hidden", "true"), h.Data("slot", "questionnaire-choice-shortcut"), h.Class(sShortcut.Class(st)))
	if shortcut != "" {
		short.Attr("data-shortcut", shortcut)
		short.Children(shortcut)
	} else {
		short.Bool("hidden", true)
	}
	base := h.Label(h.Data("slot", "questionnaire-choice"), h.Class(sChoice.Class(st, "")))
	stateAttrs(base)
	props := b.Elem()
	kids := props.ChildNodes()
	base = h.Spread(base, props.Clone().SetChildren(), lib.Cn)
	base.Children(
		input,
		h.Span(h.Aria("hidden", "true"), h.Data("slot", "questionnaire-choice-indicator"), h.Class(sIndicator.Class(st)),
			h.Span(h.Data("slot", "questionnaire-choice-indicator-dot"), h.Class(sDot.Class(st))),
			lucide.Check(h.Data("slot", "questionnaire-choice-indicator-check"), h.Class(sCheck.Class(st))),
		),
		h.Span(h.Data("slot", "questionnaire-choice-label"), h.Class(sChoiceLabel.Class(st)), h.Group(kids...)),
		short,
	)
	return base
}

func answerKeyShortcuts(shortcut string, filled bool) string {
	var keys []string
	if shortcut != "" {
		keys = append(keys, shortcut)
	}
	if filled {
		keys = append(keys, "Enter")
	}
	return strings.Join(keys, " ")
}

// Render renders the choice.
func (b *ChoiceBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- Input -------------------------------------------------------------------------------------

// InputBuilder is returned by Input.
type InputBuilder struct {
	h.Builder[InputBuilder]
	id   string
	item *ItemBuilder
}

// Input is a free-text answer. Attributes go to the <input> (h.Value sets
// its default value).
func Input(children ...any) *InputBuilder {
	b := &InputBuilder{}
	b.Builder = h.NewBuilder(b, "input")
	return b.Children(children...)
}

func (b *InputBuilder) filled() bool {
	v, _ := b.Elem().GetAttr("value")
	return strings.TrimSpace(v) != ""
}

// Element resolves the input's wrapper.
func (b *InputBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the input's wrapper as it renders inside w.
func (b *InputBuilder) ElementIn(w io.Writer) *h.Element {
	st := styleOf(w)
	it := b.item
	if it == nil {
		it = &ItemBuilder{}
	}
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	disabled := it.disabled || props.Has("disabled")
	filled := b.filled()
	typ, _ := props.GetAttr("type")
	if typ == "" {
		typ = "text"
	}
	input := h.Input(h.Type(typ), h.ID(b.id), h.Data("slot", "questionnaire-input"), h.Class(sInput.Class(st, className)))
	if filled {
		input.Attr("name", it.name)
		input.Attr("data-filled", "")
		input.Aria("keyshortcuts", "Enter")
	} else {
		input.Attr("form", "")
		input.Attr("data-empty", "")
	}
	input.Bool("disabled", disabled)
	boolAttr(input, "data-disabled", disabled)
	input = h.Spread(input, props, nil)
	return h.Div(h.Data("slot", "questionnaire-input-wrapper"), h.Class(sInputWrap.Class(st)), input)
}

// Render renders the input.
func (b *InputBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- Error -------------------------------------------------------------------------------------

// ErrorBuilder is returned by Error.
type ErrorBuilder struct {
	h.Builder[ErrorBuilder]
	item *ItemBuilder
}

// Error is the item's validation message, shown after a failed Next or
// Submit. Children replace the default message.
func Error(children ...any) *ErrorBuilder {
	b := &ErrorBuilder{}
	b.Builder = h.NewBuilder(b, "p")
	return b.Children(children...)
}

// Element resolves the error.
func (b *ErrorBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the error as it renders inside w.
func (b *ErrorBuilder) ElementIn(w io.Writer) *h.Element {
	base := h.P(h.Data("slot", "questionnaire-error"), h.Bool("hidden"), h.Class(sError.Class(styleOf(w), "")))
	if len(b.Elem().ChildNodes()) == 0 {
		msg := "Choose an answer or skip this question."
		if b.item != nil && b.item.required {
			msg = "Choose an answer to continue."
		}
		base.Children(msg)
	}
	return h.Spread(base, b.Elem(), lib.Cn)
}

// Render renders the error.
func (b *ErrorBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }

// --- navigation ----------------------------------------------------------------------------------

type navKind int

const (
	navPrevious navKind = iota
	navSkip
	navNext
	navSubmit
)

// NavBuilder is returned by Previous, Skip, Next and Submit.
type NavBuilder struct {
	h.Builder[NavBuilder]
	rooted
	kind    navKind
	variant button.Variant
	size    button.Size
}

func nav(kind navKind, variant button.Variant, children []any) *NavBuilder {
	b := &NavBuilder{kind: kind, variant: variant, size: button.SizeDefault}
	b.Builder = h.NewBuilder(b, "button")
	return b.Children(children...)
}

// Previous goes to the previous item (an outline Button by default).
func Previous(children ...any) *NavBuilder { return nav(navPrevious, button.VariantOutline, children) }

// Skip skips an optional item.
func Skip(children ...any) *NavBuilder { return nav(navSkip, button.VariantOutline, children) }

// Next validates the item and goes to the next one.
func Next(children ...any) *NavBuilder { return nav(navNext, button.VariantDefault, children) }

// Submit submits the form, on the last item.
func Submit(children ...any) *NavBuilder { return nav(navSubmit, button.VariantDefault, children) }

// AcceptChild consumes button.Variant and button.Size.
func (b *NavBuilder) AcceptChild(c any) bool {
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
func (b *NavBuilder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the button as it renders inside w.
func (b *NavBuilder) ElementIn(w io.Writer) *h.Element {
	st := styleOf(w)
	s := b.state()
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	var visible bool
	var slot, label, typ string
	var site *style.Site
	shortcut := false
	switch b.kind {
	case navPrevious:
		visible, slot, label, typ, site = s.total > 1 && !s.first, "questionnaire-previous", "Previous", "button", sPrevious
	case navSkip:
		visible, slot, label, typ, site = s.activeRequired != nil && !*s.activeRequired, "questionnaire-skip", "Skip", "button", sSkip
	case navNext:
		visible, slot, label, typ, site, shortcut = s.total > 1 && !s.last, "questionnaire-next", "Next", "button", sNext, true
	case navSubmit:
		visible, slot, label, typ, site, shortcut = s.total > 0 && s.last, "questionnaire-submit", "Submit", "submit", sSubmit, true
	}
	disabled := props.Has("disabled")
	base := h.Button(h.Type(typ), h.Data("slot", slot), h.Data("size", string(b.size)), h.Data("variant", string(b.variant)),
		h.Class(site.Class(st, button.StyleVariants(st, b.variant, b.size), className)))
	if !visible {
		base.Aria("hidden", "true").Bool("hidden", true).Bool("inert", true).TabIndex(-1)
		base.Attr("data-hidden", "")
	} else {
		base.Attr("data-visible", "")
	}
	if shortcut && visible && !disabled {
		base.Aria("keyshortcuts", "Enter")
		base.Attr("data-shortcut", "Enter")
	}
	boolAttr(base, "data-disabled", disabled)
	if s.activeItem != nil {
		base.Attr("data-status", s.activeItem.status)
	}
	base.Attr("data-questionnaire-nav", map[navKind]string{navPrevious: "previous", navSkip: "skip", navNext: "next", navSubmit: "submit"}[b.kind])
	if len(props.ChildNodes()) == 0 {
		base.Children(label)
	}
	return h.Spread(base, props, nil)
}

// Render renders the button.
func (b *NavBuilder) Render(w io.Writer) error { return b.ElementIn(w).Render(w) }
