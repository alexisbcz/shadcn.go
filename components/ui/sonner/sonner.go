// Package sonner is the shadcn/ui Sonner: an opinionated toast component
// (sonner 2.0.7). Render one Toaster per page (typically in the layout),
// then show toasts from JavaScript with the module's toast function, or
// declaratively with OnClick:
//
//	sonner.Toaster()
//
//	import { toast } from "/_shadcn/sonner.js"
//	toast("Event has been created", { description: "Sunday at 9:00 AM", action: { label: "Undo", onClick() {} } })
//	toast.success("Saved")
//
//	button.Button(sonner.OnClick(sonner.Success("Saved")), "Save")
//
// The Toaster renders sonner's <section> (empty until a toast shows) and a
// <template> carrying its options and icons; the behavior (sonner.js) is a
// port of sonner's React components: same markup, stacking, expansion on
// hover, timers, swipe to dismiss, promise toasts and the same injected
// stylesheet.
package sonner

import (
	_ "embed"
	"encoding/json"
	"io"
	stylepkg "shadcn.go/components/style"
	"strconv"
	"strings"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed sonner.js
var script string

func init() { runtime.Module("sonner", script) }

// Position is where the toasts stack.
type Position string

// Positions.
const (
	PositionTopLeft      Position = "top-left"
	PositionTopCenter    Position = "top-center"
	PositionTopRight     Position = "top-right"
	PositionBottomLeft   Position = "bottom-left"
	PositionBottomCenter Position = "bottom-center"
	PositionBottomRight  Position = "bottom-right"
)

// Theme is the toaster's color theme (sonner's own colors; shadcn's
// Toaster overrides the normal colors with the page theme's popover
// colors).
type Theme string

// Themes.
const (
	ThemeLight  Theme = "light"
	ThemeDark   Theme = "dark"
	ThemeSystem Theme = "system"
)

// Option configures a Toaster; pass options as children.
type Option func(*config)

type config struct {
	theme              Theme
	position           Position
	richColors         *bool
	expand             *bool
	closeButton        *bool
	invert             *bool
	duration           time.Duration
	visibleToasts      int
	gap                int
	offset             any
	mobileOffset       any
	dir                string
	hotkey             []string
	containerAriaLabel string
	id                 string
	swipeDirections    []string
	toastOptions       map[string]any
	plainIcons         bool
}

func boolp(v bool) *bool { return &v }

// Options.
var (
	// RichColors colors success, info, warning and error toasts.
	RichColors Option = func(c *config) { c.richColors = boolp(true) }
	// Expand shows every toast expanded, not stacked.
	Expand Option = func(c *config) { c.expand = boolp(true) }
	// CloseButton adds a close button to every toast.
	CloseButton Option = func(c *config) { c.closeButton = boolp(true) }
	// Invert inverts the toast colors.
	Invert Option = func(c *config) { c.invert = boolp(true) }
	// SonnerIcons uses sonner's built-in icons instead of shadcn's lucide
	// icons.
	SonnerIcons Option = func(c *config) { c.plainIcons = true }
)

// Duration sets how long toasts stay (default 4s).
func Duration(d time.Duration) Option { return func(c *config) { c.duration = d } }

// VisibleToasts sets how many toasts are visible at once (default 3).
func VisibleToasts(n int) Option { return func(c *config) { c.visibleToasts = n } }

// Gap sets the gap between expanded toasts in pixels (default 14).
func Gap(px int) Option { return func(c *config) { c.gap = px } }

// Offset sets the toaster's distance from the viewport edges (CSS length,
// default "24px").
func Offset(v string) Option { return func(c *config) { c.offset = v } }

// MobileOffset sets the offset on screens up to 600px wide (default "16px").
func MobileOffset(v string) Option { return func(c *config) { c.mobileOffset = v } }

// Dir sets the text direction ("ltr", "rtl", "auto"; default: the document's).
func Dir(v string) Option { return func(c *config) { c.dir = v } }

// Hotkey sets the keys that focus and expand the toasts (default
// "altKey", "KeyT"): KeyboardEvent property names or codes.
func Hotkey(keys ...string) Option { return func(c *config) { c.hotkey = keys } }

// ContainerAriaLabel labels the toasts region (default "Notifications").
func ContainerAriaLabel(v string) Option { return func(c *config) { c.containerAriaLabel = v } }

// ID makes the toaster show only the toasts created with toasterId: id.
func ID(id string) Option { return func(c *config) { c.id = id } }

// SwipeDirections sets the directions toasts can be swiped away
// ("top", "right", "bottom", "left"; default: the position's).
func SwipeDirections(dirs ...string) Option { return func(c *config) { c.swipeDirections = dirs } }

// ToastOptions sets sonner's toastOptions (className, descriptionClassName,
// style, classNames, closeButton, unstyled, duration, ...), as JSON-able
// values.
func ToastOptions(opts map[string]any) Option { return func(c *config) { c.toastOptions = opts } }

// ToasterBuilder is returned by Toaster.
type ToasterBuilder struct {
	h.Builder[ToasterBuilder]
	cfg config
}

// Toaster is shadcn's <Toaster />: sonner's Toaster with lucide icons,
// className "toaster group" and the page's popover colors. Theme defaults
// to ThemeSystem, like shadcn's when next-themes provides no theme. Accepts
// Options, a Position and a Theme; a class and a style attribute are passed
// to the toaster list like sonner's className and style props.
func Toaster(children ...any) *ToasterBuilder {
	b := &ToasterBuilder{cfg: config{theme: ThemeSystem, position: PositionBottomRight}}
	b.Builder = h.NewBuilder(b, "section")
	return b.Children(children...)
}

// AcceptChild consumes Options, Position and Theme.
func (b *ToasterBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Option:
		v(&b.cfg)
	case Position:
		b.cfg.position = v
	case Theme:
		b.cfg.theme = v
	default:
		return false
	}
	return true
}

