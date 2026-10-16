// Page-level parity harness for the documentation site (DOCS_PLAN.md, D1).
//
// Compares two servers page by page: the reference (the upstream Next.js
// docs site, parity/docs/reference.sh, default http://localhost:4000) and a
// candidate (--target, the Go docs server). The default bar is visual
// closeness (--tolerance, 1% of content pixels); --exact is strict equality.
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
//   --tolerance F         fraction of content pixels allowed to differ (default 0.01)
//   --exact               strict bar: 0 pixels, identical DOM, animations and frames
//   --keep                keep images and DOM dumps of passing captures too
//   --out DIR             reports (default parity/docs/out)
//   --rerun FILE          only the failed jobs listed in a previous summary.json
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
//               DevTools protocol (as in parity/run.mjs, see D7); the
//               running animations' definitions (keyframes, duration,
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
//   N5 Generated ids: React useId values (_R_..._, _r_..._, «r..», :r..:),
//      recharts' lodash uniqueId ids (recharts<N>-clip, recharts-radial-line-<N>),
//      also inside <style> text (ChartStyle's [data-chart=chart-_r_7_]), and the ids
//      of shadcn.go's uid package (radix-|base-ui-|aria-|react-aria- + 10
//      base32 characters) are replaced in every attribute by #idN in order
//      of first appearance, keeping any prefix and suffix
//      ("radix-_R_43a_-trigger-cli" -> "radix-#id3-trigger-cli"), numbered
//      separately in each frame's document. Content
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
//   N9 Library scratch elements outside the visible document are dropped:
//      recharts' <span id="recharts_measurement_span"> (an off-screen text
//      measuring span whose content is the last string measured, which
//      depends on render order).
//
// Nondeterminism neutralized identically on both sides. Each item was found
// by the self-comparison (--self: the reference against itself in two
// independent browser contexts) and is applied to the reference and the
// candidate in the same way.
//
//   D1  Time: Date is fixed at PARITY_NOW (Date.now() and new Date() return
//       it; timers and performance.now keep running). Playwright's
//       page.clock is deliberately not used: it replaces setTimeout and
//       requestAnimationFrame with its own fakes, defeating D6 and D10. Server-rendered dates on the candidate must use the
//       same instant: every request carries "x-parity-now: <ISO time>" (stub
//       contract for the Go server; the reference is prerendered at build).
//   D2  Math.random is a seeded PRNG (mulberry32, seed 1), fresh per document.
//   D3  Clipboard: navigator.clipboard.writeText/readText and the "copy"
//       event are captured per page (window.__parityClipboard) instead of the
//       browser-wide system clipboard, which parallel jobs would share.
//   D4  Lazy content: <img loading=lazy> and <iframe loading=lazy> are loaded
//       eagerly (loading set to eager, the attribute restored once loaded),
//       then every image is decoded before a capture.
//   D5  Fonts: document.fonts.ready in every frame before a capture.
//   D6  Late effects: a capture waits for network idle (no request in flight
//       for 150ms) and, in every frame, for a window without DOM mutations
//       and without page calls to requestAnimationFrame of 200ms (600ms after
//       a page load) AND 8 animation frames (max 15s), so rAF-driven JS
//       animations and rAF-based delays (recharts) are waited for even on a
//       loaded machine. A change of the document's scroll size also resets
//       the window (content-visibility:auto content, such as the closed
//       CodeCollapsibleWrapper, gets its real height only when it nears the
//       viewport). Lazily added frames repeat the wait.
//   D7  Animations: CSS transitions started by the page load are awaited
//       (max 1s); then the document timeline stays frozen (DevTools
//       Animation.setPlaybackRate 0) for the rest of the job, so nothing
//       progresses with wall time (a transition interrupted a few real
//       milliseconds after it started would otherwise reverse from a
//       timing-dependent point). Before each capture, finite animations are
//       finished and infinite ones (spinners, shimmer, marquees) paused at 0,
//       in every frame, repeatedly until none is left running. Scroll-driven
//       animations are left to their scroll position. The frozen timeline
//       also stops requestAnimationFrame timestamps, so page callbacks get
//       performance.now() instead: JS animations (recharts) run in real time
//       and the settle wait (D6) lets them finish.
//   D8  Rendering: Chromium runs with --disable-partial-raster (a re-raster
//       never keeps anti-aliased edges from an earlier paint, which left
//       1-gray-level differences in iframes) and --force-color-profile=srgb;
//       a full repaint is forced in every frame before each screenshot (as
//       in parity/run.mjs); the caret is hidden; deviceScaleFactor 1; one
//       Chromium build (the headless shell parity/run.mjs uses); timezone
//       UTC, locale en-US; reducedMotion no-preference. A capture that
//       differs is settled and taken again up to twice (a real mismatch is
//       deterministic); the last one is reported.
//   D9  Scroll: every job starts at the top in a fresh context (no
//       sessionStorage, so the sidebar scroll-restore script sees no state);
//       scroll steps use instant scrolling; a step target outside the
//       viewport is centered instantly, the page settled (content above can
//       grow when it scrolls into view), and repeated until the target stays
//       in view, before Playwright acts on it (its own scroll-into-view
//       retries other alignments when the first attempt is obscured, which
//       depends on timing).
//   D10 Long timers: setTimeout with a delay of 1s or more is held (it never
//       fires by itself) and setInterval of 1s or more never fires, from the
//       first script on. This freezes autoplay carousels, toast auto-dismiss,
//       "Copied!" states that revert after 2s, and anything else on a
//       wall-clock schedule, at the same point on both sides. A scenario can
//       run the held timers explicitly with the "timers" step.
//   D11 Iframe load race: once an iframe has loaded, a "load" event is
//       dispatched on it once more (after hydration), so a listener attached
//       late sees it on both sides (upstream's ChartIframe otherwise stays at
//       opacity-0 whenever the frame loads before React hydrates).
//   D12 Theme switch transitions: next-themes disables transitions with a
//       temporary style removed after getComputedStyle(document.body), which
//       forces no style update, so whether colors then transition depends on
//       frame timing. The harness forces a style update before that style is
//       removed: theme changes never transition, on both sides.
//   D13 Full-page screenshots fire "resize" in every frame (Chromium's
//       captureBeyondViewport) although no size changes; resize events are
//       swallowed (capture-phase listener installed first) while the harness
//       takes one, so no handler re-renders during the capture.
//   D14 Hydration: each frame of a Next.js page (self.__next_f present) is
//       waited for until React has hydrated it (<html> carries its fiber),
//       before the DOM-quiet wait; under load hydration can come seconds
//       after the load event. A page declaring window.__shadcnReady (the Go
//       runtime) is waited for until it is true.
//   D15 Focus taken by an iframe during load (an autofocus inside a block
//       preview): which frame ends up focused depends on load order, so if
//       the top document's active element is an iframe after load, it is
//       blurred.
//   D16 Iframe size race: an iframe's first render races the parent's
//       layout (react-resizable-panels sizes the panel holding a block
//       preview on mount; the frame may hydrate at width 0, and its
//       useIsMobile-driven state then differs). After the page settled,
//       every iframe is reloaded once and the page settled again.
//
// Harness errors (a timeout on an overloaded machine) retry the job once;
// a difference is never retried. A scenario whose first target the
// reference cannot act on (covered by a modal the page opens, e.g. /create)
// is skipped and listed under "skipped" in the report.
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
//   timers                             run the held timers (D10) once, in every frame
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
const VALUED = ["--tolerance", "--rerun", "--quiet", "--ref", "--target", "--group", "--grep", "--shard", "--themes", "--viewports", "--only", "--jobs", "--out"]
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
// Default bar: visually close. A capture passes when at most --tolerance
// (default 1%) of its content pixels (pixels differing from the page
// background, the top-left pixel, in either image) differ, and the state
// matches. DOM differences are reported but do not fail; animation frames
// are not sampled. Screenshots of different sizes are compared on the
// larger size, the smaller one padded with its background. --exact restores
// the strict bar described above.
const EXACT = flag("--exact")
const TOLERANCE = EXACT ? 0 : Number(opt("--tolerance", "0.01"))
const tooDifferent = (cmp) => cmp.pixels > Math.floor(TOLERANCE * (cmp.total ?? Infinity))
const withFrames = EXACT && !flag("--no-frames")
const only = opt("--only") ? new Set(opt("--only").split(",")) : null
const outRoot = path.resolve(opt("--out", path.join(here, "out")))
const QUIET = Number(opt("--quiet", "200"))
const QUIET_LOAD = 600
export const LAUNCH_ARGS = ["--disable-partial-raster", "--force-color-profile=srgb"]
export const PARITY_NOW = "2026-01-15T12:00:00.000Z"

