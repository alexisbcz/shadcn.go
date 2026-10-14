// Parity runner: proves a Go example renders and behaves exactly like the
// upstream React example.
//
//   node parity/run.mjs button-demo button-size     # named examples
//   node parity/run.mjs --component button           # every ported example of a component
//   node parity/run.mjs --all                        # every ported example
//
// Options:
//   --themes light,dark   themes to check (default both)
//   --video               also record both sides in real time (webm, for humans)
//   --no-auto             only run scenario files, not the automatic interactions
//   --jobs N              examples checked in parallel (default 4)
//   --keep                keep images of passing captures too
//
// For every example, theme and capture it compares, React vs Go:
//   - pixels: full-viewport screenshots must be identical (0 differing pixels)
//   - DOM: normalized markup of <body> (portals included) must be identical
//   - animations: the running animations' keyframes and timing must match,
//     and frames sampled across their duration must be pixel-identical
//
// Captures: "initial" (after load, animations settled), then the automatic
// interactions (hover, keyboard focus, open a trigger, Escape) unless
// --no-auto, then the steps in parity/scenarios/<example>.json if present.
// Results go to parity/out/<example>/report.json (plus images and DOM dumps
// of every capture that differs); the exit status is non-zero if anything
// differs. node parity/status.mjs summarizes all reports.
import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import { spawn, execFileSync } from "node:child_process"
import { chromium } from "playwright"
import pixelmatch from "pixelmatch"
import { PNG } from "pngjs"

const repo = path.resolve(import.meta.dirname, "..")
const outRoot = path.join(repo, "parity/out")

// --- arguments -------------------------------------------------------------------

const args = process.argv.slice(2)
const flag = (n) => args.includes(n)
const opt = (n, d) => {
  const i = args.indexOf(n)
  return i >= 0 ? args[i + 1] : d
}
const themes = opt("--themes", "light,dark").split(",")
const jobs = Number(opt("--jobs", "4"))
const video = flag("--video")
const auto = !flag("--no-auto")
const keep = flag("--keep")
const positional = args.filter((a, i) => !a.startsWith("--") && !["--themes", "--jobs", "--component"].includes(args[i - 1]))

// --- machine-wide slots ---------------------------------------------------------------

// Many agents run the harness at once. Each run starts a Chromium and a Go
// server; unbounded, they oversubscribe CPU and memory and every run gets
// slower. A run holds one of PARITY_SLOTS (default 4) slots for its whole
// duration; others wait. Slots of dead processes are reclaimed.
async function acquireSlot() {
  const dir = path.join(repo, "parity/.slots")
  fs.mkdirSync(dir, { recursive: true })
  const slots = Number(process.env.PARITY_SLOTS ?? 4)
  let waited = false
  for (;;) {
    for (let i = 0; i < slots; i++) {
      const slot = path.join(dir, `slot-${i}`)
      try {
        fs.mkdirSync(slot)
        fs.writeFileSync(path.join(slot, "pid"), String(process.pid))
        const release = () => fs.rmSync(slot, { recursive: true, force: true })
        process.on("exit", release)
        return release
      } catch {
        let pid = 0
        try { pid = Number(fs.readFileSync(path.join(slot, "pid"), "utf8")) } catch {}
        let alive = false
        try { if (pid) { process.kill(pid, 0); alive = true } } catch {}
        // A slot without pid file may be mid-creation: only reclaim it if old.
        const stale = !pid ? Date.now() - fs.statSync(slot, { throwIfNoEntry: false })?.mtimeMs > 10000 : !alive
        if (stale) fs.rmSync(slot, { recursive: true, force: true })
      }
    }
    if (!waited) {
      console.error(`waiting for a free harness slot (${slots} in use)...`)
      waited = true
    }
    await new Promise((r) => setTimeout(r, 1000))
  }
}

// --- server ----------------------------------------------------------------------

// Several agents may run the harness at once: regenerating the example
// index and building happen under a lock directory.
function withLock(fn) {
  const lock = path.join(repo, "parity/.lock")
  const deadline = Date.now() + 120000
  for (;;) {
    try {
      fs.mkdirSync(lock)
      break
    } catch {
      if (Date.now() > deadline) fs.rmSync(lock, { recursive: true, force: true })
      Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, 200)
    }
  }
  try {
    return fn()
  } finally {
    fs.rmSync(lock, { recursive: true, force: true })
  }
}

