package cn

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/dop251/goja"
)

// The fuzz oracle is the upstream JavaScript itself: testdata/upstream holds
// an unmodified copy of cn@0.2.2's dist/{tables,engine,index}.js (MIT,
// testdata/upstream/LICENSE), run in goja. FuzzCn compares the Go port with
// it on arbitrary inputs, seeded with the recorded corpus.

var importExport = regexp.MustCompile(`(?m)^(import|export) .*$`)

// jsOracle runs upstream cn in a goja VM. It is not safe for concurrent use;
// callers hold mu.
type jsOracle struct {
	mu       sync.Mutex
	cn       goja.Callable // cn(...args), the cached clsx wrapper
	uncached goja.Callable // engine.mergeUncached(string)
	vm       *goja.Runtime
}

var loadOracle = sync.OnceValues(func() (*jsOracle, error) {
	var src strings.Builder
	for _, f := range []string{"tables.js", "engine.js", "index.js"} {
		b, err := os.ReadFile("testdata/upstream/" + f)
		if err != nil {
			return nil, err
		}
		src.WriteString(importExport.ReplaceAllString(string(b), ""))
		src.WriteString("\n")
	}
	src.WriteString("var __uncached = createEngine(tables_generated_default, undefined, { cacheSize: 0 }).mergeUncached;\n")
	vm := goja.New()
	if _, err := vm.RunString("(function(){\n" + src.String() + "\nglobalThis.__cn = cn; globalThis.__uncached = __uncached;\n})()"); err != nil {
		return nil, err
	}
	o := &jsOracle{vm: vm}
	var ok bool
	if o.cn, ok = goja.AssertFunction(vm.Get("__cn")); !ok {
		return nil, os.ErrInvalid
	}
	if o.uncached, ok = goja.AssertFunction(vm.Get("__uncached")); !ok {
		return nil, os.ErrInvalid
	}
	return o, nil
})

func (o *jsOracle) call(args []string) (cached, uncached string, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	vals := make([]goja.Value, len(args))
	for i, a := range args {
		vals[i] = o.vm.ToValue(a)
	}
	v, err := o.cn(goja.Undefined(), vals...)
	if err != nil {
		return "", "", err
	}
	u, err := o.uncached(goja.Undefined(), o.vm.ToValue(Join(args...)))
	if err != nil {
		return "", "", err
	}
	return v.String(), u.String(), nil
}

func mustOracle(t testing.TB) *jsOracle {
	t.Helper()
	o, err := loadOracle()
	if err != nil {
		t.Fatal(err)
	}
	return o
}

// TestUpstreamCopy checks the vendored oracle is the upstream package the
// tables were generated from.
func TestUpstreamCopy(t *testing.T) {
	b, err := os.ReadFile("testdata/upstream/tables.js")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(b)
	if got := hex.EncodeToString(sum[:]); got != tablesSHA256 {
		t.Fatalf("testdata/upstream/tables.js sha256 %s, tables_gen.go %s", got, tablesSHA256)
	}
	for _, f := range []string{"engine.js", "index.js", "tables.js"} {
		up, err := os.ReadFile("../ui/apps/v4/node_modules/cn/dist/" + f)
		if err != nil {
			t.Skip("upstream checkout not available:", err)
		}
		mine, _ := os.ReadFile("testdata/upstream/" + f)
		if string(up) != string(mine) {
			t.Errorf("testdata/upstream/%s differs from ui/apps/v4/node_modules/cn/dist/%s", f, f)
		}
	}
}

// TestOracleAgreesWithCorpus checks the goja oracle reproduces node's
// recorded results on a slice of the corpus, so the fuzz oracle is trusted.
func TestOracleAgreesWithCorpus(t *testing.T) {
	c := mustCorpus(t)
	o := mustOracle(t)
	// Every case is covered by TestCorpus (Go = node) and the FuzzCn seeds
	// (Go = goja); a sample suffices here. CN_FULL=1 checks them all.
	step := 11
	if os.Getenv("CN_FULL") != "" {
		step = 1
	}
	for i := 0; i < len(c.Cases); i += step {
		tc := c.Cases[i]
		got, unc, err := o.call(tc[1:])
		if err != nil {
			t.Fatal(err)
		}
		if got != tc[0] || unc != tc[0] {
			t.Fatalf("goja oracle differs from node on %q: %q / %q, want %q", tc[1:], got, unc, tc[0])
		}
	}
}

// FuzzCn requires the Go port to agree with the upstream JavaScript on
// arbitrary argument lists. Seeds: the recorded corpus.
//
//	go test ./cn -run '^$' -fuzz FuzzCn -fuzztime 5m
func FuzzCn(f *testing.F) {
	c := mustCorpus(f)
	// When fuzzing (or with CN_FULL=1) every corpus case is a seed; a plain
	// go test runs a sample of them against the oracle.
	all := os.Getenv("CN_FULL") != ""
	if fl := flag.Lookup("test.fuzz"); fl != nil && fl.Value.String() != "" {
		all = true
	}
	for i, tc := range c.Cases {
		if !all && i%11 != 0 {
			continue
		}
		args := tc[1:]
		var a, b, rest string
		a = args[0]
		if len(args) > 1 {
			b = args[1]
		}
		if len(args) > 2 {
			rest = strings.Join(args[2:], "\x00")
		}
		f.Add(a, b, rest)
	}
	o := mustOracle(f)
	f.Fuzz(func(t *testing.T, a, b, rest string) {
		args := []string{a, b}
		if rest != "" {
			args = append(args, strings.Split(rest, "\x00")...)
		}
		got := Cn(args...)
		unc := MergeUncached(Join(args...))
		if got != unc {
			t.Fatalf("cached %q != uncached %q for %q", got, unc, args)
		}
		for _, s := range args {
			if !utf8.ValidString(s) {
				// JavaScript strings cannot hold invalid UTF-8; the port
				// only has to not panic.
				return
			}
		}
		want, wantUnc, err := o.call(args)
		if err != nil {
			t.Fatal(err)
		}
		if want != wantUnc {
			t.Fatalf("upstream cached %q != uncached %q", want, wantUnc)
		}
		if got != want {
			t.Fatalf("Cn(%q)\n got  %q\n want %q", args, got, want)
		}
	})
}
