package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// Typeset serves /typeset. TODO.
func (s *Site) Typeset(r *http.Request) (h.Node, error) {
	return s.Document(r, Meta{}, s.AppLayout(r)), nil
}
