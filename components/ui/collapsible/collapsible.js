// collapsible.js: Radix Collapsible (@radix-ui/react-collapsible 1.1.12).
//
// - Clicking a trigger toggles the collapsible (a disabled trigger is a
//   disabled <button>: it gets no clicks).
// - The root, the triggers and the content carry data-state; the triggers
//   aria-expanded. The content is handled by collapsible-content.js
//   (presence, size variables, children mounted only while open).
// - "shadcn:open-change" (cancelable, detail: { open }) is dispatched on the
//   root before every change.
import { define, emit } from "./core.js"
import { createCollapsibleContent } from "./collapsible-content.js"

const ROOT = '[data-slot="collapsible"]'
const own = (root, el) => el.parentElement?.closest(ROOT) === root

define("collapsible", {
  selector: ROOT,
  mount(root) {
    const cleanups = []
    const triggers = () => [...root.querySelectorAll('[data-slot="collapsible-trigger"]')].filter((el) => own(root, el))
    const contents = [...root.querySelectorAll('[data-slot="collapsible-content"]')].filter((el) => own(root, el))
    let open = root.getAttribute("data-state") === "open"

    const panels = contents.map((el) => {
      const tpl = el.querySelector(':scope > template[data-slot="collapsible-content-template"]')
      const fragment = tpl?.content ?? null
      tpl?.remove()
      return createCollapsibleContent(el, { open, fragment })
    })

    const setOpen = (next) => {
      if (next === open) return
      if (!emit(root, "shadcn:open-change", { open: next })) return
      open = next
      root.setAttribute("data-state", open ? "open" : "closed")
      for (const t of triggers()) {
        t.setAttribute("aria-expanded", String(open))
        t.setAttribute("data-state", open ? "open" : "closed")
      }
      for (const p of panels) p.setOpen(open)
    }

    const onClick = (event) => {
      const t = event.target.closest?.('[data-slot="collapsible-trigger"]')
      if (!t || !own(root, t) || event.defaultPrevented) return
      setOpen(!open)
    }
    root.addEventListener("click", onClick)
    cleanups.push(() => root.removeEventListener("click", onClick))

    return () => {
      panels.forEach((p) => p.destroy())
      cleanups.forEach((fn) => fn())
    }
  },
})
