// drawer.js: vaul (1.1.2) for the server-rendered markup of
// components/ui/drawer. vaul wraps Radix Dialog (@radix-ui/react-dialog
// 1.1.15): this is the Radix Dialog behavior (as in dialog.js, built on the
// same runtime primitives) with vaul's on top.
//
// Markup contract (rendered by the Go package):
//   - <template data-vaul-root data-content-id="<id>" data-direction="bottom"
//       [data-modal="false"] [data-default-open] [data-dismissible="false"]
//       [data-scale-background] [data-no-background-color] [data-no-body-styles]
//       [data-handle-only] [data-auto-focus] [data-close-threshold="0.25"]
//       [data-snap-points='["0.5","355px"]'] [data-fade-from-index="1"]>
//     holds the overlay and the content exactly as React renders them when
//     open; they are cloned into <body> on open.
//   - Triggers are the elements with aria-controls="<id>".
//   - Markers removed on mount: data-radix-dialog-overlay,
//     data-radix-dialog-close="<id>".
//
// Scripts can observe and drive a drawer:
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getDrawer } from "/_shadcn/drawer.js"
//   getDrawer(contentIdOrTemplateOrTrigger)?.setOpen(true)   // .open, .toggle()

import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createFocusScope } from "./focus-scope.js"
import { acquireFocusGuards } from "./focus-guards.js"
import { lockScroll } from "./scroll-lock.js"
import { hideOthers } from "./hide-others.js"
import { injectStyle } from "./style.js"

// vaul's injected stylesheet, verbatim.
const VAUL_CSS = "[data-vaul-drawer]{touch-action:none;will-change:transform;transition:transform .5s cubic-bezier(.32, .72, 0, 1);animation-duration:.5s;animation-timing-function:cubic-bezier(0.32,0.72,0,1)}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=bottom][data-state=open]{animation-name:slideFromBottom}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=bottom][data-state=closed]{animation-name:slideToBottom}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=top][data-state=open]{animation-name:slideFromTop}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=top][data-state=closed]{animation-name:slideToTop}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=left][data-state=open]{animation-name:slideFromLeft}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=left][data-state=closed]{animation-name:slideToLeft}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=right][data-state=open]{animation-name:slideFromRight}[data-vaul-drawer][data-vaul-snap-points=false][data-vaul-drawer-direction=right][data-state=closed]{animation-name:slideToRight}[data-vaul-drawer][data-vaul-snap-points=true][data-vaul-drawer-direction=bottom]{transform:translate3d(0,var(--initial-transform,100%),0)}[data-vaul-drawer][data-vaul-snap-points=true][data-vaul-drawer-direction=top]{transform:translate3d(0,calc(var(--initial-transform,100%) * -1),0)}[data-vaul-drawer][data-vaul-snap-points=true][data-vaul-drawer-direction=left]{transform:translate3d(calc(var(--initial-transform,100%) * -1),0,0)}[data-vaul-drawer][data-vaul-snap-points=true][data-vaul-drawer-direction=right]{transform:translate3d(var(--initial-transform,100%),0,0)}[data-vaul-drawer][data-vaul-delayed-snap-points=true][data-vaul-drawer-direction=top]{transform:translate3d(0,var(--snap-point-height,0),0)}[data-vaul-drawer][data-vaul-delayed-snap-points=true][data-vaul-drawer-direction=bottom]{transform:translate3d(0,var(--snap-point-height,0),0)}[data-vaul-drawer][data-vaul-delayed-snap-points=true][data-vaul-drawer-direction=left]{transform:translate3d(var(--snap-point-height,0),0,0)}[data-vaul-drawer][data-vaul-delayed-snap-points=true][data-vaul-drawer-direction=right]{transform:translate3d(var(--snap-point-height,0),0,0)}[data-vaul-overlay][data-vaul-snap-points=false]{animation-duration:.5s;animation-timing-function:cubic-bezier(0.32,0.72,0,1)}[data-vaul-overlay][data-vaul-snap-points=false][data-state=open]{animation-name:fadeIn}[data-vaul-overlay][data-state=closed]{animation-name:fadeOut}[data-vaul-animate=false]{animation:none!important}[data-vaul-overlay][data-vaul-snap-points=true]{opacity:0;transition:opacity .5s cubic-bezier(.32, .72, 0, 1)}[data-vaul-overlay][data-vaul-snap-points=true]{opacity:1}[data-vaul-drawer]:not([data-vaul-custom-container=true])::after{content:'';position:absolute;background:inherit;background-color:inherit}[data-vaul-drawer][data-vaul-drawer-direction=top]::after{top:initial;bottom:100%;left:0;right:0;height:200%}[data-vaul-drawer][data-vaul-drawer-direction=bottom]::after{top:100%;bottom:initial;left:0;right:0;height:200%}[data-vaul-drawer][data-vaul-drawer-direction=left]::after{left:initial;right:100%;top:0;bottom:0;width:200%}[data-vaul-drawer][data-vaul-drawer-direction=right]::after{left:100%;right:initial;top:0;bottom:0;width:200%}[data-vaul-overlay][data-vaul-snap-points=true]:not([data-vaul-snap-points-overlay=true]):not(\n[data-state=closed]\n){opacity:0}[data-vaul-overlay][data-vaul-snap-points-overlay=true]{opacity:1}[data-vaul-handle]{display:block;position:relative;opacity:.7;background:#e2e2e4;margin-left:auto;margin-right:auto;height:5px;width:32px;border-radius:1rem;touch-action:pan-y}[data-vaul-handle]:active,[data-vaul-handle]:hover{opacity:1}[data-vaul-handle-hitarea]{position:absolute;left:50%;top:50%;transform:translate(-50%,-50%);width:max(100%,2.75rem);height:max(100%,2.75rem);touch-action:inherit}@media (hover:hover) and (pointer:fine){[data-vaul-drawer]{user-select:none}}@media (pointer:fine){[data-vaul-handle-hitarea]:{width:100%;height:100%}}@keyframes fadeIn{from{opacity:0}to{opacity:1}}@keyframes fadeOut{to{opacity:0}}@keyframes slideFromBottom{from{transform:translate3d(0,var(--initial-transform,100%),0)}to{transform:translate3d(0,0,0)}}@keyframes slideToBottom{to{transform:translate3d(0,var(--initial-transform,100%),0)}}@keyframes slideFromTop{from{transform:translate3d(0,calc(var(--initial-transform,100%) * -1),0)}to{transform:translate3d(0,0,0)}}@keyframes slideToTop{to{transform:translate3d(0,calc(var(--initial-transform,100%) * -1),0)}}@keyframes slideFromLeft{from{transform:translate3d(calc(var(--initial-transform,100%) * -1),0,0)}to{transform:translate3d(0,0,0)}}@keyframes slideToLeft{to{transform:translate3d(calc(var(--initial-transform,100%) * -1),0,0)}}@keyframes slideFromRight{from{transform:translate3d(var(--initial-transform,100%),0,0)}to{transform:translate3d(0,0,0)}}@keyframes slideToRight{to{transform:translate3d(var(--initial-transform,100%),0,0)}}"

