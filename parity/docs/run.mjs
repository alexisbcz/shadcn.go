// Page-level parity harness for the documentation site (DOCS_PLAN.md, D1).
//
// Compares two servers page by page: the reference (the upstream Next.js
// docs site, parity/docs/reference.sh, default http://localhost:4000) and a
// candidate (--target, the Go docs server). The bar is exact equality.
//
// Usage
//
//   node parity/docs/run.mjs --target http://localhost:8080 /docs/components/radix/button
//   node parity/docs/run.mjs --target URL --group docs           # every docs page (parity/docs/pages.mjs)
//   node parity/docs/run.mjs --target URL --all                  # docs + app routes + iframe views
//   node parity/docs/run.mjs --self --all                        # the reference against itself (determinism)
//   node parity/docs/run.mjs --target URL --grep '^/docs/components/radix/'
//
// Options
//
//   --ref URL             reference server (default http://localhost:4000)
//   --target URL          candidate server (required unless --self)
//   --self                candidate = reference, in independent browser contexts
//   --all | --group g,..  page groups from pages.mjs: docs, app, view
//   --grep RE             keep only paths matching RE
//   --shard i/n           run the i-th of n slices of the page list (1-based)
//   --themes light,dark   default both
//   --viewports desktop,laptop,mobile
//                         desktop 1440x900, laptop 1024x768, mobile 390x844 (touch, mobile emulation)
//   --no-scenarios        initial captures only
//   --only label,..       run only these scenarios (plus "initial")
//   --no-frames           skip the animation comparison in scenarios
//   --jobs N              page jobs in parallel (default 3); each job is one page x theme x viewport
//   --keep                keep images and DOM dumps of passing captures too
//   --out DIR             reports (default parity/docs/out)
//   --quiet MS            DOM-quiet window of the settle wait (default 200)
//
// What is compared, per page, theme, viewport and capture
//
//   pixels      screenshots, 0 differing pixels allowed (pixelmatch threshold 0,
//               anti-aliasing included). "initial" is a full-page screenshot;
//               scenario captures show the viewport (where the interaction
//               is) unless the scenario sets full: true. A difference is
//               located: the report lists the differing regions (32px grid
//               clusters, document coordinates) and the deepest reference
//               element covering each region.
//   DOM         the normalized document (rules below) must be identical; each
//               differing line is reported with its ancestor path.
//   state       URL path + hash after the steps, window scroll position,
//               document.title, HTTP status of the first load, clipboard
//               text after copy actions.
//   animations  in scenarios: the document timeline is frozen through the
//               DevTools protocol before the steps (as in parity/run.mjs);
//               the running animations' definitions (keyframes, duration,
//               delay, easing, iterations, direction, fill, target) must be
//               identical, then frames at 0, 25, 50 and 75% of the longest
//               animation are compared pixel-exactly. Scroll-driven
//               animations (animation-timeline: scroll()/view()) are compared
//               by definition and current progress, never seeked.
//   errors      uncaught page errors on the candidate fail the job.
//
// DOM normalization rules (snapshotDOM). Everything else is compared verbatim.
//
//   N1 Elements dropped: <script> (all of them: the __next_f RSC payload
//      pushes, Next's chunk loaders, inline bootstrap scripts; their effects
//      are compared through the DOM and the pixels), <next-route-announcer>
//      (Next's client-navigation live region, outside the visible document),
//      <noscript>, <template>, and a <div hidden> whose only content is
//      comments (the <!--$--><!--/$--> Suspense boundary outlet Next streams
//      for metadata).
//   N2 In <head>, only what describes the document is kept: <title>, <meta>
//      (except <meta name="next-size-adjust">, a next/font marker), and
//      <link> with rel canonical, icon, shortcut icon, apple-touch-icon,
//      manifest, author, alternate. Asset delivery (<link rel=preload,
//      stylesheet, modulepreload, preconnect, dns-prefetch, expect>, <style>)
//      is dropped: the resulting styles are compared through the pixels.
//      Anywhere else, <link rel=preload|modulepreload|stylesheet|expect> is
//      dropped too (React 19 hoists resource hints).
//   N3 URLs: the hashed file name of every /_next/static/<dir>/<file> becomes
//      /_next/static/<dir>/[hash].<ext> (build hashes, e.g. fonts and static
//      images). Each side's own origin (http://host:port) becomes {origin}.
//      Other URLs (including /_next/image?url=...&w=...&q=...) are verbatim.
//   N4 next/font CSS-module classes "<family>_<hex8>-module__<hash>__<local>"
//      become "next-font:<family>:<local>" (the class hash is a build hash;
//      which font is used is checked through the pixels).
//   N5 Generated ids: React useId values (_R_..._, «r..», :r..:) and the ids
//      of shadcn.go's uid package (radix-|base-ui-|aria-|react-aria- + 10
//      base32 characters) are replaced in every attribute by #idN in order
//      of first appearance, keeping any prefix and suffix
//      ("radix-_R_43a_-trigger-cli" -> "radix-#id3-trigger-cli"). Content
//      ids (heading slugs, example ids such as "email") are kept verbatim.
//   N6 class: tokens sorted (order is not significant; duplicates kept).
//      style: parsed declarations, sorted. Attributes sorted by name.
//      nonce and data-parity-* attributes ignored (the Go pages run under a
//      nonce CSP).
//   N7 Text: whitespace collapsed and trimmed, except inside <pre> and
//      <textarea>, where text is compared verbatim. Comments ignored.
//   N8 Live state added: .value of inputs, .checked of checkboxes/radios,
//      (focused) on document.activeElement; same-origin <iframe> documents
//      are snapshotted inline under their <iframe>.
//
// Nondeterminism neutralized identically on both sides (found by the
// self-comparison, --self; each item is applied to the reference and the
// candidate in the same way)
//
//   D1 Time: Date is fixed at PARITY_NOW (page.clock.setFixedTime; timers
//      keep running). Server-rendered dates on the candidate must use the
//      same instant: every request carries "x-parity-now: <ISO time>" (stub
//      contract for the Go server; the reference is prerendered at build).
//   D2 Math.random is a seeded PRNG (mulberry32, seed 1), fresh per document.
//   D3 Clipboard: navigator.clipboard.writeText/readText and the "copy"
//      event are captured per page (window.__parityClipboard) instead of the
//      browser-wide system clipboard, which parallel jobs would share.
//   D4 Lazy content: <img loading=lazy> and <iframe loading=lazy> are loaded
//      eagerly (loading set to eager, then the attribute restored once
//      loaded), then every image is decoded before a capture.
//   D5 Fonts: document.fonts.ready in every frame before a capture.
//   D6 Hydration / late effects: a capture waits for network idle (no request
//      in flight for 150ms) and for 200ms without DOM mutations (max 10s) in
//      every frame.
//   D7 Animations: CSS transitions started by the page load are awaited
//      (max 1s); before each capture, finite animations are finished and
//      infinite ones (spinners, shimmer, marquees) paused at 0, in every
//      frame. Scroll-driven animations are left to their scroll position.
//   D8 Rendering: a forced full repaint before each screenshot (as in
//      parity/run.mjs); the caret is hidden; deviceScaleFactor 1; one
//      Chromium build (the headless shell parity/run.mjs uses); timezone
//      UTC, locale en-US; reducedMotion no-preference.
//   D9 Scroll: every job starts at the top in a fresh context (no
//      sessionStorage, so the sidebar scroll-restore script sees no state);
//      scroll steps use instant scrolling.
//
// Steps (scenarios.mjs and parity/docs/scenarios/*.json)
//
//   hover, click, dblclick, focus      target (CSS, Playwright syntax), nth, force,
//                                      nav (wait for the URL to change)
//   press                              key (Playwright key, e.g. "ControlOrMeta+k"), optional target
//   type                               text (into the focused element)
//   mouse                              x, y: move the pointer
//   scroll                             y: wheel delta, or target + y: element scrollTop
//   scrollTo                           y (px) or fraction (of the scrollable height): window, instant
//   wait                               ms
//   settle                             network idle + DOM quiet + fonts + images
//   clipboard                          record the captured clipboard text in the capture
//   eval                               js
//
// Targets resolve to the nth VISIBLE match.
//
// Output: <out>/<page>/<theme>-<viewport>/report.json, plus for each
// differing capture <label>-ref.png, <label>-target.png, <label>-diff.png,
// <label>-{ref,target}.dom.txt; <out>/summary.json lists every job. The exit
// status is 1 if anything differs, 2 on a harness error.
import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import { chromium } from "playwright"
import pixelmatch from "pixelmatch"
import { PNG } from "pngjs"

