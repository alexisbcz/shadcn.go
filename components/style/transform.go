package style

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"shadcn.go/cn"
)

// This file ports upstream's transformStyleMap
// (packages/shadcn/src/styles/transform-style-map.ts) and the prettier pass
// of the registry build to the reduced trees of tree.go. Every quirk is kept:
// tokens are only replaced in the places upstream looks at, and the
// replacement depends on what was replaced before in the same file.

// allowlist is upstream's ALLOWLIST: tokens that are never removed.
var allowlist = map[string]bool{
	"cn-menu-target":      true,
	"cn-menu-translucent": true,
	"cn-logical-sides":    true,
	"cn-rtl-flip":         true,
	"cn-font-heading":     true,
}

// cnTokenRe is upstream's /\bcn-[\w-]+\b/g. Go's \b and \w are ASCII, like
// JavaScript's without the u flag, and leftmost-first matching picks the
// same match as JavaScript's backtracking.
var cnTokenRe = regexp.MustCompile(`\bcn-[\w-]+\b`)

// extractCnClasses is upstream's extractCnClasses.
func extractCnClasses(s string) []string { return cnTokenRe.FindAllString(s, -1) }

// isJSSpace reports whether r is whitespace for JavaScript's \s and trim().
func isJSSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xa0, 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff:
		return true
	}
	return r >= 0x2000 && r <= 0x200a
}

func trimJSSpace(s string) string { return trimRightJSSpace(trimLeftJSSpace(s)) }

func trimLeftJSSpace(s string) string { return strings.TrimLeftFunc(s, isJSSpace) }

func trimRightJSSpace(s string) string { return strings.TrimRightFunc(s, isJSSpace) }

// collapseJSSpace is s.replace(/\s+/g, " ").
func collapseJSSpace(s string) string {
	var b strings.Builder
	space := false
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && n == 1 {
			// Not valid UTF-8: keep the byte as is.
			if space {
				b.WriteByte(' ')
				space = false
			}
			b.WriteByte(s[i])
			i++
			continue
		}
		if isJSSpace(r) {
			space = true
		} else {
			if space {
				b.WriteByte(' ')
				space = false
			}
			b.WriteString(s[i : i+n])
		}
		i += n
	}
	if space {
		b.WriteByte(' ')
	}
	return b.String()
}

// removeCnClasses is upstream's removeCnClasses: drops cn-* tokens (except
// allowlisted ones), collapses whitespace and trims.
func removeCnClasses(s string) string {
	s = cnTokenRe.ReplaceAllStringFunc(s, func(m string) string {
		if allowlist[m] {
			return m
		}
		return ""
	})
	return trimJSSpace(collapseJSSpace(s))
}

// mergeClasses is upstream's mergeClasses: twMerge from the cn package.
func mergeClasses(newClasses, existing string) string { return cn.Cn(newClasses, existing) }

// transformer holds the state of one transformStyleMap run (one file).
type transformer struct {
	styleMap map[string]string
	matched  map[string]bool
}

// stylesFor returns the style-map classes of the unmatched, non-allowlisted
// tokens, joined, or "" when there are none.
func (t *transformer) stylesFor(unmatched []string) string {
	var parts []string
	for _, c := range unmatched {
		if allowlist[c] {
			continue
		}
		if v := t.styleMap[c]; v != "" {
			parts = append(parts, v)
		}
	}
	return strings.Join(parts, " ")
}

func (t *transformer) unmatched(tokens []string) []string {
	var out []string
	for _, c := range tokens {
		if !t.matched[c] {
			out = append(out, c)
		}
	}
	return out
}

// transformStyleMap runs upstream's transformStyleMap on the roots of one
// file: cva calls, then JSX className attributes, then className
// properties, then mergeProps calls.
func transformStyleMap(roots []*node, styleMap map[string]string) {
	t := &transformer{styleMap: styleMap, matched: map[string]bool{}}
	walkRoots(roots, t.cvaCall)
	walkRoots(roots, t.classNameAttribute)
	walkRoots(roots, t.classNameProperty)
	walkRoots(roots, t.mergePropsCall)
}

// applyStyleToCvaString is upstream's applyStyleToCvaString.
func (t *transformer) applyStyleToCvaString(s *node) {
	value := s.str
	tokens := extractCnClasses(value)
	if len(tokens) == 0 {
		return
	}
	unmatched := t.unmatched(tokens)
	if len(unmatched) == 0 {
		s.str = removeCnClasses(value)
		return
	}
	if tw := t.stylesFor(unmatched); tw != "" {
		s.str = removeCnClasses(mergeClasses(tw, value))
		for _, c := range unmatched {
			t.matched[c] = true
		}
	} else {
		s.str = removeCnClasses(value)
	}
}

