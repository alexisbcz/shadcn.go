package content

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestSourceInSync checks that the vendored content (source/, see gen/sync)
// is identical to ui/apps/v4/content/docs, when the upstream clone exists.
func TestSourceInSync(t *testing.T) {
	if _, err := os.Stat(contentRoot); err != nil {
		t.Skip("upstream clone not present")
	}
	upstream := map[string][]byte{}
	filepath.WalkDir(contentRoot, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(contentRoot, p)
			upstream[filepath.ToSlash(rel)], _ = os.ReadFile(p)
		}
		return nil
	})
	n := 0
	fs.WalkDir(Source(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		n++
		b, _ := fs.ReadFile(Source(), p)
		if want, ok := upstream[p]; !ok || !bytes.Equal(b, want) {
			t.Errorf("%s differs from upstream; run go run ./docs/content/gen/sync", p)
		}
		return nil
	})
	if n != len(upstream) {
		t.Errorf("vendored %d files, upstream has %d; run go run ./docs/content/gen/sync", n, len(upstream))
	}
}
