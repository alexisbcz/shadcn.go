// The shimmer-once docs preview: Replay remounts the paragraph (React's key
// change), which restarts its one-shot shimmer animation.
import { define } from "./core.js"

define("example-shimmer-replay", {
  selector: "[data-example-shimmer-replay]",
  mount(root) {
    const button = root.querySelector('[data-slot="button"]')
    const onClick = () => {
      const p = root.querySelector("p")
      p.replaceWith(p.cloneNode(true))
    }
    button.addEventListener("click", onClick)
    return () => button.removeEventListener("click", onClick)
  },
})
