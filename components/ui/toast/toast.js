// toast.js: Base UI's toast manager (@base-ui/react 1.6 toast/*) with
// styles/base-nova/ui/toast.tsx's markup. Toasts are prepended to the
// viewport (newest first) and stacked with CSS variables like Base UI:
// --toast-index, --toast-offset-y (sum of the newer toasts' heights),
// --toast-height and the viewport's --toast-frontmost-height; toasts past
// the limit (3) get data-limited and inert; hovering or focusing the
// viewport expands the stack (data-expanded) and pauses the timers
// (5s by default); data-starting-style / data-ending-style drive the
// enter and exit transitions; a swipe down or right dismisses a toast.

import { define, uid } from "./core.js"

const LIMIT = 3
const TIMEOUT = 5000
const SWIPE_THRESHOLD = 40

const VIEWPORT_CLASS =
  "pointer-events-none fixed inset-x-4 bottom-4 z-50 mx-auto w-auto max-w-sm outline-none sm:right-4 sm:left-auto sm:mx-0 sm:w-full"
const TOAST_CLASS = [
  "group/toast pointer-events-auto absolute right-0 bottom-0 z-[calc(1000-var(--toast-index))] w-full origin-bottom rounded-2xl border bg-popover text-popover-foreground shadow-lg will-change-transform outline-none select-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
  "[--gap:0.75rem] [--height:var(--toast-frontmost-height,var(--toast-height))] [--offset-y:calc(var(--toast-offset-y)*-1+calc(var(--toast-index)*var(--gap)*-1)+var(--toast-swipe-movement-y))] [--peek:0.75rem] [--scale:calc(max(0,1-(var(--toast-index)*0.1)))] [--shrink:calc(1-var(--scale))]",
  "h-(--height) [transform:translateX(var(--toast-swipe-movement-x))_translateY(calc(var(--toast-swipe-movement-y)-(var(--toast-index)*var(--peek))-(var(--shrink)*var(--height))))_scale(var(--scale))] [transition:transform_500ms_cubic-bezier(0.22,1,0.36,1),opacity_500ms,height_150ms]",
  "after:absolute after:top-full after:left-0 after:h-[calc(var(--gap)+1px)] after:w-full after:content-['']",
  "data-expanded:h-(--toast-height) data-expanded:[transform:translateX(var(--toast-swipe-movement-x))_translateY(var(--offset-y))]",
  "data-limited:opacity-0 data-starting-style:[transform:translateY(150%)]",
  "[&[data-ending-style]:not([data-limited]):not([data-swipe-direction])]:[transform:translateY(150%)]",
  "data-ending-style:data-[swipe-direction=down]:[transform:translateY(calc(var(--toast-swipe-movement-y)+150%))]",
  "data-ending-style:data-[swipe-direction=left]:[transform:translateX(calc(var(--toast-swipe-movement-x)-150%))_translateY(var(--offset-y))]",
  "data-ending-style:data-[swipe-direction=right]:[transform:translateX(calc(var(--toast-swipe-movement-x)+150%))_translateY(var(--offset-y))]",
  "data-ending-style:data-[swipe-direction=up]:[transform:translateY(calc(var(--toast-swipe-movement-y)-150%))]",
  "data-expanded:data-ending-style:data-[swipe-direction=down]:[transform:translateY(calc(var(--toast-swipe-movement-y)+150%))]",
  "data-expanded:data-ending-style:data-[swipe-direction=left]:[transform:translateX(calc(var(--toast-swipe-movement-x)-150%))_translateY(var(--offset-y))]",
  "data-expanded:data-ending-style:data-[swipe-direction=right]:[transform:translateX(calc(var(--toast-swipe-movement-x)+150%))_translateY(var(--offset-y))]",
  "data-expanded:data-ending-style:data-[swipe-direction=up]:[transform:translateY(calc(var(--toast-swipe-movement-y)-150%))]",
].join(" ")
const CONTENT_CLASS =
  "flex h-full items-center gap-3 overflow-hidden p-4 transition-opacity duration-250 ease-[cubic-bezier(0.22,1,0.36,1)] data-behind:opacity-0 data-expanded:opacity-100"
