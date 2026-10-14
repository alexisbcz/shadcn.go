package mdx

import (
	"fmt"
	"regexp"
	"strings"

	"shadcn.go/docs/hast"
)

// build converts the block tree to mdast: inline content is parsed, and
// JSX flow tags are nested into elements within each container, as
// mdast-util-mdx-jsx does.
func build(doc *block) (*Node, error) {
	footnotes := map[string]bool{}
	var collect func(b *block)
	collect = func(b *block) {
		if b.kind == kFootnote {
			footnotes[normalizeIdentifier(b.label)] = true
		}
		for _, c := range b.children {
			collect(c)
		}
	}
	collect(doc)
	bd := &builder{footnotes: footnotes}
	kids, err := bd.children(doc, "root")
	if err != nil {
		return nil, err
	}
	root := &Node{Type: "root", Children: kids}
	if err := checkAutolinkTransform(root); err != nil {
		return nil, err
	}
	if err := checkVerified(root); err != nil {
		return nil, err
	}
	return root, nil
}

type builder struct {
	footnotes map[string]bool
}

func (bd *builder) children(b *block, container string) ([]*Node, error) {
	type frame struct {
		el   *Node
		list []*Node
	}
	stack := []frame{{}}
	add := func(n *Node) {
		top := &stack[len(stack)-1]
		top.list = append(top.list, n)
	}
	for _, c := range b.children {
		if c.kind == kJSX {
			for _, it := range c.items {
				if it.expr != nil {
					add(it.expr)
					continue
				}
				t := it.tag
				switch {
				case t.closing:
					if len(stack) < 2 {
						return nil, fmt.Errorf("mdx: line %d: unexpected closing tag </%s>", c.line, t.name)
					}
					top := stack[len(stack)-1]
					if top.el.Name != t.name {
						return nil, fmt.Errorf("mdx: line %d: unexpected closing tag </%s>, expected </%s>", c.line, t.name, top.el.Name)
					}
					top.el.Children = top.list
					if top.el.Children == nil {
						top.el.Children = []*Node{}
					}
					stack = stack[:len(stack)-1]
					add(top.el)
				case t.selfClosing:
					add(&Node{Type: "mdxJsxFlowElement", Name: t.name, Attrs: t.attrs, Children: []*Node{}})
				default:
					stack = append(stack, frame{el: &Node{Type: "mdxJsxFlowElement", Name: t.name, Attrs: t.attrs}})
				}
			}
			continue
		}
		n, err := bd.node(c)
		if err != nil {
			return nil, err
		}
		if n != nil {
			add(n)
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("mdx: expected a closing tag for <%s> before the end of %s", stack[len(stack)-1].el.Name, container)
	}
	return stack[0].list, nil
}

// gapBefore reports whether a blank line separates children[i] from the
// child before it.
func gapBefore(children []*block, i int) bool {
	return children[i].line > children[i-1].lastLine+1
}

var taskRe = regexp.MustCompile(`^\[[ xX]\][ \t]`)

func (bd *builder) node(b *block) (*Node, error) {
	switch b.kind {
	case kParagraph:
		// Continuation lines keep their indentation: code spans and
		// expressions keep it (micromark), text drops it (parseInline).
		text := strings.Join(b.lines, "\n")
		if definitionRe.MatchString(text) {
			return nil, fmt.Errorf("mdx: line %d: link reference definitions are not supported", b.line)
		}
		kids, err := parseInline(text, bd.footnotes, false)
		if err != nil {
			return nil, err
		}
		return &Node{Type: "paragraph", Children: kids}, nil
	case kHeading:
		kids, err := parseInline(strings.TrimSpace(b.lines[0]), bd.footnotes, false)
		if err != nil {
			return nil, err
		}
		return &Node{Type: "heading", Depth: b.level, Children: kids}, nil
	case kThematic:
		return &Node{Type: "thematicBreak"}, nil
	case kFence:
		n := &Node{Type: "code", Value: strings.Join(b.lines, "\n")}
		if b.info != "" {
			lang, meta, _ := strings.Cut(b.info, " ")
			if i := strings.IndexAny(b.info, " \t"); i >= 0 {
				lang, meta = b.info[:i], b.info[i+1:]
			}
			n.Lang, n.HasLang = unescapeString(lang), true
			if meta = strings.TrimSpace(meta); meta != "" {
				n.Meta, n.HasMeta = unescapeString(meta), true
			}
		}
		return n, nil
	case kBlockquote:
		kids, err := bd.children(b, "blockquote")
		if err != nil {
			return nil, err
		}
		return &Node{Type: "blockquote", Children: kids}, nil
	case kList:
		n := &Node{Type: "list", Ordered: b.ordered, Start: b.start}
		for i, it := range b.children {
			if i > 0 && gapBefore(b.children, i) {
				n.Spread = true
			}
			item, err := bd.item(it)
			if err != nil {
				return nil, err
			}
			n.Children = append(n.Children, item)
		}
		return n, nil
	case kFootnote:
		kids, err := bd.children(b, "footnoteDefinition")
		if err != nil {
			return nil, err
		}
		return &Node{Type: "footnoteDefinition", Identifier: normalizeIdentifier(b.label), Label: b.label, Children: kids}, nil
	case kTable:
		t := &Node{Type: "table"}
		for _, a := range b.align {
			t.Align = append(t.Align, a)
		}
		for _, row := range b.rows {
			r := &Node{Type: "tableRow"}
			for _, cell := range row {
				kids, err := parseInline(cell, bd.footnotes, true)
				if err != nil {
					return nil, err
				}
				r.Children = append(r.Children, &Node{Type: "tableCell", Children: kids})
			}
			t.Children = append(t.Children, r)
		}
		return t, nil
	case kESM:
		value := strings.Join(b.lines, "\n")
		imports, err := parseImports(value)
		if err != nil {
			return nil, err
		}
		return &Node{Type: "mdxjsEsm", Value: value, Imports: imports}, nil
	}
	return nil, fmt.Errorf("mdx: unexpected block kind %d", b.kind)
}

func (bd *builder) item(b *block) (*Node, error) {
	kids, err := bd.children(b, "listItem")
	if err != nil {
		return nil, err
	}
	n := &Node{Type: "listItem", Children: kids}
	for i := range b.children {
		if i > 0 && gapBefore(b.children, i) {
			n.Spread = true
		}
	}
	if len(kids) > 0 && kids[0].Type == "paragraph" && len(kids[0].Children) > 0 && kids[0].Children[0].Type == "text" && taskRe.MatchString(kids[0].Children[0].Value) {
		return nil, fmt.Errorf("mdx: task list items are not supported")
	}
	return n, nil
}

// definitionRe matches the start of a paragraph that may be a link
// reference definition (label possibly over several lines, with escapes):
// such paragraphs are rejected.
var definitionRe = regexp.MustCompile(`^\[(?:[^\[\]\\]|\\[\s\S])+\]:`)

var identRe = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

var importRe = regexp.MustCompile(`^import\s*\{\s*([^}]*)\}\s*from\s*["']([^"']+)["'];?\s*$`)

// parseImports reads `import { A, B as C } from "source"` statements, the
// only ESM the docs use.
func parseImports(value string) ([]hast.Import, error) {
	var out []hast.Import
	for _, stmt := range strings.Split(value, "\n") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		m := importRe.FindStringSubmatch(stmt)
		if m == nil {
			return nil, fmt.Errorf("mdx: unsupported ESM statement %q", stmt)
		}
		for _, spec := range strings.Split(m[1], ",") {
			spec = strings.TrimSpace(spec)
			if spec == "" {
				continue
			}
			imported, local, found := strings.Cut(spec, " as ")
			if !found {
				local = imported
			}
			imported, local = strings.TrimSpace(imported), strings.TrimSpace(local)
			// Plain identifiers only: acorn (no TypeScript) rejects
			// `type X` specifiers, which MDX then fails on.
			if !identRe.MatchString(imported) || !identRe.MatchString(local) || imported == "type" || local == "type" {
				return nil, fmt.Errorf("mdx: unsupported import specifier %q", spec)
			}
			out = append(out, hast.Import{Local: local, Imported: imported, Source: m[2]})
		}
	}
	return out, nil
}

