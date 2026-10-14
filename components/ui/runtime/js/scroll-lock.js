// scroll-lock.js: port of react-remove-scroll 2.7.1 (with
// react-remove-scroll-bar 2.3.8 and react-style-singleton 2.2.3), the
// <RemoveScroll> Radix wraps around modal layers (Dialog overlay, modal
// Popover / Menu / Select content).
//
// While a lock is held:
//   - <body> gets data-scroll-locked="<count>" and a <style> in <head>
//     (identical CSS, injected once for all locks, with the page's CSP
//     nonce) hides the page scrollbar and compensates its width;
//   - wheel and touch scrolling is blocked everywhere except inside the lock
//     node and its shards, where it is allowed only while the scrolled
//     element can actually scroll (no scroll chaining to the page);
//   - pinch-zoom is blocked unless allowPinchZoom.
//
// API
//
//   import { lockScroll } from "./scroll-lock.js"
//
//   const release = lockScroll(lockNode, {
//     shards: [],             // other elements (or () => element) where
//                             // scrolling stays allowed (Radix Dialog passes
//                             // its content, a sibling of the overlay)
//     allowPinchZoom: false,
//     removeScrollBar: true,  // hide the page scrollbar (body styles)
//     noIsolation: false,     // true: do not block events outside the lock
//     inert: false,           // true: pointer-events none outside the lock
//     gapMode: "margin",      // "margin" | "padding": how to fill the gap
//     noRelative: false,      // true: do not make body position: relative
//   })
//   release()
//
// Call it when Radix mounts the RemoveScroll (for a Dialog: when the
// overlay mounts) and release when it unmounts.

import { getNonce } from "./style.js"

// --- react-style-singleton -----------------------------------------------------------

function stylesheetSingleton() {
  let counter = 0
  let stylesheet = null
  return {
    add(style) {
      if (counter == 0) {
        stylesheet = document.createElement("style")
        stylesheet.type = "text/css"
        const nonce = getNonce()
        if (nonce) stylesheet.setAttribute("nonce", nonce)
        stylesheet.appendChild(document.createTextNode(style))
        ;(document.head || document.getElementsByTagName("head")[0]).appendChild(stylesheet)
      }
      counter++
    },
    remove() {
      counter--
      if (!counter && stylesheet) {
        stylesheet.parentNode && stylesheet.parentNode.removeChild(stylesheet)
        stylesheet = null
      }
    },
  }
}

// --- react-remove-scroll-bar -----------------------------------------------------------

export const zeroRightClassName = "right-scroll-bar-position"
export const fullWidthClassName = "width-before-scroll-bar"
export const noScrollbarsClassName = "with-scroll-bars-hidden"
export const removedBarSizeVariable = "--removed-body-scroll-bar-size"
export const lockAttribute = "data-scroll-locked"

const barStyle = stylesheetSingleton()

const parse = (x) => parseInt(x || "", 10) || 0
const getOffset = (gapMode) => {
  const cs = window.getComputedStyle(document.body)
  const left = cs[gapMode === "padding" ? "paddingLeft" : "marginLeft"]
  const top = cs[gapMode === "padding" ? "paddingTop" : "marginTop"]
  const right = cs[gapMode === "padding" ? "paddingRight" : "marginRight"]
  return [parse(left), parse(top), parse(right)]
}

/** Measures the page scrollbar gap (react-remove-scroll-bar getGapWidth). */
export function getGapWidth(gapMode = "margin") {
  const offsets = getOffset(gapMode)
  const documentWidth = document.documentElement.clientWidth
  const windowWidth = window.innerWidth
  return {
    left: offsets[0],
    top: offsets[1],
    right: offsets[2],
    gap: Math.max(0, windowWidth - documentWidth + offsets[2] - offsets[0]),
  }
}

