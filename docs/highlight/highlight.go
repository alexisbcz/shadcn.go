// Package highlight highlights code blocks exactly like the shadcn/ui docs:
// rehype-pretty-code 0.14.1 with shiki 3.23.0, the themes vesper (dark) and
// github-light-default (light), and the docs' own shiki transformer
// (ui/apps/v4/lib/highlight-code.ts, used by source.config.ts).
//
// Shiki itself runs unmodified in goja with its JavaScript regex engine (no
// WASM); see engine.go for the goja gaps the prelude fills. The markup that
// rehype-pretty-code builds around shiki's output (figure, title, data
// attributes, line numbers, highlighted lines and characters) is ported to
// Go in this file and chars.go.
package highlight

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"shadcn.go/docs/hast"
)

// Themes are the docs' themes, in rehype-pretty-code's option order.
var Themes = []string{"vesper", "github-light-default"}

// themeNames is getThemeNames(theme) joined with spaces (data-theme).
const themeNames = "vesper github-light-default"

// Highlighter highlights code. It is safe for concurrent use: it keeps a
// pool of JavaScript runtimes (goja runtimes are single-threaded) and
// memoizes shiki's output per (code, language).
type Highlighter struct {
	idle    chan *engine
	mu      sync.Mutex
	created int
	max     int
	cache   sync.Map        // [github NUL] lang NUL code -> *hast.Element (pristine)
	loaded  map[string]bool // languages (ids and aliases) the bundle loads
}

// New creates a Highlighter that runs up to min(GOMAXPROCS, 4) shiki
// runtimes in parallel. Creating a runtime takes about 0.2 s; grammars
// compile their regexes lazily on first use in each runtime.
func New() (*Highlighter, error) {
	return NewWithRuntimes(min(runtime.GOMAXPROCS(0), 4))
}

// NewWithRuntimes creates a Highlighter with at most n runtimes.
func NewWithRuntimes(n int) (*Highlighter, error) {
	if n < 1 {
		n = 1
	}
	e, err := newEngine()
	if err != nil {
		return nil, err
	}
	h := &Highlighter{idle: make(chan *engine, n), max: n, created: 1, loaded: e.langs}
	h.idle <- e
	return h, nil
}

func (h *Highlighter) acquire() (*engine, error) {
	select {
	case e := <-h.idle:
		return e, nil
	default:
	}
	h.mu.Lock()
	if h.created < h.max {
		h.created++
		h.mu.Unlock()
		e, err := newEngine()
		if err != nil {
			h.mu.Lock()
			h.created--
			h.mu.Unlock()
			return nil, err
		}
		return e, nil
	}
	h.mu.Unlock()
	return <-h.idle, nil
}

// codeToHast returns shiki's <pre> for code (a copy the caller may modify),
// with rehype-pretty-code's options or, with github, lib/highlight-code.ts's.
func (h *Highlighter) codeToHast(code, lang string, github bool) (*hast.Element, error) {
	key := lang + "\x00" + code
	if github {
		key = "github\x00" + key
	}
	if v, ok := h.cache.Load(key); ok {
		return v.(*hast.Element).Clone(), nil
	}
	e, err := h.acquire()
	if err != nil {
		return nil, err
	}
	pre, err := e.codeToHast(code, lang, github)
	h.idle <- e
	if err != nil {
		return nil, err
	}
	h.cache.Store(key, pre.Clone())
	return pre, nil
}

var (
	defaultOnce sync.Once
	defaultH    *Highlighter
	defaultErr  error
)

// Default returns a process-wide Highlighter.
func Default() (*Highlighter, error) {
	defaultOnce.Do(func() { defaultH, defaultErr = New() })
	return defaultH, defaultErr
}

// CodeElement builds the <pre><code> element mdast-util-to-hast produces for
// a fenced code block (the input of rehype-pretty-code): the code gets a
// trailing newline, the language becomes class="language-<lang>" and the
// meta string is kept as node data.
func CodeElement(code, lang, meta string) *hast.Element {
	c := &hast.Element{Tag: "code"}
	if lang != "" {
		c.Set("class", "language-"+lang)
	}
	if meta != "" {
		c.Data = map[string]any{"meta": meta}
	}
	value := ""
	if code != "" {
		value = code + "\n"
	}
	c.Children = []hast.Node{&hast.Text{Value: value}}
	return &hast.Element{Tag: "pre", Children: []hast.Node{c}}
}

// Block highlights one fenced code block and returns what rehype-pretty-code
// turns it into: a <figure data-rehype-pretty-code-figure> (or the
// unchanged <pre> when the block has no language).
func (h *Highlighter) Block(code, lang, meta string) (*hast.Element, error) {
	pre := CodeElement(code, lang, meta)
	if err := h.Apply(pre); err != nil {
		return nil, err
	}
	return pre, nil
}

