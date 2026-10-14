package content

import (
	"fmt"

	"shadcn.go/docs/mdx"
)

// RemarkTree parses one MDX page and runs the remark stages that precede
// fumadocs' remarkStructure (remarkMarkAndUnravel, remarkHeading), without
// the rehype stages and highlighting. The search index (package search)
// structures this tree.
func RemarkTree(src string) (*mdx.Node, *Frontmatter, error) {
	matter, body := SplitFrontmatter(src)
	fm, err := ParseFrontmatter(matter)
	if err != nil {
		return nil, nil, err
	}
	root, err := mdx.Parse(body)
	if err != nil {
		return nil, nil, fmt.Errorf("mdx: %w", err)
	}
	unravel(root)
	headingIDs(root)
	if fm != nil && fm.Title == "" {
		for _, c := range root.Children {
			if c.Type == "heading" && c.Depth == 1 {
				fm.Title = flatten(c)
				break
			}
		}
	}
	return root, fm, nil
}
