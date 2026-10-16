// basedrawer.js: Base UI's Drawer (@base-ui/react 1.6 drawer/*) over the
// markup of basedrawer.go. A trigger ([data-base-drawer-trigger]) opens
// the drawer of the next template[data-base-drawer] sibling: the template's
// backdrop and viewport are mounted at the end of <body> and driven like
// Base UI drives them:
//
//   - data-open / data-closed, data-starting-style on mount (one frame),
//     data-ending-style while closing; unmounted once transitions end.
//   - popup CSS variables: --drawer-swipe-movement-x/y, --nested-drawers,
//     --drawer-height (kept while a nested drawer is open),
//     --drawer-frontmost-height, --drawer-snap-point-offset,
//     --drawer-swipe-strength; backdrop --drawer-swipe-progress.
//   - nested drawers (a trigger inside an open popup): no backdrop; every
//     ancestor counts its open descendants (data-nested-drawer-open).
//   - swipe to dismiss in the swipe direction (half the size, or a fast
//     flick), snap points for down/up drawers, swipe handles.
//   - Escape and outside presses on the backdrop close the topmost drawer
//     (unless disablePointerDismissal); a modal drawer locks page scroll.

import { uid } from "./core.js"

const MIN_SWIPE_THRESHOLD = 10
const FAST_SWIPE_VELOCITY = 0.5
const SNAP_VELOCITY_THRESHOLD = 0.5
const SNAP_VELOCITY_MULTIPLIER = 300
const MAX_SNAP_VELOCITY = 4
const MIN_SWIPE_RELEASE_VELOCITY = 0.2
const MAX_SWIPE_RELEASE_VELOCITY = 4
const MIN_SWIPE_RELEASE_DURATION_MS = 80
const MAX_SWIPE_RELEASE_DURATION_MS = 360

const clamp = (v, lo, hi) => Math.min(Math.max(v, lo), hi)
const isMobile = () => matchMedia("(max-width: 767px)").matches

/** @type {Drawer[]} open drawers, outermost first */
const stack = []

function templateFor(trigger) {
  for (let n = trigger.nextElementSibling; n; n = n.nextElementSibling) {
    if (n.matches("template[data-base-drawer]")) return n
  }
  return null
}

class Drawer {
  constructor(trigger, template) {
    this.trigger = trigger
    this.template = template
    this.cfg = JSON.parse(template.dataset.baseDrawer)
    this.parent = stack.find((d) => d.popup.contains(trigger)) ?? null
    this.nested = !!this.parent
    this.nestedCount = 0
    this.height = 0
    this.frontmostHeight = 0
    this.activeSnap = undefined
  }

  get direction() {
    return (isMobile() && this.cfg.mobileSwipeDirection) || this.cfg.swipeDirection || "down"
  }

  open() {
    const frag = this.template.content.cloneNode(true)
    this.backdrop = frag.querySelector('[data-slot="drawer-overlay"]')
    this.viewport = frag.querySelector('[data-slot="drawer-viewport"]')
    this.popup = frag.querySelector('[data-slot="drawer-popup"]')
    if (this.nested && this.backdrop) {
      this.backdrop.remove()
      this.backdrop = null
    }
    const dir = this.direction
    this.dir = dir
    this.axis = dir === "left" || dir === "right" ? "x" : "y"
    this.popup.dataset.swipeDirection = dir
    this.popup.dataset.swipeAxis = this.axis
    const handle = this.popup.querySelector(':scope > [data-slot="drawer-swipe-handle"]')
    if (handle?.hasAttribute("data-mobile-only")) {
      if (isMobile()) handle.removeAttribute("data-mobile-only")
      else handle.remove()
    }
    const title = this.popup.querySelector('[data-slot="drawer-title"]')
    const desc = this.popup.querySelector('[data-slot="drawer-description"]')
    if (title) this.popup.setAttribute("aria-labelledby", (title.id ||= uid("base-ui")))
    if (desc) this.popup.setAttribute("aria-describedby", (desc.id ||= uid("base-ui")))
    this.popup.id = uid("base-ui")
    for (const el of this.parts()) {
      el.setAttribute("data-open", "")
      el.setAttribute("data-starting-style", "")
    }
    this.viewport.removeAttribute("data-starting-style")
    if (this.backdrop) {
      this.backdrop.setAttribute("aria-hidden", "true")
      Object.assign(this.backdrop.style, { userSelect: "none" })
      this.backdrop.style.setProperty("--drawer-swipe-progress", "0")
      this.backdrop.style.setProperty("--drawer-swipe-strength", "1")
    }
    this.setVars()
    document.body.append(...[this.backdrop, this.viewport].filter(Boolean))
    stack.push(this)
    if (this.cfg.modal && !this.nested) lockScroll(true)
    this.trigger.setAttribute("aria-expanded", "true")
    this.trigger.setAttribute("aria-controls", this.popup.id)
    this.trigger.setAttribute("data-popup-open", "")
    this.measure()
    this.ro = new ResizeObserver(() => this.measure())
    this.ro.observe(this.popup)
    this.popup.focus({ preventScroll: true })
    requestAnimationFrame(() => requestAnimationFrame(() => {
      for (const el of this.parts()) el.removeAttribute("data-starting-style")
    }))
    for (let p = this.parent; p; p = p.parent) {
      p.nestedCount += 1
      p.setVars()
    }
    this.bindSwipe()
  }

