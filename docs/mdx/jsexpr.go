package mdx

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// The JavaScript subset of MDX expressions ({...} in content and in JSX
// attribute values). upstream parses them with acorn; this parser accepts
// the expressions the docs use (literals, identifiers, arrays, member
// access, calls, arrow functions with array destructuring, and JSX) and
// reports anything else as an error.

// Expr is a parsed JavaScript expression.
type Expr interface{ isExpr() }

type (
	// Empty is an expression container with only comments or whitespace.
	Empty struct{}
	// Lit is a string, number, boolean, null or undefined literal.
	Lit struct{ Value any } // string, float64, bool, nil; Undefined{} for undefined
	// Ident is an identifier reference.
	Ident struct{ Name string }
	// Array is an array literal.
	Array struct{ Elems []Expr }
	// Member is obj.prop.
	Member struct {
		Obj  Expr
		Prop string
	}
	// Call is callee(args).
	Call struct {
		Callee Expr
		Args   []Expr
	}
	// Arrow is (params) => body.
	Arrow struct {
		Params []Pattern
		Body   Expr
	}
	// JSXElement is <Name attrs>children</Name> (Name "" for fragments).
	JSXElement struct {
		Name     string
		Attrs    []JSXAttr
		Children []Expr // JSXText, Expr (from {expr}) or *JSXElement
	}
	// JSXText is text inside a JSX element, already cleaned of
	// insignificant whitespace as JSX compilers do.
	JSXText struct{ Value string }
)

// Undefined is the value of `undefined`.
type Undefined struct{}

// JSXAttr is name="v", name={expr} or a bare name (Value nil = true).
type JSXAttr struct {
	Name  string
	Value Expr
}

// Pattern is a binding pattern in arrow function parameters.
type Pattern struct {
	Name  string    // identifier, when Elems is nil
	Elems []Pattern // array pattern
	Rest  string    // ...rest in an array pattern
	Array bool
}

func (Empty) isExpr()       {}
func (Lit) isExpr()         {}
func (Ident) isExpr()       {}
func (*Array) isExpr()      {}
func (*Member) isExpr()     {}
func (*Call) isExpr()       {}
func (*Arrow) isExpr()      {}
func (*JSXElement) isExpr() {}
func (JSXText) isExpr()     {}

type jsParser struct {
	src string
	pos int
}

func (p *jsParser) errorf(format string, args ...any) error {
	return fmt.Errorf("mdx: unsupported or invalid expression at offset %d (%q): %s", p.pos, snippet(p.src, p.pos), fmt.Sprintf(format, args...))
}

func snippet(s string, i int) string {
	end := min(i+40, len(s))
	if i > len(s) {
		i = len(s)
	}
	return s[i:end]
}

// skip skips whitespace and comments.
func (p *jsParser) skip() error {
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			p.pos++
		case strings.HasPrefix(p.src[p.pos:], "//"):
			i := strings.IndexByte(p.src[p.pos:], '\n')
			if i < 0 {
				p.pos = len(p.src)
			} else {
				p.pos += i
			}
		case strings.HasPrefix(p.src[p.pos:], "/*"):
			i := strings.Index(p.src[p.pos+2:], "*/")
			if i < 0 {
				return p.errorf("unterminated comment")
			}
			p.pos += i + 4
		default:
			return nil
		}
	}
	return nil
}

func (p *jsParser) peek() byte {
	if p.pos < len(p.src) {
		return p.src[p.pos]
	}
	return 0
}

func (p *jsParser) eat(s string) bool {
	if strings.HasPrefix(p.src[p.pos:], s) {
		p.pos += len(s)
		return true
	}
	return false
}

// parseExpressionBody parses the content of an MDX expression starting at
// src[start] (just after the opening brace) and returns the expression and
// the offset of the closing brace.
func parseExpressionBody(src string, start int) (Expr, int, error) {
	p := &jsParser{src: src, pos: start}
	if err := p.skip(); err != nil {
		return nil, 0, err
	}
	if p.peek() == '}' {
		return Empty{}, p.pos, nil
	}
	e, err := p.expr()
	if err != nil {
		return nil, 0, err
	}
	if err := p.skip(); err != nil {
		return nil, 0, err
	}
	if p.peek() != '}' {
		return nil, 0, p.errorf("expected }")
	}
	return e, p.pos, nil
}

