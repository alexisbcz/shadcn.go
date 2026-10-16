// roving-focus.js: port of @radix-ui/react-roving-focus 1.1.11
// (RovingFocusGroup + RovingFocusGroupItem), the keyboard navigation of
// Tabs, ToggleGroup, RadioGroup, Toolbar, Menubar, menus, ...
//
// Markup (server-rendered, as Radix renders it once its effects ran):
//   group: tabindex="0" (or "-1" when no item is focusable), style="outline:none",
//          data-orientation when an orientation is set
//   items: tabindex="-1" (no item is the current tab stop until one gets focus),
//          data-radix-collection-item="", data-orientation when set
//
// Behavior, identical to Radix:
//   - Tab into the group focuses the group itself (it is the tab stop until an
//     item was focused); the group then dispatches the cancelable, non-bubbling
//     "rovingFocusGroup.onEntryFocus" event and, unless it was prevented, moves
//     focus to the active item, else the current tab stop, else the first
//     focusable item. A mouse click on the group does not do that.
//   - An item that gets focus (or a mousedown) becomes the current tab stop:
//     tabindex="0", every other item tabindex="-1".
//   - Shift+Tab from an item sets the group's tabindex to -1 until focus
//     leaves it, so focus goes before the group instead of to it.
//   - Arrow keys (per orientation and reading direction), Home/End and
//     PageUp/PageDown move focus among focusable items in a setTimeout, wrapping
//     when loop is true. Keys with a modifier are ignored. Only keydowns whose
//     target is the item itself count.
//   - mousedown on a non-focusable item is prevented.
//
//   import { createRovingFocusGroup } from "./roving-focus.js"
//   const group = createRovingFocusGroup(listEl, {
//     orientation: "horizontal",       // "horizontal" | "vertical" | undefined (both)
//     dir: "ltr",                      // default: getDirection(listEl)
//     loop: true,                      // default false, like Radix
//     focusable: (item) => !item.disabled,       // default: true
//     active: (item) => item.getAttribute("aria-selected") === "true", // default: false
//     onEntryFocus: (event) => {},     // may call event.preventDefault()
//     onCurrentTabStopChange: (item) => {},
//     preventScrollOnEntryFocus: false,
//     groupTabIndex: undefined,        // a number pins the group's own tabindex (Radix
//                                      // menus: FocusScope's tabIndex=-1 overrides the
//                                      // group's computed one); items stay managed
//   })
//   group.update()      // re-read focusable/active after the component changed state
//   group.getItems()    // focusable or not, document order
//   group.destroy()
//
// Ordering: Radix composes the component's own handlers before its own and
// skips its own when the component called preventDefault. Attach component
// listeners to the items (or to the group before calling
// createRovingFocusGroup): the group listens on the group element, so item
// listeners run first; roving item handlers skip defaultPrevented events.

import { createCollection } from "./collection.js"
import { getDirection, getDirectionAwareKey } from "./direction.js"

export const ENTRY_FOCUS = "rovingFocusGroup.onEntryFocus"
const EVENT_OPTIONS = { bubbles: false, cancelable: true }

const MAP_KEY_TO_FOCUS_INTENT = {
  ArrowLeft: "prev",
  ArrowUp: "prev",
  ArrowRight: "next",
  ArrowDown: "next",
  PageUp: "first",
  Home: "first",
  PageDown: "last",
  End: "last",
}

/** Radix's getFocusIntent. */
export function getFocusIntent(event, orientation, dir) {
  const key = getDirectionAwareKey(event.key, dir)
  if (orientation === "vertical" && ["ArrowLeft", "ArrowRight"].includes(key)) return undefined
  if (orientation === "horizontal" && ["ArrowUp", "ArrowDown"].includes(key)) return undefined
  return MAP_KEY_TO_FOCUS_INTENT[key]
}

/** Focuses the first candidate that takes focus (Radix's focusFirst). */
export function focusFirst(candidates, preventScroll = false) {
  const PREVIOUSLY_FOCUSED_ELEMENT = document.activeElement
  for (const candidate of candidates) {
    if (candidate === PREVIOUSLY_FOCUSED_ELEMENT) return
    candidate.focus({ preventScroll })
    if (document.activeElement !== PREVIOUSLY_FOCUSED_ELEMENT) return
  }
}

/** Rotates array so it starts at startIndex (Radix's wrapArray). */
export function wrapArray(array, startIndex) {
  return array.map((_, index) => array[(startIndex + index) % array.length])
}

/**
 * Makes root a roving focus group.
 * @param {HTMLElement} root
 * @param {{
 *   orientation?: "horizontal" | "vertical",
 *   dir?: "ltr" | "rtl",
 *   loop?: boolean,
 *   focusable?: (item: HTMLElement) => boolean,
 *   active?: (item: HTMLElement) => boolean,
 *   currentTabStop?: HTMLElement | null,
 *   onCurrentTabStopChange?: (item: HTMLElement | null) => void,
 *   onEntryFocus?: (event: Event) => void,
 *   preventScrollOnEntryFocus?: boolean,
 *   manageTabIndex?: boolean,
 *   groupTabIndex?: number,
 * }} [options]
 */