// IsBlockCode reports whether el is a <pre> with a single <code> child.
func IsBlockCode(el *hast.Element) bool {
	if el.Tag != "pre" || len(el.Children) != 1 {
		return false
	}
	c, ok := el.Children[0].(*hast.Element)
	return ok && c.Tag == "code"
}

// specialLanguages are the names shiki handles without a grammar.
var specialLanguages = map[string]bool{"text": true, "txt": true, "plain": true, "plaintext": true, "ansi": true}

var (
	titleRe   = regexp.MustCompile(`title="([^"]*)"`)
	captionRe = regexp.MustCompile(`caption="([^"]*)"`)
)

// Apply rewrites a <pre><code> element in place, like rehype-pretty-code's
// visitor for block code.
func (h *Highlighter) Apply(pre *hast.Element) error {
	if !IsBlockCode(pre) {
		return nil
	}
	codeEl := pre.Children[0].(*hast.Element)
	meta, _ := codeEl.Data["meta"].(string)
	if v, ok := codeEl.Get("metastring"); ok && codeEl.Data["meta"] == nil {
		meta = v
	}
	// parseBlockMetaString
	var title, caption string
	var hasTitle, hasCaption bool
	if m := titleRe.FindStringSubmatch(meta); m != nil {
		title, hasTitle = m[1], true
		meta = strings.Replace(meta, m[0], "", 1)
	}
	if m := captionRe.FindStringSubmatch(meta); m != nil {
		caption, hasCaption = m[1], true
		meta = strings.Replace(meta, m[0], "", 1)
	}
	lang := ""
	if cls := codeEl.Classes(); len(cls) > 0 && strings.HasPrefix(cls[0], "language-") {
		lang = strings.Replace(cls[0], "language-", "", 1)
	}
	if lang == "" || lang == "math" {
		return nil
	}

	var lineNumbers []int
	if meta != "" {
		for _, m := range matchBraces(meta) {
			if m != "" {
				lineNumbers = append(lineNumbers, parseRange(m)...)
			}
		}
	}
	lineIDs := map[int]string{}
	for _, n := range lineNumbers {
		if id := lineID(n, meta); id != "" {
			lineIDs[n] = id
		}
	}
	var charsList []string
	var charsRanges [][]int
	charsIDs := map[string]string{}
	if meta != "" {
		for _, m := range matchChars(meta) {
			charsList = append(charsList, m.chars)
			if m.rest == "" {
				charsRanges = append(charsRanges, []int{})
			} else {
				rng, id, _ := strings.Cut(m.rest, "#")
				if rng != "" {
					charsRanges = append(charsRanges, parseRange(rng))
				}
				if id != "" {
					charsIDs[m.chars] = id
				}
			}
		}
	}

	text, ok := codeEl.Children[0].(*hast.Text)
	if !ok {
		return nil
	}
	stripped := strings.TrimSuffix(text.Value, "\n")

	hlLang := lang
	switch {
	case h.loaded[lang] || specialLanguages[lang]:
	case bundledLanguages[lang]:
		// Upstream loads any language of shiki's bundle on demand.
		return fmt.Errorf("highlight: language %q is not in the vendored shiki bundle (see gen/entry.mjs)", lang)
	default:
		// shiki does not know the language: rehype-pretty-code catches the
		// error and highlights as plain text. Other errors are not
		// swallowed, so a failure of the goja runtime never silently turns
		// into plain text.
		hlLang = "plaintext"
	}
	tree, err := h.codeToHast(stripped, hlLang, false)
	if err != nil {
		return fmt.Errorf("highlight: %w", err)
	}
	// The docs' transformer (lib/highlight-code.ts) runs on shiki's <code>.
	if code, ok := firstChild(tree); ok {
		docsTransformer(code, stripped)
	}

	lineCounter := 0
	lineNumbersMaxDigits := 0
	opts := &charsOptions{ranges: charsRanges, ids: charsIDs, counter: map[string]int{}}
	rmeta := reverse(meta)
	var visit func(el *hast.Element)
	visit = func(el *hast.Element) {
		if el.Tag == "code" && showLineNumbersRe(rmeta) {
			el.Set("data-line-numbers", "")
			if start, ok := lineNumbersStart(rmeta); ok {
				startAt := start - 1
				lineNumbersMaxDigits = startAt
				el.Set("style", "counter-set: line "+strconv.Itoa(startAt)+";")
			}
		}
		if cls := el.Classes(); len(cls) > 0 && cls[0] == "line" {
			if hast.ToString(el) == "" {
				el.Children = []hast.Node{&hast.Text{Value: " "}}
			}
			replaceLineClass(el)
			lineCounter++
			if contains(lineNumbers, lineCounter) {
				el.Set("data-highlighted-line", "")
				if id := lineIDs[lineCounter]; id != "" {
					el.Set("data-highlighted-line-id", id)
				}
			}
			charsHighlighter(el, charsList, opts)
			lineNumbersMaxDigits++
		}
		for i := 0; i < len(el.Children); i++ {
			if c, ok := el.Children[i].(*hast.Element); ok {
				visit(c)
			}
		}
	}
	visit(tree)

	// apply()
	pre.Tag = "figure"
	pre.Set("data-rehype-pretty-code-figure", "")
	codeData := codeEl.Data
	var children []hast.Node
	shikiPre := tree
	if cls := shikiPre.Classes(); len(cls) > 0 && contains(cls, "shiki") {
		var kept []string
		for _, c := range cls {
			if c != "shiki" && c != "shiki-themes" && !contains(Themes, c) {
				kept = append(kept, c)
			}
		}
		if len(kept) > 0 {
			shikiPre.Set("class", strings.Join(kept, " "))
		} else {
			shikiPre.Unset("class")
		}
	}
	shikiPre.Set("data-language", lang)
	shikiPre.Set("data-theme", themeNames)
	code, ok := firstChild(shikiPre)
	if ok {
		code.Set("data-language", lang)
		code.Set("data-theme", themeNames)
		code.Data = codeData
		if v, ok := code.Get("style"); ok {
			code.Set("style", v+"display: grid;")
		} else {
			code.Set("style", "display: grid;")
		}
		if code.HasKey("data-line-numbers") {
			code.Set("data-line-numbers-max-digits", strconv.Itoa(len(strconv.Itoa(lineNumbersMaxDigits))))
		}
		if hasTitle && title != "" {
			tag := "figcaption"
			if hasCaption && caption != "" {
				tag = "div"
			}
			t := hast.El(tag, "data-rehype-pretty-code-title", "", "data-language", lang, "data-theme", themeNames)
			t.Children = []hast.Node{&hast.Text{Value: title}}
			children = append(children, t)
		}
		children = append(children, shikiPre)
		if hasCaption && caption != "" {
			c := hast.El("figcaption", "data-rehype-pretty-code-caption", "", "data-language", lang, "data-theme", themeNames)
			c.Children = []hast.Node{&hast.Text{Value: caption}}
			children = append(children, c)
		}
	}
	pre.Children = children
	return nil
}

