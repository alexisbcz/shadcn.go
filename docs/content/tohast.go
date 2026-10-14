package content

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"

	"shadcn.go/docs/hast"
	"shadcn.go/docs/mdx"
)

// toHast is mdast-util-to-hast 13.2.0 as remark-rehype runs it inside MDX:
// MDX nodes pass through (their children converted, without the newlines
// other parents get), and GFM footnotes produce the footer section.
func toHast(root *mdx.Node) ([]hast.Node, error) {
	st := &hastState{footnoteByID: map[string]*mdx.Node{}, footnoteCounts: map[string]int{}}
	var collect func(n *mdx.Node)
	collect = func(n *mdx.Node) {
		if n.Type == "footnoteDefinition" {
			id := strings.ToUpper(n.Identifier)
			if _, ok := st.footnoteByID[id]; !ok {
				st.footnoteByID[id] = n
			}
		}
		for _, c := range n.Children {
			collect(c)
		}
	}
	collect(root)
	kids, err := st.all(root)
	if err != nil {
		return nil, err
	}
	out := wrap(kids, false)
	if foot, err := st.footer(); err != nil {
		return nil, err
	} else if foot != nil {
		out = append(out, &hast.Text{Value: "\n"}, foot)
	}
	return out, nil
}

type hastState struct {
	footnoteByID   map[string]*mdx.Node
	footnoteCounts map[string]int
	footnoteOrder  []string
}

func text(s string) *hast.Text { return &hast.Text{Value: s} }

func wrap(nodes []hast.Node, loose bool) []hast.Node {
	var out []hast.Node
	if loose {
		out = append(out, text("\n"))
	}
	for i, n := range nodes {
		if i > 0 {
			out = append(out, text("\n"))
		}
		out = append(out, n)
	}
	if loose && len(nodes) > 0 {
		out = append(out, text("\n"))
	}
	return out
}

func (st *hastState) all(parent *mdx.Node) ([]hast.Node, error) {
	var values []hast.Node
	for i, c := range parent.Children {
		res, err := st.one(c, parent)
		if err != nil {
			return nil, err
		}
		if len(res) == 0 {
			continue
		}
		if i > 0 && parent.Children[i-1].Type == "break" && len(res) == 1 {
			switch r := res[0].(type) {
			case *hast.Text:
				r.Value = strings.TrimLeft(r.Value, "\t ")
			case *hast.Element:
				if len(r.Children) > 0 {
					if t, ok := r.Children[0].(*hast.Text); ok {
						t.Value = strings.TrimLeft(t.Value, "\t ")
					}
				}
			}
		}
		values = append(values, res...)
	}
	return values, nil
}

func el(tag string, children []hast.Node, attrs ...string) *hast.Element {
	e := hast.El(tag, attrs...)
	e.Children = children
	return e
}

func (st *hastState) one(n *mdx.Node, parent *mdx.Node) ([]hast.Node, error) {
	kids := func() ([]hast.Node, error) { return st.all(n) }
	single := func(e hast.Node, err error) ([]hast.Node, error) {
		if err != nil {
			return nil, err
		}
		return []hast.Node{e}, nil
	}
	switch n.Type {
	case "paragraph", "emphasis", "strong", "delete":
		tag := map[string]string{"paragraph": "p", "emphasis": "em", "strong": "strong", "delete": "del"}[n.Type]
		c, err := kids()
		return single(el(tag, c), err)
	case "heading":
		c, err := kids()
		if err != nil {
			return nil, err
		}
		h := el("h"+strconv.Itoa(n.Depth), c)
		if id, ok := n.Data["id"].(string); ok {
			h.Set("id", id)
		}
		return []hast.Node{h}, nil
	case "text":
		return []hast.Node{text(trimLines(n.Value))}, nil
	case "inlineCode":
		v := strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ").Replace(n.Value)
		return []hast.Node{el("code", []hast.Node{text(v)})}, nil
	case "break":
		return []hast.Node{el("br", nil), text("\n")}, nil
	case "thematicBreak":
		return []hast.Node{el("hr", nil)}, nil
	case "blockquote":
		c, err := kids()
		return single(el("blockquote", wrap(c, true)), err)
	case "code":
		value := ""
		if n.Value != "" {
			value = n.Value + "\n"
		}
		code := el("code", []hast.Node{text(value)})
		if n.HasLang && n.Lang != "" {
			code.Set("class", "language-"+n.Lang)
		}
		if n.HasMeta && n.Meta != "" {
			code.Data = map[string]any{"meta": n.Meta}
		}
		return []hast.Node{el("pre", []hast.Node{code})}, nil
	case "link":
		c, err := kids()
		if err != nil {
			return nil, err
		}
		a := el("a", c, "href", normalizeURI(n.URL))
		if n.HasTitle {
			a.Set("title", n.Title)
		}
		return []hast.Node{a}, nil
	case "image":
		img := el("img", nil, "src", normalizeURI(n.URL), "alt", n.Alt)
		if n.HasTitle {
			img.Set("title", n.Title)
		}
		for _, k := range []string{"src", "width", "height"} {
			if v, ok := n.Data[k].(string); ok {
				img.Set(k, v)
			}
		}
		return []hast.Node{img}, nil
	case "list":
		return st.list(n)
	case "listItem":
		return st.listItem(n, parent)
	case "table":
		return st.table(n)
	case "footnoteReference":
		return st.footnoteReference(n)
	case "footnoteDefinition":
		return nil, nil
	case "mdxJsxFlowElement", "mdxJsxTextElement":
		c, err := kids()
		if err != nil {
			return nil, err
		}
		if c == nil {
			c = []hast.Node{}
		}
		return []hast.Node{&hast.MdxJsx{Name: n.Name, Attrs: n.Attrs, Children: c, Flow: n.Type == "mdxJsxFlowElement"}}, nil
	case "mdxFlowExpression", "mdxTextExpression":
		return []hast.Node{&hast.MdxExpression{Value: n.Value, Flow: n.Type == "mdxFlowExpression", Expr: n.Expr}}, nil
	case "mdxjsEsm":
		return []hast.Node{&hast.MdxEsm{Value: n.Value, Imports: n.Imports}}, nil
	}
	return nil, fmt.Errorf("content: unsupported mdast node %q", n.Type)
}

