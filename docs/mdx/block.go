package mdx

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// The block parser follows the CommonMark algorithm (as micromark
// implements it), with MDX's changes: no indented code (so indentation
// never prevents a block from starting), no HTML, and JSX tags, expressions
// and ESM as flow constructs. JSX flow tags are leaves here; they are
// nested into elements when the tree is built (build.go), like
// mdast-util-mdx-jsx does.

type blockKind int

const (
	kDocument blockKind = iota
	kBlockquote
	kList
	kItem
	kFootnote
	kParagraph
	kHeading
	kThematic
	kFence
	kTable
	kJSX  // one line of JSX flow tags and flow expressions
	kESM  // import/export statements
	kPend // a multi-line JSX tag or expression being read
)

type block struct {
	kind     blockKind
	parent   *block
	children []*block
	open     bool
	line     int // line number where the block starts

	lastLineBlank bool
	// lastLine is the last line holding content of the block (or of one of
	// its descendants): blank lines do not count, except inside fenced code
	// and multi-line JSX, which keep them. A blank line between two children
	// is a gap between the first's lastLine and the second's line: that is
	// what makes a list or list item spread (mdast-util-from-markdown's
	// prepareList, where blank lines that end a nested container belong to
	// the outer one).
	lastLine int

	// list, item
	ordered bool
	bullet  byte // '-', '*', '+' or the ordered delimiter '.' / ')'
	start   int
	// size is the indentation a continuation line needs: the marker's own
	// indentation plus the marker and the whitespace after it (micromark's
	// containerState.size; MDX has no indented code, so all of that
	// whitespace counts).
	size int
	// item: the first line held only the marker (blank item start).
	startsBlank bool
	// paragraph, heading, ESM, pending text
	lines []string
	// paragraph: whether each line was a lazy continuation line
	lazy []bool
	// heading
	level int
	// fence
	fenceChar   byte
	fenceLen    int
	fenceIndent int
	info        string
	closed      bool
	// table
	align     []string
	rows      [][]string
	delimLine int
	// jsx line
	items []flowItem
	// footnote
	label string
}

// flowItem is one JSX tag or expression of a flow line.
type flowItem struct {
	tag  *jsxTag
	expr *Node // mdxFlowExpression
}

type parser struct {
	lines []string
	doc   *block
	tip   *block

	// current line
	ln              string
	lineNo          int
	offset          int // byte offset in ln
	column          int
	nextNonspace    int
	nextNonspaceCol int
	indent          int
	blank           bool
	partialTab      bool
	allClosed       bool
	lastMatched     *block
	oldTip          *block
	err             error // set by continueBlock
	// contentStart is the offset after the container prefixes, before the
	// line's indentation (kept in paragraph continuation lines).
	contentStart int
}

