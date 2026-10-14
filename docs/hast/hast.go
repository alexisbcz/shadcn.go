// Package hast is a small HTML syntax tree for the docs content pipeline,
// modelled on hast (the tree unified, rehype and MDX use). Properties are
// kept as HTML attributes, in insertion order, so trees serialize exactly
// like hast-util-to-html serializes upstream's trees.
//
// It also carries the MDX nodes that pass through rehype untouched: JSX
// elements, expressions and ESM.
package hast

import (
	"strconv"
	"strings"
)

// Node is a *Element, *Text, *MdxJsx, *MdxExpression or *MdxEsm.
type Node interface{ isNode() }

// Text is a text node.
type Text struct{ Value string }

// Prop is one attribute. A Prop with Unset true keeps its position (hast
// keeps a key set to undefined in place) but is not serialized.
type Prop struct {
	Name  string
	Value string
	// Bool marks a boolean attribute that serializes as its bare name.
	Bool  bool
	Unset bool
}

// Element is an HTML element.
type Element struct {
	Tag      string
	Props    []Prop
	Children []Node
	// Data holds values that do not serialize (hast's node.data).
	Data map[string]any
}

// AttrKind distinguishes how a JSX attribute value was written.
type AttrKind int

const (
	// AttrString is name="value".
	AttrString AttrKind = iota
	// AttrBool is a bare name (value true).
	AttrBool
	// AttrExpression is name={expression}.
	AttrExpression
)

// JsxAttr is one attribute of an MDX JSX element.
type JsxAttr struct {
	Name  string
	Kind  AttrKind
	Value string // the string, or the expression source
	// Expr is the parsed expression, for AttrExpression.
	Expr Expr
}

// MdxJsx is a JSX element written in MDX: <Name attr="x">children</Name>.
// Flow is true for block-level (mdxJsxFlowElement) elements.
type MdxJsx struct {
	Name     string
	Attrs    []JsxAttr
	Children []Node
	Flow     bool
}

// MdxExpression is {expression} in MDX.
type MdxExpression struct {
	Value string
	Flow  bool
	Expr  Expr
}

// MdxEsm is an import or export statement.
type MdxEsm struct {
	Value string
	// Imports lists the bindings of an import statement.
	Imports []Import
}

// Import is one binding of an ESM import: `import { Local } from "Source"`.
type Import struct{ Local, Imported, Source string }

// Expr is a parsed JavaScript expression (see package mdx).
type Expr interface{}

func (*Element) isNode()       {}
func (*Text) isNode()          {}
func (*MdxJsx) isNode()        {}
func (*MdxExpression) isNode() {}
func (*MdxEsm) isNode()        {}

// El builds an element with attribute pairs ("name", "value", ...).
func El(tag string, attrs ...string) *Element {
	e := &Element{Tag: tag}
	for i := 0; i+1 < len(attrs); i += 2 {
		e.Set(attrs[i], attrs[i+1])
	}
	return e
}

func (e *Element) index(name string) int {
	for i := range e.Props {
		if e.Props[i].Name == name {
			return i
		}
	}
	return -1
}

// Get returns the value of a set attribute.
func (e *Element) Get(name string) (string, bool) {
	if i := e.index(name); i >= 0 && !e.Props[i].Unset {
		return e.Props[i].Value, true
	}
	return "", false
}

// Has reports whether the attribute is set.
func (e *Element) Has(name string) bool { _, ok := e.Get(name); return ok }

// HasKey reports whether the attribute key exists, even when unset (like
// Object.hasOwn on hast properties).
func (e *Element) HasKey(name string) bool { return e.index(name) >= 0 }

// Set sets an attribute, keeping its position when it already exists.
func (e *Element) Set(name, value string) *Element {
	if i := e.index(name); i >= 0 {
		e.Props[i] = Prop{Name: name, Value: value}
		return e
	}
	e.Props = append(e.Props, Prop{Name: name, Value: value})
	return e
}

// Unset removes the attribute's value but keeps its position.
func (e *Element) Unset(name string) {
	if i := e.index(name); i >= 0 {
		e.Props[i] = Prop{Name: name, Unset: true}
	}
}

// Delete removes the attribute entirely.
func (e *Element) Delete(name string) {
	if i := e.index(name); i >= 0 {
		e.Props = append(e.Props[:i], e.Props[i+1:]...)
	}
}

// Classes returns the class list.
func (e *Element) Classes() []string {
	v, _ := e.Get("class")
	return strings.Fields(v)
}

// Clone returns a deep copy.
func (e *Element) Clone() *Element {
	c := &Element{Tag: e.Tag, Props: append([]Prop(nil), e.Props...)}
	for _, ch := range e.Children {
		c.Children = append(c.Children, CloneNode(ch))
	}
	if e.Data != nil {
		c.Data = map[string]any{}
		for k, v := range e.Data {
			c.Data[k] = v
		}
	}
	return c
}

// CloneNode deep-copies a node.
func CloneNode(n Node) Node {
	switch x := n.(type) {
	case *Element:
		return x.Clone()
	case *Text:
		t := *x
		return &t
	case *MdxJsx:
		c := *x
		c.Attrs = append([]JsxAttr(nil), x.Attrs...)
		c.Children = nil
		for _, ch := range x.Children {
			c.Children = append(c.Children, CloneNode(ch))
		}
		return &c
	case *MdxExpression:
		c := *x
		return &c
	case *MdxEsm:
		c := *x
		return &c
	}
	return n
}

// ToString returns the text content (hast-util-to-string).
func ToString(n Node) string {
	var b strings.Builder
	var walk func(Node)
	walk = func(n Node) {
		switch x := n.(type) {
		case *Text:
			b.WriteString(x.Value)
		case *Element:
			for _, c := range x.Children {
				walk(c)
			}
		case *MdxJsx:
			for _, c := range x.Children {
				walk(c)
			}
		}
	}
	walk(n)
	return b.String()
}

// Children returns the children of a parent node.
func Children(n Node) []Node {
	switch x := n.(type) {
	case *Element:
		return x.Children
	case *MdxJsx:
		return x.Children
	}
	return nil
}

// Itoa is strconv.Itoa, for callers building numeric attributes.
func Itoa(i int) string { return strconv.Itoa(i) }