func (st *hastState) list(n *mdx.Node) ([]hast.Node, error) {
	results, err := st.all(n)
	if err != nil {
		return nil, err
	}
	tag := "ul"
	if n.Ordered {
		tag = "ol"
	}
	l := el(tag, wrap(results, true))
	if n.Ordered && n.Start != 1 {
		l.Props = append([]hast.Prop{{Name: "start", Value: strconv.Itoa(n.Start)}}, l.Props...)
	}
	return []hast.Node{l}, nil
}

func listLoose(list *mdx.Node) bool {
	if list.Spread {
		return true
	}
	for _, it := range list.Children {
		if it.Spread {
			return true
		}
	}
	return false
}

func (st *hastState) listItem(n *mdx.Node, parent *mdx.Node) ([]hast.Node, error) {
	results, err := st.all(n)
	if err != nil {
		return nil, err
	}
	loose := n.Spread
	if parent != nil {
		loose = listLoose(parent)
	}
	var children []hast.Node
	isP := func(h hast.Node) bool { e, ok := h.(*hast.Element); return ok && e.Tag == "p" }
	for i, c := range results {
		if loose || i != 0 || !isP(c) {
			children = append(children, text("\n"))
		}
		if isP(c) && !loose {
			children = append(children, c.(*hast.Element).Children...)
		} else {
			children = append(children, c)
		}
	}
	if len(results) > 0 {
		tail := results[len(results)-1]
		if loose || !isP(tail) {
			children = append(children, text("\n"))
		}
	}
	return []hast.Node{el("li", children)}, nil
}

func (st *hastState) table(n *mdx.Node) ([]hast.Node, error) {
	var rows []hast.Node
	for i, r := range n.Children {
		tag := "td"
		if i == 0 {
			tag = "th"
		}
		var cells []hast.Node
		for ci := 0; ci < len(n.Align); ci++ {
			cell := el(tag, nil)
			if a := n.Align[ci]; a != "" {
				cell.Set("align", a)
			}
			if ci < len(r.Children) {
				kids, err := st.all(r.Children[ci])
				if err != nil {
					return nil, err
				}
				cell.Children = kids
			}
			cells = append(cells, cell)
		}
		rows = append(rows, el("tr", wrap(cells, true)))
	}
	var content []hast.Node
	if len(rows) > 0 {
		content = append(content, el("thead", wrap(rows[:1], true)))
	}
	if len(rows) > 1 {
		content = append(content, el("tbody", wrap(rows[1:], true)))
	}
	return []hast.Node{el("table", wrap(content, true))}, nil
}

const clobberPrefix = "user-content-"

func (st *hastState) footnoteReference(n *mdx.Node) ([]hast.Node, error) {
	id := strings.ToUpper(n.Identifier)
	safeID := normalizeURI(strings.ToLower(id))
	index := -1
	for i, x := range st.footnoteOrder {
		if x == id {
			index = i
		}
	}
	var counter int
	reuse, seen := st.footnoteCounts[id]
	if !seen {
		reuse = 0
		st.footnoteOrder = append(st.footnoteOrder, id)
		counter = len(st.footnoteOrder)
	} else {
		counter = index + 1
	}
	reuse++
	st.footnoteCounts[id] = reuse
	refID := clobberPrefix + "fnref-" + safeID
	if reuse > 1 {
		refID += "-" + strconv.Itoa(reuse)
	}
	a := el("a", []hast.Node{text(strconv.Itoa(counter))}, "href", "#"+clobberPrefix+"fn-"+safeID, "id", refID)
	a.Props = append(a.Props, hast.Prop{Name: "data-footnote-ref", Bool: true}, hast.Prop{Name: "aria-describedby", Value: "footnote-label"})
	return []hast.Node{el("sup", []hast.Node{a})}, nil
}

