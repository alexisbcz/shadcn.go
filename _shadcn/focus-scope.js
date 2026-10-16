// focus-scope.js: port of Radix FocusScope (@radix-ui/react-focus-scope
// 1.1.7).
//
// A focus scope moves focus into its container when it mounts (first
// tabbable element, links excluded, inputs selected; else the container
// itself, which therefore renders tabindex="-1"), gives it back when it
// unmounts (to the element focused before), optionally loops Tab /
// Shift+Tab at the edges, and optionally traps focus inside (focus leaving
// the container is pulled back to the last focused element inside). Nested
// scopes pause the ones below them.
//
// API
//
//   import { createFocusScope } from "./focus-scope.js"
//
//   const scope = createFocusScope(container, {
//     loop: false,              // Tab / Shift+Tab wrap at the edges
//     trapped: false,           // focus cannot leave the container
//     onMountAutoFocus(event) {},   // CustomEvent "focusScope.autoFocusOnMount";
//                                   // preventDefault() to focus something else
//     onUnmountAutoFocus(event) {}, // CustomEvent "focusScope.autoFocusOnUnmount";
//                                   // preventDefault() to skip restoring focus
//   })
//   scope.setTrapped(false)     // like changing the `trapped` prop
//   scope.destroy()             // unmount: on the next task, dispatches the
//                               // unmount event and restores focus
//
// Create it when Radix mounts the FocusScope (container in the DOM, after
// the effects of the components inside it ran), destroy it when it unmounts.
//
// Helpers (also used by menus, popovers, ...):
//   getTabbableCandidates(container), getTabbableEdges(container),
//   focusFirst(candidates, { select }), focus(element, { select }),
//   removeLinks(elements)

const AUTOFOCUS_ON_MOUNT = "focusScope.autoFocusOnMount"
const AUTOFOCUS_ON_UNMOUNT = "focusScope.autoFocusOnUnmount"
const EVENT_OPTIONS = { bubbles: false, cancelable: true }

/**
 * @param {HTMLElement} container
 * @param {{ loop?: boolean, trapped?: boolean, onMountAutoFocus?: (e: Event) => void, onUnmountAutoFocus?: (e: Event) => void }} [props]
 */
export function createFocusScope(container, props = {}) {
  const loop = !!props.loop
  let trapped = !!props.trapped
  const onMountAutoFocus = (e) => props.onMountAutoFocus?.(e)
  const onUnmountAutoFocus = (e) => props.onUnmountAutoFocus?.(e)
  let lastFocusedElement = null

  const focusScope = {
    paused: false,
    pause() {
      this.paused = true
    },
    resume() {
      this.paused = false
    },
  }

  // Effect 1: the trap.
  let untrap = null
  const trap = () => {
    if (!trapped) return
    const handleFocusIn = (event) => {
      if (focusScope.paused || !container) return
      const target = event.target
      if (container.contains(target)) lastFocusedElement = target
      else focus(lastFocusedElement, { select: true })
    }
    const handleFocusOut = (event) => {
      if (focusScope.paused || !container) return
      const relatedTarget = event.relatedTarget
      // A null relatedTarget means focus left the page or went to body
      // (e.g. the focused element was removed): handled by the observer.
      if (relatedTarget === null) return
      if (!container.contains(relatedTarget)) focus(lastFocusedElement, { select: true })
    }
    const handleMutations = (mutations) => {
      const focusedElement = document.activeElement
      if (focusedElement !== document.body) return
      for (const mutation of mutations) {
        if (mutation.removedNodes.length > 0) focus(container)
      }
    }
    document.addEventListener("focusin", handleFocusIn)
    document.addEventListener("focusout", handleFocusOut)
    const mutationObserver = new MutationObserver(handleMutations)
    mutationObserver.observe(container, { childList: true, subtree: true })
    untrap = () => {
      document.removeEventListener("focusin", handleFocusIn)
      document.removeEventListener("focusout", handleFocusOut)
      mutationObserver.disconnect()
      untrap = null
    }
  }
  trap()

  // Effect 2: auto focus on mount.
  focusScopesStack.add(focusScope)
  const previouslyFocusedElement = document.activeElement
  const hasFocusedCandidate = container.contains(previouslyFocusedElement)
  if (!hasFocusedCandidate) {
    const mountEvent = new CustomEvent(AUTOFOCUS_ON_MOUNT, EVENT_OPTIONS)
    container.addEventListener(AUTOFOCUS_ON_MOUNT, onMountAutoFocus)
    container.dispatchEvent(mountEvent)
    if (!mountEvent.defaultPrevented) {
      focusFirst(removeLinks(getTabbableCandidates(container)), { select: true })
      if (document.activeElement === previouslyFocusedElement) focus(container)
    }
  }

  // Tab / Shift+Tab at the edges.
  const handleKeyDown = (event) => {
    if (!loop && !trapped) return
    if (focusScope.paused) return
    const isTabKey = event.key === "Tab" && !event.altKey && !event.ctrlKey && !event.metaKey
    const focusedElement = document.activeElement
    if (isTabKey && focusedElement) {
      const container2 = event.currentTarget
      const [first, last] = getTabbableEdges(container2)
      const hasTabbableElementsInside = first && last
      if (!hasTabbableElementsInside) {
        if (focusedElement === container2) event.preventDefault()
      } else {
        if (!event.shiftKey && focusedElement === last) {
          event.preventDefault()
          if (loop) focus(first, { select: true })
        } else if (event.shiftKey && focusedElement === first) {
          event.preventDefault()
          if (loop) focus(last, { select: true })
        }
      }
    }
  }
  container.addEventListener("keydown", handleKeyDown)

  let destroyed = false
  return {
    focusScope,
    setTrapped(next) {
      next = !!next
      if (next === trapped) return
      trapped = next
      untrap?.()
      trap()
    },
    destroy() {
      if (destroyed) return
      destroyed = true
      untrap?.()
      container.removeEventListener("keydown", handleKeyDown)
      container.removeEventListener(AUTOFOCUS_ON_MOUNT, onMountAutoFocus)
      setTimeout(() => {
        const unmountEvent = new CustomEvent(AUTOFOCUS_ON_UNMOUNT, EVENT_OPTIONS)
        container.addEventListener(AUTOFOCUS_ON_UNMOUNT, onUnmountAutoFocus)
        container.dispatchEvent(unmountEvent)
        if (!unmountEvent.defaultPrevented) focus(previouslyFocusedElement ?? document.body, { select: true })
        container.removeEventListener(AUTOFOCUS_ON_UNMOUNT, onUnmountAutoFocus)
        focusScopesStack.remove(focusScope)
      }, 0)
    },
  }
}

