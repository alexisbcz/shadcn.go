// Package cva builds class names from variants, like the
// class-variance-authority library used by shadcn/ui:
//
//	var buttonVariants = cva.New("inline-flex items-center", cva.Config{
//		Variants: cva.Variants{
//			"variant": {"default": "bg-primary", "outline": "border"},
//			"size":    {"default": "h-9 px-4", "sm": "h-8 px-3"},
//		},
//		Defaults: cva.Selection{"variant": "default", "size": "default"},
//	})
//
//	buttonVariants.Classes(cva.Selection{"size": "sm"}) // "inline-flex items-center bg-primary h-8 px-3"
//
// Keep the class strings literal in Go source so Tailwind's scanner finds
// them.
package cva

import (
	"slices"
	"sort"
	"strings"
)

// Variants maps a variant name to its options and their classes.
type Variants map[string]map[string]string

// Selection picks an option for some variants.
type Selection map[string]string

// Compound adds Class when every variant in When has the given option.
type Compound struct {
	When  Selection
	Class string
}

// Config holds the variants, their defaults and compound variants.
type Config struct {
	Variants Variants
	Defaults Selection
	Compound []Compound
	// Order lists variant names in the order their classes are emitted, as
	// the key order of the JavaScript object would. Variants not listed
	// follow in sorted order.
	Order []string
}

// CVA is a compiled class-variance definition. It is safe for concurrent
// use.
type CVA struct {
	base  string
	cfg   Config
	names []string // variant names in emission order
}

// New returns a CVA with the given base classes.
func New(base string, cfg Config) *CVA {
	names := make([]string, 0, len(cfg.Variants))
	for _, n := range cfg.Order {
		if _, ok := cfg.Variants[n]; ok {
			names = append(names, n)
		}
	}
	var rest []string
	for n := range cfg.Variants {
		if !slices.Contains(names, n) {
			rest = append(rest, n)
		}
	}
	sort.Strings(rest)
	names = append(names, rest...)
	return &CVA{base: base, cfg: cfg, names: names}
}

// Classes returns the base classes, the classes of the selected (or
// default) option of each variant, matching compound classes and finally
// extra. Unknown options select nothing. The result is not merged: pass it
// through lib.Cn (package cn) to resolve conflicts, as shadcn's cn() does.
func (c *CVA) Classes(sel Selection, extra ...string) string {
	parts := make([]string, 0, 2+len(c.names)+len(c.cfg.Compound)+len(extra))
	parts = append(parts, c.base)
	for _, name := range c.names {
		if cls := c.cfg.Variants[name][c.option(sel, name)]; cls != "" {
			parts = append(parts, cls)
		}
	}
	for _, cp := range c.cfg.Compound {
		if c.matches(sel, cp.When) {
			parts = append(parts, cp.Class)
		}
	}
	parts = append(parts, extra...)
	return strings.Join(strings.Fields(strings.Join(parts, " ")), " ")
}

// option returns the option in effect for a variant: the selected one, or
// the default.
func (c *CVA) option(sel Selection, name string) string {
	if v, ok := sel[name]; ok && v != "" {
		return v
	}
	return c.cfg.Defaults[name]
}

func (c *CVA) matches(sel, when Selection) bool {
	for name, want := range when {
		if c.option(sel, name) != want {
			return false
		}
	}
	return true
}

// Options returns the options of a variant in sorted order.
func (c *CVA) Options(variant string) []string {
	opts := make([]string, 0, len(c.cfg.Variants[variant]))
	for o := range c.cfg.Variants[variant] {
		opts = append(opts, o)
	}
	sort.Strings(opts)
	return opts
}
