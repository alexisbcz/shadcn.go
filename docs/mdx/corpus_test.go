package mdx

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"testing"
)

// corpusEntry is one input of docs/content/testdata/corpus.json.gz
// (docs/oracle/corpus.mjs): CommonMark spec examples, docs blocks and
// seeded mutations of them, with upstream's mdast or parse error.
type corpusEntry struct {
	Kind       string         `json:"kind"`
	Src        string         `json:"src"`
	Mdast      map[string]any `json:"mdast"`
	MdastError string         `json:"mdastError"`
}

// TestCorpusMdast parses every corpus input. The parser may reject an
// input (it only supports the constructs the docs use), but when it
// accepts one, upstream must accept it too and the mdast must be identical.
func TestCorpusMdast(t *testing.T) {
	f, err := os.Open("../content/testdata/corpus.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var list []corpusEntry
	if err := json.NewDecoder(zr).Decode(&list); err != nil {
		t.Fatal(err)
	}
	stats := map[string]*[4]int{} // accepted+identical, rejected, accepted but upstream errors, different
	failed := 0
	for i, e := range list {
		s := stats[e.Kind]
		if s == nil {
			s = &[4]int{}
			stats[e.Kind] = s
		}
		root, err := parseSafe(e.Src)
		if err != nil {
			s[1]++
			continue
		}
		if e.MdastError != "" {
			s[2]++
			failed++
			if failed <= corpusLimit {
				t.Errorf("input %d (%s) accepted, upstream fails (%s):\n%q", i, e.Kind, e.MdastError, e.Src)
			}
			continue
		}
		gb, _ := json.Marshal(canon(root))
		var g any
		json.Unmarshal(gb, &g)
		if d := diff(g, any(e.Mdast), "root"); d != "" {
			s[3]++
			failed++
			if failed <= corpusLimit {
				t.Errorf("input %d (%s) differs: %s\n%q", i, e.Kind, d, e.Src)
			}
			continue
		}
		s[0]++
	}
	for _, k := range []string{"spec", "curated", "docs", "mutation"} {
		if s := stats[k]; s != nil {
			t.Logf("%-8s identical %5d, rejected %5d, accepted but upstream fails %d, different %d", k, s[0], s[1], s[2], s[3])
		}
	}
}

// parseSafe turns a parser panic into an error, so one bad input reports
// instead of aborting the run.
func parseSafe(src string) (n *Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			n, err = nil, panicError{r}
		}
	}()
	return Parse(src)
}

type panicError struct{ v any }

func (p panicError) Error() string { return "panic" }

// corpusLimit caps the failures reported (CORPUS_ALL=1 reports all).
var corpusLimit = func() int {
	if os.Getenv("CORPUS_ALL") != "" {
		return 1 << 30
	}
	return 20
}()
