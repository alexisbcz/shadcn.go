// The combobox examples keep the picked value in React state: the trigger
// shows its label, combobox-demo marks it with a check in the list (live,
// and in the list the popover mounts on its next opening), and picking an
// item closes the popover (or drawer). Mounted on the trigger
// (data-example-combobox="demo" | "status").
//
// combobox-responsive also swaps its Popover for the Drawer version kept in
// template[data-example-combobox-media] while the media query does not match
// (useMediaQuery).
import { define } from "./core.js"
import { getPopover } from "./popover.js"
import { getDrawer } from "./drawer.js"
import { getDropdownMenu } from "./dropdown-menu.js"

define("example-combobox", {
  selector: "[data-example-combobox]",
  mount(trigger) {
    const kind = trigger.getAttribute("data-example-combobox")
    const options = (trigger.getAttribute("data-example-combobox-options") ?? "")
      .split("\n")
      .map((line) => line.split("\t"))
    trigger.removeAttribute("data-example-combobox")
    trigger.removeAttribute("data-example-combobox-options")
    const contentId = trigger.getAttribute("aria-controls")
    const icon = trigger.querySelector("svg")
    const labelOf = (value) => options.find(([v]) => v === value)?.[1]
    let current = ""

    // combobox-demo: <Check className={cn("ml-auto", value === framework.value ? "opacity-100" : "opacity-0")} />
    const markChecks = (root) => {
      for (const item of root.querySelectorAll("[cmdk-item]")) {
        const check = item.querySelector("svg")
        if (!check) continue
        const on = item.getAttribute("data-value") === current
        check.classList.remove(on ? "opacity-0" : "opacity-100")
        check.classList.add(on ? "opacity-100" : "opacity-0")
      }
    }

    const onSelect = (e) => {
      const content = e.target.closest?.('[data-slot="popover-content"], [data-slot="drawer-content"]')
      if (!content || content.id !== contentId) return
      const value = e.detail.value
      if (kind === "demo") {
        current = value === current ? "" : value
        trigger.replaceChildren(current ? labelOf(current) : "Select framework...", icon)
        markChecks(content)
        const tpl = document.querySelector(`template[data-content-id="${CSS.escape(contentId)}"]`)
        if (tpl) markChecks(tpl.content)
      } else {
        const label = labelOf(value)
        trigger.replaceChildren(label ?? "+ Set status")
      }
      ;(getPopover(contentId) ?? getDrawer(contentId))?.setOpen(false)
    }
    document.addEventListener("shadcn:select", onSelect)
    return () => document.removeEventListener("shadcn:select", onSelect)
  },
})

define("example-combobox-media", {
  selector: "template[data-example-combobox-media]",
  mount(tpl) {
    const parent = tpl.parentElement
    if (!parent) return
    const query = window.matchMedia(tpl.getAttribute("data-example-combobox-media"))
    const matching = [...parent.childNodes].slice(0, [...parent.childNodes].indexOf(tpl))
    let other = null
    const show = (matches) => {
      if (matches) {
        if (!other) return
        other.forEach((n) => n.remove())
        other = null
        matching.forEach((n) => parent.insertBefore(n, tpl))
      } else {
        if (other) return
        matching.forEach((n) => n.remove())
        other = [...tpl.content.cloneNode(true).childNodes]
        other.forEach((n) => parent.insertBefore(n, tpl))
      }
    }
    const onChange = (e) => show(e.matches)
    query.addEventListener("change", onChange)
    queueMicrotask(() => show(query.matches))
    return () => query.removeEventListener("change", onChange)
  },
})

// combobox-dropdown-menu: the label filter input has autoFocus (React
// focuses it when the submenu mounts), and picking a label sets it and
// closes the menu.
define("example-combobox-autofocus", {
  selector: "[data-example-combobox-autofocus]",
  mount(input) {
    input.removeAttribute("data-example-combobox-autofocus")
    input.focus()
  },
})

let labelEl = null
define("example-combobox-label", {
  selector: "[data-example-combobox-label]",
  mount(el) {
    el.removeAttribute("data-example-combobox-label")
    labelEl = el
  },
})

document.addEventListener("shadcn:select", (e) => {
  const item = e.target
  if (!(item instanceof Element) || !item.closest("[data-slot=dropdown-menu-sub-content]") || !item.hasAttribute("cmdk-item")) return
  if (labelEl) labelEl.textContent = e.detail.value
  const tpl = [...document.querySelectorAll("template[data-content-id]")].find((t) => getDropdownMenu(t))
  ;(tpl && getDropdownMenu(tpl))?.setOpen(false)
})