// --- browser-side helpers (serialized into the page) ---------------------------------------

// D2 + D3: installed before any page script, in every frame.
const PRELUDE = `(() => {
  // D1: a fixed Date. (Playwright's page.clock is not used: it replaces
  // setTimeout and requestAnimationFrame with its own fakes, over D6/D10.)
  const RealDate = Date
  const NOW = ${Date.parse(PARITY_NOW)}
  function FixedDate(...a) {
    if (!new.target) return new RealDate(NOW).toString()
    return a.length ? new RealDate(...a) : new RealDate(NOW)
  }
  FixedDate.prototype = RealDate.prototype
  FixedDate.now = () => NOW
  FixedDate.parse = RealDate.parse
  FixedDate.UTC = RealDate.UTC
  window.Date = FixedDate
  let s = 1
  Math.random = function () {
    s |= 0; s = (s + 0x6d2b79f5) | 0
    let t = Math.imul(s ^ (s >>> 15), 1 | s)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
  // D10: timers of 1s or more are held (never fire by themselves); the
  // "timers" step flushes them. Intervals of 1s or more never fire.
  const realSetTimeout = window.setTimeout.bind(window)
  window.__paritySetTimeout = realSetTimeout // for the harness's own waits
  // D6: page code's last requestAnimationFrame call, so JS animation loops
  // (recharts' rAF timeouts and tweens) count as activity for the settle wait.
  const realRAF = window.requestAnimationFrame.bind(window)
  window.__parityRAF = realRAF
  window.__parityRafAt = 0
  // The callback gets performance.now() as its timestamp: with the document
  // timeline frozen (D7) Chromium's rAF timestamps stop advancing, and a JS
  // animation computing progress from them (recharts) would never end.
  window.requestAnimationFrame = function (cb) {
    window.__parityRafAt = performance.now()
    return realRAF(() => cb(performance.now()))
  }
  const realClearTimeout = window.clearTimeout.bind(window)
  const realSetInterval = window.setInterval.bind(window)
  const realClearInterval = window.clearInterval.bind(window)
  const held = new Map()
  let nextHeld = 1e9
  window.setTimeout = function (fn, ms, ...rest) {
    if (Number(ms) >= 1000 && typeof fn === "function") {
      const id = nextHeld++
      held.set(id, () => fn(...rest))
      return id
    }
    return realSetTimeout(fn, ms, ...rest)
  }
  window.clearTimeout = function (id) { if (held.has(id)) held.delete(id); else realClearTimeout(id) }
  window.setInterval = function (fn, ms, ...rest) {
    if (Number(ms) >= 1000) { const id = nextHeld++; held.set(id, null); return id }
    return realSetInterval(fn, ms, ...rest)
  }
  window.clearInterval = function (id) { if (held.has(id)) held.delete(id); else realClearInterval(id) }
  window.__parityFlushTimers = () => {
    const due = [...held.entries()].filter(([, f]) => f)
    for (const [id] of due) held.delete(id)
    for (const [, f] of due) { try { f() } catch (e) { realSetTimeout(() => { throw e }) } }
    return due.length
  }
  // D12: next-themes disables transitions during a theme change with a
  // temporary "transition:none!important" <style>, removed 1ms later after
  // getComputedStyle(document.body), which does not force a style update:
  // whether transitions then run depends on frame timing. Force the update
  // before such a style is removed, so the change never transitions.
  const isNoTransition = (n) => n && n.nodeName === "STYLE" && /transition:none!important/.test(n.textContent || "")
  const realRemoveChild = Node.prototype.removeChild
  Node.prototype.removeChild = function (child) {
    if (isNoTransition(child)) void document.documentElement.offsetHeight
    return realRemoveChild.call(this, child)
  }
  const realRemove = Element.prototype.remove
  Element.prototype.remove = function () {
    if (isNoTransition(this)) void document.documentElement.offsetHeight
    return realRemove.call(this)
  }
  // D13: a full-page screenshot (captureBeyondViewport) fires "resize" in
  // every frame although no size changes; handlers then re-render or
  // re-initialize (carousels, charts) while the capture is taken. Resize
  // events are swallowed while the harness takes such a screenshot.
  window.__parityHoldResize = false
  const gate = (e) => { if (window.__parityHoldResize) e.stopImmediatePropagation() }
  window.addEventListener("resize", gate, true)
  try { window.visualViewport && window.visualViewport.addEventListener("resize", gate, true) } catch {}
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
  // One id map per document: ids never refer across frames, and each frame
  // has its own generators (recharts' counter restarts in every iframe).
  let ids = new Map()
  const idOf = (v) => {
    if (!ids.has(v)) ids.set(v, `#id${ids.size}`)
    return ids.get(v)
  }
  const REACT_ID = /_[Rr]_[0-9A-Za-z]+_|«[^»]+»|:r[0-9a-z]+:|\brecharts\d+(?=-)|\brecharts(?:-[a-z]+)+-\d+\b/g
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
    if (el.id && !/_[Rr]_|«|:r/.test(el.id)) s += `#${el.id}`
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
        let t = verbatim ? child.textContent : child.textContent.replace(/\s+/g, " ").trim()
        if (node.nodeName === "STYLE") t = genIds(t)
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
      if (child.id === "recharts_measurement_span") continue
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
    const outer = ids
    ids = new Map()
    try {
      walkDoc(doc, depth, trail)
    } finally {
      ids = outer
    }
  }
  const walkDoc = (doc, depth, trail) => {
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
  const rf = window.__parityRAF ?? requestAnimationFrame
  const raf2 = () => new Promise((r) => rf(() => rf(r)))
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
  const st = window.__paritySetTimeout ?? setTimeout
  const running = document.getAnimations().filter((a) => a.constructor.name === "CSSTransition")
  return Promise.race([Promise.all(running.map((a) => a.finished.catch(() => {}))), new Promise((r) => st(r, 1000))])
}