async function startServer() {
  const bin = path.join(os.tmpdir(), `shadcn-parity-${process.pid}`)
  withLock(() => {
    execFileSync("go", ["generate", "./examples/all"], { cwd: repo, stdio: "inherit" })
    execFileSync("go", ["build", "-o", bin, "./cmd/parity"], { cwd: repo, stdio: "inherit" })
  })
  const proc = spawn(bin, ["-addr", "127.0.0.1:0", "-root", repo], { cwd: repo, stdio: ["ignore", "pipe", "pipe"] })
  let stderr = ""
  proc.stderr.on("data", (d) => (stderr += d))
  const url = await new Promise((resolve, reject) => {
    let buf = ""
    proc.stdout.on("data", (d) => {
      buf += d
      const m = buf.match(/listening on (http:\/\/\S+)/)
      if (m) resolve(m[1])
    })
    proc.on("exit", (code) => reject(new Error(`parity server exited (${code}):\n${stderr}`)))
  })
  return { url, stop: () => { proc.kill(); fs.rmSync(bin, { force: true }) } }
}

function chromiumPath() {
  if (process.env.PARITY_CHROMIUM) return process.env.PARITY_CHROMIUM
  const cache = path.join(os.homedir(), "Library/Caches/ms-playwright")
  const linuxCache = path.join(os.homedir(), ".cache/ms-playwright")
  for (const dir of [cache, linuxCache]) {
    if (!fs.existsSync(dir)) continue
    const shells = fs.readdirSync(dir).filter((d) => d.startsWith("chromium_headless_shell-")).sort().reverse()
    for (const s of shells) {
      for (const sub of fs.readdirSync(path.join(dir, s))) {
        const p = path.join(dir, s, sub, "chrome-headless-shell")
        if (fs.existsSync(p)) return p
      }
    }
  }
  return undefined // let Playwright use its own (npx playwright install chromium)
}

// --- page helpers (run in the browser) -----------------------------------------------------------

// Serializes <body> into comparable lines. Ids are replaced by their order
// of first appearance (React's useId and the Go runtime generate different
// ids), attribute order and class order are ignored, inline styles are
// compared by their parsed declarations.
function snapshotDOM() {
  const ID_REFS = ["id", "for", "aria-controls", "aria-labelledby", "aria-describedby", "aria-activedescendant", "aria-owns", "aria-details", "aria-errormessage", "list", "form", "headers", "data-testid"]
  const ids = new Map()
  const idOf = (v) => {
    if (!ids.has(v)) ids.set(v, `#id${ids.size}`)
    return ids.get(v)
  }
  const IGNORE = new Set(["nonce", "data-parity-root"])
  const lines = []
  const walk = (node, depth, pathStr) => {
    for (const child of node.childNodes) {
      if (child.nodeType === 3) {
        const t = child.textContent.replace(/\s+/g, " ").trim()
        if (t) lines.push(`${"  ".repeat(depth)}"${t}"`)
        continue
      }
      if (child.nodeType !== 1) continue
      const tag = child.tagName.toLowerCase()
      if (tag === "script" || tag === "template" && child.hasAttribute("data-parity-ignore")) continue
      const attrs = []
      for (const a of child.attributes) {
        if (IGNORE.has(a.name)) continue
        let v = a.value
        if (a.name === "class") v = [...new Set(v.split(/\s+/).filter(Boolean))].sort().join(" ")
        else if (a.name === "style") v = child.style.cssText.split(";").map((d) => d.trim()).filter(Boolean).sort().join("; ")
        else if (ID_REFS.includes(a.name)) v = v.split(/\s+/).filter(Boolean).map(idOf).join(" ")
        else if (a.name === "href" && v.startsWith("#") && v.length > 1) v = idOf(v.slice(1))
        attrs.push(`${a.name}=${JSON.stringify(v)}`)
      }
      attrs.sort()
      // Live properties that are not reflected as attributes.
      if (tag === "input" || tag === "textarea" || tag === "select") {
        attrs.push(`.value=${JSON.stringify(child.value)}`)
        if (child.type === "checkbox" || child.type === "radio") attrs.push(`.checked=${child.checked}`)
      }
      if (document.activeElement === child) attrs.push("(focused)")
      lines.push(`${"  ".repeat(depth)}<${tag}${attrs.length ? " " + attrs.join(" ") : ""}>`)
      walk(child.shadowRoot ?? child, depth + 1)
    }
  }
  walk(document.body, 0)
  return lines
}

