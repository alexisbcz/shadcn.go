package content

import (
	"embed"
	"io/fs"
)

// source holds the docs content, vendored from ui/apps/v4/content/docs by
// gen/sync: the 316 MDX pages and the meta.json files that order the page
// tree.
//
//go:embed all:source
var source embed.FS

// Source returns the vendored docs content (paths like
// "components/radix/button.mdx", "(root)/meta.json").
func Source() fs.FS {
	sub, err := fs.Sub(source, "source")
	if err != nil {
		panic(err)
	}
	return sub
}

// LoadDocs parses every page of the vendored docs content.
func LoadDocs(opts Options) ([]*Page, error) { return Load(Source(), ".", opts) }
