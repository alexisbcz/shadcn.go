package site

import "strings"

// The reference stylesheet goes through lightningcss (Next.js), which also
// compiles CSS nesting away; Tailwind's output keeps it (.a{&:hover{...}},
// .a{@media (...){...}}). Chromium matches nested rules much more slowly:
// every forced style recalculation of a large docs page takes about four
// times longer (recharts measures text with getBoundingClientRect, so the
// chart page took 13s instead of 3s to load). flattenNesting rewrites the
// nested rules as plain ones, keeping their order:
//
//	.a{color:red;&:hover{color:blue}@media (x){&:focus{color:green}}}
//	-> .a{color:red}.a:hover{color:blue}@media (x){.a:focus{color:green}}
//
// A parent selector list is expanded (.a,.b{&:hover{}} -> .a:hover,.b:hover);
// a complex parent referenced after the start of a nested selector is
// wrapped in :is() (.x &, with & = .a .b -> .x :is(.a .b)), as nesting
// matches it.

// opaqueAtRules keep their body as is (no style rules inside).
var opaqueAtRules = map[string]bool{
	"keyframes": true, "-webkit-keyframes": true, "font-face": true, "property": true,
	"page": true, "counter-style": true, "font-feature-values": true, "font-palette-values": true,
}

// flattenNesting flattens the nested style rules of a (minified) stylesheet.
func flattenNesting(css string) string {
	var b strings.Builder
	flattenItems(&b, css, nil)
	return b.String()
}

// cssItem is one item of a block body: a declaration (or at-rule
// statement) ending with ";", or a block "prelude{body}".
type cssItem struct {
	prelude string // block prelude, or the declaration text
	body    string
	block   bool
}

// splitItems splits a block body into its top-level items.
func splitItems(s string) []cssItem {
	var items []cssItem
	start, depth, paren := 0, 0, 0
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\\':
			i++
		case '"', '\'':
			quote = c
		case '/':
			if i+1 < len(s) && s[i+1] == '*' {
				if end := strings.Index(s[i+2:], "*/"); end >= 0 {
					i += end + 3
				} else {
					i = len(s)
				}
			}
		case '(':
			paren++
		case ')':
			if paren > 0 {
				paren--
			}
		case '{':
			if paren == 0 {
				if depth == 0 {
					items = append(items, cssItem{prelude: s[start:i], block: true})
					start = i + 1
				}
				depth++
			}
		case '}':
			if paren == 0 && depth > 0 {
				depth--
				if depth == 0 {
					items[len(items)-1].body = s[start:i]
					start = i + 1
				}
			}
		case ';':
			if depth == 0 && paren == 0 {
				if t := strings.TrimSpace(s[start:i]); t != "" {
					items = append(items, cssItem{prelude: t})
				}
				start = i + 1
			}
		}
	}
	if depth == 0 {
		if t := strings.TrimSpace(s[start:]); t != "" {
			items = append(items, cssItem{prelude: t})
		}
	}
	return items
}

// flattenItems writes the items of body, inside the style rule(s) parent
// (nil at the top level or directly inside a grouping at-rule).
func flattenItems(b *strings.Builder, body string, parent []string) {
	var decls []string
	flush := func() {
		if len(decls) > 0 && parent != nil {
			b.WriteString(strings.Join(parent, ","))
			b.WriteString("{")
			b.WriteString(strings.Join(decls, ";"))
			b.WriteString("}")
		} else if len(decls) > 0 {
			// Statements at the top level (@import, @layer a,b;, @charset).
			for _, d := range decls {
				b.WriteString(d)
				b.WriteString(";")
			}
		}
		decls = decls[:0]
	}
	for _, it := range splitItems(body) {
		if !it.block {
			decls = append(decls, it.prelude)
			continue
		}
		flush()
		prelude := strings.TrimSpace(it.prelude)
		if strings.HasPrefix(prelude, "@") {
			name := prelude[1:]
			if i := strings.IndexAny(name, " ({"); i >= 0 {
				name = name[:i]
			}
			if opaqueAtRules[strings.ToLower(name)] {
				b.WriteString(prelude)
				b.WriteString("{")
				b.WriteString(it.body)
				b.WriteString("}")
				continue
			}
			// A grouping rule: its body holds rules (and, inside a style
			// rule, declarations for the parent).
			var inner strings.Builder
			flattenItems(&inner, it.body, parent)
			if inner.Len() > 0 || parent == nil {
				b.WriteString(prelude)
				b.WriteString("{")
				b.WriteString(inner.String())
				b.WriteString("}")
			}
			continue
		}
		sels := splitSelectors(prelude)
		if parent != nil {
			sels = resolveNested(sels, parent)
		}
		flattenItems(b, it.body, sels)
	}
	flush()
}

// splitSelectors splits a selector list at its top-level commas.
func splitSelectors(s string) []string {
	var out []string
	start, paren, bracket := 0, 0, 0
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\\':
			i++
		case '"', '\'':
			quote = c
		case '(':
			paren++
		case ')':
			paren--
		case '[':
			bracket++
		case ']':
			bracket--
		case ',':
			if paren == 0 && bracket == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	return append(out, strings.TrimSpace(s[start:]))
}

// resolveNested resolves nested selectors against their parent list.
func resolveNested(nested, parent []string) []string {
	var out []string
	for _, n := range nested {
		if !hasNestingSelector(n) {
			// A relative nested selector: "& n" (or "& > n").
			n = "& " + n
		}
		for _, p := range parent {
			out = append(out, substituteNesting(n, p))
		}
	}
	return out
}

func hasNestingSelector(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '&':
			return true
		}
	}
	return false
}

// substituteNesting replaces each & of n by p (by :is(p) where a textual
// replacement would change the meaning: p complex and & not leading).
func substituteNesting(n, p string) string {
	var b strings.Builder
	complexParent := isComplexSelector(p)
	for i := 0; i < len(n); i++ {
		c := n[i]
		if c == '\\' && i+1 < len(n) {
			b.WriteByte(c)
			b.WriteByte(n[i+1])
			i++
			continue
		}
		if c != '&' {
			b.WriteByte(c)
			continue
		}
		if i == 0 || !complexParent {
			b.WriteString(p)
		} else {
			b.WriteString(":is(" + p + ")")
		}
	}
	return b.String()
}

// isComplexSelector reports whether s has a top-level combinator.
func isComplexSelector(s string) bool {
	paren, bracket := 0, 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\':
			i++
		case '(':
			paren++
		case ')':
			paren--
		case '[':
			bracket++
		case ']':
			bracket--
		case ' ', '>', '+', '~':
			if paren == 0 && bracket == 0 {
				return true
			}
		}
	}
	return false
}
