// Package search is the documentation search of the Go docs site: the
// index behind app/api/search/route.ts (fumadocs-core's createFromSource,
// an Orama full-text index of every page's structured data) and the
// structured data itself (fumadocs' remarkStructure).
//
// Both run the upstream JavaScript, fumadocs-core 16.10.5 with Orama
// 3.1.18 and mdast-util-to-markdown, bundled unmodified by gen/build.mjs
// (dist/search.js, licenses in dist/LICENSES) and run in goja. The Go side
// supplies the mdast (package mdx, with content.RemarkTree's remark
// stages), so results match the reference's /api/search.
package search

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/dop251/goja"

	"shadcn.go/docs/hast"
	"shadcn.go/docs/mdx"
)

//go:embed dist/search.js
var searchJS string

// Index is one page of the search index (fumadocs' AdvancedIndex).
type Index struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description,omitempty"`
	URL            string          `json:"url"`
	Breadcrumbs    []string        `json:"breadcrumbs,omitempty"`
	StructuredData json.RawMessage `json:"structuredData"`
}

// Engine runs the bundled search code. It is safe for concurrent use
// (calls are serialized: goja runtimes are single-threaded).
type Engine struct {
	mu        sync.Mutex
	vm        *goja.Runtime
	structure goja.Callable
	index     goja.Callable
	search    goja.Callable
	indexed   bool
}

var (
	progOnce sync.Once
	prog     *goja.Program
	progErr  error
)

// New loads the bundle in a fresh runtime.
func New() (*Engine, error) {
	progOnce.Do(func() { prog, progErr = goja.Compile("search.js", searchJS, false) })
	if progErr != nil {
		return nil, progErr
	}
	vm := goja.New()
	if _, err := vm.RunProgram(prog); err != nil {
		return nil, fmt.Errorf("search: loading bundle: %w", err)
	}
	api := vm.Get("DocsSearch").ToObject(vm)
	e := &Engine{vm: vm}
	var ok bool
	if e.structure, ok = goja.AssertFunction(api.Get("structure")); !ok {
		return nil, errors.New("search: structure missing")
	}
	e.index, _ = goja.AssertFunction(api.Get("index"))
	e.search, _ = goja.AssertFunction(api.Get("search"))
	return e, nil
}

// Structure returns fumadocs' structured data ({headings, contents}) of a
// page's remark tree (content.RemarkTree).
func (e *Engine) Structure(root *mdx.Node) (json.RawMessage, error) {
	b, err := json.Marshal(Mdast(root))
	if err != nil {
		return nil, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	v, err := e.structure(goja.Undefined(), e.vm.ToValue(string(b)))
	if err != nil {
		return nil, err
	}
	return json.RawMessage(v.String()), nil
}

// SetIndex (re)builds the Orama database from the pages, in the order of
// fumadocs' loader.getPages().
func (e *Engine) SetIndex(pages []Index) error {
	b, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, err := e.index(goja.Undefined(), e.vm.ToValue(string(b))); err != nil {
		return err
	}
	e.indexed = true
	return nil
}

// Search returns the JSON body /api/search?query=... answers ("[]" before
// SetIndex).
func (e *Engine) Search(query string) ([]byte, error) {
	if query == "" {
		return []byte("[]"), nil
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.indexed {
		return []byte("[]"), nil
	}
	v, err := e.search(goja.Undefined(), e.vm.ToValue(query))
	if err != nil {
		return nil, err
	}
	p, ok := v.Export().(*goja.Promise)
	if !ok {
		return nil, errors.New("search: search did not return a promise")
	}
	switch p.State() {
	case goja.PromiseStateFulfilled:
		return []byte(p.Result().String()), nil
	case goja.PromiseStateRejected:
		return nil, fmt.Errorf("search: %v", p.Result())
	}
	return nil, errors.New("search: search did not settle")
}

// Mdast converts a parsed tree to mdast's JSON shape (positions omitted).
// Heading ids (remarkHeading) go to data.hProperties.id.
func Mdast(n *mdx.Node) map[string]any {
	o := map[string]any{"type": n.Type}
	switch n.Type {
	case "text", "inlineCode", "mdxFlowExpression", "mdxTextExpression", "mdxjsEsm", "code", "html":
		o["value"] = n.Value
	}
	switch n.Type {
	case "heading":
		o["depth"] = n.Depth
		if id, ok := n.Data["id"].(string); ok {
			o["data"] = map[string]any{"hProperties": map[string]any{"id": id}}
		}
	case "list":
		o["ordered"] = n.Ordered
		if n.Ordered {
			o["start"] = n.Start
		} else {
			o["start"] = nil
		}
		o["spread"] = n.Spread
	case "listItem":
		o["spread"] = n.Spread
		if n.Checked != nil {
			o["checked"] = *n.Checked
		} else {
			o["checked"] = nil
		}
	case "code":
		o["lang"], o["meta"] = nil, nil
		if n.HasLang {
			o["lang"] = n.Lang
		}
		if n.HasMeta {
			o["meta"] = n.Meta
		}
	case "link", "image":
		o["url"] = n.URL
		o["title"] = nil
		if n.HasTitle {
			o["title"] = n.Title
		}
		if n.Type == "image" {
			o["alt"] = n.Alt
		}
	case "table":
		a := []any{}
		for _, x := range n.Align {
			if x == "" {
				a = append(a, nil)
			} else {
				a = append(a, x)
			}
		}
		o["align"] = a
	case "footnoteReference", "footnoteDefinition":
		o["identifier"] = n.Identifier
		o["label"] = n.Label
	case "mdxJsxFlowElement", "mdxJsxTextElement":
		if n.Name != "" {
			o["name"] = n.Name
		} else {
			o["name"] = nil
		}
		attrs := []any{}
		for _, a := range n.Attrs {
			m := map[string]any{"type": "mdxJsxAttribute", "name": a.Name}
			switch a.Kind {
			case hast.AttrString:
				m["value"] = a.Value
			case hast.AttrBool:
				m["value"] = nil
			case hast.AttrExpression:
				m["value"] = map[string]any{"type": "mdxJsxAttributeValueExpression", "value": a.Value}
			}
			attrs = append(attrs, m)
		}
		o["attributes"] = attrs
	}
	switch n.Type {
	case "text", "inlineCode", "code", "thematicBreak", "mdxFlowExpression", "mdxTextExpression", "mdxjsEsm", "image", "footnoteReference", "break", "html":
	default:
		kids := make([]any, 0, len(n.Children))
		for _, c := range n.Children {
			kids = append(kids, Mdast(c))
		}
		o["children"] = kids
	}
	return o
}
