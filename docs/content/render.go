package content

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"shadcn.go/cn"
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/docs/hast"
	"shadcn.go/docs/mdx"
	h "shadcn.go/libhtml"
)

// Rendering reproduces what React renders for a compiled MDX page with the
// docs' components (ui/apps/v4/mdx-components.tsx):
//
//   - elements produced from Markdown go through the docs' element mapping
//     (headings with anchors, tables in a scroll wrapper, <pre> with the
//     code-block classes, figure captions with a language icon, code with
//     its copy button or package-manager command);
//   - JSX written in the MDX: lowercase tags are plain elements (MDX does not
//     map explicit JSX), capitalized ones are components, looked up in the
//     docs' wrappers and then in the Registry;
//   - expressions ({" "}, comments, the arrays the docs map over) are
//     evaluated.

// Prop is one prop of a component, in source order. Value is a string,
// float64, bool, nil, []any, or an h.Node (a JSX element passed as a prop).
type Prop struct {
	Name  string
	Value any
}

// Props are a component's props, in order (children excluded).
type Props []Prop

// Get returns a prop value.
func (ps Props) Get(name string) (any, bool) {
	for _, p := range ps {
		if p.Name == name {
			return p.Value, true
		}
	}
	return nil, false
}

// String returns a string prop ("" when absent or not a string).
func (ps Props) String(name string) string {
	v, _ := ps.Get(name)
	s, _ := v.(string)
	return s
}

// Without returns the props without the named ones.
func (ps Props) Without(names ...string) Props {
	var out Props
	for _, p := range ps {
		skip := false
		for _, n := range names {
			if p.Name == n {
				skip = true
			}
		}
		if !skip {
			out = append(out, p)
		}
	}
	return out
}

// Set returns the props with name set to v, keeping its position when it
// exists (like a JavaScript object spread), appended otherwise.
func (ps Props) Set(name string, v any) Props {
	out := append(Props(nil), ps...)
	for i := range out {
		if out[i].Name == name {
			out[i].Value = v
			return out
		}
	}
	return append(out, Prop{Name: name, Value: v})
}

// Component renders a JSX component. children are already rendered.
type Component func(ctx *Context, props Props, children []h.Node) (h.Node, error)

// Context is passed to components.
type Context struct {
	Page     *Page
	Registry *Registry
	// Node is the JSX element being rendered, when it comes from the MDX
	// (nil for components called by the element mapping).
	Node *hast.MdxJsx
}

// Components lists the components the docs content needs from a Registry:
// the ones ui/apps/v4/mdx-components.tsx imports from the app (its own
// wrappers, such as Steps, Tab or the className defaults of Tabs and Button,
// are applied here before the registry is called), with the props they
// receive:
//
//   - Callout, CodeTabs, ComponentPreview, ComponentSource,
//     CodeCollapsibleWrapper, DirectoryList, Kbd, Accordion*, Alert*,
//     AspectRatio, Button, Tabs, TabsList, TabsTrigger, TabsContent: the
//     JSX props (Button, Accordion and Tabs* with the wrapper's className);
//   - CopyButton {value, src}: before each highlighted code block's <code>;
//   - CodeBlockCommand {__npm__, __yarn__, __pnpm__, __bun__}: replaces the
//     <code> of package manager commands;
//   - ComponentsList {componentsFolder (PageTreeRef), currentBase, variant};
//   - Image: next/image with the wrapper's props (className, src, width,
//     height as numbers, alt, ...);
//   - Link: next/link (the default renders <a>, href last, like Next.js).
var Components = []string{
	"Accordion", "AccordionContent", "AccordionItem", "AccordionTrigger", "Alert", "AlertDescription",
	"AlertTitle", "AspectRatio", "Button", "Callout", "CodeBlockCommand", "CodeCollapsibleWrapper",
	"CodeTabs", "ComponentPreview", "ComponentSource", "ComponentsList", "CopyButton", "DirectoryList",
	"Image", "Kbd", "Link", "Tabs", "TabsContent", "TabsList", "TabsTrigger",
}

// Registry resolves the custom components the docs content uses. Names are
// the JSX names ("Callout", "ComponentPreview", "Tabs", ...); components
// imported in the MDX are registered as "<module>:<name>" (for example
// "lucide-react:InfoIcon"). Unregistered components render as placeholders
// (see Placeholder).
type Registry struct{ components map[string]Component }

