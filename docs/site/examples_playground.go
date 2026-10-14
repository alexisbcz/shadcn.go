package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// ExamplePlayground serves /examples/playground (app/(app)/examples/playground). TODO.
func (s *Site) ExamplePlayground(r *http.Request) (h.Node, error) {
	return s.examplesLayout(r, "/examples/playground"), nil
}
