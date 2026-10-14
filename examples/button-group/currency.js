// The button-group-select example's currency state: the select trigger
// renders the selected value as its only child.
import { define } from "./core.js"

define("example-button-group-currency", {
  selector: "[data-example-currency]",
  mount(trigger) {
    trigger.removeAttribute("data-example-currency")
    const onChange = (e) => {
      if (e.target.getAttribute?.("data-content-id") !== trigger.getAttribute("aria-controls")) return
      if (typeof e.detail?.value === "string") trigger.textContent = e.detail.value
    }
    document.addEventListener("shadcn:value-change", onChange)
    return () => document.removeEventListener("shadcn:value-change", onChange)
  },
})
