// survey.mjs lists every MDX/Markdown construct the docs content uses, from
// upstream's own parser (remark-parse + remark-mdx + remark-gfm), so the Go
// parser can support exactly that set.
import fs from "node:fs"
import path from "node:path"
import { contentDir, listPages, fromMdx, fromCore, frontmatter } from "./lib.mjs"

const { createProcessor } = await fromMdx("@mdx-js/mdx")
const { remarkGfm } = await fromCore("fumadocs-core/mdx-plugins")
const proc = createProcessor({ remarkPlugins: [remarkGfm] })

const counts = {}
const inc = (k) => (counts[k] = (counts[k] || 0) + 1)
const examples = {}
for (const rel of listPages()) {
  const src = fs.readFileSync(path.join(contentDir, rel), "utf8")
  const fm = await frontmatter(src)
  const tree = proc.parse(fm.content)
  const walk = (n, parentType) => {
    inc("node " + n.type + (n.type === "heading" ? n.depth : "") + (n.type === "list" ? (n.ordered ? " ordered" : " bullet") + (n.spread ? " spread" : "") : ""))
    if (n.type === "mdxJsxFlowElement" || n.type === "mdxJsxTextElement") {
      inc(`jsx ${n.type === "mdxJsxFlowElement" ? "flow" : "text"} <${n.name ?? "fragment"}>`)
      for (const a of n.attributes) {
        if (a.type === "mdxJsxExpressionAttribute") { inc("attr spread"); examples["attr spread"] = a.value; continue }
        const kind = a.value === null || a.value === undefined ? "boolean" : typeof a.value === "string" ? "string" : "expression"
        inc(`attr ${kind}`)
        if (kind === "expression") { const k = "attrexpr " + a.name + "=" + a.value.value; inc(k) }
      }
    }
    if (n.type === "mdxFlowExpression" || n.type === "mdxTextExpression") { inc("expr " + JSON.stringify(n.value)) }
    if (n.type === "mdxjsEsm") inc("esm " + n.value)
    if (n.type === "code") inc(`code lang=${n.lang} metaWords=${(n.meta || "").split(/\s+/).filter(Boolean).map((w) => w.replace(/=.*/, "=").replace(/^\/.*/, "/re/").replace(/^\{.*/, "{range}")).join(",")}`)
    if (n.type === "link" || n.type === "image") { if (n.title) inc(n.type + " with title") }
    if (n.type === "linkReference" || n.type === "definition" || n.type === "imageReference" || n.type === "footnoteReference" || n.type === "html") examples[n.type] = (examples[n.type] || []).concat(rel)
    if (n.type === "table") inc("table align " + JSON.stringify(n.align))
    if (n.type === "listItem" && n.checked != null) inc("task item")
    if (n.type === "list" && n.ordered && n.start !== 1) inc("ordered start " + n.start)
    for (const c of n.children || []) walk(c, n.type)
  }
  walk(tree)
}
for (const k of Object.keys(counts).sort()) console.log(String(counts[k]).padStart(6), k)
console.log(JSON.stringify(examples, null, 1))