// NewRegistry returns a registry with the defaults: Link renders <a> like
// next/link, the icons the docs import from lucide-react and
// @tabler/icons-react render their SVG; every other component renders a
// placeholder.
func NewRegistry() *Registry {
	r := &Registry{components: map[string]Component{}}
	r.Register("Link", nextLink)
	for name, fn := range defaultIcons {
		r.Register(name, iconComponent(fn))
	}
	return r
}

var defaultIcons = map[string]func(...any) *h.Element{
	"lucide-react:InfoIcon":              lucide.Info,
	"lucide-react:TriangleAlertIcon":     lucide.TriangleAlert,
	"lucide-react:ClipboardListIcon":     lucide.ClipboardList,
	"lucide-react:ExternalLinkIcon":      lucide.ExternalLink,
	"@tabler/icons-react:IconInfoCircle": tabler.IconInfoCircle,
}

// iconComponent renders an icon component: className adds classes, other
// props become attributes of the <svg>.
func iconComponent(fn func(...any) *h.Element) Component {
	return func(_ *Context, props Props, children []h.Node) (h.Node, error) {
		var args []any
		for _, p := range props {
			switch p.Name {
			case "className":
				if s, ok := p.Value.(string); ok {
					args = append(args, h.Class(s))
				}
			default:
				return nil, fmt.Errorf("content: unsupported icon prop %q", p.Name)
			}
		}
		for _, c := range children {
			args = append(args, c)
		}
		return fn(args...), nil
	}
}

// Register sets the component for a name.
func (r *Registry) Register(name string, c Component) { r.components[name] = c }

// Lookup returns the registered component.
func (r *Registry) Lookup(name string) (Component, bool) {
	c, ok := r.components[name]
	return c, ok
}