// Describes the running animations so both sides can be compared exactly.
function describeAnimations() {
  return document.getAnimations().map((a) => {
    const t = a.effect?.getTiming?.() ?? {}
    const target = a.effect?.target
    return {
      type: a.constructor.name,
      timeline: a.timeline && a.timeline !== document.timeline ? a.timeline.constructor.name : "document",
      name: a.animationName ?? a.transitionProperty ?? "",
      target: target ? `${target.tagName.toLowerCase()}[data-slot=${target.getAttribute("data-slot") ?? ""}]` : "",
      duration: t.duration,
      delay: t.delay,
      easing: t.easing,
      iterations: t.iterations,
      direction: t.direction,
      fill: t.fill,
      keyframes: a.effect?.getKeyframes?.().map(({ computedOffset, ...k }) => k) ?? [],
    }
  }).sort((x, y) => JSON.stringify(x).localeCompare(JSON.stringify(y)))
}

// Forces a full re-raster: Chromium's partial raster can leave anti-aliased
// edges one gray level off depending on paint history, which differs between
// two pages showing the same thing. The DOM is left unchanged.
async function forceRepaint() {
  const raf2 = () => new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r)))
  const b = document.body
  const hadStyle = b.hasAttribute("style")
  const prev = b.style.willChange
  b.style.willChange = "transform"
  await raf2()
  b.style.willChange = prev
  if (!hadStyle && !b.getAttribute("style")) b.removeAttribute("style")
  await raf2()
}

// Waits (up to 1s) for CSS transitions started by the page load itself.
// Chromium starts transitions from UA default styles on some server-rendered
// elements (inside size containers); a client-rendered page never does.
function settleLoadTransitions() {
  const running = document.getAnimations().filter((a) => a.constructor.name === "CSSTransition")
  return Promise.race([Promise.all(running.map((a) => a.finished.catch(() => {}))), new Promise((r) => setTimeout(r, 1000))])
}

// Finishes finite animations and parks infinite ones at 0, so a capture
// shows the settled state on both sides.
// Scroll-driven animations (animation-timeline: scroll()/view()) have no
// absolute time: their state follows the scroll position, which the
// screenshots already compare. Seeking them throws, so they are skipped by
// settle/seek/end and compared only by their definition.
function settleAnimations() {
  for (const a of document.getAnimations()) {
    if (a.timeline && a.timeline !== document.timeline) continue
    const end = a.effect?.getComputedTiming?.().endTime
    if (Number.isFinite(end)) a.finish()
    else {
      a.pause()
      a.currentTime = 0
    }
  }
}

function seekAnimations(t) {
  for (const a of document.getAnimations()) {
    if (a.timeline && a.timeline !== document.timeline) continue
    a.pause()
    a.currentTime = t
  }
}

function maxAnimationEnd() {
  let max = 0
  for (const a of document.getAnimations()) {
    if (a.timeline && a.timeline !== document.timeline) continue
    const end = a.effect?.getComputedTiming?.().endTime
    if (Number.isFinite(end)) max = Math.max(max, end)
  }
  return max
}

// --- one side ------------------------------------------------------------------------

class Side {
  constructor(browser, base, kind, name, theme, dir) {
    Object.assign(this, { browser, base, kind, name, theme, dir })
  }

  async open() {
    this.context = await this.browser.newContext({
      viewport: { width: 1024, height: 768 },
      deviceScaleFactor: 1,
      colorScheme: this.theme,
      reducedMotion: "no-preference",
      ...(video ? { recordVideo: { dir: path.join(this.dir, "video", this.kind), size: { width: 1024, height: 768 } } } : {}),
    })
    this.page = await this.context.newPage()
    this.errors = []
    this.page.on("pageerror", (e) => this.errors.push(String(e)))
    this.page.on("console", (m) => m.type() === "error" && this.errors.push(m.text()))
    this.cdp = await this.context.newCDPSession(this.page)
    await this.cdp.send("Animation.enable")
    await this.page.goto(`${this.base}/${this.kind}/${this.name}?theme=${this.theme}`)
    if (this.kind === "react") {
      await this.page.waitForFunction(() => window.__parityReady || window.__parityError, null, { timeout: 15000 })
      const err = await this.page.evaluate(() => window.__parityError)
      if (err) throw new Error(`react: ${err}`)
    } else {
      await this.page.waitForFunction(() => window.__shadcnReady, null, { timeout: 15000 })
    }
    await this.page.evaluate(() => document.fonts.ready)
    await this.page.evaluate(settleLoadTransitions)
    await this.idle()
  }

