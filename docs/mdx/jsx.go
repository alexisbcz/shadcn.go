package mdx

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode/utf8"

	"shadcn.go/docs/hast"
)

// jsxTag is one MDX JSX tag: <Name attrs>, </Name>, <Name /> or a fragment.
type jsxTag struct {
	name        string
	closing     bool
	selfClosing bool
	attrs       []hast.JsxAttr
}

// errIncomplete reports that the text ended inside a tag or expression.
var errIncomplete = errors.New("incomplete")

// errNotTag reports a "<" that does not start a tag (followed by
// whitespace), which micromark treats as plain text.
var errNotTag = errors.New("not a tag")

func isSpace(c byte) bool { return c == ' ' || c == '\t' }

func isWS(c byte) bool { return c == ' ' || c == '\t' || c == '\n' }

// parseTag parses a JSX tag at s[i] == '<' and returns it with the offset
// after its '>'. Whitespace, including line endings, may appear between
// the parts of a tag.
func parseTag(s string, i int) (*jsxTag, int, error) {
	t := &jsxTag{}
	i++ // <
	if i >= len(s) {
		return nil, 0, errIncomplete
	}
	if isWS(s[i]) {
		return nil, 0, errNotTag
	}
	skip := func() {
		for i < len(s) && (isWS(s[i]) || isUnicodeWS(s, i)) {
			_, size := utf8.DecodeRuneInString(s[i:])
			i += size
		}
	}
	readName := func() (string, error) {
		start := i
		r, size := utf8.DecodeRuneInString(s[i:])
		if size == 0 {
			return "", errIncomplete
		}
		if !isIDStart(r) {
			return "", fmt.Errorf("mdx: unexpected %q in JSX tag name", r)
		}
		i += size
		for i < len(s) {
			r, size := utf8.DecodeRuneInString(s[i:])
			if !(isIDCont(r) || r == '-') {
				break
			}
			i += size
		}
		return s[start:i], nil
	}
	skip()
	if i >= len(s) {
		return nil, 0, errIncomplete
	}
	if s[i] == '/' {
		t.closing = true
		i++
		skip()
	}
	if i >= len(s) {
		return nil, 0, errIncomplete
	}
	if s[i] != '>' {
		name, err := readName()
		if err != nil {
			return nil, 0, err
		}
		skip()
		for i < len(s) && (s[i] == '.' || s[i] == ':') {
			sep := s[i]
			i++
			skip()
			part, err := readName()
			if err != nil {
				return nil, 0, err
			}
			name += string(sep) + part
			skip()
		}
		t.name = name
	}
	// Attributes.
	for {
		skip()
		if i >= len(s) {
			return nil, 0, errIncomplete
		}
		switch c := s[i]; {
		case c == '>':
			return t, i + 1, nil
		case c == '/':
			if t.closing {
				return nil, 0, errors.New("mdx: unexpected self-closing slash in closing tag")
			}
			i++
			skip()
			if i >= len(s) {
				return nil, 0, errIncomplete
			}
			if s[i] != '>' {
				return nil, 0, errors.New("mdx: expected > after self-closing slash")
			}
			t.selfClosing = true
			return t, i + 1, nil
		case c == '{':
			return nil, 0, errors.New("mdx: JSX spread attributes are not supported")
		default:
			if t.closing {
				return nil, 0, errors.New("mdx: unexpected attribute in closing tag")
			}
			name, err := readName()
			if err != nil {
				return nil, 0, err
			}
			skip()
			if i < len(s) && s[i] == ':' {
				i++
				skip()
				local, err := readName()
				if err != nil {
					return nil, 0, err
				}
				name += ":" + local
				skip()
			}
			if i >= len(s) {
				return nil, 0, errIncomplete
			}
			if s[i] != '=' {
				t.attrs = append(t.attrs, hast.JsxAttr{Name: name, Kind: hast.AttrBool})
				continue
			}
			i++
			skip()
			if i >= len(s) {
				return nil, 0, errIncomplete
			}
			switch s[i] {
			case '"', '\'':
				q := s[i]
				end := strings.IndexByte(s[i+1:], q)
				if end < 0 {
					return nil, 0, errIncomplete
				}
				v := s[i+1 : i+1+end]
				if strings.Contains(v, "\n") {
					return nil, 0, errors.New("mdx: multi-line attribute values are not supported")
				}
				t.attrs = append(t.attrs, hast.JsxAttr{Name: name, Kind: hast.AttrString, Value: decodeEntities(v)})
				i += end + 2
			case '{':
				e, end, err := parseExpressionBody(s, i+1)
				if err != nil {
					if isUnterminated(s, i+1) {
						return nil, 0, errIncomplete
					}
					return nil, 0, err
				}
				if _, ok := e.(Empty); ok {
					return nil, 0, errors.New("mdx: empty attribute expression")
				}
				t.attrs = append(t.attrs, hast.JsxAttr{Name: name, Kind: hast.AttrExpression, Value: exprValue(s[i+1 : end]), Expr: e})
				i = end + 1
			default:
				return nil, 0, fmt.Errorf("mdx: unexpected %q before attribute value", s[i])
			}
		}
	}
}

