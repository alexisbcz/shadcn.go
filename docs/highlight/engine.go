package highlight

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/dop251/goja"

	"shadcn.go/docs/hast"
)

// ShikiVersion is the vendored shiki version: the one rehype-pretty-code
// 0.14.1 resolves to in ui/apps/v4 (pnpm lockfile).
const ShikiVersion = "3.23.0"

// shikiJS is shiki 3.23.0 (@shikijs/core, @shikijs/engine-javascript with
// oniguruma-to-es 4.3.6), the grammars and the two themes the docs use,
// bundled unmodified by gen/build.sh. See dist/LICENSES.
//
//go:embed dist/shiki.js
var shikiJS string

// literalFixes rewrites the three regular-expression literals of
// oniguruma-to-es and oniguruma-parser that goja cannot parse: they use
// \p{AHex} (ASCII hex digit), which goja's regexp engine does not know.
// [0-9A-Fa-f] matches the same characters. Literals are compiled when the
// script is parsed, before the prelude can intercept RegExp, so they are
// rewritten in the source; every other \p{...} is rewritten at run time by
// the prelude. The count is checked so a changed bundle fails loudly.
var literalFixes = []struct {
	old, new string
	n        int
}{
	{`/^\\x[89A-Fa-f]\p{AHex}/u`, `/^\\x[89A-Fa-f][0-9A-Fa-f]/u`, 1},
	{`/^(?:\\u(?!\p{AHex}{4})|\\x(?!\p{AHex}{1,2}|\{\p{AHex}{1,8}\}))/u`, `/^(?:\\u(?![0-9A-Fa-f]{4})|\\x(?![0-9A-Fa-f]{1,2}|\{[0-9A-Fa-f]{1,8}\}))/u`, 1},
	{`/^\\x\{\s*(?<hex>\p{AHex}+)/u`, `/^\\x\{\s*(?<hex>[0-9A-Fa-f]+)/u`, 1},
}

