// The bubble-collapsible docs example's state: the bubble shows the full
// text and "Show less" while its collapsible is open, like the example's
// React state.
import { define } from "./core.js"

define("example-bubble-collapsible", {
  selector: "[data-example-bubble-collapsible]",
  mount(text) {
    const root = text.closest('[data-slot="collapsible"]')
    if (!root) return
    const full = text.dataset.full
    const preview = text.textContent
    const label = () => {
      const btn = root.querySelector('[data-slot="collapsible-trigger"]')
      return btn && Array.from(btn.childNodes).find((n) => n.nodeType === 3 && n.textContent.trim())
    }
    const onChange = (e) => {
      if (e.target !== root) return
      const open = e.detail.open
      text.textContent = open ? full : preview
      const l = label()
      if (l) l.textContent = open ? "Show less" : "Show more"
    }
    root.addEventListener("shadcn:open-change", onChange)
    return () => root.removeEventListener("shadcn:open-change", onChange)
  },
})
