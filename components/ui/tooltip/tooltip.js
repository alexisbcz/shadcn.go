// tooltip.js: Radix Tooltip behavior (@radix-ui/react-tooltip 1.2.8) for
// the server-rendered markup of components/ui/tooltip.
//
// Markup contract (rendered by the Go package):
//   - <template data-radix-tooltip data-content-id="<id>"
//       data-provider-id="<provider>" data-delay-duration="<ms>"
//       data-skip-delay-duration="<ms>" [data-disable-hoverable-content]
//       [data-default-open]> holds the popper wrapper (see popper.js) with
//     the content as React renders it, before it is positioned. Inside it,
//     the visually hidden copy (span[role=tooltip]) carries
//     data-radix-tooltip-copy when it must be filled with a copy of the
//     content's children (Radix renders them twice).
//   - The trigger carries data-radix-tooltip-trigger="<id>", removed on
//     mount.
//
// Tooltips sharing a data-provider-id share TooltipProvider's state: the
// open delay is skipped while another tooltip was open less than
// skipDelayDuration ago.
//
// Scripts can observe and drive a tooltip:
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getTooltip } from "/_shadcn/tooltip.js"
//   getTooltip(contentIdOrTemplate)?.setOpen(true)

import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createPopper, openOnLoad } from "./popper.js"

const TOOLTIP_OPEN = "tooltip.open"
const tooltips = new Map() // content id -> controller
const providers = new Map() // provider id -> provider state

/** Returns the controller of a tooltip by content id or template. */
export function getTooltip(ref) {
  if (typeof ref === "string") return tooltips.get(ref)
  if (ref instanceof Element) {
    const id = ref.getAttribute("data-content-id")
    return id ? tooltips.get(id) : undefined
  }
  return undefined
}

// TooltipProvider.
function getProvider(id, skipDelayDuration) {
  let p = providers.get(id)
  if (!p) {
    p = {
      isOpenDelayed: true,
      isPointerInTransit: false,
      skipDelayTimer: 0,
      skipDelayDuration,
      onOpen() {
        window.clearTimeout(p.skipDelayTimer)
        p.isOpenDelayed = false
      },
      onClose() {
        window.clearTimeout(p.skipDelayTimer)
        p.skipDelayTimer = window.setTimeout(() => (p.isOpenDelayed = true), p.skipDelayDuration)
      },
    }
    providers.set(id, p)
  }
  return p
}


// Finds the elements carrying marker attribute name with value, removes the
// attribute and returns them (other components' markers stay in place).
function takeOwnMarked(name, value) {
  return [...document.querySelectorAll(`[${name}="${CSS.escape(value)}"]`)].map((el) => {
    el.removeAttribute(name)
    return el
  })
}