// verified lists the node types the docs content uses (checked against
// upstream's mdast for all 316 pages). Parse rejects trees with any other
// node type, rather than risk output that differs from upstream.
var verified = map[string]bool{
	"root": true, "paragraph": true, "heading": true, "text": true, "emphasis": true, "strong": true,
	"inlineCode": true, "link": true, "image": true, "code": true, "list": true, "listItem": true,
	"table": true, "tableRow": true, "tableCell": true, "blockquote": true, "thematicBreak": true,
	"footnoteDefinition": true, "footnoteReference": true, "mdxJsxFlowElement": true,
	"mdxJsxTextElement": true, "mdxFlowExpression": true, "mdxTextExpression": true, "mdxjsEsm": true,
}

func checkVerified(n *Node) error {
	if !verified[n.Type] {
		return fmt.Errorf("mdx: %s nodes are not supported (the docs content uses none)", n.Type)
	}
	if (n.Type == "link" || n.Type == "image") && n.HasTitle {
		return fmt.Errorf("mdx: %s titles are not supported (the docs content uses none)", n.Type)
	}
	if n.Type == "link" && n.Data != nil && n.Data["literal"] != nil && n.Data["literal"] != "http" {
		return fmt.Errorf("mdx: %s autolink literals are not supported (the docs content uses only http(s) ones)", n.Data["literal"])
	}
	for _, c := range n.Children {
		if err := checkVerified(c); err != nil {
			return err
		}
	}
	return nil
}
