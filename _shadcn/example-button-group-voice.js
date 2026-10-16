// The button-group-input-group example's state: the voice button toggles
// voice mode, which disables the input and changes its placeholder.
import { define } from "./core.js"

define("example-button-group-voice", {
  selector: '[data-slot="input-group"] button[data-active][aria-pressed]',
  mount(button) {
    const input = button.closest('[data-slot="input-group"]')?.querySelector('[data-slot="input-group-control"]')
    const onClick = () => {
      const on = button.getAttribute("aria-pressed") !== "true"
      button.setAttribute("data-active", String(on))
      button.setAttribute("aria-pressed", String(on))
      if (!input) return
      input.placeholder = on ? "Record and send audio..." : "Send a message..."
      input.disabled = on
    }
    button.addEventListener("click", onClick)
    return () => button.removeEventListener("click", onClick)
  },
})
