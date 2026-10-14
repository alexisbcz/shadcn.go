// Package mdx parses the MDX of the shadcn/ui docs into an mdast tree,
// exactly as upstream's parser does (micromark with mdxjs and GFM, through
// mdast-util-from-markdown, as configured by fumadocs-mdx 15.0.12).
//
// It supports the constructs the 316 docs pages use: CommonMark blocks and
// inlines without indented code, HTML and autolinks (MDX disables them),
// GFM tables, strikethrough, autolink literals, task lists and footnotes,
// and MDX JSX elements, expressions and ESM imports. Anything else is an
// error, never silently different output.
package mdx

import "shadcn.go/docs/hast"

// Node is an mdast node. Type is the mdast type name ("paragraph",
// "heading", "mdxJsxFlowElement", ...); fields beyond Type and Children are
// used by the node types noted.
type Node struct {
	Type     string
	Children []*Node

	// text, inlineCode, code, mdxFlowExpression, mdxTextExpression,
	// mdxjsEsm.
	Value string
	// heading
	Depth int
	// list
	Ordered bool
	Start   int
	// list, listItem
	Spread bool
	// listItem (GFM task list item): nil when not a task.
	Checked *bool
	// code
	Lang, Meta string
	HasLang    bool
	HasMeta    bool
	// link, image
	URL, Title string
	HasTitle   bool
	// image
	Alt string
	// table
	Align []string
	// footnoteReference, footnoteDefinition
	Identifier, Label string
	// mdxJsxFlowElement, mdxJsxTextElement ("" for fragments)
	Name  string
	Attrs []hast.JsxAttr
	// mdxjsEsm
	Imports []hast.Import
	// Expr is the parsed expression of mdxFlowExpression and
	// mdxTextExpression.
	Expr hast.Expr

	// Data mirrors mdast's node.data (hProperties, etc.).
	Data map[string]any
}

// Text returns the concatenated text values of the node's descendants,
// like fumadocs' flattenNode.
func (n *Node) Text() string {
	if n.Children == nil {
		return n.Value
	}
	s := ""
	for _, c := range n.Children {
		s += c.Text()
	}
	return s
}