func (p *jsParser) expr() (Expr, error) {
	// Arrow function: (params) => body, or ident => body.
	save := p.pos
	if params, ok := p.tryArrowParams(); ok {
		if err := p.skip(); err != nil {
			return nil, err
		}
		body, err := p.expr()
		if err != nil {
			return nil, err
		}
		return &Arrow{Params: params, Body: body}, nil
	}
	p.pos = save
	return p.postfix()
}

func (p *jsParser) tryArrowParams() ([]Pattern, bool) {
	if p.peek() == '(' {
		p.pos++
		var params []Pattern
		for {
			if p.skip() != nil {
				return nil, false
			}
			if p.peek() == ')' {
				p.pos++
				break
			}
			pat, ok := p.pattern()
			if !ok {
				return nil, false
			}
			params = append(params, pat)
			if p.skip() != nil {
				return nil, false
			}
			if p.peek() == ',' {
				p.pos++
				continue
			}
			if p.peek() != ')' {
				return nil, false
			}
		}
		if p.skip() != nil || !p.eat("=>") {
			return nil, false
		}
		return params, true
	}
	if name := p.ident(); name != "" {
		if p.skip() != nil || !p.eat("=>") {
			return nil, false
		}
		return []Pattern{{Name: name}}, true
	}
	return nil, false
}

func (p *jsParser) pattern() (Pattern, bool) {
	if p.peek() == '[' {
		p.pos++
		pat := Pattern{Array: true}
		for {
			if p.skip() != nil {
				return pat, false
			}
			if p.peek() == ']' {
				p.pos++
				return pat, true
			}
			if p.eat("...") {
				name := p.ident()
				if name == "" {
					return pat, false
				}
				pat.Rest = name
			} else {
				el, ok := p.pattern()
				if !ok {
					return pat, false
				}
				pat.Elems = append(pat.Elems, el)
			}
			if p.skip() != nil {
				return pat, false
			}
			if p.peek() == ',' {
				p.pos++
			}
		}
	}
	name := p.ident()
	return Pattern{Name: name}, name != ""
}

func isIDStart(r rune) bool {
	return r == '$' || r == '_' || unicode.IsLetter(r) || unicode.Is(unicode.Nl, r) || unicode.Is(unicode.Other_ID_Start, r)
}

func isIDCont(r rune) bool {
	return isIDStart(r) || unicode.IsDigit(r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) ||
		unicode.Is(unicode.Pc, r) || unicode.Is(unicode.Nd, r) || r == '‌' || r == '‍'
}

func (p *jsParser) ident() string {
	start := p.pos
	r, size := utf8.DecodeRuneInString(p.src[p.pos:])
	if size == 0 || !isIDStart(r) {
		return ""
	}
	p.pos += size
	for p.pos < len(p.src) {
		r, size := utf8.DecodeRuneInString(p.src[p.pos:])
		if !isIDCont(r) {
			break
		}
		p.pos += size
	}
	return p.src[start:p.pos]
}

func (p *jsParser) postfix() (Expr, error) {
	e, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		save := p.pos
		if err := p.skip(); err != nil {
			return nil, err
		}
		switch {
		case p.peek() == '.' && !strings.HasPrefix(p.src[p.pos:], "..."):
			p.pos++
			if err := p.skip(); err != nil {
				return nil, err
			}
			name := p.ident()
			if name == "" {
				return nil, p.errorf("expected property name")
			}
			e = &Member{Obj: e, Prop: name}
		case p.peek() == '(':
			p.pos++
			var args []Expr
			for {
				if err := p.skip(); err != nil {
					return nil, err
				}
				if p.peek() == ')' {
					p.pos++
					break
				}
				a, err := p.expr()
				if err != nil {
					return nil, err
				}
				args = append(args, a)
				if err := p.skip(); err != nil {
					return nil, err
				}
				if p.peek() == ',' {
					p.pos++
				} else if p.peek() != ')' {
					return nil, p.errorf("expected , or )")
				}
			}
			e = &Call{Callee: e, Args: args}
		default:
			p.pos = save
			return e, nil
		}
	}
}