// vaul injects its stylesheet when the module loads; here, when the first
// drawer mounts, so pages without drawers do not get its global keyframes.
let cssInjected = false
function ensureCSS() {
  if (cssInjected) return
  cssInjected = true
  injectStyle(VAUL_CSS)
}

const TRANSITIONS = { DURATION: 0.5, EASE: [0.32, 0.72, 0, 1] }
const EASE = `cubic-bezier(${TRANSITIONS.EASE.join(",")})`
const VELOCITY_THRESHOLD = 0.4
const SCROLL_LOCK_TIMEOUT = 100
const BORDER_RADIUS = 8
const WINDOW_TOP_OFFSET = 26
const DRAG_CLASS = "vaul-dragging"
const LONG_HANDLE_PRESS_TIMEOUT = 250
const DOUBLE_TAP_TIMEOUT = 120

// --- vaul helpers ------------------------------------------------------------------------------

const styleCache = new WeakMap()
function set(el, styles, ignoreCache = false) {
  if (!el || !(el instanceof HTMLElement)) return
  const originalStyles = {}
  Object.entries(styles).forEach(([key, value]) => {
    if (key.startsWith("--")) {
      el.style.setProperty(key, value)
      return
    }
    originalStyles[key] = el.style[key]
    el.style[key] = value
  })
  if (ignoreCache) return
  styleCache.set(el, originalStyles)
}
function reset(el, prop) {
  if (!el || !(el instanceof HTMLElement)) return
  const originalStyles = styleCache.get(el)
  if (!originalStyles) return
  el.style[prop] = originalStyles[prop]
}
const isVertical = (direction) => direction === "top" || direction === "bottom"
function getTranslate(element, direction) {
  if (!element) return null
  const style = window.getComputedStyle(element)
  const transform = style.transform || style.webkitTransform || style.mozTransform
  let mat = transform.match(/^matrix3d\((.+)\)$/)
  if (mat) return parseFloat(mat[1].split(", ")[isVertical(direction) ? 13 : 12])
  mat = transform.match(/^matrix\((.+)\)$/)
  return mat ? parseFloat(mat[1].split(", ")[isVertical(direction) ? 5 : 4]) : null
}
const dampenValue = (v) => 8 * (Math.log(v + 1) - 2)
function assignStyle(element, style) {
  if (!element) return () => {}
  const prevStyle = element.style.cssText
  Object.assign(element.style, style)
  return () => {
    element.style.cssText = prevStyle
  }
}
const getScale = () => (window.innerWidth - WINDOW_TOP_OFFSET) / window.innerWidth

// --- drawers -----------------------------------------------------------------------------------

const drawers = new Map() // content id -> controller

/** Returns the controller of a drawer by content id, template or trigger. */
export function getDrawer(ref) {
  if (typeof ref === "string") return drawers.get(ref)
  if (ref instanceof Element) {
    const id = ref.getAttribute("data-content-id") ?? ref.getAttribute("aria-controls")
    return id ? drawers.get(id) : undefined
  }
  return undefined
}

const getState = (open) => (open ? "open" : "closed")

