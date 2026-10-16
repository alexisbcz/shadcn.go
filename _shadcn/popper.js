// popper.js: port of Radix Popper (@radix-ui/react-popper 1.2.8:
// PopperContent, PopperAnchor, PopperArrow) on the vendored floating-ui.js
// (@floating-ui/dom 1.7.4, driven like @floating-ui/react-dom 2.1.6's
// useFloating). Popover, Tooltip, HoverCard, Select (popper position),
// DropdownMenu, ContextMenu, Menubar, NavigationMenu... content is
// positioned by it.
//
// Markup contract (server-rendered by components/ui/popper, exactly what
// React renders before the first position is known):
//
//   <div data-radix-popper-content-wrapper=""
//        style="position: fixed; left: 0px; top: 0px; transform: translate(0, -200%); min-width: max-content;"
//        data-radix-popper='{"side":"bottom","sideOffset":4,...}'>   <- options marker
//     <div data-side="bottom" data-align="center" style="...; animation: none" ...>  <- the content
//       ...
//       <span data-radix-popper-arrow="" style="...">  <- optional PopperArrow wrapper
//         <svg width="10" height="5" viewBox="0 0 30 10" preserveAspectRatio="none" style="display: block">
//           <polygon points="0,0 30,0 15,10"></polygon></svg></span>
//     </div>
//   </div>
//
// The markers (data-radix-popper, data-radix-popper-arrow) are removed when
// the popper is created, so the live DOM matches React's.
//
// API
//
//   import { createPopper } from "./popper.js"
//
//   const popper = createPopper(wrapper, {
//     anchor,            // the reference element (PopperAnchor: the trigger,
//                        // or a custom anchor); may be set later
//     onPlaced() {},     // first time the content is positioned
//     // Any option of the marker may be overridden here: side, sideOffset,
//     // align, alignOffset, arrowPadding, avoidCollisions,
//     // collisionBoundary (elements or CSS selectors), collisionPadding
//     // (number or {top,right,bottom,left}), sticky ("partial" | "always"),
//     // hideWhenDetached, updatePositionStrategy ("optimized" | "always").
//   })
//   popper.content          // the content element (wrapper's first child)
//   popper.isPositioned     // whether a position was computed
//   popper.placedSide / .placedAlign
//   popper.setAnchor(el)    // PopperAnchor changed
//   popper.update()         // recompute now (autoUpdate does it on scroll,
//                           // resize, layout shift and element resize)
//   popper.destroy()        // unmount: stops observing
//
//   import { openOnLoad } from "./popper.js"
//   openOnLoad(templateEl, () => setOpen(true))
//                           // defaultOpen: React mounts every portal open on
//                           // load in tree order, whatever the component;
//                           // callbacks queued in one task run in a
//                           // microtask, sorted by the document order of
//                           // their templates.
//
// Every update renders exactly what PopperContent renders: the wrapper's
// style (fixed position, translate, min-width, z-index copied from the
// content, --radix-popper-transform-origin, hidden when detached), the
// size middleware's --radix-popper-available-width/height and
// --radix-popper-anchor-width/height on the wrapper, the content's
// data-side / data-align and its `animation: none` until positioned (so
// enter animations start with the right side), and the arrow's style.

import { computePosition, autoUpdate, offset, shift, limitShift, flip, size, arrow as arrowMiddleware, hide } from "./floating-ui.js"
import { takeMarked } from "./portal.js"

export const SIDE_OPTIONS = ["top", "right", "bottom", "left"]
export const ALIGN_OPTIONS = ["start", "center", "end"]
export const OPTIONS_ATTR = "data-radix-popper"
export const ARROW_ATTR = "data-radix-popper-arrow"

const DEFAULTS = {
  side: "bottom",
  sideOffset: 0,
  align: "center",
  alignOffset: 0,
  arrowPadding: 0,
  avoidCollisions: true,
  collisionBoundary: [],
  collisionPadding: 0,
  sticky: "partial",
  hideWhenDetached: false,
  updatePositionStrategy: "optimized",
}

const OPPOSITE_SIDE = { top: "bottom", right: "left", bottom: "top", left: "right" }