import { allPages } from "./pages.mjs"
import { scenarios as chromeScenarios } from "./scenarios.mjs"

const here = import.meta.dirname

// --- arguments -------------------------------------------------------------------------

const args = process.argv.slice(2)
const VALUED = ["--quiet", "--ref", "--target", "--group", "--grep", "--shard", "--themes", "--viewports", "--only", "--jobs", "--out"]
const flag = (n) => args.includes(n)
const opt = (n, d) => (args.includes(n) ? args[args.indexOf(n) + 1] : d)
const positional = args.filter((a, i) => !a.startsWith("--") && !VALUED.includes(args[i - 1]))

const REF = opt("--ref", "http://localhost:4000").replace(/\/$/, "")
const SELF = flag("--self")
const TARGET = (SELF ? REF : opt("--target", "")).replace(/\/$/, "")
const themes = opt("--themes", "light,dark").split(",")
export const VIEWPORTS = {
  desktop: { viewport: { width: 1440, height: 900 } },
  laptop: { viewport: { width: 1024, height: 768 } },
  mobile: { viewport: { width: 390, height: 844 }, isMobile: true, hasTouch: true },
}
const viewports = opt("--viewports", "desktop,laptop,mobile").split(",")
const jobsN = Number(opt("--jobs", "3"))
const keep = flag("--keep")
const withScenarios = !flag("--no-scenarios")
const withFrames = !flag("--no-frames")
const only = opt("--only") ? new Set(opt("--only").split(",")) : null
const outRoot = path.resolve(opt("--out", path.join(here, "out")))
const QUIET = Number(opt("--quiet", "200"))
const QUIET_LOAD = 600
export const PARITY_NOW = "2026-01-15T12:00:00.000Z"

// --- browser-side helpers (serialized into the page) ---------------------------------------

// D2 + D3: installed before any page script, in every frame.
const PRELUDE = `(() => {
  let s = 1
  Math.random = function () {
    s |= 0; s = (s + 0x6d2b79f5) | 0
    let t = Math.imul(s ^ (s >>> 15), 1 | s)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
  window.__parityClipboard = null
  const clip = {
    writeText: (t) => { window.__parityClipboard = String(t); return Promise.resolve() },
    readText: () => Promise.resolve(window.__parityClipboard ?? ""),
    write: () => Promise.resolve(),
    read: () => Promise.resolve([]),
  }
  try { Object.defineProperty(Navigator.prototype, "clipboard", { get: () => clip, configurable: true }) } catch {}
  document.addEventListener("copy", (e) => {
    setTimeout(() => {
      if (e.clipboardData && e.clipboardData.getData("text/plain")) window.__parityClipboard = e.clipboardData.getData("text/plain")
      else if (!e.defaultPrevented) window.__parityClipboard = String(document.getSelection() ?? "")
    })
  }, true)
})()`

