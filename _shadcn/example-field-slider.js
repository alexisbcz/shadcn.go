// The field-slider example's state: the description shows the slider's
// current range, like the example's React state.
import { define } from "./core.js"

define("example-field-slider", {
  selector: '[data-slot="field"]:has(> [data-slot="slider"][aria-label="Price Range"])',
  mount(field) {
    const slider = field.querySelector(':scope > [data-slot="slider"]')
    const spans = field.querySelectorAll(':scope > [data-slot="field-description"] > span.tabular-nums')
    const onChange = (e) => {
      const value = e.detail.value
      spans.forEach((span, i) => {
        if (span.textContent !== String(value[i])) span.textContent = String(value[i])
      })
    }
    slider.addEventListener("shadcn:value-change", onChange)
    return () => slider.removeEventListener("shadcn:value-change", onChange)
  },
})
