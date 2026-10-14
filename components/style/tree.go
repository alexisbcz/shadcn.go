package style

// node is a reduced syntax tree of the class-bearing parts of a TSX source:
// enough of it to run upstream's transformStyleMap and the class sorting of
// prettier-plugin-tailwindcss exactly as they run on the full source. The
// render-time API builds these trees from Go calls; the differential test
// decodes them from the data extracted from upstream (testdata/refgen).
type node struct {
	kind kind

	str string // kStr: value; kID: name; kGeneric: syntax kind
	tpl bool   // kStr: a no-substitution template literal (`...`)

	quasis [][2]string // kTpl: [raw, cooked] parts
	exprs  []*node     // kTpl: substitutions

	callee     string // kCall: callee when it is an identifier
	calleeRoot string // kCall: identifier at the root of the callee chain
	calleeNode *node  // kCall: the callee expression
	args       []*node

	name  string // kProp: name as written; kAttr: attribute name
	ident bool   // kProp: the name is an identifier
	key   *node  // kProp: the name, when it is a string literal
	value *node  // kProp: initializer; kAttr: value (nil when absent); kJsxExpr: expression (nil when empty)

	op          string // kBinary: operator
	left, right *node  // kBinary

	children []*node // kObject: properties; kJSX: attributes then children; kGeneric: children
	nAttrs   int     // kJSX: how many of children are attributes
}

type kind uint8

const (
	kStr     kind = iota // string literal
	kTpl                 // template literal with substitutions
	kCall                // call expression
	kObject              // object literal
	kProp                // property assignment
	kAttr                // JSX attribute
	kJsxExpr             // JSX expression container
	kJSX                 // JSX element
	kBinary              // binary expression
	kID                  // identifier
	kGeneric             // anything else
)

// isStringLiteralLike is upstream's isStringLiteralLike: a string literal or
// a no-substitution template literal.
func (n *node) isStringLiteralLike() bool { return n != nil && n.kind == kStr }

// isStringLiteral is ts-morph's Node.isStringLiteral (not templates).
func (n *node) isStringLiteral() bool { return n != nil && n.kind == kStr && !n.tpl }

// isCnCall is upstream's isCnCall.
func (n *node) isCnCall() bool { return n != nil && n.kind == kCall && n.callee == "cn" }

// kids returns the children of n in source order, as ts-morph's
// forEachDescendant visits them.
func (n *node) kids(visit func(*node)) {
	switch n.kind {
	case kTpl:
		for _, e := range n.exprs {
			visit(e)
		}
	case kCall:
		if n.calleeNode != nil {
			visit(n.calleeNode)
		}
		for _, a := range n.args {
			visit(a)
		}
	case kProp:
		if n.key != nil {
			visit(n.key)
		}
		if n.value != nil {
			visit(n.value)
		}
	case kAttr, kJsxExpr:
		if n.value != nil {
			visit(n.value)
		}
	case kBinary:
		visit(n.left)
		visit(n.right)
	case kObject, kJSX, kGeneric:
		for _, c := range n.children {
			visit(c)
		}
	}
}

// forEachDescendant visits the descendants of n in document order (pre-order,
// n excluded). Like ts-morph, children are read after fn returns, so a node
// fn rewrites is traversed in its new form.
func (n *node) forEachDescendant(fn func(*node)) {
	n.kids(func(c *node) {
		fn(c)
		c.forEachDescendant(fn)
	})
}

// walkRoots visits every root and its descendants in document order.
func walkRoots(roots []*node, fn func(*node)) {
	for _, r := range roots {
		fn(r)
		r.forEachDescendant(fn)
	}
}

func strNode(v string) *node { return &node{kind: kStr, str: v} }