// Serializes the document into comparable lines, see "DOM normalization
// rules" above. Each line carries the path of its ancestors for reports.
function snapshotDOM(origins) {
  const ids = new Map()
  const idOf = (v) => {
    if (!ids.has(v)) ids.set(v, `#id${ids.size}`)
    return ids.get(v)
  }
  const REACT_ID = /_R_[0-9A-Za-z]+_|«[^»]+»|:r[0-9a-z]+:/g
  const GO_UID = /(^|[\s#])((?:radix|base-ui|aria|react-aria)-)([a-z2-7]{10})(?=$|[\s-])/g
  const genIds = (v) => v.replace(REACT_ID, (m) => idOf(m)).replace(GO_UID, (_, pre, prefix, core) => pre + prefix + idOf(core))
  const url = (v) => {
    for (const o of origins) v = v.split(o).join("{origin}")
    return v.replace(/\/_next\/static\/([\w-]+)\/[^\s"'),?#]+?(\.[a-z0-9]+)(?=[\s"'),?#]|$)/g, "/_next/static/$1/[hash]$2")
  }
  const fontClass = (c) => c.replace(/^([a-z0-9_]+?)_[0-9a-f]{8}-module__[\w-]+__(\w+)$/, "next-font:$1:$2")
  const HEAD_LINK = new Set(["canonical", "icon", "shortcut icon", "apple-touch-icon", "manifest", "author", "alternate"])
  const DROP_LINK = /(^|\s)(preload|modulepreload|stylesheet|expect|preconnect|dns-prefetch)(\s|$)/
  const lines = []
  const paths = []
  const label = (el) => {
    let s = el.tagName.toLowerCase()
    const slot = el.getAttribute("data-slot")
    if (slot) s += `[data-slot=${slot}]`
    if (el.id && !/_R_|«|:r/.test(el.id)) s += `#${el.id}`
    return s
  }
  const keepHead = (el) => {
    const tag = el.tagName.toLowerCase()
    if (tag === "title") return true
    if (tag === "meta") return el.getAttribute("name") !== "next-size-adjust"
    if (tag === "link") return HEAD_LINK.has((el.getAttribute("rel") ?? "").toLowerCase())
    return false
  }
  const walk = (node, depth, trail, inHead, verbatim, doc) => {
    for (const child of node.childNodes) {
      if (child.nodeType === 3) {
        const t = verbatim ? child.textContent : child.textContent.replace(/\s+/g, " ").trim()
        if (t) {
          lines.push(`${"  ".repeat(depth)}${JSON.stringify(t)}`)
          paths.push(trail)
        }
        continue
      }
      if (child.nodeType !== 1) continue
      const tag = child.tagName.toLowerCase()
      if (tag === "script" || tag === "noscript" || tag === "template" || tag === "next-route-announcer") continue
      if (inHead && !keepHead(child)) continue
      if (tag === "link" && DROP_LINK.test(child.getAttribute("rel") ?? "")) continue
      if (tag === "div" && child.hidden && !child.children.length && !child.textContent.trim()) continue
      const attrs = []
      for (const a of child.attributes) {
        if (a.name === "nonce" || a.name.startsWith("data-parity")) continue
        let v = a.value
        if (a.name === "class") v = v.split(/\s+/).filter(Boolean).map(fontClass).sort().join(" ")
        else if (a.name === "style") v = child.style.cssText.split(";").map((d) => d.trim()).filter(Boolean).sort().join("; ")
        v = url(genIds(v))
        attrs.push(`${a.name}=${JSON.stringify(v)}`)
      }
      attrs.sort()
      if (tag === "input" || tag === "textarea" || tag === "select") {
        attrs.push(`.value=${JSON.stringify(child.value)}`)
        if (child.type === "checkbox" || child.type === "radio") attrs.push(`.checked=${child.checked}`)
      }
      if (doc.activeElement === child && child !== doc.body) attrs.push("(focused)")
      const here = `${trail} > ${label(child)}`
      lines.push(`${"  ".repeat(depth)}<${tag}${attrs.length ? " " + attrs.join(" ") : ""}>`)
      paths.push(here)
      if (tag === "iframe") {
        let sub = null
        try { sub = child.contentDocument } catch {}
        if (sub?.documentElement) {
          lines.push(`${"  ".repeat(depth + 1)}#document ${JSON.stringify(url(sub.location.pathname + sub.location.search))}`)
          paths.push(here)
          walkRoot(sub, depth + 1, here)
        }
        continue
      }
      walk(child.shadowRoot ?? child, depth + 1, here, inHead || tag === "head", verbatim || tag === "pre" || tag === "textarea", doc)
    }
  }
  const walkRoot = (doc, depth, trail) => {
    const html = doc.documentElement
    const attrs = [...html.attributes].map((a) => {
      let v = a.value
      if (a.name === "class") v = v.split(/\s+/).filter(Boolean).map(fontClass).sort().join(" ")
      else if (a.name === "style") v = html.style.cssText.split(";").map((d) => d.trim()).filter(Boolean).sort().join("; ")
      return `${a.name}=${JSON.stringify(url(genIds(v)))}`
    }).sort()
    lines.push(`${"  ".repeat(depth)}<html ${attrs.join(" ")}>`)
    paths.push(`${trail} > html`)
    walk(html, depth + 1, `${trail} > html`, false, false, doc)
  }
  walkRoot(document, 0, "")
  return { lines, paths }
}

// Describes the running animations of this document.
function describeAnimations() {
  const pathOf = (el) => {
    const out = []
    for (let e = el; e && e.nodeType === 1 && out.length < 4; e = e.parentElement) {
      let s = e.tagName.toLowerCase()
      const slot = e.getAttribute("data-slot")
      if (slot) s += `[data-slot=${slot}]`
      out.unshift(s)
    }
    return out.join(" > ")
  }
  return document.getAnimations().map((a) => {
    const t = a.effect?.getTiming?.() ?? {}
    const target = a.effect?.target
    const scroll = a.timeline && a.timeline !== document.timeline
    return {
      type: a.constructor.name,
      name: a.animationName ?? a.transitionProperty ?? "",
      target: target ? pathOf(target) + (a.effect.pseudoElement ?? "") : "",
      timeline: scroll ? a.timeline.constructor.name : "document",
      progress: scroll ? String(a.currentTime) : undefined,
      duration: scroll ? undefined : t.duration,
      delay: t.delay,
      easing: t.easing,
      iterations: t.iterations,
      direction: t.direction,
      fill: t.fill,
      keyframes: a.effect?.getKeyframes?.().map(({ computedOffset, ...k }) => k) ?? [],
    }
  }).sort((x, y) => JSON.stringify(x).localeCompare(JSON.stringify(y)))
}

function settleAnimations() {
  const doc = (a) => !a.timeline || a.timeline === document.timeline
  for (const a of document.getAnimations()) {
    if (!doc(a)) continue
    const end = a.effect?.getComputedTiming?.().endTime
    if (Number.isFinite(end)) a.finish()
    else {
      a.pause()
      a.currentTime = 0
    }
  }
}

// Animations on the document timeline that are not settled yet.
function pendingAnimations() {
  let n = 0
  for (const a of document.getAnimations()) {
    if (a.timeline && a.timeline !== document.timeline) continue
    const end = a.effect?.getComputedTiming?.().endTime
    if (Number.isFinite(end) ? a.playState !== "finished" : a.playState !== "paused" || a.currentTime !== 0) n++
  }
  return n
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

async function forceRepaint() {
  const raf2 = () => new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r)))
  const b = document.body
  if (!b) return
  const hadStyle = b.hasAttribute("style")
  const prev = b.style.willChange
  b.style.willChange = "transform"
  await raf2()
  b.style.willChange = prev
  if (!hadStyle && !b.getAttribute("style")) b.removeAttribute("style")
  await raf2()
}

function settleLoadTransitions() {
  const running = document.getAnimations().filter((a) => a.constructor.name === "CSSTransition")
  return Promise.race([Promise.all(running.map((a) => a.finished.catch(() => {}))), new Promise((r) => setTimeout(r, 1000))])
}

// D4: load lazy images and iframes now; restore the attribute afterwards.
function eagerLoad() {
  const lazy = [...document.querySelectorAll('img[loading="lazy"], iframe[loading="lazy"]')]
  for (const el of lazy) {
    el.setAttribute("data-parity-lazy", "")
    el.loading = "eager"
  }
  return lazy.length
}

async function restoreLazyAndDecode() {
  const imgs = [...document.images]
  await Promise.race([
    Promise.all(imgs.map((i) => (i.complete ? null : new Promise((r) => { i.addEventListener("load", r, { once: true }); i.addEventListener("error", r, { once: true }) })))),
    new Promise((r) => setTimeout(r, 10000)),
  ])
  await Promise.all(imgs.map((i) => (i.complete && i.naturalWidth ? i.decode().catch(() => {}) : null)))
  for (const el of document.querySelectorAll("[data-parity-lazy]")) {
    el.setAttribute("loading", "lazy")
    el.removeAttribute("data-parity-lazy")
  }
  await document.fonts.ready
}

// D6: resolves after `quiet` ms without DOM mutations (or `max` ms).
function domQuiet({ quiet, max }) {
  return new Promise((resolve) => {
    let timer
    const start = performance.now()
    const done = () => { obs.disconnect(); clearTimeout(timer); clearTimeout(cap); resolve(Math.round(performance.now() - start)) }
    const obs = new MutationObserver(() => { clearTimeout(timer); timer = setTimeout(done, quiet) })
    obs.observe(document, { subtree: true, childList: true, attributes: true, characterData: true })
    timer = setTimeout(done, quiet)
    const cap = setTimeout(done, max)
  })
}

// Deepest element covering each point (document coordinates), for reports.
function locate(points) {
  const all = [...document.body.querySelectorAll("*")]
  const pathOf = (el) => {
    const out = []
    for (let e = el; e && e !== document.body && out.length < 7; e = e.parentElement) {
      let s = e.tagName.toLowerCase()
      const slot = e.getAttribute("data-slot")
      if (slot) s += `[data-slot=${slot}]`
      else if (e.id && !/_R_/.test(e.id)) s += `#${e.id}`
      else if (typeof e.className === "string" && e.className) s += "." + e.className.trim().split(/\s+/).slice(0, 2).join(".")
      out.unshift(s)
    }
    return out.join(" > ")
  }
  return points.map(({ x, y, vx, vy }) => {
    let best = null
    let area = Infinity
    for (const el of all) {
      const r = el.getBoundingClientRect()
      const left = vx ? r.left : r.left + scrollX
      const top = vy ? r.top : r.top + scrollY
      if (x >= left && x < left + r.width && y >= top && y < top + r.height && r.width * r.height < area) {
        const cs = getComputedStyle(el)
        if (cs.visibility === "hidden" || cs.display === "none") continue
        best = el
        area = r.width * r.height
      }
    }
    if (!best) return "(none)"
    const text = (best.textContent ?? "").replace(/\s+/g, " ").trim().slice(0, 40)
    return pathOf(best) + (text ? ` "${text}"` : "")
  })
}

// --- one side ---------------------------------------------------------------------------

function chromiumPath() {
  if (process.env.PARITY_CHROMIUM) return process.env.PARITY_CHROMIUM
  for (const dir of [path.join(os.homedir(), "Library/Caches/ms-playwright"), path.join(os.homedir(), ".cache/ms-playwright")]) {
    if (!fs.existsSync(dir)) continue
    const shells = fs.readdirSync(dir).filter((d) => d.startsWith("chromium_headless_shell-")).sort().reverse()
    for (const s of shells) {
      for (const sub of fs.readdirSync(path.join(dir, s))) {
        const p = path.join(dir, s, sub, "chrome-headless-shell")
        if (fs.existsSync(p)) return p
      }
    }
  }
  return undefined
}

class Side {
  constructor(browser, base, role, theme, vpName) {
    Object.assign(this, { browser, base, role, theme, vpName })
    this.origin = new URL(base).origin
  }

  async open(pathname) {
    this.context = await this.browser.newContext({
      ...VIEWPORTS[this.vpName],
      deviceScaleFactor: 1,
      colorScheme: this.theme,
      reducedMotion: "no-preference",
      timezoneId: "UTC",
      locale: "en-US",
      extraHTTPHeaders: { "x-parity-now": PARITY_NOW },
    })
    await this.context.addInitScript(PRELUDE)
    this.page = await this.context.newPage()
    await this.page.clock.setFixedTime(new Date(PARITY_NOW))
    this.errors = []
    this.page.on("pageerror", (e) => this.errors.push(String(e)))
    this.page.on("console", (m) => m.type() === "error" && this.errors.push(m.text()))
    // In-flight requests, for a network-idle wait shorter than Playwright's
    // fixed 500ms "networkidle".
    this.inflight = new Set()
    this.lastNet = Date.now()
    const done = (r) => { this.inflight.delete(r); this.lastNet = Date.now() }
    this.page.on("request", (r) => { this.inflight.add(r); this.lastNet = Date.now() })
    this.page.on("requestfinished", done)
    this.page.on("requestfailed", done)
    this.cdp = await this.context.newCDPSession(this.page)
    await this.cdp.send("Animation.enable")
    const res = await this.page.goto(this.base + pathname, { waitUntil: "load", timeout: 60000 })
    this.status = res?.status() ?? 0
    await this.settle({ load: true })
    // D7: the document timeline stays frozen from here on; captures finish
    // or park animations explicitly, so no animation progresses with wall
    // time (a transition interrupted after a few real milliseconds would
    // otherwise be reversed from a timing-dependent point).
    await this.freeze(true)
  }

  // Resolves once no request has been in flight for `quiet` ms (max 20s).
  async netIdle(quiet = 150) {
    const deadline = Date.now() + 20000
    while (Date.now() < deadline) {
      if (this.inflight.size === 0 && Date.now() - this.lastNet >= quiet) return
      await new Promise((r) => setTimeout(r, 25))
    }
  }

  frames() {
    return this.page.frames().filter((f) => !f.isDetached())
  }

  async eachFrame(fn, arg) {
    return Promise.all(this.frames().map((f) => f.evaluate(fn, arg).catch(() => null)))
  }

  async settle({ load = false } = {}) {
    const p = this.page
    await this.netIdle()
    for (let round = 0; round < 3; round++) {
      const n = this.frames().length
      const lazy = (await this.eachFrame(eagerLoad)).reduce((a, b) => a + (b ?? 0), 0)
      if (lazy) await this.netIdle()
      await Promise.all(this.frames().map((f) => f.waitForLoadState("load", { timeout: 20000 }).catch(() => {})))
      await this.eachFrame(restoreLazyAndDecode)
      if (load) await this.eachFrame(settleLoadTransitions)
      await this.eachFrame(domQuiet, { quiet: load ? QUIET_LOAD : QUIET, max: 10000 })
      await this.netIdle()
      if (this.frames().length === n && !lazy) break
    }
    await this.idle()
  }

  async idle() {
    await this.eachFrame(() => new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r))))
  }

  async freeze(on) {
    await this.cdp.send("Animation.setPlaybackRate", { playbackRate: on ? 0 : 1 })
  }

  locator(s) {
    return this.page.locator(s.target).filter({ visible: true }).nth(s.nth ?? 0)
  }

  async step(s) {
    const p = this.page
    const target = s.target ? this.locator(s) : null
    const force = s.force ? { force: true } : {}
    const before = p.url()
    const t = { timeout: 10000 }
    switch (s.do) {
      case "hover": await target.hover({ ...force, ...t }); break
      case "click": await target.click({ ...force, ...t }); break
      case "dblclick": await target.dblclick({ ...force, ...t }); break
      case "focus": await target.focus(t); break
      case "press": await (target ? target.press(s.key, t) : p.keyboard.press(s.key)); break
      case "type": await (target ? target.pressSequentially(s.text, t) : p.keyboard.type(s.text)); break
      case "mouse": await p.mouse.move(s.x, s.y, { steps: s.steps ?? 1 }); break
      case "wait": await p.waitForTimeout(s.ms ?? 100); break
      case "scroll": await (target ? target.evaluate((el, y) => (el.scrollTop = y), s.y ?? 0) : p.mouse.wheel(0, s.y ?? 0)); break
      case "scrollTo":
        await p.evaluate(({ y, fraction }) => {
          const max = document.documentElement.scrollHeight - innerHeight
          window.scrollTo({ top: fraction != null ? Math.round(fraction * max) : y, behavior: "instant" })
        }, s)
        break
      case "settle": await this.settle(); break
      case "clipboard": this.clipboard = await p.evaluate(() => window.__parityClipboard); break
      case "eval": await p.evaluate(s.js); break
      default: throw new Error(`unknown step ${JSON.stringify(s)}`)
    }
    if (s.nav) {
      await p.waitForURL((u) => u.toString() !== before, { timeout: 15000 })
      await this.settle()
      await this.freeze(true) // a full navigation starts a new document timeline
    }
  }

  async visible(selector) {
    return (await this.page.locator(selector).filter({ visible: true }).count()) > 0
  }

  async state() {
    const p = this.page
    const s = await p.evaluate(() => ({ path: location.pathname + location.search + location.hash, scroll: [Math.round(scrollX), Math.round(scrollY)], title: document.title }))
    if (this.clipboard !== undefined) s.clipboard = this.clipboard
    return s
  }

  async shot(full) {
    await this.page.evaluate(forceRepaint)
    return this.page.screenshot({ fullPage: full, animations: "allow", caret: "hide", timeout: 60000 })
  }

  async snapshot(origins) {
    return this.page.evaluate(snapshotDOM, origins)
  }

  async close() {
    await this.context?.close().catch(() => {})
  }
}

