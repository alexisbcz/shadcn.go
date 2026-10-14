package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// ExampleAuthentication serves /examples/authentication (app/(app)/examples/authentication). TODO.
func (s *Site) ExampleAuthentication(r *http.Request) (h.Node, error) {
	return s.examplesLayout(r, "/examples/authentication"), nil
}
