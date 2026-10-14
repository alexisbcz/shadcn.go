// extract.mjs records, for every generated style directory, the class-bearing
// parts of each upstream base source and of its generated counterpart, as
// reduced syntax trees (see "Format" below). The Go differential test
// (diff_test.go) replays upstream's transformStyle + prettier class sorting
// on the base trees and requires the generated trees exactly.
//
// Verification tooling only (Node.js + the upstream clone in ./ui):
//
//	node components/style/testdata/refgen/extract.mjs            # all 24 directories
//	node components/style/testdata/refgen/extract.mjs radix-nova # one
//
// Sources, for <base>-<style>:
//   - ui:     ui/apps/v4/registry/bases/<base>/ui/*.tsx -> ui/apps/v4/styles/<base>-<style>/ui/*.tsx
//   - others: every other file the build transforms (examples, blocks, hooks,
//     lib, components), base -> the build's transform output in
//     ui/apps/v4/node_modules/.cache/build-registry/transforms/<base>-<style>/.
//
// Format: { dir, files: [{ path, base: [root...], gen: [root...] }] }.
// A root is a class context, outermost first, in document order: a JSX
// className attribute, a call whose callee is (or, through call/member
// chains, starts with) cn, cva or mergeProps, or a `className:` property.
// Nodes:
//   {k:"str", v, t?:1}                string literal (t: no-substitution template)
//   {k:"tpl", q:[[raw,cooked]...], e:[node]}  template with substitutions
//   {k:"call", f:callee identifier|null, r:root identifier|null, a:[node]}
//   {k:"obj", p:[node]}                object literal
//   {k:"pa", n:name text, id:1|0, v:node, nk?:str}  property assignment
//                                      (nk: the key, when it is a string)
//   {k:"attr", n:name, v:null|node}    JSX attribute (value: str or jx)
//   {k:"jx", e:node|null}              JSX expression container
//   {k:"jsx", a:[node], c:[node]}      JSX element (tag names are not kept:
//                                      the icon transform renames them)
//   {k:"bin", op, l, r}                binary expression
//   {k:"id", v}                        identifier
//   {k:"g", s:kind, c:[node]}          any other node, children kept
// Parentheses are unwrapped (prettier adds and removes them).

import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"
import { createRequire } from "node:module"
import { fileURLToPath } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
const repo = path.resolve(here, "../../../..")
const v4 = path.join(repo, "ui/apps/v4")
const require = createRequire(path.join(v4, "package.json"))
const ts = require("typescript")

const BASES = ["radix", "base", "aria"]
const STYLES = ["nova", "vega", "lyra", "maia", "mira", "luma", "sera", "rhea"]
const cacheRoot = path.join(v4, "node_modules/.cache/build-registry/transforms")
const outDir = path.join(here, "../ref")

function calleeIdent(expr) {
  return ts.isIdentifier(expr) ? expr.text : null
}

function calleeRoot(expr) {
  let e = expr
  for (;;) {
    if (ts.isCallExpression(e)) e = e.expression
    else if (ts.isPropertyAccessExpression(e) || ts.isElementAccessExpression(e)) e = e.expression
    else break
  }
  return ts.isIdentifier(e) ? e.text : null
}

const ROOT_CALLEES = new Set(["cn", "cva", "mergeProps"])

function isRoot(node) {
  if (ts.isJsxAttribute(node)) return node.name.getText() === "className"
  if (ts.isCallExpression(node)) {
    const f = calleeIdent(node.expression)
    const r = calleeRoot(node.expression)
    return ROOT_CALLEES.has(f) || r === "cn" || r === "cva"
  }
  if (ts.isPropertyAssignment(node)) return node.name.getText() === "className"
  return false
}

function children(node) {
  const out = []
  ts.forEachChild(node, (c) => {
    out.push(c)
  })
  return out
}