// --- comparison ---------------------------------------------------------------------------

function comparePNG(a, b) {
  const A = PNG.sync.read(a)
  const B = PNG.sync.read(b)
  if (A.width !== B.width || A.height !== B.height) {
    return { pixels: Math.max(A.width * A.height, B.width * B.height), size: `${A.width}x${A.height} vs ${B.width}x${B.height}`, diff: null, regions: [] }
  }
  const diff = new PNG({ width: A.width, height: A.height })
  const pixels = pixelmatch(A.data, B.data, diff.data, A.width, A.height, { threshold: 0, includeAA: true })
  return { pixels, diff: pixels ? PNG.sync.write(diff) : null, regions: pixels ? regionsOf(A, B) : [] }
}

// Clusters differing pixels on a 32px grid (8-connected cells), largest first.
function regionsOf(A, B) {
  const C = 32
  const cw = Math.ceil(A.width / C)
  const cells = new Map()
  const a = A.data
  const b = B.data
  for (let y = 0; y < A.height; y++) {
    for (let x = 0; x < A.width; x++) {
      const i = (y * A.width + x) * 4
      if (a[i] !== b[i] || a[i + 1] !== b[i + 1] || a[i + 2] !== b[i + 2] || a[i + 3] !== b[i + 3]) {
        const k = Math.floor(y / C) * cw + Math.floor(x / C)
        const c = cells.get(k)
        if (c) {
          c.n++
          c.x0 = Math.min(c.x0, x); c.x1 = Math.max(c.x1, x); c.y0 = Math.min(c.y0, y); c.y1 = Math.max(c.y1, y)
        } else cells.set(k, { n: 1, x0: x, x1: x, y0: y, y1: y })
      }
    }
  }
  const seen = new Set()
  const regions = []
  for (const k of cells.keys()) {
    if (seen.has(k)) continue
    const r = { pixels: 0, x0: Infinity, y0: Infinity, x1: -1, y1: -1 }
    const stack = [k]
    seen.add(k)
    while (stack.length) {
      const c = stack.pop()
      const v = cells.get(c)
      r.pixels += v.n
      r.x0 = Math.min(r.x0, v.x0); r.x1 = Math.max(r.x1, v.x1); r.y0 = Math.min(r.y0, v.y0); r.y1 = Math.max(r.y1, v.y1)
      const cx = c % cw
      for (const dy of [-1, 0, 1]) for (const dx of [-1, 0, 1]) {
        if (cx + dx < 0 || cx + dx >= cw) continue
        const n = c + dy * cw + dx
        if (cells.has(n) && !seen.has(n)) { seen.add(n); stack.push(n) }
      }
    }
    regions.push(r)
  }
  return regions.sort((x, y) => y.pixels - x.pixels).slice(0, 8).map((r) => ({ pixels: r.pixels, box: [r.x0, r.y0, r.x1 - r.x0 + 1, r.y1 - r.y0 + 1] }))
}

