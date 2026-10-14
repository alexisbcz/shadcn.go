#!/usr/bin/env node
// Builds cn/testdata/corpus.json.gz: cn() calls and their results computed by
// the real cn npm package, for the differential test of the Go port.
//
//   node cn/gen/corpus.mjs
//
// Cases come from three sources:
//
//  1. registry: every cn() call in the upstream registry (new-york-v4 ui,
//     examples, blocks, charts, internal, hooks, lib), its arguments
//     enumerated statically: string literals, both branches of conditionals,
//     `x && "..."` (falsy or the string), template literals, cva() variant
//     functions (evaluated with the real class-variance-authority for every
//     variant combination), and `className` filled with each className the
//     registry passes to that component (`<CardHeader className="...">` feeds
//     the cn() calls of CardHeader in ui/card.tsx). Every string literal and
//     className in the registry is also merged on its own.
//  2. pairs: registry class strings merged pairwise (sampled).
//  3. random: a seeded random corpus of tokens from the registry plus
//     synthetic tokens that exercise every engine path (variants, important,
//     negatives, arbitrary values and properties, postfix modifiers,
//     container queries, number formats, unicode, odd whitespace).
//
// Results are computed with cn() (cached, as shadcn calls it) and checked
// against an uncached engine, so the corpus also proves caching never changes
// a result.
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs"
import { gzipSync } from "node:zlib"
import { createRequire } from "node:module"
import { dirname, join, relative, resolve, basename } from "node:path"
import { fileURLToPath, pathToFileURL } from "node:url"

const here = dirname(fileURLToPath(import.meta.url))
const root = resolve(here, "..", "..")
const require = createRequire(join(root, "ui", "package.json"))
const ts = require("typescript")
const cnDist = join(root, "ui/apps/v4/node_modules/cn/dist")
const { cn, createEngine } = await import(pathToFileURL(join(cnDist, "index.js")).href)
const tables = (await import(pathToFileURL(join(cnDist, "tables.js")).href)).default
const cnPkg = JSON.parse(readFileSync(join(cnDist, "..", "package.json"), "utf8"))
const { cva } = await import(pathToFileURL(join(root, "parity/node_modules/class-variance-authority/dist/index.mjs")).href)
const uncached = createEngine(tables, undefined, { cacheSize: 0 }).mergeUncached

const registry = join(root, "ui/apps/v4/registry/new-york-v4")

// --- deterministic randomness ---------------------------------------------------
let seed = 0x5eed1234
const rnd = () => {
  seed |= 0
  seed = (seed + 0x6d2b79f5) | 0
  let t = Math.imul(seed ^ (seed >>> 15), 1 | seed)
  t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296
}
const pick = (a) => a[Math.floor(rnd() * a.length)]
const sample = (a, n) => {
  if (a.length <= n) return a
  const out = new Set()
  while (out.size < n) out.add(a[Math.floor(rnd() * a.length)])
  return [...out]
}

// --- source files ---------------------------------------------------------------
const files = []
const walk = (d) => {
  for (const f of readdirSync(d).sort()) {
    const p = join(d, f)
    if (statSync(p).isDirectory()) walk(p)
    else if (/\.(tsx?|jsx?)$/.test(f)) files.push(p)
  }
}
walk(registry)

const parse = (p) => ts.createSourceFile(p, readFileSync(p, "utf8"), ts.ScriptTarget.Latest, true, p.endsWith("x") ? ts.ScriptKind.TSX : ts.ScriptKind.TS)
const sources = files.map((p) => ({ path: p, rel: relative(registry, p), sf: parse(p) }))

const allStrings = new Set()
const CAP = 64 // max values per argument / combos per call

const uniq = (a) => [...new Set(a)]

