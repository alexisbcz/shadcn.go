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
const outdir = path.join(here, "dist")

if (!fs.existsSync(path.join(app, "node_modules/react"))) {
  console.error("upstream dependencies missing: run the pnpm install described in build.mjs")
  process.exit(1)
}

const names = fs.readdirSync(examplesDir).filter((f) => f.endsWith(".tsx")).map((f) => f.slice(0, -4)).sort()

// Examples that toast need sonner's <Toaster />, which the upstream docs
// mount once in their layout.
const usesSonner = names.filter((n) => /from "sonner"/.test(fs.readFileSync(path.join(examplesDir, n + ".tsx"), "utf8")))

const entry = `
import * as React from "react"
import { createRoot } from "react-dom/client"
import { ThemeProvider } from "next-themes"
const examples = {
${names.map((n) => `  ${JSON.stringify(n)}: () => import(${JSON.stringify(path.join(examplesDir, n + ".tsx"))}),`).join("\n")}
}
const sonner = new Set(${JSON.stringify(usesSonner)})

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
  createRoot(document.getElementById("preview")).render(
    <ThemeProvider attribute="class" forcedTheme={theme} enableSystem={false} disableTransitionOnChange>
      <Ready><Example /></Ready>
    </ThemeProvider>
  )
  if (sonner.has(name)) {
    const { Toaster } = await import(${JSON.stringify(path.join(app, "registry/new-york-v4/ui/sonner.tsx"))})
    createRoot(document.getElementById("parity-extra")).render(
      <ThemeProvider attribute="class" forcedTheme={theme} enableSystem={false}><Toaster /></ThemeProvider>
    )
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
  "next/font/google": [...new Set(names.flatMap((n) => [...fs.readFileSync(path.join(examplesDir, n + ".tsx"), "utf8").matchAll(/import\s*\{([^}]*)\}\s*from\s*"next\/font\/google"/g)].flatMap((m) => m[1].split(",").map((x) => x.trim()).filter(Boolean))))]
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

fs.rmSync(outdir, { recursive: true, force: true })
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
fs.writeFileSync(path.join(outdir, "examples.json"), JSON.stringify(names))
const bytes = Object.values(result.metafile.outputs).reduce((n, o) => n + o.bytes, 0)
console.log(`reference: ${names.length} examples, ${(bytes / 1e6).toFixed(1)} MB in ${Object.keys(result.metafile.outputs).length} files`)