func firstChild(el *hast.Element) (*hast.Element, bool) {
	if len(el.Children) == 0 {
		return nil, false
	}
	c, ok := el.Children[0].(*hast.Element)
	return c, ok
}

// docsTransformer is the `code` hook of the transformer in
// ui/apps/v4/lib/highlight-code.ts: it stores the raw source, and package
// manager variants of install commands, as properties of <code>.
func docsTransformer(code *hast.Element, raw string) {
	code.Set("__raw__", raw)
	set := func(npm, yarn, pnpm, bun string) {
		code.Set("__npm__", raw)
		code.Set("__yarn__", yarn)
		code.Set("__pnpm__", pnpm)
		code.Set("__bun__", bun)
	}
	r := func(old, new string) string { return strings.Replace(raw, old, new, 1) }
	switch {
	case strings.HasPrefix(raw, "npm install"):
		set(raw, r("npm install", "yarn add"), r("npm install", "pnpm add"), r("npm install", "bun add"))
	case strings.HasPrefix(raw, "npx create-"):
		set(raw, r("npx create-", "yarn create "), r("npx create-", "pnpm create "), r("npx", "bunx --bun"))
	case strings.HasPrefix(raw, "npm create"):
		set(raw, r("npm create", "yarn create"), r("npm create", "pnpm create"), r("npm create", "bun create"))
	case strings.HasPrefix(raw, "npx"):
		set(raw, r("npx", "yarn dlx"), r("npx", "pnpm dlx"), r("npx", "bunx --bun"))
	case strings.HasPrefix(raw, "npm run"):
		set(raw, r("npm run", "yarn"), r("npm run", "pnpm"), r("npm run", "bun"))
	}
}

func replaceLineClass(el *hast.Element) {
	cls := el.Classes()
	if !contains(cls, "line") {
		return
	}
	var kept []string
	for _, c := range cls {
		if c != "line" {
			kept = append(kept, c)
		}
	}
	if len(kept) > 0 {
		el.Set("class", strings.Join(kept, " "))
	} else {
		el.Unset("class")
	}
	el.Set("data-line", "")
}