  async idle() {
    await this.page.evaluate(() => new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r))))
  }

  async freeze(on) {
    await this.cdp.send("Animation.setPlaybackRate", { playbackRate: on ? 0 : 1 })
  }

  async step(s) {
    const p = this.page
    const target = s.target ? p.locator(s.target).nth(s.nth ?? 0) : null
    const force = s.force ? { force: true } : {}
    const center = async () => {
      const box = await target.boundingBox()
      return { x: box.x + box.width / 2, y: box.y + box.height / 2 }
    }
    switch (s.do) {
      case "hover": await target.hover(force); break
      case "click": await target.click(force); break
      case "dblclick": await target.dblclick(force); break
      case "down": await p.mouse.down({ button: s.button ?? "left" }); break
      case "up": await p.mouse.up({ button: s.button ?? "left" }); break
      case "drag": {
        const c = await center()
        await p.mouse.move(c.x, c.y)
        await p.mouse.down()
        await p.mouse.move(c.x + (s.dx ?? 0), c.y + (s.dy ?? 0), { steps: s.steps ?? 10 })
        await p.mouse.up()
        break
      }
      case "focus": await target.focus(); break
      case "press": await (target ?? p.keyboard).press(s.key); break
      case "type": await (target ? target.pressSequentially(s.text) : p.keyboard.type(s.text)); break
      case "mouse": await p.mouse.move(s.x, s.y, { steps: s.steps ?? 1 }); break
      case "wait": await p.waitForTimeout(s.ms ?? 100); break
      case "scroll": await (target ? target.evaluate((el, y) => (el.scrollTop = y), s.y ?? 0) : p.mouse.wheel(0, s.y ?? 0)); break
      case "eval": await p.evaluate(s.js); break
      default: throw new Error(`unknown step ${JSON.stringify(s)}`)
    }
  }

  async exists(selector) {
    return (await this.page.locator(selector).count()) > 0
  }

  async shot() {
    await this.page.evaluate(forceRepaint)
    return this.page.screenshot({ fullPage: true, animations: "allow" })
  }

  async capture() {
    return {
      png: await this.shot(),
      dom: await this.page.evaluate(snapshotDOM),
    }
  }

  async close() {
    await this.context.close()
  }
}

// --- comparison ------------------------------------------------------------------------

function comparePNG(a, b) {
  const A = PNG.sync.read(a)
  const B = PNG.sync.read(b)
  if (A.width !== B.width || A.height !== B.height) {
    return { pixels: Math.max(A.width * A.height, B.width * B.height), size: `${A.width}x${A.height} vs ${B.width}x${B.height}`, diff: null }
  }
  const diff = new PNG({ width: A.width, height: A.height })
  const pixels = pixelmatch(A.data, B.data, diff.data, A.width, A.height, { threshold: 0, includeAA: true })
  return { pixels, diff: PNG.sync.write(diff) }
}

function diffLines(a, b, max = 40) {
  // Small LCS-free diff: walk both, report the first divergences in order.
  const out = []
  let i = 0
  let j = 0
  while ((i < a.length || j < b.length) && out.length < max) {
    if (a[i] === b[j]) {
      i++, j++
      continue
    }
    const k = b.indexOf(a[i], j)
    const l = a.indexOf(b[j], i)
    if (k >= 0 && (l < 0 || k - j <= l - i)) {
      while (j < k && out.length < max) out.push(`+ go:    ${b[j++]}`)
    } else if (l >= 0) {
      while (i < l && out.length < max) out.push(`- react: ${a[i++]}`)
    } else {
      if (i < a.length) out.push(`- react: ${a[i++]}`)
      if (j < b.length) out.push(`+ go:    ${b[j++]}`)
    }
  }
  return out
}

// --- scenarios ------------------------------------------------------------------------

const INTERACTIVE = 'button:not([disabled]), a[href], input:not([disabled]), textarea:not([disabled]), select:not([disabled]), [role="button"], [role="tab"], [role="switch"], [role="checkbox"], [role="combobox"], [role="slider"], [tabindex="0"]'
const TRIGGER = '[data-slot$="-trigger"], [aria-haspopup]:not([aria-haspopup="false"]), [aria-expanded]'