func (p *jsParser) primary() (Expr, error) {
	if err := p.skip(); err != nil {
		return nil, err
	}
	c := p.peek()
	switch {
	case c == '"' || c == '\'':
		s, err := p.stringLit()
		if err != nil {
			return nil, err
		}
		return Lit{Value: s}, nil
	case c >= '0' && c <= '9':
		start := p.pos
		for p.pos < len(p.src) && (p.src[p.pos] >= '0' && p.src[p.pos] <= '9' || p.src[p.pos] == '.') {
			p.pos++
		}
		f, err := strconv.ParseFloat(p.src[start:p.pos], 64)
		if err != nil {
			return nil, p.errorf("bad number")
		}
		return Lit{Value: f}, nil
	case c == '[':
		p.pos++
		arr := &Array{}
		for {
			if err := p.skip(); err != nil {
				return nil, err
			}
			if p.peek() == ']' {
				p.pos++
				return arr, nil
			}
			e, err := p.expr()
			if err != nil {
				return nil, err
			}
			arr.Elems = append(arr.Elems, e)
			if err := p.skip(); err != nil {
				return nil, err
			}
			if p.peek() == ',' {
				p.pos++
			} else if p.peek() != ']' {
				return nil, p.errorf("expected , or ]")
			}
		}
	case c == '(':
		p.pos++
		e, err := p.expr()
		if err != nil {
			return nil, err
		}
		if err := p.skip(); err != nil {
			return nil, err
		}
		if !p.eat(")") {
			return nil, p.errorf("expected )")
		}
		return e, nil
	case c == '<':
		return p.jsxElement()
	}
	name := p.ident()
	if name == "" {
		return nil, p.errorf("unexpected character")
	}
	switch name {
	case "true":
		return Lit{Value: true}, nil
	case "false":
		return Lit{Value: false}, nil
	case "null":
		return Lit{Value: nil}, nil
	case "undefined":
		return Lit{Value: Undefined{}}, nil
	}
	return Ident{Name: name}, nil
}

func (p *jsParser) stringLit() (string, error) {
	q := p.src[p.pos]
	p.pos++
	var b strings.Builder
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		switch c {
		case q:
			p.pos++
			return b.String(), nil
		case '\\':
			if p.pos+1 >= len(p.src) {
				return "", p.errorf("bad escape")
			}
			e := p.src[p.pos+1]
			p.pos += 2
			switch e {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '\\', '\'', '"', '`':
				b.WriteByte(e)
			default:
				return "", p.errorf("unsupported escape \\%c", e)
			}
		case '\n':
			return "", p.errorf("unterminated string")
		default:
			b.WriteByte(c)
			p.pos++
		}
	}
	return "", p.errorf("unterminated string")
}

