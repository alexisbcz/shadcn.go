// upstream.mjs runs upstream's own createStyleMap and transformStyle (from
// the built shadcn package, ui/packages/shadcn/dist/utils) and records their
// results for the Go tests. Verification tooling only.
//
//	node components/style/testdata/refgen/upstream.mjs
//
// It writes testdata/ref/upstream.json.gz:
//   - stylemaps: createStyleMap of each registry/styles/style-<name>.css;
//   - cases: the sources of upstream's transform-style-map.test.ts and
//     create-style-map.test.ts, plus extra edge cases (cases.tsx), each with
//     its style map and the reduced trees (see extract.mjs) of the source
//     before and after transformStyle (and after prettier), and the
//     createStyleMap edge cases of cases.css.
import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"
import { createRequire } from "node:module"
import { fileURLToPath, pathToFileURL } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
const repo = path.resolve(here, "../../../..")
const pkg = path.join(repo, "ui/packages/shadcn")
const require = createRequire(path.join(pkg, "package.json"))
const ts = require("typescript")
const { createStyleMap, transformStyle } = await import(pathToFileURL(path.join(pkg, "dist/utils/index.js")).href)
const { rootsOf } = await import(pathToFileURL(path.join(here, "trees.mjs")).href)

const STYLES = ["nova", "vega", "lyra", "maia", "mira", "luma", "sera", "rhea"]
const stylemaps = {}
for (const s of STYLES) {
  const css = fs.readFileSync(path.join(repo, `ui/apps/v4/registry/styles/style-${s}.css`), "utf8")
  stylemaps[s] = createStyleMap(css)
}

// The test files' cases: every `it(...)` whose body declares `source` (a
// template literal) and calls applyTransform(source, <map>), and every
// createStyleMap(css) case.
function testCases(file, baseStyleMap) {
  const text = fs.readFileSync(file, "utf8")
  const sf = ts.createSourceFile(file, text, ts.ScriptTarget.Latest, true, ts.ScriptKind.TS)
  const out = []
  const visit = (node) => {
    if (ts.isCallExpression(node) && node.expression.getText() === "it" && node.arguments.length >= 2) {
      const name = node.arguments[0].text
      const vars = {}
      let mapExpr = null
      const body = node.arguments[1]
      const walk = (n) => {
        if (ts.isVariableDeclaration(n) && n.initializer) {
          const init = n.initializer
          if (ts.isNoSubstitutionTemplateLiteral(init) || ts.isStringLiteral(init)) vars[n.name.getText()] = init.text
          else if (ts.isObjectLiteralExpression(init)) vars[n.name.getText()] = Function(`return (${init.getText()})`)()
        }
        if (ts.isCallExpression(n) && n.expression.getText() === "applyTransform") mapExpr = n.arguments[1].getText()
        ts.forEachChild(n, walk)
      }
      walk(body)
      if (vars.source !== undefined && mapExpr) {
        out.push({ name, source: vars.source, styleMap: mapExpr === "baseStyleMap" ? baseStyleMap : (vars[mapExpr] ?? Function(`return (${mapExpr})`)()) })
      }
      if (vars.css !== undefined) out.push({ name, css: vars.css })
    }
    ts.forEachChild(node, visit)
  }
  visit(sf)
  return out
}

const srcDir = path.join(pkg, "src/styles")
const cases = [
  ...testCases(path.join(srcDir, "transform-style-map.test.ts"), { "cn-foo": "bg-background gap-4 rounded-xl" }),
  ...testCases(path.join(srcDir, "create-style-map.test.ts")),
]
// Extra edge cases, one source per "// ---" section of cases.tsx, run with
// the nova map.
const extra = fs.readFileSync(path.join(here, "cases.tsx"), "utf8").split(/^\/\/ ---.*$/m)
extra.forEach((source, i) => {
  if (source.trim()) cases.push({ name: `cases.tsx #${i}`, source, styleMap: stylemaps.nova })
})
// createStyleMap edge cases, one stylesheet per "/* --- name */" section of
// cases.css.
const cssCases = fs.readFileSync(path.join(here, "cases.css"), "utf8").split(/^\/\* --- (.*) \*\/$/m)
for (let i = 1; i < cssCases.length; i += 2) cases.push({ name: `cases.css: ${cssCases[i]}`, css: cssCases[i + 1] })

// Prettier as the registry build runs it (ui/apps/v4/package.json config,
// with prettier-plugin-tailwindcss), on a path inside ui/apps/v4.
const v4 = path.join(repo, "ui/apps/v4")
const v4require = createRequire(path.join(v4, "package.json"))
const prettier = (await import(pathToFileURL(v4require.resolve("prettier")).href)).default
const prettierConfig = await prettier.resolveConfig(path.join(v4, "package.json"))
// Plugins are named in the config; resolve them from ui/apps/v4.
prettierConfig.plugins = prettierConfig.plugins.map((p) => v4require.resolve(p))

