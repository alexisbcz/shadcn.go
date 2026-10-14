// Package chart is the shadcn/ui Chart: recharts charts styled by a chart
// config, with shadcn's tooltip and legend content.
//
// Upstream, ChartContainer wraps recharts elements. Here the recharts
// element tree is described with the same names and props, and the client
// runtime (chart.js, a port of recharts 3.8.0) renders the exact SVG
// recharts renders, at the container's measured size:
//
//	config := chart.Config{
//		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
//	}
//	chart.Container(config, chart.Nonce(secure.Nonce(r)),
//		chart.AreaChart(chart.Data(rows), chart.Margin{Left: 12, Right: 12},
//			chart.CartesianGrid(chart.Vertical(false)),
//			chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.AxisLine(false),
//				chart.TickMargin(8), chart.TickFormatter(chart.Slice(0, 3))),
//			chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
//			chart.Area(chart.DataKey("desktop"), chart.Type("natural"),
//				chart.Fill("var(--color-desktop)"), chart.FillOpacity(0.4), chart.Stroke("var(--color-desktop)")),
//		),
//	)
//
// Props use recharts' names (P("dataKey", "month") is DataKey("month"));
// function props (tick formatters, custom shapes and dots) are named
// functions of the runtime (Slice, ConfigLabel, DateFormat, ...), see Fn.
// Other chart types register their element names with the runtime (see
// chart.js: registerChart, registerFn).
//
// The container's <style> sets the --color-<key> variables of the config;
// under a nonce-based CSP pass the request's nonce with Nonce (otherwise
// the runtime recreates the element with the page's nonce).
package chart

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed chart.js
var script string

func init() { runtime.Module("chart", script) }

// Theme sets a config item's color per theme.
type Theme struct {
	Light string
	Dark  string
}

// Icon renders an icon (any lucide or tabler icon function).
type Icon func(children ...any) *h.Element

// ConfigItem is one entry of a chart config: the label and color (or
// per-theme colors) of a data key, and optionally an icon shown in the
// legend and tooltip.
type ConfigItem struct {
	Key   string
	Label string
	Color string
	Theme *Theme
	Icon  Icon
}

// Config is shadcn's ChartConfig, in declaration order.
type Config []ConfigItem

// Get returns the item for key.
func (c Config) Get(key string) (ConfigItem, bool) {
	for _, it := range c {
		if it.Key == key {
			return it, true
		}
	}
	return ConfigItem{}, false
}

var (
	safeKey   = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	safeColor = regexp.MustCompile(`^[^;{}<>\\\n\r"'` + "`" + `]*$`)
)

// css returns ChartStyle's stylesheet for the chart id, exactly as
// upstream builds it.
func (c Config) css(id string) (string, error) {
	var items []ConfigItem
	for _, it := range c {
		if it.Theme != nil || it.Color != "" {
			if !safeKey.MatchString(it.Key) {
				return "", fmt.Errorf("chart: config key %q: only letters, digits, - and _ are allowed", it.Key)
			}
			items = append(items, it)
		}
	}
	if len(items) == 0 {
		return "", nil
	}
	themes := []struct{ name, prefix string }{{"light", ""}, {"dark", ".dark"}}
	blocks := make([]string, 0, len(themes))
	for _, t := range themes {
		lines := make([]string, 0, len(items))
		for _, it := range items {
			color := it.Color
			if it.Theme != nil {
				if t.name == "light" {
					color = it.Theme.Light
				} else {
					color = it.Theme.Dark
				}
			}
			if color == "" {
				lines = append(lines, "")
				continue
			}
			if !safeColor.MatchString(color) {
				return "", fmt.Errorf("chart: config color %q for %q is not a plain CSS value", color, it.Key)
			}
			lines = append(lines, "  --color-"+it.Key+": "+color+";")
		}
		blocks = append(blocks, "\n"+t.prefix+" [data-chart="+id+"] {\n"+strings.Join(lines, "\n")+"\n}\n")
	}
	return strings.Join(blocks, "\n"), nil
}

// --- spec tree -------------------------------------------------------------------

// Node is a recharts element (AreaChart, XAxis, Bar, LabelList, ...) or an
// SVG element inside the chart (defs, linearGradient, stop).
type Node struct {
	Type     string
	Props    map[string]any
	Children []*Node
}