// shadcn's style prop: the normal colors follow the page theme.
var shadcnStyle = [][2]string{
	{"--normal-bg", "var(--popover)"},
	{"--normal-text", "var(--popover-foreground)"},
	{"--normal-border", "var(--border)"},
	{"--border-radius", "var(--radius)"},
}

// Element resolves the toaster.
func (b *ToasterBuilder) Element() *h.Element { return b.element(false) }

func (b *ToasterBuilder) element(styled bool) *h.Element {
	c := b.cfg
	if styled && c.toastOptions == nil {
		// The bases' toastOptions={{ classNames: { toast: "cn-toast" } }}
		// (kept as is by the style transform), before {...props}.
		c.toastOptions = map[string]any{"classNames": map[string]any{"toast": "cn-toast"}}
	}
	user := b.Elem()
	className := "toaster group"
	if uc := user.ClassName(); uc != "" {
		className = uc // {...props} after className, like upstream
	}
	style := map[string]string{}
	var styleKeys []string
	for _, kv := range shadcnStyle {
		style[kv[0]] = kv[1]
		styleKeys = append(styleKeys, kv[0])
	}
	if us, ok := user.GetAttr("style"); ok {
		style = map[string]string{}
		styleKeys = nil
		for _, decl := range strings.Split(us, ";") {
			k, v, ok := strings.Cut(decl, ":")
			if !ok {
				continue
			}
			k, v = strings.TrimSpace(k), strings.TrimSpace(v)
			style[k] = v
			styleKeys = append(styleKeys, k)
		}
	}
	styleJSON := orderedJSON(styleKeys, style)

	hotkey := c.hotkey
	if hotkey == nil {
		hotkey = []string{"altKey", "KeyT"}
	}
	label := c.containerAriaLabel
	if label == "" {
		label = "Notifications"
	}
	hotkeyLabel := strings.NewReplacer("Key", "", "Digit", "").Replace(strings.Join(hotkey, "+"))

	jsonAttr := func(name string, v any) h.Prop {
		bs, _ := json.Marshal(v)
		return h.Data(name, string(bs))
	}
	boolAttr := func(name string, v *bool) any {
		if v == nil {
			return nil
		}
		return h.Data(name, strconv.FormatBool(*v))
	}
	var icons []any
	if !c.plainIcons {
		icons = []any{
			lucide.CircleCheck(h.Class("size-4"), h.Data("sonner-icon", "success")),
			lucide.Info(h.Class("size-4"), h.Data("sonner-icon", "info")),
			lucide.TriangleAlert(h.Class("size-4"), h.Data("sonner-icon", "warning")),
			lucide.OctagonX(h.Class("size-4"), h.Data("sonner-icon", "error")),
			lucide.LoaderCircle(h.Class("size-4 animate-spin"), h.Data("sonner-icon", "loading")),
		}
	}
	tpl := h.Template(
		h.Data("parity-ignore", ""),
		h.Data("sonner-toaster", ""),
		h.Data("theme", string(c.theme)),
		h.Data("position", string(c.position)),
		h.Data("class", className),
		h.Data("style", styleJSON),
		boolAttr("rich-colors", c.richColors),
		boolAttr("expand", c.expand),
		boolAttr("close-button", c.closeButton),
		boolAttr("invert", c.invert),
		h.If(c.duration > 0, h.Data("duration", strconv.FormatInt(c.duration.Milliseconds(), 10))),
		h.If(c.visibleToasts > 0, h.Data("visible-toasts", strconv.Itoa(c.visibleToasts))),
		h.If(c.gap > 0, h.Data("gap", strconv.Itoa(c.gap))),
		h.If(c.offset != nil, jsonAttr("offset", c.offset)),
		h.If(c.mobileOffset != nil, jsonAttr("mobile-offset", c.mobileOffset)),
		h.If(c.dir != "", h.Data("dir", c.dir)),
		h.If(c.hotkey != nil, jsonAttr("hotkey", c.hotkey)),
		h.If(c.containerAriaLabel != "", h.Data("container-aria-label", c.containerAriaLabel)),
		h.If(c.id != "", h.Data("id", c.id)),
		h.If(c.swipeDirections != nil, jsonAttr("swipe-directions", c.swipeDirections)),
		h.If(c.toastOptions != nil, jsonAttr("toast-options", c.toastOptions)),
		h.Group(icons...),
	)
	return h.El("section",
		h.Aria("label", label+" "+hotkeyLabel),
		h.TabIndex(-1),
		h.Aria("live", "polite"),
		h.Aria("relevant", "additions text"),
		h.Aria("atomic", "false"),
		tpl,
	)
}