// Fuzz cases (deterministic): random class strings, separators, tokens of a
// random style and site shapes, through transformStyle and prettier like the
// other cases.
let seed = 20261014
const rand = (n) => {
  seed = (Math.imul(seed, 1103515245) + 12345) >>> 0
  return (seed >>> 16) % n
}
const pick = (a) => a[rand(a.length)]
const pool = [...new Set(Object.values(stylemaps).flatMap((m) => Object.values(m).join(" ").split(/\s+/)))].filter(Boolean).sort()
const tokenPool = Object.keys(stylemaps.nova).sort()
const oddPool = ["unknown-a", "foo", "...", "…", "cn-not-a-token", "cn-menu-target", "cn-font-heading", "cn-rtl-flip"]
const seps = [" ", "  ", "\t", "\n", " \n  ", "\f", "\r", " ", " ", "\v"]
function classString() {
  const n = rand(7)
  const parts = []
  for (let i = 0; i < n; i++) {
    const r = rand(10)
    parts.push(r < 5 ? pick(pool) : r < 8 ? pick(tokenPool) : pick(oddPool))
    if (rand(6) === 0 && parts.length) parts.push(parts[rand(parts.length)]) // duplicates
  }
  let s = rand(5) === 0 ? pick(seps) : ""
  parts.forEach((p, i) => (s += (i ? pick(seps) : "") + p))
  if (rand(5) === 0) s += pick(seps)
  return s
}
const lit = (s) => JSON.stringify(s)
for (let f = 0; f < 96; f++) {
  const st = STYLES[f % STYLES.length]
  const lines = []
  for (let i = 0; i < 40; i++) {
    const a = classString(), b = classString()
    switch (rand(8)) {
      case 0: lines.push(`      <div className={${lit(a)}} />`); break
      case 1: lines.push(`      <div className={cn(${lit(a)}, className)} />`); break
      case 2: lines.push(`      <div className={cn(className, ${lit(a)}, ${lit(b)})} />`); break
      case 3: lines.push(`      <div className={cn(${lit(a)}, ${lit(b)}, x && ${lit(classString())})} />`); break
      case 4: lines.push(`      <div className={x ? cn(${lit(a)}) : cn(${lit(b)}, className)} />`); break
      case 5: lines.push(`      <div className={\`${a.replace(/[`\\$]/g, "")} \${className} ${b.replace(/[`\\$]/g, "")}\`} />`); break
      case 6: lines.push(`      <div {...mergeProps({ className: cn(${lit(a)}, className) }, { className: cn(className, ${lit(b)}) })} />`); break
      case 7: lines.push(`      <div {...{ className: cn(${lit(a)}), other: ${lit(b)} }} />`); break
    }
  }
  const variants = []
  for (let i = 0; i < 4; i++) variants.push(`    v${i}: { a: ${lit(classString())}, b: ${lit(classString())} },`)
  const source = `import { cva } from "class-variance-authority"
import { cn } from "@/lib/utils"

const v = cva(${lit(classString())}, {
  variants: {
${variants.join("\n")}
  },
})

export function F({ className, x }: { className?: string; x?: boolean }) {
  return (
    <>
${lines.join("\n")}
    </>
  )
}
`
  cases.push({ name: `fuzz #${f} (${st})`, source, styleMap: stylemaps[st], style: st })
}

// The token regex and whitespace handling: upstream's extractCnClasses and
// removeCnClasses (with ALLOWLIST), transpiled from its source, on random
// strings.
const tsm = fs.readFileSync(path.join(srcDir, "transform-style-map.ts"), "utf8")
const tsf = ts.createSourceFile("t.ts", tsm, ts.ScriptTarget.Latest, true)
let fnText = ""
for (const st of tsf.statements) {
  if ((ts.isFunctionDeclaration(st) && ["extractCnClasses", "removeCnClasses"].includes(st.name?.text)) ||
      (ts.isVariableStatement(st) && st.getText().includes("ALLOWLIST"))) fnText += st.getText() + "\n"
}
const js = ts.transpile(fnText + "return { extractCnClasses, removeCnClasses }", { target: ts.ScriptTarget.ES2022 })
const { extractCnClasses, removeCnClasses } = Function(js)()
const alphabet = ["c", "n", "-", "cn-", "cn-", "a", "_", "1", " ", "\t", "\n", " ", " ", "﻿", "᠎", "\v", "é", ".", ":", "[", "/", "cn-menu-target", "cn-font-heading"]
const regex = []
for (let i = 0; i < 4000; i++) {
  let s = ""
  const n = rand(14)
  for (let j = 0; j < n; j++) s += pick(alphabet)
  regex.push({ s, extract: extractCnClasses(s), removed: removeCnClasses(s) })
}

for (const c of cases) {
  if (c.css !== undefined) {
    c.result = createStyleMap(c.css)
    continue
  }
  const out = await transformStyle(c.source, { styleMap: c.styleMap })
  c.output = out
  c.base = rootsOf("case.tsx", c.source)
  c.gen = rootsOf("case.tsx", out)
  try {
    const formatted = await prettier.format(out, { ...prettierConfig, filepath: path.join(v4, "registry/case.tsx") })
    c.formatted = rootsOf("case.tsx", formatted)
  } catch (e) {
    c.formatError = String(e.message).split("\n")[0]
  }
}

const outFile = path.join(here, "../ref/upstream.json.gz")
fs.writeFileSync(outFile, zlib.gzipSync(JSON.stringify({ stylemaps, cases, regex }), { level: 9 }))
console.log(
  `stylemaps: ${Object.entries(stylemaps).map(([k, v]) => `${k} ${Object.keys(v).length}`).join(", ")}; ` +
    `${cases.filter((c) => c.source).length} transform cases, ${cases.filter((c) => c.css).length} createStyleMap cases, ${regex.length} regex cases`
)
