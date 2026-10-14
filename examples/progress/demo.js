// The progress-demo example: the value goes from 13 to 66 after 500ms,
// like the example's React state and effect.
import { define } from "./core.js"
import { setProgress } from "./progress.js"

define("example-progress-demo", {
  selector: '[data-slot="progress"][data-example="progress-demo"]',
  mount(root) {
    root.removeAttribute("data-example")
    const timer = setTimeout(() => setProgress(root, 66), 500)
    return () => clearTimeout(timer)
  },
})
