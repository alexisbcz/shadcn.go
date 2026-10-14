package style

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// CreateStyleMap is upstream's createStyleMap
// (packages/shadcn/src/styles/create-style-map.ts): it maps every cn-* token
// of a style stylesheet to the Tailwind classes its rules @apply.
//
//   - Every rule is visited, nested ones included, in document order, and
//     only its own @apply at-rules count (joined with spaces, in order).
//   - A rule applies to each of its comma-separated selectors, after "&" and
//     the spaces around it are removed. In each selector the subject is the
//     last class starting with "cn-", wherever it is (inside :has() or
//     :not() too).
//   - When several rules apply to one token, a later rule's classes come
//     first: "later earlier".
//
// It returns an error when the stylesheet cannot be parsed (an unclosed
// block, string or comment).
func CreateStyleMap(css string) (map[string]string, error) {
	p := &cssParser{src: css}
	rules, err := p.parseBlock(false)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	var walk func(items []cssItem)
	walk = func(items []cssItem) {
		for _, it := range items {
			if it.isRule {
				applyRule(it, result)
			}
			walk(it.children)
		}
	}
	walk(rules)
	return result, nil
}

func applyRule(rule cssItem, result map[string]string) {
	var applies []string
	for _, c := range rule.children {
		if !c.isRule && c.atName == "apply" {
			if v := trimJSSpace(c.params); v != "" {
				applies = append(applies, v)
			}
		}
	}
	if len(applies) == 0 {
		return
	}
	tw := strings.Join(applies, " ")
	for _, sel := range splitList(rule.selector, ',') {
		for _, s := range splitSelectors(normalizeSelector(sel)) {
			cls := subjectClass(s)
			if !strings.HasPrefix(cls, "cn-") {
				continue
			}
			if old, ok := result[cls]; ok && old != "" {
				result[cls] = tw + " " + old
			} else {
				result[cls] = tw
			}
		}
	}
}

// normalizeSelector is selector.replace(/\s*&\s*/g, "").trim().
func normalizeSelector(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if j := strings.IndexByte(s[i:], '&'); j >= 0 {
			pre := trimRightJSSpace(s[i : i+j])
			b.WriteString(pre)
			i += j + 1
			rest := trimLeftJSSpace(s[i:])
			i = len(s) - len(rest)
			continue
		}
		b.WriteString(s[i:])
		break
	}
	return trimJSSpace(b.String())
}

// --- a postcss-compatible parser for what style files contain --------------

type cssItem struct {
	isRule   bool
	selector string // rule
	atName   string // at-rule name
	params   string // at-rule params
	children []cssItem
}

type cssParser struct {
	src string
	pos int
}

func (p *cssParser) errorf(format string, args ...any) error {
	line := 1 + strings.Count(p.src[:min(p.pos, len(p.src))], "\n")
	return fmt.Errorf("style: css line %d: %s", line, fmt.Sprintf(format, args...))
}

// parseBlock parses items until the end of input or, when nested, the
// closing brace.
func (p *cssParser) parseBlock(nested bool) ([]cssItem, error) {
	var items []cssItem
	var buf strings.Builder
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		switch {
		case c == '/' && strings.HasPrefix(p.src[p.pos:], "/*"):
			end := strings.Index(p.src[p.pos+2:], "*/")
			if end < 0 {
				return nil, p.errorf("unclosed comment")
			}
			p.pos += end + 4
		case c == '"' || c == '\'':
			start := p.pos
			if err := p.skipString(); err != nil {
				return nil, err
			}
			buf.WriteString(p.src[start:p.pos])
		case c == '\\' && p.pos+1 < len(p.src):
			buf.WriteString(p.src[p.pos : p.pos+2])
			p.pos += 2
		case c == '(':
			start := p.pos
			if err := p.skipParens(); err != nil {
				return nil, err
			}
			buf.WriteString(p.src[start:p.pos])
		case c == ';':
			p.pos++
			if it, ok := statement(buf.String()); ok {
				items = append(items, it)
			}
			buf.Reset()
		case c == '{':
			p.pos++
			head := trimJSSpace(buf.String())
			buf.Reset()
			children, err := p.parseBlock(true)
			if err != nil {
				return nil, err
			}
			if strings.HasPrefix(head, "@") {
				name, params := splitAtRule(head)
				items = append(items, cssItem{atName: name, params: params, children: children})
			} else {
				items = append(items, cssItem{isRule: true, selector: head, children: children})
			}
		case c == '}':
			if !nested {
				return nil, p.errorf("unexpected }")
			}
			p.pos++
			if it, ok := statement(buf.String()); ok {
				items = append(items, it)
			}
			return items, nil
		default:
			buf.WriteByte(c)
			p.pos++
		}
	}
	if nested {
		return nil, p.errorf("unclosed block")
	}
	if it, ok := statement(buf.String()); ok {
		items = append(items, it)
	}
	return items, nil
}