async function autoScenario(react) {
  const steps = []
  if (await react.exists(`#preview :is(${INTERACTIVE})`)) {
    steps.push({ label: "hover", steps: [{ do: "hover", target: `#preview :is(${INTERACTIVE})`, force: true }], frames: "auto" })
    steps.push({ label: "focus-visible", steps: [{ do: "mouse", x: 0, y: 0 }, { do: "press", key: "Tab" }], frames: "auto" })
  }
  if (await react.exists(`#preview :is(${TRIGGER})`)) {
    steps.push({ label: "open", steps: [{ do: "click", target: `#preview :is(${TRIGGER})`, force: true }], frames: "auto" })
    steps.push({ label: "escape", steps: [{ do: "press", key: "Escape" }], frames: "auto" })
  }
  return steps
}

function scenarioFile(name) {
  const f = path.join(repo, "parity/scenarios", name + ".json")
  return fs.existsSync(f) ? JSON.parse(fs.readFileSync(f, "utf8")) : null
}

// --- one example ------------------------------------------------------------------------

async function checkExample(browser, base, name) {
  const dir = path.join(outRoot, name)
  fs.rmSync(dir, { recursive: true, force: true })
  fs.mkdirSync(dir, { recursive: true })
  const report = { name, pass: true, themes: {} }

  for (const theme of themes) {
    const tdir = path.join(dir, theme)
    fs.mkdirSync(tdir, { recursive: true })
    const react = new Side(browser, base, "react", name, theme, tdir)
    const go = new Side(browser, base, "go", name, theme, tdir)
    const result = { captures: [], errors: {} }
    report.themes[theme] = result
    try {
      await Promise.all([react.open(), go.open()])

      const record = async (label, frames) => {
        const cap = { label, pass: true }
        if (frames) {
          const [ra, ga] = await Promise.all([react.page.evaluate(describeAnimations), go.page.evaluate(describeAnimations)])
          cap.animations = { react: ra, go: ga, equal: JSON.stringify(ra) === JSON.stringify(ga) }
          if (!cap.animations.equal) cap.pass = false
          const end = Math.max(await react.page.evaluate(maxAnimationEnd), await go.page.evaluate(maxAnimationEnd))
          const times = end > 0 ? [0, 0.25, 0.5, 0.75].map((f) => Math.round(f * end)) : []
          cap.frames = []
          for (const t of times) {
            // A real mismatch is deterministic; under heavy load Chromium
            // occasionally composites a seeked frame late on one side, so a
            // differing frame is re-seeked and re-captured up to twice.
            let r, g, cmp
            for (let attempt = 0; attempt < 3; attempt++) {
              await Promise.all([react.page.evaluate(seekAnimations, t), go.page.evaluate(seekAnimations, t)])
              // No forced repaint here: toggling will-change on <body> while
              // composited animations are paused at a seeked time makes the
              // compositor's frame nondeterministic. Two frames let it commit.
              await Promise.all([react.idle(), go.idle()])
              ;[r, g] = await Promise.all([
                react.page.screenshot({ fullPage: true, animations: "allow" }),
                go.page.screenshot({ fullPage: true, animations: "allow" }),
              ])
              cmp = comparePNG(r, g)
              if (cmp.pixels === 0) break
            }
            cap.frames.push({ t, pixels: cmp.pixels })
            if (cmp.pixels > 0 || keep) {
              fs.writeFileSync(path.join(tdir, `${label}@${t}ms-react.png`), r)
              fs.writeFileSync(path.join(tdir, `${label}@${t}ms-go.png`), g)
              if (cmp.diff) fs.writeFileSync(path.join(tdir, `${label}@${t}ms-diff.png`), cmp.diff)
            }
            if (cmp.pixels > 0) cap.pass = false
          }
        }
        await Promise.all([react.page.evaluate(settleAnimations), go.page.evaluate(settleAnimations)])
        await Promise.all([react.idle(), go.idle()])
        let r, g, cmp
        for (let attempt = 0; attempt < 3; attempt++) {
          ;[r, g] = await Promise.all([react.capture(), go.capture()])
          cmp = comparePNG(r.png, g.png)
          if (cmp.pixels === 0) break
        }
        cap.pixels = cmp.pixels
        if (cmp.size) cap.size = cmp.size
        cap.dom = diffLines(r.dom, g.dom)
        if (cmp.pixels > 0 || cap.dom.length > 0) cap.pass = false
        if (!cap.pass || keep) {
          fs.writeFileSync(path.join(tdir, `${label}-react.png`), r.png)
          fs.writeFileSync(path.join(tdir, `${label}-go.png`), g.png)
          if (cmp.diff) fs.writeFileSync(path.join(tdir, `${label}-diff.png`), cmp.diff)
          fs.writeFileSync(path.join(tdir, `${label}-react.dom.txt`), r.dom.join("\n"))
          fs.writeFileSync(path.join(tdir, `${label}-go.dom.txt`), g.dom.join("\n"))
        }
        result.captures.push(cap)
        if (!cap.pass) report.pass = false
      }

      await record("initial", false)

      const scenarios = [...(auto ? await autoScenario(react) : []), ...(scenarioFile(name) ?? [])]
      for (const sc of scenarios) {
        if (sc.reset) {
          await Promise.all([react.close(), go.close()])
          await Promise.all([react.open(), go.open()])
        }
        const frames = sc.frames !== false
        if (frames) await Promise.all([react.freeze(true), go.freeze(true)])
        for (const s of sc.steps) {
          await Promise.all([react.step(s), go.step(s)])
        }
        await Promise.all([react.idle(), go.idle()])
        await record(sc.label, frames)
        if (frames) await Promise.all([react.freeze(false), go.freeze(false)])
      }

      if (video) {
        // Replay the scenarios at real speed on fresh pages for the videos.
        for (const side of [react, go]) {
          await side.page.reload()
          await side.page.waitForTimeout(300)
          for (const sc of scenarios) {
            for (const s of sc.steps) await side.step(s)
            await side.page.waitForTimeout(600)
          }
        }
      }
    } catch (e) {
      result.errors.run = String(e.stack ?? e)
      report.pass = false
    } finally {
      result.errors.react = react.errors
      result.errors.go = go.errors
      if (go.errors?.length) report.pass = false
      await Promise.allSettled([react.close?.(), go.close?.()])
    }
  }
  fs.writeFileSync(path.join(dir, "report.json"), JSON.stringify(report, null, 2))
  return report
}

