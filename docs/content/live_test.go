package content

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	h "shadcn.go/libhtml"
)

// TestLiveReference compares every page's rendering with the HTML the
// running reference site serves (DOCS_REFERENCE, e.g. http://localhost:4000;
// skipped when unset). It takes the content of the page's
// <div class="typeset ..."> (the MDX output; on component pages its first
// child, the base switcher, is page chrome and is dropped) and aligns it
// with the Go rendering, where components without an implementation are
// placeholders: every node outside a placeholder must be identical, in
// order; the live nodes a placeholder stands for are skipped (the next Go
// node must match right after them), and a run of N placeholders must stand
// for exactly N live nodes. Markdown inside a placeholder (for
// example the paragraphs in a Callout, the steps in Steps) is looked up in
// the skipped live nodes and reported.
func TestLiveReference(t *testing.T) {
	base := os.Getenv("DOCS_REFERENCE")
	if base == "" {
		t.Skip("DOCS_REFERENCE not set")
	}
	pages := parsedPages(t)
	reg := NewRegistry()
	type result struct {
		rel                      string
		err                      string
		compared, skipped        int
		innerFound, innerMissing int
		relAdded                 int
		notServed                string
		toc                      int
	}
	rels := oraclePages(t)
	results := make([]result, len(rels))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	for i, rel := range rels {
		wg.Add(1)
		go func(i int, rel string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			res := result{rel: rel}
			defer func() { results[i] = res }()
			p := pages[rel]
			doc, err := fetchPage(base + p.URL)
			if ns, ok := err.(notServedError); ok {
				res.notServed = string(ns)
				return
			}
			if err != nil {
				res.err = err.Error()
				return
			}
			if err := checkHeader(p, doc); err != nil {
				res.err = err.Error()
				return
			}
			n, err := checkTOC(p, doc)
			if err != nil {
				res.err = err.Error()
				return
			}
			res.toc = n
			live, err := typeset(doc)
			if ns, ok := err.(notServedError); ok {
				res.notServed = string(ns)
				return
			}
			if err != nil {
				res.err = err.Error()
				return
			}
			if len(p.Slugs) >= 3 && p.Slugs[0] == "components" {
				live = dropFirstElement(live)
			}
			node, err := Render(p, reg)
			if err != nil {
				res.err = err.Error()
				return
			}
			s, err := h.RenderString(node)
			if err != nil {
				res.err = err.Error()
				return
			}
			got, err := parseHTML(s)
			if err != nil {
				res.err = err.Error()
				return
			}
			a := &aligner{}
			if err := a.align(got, live, "typeset"); err != nil {
				res.err = err.Error()
			}
			res.compared, res.skipped, res.innerFound, res.innerMissing, res.relAdded = a.compared, a.skipped, a.innerFound, a.innerMissing, a.relAdded
		}(i, rel)
	}
	wg.Wait()
	var compared, skipped, found, missing, relAdded, failed, notServed, toc int
	for _, r := range results {
		if r.notServed != "" {
			notServed++
			t.Logf("not served by the docs page route (compared with React only): %s: %s", r.rel, r.notServed)
			continue
		}
		compared += r.compared
		skipped += r.skipped
		found += r.innerFound
		missing += r.innerMissing
		relAdded += r.relAdded
		toc += r.toc
		if r.err != "" {
			failed++
			if failed <= 10 {
				t.Errorf("%s: %s", r.rel, r.err)
			}
		}
	}
	t.Logf("%d pages: %d served, %d identical, %d different", len(rels), len(rels)-notServed, len(rels)-notServed-failed, failed)
	t.Logf("title, description (h1, <title>, meta and header paragraph) and TOC identical on every served page; %d TOC entries compared", toc)
	t.Logf("%d nodes compared outside placeholders; %d live subtrees stand for placeholders", compared, skipped)
	t.Logf("Markdown inside placeholders: %d subtrees found identical in the live component output, %d not rendered there (e.g. inactive tabs)", found, missing)
	t.Logf("%d links where libhtml adds rel=\"noopener noreferrer\"", relAdded)
}

// notServedError reports a page whose URL upstream serves with something
// else: a redirect in next.config.mjs, or another route.
type notServedError string

func (e notServedError) Error() string { return string(e) }

var noRedirect = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

func fetchPage(url string) (string, error) {
	resp, err := noRedirect.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return "", notServedError(fmt.Sprintf("%s redirects to %s", url, resp.Header.Get("Location")))
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	s := string(b)
	if !strings.Contains(s, docsTypeset) {
		return "", notServedError(fmt.Sprintf("%s is not rendered by the docs page route (no docs typeset div)", url))
	}
	return s, nil
}

const docsTypeset = `<div class="typeset w-full flex-1`

