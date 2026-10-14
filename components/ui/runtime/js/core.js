// core.js: the shadcn.go client runtime.
//
// Server-rendered components carry their state in the DOM (data-slot,
// data-state, aria-*), exactly like the React components render it.
// Behaviors attach to that markup:
//
//   import { define } from "./core.js"
//   define("collapsible", {
//     selector: '[data-slot="collapsible"]',
//     mount(el) {
//       ...                  // wire up events, return a cleanup function
//       return () => { ... }
//     },
//   })
//
// mount runs once per matching element, for elements present when the
// page loads and for elements added later (any DOM insertion, including
// fetched HTML). The returned cleanup runs when the element leaves the DOM.

const behaviors = []
const mounted = new WeakMap() // element -> Map(behavior name -> cleanup)
let started = false

/**
 * Registers a behavior. Safe to call before or after start().
 * @param {string} name
 * @param {{ selector: string, mount: (el: Element) => (void | (() => void)) }} behavior
 */
export function define(name, behavior) {
  const b = { name, ...behavior }
  behaviors.push(b)
  if (started) scan(document, [b])
}

function mountOne(el, b) {
  let map = mounted.get(el)
  if (!map) mounted.set(el, (map = new Map()))
  if (map.has(b.name)) return
  map.set(b.name, null)
  try {
    const cleanup = b.mount(el)
    if (typeof cleanup === "function") map.set(b.name, cleanup)
  } catch (err) {
    console.error(`[shadcn] ${b.name}:`, err)
  }
}

function scan(root, list = behaviors) {
  for (const b of list) {
    if (root.nodeType === 1 && root.matches(b.selector)) mountOne(root, b)
    root.querySelectorAll?.(b.selector).forEach((el) => mountOne(el, b))
  }
}

function unmount(root) {
  const visit = (el) => {
    const map = mounted.get(el)
    if (!map) return
    for (const cleanup of map.values()) {
      try {
        cleanup?.()
      } catch (err) {
        console.error("[shadcn] cleanup:", err)
      }
    }
    mounted.delete(el)
  }
  if (root.nodeType !== 1) return
  visit(root)
  root.querySelectorAll("*").forEach(visit)
}

/** Starts the runtime. Called by the generated index module. */
export function start() {
  if (started) return
  started = true
  const run = () => {
    scan(document)
    window.__shadcnReady = true
    new MutationObserver((records) => {
      for (const r of records) {
        // Elements moved within the document (portals) are removed and
        // re-added in the same batch: only unmount what is really gone.
        r.removedNodes.forEach((n) => {
          if (n.nodeType === 1 && !n.isConnected) unmount(n)
        })
        r.addedNodes.forEach((n) => {
          if (n.nodeType === 1 && n.isConnected) scan(n)
        })
      }
    }).observe(document.documentElement, { childList: true, subtree: true })
  }
  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", run, { once: true })
  else run()
}

// --- shared helpers ------------------------------------------------------------

let idCounter = 0

/** Returns a document-unique id with the given prefix, like React's useId. */
export function uid(prefix = "radix") {
  idCounter += 1
  return `${prefix}-${idCounter.toString(36)}`
}

/**
 * Dispatches a cancelable CustomEvent and reports whether it was not
 * cancelled. Components announce state changes this way
 * (e.g. "shadcn:open-change") so applications can observe or veto them.
 */
export function emit(el, type, detail) {
  return el.dispatchEvent(new CustomEvent(type, { bubbles: true, cancelable: true, detail }))
}

/** Returns the closest ancestor (or self) whose data-slot is slot. */
export function closestSlot(el, slot) {
  return el.closest(`[data-slot="${slot}"]`)
}

/** Returns the descendants of root whose data-slot is slot, not crossing into nested roots with the same rootSlot. */
export function slots(root, slot, rootSlot) {
  const all = [...root.querySelectorAll(`[data-slot="${slot}"]`)]
  if (!rootSlot) return all
  return all.filter((el) => el.parentElement?.closest(`[data-slot="${rootSlot}"]`) === root)
}

/** Sets or removes an attribute: null/undefined/false remove it, true sets "". */
export function attr(el, name, value) {
  if (value === null || value === undefined || value === false) el.removeAttribute(name)
  else el.setAttribute(name, value === true ? "" : String(value))
}
