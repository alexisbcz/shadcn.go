// Package mdx implements the custom components the docs MDX content uses
// (the JSX names of ui/apps/v4/mdx-components.tsx: ComponentPreview,
// CodeTabs, ComponentSource, Callout, Tabs*, ...), as Go functions.
//
// Extension point: one file per component (or component family). Each file
// registers its functions from init with Register; Registry builds the
// content.Registry used to render a page. To implement a component, replace
// the placeholder body in its file; no other file needs to change.
//
// Components receive an *Env with the request being served, the docs page
// tree and the repository root, besides content's Context (the page being
// rendered), the props in source order and the rendered children. The
// wrappers of mdx-components.tsx (Steps, Step, Tab, the className defaults
// of Tabs*, Button and Accordion, LinkedCard -> Link, ComponentsList's
// props, Image's props) are applied by package content before the call.
package mdx

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sort"
	"sync"

	"shadcn.go/docs/content"
	"shadcn.go/docs/site/pagetree"
	h "shadcn.go/libhtml"
)

// Env is what components know about the site and the request.
type Env struct {
	// Request is the request being served (CSP nonce, query, ...).
	Request *http.Request
	// Tree is the docs page tree (source.pageTree).
	Tree *pagetree.Node
	// Root is the repository root (upstream sources live under
	// Root/ui/apps/v4).
	Root string
}

// Func is a docs component.
type Func func(env *Env, ctx *content.Context, props content.Props, children []h.Node) (h.Node, error)

var (
	mu    sync.RWMutex
	funcs = map[string]Func{}
)

// Register sets the implementation of a component name (a JSX name such as
// "Callout", or "<module>:<name>" for components imported by a page).
// Called from init; a later registration replaces an earlier one.
func Register(name string, f Func) {
	mu.Lock()
	defer mu.Unlock()
	funcs[name] = f
}

// Names returns the registered names, sorted.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(funcs))
	for n := range funcs {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Registry returns a content registry (content.NewRegistry's defaults plus
// every registered component) bound to env.
func Registry(env *Env) *content.Registry {
	reg := content.NewRegistry()
	mu.RLock()
	defer mu.RUnlock()
	for name, f := range funcs {
		f := f
		reg.Register(name, func(ctx *content.Context, props content.Props, children []h.Node) (h.Node, error) {
			return f(env, ctx, props, children)
		})
	}
	return reg
}

// placeholder renders content's generic placeholder: used by components not
// implemented yet.
func placeholder(name string) Func {
	return func(_ *Env, _ *content.Context, props content.Props, children []h.Node) (h.Node, error) {
		return content.Placeholder(name, props, children), nil
	}
}

// className returns the className prop.
func className(props content.Props) string { return props.String("className") }

func nodes(children []h.Node) []any {
	out := make([]any, len(children))
	for i, c := range children {
		out[i] = c
	}
	return out
}

// newID returns a unique element id (Radix useId style).
func newID(prefix string) string {
	var b [6]byte
	_, _ = rand.Read(b[:])
	return prefix + "-" + hex.EncodeToString(b[:])
}