// D4: load lazy images and iframes now; restore the attribute afterwards.
function eagerLoad() {
  const done = (window.__parityEager ??= new WeakSet())
  let n = 0
  for (const el of document.querySelectorAll('img[loading="lazy"], iframe[loading="lazy"]')) {
    if (done.has(el)) continue
    done.add(el)
    el.setAttribute("data-parity-lazy", "")
    el.loading = "eager"
    n++
  }
  return n
}

async function restoreLazyAndDecode() {
  const imgs = [...document.images]
  await Promise.race([
    Promise.all(imgs.map((i) => (i.complete ? null : new Promise((r) => { i.addEventListener("load", r, { once: true }); i.addEventListener("error", r, { once: true }) })))),
    new Promise((r) => (window.__paritySetTimeout ?? setTimeout)(r, 10000)),
  ])
  await Promise.all(imgs.map((i) => (i.complete && i.naturalWidth ? i.decode().catch(() => {}) : null)))
  for (const el of document.querySelectorAll("[data-parity-lazy]")) {
    el.setAttribute("loading", "lazy")
    el.removeAttribute("data-parity-lazy")
  }
  await document.fonts.ready
}

// D14: resolves once the document is interactive: on a Next.js page (it
// carries self.__next_f) when React has hydrated (<html> has its fiber);
// on a page that declares window.__shadcnReady when that is true.
function waitHydrated(max) {
  const st = window.__paritySetTimeout ?? setTimeout
  const ok = () => {
    if (self.__next_f) return Object.keys(document.documentElement).some((k) => k.startsWith("__reactFiber$"))
    if ("__shadcnReady" in window) return !!window.__shadcnReady
    return document.readyState !== "loading"
  }
  const start = performance.now()
  return new Promise((resolve) => {
    const poll = () => (ok() || performance.now() - start > max ? resolve(ok()) : st(poll, 25))
    poll()
  })
}