func (st *hastState) footer() (hast.Node, error) {
	var items []hast.Node
	for ri, id := range st.footnoteOrder {
		def := st.footnoteByID[id]
		if def == nil {
			continue
		}
		content, err := st.all(def)
		if err != nil {
			return nil, err
		}
		safeID := normalizeURI(strings.ToLower(strings.ToUpper(def.Identifier)))
		var back []hast.Node
		counts := st.footnoteCounts[id]
		for rr := 1; rr <= counts; rr++ {
			if len(back) > 0 {
				back = append(back, text(" "))
			}
			href := "#" + clobberPrefix + "fnref-" + safeID
			label := "Back to reference " + strconv.Itoa(ri+1)
			if rr > 1 {
				href += "-" + strconv.Itoa(rr)
				label += "-" + strconv.Itoa(rr)
			}
			kids := []hast.Node{text("↩")}
			if rr > 1 {
				kids = append(kids, el("sup", []hast.Node{text(strconv.Itoa(rr))}))
			}
			back = append(back, el("a", kids, "href", href, "data-footnote-backref", "", "aria-label", label, "class", "data-footnote-backref"))
		}
		if len(content) > 0 {
			if tail, ok := content[len(content)-1].(*hast.Element); ok && tail.Tag == "p" {
				if len(tail.Children) > 0 {
					if tt, ok := tail.Children[len(tail.Children)-1].(*hast.Text); ok {
						tt.Value += " "
					} else {
						tail.Children = append(tail.Children, text(" "))
					}
				} else {
					tail.Children = append(tail.Children, text(" "))
				}
				tail.Children = append(tail.Children, back...)
			} else {
				content = append(content, back...)
			}
		} else {
			content = append(content, back...)
		}
		items = append(items, el("li", wrap(content, true), "id", clobberPrefix+"fn-"+safeID))
	}
	if len(items) == 0 {
		return nil, nil
	}
	section := &hast.Element{Tag: "section", Props: []hast.Prop{{Name: "data-footnotes", Bool: true}, {Name: "class", Value: "footnotes"}}}
	section.Children = []hast.Node{
		el("h2", []hast.Node{text("Footnotes")}, "class", "sr-only", "id", "footnote-label"),
		text("\n"),
		el("ol", wrap(items, true)),
		text("\n"),
	}
	return section, nil
}

// trimLines is trim-lines 3.0.1: spaces and tabs around line endings are
// removed.
func trimLines(s string) string {
	if !strings.ContainsAny(s, "\r\n") {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if i > 0 {
			l = strings.TrimLeft(l, " \t")
		}
		if i < len(lines)-1 {
			l = strings.TrimRight(l, " \t")
		}
		lines[i] = l
	}
	return strings.Join(lines, "\n")
}

// normalizeURI is micromark-util-sanitize-uri's normalizeUri: characters
// that are not allowed in URLs are percent-encoded (encodeURIComponent),
// existing percent escapes are kept.
func normalizeURI(v string) string {
	units := utf16.Encode([]rune(v))
	var b strings.Builder
	isAlnum := func(i int) bool {
		if i >= len(units) {
			return false
		}
		c := units[i]
		return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
	}
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch {
		case c == '%' && isAlnum(i+1) && isAlnum(i+2):
			b.WriteRune(rune(c))
			b.WriteRune(rune(units[i+1]))
			b.WriteRune(rune(units[i+2]))
			i += 2
		case c < 128:
			if allowedURLChar(byte(c)) {
				b.WriteByte(byte(c))
			} else {
				b.WriteString(encodeURIComponent(string(rune(c))))
			}
		case c >= 0xd800 && c < 0xe000:
			if c < 0xdc00 && i+1 < len(units) && units[i+1] >= 0xdc00 && units[i+1] < 0xe000 {
				b.WriteString(encodeURIComponent(string(utf16.DecodeRune(rune(c), rune(units[i+1])))))
				i++
			} else {
				b.WriteString(encodeURIComponent("�"))
			}
		default:
			b.WriteString(encodeURIComponent(string(rune(c))))
		}
	}
	return b.String()
}

// allowedURLChar is /[!#$&-;=?-Z_a-z~]/.
func allowedURLChar(c byte) bool {
	return c == '!' || c == '#' || c == '$' || (c >= '&' && c <= ';') || c == '=' || (c >= '?' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') || c == '~'
}

func encodeURIComponent(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 128 && (c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.IndexByte("-_.!~*'()", c) >= 0) {
			b.WriteByte(c)
			continue
		}
		fmt.Fprintf(&b, "%%%02X", c)
	}
	return b.String()
}
