// Lists every page path of the upstream documentation site (ui/apps/v4),
// the input of the page-level parity harness (parity/docs/run.mjs).
//
//   node parity/docs/pages.mjs                     # print every path, one per line
//   node parity/docs/pages.mjs --group docs        # only one group (docs, app, view)
//   node parity/docs/pages.mjs --check             # cross-check against the build and the live server
//   node parity/docs/pages.mjs --json              # [{ path, group, source }]
//
// The list is computed from the sources, not from the build:
//
//   docs  every .mdx file under content/docs, routed the way fumadocs-core's
//         loader({ baseUrl: "/docs" }) routes them: the path relative to
//         content/docs without the extension, folder names in parentheses
//         ("(root)") are route groups and are dropped, and index.mdx stands
//         for its folder. Plus app/(app)/docs/changelog/page.tsx.
//   app   every app/**/page.tsx without a dynamic segment, route groups
//         dropped; the dynamic app routes expanded with their own
//         generateStaticParams data: /blocks/[...categories] from
//         lib/categories.ts, /charts/[type] from the chartTypes array of its
//         page. Plus one unknown path, to compare the 404 page.
//   view  the iframe views (/view/[style]/[name], /preview/[base]/[name],
//         /preview/typeset/[name]): their generateStaticParams read the
//         generated registry indexes, so they are taken from the build output
//         (.next/prerender-manifest.json) instead of recomputed.
//
// --check verifies that (1) every computed docs/app path is prerendered by
// the build (.next/prerender-manifest.json) and every prerendered HTML route
// is in the list, (2) the list contains every <loc> of the reference's
// /sitemap.xml, and (3) the reference serves every path with HTTP 200 (the
// 404 probe with 404). It exits non-zero on any mismatch.
import fs from "node:fs"
import path from "node:path"

const repo = path.resolve(import.meta.dirname, "../..")
const app = path.join(repo, "ui/apps/v4")

export const NOT_FOUND_PROBE = "/parity-not-found"

const isGroup = (seg) => /^\(.*\)$/.test(seg)

function walk(dir, ext) {
  const out = []
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, e.name)
    if (e.isDirectory()) out.push(...walk(p, ext))
    else if (e.name.endsWith(ext)) out.push(p)
  }
  return out.sort()
}

export function docsPaths() {
  const root = path.join(app, "content/docs")
  const out = []
  for (const f of walk(root, ".mdx")) {
    const segs = path.relative(root, f).replace(/\.mdx$/, "").split(path.sep).filter((s) => !isGroup(s))
    if (segs.at(-1) === "index") segs.pop()
    out.push({ path: "/docs" + (segs.length ? "/" + segs.join("/") : ""), group: "docs", source: path.relative(repo, f) })
  }
  out.push({ path: "/docs/changelog", group: "docs", source: "ui/apps/v4/app/(app)/docs/changelog/page.tsx" })
  return out
}

export function appPaths() {
  const root = path.join(app, "app/(app)")
  const out = []
  for (const f of walk(root, "page.tsx")) {
    if (path.basename(f) !== "page.tsx") continue
    const segs = path.relative(root, path.dirname(f)).split(path.sep).filter((s) => s && !isGroup(s))
    if (segs[0] === "docs") continue // the docs group
    const src = path.relative(repo, f)
    if (!segs.some((s) => s.startsWith("["))) {
      out.push({ path: "/" + segs.join("/"), group: "app", source: src })
      continue
    }
    const route = segs.join("/")
    if (route === "blocks/[...categories]") {
      const cats = fs.readFileSync(path.join(app, "lib/categories.ts"), "utf8")
      for (const m of cats.matchAll(/slug:\s*"([^"]+)"/g)) out.push({ path: `/blocks/${m[1]}`, group: "app", source: src })
    } else if (route === "charts/[type]") {
      const page = fs.readFileSync(f, "utf8")
      const arr = page.match(/const chartTypes = \[([^\]]*)\]/)
      for (const m of arr[1].matchAll(/"([^"]+)"/g)) out.push({ path: `/charts/${m[1]}`, group: "app", source: src })
    } else {
      throw new Error(`pages.mjs: unhandled dynamic app route ${route} (${src})`)
    }
  }
  out.push({ path: NOT_FOUND_PROBE, group: "app", source: "app/(app)/../_not-found (404 page)" })
  return out
}

function prerendered() {
  const m = JSON.parse(fs.readFileSync(path.join(app, ".next/prerender-manifest.json"), "utf8"))
  return Object.keys(m.routes)
}

// Prerendered routes that are not HTML pages (files, feeds, images, the
// framework's internal error pages).
const NON_PAGE = /^\/(_global-error|_not-found|robots\.txt|rss\.xml|sitemap\.xml|r\/registries\.json)$|\.(jpg|png|json|txt|xml)$|^\/og$/

