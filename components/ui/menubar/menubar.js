// menubar.js: Radix Menubar behavior (@radix-ui/react-menubar 1.1.16) for
// the server-rendered markup of components/ui/menubar, on the Menu
// primitive (runtime/js/menu.js) and the roving focus group.
//
// Markup contract: the [role="menubar"] element (markers, removed on mount:
// data-radix-menubar-loop="false", data-radix-menubar-dir="rtl") contains,
// per menu, a trigger (id = the template's data-trigger-id) and a
// <template data-radix-menubar-menu data-value data-content-id
// data-trigger-id ...> holding the content (see menu.js).
//
//   menubarEl.addEventListener("shadcn:value-change", (e) => e.detail.value)  // cancelable
//   import { getMenubar } from "/_shadcn/menubar.js"
//   getMenubar(menubarEl)?.setValue("file")     // .value; "" closes

import { define, emit } from "./core.js"
import { createMenu } from "./menu.js"
import { createRovingFocusGroup, wrapArray } from "./roving-focus.js"
import { openOnLoad } from "./popper.js"

const menubars = new WeakMap() // menubar element -> controller

/** Returns the controller of a menubar element. */
export function getMenubar(el) {
  return menubars.get(el)
}

function createMenubar(root) {
  const loop = root.getAttribute("data-radix-menubar-loop") !== "false"
  const dir = root.getAttribute("data-radix-menubar-dir") === "rtl" ? "rtl" : "ltr"
  root.removeAttribute("data-radix-menubar-loop")
  root.removeAttribute("data-radix-menubar-dir")
  let value = ""
  const cleanups = []

  // RovingFocusGroup over the triggers (Collection.Slot on the menubar).
  const roving = createRovingFocusGroup(root, {
    orientation: "horizontal",
    loop,
    dir,
    focusable: (item) => !item.hasAttribute("data-disabled"),
  })
  cleanups.push(() => roving.destroy())

  const menus = []
  for (const tpl of root.querySelectorAll("template[data-radix-menubar-menu]")) {
    if (tpl.parentElement?.closest('[role="menubar"]') !== root) continue
    menus.push(createMenubarMenu(tpl))
  }
  const triggerItems = () => menus.filter((m) => m.trigger()).sort((a, b) => (a.trigger().compareDocumentPosition(b.trigger()) & Node.DOCUMENT_POSITION_FOLLOWING ? -1 : 1))

  function setValue(next) {
    if (next === value) return
    if (!emit(root, "shadcn:value-change", { value: next })) return
    value = next
    // React's commit runs the effect cleanups of the menus that unmount (a
    // closing menu without exit animation) before the new menu's effects.
    for (const m of menus) if (m.value !== value) m.sync()
    for (const m of menus) if (m.value === value) m.sync()
  }
  const onMenuOpen = (v) => {
    setValue(v)
    roving.setCurrentTabStop(menus.find((m) => m.value === v)?.trigger() ?? null)
  }
  const onMenuClose = () => setValue("")
  const onMenuToggle = (v) => {
    setValue(value ? "" : v)
    roving.setCurrentTabStop(menus.find((m) => m.value === v)?.trigger() ?? null)
  }

  function createMenubarMenu(tpl) {
    const mvalue = tpl.getAttribute("data-value")
    const contentId = tpl.getAttribute("data-content-id")
    const triggerId = tpl.getAttribute("data-trigger-id")
    const trigger = () => document.getElementById(triggerId)
    let open = false
    let wasKeyboardTriggerOpen = false
    let hasInteractedOutside = false

    const menu = createMenu(tpl, {
      modal: false,
      dir,
      anchor: trigger,
      layerParent: trigger,
      onOpenChange: (o) => {
        if (!o) onMenuClose()
      },
      onCloseAutoFocus(event) {
        const menubarOpen = Boolean(value)
        if (!menubarOpen && !hasInteractedOutside) trigger()?.focus()
        hasInteractedOutside = false
        event.preventDefault()
      },
      onFocusOutside(event) {
        const target = event.target
        const isMenubarTrigger = menus.some((m) => m.trigger()?.contains(target))
        if (isMenubarTrigger) event.preventDefault()
      },
      onInteractOutside() {
        hasInteractedOutside = true
      },
      onEntryFocus(event) {
        if (!wasKeyboardTriggerOpen) event.preventDefault()
      },
      onKeyDown(event) {
        if (!["ArrowRight", "ArrowLeft"].includes(event.key)) return
        const content = menu.content
        const target = event.target
        const targetIsSubTrigger = target.hasAttribute("data-radix-menubar-subtrigger")
        const isKeyDownInsideSubMenu = target.closest("[data-radix-menubar-content]") !== content
        const prevMenuKey = dir === "rtl" ? "ArrowRight" : "ArrowLeft"
        const isPrevKey = prevMenuKey === event.key
        const isNextKey = !isPrevKey
        if (isNextKey && targetIsSubTrigger) return
        if (isKeyDownInsideSubMenu && isPrevKey) return
        const items = triggerItems().filter((m) => !m.trigger().hasAttribute("data-disabled"))
        let candidateValues = items.map((m) => m.value)
        if (isPrevKey) candidateValues.reverse()
        const currentIndex = candidateValues.indexOf(mvalue)
        candidateValues = loop ? wrapArray(candidateValues, currentIndex + 1) : candidateValues.slice(currentIndex + 1)
        const [nextValue] = candidateValues
        if (nextValue) onMenuOpen(nextValue)
      },
    })

    // MenubarTrigger handlers (on the trigger: before the group's).
    const t0 = trigger()
    const disabled = () => trigger()?.hasAttribute("data-disabled")
    const onPointerDown = (event) => {
      if (event.defaultPrevented || disabled()) return
      if (event.button === 0 && event.ctrlKey === false) {
        const wasOpen = open
        onMenuOpen(mvalue)
        if (!wasOpen) event.preventDefault()
      }
    }
    const onPointerEnter = (event) => {
      if (event.defaultPrevented) return
      const menubarOpen = Boolean(value)
      if (menubarOpen && !open) {
        // React re-renders after the handler: the focus comes first, then
        // the new menu mounts (and its FocusScope takes focus).
        trigger()?.focus()
        onMenuOpen(mvalue)
      }
    }
    const onKeyDown = (event) => {
      if (event.defaultPrevented || disabled()) return
      // The flag is set in the same handler as the state change, before
      // React renders the menu (whose entry focus reads it).
      const isOpenKey = ["Enter", " ", "ArrowDown"].includes(event.key)
      if (isOpenKey) wasKeyboardTriggerOpen = true
      if (["Enter", " "].includes(event.key)) onMenuToggle(mvalue)
      if (event.key === "ArrowDown") onMenuOpen(mvalue)
      if (isOpenKey) event.preventDefault()
    }
    const onFocus = () => trigger()?.setAttribute("data-highlighted", "")
    const onBlur = () => trigger()?.removeAttribute("data-highlighted")
    if (t0) {
      t0.addEventListener("pointerdown", onPointerDown)
      t0.addEventListener("pointerenter", onPointerEnter)
      t0.addEventListener("keydown", onKeyDown)
      t0.addEventListener("focus", onFocus)
      t0.addEventListener("blur", onBlur)
      cleanups.push(() => {
        t0.removeEventListener("pointerdown", onPointerDown)
        t0.removeEventListener("pointerenter", onPointerEnter)
        t0.removeEventListener("keydown", onKeyDown)
        t0.removeEventListener("focus", onFocus)
        t0.removeEventListener("blur", onBlur)
      })
    }
    cleanups.push(() => menu.destroy())

    return {
      value: mvalue,
      trigger,
      sync() {
        const next = value === mvalue
        if (next === open) return
        open = next
        const t = trigger()
        if (t) {
          t.setAttribute("aria-expanded", String(open))
          if (open) t.setAttribute("aria-controls", contentId)
          else t.removeAttribute("aria-controls")
          t.setAttribute("data-state", open ? "open" : "closed")
        }
        if (!open) wasKeyboardTriggerOpen = false
        menu.setOpen(open)
      },
      defaultOpen: tpl.hasAttribute("data-default-open"),
      tpl,
    }
  }

  const controller = {
    get value() {
      return value
    },
    setValue: (v) => (v ? onMenuOpen(v) : onMenuClose()),
    dispose() {
      for (const c of cleanups.reverse()) c()
      menubars.delete(root)
    },
  }
  menubars.set(root, controller)
  const initial = menus.find((m) => m.defaultOpen)
  if (initial) {
    // defaultValue opens the menu without making its trigger the tab stop.
    openOnLoad(initial.tpl, () => setValue(initial.value))
  }
  return controller
}

define("menubar", {
  selector: '[role="menubar"]',
  mount(root) {
    const m = createMenubar(root)
    return () => m.dispose()
  },
})
