// input-group.js: InputGroupAddon's click handler. A click on an addon
// that does not land on a button focuses the first <input> of the group,
// like upstream's onClick.
import { define } from "./core.js"

define("input-group-addon", {
  selector: '[data-slot="input-group-addon"]',
  mount(el) {
    const onClick = (event) => {
      if (event.target instanceof Element && event.target.closest("button")) return
      el.parentElement?.querySelector("input")?.focus()
    }
    el.addEventListener("click", onClick)
    return () => el.removeEventListener("click", onClick)
  },
})