// prelude fills the goja gaps shiki's JavaScript engine runs into:
//
//   - Unicode property escapes. RegExp is wrapped so that, with the u flag,
//     every \p{Name} and \P{Name} becomes an explicit character class
//     computed in Go (see uprops.go): same code points, syntax goja accepts.
//   - Modifier groups. goja accepts .NET-style "(?i:...)", which ECMAScript
//     2018 does not have. oniguruma-to-es probes for it to pick its output
//     target; the wrapper rejects it, as V8 without modifiers would, so
//     oniguruma-to-es emits plain ES2018 patterns (its most portable target).
const prelude = `(function (g, uprop, dCompile, dExec) {
  const Native = g.RegExp
  function rewrite(src, flags) {
    const unicode = flags.indexOf("u") >= 0
    let out = "", inClass = false
    for (let i = 0; i < src.length; i++) {
      const c = src[i]
      if (c === "\\") {
        const n = src[i + 1]
        if (unicode && (n === "p" || n === "P") && src[i + 2] === "{") {
          const end = src.indexOf("}", i + 3)
          if (end < 0) throw new SyntaxError("Invalid regular expression: /" + src + "/: Invalid property name")
          const body = uprop(src.slice(i + 3, end), n === "P")
          out += inClass ? body : "[" + body + "]"
          i = end
          continue
        }
        out += c + (n === undefined ? "" : n)
        i++
        continue
      }
      if (inClass) {
        if (c === "]") inClass = false
      } else if (c === "[") {
        inClass = true
        out += c
        if (src[i + 1] === "^") { out += "^"; i++ }
        if (src[i + 1] === "]") { out += "]"; i++; inClass = false }
        continue
      } else if (c === "(" && src[i + 1] === "?" && /^[imsx]*(?:-[imsx]*)?:/.test(src.slice(i + 2)) && /^[imsx-]/.test(src[i + 2])) {
        throw new SyntaxError("Invalid regular expression: /" + src + "/: Invalid group")
      }
      out += c
    }
    return out
  }
  const meta = new WeakMap()
  function RegExp(pattern, flags) {
    const isRe = pattern instanceof Native
    if (!new.target && isRe && flags === undefined && pattern.constructor === RegExp) return pattern
    const src = isRe ? pattern.source : pattern === undefined ? "" : String(pattern)
    let fl = flags === undefined ? (isRe ? pattern.flags : "") : String(flags)
    const hasIndices = fl.indexOf("d") >= 0
    if (hasIndices) fl = fl.replace("d", "")
    const rewritten = rewrite(src, fl)
    const re = Reflect.construct(Native, [rewritten, fl], new.target || RegExp)
    if (hasIndices) {
      const c = dCompile(rewritten, fl)
      meta.set(re, { id: c.id, names: c.names, global: fl.indexOf("g") >= 0 || fl.indexOf("y") >= 0 })
    }
    return re
  }
  const proto = Native.prototype
  const nativeExec = proto.exec
  const flagsGetter = Object.getOwnPropertyDescriptor(proto, "flags").get
  Object.defineProperty(proto, "hasIndices", { configurable: true, get() { return meta.has(this) } })
  Object.defineProperty(proto, "flags", { configurable: true, get() { const f = flagsGetter.call(this); return meta.has(this) ? "d" + f : f } })
  Object.defineProperty(proto, "exec", { writable: true, configurable: true, value: function exec(str) {
    const m = meta.get(this)
    if (!m) return nativeExec.call(this, str)
    const s = String(str)
    let last = 0
    if (m.global) {
      last = Math.min(Math.max(Math.trunc(Number(this.lastIndex)) || 0, 0), Number.MAX_SAFE_INTEGER)
      if (last > s.length) { this.lastIndex = 0; return null }
    }
    const r = dExec(m.id, s, last)
    if (r === null) { if (m.global) this.lastIndex = 0; return null }
    const res = [], indices = []
    let groups, igroups
    if (m.names.some((n) => n)) { groups = Object.create(null); igroups = Object.create(null) }
    for (let i = 0; i < r.length / 2; i++) {
      const a = r[2 * i], b = r[2 * i + 1]
      res.push(a < 0 ? undefined : s.slice(a, b))
      indices.push(a < 0 ? undefined : [a, b])
      if (groups && m.names[i]) { groups[m.names[i]] = res[i]; igroups[m.names[i]] = indices[i] }
    }
    res.index = r[0]
    res.input = s
    res.groups = groups
    indices.groups = igroups
    res.indices = indices
    if (m.global) this.lastIndex = r[1]
    return res
  } })
  for (const k of Object.getOwnPropertyNames(Native)) {
    if (k !== "length" && k !== "name" && k !== "prototype") Object.defineProperty(RegExp, k, Object.getOwnPropertyDescriptor(Native, k))
  }
  for (const s of Object.getOwnPropertySymbols(Native)) Object.defineProperty(RegExp, s, Object.getOwnPropertyDescriptor(Native, s))
  Object.defineProperty(RegExp, "prototype", { value: Native.prototype })
  Object.defineProperty(Native.prototype, "constructor", { value: RegExp, writable: true, configurable: true, enumerable: false })
  g.RegExp = RegExp
})`

// glue creates the highlighter the way rehype-pretty-code does
// (getSingletonHighlighter with both themes), with every grammar the docs
// use loaded, and serializes shiki's hast as nested arrays so property
// order survives the trip to Go.
const glue = `(function (S) {
  const h = S.createHighlighterCoreSync({
    engine: S.createJavaScriptRegexEngine(),
    themes: [S.themes.vesper, S.themes["github-light-default"]],
    langs: Object.values(S.langs),
  })
  const enc = (n) => n.type === "text" ? n.value
    : n.type === "element" ? [n.tagName, Object.entries(n.properties || {}).map(([k, v]) => [k, Array.isArray(v) ? v.join(" ") : String(v)]), (n.children || []).map(enc)]
    : n.type === "root" ? ["#root", [], n.children.map(enc)] : null
  return {
    hast(code, lang) {
      return JSON.stringify(enc(h.codeToHast(code, { lang, themes: { dark: "vesper", light: "github-light-default" }, defaultColor: false })))
    },
    langs() { return JSON.stringify(h.getLoadedLanguages()) },
  }
})`

type engine struct {
	vm    *goja.Runtime
	hast  goja.Callable
	langs map[string]bool
}

var (
	progOnce sync.Once
	prog     *goja.Program
	progErr  error
)

func compiled() (*goja.Program, error) {
	progOnce.Do(func() {
		src := shikiJS
		for _, f := range literalFixes {
			if c := strings.Count(src, f.old); c != f.n {
				progErr = fmt.Errorf("highlight: bundle changed: %d occurrences of %s, want %d", c, f.old, f.n)
				return
			}
			src = strings.ReplaceAll(src, f.old, f.new)
		}
		prog, progErr = goja.Compile("shiki.js", src, false)
	})
	return prog, progErr
}