  parts() {
    return [this.backdrop, this.viewport, this.popup].filter(Boolean)
  }

  measure() {
    if (!this.popup) return
    const h = this.popup.offsetHeight
    if (this.nestedCount > 0 && this.height > 0) return
    if (h === this.height) return
    this.height = h
    this.setVars()
    this.notifyFrontmost()
  }

  notifyFrontmost() {
    // Each ancestor learns the height of the frontmost open drawer.
    const front = stack[stack.length - 1]
    for (const d of stack) {
      d.frontmostHeight = front?.height ?? 0
      d.setVars()
    }
  }

  snapPoints() {
    const points = this.cfg.snapPoints
    if (!points?.length || this.axis !== "y" || !this.height) return []
    const vh = this.viewport.offsetHeight || document.documentElement.clientHeight
    const rem = parseFloat(getComputedStyle(document.documentElement).fontSize) || 16
    const max = Math.min(this.height, vh)
    const out = []
    for (const value of points) {
      let px = null
      if (typeof value === "number") px = value <= 1 ? clamp(value, 0, 1) * vh : value
      else if (value.endsWith("px")) px = parseFloat(value)
      else if (value.endsWith("rem")) px = parseFloat(value) * rem
      if (px == null || !Number.isFinite(px)) continue
      const height = clamp(px, 0, max)
      out.push({ value, height, offset: Math.max(0, this.height - height) })
    }
    return out
  }

  activeSnapPoint() {
    const points = this.snapPoints()
    if (!points.length) return null
    if (this.activeSnap === undefined) return points[0]
    return points.find((p) => Object.is(p.value, this.activeSnap)) ?? points[0]
  }

  // snapProgress is the backdrop's swipe progress for a snap point offset:
  // 0 at the fullest snap point, 1 at the next one (Base UI's
  // snapPointRange), or null without snap points.
  snapProgress(offset) {
    const offsets = this.snapPoints().map((p) => p.offset).sort((a, b) => a - b)
    if (offsets.length < 2) return null
    let range = offsets[1] - offsets[0]
    if (!(range > 0)) range = offsets[offsets.length - 1] - offsets[0]
    if (!(range > 0)) return null
    return clamp((offset - offsets[0]) / range, 0, 1)
  }

  setBackdropProgress(progress) {
    const b = this.backdrop
    if (!b || this.nested) return
    b.style.setProperty("--drawer-swipe-progress", `${progress > 0 ? progress : 0}`)
    if (progress > 0 && this.frontmostHeight) b.style.setProperty("--drawer-height", `${this.frontmostHeight}px`)
    else b.style.removeProperty("--drawer-height")
  }

  setVars(movement = { x: 0, y: 0 }) {
    const p = this.popup
    if (!p) return
    const s = p.style
    s.setProperty("--drawer-swipe-movement-x", `${movement.x}px`)
    s.setProperty("--drawer-swipe-movement-y", `${movement.y}px`)
    s.setProperty("--drawer-swipe-progress", "0")
    s.setProperty("--nested-drawers", String(this.nestedCount))
    const snap = this.activeSnapPoint()
    let offset = snap ? snap.offset : 0
    if (this.dir === "up") offset = -offset
    s.setProperty("--drawer-snap-point-offset", `${offset}px`)
    if (!s.getPropertyValue("--drawer-swipe-strength")) s.setProperty("--drawer-swipe-strength", "1")
    if (this.nestedCount > 0 && this.height) s.setProperty("--drawer-height", `${this.height}px`)
    else s.removeProperty("--drawer-height")
    if (this.frontmostHeight) s.setProperty("--drawer-frontmost-height", `${this.frontmostHeight}px`)
    else s.removeProperty("--drawer-frontmost-height")
    if (!this.swiping) this.setBackdropProgress(snap ? this.snapProgress(snap.offset) ?? 0 : 0)
    p.toggleAttribute("data-nested-drawer-open", this.nestedCount > 0)
    p.toggleAttribute("data-expanded", !!snap && snap.value === 1)
  }