/** Reads (and removes) the options marker of a wrapper. */
export function readOptions(wrapper) {
  const raw = wrapper.getAttribute(OPTIONS_ATTR)
  if (raw === null) return {}
  wrapper.removeAttribute(OPTIONS_ATTR)
  try {
    return JSON.parse(raw) ?? {}
  } catch {
    return {}
  }
}

/** Radix's getSideAndAlignFromPlacement. */
export function getSideAndAlignFromPlacement(placement) {
  const [side, align = "center"] = placement.split("-")
  return [side, align]
}

// React's style handling for a numeric value: px for lengths, as is for 0.
const px = (v) => (typeof v === "number" ? (v === 0 ? "0" : `${v}px`) : v)

// Applies a style object like React does on re-render: properties no longer
// present are cleared, undefined values are cleared.
function applyStyle(el, prev, next) {
  for (const k of Object.keys(prev)) {
    if (!(k in next) || next[k] === undefined) setStyle(el, k, "")
  }
  for (const [k, v] of Object.entries(next)) {
    if (v === undefined) continue
    if (prev[k] === v) continue
    setStyle(el, k, v)
  }
  return next
}

function setStyle(el, name, value) {
  if (name.startsWith("--")) el.style.setProperty(name, value === "" ? "" : String(value))
  else el.style[name] = value === "" ? "" : px(value)
}

function getDPR(element) {
  const win = element.ownerDocument.defaultView || window
  return win.devicePixelRatio || 1
}

function roundByDPR(element, value) {
  const dpr = getDPR(element)
  return Math.round(value * dpr) / dpr
}

// Radix's transformOrigin middleware.
const transformOrigin = (options) => ({
  name: "transformOrigin",
  options,
  fn(data) {
    const { placement, rects, middlewareData } = data
    const cannotCenterArrow = middlewareData.arrow?.centerOffset !== 0
    const isArrowHidden = cannotCenterArrow
    const arrowWidth = isArrowHidden ? 0 : options.arrowWidth
    const arrowHeight = isArrowHidden ? 0 : options.arrowHeight
    const [placedSide, placedAlign] = getSideAndAlignFromPlacement(placement)
    const noArrowAlign = { start: "0%", center: "50%", end: "100%" }[placedAlign]
    const arrowXCenter = (middlewareData.arrow?.x ?? 0) + arrowWidth / 2
    const arrowYCenter = (middlewareData.arrow?.y ?? 0) + arrowHeight / 2
    let x = ""
    let y = ""
    if (placedSide === "bottom") {
      x = isArrowHidden ? noArrowAlign : `${arrowXCenter}px`
      y = `${-arrowHeight}px`
    } else if (placedSide === "top") {
      x = isArrowHidden ? noArrowAlign : `${arrowXCenter}px`
      y = `${rects.floating.height + arrowHeight}px`
    } else if (placedSide === "right") {
      x = `${-arrowHeight}px`
      y = isArrowHidden ? noArrowAlign : `${arrowYCenter}px`
    } else if (placedSide === "left") {
      x = `${rects.floating.width + arrowHeight}px`
      y = isArrowHidden ? noArrowAlign : `${arrowYCenter}px`
    }
    return { data: { x, y } }
  },
})

function resolveBoundary(b) {
  if (typeof b === "string") return document.querySelector(b)
  return b ?? null
}

/**
 * Positions the content of a popper wrapper against an anchor.
 * @param {HTMLElement} wrapper the [data-radix-popper-content-wrapper] element
 * @param {object} [props] see the file header
 */