const ACTION_CLASS =
  "group/button inline-flex items-center justify-center border bg-clip-padding font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 border-border bg-background hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5 shrink-0"
const CLOSE_CLASS =
  "group/button inline-flex items-center justify-center border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 hover:bg-muted aria-expanded:bg-muted aria-expanded:text-foreground dark:hover:bg-muted/50 size-7 rounded-[min(var(--radius-md),12px)] in-data-[slot=button-group]:rounded-lg relative shrink-0 text-muted-foreground after:absolute after:-inset-2 after:content-[''] hover:text-foreground"

const svg = (name, body, cls = "") =>
  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-${name}${cls ? " " + cls : ""}" aria-hidden="true">${body}</svg>`
const ICONS = {
  success: svg("circle-check", '<circle cx="12" cy="12" r="10"></circle><path d="m9 12 2 2 4-4"></path>'),
  info: svg("info", '<circle cx="12" cy="12" r="10"></circle><path d="M12 16v-4"></path><path d="M12 8h.01"></path>'),
  warning: svg("triangle-alert", '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"></path><path d="M12 9v4"></path><path d="M12 17h.01"></path>'),
  error: svg("octagon-x", '<path d="m15 9-6 6"></path><path d="M2.586 16.726A2 2 0 0 1 2 15.312V8.688a2 2 0 0 1 .586-1.414l4.688-4.688A2 2 0 0 1 8.688 2h6.624a2 2 0 0 1 1.414.586l4.688 4.688A2 2 0 0 1 22 8.688v6.624a2 2 0 0 1-.586 1.414l-4.688 4.688a2 2 0 0 1-1.414.586H8.688a2 2 0 0 1-1.414-.586z"></path><path d="m9 9 6 6"></path>', "text-destructive"),
  loading: svg("loader-circle", '<path d="M21 12a9 9 0 1 1-6.219-8.56"></path>', "animate-spin"),
}
const X_ICON = svg("x", '<path d="M18 6 6 18"></path><path d="m6 6 12 12"></path>')

/** @type {Map<string, {el: HTMLElement, opts: object, timer: any, remaining: number, started: number}>} */
const toasts = new Map()
let order = [] // ids, newest first
let expanded = false

function viewport() {
  let vp = document.querySelector('[data-slot="toast-viewport"]')
  if (!vp) {
    const portal = document.createElement("div")
    portal.setAttribute("data-base-ui-portal", "")
    portal.dataset.slot = "toast-portal"
    vp = document.createElement("div")
    for (const [k, v] of Object.entries({
      tabindex: "-1",
      role: "region",
      "aria-live": "polite",
      "aria-atomic": "false",
      "aria-relevant": "additions text",
      "aria-label": "Notifications",
      "data-slot": "toast-viewport",
      class: VIEWPORT_CLASS,
    }))
      vp.setAttribute(k, v)
    portal.append(vp)
    document.body.append(portal)
  }
  if (!vp.dataset.toastBound) {
    vp.dataset.toastBound = "1"
    const expand = (on) => {
      if (expanded === on) return
      expanded = on
      for (const t of toasts.values()) on ? pause(t) : resume(t)
      layout()
    }
    vp.addEventListener("mouseenter", () => expand(true))
    vp.addEventListener("mouseleave", () => {
      if (!vp.contains(document.activeElement)) expand(false)
    })
    vp.addEventListener("focusin", () => expand(true))
    vp.addEventListener("focusout", (e) => {
      if (!vp.contains(e.relatedTarget) && !vp.matches(":hover")) expand(false)
    })
  }
  return vp
}

function escape(s) {
  return String(s).replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c])
}