func newEngine() (*engine, error) {
	p, err := compiled()
	if err != nil {
		return nil, err
	}
	vm := goja.New()
	pre, err := vm.RunString(prelude)
	if err != nil {
		return nil, err
	}
	preFn, _ := goja.AssertFunction(pre)
	uprop := func(name string, negate bool) string {
		s, err := propClass(name, negate)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}
		return s
	}
	reg := &dRegistry{}
	dCompile := func(source, flags string) map[string]any {
		id, err := reg.compile(source, flags)
		if err != nil {
			panic(vm.NewGoError(fmt.Errorf("SyntaxError: Invalid regular expression: /%s/%s: %v", source, flags, err)))
		}
		return map[string]any{"id": id, "names": reg.list[id].names}
	}
	dExec := func(id int, s string, start int) goja.Value {
		r, err := reg.exec(id, s, start)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		if r == nil {
			return goja.Null()
		}
		return vm.ToValue(r)
	}
	if _, err := preFn(goja.Undefined(), vm.GlobalObject(), vm.ToValue(uprop), vm.ToValue(dCompile), vm.ToValue(dExec)); err != nil {
		return nil, err
	}
	if _, err := vm.RunProgram(p); err != nil {
		return nil, fmt.Errorf("highlight: loading shiki: %w", err)
	}
	g, err := vm.RunString(glue)
	if err != nil {
		return nil, err
	}
	gfn, _ := goja.AssertFunction(g)
	api, err := gfn(goja.Undefined(), vm.Get("Shiki"))
	if err != nil {
		return nil, fmt.Errorf("highlight: creating highlighter: %w", err)
	}
	obj := api.ToObject(vm)
	hastFn, _ := goja.AssertFunction(obj.Get("hast"))
	langsFn, _ := goja.AssertFunction(obj.Get("langs"))
	lv, err := langsFn(obj)
	if err != nil {
		return nil, err
	}
	var names []string
	if err := json.Unmarshal([]byte(lv.String()), &names); err != nil {
		return nil, err
	}
	e := &engine{vm: vm, hast: hastFn, langs: map[string]bool{}}
	for _, n := range names {
		e.langs[n] = true
	}
	return e, nil
}

// codeToHast runs shiki's codeToHast with the docs' themes and returns the
// root's single <pre> element.
func (e *engine) codeToHast(code, lang string) (*hast.Element, error) {
	v, err := e.hast(goja.Undefined(), e.vm.ToValue(code), e.vm.ToValue(lang))
	if err != nil {
		return nil, err
	}
	var raw any
	if err := json.Unmarshal([]byte(v.String()), &raw); err != nil {
		return nil, err
	}
	n, err := decode(raw)
	if err != nil {
		return nil, err
	}
	root, ok := n.(*hast.Element)
	if !ok || len(root.Children) != 1 {
		return nil, errors.New("highlight: unexpected shiki output")
	}
	pre, ok := root.Children[0].(*hast.Element)
	if !ok || pre.Tag != "pre" {
		return nil, errors.New("highlight: unexpected shiki output")
	}
	return pre, nil
}

func decode(v any) (hast.Node, error) {
	switch x := v.(type) {
	case string:
		return &hast.Text{Value: x}, nil
	case []any:
		if len(x) != 3 {
			return nil, errors.New("highlight: bad node")
		}
		tag, _ := x[0].(string)
		el := &hast.Element{Tag: tag}
		props, _ := x[1].([]any)
		for _, p := range props {
			kv, _ := p.([]any)
			if len(kv) != 2 {
				return nil, errors.New("highlight: bad property")
			}
			k, _ := kv[0].(string)
			val, _ := kv[1].(string)
			el.Set(attrName(k), val)
		}
		kids, _ := x[2].([]any)
		for _, k := range kids {
			c, err := decode(k)
			if err != nil {
				return nil, err
			}
			el.Children = append(el.Children, c)
		}
		return el, nil
	}
	return nil, fmt.Errorf("highlight: bad node %T", v)
}

// attrName maps the hast property names shiki emits to HTML attribute
// names, as hast-util-to-html does.
func attrName(k string) string {
	switch k {
	case "className":
		return "class"
	case "tabindex", "tabIndex":
		return "tabindex"
	}
	return k
}
