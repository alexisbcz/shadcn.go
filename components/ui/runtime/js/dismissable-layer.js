// dismissable-layer.js: port of Radix DismissableLayer
// (@radix-ui/react-dismissable-layer 1.1.11), with use-escape-keydown.
//
// A layer is a piece of UI that floats above the page and can be dismissed:
// Dialog, AlertDialog, Sheet, Popover, Menu, Select, Tooltip, HoverCard,
// NavigationMenu content. Layers stack in the order they open. The layer
// on top handles Escape; a pointerdown or a focus move outside a layer
// (and outside the layers opened from inside it) dismisses it; a layer
// opened with disableOutsidePointerEvents makes <body> pointer-events: none
// and every layer below it pointer-events: none, the ones above auto.
//
// API
//
//   import { createDismissableLayer, createDismissableLayerBranch, layerOf } from "./dismissable-layer.js"
//
//   const layer = createDismissableLayer(node, {
//     disableOutsidePointerEvents: false,
//     onEscapeKeyDown(event) {},       // KeyboardEvent; preventDefault() keeps the layer
//     onPointerDownOutside(event) {},  // CustomEvent "dismissableLayer.pointerDownOutside",
//                                      // event.detail.originalEvent is the PointerEvent;
//                                      // preventDefault() keeps the layer
//     onFocusOutside(event) {},        // CustomEvent "dismissableLayer.focusOutside",
//                                      // detail.originalEvent is the FocusEvent
//     onInteractOutside(event) {},     // called after either of the two above
//     onDismiss() {},                  // the layer should close
//     parent: element,                 // optional: an element of the layer this
//                                      // one was opened from (typically its
//                                      // trigger). Radix relies on the React
//                                      // tree for this: pointer and focus events
//                                      // inside a child layer count as inside
//                                      // every ancestor layer. Defaults to the
//                                      // topmost layer containing the node.
//   })
//   layer.update({ onDismiss, ... })  // replace callbacks (like new props)
//   layer.setDisableOutsidePointerEvents(on)
//                                     // like changing that prop while mounted
//                                     // (Radix menus pass `open`): React re-runs
//                                     // the effect, so turning it off restores
//                                     // body pointer-events (when this was the
//                                     // last such layer) but keeps the layer in
//                                     // the disabled set until it unmounts
//   layer.node                        // the element
//   layer.destroy()                   // unmount: removes the layer and
//                                     // restores body pointer-events
//
//   const removeBranch = createDismissableLayerBranch(node)
//                                     // DismissableLayer.Branch: events in
//                                     // node never count as outside any layer
//
//   layerOf(element)                  // the topmost layer containing element
//
// Call createDismissableLayer when Radix mounts the layer (after the node is
// in the DOM) and destroy when it unmounts (after its exit animation).

const CONTEXT_UPDATE = "dismissableLayer.update"
const POINTER_DOWN_OUTSIDE = "dismissableLayer.pointerDownOutside"
const FOCUS_OUTSIDE = "dismissableLayer.focusOutside"

let originalBodyPointerEvents

// DismissableLayerContext: module-level, like the React context default.
const context = {
  layers: new Set(), // layer records, in mount order
  layersWithOutsidePointerEventsDisabled: new Set(),
  branches: new Set(), // elements
}

/** Returns the topmost layer whose node contains element, or undefined. */
export function layerOf(element) {
  const layers = [...context.layers]
  for (let i = layers.length - 1; i >= 0; i--) {
    if (layers[i].node.contains(element)) return layers[i].handle
  }
  return undefined
}

function dispatchUpdate() {
  document.dispatchEvent(new CustomEvent(CONTEXT_UPDATE))
  updateStyles()
}