// jsxElement parses a JSX element or fragment at '<'.
func (p *jsParser) jsxElement() (*JSXElement, error) {
	p.pos++ // <
	if err := p.skip(); err != nil {
		return nil, err
	}
	el := &JSXElement{}
	if p.peek() != '>' {
		name, err := p.jsxName()
		if err != nil {
			return nil, err
		}
		el.Name = name
		for {
			if err := p.skip(); err != nil {
				return nil, err
			}
			if p.eat("/>") {
				return el, nil
			}
			if p.peek() == '>' {
				break
			}
			if p.peek() == '{' {
				return nil, p.errorf("spread attributes are not supported")
			}
			an, err := p.jsxName()
			if err != nil {
				return nil, err
			}
			attr := JSXAttr{Name: an}
			if err := p.skip(); err != nil {
				return nil, err
			}
			if p.eat("=") {
				if err := p.skip(); err != nil {
					return nil, err
				}
				switch p.peek() {
				case '"', '\'':
					q := p.src[p.pos]
					end := strings.IndexByte(p.src[p.pos+1:], q)
					if end < 0 {
						return nil, p.errorf("unterminated attribute")
					}
					attr.Value = Lit{Value: p.src[p.pos+1 : p.pos+1+end]}
					p.pos += end + 2
				case '{':
					p.pos++
					e, end, err := parseExpressionBody(p.src, p.pos)
					if err != nil {
						return nil, err
					}
					if _, ok := e.(Empty); ok {
						return nil, p.errorf("empty attribute expression")
					}
					attr.Value = e
					p.pos = end + 1
				case '<':
					e, err := p.jsxElement()
					if err != nil {
						return nil, err
					}
					attr.Value = e
				default:
					return nil, p.errorf("bad attribute value")
				}
			}
			el.Attrs = append(el.Attrs, attr)
		}
	}
	p.pos++ // >
	// Children.
	for {
		if p.pos >= len(p.src) {
			return nil, p.errorf("unterminated JSX element")
		}
		switch p.src[p.pos] {
		case '<':
			if strings.HasPrefix(strings.TrimLeft(p.src[p.pos+1:], " \t\n"), "/") {
				p.pos++
				if err := p.skip(); err != nil {
					return nil, err
				}
				p.pos++ // /
				if err := p.skip(); err != nil {
					return nil, err
				}
				if el.Name != "" {
					name, err := p.jsxName()
					if err != nil {
						return nil, err
					}
					if name != el.Name {
						return nil, p.errorf("closing tag </%s> does not match <%s>", name, el.Name)
					}
					if err := p.skip(); err != nil {
						return nil, err
					}
				}
				if !p.eat(">") {
					return nil, p.errorf("expected >")
				}
				return el, nil
			}
			child, err := p.jsxElement()
			if err != nil {
				return nil, err
			}
			el.Children = append(el.Children, child)
		case '{':
			p.pos++
			e, end, err := parseExpressionBody(p.src, p.pos)
			if err != nil {
				return nil, err
			}
			p.pos = end + 1
			if _, ok := e.(Empty); !ok {
				el.Children = append(el.Children, e)
			}
		default:
			start := p.pos
			for p.pos < len(p.src) && p.src[p.pos] != '<' && p.src[p.pos] != '{' {
				p.pos++
			}
			if t := cleanJSXText(p.src[start:p.pos]); t != "" {
				el.Children = append(el.Children, JSXText{Value: t})
			}
		}
	}
}

func (p *jsParser) jsxName() (string, error) {
	start := p.pos
	if p.ident() == "" {
		return "", p.errorf("expected a JSX name")
	}
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == '-' {
			p.pos++
			continue
		}
		if c == '.' || c == ':' {
			p.pos++
			if p.ident() == "" {
				return "", p.errorf("expected a JSX name")
			}
			continue
		}
		r, size := utf8.DecodeRuneInString(p.src[p.pos:])
		if size > 0 && isIDCont(r) {
			p.pos += size
			continue
		}
		break
	}
	return p.src[start:p.pos], nil
}

// cleanJSXText applies the JSX whitespace rules (Babel's
// cleanJSXElementLiteralChild): lines are trimmed, lines that are only
// whitespace are dropped, and the remaining lines are joined by a space.
func cleanJSXText(s string) string {
	lines := strings.Split(s, "\n")
	lastNonEmpty := 0
	for i, l := range lines {
		if strings.TrimLeft(l, " \t\r") != "" {
			lastNonEmpty = i
		}
	}
	var out strings.Builder
	for i, l := range lines {
		isFirst, isLast := i == 0, i == len(lines)-1
		isLastNonEmpty := i == lastNonEmpty
		trimmed := strings.ReplaceAll(l, "\t", " ")
		if !isFirst {
			trimmed = strings.TrimLeft(trimmed, " ")
		}
		if !isLast {
			trimmed = strings.TrimRight(trimmed, " ")
		}
		if trimmed != "" {
			if !isLastNonEmpty {
				trimmed += " "
			}
			out.WriteString(trimmed)
		}
	}
	return out.String()
}
