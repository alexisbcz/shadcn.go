// Package uid generates the element ids Radix creates with useId: unique
// per page, safe to call concurrently from any request.
package uid

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

var enc = base32.StdEncoding.WithPadding(base32.NoPadding)

// New returns a random id like "radix-k3j9f2a1qd", with the given prefix
// ("radix" when empty). 50 random bits make collisions within a page
// negligible.
func New(prefix string) string {
	if prefix == "" {
		prefix = "radix"
	}
	var b [7]byte
	_, _ = rand.Read(b[:])
	return prefix + "-" + strings.ToLower(enc.EncodeToString(b[:]))[:10]
}