// Parse parses MDX source (without frontmatter) into an mdast root.
func Parse(src string) (*Node, error) {
	if strings.ContainsRune(src, '\r') {
		return nil, fmt.Errorf("mdx: carriage returns are not supported")
	}
	lines := strings.Split(src, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	p := &parser{lines: lines}
	p.doc = &block{kind: kDocument, open: true}
	p.tip = p.doc
	for i, l := range lines {
		p.lineNo = i + 1
		if err := p.incorporate(l); err != nil {
			return nil, fmt.Errorf("mdx: line %d: %w", p.lineNo, err)
		}
	}
	for p.tip != nil {
		if err := p.finalize(p.tip); err != nil {
			return nil, err
		}
	}
	return build(p.doc)
}

// --- line cursor ------------------------------------------------------------

func (p *parser) findNextNonspace() {
	i, cols := p.offset, p.column
	for i < len(p.ln) {
		c := p.ln[i]
		if c == ' ' {
			i++
			cols++
		} else if c == '\t' {
			i++
			cols += 4 - cols%4
		} else {
			break
		}
	}
	p.blank = i >= len(p.ln)
	p.nextNonspace = i
	p.nextNonspaceCol = cols
	p.indent = cols - p.column
}

func (p *parser) advanceNextNonspace() {
	p.offset = p.nextNonspace
	p.column = p.nextNonspaceCol
	p.partialTab = false
}

// advanceOffset advances by count bytes, or count columns when columns.
func (p *parser) advanceOffset(count int, columns bool) {
	for count > 0 && p.offset < len(p.ln) {
		c := p.ln[p.offset]
		if c == '\t' {
			charsToTab := 4 - p.column%4
			if columns {
				p.partialTab = charsToTab > count
				adv := min(charsToTab, count)
				p.column += adv
				if !p.partialTab {
					p.offset++
				}
				count -= adv
			} else {
				p.partialTab = false
				p.column += charsToTab
				p.offset++
				count--
			}
		} else {
			p.partialTab = false
			p.offset++
			p.column++
			count--
		}
	}
}

// rest returns the remainder of the line, expanding a partially consumed
// tab into spaces.
func (p *parser) rest() string {
	if p.partialTab {
		p.offset++
		charsToTab := 4 - p.column%4
		p.partialTab = false
		return strings.Repeat(" ", charsToTab) + p.ln[p.offset:]
	}
	return p.ln[p.offset:]
}

// --- main loop --------------------------------------------------------------

func acceptsLines(k blockKind) bool {
	return k == kParagraph || k == kFence || k == kTable || k == kESM || k == kPend
}

func canContain(parent blockKind, child blockKind) bool {
	switch parent {
	case kDocument, kBlockquote, kItem, kFootnote:
		return child != kItem
	case kList:
		return child == kItem
	}
	return false
}

func (p *parser) incorporate(ln string) error {
	p.ln = ln
	p.offset, p.column = 0, 0
	p.contentStart = 0
	p.blank, p.partialTab = false, false
	p.oldTip = p.tip
	container := p.doc
	allMatched := true

	for {
		last := lastChild(container)
		if last == nil || !last.open {
			break
		}
		container = last
		p.findNextNonspace()
		res := p.continueBlock(container)
		if p.err != nil {
			return p.err
		}
		switch res {
		case 0:
		case 1:
			allMatched = false
		case 2:
			// A closing fence: the line belongs to the fence.
			p.touch(container)
			return nil
		}
		if !allMatched {
			container = container.parent
			break
		}
	}
	p.allClosed = container == p.oldTip
	p.lastMatched = container

	matchedLeaf := container.kind != kParagraph && acceptsLines(container.kind)
	for !matchedLeaf {
		p.findNextNonspace()
		res, err := p.tryStarts(container)
		if err != nil {
			return err
		}
		if res == 0 {
			p.contentStart = p.offset
			p.advanceNextNonspace()
			break
		}
		container = p.tip
		if res == 2 {
			matchedLeaf = true
		}
	}

	if !p.allClosed && !p.blank && p.tip.kind == kParagraph {
		// Lazy paragraph continuation.
		p.tip.lines = append(p.tip.lines, p.rest())
		p.tip.lazy = append(p.tip.lazy, true)
		p.touch(p.tip)
		return nil
	}
	defer func() {
		// The line's content belongs to the tip and its ancestors.
		if !p.blank {
			p.touch(p.tip)
			return
		}
		// A blank line inside a block quote (a line with just ">") belongs
		// to the quote; a blank line inside fenced code or a multi-line
		// JSX tag belongs to it.
		for c := p.tip; c != nil; c = c.parent {
			if c.kind == kBlockquote || (c.kind == kFence || c.kind == kPend) && c.line < p.lineNo {
				p.touch(c)
				return
			}
		}
	}()
	if err := p.closeUnmatched(); err != nil {
		return err
	}
	if p.blank && lastChild(container) != nil {
		lastChild(container).lastLineBlank = true
	}
	k := container.kind
	lastLineBlank := p.blank && !(k == kBlockquote || (k == kFence) || (k == kItem && len(container.children) == 0 && container.line == p.lineNo))
	for c := container; c != nil; c = c.parent {
		c.lastLineBlank = lastLineBlank
	}
	switch {
	case acceptsLines(k):
		return p.addLine(container)
	case p.offset < len(p.ln) && !p.blank:
		para := p.addChild(kParagraph)
		p.advanceNextNonspace()
		para.lines = append(para.lines, p.rest())
		para.lazy = append(para.lazy, false)
	case p.blank && k == kTable:
	}
	return nil
}

// touch records that line p.lineNo holds content of b and its ancestors.
func (p *parser) touch(b *block) {
	for ; b != nil; b = b.parent {
		b.lastLine = p.lineNo
	}
}

func lastChild(b *block) *block {
	if len(b.children) == 0 {
		return nil
	}
	return b.children[len(b.children)-1]
}

func (p *parser) closeUnmatched() error {
	if !p.allClosed {
		for p.oldTip != p.lastMatched {
			parent := p.oldTip.parent
			if err := p.finalize(p.oldTip); err != nil {
				return err
			}
			p.oldTip = parent
		}
		p.allClosed = true
	}
	return nil
}

func (p *parser) addChild(kind blockKind) *block {
	for !canContain(p.tip.kind, kind) {
		_ = p.finalize(p.tip)
	}
	b := &block{kind: kind, parent: p.tip, open: true, line: p.lineNo, lastLine: p.lineNo}
	p.tip.children = append(p.tip.children, b)
	p.tip = b
	return b
}

// continueBlock: 0 = matched, 1 = not matched, 2 = line consumed.
func (p *parser) continueBlock(b *block) int {
	switch b.kind {
	case kBlockquote:
		if p.ln[p.nextNonspace:] != "" && p.ln[p.nextNonspace] == '>' {
			p.advanceNextNonspace()
			p.advanceOffset(1, false)
			if p.offset < len(p.ln) && (p.ln[p.offset] == ' ' || p.ln[p.offset] == '\t') {
				p.advanceOffset(1, true)
			}
			return 0
		}
		return 1
	case kItem:
		if p.blank {
			if len(b.children) == 0 {
				return 1 // a blank item cannot start with a blank line
			}
			p.advanceNextNonspace()
			return 0
		}
		if p.indent >= b.size {
			if strings.ContainsRune(p.ln[p.offset:p.nextNonspace], '\t') {
				p.err = fmt.Errorf("tabs in list item indentation are not supported")
			}
			p.advanceOffset(b.size, true)
			return 0
		}
		return 1
	case kFootnote:
		if p.blank {
			p.advanceNextNonspace()
			return 0
		}
		if p.indent >= 4 {
			p.advanceOffset(4, true)
			return 0
		}
		return 1
	case kList, kDocument:
		return 0
	case kParagraph:
		if p.blank {
			return 1
		}
		return 0
	case kHeading, kThematic, kJSX:
		return 1
	case kFence:
		if !p.blank {
			c := p.ln[p.nextNonspace]
			if c == b.fenceChar {
				n := 0
				for i := p.nextNonspace; i < len(p.ln) && p.ln[i] == c; i++ {
					n++
				}
				if n >= b.fenceLen && strings.TrimRight(p.ln[p.nextNonspace+n:], " \t") == "" {
					b.closed = true
					_ = p.finalize(b)
					return 2
				}
			}
		}
		// Skip the fence's indentation.
		i := b.fenceIndent
		for i > 0 && p.offset < len(p.ln) && (p.ln[p.offset] == ' ' || p.ln[p.offset] == '\t') {
			p.advanceOffset(1, true)
			i--
		}
		return 0
	case kTable:
		// A body row is any line that does not start another construct
		// (micromark tries the table last); lazy lines never get here.
		if p.blank || p.startsConstruct(b) {
			return 1
		}
		return 0
	case kESM:
		if p.blank {
			return 1
		}
		return 0
	case kPend:
		return 0
	}
	return 1
}

func (p *parser) addLine(b *block) error {
	// The opening line of a fence, pending tag or table (and a table's
	// delimiter row) is not content.
	if (b.kind == kFence || b.kind == kPend) && p.lineNo == b.line || b.kind == kTable && p.lineNo == b.delimLine {
		return nil
	}
	switch b.kind {
	case kParagraph:
		if p.partialTab || strings.ContainsRune(p.ln[p.contentStart:p.offset], '\t') {
			return fmt.Errorf("tabs indenting a paragraph line are not supported")
		}
		b.lines = append(b.lines, p.ln[p.contentStart:])
		b.lazy = append(b.lazy, false)
	case kFence:
		b.lines = append(b.lines, p.rest())
	case kESM:
		b.lines = append(b.lines, p.rest())
	case kTable:
		if p.blank {
			return nil
		}
		b.rows = append(b.rows, splitRow(p.rest()))
	case kPend:
		b.lines = append(b.lines, p.rest())
		return p.tryCompletePending(b)
	}
	return nil
}

// --- block starts -----------------------------------------------------------

var (
	reATX        = regexp.MustCompile(`^#{1,6}(?:[ \t]+|$)`)
	reThematic   = regexp.MustCompile(`^(?:\*[ \t]*){3,}$|^(?:_[ \t]*){3,}$|^(?:-[ \t]*){3,}$`)
	reFenceStart = regexp.MustCompile("^`{3,}|^~{3,}")
	reBullet     = regexp.MustCompile(`^[*+-]`)
	reOrdered    = regexp.MustCompile(`^(\d{1,9})([.)])`)
	reSetext     = regexp.MustCompile(`^(?:=+|-+)[ \t]*$`)
	reFootDef    = regexp.MustCompile(`^\[\^([^\]\s]+)\]:`)
	reDefinition = regexp.MustCompile(`^\[[^\]]+\]:`)
)

// tryStarts tries each block start: 0 = none, 1 = container started,
// 2 = leaf started.
func (p *parser) tryStarts(container *block) (int, error) {
	rest := p.ln[p.nextNonspace:]
	if p.blank {
		return 0, nil
	}
	interrupting := p.tip.kind == kParagraph && container == p.tip

	// Block quote.
	if rest[0] == '>' {
		p.advanceNextNonspace()
		p.advanceOffset(1, false)
		if p.offset < len(p.ln) && (p.ln[p.offset] == ' ' || p.ln[p.offset] == '\t') {
			p.advanceOffset(1, true)
		}
		if err := p.closeUnmatched(); err != nil {
			return 0, err
		}
		p.addChild(kBlockquote)
		return 1, nil
	}
	// ATX heading.
	if m := reATX.FindString(rest); m != "" {
		p.advanceNextNonspace()
		if err := p.closeUnmatched(); err != nil {
			return 0, err
		}
		h := p.addChild(kHeading)
		h.level = strings.Count(strings.TrimRight(m, " \t"), "#")
		text := rest[len(m):]
		// Remove the closing sequence.
		t := strings.TrimRight(text, " \t")
		if j := strings.LastIndexFunc(t, func(r rune) bool { return r != '#' }); j < 0 {
			t = ""
		} else if j < len(t)-1 && (t[j] == ' ' || t[j] == '\t') {
			t = strings.TrimRight(t[:j+1], " \t")
		}
		h.lines = []string{t}
		p.offset = len(p.ln)
		return 2, nil
	}
	// Fenced code.
	if m := reFenceStart.FindString(rest); m != "" && !(m[0] == '`' && strings.ContainsRune(rest[len(m):], '`')) {
		indent := p.indent
		p.advanceNextNonspace()
		p.advanceOffset(len(m), false)
		if err := p.closeUnmatched(); err != nil {
			return 0, err
		}
		f := p.addChild(kFence)
		f.fenceChar, f.fenceLen, f.fenceIndent = m[0], len(m), indent
		f.info = strings.TrimSpace(rest[len(m):])
		p.offset = len(p.ln)
		return 2, nil
	}
	// Setext underline: not supported (the docs use none).
	if interrupting && reSetext.MatchString(rest) && !reThematic.MatchString(rest) {
		return 0, fmt.Errorf("setext headings are not supported")
	}
	if interrupting && reSetext.MatchString(rest) && strings.Trim(rest, "- \t") == "" && strings.Count(rest, "-") >= 1 {
		return 0, fmt.Errorf("setext headings are not supported")
	}
	// Thematic break.
	if reThematic.MatchString(rest) {
		p.advanceNextNonspace()
		if err := p.closeUnmatched(); err != nil {
			return 0, err
		}
		p.addChild(kThematic)
		p.offset = len(p.ln)
		return 2, nil
	}
	// GFM footnote definition.
	if m := reFootDef.FindStringSubmatch(rest); m != nil && !interrupting {
		p.advanceNextNonspace()
		p.advanceOffset(len(m[0]), false)
		if err := p.closeUnmatched(); err != nil {
			return 0, err
		}
		f := p.addChild(kFootnote)
		f.label = m[1]
		// Skip whitespace after the colon.
		p.findNextNonspace()
		p.advanceNextNonspace()
		return 1, nil
	}
	// List item.
	if ok, err := p.tryListItem(container, interrupting); ok || err != nil {
		return 1, err
	}
	// MDX: JSX flow tags and expressions.
	if rest[0] == '<' || rest[0] == '{' {
		ok, err := p.tryFlowMDX(interrupting)
		if ok || err != nil {
			return 2, err
		}
	}
	// MDX: ESM.
	if !interrupting && container == p.doc && p.indent == 0 && (strings.HasPrefix(rest, "import ") || strings.HasPrefix(rest, "export ")) {
		if err := p.closeUnmatched(); err != nil {
			return 0, err
		}
		p.addChild(kESM)
		return 2, nil
	}
	// GFM table: a delimiter row turns the paragraph's last line into the
	// table head.
	if interrupting {
		if ok, err := p.tryTable(container); ok || err != nil {
			return 2, err
		}
	}
	if !interrupting && reDefinition.MatchString(rest) {
		return 0, fmt.Errorf("link reference definitions are not supported")
	}
	return 0, nil
}

func (p *parser) tryListItem(container *block, interrupting bool) (bool, error) {
	rest := p.ln[p.nextNonspace:]
	var ordered bool
	var bullet byte
	var start int
	var markerLen int
	if m := reBullet.FindString(rest); m != "" {
		bullet = m[0]
		markerLen = 1
	} else if m := reOrdered.FindStringSubmatch(rest); m != nil && (container.kind != kParagraph || m[1] == "1") {
		ordered = true
		start, _ = strconv.Atoi(m[1])
		bullet = m[2][0]
		markerLen = len(m[0])
		if interrupting && start != 1 {
			return false, nil
		}
	} else {
		return false, nil
	}
	after := rest[markerLen:]
	if after != "" && after[0] != ' ' && after[0] != '\t' {
		return false, nil
	}
	if interrupting && strings.TrimSpace(after) == "" {
		return false, nil
	}
	if !ordered && (bullet == '*' || bullet == '-') && reThematic.MatchString(rest) {
		return false, nil
	}
	if strings.ContainsRune(p.ln[p.offset:p.nextNonspace], '\t') || strings.HasPrefix(after, "\t") || strings.HasPrefix(strings.TrimLeft(after, " "), "\t") {
		return false, fmt.Errorf("tabs around list item markers are not supported")
	}
	markerOffset := p.indent
	p.advanceNextNonspace()
	p.advanceOffset(markerLen, true)
	spacesStartOffset := p.offset
	for p.offset < len(p.ln) && p.ln[p.offset] == ' ' {
		p.advanceOffset(1, true)
	}
	blankItem := p.offset >= len(p.ln)
	spaces := p.offset - spacesStartOffset
	if blankItem {
		// micromark's onBlank: the content starts one column after the
		// marker.
		p.offset = spacesStartOffset
		p.column = p.column - spaces
		spaces = 1
	}
	if err := p.closeUnmatched(); err != nil {
		return false, err
	}
	// A new list unless the tip is a list of the same type.
	if p.tip.kind != kList || p.tip.ordered != ordered || p.tip.bullet != bullet {
		l := p.addChild(kList)
		l.ordered, l.bullet, l.start = ordered, bullet, start
	}
	it := p.addChild(kItem)
	it.ordered, it.bullet, it.start = ordered, bullet, start
	it.size = markerOffset + markerLen + spaces
	it.startsBlank = blankItem
	return true, nil
}

// --- MDX flow ---------------------------------------------------------------

// tryFlowMDX handles a line starting with "<" or "{": when it holds only
// complete JSX tags and expressions (and whitespace), it is a JSX flow
// line; when a tag or expression continues on the next lines, a pending
// block collects them.
func (p *parser) tryFlowMDX(interrupting bool) (bool, error) {
	text := p.ln[p.nextNonspace:]
	items, complete, ok, err := scanFlowLine(text, false)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if !complete && !p.allClosed && p.oldTip.kind == kParagraph {
		// micromark: the tag's next line is lazy too, which is an error
		// ("Unexpected lazy line in container").
		return false, fmt.Errorf("a multi-line JSX tag or expression cannot start on a lazy line")
	}
	if err := p.closeUnmatched(); err != nil {
		return false, err
	}
	if !complete {
		pend := p.addChild(kPend)
		pend.lines = []string{text}
		p.offset = len(p.ln)
		return true, nil
	}
	j := p.addChild(kJSX)
	j.items = items
	j.open = false
	p.tip = j.parent
	p.touch(j)
	p.offset = len(p.ln)
	return true, nil
}

func (p *parser) tryCompletePending(b *block) error {
	text := strings.Join(b.lines, "\n")
	items, complete, ok, err := scanFlowLine(text, true)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("multi-line JSX tag or expression followed by text is not supported")
	}
	if complete {
		b.kind = kJSX
		b.items = items
		b.lines = nil
		return p.finalize(b)
	}
	return nil
}

