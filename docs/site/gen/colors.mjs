// Extracts the Tailwind palette of the colors page (registry/_legacy-colors.ts)
// into docs/site/data/colors.json, in source order.
//
//   node docs/site/gen/colors.mjs
import fs from "node:fs"
import path from "node:path"
import { execFileSync } from "node:child_process"

const root = path.resolve(import.meta.dirname, "../../..")
const src = path.join(root, "ui/apps/v4/registry/_legacy-colors.ts")
const esbuild = path.join(root, "parity/node_modules/.bin/esbuild")
const js = execFileSync(esbuild, [src, "--format=esm", "--log-level=error"], { encoding: "utf8" })
const tmp = path.join(import.meta.dirname, ".colors.tmp.mjs")
fs.writeFileSync(tmp, js)
const { colors } = await import(tmp)
fs.rmSync(tmp)
const palettes = Object.entries(colors).filter(([, v]) => Array.isArray(v)).map(([name, list]) => ({ name, colors: list }))
fs.writeFileSync(path.join(root, "docs/site/data/colors.json"), JSON.stringify(palettes, null, 1) + "\n")
console.log(`${palettes.length} palettes`)
