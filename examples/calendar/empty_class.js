// React renders className="" as class=""; libhtml drops empty class
// attributes, so the calendar-hijri wrapper (the Vazirmatn font's class
// name, empty in the reference build) gets it here.
import { define } from "./core.js"

define("example-empty-class", {
  selector: "[data-example-empty-class]",
  mount(el) {
    el.removeAttribute("data-example-empty-class")
    if (!el.hasAttribute("class")) el.setAttribute("class", "")
  },
})
