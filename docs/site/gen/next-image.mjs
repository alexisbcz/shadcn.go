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
// remote/ belongs to next-image-remote.mjs.
for (const d of fs.existsSync(out) ? fs.readdirSync(out) : []) {
  if (d !== "remote") fs.rmSync(path.join(out, d), { recursive: true, force: true })
}
for (const { src, w } of want.values()) {
  const url = `${ref}/_next/image?url=${encodeURIComponent(src)}&w=${w}&q=75`
  const res = await fetch(url, { headers: { accept: "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8" } })
  if (!res.ok) throw new Error(`${url}: ${res.status}`)
  const ext = res.headers.get("content-type") === "image/webp" ? ".webp" : ""
  const file = path.join(out, String(w), src + ext)
  fs.mkdirSync(path.dirname(file), { recursive: true })
  fs.writeFileSync(file, Buffer.from(await res.arrayBuffer()))
}
console.log(`${want.size} images in ${path.relative(root, out)}`)
