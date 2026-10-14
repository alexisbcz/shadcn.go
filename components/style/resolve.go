package style

import (
	"fmt"
	"strconv"
	"sync"

	"shadcn.go/cn"
	"shadcn.go/cva"
)

// A File holds the classes of one upstream source file, written as in
// registry/bases/<base>/: its cva() definitions and its class sites. It
// resolves them for each style exactly as upstream's transformStyle and the
// prettier pass of its build do for that file.
//
// The file matters because upstream's transform has state: a token that a
// cva() call of the file has already been given styles for is only removed,
// not styled again, in the file's className attributes and cn() calls. So
// declare one File per ported upstream file, its CVA definitions first and
// as package-level variables (they must all exist before the first render),
// then its sites:
//
//	var file = style.NewFile("ui/button.tsx")
//
//	var buttonVariants = file.CVA("cn-button group/button inline-flex shrink-0 ...",
//		[]style.Variant{
//			{Name: "variant", Options: []style.Option{
//				{Name: "default", Class: "cn-button-variant-default"},
//				{Name: "outline", Class: "cn-button-variant-outline"},
//			}},
//			{Name: "size", Options: []style.Option{
//				{Name: "default", Class: "cn-button-size-default"},
//				{Name: "sm", Class: "cn-button-size-sm"},
//			}},
//		},
//		cva.Selection{"variant": "default", "size": "default"})
//
//	// className={cn(buttonVariants({ variant, size, className }))}
//	cls := lib.Cn(buttonVariants.Classes(st, cva.Selection{"variant": v, "size": s}, className))
//
// Sites mirror how the base source passes classes:
//
//	file.Cn("cn-card-header grid gap-1", style.Arg) // className={cn("cn-card-header ...", className)}
//	file.Attr("cn-card-title")                      // className="cn-card-title"
//
// Resolution happens once per style, on first use, and is cached.
type File struct {
	name string

	mu     sync.Mutex
	cvas   []*CVA
	frozen bool

	once    [8]sync.Once
	matched [8]map[string]bool
}

// NewFile returns an empty File. The name (the upstream path, such as
// "ui/button.tsx") is only used in messages.
func NewFile(name string) *File { return &File{name: name} }

// Option is a variant option and its classes.
type Option struct{ Name, Class string }

// Variant is a cva() variant: its options, in source order.
type Variant struct {
	Name    string
	Options []Option
}

// CVA is a cva() definition written with base tokens. Use Classes, or For to
// get the class-variance-authority definition of a style.
type CVA struct {
	file      *File
	base      string
	variants  []Variant
	defaults  cva.Selection
	compounds []cva.Compound
	resolved  [8]*cva.CVA
}

// CVA declares the file's next cva(base, { variants, defaultVariants,
// compoundVariants }) call. Variants and options are in source order, which
// matters: upstream styles a token only the first time the file uses it.
func (f *File) CVA(base string, variants []Variant, defaults cva.Selection, compounds ...cva.Compound) *CVA {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.frozen {
		panic(fmt.Sprintf("style: %s: CVA declared after the file was first resolved; declare cva definitions as package-level variables", f.name))
	}
	c := &CVA{file: f, base: base, variants: variants, defaults: defaults, compounds: compounds}
	f.cvas = append(f.cvas, c)
	return c
}

// resolve runs the cva pass of upstream's transform over all of the file's
// definitions, for one style, and prettier's sorting on them.
func (f *File) resolve(st Style) (int, map[string]bool) {
	i := st.index()
	if i < 0 {
		panic("style: invalid style " + strconv.Quote(string(st)))
	}
	f.once[i].Do(func() {
		f.mu.Lock()
		f.frozen = true
		cvas := f.cvas
		f.mu.Unlock()

		t := &transformer{styleMap: st.Map(), matched: map[string]bool{}}
		for _, c := range cvas {
			call, base, opts, comps := c.tree()
			t.cvaCall(call)
			sortTree([]*node{call})
			c.resolved[i] = c.build(base, opts, comps)
		}
		f.matched[i] = t.matched
	})
	return i, f.matched[i]
}

