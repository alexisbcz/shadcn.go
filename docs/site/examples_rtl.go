package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// ExampleRTL serves /examples/rtl (app/(app)/examples/rtl). TODO.
func (s *Site) ExampleRTL(r *http.Request) (h.Node, error) {
	return s.examplesLayout(r, "/examples/rtl"), nil
}