// Prop is one recharts prop.
type Prop struct {
	Key   string
	Value any
}

// P returns the prop key=value, for props without a dedicated helper.
func P(key string, value any) Prop { return Prop{key, value} }

// Props sets several props at once.
type Props map[string]any

// MarshalJSON encodes the node for the runtime.
func (n *Node) MarshalJSON() ([]byte, error) {
	type out struct {
		Type     string         `json:"type"`
		Props    map[string]any `json:"props,omitempty"`
		Children []*Node        `json:"children,omitempty"`
	}
	return json.Marshal(out{n.Type, n.Props, n.Children})
}

// El returns an element of the given type: a recharts element name or an
// SVG tag (for <defs> content), with props and children.
func El(typ string, args ...any) *Node {
	n := &Node{Type: typ, Props: map[string]any{}}
	n.add(args)
	return n
}

func (n *Node) add(args []any) {
	for _, a := range args {
		switch v := a.(type) {
		case nil:
		case Prop:
			n.Props[v.Key] = v.Value
		case Props:
			for k, x := range v {
				n.Props[k] = x
			}
		case *Node:
			if v != nil {
				n.Children = append(n.Children, v)
			}
		case []*Node:
			for _, c := range v {
				if c != nil {
					n.Children = append(n.Children, c)
				}
			}
		case []any:
			n.add(v)
		case Margin:
			n.Props["margin"] = v.value()
		default:
			panic(fmt.Sprintf("chart: unexpected %s child %T", n.Type, a))
		}
	}
}

// Margin is a chart's margin. Upstream, sides left out are 0 once a margin
// is given; without a margin recharts uses 5 on every side.
type Margin struct {
	Top, Right, Bottom, Left float64
}

func (m Margin) value() map[string]float64 {
	return map[string]float64{"top": m.Top, "right": m.Right, "bottom": m.Bottom, "left": m.Left}
}

// Row is one data entry.
type Row = map[string]any

// Fn is a function prop: a function of the runtime, by name, with
// arguments. Function props upstream are inline closures; the runtime
// provides the ones the shadcn charts use, and registerFn adds more.
type Fn struct {
	Name string
	Args []any
	icon Icon // for IconDot: rendered into a <template>
}

// MarshalJSON encodes the function reference.
func (f Fn) MarshalJSON() ([]byte, error) {
	args := f.Args
	if args == nil {
		args = []any{}
	}
	return json.Marshal(map[string]any{"$fn": f.Name, "args": args})
}

// Field reads a field of the data entry in a shape or dot's props
// (payload.fill upstream is Field("fill")).
type Field string

// MarshalJSON encodes the field reference.
func (f Field) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"$field": string(f)})
}

// Slice is (value) => value.slice(start, end).
func Slice(start, end int) Fn { return Fn{Name: "slice", Args: []any{start, end}} }

// ConfigLabel is (value) => chartConfig[value]?.label.
func ConfigLabel() Fn { return Fn{Name: "configLabel"} }

// DateFormat is (value) => new Date(value).toLocaleDateString(locale, options).
func DateFormat(locale string, options map[string]string) Fn {
	return Fn{Name: "date", Args: []any{locale, options}}
}

// ActiveIndexShape is a Bar shape that renders the Rectangle at index with
// extra props: ({index, ...props}) => index === i ? <Rectangle {...props} {...extra}/> : <Rectangle {...props}/>.
func ActiveIndexShape(index int, extra Props) Fn {
	return Fn{Name: "activeIndexShape", Args: []any{index, map[string]any(extra)}}
}

// DotShape is a Line dot rendering <Dot cx cy {...props}/>; Field values
// read the data entry.
func DotShape(props Props) Fn { return Fn{Name: "dot", Args: []any{map[string]any(props)}} }

// IconDot is a Line dot rendering icon centered on the point, size wide,
// with props: <Icon x={cx - size/2} y={cy - size/2} width={size} height={size} {...props}/>.
func IconDot(icon Icon, size float64, props Props) Fn {
	return Fn{Name: "iconDot", Args: []any{"", size, map[string]any(props)}, icon: icon}
}

// --- recharts elements ----------------------------------------------------------

// AreaChart, BarChart, LineChart and ComposedChart are the chart roots.
func AreaChart(args ...any) *Node     { return El("AreaChart", args...) }
func BarChart(args ...any) *Node      { return El("BarChart", args...) }
func LineChart(args ...any) *Node     { return El("LineChart", args...) }
func ComposedChart(args ...any) *Node { return El("ComposedChart", args...) }

