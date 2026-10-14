package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// ExampleDashboard serves /examples/dashboard (app/(app)/examples/dashboard). TODO.
func (s *Site) ExampleDashboard(r *http.Request) (h.Node, error) {
	return s.examplesLayout(r, "/examples/dashboard"), nil
}
