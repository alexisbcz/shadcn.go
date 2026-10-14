// order.mjs checks components/style/order.txt.gz (produced in Go by
// components/style/gen, in goja) against Tailwind's getClassOrder run by
// Node.js, with the design system loaded the way prettier-plugin-tailwindcss
// loads it for ui/apps/v4 (tailwindStylesheet ./app/globals.css, packages
// resolved from node_modules). Verification tooling only.
//
//	node components/style/testdata/refgen/order.mjs
import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"
import { createRequire } from "node:module"
import { fileURLToPath } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
const repo = path.resolve(here, "../../../..")
const v4 = path.join(repo, "ui/apps/v4")
const require = createRequire(path.join(v4, "package.json"))
const tw = require("tailwindcss")

const text = zlib.gunzipSync(fs.readFileSync(path.join(repo, "components/style/order.txt.gz"))).toString()
// Line 1 is a header, "# unknown" separates the two lists.
const lines = text.split("\n").slice(1, -1)
const known = []
const unknown = []
let inUnknown = false
for (const l of lines) {
  if (l === "# unknown") inUnknown = true
  else (inUnknown ? unknown : known).push(l)
}

function resolveCss(id, base) {
  if (id.startsWith("./") || id.startsWith("../")) {
    const p = path.resolve(base, id)
    return fs.existsSync(p) ? p : p + ".css"
  }
  const [pkg, ...rest] = id.split("/")
  const name = pkg.startsWith("@") ? `${pkg}/${rest.shift()}` : pkg
  let pj = null
  for (let d = base; ; d = path.dirname(d)) {
    const c = path.join(d, "node_modules", name, "package.json")
    if (fs.existsSync(c)) {
      pj = fs.realpathSync(c)
      break
    }
    if (path.dirname(d) === d) throw new Error(`cannot resolve ${id} from ${base}`)
  }
  const meta = JSON.parse(fs.readFileSync(pj, "utf8"))
  const sub = rest.length ? `./${rest.join("/")}` : "."
  let target = null
  const exp = meta.exports && meta.exports[sub]
  if (typeof exp === "string") target = exp
  else if (exp && typeof exp === "object") target = exp.style ?? exp.default ?? exp.import
  if (!target && sub === ".") target = meta.style ?? "index.css"
  if (!target) target = sub
  return path.join(path.dirname(pj), target)
}

const globals = path.join(v4, "app/globals.css")
const ds = await tw.__unstable__loadDesignSystem(fs.readFileSync(globals, "utf8"), {
  base: path.dirname(globals),
  async loadStylesheet(id, base) {
    const p = resolveCss(id, base)
    return { path: p, base: path.dirname(p), content: fs.readFileSync(p, "utf8") }
  },
  async loadModule() {
    throw new Error("unsupported")
  },
})

const all = [...known, ...unknown]
const order = ds.getClassOrder(all)
const got = order.filter(([, o]) => o !== null).sort((a, b) => (a[1] < b[1] ? -1 : a[1] > b[1] ? 1 : 0)).map(([n]) => n)
const gotUnknown = order.filter(([, o]) => o === null).map(([n]) => n)
let bad = 0
if (got.length !== known.length) {
  console.log(`known count differs: node ${got.length}, go ${known.length}`)
  bad++
}
for (let i = 0; i < Math.min(got.length, known.length); i++) {
  if (got[i] !== known[i]) {
    if (bad++ < 10) console.log(`rank ${i}: node ${got[i]}, go ${known[i]}`)
  }
}
const u = new Set(unknown)
for (const n of gotUnknown) if (!u.has(n)) if (bad++ < 20) console.log(`unknown in node only: ${n}`)
console.log(`${known.length} ordered classes, ${unknown.length} unknown: ${bad} differences`)
process.exit(bad ? 1 : 0)
