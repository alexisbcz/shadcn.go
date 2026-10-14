// hover-card.js: Radix HoverCard behavior (@radix-ui/react-hover-card
// 1.1.15) for the server-rendered markup of components/ui/hover-card.
//
// Markup contract (rendered by the Go package):
//   - <template data-radix-hover-card data-content-id="<id>"
//       data-open-delay="<ms>" data-close-delay="<ms>" [data-default-open]>
//     holds the popper wrapper (see popper.js) with the content as React
//     renders it, before it is positioned. It is cloned into <body> on open.
//   - The trigger carries data-radix-hover-card-trigger="<id>", removed on
//     mount.
//
// Scripts can observe and drive a hover card:
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getHoverCard } from "/_shadcn/hover-card.js"
//   getHoverCard(contentIdOrTemplate)?.setOpen(true)

import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createPopper, openOnLoad } from "./popper.js"

const cards = new Map() // content id -> controller
let originalBodyUserSelect

/** Returns the controller of a hover card by content id or template. */
export function getHoverCard(ref) {
  if (typeof ref === "string") return cards.get(ref)
  if (ref instanceof Element) {
    const id = ref.getAttribute("data-content-id")
    return id ? cards.get(id) : undefined
  }
  return undefined
}

const getState = (open) => (open ? "open" : "closed")

// Radix's getTabbableNodes.
function getTabbableNodes(container) {
  const nodes = []
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_ELEMENT, {
    acceptNode: (node) => (node.tabIndex >= 0 ? NodeFilter.FILTER_ACCEPT : NodeFilter.FILTER_SKIP),
  })
  while (walker.nextNode()) nodes.push(walker.currentNode)
  return nodes
}

const excludeTouch = (handler) => (event) => (event.pointerType === "touch" ? undefined : handler())


// Finds the elements carrying marker attribute name with value, removes the
// attribute and returns them (other components' markers stay in place).
function takeOwnMarked(name, value) {
  return [...document.querySelectorAll(`[${name}="${CSS.escape(value)}"]`)].map((el) => {
    el.removeAttribute(name)
    return el
  })
}

