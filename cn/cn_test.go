package cn

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"testing"
)

type corpus struct {
	CN     string         `json:"cn"`
	Counts map[string]int `json:"counts"`
	Cases  [][]string     `json:"cases"`
}

var loadCorpus = sync.OnceValues(func() (*corpus, error) {
	f, err := os.Open("testdata/corpus.json.gz")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	var c corpus
	if err := json.NewDecoder(zr).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
})

func mustCorpus(t testing.TB) *corpus {
	t.Helper()
	c, err := loadCorpus()
	if err != nil {
		t.Fatal(err)
	}
	if c.CN != Version {
		t.Fatalf("corpus was computed with cn@%s, tables are cn@%s: run node cn/gen/corpus.mjs", c.CN, Version)
	}
	return c
}

// TestCorpus requires byte-identical results to the JavaScript package on
// every case of the recorded corpus (cn/gen/corpus.mjs): every cn() call of
// the upstream registry with its argument combinations, registry class
// strings merged pairwise, and a randomized corpus.
func TestCorpus(t *testing.T) {
	c := mustCorpus(t)
	if len(c.Cases) < 50000 {
		t.Fatalf("corpus has only %d cases", len(c.Cases))
	}
	fails := 0
	for _, tc := range c.Cases {
		want, args := tc[0], tc[1:]
		got := Cn(args...)
		unc := MergeUncached(Join(args...))
		if got != want || unc != want {
			fails++
			if fails <= 20 {
				t.Errorf("Cn(%q)\n got  %q\n uncached %q\n want %q", args, got, unc, want)
			}
		}
	}
	if fails > 0 {
		t.Fatalf("%d of %d cases differ", fails, len(c.Cases))
	}
	t.Logf("%d cases identical (%v)", len(c.Cases), c.Counts)
}

// TestCorpusConcurrent runs the corpus from several goroutines at once to
// check the per-call scratch state and the cache are race free.
func TestCorpusConcurrent(t *testing.T) {
	c := mustCorpus(t)
	var wg sync.WaitGroup
	errs := make(chan string, 8)
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := g; i < len(c.Cases); i += 8 * 7 {
				tc := c.Cases[i]
				if got := Cn(tc[1:]...); got != tc[0] {
					errs <- got
					return
				}
			}
		}(g)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Errorf("concurrent mismatch: %q", e)
	}
}

// TestTablesUpToDate checks tables_gen.go was generated from the cn package
// the upstream app depends on.
func TestTablesUpToDate(t *testing.T) {
	b, err := os.ReadFile("../ui/apps/v4/node_modules/cn/dist/tables.js")
	if err != nil {
		t.Skip("upstream checkout not available:", err)
	}
	sum := sha256.Sum256(b)
	if got := hex.EncodeToString(sum[:]); got != tablesSHA256 {
		t.Fatalf("tables_gen.go is stale (sha256 %s, upstream %s): run go generate ./cn", tablesSHA256, got)
	}
}

func TestBasics(t *testing.T) {
	for _, tc := range []struct {
		in   []string
		want string
	}{
		{nil, ""},
		{[]string{"", ""}, ""},
		{[]string{"  p-2  "}, "p-2"},
		{[]string{"px-2 py-1 bg-red-500", "p-3 bg-blue-500"}, "p-3 bg-blue-500"},
		{[]string{"hover:bg-accent", "bg-primary"}, "hover:bg-accent bg-primary"},
		{[]string{"rounded-md", "", "rounded-full"}, "rounded-full"},
		{[]string{"p-2\tp-4\n"}, "p-4"},
		{[]string{"a  b"}, "a b"},
		{[]string{"hover:focus:p-2 focus:hover:p-4"}, "focus:hover:p-4"},
		{[]string{"!p-2 p-4!"}, "p-4!"},
		{[]string{"text-sm/6 leading-7"}, "text-sm/6 leading-7"},
		{[]string{"leading-7 text-sm/6"}, "text-sm/6"},
	} {
		if got := Cn(tc.in...); got != tc.want {
			t.Errorf("Cn(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestJSNumber(t *testing.T) {
	for _, tc := range []struct {
		in       string
		num, int bool
	}{
		{"1", true, true}, {"1.5", true, false}, {".5", true, false}, {"5.", true, true},
		{"1e3", true, true}, {"1e-3", true, false}, {"+5", true, true}, {"-5", true, true},
		{"0x10", true, true}, {"0X1f", true, true}, {"0b101", true, true}, {"0o17", true, true},
		{"-0x10", false, false}, {"0x", false, false}, {"0b2", false, false},
		{"Infinity", true, false}, {"-Infinity", true, false}, {"infinity", false, false},
		{"NaN", false, false}, {"1_000", false, false}, {"1e", false, false}, {".", false, false},
		{"\u00a01\u3000", true, true}, {"\ufeff2", true, true}, {"\u200b1", false, false},
		{"1e400", true, false}, {"", false, false}, {"e5", false, false}, {"07", true, true},
		{"0.0", true, true}, {"-0", true, true}, {"1/2", false, false},
	} {
		if got := isNumStr(tc.in); got != tc.num {
			t.Errorf("isNumStr(%q) = %v, want %v", tc.in, got, tc.num)
		}
		if got := isIntStr(tc.in); got != tc.int {
			t.Errorf("isIntStr(%q) = %v, want %v", tc.in, got, tc.int)
		}
	}
}

func TestJoin(t *testing.T) {
	if got := Join("a", "", "b c", ""); got != "a b c" {
		t.Fatal(got)
	}
	if !strings.Contains(Version, ".") {
		t.Fatal(Version)
	}
}

func BenchmarkCnCached(b *testing.B) {
	base := "inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2 has-[>svg]:px-3"
	for b.Loop() {
		Cn(base, "w-full px-2")
	}
}

func BenchmarkMergeUncached(b *testing.B) {
	base := "inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2 has-[>svg]:px-3 w-full px-2"
	for b.Loop() {
		MergeUncached(base)
	}
}
