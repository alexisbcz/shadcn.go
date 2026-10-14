package site

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

// TestSearchMatchesReference compares /api/search with the reference
// server's (SEARCH_REF, e.g. http://localhost:4000).
func TestSearchMatchesReference(t *testing.T) {
	ref := os.Getenv("SEARCH_REF")
	if ref == "" {
		t.Skip("SEARCH_REF not set")
	}
	s, err := New(Options{Root: "../..", Logf: t.Logf})
	if err != nil {
		t.Fatal(err)
	}
	t0 := time.Now()
	if _, err := s.searchJSON("x"); err != nil {
		t.Fatal(err)
	}
	t.Logf("index built in %v", time.Since(t0))
	for _, q := range []string{"accordion", "button", "dark mode", "registry", "instal", "theming css variables", "tailwind v4", "form", "sidebar", "mcp", "cli", "toast", "xyzzy", "radix", "data table"} {
		got, err := s.searchJSON(q)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Get(ref + "/api/search?query=" + url.QueryEscape(q))
		if err != nil {
			t.Fatal(err)
		}
		var want []map[string]any
		json.NewDecoder(resp.Body).Decode(&want)
		resp.Body.Close()
		if !reflect.DeepEqual(norm(got), norm(want)) {
			n := 0
			for i := 0; i < len(got) || i < len(want); i++ {
				var g, w any
				if i < len(got) {
					g = got[i]
				}
				if i < len(want) {
					w = want[i]
				}
				if !reflect.DeepEqual(g, w) {
					t.Errorf("%q [%d]:\n got  %v\n want %v", q, i, g, w)
					if n++; n > 3 {
						break
					}
				}
			}
			t.Errorf("%q: %d results, want %d", q, len(got), len(want))
		}
	}
}

func norm(v []map[string]any) any {
	b, _ := json.Marshal(v)
	var out any
	json.Unmarshal(b, &out)
	return out
}