function createHoverCard(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const openDelay = Number(tpl.getAttribute("data-open-delay") ?? 700)
  const closeDelay = Number(tpl.getAttribute("data-close-delay") ?? 300)
  const trigger = takeOwnMarked("data-radix-hover-card-trigger", contentId)[0] ?? null

  let open = false
  let openTimer = 0
  let closeTimer = 0
  let hasSelection = false
  let isPointerDownOnContent = false
  let mounted = null

  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    trigger?.setAttribute("data-state", getState(open))
    if (open) {
      if (mounted) {
        mounted.content.setAttribute("data-state", "open")
        mounted.rendered()
        mounted.presence.setPresent(true)
      } else mountContent()
    } else if (mounted) {
      mounted.content.setAttribute("data-state", "closed")
      mounted.rendered()
      mounted.presence.setPresent(false)
    }
  }

  const handleOpen = () => {
    clearTimeout(closeTimer)
    openTimer = window.setTimeout(() => setOpen(true), openDelay)
  }
  const handleClose = () => {
    clearTimeout(openTimer)
    if (!hasSelection && !isPointerDownOnContent) {
      closeTimer = window.setTimeout(() => setOpen(false), closeDelay)
    }
  }
  const handleDismiss = () => setOpen(false)

  // HoverCardTrigger handlers.
  const triggerHandlers = [
    ["pointerenter", excludeTouch(handleOpen)],
    ["pointerleave", excludeTouch(handleClose)],
    ["focus", handleOpen],
    ["blur", handleClose],
    ["touchstart", (event) => event.preventDefault()],
  ].map(([type, fn]) => [type, (event) => !event.defaultPrevented && fn(event)])
  if (trigger) for (const [type, fn] of triggerHandlers) trigger.addEventListener(type, fn)

  function mountContent() {
    const nodes = instantiate(tpl)
    const wrapper = nodes[0]
    const content = wrapper.firstElementChild
    const state = { wrapper, content, cleanups: [], presence: null, rendered: () => {} }
    mounted = state
    const unmount = mount(nodes)

    // useEffect without deps: after every render, the content's tabbable
    // descendants leave the tab order.
    const removeTabbables = () => getTabbableNodes(content).forEach((t) => t.setAttribute("tabindex", "-1"))

    // containSelection state.
    let containSelection = false
    let restoreBody = null
    const setContainSelection = (on) => {
      if (on === containSelection) return
      containSelection = on
      content.style.userSelect = on ? "text" : ""
      content.style.webkitUserSelect = on ? "text" : ""
      restoreBody?.()
      restoreBody = null
      if (on) {
        const body = document.body
        originalBodyUserSelect = body.style.userSelect || body.style.webkitUserSelect
        body.style.userSelect = "none"
        body.style.webkitUserSelect = "none"
        restoreBody = () => {
          body.style.userSelect = originalBodyUserSelect
          body.style.webkitUserSelect = originalBodyUserSelect
        }
      }
      removeTabbables()
    }
    state.rendered = removeTabbables

    // HoverCardContent's handlers (composed with the user's).
    const onEnter = excludeTouch(handleOpen)
    const onLeave = excludeTouch(handleClose)
    const onPointerDown = (event) => {
      if (event.defaultPrevented) return
      if (event.currentTarget.contains(event.target)) setContainSelection(true)
      hasSelection = false
      isPointerDownOnContent = true
    }
    const guard = (fn) => (event) => !event.defaultPrevented && fn(event)
    const enter = guard(onEnter)
    const leave = guard(onLeave)
    content.addEventListener("pointerenter", enter)
    content.addEventListener("pointerleave", leave)
    content.addEventListener("pointerdown", onPointerDown)
    state.cleanups.push(() => {
      content.removeEventListener("pointerenter", enter)
      content.removeEventListener("pointerleave", leave)
      content.removeEventListener("pointerdown", onPointerDown)
    })

    // Effects, children first: PopperContent, DismissableLayer,
    // HoverCardContentImpl.
    const popper = createPopper(wrapper, { anchor: trigger })
    state.cleanups.push(() => popper.destroy())

    const layer = createDismissableLayer(content, {
      disableOutsidePointerEvents: false,
      parent: trigger,
      onFocusOutside: (event) => event.preventDefault(),
      onDismiss: handleDismiss,
    })
    state.cleanups.push(() => layer.destroy())

    state.cleanups.push(() => {
      restoreBody?.()
      restoreBody = null
    })
    const handlePointerUp = () => {
      setContainSelection(false)
      isPointerDownOnContent = false
      setTimeout(() => {
        const sel = document.getSelection()?.toString() !== ""
        if (sel) hasSelection = true
      })
    }
    document.addEventListener("pointerup", handlePointerUp)
    state.cleanups.push(() => {
      document.removeEventListener("pointerup", handlePointerUp)
      hasSelection = false
      isPointerDownOnContent = false
    })
    removeTabbables()

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

  const controller = {
    get open() {
      return open
    },
    setOpen,
    dispose() {
      clearTimeout(openTimer)
      clearTimeout(closeTimer)
      if (trigger) for (const [type, fn] of triggerHandlers) trigger.removeEventListener(type, fn)
      if (mounted) {
        const state = mounted
        state.presence.destroy()
        state.wrapper.remove()
        for (const c of state.cleanups.reverse()) c()
        mounted = null
      }
      if (cards.get(contentId) === controller) cards.delete(contentId)
    },
  }
  cards.set(contentId, controller)
  if (tpl.hasAttribute("data-default-open")) openOnLoad(tpl, () => setOpen(true))
  return controller
}

define("hover-card", {
  selector: "template[data-radix-hover-card]",
  mount(tpl) {
    const c = createHoverCard(tpl)
    return () => c.dispose()
  },
})
