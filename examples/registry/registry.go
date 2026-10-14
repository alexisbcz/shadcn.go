// Package registry indexes the Go ports of the upstream shadcn/ui examples
// (ui/apps/v4/registry/new-york-v4/examples). Each examples/<component>
// package registers its examples from init; examples/all links them all.
package registry

import (
	"fmt"
	"net/http"
	"slices"
	"sync"

	h "shadcn.go/libhtml"
)

// Example is one ported example. Name is the upstream file name without
// extension (e.g. "button-demo"); Render builds the same markup the React
// example renders.
type Example struct {
	Name      string
	Component string // upstream ui component the example demonstrates
	Render    func(r *http.Request) h.Node
	// Extra, if set, renders page-level companions outside the preview,
	// such as sonner's <Toaster />, which the upstream docs mount in their
	// layout.
	Extra func(r *http.Request) h.Node
}

var (
	mu       sync.RWMutex
	examples = map[string]Example{}
)

// Register adds examples. Names must be unique.
func Register(list ...Example) {
	mu.Lock()
	defer mu.Unlock()
	for _, e := range list {
		if _, dup := examples[e.Name]; dup {
			panic(fmt.Sprintf("registry: example %q registered twice", e.Name))
		}
		examples[e.Name] = e
	}
}

// Get returns an example by name.
func Get(name string) (Example, bool) {
	mu.RLock()
	defer mu.RUnlock()
	e, ok := examples[name]
	return e, ok
}

// All returns every example, sorted by name.
func All() []Example {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Example, 0, len(examples))
	for _, e := range examples {
		out = append(out, e)
	}
	slices.SortFunc(out, func(a, b Example) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
	return out
}

// Static adapts an example that does not need the request.
func Static(fn func() h.Node) func(*http.Request) h.Node {
	return func(*http.Request) h.Node { return fn() }
}