func isUnicodeWS(s string, i int) bool {
	r, _ := utf8.DecodeRuneInString(s[i:])
	return r == '\u00a0' || r == '\u2028' || r == '\u2029' || r == '\ufeff'
}

// isUnterminated reports whether an expression starting at s[i] runs to the
// end of s without its closing brace (so more lines may complete it).
func isUnterminated(s string, i int) bool {
	depth := 0
	for ; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			if depth == 0 {
				return false
			}
			depth--
		}
	}
	return true
}

// exprValue is the value mdast keeps for an expression: the source between
// the braces, where each continuation line loses up to two leading spaces
// or tabs (micromark-factory-mdx-expression's line prefix).
func exprValue(raw string) string {
	if !strings.Contains(raw, "\n") {
		return raw
	}
	lines := strings.Split(raw, "\n")
	for i := 1; i < len(lines); i++ {
		l := lines[i]
		n := 0
		for n < 2 && n < len(l) && isSpace(l[n]) {
			n++
		}
		lines[i] = l[n:]
	}
	return strings.Join(lines, "\n")
}

var entityRe = regexp.MustCompile(`&(?:#[0-9]+|#[xX][0-9a-fA-F]+|[a-zA-Z][a-zA-Z0-9]*);`)

// decodeEntities decodes terminated character references, as
// parse-entities does for JSX attribute values.
func decodeEntities(s string) string {
	if !strings.Contains(s, "&") {
		return s
	}
	return entityRe.ReplaceAllStringFunc(s, func(m string) string {
		d := html.UnescapeString(m)
		if d == m {
			return m
		}
		return d
	})
}

// scanFlowLine reads JSX flow tags and flow expressions at the start of s
// (which may span several lines). complete is false when s ends inside a
// tag or expression; ok is false when s is not a JSX flow line (text
// follows a tag on its last line, or "<" starts no tag).
func scanFlowLine(s string, multi bool) (items []flowItem, complete, ok bool, err error) {
	i := 0
	for {
		for i < len(s) && isSpace(s[i]) {
			i++
		}
		if i >= len(s) {
			return items, true, len(items) > 0, nil
		}
		switch s[i] {
		case '<':
			t, end, err := parseTag(s, i)
			if err == errIncomplete {
				return nil, false, true, nil
			}
			if err == errNotTag {
				return nil, false, false, nil
			}
			if err != nil {
				return nil, false, false, err
			}
			items = append(items, flowItem{tag: t})
			i = end
		case '{':
			e, end, err := parseExpressionBody(s, i+1)
			if err != nil {
				if isUnterminated(s, i+1) {
					return nil, false, true, nil
				}
				return nil, false, false, err
			}
			items = append(items, flowItem{expr: &Node{Type: "mdxFlowExpression", Value: exprValue(s[i+1 : end]), Expr: e}})
			i = end + 1
		case '\n':
			// A tag or expression ended at a line end: the flow construct ends
			// here, so anything on later lines is not part of it.
			if multi {
				return nil, false, false, errors.New("mdx: unexpected content after a multi-line JSX tag")
			}
			return items, true, true, nil
		default:
			return nil, false, false, nil
		}
	}
}

// tabIndented reports whether a continuation line of raw starts its
// indentation with a tab or has one in the first two columns: micromark
// strips up to two columns of indentation from expression lines, counting
// tabs as virtual spaces, which is not reproduced.
func tabIndented(raw string) bool {
	for _, l := range strings.Split(raw, "\n")[1:] {
		for i := 0; i < len(l) && i < 3; i++ {
			if l[i] == '\t' {
				return true
			}
			if l[i] != ' ' {
				break
			}
		}
	}
	return false
}
