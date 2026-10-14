package style

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestGeneratedFilesCurrent runs `go run ./gen -check`: the embedded style
// files must equal upstream's, and the class-order table must equal what
// Tailwind's getClassOrder gives (computed again in goja) and list every
// class of the style files and the test data. It needs the upstream clone
// in ./ui and is skipped without it.
func TestGeneratedFilesCurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	if _, err := os.Stat("../../ui/apps/v4/registry/styles/style-nova.css"); err != nil {
		t.Skip("upstream clone not found in ./ui")
	}
	out, err := exec.Command("go", "run", "./gen", "-check", "-root", "../..").CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	t.Log(strings.TrimSpace(string(out)))
}

// TestStyleClassesOrdered checks that every class the style files apply has
// a place in the order table, so sorting never meets an unlisted class
// coming from a style.
func TestStyleClassesOrdered(t *testing.T) {
	for _, st := range All {
		for tok, classes := range st.Map() {
			for _, c := range strings.Fields(classes) {
				orderOnce.Do(loadOrder)
				if _, ok := orderRanks[c]; !ok {
					t.Errorf("%s: %s: class %q missing from order.txt.gz", st, tok, c)
				}
			}
		}
	}
}