// Statically evaluates literal expressions (for cva configs).
const evalLiteral = (n) => {
  if (ts.isStringLiteral(n) || ts.isNoSubstitutionTemplateLiteral(n)) return n.text
  if (ts.isNumericLiteral(n)) return Number(n.text)
  if (n.kind === ts.SyntaxKind.TrueKeyword) return true
  if (n.kind === ts.SyntaxKind.FalseKeyword) return false
  if (n.kind === ts.SyntaxKind.NullKeyword) return null
  if (ts.isArrayLiteralExpression(n)) return n.elements.map(evalLiteral)
  if (ts.isObjectLiteralExpression(n)) {
    const o = {}
    for (const p of n.properties) {
      if (!ts.isPropertyAssignment(p)) continue
      const k = ts.isIdentifier(p.name) || ts.isStringLiteral(p.name) || ts.isNumericLiteral(p.name) ? p.name.text : undefined
      if (k !== undefined) o[k] = evalLiteral(p.initializer)
    }
    return o
  }
  if (ts.isParenthesizedExpression(n) || ts.isAsExpression(n) || ts.isSatisfiesExpression?.(n)) return evalLiteral(n.expression)
  return undefined
}

// cva variant functions by name: every output for every variant combination.
const cvaFns = new Map()
for (const { sf } of sources) {
  const visit = (n) => {
    if (ts.isVariableDeclaration(n) && n.initializer && ts.isCallExpression(n.initializer) && ts.isIdentifier(n.initializer.expression) && n.initializer.expression.text === "cva" && ts.isIdentifier(n.name)) {
      const [baseN, cfgN] = n.initializer.arguments
      const base = baseN ? evalLiteral(baseN) : undefined
      const cfg = cfgN ? evalLiteral(cfgN) : undefined
      const fn = cva(base, cfg)
      const variants = cfg?.variants ?? {}
      let combos = [{}]
      for (const [k, vals] of Object.entries(variants)) {
        const next = []
        for (const c of combos) {
          next.push(c) // default
          for (const v of Object.keys(vals)) next.push({ ...c, [k]: v === "true" ? true : v === "false" ? false : v })
        }
        combos = next
      }
      const outs = uniq(combos.map((c) => fn(c)))
      cvaFns.set(n.name.text, outs)
      outs.forEach((s) => allStrings.add(s))
    }
    ts.forEachChild(n, visit)
  }
  visit(sf)
}

// className values passed to components, keyed "uiFile#Tag".
const classNamesFor = new Map()
const addClassName = (key, v) => {
  if (!classNamesFor.has(key)) classNamesFor.set(key, new Set())
  classNamesFor.get(key).add(v)
}