// element parses the element whose start tag begins with prefix.
func element(doc, prefix string) (*dnode, error) {
	start := strings.Index(doc, prefix)
	if start < 0 {
		return nil, fmt.Errorf("no %s", prefix)
	}
	end, err := elementEnd(doc, start)
	if err != nil {
		return nil, err
	}
	nodes, err := parseHTML(doc[start:end])
	if err != nil {
		return nil, err
	}
	if len(nodes) != 1 {
		return nil, fmt.Errorf("%s: expected one element", prefix)
	}
	return nodes[0], nil
}

func typeset(doc string) ([]*dnode, error) {
	n, err := element(doc, docsTypeset)
	if err != nil {
		return nil, err
	}
	return n.children, nil
}

func textOf(n *dnode) string {
	if n.tag == "" {
		return n.text
	}
	var b strings.Builder
	for _, c := range n.children {
		b.WriteString(textOf(c))
	}
	return b.String()
}

// checkHeader compares the page title and description with what the docs
// page shows: the <h1>, the description paragraph, <title> (template
// "%s - shadcn/ui") and the description meta tag.
func checkHeader(p *Page, doc string) error {
	fm := p.Frontmatter
	h1, err := element(doc, `<h1 class="scroll-m-24 text-3xl font-semibold tracking-tight sm:text-3xl">`)
	if err != nil {
		return err
	}
	if got := textOf(h1); got != fm.Title {
		return fmt.Errorf("h1 %q, frontmatter title %q", got, fm.Title)
	}
	desc, err := element(doc, `<p class="text-[1.05rem] text-muted-foreground sm:text-base sm:text-balance md:max-w-[80%]">`)
	if err != nil {
		return err
	}
	if got := textOf(desc); got != fm.Description {
		return fmt.Errorf("description %q, frontmatter %q", got, fm.Description)
	}
	title, err := element(doc, `<title>`)
	if err != nil {
		return err
	}
	if got := textOf(title); got != fm.Title+" - shadcn/ui" {
		return fmt.Errorf("<title> %q, frontmatter title %q", got, fm.Title)
	}
	const metaPrefix = `<meta name="description" content="`
	i := strings.Index(doc, metaPrefix)
	if i < 0 {
		return fmt.Errorf("no description meta")
	}
	end := strings.IndexByte(doc[i+len(metaPrefix):], '"')
	if got := html.UnescapeString(doc[i+len(metaPrefix) : i+len(metaPrefix)+end]); got != fm.Description {
		return fmt.Errorf("meta description %q, frontmatter %q", got, fm.Description)
	}
	return nil
}

// checkTOC compares the TOC with the "On This Page" list of the page
// (DocsTableOfContents): one link per entry, with the URL, the depth and
// the title markup.
func checkTOC(p *Page, doc string) (int, error) {
	const list = `<div class="flex flex-col gap-2 p-4 pt-0 text-sm"><p class="h-6 bg-background text-xs font-medium text-muted-foreground">On This Page</p>`
	if !strings.Contains(doc, list) {
		if len(p.TOC) != 0 {
			return 0, fmt.Errorf("TOC has %d entries, the page shows none", len(p.TOC))
		}
		return 0, nil
	}
	n, err := element(doc, list)
	if err != nil {
		return 0, err
	}
	var links []*dnode
	for _, c := range n.children[1:] {
		if c.tag == "a" {
			links = append(links, c)
		}
	}
	if len(links) != len(p.TOC) {
		return 0, fmt.Errorf("TOC has %d entries, the page shows %d", len(p.TOC), len(links))
	}
	for i, item := range p.TOC {
		l := links[i]
		if l.attrs["href"] != item.URL || l.attrs["data-depth"] != fmt.Sprint(item.Depth) {
			return 0, fmt.Errorf("TOC entry %d: %s depth %d, page shows %s depth %s", i, item.URL, item.Depth, l.attrs["href"], l.attrs["data-depth"])
		}
		title, err := RenderTOCTitle(p, item, nil)
		if err != nil {
			return 0, err
		}
		s, err := h.RenderString(title)
		if err != nil {
			return 0, err
		}
		g, err := parseHTML(s)
		if err != nil {
			return 0, err
		}
		if canon(g) != canon(l.children) {
			return 0, fmt.Errorf("TOC entry %s: title %s, page shows %s", item.URL, canon(g), canon(l.children))
		}
	}
	return len(p.TOC), nil
}