// cvaCall is applyToCvaCalls for one node.
func (t *transformer) cvaCall(n *node) {
	if n.kind != kCall || n.callee != "cva" {
		return
	}
	if len(n.args) > 0 && n.args[0].isStringLiteral() {
		t.applyStyleToCvaString(n.args[0])
	}
	if len(n.args) < 2 || n.args[1].kind != kObject {
		return
	}
	var variants *node
	for _, p := range n.args[1].children {
		if p.kind == kProp && p.ident && p.name == "variants" {
			variants = p
			break
		}
	}
	if variants == nil || variants.value == nil || variants.value.kind != kObject {
		return
	}
	for _, typeProp := range variants.value.children {
		if typeProp.kind != kProp || typeProp.value == nil || typeProp.value.kind != kObject {
			continue
		}
		for _, variantProp := range typeProp.value.children {
			if variantProp.kind != kProp {
				continue
			}
			if v := variantProp.value; v != nil && v.isStringLiteral() {
				t.applyStyleToCvaString(v)
			}
		}
	}
}

// classNameAttribute is applyToClassNameAttributes for one node.
func (t *transformer) classNameAttribute(n *node) {
	if n.kind != kAttr || n.name != "className" {
		return
	}
	init := n.value
	if init == nil {
		return
	}
	if nested := nestedCnCalls(init); len(nested) > 0 {
		for _, call := range nested {
			t.applyStyleToCnCall(call, false)
		}
		return
	}
	tokens := extractCnClassesFromAttribute(init)
	if len(tokens) == 0 {
		return
	}
	unmatched := t.unmatched(tokens)
	if len(unmatched) == 0 {
		cleanCnClassesFromAttribute(init)
		return
	}
	if tw := t.stylesFor(unmatched); tw != "" {
		applyClassesToElement(n, tw)
	} else {
		cleanCnClassesFromAttribute(init)
	}
}

// nestedCnCalls is getNestedCnCallsFromAttribute.
func nestedCnCalls(init *node) []*node {
	if init.kind != kJsxExpr {
		return nil
	}
	expr := init.value
	if expr == nil || expr.isCnCall() {
		return nil
	}
	var calls []*node
	expr.forEachDescendant(func(d *node) {
		if d.isCnCall() {
			calls = append(calls, d)
		}
	})
	return calls
}

// extractCnClassesFromAttribute is upstream's function of the same name.
func extractCnClassesFromAttribute(init *node) []string {
	if init.isStringLiteralLike() {
		return extractCnClasses(init.str)
	}
	if init.kind != kJsxExpr || init.value == nil {
		return nil
	}
	expr := init.value
	if expr.isStringLiteralLike() {
		return extractCnClasses(expr.str)
	}
	var out []string
	if expr.isCnCall() {
		for _, a := range expr.args {
			if a.isStringLiteralLike() {
				out = append(out, extractCnClasses(a.str)...)
			}
		}
	}
	return out
}

// cleanCnClassesFromAttribute is upstream's function of the same name.
func cleanCnClassesFromAttribute(init *node) {
	if init.isStringLiteralLike() {
		init.str = removeCnClasses(init.str)
		return
	}
	if init.kind != kJsxExpr || init.value == nil {
		return
	}
	expr := init.value
	if expr.isStringLiteralLike() {
		expr.str = removeCnClasses(expr.str)
		return
	}
	if expr.isCnCall() {
		for _, a := range expr.args {
			if a.isStringLiteralLike() {
				a.str = removeCnClasses(a.str)
			}
		}
		removeEmptyArgumentsFromCnCall(expr)
	}
}

// removeEmptyArgumentsFromCnCall drops the string arguments that are empty
// after trimming. Upstream rewrites the call's text; the result is the same
// call without those arguments.
func removeEmptyArgumentsFromCnCall(call *node) {
	if !call.isCnCall() {
		return
	}
	kept := make([]*node, 0, len(call.args))
	for _, a := range call.args {
		if a.isStringLiteralLike() && trimJSSpace(a.str) == "" {
			continue
		}
		kept = append(kept, a)
	}
	if len(kept) != len(call.args) {
		call.args = kept
	}
}

