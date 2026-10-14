// popover.js: Radix Popover behavior (@radix-ui/react-popover 1.1.15) for
// the server-rendered markup of components/ui/popover.
//
// Markup contract (rendered by the Go package):
//   - <template data-radix-popover data-content-id="<id>" [data-modal]
//       [data-default-open]> holds the popper wrapper (see popper.js) with
//     the content exactly as React renders it when open, before it is
//     positioned. It is cloned into <body> on open.
//   - Triggers are the elements with aria-controls="<id>" (PopoverTrigger
//     always renders it).
//   - Markers removed on mount: data-radix-popover-anchor="<id>" (a custom
//     PopoverAnchor), data-radix-popover-close="<id>" (PopoverClose).
//
// Scripts can observe and drive a popover:
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getPopover } from "/_shadcn/popover.js"
//   getPopover(contentIdOrTemplateOrTrigger)?.setOpen(true)   // .open, .toggle()

import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createFocusScope } from "./focus-scope.js"
import { acquireFocusGuards } from "./focus-guards.js"
import { lockScroll } from "./scroll-lock.js"
import { hideOthers } from "./hide-others.js"
import { createPopper, openOnLoad } from "./popper.js"

const popovers = new Map() // content id -> controller

/** Returns the controller of a popover by content id, template or trigger. */
export function getPopover(ref) {
  if (typeof ref === "string") return popovers.get(ref)
  if (ref instanceof Element) {
    const id = ref.getAttribute("data-content-id") ?? ref.getAttribute("aria-controls")
    return id ? popovers.get(id) : undefined
  }
  return undefined
}

const getState = (open) => (open ? "open" : "closed")


// Finds the elements carrying marker attribute name with value, removes the
// attribute and returns them (other components' markers stay in place).
function takeOwnMarked(name, value) {
  return [...document.querySelectorAll(`[${name}="${CSS.escape(value)}"]`)].map((el) => {
    el.removeAttribute(name)
    return el
  })
}