// Render renders the toaster.
func (b *ToasterBuilder) Render(w io.Writer) error {
	_, styled := stylepkg.FromWriter(w)
	return b.element(styled).Render(w)
}

// orderedJSON encodes a string map as a JSON object with keys in order.
func orderedJSON(keys []string, m map[string]string) string {
	var sb strings.Builder
	sb.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		vb, _ := json.Marshal(m[k])
		sb.Write(kb)
		sb.WriteByte(':')
		sb.Write(vb)
	}
	sb.WriteByte('}')
	return sb.String()
}

// --- declarative toasts ------------------------------------------------------------------------

// Spec describes a toast for OnClick.
type Spec struct {
	Type        string       `json:"type,omitempty"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
	Action      string       `json:"action,omitempty"` // label; clicking dispatches "shadcn:toast-action" on document
	Cancel      string       `json:"cancel,omitempty"` // label; clicking dispatches "shadcn:toast-cancel" on document
	DurationMS  int          `json:"duration,omitempty"`
	ID          string       `json:"id,omitempty"`
	CloseButton *bool        `json:"closeButton,omitempty"`
	RichColors  *bool        `json:"richColors,omitempty"`
	Promise     *PromiseSpec `json:"promise,omitempty"`
}

// PromiseSpec is a promise toast that settles after MS milliseconds.
type PromiseSpec struct {
	Loading string `json:"loading"`
	Success string `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	MS      int    `json:"ms,omitempty"`
	Reject  bool   `json:"reject,omitempty"`
}

// Message is a default toast.
func Message(msg string) Spec { return Spec{Message: msg} }

// Success is a success toast.
func Success(msg string) Spec { return Spec{Type: "success", Message: msg} }

// Info is an info toast.
func Info(msg string) Spec { return Spec{Type: "info", Message: msg} }

// Warning is a warning toast.
func Warning(msg string) Spec { return Spec{Type: "warning", Message: msg} }

// Error is an error toast.
func Error(msg string) Spec { return Spec{Type: "error", Message: msg} }

// Loading is a loading toast.
func Loading(msg string) Spec { return Spec{Type: "loading", Message: msg} }

// WithDescription returns s with a description.
func (s Spec) WithDescription(d string) Spec { s.Description = d; return s }

// WithAction returns s with an action button.
func (s Spec) WithAction(label string) Spec { s.Action = label; return s }

// OnClick returns the attribute that shows the toast when the element is
// clicked.
func OnClick(s Spec) h.Prop {
	bs, _ := json.Marshal(s)
	return h.Data("sonner-trigger", string(bs))
}