// --- main --------------------------------------------------------------------------------

const releaseSlot = await acquireSlot()
const server = await startServer()
let exitCode = 0
try {
  const list = await (await fetch(server.url + "/api/examples")).json()
  let names
  if (flag("--all")) names = list.filter((e) => e.ported).map((e) => e.name)
  else if (opt("--component")) names = list.filter((e) => e.ported && e.component === opt("--component")).map((e) => e.name)
  else names = positional
  const unknown = names.filter((n) => !list.some((e) => e.name === n && e.ported))
  if (unknown.length) throw new Error(`not ported or unknown: ${unknown.join(", ")}`)
  if (!names.length) throw new Error("no examples selected (names, --component c or --all)")

  const browser = await chromium.launch({ executablePath: chromiumPath() })
  const reports = []
  const queue = [...names]
  await Promise.all(Array.from({ length: Math.min(jobs, queue.length) }, async () => {
    while (queue.length) {
      const name = queue.shift()
      const r = await checkExample(browser, server.url, name)
      reports.push(r)
      const caps = Object.values(r.themes).flatMap((t) => t.captures)
      const bad = caps.filter((c) => !c.pass).map((c) => c.label)
      console.log(`${r.pass ? "PASS" : "FAIL"} ${name} (${caps.length} captures${bad.length ? `; differs: ${[...new Set(bad)].join(", ")}` : ""})`)
      for (const [theme, t] of Object.entries(r.themes)) {
        if (t.errors.run) console.log(`     ${theme}: ${t.errors.run.split("\n")[0]}`)
        if (t.errors.go?.length) console.log(`     ${theme}: go page errors: ${t.errors.go.slice(0, 3).join(" | ")}`)
      }
    }
  }))
  await browser.close()

  const failed = reports.filter((r) => !r.pass)
  console.log(`\n${reports.length - failed.length}/${reports.length} examples pixel-, DOM- and animation-identical. Reports: parity/out/`)
  if (failed.length) exitCode = 1
} catch (e) {
  console.error(e.message ?? e)
  exitCode = 2
} finally {
  server.stop()
  releaseSlot()
}
process.exit(exitCode)
