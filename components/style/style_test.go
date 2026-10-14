package style

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestStyles(t *testing.T) {
	if len(All) != 8 {
		t.Fatalf("%d styles", len(All))
	}
	for _, st := range All {
		if s, ok := Parse(string(st)); !ok || s != st {
			t.Errorf("Parse(%q)", st)
		}
		if !strings.HasPrefix(st.CSS(), ".style-"+string(st)+" {") {
			t.Errorf("%s: CSS starts with %q", st, st.CSS()[:20])
		}
		if len(st.Map()) < 400 {
			t.Errorf("%s: %d tokens", st, len(st.Map()))
		}
	}
	if _, ok := Parse("new-york"); ok {
		t.Error("Parse accepted an unknown style")
	}
	if Style("x").CSS() != "" || Style("x").Map() != nil {
		t.Error("invalid style has data")
	}
}

func TestDefaultAndContext(t *testing.T) {
	if Default() != Nova {
		t.Fatalf("default %s", Default())
	}
	SetDefault(Vega)
	defer SetDefault(Nova)
	if Default() != Vega || Style("").Or() != Vega || Lyra.Or() != Lyra {
		t.Fatal("SetDefault")
	}
	if FromContext(context.Background()) != Vega {
		t.Error("FromContext without a style")
	}
	if FromContext(NewContext(context.Background(), Mira)) != Mira {
		t.Error("FromContext")
	}
	h := Middleware(func(r *http.Request) (Style, bool) { return Parse(r.URL.Query().Get("style")) })(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(FromContext(r.Context()))) }))
	for q, want := range map[string]string{"?style=luma": "luma", "?style=bogus": "vega", "": "vega"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/"+q, nil))
		if rec.Body.String() != want {
			t.Errorf("%q: got %q, want %q", q, rec.Body.String(), want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("SetDefault accepted an invalid style")
		}
	}()
	SetDefault("bogus")
}

func TestCreateStyleMapErrors(t *testing.T) {
	for _, css := range []string{".cn-a { @apply p-1;", ".cn-a { @apply p-1; } }", "/* x", ".cn-a { content: \"x }"} {
		if _, err := CreateStyleMap(css); err == nil {
			t.Errorf("%q: no error", css)
		}
	}
}

func TestCreateStyleMapSelectors(t *testing.T) {
	got, err := CreateStyleMap(`
.style-x {
  .cn-a, .cn-b:hover { @apply p-1; }
  .cn-a:not(.cn-c) { @apply m-1; }          /* subject: the last cn- class, cn-c */
  .cn-d .cn-e, .x .y { @apply z-1; }
  @media (min-width: 1px) { .cn-f { @apply flex; } }
  .cn-g { @apply  ; color: red; }
  .cn-h\:i { @apply block; }
  .cn-j { &.cn-k { @apply hidden; } }
  .cn-a { @apply  p-2  ; @apply m-2; }
}`)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"cn-a": "p-2 m-2 p-1", "cn-b": "p-1", "cn-c": "m-1", "cn-e": "z-1",
		"cn-f": "flex", "cn-h:i": "block", "cn-k": "hidden",
	}
	if len(got) != len(want) {
		t.Errorf("got %v", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestSiteErrors(t *testing.T) {
	f := NewFile("x.tsx")
	s := f.Cn("cn-card p-1", Arg)
	mustPanic(t, "wrong value count", func() { s.Class(Nova) })
	mustPanic(t, "invalid style", func() { s.Class("bogus", "") })
	_ = s.Class(Nova, "")
	mustPanic(t, "CVA after resolution", func() { f.CVA("cn-button", nil, nil) })
	mustPanic(t, "bad argument", func() { f.Cn(42) })
	mustPanic(t, "bad PropCn argument", func() { f.PropCn(1.5) })
	mustPanic(t, "Enclosed after resolution", func() { s.Enclosed(1) })
	mustPanic(t, "Enclosed on an attribute", func() { f.Attr("p-1").Enclosed(1) })
	mustPanic(t, "negative Enclosed", func() { f.Cn("p-1").Enclosed(-1) })
}

// TestSortPasses shows why sites record their context: prettier sorts a
// className={cn(...)} literal twice (attribute, call), a className: cn(...)
// property literal once, and the passes can differ. Here the first pass
// trims a leading U+2003, which makes the class known to Tailwind, so the
// second pass moves it after the unknown "foo" (upstream's output for the
// fuzz case this comes from, testdata/ref/upstream.json.gz).
func TestSortPasses(t *testing.T) {
	const lit = "\u2003group-data-popup-open/navigation-menu-trigger:rotate-180 \n  foo"
	f := NewFile("passes.tsx")
	attr, prop, deep := f.Cn(lit, Arg), f.PropCn(lit, Arg), f.PropCn(lit).Enclosed(1)
	if got := attr.Args(Rhea)[0]; got != "foo group-data-popup-open/navigation-menu-trigger:rotate-180" {
		t.Errorf("attribute: %q", got)
	}
	if got := prop.Args(Rhea)[0]; got != "group-data-popup-open/navigation-menu-trigger:rotate-180 foo" {
		t.Errorf("property: %q", got)
	}
	if got := deep.Args(Rhea)[0]; got != attr.Args(Rhea)[0] {
		t.Errorf("enclosed property: %q", got)
	}
}

func mustPanic(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: no panic", what)
		}
	}()
	fn()
}

// TestFileClaims checks the per-file state: a token a cva() call styled is
// only removed from the file's other sites.
func TestFileClaims(t *testing.T) {
	f := NewFile("claims.tsx")
	f.CVA("cn-button flex", nil, nil)
	site := f.Cn("cn-button p-1", Arg)
	if got := site.Args(Nova); len(got) != 2 || got[0] != "p-1" {
		t.Errorf("claimed token: got %v", got)
	}
	other := Cn("cn-button p-1", Arg)
	if got := other.Args(Nova)[0].(string); !strings.HasPrefix(got, "rounded-lg border") {
		t.Errorf("standalone site: got %q", got)
	}
}

func TestConcurrentResolution(t *testing.T) {
	f := NewFile("concurrent.tsx")
	v := f.CVA("cn-badge", []Variant{{Name: "variant", Options: []Option{{Name: "default", Class: "cn-badge-variant-default"}}}}, nil)
	s := f.Cn("cn-badge-variant-outline flex", Arg)
	var wg sync.WaitGroup
	results := make([][]string, 16)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for _, st := range All {
				results[i] = append(results[i], s.Class(st, "x"), v.Classes(st, nil), Resolve(st, "cn-card"))
			}
		}(i)
	}
	wg.Wait()
	for i := 1; i < len(results); i++ {
		if strings.Join(results[i], "|") != strings.Join(results[0], "|") {
			t.Fatal("results differ between goroutines")
		}
	}
}

func BenchmarkSiteClass(b *testing.B) {
	s := Cn("cn-card group/card flex flex-col", Arg)
	s.Class(Nova, "")
	b.ReportAllocs()
	for b.Loop() {
		s.Class(Nova, "rounded-none")
	}
}

func BenchmarkResolve(b *testing.B) {
	Resolve(Nova, "cn-card-title")
	b.ReportAllocs()
	for b.Loop() {
		Resolve(Nova, "cn-card-title")
	}
}