func contains[T comparable](list []T, v T) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// reverse is s.split("").reverse().join("") on UTF-16 code units; meta
// strings in the docs are ASCII, where that equals reversing bytes.
func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// showLineNumbersRe is /srebmuNeniLwohs(?!(.*)(\/))/.test(reversedMeta):
// "showLineNumbers" not followed (in the reversed string) by a later "/",
// i.e. not preceded by a "/" in the meta string.
func showLineNumbersRe(r string) bool {
	const word = "srebmuNeniLwohs"
	for i := 0; ; {
		j := strings.Index(r[i:], word)
		if j < 0 {
			return false
		}
		at := i + j
		rest := r[at+len(word):]
		// (?!(.*)(\/)): no "/" after it on the same line.
		line := rest
		if k := strings.IndexAny(line, "\n\r\u2028\u2029"); k >= 0 {
			line = line[:k]
		}
		if !strings.Contains(line, "/") {
			return true
		}
		i = at + 1
	}
}

// lineNumbersStart is the match of /(?:\}(\d+){)?srebmuNeniLwohs(?!(.*)(\/))/
// on the reversed meta: showLineNumbers{N} sets the first line number.
func lineNumbersStart(r string) (int, bool) {
	re := regexp.MustCompile(`(?:\}(\d+)\{)?srebmuNeniLwohs`)
	for _, loc := range re.FindAllStringSubmatchIndex(r, -1) {
		rest := r[loc[1]:]
		if k := strings.IndexAny(rest, "\n\r"); k >= 0 {
			rest = rest[:k]
		}
		if strings.Contains(rest, "/") {
			continue
		}
		if loc[2] < 0 {
			return 0, false
		}
		n, err := strconv.Atoi(reverse(r[loc[2]:loc[3]]))
		return n, err == nil
	}
	return 0, false
}

func isWord(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}

// matchBraces returns the captures of meta.matchAll(/\B\{(.*?)\}\B/g).
func matchBraces(meta string) []string {
	var out []string
	nonBoundary := func(i int) bool { // \B at position i
		before := i > 0 && isWord(meta[i-1])
		after := i < len(meta) && isWord(meta[i])
		return before == after
	}
	for i := 0; i < len(meta); {
		if meta[i] != '{' || !nonBoundary(i) {
			i++
			continue
		}
		matched := false
		for j := i + 1; j < len(meta) && meta[j] != '\n' && meta[j] != '\r'; j++ {
			if meta[j] == '}' && nonBoundary(j+1) {
				out = append(out, meta[i+1:j])
				i = j + 1
				matched = true
				break
			}
		}
		if !matched {
			i++
		}
	}
	return out
}

type charsMatch struct{ chars, rest string }

// matchChars is meta.matchAll(/(?<delimiter>["/])(?<chars>.*?)\k<delimiter>(?<charsIdAndOrRange>\S*)/g).
func matchChars(meta string) []charsMatch {
	var out []charsMatch
	for i := 0; i < len(meta); {
		d := meta[i]
		if d != '"' && d != '/' {
			i++
			continue
		}
		end := -1
		for j := i + 1; j < len(meta) && meta[j] != '\n' && meta[j] != '\r'; j++ {
			if meta[j] == d {
				end = j
				break
			}
		}
		if end < 0 {
			i++
			continue
		}
		k := end + 1
		for k < len(meta) && !isJSSpace(meta[k]) {
			k++
		}
		out = append(out, charsMatch{chars: meta[i+1 : end], rest: meta[end+1 : k]})
		if k == i {
			k++
		}
		i = k
	}
	return out
}

func isJSSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

var lineIDRe = regexp.MustCompile(`\{[^}]+\}#[a-zA-Z0-9]+`)

func lineID(n int, meta string) string {
	for _, seg := range lineIDRe.FindAllString(meta, -1) {
		parts := strings.Split(seg, "#")
		rng, id := parts[0], parts[1]
		if rng == "" || id == "" {
			continue
		}
		m := regexp.MustCompile(`\{(.*?)\}`).FindStringSubmatch(rng)
		if m != nil && m[1] != "" && contains(parseRange(m[1]), n) {
			return id
		}
	}
	return ""
}

var (
	numRe   = regexp.MustCompile(`^-?\d+$`)
	rangeRe = regexp.MustCompile(`^(-?\d+)(-|\.\.\.?|\x{2025}|\x{2026}|\x{22EF})(-?\d+)$`)
)

// parseRange is parse-numeric-range 1.3.0.
func parseRange(s string) []int {
	var res []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if numRe.MatchString(part) {
			n, _ := strconv.Atoi(part)
			res = append(res, n)
			continue
		}
		m := rangeRe.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		lhs, _ := strconv.Atoi(m[1])
		rhs, _ := strconv.Atoi(m[3])
		incr := 1
		if lhs >= rhs {
			incr = -1
		}
		if m[2] == "-" || m[2] == ".." || m[2] == "‥" {
			rhs += incr
		}
		for i := lhs; i != rhs; i += incr {
			res = append(res, i)
		}
	}
	return res
}
