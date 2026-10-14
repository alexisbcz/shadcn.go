// Package ui is a Go port of the shadcn/ui components (new-york-v4 style),
// built on libhtml.
//
// Every component takes its children JSX-style and also offers builder
// methods:
//
//	ui.Button(ui.ButtonVariantOutline, ui.ButtonSizeSm, h.Class("w-full"), "Save")
//	ui.Button("Save").Variant(ui.ButtonVariantOutline).Size(ui.ButtonSizeSm).Class("w-full")
//
// The Tailwind classes are kept literal in this package, so point Tailwind's
// scanner at it (@source "path/to/components/ui/*.go").
package ui

import (
	"shadcn.go/twmerge"
)

// Cn combines class lists and resolves Tailwind conflicts, later classes
// winning. It is shadcn's cn() helper.
func Cn(classes ...string) string { return twmerge.Merge(classes...) }

// asChildFlag is the type of AsChild.
type asChildFlag struct{}

// AsChild, passed as a child, makes a component render its single element
// child instead of its own element, merging its classes and attributes into
// it (Radix's asChild):
//
//	ui.Button(ui.AsChild, h.A(h.Href("/login"), "Login"))
var AsChild = asChildFlag{}
