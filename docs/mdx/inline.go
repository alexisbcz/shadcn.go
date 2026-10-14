package mdx

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Inline parsing follows micromark's text tokenizer and its resolvers:
// code spans, character escapes and references, links and images, attention
// (emphasis and strong, CommonMark's delimiter algorithm), GFM
// strikethrough, autolink literals and footnote calls, and MDX JSX tags and
// expressions.

type item struct {
	n          *Node
	prev, next *item
	// delimiter run for * _ ~
	delim *delim
	// JSX tag placeholder
	tag *jsxTag
	// bracket opener placeholder "[" or "!["
	bracket bool
}

type delim struct {
	it       *item
	char     byte
	count    int // remaining
	origLen  int
	canOpen  bool
	canClose bool
	prev     *delim
	next     *delim
}

type bracketEntry struct {
	it           *item
	image        bool
	active       bool
	prev         *bracketEntry
	prevDel      *delim
	index        int // position in src after the bracket
	bracketAfter bool
}

type inlineParser struct {
	src       string
	pos       int
	head      *item
	tail      *item
	delims    *delim
	brackets  *bracketEntry
	footnotes map[string]bool
	inTable   bool
}

func (p *inlineParser) append(it *item) *item {
	it.prev = p.tail
	if p.tail != nil {
		p.tail.next = it
	} else {
		p.head = it
	}
	p.tail = it
	return it
}

func (p *inlineParser) text(s string) {
	if s == "" {
		return
	}
	if p.tail != nil && p.tail.n != nil && p.tail.n.Type == "text" && p.tail.delim == nil && !p.tail.bracket {
		p.tail.n.Value += s
		return
	}
	p.append(&item{n: &Node{Type: "text", Value: s}})
}

func (p *inlineParser) remove(it *item) {
	if it.prev != nil {
		it.prev.next = it.next
	} else {
		p.head = it.next
	}
	if it.next != nil {
		it.next.prev = it.prev
	} else {
		p.tail = it.prev
	}
	it.prev, it.next = nil, nil
}

// parseInline parses paragraph (or heading, cell) content.
func parseInline(src string, footnotes map[string]bool, inTable bool) ([]*Node, error) {
	p := &inlineParser{src: src, footnotes: footnotes, inTable: inTable}
	if err := p.run(); err != nil {
		return nil, err
	}
	p.processEmphasis(nil)
	nodes := p.collect(p.head, nil)
	nodes, err := nestTags(nodes)
	if err != nil {
		return nil, err
	}
	nodes = mergeText(nodes)
	return nodes, nil
}

// collect converts the item list between from and to (exclusive) to nodes.
func (p *inlineParser) collect(from, to *item) []*Node {
	var out []*Node
	for it := from; it != nil && it != to; it = it.next {
		switch {
		case it.tag != nil:
			out = append(out, &Node{Type: "\x00tag", Data: map[string]any{"tag": it.tag}})
		case it.delim != nil || it.bracket:
			out = append(out, &Node{Type: "text", Value: it.n.Value})
		default:
			out = append(out, it.n)
		}
	}
	return out
}

func mergeText(nodes []*Node) []*Node {
	var out []*Node
	for _, n := range nodes {
		if n.Type == "text" && n.Value == "" {
			continue
		}
		if n.Type == "text" && len(out) > 0 && out[len(out)-1].Type == "text" {
			out[len(out)-1] = &Node{Type: "text", Value: out[len(out)-1].Value + n.Value}
			continue
		}
		if len(n.Children) > 0 && n.Type != "mdxJsxTextElement" {
			n.Children = mergeText(n.Children)
		} else if n.Type == "mdxJsxTextElement" {
			n.Children = mergeText(n.Children)
		}
		out = append(out, n)
	}
	return out
}

