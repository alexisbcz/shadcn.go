// toggle.js: Radix Toggle (@radix-ui/react-toggle 1.1.10) behavior.
//
// A click on an enabled toggle flips aria-pressed and data-state ("on" /
// "off"). The cancelable "shadcn:pressed-change" event (detail: { pressed })
// is dispatched first; cancelling it keeps the state (controlled usage).
import { define, emit } from "./core.js"

define("toggle", {
  selector: '[data-slot="toggle"]',
  mount(el) {
    const onClick = (event) => {
      // Radix composeEventHandlers: a consumer that prevented the click keeps the state.
      if (event.defaultPrevented || el.disabled) return
      const pressed = el.getAttribute("aria-pressed") !== "true"
      if (!emit(el, "shadcn:pressed-change", { pressed })) return
      el.setAttribute("aria-pressed", String(pressed))
      el.setAttribute("data-state", pressed ? "on" : "off")
    }
    el.addEventListener("click", onClick)
    return () => el.removeEventListener("click", onClick)
  },
})
