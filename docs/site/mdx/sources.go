package mdx

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"strings"
	"sync"

	"shadcn.go/docs/hast"
	"shadcn.go/docs/highlight"
	h "shadcn.go/libhtml"
)

// sources.json.gz holds the code the docs show through ComponentPreview's
// code panel and ComponentSource, formatted by upstream's own pipeline
// (lib/registry.ts + lib/format-code.ts: transformStyle, transformIcons,
// transformMenu, transformRender). Regenerate it with gen/sources.mts.
//
//go:embed sources.json.gz
var sourcesGz []byte

var (
	sourcesOnce sync.Once
	sources     map[string]string
)

func loadSources() map[string]string {
	sourcesOnce.Do(func() {
		sources = map[string]string{}
		r, err := gzip.NewReader(bytes.NewReader(sourcesGz))
		if err != nil {
			return
		}
		_ = json.NewDecoder(r).Decode(&sources)
	})
	return sources
}

// Source returns the formatted source shown for name in styleName (the
// example demo, or else the registry item's first file), like
// ComponentSource's name prop.
func Source(styleName, name string) (string, bool) {
	if styleName == "" {
		styleName = "base-nova"
	}
	s, ok := loadSources()[styleName+"\x00"+name]
	return s, ok
}

// SourceFile returns the formatted source of a file under ui/apps/v4 (the
// src prop of ComponentSource).
func SourceFile(styleName, src string) (string, bool) {
	if styleName == "" {
		styleName = "base-nova"
	}
	s, ok := loadSources()["src\x00"+styleName+"\x00"+src]
	return s, ok
}

// highlightCode is lib/highlight-code.ts's highlightCode: the HTML of the
// highlighted <pre>, inserted as is (upstream: dangerouslySetInnerHTML).
// The HTML is produced by our own highlighter from repository sources.
// HighlightCode is highlightCode (lib/highlight-code.ts) for use outside
// MDX components, e.g. the block viewer's code view.
func HighlightCode(code, lang string) h.Node { return highlightCode(code, lang) }

func highlightCode(code, lang string) h.Node {
	hl, err := highlight.Default()
	if err == nil {
		if pre, err := hl.CodeTree(code, lang); err == nil {
			return h.UnsafeHTML(hast.ToHTML(pre))
		}
	}
	return h.Pre(h.Code(code))
}

// firstLines keeps the first n lines of code (ComponentSource's maxLines).
func firstLines(code string, n int) string {
	lines := strings.Split(code, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
