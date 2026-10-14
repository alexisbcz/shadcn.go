// dialog.js: Radix Dialog behavior (@radix-ui/react-dialog 1.1.15, and the
// AlertDialog variant of @radix-ui/react-alert-dialog 1.1.15) for the
// server-rendered markup of components/ui/dialog/primitive (used by the
// shadcn Dialog, AlertDialog and Sheet).
//
// Markup contract (rendered by the Go package):
//   - <template data-radix-dialog="" | "alert" data-content-id="<id>"
//       [data-modal="false"] [data-default-open] [data-in-place]>
//     holds the portalled parts (overlay, content) exactly as React renders
//     them when open. Its children are cloned into <body> on open (or next
//     to the template with data-in-place).
//   - Triggers are the elements with aria-controls="<id>" (Radix's
//     DialogTrigger always renders it).
//   - Inside the template, markers the behavior removes when cloning:
//     data-radix-dialog-overlay (the overlay), data-radix-dialog-close="<id>"
//     (Close buttons), data-radix-dialog-cancel (AlertDialog Cancel). Close
//     buttons outside the content carry the same marker, removed on mount.
//
// Scripts can observe and drive a dialog:
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getDialog } from "/_shadcn/dialog.js"
//   getDialog(contentIdOrTemplateOrTrigger)?.setOpen(true)   // .open, .toggle()

import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createFocusScope } from "./focus-scope.js"
import { acquireFocusGuards } from "./focus-guards.js"
import { lockScroll } from "./scroll-lock.js"
import { hideOthers } from "./hide-others.js"

const dialogs = new Map() // content id -> controller

/** Returns the controller of a dialog by content id, template or trigger. */
export function getDialog(ref) {
  if (typeof ref === "string") return dialogs.get(ref)
  if (ref instanceof Element) {
    const id = ref.getAttribute("data-content-id") ?? ref.getAttribute("aria-controls")
    return id ? dialogs.get(id) : undefined
  }
  return undefined
}

const getState = (open) => (open ? "open" : "closed")

