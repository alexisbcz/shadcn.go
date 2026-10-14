// render.mjs renders every docs page the way the docs site's React does:
// the compiled MDX with the real ui/apps/v4/mdx-components.tsx (bundled with
// esbuild), React 19's renderToStaticMarkup, and the real icon packages the
// pages import. The components mdx-components.tsx imports from the app
// (Callout, ComponentPreview, the new-york-v4 UI, next/image, ...) are
// replaced by placeholders, written exactly like content.Placeholder in Go:
//
//   <mdx-component data-name="X" data-props='[["name",value],...]'>
//     <template data-prop="icon">...</template>children
//   </mdx-component>
//
// next/link is replaced by an equivalent <a> (href last), which the Go
// default registry also renders.
//
// Output: docs/content/testdata/react.json.gz ({"<page>.mdx": html})
// Usage: node docs/oracle/render.mjs [page.mdx ...]
//        node docs/oracle/render.mjs --corpus
import fs from "node:fs"
import os from "node:os"
import zlib from "node:zlib"
import path from "node:path"
import { createRequire } from "node:module"
import { pathToFileURL } from "node:url"
import { repo, app, contentDir, listPages, fromMdx, mdxOptions, frontmatter } from "./lib.mjs"

const reactDir = fs.realpathSync(path.join(app, "node_modules/react"))
const reactDomDir = fs.realpathSync(path.join(app, "node_modules/react-dom"))
const reactPath = path.join(reactDir, "index.js")
const jsxPath = path.join(reactDir, "jsx-runtime.js")

const esbuild = createRequire(path.join(repo, "ui/node_modules/.pnpm/esbuild@0.25.11/node_modules/esbuild/package.json"))("esbuild")

// Placeholder components (shared with the bundle through globalThis).
const React = (await import(pathToFileURL(reactPath).href)).default
const { renderToStaticMarkup } = await import(pathToFileURL(path.join(reactDomDir, "server.node.js")).href)
globalThis.__placeholder = (name) =>
  function Placeholder(props) {
    const pairs = []
    const templates = []
    for (const [k, v] of Object.entries(props)) {
      if (k === "children" || v === undefined) continue
      // Objects (the page tree folder ComponentsList receives) are data the
      // docs chrome owns; placeholders leave them out.
      if (v !== null && typeof v === "object" && !Array.isArray(v) && !React.isValidElement(v)) continue
      if (React.isValidElement(v)) {
        templates.push(React.createElement("template", { key: "t-" + k, "data-prop": k }, v))
        continue
      }
      pairs.push([k, typeof v === "number" && !Number.isFinite(v) ? null : v])
    }
    return React.createElement("mdx-component", { "data-name": name, "data-props": JSON.stringify(pairs) }, ...templates, props.children)
  }
globalThis.__nextLink = function Link({ href, children, ...rest }) {
  return React.createElement("a", { ...rest, href }, children)
}

const stubs = {
  "@/lib/source": `export const source = { pageTree: { children: [{ $id: "components", type: "folder", name: "Components", children: [] }] } }`,
  "next/image": `export default globalThis.__placeholder("Image")`,
  "next/link": `export default globalThis.__nextLink`,
  "@/components/components-list": `export const ComponentsList = globalThis.__placeholder("ComponentsList")`,
}
const named = {
  "@/components/callout": ["Callout"],
  "@/components/code-block-command": ["CodeBlockCommand"],
  "@/components/code-collapsible-wrapper": ["CodeCollapsibleWrapper"],
  "@/components/code-tabs": ["CodeTabs"],
  "@/components/component-preview": ["ComponentPreview"],
  "@/components/component-source": ["ComponentSource"],
  "@/components/copy-button": ["CopyButton"],
  "@/components/directory-list": ["DirectoryList"],
  "@/registry/new-york-v4/ui/accordion": ["Accordion", "AccordionContent", "AccordionItem", "AccordionTrigger"],
  "@/registry/new-york-v4/ui/alert": ["Alert", "AlertDescription", "AlertTitle"],
  "@/registry/new-york-v4/ui/aspect-ratio": ["AspectRatio"],
  "@/registry/new-york-v4/ui/button": ["Button"],
  "@/registry/new-york-v4/ui/kbd": ["Kbd"],
  "@/registry/new-york-v4/ui/tabs": ["Tabs", "TabsContent", "TabsList", "TabsTrigger"],
}
for (const [mod, names] of Object.entries(named)) {
  stubs[mod] = names.map((n) => `export const ${n} = globalThis.__placeholder(${JSON.stringify(n)})`).join("\n")
}