  close() {
    if (this.closing) return
    this.closing = true
    // Nested drawers close with their parent.
    for (const d of [...stack].reverse()) if (d.parent === this) d.close()
    const i = stack.indexOf(this)
    if (i >= 0) stack.splice(i, 1)
    for (let p = this.parent; p; p = p.parent) {
      p.nestedCount = Math.max(0, p.nestedCount - 1)
      p.setVars()
    }
    this.notifyFrontmost()
    for (const el of this.parts()) {
      el.removeAttribute("data-open")
      el.setAttribute("data-closed", "")
      el.setAttribute("data-ending-style", "")
    }
    if (this.backdrop) this.backdrop.style.pointerEvents = "none"
    this.trigger.setAttribute("aria-expanded", "false")
    this.trigger.removeAttribute("aria-controls")
    this.trigger.removeAttribute("data-popup-open")
    if (this.cfg.modal && !this.nested) lockScroll(false)
    this.unbindSwipe?.()
    this.ro?.disconnect()
    const done = () => {
      for (const el of [this.backdrop, this.viewport]) el?.remove()
      this.popup = null
    }
    const wait = () =>
      Promise.all(
        this.parts().flatMap((el) => el.getAnimations({ subtree: true }).map((a) => a.finished.catch(() => {}))),
      )
    requestAnimationFrame(() => wait().then(done))
    if (this.trigger.isConnected) this.trigger.focus({ preventScroll: true })
  }

  // --- swipe ------------------------------------------------------------------

  bindSwipe() {
    const popup = this.popup
    let start = null
    let last = null
    let velocity = { x: 0, y: 0 }
    const sign = this.dir === "left" || this.dir === "up" ? -1 : 1
    const onDown = (e) => {
      if (e.button !== 0 || this.nestedCount > 0) return
      if (e.target.closest('button, a, input, select, textarea, label, [role="button"]')) return
      // A mouse drag inside the content selects text instead (as in Base
      // UI); touch swipes start anywhere.
      if (e.pointerType !== "touch" && e.target.closest("[data-drawer-content]")) return
      start = { x: e.clientX, y: e.clientY, t: performance.now() }
      last = { ...start }
      velocity = { x: 0, y: 0 }
      this.swiping = false
    }
    const movement = (e) => {
      const dx = e.clientX - start.x
      const dy = e.clientY - start.y
      const delta = this.axis === "x" ? dx : dy
      let along = delta * sign
      const snap = this.snapPoints().length > 0 && this.axis === "y"
      if (snap) {
        // Snap drawers move both ways; past the top they resist.
        const base = this.activeSnapPoint()?.offset ?? 0
        const next = base + along
        if (next < 0) along = -Math.sqrt(-next) - base
      } else if (along < 0) {
        along = -Math.sqrt(-along)
      }
      return { delta, along }
    }
    const onMove = (e) => {
      if (!start) return
      const now = performance.now()
      const dt = Math.max(1, now - last.t)
      velocity = { x: (e.clientX - last.x) / dt, y: (e.clientY - last.y) / dt }
      last = { x: e.clientX, y: e.clientY, t: now }
      const { delta, along } = movement(e)
      if (!this.swiping) {
        if (Math.abs(delta) < MIN_SWIPE_THRESHOLD) return
        this.swiping = true
        popup.setPointerCapture?.(e.pointerId)
        for (const el of this.parts()) el.setAttribute("data-swiping", "")
        if (this.parent) this.parent.popup?.setAttribute("data-nested-drawer-swiping", "")
      }
      const v = along * sign
      const m = this.axis === "x" ? { x: v, y: 0 } : { x: 0, y: v }
      this.setVars(m)
      const size = this.axis === "x" ? popup.offsetWidth : popup.offsetHeight
      const snap = this.activeSnapPoint()
      let progress = clamp(Math.max(0, along) / (size || 1), 0, 1)
      if (snap) progress = this.snapProgress(clamp(snap.offset + along, 0, this.height)) ?? progress
      this.setBackdropProgress(progress)
      if (this.parent?.popup) this.parent.popup.style.setProperty("--drawer-swipe-progress", `${progress}`)
    }
    const onUp = (e) => {
      if (!start) return
      const s = start
      start = null
      if (!this.swiping) return
      this.swiping = false
      for (const el of this.parts()) el.removeAttribute("data-swiping")
      this.parent?.popup?.removeAttribute("data-nested-drawer-swiping")
      this.parent?.popup?.style.setProperty("--drawer-swipe-progress", "0")
      const dx = e.clientX - s.x
      const dy = e.clientY - s.y
      const delta = this.axis === "x" ? dx : dy
      const along = delta * sign
      const v = (this.axis === "x" ? velocity.x : velocity.y) * sign
      const size = this.axis === "x" ? popup.offsetWidth : popup.offsetHeight
      const points = this.snapPoints()
      let close = false
      if (points.length) {
        const current = this.activeSnapPoint()?.offset ?? 0
        const dragTarget = clamp(current + along, 0, this.height)
        const vOffset = Math.abs(v) >= SNAP_VELOCITY_THRESHOLD ? clamp(v, -MAX_SNAP_VELOCITY, MAX_SNAP_VELOCITY) * SNAP_VELOCITY_MULTIPLIER : 0
        const target = clamp(dragTarget + vOffset, 0, this.height)
        if (v >= FAST_SWIPE_VELOCITY && along > 0) close = true
        else {
          let best = points[0]
          for (const p of points) if (Math.abs(target - p.offset) < Math.abs(target - best.offset)) best = p
          if (Math.abs(target - this.height) < Math.abs(target - best.offset)) close = true
          else this.activeSnap = best.value
        }
      } else if (along > 0 && (v >= FAST_SWIPE_VELOCITY || along > Math.max(size * 0.5, MIN_SWIPE_THRESHOLD))) {
        close = true
      }
      if (close) {
        const remaining = Math.max(0, size - along - (points.length ? this.activeSnapPoint()?.offset ?? 0 : 0))
        let strength = 1
        if (v > MIN_SWIPE_RELEASE_VELOCITY && remaining > 0) {
          const cv = clamp(v, MIN_SWIPE_RELEASE_VELOCITY, MAX_SWIPE_RELEASE_VELOCITY)
          const ms = clamp(remaining / cv, MIN_SWIPE_RELEASE_DURATION_MS, MAX_SWIPE_RELEASE_DURATION_MS)
          strength = clamp(0.1 + ((ms - MIN_SWIPE_RELEASE_DURATION_MS) / (MAX_SWIPE_RELEASE_DURATION_MS - MIN_SWIPE_RELEASE_DURATION_MS)) * 0.9, 0.1, 1)
        }
        for (const el of this.parts()) el.style.setProperty("--drawer-swipe-strength", `${strength}`)
        // Close from where the swipe left the popup.
        const m = this.axis === "x" ? { x: along * sign, y: 0 } : { x: 0, y: along * sign }
        this.setVars(m)
        this.close()
      } else {
        this.setVars()
      }
    }
    popup.addEventListener("pointerdown", onDown)
    window.addEventListener("pointermove", onMove)
    window.addEventListener("pointerup", onUp)
    window.addEventListener("pointercancel", onUp)
    this.unbindSwipe = () => {
      popup.removeEventListener("pointerdown", onDown)
      window.removeEventListener("pointermove", onMove)
      window.removeEventListener("pointerup", onUp)
      window.removeEventListener("pointercancel", onUp)
    }
  }
}

