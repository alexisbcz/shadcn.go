// Extracts the data of the typeset page (app/(app)/(typeset)) into
// docs/site/data/typeset/typeset.json: the fonts (lib/fonts.ts and their
// lib/font-definitions.ts entries), the content fixtures, the chat question
// and the framework templates. It also vendors, from the reference build,
// the next/font CSS of the preview fonts (app/(view)/preview/fonts.ts) to
// static/typeset-fonts.css (font files under static/fonts/, URLs rewritten
// to /fonts/) and the class list those fonts are applied with.
//
//   node docs/site/gen/typeset.mjs [http://localhost:4000]
import fs from "node:fs"
import path from "node:path"
import { execFileSync } from "node:child_process"

const ref = process.argv[2] ?? "http://localhost:4000"
const root = path.resolve(import.meta.dirname, "../../..")
const app = path.join(root, "ui/apps/v4")
const esbuild = path.join(root, "parity/node_modules/.bin/esbuild")
const entry = path.join(import.meta.dirname, ".typeset.entry.ts")
fs.writeFileSync(entry, `
export { FONTS } from "@/app/(app)/(typeset)/lib/fonts"
export { FONT_DEFINITIONS } from "@/lib/font-definitions"
export { FIXTURES, CONTENT_OPTIONS } from "@/app/(app)/(typeset)/lib/fixtures"
export { CHAT_QUESTION } from "@/app/(app)/(typeset)/lib/fixtures/chat"
export { TEMPLATES } from "@/app/(app)/(create)/lib/templates"
`)
const js = execFileSync(esbuild, [entry, "--bundle", "--format=esm", "--platform=node", "--log-level=error", `--tsconfig=${path.join(app, "tsconfig.json")}`], { encoding: "utf8", cwd: app })
fs.rmSync(entry)
const tmp = path.join(import.meta.dirname, ".typeset.tmp.mjs")
fs.writeFileSync(tmp, js)
const m = await import(tmp)
fs.rmSync(tmp)

const data = {
  fonts: m.FONTS,
  definitions: m.FONT_DEFINITIONS.filter((d) => m.FONTS.some((f) => f.id === d.name)).map((d) => ({
    name: d.name, family: d.family, import: d.import, dependency: d.dependency, subsets: d.subsets, ...(d.weight ? { weight: d.weight } : {}),
  })),
  items: m.CONTENT_OPTIONS.map((o) => ({ ...o, html: m.FIXTURES[o.value] })),
  chatQuestion: m.CHAT_QUESTION,
  templates: m.TEMPLATES.map((t) => ({ value: t.value, title: t.title })),
}

// The preview fonts: the stylesheet of the preview page that is not on the
// docs pages, and the next/font classes of the preview layout's wrapper.
const page = await (await fetch(ref + "/preview/typeset/docs")).text()
const docs = await (await fetch(ref + "/docs")).text()
const sheets = (html) => [...html.matchAll(/<link rel="stylesheet" href="([^"]+)"/g)].map((x) => x[1])
const own = sheets(page).filter((s) => !sheets(docs).includes(s))
let css = ""
for (const s of own) {
  const text = await (await fetch(ref + s)).text()
  if (text.includes("@font-face")) css += text
}
const outFonts = path.join(root, "docs/site/static/fonts")
for (const [, name] of css.matchAll(/url\(\.\.\/media\/([^)]+)\)/g)) {
  const file = path.join(outFonts, name)
  if (!fs.existsSync(file)) fs.writeFileSync(file, Buffer.from(await (await fetch(`${ref}/_next/static/media/${name}`)).arrayBuffer()))
}
css = css.replace(/url\(\.\.\/media\//g, "url(/fonts/")
fs.writeFileSync(path.join(root, "docs/site/static/typeset-fonts.css"), css)
data.fontClasses = page.match(/<div class="(geist_[^"]*)"/)[1]

// The hugeicons of the customizer and toolbar, as HugeiconsIcon renders them
// with strokeWidth={2}: inner markup by name.
const hugeicons = await import(path.join(root, "ui/node_modules/.pnpm/@hugeicons+core-free-icons@1.2.1/node_modules/@hugeicons/core-free-icons/dist/esm/index.js"))
const kebab = (k) => k.replace(/[A-Z]/g, (c) => "-" + c.toLowerCase())
const esc = (s) => String(s).replace(/&/g, "&amp;").replace(/"/g, "&quot;")
data.icons = {}
for (const n of ["ArrowHorizontalIcon", "ParagraphSpacingIcon", "TextSmallcapsIcon", "Menu09Icon", "SquareLock01Icon", "SquareUnlock01Icon", "LinkSquare02Icon"]) {
  data.icons[n] = hugeicons[n].map(([tag, attrs]) => `<${tag}` + Object.entries(attrs).filter(([k]) => k !== "key").map(([k, v]) => ` ${kebab(k)}="${esc(k === "strokeWidth" ? 2 : v)}"`).join("") + `></${tag}>`).join("")
}

fs.writeFileSync(path.join(root, "docs/site/data/typeset/typeset.json"), JSON.stringify(data, null, 1) + "\n")
console.log(`${data.fonts.length} fonts, ${data.items.length} items, ${css.length} bytes of font css`)