function createTooltip(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const provider = getProvider(tpl.getAttribute("data-provider-id") ?? "default", Number(tpl.getAttribute("data-skip-delay-duration") ?? 300))
  const delayDuration = Number(tpl.getAttribute("data-delay-duration") ?? 0)
  const disableHoverableContent = tpl.hasAttribute("data-disable-hoverable-content")

  const trigger = takeOwnMarked("data-radix-tooltip-trigger", contentId)[0] ?? null

  let open = tpl.hasAttribute("data-default-open")
  let wasOpenDelayed = false
  let stateAttribute = open ? "instant-open" : "closed"
  let openTimer = 0
  let mounted = null

  const renderState = () => {
    if (trigger) {
      if (open) trigger.setAttribute("aria-describedby", contentId)
      else trigger.removeAttribute("aria-describedby")
      trigger.setAttribute("data-state", stateAttribute)
    }
    if (mounted) mounted.content.setAttribute("data-state", stateAttribute)
  }

  // useControllableState: the value changes now, onChange runs after the
  // render (in an effect), before a freshly mounted content's effects.
  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    stateAttribute = open ? (wasOpenDelayed ? "delayed-open" : "instant-open") : "closed"
    renderState()
    if (open) {
      if (mounted) {
        mounted.presence.setPresent(true)
        onChange(true)
      } else {
        onChange(true)
        if (open && !mounted) mountContent()
      }
    } else {
      if (mounted) mounted.presence.setPresent(false)
      onChange(false)
    }
  }

  function onChange(value) {
    if (value) {
      provider.onOpen()
      document.dispatchEvent(new CustomEvent(TOOLTIP_OPEN))
    } else {
      provider.onClose()
    }
  }

  const handleOpen = () => {
    window.clearTimeout(openTimer)
    openTimer = 0
    wasOpenDelayed = false
    setOpen(true)
  }
  const handleClose = () => {
    window.clearTimeout(openTimer)
    openTimer = 0
    setOpen(false)
  }
  const handleDelayedOpen = () => {
    window.clearTimeout(openTimer)
    openTimer = window.setTimeout(() => {
      wasOpenDelayed = true
      setOpen(true)
      openTimer = 0
    }, delayDuration)
  }
  const onTriggerEnter = () => {
    if (provider.isOpenDelayed) handleDelayedOpen()
    else handleOpen()
  }
  const onTriggerLeave = () => {
    if (disableHoverableContent) {
      handleClose()
    } else {
      window.clearTimeout(openTimer)
      openTimer = 0
    }
  }

  // TooltipTrigger handlers.
  let isPointerDown = false
  let hasPointerMoveOpened = false
  const handlePointerUp = () => (isPointerDown = false)
  const triggerHandlers = {
    pointermove(event) {
      if (event.pointerType === "touch") return
      if (!hasPointerMoveOpened && !provider.isPointerInTransit) {
        onTriggerEnter()
        hasPointerMoveOpened = true
      }
    },
    pointerleave() {
      onTriggerLeave()
      hasPointerMoveOpened = false
    },
    pointerdown() {
      if (open) handleClose()
      isPointerDown = true
      document.addEventListener("pointerup", handlePointerUp, { once: true })
    },
    // React's onFocus/onBlur are focusin/focusout: focus moving into a
    // child of the trigger (a button inside an addon) opens it too.
    focusin() {
      if (!isPointerDown) handleOpen()
    },
    focusout() {
      handleClose()
    },
    click() {
      handleClose()
    },
  }
  const wrapHandler = (fn) => (event) => {
    if (!event.defaultPrevented) fn(event)
  }
  const boundHandlers = Object.entries(triggerHandlers).map(([type, fn]) => [type, wrapHandler(fn)])
  if (trigger) for (const [type, fn] of boundHandlers) trigger.addEventListener(type, fn)

  function mountContent() {
    const nodes = instantiate(tpl)
    const wrapper = nodes[0]
    const content = wrapper.firstElementChild
    // VisuallyHidden copy of the children (the Arrow renders nothing there).
    const copy = takeMarked(content, "data-radix-tooltip-copy")[0]?.el
    if (copy) {
      for (const child of [...content.childNodes]) {
        if (child === copy) continue
        if (child.nodeType === 1 && child.matches("[data-radix-popper-arrow]")) continue
        copy.appendChild(child.cloneNode(true))
      }
    }
    const state = { wrapper, content, cleanups: [], presence: null }
    mounted = state
    content.setAttribute("data-state", stateAttribute)
    const unmount = mount(nodes)

    // Effects, children first: PopperContent, DismissableLayer,
    // TooltipContentImpl, TooltipContentHoverable.
    const popper = createPopper(wrapper, { anchor: trigger })
    state.cleanups.push(() => popper.destroy())

    const layer = createDismissableLayer(content, {
      disableOutsidePointerEvents: false,
      parent: trigger,
      onFocusOutside: (event) => event.preventDefault(),
      onDismiss: handleClose,
    })
    state.cleanups.push(() => layer.destroy())

    document.addEventListener(TOOLTIP_OPEN, handleClose)
    state.cleanups.push(() => document.removeEventListener(TOOLTIP_OPEN, handleClose))
    if (trigger) {
      const handleScroll = (event) => {
        const target = event.target
        if (target?.contains?.(trigger)) handleClose()
      }
      window.addEventListener("scroll", handleScroll, { capture: true })
      state.cleanups.push(() => window.removeEventListener("scroll", handleScroll, { capture: true }))
    }

    if (!disableHoverableContent) state.cleanups.push(hoverable(content))

    const presence = createPresence(content, {
      present: true,
      onExitComplete: () => {
        presence.destroy()
        unmount()
        for (const c of state.cleanups.reverse()) c()
        if (mounted === state) mounted = null
      },
    })
    state.presence = presence
  }

  // TooltipContentHoverable: a grace area lets the pointer travel between
  // the trigger and the content.
  function hoverable(content) {
    let graceArea = null
    const handleTrackPointerGrace = (event) => {
      const target = event.target
      const pointerPosition = { x: event.clientX, y: event.clientY }
      const hasEnteredTarget = trigger?.contains(target) || content?.contains(target)
      const isPointerOutsideGraceArea = !isPointInPolygon(pointerPosition, activeArea)
      if (hasEnteredTarget) {
        handleRemoveGraceArea()
      } else if (isPointerOutsideGraceArea) {
        handleRemoveGraceArea()
        handleClose()
      }
    }
    // The grace area is React state: the document listener follows it in an
    // effect, after a render scheduled outside the event that created it, so
    // the pointermove dispatched with that pointerleave is not tracked.
    let listening = false
    let activeArea = null
    let syncTimer = 0
    const syncListener = () => {
      syncTimer = 0
      activeArea = graceArea
      const want = !!graceArea && !disposed
      if (want && !listening) document.addEventListener("pointermove", handleTrackPointerGrace)
      if (!want && listening) document.removeEventListener("pointermove", handleTrackPointerGrace)
      listening = want
    }
    let disposed = false
    const setGraceArea = (area) => {
      graceArea = area
      if (!syncTimer) syncTimer = window.setTimeout(syncListener, 0)
    }
    const handleRemoveGraceArea = () => {
      setGraceArea(null)
      provider.isPointerInTransit = false
    }
    const handleCreateGraceArea = (event, hoverTarget) => {
      const currentTarget = event.currentTarget
      const exitPoint = { x: event.clientX, y: event.clientY }
      const exitSide = getExitSideFromRect(exitPoint, currentTarget.getBoundingClientRect())
      const paddedExitPoints = getPaddedExitPoints(exitPoint, exitSide)
      const hoverTargetPoints = getPointsFromRect(hoverTarget.getBoundingClientRect())
      setGraceArea(getHull([...paddedExitPoints, ...hoverTargetPoints]))
      provider.isPointerInTransit = true
    }
    const handleTriggerLeave = (event) => handleCreateGraceArea(event, content)
    const handleContentLeave = (event) => handleCreateGraceArea(event, trigger)
    trigger?.addEventListener("pointerleave", handleTriggerLeave)
    content.addEventListener("pointerleave", handleContentLeave)
    return () => {
      trigger?.removeEventListener("pointerleave", handleTriggerLeave)
      content.removeEventListener("pointerleave", handleContentLeave)
      handleRemoveGraceArea()
      disposed = true
      window.clearTimeout(syncTimer)
      syncListener()
    }
  }

  const controller = {
    get open() {
      return open
    },
    setOpen,
    dispose() {
      window.clearTimeout(openTimer)
      if (trigger) for (const [type, fn] of boundHandlers) trigger.removeEventListener(type, fn)
      if (mounted) {
        const state = mounted
        state.presence.destroy()
        state.wrapper.remove()
        for (const c of state.cleanups.reverse()) c()
        mounted = null
      }
      if (tooltips.get(contentId) === controller) tooltips.delete(contentId)
    },
  }
  tooltips.set(contentId, controller)
  if (open) openOnLoad(tpl, () => open && !mounted && mountContent())
  return controller
}