let locked = 0
let saved = ""
function lockScroll(on) {
  const html = document.documentElement
  if (on) {
    if (locked++ === 0) {
      saved = html.style.overflow
      html.style.overflow = "hidden"
    }
  } else if (locked > 0 && --locked === 0) {
    html.style.overflow = saved
  }
}

document.addEventListener("click", (e) => {
  const t = e.target instanceof Element ? e.target : null
  if (!t) return
  const trigger = t.closest("[data-base-drawer-trigger]")
  if (trigger) {
    const template = templateFor(trigger)
    if (!template) return
    if (stack.some((d) => d.trigger === trigger)) return
    new Drawer(trigger, template).open()
    return
  }
  const close = t.closest("[data-base-drawer-close]")
  if (close) [...stack].reverse().find((d) => d.popup?.contains(close))?.close()
})

// Outside presses: a press on the topmost drawer's backdrop or viewport
// (outside its popup) closes it.
document.addEventListener("pointerdown", (e) => {
  const top = stack[stack.length - 1]
  if (!top || top.cfg.disablePointerDismissal || e.button !== 0) return
  const t = e.target
  if (!(t instanceof Element)) return
  if (top.cfg.modal) {
    const backdrop = top.nested ? stack.find((d) => !d.nested)?.backdrop : top.backdrop
    if (t === backdrop || t === top.viewport) top.close()
  } else if (!top.popup.contains(t) && !top.trigger.contains(t)) {
    top.close()
  }
})

document.addEventListener("keydown", (e) => {
  if (e.key !== "Escape") return
  const top = stack[stack.length - 1]
  if (!top) return
  e.preventDefault()
  top.close()
})