// newCnCall is the call upstream writes as `cn(${JSON.stringify(classes)}, ...)`.
func newCnCall(args ...*node) *node {
	return &node{kind: kCall, callee: "cn", calleeRoot: "cn", calleeNode: &node{kind: kID, str: "cn"}, args: args}
}

// applyClassesToElement is upstream's function of the same name, on the
// element's className attribute.
func applyClassesToElement(attr *node, tw string) {
	init := attr.value
	if init == nil {
		attr.value = &node{kind: kJsxExpr, value: newCnCall(strNode(tw))}
		return
	}
	if init.isStringLiteralLike() {
		init.str = removeCnClasses(mergeClasses(tw, init.str))
		return
	}
	if init.kind != kJsxExpr {
		return
	}
	expr := init.value
	if expr == nil {
		attr.value = &node{kind: kJsxExpr, value: newCnCall(strNode(tw))}
		return
	}
	if expr.isStringLiteralLike() {
		expr.str = removeCnClasses(mergeClasses(tw, expr.str))
		return
	}
	if expr.isCnCall() {
		if applyToFirstArgument(expr, tw) {
			return
		}
		attr.value = &node{kind: kJsxExpr, value: newCnCall(rebuiltArguments(expr, tw)...)}
		return
	}
	attr.value = &node{kind: kJsxExpr, value: newCnCall(strNode(tw), expr)}
}

// applyToFirstArgument merges tw into the call's first argument when it is a
// string, cleans the other string arguments and drops empty ones. It
// reports false, changing nothing, when the first argument is not a string.
func applyToFirstArgument(call *node, tw string) bool {
	if len(call.args) == 0 || !call.args[0].isStringLiteralLike() {
		return false
	}
	first := call.args[0]
	first.str = removeCnClasses(mergeClasses(tw, first.str))
	for _, a := range call.args[1:] {
		if a.isStringLiteralLike() {
			a.str = removeCnClasses(a.str)
		}
	}
	removeEmptyArgumentsFromCnCall(call)
	return true
}

// rebuiltArguments are the arguments upstream writes when the first
// argument is not a string: the style classes first, then every argument,
// strings cleaned and dropped when empty.
func rebuiltArguments(call *node, tw string) []*node {
	args := []*node{strNode(tw)}
	for _, a := range call.args {
		if a.isStringLiteralLike() {
			if c := removeCnClasses(a.str); c != "" {
				args = append(args, strNode(c))
			}
			continue
		}
		args = append(args, a)
	}
	return args
}

// classNameProperty is applyToClassNameProperties for one node.
func (t *transformer) classNameProperty(n *node) {
	if n.kind != kProp || n.name != "className" {
		return
	}
	if n.value != nil && n.value.isCnCall() {
		t.applyStyleToCnCall(n.value, false)
	}
}

// mergePropsCall is applyToMergePropsCalls for one node.
func (t *transformer) mergePropsCall(n *node) {
	if n.kind != kCall || n.callee != "mergeProps" {
		return
	}
	for _, arg := range n.args {
		if arg.kind != kObject {
			continue
		}
		var prop *node
		for _, p := range arg.children {
			if p.kind == kProp && p.ident && p.name == "className" {
				prop = p
				break
			}
		}
		if prop == nil || prop.value == nil {
			continue
		}
		if prop.value.isCnCall() {
			t.applyStyleToCnCall(prop.value, true)
		}
	}
}

// applyStyleToCnCall is upstream's function of the same name.
func (t *transformer) applyStyleToCnCall(call *node, markMatches bool) {
	var tokens []string
	for _, a := range call.args {
		if a.isStringLiteralLike() {
			tokens = append(tokens, extractCnClasses(a.str)...)
		}
	}
	if len(tokens) == 0 {
		return
	}
	unmatched := t.unmatched(tokens)
	if len(unmatched) == 0 {
		cleanCnClassesFromCnCall(call)
		return
	}
	tw := t.stylesFor(unmatched)
	if tw == "" {
		cleanCnClassesFromCnCall(call)
		return
	}
	if !applyToFirstArgument(call, tw) {
		// Upstream replaces the call with a new one; the node is rewritten in
		// place so its parent sees the new call.
		call.args = rebuiltArguments(call, tw)
	}
	if markMatches {
		for _, c := range unmatched {
			t.matched[c] = true
		}
	}
}

