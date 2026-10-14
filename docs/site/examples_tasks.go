package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// ExampleTasks serves /examples/tasks (app/(app)/examples/tasks). TODO.
func (s *Site) ExampleTasks(r *http.Request) (h.Node, error) {
	return s.examplesLayout(r, "/examples/tasks"), nil
}