function getExitSideFromRect(point, rect) {
  const top = Math.abs(rect.top - point.y)
  const bottom = Math.abs(rect.bottom - point.y)
  const right = Math.abs(rect.right - point.x)
  const left = Math.abs(rect.left - point.x)
  switch (Math.min(top, bottom, right, left)) {
    case left:
      return "left"
    case right:
      return "right"
    case top:
      return "top"
    case bottom:
      return "bottom"
    default:
      throw new Error("unreachable")
  }
}

function getPaddedExitPoints(exitPoint, exitSide, padding = 5) {
  const paddedExitPoints = []
  switch (exitSide) {
    case "top":
      paddedExitPoints.push({ x: exitPoint.x - padding, y: exitPoint.y + padding }, { x: exitPoint.x + padding, y: exitPoint.y + padding })
      break
    case "bottom":
      paddedExitPoints.push({ x: exitPoint.x - padding, y: exitPoint.y - padding }, { x: exitPoint.x + padding, y: exitPoint.y - padding })
      break
    case "left":
      paddedExitPoints.push({ x: exitPoint.x + padding, y: exitPoint.y - padding }, { x: exitPoint.x + padding, y: exitPoint.y + padding })
      break
    case "right":
      paddedExitPoints.push({ x: exitPoint.x - padding, y: exitPoint.y - padding }, { x: exitPoint.x - padding, y: exitPoint.y + padding })
      break
  }
  return paddedExitPoints
}

