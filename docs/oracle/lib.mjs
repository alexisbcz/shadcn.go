// Shared helpers for the docs oracles: they run upstream's actual MDX
// pipeline (fumadocs-mdx 15.0.12 with ui/apps/v4/source.config.ts) in Node.
// Verification tooling only; the Go packages never run this.
import { createRequire } from "node:module"
import fs from "node:fs"
import path from "node:path"
import { fileURLToPath, pathToFileURL } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
export const repo = path.resolve(here, "../..")
export const app = path.join(repo, "ui/apps/v4")
export const contentDir = path.join(app, "content/docs")
const pnpm = path.join(repo, "ui/node_modules/.pnpm")

// from resolves a package the way the package in `dir` would.
function from(dir) {
  const req = createRequire(path.join(fs.realpathSync(dir), "package.json"))
  return async (name) => {
    let p
    try {
      p = req.resolve(name)
    } catch {
      // ESM-only packages without a "require" condition.
      const base = path.join(fs.realpathSync(dir), "..", name)
      const pkg = JSON.parse(fs.readFileSync(path.join(base, "package.json"), "utf8"))
      let entry = pkg.exports?.["."]
      while (entry && typeof entry === "object") entry = entry.import ?? entry.default ?? entry.node
      p = path.join(base, entry ?? pkg.module ?? pkg.main ?? "index.js")
    }
    return import(pathToFileURL(p).href)
  }
}

export const fromApp = from(path.join(app, "node_modules/next"))
export const fromMdx = from(path.join(app, "node_modules/fumadocs-mdx"))
export const fromCore = from(path.join(app, "node_modules/fumadocs-core"))
export const fromPretty = from(path.join(app, "node_modules/rehype-pretty-code"))
export const pnpmImport = (dir, sub) => import(pathToFileURL(path.join(pnpm, dir, "node_modules", sub)).href)

export const { toHtml } = await pnpmImport("hast-util-to-html@9.0.5", "hast-util-to-html/index.js")

// listPages returns every MDX file under content/docs, relative and sorted.
export function listPages() {
  const out = []
  const walk = (d) => {
    for (const e of fs.readdirSync(d, { withFileTypes: true })) {
      const p = path.join(d, e.name)
      if (e.isDirectory()) walk(p)
      else if (e.name.endsWith(".mdx")) out.push(path.relative(contentDir, p))
    }
  }
  walk(contentDir)
  return out.sort()
}

// mdxOptions returns the MDX options fumadocs-mdx derives from
// source.config.ts (applyMdxPreset with the "bundler" environment, which is
// what next build uses).
export async function mdxOptions() {
  const mdxDist = path.join(fs.realpathSync(path.join(app, "node_modules/fumadocs-mdx")), "dist")
  const preset = await import(pathToFileURL(path.join(mdxDist, "preset-NZprPlqc.js")).href)
  const include = await import(pathToFileURL(path.join(mdxDist, "remark-include-DJN_j9Y9.js")).href)
  const config = await import(pathToFileURL(path.join(app, ".source/source.config.mjs")).href)
  const global = config.default
  const fn = global.mdxOptions
  const opts = await preset.t(typeof fn === "function" ? await fn() : fn)("bundler")
  return { opts, remarkInclude: include.t }
}

export async function frontmatter(source) {
  const m = await fromCore("fumadocs-core/content/md/frontmatter")
  return m.frontmatter(source)
}