// Possible values of an expression used as a cn() argument.
const values = (n, slot) => {
  if (!n) return [""]
  if (ts.isStringLiteral(n) || ts.isNoSubstitutionTemplateLiteral(n)) return [n.text]
  if (ts.isParenthesizedExpression(n) || ts.isAsExpression(n) || ts.isNonNullExpression(n)) return values(n.expression, slot)
  if (ts.isConditionalExpression(n)) return uniq([...values(n.whenTrue, slot), ...values(n.whenFalse, slot)]).slice(0, CAP)
  if (ts.isBinaryExpression(n)) {
    const op = n.operatorToken.kind
    if (op === ts.SyntaxKind.AmpersandAmpersandToken) return uniq(["", ...values(n.right, slot)]).slice(0, CAP)
    if (op === ts.SyntaxKind.BarBarToken || op === ts.SyntaxKind.QuestionQuestionToken) return uniq([...values(n.left, slot), ...values(n.right, slot)]).slice(0, CAP)
    if (op === ts.SyntaxKind.PlusToken) {
      const out = []
      for (const a of values(n.left, slot)) for (const b of values(n.right, slot)) out.push(a + b)
      return uniq(out).slice(0, CAP)
    }
  }
  if (ts.isTemplateExpression(n)) {
    let outs = [n.head.text]
    for (const span of n.templateSpans) {
      const vs = values(span.expression, slot)
      const next = []
      for (const o of outs) for (const v of vs.length ? vs : [""]) next.push(o + v + span.literal.text)
      outs = sample(uniq(next), CAP)
    }
    return outs
  }
  if (ts.isCallExpression(n)) {
    const callee = n.expression
    if (ts.isIdentifier(callee) && cvaFns.has(callee.text)) {
      // cva fn ({ ..., className }): every combination, plus className slot
      const outs = cvaFns.get(callee.text)
      let cls = [""]
      const arg = n.arguments[0]
      if (arg && ts.isObjectLiteralExpression(arg) && arg.properties.some((p) => p.name?.text === "className" || p.name?.text === "class")) cls = slot()
      const res = []
      for (const o of outs) for (const c of cls) res.push(c ? o + " " + c : o)
      return sample(uniq(res), CAP * 2)
    }
    if (ts.isIdentifier(callee) && callee.text === "cn") return cnCombos(n, slot)
    return [""]
  }
  if (ts.isIdentifier(n) && (n.text === "className" || n.text === "class")) return slot()
  if (ts.isArrayLiteralExpression(n)) {
    // clsx arrays: join each element's first few values
    let outs = [""]
    for (const e of n.elements) {
      const vs = values(e, slot)
      const next = []
      for (const o of outs) for (const v of vs) next.push(o && v ? o + " " + v : o || v)
      outs = sample(uniq(next), CAP)
    }
    return outs
  }
  // Unknown expression (props, variables): falsy, or any literal inside it.
  const lits = [""]
  const visit = (m) => {
    if (ts.isStringLiteral(m) || ts.isNoSubstitutionTemplateLiteral(m)) lits.push(m.text)
    ts.forEachChild(m, visit)
  }
  visit(n)
  return uniq(lits).slice(0, CAP)
}

// All argument lists of one cn() call (cartesian product, sampled).
const cnArgLists = (call, slot) => {
  let lists = [[]]
  for (const a of call.arguments) {
    const vs = values(a, slot)
    const next = []
    for (const l of lists) for (const v of vs) next.push([...l, v])
    lists = next.length > CAP * 4 ? sample(next, CAP * 4) : next
  }
  return lists
}
const cnCombos = (call, slot) => uniq(cnArgLists(call, slot).map((l) => l.filter(Boolean).join(" ")))

// Pass 1: classNames given to imported ui components, and every literal.
const uiImports = (sf) => {
  const m = new Map() // local name -> "file#Export"
  for (const st of sf.statements) {
    if (!ts.isImportDeclaration(st) || !st.importClause?.namedBindings || !ts.isNamedImports(st.importClause.namedBindings)) continue
    const spec = st.moduleSpecifier.text
    const mm = spec.match(/registry\/new-york-v4\/ui\/([\w-]+)$/)
    if (!mm) continue
    for (const el of st.importClause.namedBindings.elements) m.set(el.name.text, mm[1] + "#" + (el.propertyName ?? el.name).text)
  }
  return m
}
for (const { sf, rel } of sources) {
  const imports = uiImports(sf)
  const self = rel.startsWith("ui/") ? basename(rel).replace(/\.tsx?$/, "") : null
  const visit = (n) => {
    if (ts.isStringLiteral(n) || ts.isNoSubstitutionTemplateLiteral(n)) allStrings.add(n.text)
    if (ts.isTemplateExpression(n)) {
      allStrings.add(n.head.text)
      n.templateSpans.forEach((s) => allStrings.add(s.literal.text))
    }
    if ((ts.isJsxOpeningElement(n) || ts.isJsxSelfClosingElement(n)) && ts.isIdentifier(n.tagName)) {
      const tag = n.tagName.text
      const key = imports.get(tag) ?? (self ? self + "#" + tag : null)
      if (key) {
        for (const attr of n.attributes.properties) {
          if (!ts.isJsxAttribute(attr) || attr.name.text !== "className" || !attr.initializer) continue
          const init = attr.initializer
          const vs = ts.isStringLiteral(init) ? [init.text] : ts.isJsxExpression(init) ? values(init.expression, () => [""]) : []
          vs.forEach((v) => v && addClassName(key, v))
        }
      }
    }
    ts.forEachChild(n, visit)
  }
  visit(sf)
}
const allStringList = [...allStrings].filter((s) => s.length < 2000)

