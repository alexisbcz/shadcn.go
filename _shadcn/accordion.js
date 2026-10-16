// accordion.js: Radix Accordion (@radix-ui/react-accordion 1.2.12).
//
// - Each item is a Radix Collapsible: clicking its trigger (Enter/Space on
//   the <button>) asks to open or close it. Single: opening an item closes
//   the others; the open item closes only when collapsible (otherwise its
//   trigger has aria-disabled="true"). Multiple: items toggle independently.
// - Keyboard (on the root, from an enabled trigger of this accordion): Home,
//   End, ArrowDown/ArrowUp (vertical) or ArrowRight/ArrowLeft (horizontal,
//   per the Dir option, ltr by default) move focus among enabled triggers, wrapping.
// - Contents are collapsible-content.js instances: presence, size
//   variables, children mounted only while open.
// - "shadcn:value-change" (cancelable; detail: { value }, a string for
//   single, an array for multiple) is dispatched on the root before every
//   change.
import { define, emit } from "./core.js"
import { createCollection } from "./collection.js"
import { createCollapsibleContent } from "./collapsible-content.js"

const ACCORDION_KEYS = ["Home", "End", "ArrowDown", "ArrowUp", "ArrowLeft", "ArrowRight"]

// Radix props with no DOM trace, server-rendered as data attributes: read
// once and removed, so the DOM matches React's; remembered in case the
// element is mounted again.
const savedOptions = new WeakMap()
function takeOption(el, name) {
  let saved = savedOptions.get(el)
  if (!saved) savedOptions.set(el, (saved = {}))
  if (el.hasAttribute(name)) {
    saved[name] = el.getAttribute(name)
    el.removeAttribute(name)
  }
  return saved[name] ?? null
}

const state = (open) => (open ? "open" : "closed")

define("accordion", {
  selector: '[data-slot="accordion"]',
  mount(root) {
    const multiple = takeOption(root, "data-type") === "multiple"
    const collapsible = multiple || takeOption(root, "data-collapsible") !== null
    const localDir = takeOption(root, "data-dir")
    const orientation = root.getAttribute("data-orientation") === "horizontal" ? "horizontal" : "vertical"
    const collection = createCollection(root, { filter: (el) => el.closest('[data-slot="accordion"]') === root })

    const items = collection.getItems().map((trigger) => {
      const content = document.getElementById(trigger.getAttribute("aria-controls"))
      let item = trigger.parentElement
      while (item && item !== root && !(content && item.contains(content))) item = item.parentElement
      if (!item || item === root) item = trigger.parentElement?.parentElement
      const header = trigger.parentElement
      const value = takeOption(item, "data-value") ?? ""
      const open = item.getAttribute("data-state") === "open"
      let panel = null
      if (content) {
        const tpl = content.querySelector(':scope > template[data-slot="accordion-content-template"]')
        const fragment = tpl?.content ?? null
        tpl?.remove()
        panel = createCollapsibleContent(content, { open, fragment })
      }
      return { trigger, item, header, value, open, panel }
    })

    let value = items.filter((it) => it.open && it.value).map((it) => it.value)
    if (!multiple) value = value.slice(0, 1)

    const apply = () => {
      for (const it of items) {
        const open = !!it.value && value.includes(it.value)
        it.item.setAttribute("data-state", state(open))
        if (it.header !== it.item) it.header.setAttribute("data-state", state(open))
        it.trigger.setAttribute("aria-expanded", String(open))
        it.trigger.setAttribute("data-state", state(open))
        if (open && !collapsible) it.trigger.setAttribute("aria-disabled", "true")
        else it.trigger.removeAttribute("aria-disabled")
        if (open !== it.open) {
          it.open = open
          it.panel?.setOpen(open)
        }
      }
    }

    const setValue = (next) => {
      const same = next.length === value.length && next.every((v, i) => v === value[i])
      if (same) return
      if (!emit(root, "shadcn:value-change", { value: multiple ? [...next] : (next[0] ?? "") })) return
      value = next
      apply()
    }

    // CollapsibleTrigger onClick -> onOpenToggle -> onOpenChange(!open).
    const onClick = (event) => {
      if (event.defaultPrevented) return
      const it = items.find((x) => x.trigger === event.target || x.trigger.contains(event.target))
      if (!it || !collection.has(it.trigger)) return
      if (!it.open) {
        // onItemOpen
        setValue(multiple ? [...value, it.value] : [it.value])
      } else if (multiple) {
        setValue(value.filter((v) => v !== it.value))
      } else if (collapsible) {
        setValue([])
      }
    }

    const onKeyDown = (event) => {
      if (event.defaultPrevented || !ACCORDION_KEYS.includes(event.key)) return
      const target = event.target
      const triggerCollection = collection.getItems().filter((el) => !el.disabled)
      const triggerIndex = triggerCollection.indexOf(target)
      const triggerCount = triggerCollection.length
      if (triggerIndex === -1) return
      event.preventDefault()
      let nextIndex = triggerIndex
      const homeIndex = 0
      const endIndex = triggerCount - 1
      const moveNext = () => {
        nextIndex = triggerIndex + 1
        if (nextIndex > endIndex) nextIndex = homeIndex
      }
      const movePrev = () => {
        nextIndex = triggerIndex - 1
        if (nextIndex < homeIndex) nextIndex = endIndex
      }
      // useDirection(dir): the dir prop, else a DirectionProvider (none in
      // the port), else "ltr"; the DOM's dir attribute does not count.
      const isDirectionLTR = localDir !== "rtl"
      switch (event.key) {
        case "Home":
          nextIndex = homeIndex
          break
        case "End":
          nextIndex = endIndex
          break
        case "ArrowRight":
          if (orientation === "horizontal") isDirectionLTR ? moveNext() : movePrev()
          break
        case "ArrowDown":
          if (orientation === "vertical") moveNext()
          break
        case "ArrowLeft":
          if (orientation === "horizontal") isDirectionLTR ? movePrev() : moveNext()
          break
        case "ArrowUp":
          if (orientation === "vertical") movePrev()
          break
      }
      const clampedIndex = nextIndex % triggerCount
      triggerCollection[clampedIndex]?.focus()
    }

    root.addEventListener("click", onClick)
    root.addEventListener("keydown", onKeyDown)
    return () => {
      root.removeEventListener("click", onClick)
      root.removeEventListener("keydown", onKeyDown)
      items.forEach((it) => it.panel?.destroy())
      collection.destroy()
    }
  },
})
