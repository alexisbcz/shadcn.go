// tabs.js: Radix Tabs (@radix-ui/react-tabs 1.1.13).
//
// - The list is a roving focus group (orientation of the tabs, loop unless
//   data-loop="false" on the list); Tab enters it on the selected trigger.
// - A trigger is selected on primary-button mousedown without Ctrl, on
//   Enter/Space, and, in automatic activation mode (default; the root has
//   data-activation-mode="manual" otherwise), when it gets focus.
// - Panels are mounted only while selected (Radix Presence): an inactive
//   panel is hidden and empty; on activation it gets a fresh copy of its
//   content from the template the server rendered.
// - "shadcn:value-change" (cancelable, detail: { value }) is dispatched on
//   the root before every change.
import { define, emit } from "./core.js"
import { createRovingFocusGroup } from "./roving-focus.js"
import { createPresence } from "./presence.js"

const own = (root, el) => el.closest('[data-slot="tabs"]') === root

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

define("tabs", {
  selector: '[data-slot="tabs"]',
  mount(root) {
    const cleanups = []
    const on = (el, type, fn) => {
      el.addEventListener(type, fn)
      cleanups.push(() => el.removeEventListener(type, fn))
    }
    const lists = [...root.querySelectorAll('[data-slot="tabs-list"]')].filter((el) => own(root, el))
    const triggers = () => [...root.querySelectorAll('[role="tab"]')].filter((el) => own(root, el))
    const contents = [...root.querySelectorAll('[role="tabpanel"]')].filter((el) => own(root, el))
    const valueOf = (el, kind) => {
      const id = el.id
      const i = id.indexOf(`-${kind}-`)
      return i < 0 ? id : id.slice(i + kind.length + 2)
    }
    const automatic = takeOption(root, "data-activation-mode") !== "manual"

    const selectedTrigger = triggers().find((t) => t.getAttribute("aria-selected") === "true")
    const activeContent = contents.find((c) => c.getAttribute("data-state") === "active")
    let value = selectedTrigger ? valueOf(selectedTrigger, "trigger") : activeContent ? valueOf(activeContent, "content") : ""

    // Panels: content templates and presence.
    const panels = contents.map((el) => {
      const tpl = el.querySelector(':scope > template[data-slot="tabs-content-template"]')
      const fragment = tpl?.content ?? null
      tpl?.remove()
      const selected = el.getAttribute("data-state") === "active"
      const panel = { el, fragment, preventMountAnimation: selected && el.style.animationDuration === "0s" }
      panel.presence = createPresence(el, {
        present: selected,
        onExitComplete() {
          el.hidden = true
          el.replaceChildren()
        },
      })
      return panel
    })

    const setValue = (next) => {
      if (next === value) return
      if (!emit(root, "shadcn:value-change", { value: next })) return
      value = next
      for (const t of triggers()) {
        const selected = valueOf(t, "trigger") === value
        t.setAttribute("aria-selected", String(selected))
        t.setAttribute("data-state", selected ? "active" : "inactive")
      }
      for (const p of panels) {
        const selected = valueOf(p.el, "content") === value
        if (p.preventMountAnimation) {
          // Radix drops animationDuration "0s" on the first re-render.
          p.preventMountAnimation = false
          p.el.style.animationDuration = ""
        }
        p.el.setAttribute("data-state", selected ? "active" : "inactive")
        const wasPresent = p.presence.isPresent
        p.presence.setPresent(selected)
        if (selected && !wasPresent) {
          p.el.hidden = false
          if (p.fragment) p.el.replaceChildren(p.fragment.cloneNode(true))
        }
      }
      for (const r of rovings) r.update()
    }

    for (const t of triggers()) {
      on(t, "mousedown", (event) => {
        if (!t.disabled && event.button === 0 && event.ctrlKey === false) setValue(valueOf(t, "trigger"))
        else event.preventDefault()
      })
      on(t, "keydown", (event) => {
        if ([" ", "Enter"].includes(event.key)) setValue(valueOf(t, "trigger"))
      })
      on(t, "focusin", () => {
        const isSelected = valueOf(t, "trigger") === value
        if (!isSelected && !t.disabled && automatic) setValue(valueOf(t, "trigger"))
      })
    }

    const rovings = lists.map((list) =>
      createRovingFocusGroup(list, {
        orientation: root.getAttribute("data-orientation") || "horizontal",
        loop: takeOption(list, "data-loop") !== "false",
        focusable: (t) => !t.disabled,
        active: (t) => t.getAttribute("aria-selected") === "true",
      }),
    )

    return () => {
      rovings.forEach((r) => r.destroy())
      panels.forEach((p) => p.presence.destroy())
      cleanups.forEach((fn) => fn())
    }
  },
})
