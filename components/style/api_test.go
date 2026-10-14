package style

import (
	"fmt"
	"testing"

	"shadcn.go/cva"
)

// TestAPIDifferential ports every class site of every base source of the 24
// directories to the public API (File.CVA, File.Cn, File.Attr), the way a Go
// component is written, and requires the generated source's strings.
func TestAPIDifferential(t *testing.T) {
	var summary []string
	for _, b := range bases {
		for _, st := range []Style{Nova, Vega, Lyra, Maia, Mira, Luma, Sera, Rhea} {
			dir := b + "-" + string(st)
			ref, err := loadRef(dir)
			if err != nil {
				t.Fatal(err)
			}
			var stats apiStats
			for _, f := range ref.Files {
				checkFileAPI(t, dir, f, st, &stats)
			}
			summary = append(summary, fmt.Sprintf("%-11s cva %4d (strings %5d)  cn sites %5d  attr sites %4d  strings %6d  verbatim %5d  mismatches %d",
				dir, stats.cvas, stats.cvaStrings, stats.cnSites, stats.attrSites, stats.strings, stats.verbatim, stats.mismatches))
			if stats.mismatches > 0 {
				t.Errorf("%s: %d mismatches", dir, stats.mismatches)
			}
		}
	}
	for _, s := range summary {
		t.Log(s)
	}
}

type apiStats struct {
	cvas, cvaStrings, cnSites, attrSites, strings, mismatches int
	// verbatim counts the strings outside the API's sites (dynamic cn()
	// arguments and everything else in a class root) that a Go author
	// copies from the base source: they must equal the generated source.
	verbatim int
	// skipVerbatim skips that check, for sources that are not in
	// prettier's class order (the fuzz cases); upstream's are.
	skipVerbatim bool
	// skipCVAClasses skips comparing CVA.For(st).Classes with a plain join
	// of the strings: shadcn.go/cva normalizes whitespace with
	// strings.Fields, which class-variance-authority does not; that only
	// shows with whitespace-only or non-ASCII whitespace values (fuzz).
	skipCVAClasses bool
}

type cvaPair struct {
	def      *CVA
	base     *node // base cva() call
	gen      *node // generated cva() call
	variants []Variant
}

