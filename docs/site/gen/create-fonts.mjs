// Vendors the preview fonts of the reference build (app/(view)/preview/fonts.ts:
// next/font/google, one variable class per font) for /preview pages:
// static/create-fonts.css with the font files under static/fonts.
//
// Usage: node docs/site/gen/create-fonts.mjs [http://localhost:4000]
import fs from "node:fs"
import path from "node:path"

const ref = process.argv[2] ?? "http://localhost:4000"
const root = path.resolve(import.meta.dirname, "../../..")
const html = await (await fetch(`${ref}/preview/base/preview-02`)).text()
const sheets = [...html.matchAll(/href="(\/_next\/static\/chunks\/[^"]+\.css)"/g)].map((m) => m[1])
let css
for (const s of sheets) {
  const text = await (await fetch(ref + s)).text()
  if (text.includes("inter_") && text.includes("-module__") && text.startsWith("@font-face")) css = text
}
if (!css) throw new Error("preview font stylesheet not found")
const fonts = path.join(root, "docs/site/static/fonts")
for (const [, file] of css.matchAll(/url\(\.\.\/media\/([^)]+)\)/g)) {
  const dst = path.join(fonts, file)
  if (fs.existsSync(dst)) continue
  const res = await fetch(`${ref}/_next/static/media/${file}`)
  if (!res.ok) throw new Error(`${file}: ${res.status}`)
  fs.writeFileSync(dst, Buffer.from(await res.arrayBuffer()))
}
css = css.replace(/url\(\.\.\/media\//g, "url(/fonts/")
fs.writeFileSync(path.join(root, "docs/site/static/create-fonts.css"), css)
const classes = [...new Set([...css.matchAll(/\.([a-z_0-9]+-module__[A-Za-z0-9_-]+__variable)\{/g)].map((m) => m[1]))]
fs.writeFileSync(path.join(root, "docs/site/data/create/font-classes.json"), JSON.stringify(classes) + "\n")
console.log(`${classes.length} font classes`)
