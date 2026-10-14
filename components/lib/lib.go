// Package lib holds the helpers shared by every shadcn/ui component, like
// lib/utils.ts upstream.
package lib

import (
	"shadcn.go/cn"
	"time"
)

// Cn combines class lists and resolves Tailwind conflicts, later classes
// winning. It is shadcn's cn() helper: upstream imports it from the cn npm
// package, and package cn is a byte-for-byte port of it. Empty arguments are
// skipped, like falsy arguments upstream.
func Cn(classes ...string) string { return cn.Cn(classes...) }

// AsChildFlag is the type of AsChild.
type AsChildFlag struct{}

// AsChild, passed as a child, makes a component render its single element
// child instead of its own element, merging its classes and attributes into
// it (Radix's asChild):
//
//	button.Button(lib.AsChild, h.A(h.Href("/login"), "Login"))
var AsChild = AsChildFlag{}

// now is the clock behind Now; SetNow pins it.
var now = time.Now

// Now returns the current time as components and examples should see it:
// time.Now, unless SetNow pinned it (the docs server does, with -now, so
// server-rendered dates match a test browser whose clock is frozen).
func Now() time.Time { return now() }

// SetNow pins Now to t. A zero t restores time.Now. Call it at startup,
// before serving requests.
func SetNow(t time.Time) {
	if t.IsZero() {
		now = time.Now
		return
	}
	now = func() time.Time { return t }
}