// Pass 2: cn() calls.
const cases = new Map() // JSON(args) -> args
const addCase = (args) => {
  const k = JSON.stringify(args)
  if (!cases.has(k)) cases.set(k, args)
}
const stats = { registryCalls: 0 }
for (const { sf, rel } of sources) {
  const self = rel.startsWith("ui/") ? basename(rel).replace(/\.tsx?$/, "") : null
  const visitFn = (n, fnName) => {
    let name = fnName
    if (ts.isFunctionDeclaration(n) && n.name) name = n.name.text
    if (ts.isVariableDeclaration(n) && ts.isIdentifier(n.name) && n.initializer && (ts.isArrowFunction(n.initializer) || ts.isFunctionExpression(n.initializer) || ts.isCallExpression(n.initializer))) name = n.name.text
    if (ts.isCallExpression(n) && ts.isIdentifier(n.expression) && n.expression.text === "cn") {
      stats.registryCalls++
      const given = self && name ? classNamesFor.get(self + "#" + name) : undefined
      const slotVals = uniq(["", ...sample([...(given ?? [])], CAP), ...sample(allStringList, 4)])
      const slot = () => slotVals
      for (const l of cnArgLists(n, slot)) addCase(l)
    }
    ts.forEachChild(n, (c) => visitFn(c, name))
  }
  visitFn(sf, null)
}
// Every literal and className on its own, and cva outputs with each className.
for (const s of allStringList) addCase([s])
for (const set of classNamesFor.values()) for (const s of set) addCase([s])
const registryCount = cases.size

