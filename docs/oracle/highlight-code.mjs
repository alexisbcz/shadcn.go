// highlight-code.mjs runs the docs' own highlightCode (ui/apps/v4/lib/
// highlight-code.ts, bundled with esbuild, with the app's shiki 3.23.0) on a
// corpus of the sources the docs components show: every radix example
// (examples/radix/*.tsx), every radix-nova component (styles/radix-nova/ui),
// the style sheets (registry/styles/*.css, as css) and the radix-nova
// registry items (public/r/styles/radix-nova/*.json, as json).
//
// Output: docs/highlight/testdata/code.json.gz ([{path, lang, code, html}])
// Usage: node docs/oracle/highlight-code.mjs [--extended]
import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import zlib from "node:zlib"
import { createRequire } from "node:module"
import { pathToFileURL } from "node:url"
import { repo, app } from "./lib.mjs"

const esbuild = createRequire(path.join(repo, "ui/node_modules/.pnpm/esbuild@0.25.11/node_modules/esbuild/package.json"))("esbuild")
const out = path.join(os.tmpdir(), "shadcn-go-docs-oracle", "highlight-code.mjs")
const shikiDir = fs.realpathSync(path.join(app, "node_modules/shiki"))
const lruDir = fs.realpathSync(path.join(app, "node_modules/lru-cache"))
await esbuild.build({
  entryPoints: [path.join(app, "lib/highlight-code.ts")],
  bundle: true,
  format: "esm",
  platform: "node",
  outfile: out,
  logLevel: "error",
  plugins: [
    {
      name: "externals",
      setup(build) {
        build.onResolve({ filter: /^shiki$/ }, () => ({ path: path.join(shikiDir, "dist/index.mjs"), external: true }))
        build.onResolve({ filter: /^lru-cache$/ }, () => {
          const pkg = JSON.parse(fs.readFileSync(path.join(lruDir, "package.json"), "utf8"))
          let entry = pkg.exports["."]
          while (typeof entry === "object") entry = entry.import ?? entry.default
          return { path: path.join(lruDir, entry), external: true }
        })
      },
    },
  ],
})
const { highlightCode } = await import(pathToFileURL(out).href)

const extended = process.argv[2] === "--extended"
const inputs = []
const add = (dir, ext, lang) => {
  for (const f of fs.readdirSync(path.join(app, dir)).sort()) {
    if (f.endsWith(ext)) inputs.push({ path: path.join("ui/apps/v4", dir, f), lang, file: path.join(app, dir, f) })
  }
}
if (!extended) {
  add("examples/radix", ".tsx", "tsx")
  add("styles/radix-nova/ui", ".tsx", "tsx")
  add("registry/styles", ".css", "css")
  add("public/r/styles/radix-nova", ".json", "json")
} else {
  // --extended: the other sources the docs components highlight, beyond
  // the radix corpus above: the base and aria examples, the components of
  // the styles ComponentSource uses on the docs pages (base-nova,
  // aria-nova and the three rhea styles), and a seeded sample of the other
  // registry items' files (blocks, charts, other styles) as BlockDisplay
  // and ChartDisplay highlight them (tsx). Written to code-extended.json.gz
  // with the SHA-256 of the expected HTML instead of the HTML, to keep the
  // testdata small.
  add("examples/base", ".tsx", "tsx")
  add("examples/aria", ".tsx", "tsx")
  for (const st of ["base-nova", "aria-nova", "radix-rhea", "base-rhea", "aria-rhea"]) add(`styles/${st}/ui`, ".tsx", "tsx")
  let seed = 7
  const rand = () => ((seed = (Math.imul(seed, 1103515245) + 12345) >>> 0) / 4294967296)
  const pool = []
  const rdir = path.join(app, "public/r/styles")
  for (const st of fs.readdirSync(rdir).sort()) {
    if (!fs.statSync(path.join(rdir, st)).isDirectory()) continue
    for (const f of fs.readdirSync(path.join(rdir, st)).sort()) {
      if (!f.endsWith(".json")) continue
      let item
      try {
        item = JSON.parse(fs.readFileSync(path.join(rdir, st, f), "utf8"))
      } catch {
        continue
      }
      for (const file of item.files ?? []) {
        if (file.content) pool.push({ path: `ui/apps/v4/public/r/styles/${st}/${f}#${file.path}`, lang: "tsx", code: file.content })
      }
    }
  }
  for (const it of pool) if (rand() < 0.05) inputs.push(it)
}
const seenCode = new Set()
for (let i = inputs.length - 1; i >= 0; i--) {
  if (!inputs[i].code) inputs[i].code = fs.readFileSync(inputs[i].file, "utf8")
}
// The extended corpus keeps each distinct source once.
const distinct = inputs.filter((it) => {
  const k = it.lang + "\0" + it.code
  if (extended && seenCode.has(k)) return false
  seenCode.add(k)
  return true
})
inputs.length = 0
inputs.push(...distinct)

// Warm-up (shiki's tokenize time limit, see build.mjs).
for (const it of inputs) await highlightCode(it.code + "\n// warm-up", it.lang)

const { createHash } = await import("node:crypto")
const list = []
for (const it of inputs) {
  const html = await highlightCode(it.code, it.lang)
  list.push(extended ? { path: it.path, lang: it.lang, code: it.code, sha256: createHash("sha256").update(html).digest("hex") } : { path: it.path, lang: it.lang, code: it.code, html })
}
const dest = path.join(repo, "docs/highlight/testdata", extended ? "code-extended.json.gz" : "code.json.gz")
fs.writeFileSync(dest, zlib.gzipSync(JSON.stringify(list) + "\n", { level: 9 }))
console.log(`${list.length} sources highlighted`)
