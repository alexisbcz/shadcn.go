// The card-spacing docs preview: a single toggle group (controlled, never
// empty) sets the card's --card-spacing class.
import { define } from "./core.js"

define("example-card-spacing", {
  selector: "[data-example-card-spacing]",
  mount(root) {
    const card = root.querySelector('[data-slot="card"]')
    const onChange = (e) => {
      const item = e.detail?.item
      if (!item) return
      if (!e.detail.pressed) {
        e.preventDefault()
        return
      }
      const v = item.getAttribute("data-spacing-class") ?? ""
      for (const it of root.querySelectorAll("[data-spacing-class]")) card.classList.remove(it.getAttribute("data-spacing-class"))
      card.classList.add(v)
    }
    root.addEventListener("shadcn:value-change", onChange)
    return () => root.removeEventListener("shadcn:value-change", onChange)
  },
})
