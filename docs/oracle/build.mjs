// build.mjs runs every docs page through upstream's actual MDX pipeline
// (fumadocs-mdx 15.0.12 with ui/apps/v4/source.config.ts: remark-gfm,
// fumadocs' remark plugins, rehype-toc and rehype-pretty-code with the docs'
// shiki transformer) and writes the expected outputs the Go packages are
// tested against:
//
//   docs/content/testdata/oracle/<page>.json   frontmatter, TOC, hast (HTML)
//   docs/highlight/testdata/blocks.json.gz     every code block: input and
//                                             rehype-pretty-code's markup
//
// Usage: node docs/oracle/build.mjs [page.mdx ...]
// Verification tooling only; nothing here runs in the Go packages.
import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"
import { repo, contentDir, listPages, fromMdx, fromPretty, mdxOptions, frontmatter, toHtml } from "./lib.mjs"

const { createProcessor } = await fromMdx("@mdx-js/mdx")
const { VFile } = await fromMdx("vfile")
const { opts, remarkInclude } = await mdxOptions()

// Preload every grammar the docs use into the highlighter rehype-pretty-code
// shares (shiki's singleton), so a block's output never depends on which
// pages were processed before it. The Go highlighter loads the same set.
const shiki = await fromPretty("shiki")
const hl = await shiki.getSingletonHighlighter({ themes: ["vesper", "github-light-default"], langs: ["plaintext"] })
const allLangs = ["astro", "css", "diff", "html", "javascript", "json", "jsx", "markdown", "postcss", "regexp", "shellscript", "toml", "tsx", "typescript"]
for (const l of allLangs) await hl.loadLanguage(l)

// mdxToElements converts MDX nodes into placeholder elements so the tree
// can be serialized by hast-util-to-html (the Go side does the same).
function attrsJSON(attrs) {
  return JSON.stringify(
    attrs.map((a) => {
      if (a.type !== "mdxJsxAttribute") throw new Error("unsupported JSX spread attribute")
      if (a.value === null || a.value === undefined) return { name: a.name }
      if (typeof a.value === "string") return { name: a.name, value: a.value }
      return { name: a.name, expr: a.value.value }
    })
  )
}
function mdxToElements(node) {
  if (node.type === "mdxJsxFlowElement" || node.type === "mdxJsxTextElement") {
    return {
      type: "element",
      tagName: "mdx-jsx",
      properties: { dataName: node.name ?? "", dataFlow: String(node.type === "mdxJsxFlowElement"), dataAttrs: attrsJSON(node.attributes) },
      children: node.children.map(mdxToElements),
    }
  }
  if (node.type === "mdxFlowExpression" || node.type === "mdxTextExpression") {
    return { type: "element", tagName: "mdx-expr", properties: { dataFlow: String(node.type === "mdxFlowExpression"), dataValue: node.value }, children: [] }
  }
  if (node.type === "mdxjsEsm") {
    return { type: "element", tagName: "mdx-esm", properties: { dataValue: node.value }, children: [] }
  }
  if (node.children) return { ...node, children: node.children.map(mdxToElements) }
  return node
}

const blocks = new Map()
let captured

function capturePre() {
  return (tree) => {
    const pres = []
    const visit = (n) => {
      if (n.type === "element" && n.tagName === "pre" && n.children.length === 1 && n.children[0].tagName === "code") {
        const code = n.children[0]
        const cls = code.properties.className?.[0]
        pres.push({
          el: n,
          lang: typeof cls === "string" && cls.startsWith("language-") ? cls.slice(9) : "",
          meta: code.data?.meta ?? "",
          // mdast-util-to-hast appended a newline to non-empty code.
          code: (code.children[0]?.value ?? "").replace(/\n$/, ""),
        })
      }
      for (const c of n.children || []) visit(c)
    }
    visit(tree)
    captured = { pres }
  }
}
function captureFinal() {
  return (tree) => {
    // fumadocs-mdx exports (frontmatter, toc, structuredData) are ESM nodes
    // with an empty value; they carry no markup.
    const root = mdxToElements({ ...tree, children: tree.children.filter((c) => !(c.type === "mdxjsEsm" && c.value === "")) })
    captured.hast = toHtml(root)
    for (const p of captured.pres) {
      const key = JSON.stringify([p.lang, p.meta, p.code])
      if (!blocks.has(key)) blocks.set(key, { lang: p.lang, meta: p.meta, code: p.code, html: toHtml(p.el) })
    }
  }
}

const rehypePlugins = [...opts.rehypePlugins]
// opts.rehypePlugins is [rehypeToc, [rehypePrettyCode, ...]]: capture the
// <pre> elements just before rehype-pretty-code rewrites them in place.
const rpcIndex = rehypePlugins.findIndex((p) => Array.isArray(p) && p[0].name === "rehypePrettyCode")
if (rpcIndex < 0) throw new Error("rehype-pretty-code not found in the pipeline")
rehypePlugins.splice(rpcIndex, 0, capturePre)
rehypePlugins.push(captureFinal)

const processor = createProcessor({
  development: false,
  ...opts,
  outputFormat: "function-body",
  rehypePlugins,
  remarkPlugins: [remarkInclude, ...(opts.remarkPlugins ?? [])],
  format: "mdx",
})

const { run } = await fromMdx("@mdx-js/mdx")
const runtime = await import(path.join(repo, "ui/apps/v4/node_modules/react/jsx-runtime.js"))
const React = (await import(path.join(repo, "ui/apps/v4/node_modules/react/index.js"))).default
const { renderToStaticMarkup } = await import(path.join(repo, "ui/apps/v4/node_modules/react-dom/server.node.js"))

const outDir = path.join(repo, "docs/content/testdata/oracle")
const pages = process.argv.slice(2).length ? process.argv.slice(2) : listPages()
for (const rel of pages) {
  const file = path.join(contentDir, rel)
  const source = fs.readFileSync(file, "utf8")
  const fm = await frontmatter(source)
  const vfile = new VFile({ value: fm.content, path: file, data: { frontmatter: fm.data } })
  const compiled = await processor.process(vfile)
  // The compiled module exports the TOC fumadocs computes (rehype-toc),
  // with JSX titles; render each title to HTML.
  const mod = await run(String(compiled), { ...runtime, baseUrl: "file://" + file })
  const toc = (mod.toc ?? []).map((t) => ({ depth: t.depth, url: t.url, title: renderToStaticMarkup(React.createElement(React.Fragment, null, t.title)) }))
  const out = {
    path: rel,
    frontmatter: JSON.parse(JSON.stringify(fm.data)),
    toc,
    hast: captured.hast,
  }
  const dest = path.join(outDir, rel.replace(/\.mdx$/, ".json"))
  fs.mkdirSync(path.dirname(dest), { recursive: true })
  fs.writeFileSync(dest, JSON.stringify(out, null, 1) + "\n")
}

if (process.argv.slice(2).length === 0) {
  const list = [...blocks.values()].sort((a, b) => (a.lang + a.meta + a.code < b.lang + b.meta + b.code ? -1 : 1))
  const dest = path.join(repo, "docs/highlight/testdata/blocks.json.gz")
  fs.mkdirSync(path.dirname(dest), { recursive: true })
  fs.writeFileSync(dest, zlib.gzipSync(JSON.stringify(list, null, 1) + "\n", { level: 9 }))
  console.log(`${pages.length} pages, ${list.length} distinct code blocks`)
}
