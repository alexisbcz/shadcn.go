// label.js: Radix Label behavior (@radix-ui/react-label 2.1.7). A
// mousedown on the label that is a double (or further) click is
// prevented, so the label's text is not selected, unless it lands on a
// control nested in the label.
import { define } from "./core.js"

define("label", {
  selector: '[data-slot="label"]',
  mount(el) {
    const onMouseDown = (event) => {
      const target = event.target
      if (target instanceof Element && target.closest("button, input, select, textarea")) return
      if (!event.defaultPrevented && event.detail > 1) event.preventDefault()
    }
    el.addEventListener("mousedown", onMouseDown)
    return () => el.removeEventListener("mousedown", onMouseDown)
  },
})
