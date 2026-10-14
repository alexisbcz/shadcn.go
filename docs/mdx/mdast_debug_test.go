package mdx

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestMdastOracle compares the parser's mdast with upstream's for every
// page. The oracle files are produced by docs/oracle/mdast.mjs into the
// directory named by MDAST_ORACLE (they are not committed: the committed
// proof is the hast comparison in package content).
func TestMdastOracle(t *testing.T) {
	dir := os.Getenv("MDAST_ORACLE")
	if dir == "" {
		t.Skip("MDAST_ORACLE not set")
	}
	content := "../../ui/apps/v4/content/docs"
	var pages []string
	filepath.WalkDir(content, func(p string, d os.DirEntry, err error) error {
		if strings.HasSuffix(p, ".mdx") {
			rel, _ := filepath.Rel(content, p)
			pages = append(pages, rel)
		}
		return nil
	})
	sort.Strings(pages)
	fails := 0
	for _, rel := range pages {
		if f := os.Getenv("PAGE"); f != "" && !strings.Contains(rel, f) {
			continue
		}
		src, _ := os.ReadFile(filepath.Join(content, rel))
		body := stripFrontmatter(string(src))
		want, err := os.ReadFile(filepath.Join(dir, strings.TrimSuffix(rel, ".mdx")+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var w any
		json.Unmarshal(want, &w)
		root, err := Parse(body)
		if err != nil {
			fails++
			t.Errorf("%s: %v", rel, err)
			continue
		}
		g := canon(root)
		gb, _ := json.Marshal(g)
		var gg any
		json.Unmarshal(gb, &gg)
		if d := diff(gg, w, "root"); d != "" {
			fails++
			if fails <= 15 {
				t.Errorf("%s: %s", rel, d)
			}
		}
	}
	t.Logf("%d pages, %d differ", len(pages), fails)
}

var fmRe = regexp.MustCompile(`(?s)^---\r?\n(.+?)\r?\n---\r?\n?`)

func stripFrontmatter(s string) string {
	if m := fmRe.FindStringIndex(s); m != nil {
		return s[m[1]:]
	}
	return s
}

func canon(n *Node) map[string]any {
	o := map[string]any{"type": n.Type}
	switch n.Type {
	case "text", "inlineCode", "mdxFlowExpression", "mdxTextExpression", "mdxjsEsm", "code":
		o["value"] = n.Value
	}
	switch n.Type {
	case "heading":
		o["depth"] = n.Depth
	case "list":
		o["ordered"] = n.Ordered
		if n.Ordered {
			o["start"] = n.Start
		}
		o["spread"] = n.Spread
	case "listItem":
		o["spread"] = n.Spread
	case "code":
		if n.HasLang {
			o["lang"] = n.Lang
		}
		if n.HasMeta {
			o["meta"] = n.Meta
		}
	case "link", "image":
		o["url"] = n.URL
		if n.HasTitle {
			o["title"] = n.Title
		}
		if n.Type == "image" {
			o["alt"] = n.Alt
		}
	case "table":
		var a []any
		for _, x := range n.Align {
			if x == "" {
				a = append(a, nil)
			} else {
				a = append(a, x)
			}
		}
		o["align"] = a
	case "footnoteReference", "footnoteDefinition":
		o["identifier"] = n.Identifier
		o["label"] = n.Label
	case "mdxJsxFlowElement", "mdxJsxTextElement":
		if n.Name != "" {
			o["name"] = n.Name
		}
		var attrs []any
		for _, a := range n.Attrs {
			m := map[string]any{"name": a.Name}
			switch a.Kind {
			case 0:
				m["value"] = a.Value
			case 2:
				m["expr"] = a.Value
			}
			attrs = append(attrs, m)
		}
		if attrs == nil {
			attrs = []any{}
		}
		o["attributes"] = attrs
	}
	switch n.Type {
	case "text", "inlineCode", "code", "thematicBreak", "mdxFlowExpression", "mdxTextExpression", "mdxjsEsm", "image", "footnoteReference", "break":
	default:
		kids := []any{}
		for _, c := range n.Children {
			kids = append(kids, canon(c))
		}
		o["children"] = kids
	}
	return o
}

func diff(a, b any, path string) string {
	switch x := a.(type) {
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok {
			return path + ": type differs"
		}
		keys := map[string]bool{}
		for k := range x {
			keys[k] = true
		}
		for k := range y {
			keys[k] = true
		}
		var ks []string
		for k := range keys {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		// type first
		for _, k := range append([]string{"type"}, ks...) {
			if d := diff(x[k], y[k], path+"."+k); d != "" {
				if k == "children" || k == "type" {
					return fmt.Sprintf("%s [got type %v]", d, x["type"])
				}
				return d
			}
		}
		return ""
	case []any:
		y, ok := b.([]any)
		if !ok {
			return path + ": got array, want " + short(b)
		}
		for i := 0; i < len(x) && i < len(y); i++ {
			if d := diff(x[i], y[i], fmt.Sprintf("%s[%d]", path, i)); d != "" {
				return d
			}
		}
		if len(x) != len(y) {
			var extra any
			if len(x) > len(y) {
				extra = x[len(y)]
			} else {
				extra = y[len(x)]
			}
			return fmt.Sprintf("%s: got %d items, want %d (first extra: %s)", path, len(x), len(y), short(extra))
		}
		return ""
	default:
		if fmt.Sprint(a) != fmt.Sprint(b) {
			return fmt.Sprintf("%s: got %s, want %s", path, short(a), short(b))
		}
		return ""
	}
}

func short(v any) string {
	b, _ := json.Marshal(v)
	s := string(b)
	if len(s) > 300 {
		s = s[:300] + "..."
	}
	return s
}
