package style

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"strings"
	"sync"
)

// order.txt.gz is Tailwind's class order for every class the style system
// has to sort: line 1 is a header, then the classes Tailwind knows in
// getClassOrder order, then "# unknown" and the classes it does not know
// (getClassOrder gives them null). gen/main.go produces it with the vendored
// Tailwind compiler running in goja, from upstream's app/globals.css;
// testdata/refgen/order.mjs checks it against Node.js.
//
//go:embed order.txt.gz
var orderData []byte

var (
	orderOnce  sync.Once
	orderRanks map[string]int32 // known classes: rank >= 0; unknown: -1

	missMu sync.Mutex
	misses map[string]bool
)

func loadOrder() {
	orderRanks = map[string]int32{}
	zr, err := gzip.NewReader(bytes.NewReader(orderData))
	if err != nil {
		panic("style: corrupt order.txt.gz: " + err.Error())
	}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	var rank int32
	unknown := false
	for line := 0; sc.Scan(); line++ {
		t := sc.Text()
		switch {
		case line == 0:
			continue
		case t == "# unknown":
			unknown = true
		case unknown:
			orderRanks[t] = -1
		default:
			orderRanks[t] = rank
			rank++
		}
	}
	if err := sc.Err(); err != nil {
		panic("style: corrupt order.txt.gz: " + err.Error())
	}
}

// classRank returns the class's position in Tailwind's order, or -1 when
// Tailwind does not know the class (getClassOrder returns null). A class
// missing from the table is treated as unknown and recorded, see
// UnorderedClasses.
func classRank(c string) int32 {
	orderOnce.Do(loadOrder)
	if r, ok := orderRanks[c]; ok {
		return r
	}
	missMu.Lock()
	if misses == nil {
		misses = map[string]bool{}
	}
	misses[c] = true
	missMu.Unlock()
	return -1
}

// UnorderedClasses returns the classes sorted so far that are missing from
// the embedded class-order table, so their position may differ from
// upstream's. It is empty when every class came from upstream's sources, its
// style files or the Go sources the generator scanned; when it is not,
// re-run `go generate ./components/style` to add them.
func UnorderedClasses() []string {
	missMu.Lock()
	defer missMu.Unlock()
	out := make([]string, 0, len(misses))
	for c := range misses {
		out = append(out, c)
	}
	sortStrings(out)
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && strings.Compare(s[j-1], s[j]) > 0; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
