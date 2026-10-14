// Extracts the examples theme selector's data (lib/themes.ts THEMES and
// registry/_legacy-base-colors.ts baseColorsOKLCH) into
// docs/site/data/examples/theme/themes.json, in source order.
//
//   node docs/site/gen/themes.mjs
import fs from "node:fs"
import path from "node:path"
import { execFileSync } from "node:child_process"

const root = path.resolve(import.meta.dirname, "../../..")
const src = path.join(root, "ui/apps/v4/registry/_legacy-base-colors.ts")
const esbuild = path.join(root, "parity/node_modules/.bin/esbuild")
const js = execFileSync(esbuild, [src, "--format=esm", "--log-level=error"], { encoding: "utf8" })
const tmp = path.join(import.meta.dirname, ".themes.tmp.mjs")
fs.writeFileSync(tmp, js)
const { baseColors, baseColorsOKLCH } = await import(tmp)
fs.rmSync(tmp)
// lib/themes.ts
const themes = baseColors
  .filter((theme) => !["slate", "stone", "gray", "zinc"].includes(theme.name))
  .sort((a, b) => a.name.localeCompare(b.name))
  .map((t) => ({
    name: t.name,
    label: t.label,
    hsl: { light: Object.entries(t.cssVars.light), dark: Object.entries(t.cssVars.dark) },
    oklch: baseColorsOKLCH[t.name] && {
      light: Object.entries(baseColorsOKLCH[t.name].light),
      dark: Object.entries(baseColorsOKLCH[t.name].dark),
    },
  }))
fs.writeFileSync(path.join(root, "docs/site/data/examples/theme/themes.json"), JSON.stringify(themes, null, 1) + "\n")
console.log(themes.map((t) => t.name + (t.oklch ? "" : "(no oklch)")).join(" "))