func checkFileAPI(t *testing.T, dir string, f refFile, st Style, stats *apiStats) {
	base, gen := decodeRoots(f.Base), decodeRoots(f.Gen)
	file := NewFile(f.Path)
	fail := func(format string, args ...any) {
		stats.mismatches++
		if stats.mismatches <= 5 {
			t.Errorf("%s/%s: %s", dir, f.Path, fmt.Sprintf(format, args...))
		}
	}

	// Pair the base and generated nodes of the file in document order,
	// without descending into sites (the transform rewrites their insides).
	// depth counts the cn()/cva() calls enclosing a node: each sorts the
	// node's literals once more (see Site.Enclosed).
	var bn, gn []*node
	depth := map[*node]int{}
	var pair func(b, g *node, d int)
	pair = func(b, g *node, d int) {
		bn, gn = append(bn, b), append(gn, g)
		depth[b] = d
		if isSite(b) {
			return
		}
		if isSortCall(b) {
			d++
		}
		bk, gk := childList(b), childList(g)
		if len(bk) != len(gk) {
			fail("structure differs at %s", describe(b))
			return
		}
		for i := range bk {
			pair(bk[i], gk[i], d)
		}
	}
	for i := range base {
		pair(base[i], gen[i], 0)
	}

	// Nodes the API resolves: cva() calls and className="..." sites
	// (covered), cn() sites (sites, keyed by the base call).
	covered := map[*node]bool{}
	sites := map[*node]*Site{}

	// 1. cva() definitions, in document order.
	var cvas []cvaPair
	for i, n := range bn {
		if n.kind != kCall || n.callee != "cva" {
			continue
		}
		if i >= len(gn) || gn[i].kind != kCall || gn[i].callee != "cva" {
			fail("cva call %d has no generated counterpart", i)
			continue
		}
		baseStr, variants, compounds, ok := cvaDefinition(n)
		if !ok {
			fail("cva call %d has a shape the API does not express", i)
			continue
		}
		covered[n] = true
		cvas = append(cvas, cvaPair{def: file.CVA(baseStr, variants, nil, compounds...), base: n, gen: gn[i], variants: variants})
	}
	for _, c := range cvas {
		stats.cvas++
		gotBase, gotVariants, gotCompounds := c.def.Resolved(st)
		wantBase, wantVariants, wantCompounds, _ := cvaDefinition(c.gen)
		cmp := func(what, g, w string) {
			stats.cvaStrings++
			stats.strings++
			if g != w {
				fail("cva %s: got %q, want %q", what, g, w)
			}
		}
		cmp("base", gotBase, wantBase)
		if len(gotVariants) != len(wantVariants) || len(gotCompounds) != len(wantCompounds) {
			fail("cva shape differs")
			continue
		}
		for vi, v := range wantVariants {
			for oi, o := range v.Options {
				cmp(fmt.Sprintf("%s=%s", v.Name, o.Name), gotVariants[vi].Options[oi].Class, o.Class)
			}
		}
		for i, cv := range wantCompounds {
			cmp(fmt.Sprintf("compound %d", i), gotCompounds[i].Class, cv.Class)
		}
		// The class-variance-authority definition carries the same strings.
		def := c.def.For(st)
		if len(wantVariants) > 0 && len(wantVariants[0].Options) > 0 && !stats.skipCVAClasses {
			sel := cva.Selection{wantVariants[0].Name: wantVariants[0].Options[0].Name}
			for _, v := range wantVariants[1:] {
				sel[v.Name] = "\x00none"
			}
			want := joinNonEmpty(wantBase, wantVariants[0].Options[0].Class)
			for _, cv := range wantCompounds {
				if matches(cv.When, sel) {
					want = joinNonEmpty(want, cv.Class)
				}
			}
			if got := def.Classes(sel); got != want {
				fail("cva.Classes(%v): got %q, want %q", sel, got, want)
			}
		}
	}

	// 2. Sites: className attributes and className: cn() properties.
	var check func(site *Site, want *node)
	check = func(site *Site, want *node) {
		args := site.Args(st)
		var wantArgs []*node
		switch {
		case want.kind == kStr:
			wantArgs = []*node{want}
		case want.kind == kCall:
			wantArgs = want.args
		default:
			fail("unexpected generated site %s", describe(want))
			return
		}
		if len(args) != len(wantArgs) {
			fail("site args: got %d, want %d", len(args), len(wantArgs))
			return
		}
		for i, a := range args {
			if s, ok := a.(string); ok {
				stats.strings++
				if !wantArgs[i].isStringLiteralLike() || wantArgs[i].str != s {
					fail("site arg %d: got %q, want %s", i, s, describe(wantArgs[i]))
				}
			} else if wantArgs[i].isStringLiteralLike() {
				fail("site arg %d: got a dynamic argument, want %s", i, describe(wantArgs[i]))
			}
		}
	}
	// cnSite declares the site of a cn() call: in a className attribute
	// (Cn) or a className property (PropCn), inside enclosed more cn() or
	// cva() calls.
	cnSite := func(call *node, inAttr bool, enclosed int) *Site {
		var args []any
		for _, a := range call.args {
			if a.isStringLiteralLike() {
				args = append(args, a.str)
			} else {
				args = append(args, Arg)
			}
		}
		var s *Site
		if inAttr {
			s = file.Cn(args...)
		} else {
			s = file.PropCn(args...)
		}
		if enclosed > 0 {
			s.Enclosed(enclosed)
		}
		sites[call] = s
		return s
	}
	for i, n := range bn {
		if i >= len(gn) {
			break
		}
		g := gn[i]
		switch {
		case n.kind == kAttr && n.name == "className" && n.value != nil:
			if g.kind != kAttr || g.value == nil {
				fail("attribute %d has no generated counterpart", i)
				continue
			}
			switch {
			case n.value.kind == kStr:
				stats.attrSites++
				covered[n] = true
				check(file.Attr(n.value.str), g.value)
			case n.value.kind == kJsxExpr && n.value.value != nil && n.value.value.isCnCall():
				stats.cnSites++
				if g.value.kind != kJsxExpr {
					fail("attribute %d: generated value is not an expression", i)
					continue
				}
				check(cnSite(n.value.value, true, depth[n]), g.value.value)
			case n.value.kind == kJsxExpr && n.value.value != nil && n.value.value.isStringLiteralLike():
				stats.attrSites++
				covered[n] = true
				check(file.Attr(n.value.value.str), g.value.value)
			}
		case n.kind == kProp && n.name == "className" && n.value != nil && n.value.isCnCall():
			if g.kind != kProp || g.value == nil || g.value.kind != kCall {
				fail("property %d has no generated counterpart", i)
				continue
			}
			stats.cnSites++
			check(cnSite(n.value, false, depth[n]), g.value)
		}
	}
	// cn() calls nested in className expressions (composeRenderProps, ...).
	for i, n := range bn {
		if n.kind != kAttr || n.name != "className" || n.value == nil || i >= len(gn) {
			continue
		}
		if gn[i].kind != kAttr || gn[i].value == nil {
			continue
		}
		bc, gc := nestedCnCalls(n.value), nestedCnCalls(gn[i].value)
		within := sortCallsAbove(n.value)
		if len(bc) != len(gc) {
			fail("attribute %d: %d nested cn() calls, generated %d", i, len(bc), len(gc))
			continue
		}
		for j := range bc {
			stats.cnSites++
			check(cnSite(bc[j], true, depth[n]+within[bc[j]]), gc[j])
		}
	}

	// 3. Everything else in the class roots is written verbatim in Go: the
	// dynamic arguments of cn() sites (the style.Arg values), className
	// expressions with no site (templates, variables), cn() calls outside
	// className contexts, other attributes and properties. Upstream's build
	// must leave all of it unchanged, in every style.
	var verbatim func(b, g *node, path string) string
	verbatim = func(b, g *node, path string) string {
		if b == nil || g == nil {
			if b != g {
				return fmt.Sprintf("%s: %s, generated %s", path, describe(b), describe(g))
			}
			return ""
		}
		if covered[b] {
			return ""
		}
		if site := sites[b]; site != nil {
			parts := site.resolve(st)
			if g.kind != kCall || len(g.args) != len(parts) {
				return fmt.Sprintf("%s: cn() site, generated %s", path, describe(g))
			}
			var dyn []*node
			for _, a := range b.args {
				if !a.isStringLiteralLike() {
					dyn = append(dyn, a)
				}
			}
			for k, p := range parts {
				if p.slot < 0 {
					continue
				}
				if d := verbatim(dyn[p.slot], g.args[k], fmt.Sprintf("%s.cn(arg %d)", path, k)); d != "" {
					return d
				}
			}
			return ""
		}
		same := b.kind == g.kind && b.str == g.str && b.tpl == g.tpl && b.callee == g.callee &&
			b.calleeRoot == g.calleeRoot && b.name == g.name && b.ident == g.ident && b.op == g.op &&
			len(b.quasis) == len(g.quasis)
		for i := 0; same && i < len(b.quasis); i++ {
			same = b.quasis[i] == g.quasis[i]
		}
		if !same {
			return fmt.Sprintf("%s: %s, generated %s", path, describe(b), describe(g))
		}
		if b.kind == kStr {
			stats.verbatim++
		}
		stats.verbatim += len(b.quasis)
		bk, gk := childList(b), childList(g)
		if len(bk) != len(gk) {
			return fmt.Sprintf("%s: %d children, generated %d", path, len(bk), len(gk))
		}
		for i := range bk {
			if d := verbatim(bk[i], gk[i], fmt.Sprintf("%s/%s[%d]", path, describe(b), i)); d != "" {
				return d
			}
		}
		return ""
	}
	for i := range base {
		if stats.skipVerbatim {
			break
		}
		if d := verbatim(base[i], gen[i], fmt.Sprintf("line %d", f.Base[i].Line)); d != "" {
			fail("not reproducible verbatim: %s", d)
		}
	}
}