// D15: when an iframe took the focus during load, which of several frames
// ends up focused depends on their load order: drop that focus.
function dropFrameFocus() {
  const a = document.activeElement
  if (a && a.tagName === "IFRAME") {
    try { a.contentDocument?.activeElement?.blur() } catch {}
    a.blur()
    return 1
  }
  return 0
}

// D16: reloads every iframe once the parent settled.
function reloadIframes() {
  let n = 0
  for (const f of document.querySelectorAll("iframe")) {
    try { f.contentWindow.location.reload(); n++ } catch {}
  }
  return n
}

// D11: re-dispatches "load" once on every iframe that finished loading.
// Pages that attach their onLoad listener during hydration (React) miss the
// event when the frame loaded first: a race the reference itself loses at
// random (e.g. ChartIframe stays at opacity-0).
function redispatchIframeLoad() {
  const seen = (window.__parityLoaded ??= new WeakSet())
  let n = 0
  for (const f of document.querySelectorAll("iframe")) {
    let ready = false
    try { ready = f.contentDocument?.readyState === "complete" && f.contentDocument.location.href !== "about:blank" } catch {}
    if (!ready || seen.has(f)) continue
    seen.add(f)
    f.dispatchEvent(new Event("load"))
    n++
  }
  return n
}

// D6: resolves after `quiet` ms without DOM mutations (or `max` ms).
// Counted both in milliseconds and in animation frames, so a rAF-driven JS
// animation (recharts) is seen as running even on a loaded machine.
function domQuiet({ quiet, frames, max }) {
  const st = window.__paritySetTimeout ?? setTimeout
  return new Promise((resolve) => {
    const start = performance.now()
    let last = start
    let ticks = 0
    let finished = false
    const raf = window.__parityRAF ?? requestAnimationFrame
    const obs = new MutationObserver(() => { last = performance.now(); ticks = 0 })
    obs.observe(document, { subtree: true, childList: true, attributes: true, characterData: true })
    let rafSeen = window.__parityRafAt ?? 0
    const check = () => {
      if (finished) return
      const now = performance.now()
      if ((window.__parityRafAt ?? 0) !== rafSeen) { rafSeen = window.__parityRafAt; last = Math.max(last, rafSeen); ticks = 0 }
      if ((now - last >= quiet && ticks >= frames) || now - start >= max) {
        finished = true
        obs.disconnect()
        resolve(now - start >= max ? -1 : Math.round(now - start))
      }
    }
    // Layout changes without DOM mutations count too (content-visibility:
    // auto content gets its real size when it nears the viewport).
    const de = document.documentElement
    let size = de.scrollHeight + "x" + de.scrollWidth
    const tick = () => {
      const now = de.scrollHeight + "x" + de.scrollWidth
      if (now !== size) { size = now; last = performance.now(); ticks = 0 }
      ticks++
      check()
      if (!finished) raf(tick)
    }
    raf(tick)
    const poll = () => { check(); if (!finished) st(poll, 50) }
    st(poll, 50)
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
      else if (e.id && !/_[Rr]_/.test(e.id)) s += `#${e.id}`
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

export function chromiumPath() {
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

export class Side {
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
    const res = await this.page.goto(this.base + pathname, { waitUntil: "domcontentloaded", timeout: 180000 })
    await this.page.waitForLoadState("load", { timeout: 120000 }).catch(() => this.timeout("load"))
    this.status = res?.status() ?? 0
    await this.settle({ load: true })
    // D16: an iframe's first render races the parent's layout (e.g. a
    // react-resizable-panels panel sized on mount: the frame hydrates at
    // width 0 or at its final width, and useIsMobile-driven state differs).
    // Reloading every iframe once the parent settled makes each frame start
    // at its final size.
    if (await this.page.evaluate(reloadIframes)) {
      await this.page.waitForTimeout(50)
      await this.settle({ load: true })
    }
    await this.page.evaluate(dropFrameFocus)
    // D7: the document timeline stays frozen from here on; captures finish
    // or park animations explicitly, so no animation progresses with wall
    // time (a transition interrupted after a few real milliseconds would
    // otherwise be reversed from a timing-dependent point).
    await this.freeze(true)
  }

  timeout(what) {
    ;(this.settleTimeouts ??= []).push(what)
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
    for (let round = 0; round < 6; round++) {
      const n = this.frames().length
      const lazy = (await this.eachFrame(eagerLoad)).reduce((a, b) => a + (b ?? 0), 0)
      if (lazy) await this.netIdle()
      await Promise.all(this.frames().map((f) => f.waitForLoadState("load", { timeout: 20000 }).catch(() => {})))
      const hydrated = await this.eachFrame(waitHydrated, 30000)
      if (hydrated.includes(false)) this.timeout(`hydration of ${this.frames()[hydrated.indexOf(false)]?.url()}`)
      await this.eachFrame(restoreLazyAndDecode)
      const reloaded = (await this.eachFrame(redispatchIframeLoad)).reduce((a, b) => a + (b ?? 0), 0)
      if (load) await this.eachFrame(settleLoadTransitions)
      const q = await this.eachFrame(domQuiet, { quiet: load ? QUIET_LOAD : QUIET, frames: 8, max: 15000 })
      if (q.includes(-1)) this.timeout(`dom-quiet of ${this.frames()[q.indexOf(-1)]?.url()}`)
      await this.netIdle()
      if (this.frames().length === n && !lazy && !reloaded) break
    }
    await this.idle()
  }

  async idle() {
    await this.eachFrame(() => new Promise((r) => { const rf = window.__parityRAF ?? requestAnimationFrame; rf(() => rf(r)) }))
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
    // Playwright scrolls a target into view itself, with alignment retries
    // that depend on timing (an obscured first attempt is retried with
    // another block alignment). Center a target that is not fully in the
    // viewport first, instantly, so both sides scroll to the same place.
    // Content can change size when it scrolls into view (lazily rendered
    // previews, IntersectionObserver effects), so settle after scrolling and
    // repeat until the target stays in view.
    if (target && ["hover", "click", "dblclick", "focus"].includes(s.do)) await this.bringIntoView(s)
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
      case "timers": await this.eachFrame(() => window.__parityFlushTimers?.()); break
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

  async bringIntoView(s) {
    const target = this.locator(s)
    const t = { timeout: 10000 }
    for (let i = 0; i < 4; i++) {
      const scrolled = await target.evaluate((el) => {
        const r = el.getBoundingClientRect()
        if (r.top >= 0 && r.left >= 0 && r.bottom <= innerHeight && r.right <= innerWidth) return false
        el.scrollIntoView({ block: "center", inline: "center", behavior: "instant" })
        return true
      }, null, t)
      if (!scrolled) break
      await this.settle()
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
    await this.eachFrame(forceRepaint)
    if (!full) return this.page.screenshot({ animations: "allow", caret: "hide", timeout: 60000 })
    await this.eachFrame(() => { window.__parityHoldResize = true })
    try {
      return await this.page.screenshot({ fullPage: true, animations: "allow", caret: "hide", timeout: 60000 })
    } finally {
      await this.idle()
      await this.eachFrame(() => { window.__parityHoldResize = false })
    }
  }

  async snapshot(origins) {
    return this.page.evaluate(snapshotDOM, origins)
  }

  async close() {
    await this.context?.close().catch(() => {})
  }
}

// --- comparison ---------------------------------------------------------------------------

function padPNG(P, width, height) {
  if (P.width === width && P.height === height) return P
  const Q = new PNG({ width, height })
  const bg = P.data.readUInt32BE(0)
  for (let i = 0; i < Q.data.length; i += 4) Q.data.writeUInt32BE(bg, i)
  PNG.bitblt(P, Q, 0, 0, P.width, P.height, 0, 0)
  return Q
}

function comparePNG(a, b) {
  let A = PNG.sync.read(a)
  let B = PNG.sync.read(b)
  const size = A.width !== B.width || A.height !== B.height ? `${A.width}x${A.height} vs ${B.width}x${B.height}` : undefined
  if (!EXACT && size) {
    const w = Math.max(A.width, B.width)
    const hh = Math.max(A.height, B.height)
    A = padPNG(A, w, hh)
    B = padPNG(B, w, hh)
  }
  if (A.width !== B.width || A.height !== B.height) {
    return { pixels: Math.max(A.width * A.height, B.width * B.height), size: `${A.width}x${A.height} vs ${B.width}x${B.height}`, diff: null, regions: [] }
  }
  const diff = new PNG({ width: A.width, height: A.height })
  const pixels = pixelmatch(A.data, B.data, diff.data, A.width, A.height, { threshold: 0, includeAA: true })
  let total
  if (!EXACT) {
    let content = 0
    const bgA = A.data.readUInt32BE(0)
    const bgB = B.data.readUInt32BE(0)
    for (let i = 0; i < A.data.length; i += 4) {
      if (A.data.readUInt32BE(i) !== bgA || B.data.readUInt32BE(i) !== bgB) content++
    }
    total = Math.max(content, 1)
  }
  return { pixels, total, size, diff: pixels ? PNG.sync.write(diff) : null, regions: pixels ? regionsOf(A, B) : [] }
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
    // A real mismatch is deterministic: a capture that differs is taken
    // again (after another settle) up to twice, and the last one is reported.
    for (let attempt = 0; attempt < 3; attempt++) {
      if (attempt) {
        await both((s) => s.eachFrame(domQuiet, { quiet: QUIET, frames: 8, max: 5000 }))
        await both((s) => s.eachFrame(settleAnimations))
        await both((s) => s.idle())
      }
      ;[rd, gd] = await both((s) => s.snapshot(origins))
      ;[r, g] = await both((s) => s.shot(full))
      cmp = comparePNG(r, g)
      if (EXACT ? cmp.pixels === 0 && rd.lines.join("\n") === gd.lines.join("\n") : !tooDifferent(cmp)) break
      cap.retries = attempt + 1
    }
    cap.pixels = cmp.pixels
    if (cmp.total) cap.content = cmp.total
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
    if (EXACT ? cmp.pixels || cap.dom.length || cap.state : tooDifferent(cmp) || cap.state) cap.pass = false
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
          await Promise.allSettled([ref, tgt].map((s) => s.bringIntoView(first)))
          const trial = await Promise.allSettled([ref, tgt].map((s) => s.locator(first).click({ trial: true, timeout: 5000 })))
          if (trial[0].status === "rejected" && /intercepts pointer events/.test(String(trial[0].reason?.message))) {
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
    if (ref.settleTimeouts || tgt.settleTimeouts) report.settleTimeouts = { ref: ref.settleTimeouts ?? [], target: tgt.settleTimeouts ?? [] }
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
  let pages
  const jobs = []
  if (opt("--rerun")) {
    // Only the failed jobs of a previous run's summary.json.
    const prev = JSON.parse(fs.readFileSync(opt("--rerun"), "utf8"))
    const known = new Map(allPages().map((p) => [p.path, p]))
    for (const f of prev.failures) jobs.push({ page: known.get(f.path) ?? { path: f.path, group: "app" }, theme: f.theme, vp: f.viewport })
    pages = [...new Set(jobs.map((j) => j.page.path))]
  } else {
    pages = selectPages()
    for (const p of pages) for (const vp of viewports) for (const theme of themes) jobs.push({ page: p, theme, vp })
  }
  console.log(`${opt("--rerun") ? `rerun of ${opt("--rerun")}: ${pages.length} pages,` : `${pages.length} pages x ${themes.length} themes x ${viewports.length} viewports =`} ${jobs.length} jobs; ref ${REF} vs target ${TARGET}${SELF ? " (self)" : ""}`)
  fs.mkdirSync(outRoot, { recursive: true })
  // D8: no partial raster, so a re-raster never keeps anti-aliased edges
  // from an earlier paint (1-gray-level differences in iframes otherwise).
  const browser = await chromium.launch({ executablePath: chromiumPath(), args: LAUNCH_ARGS })
  const reports = []
  const queue = [...jobs]
  const t0 = Date.now()
  await Promise.all(Array.from({ length: Math.min(jobsN, queue.length) }, async () => {
    while (queue.length) {
      const j = queue.shift()
      let r
      // A harness error (a timeout on an overloaded machine, a crashed
      // context) is retried once; a difference never is.
      let retried = false
      for (let attempt = 0; attempt < 2; attempt++) {
        try {
          r = await runJob(browser, j.page, j.theme, j.vp)
        } catch (e) {
          r = { path: j.page.path, theme: j.theme, viewport: j.vp, pass: false, captures: [], errors: { run: String(e.stack ?? e) } }
        }
        if (!r.errors?.run || attempt === 1) break
        console.log(`     retrying ${j.page.path} ${j.theme}/${j.vp} after a harness error: ${r.errors.run.split("\n")[0]}`)
        retried = true
      }
      if (retried) r.retried = true
      reports.push(r)
      console.log(`[${reports.length}/${jobs.length}] ${summaryLine(r)}`)
    }
  }))
  await browser.close()
  const failed = reports.filter((r) => !r.pass)
  const captures = reports.reduce((n, r) => n + r.captures.length, 0)
  const summary = {
    ref: REF, target: TARGET, self: SELF, jobs: reports.length, captures, failed: failed.length, seconds: Math.round((Date.now() - t0) / 1000),
    retriedCaptures: reports.reduce((n, r) => n + r.captures.filter((c) => c.pass && c.retries).length, 0),
    retriedJobs: reports.filter((r) => r.retried).length,
    skippedScenarios: reports.reduce((n, r) => n + (r.skipped?.length ?? 0), 0),
    failures: failed.map((r) => ({ path: r.path, theme: r.theme, viewport: r.viewport, captures: r.captures.filter((c) => !c.pass).map((c) => c.label), error: r.errors?.run?.split("\n")[0] })),
  }
  fs.writeFileSync(path.join(outRoot, "summary.json"), JSON.stringify(summary, null, 2))
  console.log(`\n${reports.length - failed.length}/${reports.length} jobs ${EXACT ? "identical" : `within ${TOLERANCE * 100}% of the reference`} (${captures} captures, ${summary.seconds}s). Reports: ${path.relative(process.cwd(), outRoot) || "."}/`)
  return failed.length ? 1 : 0
}

if (import.meta.url === `file://${process.argv[1]}`) {
  main().then((code) => process.exit(code), (e) => { console.error(e.stack ?? e); process.exit(2) })
}
