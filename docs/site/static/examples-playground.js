// The playground example's React state (app/(app)/examples/playground):
// - preset-selector / model-selector: picking an item shows it on the
//   trigger, moves the check mark and closes the popover;
// - the temperature / maximum length / top p readouts follow their slider;
// - preset-actions: the menu items open the filter dialog and the delete
//   alert dialog, whose Delete button closes it with a toast.
import { define } from "./core.js"
import { getPopover } from "./popover.js"
import { getDialog } from "./dialog.js"
import { toast } from "./sonner.js"

define("examples-playground", {
  selector: "[data-playground-picker]",
  mount(trigger) {
    const contentId = trigger.getAttribute("aria-controls")
    const icon = trigger.querySelector("svg")
    const mark = (root, name) => {
      for (const item of root.querySelectorAll("[data-playground-pick]")) {
        const check = item.querySelector("svg")
        const on = item.getAttribute("data-playground-pick") === name
        check?.classList.remove(on ? "opacity-0" : "opacity-100")
        check?.classList.add(on ? "opacity-100" : "opacity-0")
      }
    }
    const onSelect = (e) => {
      const content = e.target.closest?.('[data-slot="popover-content"]')
      if (!content || content.id !== contentId) return
      const item = e.target.closest("[data-playground-pick]")
      if (!item) return
      const name = item.getAttribute("data-playground-pick")
      trigger.replaceChildren(name, icon)
      mark(content, name)
      const tpl = document.querySelector(`template[data-content-id="${CSS.escape(contentId)}"]`)
      if (tpl) mark(tpl.content, name)
      getPopover(contentId)?.setOpen(false)
    }
    document.addEventListener("shadcn:select", onSelect)
    return () => document.removeEventListener("shadcn:select", onSelect)
  },
})

define("examples-playground-slider", {
  selector: "[data-playground-slider]",
  mount(root) {
    const id = root.getAttribute("data-playground-slider")
    const out = document.querySelector(`[data-playground-slider-value="${CSS.escape(id)}"]`)
    const onChange = (e) => {
      if (out) out.textContent = e.detail.value.join("")
    }
    root.addEventListener("shadcn:value-change", onChange)
    return () => root.removeEventListener("shadcn:value-change", onChange)
  },
})

const dialogOf = (name) =>
  getDialog(document.querySelector(`[data-playground-dialog="${name}"] template[data-radix-dialog]`))

define("examples-playground-actions", {
  selector: "[data-playground-open]",
  mount(item) {
    const onSelect = () => {
      const name = item.getAttribute("data-playground-open")
      // The menu restores focus to its trigger as it closes; open after.
      setTimeout(() => dialogOf(name)?.setOpen(true))
    }
    item.addEventListener("shadcn:select", onSelect)
    return () => item.removeEventListener("shadcn:select", onSelect)
  },
})

document.addEventListener("click", (e) => {
  const b = e.target instanceof Element ? e.target.closest("[data-playground-delete]") : null
  if (!b) return
  dialogOf("delete")?.setOpen(false)
  toast.success("This preset has been deleted.")
})