// Each layer's style.pointerEvents, as DismissableLayer renders it.
function updateStyles() {
  const layers = [...context.layers]
  const [highestDisabled] = [...context.layersWithOutsidePointerEventsDisabled].slice(-1)
  const highestDisabledIndex = layers.indexOf(highestDisabled)
  const isBodyPointerEventsDisabled = context.layersWithOutsidePointerEventsDisabled.size > 0
  layers.forEach((layer, index) => {
    const isPointerEventsEnabled = index >= highestDisabledIndex
    const value = isBodyPointerEventsDisabled ? (isPointerEventsEnabled ? "auto" : "none") : ""
    if (layer.node.style.pointerEvents !== value) layer.node.style.pointerEvents = value
  })
}

function handleAndDispatchCustomEvent(name, handler, detail) {
  const target = detail.originalEvent.target
  const event = new CustomEvent(name, { bubbles: false, cancelable: true, detail })
  if (handler) target.addEventListener(name, handler, { once: true })
  target.dispatchEvent(event)
}

/**
 * @param {HTMLElement} node
 * @param {object} [props] see the file header
 */
export function createDismissableLayer(node, props = {}) {
  const ownerDocument = node.ownerDocument ?? globalThis.document
  let opts = { ...props }
  let disableOutsidePointerEvents = !!opts.disableOutsidePointerEvents
  const parentLayer = opts.parent ? (opts.parent.node ? recordOf(opts.parent) : recordOfElement(opts.parent)) : recordOfElement(node.parentElement)

  const record = { node, parent: parentLayer, pointerInside: false, focusInside: false, handle: null }
  const ancestry = () => {
    const out = []
    for (let r = record; r; r = r.parent) out.push(r)
    return out
  }
  const cleanups = []

  // React's onPointerDownCapture / onFocusCapture / onBlurCapture travel the
  // React tree: an event in this layer (or a layer opened from it) sets the
  // flags of this layer and all its ancestors.
  const onPointerDownCapture = () => ancestry().forEach((r) => (r.pointerInside = true))
  const onFocusCapture = () => ancestry().forEach((r) => (r.focusInside = true))
  const onBlurCapture = () => ancestry().forEach((r) => (r.focusInside = false))
  node.addEventListener("pointerdown", onPointerDownCapture, true)
  node.addEventListener("focusin", onFocusCapture, true)
  node.addEventListener("focusout", onBlurCapture, true)
  cleanups.push(() => {
    node.removeEventListener("pointerdown", onPointerDownCapture, true)
    node.removeEventListener("focusin", onFocusCapture, true)
    node.removeEventListener("focusout", onBlurCapture, true)
  })

  const isPointerEventsEnabled = () => {
    const layers = [...context.layers]
    const [highest] = [...context.layersWithOutsidePointerEventsDisabled].slice(-1)
    return layers.indexOf(record) >= layers.indexOf(highest)
  }

  // usePointerDownOutside
  const handlePointerDownOutside = (event) => {
    const target = event.target
    const isPointerDownOnBranch = [...context.branches].some((branch) => branch.contains(target))
    if (!isPointerEventsEnabled() || isPointerDownOnBranch) return
    opts.onPointerDownOutside?.(event)
    opts.onInteractOutside?.(event)
    if (!event.defaultPrevented) opts.onDismiss?.()
  }
  let handleClick = () => {}
  const handlePointerDown = (event) => {
    if (event.target && !record.pointerInside) {
      const eventDetail = { originalEvent: event }
      const dispatch = () => handleAndDispatchCustomEvent(POINTER_DOWN_OUTSIDE, handlePointerDownOutside, eventDetail)
      // On touch devices, wait for the click: the touch may become a scroll.
      if (event.pointerType === "touch") {
        ownerDocument.removeEventListener("click", handleClick)
        handleClick = dispatch
        ownerDocument.addEventListener("click", handleClick, { once: true })
      } else {
        dispatch()
      }
    } else {
      ownerDocument.removeEventListener("click", handleClick)
    }
    record.pointerInside = false
  }
  // Registered on the next task so the event that opened the layer does not
  // dismiss it.
  const timerId = window.setTimeout(() => ownerDocument.addEventListener("pointerdown", handlePointerDown), 0)
  cleanups.push(() => {
    window.clearTimeout(timerId)
    ownerDocument.removeEventListener("pointerdown", handlePointerDown)
    ownerDocument.removeEventListener("click", handleClick)
  })

  // useFocusOutside
  const handleFocusOutside = (event) => {
    const target = event.target
    const isFocusInBranch = [...context.branches].some((branch) => branch.contains(target))
    if (isFocusInBranch) return
    opts.onFocusOutside?.(event)
    opts.onInteractOutside?.(event)
    if (!event.defaultPrevented) opts.onDismiss?.()
  }
  const handleFocus = (event) => {
    if (event.target && !record.focusInside) {
      handleAndDispatchCustomEvent(FOCUS_OUTSIDE, handleFocusOutside, { originalEvent: event })
    }
  }
  ownerDocument.addEventListener("focusin", handleFocus)
  cleanups.push(() => ownerDocument.removeEventListener("focusin", handleFocus))

  // useEscapeKeydown
  const handleKeyDown = (event) => {
    if (event.key !== "Escape") return
    const isHighestLayer = [...context.layers].indexOf(record) === context.layers.size - 1
    if (!isHighestLayer) return
    opts.onEscapeKeyDown?.(event)
    if (!event.defaultPrevented && opts.onDismiss) {
      event.preventDefault()
      opts.onDismiss()
    }
  }
  ownerDocument.addEventListener("keydown", handleKeyDown, { capture: true })
  cleanups.push(() => ownerDocument.removeEventListener("keydown", handleKeyDown, { capture: true }))

  // Mount effect: body pointer events, then register the layer.
  if (disableOutsidePointerEvents) {
    if (context.layersWithOutsidePointerEventsDisabled.size === 0) {
      originalBodyPointerEvents = ownerDocument.body.style.pointerEvents
      ownerDocument.body.style.pointerEvents = "none"
    }
    context.layersWithOutsidePointerEventsDisabled.add(record)
  }
  context.layers.add(record)
  dispatchUpdate()

  let destroyed = false
  record.handle = {
    node,
    update(next) {
      opts = { ...opts, ...next }
    },
    setDisableOutsidePointerEvents(on) {
      on = !!on
      if (destroyed || on === disableOutsidePointerEvents) return
      // Effect cleanup with the previous value, then the effect with the new one.
      if (disableOutsidePointerEvents && context.layersWithOutsidePointerEventsDisabled.size === 1) {
        ownerDocument.body.style.pointerEvents = originalBodyPointerEvents
      }
      disableOutsidePointerEvents = on
      if (on) {
        if (context.layersWithOutsidePointerEventsDisabled.size === 0) {
          originalBodyPointerEvents = ownerDocument.body.style.pointerEvents
          ownerDocument.body.style.pointerEvents = "none"
        }
        context.layersWithOutsidePointerEventsDisabled.add(record)
      }
      context.layers.add(record)
      dispatchUpdate()
    },
    destroy() {
      if (destroyed) return
      destroyed = true
      for (const c of cleanups) c()
      if (disableOutsidePointerEvents && context.layersWithOutsidePointerEventsDisabled.size === 1) {
        ownerDocument.body.style.pointerEvents = originalBodyPointerEvents
      }
      context.layers.delete(record)
      context.layersWithOutsidePointerEventsDisabled.delete(record)
      dispatchUpdate()
    },
  }
  return record.handle
}

function recordOf(handle) {
  for (const r of context.layers) if (r.handle === handle) return r
  return undefined
}

function recordOfElement(el) {
  if (!el) return undefined
  const layers = [...context.layers]
  for (let i = layers.length - 1; i >= 0; i--) {
    if (layers[i].node.contains(el)) return layers[i]
  }
  return undefined
}

/**
 * DismissableLayer.Branch: pointer and focus events inside node are never
 * outside any layer.
 * @param {Element} node
 * @returns {() => void}
 */
export function createDismissableLayerBranch(node) {
  context.branches.add(node)
  return () => context.branches.delete(node)
}
