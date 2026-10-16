// The button-group-select example's currency state: the select trigger
// renders the selected value as its only child.
import { define } from "./core.js"

define("example-button-group-currency", {
  selector: "[data-example-currency]",
  mount(trigger) {
    trigger.removeAttribute("data-example-currency")
    const onChange = (e) => {
      if (e.target.getAttribute?.("data-content-id") !== trigger.getAttribute("aria-controls")) return
      if (typeof e.detail?.value !== "string") return
      // Replace the text, keep the trigger's chevron icon.
      const text = [...trigger.childNodes].find((n) => n.nodeType === Node.TEXT_NODE)
      if (text) text.data = e.detail.value
      else trigger.prepend(e.detail.value)
    }
    document.addEventListener("shadcn:value-change", onChange)
    return () => document.removeEventListener("shadcn:value-change", onChange)
  },
})