function ser(node) {
  while (ts.isParenthesizedExpression(node)) node = node.expression
  if (ts.isStringLiteral(node)) return { k: "str", v: node.text }
  if (ts.isNoSubstitutionTemplateLiteral(node)) return { k: "str", v: node.text, t: 1 }
  if (ts.isTemplateExpression(node)) {
    const q = [[node.head.rawText ?? node.head.text, node.head.text]]
    const e = []
    for (const span of node.templateSpans) {
      e.push(ser(span.expression))
      q.push([span.literal.rawText ?? span.literal.text, span.literal.text])
    }
    return { k: "tpl", q, e }
  }
  if (ts.isCallExpression(node)) {
    return {
      k: "call",
      f: calleeIdent(node.expression),
      r: calleeRoot(node.expression),
      c: ser(node.expression),
      a: node.arguments.map(ser),
    }
  }
  if (ts.isObjectLiteralExpression(node)) return { k: "obj", p: node.properties.map(ser) }
  if (ts.isPropertyAssignment(node)) {
    const out = { k: "pa", n: node.name.getText(), id: ts.isIdentifier(node.name) ? 1 : 0, v: ser(node.initializer) }
    // A string key is a literal prettier sorts too (clsx object syntax).
    if (ts.isStringLiteral(node.name) || ts.isNoSubstitutionTemplateLiteral(node.name)) out.nk = ser(node.name)
    return out
  }
  if (ts.isJsxAttribute(node)) {
    return { k: "attr", n: node.name.getText(), v: node.initializer ? ser(node.initializer) : null }
  }
  if (ts.isJsxExpression(node)) return { k: "jx", e: node.expression ? ser(node.expression) : null }
  if (ts.isJsxElement(node)) {
    return {
      k: "jsx",
      a: node.openingElement.attributes.properties.map(ser),
      c: node.children.map(ser).filter(Boolean),
    }
  }
  if (ts.isJsxSelfClosingElement(node)) {
    return { k: "jsx", a: node.attributes.properties.map(ser), c: [] }
  }
  if (ts.isJsxText(node)) {
    // Prettier rewraps JSX text; whitespace-only text is not content.
    const t = node.text.replace(/\s+/g, " ").trim()
    return t ? { k: "g", s: "JsxText" } : null
  }
  if (ts.isBinaryExpression(node)) {
    return { k: "bin", op: node.operatorToken.getText(), l: ser(node.left), r: ser(node.right) }
  }
  if (ts.isIdentifier(node)) return { k: "id", v: node.text }
  const c = children(node).map(ser).filter(Boolean)
  const out = { k: "g", s: ts.SyntaxKind[node.kind] }
  if (c.length) out.c = c
  return out
}

function roots(file, text) {
  const sf = ts.createSourceFile(file, text, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX)
  const out = []
  const visit = (node) => {
    if (isRoot(node)) {
      const { line } = sf.getLineAndCharacterOfPosition(node.getStart(sf))
      out.push({ line: line + 1, n: ser(node) })
      return
    }
    ts.forEachChild(node, visit)
  }
  visit(sf)
  return out
}

function walk(dir, rel = "") {
  const out = []
  for (const e of fs.readdirSync(path.join(dir, rel), { withFileTypes: true })) {
    const r = rel ? `${rel}/${e.name}` : e.name
    if (e.isDirectory()) out.push(...walk(dir, r))
    else if (/\.(tsx|ts)$/.test(e.name)) out.push(r)
  }
  return out.sort()
}

function extractDir(name) {
  const [base] = name.split("-")
  const baseRoot = path.join(v4, "registry/bases", base)
  const files = []
  const pairs = []
  for (const f of walk(path.join(v4, "styles", name, "ui"))) {
    pairs.push([`ui/${f}`, path.join(baseRoot, "ui", f), path.join(v4, "styles", name, "ui", f)])
  }
  for (const f of walk(path.join(cacheRoot, name))) {
    if (f.startsWith("ui/")) continue
    pairs.push([f, path.join(baseRoot, f), path.join(cacheRoot, name, f)])
  }
  for (const [rel, b, g] of pairs) {
    const bt = fs.readFileSync(b, "utf8")
    const gt = fs.readFileSync(g, "utf8")
    files.push({ path: rel, base: roots(b, bt), gen: roots(g, gt) })
  }
  return { dir: name, files }
}

const only = process.argv.slice(2)
const dirs = []
for (const b of BASES) for (const s of STYLES) dirs.push(`${b}-${s}`)
fs.mkdirSync(outDir, { recursive: true })
for (const d of dirs) {
  if (only.length && !only.includes(d)) continue
  const data = extractDir(d)
  const json = JSON.stringify(data)
  fs.writeFileSync(path.join(outDir, `${d}.json.gz`), zlib.gzipSync(json, { level: 9 }))
  const roots = data.files.reduce((n, f) => n + f.base.length, 0)
  console.log(`${d}: ${data.files.length} files, ${roots} roots, ${json.length} bytes`)
}
