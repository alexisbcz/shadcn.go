// context-menu.js: Radix ContextMenu behavior (@radix-ui/react-context-menu
// 2.2.16) for the server-rendered markup of components/ui/context-menu, on
// the Menu primitive (runtime/js/menu.js).
//
// Markup contract: <template data-radix-context-menu data-content-id ...>
// (see menu.js) holds the content; the trigger carries the marker
// data-radix-context-menu-trigger="<content id>", removed on mount. The menu
// opens at the pointer on contextmenu (right click, or a 700ms long press
// with touch or pen), positioned against a 0x0 virtual anchor there.
//
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getContextMenu } from "/_shadcn/context-menu.js"
//   getContextMenu(contentIdOrTemplate)?.setOpen(false)

import { define, emit } from "./core.js"
import { createMenu } from "./menu.js"

const menus = new Map() // content id -> controller

/** Returns the controller of a context menu by content id or template. */
export function getContextMenu(ref) {
  if (typeof ref === "string") return menus.get(ref)
  if (ref instanceof Element) return menus.get(ref.getAttribute("data-content-id") ?? "")
  return undefined
}

const whenTouchOrPen = (handler) => (event) => (event.pointerType !== "mouse" ? handler(event) : undefined)

function createContextMenu(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const modal = tpl.getAttribute("data-modal") !== "false"
  const triggers = [...document.querySelectorAll(`[data-radix-context-menu-trigger="${CSS.escape(contentId)}"]`)]
  triggers.forEach((t) => t.removeAttribute("data-radix-context-menu-trigger"))
  let open = false
  let hasInteractedOutside = false
  let point = { x: 0, y: 0 }
  const virtualAnchor = { getBoundingClientRect: () => DOMRect.fromRect({ width: 0, height: 0, ...point }) }

  const menu = createMenu(tpl, {
    modal,
    anchor: () => virtualAnchor,
    layerParent: () => triggers[triggers.length - 1] ?? null,
    onOpenChange: (next) => setOpen(next),
    onCloseAutoFocus(event) {
      if (!event.defaultPrevented && hasInteractedOutside) event.preventDefault()
      hasInteractedOutside = false
    },
    onInteractOutside(event) {
      if (!event.defaultPrevented && !modal) hasInteractedOutside = true
    },
  })

  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    for (const t of triggers) t.setAttribute("data-state", open ? "open" : "closed")
    menu.setOpen(open)
  }

  const cleanups = []
  for (const trigger of triggers) {
    let longPressTimer = 0
    const clearLongPress = () => window.clearTimeout(longPressTimer)
    const disabled = () => trigger.hasAttribute("data-disabled")
    const handleOpen = (event) => {
      point = { x: event.clientX, y: event.clientY }
      setOpen(true)
    }
    const onContextMenu = (event) => {
      if (disabled() || event.defaultPrevented) return
      clearLongPress()
      handleOpen(event)
      event.preventDefault()
    }
    const onPointerDown = whenTouchOrPen((event) => {
      if (disabled() || event.defaultPrevented) return
      clearLongPress()
      longPressTimer = window.setTimeout(() => handleOpen(event), 700)
    })
    const onClear = whenTouchOrPen((event) => {
      if (disabled() || event.defaultPrevented) return
      clearLongPress()
    })
    trigger.addEventListener("contextmenu", onContextMenu)
    trigger.addEventListener("pointerdown", onPointerDown)
    trigger.addEventListener("pointermove", onClear)
    trigger.addEventListener("pointercancel", onClear)
    trigger.addEventListener("pointerup", onClear)
    cleanups.push(() => {
      clearLongPress()
      trigger.removeEventListener("contextmenu", onContextMenu)
      trigger.removeEventListener("pointerdown", onPointerDown)
      trigger.removeEventListener("pointermove", onClear)
      trigger.removeEventListener("pointercancel", onClear)
      trigger.removeEventListener("pointerup", onClear)
    })
  }

  const controller = {
    get open() {
      return open
    },
    setOpen,
    dispose() {
      cleanups.forEach((c) => c())
      menu.destroy()
      if (menus.get(contentId) === controller) menus.delete(contentId)
    },
  }
  menus.set(contentId, controller)
  return controller
}

define("context-menu", {
  selector: "template[data-radix-context-menu]",
  mount(tpl) {
    const m = createContextMenu(tpl)
    return () => m.dispose()
  },
})
