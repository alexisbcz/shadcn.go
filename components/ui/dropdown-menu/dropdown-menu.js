// dropdown-menu.js: Radix DropdownMenu behavior (@radix-ui/react-dropdown-menu
// 2.1.16) for the server-rendered markup of components/ui/dropdown-menu, on
// the Menu primitive (runtime/js/menu.js).
//
// Markup contract: <template data-radix-dropdown-menu data-content-id
// data-trigger-id ...> (see menu.js) holds the content; the trigger is the
// element whose id is data-trigger-id.
//
// Scripts can observe and drive a menu (and set data-anchor, a selector,
// on the template to anchor the content on another element):
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)  // cancelable
//   import { getDropdownMenu } from "/_shadcn/dropdown-menu.js"
//   getDropdownMenu(contentIdOrTemplateOrTrigger)?.setOpen(true)   // .open, .toggle()

import { define, emit } from "./core.js"
import { createMenu } from "./menu.js"
import { openOnLoad } from "./popper.js"

const menus = new Map() // content id -> controller

/** Returns the controller of a dropdown menu by content id, template or trigger. */
export function getDropdownMenu(ref) {
  if (typeof ref === "string") return menus.get(ref)
  if (ref instanceof Element) {
    if (ref.hasAttribute("data-content-id")) return menus.get(ref.getAttribute("data-content-id"))
    for (const m of menus.values()) if (m.trigger === ref) return m
  }
  return undefined
}

function createDropdownMenu(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const triggerId = tpl.getAttribute("data-trigger-id")
  const modal = tpl.getAttribute("data-modal") !== "false"
  const trigger = () => document.getElementById(triggerId)
  let open = false
  let hasInteractedOutside = false

  const menu = createMenu(tpl, {
    modal,
    // A data-anchor selector on the template anchors the content elsewhere
    // (Base UI Menu's anchor prop); read when the menu opens.
    anchor: () => (tpl.dataset.anchor && document.querySelector(tpl.dataset.anchor)) || trigger(),
    layerParent: trigger,
    onOpenChange: (next) => setOpen(next),
    onCloseAutoFocus(event) {
      if (!hasInteractedOutside) trigger()?.focus()
      hasInteractedOutside = false
      event.preventDefault()
    },
    onInteractOutside(event) {
      const originalEvent = event.detail.originalEvent
      const ctrlLeftClick = originalEvent.button === 0 && originalEvent.ctrlKey === true
      const isRightClick = originalEvent.button === 2 || ctrlLeftClick
      if (!modal || isRightClick) hasInteractedOutside = true
    },
  })

  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    const t = trigger()
    if (t) {
      t.setAttribute("aria-expanded", String(open))
      if (open) t.setAttribute("aria-controls", contentId)
      else t.removeAttribute("aria-controls")
      t.setAttribute("data-state", open ? "open" : "closed")
    }
    menu.setOpen(open)
  }

  // DropdownMenuTrigger handlers (composed after the user's).
  const disabled = (t) => t.hasAttribute("data-disabled")
  const onPointerDown = (event) => {
    const t = event.currentTarget
    if (event.defaultPrevented || disabled(t)) return
    if (event.button === 0 && event.ctrlKey === false) {
      const wasOpen = open
      setOpen(!open)
      if (!wasOpen) event.preventDefault()
    }
  }
  const onKeyDown = (event) => {
    const t = event.currentTarget
    if (event.defaultPrevented || disabled(t)) return
    if (["Enter", " "].includes(event.key)) setOpen(!open)
    if (event.key === "ArrowDown") setOpen(true)
    if (["Enter", " ", "ArrowDown"].includes(event.key)) event.preventDefault()
  }
  let wired = null
  const wire = () => {
    const t = trigger()
    if (t === wired) return
    unwire()
    if (!t) return
    t.addEventListener("pointerdown", onPointerDown)
    t.addEventListener("keydown", onKeyDown)
    wired = t
  }
  const unwire = () => {
    wired?.removeEventListener("pointerdown", onPointerDown)
    wired?.removeEventListener("keydown", onKeyDown)
    wired = null
  }
  wire()

  const controller = {
    get open() {
      return open
    },
    get trigger() {
      return trigger()
    },
    setOpen,
    toggle: () => setOpen(!open),
    dispose() {
      unwire()
      menu.destroy()
      if (menus.get(contentId) === controller) menus.delete(contentId)
    },
  }
  menus.set(contentId, controller)
  if (tpl.hasAttribute("data-default-open")) {
    // The server rendered the trigger open already.
    open = false
    openOnLoad(tpl, () => setOpen(true))
  }
  return controller
}

define("dropdown-menu", {
  selector: "template[data-radix-dropdown-menu]",
  mount(tpl) {
    const m = createDropdownMenu(tpl)
    return () => m.dispose()
  },
})