// Line diff (Myers-free greedy walk with common prefix/suffix trimming).
function diffLines(A, B, max = 40) {
  const a = A.lines
  const b = B.lines
  let s = 0
  while (s < a.length && s < b.length && a[s] === b[s]) s++
  let ea = a.length
  let eb = b.length
  while (ea > s && eb > s && a[ea - 1] === b[eb - 1]) ea--, eb--
  const out = []
  let i = s
  let j = s
  const at = (P, k) => (P[k] ? `   @ ${P[k]}` : "")
  while ((i < ea || j < eb) && out.length < max) {
    if (i < ea && j < eb && a[i] === b[j]) {
      i++, j++
      continue
    }
    let k = -1
    for (let q = j; q < Math.min(eb, j + 400); q++) if (a[i] === b[q]) { k = q; break }
    let l = -1
    for (let q = i; q < Math.min(ea, i + 400); q++) if (b[j] === a[q]) { l = q; break }
    if (k >= 0 && (l < 0 || k - j <= l - i)) {
      while (j < k && out.length < max) { out.push(`+ target: ${b[j].trim()}${at(B.paths, j)}`); j++ }
    } else if (l >= 0) {
      while (i < l && out.length < max) { out.push(`- ref:    ${a[i].trim()}${at(A.paths, i)}`); i++ }
    } else {
      if (i < ea) { out.push(`- ref:    ${a[i].trim()}${at(A.paths, i)}`); i++ }
      if (j < eb) { out.push(`+ target: ${b[j].trim()}${at(B.paths, j)}`); j++ }
    }
  }
  return out
}