function render(t) {
  const { opts, el } = t
  const type = opts.type
  if (type) el.dataset.type = type
  else el.removeAttribute("data-type")
  el.setAttribute("role", opts.priority === "high" ? "alertdialog" : "dialog")
  if (opts.priority === "high") el.setAttribute("aria-hidden", "true")
  const dt = type ? ` data-type="${type}"` : ""
  const titleId = opts.title ? (t.titleId ||= uid("base-ui")) : null
  const descId = opts.description ? (t.descId ||= uid("base-ui")) : null
  if (titleId) el.setAttribute("aria-labelledby", titleId)
  else el.removeAttribute("aria-labelledby")
  if (descId) el.setAttribute("aria-describedby", descId)
  else el.removeAttribute("aria-describedby")
  const icon = ICONS[type]
  el.innerHTML =
    `<div data-slot="toast-content" class="${CONTENT_CLASS}">` +
    (icon ? `<span data-slot="toast-icon" class="shrink-0 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4">${icon}</span>` : "") +
    `<div class="flex min-w-0 flex-1 flex-col gap-1">` +
    (opts.title ? `<h2${dt} data-slot="toast-title" id="${titleId}" class="text-sm font-medium">${escape(opts.title)}</h2>` : "") +
    (opts.description ? `<p${dt} data-slot="toast-description" id="${descId}" class="text-sm text-muted-foreground">${escape(opts.description)}</p>` : "") +
    `</div>` +
    (opts.action ? `<button type="button" tabindex="0"${dt} data-slot="toast-action" class="${ACTION_CLASS}">${escape(opts.action)}</button>` : "") +
    `<button type="button" tabindex="0" data-slot="toast-close"${dt} aria-hidden="true" aria-label="Close toast" class="${CLOSE_CLASS}">${X_ICON}</button>` +
    `</div>`
  el.querySelector('[data-slot="toast-close"]').addEventListener("click", () => close(t.id))
  el.querySelector('[data-slot="toast-action"]')?.addEventListener("click", () => (opts.onAction ? opts.onAction() : close(t.id)))
  measure(t)
}

function measure(t) {
  const el = t.el
  const h = el.style.height
  el.style.height = "auto"
  t.height = el.offsetHeight
  el.style.height = h
  el.style.setProperty("--toast-height", `${t.height}px`)
}

function layout() {
  const vp = viewport()
  let offset = 0
  const live = order.filter((id) => !toasts.get(id).closing)
  live.forEach((id, index) => {
    const t = toasts.get(id)
    const el = t.el
    el.style.setProperty("--toast-index", String(index))
    el.style.setProperty("--toast-offset-y", `${offset}px`)
    offset += t.height
    const limited = index >= LIMIT
    el.toggleAttribute("data-limited", limited)
    el.toggleAttribute("inert", limited)
    el.toggleAttribute("data-expanded", expanded)
    const content = el.querySelector('[data-slot="toast-content"]')
    content?.toggleAttribute("data-behind", index > 0)
    content?.toggleAttribute("data-expanded", expanded)
  })
  const front = live[0] && toasts.get(live[0])
  if (front) vp.style.setProperty("--toast-frontmost-height", `${front.height}px`)
}

function schedule(t) {
  clearTimeout(t.timer)
  if (t.opts.type === "loading" || t.remaining <= 0) return
  t.started = performance.now()
  t.timer = setTimeout(() => close(t.id), t.remaining)
}
function pause(t) {
  if (!t.timer) return
  clearTimeout(t.timer)
  t.timer = null
  t.remaining -= performance.now() - t.started
}
function resume(t) {
  if (t.timer || t.closing) return
  schedule(t)
}

