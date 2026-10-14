// sources.mts runs the docs' own source pipeline (lib/registry.ts getDemoItem
// / getRegistryItem, lib/read-file.ts and lib/format-code.ts, which applies
// transformStyle, transformIcons, transformMenu and transformRender) on every
// source a docs page shows through ComponentPreview's code panel or
// ComponentSource, and writes the formatted code (before highlighting) to
// docs/site/mdx/sources.json.gz:
//   {"<styleName>\x00<name>": code, "src\x00<styleName>\x00<src>": code}
//
// Generation tooling only (needs ui/ installed); the Go server embeds the
// output. Usage, from the repository root:
//   (cd ui/apps/v4 && NODE_OPTIONS=--conditions=react-server \
//     ../../node_modules/.bin/tsx ../../../docs/site/mdx/gen/sources.mts)
import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"

const app = process.cwd()
const repo = path.resolve(app, "../..", "..")
const { formatCode } = await import(app + "/lib/format-code.ts")
const { getDemoItem, getRegistryItem } = await import(app + "/lib/registry.ts")
const { readFileFromRoot } = await import(app + "/lib/read-file.ts")

function walk(dir: string, out: string[] = []) {
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, e.name)
    if (e.isDirectory()) walk(p, out)
    else if (p.endsWith(".mdx")) out.push(p)
  }
  return out
}

const attr = (a: string, k: string) => a.match(new RegExp(`\\b${k}="([^"]*)"`))?.[1]
const wanted = new Map<string, () => Promise<string | undefined>>()
for (const f of walk(path.join(app, "content"))) {
  const s = fs.readFileSync(f, "utf8")
  for (const m of s.matchAll(/<(ComponentPreview|ComponentSource)\b([\s\S]*?)\/?>/g)) {
    const [, tag, a] = m
    if (tag === "ComponentPreview" && attr(a, "type") === "block") continue
    const styleName = attr(a, "styleName") ?? "base-nova"
    const name = attr(a, "name")
    const src = tag === "ComponentSource" ? attr(a, "src") : undefined
    if (src) {
      wanted.set(`src\x00${styleName}\x00${src}`, async () => readFileFromRoot(src))
    } else if (name) {
      wanted.set(`${styleName}\x00${name}`, async () => {
        const item = (await getDemoItem(name, styleName)) ?? (await getRegistryItem(name, styleName))
        return item?.files?.[0]?.content
      })
    }
  }
}

const out: Record<string, string> = {}
for (const [key, load] of [...wanted].sort(([a], [b]) => (a < b ? -1 : 1))) {
  let code: string | undefined
  try {
    code = await load()
  } catch {
    code = undefined
  }
  if (!code) continue
  const styleName = key.startsWith("src\x00") ? key.split("\x00")[1] : key.split("\x00")[0]
  code = await formatCode(code, styleName)
  code = code.replaceAll("/* eslint-disable react/no-children-prop */\n", "")
  out[key] = code
}
const dest = path.join(repo, "docs/site/mdx/sources.json.gz")
fs.writeFileSync(dest, zlib.gzipSync(JSON.stringify(out), { level: 9 }))
console.log(`${Object.keys(out).length}/${wanted.size} sources -> ${dest}`)