export function createPopper(wrapper, props = {}) {
  const { anchor: anchorProp, onPlaced, ...overrides } = props
  const opts = { ...DEFAULTS, ...readOptions(wrapper), ...overrides }
  const content = wrapper.firstElementChild
  const arrow = props.arrow ?? takeMarked(content, ARROW_ATTR)[0]?.el ?? null
  const arrowSvg = arrow?.firstElementChild ?? null

  let anchor = anchorProp ?? null
  let destroyed = false
  let cleanupAutoUpdate = null
  let arrowSize = arrow ? { width: arrow.offsetWidth, height: arrow.offsetHeight } : undefined
  let data = { x: 0, y: 0, placement: opts.side + (opts.align !== "center" ? "-" + opts.align : ""), middlewareData: {}, isPositioned: false }
  let prevWrapperStyle = {}
  let prevContentStyle = {}
  let prevArrowStyle = {}
  let placedNotified = false
  let generation = 0

  // useLayoutEffect: z-index copied from the content.
  const contentZIndex = window.getComputedStyle(content).zIndex

  const middleware = () => {
    const arrowWidth = arrowSize?.width ?? 0
    const arrowHeight = arrowSize?.height ?? 0
    const { sideOffset, alignOffset, arrowPadding, avoidCollisions, collisionBoundary, sticky, hideWhenDetached } = opts
    const collisionPaddingProp = opts.collisionPadding
    const collisionPadding =
      typeof collisionPaddingProp === "number" ? collisionPaddingProp : { top: 0, right: 0, bottom: 0, left: 0, ...collisionPaddingProp }
    const boundary = (Array.isArray(collisionBoundary) ? collisionBoundary : [collisionBoundary]).map(resolveBoundary)
    const hasExplicitBoundaries = boundary.length > 0
    const detectOverflowOptions = {
      padding: collisionPadding,
      boundary: boundary.filter((b) => b !== null),
      // with `strategy: 'fixed'`, this is the only way to get it to respect boundaries
      altBoundary: hasExplicitBoundaries,
    }
    return [
      offset({ mainAxis: sideOffset + arrowHeight, alignmentAxis: alignOffset }),
      avoidCollisions &&
        shift({
          mainAxis: true,
          crossAxis: false,
          limiter: sticky === "partial" ? limitShift() : void 0,
          ...detectOverflowOptions,
        }),
      avoidCollisions && flip({ ...detectOverflowOptions }),
      size({
        ...detectOverflowOptions,
        apply: ({ elements, rects, availableWidth, availableHeight }) => {
          const { width: anchorWidth, height: anchorHeight } = rects.reference
          const contentStyle = elements.floating.style
          contentStyle.setProperty("--radix-popper-available-width", `${availableWidth}px`)
          contentStyle.setProperty("--radix-popper-available-height", `${availableHeight}px`)
          contentStyle.setProperty("--radix-popper-anchor-width", `${anchorWidth}px`)
          contentStyle.setProperty("--radix-popper-anchor-height", `${anchorHeight}px`)
        },
      }),
      arrow && arrowMiddleware({ element: arrow, padding: arrowPadding }),
      transformOrigin({ arrowWidth, arrowHeight }),
      hideWhenDetached && hide({ strategy: "referenceHidden", ...detectOverflowOptions }),
    ]
  }

  const render = () => {
    const { middlewareData, isPositioned } = data
    const [placedSide, placedAlign] = getSideAndAlignFromPlacement(data.placement)
    // useFloating's floatingStyles (strategy fixed, transform: true).
    const x = roundByDPR(wrapper, data.x)
    const y = roundByDPR(wrapper, data.y)
    const floatingStyles = {
      position: "fixed",
      left: 0,
      top: 0,
      transform: "translate(" + x + "px, " + y + "px)",
      ...(getDPR(wrapper) >= 1.5 && { willChange: "transform" }),
    }
    prevWrapperStyle = applyStyle(wrapper, prevWrapperStyle, {
      ...floatingStyles,
      transform: isPositioned ? floatingStyles.transform : "translate(0, -200%)",
      minWidth: "max-content",
      zIndex: contentZIndex,
      "--radix-popper-transform-origin": [middlewareData.transformOrigin?.x, middlewareData.transformOrigin?.y].join(" "),
      ...(middlewareData.hide?.referenceHidden && { visibility: "hidden", pointerEvents: "none" }),
    })
    if (content.getAttribute("data-side") !== placedSide) content.setAttribute("data-side", placedSide)
    if (content.getAttribute("data-align") !== placedAlign) content.setAttribute("data-align", placedAlign)
    prevContentStyle = applyStyle(content, prevContentStyle, { animation: !isPositioned ? "none" : undefined })

    if (arrow) {
      const arrowX = middlewareData.arrow?.x
      const arrowY = middlewareData.arrow?.y
      const shouldHideArrow = middlewareData.arrow?.centerOffset !== 0
      const baseSide = OPPOSITE_SIDE[placedSide]
      prevArrowStyle = applyStyle(arrow, prevArrowStyle, {
        position: "absolute",
        left: arrowX,
        top: arrowY,
        [baseSide]: 0,
        transformOrigin: { top: "", right: "0 0", bottom: "center 0", left: "100% 0" }[placedSide],
        transform: {
          top: "translateY(100%)",
          right: "translateY(50%) rotate(90deg) translateX(-50%)",
          bottom: "rotate(180deg)",
          left: "translateY(50%) rotate(-90deg) translateX(50%)",
        }[placedSide],
        visibility: shouldHideArrow ? "hidden" : undefined,
      })
    }
  }

  // The server rendered the pre-position state; adopt it as React's
  // previous style so later renders clear what React would clear.
  const adopt = (el) => {
    const out = {}
    for (let i = 0; i < el.style.length; i++) {
      const name = el.style[i]
      const key = name.startsWith("--") ? name : name.replace(/-([a-z])/g, (_, c) => c.toUpperCase())
      out[key] = el.style.getPropertyValue(name)
    }
    return out
  }
  prevWrapperStyle = adopt(wrapper)
  prevContentStyle = { animation: content.style.animation || undefined }
  if (arrow) prevArrowStyle = adopt(arrow)
  render()

  const update = () => {
    if (destroyed || !anchor || !wrapper.isConnected) return
    const gen = generation
    computePosition(anchor, wrapper, {
      placement: opts.side + (opts.align !== "center" ? "-" + opts.align : ""),
      strategy: "fixed",
      middleware: middleware(),
    }).then((result) => {
      if (destroyed || gen !== generation) return
      data = { ...result, isPositioned: true }
      render()
      if (!placedNotified) {
        placedNotified = true
        onPlaced?.()
      }
    })
  }

  const start = () => {
    cleanupAutoUpdate?.()
    cleanupAutoUpdate = null
    if (destroyed || !anchor) return
    cleanupAutoUpdate = autoUpdate(anchor, wrapper, update, { animationFrame: opts.updatePositionStrategy === "always" })
  }

  // useSize(arrow): the arrow's size feeds the offset and arrow middleware.
  let arrowObserver = null
  if (arrow && typeof ResizeObserver === "function") {
    arrowObserver = new ResizeObserver((entries) => {
      if (!Array.isArray(entries) || !entries.length) return
      const entry = entries[0]
      let width
      let height
      if ("borderBoxSize" in entry) {
        const borderSizeEntry = entry.borderBoxSize
        const borderSize = Array.isArray(borderSizeEntry) ? borderSizeEntry[0] : borderSizeEntry
        width = borderSize.inlineSize
        height = borderSize.blockSize
      } else {
        width = arrow.offsetWidth
        height = arrow.offsetHeight
      }
      if (arrowSize && arrowSize.width === width && arrowSize.height === height) return
      arrowSize = { width, height }
      start()
    })
    arrowObserver.observe(arrow, { box: "border-box" })
  }

  start()

  return {
    wrapper,
    content,
    arrow,
    arrowSvg,
    get isPositioned() {
      return data.isPositioned
    },
    get placedSide() {
      return getSideAndAlignFromPlacement(data.placement)[0]
    },
    get placedAlign() {
      return getSideAndAlignFromPlacement(data.placement)[1]
    },
    get anchor() {
      return anchor
    },
    setAnchor(el) {
      if (el === anchor) return
      anchor = el
      start()
    },
    update,
    destroy() {
      if (destroyed) return
      destroyed = true
      generation++
      cleanupAutoUpdate?.()
      arrowObserver?.disconnect()
    },
  }
}

// --- defaultOpen ordering ---------------------------------------------------------------

let loadQueue = null

/**
 * Runs open (a component's "open now" for defaultOpen) after the current
 * task's other mounts, in document order of el among all queued elements.
 * @param {Element} el
 * @param {() => void} open
 */
export function openOnLoad(el, open) {
  if (!loadQueue) {
    loadQueue = []
    queueMicrotask(() => {
      const q = loadQueue
      loadQueue = null
      q.sort((a, b) => (a.el.compareDocumentPosition(b.el) & Node.DOCUMENT_POSITION_FOLLOWING ? -1 : 1))
      for (const item of q) item.open()
    })
  }
  loadQueue.push({ el, open })
}