function animDiff(ra, ga) {
  const out = []
  const key = (x) => `${x.type} ${x.name} @ ${x.target}`
  const n = Math.max(ra.length, ga.length)
  for (let i = 0; i < n && out.length < 20; i++) {
    const r = ra[i]
    const g = ga[i]
    if (JSON.stringify(r) === JSON.stringify(g)) continue
    if (!r || !g) { out.push(`${r ? "- ref" : "+ target"}: ${key(r ?? g)}`); continue }
    const fields = Object.keys({ ...r, ...g }).filter((f) => JSON.stringify(r[f]) !== JSON.stringify(g[f]))
    out.push(`${key(r)}${key(r) !== key(g) ? ` vs ${key(g)}` : ""}: ` + fields.map((f) => `${f} ${JSON.stringify(r[f])} vs ${JSON.stringify(g[f])}`).join("; ").slice(0, 600))
  }
  return out
}

// --- one job: page x theme x viewport ---------------------------------------------------------

const slug = (p) => (p === "/" ? "_root" : p.replace(/^\//, "").replace(/\//g, "__"))

function pageScenarios(page) {
  const f = path.join(here, "scenarios", slug(page.path) + ".json")
  return fs.existsSync(f) ? JSON.parse(fs.readFileSync(f, "utf8")) : []
}

async function runJob(browser, page, theme, vp) {
  const dir = path.join(outRoot, slug(page.path), `${theme}-${vp}`)
  fs.rmSync(dir, { recursive: true, force: true })
  fs.mkdirSync(dir, { recursive: true })
  const report = { path: page.path, group: page.group, theme, viewport: vp, ref: REF, target: TARGET, pass: true, captures: [], errors: {} }
  const ref = new Side(browser, REF, "ref", theme, vp)
  const tgt = new Side(browser, TARGET, "target", theme, vp)
  const origins = [ref.origin, tgt.origin]
  const both = (fn) => Promise.all([fn(ref), fn(tgt)])

  const save = (name, data) => fs.writeFileSync(path.join(dir, name), data)

  const record = async (label, { frames = false, full = false } = {}) => {
    const cap = { label, pass: true }
    const tStart = Date.now()
    if (frames) {
      const [ra, ga] = await both(async (s) => (await s.eachFrame(describeAnimations)).flatMap((x, i) => (x ?? []).map((a) => ({ frame: i, ...a }))))
      cap.animations = { equal: JSON.stringify(ra) === JSON.stringify(ga), count: [ra.length, ga.length] }
      if (!cap.animations.equal) {
        cap.pass = false
        cap.animations.diff = animDiff(ra, ga)
        save(`${label}.animations.json`, JSON.stringify({ ref: ra, target: ga }, null, 2))
      }
      const ends = await both(async (s) => Math.max(0, ...(await s.eachFrame(maxAnimationEnd)).map((x) => x ?? 0)))
      const end = Math.max(...ends)
      const times = end > 0 ? [0, 0.25, 0.5, 0.75].map((f) => Math.round(f * end)) : []
      cap.frames = []
      for (const t of times) {
        let r, g, cmp
        for (let attempt = 0; attempt < 3; attempt++) {
          await both((s) => s.eachFrame(seekAnimations, t))
          await both((s) => s.idle())
          ;[r, g] = await both((s) => s.page.screenshot({ animations: "allow", caret: "hide" }))
          cmp = comparePNG(r, g)
          if (cmp.pixels === 0) break
        }
        const fr = { t, pixels: cmp.pixels }
        if (cmp.pixels) {
          fr.regions = cmp.regions
          fr.at = await ref.page.evaluate(locate, cmp.regions.map((x) => ({ x: x.box[0] + x.box[2] / 2, y: x.box[1] + x.box[3] / 2, vx: true, vy: true }))).catch(() => [])
          cap.pass = false
        }
        cap.frames.push(fr)
        if (cmp.pixels || keep) {
          save(`${label}@${t}ms-ref.png`, r)
          save(`${label}@${t}ms-target.png`, g)
          if (cmp.diff) save(`${label}@${t}ms-diff.png`, cmp.diff)
        }
      }
    }
    // Settle until nothing is left running: finishing an animation can start
    // another (exit animations unmounting, transitions on the next state).
    for (let i = 0; i < 6; i++) {
      await both((s) => s.eachFrame(settleAnimations))
      await both((s) => s.idle())
      const pending = await both(async (s) => (await s.eachFrame(pendingAnimations)).reduce((a, b) => a + (b ?? 0), 0))
      if (!pending[0] && !pending[1]) break
    }
    let r, g, cmp, rd, gd
    for (let attempt = 0; attempt < 3; attempt++) {
      ;[rd, gd] = await both((s) => s.snapshot(origins))
      ;[r, g] = await both((s) => s.shot(full))
      cmp = comparePNG(r, g)
      if (cmp.pixels === 0) break
    }
    cap.pixels = cmp.pixels
    if (cmp.size) cap.size = cmp.size
    if (cmp.pixels) {
      cap.regions = cmp.regions
      const pts = cmp.regions.map((x) => ({ x: x.box[0] + x.box[2] / 2, y: x.box[1] + x.box[3] / 2, vx: !full, vy: !full }))
      cap.at = await ref.page.evaluate(locate, pts).catch(() => [])
    }
    cap.dom = diffLines(rd, gd)
    const [rs, gs] = await both((s) => s.state())
    if (JSON.stringify(rs) !== JSON.stringify(gs)) cap.state = { ref: rs, target: gs }
    else cap.stateEqual = rs
    if (cmp.pixels || cap.dom.length || cap.state) cap.pass = false
    if (!cap.pass || keep) {
      save(`${label}-ref.png`, r)
      save(`${label}-target.png`, g)
      if (cmp.diff) save(`${label}-diff.png`, cmp.diff)
      save(`${label}-ref.dom.txt`, rd.lines.join("\n"))
      save(`${label}-target.dom.txt`, gd.lines.join("\n"))
    }
    cap.ms = Date.now() - tStart
    report.captures.push(cap)
    if (!cap.pass) report.pass = false
    return cap
  }

  try {
    const tOpen = Date.now()
    await both((s) => s.open(page.path))
    report.openMs = Date.now() - tOpen
    if (ref.status !== tgt.status) report.status = { ref: ref.status, target: tgt.status }
    await record("initial", { full: true })

    if (withScenarios) {
      const list = [...chromeScenarios, ...pageScenarios(page)].filter((sc) =>
        (!sc.viewports || sc.viewports.includes(vp)) && (!sc.groups || sc.groups.includes(page.group)) && (!only || only.has(sc.label)))
      for (const sc of list) {
        if (sc.reset) {
          await both((s) => s.close())
          await both((s) => s.open(page.path))
        }
        if (sc.when && !(await ref.visible(sc.when))) continue
        // A scenario whose first target the reference cannot act on (e.g.
        // covered by a modal the page opens) is skipped. The trial runs on
        // both sides, so any scrolling it does stays symmetric.
        const first = sc.steps[0]
        if (first?.target && ["click", "hover", "dblclick"].includes(first.do) && !first.force) {
          const trial = await Promise.allSettled([ref, tgt].map((s) => s.locator(first).click({ trial: true, timeout: 2000 })))
          if (trial[0].status === "rejected") {
            report.skipped = [...(report.skipped ?? []), `${sc.label}: ${String(trial[0].reason?.message ?? "").split("\n").find((l) => l.includes("intercepts")) ?? "not actionable"}`.slice(0, 300)]
            continue
          }
        }
        const frames = withFrames && sc.frames !== false
        const tSteps = Date.now()
        try {
          for (const st of sc.steps) {
            const res = await Promise.allSettled([ref.step(st), tgt.step(st)])
            const bad = res.map((x, i) => (x.status === "rejected" ? `${i ? "target" : "ref"}: ${String(x.reason?.message ?? x.reason).split("\n")[0]}` : null)).filter(Boolean)
            if (bad.length) throw new Error(`step ${JSON.stringify(st)} failed: ${bad.join(" | ")}`)
          }
          await both((s) => s.idle())
          const stepsMs = Date.now() - tSteps
          ;(await record(sc.label, { frames, full: !!sc.full })).stepsMs = stepsMs
        } catch (e) {
          report.captures.push({ label: sc.label, pass: false, error: String(e.message ?? e) })
          report.pass = false
          // Continue from a clean state.
          await both((s) => s.close())
          await both((s) => s.open(page.path))
        } finally {
          ref.clipboard = tgt.clipboard = undefined
        }
      }
    }
  } catch (e) {
    report.errors.run = String(e.stack ?? e)
    report.pass = false
  } finally {
    report.errors.ref = ref.errors ?? []
    report.errors.target = tgt.errors ?? []
    if (report.status) report.pass = false
    if (!SELF && tgt.errors?.some((e) => !ref.errors.includes(e))) report.pass = false
    await both((s) => s.close())
  }
  save("report.json", JSON.stringify(report, null, 2))
  return report
}

// --- main ---------------------------------------------------------------------------------------

function selectPages() {
  let pages
  if (positional.length) {
    const known = new Map(allPages().map((p) => [p.path, p]))
    pages = positional.map((p) => known.get(p) ?? { path: p, group: p.startsWith("/docs") ? "docs" : "app", source: "argument" })
  } else if (flag("--all")) pages = allPages()
  else if (opt("--group")) pages = allPages(opt("--group").split(","))
  else throw new Error("no pages selected: give paths, --group g or --all")
  if (opt("--grep")) {
    const re = new RegExp(opt("--grep"))
    pages = pages.filter((p) => re.test(p.path))
  }
  if (opt("--shard")) {
    const [i, n] = opt("--shard").split("/").map(Number)
    pages = pages.filter((_, k) => k % n === i - 1)
  }
  return pages
}

function summaryLine(r) {
  const bad = r.captures.filter((c) => !c.pass)
  const parts = bad.map((c) => {
    if (c.error) return `${c.label}: ${c.error.slice(0, 200)}`
    const what = []
    if (c.pixels) what.push(`${c.pixels}px${c.at?.length ? ` at ${c.at[0]}` : ""}`)
    if (c.dom?.length) what.push(`dom ${c.dom[0].slice(0, 200)}`)
    if (c.state) what.push(`state ${JSON.stringify(c.state).slice(0, 200)}`)
    if (c.animations && !c.animations.equal) what.push(`animations ${c.animations.diff?.[0] ?? ""}`)
    if (c.frames?.some((f) => f.pixels)) what.push(`frames ${c.frames.filter((f) => f.pixels).map((f) => `${f.t}ms:${f.pixels}px`).join(",")}`)
    return `${c.label}: ${what.join("; ")}`
  })
  if (r.status) parts.unshift(`status ${r.status.ref} vs ${r.status.target}`)
  if (r.errors.run) parts.unshift(`error ${r.errors.run.split("\n")[0]}`)
  return `${r.pass ? "PASS" : "FAIL"} ${r.path} ${r.theme}/${r.viewport} (${r.captures.length} captures)${parts.length ? "\n     " + parts.join("\n     ") : ""}`
}

async function main() {
  if (!TARGET) throw new Error("--target URL is required (or --self)")
  const pages = selectPages()
  const jobs = []
  for (const p of pages) for (const vp of viewports) for (const theme of themes) jobs.push({ page: p, theme, vp })
  console.log(`${pages.length} pages x ${themes.length} themes x ${viewports.length} viewports = ${jobs.length} jobs; ref ${REF} vs target ${TARGET}${SELF ? " (self)" : ""}`)
  fs.mkdirSync(outRoot, { recursive: true })
  const browser = await chromium.launch({ executablePath: chromiumPath() })
  const reports = []
  const queue = [...jobs]
  const t0 = Date.now()
  await Promise.all(Array.from({ length: Math.min(jobsN, queue.length) }, async () => {
    while (queue.length) {
      const j = queue.shift()
      let r
      try {
        r = await runJob(browser, j.page, j.theme, j.vp)
      } catch (e) {
        r = { path: j.page.path, theme: j.theme, viewport: j.vp, pass: false, captures: [], errors: { run: String(e.stack ?? e) } }
      }
      reports.push(r)
      console.log(`[${reports.length}/${jobs.length}] ${summaryLine(r)}`)
    }
  }))
  await browser.close()
  const failed = reports.filter((r) => !r.pass)
  const captures = reports.reduce((n, r) => n + r.captures.length, 0)
  const summary = {
    ref: REF, target: TARGET, self: SELF, jobs: reports.length, captures, failed: failed.length, seconds: Math.round((Date.now() - t0) / 1000),
    failures: failed.map((r) => ({ path: r.path, theme: r.theme, viewport: r.viewport, captures: r.captures.filter((c) => !c.pass).map((c) => c.label), error: r.errors?.run?.split("\n")[0] })),
  }
  fs.writeFileSync(path.join(outRoot, "summary.json"), JSON.stringify(summary, null, 2))
  console.log(`\n${reports.length - failed.length}/${reports.length} jobs identical (${captures} captures, ${summary.seconds}s). Reports: ${path.relative(process.cwd(), outRoot) || "."}/`)
  return failed.length ? 1 : 0
}

if (import.meta.url === `file://${process.argv[1]}`) {
  main().then((code) => process.exit(code), (e) => { console.error(e.stack ?? e); process.exit(2) })
}
