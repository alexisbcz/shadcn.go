package style

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestRTLClass checks RTLClass against upstream's ui-rtl output: every
// double-quoted class string of styles/radix-nova/ui/*.tsx maps to the
// string at the same position in ui-rtl/*.tsx (compared as class sets).
func TestRTLClass(t *testing.T) {
	dir := "../../ui/apps/v4/styles/radix-nova/"
	entries, err := os.ReadDir(dir + "ui")
	if err != nil {
		t.Skip("upstream clone not available")
	}
	str := regexp.MustCompile(`"([^"\n]*)"`)
	set := func(s string) string {
		f := strings.Fields(s)
		sort.Strings(f)
		return strings.Join(f, " ")
	}
	checked := 0
	for _, e := range entries {
		a, _ := os.ReadFile(dir + "ui/" + e.Name())
		b, err := os.ReadFile(dir + "ui-rtl/" + e.Name())
		if err != nil {
			continue
		}
		as, bs := str.FindAllStringSubmatch(string(a), -1), str.FindAllStringSubmatch(string(b), -1)
		if len(as) != len(bs) {
			continue
		}
		for i := range as {
			x, y := as[i][1], bs[i][1]
			if strings.Contains(x, "@/") || !strings.Contains(x, "-") {
				continue
			}
			if got := set(RTLClass(x)); got != set(y) {
				t.Errorf("%s: RTLClass(%q)\n got %q\nwant %q", e.Name(), x, got, set(y))
			}
			checked++
		}
	}
	t.Logf("%d class strings checked", checked)
}
