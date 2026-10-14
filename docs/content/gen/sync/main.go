// Command sync vendors the docs content (ui/apps/v4/content/docs: the MDX
// pages and the meta.json files of the page tree) into docs/content/source,
// which package content embeds.
//
// Usage, from the repository root: go run ./docs/content/gen/sync
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	src := filepath.Join("ui", "apps", "v4", "content", "docs")
	dst := filepath.Join("docs", "content", "source")
	if err := os.RemoveAll(dst); err != nil {
		panic(err)
	}
	n := 0
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out := filepath.Join(dst, rel)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		n++
		return os.WriteFile(out, b, 0o644)
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d files vendored into %s\n", n, dst)
}