const getStyles = ({ left, top, right, gap }, allowRelative, gapMode = "margin", important) =>
  `
  .${noScrollbarsClassName} {
   overflow: hidden ${important};
   padding-right: ${gap}px ${important};
  }
  body[${lockAttribute}] {
    overflow: hidden ${important};
    overscroll-behavior: contain;
    ${[
      allowRelative && `position: relative ${important};`,
      gapMode === "margin" &&
        `
    padding-left: ${left}px;
    padding-top: ${top}px;
    padding-right: ${right}px;
    margin-left:0;
    margin-top:0;
    margin-right: ${gap}px ${important};
    `,
      gapMode === "padding" && `padding-right: ${gap}px ${important};`,
    ]
      .filter(Boolean)
      .join("")}
  }

  .${zeroRightClassName} {
    right: ${gap}px ${important};
  }

  .${fullWidthClassName} {
    margin-right: ${gap}px ${important};
  }

  .${zeroRightClassName} .${zeroRightClassName} {
    right: 0 ${important};
  }

  .${fullWidthClassName} .${fullWidthClassName} {
    margin-right: 0 ${important};
  }

  body[${lockAttribute}] {
    ${removedBarSizeVariable}: ${gap}px;
  }
`

const getCurrentUseCounter = () => {
  const counter = parseInt(document.body.getAttribute(lockAttribute) || "0", 10)
  return isFinite(counter) ? counter : 0
}

function removeScrollBar({ noRelative, noImportant, gapMode = "margin" } = {}) {
  // useLockAttribute and the Style singleton are both effects of
  // <RemoveScrollBar>: the attribute first, then the style.
  document.body.setAttribute(lockAttribute, (getCurrentUseCounter() + 1).toString())
  const gap = getGapWidth(gapMode)
  barStyle.add(getStyles(gap, !noRelative, gapMode, !noImportant ? "!important" : ""))
  return () => {
    barStyle.remove()
    const newCounter = getCurrentUseCounter() - 1
    if (newCounter <= 0) document.body.removeAttribute(lockAttribute)
    else document.body.setAttribute(lockAttribute, newCounter.toString())
  }
}

// --- handleScroll ---------------------------------------------------------------------

const alwaysContainsScroll = (node) => node.tagName === "TEXTAREA"

const elementCanBeScrolled = (node, overflow) => {
  if (!(node instanceof Element)) return false
  const styles = window.getComputedStyle(node)
  return (
    styles[overflow] !== "hidden" &&
    !(styles.overflowY === styles.overflowX && !alwaysContainsScroll(node) && styles[overflow] === "visible")
  )
}
const elementCouldBeVScrolled = (node) => elementCanBeScrolled(node, "overflowY")
const elementCouldBeHScrolled = (node) => elementCanBeScrolled(node, "overflowX")

export const locationCouldBeScrolled = (axis, node) => {
  const ownerDocument = node.ownerDocument
  let current = node
  do {
    if (typeof ShadowRoot !== "undefined" && current instanceof ShadowRoot) current = current.host
    const isScrollable = elementCouldBeScrolled(axis, current)
    if (isScrollable) {
      const [, scrollHeight, clientHeight] = getScrollVariables(axis, current)
      if (scrollHeight > clientHeight) return true
    }
    current = current.parentNode
  } while (current && current !== ownerDocument.body)
  return false
}

const getVScrollVariables = ({ scrollTop, scrollHeight, clientHeight }) => [scrollTop, scrollHeight, clientHeight]
const getHScrollVariables = ({ scrollLeft, scrollWidth, clientWidth }) => [scrollLeft, scrollWidth, clientWidth]
const elementCouldBeScrolled = (axis, node) => (axis === "v" ? elementCouldBeVScrolled(node) : elementCouldBeHScrolled(node))
const getScrollVariables = (axis, node) => (axis === "v" ? getVScrollVariables(node) : getHScrollVariables(node))
const getDirectionFactor = (axis, direction) => (axis === "h" && direction === "rtl" ? -1 : 1)