// --- pairs ----------------------------------------------------------------------
const classy = allStringList.filter((s) => /^[\w!@&\[\]():\/.,%*#>~+=<'"^$|-]+(\s+[\w!@&\[\]():\/.,%*#>~+=<'"^$|-]+)*$/.test(s) && /[-:]/.test(s))
for (let i = 0; i < 6000; i++) addCase([pick(classy), pick(classy)])
for (let i = 0; i < 2000; i++) addCase([pick(classy), pick(classy), pick(classy)])
const pairsCount = cases.size - registryCount

// --- random ---------------------------------------------------------------------
const vocab = uniq(classy.flatMap((s) => s.split(/\s+/)).filter(Boolean))
const variants = ["hover", "focus", "dark", "md", "lg", "sm", "group-hover", "peer-focus", "data-[state=open]", "aria-invalid", "has-[>svg]", "[&_svg]", "[&>*]", "*", "**", "before", "after", "placeholder", "first", "last", "@md", "@container/main", "supports-[display:grid]", "not-[&:hover]", "group-data-[collapsible=icon]", "max-sm", "min-[400px]", "in-data-[side=left]", "file", "selection", "marker", "backdrop", "first-letter", "details-content", "open", "motion-safe", "print", "rtl", "ltr", "starting", "inert", "nth-3", "[@media(min-width:900px)]", "[.dark_&]", "group-[.is-active]"]
const prefixes = ["p", "px", "py", "pt", "pr", "pb", "pl", "ps", "pe", "m", "mx", "my", "mt", "-mt", "ms", "me", "gap", "gap-x", "gap-y", "space-x", "space-y", "w", "h", "size", "min-w", "max-w", "min-h", "max-h", "inset", "inset-x", "top", "-top", "left", "right", "bottom", "z", "-z", "order", "basis", "grow", "shrink", "flex", "col-span", "col-start", "row-span", "grid-cols", "grid-rows", "text", "font", "leading", "tracking", "indent", "-indent", "line-clamp", "bg", "from", "via", "to", "fill", "stroke", "stroke-w", "border", "border-x", "border-y", "border-t", "border-r", "border-b", "border-l", "border-s", "border-e", "rounded", "rounded-t", "rounded-tl", "rounded-s", "rounded-ss", "outline", "outline-offset", "ring", "ring-offset", "inset-ring", "shadow", "inset-shadow", "drop-shadow", "text-shadow", "opacity", "blur", "backdrop-blur", "brightness", "contrast", "grayscale", "hue-rotate", "-hue-rotate", "invert", "saturate", "sepia", "scale", "scale-x", "-scale-y", "rotate", "-rotate", "rotate-x", "skew", "-skew-x", "translate", "translate-x", "-translate-y", "translate-z", "origin", "perspective", "duration", "delay", "ease", "animate", "transition", "columns", "aspect", "object", "overflow", "overscroll", "scroll-m", "scroll-mx", "scroll-p", "snap", "touch", "select", "cursor", "caret", "accent", "decoration", "underline-offset", "list", "content", "align", "whitespace", "break", "hyphens", "mask", "mask-linear", "bg-linear", "bg-radial", "bg-conic", "divide", "divide-x", "divide-y", "placeholder", "zoom", "will-change", "field-sizing", "container", "sr", "not-sr", "shrink-0", "items", "justify", "self", "place", "content", "isolation", "mix-blend", "bg-blend", "table", "caption", "transform", "backface", "box", "float", "clear", "display", "font-stretch", "tab", "wrap", "text-wrap", "truncate", "underline", "line-through", "no-underline", "uppercase", "italic", "fill-current"]
const values_ = ["0", "0.5", "1", "1.5", "2", "3", "4", "5", "6", "8", "10", "12", "16", "24", "32", "40", "64", "96", "px", "auto", "full", "screen", "min", "max", "fit", "dvh", "svw", "lh", "1/2", "2/3", "3/4", "1/3.5", "0.5/2", "xs", "sm", "md", "lg", "xl", "2xl", "3xl", "1.5xl", "none", "inherit", "current", "transparent", "black", "white", "primary", "primary-foreground", "red-500", "red-500/50", "blue-600/[.3]", "zinc-950/10", "accent/50", "foreground", "muted-foreground", "border", "ring/50", "[3px]", "[10%]", "[calc(100%-2rem)]", "[min(10px,2vw)]", "[length:var(--x)]", "[color:red]", "[#fff]", "[rgb(0,0,0)]", "[hsl(0_0%_0%)]", "[oklch(0.5_0.2_100)]", "[color-mix(in_oklab,red,blue)]", "[url(/a.png)]", "[linear-gradient(red,blue)]", "[image:var(--i)]", "[0_0_0_1px_red]", "[inset_0_1px_0_red]", "[2px_2px]", "[shadow:var(--s)]", "[number:1]", "[family-name:Inter]", "[weight:500]", "[position:center]", "[percentage:50%]", "[size:10px]", "[bg-size:cover]", "[1.5]", "[100]", "[50%]", "[auto]", "[var(--x)]", "[theme(spacing.4)]", "[--my-var]", "(--my-var)", "(length:--x)", "(color:--c)", "(image:--i)", "(--x)", "[0]", "[-1px]", "[.5em]", "[2rem]", "[10vh]", "[3cqw]", "[1lh]", "[5ch]", "[1e3px]", "[clamp(1rem,2vw,3rem)]", "1e3", "1e-3", "0x10", "0b11", "0o7", "Infinity", "-Infinity", "+5", ".5", "5.", "1_000", "NaN", "1.5e2", "07", "1/0", "100%", "50.5%", "1e2%", "%", "@container", "@container/main", "@container-size/x", "@container-normal/x", "@container-size/", "linear", "radial", "conic", "linear/oklch", "linear-to-r/srgb", "[abc]", "[a:b]", "[:b]", "[a:]"]
const extras = ["[mask-type:luminance]", "[mask-type:alpha]", "[--x:1]", "[--x:2]", "[color:red]", "[color:blue]", "[&>svg]:size-4", "[a]", "[]", "[:]", "[x:]", "!", "!!", "-", "--", "/", "//", ":", "::", "!p-2", "p-2!", "!-mt-2", "-mt-2!", "!w-1/2", "w-1/2!", "hover:!p-4", "dark:hover:p-4", "hover:dark:p-4", "*:p-2", "**:p-2", "before:*:p-2", "*:before:p-2", "[&_svg]:p-2", "p-[3px]/50", "bg-red-500/[.5]", "text-sm/6", "text-sm/[1.5]", "text-[14px]/tight", "leading-6", "text-lg/7", "font-[500]", "font-[Inter]", "font-(family-name:--f)", "text-(length:--l)", "text-(color:--c)", "bg-(image:--i)", "shadow-(--s)", "shadow-[0_0_0_1px_red]", "shadow-xs/20", "drop-shadow-md/50", "inset-shadow-sm/10", "ring-[3px]", "ring-ring/50", "outline-hidden", "outline-dashed", "outline-2", "outline-red-500", "border-2", "border-dashed", "border-[3px]", "border-[red]", "border-x-[2px]", "rounded-[2px]", "rounded-t-lg", "rounded-tl-none", "size-[calc(100%-1px)]", "w-(--sidebar-width)", "h-[var(--radix-select-trigger-height)]", "min-w-[8rem]", "max-h-(--radix-dropdown-menu-content-available-height)", "origin-(--radix-popover-content-transform-origin)", "translate-x-[-50%]", "translate-y-[calc(-50%_-_2px)]", "rotate-45", "-rotate-90", "rotate-[10deg]", "scale-95", "zoom-in-95", "fade-in-0", "slide-in-from-top-2", "animate-in", "animate-out", "data-[state=open]:animate-in", "duration-200", "ease-[cubic-bezier(0.4,0,0.2,1)]", "grid-cols-[auto_1fr]", "col-span-full", "row-start-2", "col-start-1", "order-first", "-order-1", "z-50", "-z-10", "z-[100]", "opacity-0", "opacity-[.5]", "opacity-50", "blur-sm", "blur-[2px]", "grayscale", "grayscale-0", "invert", "sepia-[.5]", "flex-1", "flex-[2_2_0%]", "flex-auto", "grow-0", "shrink", "basis-1/2", "line-clamp-2", "line-clamp-none", "tracking-tight", "leading-none", "indent-4", "-indent-px", "decoration-2", "decoration-wavy", "decoration-red-500", "underline-offset-4", "caret-transparent", "accent-red-500", "fill-none", "stroke-2", "stroke-[1.5]", "stroke-red-500", "content-['']", "content-none", "align-middle", "align-[4px]", "list-disc", "list-inside", "list-image-none", "columns-3", "columns-[10rem]", "aspect-video", "aspect-[4/3]", "aspect-4/3", "object-cover", "object-[25%_75%]", "bg-cover", "bg-[length:200px]", "bg-[position:center]", "bg-[size:10px]", "bg-[url(/a.png)]", "bg-linear-to-r", "bg-linear-45", "bg-radial-[at_25%_25%]", "bg-conic-180", "from-10%", "from-[10%]", "via-50%", "to-red-500", "bg-none", "bg-fixed", "bg-repeat-x", "bg-top-left", "bg-origin-border", "bg-clip-text", "mask-none", "mask-l-from-50%", "mask-radial-at-center", "divide-x-2", "divide-y-reverse", "divide-red-500", "space-x-4", "space-x-reverse", "gap-[2px]", "gap-x-2", "inset-0", "inset-x-0", "-inset-1", "top-[50%]", "start-0", "end-0", "left-1/2", "-left-1/2", "translate-x-1/2", "w-1/2", "w-[50%]", "w-screen", "w-dvw", "h-svh", "size-full", "min-h-svh", "max-w-prose", "max-w-screen-sm", "container", "@container", "@container/main", "@container-size/x", "@container-normal", "@container-normal/x", "@sm:p-4", "@md/main:p-4", "@max-md:p-2", "sr-only", "not-sr-only", "hidden", "block", "inline-flex", "contents", "table", "flow-root", "visible", "invisible", "collapse", "static", "relative", "absolute", "sticky", "fixed", "isolate", "isolation-auto", "overflow-hidden", "overflow-x-auto", "overscroll-none", "truncate", "text-ellipsis", "text-clip", "whitespace-nowrap", "break-words", "break-all", "wrap-break-word", "text-balance", "text-pretty", "hyphens-auto", "uppercase", "normal-case", "italic", "not-italic", "font-mono", "font-sans", "font-bold", "font-stretch-condensed", "font-stretch-50%", "antialiased", "tabular-nums", "ordinal", "slashed-zero", "normal-nums", "underline", "no-underline", "line-through", "overline", "transition", "transition-all", "transition-[color,box-shadow]", "transition-none", "transform", "transform-gpu", "transform-none", "will-change-transform", "pointer-events-none", "cursor-pointer", "cursor-[url(a.cur),_pointer]", "select-none", "resize-none", "resize-y", "touch-none", "touch-pan-x", "snap-x", "snap-mandatory", "snap-start", "scroll-smooth", "scroll-m-4", "scroll-mx-2", "scroll-p-[3px]", "appearance-none", "outline-none", "field-sizing-content", "forced-color-adjust-none", "backface-hidden", "perspective-near", "perspective-[100px]", "transform-3d", "rotate-x-12", "-rotate-y-6", "translate-z-4", "scale-z-50", "scale-3d", "box-border", "box-content", "float-left", "clear-both", "table-auto", "border-collapse", "border-spacing-2", "caption-top", "mix-blend-multiply", "bg-blend-darken", "text-shadow-lg", "text-shadow-red-500/50", "text-shadow-none", "shadow-none", "shadow-red-500", "ring-0", "ring-inset", "ring-offset-2", "ring-offset-background", "inset-ring-2", "inset-ring-red-500", "inset-shadow-red-500", "accent-auto", "scheme-dark", "color-scheme-dark", "tab-4", "text-[length:12px]", "text-[color:var(--c)]", "text-[12px]", "text-[red]", "text-[#000]", "text-[rgb(1,2,3)]", "text-[var(--x)]", "text-xs", "text-2xs", "text-1.5xl", "text-red-500", "text-current", "text-inherit", "text-left", "text-start", "text-justify", "font-[number:500]", "font-[weight:500]", "font-500", "font-[family-name:x]", "leading-[1.5]", "leading-[20px]", "tracking-[0.1em]", "p-0.5", "p-1.25", "p-[0.5rem]", "p-px", "p-(--x)", "p-[--x]", "px-[calc(1rem+1px)]", "m-auto", "-m-px", "-mx-1", "mx-auto", "my-0", "ms-2", "me-[3px]", "mt-[-1px]", "w-auto", "w-fit", "w-min", "w-max", "w-3xs", "w-xl", "w-sm", "w-[200px]", "h-[1px]", "h-px", "h-0", "size-4", "size-[17px]", "size-auto", "min-w-0", "max-w-[calc(100%-2rem)]", "max-w-lg", "max-w-none", "grid", "grid-cols-2", "grid-cols-[1fr_2fr]", "grid-cols-subgrid", "grid-rows-[auto_1fr]", "grid-flow-col", "auto-rows-min", "auto-cols-fr", "gap-1.5", "items-center", "items-start", "justify-between", "justify-center-safe", "content-center", "self-end", "place-items-center", "justify-items-start", "flex-col", "flex-row-reverse", "flex-wrap", "flex-nowrap", "order-2", "col-span-2", "col-[span_2]", "row-span-3", "col-auto", "row-auto"]
const unicodeBits = ["é", "ü", "日本", "🙂", " ", " ", " ", "﻿", "　", "​", "ﬃ", "Ä", "́", "𝟘"]
const synth = () => {
  const r = rnd()
  if (r < 0.35) return pick(vocab)
  if (r < 0.55) return pick(extras)
  if (r < 0.85) {
    let t = pick(prefixes)
    if (rnd() < 0.85) t += "-" + pick(values_)
    if (rnd() < 0.1) t = "-" + t
    if (rnd() < 0.1) t += "/" + pick(["50", "[.5]", "(--o)", "2", "0", "tight", "6"])
    return t
  }
  if (r < 0.93) {
    // unicode / odd numbers inside an otherwise valid class
    const t = pick(prefixes) + "-" + pick(["[", "(", ""]) + pick(unicodeBits) + pick(values_) + pick(["]", ")", ""])
    return t
  }
  return pick(unicodeBits) + pick(["", pick(values_)])
}
const decorate = (t) => {
  const r = rnd()
  let out = t
  if (r < 0.45) {
    const n = rnd() < 0.7 ? 1 : rnd() < 0.7 ? 2 : 3
    const mods = []
    for (let i = 0; i < n; i++) mods.push(pick(variants))
    out = mods.join(":") + ":" + out
  }
  const r2 = rnd()
  if (r2 < 0.06) out = "!" + out
  else if (r2 < 0.12) out = out + "!"
  return out
}
const seps = [" ", " ", " ", " ", "  ", "\t", "\n", " \n ", "\r\n", "\f", "\v", "   "]
const randomCase = () => {
  const nArgs = rnd() < 0.6 ? 1 : 1 + Math.floor(rnd() * 4)
  const args = []
  for (let a = 0; a < nArgs; a++) {
    const nTok = Math.floor(rnd() * 9)
    let s = rnd() < 0.08 ? pick(seps) : ""
    for (let i = 0; i < nTok; i++) {
      if (i > 0) s += pick(seps)
      s += decorate(synth())
    }
    if (rnd() < 0.08) s += pick(seps)
    args.push(s)
  }
  return args
}
// Same-group conflicts are what matters, so also generate focused cases: one
// prefix, several values and variants.
const conflictCase = () => {
  const p = pick(prefixes)
  const v = rnd() < 0.5 ? pick(variants) + ":" : ""
  const toks = []
  const n = 2 + Math.floor(rnd() * 4)
  for (let i = 0; i < n; i++) {
    let t = p + "-" + pick(values_)
    if (rnd() < 0.3) t = pick(prefixes) + "-" + pick(values_)
    if (rnd() < 0.3) t = v + t
    toks.push(decorate(t))
  }
  return [toks.join(" ")]
}
for (let i = 0; i < 25000; i++) addCase(randomCase())
for (let i = 0; i < 15000; i++) addCase(conflictCase())
const randomCount = cases.size - registryCount - pairsCount

// --- results --------------------------------------------------------------------
const out = []
for (const args of cases.values()) {
  const r = cn(...args)
  const joined = args.filter(Boolean).join(" ")
  const u = uncached(joined)
  if (r !== u) throw new Error("cached and uncached results differ for " + JSON.stringify(args))
  out.push([r, ...args])
}
const doc = {
  comment: "Generated by cn/gen/corpus.mjs from the upstream registry and cn@" + cnPkg.version + "; each case is [cn(...args), ...args].",
  cn: cnPkg.version,
  counts: { registry: registryCount, pairs: pairsCount, random: randomCount, registryCnCalls: stats.registryCalls },
  cases: out,
}
const dst = join(root, "cn/testdata/corpus.json.gz")
const { cases: _, ...head } = doc
const text = JSON.stringify(head).slice(0, -1) + ',"cases":[\n' + out.map((c) => JSON.stringify(c)).join(",\n") + "\n]}\n"
writeFileSync(dst, gzipSync(text, { level: 9 }))
console.log("wrote", dst, doc.counts, out.length, "cases")