// --- tables -----------------------------------------------------------------

// tryTable is called on a non-lazy paragraph continuation line: when the
// line is a delimiter row with as many cells as the paragraph's last line
// (which may itself have been lazy, as in micromark), that line becomes the
// head of a table and the paragraph keeps the lines before it.
func (p *parser) tryTable(para *block) (bool, error) {
	align, ok := parseDelimiterRow(p.ln[p.nextNonspace:])
	if !ok {
		return false, nil
	}
	n := len(para.lines)
	head := strings.TrimSpace(para.lines[n-1])
	if strings.Trim(head, "| \t") == "" && strings.Count(head, "|") <= 1 {
		return false, nil // micromark needs more than a lone pipe
	}
	if strings.ContainsRune(head, '\t') {
		return false, fmt.Errorf("tabs in table rows are not supported")
	}
	cells := splitRow(head)
	if len(cells) != len(align) {
		return false, nil
	}
	if err := p.closeUnmatched(); err != nil {
		return false, err
	}
	var t *block
	if n == 1 {
		t = para
		t.kind = kTable
		t.lines, t.lazy = nil, nil
	} else {
		para.lines, para.lazy = para.lines[:n-1], para.lazy[:n-1]
		para.lastLine = p.lineNo - 2
		if err := p.finalize(para); err != nil {
			return false, err
		}
		t = p.addChild(kTable)
		t.line = p.lineNo - 1
	}
	t.delimLine = p.lineNo
	t.align = align
	t.rows = [][]string{cells}
	p.offset = len(p.ln)
	return true, nil
}