// tree builds the cva() call and returns it with the nodes holding the base
// classes, each option's classes and each compound's classes.
func (c *CVA) tree() (call, base *node, opts [][]*node, comps []*node) {
	base = strNode(c.base)
	variants := &node{kind: kObject}
	for _, v := range c.variants {
		obj := &node{kind: kObject}
		var row []*node
		for _, o := range v.Options {
			s := strNode(o.Class)
			row = append(row, s)
			obj.children = append(obj.children, &node{kind: kProp, name: strconv.Quote(o.Name), value: s})
		}
		opts = append(opts, row)
		variants.children = append(variants.children, &node{kind: kProp, name: strconv.Quote(v.Name), value: obj})
	}
	compounds := &node{kind: kGeneric, str: "ArrayLiteralExpression"}
	for _, cv := range c.compounds {
		s := strNode(cv.Class)
		comps = append(comps, s)
		compounds.children = append(compounds.children, &node{kind: kObject, children: []*node{
			{kind: kProp, name: "className", ident: true, value: s},
		}})
	}
	config := &node{kind: kObject, children: []*node{
		{kind: kProp, name: "variants", ident: true, value: variants},
		{kind: kProp, name: "compoundVariants", ident: true, value: compounds},
	}}
	call = &node{kind: kCall, callee: "cva", calleeRoot: "cva", calleeNode: &node{kind: kID, str: "cva"}, args: []*node{base, config}}
	return call, base, opts, comps
}

func (c *CVA) build(base *node, opts [][]*node, comps []*node) *cva.CVA {
	cfg := cva.Config{Variants: cva.Variants{}, Defaults: c.defaults}
	for i, v := range c.variants {
		m := map[string]string{}
		for j, o := range v.Options {
			m[o.Name] = opts[i][j].str
		}
		cfg.Variants[v.Name] = m
		cfg.Order = append(cfg.Order, v.Name)
	}
	for i, cv := range c.compounds {
		cfg.Compound = append(cfg.Compound, cva.Compound{When: cv.When, Class: comps[i].str})
	}
	return cva.New(base.str, cfg)
}

// For returns the definition as upstream's generated source for style st
// has it.
func (c *CVA) For(st Style) *cva.CVA {
	i, _ := c.file.resolve(st)
	return c.resolved[i]
}

// Classes is For(st).Classes(sel, extra...): the base classes, the selected
// options' classes, matching compound classes and extra, not merged (pass
// the result to lib.Cn, as upstream passes buttonVariants(...) to cn).
func (c *CVA) Classes(st Style, sel cva.Selection, extra ...string) string {
	return c.For(st).Classes(sel, extra...)
}

// argPlaceholder is the type of Arg.
type argPlaceholder struct{}

// Arg stands for a cn() argument that is not a string literal (the
// className prop, a cva() result, a condition, ...): Site.Class takes its
// value at render time. Upstream's transform leaves such arguments as they
// are; string literals nested inside them are already sorted in the base
// sources and identical in every style.
var Arg = argPlaceholder{}

// Site is a place where the base source passes classes: a cn() call in a
// className attribute or property, or a className="..." literal.
type Site struct {
	file  *File
	attr  bool  // className="..." (a literal attribute)
	args  []any // string or Arg
	slots int

	once  [8]sync.Once
	parts [8][]part
}

// part is an argument of the generated cn() call: a literal, or the
// dynamic argument numbered slot.
type part struct {
	class string
	slot  int // -1 for a literal
}

// Cn declares a cn(args...) call of the file in a className attribute or
// property (className={cn(...)}, className: cn(...), cn() calls nested in a
// className expression). Each argument is a string literal of the base
// source or Arg.
func (f *File) Cn(args ...any) *Site {
	s := &Site{file: f}
	for _, a := range args {
		switch a := a.(type) {
		case string:
			s.args = append(s.args, a)
		case argPlaceholder:
			s.args = append(s.args, a)
			s.slots++
		default:
			panic(fmt.Sprintf("style: %s: Cn arguments are string literals or style.Arg, not %T", f.name, a))
		}
	}
	return s
}

