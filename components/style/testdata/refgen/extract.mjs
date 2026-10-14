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
// Completeness: for every pair it also requires every string outside the
// roots (stringsOutsideRoots in trees.mjs) to be unchanged by the build, so
// the roots hold every string the build changes.
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
import { fileURLToPath } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
const repo = path.resolve(here, "../../../..")
const v4 = path.join(repo, "ui/apps/v4")
import { rootsOf, stringsOutsideRoots } from "./trees.mjs"
import { isDeepStrictEqual } from "node:util"

const BASES = ["radix", "base", "aria"]
const STYLES = ["nova", "vega", "lyra", "maia", "mira", "luma", "sera", "rhea"]
const cacheRoot = path.join(v4, "node_modules/.cache/build-registry/transforms")
const outDir = path.join(here, "../ref")

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
  let outside = 0
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
    files.push({ path: rel, base: rootsOf(b, bt), gen: rootsOf(g, gt) })
    const ob = stringsOutsideRoots(b, bt)
    const og = stringsOutsideRoots(g, gt)
    if (!isDeepStrictEqual(ob, og)) throw new Error(`${name}/${rel}: the build changed strings outside the class roots`)
    outside += ob.length
  }
  return { dir: name, files, outside }
}

const only = process.argv.slice(2)
const dirs = []
for (const b of BASES) for (const s of STYLES) dirs.push(`${b}-${s}`)
fs.mkdirSync(outDir, { recursive: true })
for (const d of dirs) {
  if (only.length && !only.includes(d)) continue
  const data = extractDir(d)
  const outside = data.outside
  delete data.outside
  const json = JSON.stringify(data)
  fs.writeFileSync(path.join(outDir, `${d}.json.gz`), zlib.gzipSync(json, { level: 9 }))
  const roots = data.files.reduce((n, f) => n + f.base.length, 0)
  console.log(`${d}: ${data.files.length} files, ${roots} roots, ${json.length} bytes; ${outside} strings outside the roots unchanged`)
}
