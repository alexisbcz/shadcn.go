// radio-group.js: Radix RadioGroup (@radix-ui/react-radio-group 1.3.8).
//
// - Clicking an unchecked item (or Space) checks it; Enter does nothing.
// - Arrow keys move focus (roving focus, loop unless data-loop="false") and
//   check the newly focused item; Tab enters the group on the checked item.
// - The indicator mounts when an item is checked and unmounts (after its
//   exit animation, if any: Radix Presence) when it is unchecked.
// - Each item's hidden <input type="radio"> (Radix RadioBubbleInput) is kept
//   only inside a form; its checked state follows the item and a bubbling
//   "click" is dispatched on it when it changes, so forms see the change.
// - "shadcn:value-change" (cancelable, detail: { value }) is dispatched on
//   the group before every change.
import { define, emit } from "./core.js"
import { createRovingFocusGroup } from "./roving-focus.js"
import { createPresence } from "./presence.js"

const ITEM = '[data-slot="radio-group-item"]'
const ARROW_KEYS = ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"]

// Radix props with no DOM trace (loop, rovingFocus, activationMode) are
// server-rendered as data attributes only when they differ from the default;
// takeOption reads one and removes it, so the DOM matches React's, and
// remembers it in case the element is mounted again.
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

define("radio-group", {
  selector: '[data-slot="radio-group"]',
  mount(root) {
    const cleanups = []
    const on = (el, type, fn, opts) => {
      el.addEventListener(type, fn, opts)
      cleanups.push(() => el.removeEventListener(type, fn, opts))
    }
    const items = () => [...root.querySelectorAll(ITEM)].filter((el) => el.closest('[data-slot="radio-group"]') === root)
    const state = new Map() // item -> { input, template, indicator, presence, stopped }

    for (const item of items()) {
      const s = { input: null, template: null, indicator: null, presence: null, stopped: false }
      state.set(item, s)
      const next = item.nextElementSibling
      if (next?.matches('input[type="radio"][aria-hidden="true"]')) {
        const isFormControl = !!item.getAttribute("form") || !!item.closest("form")
        if (isFormControl) {
          s.input = next
          syncInputSize(item, next, cleanups)
        } else next.remove()
      }
      const tpl = item.querySelector(':scope > template[data-slot="radio-group-indicator-template"]')
      if (tpl) {
        s.template = tpl.content
        tpl.remove()
      }
      const ind = item.querySelector(':scope > [data-slot="radio-group-indicator"]')
      if (ind) {
        s.indicator = ind
        s.presence = createPresence(ind, { present: true, onExitComplete: () => ind.remove() })
      }
    }

    const isChecked = (item) => item.getAttribute("aria-checked") === "true"

    const setChecked = (item, checked) => {
      const s = state.get(item)
      if (!s) return
      const was = isChecked(item)
      item.setAttribute("aria-checked", String(checked))
      item.setAttribute("data-state", checked ? "checked" : "unchecked")
      // Indicator (Presence present={checked}).
      if (checked) {
        if (!s.indicator || !s.indicator.isConnected) {
          const node = s.indicator ?? s.template?.firstElementChild?.cloneNode(true)
          if (node) {
            node.setAttribute("data-state", "checked")
            item.appendChild(node)
            s.indicator = node
            s.presence?.destroy()
            s.presence = createPresence(node, { present: true, onExitComplete: () => node.remove() })
          }
        } else {
          s.indicator.setAttribute("data-state", "checked")
          s.presence?.setPresent(true)
        }
      } else if (s.indicator?.isConnected) {
        s.indicator.setAttribute("data-state", "unchecked")
        if (s.presence) s.presence.setPresent(false)
        else s.indicator.remove()
      }
      // Bubble input render: React updates its defaultChecked (the
      // `checked` attribute) in the commit phase.
      if (s.input && was !== checked) s.input.defaultChecked = checked
      return was !== checked
    }

    // Bubble input effect (after every item's commit, like React effects):
    // set the checked property and dispatch a click forms can observe.
    const bubble = (item) => {
      const s = state.get(item)
      if (!s?.input) return
      s.input.checked = isChecked(item)
      s.input.dispatchEvent(new Event("click", { bubbles: !s.stopped }))
    }

    // Returns the items whose checked state changed; their bubble inputs
    // fire once the triggering click handler is done, as React's effects do.
    const setValue = (value) => {
      if (!emit(root, "shadcn:value-change", { value })) return []
      return items().filter((item) => setChecked(item, item.value === value))
    }

    let isArrowKeyPressed = false
    on(document, "keydown", (e) => {
      if (ARROW_KEYS.includes(e.key)) isArrowKeyPressed = true
    })
    on(document, "keyup", () => (isArrowKeyPressed = false))

    for (const item of items()) {
      on(item, "click", (event) => {
        const changed = isChecked(item) ? [] : setValue(item.value)
        const s = state.get(item)
        if (s?.input) {
          s.stopped = event.cancelBubble
          if (!s.stopped) event.stopPropagation()
        }
        changed.forEach(bubble)
      })
      on(item, "keydown", (event) => {
        if (event.key === "Enter") event.preventDefault()
      })
      on(item, "focusin", () => {
        if (isArrowKeyPressed) item.click()
      })
    }

    const roving = createRovingFocusGroup(root, {
      orientation: root.getAttribute("aria-orientation") || undefined,
      loop: takeOption(root, "data-loop") !== "false",
      focusable: (item) => !item.disabled,
      active: isChecked,
    })

    return () => {
      roving.destroy()
      for (const s of state.values()) s.presence?.destroy()
      cleanups.forEach((fn) => fn())
    }
  },
})

// Radix sizes the bubble input like its control (useSize).
function syncInputSize(control, input, cleanups) {
  const apply = () => {
    input.style.width = control.offsetWidth + "px"
    input.style.height = control.offsetHeight + "px"
  }
  apply()
  if (typeof ResizeObserver === "undefined") return
  const ro = new ResizeObserver(apply)
  ro.observe(control, { box: "border-box" })
  cleanups.push(() => ro.disconnect())
}
