package style

import (
	"fmt"
	"strings"
	"testing"
)

var bases = []string{"radix", "base", "aria"}

// TestDifferential replays upstream's build on every base source of the 24
// generated directories (ui/apps/v4/styles/<base>-<style>/ui and the
// examples, blocks, hooks and lib files the build transforms with them) and
// requires every class-bearing string of the generated source exactly: cva
// base strings and variant values, cn() arguments, className literals and
// every other string inside them. The data is extracted from upstream with
// the TypeScript compiler by testdata/refgen/extract.mjs.
func TestDifferential(t *testing.T) {
	total, totalFiles := 0, 0
	var summary []string
	for _, b := range bases {
		for _, st := range []Style{Nova, Vega, Lyra, Maia, Mira, Luma, Sera, Rhea} {
			dir := b + "-" + string(st)
			ref, err := loadRef(dir)
			if err != nil {
				t.Fatal(err)
			}
			compared, mismatches, uiStrings := 0, 0, 0
			for _, f := range ref.Files {
				base, gen := decodeRoots(f.Base), decodeRoots(f.Gen)
				if len(base) != len(gen) {
					t.Errorf("%s/%s: %d class roots, generated has %d", dir, f.Path, len(base), len(gen))
					mismatches++
					continue
				}
				transformStyleMap(base, st.Map())
				sortTree(base)
				for i := range base {
					n := 0
					if d := compareTrees(base[i], gen[i], "", &n); d != "" {
						mismatches++
						if mismatches <= 5 {
							t.Errorf("%s/%s:%d: %s", dir, f.Path, f.Base[i].Line, d)
						}
					}
					compared += n
					if strings.HasPrefix(f.Path, "ui/") {
						uiStrings += n
					}
				}
			}
			total += compared
			totalFiles += len(ref.Files)
			summary = append(summary, fmt.Sprintf("%-11s %4d files %6d strings (ui %5d) %d mismatches", dir, len(ref.Files), compared, uiStrings, mismatches))
		}
	}
	for _, s := range summary {
		t.Log(s)
	}
	t.Logf("total: %d files, %d strings compared", totalFiles, total)
	if u := UnorderedClasses(); len(u) > 0 {
		t.Errorf("classes missing from the order table: %v", u)
	}
}

// TestDifferentialDetects checks that the differential data discriminates:
// skipping the token transform, skipping prettier's sorting, or using
// another style's map must each produce mismatches.
func TestDifferentialDetects(t *testing.T) {
	ref, err := loadRef("radix-nova")
	if err != nil {
		t.Fatal(err)
	}
	run := func(transform, sort bool, m map[string]string) (mismatches int) {
		for _, f := range ref.Files {
			base, gen := decodeRoots(f.Base), decodeRoots(f.Gen)
			if transform {
				transformStyleMap(base, m)
			}
			if sort {
				sortTree(base)
			}
			for i := range base {
				n := 0
				if compareTrees(base[i], gen[i], "", &n) != "" {
					mismatches++
				}
			}
		}
		return mismatches
	}
	if n := run(true, true, Nova.Map()); n != 0 {
		t.Fatalf("baseline: %d mismatches", n)
	}
	for _, c := range []struct {
		name            string
		transform, sort bool
		m               map[string]string
	}{
		{"no transform", false, true, Nova.Map()},
		{"no sort", true, false, Nova.Map()},
		{"vega map", true, true, Vega.Map()},
	} {
		n := run(c.transform, c.sort, c.m)
		t.Logf("%s: %d mismatching roots", c.name, n)
		if n == 0 {
			t.Errorf("%s: no mismatches; the data does not test it", c.name)
		}
	}
}

// TestBaseSourcesSorted shows that upstream's base sources are already in
// prettier's class order: sorting them changes no string. So a literal the
// transform does not touch (one nested in a condition, say) is the same in
// every generated style, and the API passes it through unchanged (Arg).
func TestBaseSourcesSorted(t *testing.T) {
	changed, total := 0, 0
	for _, b := range bases {
		ref, err := loadRef(b + "-nova")
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range ref.Files {
			roots, orig := decodeRoots(f.Base), decodeRoots(f.Base)
			sortTree(roots)
			for i := range roots {
				n := 0
				if d := compareTrees(roots[i], orig[i], "", &n); d != "" {
					changed++
					if changed <= 5 {
						t.Errorf("%s/%s:%d: %s", b, f.Path, f.Base[i].Line, d)
					}
				}
				total += n
			}
		}
	}
	t.Logf("%d strings of the 3 bases' sources: %d change when sorted", total, changed)
}
