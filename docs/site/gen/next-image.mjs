// Vendors the next/image optimizer output of the reference build for every
// local <Image> of the docs content and the app pages (app/(app)/**.tsx), so the Go server serves the same bytes at
// /_next/image (docs/site/nextimage.go). next/image resizes with sharp and
// re-encodes to webp q75; reproducing that in Go is not pixel-identical.
//
// Usage: node docs/site/gen/next-image.mjs [http://localhost:4000]
import fs from "node:fs"
import path from "node:path"

const ref = process.argv[2] ?? "http://localhost:4000"
const root = path.resolve(import.meta.dirname, "../../..")
const content = path.join(root, "ui/apps/v4/content/docs")
const out = path.join(root, "docs/site/static/nextimage")

// images config of the reference build (.next/required-server-files.json).
const sizes = [32, 48, 64, 96, 128, 256, 384, 640, 750, 828, 1080, 1200, 1920, 2048, 3840]
const widths = (w) => [...new Set([w, w * 2].map((x) => sizes.find((s) => s >= x) ?? sizes.at(-1)))]

const app = path.join(root, "ui/apps/v4/app/(app)")
const files = [
  ...fs.readdirSync(content, { recursive: true }).filter((f) => f.endsWith(".mdx")).map((f) => path.join(content, f)),
  ...fs.readdirSync(app, { recursive: true }).filter((f) => f.endsWith(".tsx")).map((f) => path.join(app, f)),
]
const want = new Map()
for (const f of files) {
  const src = fs.readFileSync(f, "utf8")
  for (const [tag] of src.matchAll(/<Image\b[^>]*?\/>/gs)) {
    const s = tag.match(/src="(\/[^"]+)"/)?.[1]
    const w = tag.match(/width=(?:"(\d+)"|\{(\d+)\})/)
    if (!s || !w) continue
    for (const x of widths(Number(w[1] ?? w[2]))) want.set(`${x}${s}`, { src: s, w: x })
  }
}
// Plus every local image the reference's app pages request (block
// screenshots built from data, not literal <Image> tags).
const { execFileSync } = await import("node:child_process")
const appPages = execFileSync("node", [path.join(root, "parity/docs/pages.mjs"), "--group", "app"], { encoding: "utf8" })
  .split("\n").filter((p) => p.startsWith("/"))
for (const p of appPages) {
  const html = await (await fetch(ref + p)).text()
  for (const [, u, w] of html.matchAll(/\/_next\/image\?url=(%2F[^&"]+)&(?:amp;)?w=(\d+)/g)) {
    const src = decodeURIComponent(u)
    want.set(`${w}${src}`, { src, w: Number(w) })
  }
}
// The screenshots of block previews (<ComponentPreview type="block">).
for (const f of files.filter((f) => f.endsWith(".mdx"))) {
  for (const [tag] of fs.readFileSync(f, "utf8").matchAll(/<ComponentPreview\b[^>]*?\/>/gs)) {
    const name = tag.match(/name="([^"]+)"/)?.[1]
    if (!name || !/type="block"/.test(tag)) continue
    for (const theme of ["light", "dark"]) {
      for (const w of [1920, 3840]) want.set(`${w}/r/styles/new-york/${name}-${theme}.png`, { src: `/r/styles/new-york/${name}-${theme}.png`, w })
    }
  }
}
// remote/ belongs to next-image-remote.mjs.
for (const d of fs.existsSync(out) ? fs.readdirSync(out) : []) {
  if (d !== "remote") fs.rmSync(path.join(out, d), { recursive: true, force: true })
}
for (const { src, w } of want.values()) {
  const url = `${ref}/_next/image?url=${encodeURIComponent(src)}&w=${w}&q=75`
  const res = await fetch(url, { headers: { accept: "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8" } })
  if (!res.ok) {
    // Missing upstream (e.g. a block without screenshots): the Go server
    // answers 400 as well.
    console.warn(`skip ${url}: ${res.status}`)
    continue
  }
  const ext = res.headers.get("content-type") === "image/webp" ? ".webp" : ""
  const file = path.join(out, String(w), src + ext)
  fs.mkdirSync(path.dirname(file), { recursive: true })
  fs.writeFileSync(file, Buffer.from(await res.arrayBuffer()))
}
console.log(`${want.size} images in ${path.relative(root, out)}`)