export function createRovingFocusGroup(root, options = {}) {
  const opts = {
    orientation: root.getAttribute("data-orientation") || undefined,
    loop: false,
    focusable: () => true,
    active: () => false,
    preventScrollOnEntryFocus: false,
    manageTabIndex: true,
    ...options,
  }
  const dir = () => opts.dir ?? getDirection(root)
  const collection = createCollection(root)
  let currentTabStop = opts.currentTabStop ?? null
  let isTabbingBackOut = false
  let isClickFocus = false

  const getItems = () => collection.getItems()
  const focusableItems = () => getItems().filter((item) => opts.focusable(item))

  const render = () => {
    if (!opts.manageTabIndex) return
    const items = getItems()
    if (currentTabStop && !items.includes(currentTabStop)) currentTabStop = null
    const count = items.filter((item) => opts.focusable(item)).length
    setTabIndex(root, typeof opts.groupTabIndex === "number" ? opts.groupTabIndex : isTabbingBackOut || count === 0 ? -1 : 0)
    for (const item of items) setTabIndex(item, item === currentTabStop ? 0 : -1)
  }

  const setCurrentTabStop = (item) => {
    if (item === currentTabStop) return
    currentTabStop = item
    render()
    opts.onCurrentTabStopChange?.(item)
  }

  // --- group handlers ---
  const onMouseDown = (event) => {
    // Radix composes every handler with composeEventHandlers: once a
    // component handler (e.g. a Tabs trigger on right-click) or the item
    // handler (non-focusable item) prevented the mousedown, neither the item
    // nor the group handler runs, so isClickFocus stays false.
    if (event.defaultPrevented) return
    const item = collection.itemFor(event.target)
    if (item) {
      if (!opts.focusable(item)) {
        event.preventDefault()
        return
      }
      setCurrentTabStop(item)
    }
    isClickFocus = true
  }

  const onFocusIn = (event) => {
    const item = collection.itemFor(event.target)
    if (item && !event.defaultPrevented) setCurrentTabStop(item)
    const isKeyboardFocus = !isClickFocus
    if (event.target === root && isKeyboardFocus && !isTabbingBackOut) {
      const entryFocusEvent = new CustomEvent(ENTRY_FOCUS, EVENT_OPTIONS)
      root.dispatchEvent(entryFocusEvent)
      if (!entryFocusEvent.defaultPrevented) {
        const items = focusableItems()
        const activeItem = items.find((it) => opts.active(it))
        const currentItem = items.find((it) => it === currentTabStop)
        const candidates = [activeItem, currentItem, ...items].filter(Boolean)
        focusFirst(candidates, opts.preventScrollOnEntryFocus)
      }
    }
    isClickFocus = false
  }

  const onFocusOut = () => {
    if (!isTabbingBackOut) return
    isTabbingBackOut = false
    render()
  }

  const onKeyDown = (event) => {
    const item = collection.itemFor(event.target)
    if (!item || event.defaultPrevented) return
    if (event.key === "Tab" && event.shiftKey) {
      isTabbingBackOut = true
      render()
      return
    }
    if (event.target !== item) return
    const focusIntent = getFocusIntent(event, opts.orientation, dir())
    if (focusIntent === undefined) return
    if (event.metaKey || event.ctrlKey || event.altKey || event.shiftKey) return
    event.preventDefault()
    let candidates = focusableItems()
    if (focusIntent === "last") candidates.reverse()
    else if (focusIntent === "prev" || focusIntent === "next") {
      if (focusIntent === "prev") candidates.reverse()
      const currentIndex = candidates.indexOf(item)
      candidates = opts.loop ? wrapArray(candidates, currentIndex + 1) : candidates.slice(currentIndex + 1)
    }
    setTimeout(() => focusFirst(candidates))
  }

  const onEntry = (event) => opts.onEntryFocus?.(event)

  root.addEventListener("mousedown", onMouseDown)
  root.addEventListener("focusin", onFocusIn)
  root.addEventListener("focusout", onFocusOut)
  root.addEventListener("keydown", onKeyDown)
  root.addEventListener(ENTRY_FOCUS, onEntry)
  render()

  return {
    root,
    collection,
    getItems,
    focusableItems,
    /** Re-reads focusable/active state and updates tabindex. */
    update: render,
    get currentTabStop() {
      return currentTabStop
    },
    setCurrentTabStop,
    /** Updates options (orientation, dir, loop, focusable, active, ...). */
    setOptions(next) {
      Object.assign(opts, next)
      render()
    },
    destroy() {
      root.removeEventListener("mousedown", onMouseDown)
      root.removeEventListener("focusin", onFocusIn)
      root.removeEventListener("focusout", onFocusOut)
      root.removeEventListener("keydown", onKeyDown)
      root.removeEventListener(ENTRY_FOCUS, onEntry)
      collection.destroy()
    },
  }
}

function setTabIndex(el, value) {
  if (el.getAttribute("tabindex") !== String(value)) el.setAttribute("tabindex", String(value))
}