// Attr declares a className="..." attribute of the file.
func (f *File) Attr(class string) *Site {
	return &Site{file: f, attr: true, args: []any{class}}
}

func (s *Site) resolve(st Style) []part {
	i, matched := s.file.resolve(st)
	s.once[i].Do(func() {
		t := &transformer{styleMap: st.Map(), matched: matched}
		var root *node
		if s.attr {
			root = &node{kind: kAttr, name: "className", value: strNode(s.args[0].(string))}
			t.classNameAttribute(root)
			sortTree([]*node{root})
			s.parts[i] = []part{{class: root.value.str, slot: -1}}
			return
		}
		call := newCnCall()
		slot := 0
		for _, a := range s.args {
			if c, ok := a.(string); ok {
				call.args = append(call.args, strNode(c))
			} else {
				call.args = append(call.args, &node{kind: kID, str: "\x00" + strconv.Itoa(slot)})
				slot++
			}
		}
		root = &node{kind: kProp, name: "className", ident: true, value: call}
		t.classNameProperty(root)
		sortTree([]*node{root})
		parts := make([]part, 0, len(call.args))
		for _, a := range root.value.args {
			if a.kind == kStr {
				parts = append(parts, part{class: a.str, slot: -1})
			} else {
				n, _ := strconv.Atoi(a.str[1:])
				parts = append(parts, part{slot: n})
			}
		}
		s.parts[i] = parts
	})
	return s.parts[i]
}

// Args returns the arguments of the cn() call upstream generates for style
// st: strings for literals, Arg for the dynamic arguments, in order. For an
// Attr site it is the attribute's single string.
func (s *Site) Args(st Style) []any {
	parts := s.resolve(st)
	out := make([]any, len(parts))
	for i, p := range parts {
		if p.slot < 0 {
			out[i] = p.class
		} else {
			out[i] = Arg
		}
	}
	return out
}

// Class returns the site's class attribute for style st: for a Cn site,
// cn() of the generated arguments with values in place of the Arg
// placeholders, in order; for an Attr site, the generated literal. It
// panics if the number of values does not match the placeholders.
func (s *Site) Class(st Style, values ...string) string {
	if len(values) != s.slots {
		panic(fmt.Sprintf("style: %s: site has %d dynamic arguments, got %d values", s.file.name, s.slots, len(values)))
	}
	parts := s.resolve(st)
	if s.attr {
		return parts[0].class
	}
	switch len(parts) {
	case 1:
		if parts[0].slot < 0 {
			return cn.Cn(parts[0].class)
		}
		return cn.Cn(values[parts[0].slot])
	case 2:
		return cn.Cn(parts[0].value(values), parts[1].value(values))
	}
	args := make([]string, len(parts))
	for i, p := range parts {
		args[i] = p.value(values)
	}
	return cn.Cn(args...)
}

func (p part) value(values []string) string {
	if p.slot < 0 {
		return p.class
	}
	return values[p.slot]
}

var (
	looseFile = NewFile("(standalone)")
	attrCache sync.Map // attrKey -> string
)

type attrKey struct {
	st    Style
	class string
}

// Resolve returns the class string that upstream's generated source has, for
// style st, in place of a className="class" attribute of a base source (in
// a file whose cva() calls do not use the same tokens). The result is
// cached.
//
//	style.Resolve(style.Nova, "cn-card-title")  // "font-heading text-base leading-snug font-medium ..."
func Resolve(st Style, class string) string {
	k := attrKey{st, class}
	if v, ok := attrCache.Load(k); ok {
		return v.(string)
	}
	v := looseFile.Attr(class).resolve(st)[0].class
	attrCache.Store(k, v)
	return v
}

// Cn declares a standalone cn() site (in a file without cva() calls); see
// File.Cn. Declare it once, as a package-level variable.
func Cn(args ...any) *Site { return looseFile.Cn(args...) }