function createDrawer(tpl) {
  ensureCSS()
  const contentId = tpl.getAttribute("data-content-id")
  const modal = tpl.getAttribute("data-modal") !== "false"
  const direction = tpl.getAttribute("data-direction") || "bottom"
  const dismissible = tpl.getAttribute("data-dismissible") !== "false"
  const shouldScaleBackground = tpl.hasAttribute("data-scale-background")
  const setBackgroundColorOnScale = !tpl.hasAttribute("data-no-background-color")
  const noBodyStyles = tpl.hasAttribute("data-no-body-styles")
  const handleOnly = tpl.hasAttribute("data-handle-only")
  const autoFocus = tpl.hasAttribute("data-auto-focus")
  const closeThreshold = Number(tpl.getAttribute("data-close-threshold") ?? 0.25)
  const snapPoints = tpl.hasAttribute("data-snap-points")
    ? JSON.parse(tpl.getAttribute("data-snap-points")).map((p) => (/px$/.test(p) ? p : Number(p)))
    : undefined
  const hasSnapPoints = !!snapPoints && snapPoints.length > 0
  const fadeFromIndex = tpl.hasAttribute("data-fade-from-index") ? Number(tpl.getAttribute("data-fade-from-index")) : snapPoints && snapPoints.length - 1
  const defaultOpen = tpl.hasAttribute("data-default-open")

  let open = false // Radix/vaul isOpen
  let mounted = null // { nodes, content, overlay, presences, cleanups, ... } while present

  // vaul Root state and refs.
  let isDragging = false
  let openTime = null
  let dragStartTime = null
  let dragEndTime = null
  let lastTimeDragPrevented = null
  let isAllowedToDrag = false
  let pointerStart = 0
  let shouldAnimate = !defaultOpen
  let drawerHeight = 0
  let drawerWidth = 0
  let activeSnapPoint = snapPoints?.[0]
  // drawerRef.current as vaul's Root last rendered it. useSnapPoints'
  // snapToPoint is a useCallback depending on drawerRef.current, and an
  // effect snaps to the active point whenever its identity changes: on the
  // first Root render after the content mounted (the next open/close or
  // drag start, as the content mounts in the Portal's own render).
  let renderedDrawer = null
  let closeSnapTimer = null
  let scaleCleanup = null
  let bgTimer = null
  const initialBackgroundColor = document.body.style.backgroundColor

  const triggers = () => [...document.querySelectorAll(`[aria-controls="${CSS.escape(contentId)}"]`)]
  const drawerEl = () => mounted?.content ?? null
  const overlayEl = () => mounted?.overlay ?? null

  requestAnimationFrame(() => {
    shouldAnimate = true
  })
  // Non-modal: Radix's dialog still disables outside pointer events; vaul
  // turns them back on.
  if (!modal) requestAnimationFrame(() => (document.body.style.pointerEvents = "auto"))

  // Close buttons outside the portalled content.
  const outsideCloses = [...document.querySelectorAll(`[data-radix-dialog-close="${CSS.escape(contentId)}"]`)].map((el) => {
    el.removeAttribute("data-radix-dialog-close")
    el.removeAttribute("data-radix-dialog-cancel")
    return el
  })
  const onCloseClick = (event) => {
    if (!event.defaultPrevented) onOpenChange(false)
  }
  outsideCloses.forEach((el) => el.addEventListener("click", onCloseClick))

  // --- snap points (useSnapPoints) --------------------------------------------------------------

  const snapPointsOffset = () => {
    if (!snapPoints) return []
    const size = { width: window.innerWidth, height: window.innerHeight }
    return snapPoints.map((snapPoint) => {
      const isPx = typeof snapPoint === "string"
      const n = isPx ? parseInt(snapPoint, 10) : 0
      if (isVertical(direction)) {
        const height = isPx ? n : snapPoint * size.height
        return direction === "bottom" ? size.height - height : -size.height + height
      }
      const width = isPx ? n : snapPoint * size.width
      return direction === "right" ? size.width - width : -size.width + width
    })
  }
  const activeSnapPointIndex = () => (snapPoints ? snapPoints.findIndex((p) => p === activeSnapPoint) : null)
  const isLastSnapPoint = () => activeSnapPoint === snapPoints?.[snapPoints.length - 1] || null
  const shouldFade = () =>
    (snapPoints && snapPoints.length > 0 && (fadeFromIndex || fadeFromIndex === 0) && !Number.isNaN(fadeFromIndex) && snapPoints[fadeFromIndex] === activeSnapPoint) || !snapPoints
  const activeSnapPointOffset = () => {
    const i = activeSnapPointIndex()
    return i !== null && i !== -1 ? snapPointsOffset()[i] : null
  }

  function snapToPoint(dimension) {
    const offsets = snapPointsOffset()
    const newIndex = offsets.findIndex((d) => d === dimension)
    if (snapPoints && newIndex === offsets.length - 1) openTime = new Date()
    set(drawerEl(), {
      transition: `transform ${TRANSITIONS.DURATION}s ${EASE}`,
      transform: isVertical(direction) ? `translate3d(0, ${dimension}px, 0)` : `translate3d(${dimension}px, 0, 0)`,
    })
    if (offsets && newIndex !== offsets.length - 1 && fadeFromIndex !== undefined && newIndex !== fadeFromIndex && newIndex < fadeFromIndex) {
      set(overlayEl(), { transition: `opacity ${TRANSITIONS.DURATION}s ${EASE}`, opacity: "0" })
    } else {
      set(overlayEl(), { transition: `opacity ${TRANSITIONS.DURATION}s ${EASE}`, opacity: "1" })
    }
    setActiveSnapPoint(snapPoints?.[Math.max(newIndex, 0)])
  }

  function setActiveSnapPoint(p) {
    if (p === activeSnapPoint) return
    activeSnapPoint = p
    renderedDrawer = drawerEl()
    renderSnap()
    // Effect: snap to the new active point.
    const i = activeSnapPointIndex()
    const offsets = snapPointsOffset()
    if (i !== -1 && typeof offsets[i] === "number") snapToPoint(offsets[i])
  }

  // A render of vaul's Root that saw drawerRef.current = drawerAtRender,
  // then its useSnapPoints effect.
  function rootRender(drawerAtRender) {
    if (drawerAtRender === renderedDrawer) return
    renderedDrawer = drawerAtRender
    if (!activeSnapPoint) return
    const i = activeSnapPointIndex()
    const offsets = snapPointsOffset()
    if (i !== -1 && typeof offsets[i] === "number") snapToPoint(offsets[i])
  }

  function renderSnap() {
    if (!mounted) return
    const c = mounted.content
    const o = mounted.overlay
    if (c && hasSnapPoints) c.style.setProperty("--snap-point-height", `${snapPointsOffset()[activeSnapPointIndex() ?? 0]}px`)
    if (o) o.setAttribute("data-vaul-snap-points-overlay", String(open && shouldFade()))
  }

  function onReleaseSnapPoints({ draggedDistance, velocity }) {
    if (fadeFromIndex === undefined) return
    const offsets = snapPointsOffset()
    const activeOffset = activeSnapPointOffset()
    const idx = activeSnapPointIndex()
    const currentPosition = direction === "bottom" || direction === "right" ? (activeOffset ?? 0) - draggedDistance : (activeOffset ?? 0) + draggedDistance
    const isOverlaySnapPoint = idx === fadeFromIndex - 1
    const isFirst = idx === 0
    const hasDraggedUp = draggedDistance > 0
    if (isOverlaySnapPoint) set(overlayEl(), { transition: `opacity ${TRANSITIONS.DURATION}s ${EASE}` })
    if (velocity > 2 && !hasDraggedUp) {
      if (dismissible) closeDrawer()
      else snapToPoint(offsets[0])
      return
    }
    if (velocity > 2 && hasDraggedUp && snapPoints) {
      snapToPoint(offsets[snapPoints.length - 1])
      return
    }
    const closest = offsets.reduce((prev, curr) => {
      if (typeof prev !== "number" || typeof curr !== "number") return prev
      return Math.abs(curr - currentPosition) < Math.abs(prev - currentPosition) ? curr : prev
    })
    const dim = isVertical(direction) ? window.innerHeight : window.innerWidth
    if (velocity > VELOCITY_THRESHOLD && Math.abs(draggedDistance) < dim * 0.4) {
      const dragDirection = hasDraggedUp ? 1 : -1
      if (dragDirection > 0 && isLastSnapPoint() && snapPoints) {
        snapToPoint(offsets[snapPoints.length - 1])
        return
      }
      if (isFirst && dragDirection < 0 && dismissible) closeDrawer()
      if (idx === null) return
      snapToPoint(offsets[idx + dragDirection])
      return
    }
    snapToPoint(closest)
  }

  function onDragSnapPoints({ draggedDistance }) {
    const activeOffset = activeSnapPointOffset()
    if (activeOffset === null) return
    const offsets = snapPointsOffset()
    const newValue = direction === "bottom" || direction === "right" ? activeOffset - draggedDistance : activeOffset + draggedDistance
    if ((direction === "bottom" || direction === "right") && newValue < offsets[offsets.length - 1]) return
    if ((direction === "top" || direction === "left") && newValue > offsets[offsets.length - 1]) return
    set(drawerEl(), { transform: isVertical(direction) ? `translate3d(0, ${newValue}px, 0)` : `translate3d(${newValue}px, 0, 0)` })
  }

  function getSnapPointsPercentageDragged(absDraggedDistance, isDraggingDown) {
    const idx = activeSnapPointIndex()
    const offsets = snapPointsOffset()
    if (!snapPoints || typeof idx !== "number" || !offsets || fadeFromIndex === undefined) return null
    const isOverlaySnapPoint = idx === fadeFromIndex - 1
    const isOverlaySnapPointOrHigher = idx >= fadeFromIndex
    if (isOverlaySnapPointOrHigher && isDraggingDown) return 0
    if (isOverlaySnapPoint && !isDraggingDown) return 1
    if (!shouldFade() && !isOverlaySnapPoint) return null
    const target = isOverlaySnapPoint ? idx + 1 : idx - 1
    const distance = isOverlaySnapPoint ? offsets[target] - offsets[target - 1] : offsets[target + 1] - offsets[target]
    const percentageDragged = absDraggedDistance / Math.abs(distance)
    return isOverlaySnapPoint ? 1 - percentageDragged : percentageDragged
  }

  // --- open state ------------------------------------------------------------------------------

  // Radix Dialog's onOpenChange, as vaul's Root passes it.
  function onOpenChange(next) {
    if (!dismissible && !next) return
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    if (!next) closeDrawer(true)
    setIsOpen(next)
  }

  function setIsOpen(next) {
    next = !!next
    if (next === open) return
    const drawerAtRender = drawerEl()
    open = next
    for (const t of triggers()) {
      t.setAttribute("aria-expanded", String(open))
      t.setAttribute("data-state", getState(open))
    }
    if (open) {
      if (mounted) reshow()
      else mountParts()
    } else if (mounted) {
      for (const el of mounted.stateful) el.setAttribute("data-state", "closed")
      mounted.overlay?.setAttribute("data-vaul-snap-points-overlay", "false")
      mounted.content?.setAttribute("data-vaul-snap-points", "false")
      mounted.overlay?.setAttribute("data-vaul-snap-points", "false")
      for (const h of mounted.handles) h.setAttribute("data-vaul-drawer-visible", "false")
      syncAnimateAttr()
      mounted.focusScope?.setTrapped(false)
      for (const p of mounted.presences) p.setPresent(false)
    }
    onIsOpenChange()
    rootRender(drawerAtRender)
  }

  // Effects of vaul's Root reacting to isOpen.
  function onIsOpenChange() {
    if (open) {
      // Trigger enter animation without using CSS animation
      set(document.documentElement, { scrollBehavior: "auto" })
      openTime = new Date()
    } else {
      reset(document.documentElement, "scrollBehavior")
    }
    scaleBackground()
    // useControllableState onChange (an effect after the commit).
    if (open && !modal) requestAnimationFrame(() => (document.body.style.pointerEvents = "auto"))
    if (!open) document.body.style.pointerEvents = "auto"
  }

  function syncAnimateAttr() {
    if (!mounted) return
    for (const el of mounted.stateful) el.setAttribute("data-vaul-animate", String(shouldAnimate))
  }

  // useScaleBackground
  function scaleBackground() {
    if (scaleCleanup) {
      const c = scaleCleanup
      scaleCleanup = null
      c()
    }
    if (!(open && shouldScaleBackground)) return
    if (bgTimer) clearTimeout(bgTimer)
    const wrapper = document.querySelector("[data-vaul-drawer-wrapper]") || document.querySelector("[vaul-drawer-wrapper]")
    if (!wrapper) return
    if (setBackgroundColorOnScale && !noBodyStyles) assignStyle(document.body, { background: "black" })
    assignStyle(wrapper, {
      transformOrigin: isVertical(direction) ? "top" : "left",
      transitionProperty: "transform, border-radius",
      transitionDuration: `${TRANSITIONS.DURATION}s`,
      transitionTimingFunction: EASE,
    })
    const wrapperStylesCleanup = assignStyle(wrapper, {
      borderRadius: `${BORDER_RADIUS}px`,
      overflow: "hidden",
      ...(isVertical(direction)
        ? { transform: `scale(${getScale()}) translate3d(0, calc(env(safe-area-inset-top) + 14px), 0)` }
        : { transform: `scale(${getScale()}) translate3d(calc(env(safe-area-inset-top) + 14px), 0, 0)` }),
    })
    scaleCleanup = () => {
      wrapperStylesCleanup()
      bgTimer = window.setTimeout(() => {
        if (initialBackgroundColor) document.body.style.background = initialBackgroundColor
        else document.body.style.removeProperty("background")
      }, TRANSITIONS.DURATION * 1000)
    }
  }

  function closeDrawer(fromWithin) {
    cancelDrag()
    if (!fromWithin) onOpenChangeFromVaul(false)
    clearTimeout(closeSnapTimer)
    closeSnapTimer = setTimeout(() => {
      if (snapPoints) setActiveSnapPoint(snapPoints[0])
    }, TRANSITIONS.DURATION * 1000)
  }

  // vaul's own setIsOpen(false) (drag to close, handle): Radix is controlled
  // by it, so the dialog closes without going through onOpenChange.
  function onOpenChangeFromVaul(next) {
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    setIsOpen(next)
  }

  function reshow() {
    for (const el of mounted.stateful) el.setAttribute("data-state", "open")
    mounted.overlay?.setAttribute("data-vaul-snap-points-overlay", String(shouldFade()))
    mounted.content?.setAttribute("data-vaul-snap-points", String(hasSnapPoints))
    mounted.overlay?.setAttribute("data-vaul-snap-points", String(hasSnapPoints))
    for (const h of mounted.handles) h.setAttribute("data-vaul-drawer-visible", "true")
    syncAnimateAttr()
    mounted.focusScope?.setTrapped(true)
    for (const p of mounted.presences) p.setPresent(true)
  }

  function mountParts() {
    const nodes = instantiate(tpl)
    const overlays = new Set(takeMarked(nodes, "data-radix-dialog-overlay").map((m) => m.el))
    const closes = takeMarked(nodes, "data-radix-dialog-close").map((m) => m.el)
    takeMarked(nodes, "data-radix-dialog-cancel")
    const content = nodes.flatMap((n) => (n.id === contentId ? [n] : [...n.querySelectorAll(`[id="${CSS.escape(contentId)}"]`)]))[0]
    const overlay = nodes.find((n) => overlays.has(n))
    const handles = nodes.flatMap((n) => [...n.querySelectorAll("[data-vaul-handle]")])
    const state = { nodes, content, overlay, handles, stateful: [overlay, content].filter(Boolean), presences: [], cleanups: new Map(), focusScope: null, live: new Set(nodes) }
    mounted = state
    // data-vaul-animate is shouldAnimate.current when React renders.
    for (const el of state.stateful) el.setAttribute("data-vaul-animate", String(shouldAnimate))

    closes.forEach((el) => el.addEventListener("click", onCloseClick))

    if (content && hasSnapPoints) {
      content.style.setProperty("--snap-point-height", `${snapPointsOffset()[activeSnapPointIndex() ?? 0]}px`)
    }

    mount(nodes)

    if (overlay) {
      const release = lockScroll(overlay, { allowPinchZoom: true, shards: [() => content] })
      overlay.addEventListener("mouseup", onRelease)
      state.cleanups.set(overlay, [release, () => overlay.removeEventListener("mouseup", onRelease)])
    }
    if (content) state.cleanups.set(content, mountContent(content, state))
    for (const h of handles) mountHandle(h, state)

    for (const node of nodes) {
      const presence = createPresence(node, {
        present: true,
        onExitComplete: () => unmountNode(state, node, presence),
      })
      state.presences.push(presence)
    }
  }

  function mountContent(content, state) {
    const cleanups = []
    const trigger = triggers().pop()

    const layer = createDismissableLayer(content, {
      // vaul does not pass `modal` to Radix: the dialog is always modal.
      disableOutsidePointerEvents: true,
      parent: trigger,
      onDismiss: () => onOpenChange(false),
      onPointerDownOutside: (event) => {
        // vaul's handler, then DialogContentModal's.
        if (!modal || event.defaultPrevented) {
          event.preventDefault()
          return
        }
        const originalEvent = event.detail.originalEvent
        const ctrlLeftClick = originalEvent.button === 0 && originalEvent.ctrlKey === true
        const isRightClick = originalEvent.button === 2 || ctrlLeftClick
        if (isRightClick) event.preventDefault()
      },
      onFocusOutside: (event) => {
        // vaul: !modal prevents; DialogContentModal prevents too.
        event.preventDefault()
      },
    })
    cleanups.push(() => layer.destroy())

    const focusScope = createFocusScope(content, {
      loop: true,
      trapped: open,
      onMountAutoFocus: (event) => {
        if (!autoFocus) event.preventDefault()
      },
      onUnmountAutoFocus: (event) => {
        event.preventDefault()
        triggers().pop()?.focus()
      },
    })
    state.focusScope = focusScope
    cleanups.push(() => focusScope.destroy())

    cleanups.push(acquireFocusGuards())
    cleanups.push(hideOthers(content))

    // vaul's Drawer.Content pointer handlers.
    let pointerStartRef = null
    let lastKnownPointerEvent = null
    let wasBeyondThePoint = false
    const isDeltaInDirection = (delta, threshold = 0) => {
      if (wasBeyondThePoint) return true
      const deltaY = Math.abs(delta.y)
      const deltaX = Math.abs(delta.x)
      const isDeltaX = deltaX > deltaY
      const dFactor = ["bottom", "right"].includes(direction) ? 1 : -1
      if (direction === "left" || direction === "right") {
        const isReverseDirection = delta.x * dFactor < 0
        if (!isReverseDirection && deltaX >= 0 && deltaX <= threshold) return isDeltaX
      } else {
        const isReverseDirection = delta.y * dFactor < 0
        if (!isReverseDirection && deltaY >= 0 && deltaY <= threshold) return !isDeltaX
      }
      wasBeyondThePoint = true
      return true
    }
    const handleOnPointerUp = (event) => {
      pointerStartRef = null
      wasBeyondThePoint = false
      onRelease(event)
    }
    const onPointerDown = (event) => {
      if (handleOnly) return
      pointerStartRef = { x: event.pageX, y: event.pageY }
      onPress(event)
    }
    const onPointerMove = (event) => {
      lastKnownPointerEvent = event
      if (handleOnly) return
      if (!pointerStartRef) return
      const yPosition = event.pageY - pointerStartRef.y
      const xPosition = event.pageX - pointerStartRef.x
      const swipeStartThreshold = event.pointerType === "touch" ? 10 : 2
      const delta = { x: xPosition, y: yPosition }
      if (isDeltaInDirection(delta, swipeStartThreshold)) onDrag(event)
      else if (Math.abs(xPosition) > swipeStartThreshold || Math.abs(yPosition) > swipeStartThreshold) pointerStartRef = null
    }
    const onPointerUp = (event) => {
      pointerStartRef = null
      wasBeyondThePoint = false
      onRelease(event)
    }
    const onPointerOut = () => handleOnPointerUp(lastKnownPointerEvent)
    const onContextMenu = () => {
      if (lastKnownPointerEvent) handleOnPointerUp(lastKnownPointerEvent)
    }
    content.addEventListener("pointerdown", onPointerDown)
    content.addEventListener("pointermove", onPointerMove)
    content.addEventListener("pointerup", onPointerUp)
    content.addEventListener("pointerout", onPointerOut)
    content.addEventListener("contextmenu", onContextMenu)
    cleanups.push(() => {
      content.removeEventListener("pointerdown", onPointerDown)
      content.removeEventListener("pointermove", onPointerMove)
      content.removeEventListener("pointerup", onPointerUp)
      content.removeEventListener("pointerout", onPointerOut)
      content.removeEventListener("contextmenu", onContextMenu)
    })

    // Delayed snap points (transition instead of animation).
    if (hasSnapPoints) {
      requestAnimationFrame(() => {
        if (mounted === state) content.setAttribute("data-vaul-delayed-snap-points", "true")
      })
    }

    // Visual viewport (mobile keyboards): keep the drawer above it.
    let initialDrawerHeight = 0
    let previousDiffFromInitial = 0
    let keyboardIsOpen = false
    const onVisualViewportChange = () => {
      const focusedElement = document.activeElement
      if (!(isInput(focusedElement) || keyboardIsOpen)) return
      const visualViewportHeight = window.visualViewport?.height || 0
      const totalHeight = window.innerHeight
      let diffFromInitial = totalHeight - visualViewportHeight
      const dh = content.getBoundingClientRect().height || 0
      const isTallEnough = dh > totalHeight * 0.8
      if (!initialDrawerHeight) initialDrawerHeight = dh
      const offsetFromTop = content.getBoundingClientRect().top
      if (Math.abs(previousDiffFromInitial - diffFromInitial) > 60) keyboardIsOpen = !keyboardIsOpen
      const idx = activeSnapPointIndex()
      if (snapPoints && snapPoints.length > 0 && idx) diffFromInitial += snapPointsOffset()[idx] || 0
      previousDiffFromInitial = diffFromInitial
      if (dh > visualViewportHeight || keyboardIsOpen) {
        const height = content.getBoundingClientRect().height
        let newDrawerHeight = height
        if (height > visualViewportHeight) newDrawerHeight = visualViewportHeight - (isTallEnough ? offsetFromTop : WINDOW_TOP_OFFSET)
        content.style.height = `${Math.max(newDrawerHeight, visualViewportHeight - offsetFromTop)}px`
      } else {
        content.style.height = `${initialDrawerHeight}px`
      }
      if (snapPoints && snapPoints.length > 0 && !keyboardIsOpen) content.style.bottom = "0px"
      else content.style.bottom = `${Math.max(diffFromInitial, 0)}px`
    }
    window.visualViewport?.addEventListener("resize", onVisualViewportChange)
    cleanups.push(() => window.visualViewport?.removeEventListener("resize", onVisualViewportChange))

    return cleanups.reverse()
  }

  function mountHandle(handle, state) {
    let closeTimeoutId = null
    let shouldCancelInteraction = false
    const handleCancelInteraction = () => {
      if (closeTimeoutId) window.clearTimeout(closeTimeoutId)
      shouldCancelInteraction = false
    }
    const handleCycleSnapPoints = () => {
      if (isDragging || handle.hasAttribute("data-prevent-cycle") || shouldCancelInteraction) {
        handleCancelInteraction()
        return
      }
      handleCancelInteraction()
      if (!snapPoints || snapPoints.length === 0) {
        if (!dismissible) closeDrawer()
        return
      }
      const isLast = activeSnapPoint === snapPoints[snapPoints.length - 1]
      if (isLast && dismissible) {
        closeDrawer()
        return
      }
      const currentSnapIndex = snapPoints.findIndex((p) => p === activeSnapPoint)
      if (currentSnapIndex === -1) return
      setActiveSnapPoint(snapPoints[currentSnapIndex + 1])
    }
    const onClick = () => {
      if (shouldCancelInteraction) {
        handleCancelInteraction()
        return
      }
      window.setTimeout(handleCycleSnapPoints, DOUBLE_TAP_TIMEOUT)
    }
    const onPointerDown = (e) => {
      if (handleOnly) onPress(e)
      closeTimeoutId = window.setTimeout(() => (shouldCancelInteraction = true), LONG_HANDLE_PRESS_TIMEOUT)
    }
    const onPointerMove = (e) => {
      if (handleOnly) onDrag(e)
    }
    handle.addEventListener("click", onClick)
    handle.addEventListener("pointercancel", handleCancelInteraction)
    handle.addEventListener("pointerdown", onPointerDown)
    handle.addEventListener("pointermove", onPointerMove)
    void state
  }

  // --- dragging (vaul Root) --------------------------------------------------------------------

  function onPress(event) {
    if (!dismissible && !snapPoints) return
    const drawer = drawerEl()
    if (drawer && !drawer.contains(event.target)) return
    drawerHeight = drawer?.getBoundingClientRect().height || 0
    drawerWidth = drawer?.getBoundingClientRect().width || 0
    isDragging = true
    dragStartTime = new Date()
    event.target.setPointerCapture(event.pointerId)
    pointerStart = isVertical(direction) ? event.pageY : event.pageX
    rootRender(drawer) // setIsDragging(true)
  }

  function shouldDrag(el, isDraggingInDirection) {
    let element = el
    const highlightedText = window.getSelection()?.toString()
    const drawer = drawerEl()
    const swipeAmount = drawer ? getTranslate(drawer, direction) : null
    const date = new Date()
    if (element.tagName === "SELECT") return false
    if (element.hasAttribute("data-vaul-no-drag") || element.closest("[data-vaul-no-drag]")) return false
    if (direction === "right" || direction === "left") return true
    if (openTime && date.getTime() - openTime.getTime() < 500) return false
    if (swipeAmount !== null) {
      if (direction === "bottom" ? swipeAmount > 0 : swipeAmount < 0) return true
    }
    if (highlightedText && highlightedText.length > 0) return false
    if (lastTimeDragPrevented && date.getTime() - lastTimeDragPrevented.getTime() < SCROLL_LOCK_TIMEOUT && swipeAmount === 0) {
      lastTimeDragPrevented = date
      return false
    }
    if (isDraggingInDirection) {
      lastTimeDragPrevented = date
      return false
    }
    while (element) {
      if (element.scrollHeight > element.clientHeight) {
        if (element.scrollTop !== 0) {
          lastTimeDragPrevented = new Date()
          return false
        }
        if (element.getAttribute("role") === "dialog") return true
      }
      element = element.parentNode
    }
    return true
  }

  function onDrag(event) {
    const drawer = drawerEl()
    if (!drawer) return
    if (!isDragging) return
    const overlay = overlayEl()
    const directionMultiplier = direction === "bottom" || direction === "right" ? 1 : -1
    const draggedDistance = (pointerStart - (isVertical(direction) ? event.pageY : event.pageX)) * directionMultiplier
    const isDraggingInDirection = draggedDistance > 0
    const noCloseSnapPointsPreCondition = snapPoints && !dismissible && !isDraggingInDirection
    if (noCloseSnapPointsPreCondition && activeSnapPointIndex() === 0) return
    const absDraggedDistance = Math.abs(draggedDistance)
    const wrapper = document.querySelector("[data-vaul-drawer-wrapper]")
    const drawerDimension = direction === "bottom" || direction === "top" ? drawerHeight : drawerWidth
    let percentageDragged = absDraggedDistance / drawerDimension
    const snapPointPercentageDragged = getSnapPointsPercentageDragged(absDraggedDistance, isDraggingInDirection)
    if (snapPointPercentageDragged !== null) percentageDragged = snapPointPercentageDragged
    if (noCloseSnapPointsPreCondition && percentageDragged >= 1) return
    if (!isAllowedToDrag && !shouldDrag(event.target, isDraggingInDirection)) return
    drawer.classList.add(DRAG_CLASS)
    isAllowedToDrag = true
    set(drawer, { transition: "none" })
    set(overlay, { transition: "none" })
    if (snapPoints) onDragSnapPoints({ draggedDistance })
    if (isDraggingInDirection && !snapPoints) {
      const dampenedDraggedDistance = dampenValue(draggedDistance)
      const translateValue = Math.min(dampenedDraggedDistance * -1, 0) * directionMultiplier
      set(drawer, { transform: isVertical(direction) ? `translate3d(0, ${translateValue}px, 0)` : `translate3d(${translateValue}px, 0, 0)` })
      return
    }
    const opacityValue = 1 - percentageDragged
    const idx = activeSnapPointIndex()
    if (shouldFade() || (fadeFromIndex && idx === fadeFromIndex - 1)) {
      set(overlay, { opacity: `${opacityValue}`, transition: "none" }, true)
    }
    if (wrapper && overlay && shouldScaleBackground) {
      const scaleValue = Math.min(getScale() + percentageDragged * (1 - getScale()), 1)
      const borderRadiusValue = 8 - percentageDragged * 8
      const translateValue = Math.max(0, 14 - percentageDragged * 14)
      set(
        wrapper,
        {
          borderRadius: `${borderRadiusValue}px`,
          transform: isVertical(direction) ? `scale(${scaleValue}) translate3d(0, ${translateValue}px, 0)` : `scale(${scaleValue}) translate3d(${translateValue}px, 0, 0)`,
          transition: "none",
        },
        true,
      )
    }
    if (!snapPoints) {
      const translateValue = absDraggedDistance * directionMultiplier
      set(drawer, { transform: isVertical(direction) ? `translate3d(0, ${translateValue}px, 0)` : `translate3d(${translateValue}px, 0, 0)` })
    }
  }

  function resetDrawer() {
    const drawer = drawerEl()
    if (!drawer) return
    const wrapper = document.querySelector("[data-vaul-drawer-wrapper]")
    const currentSwipeAmount = getTranslate(drawer, direction)
    set(drawer, { transform: "translate3d(0, 0, 0)", transition: `transform ${TRANSITIONS.DURATION}s ${EASE}` })
    set(overlayEl(), { transition: `opacity ${TRANSITIONS.DURATION}s ${EASE}`, opacity: "1" })
    if (shouldScaleBackground && currentSwipeAmount && currentSwipeAmount > 0 && open) {
      set(
        wrapper,
        {
          borderRadius: `${BORDER_RADIUS}px`,
          overflow: "hidden",
          ...(isVertical(direction)
            ? { transform: `scale(${getScale()}) translate3d(0, calc(env(safe-area-inset-top) + 14px), 0)`, transformOrigin: "top" }
            : { transform: `scale(${getScale()}) translate3d(calc(env(safe-area-inset-top) + 14px), 0, 0)`, transformOrigin: "left" }),
          transitionProperty: "transform, border-radius",
          transitionDuration: `${TRANSITIONS.DURATION}s`,
          transitionTimingFunction: EASE,
        },
        true,
      )
    }
  }

  function cancelDrag() {
    const drawer = drawerEl()
    if (!isDragging || !drawer) return
    drawer.classList.remove(DRAG_CLASS)
    isAllowedToDrag = false
    isDragging = false
    dragEndTime = new Date()
  }

  function onRelease(event) {
    const drawer = drawerEl()
    if (!isDragging || !drawer) return
    drawer.classList.remove(DRAG_CLASS)
    isAllowedToDrag = false
    isDragging = false
    dragEndTime = new Date()
    const swipeAmount = getTranslate(drawer, direction)
    if (!event || !shouldDrag(event.target, false) || !swipeAmount || Number.isNaN(swipeAmount)) return
    if (dragStartTime === null) return
    const timeTaken = dragEndTime.getTime() - dragStartTime.getTime()
    const distMoved = pointerStart - (isVertical(direction) ? event.pageY : event.pageX)
    const velocity = Math.abs(distMoved) / timeTaken
    if (snapPoints) {
      const directionMultiplier = direction === "bottom" || direction === "right" ? 1 : -1
      onReleaseSnapPoints({ draggedDistance: distMoved * directionMultiplier, velocity })
      return
    }
    if (direction === "bottom" || direction === "right" ? distMoved > 0 : distMoved < 0) {
      resetDrawer()
      return
    }
    if (velocity > VELOCITY_THRESHOLD) {
      closeDrawer()
      return
    }
    const visibleDrawerHeight = Math.min(drawer.getBoundingClientRect().height ?? 0, window.innerHeight)
    const visibleDrawerWidth = Math.min(drawer.getBoundingClientRect().width ?? 0, window.innerWidth)
    const isHorizontalSwipe = direction === "left" || direction === "right"
    if (Math.abs(swipeAmount) >= (isHorizontalSwipe ? visibleDrawerWidth : visibleDrawerHeight) * closeThreshold) {
      closeDrawer()
      return
    }
    resetDrawer()
  }

  function unmountNode(state, node, presence) {
    presence.destroy()
    node.remove()
    state.live.delete(node)
    for (const c of state.cleanups.get(node) ?? []) c()
    state.cleanups.delete(node)
    if (state.live.size === 0 && mounted === state) mounted = null
  }

  const controller = {
    get open() {
      return open
    },
    setOpen: (next) => onOpenChange(!!next),
    toggle: () => onOpenChange(!open),
    dispose() {
      outsideCloses.forEach((el) => el.removeEventListener("click", onCloseClick))
      if (mounted) {
        const state = mounted
        for (const [i, node] of state.nodes.entries()) if (state.live.has(node)) unmountNode(state, node, state.presences[i])
      }
      if (drawers.get(contentId) === controller) drawers.delete(contentId)
    },
  }
  drawers.set(contentId, controller)
  if (defaultOpen) setIsOpen(true)
  return controller
}

const nonTextInputTypes = new Set(["checkbox", "radio", "range", "color", "file", "image", "button", "submit", "reset"])
function isInput(target) {
  return (target instanceof HTMLInputElement && !nonTextInputTypes.has(target.type)) || target instanceof HTMLTextAreaElement || (target instanceof HTMLElement && target.isContentEditable)
}

// DialogTrigger: onClick toggles (composed after the user's handlers).
document.addEventListener("click", (event) => {
  if (event.defaultPrevented) return
  const t = event.target instanceof Element ? event.target.closest("[aria-controls]") : null
  if (!t) return
  const d = drawers.get(t.getAttribute("aria-controls"))
  if (d) d.toggle()
})

define("drawer", {
  selector: "template[data-vaul-root]",
  mount(tpl) {
    const d = createDrawer(tpl)
    return () => d.dispose()
  },
})
