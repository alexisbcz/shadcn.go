// extra.mjs widens the differential corpus beyond the files the registry
// build happens to transform. It runs upstream's own transformStyle (from
// ui/packages/shadcn/dist/utils) and prettier with the ui/apps/v4 config,
// exactly as scripts/build-registry.mts does for each file, on:
//
//   - skipped:  every registry/bases/<base> .tsx/.ts source the build does
//               not transform (blocks/preview*, internal, ...), in all 8
//               styles;
//   - tokens:   synthetic sources that use every token of every style map
//               (most tokens appear in only a few upstream files), in each
//               kind of site the transform handles: className literals,
//               cn() with a string or a non-string first argument, className
//               properties, mergeProps, cva base strings and variants, and
//               pairs of tokens in one string (combined style classes);
//   - examples: the docs examples ui/apps/v4/examples/<base>/*.tsx (no
//               tokens; prettier's class sorting only), with the nova map.
//
// Verification tooling only:
//
//	node components/style/testdata/refgen/extra.mjs
//
// It writes testdata/ref/extra.json.gz: { cases: [{ set, name, style,
// base: [root], gen: [root] | null }] }, trees as in extract.mjs; gen is
// null when the formatted output's roots equal the base's (checked here).
import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"
import { createRequire } from "node:module"
import { fileURLToPath, pathToFileURL } from "node:url"
import { isDeepStrictEqual } from "node:util"

const here = path.dirname(fileURLToPath(import.meta.url))
const repo = path.resolve(here, "../../../..")
const v4 = path.join(repo, "ui/apps/v4")
const pkg = path.join(repo, "ui/packages/shadcn")
const { createStyleMap, transformStyle } = await import(pathToFileURL(path.join(pkg, "dist/utils/index.js")).href)
const { rootsOf } = await import(pathToFileURL(path.join(here, "trees.mjs")).href)

const v4require = createRequire(path.join(v4, "package.json"))
const prettier = (await import(pathToFileURL(v4require.resolve("prettier")).href)).default
const prettierConfig = await prettier.resolveConfig(path.join(v4, "package.json"))
prettierConfig.plugins = prettierConfig.plugins.map((p) => v4require.resolve(p))

const BASES = ["radix", "base", "aria"]
const STYLES = ["nova", "vega", "lyra", "maia", "mira", "luma", "sera", "rhea"]
const maps = {}
for (const s of STYLES) {
  maps[s] = createStyleMap(fs.readFileSync(path.join(v4, `registry/styles/style-${s}.css`), "utf8"))
}

// Line numbers are not data: drop them so unchanged roots compare equal.
const strip = (roots) => roots.map((r) => r.n)

const cases = []
let unchanged = 0
async function run(set, name, style, source, file) {
  const out = await transformStyle(source, { styleMap: maps[style] })
  // The build rewrites @/registry/bases/<base>/ imports; no class effect.
  const formatted = await prettier.format(out, { ...prettierConfig, filepath: path.join(v4, "registry", file) })
  const base = rootsOf(file, source)
  const gen = rootsOf(file, formatted)
  if (base.length !== gen.length) throw new Error(`${set} ${name} ${style}: ${base.length} roots, ${gen.length} after`)
  const same = isDeepStrictEqual(strip(base), strip(gen))
  if (same) unchanged++
  if (base.length) cases.push({ set, name, style, base, gen: same ? null : gen })
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

// 1. Base sources the build skips.
const cache = path.join(v4, "node_modules/.cache/build-registry/transforms")
let skipped = 0
for (const b of BASES) {
  const root = path.join(v4, "registry/bases", b)
  for (const f of walk(root)) {
    if (fs.existsSync(path.join(cache, `${b}-nova`, f))) continue
    skipped++
    const src = fs.readFileSync(path.join(root, f), "utf8")
    for (const s of STYLES) await run("skipped", `${b}/${f}`, s, src, f.endsWith(".ts") ? "case.ts" : "case.tsx")
  }
}

// 2. Synthetic sources over every token of every style map.
let rnd = 12345
const next = (n) => {
  rnd = (Math.imul(rnd, 1103515245) + 12345) >>> 0
  return (rnd >>> 16) % n
}
let tokenFiles = 0
for (const s of STYLES) {
  const tokens = Object.keys(maps[s]).sort()
  const chunks = []
  for (let i = 0; i < tokens.length; i += 40) chunks.push(tokens.slice(i, i + 40))
  for (const [ci, chunk] of chunks.entries()) {
    // Attribute and cn() sites, no cva in the file.
    const lines = chunk.map(
      (t, i) => `      <div className="${t} p-2 text-sm" />
      <div className={"${t} flex"} />
      <div className={cn("${t} items-center gap-1", className)} />
      <div className={cn(className, "${t} m-1", active && "z-${i}")} />
      <div className={cn(\`${t}\`, "${chunk[(i + 1) % chunk.length]} rounded-md")} />
      <div className={active ? cn("${t} border") : undefined} />
      <div className="${t} ${tokens[next(tokens.length)]} ${tokens[next(tokens.length)]}" />`
    )
    const props = chunk.map((t) => `  { className: cn("${t} px-3", className) },`)
    const merge = chunk.map((t) => `  mergeProps({ className: cn("${t} py-1") }, props),`)
    const src = `import { cn } from "@/lib/utils"

export function Sites({ className, active }: { className?: string; active?: boolean }) {
  return (
    <>
${lines.join("\n")}
    </>
  )
}

export const props = [
${props.join("\n")}
]

export const merged = [
${merge.join("\n")}
]
`
    await run("tokens", `sites-${ci}`, s, src, "case.tsx")
    // cva(): base strings and variants claim tokens for the rest of the file.
    const cvas = chunk.map(
      (t, i) => `const v${i} = cva("${t} inline-flex", {
  variants: {
    size: { sm: "${chunk[(i + 1) % chunk.length]} h-8", lg: "${t} h-10" },
    tone: { a: "${tokens[next(tokens.length)]}" },
  },
  compoundVariants: [{ size: "sm", tone: "a", className: "z-1 ${t}" }],
})`
    )
    const after = chunk.map((t) => `      <div className={cn("${t} w-full", className)} />`)
    const cvaSrc = `import { cva } from "class-variance-authority"
import { cn } from "@/lib/utils"

${cvas.join("\n")}

export function After({ className }: { className?: string }) {
  return (
    <>
${after.join("\n")}
    </>
  )
}
`
    await run("tokens", `cva-${ci}`, s, cvaSrc, "case.tsx")
    tokenFiles += 2
  }
}

// 3. Docs examples: prettier's sorting on their class strings.
let examples = 0
for (const b of BASES) {
  const root = path.join(v4, "examples", b)
  for (const f of walk(root)) {
    examples++
    await run("examples", `${b}/${f}`, "nova", fs.readFileSync(path.join(root, f), "utf8"), "case.tsx")
  }
}

const outFile = path.join(here, "../ref/extra.json.gz")
fs.writeFileSync(outFile, zlib.gzipSync(JSON.stringify({ cases }), { level: 9 }))
const roots = cases.reduce((n, c) => n + c.base.length, 0)
console.log(
  `skipped sources: ${skipped} (x8 styles); token files: ${tokenFiles}; examples: ${examples}; ` +
    `${cases.length} cases with roots, ${roots} roots, ${unchanged} files unchanged by the build`
)
