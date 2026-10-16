// presence.js: port of Radix Presence (@radix-ui/react-presence 1.1.5).
//
// Presence decides when an element that is going away may actually leave:
// immediately when it has no exit animation, or once its exit animation
// (the CSS animation that starts when its state attribute flips, e.g.
// data-state="closed" + tw-animate-css "animate-out") has ended. It is the
// state machine Radix runs for every Dialog, Popover, Menu, Tooltip,
// Collapsible, Accordion, Checkbox indicator, ... content.
//
// API
//
//   import { createPresence } from "./presence.js"
//
//   const presence = createPresence(node, {
//     present: true,          // initial value (default true): the node is in
//                             // the DOM and showing
//     onExitComplete() {},    // the node may leave: remove it from the DOM
//                             // (Radix unmounts it) or hide it (forceMount
//                             // users such as Collapsible set `hidden`)
//   })
//
//   // Showing -> hiding. FIRST update the node's attributes so the exit
//   // animation applies (e.g. data-state="closed"), THEN call:
//   presence.setPresent(false)  // calls onExitComplete now (no exit
//                               // animation, or display: none) or when the
//                               // animation ends
//   presence.setPresent(true)   // re-show while the exit animation runs:
//                               // cancels the pending onExitComplete
//   presence.isPresent          // true while mounted or exiting (Radix
//                               // isPresent: "mounted" | "unmountSuspended")
//   presence.state              // "mounted" | "unmountSuspended" | "unmounted"
//   presence.destroy()          // removes the listeners
//
// The node must be connected (or connected right after) when the presence
// is created and when setPresent is called: the decision reads its computed
// animation-name, exactly like Radix's useLayoutEffect does after React
// committed the new attributes.
//
// Like Radix, when the exit animation ends while not present, the node's
// inline animation-fill-mode is set to "forwards" until the next task, so
// the end frame does not flash before the node is removed or hidden.

/** Returns the computed animation-name of styles, "none" if unknown. */
export function getAnimationName(styles) {
  return styles?.animationName || "none"
}

const MACHINE = {
  mounted: { UNMOUNT: "unmounted", ANIMATION_OUT: "unmountSuspended" },
  unmountSuspended: { MOUNT: "mounted", ANIMATION_END: "unmounted" },
  unmounted: { MOUNT: "mounted" },
}

/**
 * @param {HTMLElement} node
 * @param {{ present?: boolean, onExitComplete?: () => void }} [options]
 */
export function createPresence(node, { present = true, onExitComplete } = {}) {
  const win = node.ownerDocument.defaultView ?? window
  const styles = win.getComputedStyle(node) // live object, like stylesRef
  let state = present ? "mounted" : "unmounted"
  let prevPresent = present
  let prevAnimationName = "none"
  let timeoutId

  // Radix: useEffect(() => { prevAnimationNameRef = state === "mounted" ? current : "none" }, [state])
  const syncPrevAnimationName = () => {
    prevAnimationName = state === "mounted" ? getAnimationName(styles) : "none"
  }
  syncPrevAnimationName()

  const send = (event) => {
    const next = MACHINE[state][event]
    if (!next || next === state) return
    state = next
    syncPrevAnimationName()
    if (state === "unmounted") onExitComplete?.()
  }

  const handleAnimationEnd = (event) => {
    const currentAnimationName = getAnimationName(styles)
    const isCurrentAnimation = currentAnimationName.includes(CSS.escape(event.animationName))
    if (event.target === node && isCurrentAnimation) {
      if (!prevPresent) {
        // Radix sets this right after dispatching ANIMATION_END, before
        // React re-renders and unmounts; set it first since our
        // onExitComplete runs synchronously.
        const currentFillMode = node.style.animationFillMode
        node.style.animationFillMode = "forwards"
        timeoutId = win.setTimeout(() => {
          if (node.style.animationFillMode === "forwards") node.style.animationFillMode = currentFillMode
        })
      }
      send("ANIMATION_END")
    }
  }
  const handleAnimationStart = (event) => {
    if (event.target === node) prevAnimationName = getAnimationName(styles)
  }
  node.addEventListener("animationstart", handleAnimationStart)
  node.addEventListener("animationcancel", handleAnimationEnd)
  node.addEventListener("animationend", handleAnimationEnd)

  return {
    get state() {
      return state
    },
    get isPresent() {
      return state === "mounted" || state === "unmountSuspended"
    },
    setPresent(present) {
      const wasPresent = prevPresent
      if (wasPresent === present) return
      prevPresent = present
      const currentAnimationName = getAnimationName(styles)
      if (present) {
        send("MOUNT")
      } else if (currentAnimationName === "none" || styles?.display === "none") {
        send("UNMOUNT")
      } else {
        const isAnimating = prevAnimationName !== currentAnimationName
        if (wasPresent && isAnimating) send("ANIMATION_OUT")
        else send("UNMOUNT")
      }
    },
    destroy() {
      win.clearTimeout(timeoutId)
      node.removeEventListener("animationstart", handleAnimationStart)
      node.removeEventListener("animationcancel", handleAnimationEnd)
      node.removeEventListener("animationend", handleAnimationEnd)
    },
  }
}