export const handleScroll = (axis, endTarget, event, sourceDelta, noOverscroll) => {
  const directionFactor = getDirectionFactor(axis, window.getComputedStyle(endTarget).direction)
  const delta = directionFactor * sourceDelta
  let target = event.target
  const targetInLock = endTarget.contains(target)
  let shouldCancelScroll = false
  const isDeltaPositive = delta > 0
  let availableScroll = 0
  let availableScrollTop = 0
  do {
    if (!target) break
    const [position, scroll, capacity] = getScrollVariables(axis, target)
    const elementScroll = scroll - capacity - directionFactor * position
    if (position || elementScroll) {
      if (elementCouldBeScrolled(axis, target)) {
        availableScroll += elementScroll
        availableScrollTop += position
      }
    }
    const parent = target.parentNode
    target = parent && parent.nodeType === Node.DOCUMENT_FRAGMENT_NODE ? parent.host : parent
  } while (
    (!targetInLock && target !== document.body) ||
    (targetInLock && (endTarget.contains(target) || endTarget === target))
  )
  if (isDeltaPositive && ((noOverscroll && Math.abs(availableScroll) < 1) || (!noOverscroll && delta > availableScroll))) {
    shouldCancelScroll = true
  } else if (
    !isDeltaPositive &&
    ((noOverscroll && Math.abs(availableScrollTop) < 1) || (!noOverscroll && -delta > availableScrollTop))
  ) {
    shouldCancelScroll = true
  }
  return shouldCancelScroll
}

// --- RemoveScroll side car ---------------------------------------------------------------

const nonPassive = { passive: false }
const getTouchXY = (event) => ("changedTouches" in event ? [event.changedTouches[0].clientX, event.changedTouches[0].clientY] : [0, 0])
const getDeltaXY = (event) => [event.deltaX, event.deltaY]
const extractRef = (ref) => (typeof ref === "function" ? ref() : ref && "current" in ref ? ref.current : ref)
const deltaCompare = (x, y) => x[0] === y[0] && x[1] === y[1]
const generateStyle = (id) => `
  .block-interactivity-${id} {pointer-events: none;}
  .allow-interactivity-${id} {pointer-events: all;}
`

let idCounter = 0
let lockStack = []

function getOutermostShadowParent(node) {
  let shadowParent = null
  while (node !== null) {
    if (node instanceof ShadowRoot) {
      shadowParent = node.host
      node = node.host
    }
    node = node.parentNode
  }
  return shadowParent
}

/**
 * Locks page scrolling outside lockNode (and shards) until the returned
 * function is called.
 * @param {Element} lockNode
 * @param {{ shards?: (Element | (() => Element))[], allowPinchZoom?: boolean, removeScrollBar?: boolean, noIsolation?: boolean, inert?: boolean, gapMode?: "margin" | "padding", noRelative?: boolean }} [props]
 */
