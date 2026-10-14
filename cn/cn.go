// Package cn is a Go port of the cn npm package (cn@0.2.2), the compiled
// clsx + tailwind-merge replacement that shadcn/ui's cn() helper calls:
//
//	cn.Cn("px-2 py-1 bg-red-500", "p-3 bg-blue-500") // "p-3 bg-blue-500"
//	cn.Cn("hover:bg-accent", "bg-primary")           // both kept: different variants
//
// It runs the same engine algorithm on the same generated tables
// (tables_gen.go, produced by gen/tables.mjs from the package's
// dist/tables.js), and is tested to produce byte-identical results to the
// JavaScript package (see cn_test.go and fuzz_test.go).
//
// All functions are safe for concurrent use.
package cn

//go:generate node gen/tables.mjs

import (
	"strings"
	"sync"
)

var std = newEngine()

// Cn joins the non-empty class lists with spaces and merges the result,
// later classes overriding conflicting earlier ones. It is upstream's
// cn(...) called with string arguments (clsx skips empty strings).
func Cn(classes ...string) string {
	// Skip empty arguments like clsx; keep the first two non-empty ones
	// apart so the common cn(base, className) call allocates nothing on a
	// cache hit.
	var k cacheKey
	n := 0
	for _, c := range classes {
		if c == "" {
			continue
		}
		switch n {
		case 0:
			k.a = c
		case 1:
			k.b = c
		default:
			k.b += " " + c
		}
		n++
	}
	switch n {
	case 0:
		return ""
	case 1:
		return Merge(k.a)
	}
	if v, ok := cache.get(k); ok {
		return v
	}
	out := std.mergeClassList(k.a + " " + k.b)
	cache.put(k, out)
	return out
}

// Merge merges one already-joined class string (the engine's mergeString).
func Merge(classes string) string {
	k := cacheKey{a: classes}
	if v, ok := cache.get(k); ok {
		return v
	}
	out := std.mergeClassList(classes)
	cache.put(k, out)
	return out
}

// MergeUncached is Merge without the result cache (mergeUncached).
func MergeUncached(classes string) string { return std.mergeClassList(classes) }

// Join joins the non-empty class lists with spaces without merging (twJoin
// and clsx for strings).
func Join(classes ...string) string {
	var b strings.Builder
	for _, c := range classes {
		if c == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(c)
	}
	return b.String()
}

// --- cache ----------------------------------------------------------------------

// Component class lists repeat constantly, so results are memoized. The cache
// is bounded: when full it is simply cleared. Merging is a pure function of
// the joined string, so caching never changes a result.
type cacheKey struct{ a, b string }

type boundedCache struct {
	mu sync.RWMutex
	m  map[cacheKey]string
}

const cacheSize = 8192

var cache = &boundedCache{m: make(map[cacheKey]string, 256)}

func (c *boundedCache) get(k cacheKey) (string, bool) {
	c.mu.RLock()
	v, ok := c.m[k]
	c.mu.RUnlock()
	return v, ok
}

func (c *boundedCache) put(k cacheKey, v string) {
	c.mu.Lock()
	if len(c.m) >= cacheSize {
		clear(c.m)
	}
	c.m[k] = v
	c.mu.Unlock()
}