function createDialog(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const alert = tpl.getAttribute("data-radix-dialog") === "alert"
  const modal = alert || tpl.getAttribute("data-modal") !== "false"
  const inPlace = tpl.hasAttribute("data-in-place")
  let open = false
  let mounted = null // { nodes, content, overlay, presences, cleanups... } while present

  const triggers = () => [...document.querySelectorAll(`[aria-controls="${CSS.escape(contentId)}"]`)]

  // Close buttons outside the portalled content.
  const outsideCloses = [...document.querySelectorAll(`[data-radix-dialog-close="${CSS.escape(contentId)}"]`)].map((el) => {
    el.removeAttribute("data-radix-dialog-close")
    el.removeAttribute("data-radix-dialog-cancel")
    return { el }
  })
  const onCloseClick = (event) => {
    if (!event.defaultPrevented) setOpen(false)
  }
  outsideCloses.forEach(({ el }) => el.addEventListener("click", onCloseClick))

  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    // Re-render of the parts reading context.open: triggers first (they
    // come first in the tree), then the portalled parts.
    for (const t of triggers()) {
      t.setAttribute("aria-expanded", String(open))
      t.setAttribute("data-state", getState(open))
    }
    if (open) {
      if (mounted) reshow()
      else mountParts()
    } else if (mounted) {
      for (const el of mounted.stateful) el.setAttribute("data-state", "closed")
      mounted.focusScope?.setTrapped(false)
      for (const p of mounted.presences) p.setPresent(false)
    }
  }

  // Open again while the exit animation runs: Presence goes back to
  // "mounted", the same elements stay.
  function reshow() {
    for (const el of mounted.stateful) el.setAttribute("data-state", "open")
    mounted.focusScope?.setTrapped(modal)
    for (const p of mounted.presences) p.setPresent(true)
  }

  function mountParts() {
    const nodes = instantiate(tpl)
    const overlays = new Set(takeMarked(nodes, "data-radix-dialog-overlay").map((m) => m.el))
    const closes = takeMarked(nodes, "data-radix-dialog-close").map((m) => m.el)
    const cancel = takeMarked(nodes, "data-radix-dialog-cancel").map((m) => m.el)[0]
    const content = nodes.flatMap((n) => (n.id === contentId ? [n] : [...n.querySelectorAll(`[id="${CSS.escape(contentId)}"]`)]))[0]
    const overlay = nodes.find((n) => overlays.has(n))
    const state = { nodes, stateful: [overlay, content].filter(Boolean), presences: [], cleanups: new Map(), focusScope: null, live: new Set(nodes) }
    mounted = state

    closes.forEach((el) => el.addEventListener("click", onCloseClick))

    // Portal: every child of the portal is mounted at the end of <body>.
    if (inPlace) tpl.after(...nodes)
    else mount(nodes)

    // Effects, in React's order: the overlay subtree, then the content's.
    if (overlay) {
      // DialogOverlayImpl: <RemoveScroll allowPinchZoom shards={[contentRef]}>
      const release = lockScroll(overlay, { allowPinchZoom: true, shards: [() => content] })
      state.cleanups.set(overlay, [release])
    }
    if (content) state.cleanups.set(content, mountContent(content, cancel, state))

    // Presence of each portalled child.
    for (const node of nodes) {
      const presence = createPresence(node, {
        present: true,
        onExitComplete: () => unmountNode(state, node, presence),
      })
      state.presences.push(presence)
    }
  }

  function mountContent(content, cancel, state) {
    const cleanups = []
    const trigger = triggers().pop()
    let hasInteractedOutside = false
    let hasPointerDownOutside = false

    // DialogContentImpl > FocusScope > DismissableLayer: the layer's
    // effects run first.
    const layer = createDismissableLayer(content, {
      disableOutsidePointerEvents: modal,
      parent: trigger,
      onDismiss: () => setOpen(false),
      onPointerDownOutside: (event) => {
        if (alert) return event.preventDefault()
        if (!modal) return
        const originalEvent = event.detail.originalEvent
        const ctrlLeftClick = originalEvent.button === 0 && originalEvent.ctrlKey === true
        const isRightClick = originalEvent.button === 2 || ctrlLeftClick
        if (isRightClick) event.preventDefault()
      },
      onFocusOutside: (event) => {
        if (modal) event.preventDefault()
      },
      onInteractOutside: (event) => {
        if (alert) return event.preventDefault()
        if (modal) return
        // DialogContentNonModal
        if (!event.defaultPrevented) {
          hasInteractedOutside = true
          if (event.detail.originalEvent.type === "pointerdown") hasPointerDownOutside = true
        }
        const target = event.target
        const targetIsTrigger = triggers().some((t) => t.contains(target))
        if (targetIsTrigger) event.preventDefault()
        if (event.detail.originalEvent.type === "focusin" && hasPointerDownOutside) event.preventDefault()
      },
    })
    cleanups.push(() => layer.destroy())

    const focusScope = createFocusScope(content, {
      loop: true,
      trapped: modal && open,
      onMountAutoFocus: (event) => {
        if (alert) {
          event.preventDefault()
          cancel?.focus({ preventScroll: true })
        }
      },
      onUnmountAutoFocus: (event) => {
        if (modal) {
          event.preventDefault()
          triggers().pop()?.focus()
          return
        }
        if (!event.defaultPrevented) {
          if (!hasInteractedOutside) triggers().pop()?.focus()
          event.preventDefault()
        }
        hasInteractedOutside = false
        hasPointerDownOutside = false
      },
    })
    state.focusScope = focusScope
    cleanups.push(() => focusScope.destroy())

    // DialogContentImpl: useFocusGuards(); DialogContentModal: hideOthers.
    const releaseGuards = acquireFocusGuards()
    cleanups.push(releaseGuards)
    if (modal) cleanups.push(hideOthers(content))
    // Unmount runs parent effects' cleanups first.
    return cleanups.reverse()
  }

  function unmountNode(state, node, presence) {
    presence.destroy()
    node.remove()
    state.live.delete(node)
    for (const c of state.cleanups.get(node) ?? []) c()
    state.cleanups.delete(node)
    if (state.live.size === 0 && mounted === state) mounted = null
  }

  const controller = {
    get open() {
      return open
    },
    setOpen,
    toggle: () => setOpen(!open),
    dispose() {
      outsideCloses.forEach(({ el }) => el.removeEventListener("click", onCloseClick))
      if (mounted) {
        const state = mounted
        for (const [i, node] of state.nodes.entries()) if (state.live.has(node)) unmountNode(state, node, state.presences[i])
      }
      if (dialogs.get(contentId) === controller) dialogs.delete(contentId)
    },
  }
  dialogs.set(contentId, controller)
  if (tpl.hasAttribute("data-default-open")) {
    // The server rendered the triggers open already.
    open = false
    setOpen(true)
  }
  return controller
}

// DialogTrigger: onClick toggles (composed after the user's handlers).
document.addEventListener("click", (event) => {
  if (event.defaultPrevented) return
  const t = event.target instanceof Element ? event.target.closest("[aria-controls]") : null
  if (!t) return
  const d = dialogs.get(t.getAttribute("aria-controls"))
  if (d) d.toggle()
})

define("dialog", {
  selector: "template[data-radix-dialog]",
  mount(tpl) {
    const d = createDialog(tpl)
    return () => d.dispose()
  },
})
