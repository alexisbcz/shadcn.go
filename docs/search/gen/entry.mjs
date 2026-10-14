// The JavaScript half of package search, bundled by build.mjs and run in
// goja: fumadocs-core 16.10.5's remark-structure (with the to-markdown
// extensions remark-mdx and remark-gfm register, as in the docs' MDX
// pipeline) and its advanced Orama search server (createFromSource ->
// initAdvancedSearch), unmodified. Sets globalThis.DocsSearch.
import { unified } from "unified"
import remarkMdx from "remark-mdx"
import remarkGfm from "remark-gfm"
import { VFile } from "vfile"
import { remarkStructure } from "fumadocs-core/mdx-plugins/remark-structure"
import { initAdvancedSearch } from "fumadocs-core/search/server"

const processor = unified().use(remarkMdx).use(remarkGfm).use(remarkStructure).freeze()

let server = null

globalThis.DocsSearch = {
  // structure(mdast JSON) -> structuredData JSON.
  structure(json) {
    const tree = JSON.parse(json)
    const file = new VFile()
    processor.runSync(tree, file)
    return JSON.stringify(file.data.structuredData)
  },
  // index(indexes JSON): [{ id, title, description, url, breadcrumbs, structuredData }].
  index(json) {
    server = initAdvancedSearch({ indexes: JSON.parse(json) })
  },
  // search(query) -> Promise<JSON>: the /api/search response body.
  async search(query) {
    if (!server) return "[]"
    return JSON.stringify(await server.search(query, {}))
  },
}