function getPointsFromRect(rect) {
  const { top, right, bottom, left } = rect
  return [
    { x: left, y: top },
    { x: right, y: top },
    { x: right, y: bottom },
    { x: left, y: bottom },
  ]
}

function isPointInPolygon(point, polygon) {
  const { x, y } = point
  let inside = false
  for (let i = 0, j = polygon.length - 1; i < polygon.length; j = i++) {
    const ii = polygon[i]
    const jj = polygon[j]
    const xi = ii.x
    const yi = ii.y
    const xj = jj.x
    const yj = jj.y
    const intersect = yi > y !== yj > y && x < ((xj - xi) * (y - yi)) / (yj - yi) + xi
    if (intersect) inside = !inside
  }
  return inside
}

function getHull(points) {
  const newPoints = points.slice()
  newPoints.sort((a, b) => {
    if (a.x < b.x) return -1
    else if (a.x > b.x) return 1
    else if (a.y < b.y) return -1
    else if (a.y > b.y) return 1
    else return 0
  })
  return getHullPresorted(newPoints)
}

function getHullPresorted(points) {
  if (points.length <= 1) return points.slice()
  const upperHull = []
  for (let i = 0; i < points.length; i++) {
    const p = points[i]
    while (upperHull.length >= 2) {
      const q = upperHull[upperHull.length - 1]
      const r = upperHull[upperHull.length - 2]
      if ((q.x - r.x) * (p.y - r.y) >= (q.y - r.y) * (p.x - r.x)) upperHull.pop()
      else break
    }
    upperHull.push(p)
  }
  upperHull.pop()
  const lowerHull = []
  for (let i = points.length - 1; i >= 0; i--) {
    const p = points[i]
    while (lowerHull.length >= 2) {
      const q = lowerHull[lowerHull.length - 1]
      const r = lowerHull[lowerHull.length - 2]
      if ((q.x - r.x) * (p.y - r.y) >= (q.y - r.y) * (p.x - r.x)) lowerHull.pop()
      else break
    }
    lowerHull.push(p)
  }
  lowerHull.pop()
  if (upperHull.length === 1 && lowerHull.length === 1 && upperHull[0].x === lowerHull[0].x && upperHull[0].y === lowerHull[0].y) {
    return upperHull
  } else {
    return upperHull.concat(lowerHull)
  }
}

define("tooltip", {
  selector: "template[data-radix-tooltip]",
  mount(tpl) {
    const t = createTooltip(tpl)
    return () => t.dispose()
  },
})