// cleanCnClassesFromCnCall is upstream's function of the same name.
func cleanCnClassesFromCnCall(call *node) {
	for _, a := range call.args {
		if a.isStringLiteralLike() {
			a.str = removeCnClasses(a.str)
		}
	}
	removeEmptyArgumentsFromCnCall(call)
}

// ---------------------------------------------------------------------------
// Prettier: prettier-plugin-tailwindcss sorts the class strings of className
// attributes and of every literal inside cn() and cva() calls.

var sortFunctions = map[string]bool{"cn": true, "cva": true}

// sortTree runs the plugin's JavaScript transform on the roots of one file.
func sortTree(roots []*node) {
	walkRoots(roots, func(n *node) {
		switch {
		case n.kind == kAttr && (n.name == "className" || n.name == "class") && n.value != nil:
			if n.value.kind == kStr {
				sortLiteral(n.value, collapseAll)
			} else if n.value.kind == kJsxExpr {
				sortExpression(n.value)
			}
		case n.kind == kCall && len(n.args) > 0 && sortFunctions[n.calleeRoot]:
			for _, a := range n.args {
				sortExpression(a)
			}
		}
	})
}

// sortLiteral sorts a string literal (or no-substitution template).
func sortLiteral(n *node, c collapse) {
	n.str = sortClasses(n.str, sortOpts{collapse: c})
}

// sortExpression is the plugin's walk over an expression: every string and
// template literal in it is sorted, with whitespace kept at the edges that
// touch a + concatenation.
func sortExpression(root *node) {
	var walk func(n *node, c collapse)
	walk = func(n *node, c collapse) {
		switch n.kind {
		case kStr:
			sortLiteral(n, c)
		case kTpl:
			sortTemplate(n, c)
		case kGeneric:
			if n.str == "TaggedTemplateExpression" && len(n.children) == 2 {
				if tag := n.children[0]; tagRoot(tag) != "" && sortFunctions[tagRoot(tag)] {
					if q := n.children[1]; q.kind == kStr {
						sortLiteral(q, c)
					} else if q.kind == kTpl {
						sortTemplate(q, c)
					}
				}
			}
		}
		switch n.kind {
		case kBinary:
			if n.op == "+" {
				walk(n.left, collapse{c.start, false})
				walk(n.right, collapse{false, c.end})
				return
			}
		case kTpl:
			for i, e := range n.exprs {
				// A literal directly inside ${...}: the plugin compares
				// positions with the template's quasis (see templateCollapse).
				walk(e, templateCollapse(n, i, c))
			}
			return
		}
		n.kids(func(k *node) { walk(k, c) })
	}
	walk(root, collapseAll)
}

// templateCollapse reproduces the plugin's getCollapseWhitespace for a node
// whose parent is a template literal. Its position checks compare the
// node's start with each quasi's end (so the quasi just before "${" and all
// quasis after the node match) and the node's end with each quasi's end (so
// all quasis after the node match).
func templateCollapse(t *node, i int, c collapse) collapse {
	for j, q := range t.quasis {
		raw := q[0]
		if j == i || j > i { // the quasi before the substitution, and the ones after
			if c.start {
				c.start = raw != "" && isJSSpaceStart(raw)
			}
		}
		if j > i {
			if c.end {
				c.end = raw != "" && isJSSpaceEnd(raw)
			}
		}
	}
	return c
}

func isJSSpaceStart(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return isJSSpace(r)
}

func isJSSpaceEnd(s string) bool {
	r, _ := utf8.DecodeLastRuneInString(s)
	return isJSSpace(r)
}

func tagRoot(n *node) string {
	for {
		switch {
		case n.kind == kID:
			return n.str
		case n.kind == kCall:
			return n.calleeRoot
		case n.kind == kGeneric && n.str == "PropertyAccessExpression" && len(n.children) > 0:
			n = n.children[0]
		default:
			return ""
		}
	}
}

// sortTemplate is the plugin's sortTemplateLiteral.
func sortTemplate(n *node, c collapse) {
	last := len(n.exprs)
	for i := range n.quasis {
		for k := 0; k < 2; k++ {
			s := n.quasis[i][k]
			n.quasis[i][k] = sortClasses(s, sortOpts{
				ignoreFirst: i > 0 && !(s != "" && isJSSpaceStart(s)),
				ignoreLast:  i < last && !(s != "" && isJSSpaceEnd(s)),
				collapse:    collapse{c.start && i == 0, c.end && i >= last},
			})
		}
	}
}