// elementEnd returns the offset after the element starting at s[start].
func elementEnd(s string, start int) (int, error) {
	depth := 0
	i := start
	for i < len(s) {
		lt := strings.IndexByte(s[i:], '<')
		if lt < 0 {
			break
		}
		i += lt
		switch {
		case strings.HasPrefix(s[i:], "<!--"):
			e := strings.Index(s[i:], "-->")
			i += e + 3
		case strings.HasPrefix(s[i:], "</"):
			e := strings.IndexByte(s[i:], '>')
			depth--
			i += e + 1
			if depth == 0 {
				return i, nil
			}
		default:
			n, next, selfClosing, err := parseStartTag(s, i)
			if err != nil {
				return 0, err
			}
			i = next
			if rawTextTags[n.tag] {
				e := strings.Index(s[i:], "</"+n.tag)
				i += e
				i += strings.IndexByte(s[i:], '>') + 1
				continue
			}
			if !voidTags[n.tag] && !selfClosing {
				depth++
			}
		}
	}
	return 0, fmt.Errorf("unterminated element")
}

func dropFirstElement(ns []*dnode) []*dnode {
	for i, n := range ns {
		if n.tag != "" {
			return append(append([]*dnode{}, ns[:i]...), ns[i+1:]...)
		}
	}
	return ns
}

type aligner struct {
	compared, skipped, innerFound, innerMissing, relAdded int
}

func isPlaceholder(n *dnode) bool { return n.tag == "mdx-component" }

func (a *aligner) align(got, live []*dnode, path string) error {
	j := 0
	for i := 0; i < len(got); i++ {
		g := got[i]
		if isPlaceholder(g) {
			// A run of consecutive placeholders stands for the live nodes up
			// to the one matching the next Go node.
			run := []*dnode{g}
			for i+1 < len(got) && isPlaceholder(got[i+1]) {
				i++
				run = append(run, got[i])
			}
			var anchor *dnode
			if i+1 < len(got) {
				anchor = got[i+1]
			}
			start := j
			for j < len(live) && (anchor == nil || !a.equalNoCount(anchor, live[j])) {
				j++
			}
			if anchor != nil && j >= len(live) {
				return fmt.Errorf("%s: after placeholder <%s>, no live node matches the next node %s", path, g.attrs["data-name"], short(canon([]*dnode{anchor})))
			}
			// Every component renders one element: the region must hold
			// exactly one live node per placeholder, so nothing else can hide
			// in it.
			if j-start != len(run) {
				return fmt.Errorf("%s: %d placeholder(s) starting with <%s> stand for %d live nodes", path, len(run), g.attrs["data-name"], j-start)
			}
			a.skipped += j - start
			for _, ph := range run {
				a.inner(ph, live[start:j])
			}
			continue
		}
		if j >= len(live) {
			return fmt.Errorf("%s: extra Go node %s", path, short(canon([]*dnode{g})))
		}
		l := live[j]
		if err := a.compare(g, l, path); err != nil {
			return err
		}
		j++
	}
	if j < len(live) {
		return fmt.Errorf("%s: extra live node %s", path, short(canon(live[j:j+1])))
	}
	return nil
}

func (a *aligner) compare(g, l *dnode, path string) error {
	if g.tag != l.tag {
		return fmt.Errorf("%s: got %s, live %s", path, short(canon([]*dnode{g})), short(canon([]*dnode{l})))
	}
	a.compared++
	if g.tag == "" {
		if g.text != l.text {
			return fmt.Errorf("%s: text %q, live %q", path, g.text, l.text)
		}
		return nil
	}
	if g.tag == "a" && g.attrs["target"] == "_blank" {
		if _, has := l.attrs["rel"]; !has && g.attrs["rel"] == "noopener noreferrer" {
			delete(g.attrs, "rel")
			a.relAdded++
		}
	}
	if ga, la := attrString(g), attrString(l); ga != la {
		return fmt.Errorf("%s/%s: attributes %s, live %s", path, g.tag, ga, la)
	}
	return a.align(g.children, l.children, path+"/"+g.tag)
}

func attrString(n *dnode) string {
	return canon([]*dnode{{tag: "x", attrs: n.attrs}})
}

// equalNoCount reports whether g (which may contain placeholders) matches l.
func (a *aligner) equalNoCount(g, l *dnode) bool {
	b := &aligner{}
	return b.compare(g, l, "") == nil
}

// inner looks for the placeholder's Markdown children in the live nodes
// it stands for.
func (a *aligner) inner(ph *dnode, live []*dnode) {
	for _, c := range ph.children {
		if c.tag == "" || c.tag == "template" {
			continue
		}
		if isPlaceholder(c) {
			a.inner(c, live)
			continue
		}
		if containsEqual(live, c) {
			a.innerFound++
		} else {
			a.innerMissing++
		}
	}
}

func containsEqual(ns []*dnode, want *dnode) bool {
	for _, n := range ns {
		if n.tag == want.tag && (&aligner{}).compare(want, n, "") == nil {
			return true
		}
		if containsEqual(n.children, want) {
			return true
		}
	}
	return false
}

func short(s string) string {
	if len(s) > 400 {
		return s[:400] + "…"
	}
	return s
}
