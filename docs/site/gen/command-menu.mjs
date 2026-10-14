// command-menu.mjs regenerates ../commandmenu.json, the static data of the
// search command menu (components/command-menu.tsx) that the Go site cannot
// read from Markdown: the styles group (registry/styles.tsx: name, title,
// description, the icon rendered to SVG markup, and the /create?preset=
// URL from shadcn/preset's encodePreset) and the color palettes
// (lib/colors.ts getColors()). Bundled from ./ui/apps/v4 with esbuild and
// run in Node.
//
// Usage: node docs/site/gen/command-menu.mjs
import fs from "node:fs"
import path from "node:path"
import { createRequire } from "node:module"
import { fileURLToPath, pathToFileURL } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
const repo = path.resolve(here, "../../..")
const app = path.join(repo, "ui/apps/v4")
const pnpm = path.join(repo, "ui/node_modules/.pnpm")
const esbuild = createRequire(path.join(pnpm, "esbuild@0.25.11/node_modules/esbuild/package.json"))("esbuild")

const entry = `
import { renderToStaticMarkup } from "react-dom/server"
import { STYLES } from "@/registry/styles"
import { getColors } from "@/lib/colors"
import { encodePreset } from "shadcn/preset"
export default {
  styles: STYLES.map((s) => ({
    name: s.name,
    title: s.title,
    description: s.description,
    icon: renderToStaticMarkup(s.icon),
    href: "/create?preset=" + encodePreset({ style: s.name }),
  })),
  colors: getColors().map((p) => ({
    name: p.name,
    colors: p.colors.map((c) => ({ className: c.className, hex: c.hex, oklch: c.oklch, name: c.name })),
  })),
}
`
// Written under the app so that external packages resolve from its node_modules.
const out = path.join(app, "node_modules/.cache/shadcn-go-command-menu.mjs")
await esbuild.build({
  stdin: { contents: entry, resolveDir: app, loader: "tsx" },
  absWorkingDir: app,
  bundle: true,
  platform: "node",
  format: "esm",
  jsx: "automatic",
  packages: "external",
  alias: { "@": app },
  outfile: out,
  logLevel: "error",
})
const data = (await import(pathToFileURL(out))).default
fs.rmSync(out)
fs.writeFileSync(path.join(here, "../commandmenu.json"), JSON.stringify(data, null, 1) + "\n")
console.log(`${data.styles.length} styles, ${data.colors.length} palettes`)
