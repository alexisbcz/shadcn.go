package style

import (
	"compress/gzip"
	"encoding/json"
	"maps"
	"os"
	"slices"
	"testing"
)

// upstreamRef is testdata/ref/upstream.json.gz, written by
// testdata/refgen/upstream.mjs with upstream's own createStyleMap,
// transformStyle and prettier.
type upstreamRef struct {
	StyleMaps map[string]map[string]string `json:"stylemaps"`
	Cases     []struct {
		Name      string            `json:"name"`
		CSS       *string           `json:"css"`
		Result    map[string]string `json:"result"`
		Source    string            `json:"source"`
		StyleMap  map[string]string `json:"styleMap"`
		Base      []refRoot         `json:"base"`
		Gen       []refRoot         `json:"gen"`
		Formatted []refRoot         `json:"formatted"`
		Style     Style             `json:"style"`
		FormatErr string            `json:"formatError"`
	} `json:"cases"`
	// Regex holds random strings with upstream's extractCnClasses and
	// removeCnClasses results.
	Regex []struct {
		S       string   `json:"s"`
		Extract []string `json:"extract"`
		Removed string   `json:"removed"`
	} `json:"regex"`
}

func loadUpstream(t *testing.T) *upstreamRef {
	t.Helper()
	f, err := os.Open("testdata/ref/upstream.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var u upstreamRef
	if err := json.NewDecoder(zr).Decode(&u); err != nil {
		t.Fatal(err)
	}
	return &u
}

// TestStyleMapsMatchUpstream compares CreateStyleMap of the 8 embedded
// stylesheets with upstream's createStyleMap.
func TestStyleMapsMatchUpstream(t *testing.T) {
	u := loadUpstream(t)
	for _, st := range All {
		want := u.StyleMaps[string(st)]
		got := st.Map()
		if len(want) == 0 {
			t.Fatalf("%s: no upstream map", st)
		}
		for _, k := range slices.Sorted(maps.Keys(want)) {
			if got[k] != want[k] {
				t.Errorf("%s: %s = %q, want %q", st, k, got[k], want[k])
			}
		}
		for k := range got {
			if _, ok := want[k]; !ok {
				t.Errorf("%s: extra token %s", st, k)
			}
		}
		t.Logf("%s: %d tokens identical", st, len(want))
	}
}

// TestUpstreamUnitCases runs the cases of upstream's create-style-map.test.ts
// and transform-style-map.test.ts, and the edge cases of
// testdata/refgen/cases.tsx, against upstream's results: createStyleMap,
// transformStyle alone, then transformStyle followed by prettier.
func TestUpstreamUnitCases(t *testing.T) {
	u := loadUpstream(t)
	nCSS, nTransform := 0, 0
	for _, c := range u.Cases {
		if c.CSS != nil {
			nCSS++
			got, err := CreateStyleMap(*c.CSS)
			if err != nil {
				t.Errorf("%s: %v", c.Name, err)
				continue
			}
			if !maps.Equal(got, c.Result) {
				t.Errorf("%s: got %v, want %v", c.Name, got, c.Result)
			}
			continue
		}
		nTransform++
		if c.FormatErr != "" {
			t.Errorf("%s: prettier failed upstream: %s", c.Name, c.FormatErr)
			continue
		}
		base := decodeRoots(c.Base)
		gen := decodeRoots(c.Gen)
		formatted := decodeRoots(c.Formatted)
		if len(base) != len(gen) || len(base) != len(formatted) {
			t.Errorf("%s: root counts %d/%d/%d", c.Name, len(base), len(gen), len(formatted))
			continue
		}
		transformStyleMap(base, c.StyleMap)
		for i := range base {
			n := 0
			if d := compareTrees(base[i], gen[i], "", &n); d != "" {
				t.Errorf("%s: transformStyle: line %d: %s", c.Name, c.Base[i].Line, d)
			}
		}
		sortTree(base)
		for i := range base {
			n := 0
			if d := compareTrees(base[i], formatted[i], "", &n); d != "" {
				t.Errorf("%s: prettier: line %d: %s", c.Name, c.Base[i].Line, d)
			}
		}
	}
	t.Logf("%d createStyleMap cases, %d transform cases", nCSS, nTransform)

	// The fuzz cases through the public API (File.CVA, File.Cn, File.Attr).
	api := apiStats{skipVerbatim: true, skipCVAClasses: true}
	nAPI := 0
	for _, c := range u.Cases {
		if c.Style == "" || c.FormatErr != "" {
			continue
		}
		nAPI++
		checkFileAPI(t, c.Name, refFile{Path: c.Name, Base: c.Base, Gen: c.Formatted}, c.Style, &api)
	}
	t.Logf("API on %d fuzz files: cva %d (strings %d), cn sites %d, attr sites %d, %d strings, %d mismatches",
		nAPI, api.cvas, api.cvaStrings, api.cnSites, api.attrSites, api.strings, api.mismatches)
	if nAPI < 96 || api.mismatches > 0 {
		t.Errorf("API on fuzz cases: %d files, %d mismatches", nAPI, api.mismatches)
	}
	if nCSS < 9+5 || nTransform < 33+96 {
		t.Errorf("expected at least 14 and 129 cases")
	}
}

// TestUpstreamRegex compares extractCnClasses and removeCnClasses (the
// /\bcn-[\w-]+\b/g token regex, JavaScript's \s and trim()) with upstream's
// functions, transpiled from its source, on random strings.
func TestUpstreamRegex(t *testing.T) {
	u := loadUpstream(t)
	if len(u.Regex) < 4000 {
		t.Fatalf("%d regex cases", len(u.Regex))
	}
	bad := 0
	for _, c := range u.Regex {
		got := extractCnClasses(c.S)
		if !slices.Equal(got, c.Extract) && !(len(got) == 0 && len(c.Extract) == 0) {
			bad++
			t.Errorf("extractCnClasses(%q) = %q, want %q", c.S, got, c.Extract)
		}
		if r := removeCnClasses(c.S); r != c.Removed {
			bad++
			t.Errorf("removeCnClasses(%q) = %q, want %q", c.S, r, c.Removed)
		}
		if bad > 10 {
			t.FailNow()
		}
	}
	t.Logf("%d strings: extractCnClasses and removeCnClasses identical", len(u.Regex))
}