const outBundle = path.join(os.tmpdir(), "shadcn-go-docs-oracle", "mdx-components.mjs")
await esbuild.build({
  entryPoints: [path.join(app, "mdx-components.tsx")],
  bundle: true,
  format: "esm",
  platform: "node",
  outfile: outBundle,
  jsx: "automatic",
  logLevel: "error",
  plugins: [
    {
      name: "docs-oracle",
      setup(build) {
        build.onResolve({ filter: /.*/ }, (args) => {
          if (stubs[args.path] !== undefined) return { path: args.path, namespace: "stub" }
          if (args.path === "react") return { path: reactPath, external: true }
          if (args.path === "react/jsx-runtime") return { path: jsxPath, external: true }
          if (args.path.startsWith("@/")) {
            const base = path.join(app, args.path.slice(2))
            for (const ext of [".tsx", ".ts", "/index.tsx", "/index.ts"]) if (fs.existsSync(base + ext)) return { path: base + ext }
            throw new Error("unresolved " + args.path)
          }
          if (args.path === "cn") return { path: path.join(fs.realpathSync(path.join(app, "node_modules/cn")), "dist/index.js"), external: true }
          if (args.path === "lucide-react") return { path: path.join(fs.realpathSync(path.join(app, "node_modules/lucide-react")), "dist/esm/lucide-react.js"), external: true }
          return undefined
        })
        build.onLoad({ filter: /.*/, namespace: "stub" }, (args) => ({ contents: stubs[args.path], loader: "js" }))
      },
    },
  ],
})
const { mdxComponents } = await import(pathToFileURL(outBundle).href)

const { createProcessor, run } = await fromMdx("@mdx-js/mdx")
const { VFile } = await fromMdx("vfile")
const { opts, remarkInclude } = await mdxOptions()
const processor = createProcessor({
  development: false,
  ...opts,
  outputFormat: "function-body",
  remarkPlugins: [remarkInclude, ...(opts.remarkPlugins ?? [])],
  format: "mdx",
})
const runtime = await import(pathToFileURL(jsxPath).href)

// --corpus: render every input of docs/content/testdata/corpus.json.gz
// (corpus.mjs) that upstream compiles, into
// docs/content/testdata/corpus-react.json.gz (an array aligned with the
// corpus: the HTML, or null when upstream fails to compile or run it, for
// example for an undefined component or an unresolvable import).
if (process.argv[2] === "--corpus") {
  const corpus = JSON.parse(zlib.gunzipSync(fs.readFileSync(path.join(repo, "docs/content/testdata/corpus.json.gz"))))
  const file = path.join(contentDir, "corpus.mdx")
  const out = []
  let ok = 0
  for (const it of corpus) {
    let html = null
    if (it.hastError === undefined) {
      try {
        const code = String(await processor.process(new VFile({ value: it.src, path: file, data: { frontmatter: { title: "Corpus" } } })))
        const mod = await run(code, { ...runtime, baseUrl: pathToFileURL(file).href })
        html = renderToStaticMarkup(React.createElement(mod.default, { components: mdxComponents }))
        ok++
      } catch {
        html = null
      }
    }
    out.push(html)
  }
  fs.writeFileSync(path.join(repo, "docs/content/testdata/corpus-react.json.gz"), zlib.gzipSync(JSON.stringify(out) + "\n", { level: 9 }))
  console.log(`${corpus.length} corpus inputs, ${ok} rendered`)
  process.exit(0)
}

// Warm-up pass (see build.mjs: shiki's tokenize time limit).
const pages = process.argv.slice(2).length ? process.argv.slice(2) : listPages()
const compile = async (rel) => {
  const file = path.join(contentDir, rel)
  const fm = await frontmatter(fs.readFileSync(file, "utf8"))
  return String(await processor.process(new VFile({ value: fm.content, path: file, data: { frontmatter: fm.data } })))
}
for (const rel of pages) await compile(rel)

const rendered = {}
for (const rel of pages) {
  const file = path.join(contentDir, rel)
  const code = await compile(rel)
  const mod = await run(code, { ...runtime, baseUrl: pathToFileURL(file).href })
  const html = renderToStaticMarkup(React.createElement(mod.default, { components: mdxComponents }))
  rendered[rel] = html
}
if (process.argv.slice(2).length === 0) {
  fs.writeFileSync(path.join(repo, "docs/content/testdata/react.json.gz"), zlib.gzipSync(JSON.stringify(rendered, null, 1) + "\n", { level: 9 }))
}
console.log(`${pages.length} pages rendered`)
