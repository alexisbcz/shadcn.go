import fs from "node:fs"
import { PNG } from "pngjs"
const [dir, label, x, y, w, h, out] = process.argv.slice(2)
const X = +x, Y = +y, W = +w, H = +h
const imgs = ["ref", "target", "diff"].map((k) => PNG.sync.read(fs.readFileSync(`${dir}/${label}-${k}.png`)))
const o = new PNG({ width: W * 3 + 20, height: H }); o.data.fill(255)
imgs.forEach((img, k) => { for (let j = 0; j < H; j++) for (let i = 0; i < W; i++) { const sx = X + i, sy = Y + j; if (sx >= img.width || sy >= img.height) continue; const s = (sy * img.width + sx) * 4, d = (j * o.width + k * (W + 10) + i) * 4; img.data.copy(o.data, d, s, s + 4) } })
fs.writeFileSync(out, PNG.sync.write(o))
