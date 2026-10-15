// collapsible-content.js: port of Radix CollapsibleContent
// (@radix-ui/react-collapsible 1.1.12: CollapsibleContent + CollapsibleContentImpl
// and the Presence that wraps it), shared by Collapsible and Accordion.
//
// The content element is always in the DOM (Radix renders it through a
// function-child Presence, i.e. force-mounted). While closed it is `hidden`
// and empty; while open, and while its exit animation runs, it holds its
// children. Every time the open state changes, Radix's layout effect
// disables the content's animation and transition, measures it, restores
// them (except right after mount, so an initially open content does not
// animate in) and publishes the size as
// --radix-collapsible-content-height / --radix-collapsible-content-width,
// which the accordion-down/up and collapsible-down/up keyframes use.
//
//   import { createCollapsibleContent } from "./collapsible-content.js"
//
//   const content = createCollapsibleContent(el, {
//     open: false,         // initial state, as server-rendered: open content
//                          // carries its children and the mount-prevention
//                          // styles "transition-duration:0s;animation-name:none";
//                          // closed content is hidden and empty
//     fragment,            // DocumentFragment with the children: cloned on
//                          // every open, like React remounting them
//   })
//   content.setOpen(true)  // after updating the trigger/root attributes
//   content.open           // current open state
//   content.destroy()
//
// Markup the server renders (React's first render): data-state, hidden when
// closed, style="" when closed (React leaves an empty style attribute after
// the mount effect restored the styles) or the mount-prevention styles when
// open. The size variables are measured here at mount.

import { createPresence } from "./presence.js"

const HEIGHT_VAR = "--radix-collapsible-content-height"
const WIDTH_VAR = "--radix-collapsible-content-width"

/**
 * @param {HTMLElement} node
 * @param {{ open?: boolean, fragment?: DocumentFragment | null }} [options]
 */
export function createCollapsibleContent(node, { open = false, fragment = null } = {}) {
  const win = node.ownerDocument.defaultView ?? window
  let contextOpen = open
  let isPresentState = open // CollapsibleContentImpl's isPresent state
  let height = 0 // heightRef
  let width = 0 // widthRef
  let mountAnimationPrevented = open // isMountAnimationPreventedRef = useRef(isOpen)
  // originalStylesRef: the styles React rendered before the first layout
  // effect touched them. The server renders the mount-prevention styles on
  // open content, so they are not the original ones.
  const original = open
    ? { transitionDuration: "", animationName: "" }
    : { transitionDuration: node.style.transitionDuration, animationName: node.style.animationName }
  let childrenMounted = open
  let lastOpen
  let lastPresent
  let destroyed = false
  const untilFound = node.hasAttribute("data-hidden-until-found")

  // React render of CollapsibleContentImpl.
  const render = () => {
    const isOpen = contextOpen || isPresentState
    node.setAttribute("data-state", contextOpen ? "open" : "closed")
    // data-hidden-until-found: closed content stays laid out (React Aria's
    // DisclosurePanel uses hidden="until-found" once hydrated).
    if (!isOpen && untilFound) node.setAttribute("hidden", "until-found")
    else node.hidden = !isOpen
    if (isOpen && !childrenMounted) {
      childrenMounted = true
      node.replaceChildren(...(fragment ? [fragment.cloneNode(true)] : []))
    } else if (!isOpen && childrenMounted) {
      childrenMounted = false
      node.replaceChildren()
    }
    if (height) node.style.setProperty(HEIGHT_VAR, `${height}px`)
    else node.style.removeProperty(HEIGHT_VAR)
    if (width) node.style.setProperty(WIDTH_VAR, `${width}px`)
    else node.style.removeProperty(WIDTH_VAR)
  }

  // useLayoutEffect(..., [context.open, present]).
  const layoutEffect = (present) => {
    node.style.transitionDuration = "0s"
    node.style.animationName = "none"
    const rect = node.getBoundingClientRect()
    height = rect.height
    width = rect.width
    if (!mountAnimationPrevented) {
      node.style.transitionDuration = original.transitionDuration
      node.style.animationName = original.animationName
    }
    const changed = isPresentState !== present
    isPresentState = present
    return changed
  }

  // Renders, then runs the layout effect when its dependencies changed, and
  // re-renders while it updates state, like React's synchronous layout phase.
  const commit = () => {
    if (destroyed) return
    for (;;) {
      const present = presence.isPresent
      render()
      if (present === lastPresent && contextOpen === lastOpen) break
      lastPresent = present
      lastOpen = contextOpen
      if (!layoutEffect(present)) break
    }
  }

  const presence = createPresence(node, {
    present: open,
    onExitComplete: () => commit(),
  })

  // Mount.
  lastOpen = contextOpen
  lastPresent = presence.isPresent
  layoutEffect(presence.isPresent)
  render()
  const rAF = win.requestAnimationFrame(() => (mountAnimationPrevented = false))

  return {
    get open() {
      return contextOpen
    },
    setOpen(next) {
      if (destroyed || next === contextOpen) return
      contextOpen = next
      commit()
      // Presence's layout effect runs after its child's.
      presence.setPresent(next)
      commit()
    },
    destroy() {
      destroyed = true
      win.cancelAnimationFrame(rAF)
      presence.destroy()
    },
  }
}
