// corpus.mjs runs upstream's actual pipeline on inputs the 316 docs pages do
// not contain, to check that the Go parser and pipeline either reject an
// input or produce exactly upstream's result (never something different):
//
//   - every example of the CommonMark spec (spec.json from the goldmark
//     module, CommonMark spec 0.31.2, CC-BY-SA 4.0), read from
//     $COMMONMARK_SPEC or the Go module cache;
//   - every distinct block of the docs (page bodies split on blank lines),
//     and pairs of consecutive blocks glued without the blank line;
//   - seeded random mutations of docs blocks (insert, delete or duplicate a
//     Markdown/MDX-significant character, indent or join lines).
//
// For each input it records upstream's mdast (remark-parse + remark-mdx +
// remark-gfm, canonical form of mdast.mjs) and the final HTML tree after
// all of fumadocs-mdx's remark and rehype plugins (as build.mjs does, with
// frontmatter {title: "Corpus"}), or the error upstream throws.
//
// Output: docs/content/testdata/corpus.json.gz
// Usage: node docs/oracle/corpus.mjs
// Verification tooling only; nothing here runs in the Go packages.
import fs from "node:fs"
import path from "node:path"
import zlib from "node:zlib"
import { execFileSync } from "node:child_process"
import { repo, contentDir, listPages, fromMdx, fromCore, fromPretty, mdxOptions, frontmatter, toHtml } from "./lib.mjs"

// ---- inputs ----

function specPath() {
  if (process.env.COMMONMARK_SPEC) return process.env.COMMONMARK_SPEC
  const modcache = execFileSync("go", ["env", "GOMODCACHE"], { encoding: "utf8" }).trim()
  const dir = path.join(modcache, "github.com/yuin")
  const v = fs.readdirSync(dir).filter((d) => d.startsWith("goldmark@")).sort().pop()
  if (!v) throw new Error("goldmark not in the module cache; set COMMONMARK_SPEC")
  return path.join(dir, v, "_test/spec.json")
}

