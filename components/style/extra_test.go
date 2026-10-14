package style

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
)

// extraRef is testdata/ref/extra.json.gz, written by testdata/refgen/extra.mjs
// with upstream's transformStyle and prettier: base sources the registry
// build skips, synthetic sources over every token of every style map, and
// the docs examples. Gen is nil when the build leaves the roots unchanged.
type extraRef struct {
	Cases []struct {
		Set   string    `json:"set"`
		Name  string    `json:"name"`
		Style Style     `json:"style"`
		Base  []refRoot `json:"base"`
		Gen   []refRoot `json:"gen"`
	} `json:"cases"`
}

func loadExtra(t *testing.T) *extraRef {
	t.Helper()
	f, err := os.Open("testdata/ref/extra.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var r extraRef
	if err := json.NewDecoder(zr).Decode(&r); err != nil {
		t.Fatal(err)
	}
	return &r
}

// TestStyleValuesHaveNoTokens checks an invariant the API relies on: no
// style map value contains a cn-* token. Then a className property's tokens
// are all gone (or allowlisted, which is never styled) after upstream's
// property pass, so its later mergeProps pass changes nothing and marks no
// token, and PropCn need not model it.
func TestStyleValuesHaveNoTokens(t *testing.T) {
	for _, st := range All {
		for k, v := range st.Map() {
			if tk := extractCnClasses(v); len(tk) > 0 {
				t.Errorf("%s: %s = %q contains %v", st, k, v, tk)
			}
		}
	}
}

// TestExtraCorpus replays the extra corpus through the engine (transform
// and sorting) and through the public API (File.CVA, File.Cn, File.Attr and
// the verbatim parts), and requires upstream's output exactly. It also
// checks that every style token was exercised.
func TestExtraCorpus(t *testing.T) {
	ref := loadExtra(t)
	type count struct{ files, roots, strings, mismatches int }
	counts := map[string]*count{}
	var api apiStats
	tokens := map[Style]map[string]bool{}
	for _, c := range ref.Cases {
		key := fmt.Sprintf("%-8s %s", c.Set, c.Style)
		n := counts[key]
		if n == nil {
			n = &count{}
			counts[key] = n
		}
		n.files++
		want := c.Gen
		if want == nil {
			want = c.Base
		}
		base, gen := decodeRoots(c.Base), decodeRoots(want)
		if len(base) != len(gen) {
			t.Fatalf("%s %s: %d roots, want %d", c.Name, c.Style, len(base), len(gen))
		}
		if c.Set == "tokens" {
			if tokens[c.Style] == nil {
				tokens[c.Style] = map[string]bool{}
			}
			walkRoots(base, func(x *node) {
				if x.kind == kStr {
					for _, tk := range extractCnClasses(x.str) {
						tokens[c.Style][tk] = true
					}
				}
			})
		}
		transformStyleMap(base, c.Style.Map())
		sortTree(base)
		for i := range base {
			s := 0
			if d := compareTrees(base[i], gen[i], "", &s); d != "" {
				n.mismatches++
				if n.mismatches <= 3 {
					t.Errorf("%s %s (%s): line %d: %s", c.Set, c.Name, c.Style, c.Base[i].Line, d)
				}
			}
			n.roots++
			n.strings += s
		}
		checkFileAPI(t, c.Set+"-"+string(c.Style), refFile{Path: c.Name, Base: c.Base, Gen: want}, c.Style, &api)
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	total := 0
	for _, k := range keys {
		n := counts[k]
		total += n.strings
		t.Logf("%s: %4d files %6d roots %7d strings %d mismatches", k, n.files, n.roots, n.strings, n.mismatches)
	}
	t.Logf("engine: %d strings compared; API: cva %d (strings %d), cn sites %d, attr sites %d, %d strings, verbatim %d, %d mismatches",
		total, api.cvas, api.cvaStrings, api.cnSites, api.attrSites, api.strings, api.verbatim, api.mismatches)
	if api.mismatches > 0 {
		t.Errorf("API: %d mismatches", api.mismatches)
	}
	for _, st := range All {
		for tk := range st.Map() {
			if !tokens[st][tk] {
				t.Errorf("%s: token %s not exercised by the synthetic sources", st, tk)
			}
		}
	}
	if u := UnorderedClasses(); len(u) > 0 {
		t.Errorf("classes missing from the order table: %d (%q...)", len(u), u[0])
	}
}