// Recharts components.
func CartesianGrid(args ...any) *Node { return El("CartesianGrid", args...) }
func XAxis(args ...any) *Node         { return El("XAxis", args...) }
func YAxis(args ...any) *Node         { return El("YAxis", args...) }
func Area(args ...any) *Node          { return El("Area", args...) }
func Bar(args ...any) *Node           { return El("Bar", args...) }
func Line(args ...any) *Node          { return El("Line", args...) }
func LabelList(args ...any) *Node     { return El("LabelList", args...) }
func Cell(args ...any) *Node          { return El("Cell", args...) }

// Tooltip is ChartTooltip (recharts' Tooltip); give it a TooltipContent.
func Tooltip(args ...any) *Node { return El("Tooltip", args...) }

// Legend is ChartLegend (recharts' Legend); give it a LegendContent.
func Legend(args ...any) *Node { return El("Legend", args...) }

// Defs is an SVG <defs> rendered inside the chart's <svg> (gradients);
// build its content with El("linearGradient", P("id", ...), El("stop", ...)).
func Defs(args ...any) *Node { return El("defs", args...) }

// Data sets the chart data.
func Data(rows []Row) Prop { return Prop{"data", rows} }

// Common props.
func DataKey(v string) Prop             { return Prop{"dataKey", v} }
func NameKey(v string) Prop             { return Prop{"nameKey", v} }
func LabelKey(v string) Prop            { return Prop{"labelKey", v} }
func Type(v string) Prop                { return Prop{"type", v} }
func Layout(v string) Prop              { return Prop{"layout", v} }
func StackID(v string) Prop             { return Prop{"stackId", v} }
func StackOffset(v string) Prop         { return Prop{"stackOffset", v} }
func Fill(v string) Prop                { return Prop{"fill", v} }
func FillOpacity(v float64) Prop        { return Prop{"fillOpacity", v} }
func Stroke(v string) Prop              { return Prop{"stroke", v} }
func StrokeWidth(v float64) Prop        { return Prop{"strokeWidth", v} }
func Radius(v float64) Prop             { return Prop{"radius", v} }
func Radii(tl, tr, br, bl float64) Prop { return Prop{"radius", []float64{tl, tr, br, bl}} }
func Hide() Prop                        { return Prop{"hide", true} }
func Dot(v any) Prop                    { return Prop{"dot", v} }
func ActiveDot(v any) Prop              { return Prop{"activeDot", v} }
func Shape(v any) Prop                  { return Prop{"shape", v} }
func TickLine(v any) Prop               { return Prop{"tickLine", v} }
func AxisLine(v any) Prop               { return Prop{"axisLine", v} }
func TickMargin(v float64) Prop         { return Prop{"tickMargin", v} }
func TickCount(v int) Prop              { return Prop{"tickCount", v} }
func MinTickGap(v float64) Prop         { return Prop{"minTickGap", v} }
func TickFormatter(f Fn) Prop           { return Prop{"tickFormatter", f} }
func Formatter(f Fn) Prop               { return Prop{"formatter", f} }
func Vertical(v bool) Prop              { return Prop{"vertical", v} }
func Horizontal(v bool) Prop            { return Prop{"horizontal", v} }
func Cursor(v any) Prop                 { return Prop{"cursor", v} }
func Position(v string) Prop            { return Prop{"position", v} }
func Offset(v float64) Prop             { return Prop{"offset", v} }
func ClassName(v string) Prop           { return Prop{"className", v} }
func FontSize(v float64) Prop           { return Prop{"fontSize", v} }
func AccessibilityLayer() Prop          { return Prop{"accessibilityLayer", true} }
func IsAnimationActive(v bool) Prop     { return Prop{"isAnimationActive", v} }

// --- shadcn tooltip and legend content -------------------------------------------------

// Indicator is ChartTooltipContent's indicator: "dot" (default), "line" or "dashed".
func Indicator(v string) Prop { return Prop{"indicator", v} }

// HideLabel, HideIndicator and HideIcon are ChartTooltipContent and
// ChartLegendContent's boolean props.
func HideLabel() Prop     { return Prop{"hideLabel", true} }
func HideIndicator() Prop { return Prop{"hideIndicator", true} }
func HideIcon() Prop      { return Prop{"hideIcon", true} }

