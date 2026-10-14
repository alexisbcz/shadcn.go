// mdast.mjs dumps upstream's mdast (remark-parse + remark-mdx + remark-gfm,
// before transforms) for every page as canonical JSON, to debug the Go
// parser against: node docs/oracle/mdast.mjs <outdir>
import fs from "node:fs"
import path from "node:path"
import { contentDir, listPages, fromMdx, fromCore, frontmatter } from "./lib.mjs"

const { createProcessor } = await fromMdx("@mdx-js/mdx")
const { remarkGfm } = await fromCore("fumadocs-core/mdx-plugins")
const proc = createProcessor({ remarkPlugins: [remarkGfm] })
const out = process.argv[2]

function canon(n) {
  const o = { type: n.type }
  for (const k of ["value", "depth", "ordered", "start", "spread", "checked", "lang", "meta", "url", "title", "alt", "align", "identifier", "label", "name"]) {
    if (n[k] !== undefined && n[k] !== null) o[k] = n[k]
  }
  if (n.attributes) {
    o.attributes = n.attributes.map((a) => (a.value === null || a.value === undefined ? { name: a.name } : typeof a.value === "string" ? { name: a.name, value: a.value } : { name: a.name, expr: a.value.value }))
  }
  if (n.children) o.children = n.children.map(canon)
  return o
}
for (const rel of listPages()) {
  const src = fs.readFileSync(path.join(contentDir, rel), "utf8")
  const fm = await frontmatter(src)
  const tree = proc.parse(fm.content)
  const dest = path.join(out, rel.replace(/\.mdx$/, ".json"))
  fs.mkdirSync(path.dirname(dest), { recursive: true })
  fs.writeFileSync(dest, JSON.stringify(canon(tree), null, 1) + "\n")
}