// Names returns the registered names, sorted.
func (r *Registry) Names() []string {
	var out []string
	for k := range r.components {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (r *Registry) render(ctx *Context, name string, props Props, children []h.Node) (h.Node, error) {
	if c, ok := r.components[name]; ok {
		return c(ctx, props, children)
	}
	return Placeholder(name, props, children), nil
}

// Placeholder renders a component that has no implementation yet:
//
//	<mdx-component data-name="Callout" data-props='[["title","Note"]]'>
//	  <template data-prop="icon">...</template>
//	  ...children...
//	</mdx-component>
//
// data-props holds the props that are plain values, as JSON pairs in prop
// order; props that are elements are rendered in <template data-prop>.
// Undefined props are left out.
func Placeholder(name string, props Props, children []h.Node) h.Node {
	var b strings.Builder
	b.WriteByte('[')
	var templates []any
	n := 0
	for _, p := range props {
		if node, ok := p.Value.(h.Node); ok {
			templates = append(templates, h.El("template", h.Attr("data-prop", p.Name), node))
			continue
		}
		switch p.Value.(type) {
		case mdx.Undefined, PageTreeRef:
			continue
		}
		if n > 0 {
			b.WriteByte(',')
		}
		n++
		b.WriteByte('[')
		hast.JSONString(&b, p.Name)
		b.WriteByte(',')
		writeJSON(&b, p.Value)
		b.WriteByte(']')
	}
	b.WriteByte(']')
	el := h.El("mdx-component", h.Attr("data-name", name), h.Attr("data-props", b.String()))
	el.Children(templates...)
	for _, c := range children {
		el.Children(c)
	}
	return el
}

func writeJSON(b *strings.Builder, v any) {
	switch x := v.(type) {
	case nil:
		b.WriteString("null")
	case string:
		hast.JSONString(b, x)
	case bool:
		b.WriteString(strconv.FormatBool(x))
	case float64:
		b.WriteString(jsNumber(x))
	case []any:
		b.WriteByte('[')
		for i, e := range x {
			if i > 0 {
				b.WriteByte(',')
			}
			writeJSON(b, e)
		}
		b.WriteByte(']')
	default:
		b.WriteString("null")
	}
}

func jsNumber(f float64) string {
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// nextLink renders next/link's <Link>: an <a> with the props, href last.
func nextLink(_ *Context, props Props, children []h.Node) (h.Node, error) {
	href, _ := props.Get("href")
	rest := props.Without("href", "prefetch", "replace", "scroll", "shallow", "passHref", "legacyBehavior", "locale")
	if href != nil {
		rest = append(rest, Prop{Name: "href", Value: href})
	}
	return intrinsic("a", rest, children)
}

// Render renders a page's content (what <MDX components={mdxComponents} />
// renders) with the registry's components (nil uses NewRegistry()).
func Render(p *Page, reg *Registry) (h.Node, error) {
	if reg == nil {
		reg = NewRegistry()
	}
	r := &renderer{ctx: &Context{Page: p, Registry: reg}, imports: map[string]string{}}
	for _, im := range p.Imports {
		r.imports[im.Local] = im.Source + ":" + im.Imported
	}
	// The MDX root is a fragment: whitespace-only text at its start and end
	// is dropped (hast-util-to-estree's root handler).
	nodes := p.Tree
	var kept []hast.Node
	for _, n := range nodes {
		if _, ok := n.(*hast.MdxEsm); ok {
			continue
		}
		kept = append(kept, n)
	}
	start, end := 0, len(kept)
	for start < end && isWhitespaceText(kept[start]) {
		start++
	}
	for end > start && isWhitespaceText(kept[end-1]) {
		end--
	}
	out, err := r.nodes(kept[start:end])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p.Path, err)
	}
	return h.Group(toAny(out)...), nil
}

// RenderTOCTitle renders a TOC entry's title as fumadocs exports it: the
// heading's content as plain elements (no element mapping).
func RenderTOCTitle(p *Page, item TOCItem, reg *Registry) (h.Node, error) {
	if reg == nil {
		reg = NewRegistry()
	}
	r := &renderer{ctx: &Context{Page: p, Registry: reg}, imports: map[string]string{}, plain: true}
	for _, im := range p.Imports {
		r.imports[im.Local] = im.Source + ":" + im.Imported
	}
	out, err := r.nodes(item.Title)
	if err != nil {
		return nil, err
	}
	return h.Group(toAny(out)...), nil
}

func toAny(ns []h.Node) []any {
	out := make([]any, len(ns))
	for i, n := range ns {
		out[i] = n
	}
	return out
}

func isWhitespaceText(n hast.Node) bool {
	t, ok := n.(*hast.Text)
	return ok && strings.Trim(t.Value, " \t\n\f\r") == ""
}

type renderer struct {
	ctx     *Context
	imports map[string]string
	// plain renders hast elements without the docs' element mapping.
	plain bool
}

func (r *renderer) nodes(ns []hast.Node) ([]h.Node, error) {
	var out []h.Node
	for _, n := range ns {
		res, err := r.node(n)
		if err != nil {
			return nil, err
		}
		out = append(out, res...)
	}
	return out, nil
}

func (r *renderer) node(n hast.Node) ([]h.Node, error) {
	switch x := n.(type) {
	case *hast.Text:
		if x.Value == "" {
			return nil, nil
		}
		return []h.Node{h.Text(x.Value)}, nil
	case *hast.Element:
		children := x.Children
		if tableElements[x.Tag] {
			// hast-util-to-estree drops the line breaks between table parts
			// (react-dom warns about any whitespace in tables).
			children = nil
			for _, c := range x.Children {
				if t, ok := c.(*hast.Text); ok && t.Value == "\n" {
					continue
				}
				children = append(children, c)
			}
		}
		kids, err := r.nodes(children)
		if err != nil {
			return nil, err
		}
		props := hastProps(x)
		if r.plain {
			n, err := intrinsic(x.Tag, props, kids)
			return one(n, err)
		}
		n, err := r.mapped(x, props, kids)
		return one(n, err)
	case *hast.MdxJsx:
		return r.jsx(x)
	case *hast.MdxExpression:
		v, err := r.eval(x.Expr.(mdx.Expr), nil)
		if err != nil {
			return nil, err
		}
		return valueNodes(v), nil
	case *hast.MdxEsm:
		return nil, nil
	}
	return nil, fmt.Errorf("content: cannot render %T", n)
}

var tableElements = map[string]bool{"table": true, "tbody": true, "thead": true, "tfoot": true, "tr": true}

func one(n h.Node, err error) ([]h.Node, error) {
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	return []h.Node{n}, nil
}

// hastProps converts hast attributes to props (hast-util-to-estree):
// boolean attributes become true.
func hastProps(e *hast.Element) Props {
	var ps Props
	for _, p := range e.Props {
		if p.Unset {
			continue
		}
		if p.Bool {
			ps = append(ps, Prop{Name: p.Name, Value: true})
			continue
		}
		name := p.Name
		if name == "class" {
			name = "className"
		}
		ps = append(ps, Prop{Name: name, Value: p.Value})
	}
	if e.Tag == "td" || e.Tag == "th" {
		// MDX compiles with tableCellAlignToStyle: align becomes text-align.
		if a := ps.String("align"); a != "" {
			ps = ps.Without("align").Set("style", "text-align:"+a)
		}
	}
	return ps
}

// --- the docs' element mapping (mdx-components.tsx) ---------------------------

func (r *renderer) mapped(e *hast.Element, props Props, kids []h.Node) (h.Node, error) {
	switch e.Tag {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		// const headingId = id ?? getHeadingId(children): an empty id
		// (an empty heading's slug) is kept as id="" without an anchor.
		idv, hasID := props.Get("id")
		id, _ := idv.(string)
		if !hasID || idv == nil {
			id = headingID(hast.ToString(e))
		}
		rest := props.Without("id")
		var inner []h.Node
		if id == "" {
			inner = kids
		} else {
			inner = []h.Node{headingAnchor(id, kids)}
		}
		var ps Props
		if hasID && idv != nil || id != "" {
			ps = append(ps, Prop{Name: "id", Value: id})
		}
		return intrinsic(e.Tag, append(ps, rest...), inner)
	case "table":
		t, err := intrinsic("table", props, kids)
		if err != nil {
			return nil, err
		}
		return intrinsic("div", Props{{"className", "typeset-scroll scroll-fade-x scrollbar-none *:[table]:w-full"}}, []h.Node{t})
	case "pre":
		cls := cn.Cn("no-scrollbar min-w-0 overflow-x-auto overflow-y-auto overscroll-x-contain overscroll-y-auto px-4 py-3.5 outline-none has-data-highlighted-line:px-0 has-data-line-numbers:px-0 has-data-[slot=tabs]:p-0", props.String("className"))
		ps := Props{{"data-not-typeset", true}, {"className", cls}}
		ps = append(ps, props.Without("className")...)
		return intrinsic("pre", ps, kids)
	case "figcaption":
		cls := cn.Cn("flex items-center gap-2 text-code-foreground [&_svg]:size-4 [&_svg]:text-code-foreground [&_svg]:opacity-70", props.String("className"))
		ps := Props{{"className", cls}}
		ps = append(ps, props.Without("className")...)
		var children []h.Node
		if lang, ok := props.Get("data-language"); ok {
			if s, ok := lang.(string); ok {
				children = append(children, languageIcon(s))
			}
		}
		children = append(children, kids...)
		return intrinsic("figcaption", ps, children)
	case "code":
		return r.code(e, props, kids)
	}
	return intrinsic(e.Tag, props, kids)
}

func (r *renderer) code(e *hast.Element, props Props, kids []h.Node) (h.Node, error) {
	// Inline code: a single text child (React's props.children is a string).
	if len(e.Children) == 1 {
		if _, ok := e.Children[0].(*hast.Text); ok {
			return intrinsic("code", props, kids)
		}
	}
	str := func(n string) string { return props.String(n) }
	rest := props.Without("className", "__raw__", "__src__", "__npm__", "__yarn__", "__pnpm__", "__bun__")
	if str("__npm__") != "" && str("__yarn__") != "" && str("__pnpm__") != "" && str("__bun__") != "" {
		return r.ctx.Registry.render(r.ctx, "CodeBlockCommand", Props{
			{"__npm__", str("__npm__")}, {"__yarn__", str("__yarn__")}, {"__pnpm__", str("__pnpm__")}, {"__bun__", str("__bun__")},
		}, nil)
	}
	var out []any
	if raw := str("__raw__"); raw != "" {
		cp := Props{{"value", raw}}
		if src, ok := props.Get("__src__"); ok {
			cp = append(cp, Prop{"src", src})
		}
		b, err := r.ctx.Registry.render(r.ctx, "CopyButton", cp, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	c, err := intrinsic("code", rest, kids)
	if err != nil {
		return nil, err
	}
	out = append(out, c)
	return h.Group(out...), nil
}

func headingAnchor(id string, kids []h.Node) h.Node {
	return h.A(h.Class("group no-underline"), h.Href("#"+id),
		h.Span(h.Class("underline-offset-4 group-hover:underline"), toAny(kids)),
		h.Span(h.Aria("hidden", "true"), h.Class("ml-2 text-muted-foreground opacity-0 group-hover:opacity-100"), "#"),
	)
}

// headingID is mdx-components' getHeadingId, for headings without an id.
func headingID(text string) string {
	s := strings.TrimFunc(text, isJSSpace)
	s = strings.Join(strings.FieldsFunc(s, isJSSpace), "-")
	s = strings.NewReplacer("'", "", "?", "").Replace(s)
	return strings.ToLower(s)
}

// --- JSX ------------------------------------------------------------------------

func isComponentName(name string) bool {
	if name == "" {
		return false
	}
	r := []rune(name)[0]
	return unicode.IsUpper(r) || strings.ContainsAny(name, ".")
}

func (r *renderer) jsx(x *hast.MdxJsx) ([]h.Node, error) {
	kids, err := r.nodes(x.Children)
	if err != nil {
		return nil, err
	}
	props, err := r.jsxProps(x.Attrs)
	if err != nil {
		return nil, err
	}
	if x.Name == "" { // fragment
		return kids, nil
	}
	if !isComponentName(x.Name) {
		n, err := intrinsic(x.Name, props, kids)
		return one(n, err)
	}
	ctx := *r.ctx
	ctx.Node = x
	n, err := r.component(&ctx, x.Name, props, kids)
	return one(n, err)
}

func (r *renderer) jsxProps(attrs []hast.JsxAttr) (Props, error) {
	var ps Props
	for _, a := range attrs {
		switch a.Kind {
		case hast.AttrString:
			ps = append(ps, Prop{a.Name, a.Value})
		case hast.AttrBool:
			ps = append(ps, Prop{a.Name, true})
		case hast.AttrExpression:
			v, err := r.eval(a.Expr.(mdx.Expr), nil)
			if err != nil {
				return nil, err
			}
			ps = append(ps, Prop{a.Name, v})
		}
	}
	return ps, nil
}

// component resolves a capitalized JSX name: the page's imports, the docs'
// wrappers in mdx-components.tsx, then the registry.
func (r *renderer) component(ctx *Context, name string, props Props, kids []h.Node) (h.Node, error) {
	if src, ok := r.imports[name]; ok {
		return ctx.Registry.render(ctx, src, props, kids)
	}
	withClass := func(target, base string) (h.Node, error) {
		ps := Props{{"className", cn.Cn(base, props.String("className"))}}
		ps = append(ps, props.Without("className")...)
		return ctx.Registry.render(ctx, target, ps, kids)
	}
	switch name {
	case "Step":
		return intrinsic("h3", props, kids)
	case "Steps":
		ps := Props{{"className", cn.Cn("steps mb-12 [counter-reset:step] md:ml-4 md:border-l md:pl-8 [&>h3]:step", props.String("className"))}}
		return intrinsic("div", append(ps, props.Without("className")...), kids)
	case "Tab":
		return intrinsic("div", props, kids)
	case "Image":
		ps := Props{
			{"className", cn.Cn("mt-6 rounded-2xl border", props.String("className"))},
			{"src", orEmpty(props, "src")},
			{"width", jsNumberValue(props, "width")},
			{"height", jsNumberValue(props, "height")},
			{"alt", orEmpty(props, "alt")},
		}
		ps = append(ps, props.Without("src", "className", "width", "height", "alt")...)
		return ctx.Registry.render(ctx, "Image", ps, kids)
	case "Tabs":
		return withClass("Tabs", "relative mt-6 w-full")
	case "TabsList":
		return withClass("TabsList", "justify-start gap-4 rounded-none bg-transparent px-0")
	case "TabsTrigger":
		return withClass("TabsTrigger", "not-typset rounded-none border-0 border-b-2 border-transparent bg-transparent px-0 pb-3 text-base text-muted-foreground hover:text-primary data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none! dark:data-[state=active]:border-primary dark:data-[state=active]:bg-transparent")
	case "TabsContent":
		return withClass("TabsContent", "relative [&_h3.font-heading]:text-base [&_h3.font-heading]:font-medium *:[figure]:first:mt-0 [&>.steps]:mt-6")
	case "Button":
		return withClass("Button", "not-typeset")
	case "Accordion":
		return withClass("Accordion", "not-typeset")
	case "LinkedCard":
		ps := Props{{"data-not-typeset", true}, {"className", cn.Cn("flex w-full flex-col items-center rounded-2xl bg-surface p-6 text-surface-foreground transition-colors hover:bg-surface/80 sm:p-10", props.String("className"))}}
		ps = append(ps, props.Without("className")...)
		return ctx.Registry.render(ctx, "Link", ps, kids)
	case "ComponentsList":
		// ComponentsListWrapper: <ComponentsList componentsFolder={...}
		// currentBase="base" variant={variant} />. The components folder of
		// the page tree is the docs chrome's data, passed as a PageTreeRef.
		v, ok := props.Get("variant")
		if !ok {
			v = mdx.Undefined{}
		}
		return ctx.Registry.render(ctx, "ComponentsList", Props{{"componentsFolder", PageTreeRef{ID: "components"}}, {"currentBase", "base"}, {"variant", v}}, nil)
	case "Callout", "AccordionContent", "AccordionItem", "AccordionTrigger", "Alert", "AlertTitle",
		"AlertDescription", "AspectRatio", "CodeTabs", "ComponentPreview", "ComponentSource",
		"CodeCollapsibleWrapper", "DirectoryList", "Link", "Kbd":
		return ctx.Registry.render(ctx, name, props, kids)
	}
	return nil, fmt.Errorf("content: unknown component <%s> (not in mdx-components.tsx nor imported)", name)
}

func orEmpty(ps Props, name string) any {
	v, ok := ps.Get(name)
	if !ok || v == nil || v == "" || v == false {
		return ""
	}
	if _, ok := v.(mdx.Undefined); ok {
		return ""
	}
	return v
}

// jsNumberValue is Number(props[name]).
func jsNumberValue(ps Props, name string) any {
	v, ok := ps.Get(name)
	if !ok {
		return nanValue{}
	}
	switch x := v.(type) {
	case float64:
		return x
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			if strings.TrimSpace(x) == "" {
				return 0.0
			}
			return nanValue{}
		}
		return f
	case bool:
		if x {
			return 1.0
		}
		return 0.0
	case nil:
		return 0.0
	}
	return nanValue{}
}

// PageTreeRef refers to a folder of the docs page tree (by its $id), which
// the docs chrome owns.
type PageTreeRef struct{ ID string }

// nanValue is NaN (JSON.stringify writes null).
type nanValue struct{}

// --- intrinsic elements -----------------------------------------------------------

// reactAttr maps React prop names to HTML/SVG attribute names.
var reactAttr = map[string]string{
	"className": "class", "htmlFor": "for", "tabIndex": "tabindex", "readOnly": "readonly",
	"maxLength": "maxlength", "colSpan": "colspan", "rowSpan": "rowspan", "crossOrigin": "crossorigin",
	"strokeWidth": "stroke-width", "strokeLinecap": "stroke-linecap", "strokeLinejoin": "stroke-linejoin",
	"strokeMiterlimit": "stroke-miterlimit", "strokeOpacity": "stroke-opacity", "strokeDasharray": "stroke-dasharray",
	"strokeDashoffset": "stroke-dashoffset", "fillRule": "fill-rule", "clipRule": "clip-rule", "fillOpacity": "fill-opacity",
	"clipPath": "clip-path", "stopColor": "stop-color", "stopOpacity": "stop-opacity", "fontSize": "font-size",
	"fontFamily": "font-family", "fontWeight": "font-weight", "textAnchor": "text-anchor", "xmlnsXlink": "xmlns:xlink",
	"xlinkHref": "xlink:href", "dominantBaseline": "dominant-baseline", "colorInterpolationFilters": "color-interpolation-filters",
}

// keptCamel are React props rendered with their camelCase name.
var keptCamel = map[string]bool{
	"viewBox": true, "preserveAspectRatio": true, "gradientUnits": true, "gradientTransform": true,
	"patternUnits": true, "patternContentUnits": true, "maskUnits": true, "clipPathUnits": true,
	"stdDeviation": true, "baseFrequency": true, "numOctaves": true, "filterUnits": true, "markerWidth": true,
	"markerHeight": true, "refX": true, "refY": true, "pathLength": true, "spreadMethod": true,
}

var htmlBoolean = map[string]bool{
	"allowfullscreen": true, "async": true, "autofocus": true, "autoplay": true, "checked": true,
	"controls": true, "default": true, "defer": true, "disabled": true, "formnovalidate": true,
	"hidden": true, "inert": true, "loop": true, "multiple": true, "muted": true, "nomodule": true,
	"novalidate": true, "open": true, "playsinline": true, "readonly": true, "required": true,
	"reversed": true, "scoped": true, "seamless": true, "selected": true, "itemscope": true,
}

// reactHoisted are elements React DOM's server renderer treats specially
// (document structure and resources it hoists into <head>). <title> is
// left out: the docs use it inside SVG icons, where React keeps it.
var reactHoisted = map[string]bool{"html": true, "head": true, "body": true, "meta": true, "link": true, "style": true, "script": true, "base": true}

// intrinsic renders a DOM element from React props as React does.
func intrinsic(tag string, props Props, kids []h.Node) (h.Node, error) {
	if reactHoisted[tag] {
		return nil, fmt.Errorf("content: <%s> is not supported (React DOM gives it document semantics: it adds <head>, hoists or moves it)", tag)
	}
	el := h.El(tag)
	for _, p := range props {
		name := p.Name
		if name == "key" || name == "ref" || name == "children" {
			continue
		}
		if a, ok := reactAttr[name]; ok {
			name = a
		} else if !keptCamel[name] && name != strings.ToLower(name) && !strings.HasPrefix(name, "data-") && !strings.HasPrefix(name, "aria-") {
			return nil, fmt.Errorf("content: <%s>: unsupported React prop %q", tag, name)
		}
		switch v := p.Value.(type) {
		case nil, mdx.Undefined, nanValue:
			continue
		case string:
			if name == "style" {
				v = reactStyle(v)
			}
			if err := checkAttr(name, v); err != nil {
				return nil, fmt.Errorf("content: <%s>: %w", tag, err)
			}
			el.Attr(name, v)
		case float64:
			el.Attr(name, jsNumber(v))
		case bool:
			switch {
			case strings.HasPrefix(name, "data-") || strings.HasPrefix(name, "aria-"):
				el.Attr(name, strconv.FormatBool(v))
			case htmlBoolean[strings.ToLower(name)]:
				if v {
					el.Bool(name, true)
				}
			default:
				// React drops booleans for non-boolean attributes.
			}
		default:
			return nil, fmt.Errorf("content: <%s>: unsupported value for prop %q", tag, p.Name)
		}
	}
	for _, k := range kids {
		el.Children(k)
	}
	return el, nil
}

// reactStyle renders a CSS declaration list the way React serializes the
// style object MDX builds from it: "name:value" pairs joined by ";".
func reactStyle(s string) string {
	var parts []string
	for _, decl := range strings.Split(s, ";") {
		name, value, ok := strings.Cut(decl, ":")
		if !ok {
			continue
		}
		name, value = strings.TrimSpace(name), strings.TrimSpace(value)
		if name == "" || value == "" {
			continue
		}
		if !strings.HasPrefix(name, "--") {
			name = strings.ToLower(name)
		}
		parts = append(parts, name+":"+value)
	}
	return strings.Join(parts, ";")
}

// --- expressions -------------------------------------------------------------------

type scope struct {
	vars   map[string]any
	parent *scope
}

func (s *scope) lookup(name string) (any, bool) {
	for c := s; c != nil; c = c.parent {
		if v, ok := c.vars[name]; ok {
			return v, true
		}
	}
	return nil, false
}

type closure struct {
	fn    *mdx.Arrow
	scope *scope
}

func (r *renderer) eval(e mdx.Expr, sc *scope) (any, error) {
	switch x := e.(type) {
	case mdx.Empty:
		return mdx.Undefined{}, nil
	case mdx.Lit:
		return x.Value, nil
	case mdx.Ident:
		if v, ok := sc.lookup(x.Name); ok {
			return v, nil
		}
		return nil, fmt.Errorf("content: undefined identifier %q in expression", x.Name)
	case *mdx.Array:
		var out []any
		for _, el := range x.Elems {
			v, err := r.eval(el, sc)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		if out == nil {
			out = []any{}
		}
		return out, nil
	case *mdx.Arrow:
		return &closure{fn: x, scope: sc}, nil
	case *mdx.Call:
		m, ok := x.Callee.(*mdx.Member)
		if !ok || m.Prop != "map" || len(x.Args) != 1 {
			return nil, fmt.Errorf("content: unsupported call in expression")
		}
		obj, err := r.eval(m.Obj, sc)
		if err != nil {
			return nil, err
		}
		arr, ok := obj.([]any)
		if !ok {
			return nil, fmt.Errorf("content: .map on a non-array")
		}
		fv, err := r.eval(x.Args[0], sc)
		if err != nil {
			return nil, err
		}
		fn, ok := fv.(*closure)
		if !ok {
			return nil, fmt.Errorf("content: .map needs a function")
		}
		var out []any
		for i, item := range arr {
			v, err := r.call(fn, []any{item, float64(i), arr})
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		if out == nil {
			out = []any{}
		}
		return out, nil
	case *mdx.JSXElement:
		return r.evalJSX(x, sc)
	}
	return nil, fmt.Errorf("content: unsupported expression %T", e)
}

func (r *renderer) call(fn *closure, args []any) (any, error) {
	sc := &scope{vars: map[string]any{}, parent: fn.scope}
	for i, p := range fn.fn.Params {
		var v any = mdx.Undefined{}
		if i < len(args) {
			v = args[i]
		}
		if err := bind(sc, p, v); err != nil {
			return nil, err
		}
	}
	return r.eval(fn.fn.Body, sc)
}

func bind(sc *scope, p mdx.Pattern, v any) error {
	if !p.Array {
		sc.vars[p.Name] = v
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return fmt.Errorf("content: cannot destructure a non-array")
	}
	for i, el := range p.Elems {
		var x any = mdx.Undefined{}
		if i < len(arr) {
			x = arr[i]
		}
		if err := bind(sc, el, x); err != nil {
			return err
		}
	}
	if p.Rest != "" {
		rest := []any{}
		if len(p.Elems) < len(arr) {
			rest = append(rest, arr[len(p.Elems):]...)
		}
		sc.vars[p.Rest] = rest
	}
	return nil
}

// evalJSX renders a JSX element written inside an expression.
func (r *renderer) evalJSX(x *mdx.JSXElement, sc *scope) (any, error) {
	var kids []h.Node
	for _, c := range x.Children {
		switch t := c.(type) {
		case mdx.JSXText:
			kids = append(kids, h.Text(t.Value))
		default:
			v, err := r.eval(c, sc)
			if err != nil {
				return nil, err
			}
			kids = append(kids, valueNodes(v)...)
		}
	}
	var props Props
	for _, a := range x.Attrs {
		var v any = true
		if a.Value != nil {
			var err error
			if v, err = r.eval(a.Value, sc); err != nil {
				return nil, err
			}
		}
		props = append(props, Prop{a.Name, v})
	}
	if x.Name == "" {
		return h.Group(toAny(kids)...), nil
	}
	if !isComponentName(x.Name) {
		return intrinsic(x.Name, props, kids)
	}
	return r.component(r.ctx, x.Name, props, kids)
}

// valueNodes renders a JavaScript value as React children.
func valueNodes(v any) []h.Node {
	switch x := v.(type) {
	case nil, mdx.Undefined, bool:
		return nil
	case string:
		if x == "" {
			return nil
		}
		return []h.Node{h.Text(x)}
	case float64:
		return []h.Node{h.Text(jsNumber(x))}
	case []any:
		var out []h.Node
		for _, e := range x {
			out = append(out, valueNodes(e)...)
		}
		return out
	case h.Node:
		return []h.Node{x}
	}
	return nil
}

// isJSSpace is JavaScript's \s (and String.prototype.trim's set), which
// differs from unicode.IsSpace: U+FEFF is whitespace, U+0085 is not.
func isJSSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xa0, 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff:
		return true
	}
	return r >= 0x2000 && r <= 0x200a
}

// checkAttr fails when libhtml would not write an attribute value as given:
// it filters URLs with schemes other than http(s), mailto and tel (and
// some it misreads as schemes, such as "foo):"), unsafe styles and
// script-like attributes, where React writes the value unchanged. Failing
// keeps such content from rendering differently from upstream; the docs
// content has none.
func checkAttr(name, v string) error {
	got, err := h.RenderString(h.El("x").Attr(name, v))
	if err != nil {
		return err
	}
	if want := "<x " + name + `="` + h.EscapeAttr(v) + `"></x>`; got != want {
		return fmt.Errorf("libhtml would rewrite %s=%q (as %s), React writes it unchanged", name, v, got)
	}
	return nil
}