// startsConstruct reports whether the line (at nextNonspace) starts a
// construct other than a paragraph, which ends a table body: a container
// (block quote, list item, footnote definition: no interrupt rules apply)
// or a flow construct (heading, fence, thematic break, JSX, expression,
// ESM).
func (p *parser) startsConstruct(table *block) bool {
	rest := p.ln[p.nextNonspace:]
	switch {
	case rest[0] == '>', reATX.MatchString(rest), reThematic.MatchString(rest), reFootDef.MatchString(rest):
		return true
	}
	if m := reFenceStart.FindString(rest); m != "" && !(m[0] == '`' && strings.ContainsRune(rest[len(m):], '`')) {
		return true
	}
	marker := ""
	if m := reBullet.FindString(rest); m != "" {
		marker = m
	} else if m := reOrdered.FindString(rest); m != "" {
		marker = m
	}
	if marker != "" && (len(rest) == len(marker) || rest[len(marker)] == ' ' || rest[len(marker)] == '\t') {
		return true
	}
	if rest[0] == '<' || rest[0] == '{' {
		if _, _, ok, err := scanFlowLine(rest, false); ok || err != nil {
			return true
		}
	}
	if table.parent == p.doc && p.indent == 0 && (strings.HasPrefix(rest, "import ") || strings.HasPrefix(rest, "export ")) {
		return true
	}
	return false
}