// LabelFormatter is ChartTooltipContent's labelFormatter.
func LabelFormatter(f Fn) Prop { return Prop{"labelFormatter", f} }

// LabelClassName is ChartTooltipContent's labelClassName.
func LabelClassName(v string) Prop { return Prop{"labelClassName", v} }

// Color is ChartTooltipContent's color (overrides the indicator color).
func Color(v string) Prop { return Prop{"color", v} }

// VerticalAlign is ChartLegendContent's (and Legend's) verticalAlign.
func VerticalAlign(v string) Prop { return Prop{"verticalAlign", v} }

// TooltipContent is ChartTooltipContent. Its class names are resolved here
// (cn of the same inputs as upstream) for each layout the runtime renders.
func TooltipContent(args ...any) *Node {
	n := El("ChartTooltipContent", args...)
	className, _ := n.Props["className"].(string)
	labelClassName, _ := n.Props["labelClassName"].(string)
	indicator, _ := n.Props["indicator"].(string)
	if indicator == "" {
		indicator = "dot"
	}
	n.Props["classes"] = tooltipClasses(className, labelClassName, indicator)
	return n
}

func when(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

func tooltipClasses(className, labelClassName, indicator string) map[string]any {
	layout := func(nest bool) map[string]string {
		return map[string]string{
			"indicator": lib.Cn("shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg)",
				when(indicator == "dot", "h-2.5 w-2.5"),
				when(indicator == "line", "w-1"),
				when(indicator == "dashed", "w-0 border-[1.5px] border-dashed bg-transparent"),
				when(nest && indicator == "dashed", "my-0.5")),
			"value": lib.Cn("flex flex-1 justify-between leading-none", when(nest, "items-end"), when(!nest, "items-center")),
		}
	}
	return map[string]any{
		"root":  lib.Cn("grid min-w-[8rem] items-start gap-1.5 rounded-lg border border-border/50 bg-background px-2.5 py-1.5 text-xs shadow-xl", className),
		"label": lib.Cn("font-medium", labelClassName),
		"item":  lib.Cn("flex w-full flex-wrap items-stretch gap-2 [&>svg]:h-2.5 [&>svg]:w-2.5 [&>svg]:text-muted-foreground", when(indicator == "dot", "items-center")),
		"flat":  layout(false),
		"nest":  layout(true),
	}
}

// LegendContent is ChartLegendContent.
func LegendContent(args ...any) *Node {
	n := El("ChartLegendContent", args...)
	className, _ := n.Props["className"].(string)
	n.Props["classes"] = map[string]string{
		"rootTop":    lib.Cn("flex items-center justify-center gap-4", "pb-3", className),
		"rootBottom": lib.Cn("flex items-center justify-center gap-4", "pt-3", className),
		"item":       lib.Cn("flex items-center gap-1.5 [&>svg]:h-3 [&>svg]:w-3 [&>svg]:text-muted-foreground"),
	}
	return n
}

// --- container -------------------------------------------------------------------

const containerClass = "flex aspect-video justify-center text-xs [&_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&_.recharts-dot[stroke='#fff']]:stroke-transparent [&_.recharts-layer]:outline-hidden [&_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&_.recharts-radial-bar-background-sector]:fill-muted [&_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&_.recharts-sector]:outline-hidden [&_.recharts-sector[stroke='#fff']]:stroke-transparent [&_.recharts-surface]:outline-hidden"

// ID sets the chart id (data-chart="chart-<id>"), like ChartContainer's id.
type ID string

// Nonce is the CSP nonce for the container's <style>.
type Nonce string

// InitialDimension is the size the chart renders at before the container
// is measured (ChartContainer's initialDimension; default 320x200).
type InitialDimension struct{ Width, Height float64 }

// ContainerBuilder is returned by Container.
type ContainerBuilder struct {
	h.Builder[ContainerBuilder]
	config  Config
	chart   *Node
	id      string
	nonce   string
	initial InitialDimension
}

// Container is ChartContainer: the config, then attributes (h.Class, ...)
// and the chart root element.
func Container(config Config, children ...any) *ContainerBuilder {
	b := &ContainerBuilder{config: config, initial: InitialDimension{320, 200}}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the chart element and the options.
func (b *ContainerBuilder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case *Node:
		b.chart = v
	case ID:
		b.id = string(v)
	case Nonce:
		b.nonce = string(v)
	case InitialDimension:
		b.initial = v
	default:
		return false
	}
	return true
}

type configJSON struct {
	Key   string `json:"key"`
	Label string `json:"label,omitempty"`
	Color string `json:"color,omitempty"`
	Theme *Theme `json:"theme,omitempty"`
	Icon  bool   `json:"icon,omitempty"`
}

// Element renders the container.
func (b *ContainerBuilder) Element() *h.Element {
	generated := b.id == ""
	id := b.id
	if generated {
		// React's useId for the first component that asks: _r_0_. The
		// runtime renumbers duplicates when a page has several charts.
		id = "_r_0_"
	}
	chartID := "chart-" + strings.ReplaceAll(id, ":", "")
	css, err := b.config.css(chartID)
	if err != nil {
		return h.Errored(err)
	}
	// Icons: config icons and icon dots, rendered once into templates.
	var templates []any
	cfg := make([]configJSON, 0, len(b.config))
	for _, it := range b.config {
		cfg = append(cfg, configJSON{Key: it.Key, Label: it.Label, Color: it.Color, Theme: it.Theme, Icon: it.Icon != nil})
		if it.Icon != nil {
			templates = append(templates, iconTemplate(it.Key, it.Icon))
		}
	}
	var spec *Node
	if b.chart != nil {
		n := 0
		spec = resolveIcons(b.chart, func(icon Icon) string {
			key := fmt.Sprintf("fn-%d", n)
			n++
			templates = append(templates, iconTemplate(key, icon))
			return key
		})
	}
	payload, err := json.Marshal(map[string]any{
		"config":           cfg,
		"spec":             spec,
		"generatedId":      generated,
		"initialDimension": map[string]float64{"width": b.initial.Width, "height": b.initial.Height},
	})
	if err != nil {
		return h.Errored(fmt.Errorf("chart: %w", err))
	}
	props := b.Elem()
	el := h.Spread(h.Div(), props.Clone().SetChildren(), nil)
	el.SetClass(lib.Cn(containerClass, props.ClassName()))
	el.Data("slot", "chart")
	el.Data("chart", chartID)
	var style h.Node
	if css != "" {
		st := h.StyleEl(h.UnsafeCSS(css))
		if b.nonce != "" {
			st.Nonce(b.nonce)
		}
		style = st
	}
	el.SetChildren(
		style,
		h.Div(h.Class("recharts-responsive-container"), h.Style("width: 100%; height: 100%; min-width: 0px;"),
			h.Div(h.Style("width: 0px; height: 0px; overflow: visible;")),
		),
		// The spec is JSON produced by encoding/json, which escapes <, > and
		// &, so it cannot end the script element.
		h.Script(h.Type("application/json"), h.Data("chart-spec", ""), h.UnsafeJS(string(payload))),
		templates,
	)
	return el
}

// Render renders the container.
func (b *ContainerBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

func iconTemplate(key string, icon Icon) *h.Element {
	return h.Template(h.Data("chart-icon", key), h.Data("parity-ignore", ""), icon())
}

// resolveIcons returns a copy of n whose IconDot functions reference
// templates named by register.
func resolveIcons(n *Node, register func(Icon) string) *Node {
	out := &Node{Type: n.Type, Props: make(map[string]any, len(n.Props))}
	keys := make([]string, 0, len(n.Props))
	for k := range n.Props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		out.Props[k] = resolveValue(n.Props[k], register)
	}
	for _, c := range n.Children {
		out.Children = append(out.Children, resolveIcons(c, register))
	}
	return out
}

func resolveValue(v any, register func(Icon) string) any {
	switch x := v.(type) {
	case Fn:
		if x.icon != nil {
			args := append([]any(nil), x.Args...)
			args[0] = register(x.icon)
			return Fn{Name: x.Name, Args: args}
		}
		return x
	case Props:
		return resolveValue(map[string]any(x), register)
	case map[string]any:
		m := make(map[string]any, len(x))
		for k, e := range x {
			m[k] = resolveValue(e, register)
		}
		return m
	case []any:
		s := make([]any, len(x))
		for i, e := range x {
			s[i] = resolveValue(e, register)
		}
		return s
	}
	return v
}