/** Focuses the first candidate that accepts focus. */
export function focusFirst(candidates, { select = false } = {}) {
  const previouslyFocusedElement = document.activeElement
  for (const candidate of candidates) {
    focus(candidate, { select })
    if (document.activeElement !== previouslyFocusedElement) return
  }
}

/** Returns the first and last visible tabbable elements of container. */
export function getTabbableEdges(container) {
  const candidates = getTabbableCandidates(container)
  const first = findVisible(candidates, container)
  const last = findVisible(candidates.reverse(), container)
  return [first, last]
}

/**
 * Returns the elements of container with tabIndex >= 0 that are not
 * disabled, hidden or hidden inputs, in DOM order. (Visibility and
 * tabindex order are not considered, like Radix.)
 */
export function getTabbableCandidates(container) {
  const nodes = []
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_ELEMENT, {
    acceptNode: (node) => {
      const isHiddenInput = node.tagName === "INPUT" && node.type === "hidden"
      if (node.disabled || node.hidden || isHiddenInput) return NodeFilter.FILTER_SKIP
      return node.tabIndex >= 0 ? NodeFilter.FILTER_ACCEPT : NodeFilter.FILTER_SKIP
    },
  })
  while (walker.nextNode()) nodes.push(walker.currentNode)
  return nodes
}

function findVisible(elements, container) {
  for (const element of elements) {
    if (!isHidden(element, { upTo: container })) return element
  }
}

function isHidden(node, { upTo }) {
  if (getComputedStyle(node).visibility === "hidden") return true
  while (node) {
    if (upTo !== undefined && node === upTo) return false
    if (getComputedStyle(node).display === "none") return true
    node = node.parentElement
  }
  return false
}

function isSelectableInput(element) {
  return element instanceof HTMLInputElement && "select" in element
}

/** Focuses element without scrolling, selecting an input's text if select. */
export function focus(element, { select = false } = {}) {
  if (element && element.focus) {
    const previouslyFocusedElement = document.activeElement
    element.focus({ preventScroll: true })
    if (element !== previouslyFocusedElement && isSelectableInput(element) && select) element.select()
  }
}

/** Removes <a> elements (links are not auto-focused on mount). */
export function removeLinks(items) {
  return items.filter((item) => item.tagName !== "A")
}

const focusScopesStack = createFocusScopesStack()

function createFocusScopesStack() {
  let stack = []
  return {
    add(focusScope) {
      const activeFocusScope = stack[0]
      if (focusScope !== activeFocusScope) activeFocusScope?.pause()
      stack = arrayRemove(stack, focusScope)
      stack.unshift(focusScope)
    },
    remove(focusScope) {
      stack = arrayRemove(stack, focusScope)
      stack[0]?.resume()
    },
  }
}

function arrayRemove(array, item) {
  const updatedArray = [...array]
  const index = updatedArray.indexOf(item)
  if (index !== -1) updatedArray.splice(index, 1)
  return updatedArray
}