func parseDelimiterRow(s string) ([]string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, false
	}
	if !strings.ContainsAny(s, "|:") {
		return nil, false
	}
	s = strings.TrimPrefix(s, "|")
	s = strings.TrimSuffix(s, "|")
	var align []string
	for _, c := range strings.Split(s, "|") {
		c = strings.TrimSpace(c)
		if c == "" {
			return nil, false
		}
		left := strings.HasPrefix(c, ":")
		right := strings.HasSuffix(c, ":")
		core := strings.Trim(c, ":")
		if core == "" || strings.Trim(core, "-") != "" {
			return nil, false
		}
		switch {
		case left && right:
			align = append(align, "center")
		case left:
			align = append(align, "left")
		case right:
			align = append(align, "right")
		default:
			align = append(align, "")
		}
	}
	return align, true
}

// splitRow splits a table row into raw cell contents on unescaped pipes.
func splitRow(s string) []string {
	s = strings.TrimSpace(s)
	if s == "|" {
		return []string{""}
	}
	var cells []string
	var cur strings.Builder
	started := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) && (s[i+1] == '|' || s[i+1] == '\\') {
			cur.WriteByte(c)
			cur.WriteByte(s[i+1])
			i++
			started = true
			continue
		}
		if c == '|' {
			if i == 0 {
				continue
			}
			cells = append(cells, strings.TrimSpace(cur.String()))
			cur.Reset()
			started = false
			continue
		}
		cur.WriteByte(c)
		started = true
	}
	if started && strings.TrimSpace(cur.String()) != "" {
		cells = append(cells, strings.TrimSpace(cur.String()))
	}
	return cells
}

// --- finalize ---------------------------------------------------------------

func (p *parser) finalize(b *block) error {
	if !b.open {
		if p.tip == b {
			p.tip = b.parent
		}
		return nil
	}
	b.open = false
	switch b.kind {
	case kTable:
		// Drop the delimiter row marker line handled at start.
	case kPend:
		return fmt.Errorf("unterminated JSX tag or expression starting at line %d", b.line)
	}
	p.tip = b.parent
	return nil
}
