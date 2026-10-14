// Vendors the next/image optimizer output of the reference build for the
// remote images of the docs previews (next/image in ui/apps/v4/examples/radix,
// e.g. https://avatar.vercel.sh/...), so the Go server serves the same bytes
// at /_next/image?url=https://... (docs/site/nextimage.go). Files land under
// static/nextimage/remote/<w>/<sha1 of the url>.webp.
//
// The URL list is read from the reference's own pages: every component page
// of the radix base (previews are radix-nova whatever the base).
//
// Usage: node docs/site/gen/next-image-remote.mjs [http://localhost:4000]
import crypto from "node:crypto"
import fs from "node:fs"
import path from "node:path"

const ref = process.argv[2] ?? "http://localhost:4000"
const root = path.resolve(import.meta.dirname, "../../..")
const out = path.join(root, "docs/site/static/nextimage/remote")
const pages = fs.readdirSync(path.join(root, "ui/apps/v4/content/docs/components/radix"))
  .filter((f) => f.endsWith(".mdx")).map((f) => `/docs/components/radix/${f.slice(0, -4)}`)

const want = new Map()
for (const p of pages) {
  const html = await (await fetch(ref + p)).text()
  for (const [, u, w] of html.matchAll(/\/_next\/image\?url=(https?%3A[^&"]+)&(?:amp;)?w=(\d+)/g)) {
    const src = decodeURIComponent(u)
    want.set(`${w} ${src}`, { src, w: Number(w) })
  }
}
fs.rmSync(out, { recursive: true, force: true })
for (const { src, w } of want.values()) {
  const url = `${ref}/_next/image?url=${encodeURIComponent(src)}&w=${w}&q=75`
  const res = await fetch(url, { headers: { accept: "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8" } })
  if (!res.ok) throw new Error(`${url}: ${res.status}`)
  if (res.headers.get("content-type") !== "image/webp") throw new Error(`${url}: ${res.headers.get("content-type")}`)
  const file = path.join(out, String(w), crypto.createHash("sha1").update(src).digest("hex") + ".webp")
  fs.mkdirSync(path.dirname(file), { recursive: true })
  fs.writeFileSync(file, Buffer.from(await res.arrayBuffer()))
}
console.log(`${want.size} remote images in ${path.relative(root, out)}`)
