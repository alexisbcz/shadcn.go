// Builds the React reference: one bundle that can render any upstream
// new-york-v4 example by name, resolved exactly as the upstream app resolves
// it (ui/apps/v4 with its pnpm-installed, lockfile-pinned dependencies).
//
//   node parity/reference/build.mjs          # parity/reference/dist/app.js
//
// Prerequisite (once): in ./ui,
//   npx pnpm@10.33.4 install --frozen-lockfile --ignore-scripts --filter "v4..."
import fs from "node:fs"
import path from "node:path"
import * as esbuild from "esbuild"

const here = import.meta.dirname
const repo = path.resolve(here, "../..")
const app = path.join(repo, "ui/apps/v4")
const examplesDir = path.join(app, "registry/new-york-v4/examples")
// Builds land in builds/<id> and dist is a symlink swapped atomically, so
// a rebuild never breaks a parity server that is running meanwhile.
const buildsDir = path.join(here, "builds")
const buildId = `${Date.now()}-${process.pid}`
const outdir = path.join(buildsDir, buildId)
const distLink = path.join(here, "dist")

if (!fs.existsSync(path.join(app, "node_modules/react"))) {
  console.error("upstream dependencies missing: run the pnpm install described in build.mjs")
  process.exit(1)
}

// Every upstream reference, with its kind:
//   example   registry/new-york-v4/examples/<name>.tsx
//   internal  registry/new-york-v4/internal/<name>.tsx (sidebar demos)
//   chart     registry/new-york-v4/charts/<name>.tsx
//   block     registry/new-york-v4/blocks/<name>/page.tsx (full pages)
//   fixture   parity/fixtures/<name>.tsx: demos written for this port for
//             components upstream ships without examples; they must use the
//             upstream components unmodified.
const registryDir = path.join(app, "registry/new-york-v4")
const fixturesDir = path.join(repo, "parity/fixtures")
const tsxIn = (dir) => (fs.existsSync(dir) ? fs.readdirSync(dir).filter((f) => f.endsWith(".tsx")).map((f) => f.slice(0, -4)) : [])
const refs = [
  ...tsxIn(examplesDir).map((n) => ({ name: n, kind: "example", file: path.join(examplesDir, n + ".tsx") })),
  ...tsxIn(path.join(registryDir, "internal")).map((n) => ({ name: n, kind: "internal", file: path.join(registryDir, "internal", n + ".tsx") })),
  ...tsxIn(path.join(registryDir, "charts")).map((n) => ({ name: n, kind: "chart", file: path.join(registryDir, "charts", n + ".tsx") })),
  ...fs.readdirSync(path.join(registryDir, "blocks"), { withFileTypes: true }).filter((d) => d.isDirectory() && fs.existsSync(path.join(registryDir, "blocks", d.name, "page.tsx"))).map((d) => ({ name: d.name, kind: "block", file: path.join(registryDir, "blocks", d.name, "page.tsx") })),
  ...tsxIn(fixturesDir).map((n) => ({ name: n, kind: "fixture", file: path.join(fixturesDir, n + ".tsx") })),
].sort((a, b) => a.name.localeCompare(b.name))
const dupes = refs.map((r) => r.name).filter((n, i, a) => a.indexOf(n) !== i)
if (dupes.length) throw new Error("duplicate reference names: " + dupes.join(", "))
const names = refs.map((r) => r.name)
const fileOf = Object.fromEntries(refs.map((r) => [r.name, r.file]))

// Examples that toast need sonner's <Toaster />, which the upstream docs
// mount once in their layout.
const usesThemes = names.filter((n) => /from "next-themes"/.test(fs.readFileSync(fileOf[n], "utf8")))
const usesSonner = names.filter((n) => /from "sonner"/.test(fs.readFileSync(fileOf[n], "utf8")))

const entry = `
import * as React from "react"
import { createRoot } from "react-dom/client"
import { createPortal } from "react-dom"
import { ThemeProvider } from "next-themes"
const examples = {
${names.map((n) => `  ${JSON.stringify(n)}: () => import(${JSON.stringify(fileOf[n])}),`).join("\n")}
}
const sonner = new Set(${JSON.stringify(usesSonner)})
const themed = new Set(${JSON.stringify(usesThemes)})

// Mounts next-themes' provider (and the inline <script> it renders) in a
// detached container and portals the content to target, so the page's DOM
// matches upstream, where the provider wraps the whole layout.
function mountThemed(target, theme, content) {
  createRoot(document.createElement("div")).render(
    <ThemeProvider attribute="class" forcedTheme={theme} enableSystem={false} disableTransitionOnChange>
      {createPortal(content, target)}
    </ThemeProvider>
  )
}

function Ready({ children }) {
  React.useEffect(() => {
    // Two frames: let effects that measure and set state settle.
    requestAnimationFrame(() => requestAnimationFrame(() => { window.__parityReady = true }))
  }, [])
  return children
}

async function main() {
  const name = decodeURIComponent(location.pathname.split("/").filter(Boolean).pop() || "")
  const load = examples[name]
  if (!load) throw new Error("unknown example " + name)
  const mod = await load()
  const Example = mod.default
  const theme = document.documentElement.classList.contains("dark") ? "dark" : "light"
  // Upstream mounts next-themes' provider in its layout, outside the preview,
  // so its inline <script> must not end up inside #preview: examples that
  // need the theme context get the provider in a separate root.
  const example = <Ready><Example /></Ready>
  const preview = document.getElementById("preview")
  if (themed.has(name)) mountThemed(preview, theme, example)
  else createRoot(preview).render(example)
  if (sonner.has(name)) {
    const { Toaster } = await import(${JSON.stringify(path.join(app, "registry/new-york-v4/ui/sonner.tsx"))})
    mountThemed(document.getElementById("parity-extra"), theme, <Toaster />)
  }
}
main().catch((e) => { window.__parityError = String(e && e.stack || e); console.error(e) })
`