// mulberry32: a small seeded PRNG so the corpus is reproducible.
function rng(seed) {
  return () => {
    seed |= 0
    seed = (seed + 0x6d2b79f5) | 0
    let t = Math.imul(seed ^ (seed >>> 15), 1 | seed)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

const inputs = []
const seen = new Set()
const add = (kind, src) => {
  if (seen.has(src)) return
  seen.add(src)
  inputs.push({ kind, src })
}

for (const ex of JSON.parse(fs.readFileSync(specPath(), "utf8"))) add("spec", ex.markdown)

// Hand-written edge cases around the constructs the docs use: containers
// and laziness, list spread, tables, code spans and expressions over
// several lines, JSX, ESM.
const curated = [
  "- a\n2. b\n", "- a\n- b\n", "- a\n\n- b\n", "- a\n\n  b\n", "- a\n  - b\n\n  c\n", "- a\n  - b\n\n\n  c\n",
  "- a\n\n\n- b\n", "1. a\n1. b\n", "1. a\n\n   b\n2. c\n", "- a\n-\n- b\n", "-\n  a\n", "- \n  a\n",
  "- a\n> b\n", "> a\n- b\n", "> - a\n>\n>   b\n", "> - a\n>\n> b\n", "- > a\n  >\n  b\n", "- > a\n\n  b\n",
  "- a\n  ```\n  x\n\n  y\n  ```\n", "- a\n\n  ```\n  x\n  ```\n- b\n", "1.  a\n\n    b\n", "1.     a\n", "-    a\n     b\n",
  "- a\nb\n", "> a\nb\n", "> - a\nb\n", "- a\n# b\n", "- a\n<X />\n", "- a\n<X\n/>\n", "- a\n{1}\n", "- a\n{1 +\n2}\n",
  "- a <X\n  b /> c\n", "- a <X\nb /> c\n", "- a {1 +\n  2} c\n", "- a {1 +\n2} c\n", "> a {1 +\n> 2} c\n",
  "a `b\n     c` d\n", "a `b\n\tc` d\n", "a `\nb\n` c\n", "a ` \nb \n` c\n", "a `  \n  ` c\n", "a ``b\n``\n",
  "a {1 +\n     2} d\n", "a {1 +\n  2} d\n", "a {1 +\n 2} d\n", "a <B c={1 +\n      2} /> d\n", "a <B\n   c=\"x\" /> d\n",
  "a\n  b\n", "a\n\tb\n", "a  \n  b\n", "*a\n  b*\n", "a\n  *b*\n", "[a\n  b](/c)\n", "[a](\n  /c)\n",
  "| a |\n| - |\n| b |\n2. x\n", "| a |\n| - |\n| b |\n-\n", "| a |\n| - |\n| b |\n> q\n", "| a |\n| - |\n| b |\n<X />\n",
  "| a |\n| - |\n| b |\n{1}\n", "| a |\n| - |\n| b |\n[^1]: x\n", "| a |\n| - |\n| b |\n***\n", "| a |\n| - |\n| b |\n=\n",
  "| a |\n| - |\n| b |\n[a]: /b\n", "| a |\n| - |\n| b |\n# h\n", "| a |\n| - |\n| b |\n```\nx\n```\n", "| a |\n| - |\n| b |\nc\n",
  "| a |\n| - |\n| b |\n    c\n", "text\n| a |\n| - |\n", "- a\n| b |\n| - |\n", "- a\n  | b |\n  | - |\n  | c |\nd\n",
  "> | a |\n> | - |\n| b |\n", "| a | b |\n| - |\n", "| a |\n| - | - |\n", "a | b\n- | -\nc | d\n", "| a |\n|:-:|\n| `x \\| y` |\n",
  "| a |\n| - |\n| b | c | d |\n", "| a | b |\n| - | - |\n| c |\n", "| a |\n| - |\n|\n", "| a |\n| - |\n| \\\\ |\n",
  "[a]: /b\n", "[a\nb]: /c\n", "[a\\]]: /b\n", " [a]: /b\n", "x\n[a]: /b\n", "[^a]: b\n\nc[^a]\n", "[^a]:\n  b\n\nc[^a]\n",
  "import { A } from \"b\"\n", "import { A, B as C } from \"b\"\n", "import { type A } from \"b\"\n", "import A from \"b\"\n",
  "import { A } from \"b\";\n", "import {A} from 'b'\n", "export const a = 1\n", "import { A } from \"b\"\nimport { C } from \"d\"\n",
  "<X>\n- a\n</X>\n", "<X>\n  - a\n</X>\n", "<X>a</X>\n", "<X>\n\na\n\n</X>\n", "<X\n  a=\"b\"\n/>\n", "<X a={[1,\n  2]} />\n",
  "<X>\n```\nx\n```\n</X>\n", "<X>\n  ```\n  x\n  ```\n</X>\n", "<X>\n  ```\n    x\n  ```\n</X>\n",
  "```\n  a\n ```\n", "  ```\n  a\n b\n   c\n  ```\n", "```js title=\"a b\"\nx\n```\n", "~~~\na\n~~~\n", "````\n```\n````\n",
  "# a #\n", "# a \\#\n", "#\n", "## a ##  \n", "### `a` b\n", "# a {#b}\n", "# a [!toc]\n", "# a [toc]\n",
]
// Code blocks: every rehype-pretty-code meta feature, the docs
// transformer's command forms, language aliases and unusual text.
const fence = (info, code) => "```" + info + "\n" + code + "\n```\n"
for (const [info, code] of [
  ['tsx title="a.tsx" caption="The caption"', "const a = 1"], ['tsx caption="c"', "x"], ['tsx title=""', "x"], ['tsx title="a" title="b"', "x"],
  ["tsx showLineNumbers{5}", "a\nb"], ["tsx showLineNumbers{0}", "a"], ["tsx showLineNumbers{99}", "a\nb\nc"], ["tsx showLineNumbers", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10"],
  ["tsx {1,3-4}", "a\nb\nc\nd"], ["tsx {1}#v {2}#w", "a\nb"], ["tsx {2-1}", "a\nb"], ["tsx {0}", "a"], ["tsx {9}", "a"], ["tsx {1, 2}", "a\nb"],
  ["tsx /a/", "a a a"], ["tsx /a/1-2#x", "a a a\na"], ["tsx /a/#x /b/", "a b"], ["tsx /a/2", "a a"], ["tsx /a b/", "a b a b"], ["tsx /=/", "const a = b"],
  ["tsx /const/ {2}", "const a = 1\nconst b = 2"], ["tsx /a/ /a/", "a"], ["tsx /<div>/", "<div>x</div>"], ["tsx /\\//", "a/b"],
  ["js {2}", "const a = 1\n\nconst b = 2"], ["tsx", ""], ["", "no lang"], ["math", "x^2"], ["text", "plain <b>&amp;</b>"], ["txt showLineNumbers", "a\nb"],
  ["plaintext", "x"], ["plain", "x"], ["ansi", "\u001b[31mred\u001b[0m"], ["unknownlang", "x y"], ["TSX", "const a = 1"], ["diff", "+ a\n- b\n  c"],
  ["bash", "npm install x"], ["bash", "npx create-next-app@latest"], ["bash", "npm create vite@latest"], ["bash", "npm run dev"], ["bash", "npx shadcn@latest add button"],
  ["sh", "echo $HOME"], ["zsh", "ls -la"], ["shell", "cd x"], ["js", "let a = `t ${b}`"], ["mjs", "export default 1"], ["cjs", "module.exports = 1"],
  ["ts", "type A<T> = { [K in keyof T]?: T[K] }"], ["mts", "export {}"], ["cts", "export = 1"], ["jsx", "<A b={1} />"], ["md", "# h\n\n- a"],
  ["markdown", "```js\nconst a = 1\n```"], ["html", "<!doctype html><p class=\"a\">x &amp; y</p><script>let a = 1</script><style>p{color:red}</style>"],
  ["css", "@layer base { :root { --a: oklch(1 0 0); } }"], ["postcss", "@apply flex;"], ["json", "{\"a\": [1, true, null]}"], ["toml", "[a]\nb = \"c\""],
  ["regexp", "^(?<a>\\d+)$"], ["regex", "[a-z]+"], ["astro", "---\nconst a = 1\n---\n<p>{a}</p>"],
  ["tsx", "const s = \"\\u00e9 \u00e9 \u{1F600} \u200b\"\n\tindented\n    spaces   \ntrailing   "], ["tsx", "a <b> & \"c\" 'd'"],
  ["tsx", "x".repeat(3000)], ["tsx", "// " + "é".repeat(500)], ["json", "{\n" + Array.from({ length: 200 }, (_, i) => `  "k${i}": ${i}`).join(",\n") + "\n}"],
]) add("curated", fence(info, code))

for (const s of curated) add("curated", s)

const blocks = []
for (const rel of listPages()) {
  const fm = await frontmatter(fs.readFileSync(path.join(contentDir, rel), "utf8"))
  const parts = fm.content.split(/\n{2,}/).filter((s) => s.trim() !== "")
  for (let i = 0; i < parts.length; i++) {
    blocks.push(parts[i] + "\n")
    if (i + 1 < parts.length) blocks.push(parts[i] + "\n" + parts[i + 1] + "\n")
  }
}
for (const b of blocks) add("docs", b)

const rand = rng(20261014)
const chars = ["*", "_", "`", "[", "]", "(", ")", "!", "<", ">", "{", "}", "#", "|", "-", "\\", "~", "&", ":", "/", " ", "\n", "1.", "  ", "\t", "\"", "'", "="]
const pick = (a) => a[Math.floor(rand() * a.length)]
function mutate(s) {
  const op = Math.floor(rand() * 6)
  const i = Math.floor(rand() * (s.length + 1))
  switch (op) {
    case 0:
      return s.slice(0, i) + pick(chars) + s.slice(i)
    case 1:
      return s.slice(0, i) + s.slice(i + 1)
    case 2:
      return s.slice(0, i) + s[i] + s.slice(i)
    case 3: {
      const lines = s.split("\n")
      const j = Math.floor(rand() * lines.length)
      lines[j] = " ".repeat(1 + Math.floor(rand() * 4)) + lines[j]
      return lines.join("\n")
    }
    case 4: {
      const nl = [...s.matchAll(/\n/g)].map((m) => m.index)
      if (nl.length < 2) return s + pick(chars)
      const k = nl[Math.floor(rand() * (nl.length - 1))]
      return s.slice(0, k) + s.slice(k + 1)
    }
    default:
      return s.slice(0, i) + pick(chars) + pick(chars) + s.slice(i)
  }
}
const small = [...new Set(blocks)].filter((b) => b.length < 600)
for (let n = 0; n < 4000; n++) {
  let s = pick(small)
  const times = 1 + Math.floor(rand() * 3)
  for (let k = 0; k < times; k++) s = mutate(s)
  add("mutation", s)
}

// ---- upstream ----

const { createProcessor } = await fromMdx("@mdx-js/mdx")
const { VFile } = await fromMdx("vfile")
const { remarkGfm } = await fromCore("fumadocs-core/mdx-plugins")
const parser = createProcessor({ remarkPlugins: [remarkGfm] })
const { opts, remarkInclude } = await mdxOptions()

const shiki = await fromPretty("shiki")
const hl = await shiki.getSingletonHighlighter({ themes: ["vesper", "github-light-default"], langs: ["plaintext"] })
for (const l of ["astro", "css", "diff", "html", "javascript", "json", "jsx", "markdown", "postcss", "regexp", "shellscript", "toml", "tsx", "typescript"]) await hl.loadLanguage(l)

function canon(n) {
  const o = { type: n.type }
  for (const k of ["value", "depth", "ordered", "start", "spread", "checked", "lang", "meta", "url", "title", "alt", "align", "identifier", "label", "name"]) {
    if (n[k] !== undefined && n[k] !== null) o[k] = n[k]
  }
  if (n.attributes) {
    o.attributes = n.attributes.map((a) => (a.type === "mdxJsxExpressionAttribute" ? { spread: a.value } : a.value === null || a.value === undefined ? { name: a.name } : typeof a.value === "string" ? { name: a.name, value: a.value } : { name: a.name, expr: a.value.value }))
  }
  if (n.children) o.children = n.children.map(canon)
  return o
}

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
  if (node.type === "mdxjsEsm") return { type: "element", tagName: "mdx-esm", properties: { dataValue: node.value }, children: [] }
  if (node.children) return { ...node, children: node.children.map(mdxToElements) }
  return node
}
let captured
function captureFinal() {
  return (tree) => {
    const kids = []
    for (let i = 0; i < tree.children.length; i++) {
      const c = tree.children[i]
      if (c.type === "mdxjsEsm" && c.value === "") {
        const next = tree.children[i + 1]
        if (next && next.type === "text" && next.value === "\n") i++
        continue
      }
      kids.push(c)
    }
    captured.hast = toHtml(mdxToElements({ ...tree, children: kids }))
  }
}
const { rehypeToc } = await fromCore("fumadocs-core/mdx-plugins")
function captureToc() {
  const run = rehypeToc({ exportToc: { as: "data" } })
  return (tree) => {
    const copy = structuredClone(tree)
    const f = { data: {} }
    run(copy, f)
    captured.toc = f.data.rehypeToc.map((t) => ({ depth: t.depth, url: t.url, hast: toHtml({ type: "root", children: mdxToElements(t.title).children }) }))
  }
}
const processor = createProcessor({
  development: false,
  ...opts,
  outputFormat: "function-body",
  rehypePlugins: [captureToc, ...opts.rehypePlugins, captureFinal],
  remarkPlugins: [remarkInclude, ...(opts.remarkPlugins ?? [])],
  format: "mdx",
})

const msg = (e) => String(e?.reason ?? e?.message ?? e).split("\n")[0]
const out = []
const file = path.join(contentDir, "corpus.mdx")
for (const it of inputs) {
  const r = { kind: it.kind, src: it.src }
  try {
    r.mdast = canon(parser.parse(it.src))
  } catch (e) {
    r.mdastError = msg(e)
  }
  // The full pipeline, twice: highlighting results must not depend on
  // shiki's warm-up (tokenizeTimeLimit, see build.mjs).
  for (let pass = 0; pass < 2; pass++) {
    captured = {}
    try {
      await processor.process(new VFile({ value: it.src, path: file, data: { frontmatter: { title: "Corpus" } } }))
      r.hast = captured.hast
      r.toc = captured.toc
      delete r.hastError
    } catch (e) {
      r.hastError = msg(e)
      delete r.hast
      delete r.toc
    }
  }
  out.push(r)
}
const dest = path.join(repo, "docs/content/testdata/corpus.json.gz")
fs.writeFileSync(dest, zlib.gzipSync(JSON.stringify(out) + "\n", { level: 9 }))
const count = (k) => out.filter((r) => r.kind === k).length
console.log(`${out.length} inputs (spec ${count("spec")}, curated ${count("curated")}, docs ${count("docs")}, mutation ${count("mutation")}); upstream mdast errors ${out.filter((r) => r.mdastError).length}, pipeline errors ${out.filter((r) => r.hastError).length}`)