func (p *cssParser) skipString() error {
	q := p.src[p.pos]
	for i := p.pos + 1; i < len(p.src); i++ {
		switch p.src[i] {
		case '\\':
			i++
		case q:
			p.pos = i + 1
			return nil
		case '\n':
			return p.errorf("unclosed string")
		}
	}
	return p.errorf("unclosed string")
}

func (p *cssParser) skipParens() error {
	depth := 0
	for p.pos < len(p.src) {
		switch c := p.src[p.pos]; c {
		case '"', '\'':
			if err := p.skipString(); err != nil {
				return err
			}
			continue
		case '\\':
			p.pos++
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				p.pos++
				return nil
			}
		}
		p.pos++
	}
	return p.errorf("unclosed (")
}

// statement turns the text before a ";" into an at-rule; declarations are
// not needed.
func statement(s string) (cssItem, bool) {
	s = trimJSSpace(s)
	if !strings.HasPrefix(s, "@") {
		return cssItem{}, false
	}
	name, params := splitAtRule(s)
	return cssItem{atName: name, params: params}, true
}

func splitAtRule(s string) (name, params string) {
	s = s[1:]
	i := 0
	for i < len(s) {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '(' || c == '"' || c == '\'' {
			break
		}
		i++
	}
	return s[:i], trimJSSpace(s[i:])
}

// splitList splits on sep outside quotes and parentheses (postcss list).
func splitList(s string, sep byte) []string {
	var out []string
	depth := 0
	var quote byte
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0:
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '\\':
			i++
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case c == sep && depth == 0:
			out = append(out, trimJSSpace(s[start:i]))
			start = i + 1
		}
	}
	return append(out, trimJSSpace(s[start:]))
}

// splitSelectors splits a selector list on top-level commas, outside
// brackets, parentheses and strings (postcss-selector-parser).
func splitSelectors(s string) []string {
	var out []string
	depth := 0
	var quote byte
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0:
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '\\':
			i++
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			if depth > 0 {
				depth--
			}
		case c == ',' && depth == 0:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// subjectClass returns the last class of the selector that starts with
// "cn-" (findSubjectClass walks every class, inside pseudo-class arguments
// too, in document order), or "".
func subjectClass(sel string) string {
	last := ""
	inAttr := false
	var quote byte
	for i := 0; i < len(sel); i++ {
		c := sel[i]
		switch {
		case quote != 0:
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '\\':
			i++
		case c == '[':
			inAttr = true
		case c == ']':
			inAttr = false
		case c == '.' && !inAttr:
			name, n := readIdent(sel[i+1:])
			if strings.HasPrefix(name, "cn-") {
				last = name
			}
			i += n
		}
	}
	return last
}

// readIdent reads a CSS identifier and returns its unescaped value and the
// number of bytes read.
func readIdent(s string) (string, int) {
	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '\\' && i+1 < len(s):
			j := i + 1
			for j < len(s) && j < i+7 && isHex(s[j]) {
				j++
			}
			if j > i+1 {
				cp, _ := strconv.ParseUint(s[i+1:j], 16, 32)
				if cp == 0 || cp > 0x10ffff || (cp >= 0xd800 && cp <= 0xdfff) {
					cp = 0xfffd
				}
				b.WriteRune(rune(cp))
				if j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n') {
					j++
				}
				i = j
				continue
			}
			r, n := utf8.DecodeRuneInString(s[i+1:])
			b.WriteRune(r)
			i += 1 + n
		case c == '-' || c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80:
			b.WriteByte(c)
			i++
		default:
			return b.String(), i
		}
	}
	return b.String(), i
}

func isHex(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