const shims = {
  "next/link": `import * as React from "react"
export default React.forwardRef(function Link({ href, prefetch, replace, scroll, shallow, passHref, locale, legacyBehavior, ...props }, ref) {
  return React.createElement("a", { ref, href: typeof href === "string" ? href : href?.pathname ?? "", ...props })
})`,
  // next/image's markup for an unoptimized image with fixed dimensions.
  "next/image": `import * as React from "react"
export default React.forwardRef(function Image({ src, alt, width, height, fill, priority, quality, placeholder, blurDataURL, unoptimized, loader, sizes, style, ...props }, ref) {
  return React.createElement("img", { ref, alt, loading: priority ? undefined : "lazy", width, height, decoding: "async", "data-nimg": fill ? "fill" : "1", style: { color: "transparent", ...style }, src: typeof src === "string" ? src : src?.src, ...props })
})`,
  "next/form": `import * as React from "react"
export default React.forwardRef(function Form({ action, ...props }, ref) { return React.createElement("form", { ref, action: typeof action === "string" ? action : undefined, ...props }) })`,
  // Every Google font an example imports becomes a no-op font loader; the
  // harness renders both sides with the same system fonts.
  "next/font/google": [...new Set(names.flatMap((n) => [...fs.readFileSync(fileOf[n], "utf8").matchAll(/import\s*\{([^}]*)\}\s*from\s*"next\/font\/google"/g)].flatMap((m) => m[1].split(",").map((x) => x.trim()).filter(Boolean))))]
    .map((f) => `export function ${f}() { return { className: "", style: {}, variable: "" } }`).join("\n"),
  "next/navigation": `export const useRouter = () => ({ push() {}, replace() {}, refresh() {}, back() {}, prefetch() {} })
export const usePathname = () => location.pathname
export const useSearchParams = () => new URLSearchParams(location.search)`,
}

const aliasPlugin = {
  name: "upstream-aliases",
  setup(build) {
    build.onResolve({ filter: /^next\/(link|image|form|font\/google|navigation)$/ }, (args) => ({ path: args.path, namespace: "shim" }))
    build.onLoad({ filter: /.*/, namespace: "shim" }, (args) => ({ contents: shims[args.path], loader: "jsx", resolveDir: app }))
    build.onResolve({ filter: /^@shadcn\/react\/message-scroller$/ }, () => ({ path: path.join(repo, "ui/packages/react/src/message-scroller/index.ts") }))
    build.onResolve({ filter: /^shadcn\/schema$/ }, () => ({ path: path.join(repo, "ui/packages/shadcn/src/schema/index.ts") }))
  },
}

fs.mkdirSync(outdir, { recursive: true })
const result = await esbuild.build({
  stdin: { contents: entry, loader: "jsx", resolveDir: app, sourcefile: "parity-entry.jsx" },
  bundle: true,
  splitting: true,
  format: "esm",
  outdir,
  entryNames: "app",
  chunkNames: "chunks/[name]-[hash]",
  jsx: "automatic",
  platform: "browser",
  target: "es2022",
  tsconfig: path.join(here, "tsconfig.json"),
  define: { "process.env.NODE_ENV": '"production"' },
  plugins: [aliasPlugin],
  logLevel: "warning",
  metafile: true,
})
fs.writeFileSync(path.join(outdir, "references.json"), JSON.stringify(refs.map(({ name, kind, file }) => ({ name, kind, file: path.relative(repo, file) })), null, 1))
const tmpLink = path.join(here, `.dist-${buildId}`)
fs.symlinkSync(path.join("builds", buildId), tmpLink)
if (fs.existsSync(distLink) && !fs.lstatSync(distLink).isSymbolicLink()) fs.rmSync(distLink, { recursive: true, force: true })
fs.renameSync(tmpLink, distLink)
// Keep the three newest builds: servers started from older ones may still read them.
for (const old of fs.readdirSync(buildsDir).sort().slice(0, -3)) fs.rmSync(path.join(buildsDir, old), { recursive: true, force: true })
const bytes = Object.values(result.metafile.outputs).reduce((n, o) => n + o.bytes, 0)
const kinds = Object.entries(refs.reduce((m, r) => ((m[r.kind] = (m[r.kind] ?? 0) + 1), m), {})).map(([k, n]) => `${n} ${k}`).join(", ")
console.log(`reference: ${kinds}; ${(bytes / 1e6).toFixed(1)} MB in ${Object.keys(result.metafile.outputs).length} files`)