function createPopover(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const modal = tpl.hasAttribute("data-modal")
  let open = false
  let mounted = null // { wrapper, content, presence, cleanups, focusScope, popper } while present

  const esc = CSS.escape(contentId)
  const triggers = () => [...document.querySelectorAll(`[aria-controls="${esc}"]`)]
  const trigger = () => triggers().pop() ?? null

  // PopoverAnchor: a custom anchor replaces the trigger as the popper anchor.
  const customAnchors = takeOwnMarked("data-radix-popover-anchor", contentId)
  // Radix sets the popper anchor in effects, in tree order: the trigger
  // wraps itself in a PopperAnchor until PopoverAnchor's effect reports a
  // custom anchor, then re-renders unwrapped, which remounts its element.
  // A custom anchor placed after the trigger wins; one placed before it is
  // overwritten by the trigger's first element, now detached (React keeps
  // positioning against that stale node: a 0x0 rect at the viewport origin).
  let staleTrigger = null
  const anchor = () => {
    const custom = customAnchors[customAnchors.length - 1]
    const t = trigger()
    if (!custom) return t
    if (t && custom.compareDocumentPosition(t) & Node.DOCUMENT_POSITION_FOLLOWING) {
      staleTrigger ??= t.cloneNode(true)
      return staleTrigger
    }
    return custom
  }

  // PopoverClose buttons outside the content.
  const outsideCloses = takeOwnMarked("data-radix-popover-close", contentId)
  const onCloseClick = (event) => {
    if (!event.defaultPrevented) setOpen(false)
  }
  outsideCloses.forEach((el) => el.addEventListener("click", onCloseClick))

  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    for (const t of triggers()) {
      t.setAttribute("aria-expanded", String(open))
      t.setAttribute("data-state", getState(open))
    }
    if (open) {
      if (mounted) reshow()
      else mountContent()
    } else if (mounted) {
      mounted.content.setAttribute("data-state", "closed")
      mounted.focusScope?.setTrapped(false)
      mounted.presence.setPresent(false)
    }
  }

  function reshow() {
    mounted.content.setAttribute("data-state", "open")
    mounted.focusScope?.setTrapped(modal)
    mounted.presence.setPresent(true)
  }

  function mountContent() {
    const nodes = instantiate(tpl)
    const wrapper = nodes[0]
    const content = wrapper.firstElementChild
    const closes = takeMarked(nodes, "data-radix-popover-close").map((m) => m.el)
    const state = { wrapper, content, cleanups: [], focusScope: null, presence: null, popper: null }
    mounted = state
    closes.forEach((el) => el.addEventListener("click", onCloseClick))

    // Portal: mounted at the end of <body>.
    const unmount = mount(nodes)

    // Effects, children first: PopperContent, DismissableLayer, FocusScope,
    // PopoverContentImpl (focus guards), PopoverContentModal (hideOthers).
    const popper = createPopper(wrapper, { anchor: anchor() })
    state.popper = popper
    state.cleanups.push(() => popper.destroy())

    let isRightClickOutside = false
    let hasInteractedOutside = false
    let hasPointerDownOutside = false
    const layer = createDismissableLayer(content, {
      disableOutsidePointerEvents: modal,
      parent: trigger(),
      onDismiss: () => setOpen(false),
      onPointerDownOutside: (event) => {
        if (!modal) return
        const originalEvent = event.detail.originalEvent
        const ctrlLeftClick = originalEvent.button === 0 && originalEvent.ctrlKey === true
        isRightClickOutside = originalEvent.button === 2 || ctrlLeftClick
      },
      onFocusOutside: (event) => {
        if (modal) event.preventDefault()
      },
      onInteractOutside: (event) => {
        if (modal) return
        // PopoverContentNonModal
        if (!event.defaultPrevented) {
          hasInteractedOutside = true
          if (event.detail.originalEvent.type === "pointerdown") hasPointerDownOutside = true
        }
        const target = event.target
        const targetIsTrigger = trigger()?.contains(target)
        if (targetIsTrigger) event.preventDefault()
        if (event.detail.originalEvent.type === "focusin" && hasPointerDownOutside) event.preventDefault()
      },
    })
    state.cleanups.push(() => layer.destroy())

    const focusScope = createFocusScope(content, {
      loop: true,
      trapped: modal && open,
      onUnmountAutoFocus: (event) => {
        if (modal) {
          event.preventDefault()
          if (!isRightClickOutside) trigger()?.focus()
          return
        }
        if (!event.defaultPrevented) {
          if (!hasInteractedOutside) trigger()?.focus()
          event.preventDefault()
        }
        hasInteractedOutside = false
        hasPointerDownOutside = false
      },
    })
    state.focusScope = focusScope
    state.cleanups.push(() => focusScope.destroy())

    state.cleanups.push(acquireFocusGuards())
    if (modal) {
      state.cleanups.push(hideOthers(content))
      // <RemoveScroll as={Slot} allowPinchZoom>: the lock node is the content.
      state.cleanups.push(lockScroll(content, { allowPinchZoom: true }))
    }

    const presence = createPresence(content, {
      present: true,
      onExitComplete: () => {
        presence.destroy()
        unmount()
        // Unmount runs parent effects' cleanups first.
        for (const c of state.cleanups.reverse()) c()
        if (mounted === state) mounted = null
      },
    })
    state.presence = presence
  }

  const controller = {
    get open() {
      return open
    },
    setOpen,
    toggle: () => setOpen(!open),
    dispose() {
      outsideCloses.forEach((el) => el.removeEventListener("click", onCloseClick))
      if (mounted) {
        const state = mounted
        state.presence.destroy()
        state.wrapper.remove()
        for (const c of state.cleanups.reverse()) c()
        mounted = null
      }
      if (popovers.get(contentId) === controller) popovers.delete(contentId)
    },
  }
  popovers.set(contentId, controller)
  if (tpl.hasAttribute("data-default-open")) {
    // The server rendered the triggers open already.
    openOnLoad(tpl, () => setOpen(true))
  }
  return controller
}

// PopoverTrigger: onClick toggles (composed after the user's handlers).
document.addEventListener("click", (event) => {
  if (event.defaultPrevented) return
  const t = event.target instanceof Element ? event.target.closest("[aria-controls]") : null
  if (!t) return
  const p = popovers.get(t.getAttribute("aria-controls"))
  if (p) p.toggle()
})

define("popover", {
  selector: "template[data-radix-popover]",
  mount(tpl) {
    const p = createPopover(tpl)
    return () => p.dispose()
  },
})
