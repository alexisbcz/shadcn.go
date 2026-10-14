package site

import (
	"net/http"

	h "shadcn.go/libhtml"
)

// Sera serves /sera. TODO.
func (s *Site) Sera(r *http.Request) (h.Node, error) {
	return s.Document(r, Meta{}, s.AppLayout(r)), nil
}