func isSortCall(n *node) bool {
	return n.kind == kCall && sortFunctions[n.calleeRoot]
}

// sortCallsAbove maps every node under root to the number of cn()/cva()
// calls between root and it.
func sortCallsAbove(root *node) map[*node]int {
	out := map[*node]int{}
	var walk func(n *node, d int)
	walk = func(n *node, d int) {
		out[n] = d
		if isSortCall(n) {
			d++
		}
		n.kids(func(k *node) { walk(k, d) })
	}
	walk(root, 0)
	return out
}

func isSite(n *node) bool {
	return n.kind == kCall && n.callee == "cva" ||
		n.kind == kAttr && n.name == "className" ||
		n.kind == kProp && n.name == "className"
}

func childList(n *node) []*node {
	var out []*node
	n.kids(func(c *node) { out = append(out, c) })
	return out
}

// cvaDefinition reads a cva(base, { variants, compoundVariants }) call.
func cvaDefinition(call *node) (base string, variants []Variant, compounds []cva.Compound, ok bool) {
	if len(call.args) == 0 || !call.args[0].isStringLiteral() {
		return "", nil, nil, false
	}
	base = call.args[0].str
	if len(call.args) < 2 {
		return base, nil, nil, true
	}
	if call.args[1].kind != kObject {
		return "", nil, nil, false
	}
	for _, p := range call.args[1].children {
		if p.kind != kProp {
			return "", nil, nil, false
		}
		switch p.name {
		case "variants":
			if p.value.kind != kObject {
				return "", nil, nil, false
			}
			for _, v := range p.value.children {
				if v.kind != kProp || v.value.kind != kObject {
					return "", nil, nil, false
				}
				vr := Variant{Name: unquote(v.name)}
				for _, o := range v.value.children {
					if o.kind != kProp || !o.value.isStringLiteral() {
						return "", nil, nil, false
					}
					vr.Options = append(vr.Options, Option{Name: unquote(o.name), Class: o.value.str})
				}
				variants = append(variants, vr)
			}
		case "compoundVariants":
			if p.value.kind != kGeneric || p.value.str != "ArrayLiteralExpression" {
				return "", nil, nil, false
			}
			for _, c := range p.value.children {
				if c.kind != kObject {
					return "", nil, nil, false
				}
				var cv cva.Compound
				cv.When = cva.Selection{}
				for _, q := range c.children {
					if q.kind != kProp || !q.value.isStringLiteralLike() {
						return "", nil, nil, false
					}
					if q.name == "className" || q.name == "class" {
						cv.Class = q.value.str
					} else {
						cv.When[unquote(q.name)] = q.value.str
					}
				}
				compounds = append(compounds, cv)
			}
		}
	}
	return base, variants, compounds, true
}

func unquote(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func joinNonEmpty(a, b string) string {
	switch {
	case a == "":
		return b
	case b == "":
		return a
	}
	return a + " " + b
}

func matches(when, sel cva.Selection) bool {
	for k, v := range when {
		if sel[k] != v {
			return false
		}
	}
	return true
}