function add(opts = {}) {
  const vp = viewport()
  const id = opts.id ?? uid("toast")
  const el = document.createElement("div")
  el.setAttribute("tabindex", "0")
  el.setAttribute("aria-modal", "false")
  el.dataset.slot = "toast"
  el.className = TOAST_CLASS
  for (const v of ["--toast-swipe-movement-x", "--toast-swipe-movement-y"]) el.style.setProperty(v, "0px")
  el.style.setProperty("--toast-index", "0")
  el.style.setProperty("--toast-offset-y", "0px")
  el.setAttribute("data-starting-style", "")
  const t = { id, el, opts, timer: null, remaining: opts.timeout ?? TIMEOUT, height: 0 }
  toasts.set(id, t)
  order.unshift(id)
  vp.prepend(el)
  render(t)
  bindSwipe(t)
  layout()
  requestAnimationFrame(() => requestAnimationFrame(() => el.removeAttribute("data-starting-style")))
  if (!expanded) schedule(t)
  return id
}

function update(id, opts) {
  const t = toasts.get(id)
  if (!t) return
  t.opts = { ...t.opts, ...opts }
  render(t)
  layout()
  t.remaining = t.opts.timeout ?? TIMEOUT
  if (!expanded) schedule(t)
}

function close(id) {
  const t = toasts.get(id)
  if (!t || t.closing) return
  t.closing = true
  clearTimeout(t.timer)
  t.el.setAttribute("data-ending-style", "")
  layout()
  requestAnimationFrame(() =>
    Promise.all(t.el.getAnimations({ subtree: true }).map((a) => a.finished.catch(() => {}))).then(() => {
      t.el.remove()
      toasts.delete(id)
      order = order.filter((o) => o !== id)
      layout()
    }),
  )
}

function promise(p, { loading, success, error }) {
  const id = add({ type: "loading", description: loading })
  const text = (v, arg) => (typeof v === "function" ? v(arg) : v)
  p.then(
    (data) => update(id, { type: "success", description: text(success, data) }),
    (err) => update(id, { type: "error", description: text(error, err) }),
  )
  return p
}

function bindSwipe(t) {
  const el = t.el
  let start = null
  el.addEventListener("pointerdown", (e) => {
    if (e.button !== 0 || e.target.closest("button")) return
    start = { x: e.clientX, y: e.clientY }
    el.setPointerCapture?.(e.pointerId)
  })
  el.addEventListener("pointermove", (e) => {
    if (!start) return
    const dx = Math.max(0, e.clientX - start.x)
    const dy = Math.max(0, e.clientY - start.y)
    const horizontal = dx > dy
    el.setAttribute("data-swiping", "")
    el.style.setProperty("--toast-swipe-movement-x", `${horizontal ? dx : 0}px`)
    el.style.setProperty("--toast-swipe-movement-y", `${horizontal ? 0 : dy}px`)
  })
  const end = (e) => {
    if (!start) return
    const dx = e.clientX - start.x
    const dy = e.clientY - start.y
    start = null
    el.removeAttribute("data-swiping")
    if (dx > SWIPE_THRESHOLD && dx >= dy) {
      el.dataset.swipeDirection = "right"
      close(t.id)
    } else if (dy > SWIPE_THRESHOLD) {
      el.dataset.swipeDirection = "down"
      close(t.id)
    } else {
      el.style.setProperty("--toast-swipe-movement-x", "0px")
      el.style.setProperty("--toast-swipe-movement-y", "0px")
    }
  }
  el.addEventListener("pointerup", end)
  el.addEventListener("pointercancel", end)
}

export const toast = { add, close, update, promise }

// Declarative triggers rendered by toast.go (Add, Promise).
define("toast-add", {
  selector: "[data-toast-add]",
  mount(el) {
    const onClick = () => add(JSON.parse(el.dataset.toastAdd))
    el.addEventListener("click", onClick)
    return () => el.removeEventListener("click", onClick)
  },
})
define("toast-promise", {
  selector: "[data-toast-promise]",
  mount(el) {
    const onClick = () => {
      const o = JSON.parse(el.dataset.toastPromise)
      promise(new Promise((resolve) => setTimeout(() => resolve({ name: o.name }), o.delay)), {
        loading: o.loading,
        success: (data) => o.success.replace("{name}", data?.name ?? ""),
        error: o.error,
      })
    }
    el.addEventListener("click", onClick)
    return () => el.removeEventListener("click", onClick)
  },
})