export function viewPaths() {
  return prerendered()
    .filter((r) => /^\/(view|preview)\//.test(r))
    .sort()
    .map((p) => ({ path: p, group: "view", source: "ui/apps/v4/.next/prerender-manifest.json" }))
}

// Literal (parameter-free) redirect sources of next.config.mjs. A content
// file whose URL is redirected (e.g. components/*/typography.mdx ->
// /docs/typeset) is never served as a page, so it is not listed; --check
// verifies the reference answers it with a redirect.
export function redirectSources() {
  const cfg = fs.readFileSync(path.join(app, "next.config.mjs"), "utf8")
  const block = cfg.slice(cfg.indexOf("redirects()"), cfg.indexOf("rewrites()"))
  return [...block.matchAll(/source:\s*"([^"]+)"/g)].map((m) => m[1]).filter((s) => !/[:(*]/.test(s))
}

export function allPages(groups = ["docs", "app", "view"]) {
  const redirected = new Set(redirectSources())
  const out = []
  if (groups.includes("docs")) out.push(...docsPaths())
  if (groups.includes("app")) out.push(...appPaths())
  if (groups.includes("view")) out.push(...viewPaths())
  const seen = new Set()
  return out.filter((p) => !redirected.has(p.path) && !seen.has(p.path) && seen.add(p.path))
}

async function check(base) {
  const pages = allPages()
  const listed = new Set(pages.map((p) => p.path))
  const problems = []

  const built = prerendered().filter((r) => !NON_PAGE.test(r))
  const builtSet = new Set(built)
  for (const p of pages) if (p.group !== "view" && p.path !== NOT_FOUND_PROBE && !builtSet.has(p.path)) problems.push(`computed but not prerendered: ${p.path} (${p.source})`)
  const redirected = new Set(redirectSources())
  for (const r of built) if (!listed.has(r) && !redirected.has(r)) problems.push(`prerendered but not listed: ${r}`)
  const contentRedirects = docsPaths().filter((p) => redirected.has(p.path))
  for (const p of contentRedirects) {
    const res = await fetch(base + p.path, { redirect: "manual" })
    if (res.status < 300 || res.status >= 400) problems.push(`content page shadowed by a redirect is served ${res.status}: ${p.path}`)
  }
  console.log(`redirects: ${contentRedirects.length} content pages shadowed by next.config.mjs redirects, not listed: ${contentRedirects.map((p) => p.path).join(" ")}`)
  console.log(`build: ${built.length} prerendered HTML routes, ${pages.length} listed paths`)

  const sm = await (await fetch(base + "/sitemap.xml")).text()
  const locs = [...sm.matchAll(/<loc>([^<]+)<\/loc>/g)].map((m) => new URL(m[1]).pathname.replace(/(.)\/$/, "$1"))
  // The sitemap lists every fumadocs page, including the ones a redirect
  // shadows: those are accepted (and reported), anything else is a mismatch.
  for (const l of locs) if (!listed.has(l) && !redirected.has(l)) problems.push(`in sitemap.xml but not listed: ${l}`)
  const smRedirects = locs.filter((l) => redirected.has(l))
  if (smRedirects.length) console.log(`sitemap.xml lists ${smRedirects.length} redirected URLs (upstream behavior, not pages): ${smRedirects.join(" ")}`)
  const extra = pages.filter((p) => p.group !== "view" && p.path !== NOT_FOUND_PROBE && !locs.includes(p.path)).map((p) => p.path)
  console.log(`sitemap.xml: ${locs.length} locs; listed docs/app paths not in the sitemap (expected: pages the sitemap omits): ${extra.join(" ") || "none"}`)

  let ok = 0
  const queue = [...pages]
  await Promise.all(Array.from({ length: 8 }, async () => {
    while (queue.length) {
      const p = queue.shift()
      const res = await fetch(base + p.path, { redirect: "manual" })
      await res.arrayBuffer()
      const want = p.path === NOT_FOUND_PROBE ? 404 : 200
      if (res.status !== want) problems.push(`served ${res.status} (want ${want}): ${p.path}`)
      else ok++
    }
  }))
  console.log(`server ${base}: ${ok}/${pages.length} paths served with the expected status`)

  const groups = {}
  for (const p of pages) groups[p.group] = (groups[p.group] ?? 0) + 1
  console.log("groups:", JSON.stringify(groups))
  for (const pr of problems) console.log("MISMATCH " + pr)
  console.log(problems.length ? `${problems.length} mismatches` : "pages.mjs: list, build, sitemap and server agree")
  return problems.length === 0
}

if (import.meta.url === `file://${process.argv[1]}`) {
  const args = process.argv.slice(2)
  const opt = (n, d) => (args.includes(n) ? args[args.indexOf(n) + 1] : d)
  if (args.includes("--check")) {
    process.exit((await check(opt("--ref", "http://localhost:4000"))) ? 0 : 1)
  }
  const pages = allPages(opt("--group", "docs,app,view").split(","))
  if (args.includes("--json")) console.log(JSON.stringify(pages, null, 2))
  else for (const p of pages) console.log(p.path)
}