export function lockScroll(lockNode, props = {}) {
  const { shards = [], allowPinchZoom = false, removeScrollBar: removeBar = true, noIsolation = false, inert = false, gapMode, noRelative } = props
  const id = idCounter++
  const token = {} // identity in lockStack (the Style instance upstream)
  let shouldPreventQueue = []
  let touchStart = [0, 0]
  let activeAxis
  const cleanups = []

  const shouldCancelEvent = (event, parent) => {
    if (("touches" in event && event.touches.length === 2) || (event.type === "wheel" && event.ctrlKey)) {
      return !allowPinchZoom
    }
    const touch = getTouchXY(event)
    const deltaX = "deltaX" in event ? event.deltaX : touchStart[0] - touch[0]
    const deltaY = "deltaY" in event ? event.deltaY : touchStart[1] - touch[1]
    let currentAxis
    const target = event.target
    const moveDirection = Math.abs(deltaX) > Math.abs(deltaY) ? "h" : "v"
    if ("touches" in event && moveDirection === "h" && target.type === "range") return false
    let canBeScrolledInMainDirection = locationCouldBeScrolled(moveDirection, target)
    if (!canBeScrolledInMainDirection) return true
    if (canBeScrolledInMainDirection) {
      currentAxis = moveDirection
    } else {
      currentAxis = moveDirection === "v" ? "h" : "v"
      canBeScrolledInMainDirection = locationCouldBeScrolled(moveDirection, target)
    }
    if (!canBeScrolledInMainDirection) return false
    if (!activeAxis && "changedTouches" in event && (deltaX || deltaY)) activeAxis = currentAxis
    if (!currentAxis) return true
    const cancelingAxis = activeAxis || currentAxis
    return handleScroll(cancelingAxis, parent, event, cancelingAxis === "h" ? deltaX : deltaY, true)
  }

  const shouldPrevent = (event) => {
    if (!lockStack.length || lockStack[lockStack.length - 1] !== token) return
    const delta = "deltaY" in event ? getDeltaXY(event) : getTouchXY(event)
    const sourceEvent = shouldPreventQueue.filter(
      (e) => e.name === event.type && (e.target === event.target || event.target === e.shadowParent) && deltaCompare(e.delta, delta),
    )[0]
    if (sourceEvent && sourceEvent.should) {
      if (event.cancelable) event.preventDefault()
      return
    }
    if (!sourceEvent) {
      const shardNodes = shards
        .map(extractRef)
        .filter(Boolean)
        .filter((node) => node.contains(event.target))
      const shouldStop = shardNodes.length > 0 ? shouldCancelEvent(event, shardNodes[0]) : !noIsolation
      if (shouldStop && event.cancelable) event.preventDefault()
    }
  }

  const shouldCancel = (name, delta, target, should) => {
    const event = { name, delta, target, should, shadowParent: getOutermostShadowParent(target) }
    shouldPreventQueue.push(event)
    setTimeout(() => {
      shouldPreventQueue = shouldPreventQueue.filter((e) => e !== event)
    }, 1)
  }

  const scrollTouchStart = (event) => {
    touchStart = getTouchXY(event)
    activeAxis = undefined
  }
  const scrollWheel = (event) => {
    shouldCancel(event.type, getDeltaXY(event), event.target, shouldCancelEvent(event, lockNode))
  }
  const scrollTouchMove = (event) => {
    shouldCancel(event.type, getTouchXY(event), event.target, shouldCancelEvent(event, lockNode))
  }

  // RemoveScroll's container props: capture handlers on the lock node.
  lockNode.addEventListener("scroll", scrollWheel, true)
  lockNode.addEventListener("wheel", scrollWheel, true)
  lockNode.addEventListener("touchmove", scrollTouchMove, true)
  cleanups.push(() => {
    lockNode.removeEventListener("scroll", scrollWheel, true)
    lockNode.removeEventListener("wheel", scrollWheel, true)
    lockNode.removeEventListener("touchmove", scrollTouchMove, true)
  })

  // Side car effects, in upstream order: inert classes, then the scroll
  // bar removal (a child component, so its effects ran first upstream),
  // then the lock stack and document listeners.
  let interactivityStyle = null
  if (inert) {
    interactivityStyle = stylesheetSingleton()
    interactivityStyle.add(generateStyle(id))
  }
  if (removeBar) cleanups.push(removeScrollBar({ noRelative, gapMode }))
  if (inert) {
    document.body.classList.add(`block-interactivity-${id}`)
    const allow = [lockNode, ...shards.map(extractRef)].filter(Boolean)
    allow.forEach((el) => el.classList.add(`allow-interactivity-${id}`))
    cleanups.push(() => {
      document.body.classList.remove(`block-interactivity-${id}`)
      allow.forEach((el) => el.classList.remove(`allow-interactivity-${id}`))
      interactivityStyle.remove()
    })
  }
  lockStack.push(token)
  document.addEventListener("wheel", shouldPrevent, nonPassive)
  document.addEventListener("touchmove", shouldPrevent, nonPassive)
  document.addEventListener("touchstart", scrollTouchStart, nonPassive)
  cleanups.push(() => {
    lockStack = lockStack.filter((inst) => inst !== token)
    document.removeEventListener("wheel", shouldPrevent, nonPassive)
    document.removeEventListener("touchmove", shouldPrevent, nonPassive)
    document.removeEventListener("touchstart", scrollTouchStart, nonPassive)
  })

  let released = false
  return () => {
    if (released) return
    released = true
    for (const c of cleanups.reverse()) c()
  }
}
