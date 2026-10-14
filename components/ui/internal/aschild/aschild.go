// Package aschild implements Radix's Slot (the asChild prop) for compound
// components: the part's props are merged into the single child the caller
// provides. It is the same algorithm as dialog/primitive's Slot, available
// without linking the dialog behavior.
package aschild

import (
	"fmt"
	"strings"

	"shadcn.go/components/ui/internal/tree"
	h "shadcn.go/libhtml"
)

// component is a builder-based component (embedding h.Builder): its user
// props live in Elem, its rendered element comes from Element.
type component interface {
	tree.Container
	h.Elementer
}

// Slot implements Radix's Slot (asChild): the part's props (slot) are merged
// into the single child the caller provides, the child's own props winning;
// classes are joined slot first.
//
// When the child is another component (a builder such as button.Button),
// React passes the merged props to that component, which spreads them over
// its own defaults: the slot's props then also win over the component's
// defaults (e.g. DialogTrigger's data-slot="dialog-trigger" replaces the
// Button's data-slot="button") while the component still merges classes
// with cn. Slot does the same by adding the slot's props to the
// component's user props before resolving it. This mutates the child
// builder; a builder is meant to be rendered once.
func Slot(slot *h.Element, child any) (*h.Element, error) {
	if slot.Err() != nil {
		return nil, slot.Err()
	}
	c, ok := child.(component)
	if !ok {
		return h.Slot(slot, child, nil)
	}
	user := c.Elem()
	names, err := attrNames(slot)
	if err != nil {
		return nil, err
	}
	for _, n := range names {
		ln := strings.ToLower(n)
		if ln == "class" || user.Has(n) || strings.HasPrefix(ln, "on") {
			continue
		}
		if v, ok := slot.GetAttr(n); ok {
			user.Attr(n, v)
		}
	}
	if sc := slot.ClassName(); sc != "" {
		if uc := user.ClassName(); uc != "" {
			user.SetClass(sc + " " + uc)
		} else {
			user.SetClass(sc)
		}
	}
	if sv, ok := slot.GetAttr("style"); ok && user.Has("style") {
		uv, _ := user.GetAttr("style")
		if !strings.Contains(uv, sv) {
			user.Attr("style", strings.TrimSuffix(strings.TrimSpace(sv), ";")+";"+uv)
		}
	}
	return c.Element(), nil
}

// Resolve resolves a part rendered with lib.AsChild: props (without
// children) slotted into the single child of props. part names the part in
// errors (e.g. "collapsible Trigger").
func Resolve(part string, props *h.Element) *h.Element {
	kids := props.ChildNodes()
	if len(kids) != 1 {
		return h.Errored(fmt.Errorf("%s: AsChild needs exactly one element child, got %d children", part, len(kids)))
	}
	out, err := Slot(props.Clone().SetChildren(), kids[0])
	if err != nil {
		return h.Errored(fmt.Errorf("%s: AsChild: %w", part, err))
	}
	return out
}

// attrNames lists the attribute names of el, in order, by rendering its
// start tag (libhtml escapes quotes in values, so the scan is exact).
func attrNames(el *h.Element) ([]string, error) {
	s, err := h.RenderString(el.Clone().SetChildren())
	if err != nil {
		return nil, err
	}
	end := strings.IndexByte(s, '>')
	if end < 0 {
		return nil, nil
	}
	s = s[:end]
	// Skip "<tag".
	i := strings.IndexAny(s, " \t\n")
	if i < 0 {
		return nil, nil
	}
	s = s[i:]
	var names []string
	for len(s) > 0 {
		s = strings.TrimLeft(s, " \t\n/")
		if s == "" {
			break
		}
		j := strings.IndexAny(s, "= \t\n")
		if j < 0 {
			names = append(names, s)
			break
		}
		names = append(names, s[:j])
		s = s[j:]
		if s[0] == '=' {
			// Quoted value: skip to the closing quote.
			if len(s) > 1 && s[1] == '"' {
				k := strings.IndexByte(s[2:], '"')
				if k < 0 {
					break
				}
				s = s[2+k+1:]
			} else {
				k := strings.IndexAny(s, " \t\n")
				if k < 0 {
					break
				}
				s = s[k:]
			}
		}
	}
	return names, nil
}
