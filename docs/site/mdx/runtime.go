package mdx

import (
	_ "embed"

	"shadcn.go/components/ui/runtime"
)

//go:embed docs-mdx.js
var docsScript string

func init() { runtime.Module("docs-mdx", docsScript) }
