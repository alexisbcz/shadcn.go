package content

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"shadcn.go/docs/hast"
	"shadcn.go/docs/highlight"
)

const contentRoot = "../../ui/apps/v4/content/docs"

type oraclePage struct {
	Path        string         `json:"path"`
	Slugs       []string       `json:"slugs"`
	Frontmatter map[string]any `json:"frontmatter"`
	TOC         []struct {
		Depth int    `json:"depth"`
		URL   string `json:"url"`
		Title string `json:"title"`
		Hast  string `json:"hast"`
	} `json:"toc"`
	Hast string `json:"hast"`
}

var (
	oracleOnce sync.Once
	oracleData map[string]*oraclePage
	reactData  map[string]string
)

// loadOracles reads testdata/oracle.json.gz (docs/oracle/build.mjs) and
// testdata/react.json.gz (docs/oracle/render.mjs).
func loadOracles(t *testing.T) {
	oracleOnce.Do(func() {
		read := func(name string, v any) {
			f, err := os.Open(filepath.Join("testdata", name))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			zr, err := gzip.NewReader(f)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.NewDecoder(zr).Decode(v); err != nil {
				t.Fatal(err)
			}
		}
		var list []*oraclePage
		read("oracle.json.gz", &list)
		oracleData = map[string]*oraclePage{}
		for _, o := range list {
			oracleData[o.Path] = o
		}
		read("react.json.gz", &reactData)
	})
}

func oraclePages(t *testing.T) []string {
	loadOracles(t)
	var out []string
	for p := range oracleData {
		out = append(out, p)
	}
	sort.Strings(out)
	if len(out) != 316 {
		t.Fatalf("expected 316 oracle pages, found %d", len(out))
	}
	return out
}

func loadOracle(t *testing.T, rel string) *oraclePage {
	loadOracles(t)
	o := oracleData[rel]
	if o == nil {
		t.Fatalf("%s: not in the oracle", rel)
	}
	return o
}

var (
	hlOnce sync.Once
	hl     *highlight.Highlighter
)

func testHighlighter(t *testing.T) *highlight.Highlighter {
	hlOnce.Do(func() {
		var err error
		hl, err = highlight.New()
		if err != nil {
			t.Fatal(err)
		}
	})
	return hl
}

var (
	pagesOnce sync.Once
	allPages  map[string]*Page
	pagesErr  error
)

func parsedPages(t *testing.T) map[string]*Page {
	h := testHighlighter(t)
	pagesOnce.Do(func() {
		allPages = map[string]*Page{}
		pages, err := LoadDocs(Options{Highlighter: h})
		if err != nil {
			pagesErr = err
			return
		}
		for _, p := range pages {
			allPages[p.Path] = p
		}
	})
	if pagesErr != nil {
		t.Fatal(pagesErr)
	}
	return allPages
}

// TestPagesHast checks every docs page against upstream's pipeline
// (testdata/oracle.json.gz, written by docs/oracle/build.mjs): the HTML tree after
// all remark and rehype plugins must be byte-identical, as must the TOC
// (entries and title markup), the frontmatter and the URL slugs.
func TestPagesHast(t *testing.T) {
	pages := parsedPages(t)
	if len(pages) != 316 {
		t.Fatalf("parsed %d pages, want 316", len(pages))
	}
	failed := 0
	for _, rel := range oraclePages(t) {
		o := loadOracle(t, rel)
		p := pages[rel]
		if p == nil {
			t.Errorf("%s: not parsed", rel)
			failed++
			continue
		}
		var problems []string
		if got := hast.ToHTML(withoutExports(p.Tree)...); got != o.Hast {
			problems = append(problems, "hast differs:\n got: "+around(got, o.Hast)+"\nwant: "+around(o.Hast, got))
		}
		if strings.Join(p.Slugs, "/") != strings.Join(o.Slugs, "/") {
			problems = append(problems, "slugs differ: "+strings.Join(p.Slugs, "/")+" vs "+strings.Join(o.Slugs, "/"))
		}
		if len(p.TOC) != len(o.TOC) {
			problems = append(problems, "TOC length differs")
		} else {
			for i, it := range p.TOC {
				w := o.TOC[i]
				if it.Depth != w.Depth || it.URL != w.URL || hast.ToHTML(it.Title...) != w.Hast {
					problems = append(problems, "TOC entry differs: "+it.URL+" vs "+w.URL+" / "+hast.ToHTML(it.Title...)+" vs "+w.Hast)
					break
				}
			}
		}
		if d := compareFrontmatter(p.Frontmatter, o.Frontmatter); d != "" {
			problems = append(problems, "frontmatter: "+d)
		}
		if len(problems) > 0 {
			failed++
			if failed <= 8 {
				t.Errorf("%s:\n%s", rel, strings.Join(problems, "\n"))
			}
		}
	}
	t.Logf("%d pages, %d identical, %d different", len(pages), len(pages)-failed, failed)
}

// withoutExports drops the ESM nodes fumadocs-mdx injects with an empty
// value (frontmatter, toc, structuredData exports); the oracle drops them
// too. The page's own imports are kept.
func withoutExports(nodes []hast.Node) []hast.Node {
	return nodes
}

func compareFrontmatter(fm *Frontmatter, want map[string]any) string {
	got := map[string]any{}
	var conv func(kvs []KV) map[string]any
	conv = func(kvs []KV) map[string]any {
		m := map[string]any{}
		for _, kv := range kvs {
			switch v := kv.Value.(type) {
			case []KV:
				m[kv.Key] = conv(v)
			case time.Time:
				m[kv.Key] = v.UTC().Format("2006-01-02T15:04:05.000Z")
			default:
				m[kv.Key] = v
			}
		}
		return m
	}
	got = conv(fm.Values)
	a, _ := json.Marshal(got)
	b, _ := json.Marshal(want)
	if string(a) != string(b) {
		return string(a) + " vs " + string(b)
	}
	if fm.Title != want["title"] || fm.Description != want["description"] {
		return "title/description"
	}
	return ""
}

func around(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	lo := max(i-200, 0)
	hi := min(i+200, len(a))
	return a[lo:hi]
}