// nestTags turns JSX open/close tag placeholders into elements.
func nestTags(nodes []*Node) ([]*Node, error) {
	type frame struct {
		el   *Node
		list []*Node
	}
	stack := []frame{{list: nil}}
	for _, n := range nodes {
		if len(n.Children) > 0 && n.Type != "\x00tag" {
			kids, err := nestTags(n.Children)
			if err != nil {
				return nil, err
			}
			n.Children = kids
		}
		if n.Type != "\x00tag" {
			top := &stack[len(stack)-1]
			top.list = append(top.list, n)
			continue
		}
		t := n.Data["tag"].(*jsxTag)
		switch {
		case t.closing:
			if len(stack) < 2 {
				return nil, fmt.Errorf("mdx: unexpected closing tag </%s>", t.name)
			}
			top := stack[len(stack)-1]
			if top.el.Name != t.name {
				return nil, fmt.Errorf("mdx: unexpected closing tag </%s>, expected </%s>", t.name, top.el.Name)
			}
			top.el.Children = top.list
			if top.el.Children == nil {
				top.el.Children = []*Node{}
			}
			stack = stack[:len(stack)-1]
			parent := &stack[len(stack)-1]
			parent.list = append(parent.list, top.el)
		default:
			el := &Node{Type: "mdxJsxTextElement", Name: t.name, Attrs: t.attrs, Children: []*Node{}}
			if t.selfClosing {
				top := &stack[len(stack)-1]
				top.list = append(top.list, el)
				continue
			}
			stack = append(stack, frame{el: el})
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("mdx: expected a closing tag for <%s>", stack[len(stack)-1].el.Name)
	}
	return stack[0].list, nil
}

func isASCIIPunct(c byte) bool {
	return c >= '!' && c <= '/' || c >= ':' && c <= '@' || c >= '[' && c <= '`' || c >= '{' && c <= '~'
}

func (p *inlineParser) run() error {
	src := p.src
	start := p.pos
	flush := func() {
		if p.pos > start {
			p.text(src[start:p.pos])
		}
	}
	for p.pos < len(src) {
		c := src[p.pos]
		switch c {
		case '\n':
			flush()
			// Trailing whitespace before a line ending.
			if p.tail != nil && p.tail.n != nil && p.tail.n.Type == "text" && p.tail.delim == nil && !p.tail.bracket {
				v := p.tail.n.Value
				trimmed := strings.TrimRight(v, " \t")
				if len(v)-len(trimmed) >= 2 && strings.HasSuffix(v, "  ") {
					p.tail.n.Value = trimmed
					p.append(&item{n: &Node{Type: "break"}})
					p.pos++
					start = p.pos
					continue
				}
				p.tail.n.Value = trimmed
			}
			p.text("\n")
			p.pos++
			// The indentation of the next line is not text.
			for p.pos < len(src) && (src[p.pos] == ' ' || src[p.pos] == '\t') {
				p.pos++
			}
			start = p.pos
		case '\\':
			flush()
			if p.pos+1 < len(src) && src[p.pos+1] == '\n' {
				p.append(&item{n: &Node{Type: "break"}})
				p.pos += 2
			} else if p.pos+1 < len(src) && isASCIIPunct(src[p.pos+1]) {
				p.text(src[p.pos+1 : p.pos+2])
				p.pos += 2
			} else {
				p.text("\\")
				p.pos++
			}
			start = p.pos
		case '`':
			flush()
			ok, err := p.codeSpan()
			if err != nil {
				return err
			}
			if !ok {
				n := 0
				for p.pos+n < len(src) && src[p.pos+n] == '`' {
					n++
				}
				p.text(src[p.pos : p.pos+n])
				p.pos += n
			}
			start = p.pos
		case '&':
			flush()
			if m := entityRe.FindString(src[p.pos:]); m != "" && html.UnescapeString(m) != m {
				return fmt.Errorf("mdx: character references (%s) are not supported (the docs content uses none)", m)
			} else {
				p.text("&")
				p.pos++
			}
			start = p.pos
		case '*', '_', '~':
			flush()
			p.delimRun(c)
			start = p.pos
		case '[':
			flush()
			if p.footnoteCall() {
				start = p.pos
				continue
			}
			it := p.append(&item{n: &Node{Type: "text", Value: "["}, bracket: true})
			p.pushBracket(it, false)
			p.pos++
			start = p.pos
		case '!':
			flush()
			if p.pos+1 < len(src) && src[p.pos+1] == '[' {
				it := p.append(&item{n: &Node{Type: "text", Value: "!["}, bracket: true})
				p.pushBracket(it, true)
				p.pos += 2
			} else {
				p.text("!")
				p.pos++
			}
			start = p.pos
		case ']':
			flush()
			if err := p.closeBracket(); err != nil {
				return err
			}
			start = p.pos
		case '<':
			flush()
			t, end, err := parseTag(src, p.pos)
			switch {
			case err == errNotTag:
				p.text("<")
				p.pos++
			case err == errIncomplete:
				return errors.New("mdx: unterminated JSX tag in text")
			case err != nil:
				return err
			default:
				p.append(&item{tag: t})
				p.pos = end
			}
			start = p.pos
		case '{':
			flush()
			e, end, err := parseExpressionBody(src, p.pos+1)
			if err != nil {
				return err
			}
			if tabIndented(src[p.pos+1 : end]) {
				return errors.New("mdx: tabs indenting a multi-line expression are not supported")
			}
			p.append(&item{n: &Node{Type: "mdxTextExpression", Value: exprValue(src[p.pos+1 : end]), Expr: e}})
			p.pos = end + 1
			start = p.pos
		case 'h', 'H', 'w', 'W':
			if p.autolinkLiteral(flush) {
				start = p.pos
				continue
			}
			p.pos++
		default:
			if p.emailStart(c) {
				if p.autolinkLiteral(flush) {
					start = p.pos
					continue
				}
			}
			p.pos++
		}
	}
	flush()
	// Trailing whitespace at the end of the content.
	if p.tail != nil && p.tail.n != nil && p.tail.n.Type == "text" && p.tail.delim == nil && !p.tail.bracket {
		p.tail.n.Value = strings.TrimRight(p.tail.n.Value, " \t")
	}
	return nil
}

// --- code spans ---------------------------------------------------------

func (p *inlineParser) codeSpan() (bool, error) {
	src := p.src
	n := 0
	for p.pos+n < len(src) && src[p.pos+n] == '`' {
		n++
	}
	i := p.pos + n
	for i < len(src) {
		j := strings.IndexByte(src[i:], '`')
		if j < 0 {
			return false, nil
		}
		j += i
		m := 0
		for j+m < len(src) && src[j+m] == '`' {
			m++
		}
		if m == n {
			// mdast keeps line endings (and the indentation after them) in
			// inline code; mdast-util-to-hast turns line endings into spaces.
			// One space or line ending is stripped from each side when both
			// sides have one and the code is not only spaces and line endings.
			content := src[p.pos+n : j]
			if strings.Contains(content, "\n") && strings.ContainsRune(content, '\t') {
				return false, errors.New("mdx: tabs in multi-line code spans are not supported")
			}
			pad := func(c byte) bool { return c == ' ' || c == '\n' }
			if len(content) >= 2 && pad(content[0]) && pad(content[len(content)-1]) && strings.Trim(content, " \n") != "" {
				content = content[1 : len(content)-1]
			}
			if p.inTable {
				content = strings.ReplaceAll(content, `\|`, "|")
			}
			p.append(&item{n: &Node{Type: "inlineCode", Value: content}})
			p.pos = j + m
			return true, nil
		}
		i = j + m
	}
	return false, nil
}

// --- attention --------------------------------------------------------------

// classify is micromark's classifyCharacter: 1 for whitespace (or none),
// 2 for punctuation (Unicode P or S), 0 otherwise.
func classify(r rune, ok bool) int {
	if !ok || unicode.IsSpace(r) || r == '\n' {
		return 1
	}
	if unicode.IsPunct(r) || unicode.IsSymbol(r) {
		return 2
	}
	return 0
}

func (p *inlineParser) runeBefore(i int) (rune, bool) {
	if i <= 0 {
		return 0, false
	}
	r, _ := utf8.DecodeLastRuneInString(p.src[:i])
	return r, true
}

func (p *inlineParser) runeAt(i int) (rune, bool) {
	if i >= len(p.src) {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(p.src[i:])
	return r, true
}

func (p *inlineParser) delimRun(c byte) {
	startPos := p.pos
	n := 0
	for p.pos < len(p.src) && p.src[p.pos] == c {
		n++
		p.pos++
	}
	before := classify(p.runeBefore(startPos))
	after := classify(p.runeAt(p.pos))
	var canOpen, canClose bool
	if c == '~' {
		canOpen = after == 0 || (after == 2 && before != 0)
		canClose = before == 0 || (before == 2 && after != 0)
		if n > 2 {
			p.text(p.src[startPos:p.pos])
			return
		}
	} else {
		open := after == 0 || (after == 2 && before != 0)
		closeF := before == 0 || (before == 2 && after != 0)
		if c == '*' {
			canOpen, canClose = open, closeF
		} else {
			canOpen = open && (!closeF || before == 2)
			canClose = closeF && (!open || after == 2)
		}
	}
	it := p.append(&item{n: &Node{Type: "text", Value: p.src[startPos:p.pos]}})
	d := &delim{it: it, char: c, count: n, origLen: n, canOpen: canOpen, canClose: canClose, prev: p.delims}
	if p.delims != nil {
		p.delims.next = d
	}
	p.delims = d
	it.delim = d
}

func (p *inlineParser) removeDelim(d *delim) {
	if d.prev != nil {
		d.prev.next = d.next
	}
	if d.next != nil {
		d.next.prev = d.prev
	} else {
		p.delims = d.prev
	}
	d.it.delim = nil
}

// processEmphasis is CommonMark's "process emphasis" over delimiters above
// stackBottom, plus GFM strikethrough (runs of equal length).
func (p *inlineParser) processEmphasis(stackBottom *delim) {
	openersBottom := map[string]*delim{}
	closer := p.delims
	for closer != nil && closer.prev != stackBottom {
		closer = closer.prev
	}
	for closer != nil {
		if !closer.canClose {
			closer = closer.next
			continue
		}
		key := fmt.Sprintf("%c%v%d", closer.char, closer.canOpen, closer.origLen%3)
		opener := closer.prev
		found := false
		for opener != nil && opener != stackBottom && opener != openersBottom[key] {
			if opener.char == closer.char && opener.canOpen {
				if closer.char == '~' {
					if opener.count == closer.count {
						found = true
						break
					}
				} else {
					odd := (closer.canOpen || opener.canClose) && closer.origLen%3 != 0 && (opener.origLen+closer.origLen)%3 == 0
					if !odd {
						found = true
						break
					}
				}
			}
			opener = opener.prev
		}
		oldCloser := closer
		if !found {
			openersBottom[key] = oldCloser.prev
			if !oldCloser.canOpen {
				closer = closer.next
				p.removeDelim(oldCloser)
			} else {
				closer = closer.next
			}
			continue
		}
		var use int
		var typ string
		if closer.char == '~' {
			use = closer.count
			typ = "delete"
		} else if closer.count >= 2 && opener.count >= 2 {
			use = 2
			typ = "strong"
		} else {
			use = 1
			typ = "emphasis"
		}
		opener.count -= use
		closer.count -= use
		opener.it.n.Value = opener.it.n.Value[:opener.count]
		closer.it.n.Value = closer.it.n.Value[:closer.count]
		node := &Node{Type: typ}
		// Move items between opener and closer into node.
		var kids []*item
		for it := opener.it.next; it != closer.it; it = it.next {
			kids = append(kids, it)
		}
		// Remove delimiters between opener and closer.
		for d := closer.prev; d != nil && d != opener; {
			prev := d.prev
			p.removeDelim(d)
			d = prev
		}
		emph := &item{n: node}
		node.Children = p.itemsToNodes(kids)
		// Splice: opener.it, emph, closer.it
		if len(kids) > 0 {
			first, last := kids[0], kids[len(kids)-1]
			first.prev, last.next = nil, nil
		}
		opener.it.next = emph
		emph.prev = opener.it
		emph.next = closer.it
		closer.it.prev = emph
		if opener.count == 0 {
			p.remove(opener.it)
			p.removeDelim(opener)
		}
		if closer.count == 0 {
			next := closer.next
			p.remove(closer.it)
			p.removeDelim(closer)
			closer = next
		}
	}
	for p.delims != nil && p.delims != stackBottom {
		p.removeDelim(p.delims)
	}
}

func (p *inlineParser) itemsToNodes(items []*item) []*Node {
	var out []*Node
	for _, it := range items {
		switch {
		case it.tag != nil:
			out = append(out, &Node{Type: "\x00tag", Data: map[string]any{"tag": it.tag}})
		case it.delim != nil || it.bracket:
			out = append(out, &Node{Type: "text", Value: it.n.Value})
		default:
			out = append(out, it.n)
		}
	}
	return out
}

// --- links and images -------------------------------------------------------

func (p *inlineParser) pushBracket(it *item, image bool) {
	b := &bracketEntry{it: it, image: image, active: true, prev: p.brackets, prevDel: p.delims, index: p.pos}
	if p.brackets != nil {
		p.brackets.bracketAfter = true
	}
	p.brackets = b
}

func (p *inlineParser) closeBracket() error {
	p.pos++ // ]
	opener := p.brackets
	if opener == nil {
		p.text("]")
		return nil
	}
	if !opener.active {
		p.text("]")
		p.brackets = opener.prev
		return nil
	}
	// Inline link: (dest "title")
	dest, title, hasTitle, end, ok := p.linkTail(p.pos)
	if !ok {
		if opener.image {
			// fall through to literal
		}
		p.brackets = opener.prev
		p.text("]")
		return nil
	}
	p.pos = end
	node := &Node{Type: "link", URL: dest, Title: title, HasTitle: hasTitle}
	if opener.image {
		node.Type = "image"
	}
	var kids []*item
	for it := opener.it.next; it != nil; it = it.next {
		kids = append(kids, it)
	}
	p.processEmphasis(opener.prevDel)
	kids = kids[:0]
	for it := opener.it.next; it != nil; it = it.next {
		kids = append(kids, it)
	}
	children := p.itemsToNodes(kids)
	for _, k := range kids {
		_ = k
	}
	// Truncate list after opener and replace opener with the node.
	opener.it.next = nil
	p.tail = opener.it
	opener.it.n = node
	opener.it.bracket = false
	if node.Type == "image" {
		node.Alt = altText(children)
	} else {
		node.Children = children
		if node.Children == nil {
			node.Children = []*Node{}
		}
	}
	p.brackets = opener.prev
	if !opener.image {
		for b := p.brackets; b != nil; b = b.prev {
			if !b.image {
				b.active = false
			}
		}
	}
	return nil
}

func altText(nodes []*Node) string {
	var b strings.Builder
	for _, n := range nodes {
		switch n.Type {
		case "text", "inlineCode":
			b.WriteString(n.Value)
		case "image":
			b.WriteString(n.Alt)
		default:
			b.WriteString(altText(n.Children))
		}
	}
	return b.String()
}

// linkTail parses "(destination title)" at i.
func (p *inlineParser) linkTail(i int) (dest, title string, hasTitle bool, end int, ok bool) {
	src := p.src
	if i >= len(src) || src[i] != '(' {
		return
	}
	i++
	skipWS := func() {
		for i < len(src) && (src[i] == ' ' || src[i] == '\t' || src[i] == '\n') {
			i++
		}
	}
	skipWS()
	// Destination.
	if i < len(src) && src[i] == '<' {
		j := i + 1
		for j < len(src) && src[j] != '>' && src[j] != '\n' && src[j] != '<' {
			if src[j] == '\\' && j+1 < len(src) {
				j++
			}
			j++
		}
		if j >= len(src) || src[j] != '>' {
			return
		}
		dest = src[i+1 : j]
		i = j + 1
	} else {
		start := i
		depth := 0
		for i < len(src) {
			c := src[i]
			if c == '\\' && i+1 < len(src) && isASCIIPunct(src[i+1]) {
				i += 2
				continue
			}
			if c == '(' {
				depth++
			} else if c == ')' {
				if depth == 0 {
					break
				}
				depth--
			} else if c <= ' ' {
				break
			}
			i++
		}
		dest = src[start:i]
	}
	beforeTitle := i
	skipWS()
	if i < len(src) && (src[i] == '"' || src[i] == '\'' || src[i] == '(') && i > beforeTitle {
		q := src[i]
		if q == '(' {
			q = ')'
		}
		j := i + 1
		for j < len(src) && src[j] != q {
			if src[j] == '\\' && j+1 < len(src) {
				j++
			}
			j++
		}
		if j >= len(src) {
			return
		}
		title = unescapeString(src[i+1 : j])
		hasTitle = true
		i = j + 1
		skipWS()
	}
	if i >= len(src) || src[i] != ')' {
		return
	}
	dest = unescapeString(dest)
	return dest, title, hasTitle, i + 1, true
}

var escapeRe = regexp.MustCompile(`\\([!-/:-@\[-` + "`" + `{-~])`)

func unescapeString(s string) string {
	s = escapeRe.ReplaceAllString(s, "$1")
	if strings.Contains(s, "&") {
		s = decodeEntities(s)
	}
	return s
}

// --- footnote calls -----------------------------------------------------------

func (p *inlineParser) footnoteCall() bool {
	src := p.src
	if !strings.HasPrefix(src[p.pos:], "[^") {
		return false
	}
	end := strings.IndexByte(src[p.pos+2:], ']')
	if end <= 0 {
		return false
	}
	label := src[p.pos+2 : p.pos+2+end]
	if strings.ContainsAny(label, " \t\n[") {
		return false
	}
	id := normalizeIdentifier(label)
	if !p.footnotes[id] {
		return false
	}
	p.append(&item{n: &Node{Type: "footnoteReference", Identifier: id, Label: label}})
	p.pos += end + 3
	return true
}

func normalizeIdentifier(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.ToLower(strings.ToUpper(s))
}

// --- autolink literals --------------------------------------------------------

func gfmAtext(c byte) bool {
	return c == '+' || c == '-' || c == '.' || c == '_' || isAlnum(c)
}

func isAlnum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
func isAlpha(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }

func (p *inlineParser) emailStart(c byte) bool {
	return gfmAtext(c)
}

func (p *inlineParser) prevByte() (byte, bool) {
	if p.pos == 0 {
		return 0, false
	}
	return p.src[p.pos-1], true
}

func (p *inlineParser) unbalancedLabel() bool {
	return p.brackets != nil
}

// autolinkLiteral tries www., http(s):// and email literals at p.pos.
func (p *inlineParser) autolinkLiteral(flush func()) bool {
	src := p.src
	c := src[p.pos]
	prev, hasPrev := p.prevByte()
	if p.unbalancedLabel() {
		return false
	}
	// Email first (micromark tries it first for every atext character).
	if hasPrev && (prev == '/' || gfmAtext(prev)) {
		// not an email start
	} else if gfmAtext(c) {
		if end, ok := p.email(p.pos); ok {
			flush()
			v := src[p.pos:end]
			p.append(&item{n: &Node{Type: "link", URL: "mailto:" + v, Children: []*Node{{Type: "text", Value: v}}, Data: map[string]any{"literal": "email"}}})
			p.pos = end
			return true
		}
	}
	if (c == 'h' || c == 'H') && (!hasPrev || !isAlpha(prev)) {
		lower := strings.ToLower(src[p.pos:min(p.pos+8, len(src))])
		n := 0
		if strings.HasPrefix(lower, "https://") {
			n = 8
		} else if strings.HasPrefix(lower, "http://") {
			n = 7
		}
		if n > 0 {
			i := p.pos + n
			if i < len(src) {
				r, _ := utf8.DecodeRuneInString(src[i:])
				if !(r < 0x20 || r == 0x7f || unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r)) {
					if dEnd, ok := autolinkDomain(src, i); ok {
						end := autolinkPath(src, dEnd)
						flush()
						v := src[p.pos:end]
						p.append(&item{n: &Node{Type: "link", URL: v, Children: []*Node{{Type: "text", Value: v}}, Data: map[string]any{"literal": "http"}}})
						p.pos = end
						return true
					}
				}
			}
		}
	}
	if (c == 'w' || c == 'W') && (!hasPrev || prev == '(' || prev == '*' || prev == '_' || prev == '[' || prev == ']' || prev == '~' || prev == ' ' || prev == '\t' || prev == '\n') {
		if len(src) >= p.pos+4 && strings.EqualFold(src[p.pos:p.pos+3], "www") && src[p.pos+3] == '.' && p.pos+4 < len(src) {
			if dEnd, ok := autolinkDomain(src, p.pos); ok {
				end := autolinkPath(src, dEnd)
				flush()
				v := src[p.pos:end]
				p.append(&item{n: &Node{Type: "link", URL: "http://" + v, Children: []*Node{{Type: "text", Value: v}}, Data: map[string]any{"literal": "www"}}})
				p.pos = end
				return true
			}
		}
	}
	return false
}

func (p *inlineParser) email(i int) (int, bool) {
	src := p.src
	j := i
	for j < len(src) && gfmAtext(src[j]) {
		j++
	}
	if j >= len(src) || src[j] != '@' {
		return 0, false
	}
	j++
	data, dot := false, false
	for j < len(src) {
		c := src[j]
		if c == '.' {
			// Dot trail: a dot not followed by alnum ends the domain.
			if j+1 < len(src) && isAlnum(src[j+1]) {
				dot = true
				j++
				continue
			}
			break
		}
		if c == '-' || c == '_' || isAlnum(c) {
			data = true
			j++
			continue
		}
		break
	}
	if data && dot && isAlpha(src[j-1]) {
		return j, true
	}
	return 0, false
}

func isTrailChar(c byte) bool {
	switch c {
	case '!', '"', '\'', ')', '*', ',', '.', ':', ';', '?', '_', '~':
		return true
	}
	return false
}

// trail reports whether src[i:] is trailing punctuation (micromark's trail
// partial): punctuation, character references or a closing bracket,
// followed by the end, whitespace or "<".
func trail(src string, i int) bool {
	for i < len(src) {
		c := src[i]
		switch {
		case isTrailChar(c):
			i++
		case c == '&':
			j := i + 1
			if j >= len(src) || !isAlpha(src[j]) {
				return false
			}
			for j < len(src) && isAlpha(src[j]) {
				j++
			}
			if j >= len(src) || src[j] != ';' {
				return false
			}
			i = j + 1
		case c == ']':
			i++
			if i >= len(src) || src[i] == '(' || src[i] == '[' || isWS(src[i]) {
				return true
			}
		case c == '<' || isWS(c):
			return true
		default:
			r, _ := utf8.DecodeRuneInString(src[i:])
			return unicode.IsSpace(r)
		}
	}
	return true
}

func autolinkDomain(src string, i int) (int, bool) {
	var underscoreLast, underscoreLastLast, seen bool
	for i < len(src) {
		c := src[i]
		if c == '.' || c == '_' {
			if trail(src, i) {
				break
			}
			if c == '_' {
				underscoreLast = true
			} else {
				underscoreLastLast = underscoreLast
				underscoreLast = false
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(src[i:])
		if unicode.IsSpace(r) || (c != '-' && (unicode.IsPunct(r) || unicode.IsSymbol(r))) {
			break
		}
		seen = true
		i += size
	}
	if underscoreLastLast || underscoreLast || !seen {
		return 0, false
	}
	return i, true
}

func autolinkPath(src string, i int) int {
	open, closeN := 0, 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '(':
			open++
			i++
		case c == ')' && closeN < open:
			closeN++
			i++
		case c == '!' || c == '"' || c == '&' || c == '\'' || c == ')' || c == '*' || c == ',' || c == '.' || c == ':' || c == ';' || c == '<' || c == '?' || c == ']' || c == '_' || c == '~':
			if trail(src, i) {
				return i
			}
			if c == ')' {
				closeN++
			}
			i++
		default:
			r, size := utf8.DecodeRuneInString(src[i:])
			if unicode.IsSpace(r) {
				return i
			}
			i += size
		}
	}
	return i
}

var (
	urlFindRe   = regexp.MustCompile(`(?i)(https?://|www\.)([-.\w]+)([^ \t\r\n]*)`)
	emailFindRe = regexp.MustCompile(`([-.\w+]+)@([-\w]+(?:\.[-\w]+)+)`)
)

// checkAutolinkTransform reports text that mdast-util-gfm-autolink-literal's
// transform would link (literals micromark's tokenizer did not link). The
// docs never reach it: such text is reported as unsupported rather than
// guessed.
func checkAutolinkTransform(n *Node) error {
	if n.Type == "link" {
		return nil
	}
	if n.Type == "text" {
		for _, m := range urlFindRe.FindAllStringIndex(n.Value, -1) {
			if m[0] == 0 || isTransformPrev(n.Value, m[0]) {
				return fmt.Errorf("mdx: unsupported autolink literal in %q", n.Value)
			}
		}
		for _, m := range emailFindRe.FindAllStringIndex(n.Value, -1) {
			if (m[0] == 0 || isTransformPrev(n.Value, m[0])) && (m[0] == 0 || n.Value[m[0]-1] != '/') {
				label := n.Value[strings.IndexByte(n.Value[m[0]:m[1]], '@')+m[0]+1 : m[1]]
				if last := label[len(label)-1]; last == '-' || last == '_' || (last >= '0' && last <= '9') {
					continue
				}
				return fmt.Errorf("mdx: unsupported email autolink literal in %q", n.Value)
			}
		}
	}
	for _, c := range n.Children {
		if err := checkAutolinkTransform(c); err != nil {
			return err
		}
	}
	return nil
}

func isTransformPrev(s string, i int) bool {
	r, _ := utf8.DecodeLastRuneInString(s[:i])
	return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r)
}
